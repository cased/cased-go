package workflow

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

// Create a new workflow
func New(params *cased.WorkflowParams) (*cased.Workflow, error) {
	return client().New(params)
}

// Create a new workflow
func (c Client) New(params *cased.WorkflowParams) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodPost, "/workflows", params, workflow)
	return workflow, err
}

// Retrieve a workflow by it's ID or name.
func Get(id string) (*cased.Workflow, error) {
	return client().Get(id)
}

// Retrieve a workflow by it's ID or name.
func (c Client) Get(id string) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodGet, fmt.Sprintf("/workflows/%s", id), nil, workflow)
	return workflow, err
}

// Update a workflow by it's ID or name.
func Update(id string, params *cased.WorkflowParams) (*cased.Workflow, error) {
	return client().Update(id, params)
}

// Update a workflow by it's ID or name.
func (c Client) Update(id string, params *cased.WorkflowParams) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodPatch, fmt.Sprintf("/workflows/%s", id), params, workflow)
	return workflow, err
}

// Delete a workflow by it's ID or name.
func Delete(id string) (*cased.Workflow, error) {
	return client().Delete(id)
}

// Delete a workflow by it's ID or name.
func (c Client) Delete(id string) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodDelete, fmt.Sprintf("/workflows/%s", id), nil, workflow)
	return workflow, err
}
