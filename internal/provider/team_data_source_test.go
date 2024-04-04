// team_data_source_test.go
package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTeamDataSource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccTeamDataSourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("data.paragon_team.example", "id", "team-id"),
                    resource.TestCheckResourceAttr("data.paragon_team.example", "name", "Example Team"),
                ),
            },
        },
    })
}

const testAccTeamDataSourceConfig = `
data "paragon_team" "example" {
  id = "team-id"
}
`