package provider

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ datasource.DataSource              = &organizationsDataSource{}
    _ datasource.DataSourceWithConfigure = &organizationsDataSource{}
)

// NewOrganizationsDataSource is a helper function to simplify the provider implementation.
func NewOrganizationsDataSource() datasource.DataSource {
    return &organizationsDataSource{}
}

// organizationsDataSource is the data source implementation.
type organizationsDataSource struct {
    client *client.Client
}

// organizationsDataSourceModel maps the data source schema data.
type organizationsDataSourceModel struct {
    Organizations []organizationModel `tfsdk:"organizations"`
}

// Configure adds the provider configured client to the data source.
func (d *organizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*client.Client)
    if !ok {
        tflog.Error(ctx, "Unable to prepare client")
        return
    }
    d.client = client
}

// Metadata returns the data source type name.
func (d *organizationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_organizations"
}

// Schema defines the schema for the data source.
func (d *organizationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Fetches a list of organizations.",
        Attributes: map[string]schema.Attribute{
            "organizations": schema.ListNestedAttribute{
                Description: "The list of organizations.",
                Computed:    true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: organizationAttributes(),
                },
            },
        },
    }
}

// Read refreshes the Terraform state with the latest data.
func (d *organizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    tflog.Debug(ctx, "Preparing to read organizations data source")
    var state organizationsDataSourceModel

    organizations, err := d.client.GetOrganizations(ctx)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Read Organizations",
            err.Error(),
        )
        return
    }

    var organizationModels []organizationModel
    for _, org := range organizations {
        organizationModels = append(organizationModels, mapOrganizationToModel(org))
    }

    state.Organizations = organizationModels

    // Set state
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    tflog.Debug(ctx, "Finished reading organizations data source", map[string]any{"success": true})
}