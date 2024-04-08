package provider

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ datasource.DataSource              = &teamDataSource{}
    _ datasource.DataSourceWithConfigure = &teamDataSource{}
)

// NewTeamDataSource is a helper function to simplify the provider implementation.
func NewTeamDataSource() datasource.DataSource {
    return &teamDataSource{}
}

// teamDataSource is the data source implementation.
type teamDataSource struct {
    client *client.Client
}

// teamDataSourceModel maps the data source schema data.
type teamDataSourceModel struct {
    ID             types.String `tfsdk:"id"`
    DateCreated    types.String `tfsdk:"date_created"`
    DateUpdated    types.String `tfsdk:"date_updated"`
    Name           types.String `tfsdk:"name"`
    Website        types.String `tfsdk:"website"`
    OrganizationID types.String `tfsdk:"organization_id"`
}

// Configure adds the provider configured client to the data source.
func (d *teamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
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
func (d *teamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the data source.
func (d *teamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Fetches a team by its name.",
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
                Description: "The name of the team.",
                Required:    true,
            },
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
            "website": schema.StringAttribute{
                Description: "The website of the team.",
                Computed:    true,
            },
            "organization_id": schema.StringAttribute{
                Description: "The ID of the organization the team belongs to.",
                Computed:    true,
            },
        },
    }
}

// Read refreshes the Terraform state with the latest data.
func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var config teamDataSourceModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    teamName := config.Name.ValueString()

    // Retrieve all teams using the GetTeams function.
    teams, err := d.client.GetTeams(ctx)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Read Teams",
            err.Error(),
        )
        return
    }

    // Find the team with the matching name.
    var foundTeam *client.Team
    for _, team := range teams {
        if team.Name == teamName {
            foundTeam = &team
            break
        }
    }

    if foundTeam == nil {
        resp.Diagnostics.AddError(
            "Team Not Found",
            fmt.Sprintf("Team with name '%s' not found", teamName),
        )
        return
    }

    // Map the team data to the state.
    state := teamDataSourceModel{
        ID:             types.StringValue(foundTeam.ID),
        DateCreated:    types.StringValue(foundTeam.DateCreated),
        DateUpdated:    types.StringValue(foundTeam.DateUpdated),
        Name:           types.StringValue(foundTeam.Name),
        Website:        types.StringValue(foundTeam.Website),
        OrganizationID: types.StringValue(foundTeam.OrganizationID),
    }

    // Set state
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}