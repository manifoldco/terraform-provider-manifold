package manifold

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/manifoldco/kubernetes-credentials/helpers/client"
	"github.com/manifoldco/kubernetes-credentials/primitives"
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
	cl := meta.(*client.Client)
	ctx := context.Background()

	projectLabel, _, err := getProjectInformation(cl, d.Get("project").(string), false)
	if err != nil {
		return err
	}

	rs := &primitives.ResourceSpec{
		Name:        d.Get("resource").(string),
		Credentials: credentialSpecsFromList(d.Get("credential").(*schema.Set).List()),
	}
	resource, err := cl.GetResource(ctx, projectLabel, rs)
	if err != nil {
		return err
	}

	cv, err := cl.GetResourceCredentialValues(ctx, projectLabel, rs)
	if err != nil {
		return err
	}

	credMap, err := client.FlattenResourceCredentialValues(cv)
	if err != nil {
		return err
	}

	d.SetId(resource.ID.String())
	d.Set("resource", resource.Body.Label)
	d.Set("credentials", credMap)
	if projectLabel != nil {
		d.Set("project", *projectLabel)
	}
	return nil
}
