package argocd

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateProject(ctx context.Context, name string) (*v1alpha1.AppProject, error) {
	return c.projectClient.Create(ctx, &project.ProjectCreateRequest{
		Project: &v1alpha1.AppProject{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
	})
}

func (c *Client) DeleteProject(ctx context.Context, name string) error {
	_, err := c.projectClient.Delete(ctx, &project.ProjectQuery{
		Name: name,
	})

	return err
}

func (c *Client) GetProject(ctx context.Context, name string) (*v1alpha1.AppProject, error) {
	return c.projectClient.Get(ctx, &project.ProjectQuery{
		Name: name,
	})
}
