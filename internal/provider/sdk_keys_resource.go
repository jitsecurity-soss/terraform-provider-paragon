// sdk_keys_resource.go
package provider

import (
    "context"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &sdkKeysResource{}
    _ resource.ResourceWithConfigure = &sdkKeysResource{}
)

// NewSDKKeysResource is a helper function to simplify the provider implementation.
func NewSDKKeysResource() resource.Resource {
    return &sdkKeysResource{}
}

// sdkKeysResource is the resource implementation.
type sdkKeysResource struct {
    client *client.Client
}

// sdkKeysResourceModel maps the resource schema data.
type sdkKeysResourceModel struct {
    ID            types.String `tfsdk:"id"`
    ProjectID     types.String `tfsdk:"project_id"`
    AuthType      types.String `tfsdk:"auth_type"`
    Revoked       types.Bool   `tfsdk:"revoked"`
    GeneratedDate types.String `tfsdk:"generated_date"`
    PrivateKey    types.String `tfsdk:"private_key"`
    Version       types.String `tfsdk:"version"`
}

// Configure adds the provider configured client to the resource.
func (r *sdkKeysResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *sdkKeysResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_sdk_keys"
}

// Schema defines the schema for the resource.
func (r *sdkKeysResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages SDK keys for a project.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Identifier of the SDK key.",
                Computed:    true,
            },
            "project_id": schema.StringAttribute{
                Description: "Identifier of the project.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "auth_type": schema.StringAttribute{
                Description: "Authentication type of the SDK key.",
                Computed:    true,
            },
            "revoked": schema.BoolAttribute{
                Description: "Indicates if the SDK key is revoked.",
                Computed:    true,
            },
            "generated_date": schema.StringAttribute{
                Description: "Date when the SDK key was generated.",
                Computed:    true,
            },
            "private_key": schema.StringAttribute{
                Description: "Private key of the SDK key.",
                Computed:    true,
                Sensitive:   true,
            },
            "version": schema.StringAttribute{
                Description: "Version of the SDK key.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
        },
    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *sdkKeysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan sdkKeysResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := plan.ProjectID.ValueString()
    version := plan.Version.ValueString()

    // Create new SDK key
    sdkKey, err := r.client.CreateSDKKey(ctx, projectID)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating SDK key",
            "Could not create SDK key, unexpected error: "+err.Error(),
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
    plan.ID = types.StringValue(sdkKey.ID)
    plan.AuthType = types.StringValue(sdkKey.AuthType)
    plan.Revoked = types.BoolValue(sdkKey.Revoked)
    plan.GeneratedDate = types.StringValue(sdkKey.AuthConfig.Paragon.GeneratedDate)
    plan.PrivateKey = types.StringValue(sdkKey.PrivateKey)
    plan.Version = types.StringValue(version)

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (r *sdkKeysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state sdkKeysResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ProjectID.ValueString()

    // Retrieve the SDK keys using the GetSDKKeys function
    sdkKeys, err := r.client.GetSDKKeys(ctx, projectID)
    if err != nil {
        // Check if the error indicates a 404 status code
        if strings.Contains(err.Error(), "status code: 404") {
            // If the SDK key is not found, remove the resource to trigger recreation
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error reading SDK keys",
            "Could not read SDK keys, unexpected error: "+err.Error(),
        )
        return
    }

    // Find the SDK key with the matching ID
    var sdkKey *client.SDKKey
    for _, key := range sdkKeys {
        if key.ID == state.ID.ValueString() {
            sdkKey = &key
            break
        }
    }

    if sdkKey == nil {
        // If the SDK key is not found, remove the resource to trigger recreation
        resp.State.RemoveResource(ctx)
        return
    }

    // Map the SDK key data to the state
    state.AuthType = types.StringValue(sdkKey.AuthType)
    state.Revoked = types.BoolValue(sdkKey.Revoked)
    state.GeneratedDate = types.StringValue(sdkKey.AuthConfig.Paragon.GeneratedDate)

    // Set the refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update is not supported for this resource.
func (r *sdkKeysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Retrieve values from plan
    var plan sdkKeysResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve values from state
    var state sdkKeysResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set the state to the current plan
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sdkKeysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state sdkKeysResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    projectID := state.ProjectID.ValueString()
    keyID := state.ID.ValueString()

    // Delete the SDK key using the DeleteSDKKey function
    err := r.client.DeleteSDKKey(ctx, projectID, keyID)
    if err != nil {
        // Check if the error message indicates a 404 Not Found status code
        if strings.Contains(err.Error(), "status code: 404") {
            return
        }
        resp.Diagnostics.AddError(
            "Error deleting SDK key",
            "Could not delete SDK key, unexpected error: "+err.Error(),
        )
        return
    }
}