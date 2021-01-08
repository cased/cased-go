package main

import (
	"log"
	"net/http"

	"github.com/cased/cased-go"
	casedhttp "github.com/cased/cased-go/http"
)

func main() {
	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer cased.PublishWithContext(req.Context(), cased.AuditEvent{ // nolint:errcheck
			"action": "user.login",
		})

		w.Write([]byte("Logged in!")) // nolint:errcheck
	})
	mux.Handle("/login", casedhttp.ContextMiddleware(finalHandler))

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
