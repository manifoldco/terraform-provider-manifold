package manifold_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/manifoldco/kubernetes-credentials/helpers/client"
)

func TestManifoldResource_APIToken(t *testing.T) {
	conf := `
resource "manifold_api_token" "test_token" {
  team        = "manifold-integration-ci"
  role        = "read"
  description = "New token - Terraform CI Tests"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:     testProviderPreCheck(t),
		Providers:    testProviders,
		CheckDestroy: testAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("manifold_api_token.test_token", "id"),
					resource.TestCheckResourceAttrSet("manifold_api_token.test_token", "token"),
				),
			},
		},
	})
}

func testAPITokenDestroy(s *terraform.State) error {
	cl := testProviders["manifold"].(*schema.Provider).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "manifold_api_token" {
			continue
		}

		tokensList := cl.Tokens.List(context.Background(), "api", nil)
		defer tokensList.Close()

		for tokensList.Next() {
			token, err := tokensList.Current()
			if err != nil {
				return err
			}

			if token.ID.String() == rs.Primary.ID {
				return fmt.Errorf("Token '%s' still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}
