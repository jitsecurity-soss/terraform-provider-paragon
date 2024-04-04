package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTeamsDataSource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccTeamsDataSourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("data.paragon_teams.all", "teams.#", "2"),
                    resource.TestCheckResourceAttr("data.paragon_teams.all", "teams.0.id", "team-1"),
                    resource.TestCheckResourceAttr("data.paragon_teams.all", "teams.0.name", "Team 1"),
                    resource.TestCheckResourceAttr("data.paragon_teams.all", "teams.1.id", "team-2"),
                    resource.TestCheckResourceAttr("data.paragon_teams.all", "teams.1.name", "Team 2"),
                ),
            },
        },
    })
}

const testAccTeamsDataSourceConfig = `
data "paragon_teams" "all" {}
`