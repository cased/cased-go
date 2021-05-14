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

func New(params *cased.WorkflowParams) (*cased.Workflow, error) {
	return client().New(params)
}

func (c Client) New(params *cased.WorkflowParams) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodPost, "/workflows", params, workflow)
	return workflow, err
}

func Get(id string) (*cased.Workflow, error) {
	return client().Get(id)
}

func (c Client) Get(id string) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodGet, fmt.Sprintf("/workflows/%s", id), nil, workflow)
	return workflow, err
}

func Update(id string, params *cased.WorkflowParams) (*cased.Workflow, error) {
	return client().Update(id, params)
}

func (c Client) Update(id string, params *cased.WorkflowParams) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodPatch, fmt.Sprintf("/workflows/%s", id), params, workflow)
	return workflow, err
}

func Delete(id string) (*cased.Workflow, error) {
	return client().Delete(id)
}

func (c Client) Delete(id string) (*cased.Workflow, error) {
	workflow := &cased.Workflow{}
	err := c.Endpoint.Call(http.MethodDelete, fmt.Sprintf("/workflows/%s", id), nil, workflow)
	return workflow, err
}

func List(id string) ([]*cased.Workflow, error) {
	return client().List(id)
}

func (c Client) List(id string) ([]*cased.Workflow, error) {
	workflows := []*cased.Workflow{}
	err := c.Endpoint.Call(http.MethodGet, "/workflows", nil, workflows)
	return workflows, err
}
