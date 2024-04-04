// cli_key_resource_test.go
package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCLIKeyResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccCLIKeyResourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttrSet("paragon_cli_key.example", "id"),
                    resource.TestCheckResourceAttr("paragon_cli_key.example", "name", "Example CLI Key"),
                ),
            },
        },
    })
}

const testAccCLIKeyResourceConfig = `
resource "paragon_cli_key" "example" {
  organization_id = "org-id"
  name            = "Example CLI Key"
}
`