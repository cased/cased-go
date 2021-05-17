package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cased/cased-go"
	"github.com/stretchr/testify/assert"
)

var mockedEvent = cased.Event{
	ID:     "event_1sY67fxGiiHp4f7dlcY26pL4eY5",
	APIURL: "https://api.cased.com/workflows/events/event_1sY67fxGiiHp4f7dlcY26pL4eY5",
	Result: cased.Result{
		ID:        "result_1sY67fxGiiHp4f7dlcY26pL4eY5",
		ApiURL:    "https://api.cased.com/workflows/events/event_1sY67fxGiiHp4f7dlcY26pL4eY5/result",
		State:     cased.WorkflowStateFulfilled,
		Controls:  cased.ResultControls{},
		Workflow:  nil,
		UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	Event: cased.EventPayload{
		"testing": true,
	},
	OriginalEvent: cased.EventPayload{
		"testing": true,
	},
	UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
}

var mockedEndpoint = &MockEndpoint{
	CallFunc: func(method, path string, params cased.ParamsContainer, i interface{}) error {
		data, err := json.Marshal(mockedEvent)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, &i)
	},
}

type MockEndpoint struct {
	CallFunc func(method, path string, params cased.ParamsContainer, i interface{}) error
}

func (me *MockEndpoint) Call(method, path string, params cased.ParamsContainer, i interface{}) error {
	return me.CallFunc(method, path, params, i)
}

func TestEvent_New(t *testing.T) {
	ep := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(ep)

	e, err := New(&cased.EventParams{
		Event: cased.EventPayload{
			"testing": true,
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "event_1sY67fxGiiHp4f7dlcY26pL4eY5", e.ID)
	assert.Equal(t, "https://api.cased.com/workflows/events/event_1sY67fxGiiHp4f7dlcY26pL4eY5", e.APIURL)
	assert.Equal(t, "2021-01-01T00:00:00Z", e.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", e.CreatedAt.Format(time.RFC3339Nano))
}

func TestEvent_Get(t *testing.T) {
	ep := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(ep)

	e, err := Get("event_1sY67fxGiiHp4f7dlcY26pL4eY5")

	assert.NoError(t, err)
	assert.Equal(t, "event_1sY67fxGiiHp4f7dlcY26pL4eY5", e.ID)
	assert.Equal(t, "https://api.cased.com/workflows/events/event_1sY67fxGiiHp4f7dlcY26pL4eY5", e.APIURL)
	assert.Equal(t, "2021-01-01T00:00:00Z", e.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", e.CreatedAt.Format(time.RFC3339Nano))
}
