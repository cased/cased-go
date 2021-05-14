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
	a, err := workflow.New(&cased.WorkflowParams{
		Conditions: []cased.Condition{
			{
				Field:    "hello",
				Operator: cased.ConditionOperatorEqual,
				Value:    "world",
			},
		},
	})

	if err != nil {
		panic(err)
	}
	o(a)
	fmt.Printf("%+v\n", a)

	b, err := workflow.Get(a.ID)
	if err != nil {
		panic(err)
	}
	o(b)

	c, err := workflow.Update(a.ID, &cased.WorkflowParams{
		Name: cased.String("named"),
	})
	if err != nil {
		panic(err)
	}
	o(c)

	e, err := event.New(&cased.EventParams{
		Event: cased.EventPayload{
			"hello": "world",
		},
		WorkflowID: cased.String(a.ID),
	})
	if err != nil {
		panic(err)
	}
	o(e)

	t, err := event.New(&cased.EventParams{
		Event: cased.EventPayload{
			"hello": "world",
		},
	})
	if err != nil {
		panic(err)
	}
	o(t)

	d, err := workflow.Delete(a.ID)
	if err != nil {
		panic(err)
	}
	o(d)
}

func o(i interface{}) {
	empJSON, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(empJSON))
}
