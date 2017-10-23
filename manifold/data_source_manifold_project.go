package manifold

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	manifold "github.com/manifoldco/go-manifold"
)

func dataSourceManifoldProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManifoldProjectRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project label of which you want to retrieve the data.",
			},

			"resource": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A specific resource you want to load for this project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource label you want to get the resource for",
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
					},
				},
			},

			"credentials": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The credentials linked to this project.",
			},
		},
	}
}

func dataSourceManifoldProjectRead(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*clientWrapper)
	ctx := context.Background()

	project, err := cl.getProject(ctx, d.Get("project").(string))
	if err != nil {
		return err
	}

	resources, err := cl.getProjectResources(ctx, &project.ID)
	if err != nil {
		return err
	}

	resourceMap := map[string]*manifold.Resource{}
	resourceIDMap := map[string]string{}
	for _, res := range resources {
		resourceMap[res.Body.Label] = res
		resourceIDMap[res.ID.String()] = res.Body.Label
	}

	resourceList := d.Get("resource").(*schema.Set).List()
	if len(resourceList) != 0 {
		matched := []string{}

		// validate that all the resources we requested are available
		for _, raw := range resourceList {
			data := raw.(map[string]interface{})
			rl := data["resource"].(string)
			if _, ok := resourceMap[rl]; !ok {
				return fmt.Errorf("Resource '%s' not available", rl)
			}
			matched = append(matched, rl)
		}

		// filter out the resources that we don't need
		for k := range resourceMap {
			var found bool
			for _, v := range matched {
				if v == k {
					found = true
					break
				}
			}

			if !found {
				res := resourceMap[k]
				delete(resourceMap, k)
				delete(resourceIDMap, res.ID.String())
			}
		}
	}

	resourceIDs := []manifold.ID{}
	for _, res := range resourceMap {
		resourceIDs = append(resourceIDs, res.ID)
	}

	credentials, err := cl.getCredentials(ctx, resourceIDs)
	if err != nil {
		return err
	}

	resourceCredentials := map[string]map[string]string{}
	for _, cred := range credentials {
		resID := cred.Body.ResourceID.String()
		if _, ok := resourceCredentials[resID]; !ok {
			resourceCredentials[resID] = map[string]string{}
		}

		for k, v := range cred.Body.Values {
			resourceCredentials[resID][k] = v
		}
	}

	if len(resourceList) != 0 {
		for _, rawResource := range resourceList {
			rd := rawResource.(map[string]interface{})
			credList := rd["credential"].(*schema.Set).List()
			if len(credList) == 0 {
				break
			}
			resLabel := rd["resource"].(string)
			resID := resourceMap[resLabel].ID.String()
			availableCreds := resourceCredentials[resID]

			rcred := map[string]string{}
			for _, rawCred := range credList {
				name, value, err := parseCredential(rawCred, availableCreds)
				if err != nil {
					return err
				}
				rcred[name] = value
			}

			resourceCredentials[resID] = rcred
		}
	}

	credMap := map[string]interface{}{}
	for _, data := range resourceCredentials {
		for k, v := range data {
			if _, ok := credMap[k]; ok {
				return fmt.Errorf("Key '%s' is already used, please alias your credential", k)
			}

			credMap[k] = v
		}
	}

	d.SetId(project.ID.String())
	d.Set("project", project.Body.Label)
	d.Set("credentials", credMap)
	return nil
}
