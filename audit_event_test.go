package cased

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditEventMarshalJSON(t *testing.T) {
	ae := AuditEvent{
		"action": "user.login",
	}

	data, err := json.Marshal(ae)
	assert.NoError(t, err)

	aet := AuditEvent{}
	err = json.Unmarshal(data, &aet)
	assert.NoError(t, err)

	assert.Equal(t, "user.login", aet["action"])
}

func TestAuditEventPayloadMarshalJSONWithSensitiveValue(t *testing.T) {
	ae := AuditEvent{
		"action": "user.login",
		"user":   NewSensitiveValue("John Doe", "name"),
	}
	aep := NewAuditEventPayload(ae)

	data, err := json.Marshal(aep)
	assert.NoError(t, err)

	actual := AuditEventPayload{}
	err = json.Unmarshal(data, &actual)
	assert.NoError(t, err)

	expected := []*SensitiveRange{
		{
			Begin: 0,
			End:   8,
			Label: "name",
		},
	}

	assert.Equal(t, "user.login", actual.AuditEvent["action"])
	assert.Equal(t, "John Doe", actual.AuditEvent["user"])
	assert.Equal(t, expected, actual.DotCased.PII[".user"], string(data))
}
