package provider

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
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

type organizationModel struct {
    ID                    types.String `tfsdk:"id"`
    DateCreated           types.String `tfsdk:"date_created"`
    DateUpdated           types.String `tfsdk:"date_updated"`
    Name                  types.String `tfsdk:"name"`
    Website               types.String `tfsdk:"website"`
    Type                  types.String `tfsdk:"type"`
    Purpose               types.String `tfsdk:"purpose"`
    Referral              types.String `tfsdk:"referral"`
    Size                  types.String `tfsdk:"size"`
    Role                  types.String `tfsdk:"role"`
    CompletedQualification types.Bool   `tfsdk:"completed_qualification"`
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
                    Attributes: map[string]schema.Attribute{
                        "id": schema.StringAttribute{
                            Description: "Identifier for the organization.",
                            Computed:    true,
                        },
                        "date_created": schema.StringAttribute{
                            Description: "The creation date of the organization.",
                            Computed:    true,
                        },
                        "date_updated": schema.StringAttribute{
                            Description: "The last update date of the organization.",
                            Computed:    true,
                        },
                        "name": schema.StringAttribute{
                            Description: "The name of the organization.",
                            Computed:    true,
                        },
                        "website": schema.StringAttribute{
                            Description: "The website of the organization.",
                            Computed:    true,
                        },
                        "type": schema.StringAttribute{
                            Description: "The type of the organization.",
                            Computed:    true,
                        },
                        "purpose": schema.StringAttribute{
                            Description: "The purpose of the organization.",
                            Computed:    true,
                        },
                        "referral": schema.StringAttribute{
                            Description: "The referral of the organization.",
                            Computed:    true,
                        },
                        "size": schema.StringAttribute{
                            Description: "The size of the organization.",
                            Computed:    true,
                        },
                        "role": schema.StringAttribute{
                            Description: "The role of the organization.",
                            Computed:    true,
                        },
                        "completed_qualification": schema.BoolAttribute{
                            Description: "Indicates if the organization has completed qualification.",
                            Computed:    true,
                        },
                    },
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
        organizationModels = append(organizationModels, organizationModel{
            ID:                    types.StringValue(org.ID),
            DateCreated:           types.StringValue(org.DateCreated),
            DateUpdated:           types.StringValue(org.DateUpdated),
            Name:                  types.StringValue(org.Name),
            Website:               types.StringValue(org.Website),
            Type:                  types.StringValue(org.Type),
            Purpose:               types.StringValue(org.Purpose),
            Referral:              types.StringValue(org.Referral),
            Size:                  types.StringValue(org.Size),
            Role:                  types.StringValue(org.Role),
            CompletedQualification: types.BoolValue(org.CompletedQualification),
        })
    }

    state.Organizations = organizationModels

    // Set state
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
    tflog.Debug(ctx, "Finished reading organizations data source", map[string]any{"success": true})
}