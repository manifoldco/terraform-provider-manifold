package manifold_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestManifoldDataSource_CredentialBasic(t *testing.T) {
	t.Run("with a filtered configuration", func(t *testing.T) {
		conf := `
data "manifold_credential" "my-credential" {
	project = "terraform"
	resource = "custom-resource1-1"
	key = "TOKEN_ID"
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:  testProviderPreCheck(t),
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: conf,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.manifold_credential.my-credential", "id"),
						resource.TestCheckResourceAttr("data.manifold_credential.my-credential", "project", "terraform"),
						resource.TestCheckResourceAttr("data.manifold_credential.my-credential", "resource", "custom-resource1-1"),
						testAccCheckManifoldCredential("data.manifold_credential.my-credential", "my-secret-token-id"),
					),
				},
			},
		})
	})

	t.Run("with a non existing key", func(t *testing.T) {
		t.Run("with a default value", func(t *testing.T) {
			conf := `
data "manifold_credential" "my-credential" {
	project = "terraform"
	resource = "custom-resource1-1"
	key = "NON_EXISTING_FIELD"
	default = "my-value"
}
`

			resource.Test(t, resource.TestCase{
				PreCheck:  testProviderPreCheck(t),
				Providers: testProviders,
				Steps: []resource.TestStep{
					{
						Config: conf,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet("data.manifold_credential.my-credential", "id"),
							resource.TestCheckResourceAttr("data.manifold_credential.my-credential", "project", "terraform"),
							resource.TestCheckResourceAttr("data.manifold_credential.my-credential", "resource", "custom-resource1-1"),
							testAccCheckManifoldCredential("data.manifold_credential.my-credential", "my-value"),
						),
					},
				},
			})
		})
	})
}

func testAccCheckManifoldCredential(n, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		out, ok := rn.Primary.Attributes["value"]
		if !ok {
			return fmt.Errorf("Attribute 'value' not found: %#v", rn.Primary.Attributes)
		}
		if out != value {
			return fmt.Errorf("Value should be '%s', got '%s'", value, out)
		}
		return nil
	}
}
