package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/cased/cased-go"
	"github.com/cased/cased-go/event"
	"github.com/cased/cased-go/workflow"
	"github.com/cased/cased-go/workflow/controlstates/approval"
)

const WorkflowName = "provision-database"

type task struct {
	workflow                  *cased.Workflow
	event                     *cased.Event
	c                         chan os.Signal
	promptedForApproval       bool
	promptedForAuthentication bool
	once                      sync.Once
}

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
				SelfApproval: cased.Bool(true),
				Sources: &cased.WorkflowControlsApprovalSourcesParams{
					Slack: &cased.WorkflowControlsApprovalSourcesSlackParams{
						Channel: cased.String("#cased"),
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

	task := &task{
		workflow:                  w,
		event:                     e,
		c:                         make(chan os.Signal, 1),
		promptedForApproval:       false,
		promptedForAuthentication: false,
	}

	task.resolve()
}

func (t *task) refresh() {
	newEvent, err := event.Get(t.event.ID)
	if err != nil {
		panic(err)
	}

	t.event = newEvent
	t.resolve()
}

func (t *task) resolve() {
	switch t.event.Result.State {
	case cased.WorkflowStateFulfilled:
		fmt.Println("Workflow successful!")
	case cased.WorkflowStateUnfulfilled:
		// The reason control is unfulfilled
		reason := t.event.Result.Controls.Reason
		if reason != nil && reason.State == cased.WorkflowStateUnfulfilled {
			t.reasonPrompt()
			return
		}

		authentication := t.event.Result.Controls.Authentication
		if authentication != nil && authentication.State == cased.WorkflowStateUnfulfilled {
			t.authenticationPrompt(authentication)
			return
		}

		// The reason control is unfulfilled
		approval := t.event.Result.Controls.Approval
		if approval != nil {
			t.resolveApproval(approval)
			return
		}
	case cased.WorkflowStateRejected:
		panic("Workflow rejected")
	}
}

// Prompt the user to authenticate with Cased per the configured workflow
// controls.
func (t *task) authenticationPrompt(authentication *cased.ResultControlsAuthentication) {
	if !t.promptedForAuthentication {
		fmt.Printf("To login, please visit:\n%s\n", authentication.URL)
		t.promptedForAuthentication = true
	}

	t.refresh()
}

// Prompt the user for a reason and trigger a new workflow with the provided
// reason
func (t *task) reasonPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Provide a reason: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	triggerWorkflow(t.workflow, strings.TrimSpace(input))
}

func (t *task) notifyInterrupt(approvalID string) func() {
	return func() {
		signal.Notify(t.c, os.Interrupt)

		go func() {
			select {
			case <-t.c:
				approval.Cancel(approvalID)
			}
		}()
	}
}

// Handle each approval state
func (t *task) resolveApproval(a *cased.ResultControlsApproval) {
	t.once.Do(t.notifyInterrupt(a.ID))

	switch a.State {
	case cased.ResultControlsApprovalStatePending:
		t.refresh()
	case cased.ResultControlsApprovalStateRequested:
		if !t.promptedForApproval {
			fmt.Println("Waiting for approval…")
			t.promptedForApproval = true
		}

		t.refresh()
	case cased.ResultControlsApprovalStateApproved:
		newEvent, err := event.Get(t.event.ID)
		if err != nil {
			panic(err)
		}

		t.event = newEvent
		t.refresh()
	case cased.ResultControlsApprovalStateDenied,
		cased.ResultControlsApprovalStateTimedOut,
		cased.ResultControlsApprovalStateCanceled:
		fmt.Println(a.State)
		os.Exit(1)
	default:
		panic(fmt.Sprintf("Unhandled approval state: %s", a.State))
	}
}
