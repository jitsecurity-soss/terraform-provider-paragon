// sdk_keys_resource_test.go
package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSDKKeysResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccSDKKeysResourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttrSet("paragon_sdk_keys.example", "id"),
                    resource.TestCheckResourceAttr("paragon_sdk_keys.example", "version", "1.0.0"),
                ),
            },
        },
    })
}

const testAccSDKKeysResourceConfig = `
resource "paragon_project" "test" {
  organization_id = "org-id"
  name            = "Test Project"
}

resource "paragon_sdk_keys" "example" {
  project_id = paragon_project.test.id
  version    = "1.0.0"
}
`