package manifold

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/manifoldco/go-manifold"
)

// Provider returns the configured Terraform ResourceProvider with the Manifold
// reesource and data actions.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANIFOLD_URL_PATTERN", ""),
				Description: "The pattern used for connecting to Manifold's hosts",
			},
			"team": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANIFOLD_TEAM", ""),
				Description: "The team used to connect to the Manifold API",
			},
			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANIFOLD_API_TOKEN", ""),
				Description: "API Key to use to connect to the Manifold API",
				Sensitive:   true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"manifold_api_token":  resourceManifoldToken(),
			"manifold_credential": resourceManifoldCredential(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"manifold_resource":   dataSourceManifoldResource(),
			"manifold_project":    dataSourceManifoldProject(),
			"manifold_credential": dataSourceManifoldCredential(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cfgs := []manifold.ConfigFunc{}

	if apiToken, ok := d.GetOk("api_token"); ok && apiToken.(string) != "" {
		cfgs = append(cfgs, manifold.WithAPIToken(apiToken.(string)))
	} else {
		return nil, errAPITokenRequired
	}

	if urlPattern, ok := d.GetOk("url_pattern"); ok {
		cfgs = append(cfgs, manifold.ForURLPattern(urlPattern.(string)))
	}

	cfgs = append(cfgs, manifold.WithUserAgent(fmt.Sprintf("terraform-%s", Version)))

	wrapper, err := newWrapper(d.Get("team").(string), cfgs...)
	if err != nil {
		return nil, err
	}

	return wrapper, nil
}
