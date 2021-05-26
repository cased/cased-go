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

	// Setup webhook for Cased
	webhookHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Webhook!")) // nolint:errcheck
	})
	webhookVerificationParams := &casedhttp.VerifyWebhookSignatureParams{
		// Secret can be omitted if the CASED_WEBHOOK_SECRET environment variable is set.
		Secret: cased.String("webhooks_secret_1t2WblwTxlBsYhnHI4Sm6DF2yt2"),
		TimestampExpires: 5 * time.Minute,
	}
	mux.Handle("/webhook", casedhttp.VerifyWebhookSignatureMiddleware(webhookHandler, webhookVerificationParams))

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
