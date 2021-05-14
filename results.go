package cased

type ResultState string

const (
	// Workflow
	ResultStatePending = "pending"

	// Workflow result state is unfulfilled when all controls have not been met.
	ResultStateUnfulfilled = "unfulfilled"

	// Workflow controls were not met. When the workflow result state is rejected
	// it's intended any further progress is canceled.
	ResultStateFulfilled = "fulfilled"

	// Workflow controls were not met. When the workflow result state is rejected
	// it's intended any further progress is canceled.
	ResultStateRejected = "rejected"
)

type Result struct {
	// The Result ID
	ID string `json:"id"`

	// The API URL for the result.
	ApiURL string `json:"api_url"`

	// State
	State ResultState `json:"state"`

	// Controls
	Controls Controls `json:"controls"`

	// Workflow
	Workflow Workflow `json:"workflow"`
}
