package event

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

func New(params *cased.EventParams) (*cased.Event, error) {
	return client().New(params)
}

func (c Client) New(params *cased.EventParams) (*cased.Event, error) {
	event := &cased.Event{}
	path := "/workflows/events"
	if params.WorkflowID != nil {
		path = fmt.Sprintf("/workflows/%s/events", *params.WorkflowID)
	}

	err := c.Endpoint.Call(http.MethodPost, path, params, event)
	return event, err
}

func Get(id string) (*cased.Event, error) {
	return client().Get(id)
}

func (c Client) Get(id string) (*cased.Event, error) {
	event := &cased.Event{}
	err := c.Endpoint.Call(http.MethodGet, fmt.Sprintf("/workflows/events/%s", id), nil, event)
	return event, err
}
