package endpoint

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cased/cased-go"
	"github.com/stretchr/testify/assert"
)

var mockedWorkflowsEndpoint = cased.WebhooksEndpoint{
	ID:        "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5",
	URL:       "https://app.cased.com",
	APIURL:    "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5",
	UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
}

var mockedEndpoint = &MockEndpoint{
	CallFunc: func(method, path string, params cased.ParamsContainer, i interface{}) error {
		data, err := json.Marshal(mockedWorkflowsEndpoint)
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

func TestWorkflow_New(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	we, err := New(&cased.WebhooksEndpointParams{
		URL: cased.String("workflow"),
	})

	assert.NoError(t, err)
	assert.Equal(t, "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.ID)
	assert.Equal(t, "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.APIURL)
	assert.Equal(t, "https://app.cased.com", we.URL)
	assert.Equal(t, "2021-01-01T00:00:00Z", we.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", we.CreatedAt.Format(time.RFC3339Nano))
}

func TestWorkflow_Get(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	we, err := Get("workflow_id")

	assert.NoError(t, err)
	assert.Equal(t, "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.ID)
	assert.Equal(t, "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.APIURL)
	assert.Equal(t, "https://app.cased.com", we.URL)
	assert.Equal(t, "2021-01-01T00:00:00Z", we.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", we.CreatedAt.Format(time.RFC3339Nano))
}

func TestWorkflow_Update(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	we, err := Update("workflow_id", &cased.WebhooksEndpointParams{})

	assert.NoError(t, err)
	assert.Equal(t, "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.ID)
	assert.Equal(t, "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.APIURL)
	assert.Equal(t, "https://app.cased.com", we.URL)
	assert.Equal(t, "2021-01-01T00:00:00Z", we.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", we.CreatedAt.Format(time.RFC3339Nano))
}

func TestWorkflow_Delete(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	_, err := Delete("workflow_id")
	assert.NoError(t, err)
}
