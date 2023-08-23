package argocd

import (
	"fmt"
	"strings"
)

// closeConnections closes all the connections to the ArgoCD API server.
func (c *Client) CloseConnections() {
	stringBuilder := strings.Builder{}
	for key := range c.clientConnections {
		err := c.clientConnections[key].Close()
		if err != nil {
			stringBuilder.WriteString("error closing connection to ")
			stringBuilder.WriteString(c.indexToKey[key])
			stringBuilder.WriteString(": ")
			stringBuilder.WriteString(err.Error())
			stringBuilder.WriteString("\n")
		}
	}
	fmt.Println(stringBuilder.String())
}
