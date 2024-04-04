package provider

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

func TestAccTeamsDataSource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: providerConfig + `
data "paragon_teams" "test" {}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify the number of teams returned
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.#", "2"),
                    // Verify the attributes of the first team
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.0.id", "team-1"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.0.name", "Team 1"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.0.website", "https://team1.com"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.0.organization_id", "org-1"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.0.organization.id", "org-1"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.0.organization.name", "Organization 1"),
                    // Verify the attributes of the second team
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.1.id", "team-2"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.1.name", "Team 2"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.1.website", "https://team2.com"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.1.organization_id", "org-2"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.1.organization.id", "org-2"),
                    resource.TestCheckResourceAttr("data.paragon_teams.test", "teams.1.organization.name", "Organization 2"),
                ),
            },
        },
    })
}

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
    return map[string]func() (tfprotov6.ProviderServer, error){
        "paragon": func() (tfprotov6.ProviderServer, error) {
            return providermocks.NewMockProviderWithClient(context.Background(), mockClient), nil
        },
    }
}

var mockClient = &client.Client{
    GetTeamsFunc: func(ctx context.Context) ([]client.Team, error) {
        return []client.Team{
            {
                ID:             "team-1",
                Name:           "Team 1",
                Website:        "https://team1.com",
                OrganizationID: "org-1",
                Organization: client.Organization{
                    ID:   "org-1",
                    Name: "Organization 1",
                },
            },
            {
                ID:             "team-2",
                Name:           "Team 2",
                Website:        "https://team2.com",
                OrganizationID: "org-2",
                Organization: client.Organization{
                    ID:   "org-2",
                    Name: "Organization 2",
                },
            },
        }, nil
    },
}