package provider

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/arielb135/terraform-provider-paragon/internal/client"
)

func TestAccOrganizationsDataSource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: providerConfig + `
data "paragon_organizations" "test" {}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify the number of organizations returned
                    resource.TestCheckResourceAttr("data.paragon_organizations.test", "organizations.#", "2"),
                    // Verify the attributes of the first organization
                    resource.TestCheckResourceAttr("data.paragon_organizations.test", "organizations.0.id", "org-1"),
                    resource.TestCheckResourceAttr("data.paragon_organizations.test", "organizations.0.name", "Organization 1"),
                    resource.TestCheckResourceAttr("data.paragon_organizations.test", "organizations.0.website", "https://org1.com"),
                    // Verify the attributes of the second organization
                    resource.TestCheckResourceAttr("data.paragon_organizations.test", "organizations.1.id", "org-2"),
                    resource.TestCheckResourceAttr("data.paragon_organizations.test", "organizations.1.name", "Organization 2"),
                    resource.TestCheckResourceAttr("data.paragon_organizations.test", "organizations.1.website", "https://org2.com"),
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
    GetOrganizationsFunc: func(ctx context.Context) ([]client.Organization, error) {
        return []client.Organization{
            {
                ID:      "org-1",
                Name:    "Organization 1",
                Website: "https://org1.com",
            },
            {
                ID:      "org-2",
                Name:    "Organization 2",
                Website: "https://org2.com",
            },
        }, nil
    },
}