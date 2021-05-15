package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cased/cased-go"
	"github.com/cased/cased-go/event"
	"github.com/cased/cased-go/workflow"
)

func main() {
	var a *cased.Workflow

	existing, err := workflow.Get("named")
	if err != nil {
		panic(err)
	}

	if existing.ID != "" {
		fmt.Println("Using existing")
		a = existing

		o("Existing workflow", a)
	} else {
		a, err = workflow.New(&cased.WorkflowParams{
			Conditions: []*cased.WorkflowConditionParams{
				{
					Field:    cased.String("hello"),
					Operator: cased.String(string(cased.WorkflowConditionOperatorEqual)),
					Value:    cased.String("world"),
				},
			},
			Controls: &cased.WorkflowControlsParams{
				Authentication: cased.Bool(true),
			},
		})

		if err != nil {
			panic(err)
		}

		o("Create workflow", a)
	}

	b, err := workflow.Get(a.ID)
	if err != nil {
		panic(err)
	}
	o("Get workflow", b)

	c, err := workflow.Update(a.ID, &cased.WorkflowParams{
		Name:       cased.String("named"),
		Conditions: []*cased.WorkflowConditionParams{},
	})
	if err != nil {
		panic(err)
	}
	o("Update workflow", c)

	e, err := event.New(&cased.EventParams{
		Event: cased.EventPayload{
			"hello": "world",
		},
		WorkflowID: cased.String(a.ID),
	})
	if err != nil {
		panic(err)
	}
	o("New event with workflow", e)

	t, err := event.New(&cased.EventParams{
		Event: cased.EventPayload{
			"hello": "world",
		},
	})
	if err != nil {
		panic(err)
	}
	o("New event without workflow", t)

	d, err := workflow.Delete(a.ID)
	if err != nil {
		panic(err)
	}
	o("Delete workflow", d)
}

func o(title string, i interface{}) {
	empJSON, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("===============================================")
	fmt.Println(title)
	fmt.Println("===============================================")
	fmt.Println(string(empJSON))
}
