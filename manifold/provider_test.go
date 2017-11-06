package manifold_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"github.com/manifoldco/terraform-provider-manifold/manifold"
)

var (
	testProviders map[string]terraform.ResourceProvider
)

func init() {
	testProviders = map[string]terraform.ResourceProvider{
		"manifold": manifold.Provider(),
	}
}

func testProviderPreCheck(t *testing.T) func() {
	return func() {
		if os.Getenv("MANIFOLD_API_TOKEN") == "" {
			t.Fatal("`MANIFOLD_API_TOKEN` must be set to run the provider tests")
		}
	}
}
