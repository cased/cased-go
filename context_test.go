package cased

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContextFromContext(t *testing.T) {
	expected := AuditEvent{
		"action": "user.login",
	}
	ctx := context.WithValue(context.Background(), ContextKey, expected)
	actual := GetContextFromContext(ctx)

	assert.Equal(t, expected, actual)
}
