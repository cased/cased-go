package cased

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventParams_MarshalJSON(t *testing.T) {
	ep := EventParams{
		WorkflowID: String("workflow_id"),
		Event: EventPayload{
			"user": "user_id",
		},
	}

	data, err := json.Marshal(ep)
	assert.NoError(t, err)
	assert.Equal(t, `{"user":"user_id"}`, string(data))
}
