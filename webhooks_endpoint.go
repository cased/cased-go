package cased

import "time"

type WebhooksEndpoint struct {
	// The Webhook Endpoint ID
	ID string `json:"id"`

	URL        string   `json:"url"`
	APIURL     string   `json:"api_url"`
	Secret     string   `json:"secret"`
	EventTypes []string `json:"event_types"`

	// UpdatedAt is when the workflow was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt is when the workflow was created.
	CreatedAt time.Time `json:"created_at"`
}

type WebhooksEndpointParams struct {
	Params `json:"-"`

	// Name is optional and only required if you intend to trigger workflows
	// by publishing events directly to them.
	URL        *string   `json:"url,omitempty"`
	EventTypes []*string `json:"event_types,omitempty"`
}
