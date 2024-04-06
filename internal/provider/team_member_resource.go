package provider

import (
    "context"
    "fmt"
    "regexp"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/hashicorp/terraform-plugin-framework/path"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &teamMemberResource{}
    _ resource.ResourceWithConfigure = &teamMemberResource{}
)

// NewTeamMemberResource is a helper function to simplify the provider implementation.
func NewTeamMemberResource() resource.Resource {
    return &teamMemberResource{}
}

// teamMemberResource is the resource implementation.
type teamMemberResource struct {
    client *client.Client
}

// teamMemberResourceModel maps the resource schema data.
type teamMemberResourceModel struct {
    ID             types.String `tfsdk:"id"`
    TeamID         types.String `tfsdk:"team_id"`
    Email          types.String `tfsdk:"email"`
    Role           types.String `tfsdk:"role"`
}

// Configure adds the provider configured client to the resource.
func (r *teamMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *teamMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_team_member"
}

// Schema defines the schema for the resource.
func (r *teamMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a team member.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identifier of the team member.",
                Computed:    true,
            },
            "team_id": schema.StringAttribute{
                Description: "Identifier of the team.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "email": schema.StringAttribute{
                Description: "Email address of the team member.",
                Required:    true,
				Validators: []validator.String{
					// Validate string value satisfies the regular expression for alphanumeric characters
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
						"Must be a valid email address",
					),
				},
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "role": schema.StringAttribute{
                Description: "Role of the team member (ADMIN, MEMBER, SUPPORT).",
                Required:    true,
				Validators: []validator.String{
					// Validate string value satisfies the regular expression for alphanumeric characters
					stringvalidator.OneOf("ADMIN", "MEMBER", "SUPPORT"),
				},
            },
        },
    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *teamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan teamMemberResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    teamID := plan.TeamID.ValueString()
    email := plan.Email.ValueString()
    role := plan.Role.ValueString()

    // Check if the team member already exists
    members, err := r.client.GetTeamMembers(ctx, teamID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading team members",
            "Could not read team members, unexpected error: "+err.Error(),
        )
        return
    }
    for _, member := range members {
        if member.Email == email {
            resp.Diagnostics.AddError(
                "Team member already exists",
                fmt.Sprintf("A team member with email '%s' already exists", email),
            )
            return
        }
    }

    // Check if a team invite already exists for the email
    invites, err := r.client.GetTeamInvites(ctx, teamID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading team invites",
            "Could not read team invites, unexpected error: "+err.Error(),
        )
        return
    }
    for _, invite := range invites {
        if invite.Email == email {
            resp.Diagnostics.AddError(
                "Team invite already exists",
                fmt.Sprintf("A team invite for email '%s' already exists", email),
            )
            return
        }
    }

    // Invite the team member
    invites, err = r.client.InviteTeamMember(ctx, teamID, role, email)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error inviting team member",
            "Could not invite team member, unexpected error: "+err.Error(),
        )
        return
    }

    if len(invites) == 0 {
        resp.Diagnostics.AddError(
            "Error inviting team member",
            "No team invite was created",
        )
        return
    }

    invite := invites[0]

    // Map response body to schema and populate Computed attribute values
    plan.ID = types.StringValue(invite.ID)

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}


