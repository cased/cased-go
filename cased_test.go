package cased

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCasedPublish(t *testing.T) {
	mp, restore := NewMockPublisher()
	defer restore()

	_ = Publish(AuditEvent{
		"action": "user.login",
	})

	assert.Equal(t, 1, len(mp.Events))
}

func TestCasedPublishWithSilencedPublisher(t *testing.T) {
	mp, restore := NewSilencedMockPublisher()
	defer restore()

	_ = Publish(AuditEvent{
		"action": "user.login",
	})

	assert.Same(t, mp, CurrentPublisher())
	assert.Equal(t, 0, len(mp.Events))
}

func TestCasedPublishWithContext(t *testing.T) {
	mp, restore := NewMockPublisher()
	defer restore()

	expected := AuditEvent{
		"action": "user.second",
	}
	ctx := context.WithValue(context.Background(), ContextKey, AuditEvent{
		"action": "user.first",
	})

	_ = PublishWithContext(ctx, expected)

	assert.Equal(t, 1, len(mp.Events))
	assert.Equal(t, expected, mp.Events[0])
}

func TestCasedPublishWithContextUsesPublishedAuditEvent(t *testing.T) {
	mp, restore := NewMockPublisher()
	defer restore()

	expected := AuditEvent{
		"user":     "cased",
		"user_id":  "1",
		"location": "1.1.1.1",
	}
	ctx := context.WithValue(context.Background(), ContextKey, AuditEvent{
		"user_id":  "2",
		"location": "1.1.1.1",
	})

	_ = PublishWithContext(ctx, AuditEvent{
		"user":    "cased",
		"user_id": "1",
	})

	assert.Equal(t, 1, len(mp.Events))
	assert.Equal(t, expected, mp.Events[0])
}
