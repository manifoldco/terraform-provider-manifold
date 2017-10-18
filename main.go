package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/manifoldco/terraform-provider-manifold/manifold"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: manifold.Provider})
}
