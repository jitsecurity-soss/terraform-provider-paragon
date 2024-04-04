// environment_secret_resource_test.go
package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEnvironmentSecretResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccEnvironmentSecretResourceConfig,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttrSet("paragon_environment_secret.example", "id"),
                    resource.TestCheckResourceAttr("paragon_environment_secret.example", "key", "SECRET_KEY"),
                ),
            },
        },
    })
}

const testAccEnvironmentSecretResourceConfig = `
resource "paragon_project" "test" {
  organization_id = "org-id"
  name            = "Test Project"
}

resource "paragon_environment_secret" "example" {
  project_id = paragon_project.test.id
  key        = "SECRET_KEY"
  value      = "secret_value"
}
`