// teams_data_source.go
package provider

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ datasource.DataSource              = &teamsDataSource{}
    _ datasource.DataSourceWithConfigure = &teamsDataSource{}
)

// NewTeamsDataSource is a helper function to simplify the provider implementation.
func NewTeamsDataSource() datasource.DataSource {
    return &teamsDataSource{}
}

// teamsDataSource is the data source implementation.
type teamsDataSource struct {
    client *client.Client
}

// teamsDataSourceModel maps the data source schema data.
type teamsDataSourceModel struct {
    Teams []teamModel `tfsdk:"teams"`
}

type teamModel struct {
    ID             types.String       `tfsdk:"id"`
    DateCreated    types.String       `tfsdk:"date_created"`
    DateUpdated    types.String       `tfsdk:"date_updated"`
    Name           types.String       `tfsdk:"name"`
    Website        types.String       `tfsdk:"website"`
    OrganizationID types.String       `tfsdk:"organization_id"`
    Organization   *organizationModel `tfsdk:"organization"`
}

// Configure adds the provider configured client to the data source.
func (d *teamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*client.Client)
    if !ok {
        return
    }
    d.client = client
}

// Metadata returns the data source type name.
func (d *teamsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_teams"
}

// Schema defines the schema for the data source.
func (d *teamsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Fetches a list of teams.",
        Attributes: map[string]schema.Attribute{
            "teams": schema.ListNestedAttribute{
                Description: "The list of teams.",
                Computed:    true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.StringAttribute{
                            Description: "Identifier for the team.",
                            Computed:    true,
                        },
                        "date_created": schema.StringAttribute{
                            Description: "The creation date of the team.",
                            Computed:    true,
                        },
                        "date_updated": schema.StringAttribute{
                            Description: "The last update date of the team.",
                            Computed:    true,
                        },
                        "name": schema.StringAttribute{
                            Description: "The name of the team.",
                            Computed:    true,
                        },
                        "website": schema.StringAttribute{
                            Description: "The website of the team.",
                            Computed:    true,
                        },
                        "organization_id": schema.StringAttribute{
                            Description: "The ID of the organization the team belongs to.",
                            Computed:    true,
                        },
                        "organization": schema.SingleNestedAttribute{
                            Description: "The organization the team belongs to.",
                            Computed:    true,
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
            },
        },
    }
}

// Read refreshes the Terraform state with the latest data.
func (d *teamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state teamsDataSourceModel

    teams, err := d.client.GetTeams(ctx)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Read Teams",
            err.Error(),
        )
        return
    }

    var teamModels []teamModel
    for _, team := range teams {
        teamModels = append(teamModels, teamModel{
            ID:             types.StringValue(team.ID),
            DateCreated:    types.StringValue(team.DateCreated),
            DateUpdated:    types.StringValue(team.DateUpdated),
            Name:           types.StringValue(team.Name),
            Website:        types.StringValue(team.Website),
            OrganizationID: types.StringValue(team.OrganizationID),
            Organization: &organizationModel{
                ID:                    types.StringValue(team.Organization.ID),
                DateCreated:           types.StringValue(team.Organization.DateCreated),
                DateUpdated:           types.StringValue(team.Organization.DateUpdated),
                Name:                  types.StringValue(team.Organization.Name),
                Website:               types.StringValue(team.Organization.Website),
                Type:                  types.StringValue(team.Organization.Type),
                Purpose:               types.StringValue(team.Organization.Purpose),
                Referral:              types.StringValue(team.Organization.Referral),
                Size:                  types.StringValue(team.Organization.Size),
                Role:                  types.StringValue(team.Organization.Role),
                CompletedQualification: types.BoolValue(team.Organization.CompletedQualification),
            },
        })
    }

    state.Teams = teamModels

    // Set state
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}