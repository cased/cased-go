package main

import (
	"time"

	"github.com/cased/cased-go"
)

func main() {
	p := cased.NewPublisher(
		cased.WithPublishKey("publish_live_1mY8qb355NWIa3uY00H2fk7elpT"),
	)
	cased.SetPublisher(p)
	defer cased.Flush(10 * time.Second)

	_ = cased.Publish(cased.AuditEvent{
		"action":             "user.login",
		"actor":              cased.NewSensitiveValue("dewski", "username"),
		"actor_id":           "user_1dsGbftbx1c47iU8c7BzUcKJRcD",
		"location":           cased.NewSensitiveValue("127.0.0.1", "ip-address"),
		"organization":       "Cased",
		"organization_id":    "org_1dsJHNqtyghCVpkdVy7dWi5KUxv",
		"request_id":         "27d62d1869e5f9826acc6cfd80edca90",
		"request_url":        "https://app.cased.com/saml/consume",
		"request_user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
	})
}
