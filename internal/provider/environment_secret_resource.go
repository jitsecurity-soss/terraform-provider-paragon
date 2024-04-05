// environment_secret_resource.go
package provider

import (
    "context"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &environmentSecretResource{}
    _ resource.ResourceWithConfigure = &environmentSecretResource{}
)

// NewEnvironmentSecretResource is a helper function to simplify the provider implementation.
func NewEnvironmentSecretResource() resource.Resource {
    return &environmentSecretResource{}
}

// environmentSecretResource is the resource implementation.
type environmentSecretResource struct {
    client *client.Client
}

// environmentSecretResourceModel maps the resource schema data.
type environmentSecretResourceModel struct {
    ID        types.String `tfsdk:"id"`
    ProjectID types.String `tfsdk:"project_id"`
    Key       types.String `tfsdk:"key"`
    Value     types.String `tfsdk:"value"`
    Hash      types.String `tfsdk:"hash"`
}

// Configure adds the provider configured client to the resource.
func (r *environmentSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *environmentSecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_environment_secret"
}

// Schema defines the schema for the resource.
func (r *environmentSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages an environment secret.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identifier of the environment secret.",
                Computed:    true,
            },
            "project_id": schema.StringAttribute{
                Description: "Identifier of the project.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "key": schema.StringAttribute{
                Description: "Key of the environment secret.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
            },
            "value": schema.StringAttribute{
                Description: "Value of the environment secret.",
                Required:    true,
                Sensitive:   true,
                Validators: []validator.String{
                    stringvalidator.LengthAtLeast(1),
                },
            },
            "hash": schema.StringAttribute{
                Description: "Hash of the environment secret.",
                Computed:    true,
            },
        },
    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *environmentSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan environmentSecretResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := plan.ProjectID.ValueString()
    key := plan.Key.ValueString()
    value := plan.Value.ValueString()

    // Create new environment secret
    secret, err := r.client.CreateEnvironmentSecret(ctx, projectID, key, value)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating environment secret",
            "Could not create environment secret, unexpected error: "+err.Error(),
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
    plan.ID = types.StringValue(secret.ID)
    plan.Hash = types.StringValue(secret.Hash)

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read refreshes the Terraform state with the latest data.
func (r *environmentSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state environmentSecretResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ProjectID.ValueString()

    // Retrieve the environment secrets using the GetEnvironmentSecrets function
    secrets, err := r.client.GetEnvironmentSecrets(ctx, projectID)
    if err != nil {

        if strings.Contains(err.Error(), "status code: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error reading environment secrets",
            "Could not read environment secrets, unexpected error: "+err.Error(),
        )
        return
    }

    // Find the environment secret with the matching ID
    var secret *client.EnvironmentSecret
    for _, s := range secrets {
        if s.ID == state.ID.ValueString() {
            secret = &s
            break
        }
    }

    if secret == nil {
        // If the environment secret is not found, remove the resource from the state
        resp.State.RemoveResource(ctx)
        return
    }

    // Update the state with the latest data
    state.Key = types.StringValue(secret.Key)
    state.Hash = types.StringValue(secret.Hash)

    // Set the refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *environmentSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan environmentSecretResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var state environmentSecretResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ProjectID.ValueString()
    secretID := state.ID.ValueString()
    key := plan.Key.ValueString()
    value := plan.Value.ValueString()

    // Update the environment secret using the UpdateEnvironmentSecret function
    updatedSecret, err := r.client.UpdateEnvironmentSecret(ctx, projectID, secretID, key, value)
    if err != nil {
        if strings.Contains(err.Error(), "status code: 404") {
            resp.Diagnostics.AddError(
                "Environment secret not found during update",
                "The environment secret was not found while attempting to update it. This is an unexpected error.",
            )
            return
        }
        resp.Diagnostics.AddError(
            "Error updating environment secret",
            "Could not update environment secret, unexpected error: "+err.Error(),
        )
        return
    }

    // Update the state with the updated data
    plan.Hash = types.StringValue(updatedSecret.Hash)
    plan.ID = types.StringValue(updatedSecret.ID)

    // Set the updated state
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *environmentSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state environmentSecretResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ProjectID.ValueString()
    secretID := state.ID.ValueString()

    // Delete the environment secret using the DeleteEnvironmentSecret function
    err := r.client.DeleteEnvironmentSecret(ctx, projectID, secretID)
    if err != nil {
        if !strings.Contains(err.Error(), "status code: 404") {
            resp.Diagnostics.AddError(
                "Error deleting environment secret",
                "Could not delete environment secret, unexpected error: "+err.Error(),
            )
            return
        }
    }
}