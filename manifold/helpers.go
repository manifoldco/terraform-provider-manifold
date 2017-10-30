package manifold

import (
	"github.com/hashicorp/terraform/helper/schema"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/kubernetes-credentials/helpers/client"
	"github.com/manifoldco/kubernetes-credentials/primitives"
)

func getProjectInformation(cl *client.Client, label string, required bool) (*string, *manifold.ID, error) {
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

func resourceSpecsFromList(resourceList []interface{}) []*primitives.ResourceSpec {
	resources := []*primitives.ResourceSpec{}

	for _, resource := range resourceList {
		rm := resource.(map[string]interface{})

		var creds []*primitives.CredentialSpec
		if c, ok := rm["credential"]; ok {
			creds = credentialSpecsFromList(c.(*schema.Set).List())
		}

		rs := &primitives.ResourceSpec{
			Name:        rm["resource"].(string),
			Credentials: creds,
		}

		resources = append(resources, rs)
	}

	return resources
}

func credentialSpecsFromList(credentialList []interface{}) []*primitives.CredentialSpec {
	credentials := []*primitives.CredentialSpec{}

	for _, credential := range credentialList {
		cred := credential.(map[string]interface{})

		cs := &primitives.CredentialSpec{
			Name:    cred["name"].(string),
			Key:     cred["key"].(string),
			Default: cred["default"].(string),
		}

		credentials = append(credentials, cs)
	}

	return credentials
}
