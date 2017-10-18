package manifold

import (
	"context"

	manifold "github.com/manifoldco/go-manifold"
)

type clientWrapper struct {
	client *manifold.Client
	team   string
	teamID *manifold.ID
}

func (cl *clientWrapper) EnsureTeamID() error {
	if cl.team == "" {
		// no team specified, skip it
		return nil
	}

	teamsList := cl.client.Teams.List(context.Background())
	defer teamsList.Close()

	for teamsList.Next() {
		team, err := teamsList.Current()
		if err != nil {
			return err
		}

		if team.Body.Label == cl.team {
			cl.teamID = &team.ID
			return nil
		}
	}

	return errTeamNotFound
}

func (cl *clientWrapper) getProject(ctx context.Context, label string) (*manifold.Project, error) {
	projects := cl.client.Projects.List(ctx, &manifold.ProjectsListOpts{
		Label:  ptrString(label),
		TeamID: cl.teamID,
	})
	defer projects.Close()

	for projects.Next() {
		project, err := projects.Current()
		if err != nil {
			return nil, err
		}

		if project.Body.Label == label {
			return project, nil
		}
	}

	return nil, errProjectNotFound
}

func (cl *clientWrapper) getResource(ctx context.Context, projectID manifold.ID, label string) (*manifold.Resource, error) {
	resources := cl.client.Resources.List(ctx, &manifold.ResourcesListOpts{
		ProjectID: &projectID,
		Label:     ptrString(label),
		TeamID:    cl.teamID,
	})
	defer resources.Close()

	for resources.Next() {
		resource, err := resources.Current()
		if err != nil {
			return nil, err
		}

		if resource.Body.Label == label {
			return resource, nil
		}
	}

	return nil, errResourceNotFound
}

func (cl *clientWrapper) getCredentials(ctx context.Context, resourceIDs []manifold.ID) ([]*manifold.Credential, error) {
	credList := cl.client.Credentials.List(ctx, resourceIDs)
	defer credList.Close()

	credentials := []*manifold.Credential{}
	for credList.Next() {
		cred, err := credList.Current()
		if err != nil {
			return nil, err
		}

		credentials = append(credentials, cred)
	}

	return credentials, nil
}
