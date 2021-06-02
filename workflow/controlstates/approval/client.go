package approval

import (
	"fmt"
	"net/http"

	"github.com/cased/cased-go"
)

type Client struct {
	Endpoint cased.Endpoint
}

func client() Client {
	return Client{
		Endpoint: cased.GetEndpoint(cased.WorkflowsEndpoint),
	}
}

// Cancel a workflow approval request
func Cancel(id string) (*cased.Workflow, error) {
	return client().Cancel(id)
}

// Cancel a workflow approval request
func (c Client) Cancel(id string) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	path := fmt.Sprintf("/workflows/control-states/approvals/%s/cancel", id)
	err := c.Endpoint.Call(http.MethodPost, path, nil, workflow)
	return workflow, err
}
