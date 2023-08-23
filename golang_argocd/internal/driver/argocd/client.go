package argocd

import (
	"io"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/project"
)

type Connection struct {
	Address string
	Token   string
}

type Client struct {
	projectClient     project.ProjectServiceClient
	clusterClient     cluster.ClusterServiceClient
	applicationClient application.ApplicationServiceClient
	clientConnections [3]io.Closer
	indexToKey        map[int]string
}

func NewClient(c *Connection) (*Client, error) {
	apiClient, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: c.Address,
		Insecure:   true,
		AuthToken:  c.Token,
	})
	if err != nil {
		return nil, err
	}

	clientConnections := [3]io.Closer{}

	conn, projectClient, err := apiClient.NewProjectClient()
	if err != nil {
		return nil, err
	}

	clientConnections[0] = conn
	indexToKey := make(map[int]string, 3)
	indexToKey[0] = "project"

	conn, clusterClient, err := apiClient.NewClusterClient()
	if err != nil {
		return nil, err
	}

	clientConnections[1] = conn
	indexToKey[1] = "cluster"

	conn, applicationClient, err := apiClient.NewApplicationClient()
	if err != nil {
		return nil, err
	}

	clientConnections[2] = conn
	indexToKey[2] = "application"

	return &Client{
		projectClient:     projectClient,
		clusterClient:     clusterClient,
		applicationClient: applicationClient,
		clientConnections: clientConnections,
		indexToKey:        indexToKey,
	}, nil
}
