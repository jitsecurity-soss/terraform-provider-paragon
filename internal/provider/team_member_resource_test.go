// team_member_resource_test.go
package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTeamMemberResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccTeamMemberResourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttrSet("paragon_team_member.example", "id"),
                    resource.TestCheckResourceAttr("paragon_team_member.example", "email", "user@example.com"),
                    resource.TestCheckResourceAttr("paragon_team_member.example", "role", "MEMBER"),
                ),
            },
        },
    })
}

const testAccTeamMemberResourceConfig = `
resource "paragon_team_member" "example" {
  team_id = "team-id"
  email   = "user@example.com"
  role    = "MEMBER"
}
`
