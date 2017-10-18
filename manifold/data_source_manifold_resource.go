package manifold

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	manifold "github.com/manifoldco/go-manifold"
)

func dataSourceManifoldResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManifoldResourceRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceManifoldResourceRead(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*clientWrapper)
	ctx := context.Background()

	projectLabel := d.Get("project").(string)
	project, err := cl.getProject(ctx, projectLabel)
	if err != nil {
		return err
	}

	resourceLabel := d.Get("resource").(string)
	resource, err := cl.getResource(ctx, project.ID, resourceLabel)
	if err != nil {
		return err
	}

	credentials, err := cl.getCredentials(ctx, []manifold.ID{resource.ID})
	if err != nil {
		return err
	}

	credentialList := d.Get("credential").(*schema.Set).List()
	credMap := map[string]string{}
	availableCredentials := map[string]string{}
	for _, cred := range credentials {
		for k, v := range cred.Body.Values {
			availableCredentials[k] = v
		}
	}

	// No credential filter set up, load all credentials
	if len(credentialList) == 0 {
		credMap = availableCredentials
	} else {
		for _, raw := range credentialList {
			k, v, err := parseCredential(raw, availableCredentials)
			if err != nil {
				return err
			}

			credMap[k] = v
		}
	}

	d.SetId(resource.ID.String())
	d.Set("project", project.Body.Label)
	d.Set("resource", resource.Body.Label)
	d.Set("credentials", credMap)
	return nil
}

func parseCredential(requestedCredential interface{}, creds map[string]string) (string, string, error) {
	credData := requestedCredential.(map[string]interface{})
	key := credData["key"].(string)

	var name string
	nameI, ok := credData["name"]
	if !ok || nameI.(string) == "" {
		name = key
	} else {
		name = nameI.(string)
	}

	for k, v := range creds {
		if k == key {
			return name, v, nil
		}
	}

	if def, ok := credData["default"]; ok && def.(string) != "" {
		return name, def.(string), nil
	}

	return "", "", errCredentialNotFound
}

func ptrString(s string) *string {
	return &s
}
