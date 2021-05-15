package cased

type Result struct {
	// The Result ID
	ID string `json:"id"`

	// The API URL for the result.
	ApiURL string `json:"api_url"`

	// State
	State WorkflowState `json:"state"`

	// Controls
	Controls ResultControls `json:"controls"`

	// Workflow
	Workflow Workflow `json:"workflow"`
}

type ResultControls struct {
	Authentication *ResultControlsAuthentication `json:"authentication,omitempty"`
	Reason         *ResultControlsReason         `json:"reason,omitempty"`
	Approval       *ResultControlsApproval       `json:"approval,omitempty"`
}

type ResultControlsAuthentication struct {
	State  WorkflowState                     `json:"state"`
	User   *ResultControlsAuthenticationUser `json:"user"`
	URL    string                            `json:"url"`
	ApiURL string                            `json:"api_url"`
}

type ResultControlsAuthenticationUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type ResultControlsReason struct {
	State WorkflowState `json:"state"`
}

type ResultControlsApprovalState string

const (
	ResultControlsApprovalStatePending   ResultControlsApprovalState = "pending"
	ResultControlsApprovalStateRequested ResultControlsApprovalState = "requested"
	ResultControlsApprovalStateApproved  ResultControlsApprovalState = "approved"
	ResultControlsApprovalStateDenied    ResultControlsApprovalState = "denied"
	ResultControlsApprovalStateTimedOut  ResultControlsApprovalState = "timed_out"
	ResultControlsApprovalStateCanceled  ResultControlsApprovalState = "canceled"
)

type ResultControlsApproval struct {
	State    ResultControlsApprovalState     `json:"state"`
	Requests []ResultControlsApprovalRequest `json:"requests"`
	Source   ResultControlsApprovalSource    `json:"source"`
}

type ResultControlsApprovalRequestType string

type ResultControlsApprovalRequest struct {
	ID    string                            `json:"id"`
	State WorkflowState                     `json:"state"`
	Type  ResultControlsApprovalRequestType `json:"type"`
}

type ResultControlsApprovalSource struct {
	Email bool                              `json:"email"`
	Slack ResultControlsApprovalSourceSlack `json:"slack"`
}

type ResultControlsApprovalSourceSlack struct {
	Channel string `json:"channel"`
}