// Read refreshes the Terraform state with the latest data.
func (r *teamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state teamMemberResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    teamID := state.TeamID.ValueString()
    email := state.Email.ValueString()

    // Check if the team member exists in the members list
    tflog.Debug(ctx, "Getting team members...")
    members, err := r.client.GetTeamMembers(ctx, teamID)
    if err != nil {
        if strings.Contains(err.Error(), "status code: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error reading team members",
            "Could not read team members, unexpected error: "+err.Error(),
        )
        return
    }

    for _, member := range members {
        if member.Email == email {
            tflog.Debug(ctx, "Found member email in regular list! Updating ID.")

            // Update the state with the new ID
            state.ID = types.StringValue(member.ID)
            state.Role = types.StringValue(member.Role)

            // Set the refreshed state
            diags := resp.State.Set(ctx, &state)
            resp.Diagnostics.Append(diags...)
            if resp.Diagnostics.HasError() {
                return
            }
            return
        }
    }

    // Check if the team member exists in the invites
    tflog.Debug(ctx, "Searching invites...")
    invites, err := r.client.GetTeamInvites(ctx, teamID)
    if err != nil {
        if strings.Contains(err.Error(), "status code: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error reading team invites",
            "Could not read team invites, unexpected error: "+err.Error(),
        )
        return
    }

    for _, invite := range invites {
        if invite.Email == email {
            tflog.Debug(ctx, "Found member email in invites! All good.")

            // Map the invite data to the state
            state.Role = types.StringValue(invite.Role)

            // Set the refreshed state
            diags = resp.State.Set(ctx, &state)
            resp.Diagnostics.Append(diags...)
            if resp.Diagnostics.HasError() {
                return
            }
            return
        }
    }

    // If the team member is not found in either members or invites, remove the resource from the state
    resp.State.RemoveResource(ctx)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan teamMemberResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state teamMemberResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    teamID := state.TeamID.ValueString()
    email := state.Email.ValueString()
    role := plan.Role.ValueString()

    // Check if the team member exists in the invites
    tflog.Debug(ctx, "Searching invites...")
    invites, err := r.client.GetTeamInvites(ctx, teamID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading team invites",
            "Could not read team invites, unexpected error: "+err.Error(),
        )
        return
    }

    for _, invite := range invites {
        if invite.Email == email {
            // If the user exists in the invites and the role is planned to change, trigger a recreation
            if invite.Role != role {
                resp.Diagnostics.AddAttributeError(
                    path.Root("role"),
                    "Role update not allowed",
                    "Updating the role of a team member in the invite state is not permitted. Please delete it first.",
                )
                return
            }
        }
    }

    // Check if the team member exists in the members list
    tflog.Debug(ctx, "Getting team members...")
    members, err := r.client.GetTeamMembers(ctx, teamID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading team members",
            "Could not read team members, unexpected error: "+err.Error(),
        )
        return
    }

    for _, member := range members {
        if member.Email == email {
            // If the user exists in the members and the role is planned to change, update the role
            if member.Role != role {
                updatedMember, err := r.client.UpdateTeamMemberRole(ctx, teamID, member.ID, role)
                if err != nil {
                    resp.Diagnostics.AddError(
                        "Error updating team member role",
                        "Could not update team member role, unexpected error: "+err.Error(),
                    )
                    return
                }

                // Update the state with the updated role
                state.Role = types.StringValue(updatedMember.Role)

                // Set the updated state
                diags := resp.State.Set(ctx, &state)
                resp.Diagnostics.Append(diags...)
                if resp.Diagnostics.HasError() {
                    return
                }
            }
            return
        }
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state teamMemberResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    teamID := state.TeamID.ValueString()
    memberID := state.ID.ValueString()

    // Delete the team member from the members list
    err := r.client.DeleteTeamMember(ctx, teamID, memberID)
    if err != nil {
        if err.Error() != "status code: 404" {
            resp.Diagnostics.AddError(
               "Error deleting team member",
                fmt.Sprintf("Could not delete team member from members list, unexpected error: %s\nTeam ID: %s\nMember ID: %s", err.Error(), teamID, memberID),
)
            return
        }

        // If the member is not found in the members list, try deleting from the invites
        err = r.client.DeleteTeamInvite(ctx, teamID, memberID)
        if err != nil {
            if err.Error() != "status code: 404"  {
                resp.Diagnostics.AddError(
                    "Error deleting team invite",
                    "Could not delete team invite, unexpected error: "+err.Error(),
                )
                return
            }
        }
    }

    // Remove the resource from the state
    resp.State.RemoveResource(ctx)
}