package cased

import "time"

type EventPayload map[string]interface{}

type Event struct {
	// The Event ID
	ID string `json:"id"`

	// The API URL for the workflow.
	ApiURL string `json:"api_url"`

	Result Result `json:"result"`

	Event         EventPayload `json:"event"`
	OriginalEvent EventPayload `json:"original_event"`

	// UpdatedAt is when the workflow was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt is when the workflow was created.
	CreatedAt time.Time `json:"created_at"`
}

type EventParams struct {
	Params     `json:"-"`
	WorkflowID *string `json:"-"`
	Event      EventPayload
}
