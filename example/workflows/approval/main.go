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
			Reason: cased.Bool(true),
			Approval: &cased.WorkflowControlsApprovalParams{
				Count:        cased.Int(1),
				SelfApproval: cased.Bool(false),
				Sources: &cased.WorkflowControlsApprovalSourcesParams{
					Slack: &cased.WorkflowControlsApprovalSourcesSlackParams{
						Channel: cased.String("#dewski-guard-test"),
					},
				},
			},
		},
	})

	if err != nil {
		// Handle errors that came from Cased
		if casedErr, ok := err.(*cased.Error); ok {
			switch casedErr.Code {
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
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

	resolveEvent(w, e)
}

func resolveEvent(w *cased.Workflow, e *cased.Event) {
	switch e.Result.State {
	case cased.WorkflowStateFulfilled:
		fmt.Println("Workflow successful!")
	case cased.WorkflowStateUnfulfilled:
		if e.Result.Controls.Authentication != nil {
			fmt.Println("fill this out")
		}

		// The reason control is unfulfilled
		reason := e.Result.Controls.Reason
		if reason != nil && reason.State == cased.WorkflowStateUnfulfilled {
			reasonPrompt(w, e)
			return
		}

		// The reason control is unfulfilled
		approval := e.Result.Controls.Approval
		if approval != nil {
			resolveApproval(w, e, approval)
		}
	case cased.WorkflowStateRejected:
		panic("Workflow rejected")
	}
}

func reasonPrompt(w *cased.Workflow, e *cased.Event) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Provide a reason: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	triggerWorkflow(w, strings.TrimSpace(input))
}

func resolveApproval(w *cased.Workflow, e *cased.Event, approval *cased.ResultControlsApproval) {
	switch approval.State {
	case cased.ResultControlsApprovalStatePending:
		// Waiting for other controls
	case cased.ResultControlsApprovalStateRequested:
		fmt.Println("Waiting for approval…")
		newEvent, err := event.Get(e.ID)
		if err != nil {
			panic(err)
		}

		resolveEvent(w, newEvent)
	case cased.ResultControlsApprovalStateApproved:
		newEvent, err := event.Get(e.ID)
		if err != nil {
			panic(err)
		}

		resolveEvent(w, newEvent)
	case cased.ResultControlsApprovalStateDenied,
		cased.ResultControlsApprovalStateTimedOut,
		cased.ResultControlsApprovalStateCanceled:
		panic("Denied!")
	default:
		panic(fmt.Sprintf("Unhandled approval state: %s", approval.State))
	}
}
