package manifold

import (
	"github.com/hashicorp/terraform/helper/schema"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/integrations"
	"github.com/manifoldco/go-manifold/integrations/primitives"
)

func getProjectInformation(cl *integrations.Client, label string, required bool) (*string, *manifold.ID, error) {
	if !required && label == "" {
		return nil, nil, nil
	}

	project := ptrString(label)
	projectID, err := cl.ProjectID(project)
	if err != nil {
		return nil, nil, err
	}
	if projectID == nil {
		return nil, nil, errProjectNotFound
	}

	return project, projectID, nil
}

func ptrString(str string) *string {
	return &str
}

func resourcesFromList(resourceList []interface{}) []*primitives.Resource {
	resources := []*primitives.Resource{}

	for _, resource := range resourceList {
		rm := resource.(map[string]interface{})

		var creds []*primitives.Credential
		if c, ok := rm["credential"]; ok {
			creds = credentialsFromList(c.(*schema.Set).List())
		}

		rs := &primitives.Resource{
			Name:        rm["resource"].(string),
			Credentials: creds,
		}

		resources = append(resources, rs)
	}

	return resources
}

func credentialsFromList(credentialList []interface{}) []*primitives.Credential {
	credentials := []*primitives.Credential{}

	for _, credential := range credentialList {
		cred := credential.(map[string]interface{})

		cs := &primitives.Credential{
			Name:    cred["name"].(string),
			Key:     cred["key"].(string),
			Default: cred["default"].(string),
		}

		credentials = append(credentials, cs)
	}

	return credentials
}
