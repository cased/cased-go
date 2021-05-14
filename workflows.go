package cased

import "time"

type Workflow struct {
	// The Workflow ID
	ID string `json:"id"`

	// Name of the workflow to be used to trigger the workflow.
	Name *string `json:"name,omitempty"`

	// The API URL for the workflow.
	ApiURL string `json:"api_url"`

	// Conditions are how Cased determines which workflow should run when an event
	// is published and a workflow is not specified.
	Conditions []Condition `json:"conditions"`

	// Controls specifies the controls enabled for this workflow.
	Controls Controls `json:"controls"`

	// UpdatedAt is when the workflow was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt is when the workflow was created.
	CreatedAt time.Time `json:"created_at"`
}

// WorkflowParams contains the available fields when creating and updating
// workflows.
type WorkflowParams struct {
	Params `json:"-"`

	// Name is optional and only required if you intend to trigger workflows
	// by publishing events directly to them.
	Name *string `json:"name"`

	// Conditions specify the conditions the workflow should match when events
	// are not published directly to a workflow.
	Conditions []Condition `json:"conditions,omitempty"`

	// Configure the controls necessary for the workflow to reach the fulfilled
	// state.
	Controls Controls `json:"controls"`
}

// ConditionOperator contains all condition operators available for workflows.
type ConditionOperator string

const (
	// ConditionOperatorEndsWith case-insensitive  matches "world" in
	// "hello world"
	ConditionOperatorEndsWith = "endsWith"

	// ConditionOperatorEqual case-insensitive matches "cased" both "cased"
	// or "Cased"
	ConditionOperatorEqual = "eq"

	// ConditionOperatorIncludes case-insensitive matches when the value is
	// included in the specified field for both strings and arrays.
	ConditionOperatorIncludes = "in"

	// ConditionOperatorNotEqual case-insensitive matches when the field does not
	// contain the specified value.
	ConditionOperatorNotEqual = "not"

	// ConditionOperatorRegex matches based on the provided regular expression.
	// Not currently enabled.
	ConditionOperatorRegex = "re"

	// ConditionOperatorStartsWith case-insensitive matches "hello" in
	// "hello world"
	ConditionOperatorStartsWith = "startsWith"
)

// Condition is an individual clause in one or more conditions that can be used
// to match incoming events.
//
// All conditions are evaluated ignoring the case of the value.
type Condition struct {
	// The path to the field on the event to evaluate this condition for.
	Field string `json:"field"`

	// Operator specifies the operator use to evaluate the condition. See
	// `ConditionOperator` for all available operators.
	Operator ConditionOperator `json:"operator"`

	// Value contains the value to be used to evaluate the condition based on its
	// configured operator.
	Value string `json:"value"`
}

type Controls struct {
	// Require a user to provide a reason to continue the workflow.
	Reason bool `json:"reason"`

	// Require a user to authenticate with Cased to continue the workflow.
	Authentication bool `json:"authentication"`

	// Require a user to receive approval before a workflow is fulfilled or
	// rejected.
	Approval *ApprovalControl `json:"approval,omitempty"`
}

type ApprovalControl struct {
	// The number of approvals required to fulfill the approval requirement.
	//
	// Approval count cannot exceed the number of users on your account,
	// otherwise an error will be returned.
	Count int `json:"count"`

	// Permit an approval request to allow user requesting approval the ability
	// to approve their own request. If the Authentication control is disabled,
	// any user can approve the request and this setting is ignored.
	SelfApproval bool `json:"self_approval"`

	// Determine how long the approval lasts for.
	Duration int `json:"duration"`

	// Control how long the approval request is valid for. If not supplied,
	// approval requests can be responded to indefinitely.
	Timeout *int `json:"timeout"`

	// List of responders that can include individual users and groups of users
	// who are authorized to respond to the approval request.
	Responders *ApprovalControlResponders `json:"responders,omitempty"`

	// Sources where to obtain the approval from. If not provided, defaults to
	// email.
	Sources *ApprovalControlSources `json:"sources,omitempty"`
}

// ApprovalControlResponders is the list of individual users and groups of users
// who are authorized to respond to an approval request.
type ApprovalControlResponders map[string]string

// ApprovalControlSources determines where approval requests are delivered.
type ApprovalControlSources struct {
	// Email determines if an email is delivered for the approval request.
	Email bool `json:"email,omitempty"`

	// Slack when provided, publishes a Slack message which users can respond to
	// the request.
	Slack *ApprovalControlSourceSlack `json:"slack,omitempty"`
}

// ApprovalControlSourceSlack configures which the Slack approval source.
type ApprovalControlSourceSlack struct {
	// The fully qualified Slack channel name (ie: #security).
	Channel string `json:"channel"`
}
