package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cased/cased-go"
	"github.com/cased/cased-go/event"
	"github.com/cased/cased-go/workflow"
)

const WorkflowName = "reason-required"

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

	existing, err := workflow.Get(WorkflowName)
	if err != nil {
		panic(err)
	}

	if existing.ID != "" {
		w = existing
	} else {
		w, err = workflow.New(&cased.WorkflowParams{
			Name: cased.String(WorkflowName),
			Controls: &cased.WorkflowControlsParams{
				Reason: cased.Bool(true),
			},
		})

		if err != nil {
			panic(err)
		}
	}

	return w
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

	switch e.Result.State {
	case cased.WorkflowStateFulfilled:
		fmt.Println("Workflow successful!")
	case cased.WorkflowStateUnfulfilled:
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Provide a reason: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		triggerWorkflow(w, input)
	case cased.WorkflowStateRejected:
		panic("Workflow rejected")
	}
}
