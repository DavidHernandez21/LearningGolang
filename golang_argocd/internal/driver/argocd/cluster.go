package argocd

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func (c *Client) GetClusters(ctx context.Context) ([]v1alpha1.Cluster, error) {
	cl, err := c.clusterClient.List(ctx, &cluster.ClusterQuery{})
	if err != nil {
		return nil, err
	}

	return cl.Items, nil
}
