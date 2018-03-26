package manifold_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/manifoldco/go-manifold/integrations"
	"github.com/manifoldco/go-manifold/integrations/primitives"
)

func TestManifoldResource_CustomCredential(t *testing.T) {
	conf := `
resource "manifold_credential" "test_credential" {
	project = "terraform"
	resource = "custom-resource3-1"
	key = "NEW_KEY"
	value = "my-value"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:     testProviderPreCheck(t),
		Providers:    testProviders,
		CheckDestroy: testAPICredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("manifold_credential.test_credential", "id"),
					resource.TestCheckResourceAttrSet("manifold_credential.test_credential", "value"),
				),
			},
		},
	})
}

func testAPICredentialDestroy(s *terraform.State) error {
	cl := testProviders["manifold"].(*schema.Provider).Meta().(*integrations.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "manifold_credential" {
			continue
		}
		ctx := context.Background()

		projectLabel := rs.Primary.Attributes["project"]
		resourceLabel := rs.Primary.Attributes["resource"]

		res := &primitives.Resource{Name: resourceLabel}
		resource, err := cl.GetResource(ctx, &projectLabel, res)
		if err != nil {
			return err
		}

		credentials, err := cl.Resources.GetConfig(ctx, resource.ID)
		if err != nil {
			return err
		}

		if credentials == nil {
			return nil
		}

		creds := *credentials
		if _, ok := creds[rs.Primary.Attributes["key"]]; ok {
			return fmt.Errorf("Credential '%s' still exists", rs.Primary.ID)
		}
	}

	return nil
}
