package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cased/cased-go"
	"github.com/cased/cased-go/event"
	"github.com/cased/cased-go/workflow"
)

const WorkflowName = "provision-database"

var (
	promptedForAuthentication = false
	promptedForApproval       = false
)

func main() {
	// Fetch the reason required workflow
	w := fetchWorkflow()

	// If the user supplied a reason when invoking the program, pass it to the event
	var reason string
	if len(os.Args) > 1 {
		reason = os.Args[1]
	}

	triggerWorkflow(w, reason)
}

func fetchWorkflow() *cased.Workflow {
	var w *cased.Workflow
	var err error

	w, err = workflow.Get(WorkflowName)
	if err != nil {
		// Handle errors that came from Cased
		if casedErr, ok := err.(*cased.Error); ok {
			switch casedErr.Code {
			case cased.ErrorCodeNotFound:
				w, err = createWorkflow()
				if err != nil {
					panic(err)
				}

				return w
			default:
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return w
}

func createWorkflow() (*cased.Workflow, error) {
	w, err := workflow.New(&cased.WorkflowParams{
		Name: cased.String(WorkflowName),
		Controls: &cased.WorkflowControlsParams{
			Authentication: cased.Bool(true),
			Reason:         cased.Bool(true),
			Approval: &cased.WorkflowControlsApprovalParams{
				Count:        cased.Int(1),
				SelfApproval: cased.Bool(true),
				Sources: &cased.WorkflowControlsApprovalSourcesParams{
					Slack: &cased.WorkflowControlsApprovalSourcesSlackParams{
						Channel: cased.String("#dewski-guard-test"),
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return w, nil
}

func triggerWorkflow(w *cased.Workflow, reason string) {
	e, err := event.New(&cased.EventParams{
		WorkflowID: cased.String(w.ID),
		Event: cased.EventPayload{
			"reason": reason,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created event %s\n", e.ID)

	resolveEvent(w, e)
}

func refreshEvent(w *cased.Workflow, e *cased.Event) {
	newEvent, err := event.Get(e.ID)
	if err != nil {
		panic(err)
	}

	resolveEvent(w, newEvent)
}

func resolveEvent(w *cased.Workflow, e *cased.Event) {
	switch e.Result.State {
	case cased.WorkflowStateFulfilled:
		fmt.Println("Workflow successful!")
	case cased.WorkflowStateUnfulfilled:
		// The reason control is unfulfilled
		reason := e.Result.Controls.Reason
		if reason != nil && reason.State == cased.WorkflowStateUnfulfilled {
			reasonPrompt(w, e)
			return
		}

		authentication := e.Result.Controls.Authentication
		if authentication != nil && authentication.State == cased.WorkflowStateUnfulfilled {
			authenticationPrompt(w, e, authentication)
			return
		}

		// The reason control is unfulfilled
		approval := e.Result.Controls.Approval
		if approval != nil {
			resolveApproval(w, e, approval)
			return
		}
	case cased.WorkflowStateRejected:
		panic("Workflow rejected")
	}
}

// Prompt the user to authenticate with Cased per the configured workflow
// controls.
func authenticationPrompt(w *cased.Workflow, e *cased.Event, authentication *cased.ResultControlsAuthentication) {
	if !promptedForAuthentication {
		fmt.Printf("To login, please visit:\n%s\n", authentication.URL)
		promptedForAuthentication = true
	}

	refreshEvent(w, e)
}

// Prompt the user for a reason and trigger a new workflow with the provided
// reason
func reasonPrompt(w *cased.Workflow, e *cased.Event) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Provide a reason: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	triggerWorkflow(w, strings.TrimSpace(input))
}

// Handle each approval state
func resolveApproval(w *cased.Workflow, e *cased.Event, approval *cased.ResultControlsApproval) {
	switch approval.State {
	case cased.ResultControlsApprovalStatePending:
		refreshEvent(w, e)
	case cased.ResultControlsApprovalStateRequested:
		if !promptedForApproval {
			fmt.Println("Waiting for approval…")
			promptedForApproval = true
		}

		refreshEvent(w, e)
	case cased.ResultControlsApprovalStateApproved:
		newEvent, err := event.Get(e.ID)
		if err != nil {
			panic(err)
		}

		resolveEvent(w, newEvent)
	case cased.ResultControlsApprovalStateDenied,
		cased.ResultControlsApprovalStateTimedOut,
		cased.ResultControlsApprovalStateCanceled:
		fmt.Println(approval.State)
		os.Exit(1)
	default:
		panic(fmt.Sprintf("Unhandled approval state: %s", approval.State))
	}
}
