package main

import (
	"net/http"
	"time"

	"github.com/cased/cased-go"
)

func main() {
	p := cased.NewPublisher(
		// CASED_PUBLISH_URL=https://publish.cased.com
		cased.WithPublishURL("https://publish.cased.com"),

		// CASED_PUBLISH_KEY=publish_live_1mY8qb355NWIa3uY00H2fk7elpT
		cased.WithPublishKey("publish_live_1mY8qb355NWIa3uY00H2fk7elpT"),

		// CASED_DEBUG=1
		cased.WithDebug(true),

		// CASED_SILENCE=1
		cased.WithSilence(true),

		// You can configure your own client or re-use an existing HTTP client from
		// your application.
		cased.WithHTTPClient(&http.Client{}),

		// You can configure your own transport or re-use an existing HTTP transport
		// from your application.
		cased.WithHTTPTransport(&http.Transport{}),

		// CASED_HTTP_TIMEOUT=10s
		cased.WithHTTPTimeout(10*time.Second),
		cased.WithTransport(cased.NewNoopHTTPTransport()),
	)
	cased.SetPublisher(p)
	defer cased.Flush(10 * time.Second)

	_ = cased.Publish(cased.AuditEvent{
		"action":             "user.login",
		"actor":              "dewski",
		"actor_id":           "user_1dsGbftbx1c47iU8c7BzUcKJRcD",
		"organization":       "Cased",
		"organization_id":    "org_1dsGTnNZLzgwb1alwS2szN0KUo5",
		"request_id":         "27d62d1869e5f9826acc6cfd80edca90",
		"request_url":        "https://app.cased.com/saml/consume",
		"request_user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
	})
}
