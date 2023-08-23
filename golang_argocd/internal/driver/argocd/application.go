package argocd

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func (c *Client) GetApplication(ctx context.Context, name string) (*v1alpha1.Application, error) {
	return c.applicationClient.Get(ctx, &application.ApplicationQuery{
		Name: &name,
	})
}
