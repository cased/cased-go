package webhook

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

// Create a new webhook endpoint
func New(params *cased.WebhookEndpointParams) (*cased.WebhookEndpoint, error) {
	return client().New(params)
}

// Create a new webhook endpoint
func (c Client) New(params *cased.WebhookEndpointParams) (*cased.WebhookEndpoint, error) {
	endpoint := &cased.WebhookEndpoint{}
	err := c.Endpoint.Call(http.MethodPost, "/webhooks/endpoints", params, endpoint)
	return endpoint, err
}

// Retrieve a webhook endpoint by it's ID or name.
func Get(id string) (*cased.WebhookEndpoint, error) {
	return client().Get(id)
}

// Retrieve a webhook endpoint by it's ID or name.
func (c Client) Get(id string) (*cased.WebhookEndpoint, error) {
	endpoint := &cased.WebhookEndpoint{}
	err := c.Endpoint.Call(http.MethodGet, fmt.Sprintf("/webhooks/endpoints/%s", id), nil, endpoint)
	return endpoint, err
}

// Update a webhook endpoint by it's ID or name.
func Update(id string, params *cased.WebhookEndpointParams) (*cased.WebhookEndpoint, error) {
	return client().Update(id, params)
}

// Update a webhook endpoint by it's ID or name.
func (c Client) Update(id string, params *cased.WebhookEndpointParams) (*cased.WebhookEndpoint, error) {
	endpoint := &cased.WebhookEndpoint{}
	err := c.Endpoint.Call(http.MethodPatch, fmt.Sprintf("/webhooks/endpoints/%s", id), params, endpoint)
	return endpoint, err
}

// Delete a webhook endpoint by it's ID or name.
func Delete(id string) (*cased.WebhookEndpoint, error) {
	return client().Delete(id)
}

// Delete a webhook endpoint by it's ID or name.
func (c Client) Delete(id string) (*cased.WebhookEndpoint, error) {
	endpoint := &cased.WebhookEndpoint{}
	err := c.Endpoint.Call(http.MethodDelete, fmt.Sprintf("/webhooks/endpoints/%s", id), nil, endpoint)
	return endpoint, err
}
