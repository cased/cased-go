package cased

import (
	"context"
	"time"
)

// Publisher describes the interface for structs that want to publish audit
// events to Cased.
type Publisher interface {
	Publish(event AuditEvent) error
	Options() PublisherOptions
	Flush(timeout time.Duration) bool
}

// Publish publishes an audit event to Cased.
func Publish(event AuditEvent) error {
	client := CurrentPublisher()
	if client.Options().Silence {
		Logger.Println("Audit event was silenced.")
		return nil
	}

	return client.Publish(event)
}

// PublishWithContext enriches the provided audit event with the context set in
// the request. If the same key is present in both the context and provided
// audit event, the audit event value will be preserved.
func PublishWithContext(ctx context.Context, event AuditEvent) error {
	c := GetContextFromContext(ctx)
	for key, value := range c {
		if _, ok := event[key]; ok {
			continue
		}

		event[key] = value
	}

	return Publish(event)
}

// Flush waits for audit events to be published.
func Flush(timeout time.Duration) bool {
	return CurrentPublisher().Flush(timeout)
}
