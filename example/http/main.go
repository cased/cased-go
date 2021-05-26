package main

import (
	"log"
	"net/http"
	"time"

	"github.com/cased/cased-go"
	casedhttp "github.com/cased/cased-go/http"
)

func main() {
	mux := http.NewServeMux()
	loginHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer cased.PublishWithContext(req.Context(), cased.AuditEvent{ // nolint:errcheck
			"action": "user.login",
		})

		w.Write([]byte("Logged in!")) // nolint:errcheck
	})
	mux.Handle("/login", casedhttp.ContextMiddleware(loginHandler))

	webhookHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Webhook!")) // nolint:errcheck
	})
	webhookVerificationParams := &casedhttp.VerifyWebhookSignatureParams{
		TimestampExpires: 5 * time.Minute,
	}
	mux.Handle("/webhook", casedhttp.VerifyWebhookSignatureMiddleware(webhookHandler, webhookVerificationParams))

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
