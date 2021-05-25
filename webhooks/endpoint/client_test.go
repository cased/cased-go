package endpoint

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cased/cased-go"
	"github.com/stretchr/testify/assert"
)

var mockedWebhooksEndpoint = cased.WebhooksEndpoint{
	ID:        "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5",
	URL:       "https://app.cased.com",
	APIURL:    "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5",
	UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
}

var mockedEndpoint = &MockEndpoint{
	CallFunc: func(method, path string, params cased.ParamsContainer, i interface{}) error {
		data, err := json.Marshal(mockedWebhooksEndpoint)
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

func TestWebhooksEndpoint_New(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	we, err := New(&cased.WebhooksEndpointParams{
		URL: cased.String("https://app.cased.com"),
	})

	assert.NoError(t, err)
	assert.Equal(t, "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.ID)
	assert.Equal(t, "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.APIURL)
	assert.Equal(t, "https://app.cased.com", we.URL)
	assert.Equal(t, "2021-01-01T00:00:00Z", we.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", we.CreatedAt.Format(time.RFC3339Nano))
}

func TestWebhooksEndpoint_Get(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	we, err := Get("webhooks_endpoint_id")

	assert.NoError(t, err)
	assert.Equal(t, "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.ID)
	assert.Equal(t, "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.APIURL)
	assert.Equal(t, "https://app.cased.com", we.URL)
	assert.Equal(t, "2021-01-01T00:00:00Z", we.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", we.CreatedAt.Format(time.RFC3339Nano))
}

func TestWebhooksEndpoint_Update(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	we, err := Update("webhooks_endpoint_id", &cased.WebhooksEndpointParams{})

	assert.NoError(t, err)
	assert.Equal(t, "webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.ID)
	assert.Equal(t, "https://api.cased.com/webhooks/endpoints/webhooks_endpoint_1sY67fxGiiHp4f7dlcY26pL4eY5", we.APIURL)
	assert.Equal(t, "https://app.cased.com", we.URL)
	assert.Equal(t, "2021-01-01T00:00:00Z", we.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", we.CreatedAt.Format(time.RFC3339Nano))
}

func TestWebhooksEndpoint_Delete(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	_, err := Delete("webhooks_endpoint_id")
	assert.NoError(t, err)
}
