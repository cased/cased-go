package cased

import "time"

type WebhooksEndpoint struct {
	// The Webhook Endpoint ID
	ID string `json:"id"`

	// The API URL for the webhook endpoint.
	APIURL string `json:"api_url"`

	// Secret used to sign payloads.
	Secret string `json:"secret"`

	// EventTypes to deliver to the webhook endpoint. If none are specified, all
	// event types will deliver events.
	EventTypes []string `json:"event_types"`

	// UpdatedAt is when the workflow was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt is when the workflow was created.
	CreatedAt time.Time `json:"created_at"`
}

type WebhooksEndpointParams struct {
	Params `json:"-"`

	// URL to deliver webhook events to.
	URL *string `json:"url,omitempty"`

	// EventTypes to deliver to the webhook endpoint. If none are specified, all
	// event types will deliver events.
	EventTypes []*string `json:"event_types,omitempty"`
}
