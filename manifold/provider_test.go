package manifold

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

var (
	testProviders map[string]terraform.ResourceProvider
	testProvider  clientWrapper
)

func init() {
	testProviders = map[string]terraform.ResourceProvider{
		"manifold": Provider(),
	}
}

func testProviderPreCheck(t *testing.T) func() {
	return func() {
		if os.Getenv("MANIFOLD_API_KEY") == "" {
			t.Fatal("`MANIFOLD_API_KEY` must be set to run the provider tests")
		}
	}
}
