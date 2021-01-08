package cased

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSensitiveDataProcessor(t *testing.T) {
	ae := AuditEvent{
		"action": "user.login",
		"user":   NewSensitiveValue("John Doe", "name"),
	}
	aep := NewAuditEventPayload(ae)

	expected := []*SensitiveRange{
		{
			Begin: 0,
			End:   8,
			Label: "name",
		},
	}

	assert.Equal(t, expected, aep.DotCased.PII[".user"])
}

func TestSensitiveDataProcessorWithNestedData(t *testing.T) {
	ae := AuditEvent{
		"action": "user.login",
		"location": map[string]SensitiveValue{
			"city": NewSensitiveValue("San Francisco", "city"),
		},
	}
	aep := NewAuditEventPayload(ae)

	expected := []*SensitiveRange{
		{
			Begin: 0,
			End:   13,
			Label: "city",
		},
	}

	assert.Equal(t, expected, aep.DotCased.PII[".location.city"])
}

func TestPublishedAtProcessorAddsPublishedAt(t *testing.T) {
	aep := NewAuditEventPayload(AuditEvent{})

	assert.False(t, aep.DotCased.PublishedAt.IsZero())
}
