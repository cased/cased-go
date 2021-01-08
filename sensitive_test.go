package cased

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSensitiveValueMarshalJSON(t *testing.T) {
	sv := SensitiveValue{Value: "Hello World"}
	data, err := json.Marshal(sv)

	assert.NoError(t, err)
	assert.Equal(t, `"Hello World"`, string(data))
}

func TestNewSensitiveValue(t *testing.T) {
	sv := NewSensitiveValue("Hello World", "sensitive-label")
	expectedRanges := []SensitiveRange{
		{
			Begin: 0,
			End:   11,
			Label: "sensitive-label",
		},
	}

	assert.Equal(t, "Hello World", sv.Value)
	assert.Equal(t, expectedRanges, sv.Ranges)
}
