package provider

import (
    "context"
    "fmt"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/path"

)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &projectResource{}
    _ resource.ResourceWithConfigure = &projectResource{}
)

// NewProjectResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
    return &projectResource{}
}

// projectResource is the resource implementation.
type projectResource struct {
    client *client.Client
}

// projectResourceModel maps the resource schema data.
type projectResourceModel struct {
    ID                types.String `tfsdk:"id"`
    OrganizationID    types.String `tfsdk:"organization_id"`
    Name              types.String `tfsdk:"name"`
    Title             types.String `tfsdk:"title"`
    OwnerID           types.String `tfsdk:"owner_id"`
    TeamID            types.String `tfsdk:"team_id"`
    IsConnectProject  types.Bool   `tfsdk:"is_connect_project"`
    IsHidden          types.Bool   `tfsdk:"is_hidden"`
}

// Configure adds the provider configured client to the resource.
func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a project.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identifier of the project.",
                Computed:    true,
            },
            "organization_id": schema.StringAttribute{
                Description: "Identifier of the organization.",
                Required:    true,
            },
            "name": schema.StringAttribute{
                Description: "Name of the project.",
                Required:    true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
            },
            "title": schema.StringAttribute{
                Description: "Title of the project.",
                Computed:    true,
            },
            "owner_id": schema.StringAttribute{
                Description: "Identifier of the project owner.",
                Computed:    true,
            },
            "team_id": schema.StringAttribute{
                Description: "Identifier of the team associated with the project.",
                Computed:    true,
            },
            "is_connect_project": schema.BoolAttribute{
                Description: "Indicates if the project is a Connect project.",
                Computed:    true,
            },
            "is_hidden": schema.BoolAttribute{
                Description: "Indicates if the project is hidden.",
                Computed:    true,
            },
        },
    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan projectResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

     // Get the list of teams
    teams, err := r.client.GetTeams(ctx)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error retrieving teams",
            "Could not retrieve teams, unexpected error: "+err.Error(),
        )
        return
    }

    // Check if a team with the same name already exists
    for _, team := range teams {
        if team.Name == plan.Name.ValueString() {
            resp.Diagnostics.AddError(
                "Project already exists",
                fmt.Sprintf("A project with the name '%s' already exists", plan.Name.ValueString()),
            )
            return
        }
    }

    // Create new project
    project, err := r.client.CreateProject(ctx, plan.OrganizationID.ValueString(), plan.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating project",
            "Could not create project, unexpected error: "+err.Error(),
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
    plan.ID = types.StringValue(project.ID)
    plan.Title = types.StringValue(project.Title)
    plan.OwnerID = types.StringValue(project.OwnerID)
    plan.TeamID = types.StringValue(project.TeamID)
    plan.IsConnectProject = types.BoolValue(project.IsConnectProject)
    plan.IsHidden = types.BoolValue(project.IsHidden)

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read refreshes the Terraform state with the latest data.
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state projectResourceModel

    // Get the project ID and team ID from the state
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ID.ValueString()
    teamID := state.TeamID.ValueString()

    // Retrieve the project using the GetProjectByID function
    project, err := r.client.GetProjectByID(ctx, projectID, teamID)
    if err != nil {
        if strings.Contains(err.Error(), "status code: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error reading project",
            "Could not read project, unexpected error: "+err.Error(),
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
    state.ID = types.StringValue(project.ID)
    state.Title = types.StringValue(project.Title)
    state.OwnerID = types.StringValue(project.OwnerID)
    state.TeamID = types.StringValue(project.TeamID)
    state.IsConnectProject = types.BoolValue(project.IsConnectProject)
    state.IsHidden = types.BoolValue(project.IsHidden)

    // Set the refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan projectResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state projectResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Check if the organization_id has changed
    if !plan.OrganizationID.Equal(state.OrganizationID) {
        resp.Diagnostics.AddAttributeError(
            path.Root("organization_id"),
            "Immutable Attribute Change",
            "The 'organization_id' attribute cannot be changed once the project is created. "+
                "To change the organization, delete the project and create a new one with the desired organization.",
        )
        return
    }

    projectID := state.ID.ValueString()
    teamID := state.TeamID.ValueString()

    // Check if the name has changed
    if !plan.Name.Equal(state.Name) {
        // Update the project title using the UpdateProjectTitle function
        updatedProject, err := r.client.UpdateProjectTitle(ctx, projectID, teamID, plan.Name.ValueString())
        if err != nil {
            if strings.Contains(err.Error(), "status code: 404") {
                resp.Diagnostics.AddError(
                    "Project not found during update",
                    "The project was not found while attempting to update it. This is an unexpected error.",
                )
                return
            }
            resp.Diagnostics.AddError(
                "Error updating project title",
                "Could not update project title, unexpected error: "+err.Error(),
            )
            return
        }

        // Update the state with the updated project attributes
        plan.ID = types.StringValue(updatedProject.ID)
        plan.Title = types.StringValue(updatedProject.Title)
        plan.OwnerID = types.StringValue(updatedProject.OwnerID)
        plan.TeamID = types.StringValue(updatedProject.TeamID)
        plan.IsConnectProject = types.BoolValue(updatedProject.IsConnectProject)
        plan.IsHidden = types.BoolValue(updatedProject.IsHidden)
    }

    // Set the updated state
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}



// Delete deletes the resource and removes the Terraform state on success.
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state projectResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ID.ValueString()
    teamID := state.TeamID.ValueString()

    // Delete the project using the DeleteProject function
    err := r.client.DeleteProject(ctx, projectID, teamID)
    if err != nil {
        if !strings.Contains(err.Error(), "status code: 404") {
            resp.Diagnostics.AddError(
                "Error deleting project",
                "Could not delete project, unexpected error: "+err.Error(),
            )
            return
        }
    }
}