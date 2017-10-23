package manifold

import (
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
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANIFOLD_API_KEY", ""),
				Description: "API Key to use to connect to the Manifold API",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"manifold_resource": dataSourceManifoldResource(),
			"manifold_project":  dataSourceManifoldProject(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cfgs := []manifold.ConfigFunc{}

	if apiKey, ok := d.GetOk("api_key"); ok && apiKey.(string) != "" {
		cfgs = append(cfgs, manifold.WithAPIKey(apiKey.(string)))
	} else {
		return nil, errAPIKeyRequired
	}

	if urlPattern, ok := d.GetOk("url_pattern"); ok {
		cfgs = append(cfgs, manifold.ForURLPattern(urlPattern.(string)))
	}

	cl := manifold.New(cfgs...)
	wrapper := &clientWrapper{client: cl, team: d.Get("team").(string)}
	if err := wrapper.EnsureTeamID(); err != nil {
		return nil, err
	}

	return wrapper, nil
}
