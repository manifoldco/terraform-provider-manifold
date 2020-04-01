package manifold

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/manifoldco/go-manifold/integrations"
	"github.com/manifoldco/go-manifold/integrations/primitives"
)

func dataSourceManifoldResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManifoldResourceRead,

		Schema: map[string]*schema.Schema{
			"resource": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource label you want to get the resource for",
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project label you want to get the resource for.",
			},

			"credential": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name for this credential which will be used to reference later on. Defaults to the `key` value.",
						},

						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key to fetch from the resource credentials.",
						},

						"default": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The default value for this key if it's not set.",
						},
					},
				},
			},

			"credentials": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceManifoldResourceRead(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	projectLabel, _, err := getProjectInformation(cl, d.Get("project").(string), false)
	if err != nil {
		return err
	}

	rs := &primitives.Resource{
		Name:        d.Get("resource").(string),
		Credentials: credentialsFromList(d.Get("credential").(*schema.Set).List()),
	}
	resource, err := cl.GetResource(ctx, projectLabel, rs)
	if err != nil {
		return err
	}

	cv, err := cl.GetResourceCredentialValues(ctx, projectLabel, rs)
	if err != nil {
		return err
	}

	credMap, err := integrations.FlattenResourceCredentialValues(cv)
	if err != nil {
		return err
	}

	d.SetId(resource.ID.String())

	if err = d.Set("resource", resource.Body.Label); err != nil {
		return err
	}

	if err = d.Set("credentials", credMap); err != nil {
		return err
	}

	if projectLabel != nil {
		if err = d.Set("project", *projectLabel); err != nil {
			return err
		}
	}
	return nil
}
