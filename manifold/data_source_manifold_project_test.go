package manifold

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestManifoldDataSource_Project(t *testing.T) {
	t.Run("with a basic configuration", func(t *testing.T) {
		conf := `
data "manifold_project" "manifold-terraform" {
	project = "terraform"
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:  testProviderPreCheck(t),
			Providers: testProviders,
			Steps: []resource.TestStep{
				{
					Config: conf,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.manifold_project.manifold-terraform", "id"),
						resource.TestCheckResourceAttr("data.manifold_project.manifold-terraform", "project", "terraform"),
						testAccCheckManifoldCredentialsLength("data.manifold_project.manifold-terraform", 4),
					),
				},
			},
		})
	})

	t.Run("with a selected project", func(t *testing.T) {
		t.Run("without credentials filtered", func(t *testing.T) {
			conf := `
data "manifold_project" "manifold-terraform" {
	project = "terraform"

	resource {
		resource = "custom-resource1-1"
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
							resource.TestCheckResourceAttrSet("data.manifold_project.manifold-terraform", "id"),
							resource.TestCheckResourceAttr("data.manifold_project.manifold-terraform", "project", "terraform"),
							testAccCheckManifoldCredentialsLength("data.manifold_project.manifold-terraform", 2),
						),
					},
				},
			})
		})

		t.Run("with credentials filtered", func(t *testing.T) {
			conf := `
data "manifold_project" "manifold-terraform" {
	project = "terraform"

	resource {
		resource = "custom-resource1-1"

		credential {
			key = "TOKEN_ID"
		}
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
							resource.TestCheckResourceAttrSet("data.manifold_project.manifold-terraform", "id"),
							resource.TestCheckResourceAttr("data.manifold_project.manifold-terraform", "project", "terraform"),
							testAccCheckManifoldCredentialsLength("data.manifold_project.manifold-terraform", 1),
						),
					},
				},
			})
		})
	})
}
