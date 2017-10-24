package manifold

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestManifoldDataSource_ResourceBasic(t *testing.T) {
	t.Run("with a basic configuration", func(t *testing.T) {
		conf := `
data "manifold_resource" "my-resource" {
	project = "manifold-terraform"
	resource = "custom-resource1"
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:  testProviderPreCheck(t),
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: conf,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.manifold_resource.my-resource", "id"),
						resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "project", "manifold-terraform"),
						resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "resource", "custom-resource1"),
						testAccCheckManifoldCredentialsLength("data.manifold_resource.my-resource", 2),
						testAccCheckManifoldCredentialSet("data.manifold_resource.my-resource", "TOKEN_ID"),
						testAccCheckManifoldCredentialValue("data.manifold_resource.my-resource", "TOKEN_ID", "my-secret-token-id"),
						testAccCheckManifoldCredentialSet("data.manifold_resource.my-resource", "TOKEN_SECRET"),
						testAccCheckManifoldCredentialValue("data.manifold_resource.my-resource", "TOKEN_SECRET", "my-secret-token-secret"),
					),
				},
			},
		})
	})

	t.Run("with a filtered configuration", func(t *testing.T) {
		conf := `
data "manifold_resource" "my-resource" {
	project = "manifold-terraform"
	resource = "custom-resource1"

	credential {
		key = "TOKEN_ID"
	}
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:  testProviderPreCheck(t),
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: conf,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.manifold_resource.my-resource", "id"),
						resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "project", "manifold-terraform"),
						resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "resource", "custom-resource1"),
						testAccCheckManifoldCredentialsLength("data.manifold_resource.my-resource", 1),
						testAccCheckManifoldCredentialSet("data.manifold_resource.my-resource", "TOKEN_ID"),
						testAccCheckManifoldCredentialValue("data.manifold_resource.my-resource", "TOKEN_ID", "my-secret-token-id"),
					),
				},
			},
		})
	})

	t.Run("with a named key", func(t *testing.T) {
		conf := `
data "manifold_resource" "my-resource" {
	project = "manifold-terraform"
	resource = "custom-resource1"

	credential {
		name = "my-token"
		key = "TOKEN_ID"
	}
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:  testProviderPreCheck(t),
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: conf,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.manifold_resource.my-resource", "id"),
						resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "project", "manifold-terraform"),
						resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "resource", "custom-resource1"),
						testAccCheckManifoldCredentialsLength("data.manifold_resource.my-resource", 1),
						testAccCheckManifoldCredentialSet("data.manifold_resource.my-resource", "my-token"),
					),
				},
			},
		})
	})

	t.Run("with a non existing key", func(t *testing.T) {
		t.Run("with a default value", func(t *testing.T) {
			conf := `
data "manifold_resource" "my-resource" {
	project = "manifold-terraform"
	resource = "custom-resource1"

	credential {
		key = "NON_EXISTING_FIELD"
		default = "my-value"
	}
}
`

			resource.Test(t, resource.TestCase{
				PreCheck:  testProviderPreCheck(t),
				Providers: testProviders,
				Steps: []resource.TestStep{
					{
						Config: conf,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet("data.manifold_resource.my-resource", "id"),
							resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "project", "manifold-terraform"),
							resource.TestCheckResourceAttr("data.manifold_resource.my-resource", "resource", "custom-resource1"),
							testAccCheckManifoldCredentialsLength("data.manifold_resource.my-resource", 1),
							testAccCheckManifoldCredentialSet("data.manifold_resource.my-resource", "NON_EXISTING_FIELD"),
							testAccCheckManifoldCredentialValue("data.manifold_resource.my-resource", "NON_EXISTING_FIELD", "my-value"),
						),
					},
				},
			})
		})
	})
}

func testAccCheckManifoldCredentialsLength(n string, length int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}

		out, ok := rn.Primary.Attributes["credentials.%"]
		if !ok {
			return fmt.Errorf("Attribute 'credentials' not found: %#v", rn.Primary.Attributes)
		}

		o, _ := strconv.Atoi(out)
		if o != length {
			return fmt.Errorf("Attribute 'credentials' should be of length '%d', got '%s'", length, out)
		}
		return nil
	}
}

func testAccCheckManifoldCredentialSet(n, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		_, ok = rn.Primary.Attributes["credentials."+attr]
		if !ok {
			return fmt.Errorf("Attribute '%s' not found: %#v", attr, rn.Primary.Attributes)
		}
		return nil
	}
}

func testAccCheckManifoldCredentialValue(n, attr, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		out, ok := rn.Primary.Attributes["credentials."+attr]
		if !ok {
			return fmt.Errorf("Attribute '%s' not found: %#v", attr, rn.Primary.Attributes)
		}
		if out != value {
			return fmt.Errorf("Attribute '%s' should be value '%s', got '%s'", attr, value, out)
		}
		return nil
	}
}
