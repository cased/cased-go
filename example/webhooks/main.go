package main

import (
	"fmt"

	"github.com/cased/cased-go"
	"github.com/cased/cased-go/webhooks/endpoint"
)

func main() {
	we, err := endpoint.New(&cased.WebhooksEndpointParams{
		URL: cased.String("https://cased.com"),
		EventTypes: []*string{
			cased.String("event.created"),
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", we)

	_, err = endpoint.Delete(we.ID)
	if err != nil {
		panic(err)
	}
}
