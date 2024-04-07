package provider

import (
    "context"
    "fmt"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &integrationStatusResource{}
    _ resource.ResourceWithConfigure = &integrationStatusResource{}
)

// NewIntegrationStatusResource is a helper function to simplify the provider implementation.
func NewIntegrationStatusResource() resource.Resource {
    return &integrationStatusResource{}
}

// integrationStatusResource is the resource implementation.
type integrationStatusResource struct {
    client *client.Client
}

// integrationStatusResourceModel maps the resource schema data.
type integrationStatusResourceModel struct {
    ID            types.String `tfsdk:"id"`
    ProjectID     types.String `tfsdk:"project_id"`
    IntegrationID types.String `tfsdk:"integration_id"`
    Active        types.Bool   `tfsdk:"active"`
}

// Configure adds the provider configured client to the resource.
func (r *integrationStatusResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *integrationStatusResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_integration_status"
}

// Schema defines the schema for the resource.
func (r *integrationStatusResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages the status of an integration.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identifier of the integration status.",
                Computed:    true,
            },
            "project_id": schema.StringAttribute{
                Description: "Identifier of the project.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "integration_id": schema.StringAttribute{
                Description: "Identifier of the integration.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "active": schema.BoolAttribute{
                Description: "Indicates whether the integration is active or not.",
                Required:    true,
            },
        },
    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *integrationStatusResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan integrationStatusResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := plan.ProjectID.ValueString()
    integrationID := plan.IntegrationID.ValueString()
    active := plan.Active.ValueBool()

    // Update the integration status
    integration, err := r.client.UpdateIntegrationStatus(ctx, projectID, integrationID, active)
    if err != nil {
        if strings.Contains(err.Error(), "status code: 404") {
            resp.Diagnostics.AddError(
                "Integration not found",
                fmt.Sprintf("Integration with ID '%s' not found in the project", integrationID),
            )
        } else {
            resp.Diagnostics.AddError(
                "Error updating integration status",
                "Could not update integration status, unexpected error: "+err.Error(),
            )
        }
        return
    }

    // Set the ID and active status in the state
    plan.ID = types.StringValue(integration.ID)
    plan.Active = types.BoolValue(integration.IsActive)

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read refreshes the Terraform state with the latest data.
func (r *integrationStatusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state integrationStatusResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ProjectID.ValueString()
    integrationID := state.IntegrationID.ValueString()

    // Retrieve the integration
    integration, err := r.client.GetIntegration(ctx, projectID, integrationID)
    if err != nil {
        if strings.Contains(err.Error(), "status code: 404") {
            resp.State.RemoveResource(ctx)
        } else {
            resp.Diagnostics.AddError(
                "Error retrieving integration",
                "Could not retrieve integration, unexpected error: "+err.Error(),
            )
        }
        return
    }

    // Update the state with the retrieved data
    state.Active = types.BoolValue(integration.IsActive)
    state.ID = types.StringValue(integration.ID)

    // Set the refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *integrationStatusResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan integrationStatusResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := plan.ProjectID.ValueString()
    integrationID := plan.IntegrationID.ValueString()
    active := plan.Active.ValueBool()

    // Update the integration status
    integration, err := r.client.UpdateIntegrationStatus(ctx, projectID, integrationID, active)
    if err != nil {
        if strings.Contains(err.Error(), "status code: 404") {
            resp.Diagnostics.AddError(
                "Integration not found during update",
                fmt.Sprintf("Integration with ID '%s' not found in the project", integrationID),
            )
        } else {
            resp.Diagnostics.AddError(
                "Error updating integration status",
                "Could not update integration status, unexpected error: "+err.Error(),
            )
        }
        return
    }

    // Update the state with the retrieved data
    plan.Active = types.BoolValue(integration.IsActive)
    plan.ID = types.StringValue(integration.ID)

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *integrationStatusResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state integrationStatusResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ProjectID.ValueString()
    integrationID := state.IntegrationID.ValueString()

    // Update the integration status to inactive (false)
    _, err := r.client.UpdateIntegrationStatus(ctx, projectID, integrationID, false)
    if err != nil {
        if !strings.Contains(err.Error(), "status code: 404") {
            resp.Diagnostics.AddError(
                "Error updating integration status",
                "Could not update integration status, unexpected error: "+err.Error(),
            )
            return
        }
    }
}