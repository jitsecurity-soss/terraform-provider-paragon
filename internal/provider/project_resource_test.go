// project_resource_test.go
package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccProjectResourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("paragon_project.example", "name", "Example Project"),
                    resource.TestCheckResourceAttrSet("paragon_project.example", "id"),
                ),
            },
        },
    })
}

const testAccProjectResourceConfig = `
resource "paragon_project" "example" {
  organization_id = "org-id"
  name            = "Example Project"
}
`
