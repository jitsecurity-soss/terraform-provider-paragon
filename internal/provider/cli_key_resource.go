package provider

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &cliKeyResource{}
    _ resource.ResourceWithConfigure = &cliKeyResource{}
)

// NewCLIKeyResource is a helper function to simplify the provider implementation.
func NewCLIKeyResource() resource.Resource {
    return &cliKeyResource{}
}

// cliKeyResource is the resource implementation.
type cliKeyResource struct {
    client *client.Client
}

// cliKeyResourceModel maps the resource schema data.
type cliKeyResourceModel struct {
    ID             types.String `tfsdk:"id"`
    OrganizationID types.String `tfsdk:"organization_id"`
    KeyName        types.String `tfsdk:"name"`
    Key            types.String `tfsdk:"key"`
}

// Configure adds the provider configured client to the resource.
func (r *cliKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *cliKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cli_key"
}

// Schema defines the schema for the resource.
func (r *cliKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a CLI key.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identifier of the CLI key.",
                Computed:    true,
            },
            "organization_id": schema.StringAttribute{
                Description: "Identifier of the organization.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "name": schema.StringAttribute{
                Description: "Name of the CLI key.",
                Required:    true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
            },
            "key": schema.StringAttribute{
                Description: "The CLI key.",
                Computed:    true,
                Sensitive:   true,
            },
        },
    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *cliKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan cliKeyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    keyName := plan.KeyName.ValueString()
    organizationID := plan.OrganizationID.ValueString()


    // Extract the user ID from the access token
    tflog.Debug(ctx, "Extracting user ID from the JWT token...")
    userID, err := r.client.GetUserIDFromToken()
    if err != nil {
        resp.Diagnostics.AddError(
            "Error extracting user ID from access token",
            "Could not extract user ID from access token, unexpected error: "+err.Error(),
        )
        return
    }

    // Check if a CLI key with the same user ID and name already exists
    tflog.Debug(ctx, fmt.Sprintf("Get CLI keys for org id: %s", organizationID))
    cliKeys, err := r.client.GetCLIKeys(ctx, organizationID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading CLI keys",
            "Could not read CLI keys, unexpected error: "+err.Error(),
        )
        return
    }

    tflog.Debug(ctx, fmt.Sprintf("searching if we find key with userid: %s, keyname: %s", userID, keyName))
    for _, cliKey := range cliKeys {
        tflog.Debug(ctx, fmt.Sprintf("checking userid: %s, keyname: %s", cliKey.UserID, cliKey.Name))
        if cliKey.UserID == userID && cliKey.Name == keyName {
            resp.Diagnostics.AddError(
                "CLI key already exists",
                fmt.Sprintf("A CLI key with user ID '%s' and name '%s' already exists", userID, keyName),
            )
            return
        }
    }

    // Create the CLI key
    tflog.Debug(ctx, fmt.Sprintf("Create CLI key %s", keyName))
    cliKeyResp, err := r.client.CreateCLIKey(ctx, keyName)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating CLI key",
            "Could not create CLI key, unexpected error: "+err.Error(),
        )
        return
    }

    // After successful creation, get the keys again to get the ID
    tflog.Debug(ctx, fmt.Sprintf("Getting keys again to find new ID for org:  %s", organizationID))
    cliKeys, err = r.client.GetCLIKeys(ctx, organizationID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading CLI keys",
            "Could not read CLI keys, unexpected error: "+err.Error(),
        )
        return
    }

    // Retrieve the ID of the created CLI key
    var createdCLIKey *client.CLIKey
    for _, cliKey := range cliKeys {
        if cliKey.UserID == userID && cliKey.Name == keyName {
            createdCLIKey = &cliKey
            break
        }
    }

    if createdCLIKey == nil {
        resp.Diagnostics.AddError(
            "Error retrieving created CLI key",
            "Could not retrieve the created CLI key",
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
    plan.ID = types.StringValue(createdCLIKey.ID)
    plan.Key = types.StringValue(cliKeyResp.Key)

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read refreshes the Terraform state with the latest data.
func (r *cliKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state cliKeyResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    organizationID := state.OrganizationID.ValueString()
    keyID := state.ID.ValueString()

    // Retrieve the list of CLI keys for the organization
    cliKeys, err := r.client.GetCLIKeys(ctx, organizationID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading CLI keys",
            "Could not read CLI keys, unexpected error: "+err.Error(),
        )
        return
    }

    // Find the CLI key with the matching ID
    var foundCLIKey *client.CLIKey
    for _, cliKey := range cliKeys {
        if cliKey.ID == keyID {
            foundCLIKey = &cliKey
            break
        }
    }

    if foundCLIKey == nil {
        resp.State.RemoveResource(ctx)
        return
    }

    // Update the state with the retrieved data
    state.KeyName = types.StringValue(foundCLIKey.Name)

    // Set the refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *cliKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan cliKeyResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state cliKeyResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    organizationID := state.OrganizationID.ValueString()
    keyID := state.ID.ValueString()
    newName := plan.KeyName.ValueString()

    // Update the CLI key with the new name
    updatedCLIKey, err := r.client.UpdateCLIKey(ctx, organizationID, keyID, newName)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating CLI key",
            "Could not update CLI key, unexpected error: "+err.Error(),
        )
        return
    }

    // Update the state with the updated data
    plan.ID = state.ID
    plan.KeyName = types.StringValue(updatedCLIKey.Name)
    plan.Key = state.Key

    // Set the updated state
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *cliKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state cliKeyResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    organizationID := state.OrganizationID.ValueString()
    keyID := state.ID.ValueString()

    // Delete the CLI key
    err := r.client.DeleteCLIKey(ctx, organizationID, keyID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error deleting CLI key",
            "Could not delete CLI key, unexpected error: "+err.Error(),
        )
        return
    }
}