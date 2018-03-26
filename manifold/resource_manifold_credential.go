package manifold

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/integrations"
	"github.com/manifoldco/go-manifold/integrations/primitives"
)

func resourceManifoldCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceManifoldCredentialCreate,
		Read:   resourceManifoldResourceCredentialRead,
		Update: resourceManifoldCredentialUpdate,
		Delete: resourceManifoldCredentialDelete,

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

			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The value of the credential",
			},
		},
	}
}

func resourceManifoldCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	cfg, err := configFromSchema(ctx, cl, d)
	if err != nil {
		return err
	}

	if err := patchConfig(ctx, cl, cfg); err != nil {
		return err
	}

	d.SetId(cfg.key)
	return nil
}

func resourceManifoldCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	cfg, err := configFromSchema(ctx, cl, d)
	if err != nil {
		return err
	}

	if err := patchConfig(ctx, cl, cfg); err != nil {
		return err
	}

	d.SetId(cfg.key)
	return nil
}

func resourceManifoldCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	cfg, err := configFromSchema(ctx, cl, d)
	if err != nil {
		return err
	}
	cfg.value = nil

	return patchConfig(ctx, cl, cfg)
}

type schemaConfig struct {
	resourceID manifold.ID
	key        string
	value      *string
}

func configFromSchema(ctx context.Context, cl *integrations.Client, d *schema.ResourceData) (schemaConfig, error) {
	projectLabel, _, err := getProjectInformation(cl, d.Get("project").(string), false)
	if err != nil {
		return schemaConfig{}, err
	}

	res := &primitives.Resource{Name: d.Get("resource").(string)}
	resource, err := cl.GetResource(ctx, projectLabel, res)
	if err != nil {
		return schemaConfig{}, err
	}

	return schemaConfig{
		resourceID: resource.ID,
		key:        d.Get("key").(string),
		value:      ptrString(d.Get("value").(string)),
	}, nil
}

func patchConfig(ctx context.Context, cl *integrations.Client, cfg schemaConfig) error {
	data := map[string]interface{}{
		cfg.key: cfg.value,
	}
	_, err := cl.Resources.UpdateConfig(ctx, cfg.resourceID, &data)
	return err
}
