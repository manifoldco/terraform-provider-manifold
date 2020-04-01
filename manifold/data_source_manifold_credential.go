package manifold

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/manifoldco/go-manifold/integrations"
	"github.com/manifoldco/go-manifold/integrations/primitives"
)

func dataSourceManifoldCredential() *schema.Resource {
	return &schema.Resource{
		Read: resourceManifoldResourceCredentialRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project label you want to get the resource for.",
			},

			"resource": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource label you want to get the resource for",
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

			"value": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "The value of the credential",
			},
		},
	}
}

func resourceManifoldResourceCredentialRead(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	projectLabel, _, err := getProjectInformation(cl, d.Get("project").(string), false)
	if err != nil {
		return err
	}

	var defValue string
	def, ok := d.GetOkExists("default")
	if ok {
		defValue = def.(string)
	}

	rs := &primitives.Resource{
		Name: d.Get("resource").(string),
		Credentials: []*primitives.Credential{
			{
				Key:     d.Get("key").(string),
				Default: defValue,
			},
		},
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
	cred := credMap[d.Get("key").(string)]

	d.SetId(resource.ID.String())

	if err = d.Set("resource", resource.Body.Label); err != nil {
		return err
	}

	if err = d.Set("value", cred); err != nil {
		return err
	}

	if projectLabel != nil {
		if err = d.Set("project", *projectLabel); err != nil {
			return err
		}
	}
	return nil
}
