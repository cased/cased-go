package workflow

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cased/cased-go"
	"github.com/stretchr/testify/assert"
)

var mockedWorkflow = cased.Workflow{
	ID:        "workflow_1sY67fxGiiHp4f7dlcY26pL4eY5",
	Name:      cased.String("workflow"),
	ApiURL:    "https://api.cased.com/workflows/workflow_1sY67fxGiiHp4f7dlcY26pL4eY5",
	UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
}

var mockedEndpoint = &MockEndpoint{
	CallFunc: func(method, path string, params cased.ParamsContainer, i interface{}) error {
		data, err := json.Marshal(mockedWorkflow)
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

	w, err := New(&cased.WorkflowParams{
		Name: cased.String("workflow"),
	})

	assert.NoError(t, err)
	assert.Equal(t, "workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", w.ID)
	assert.Equal(t, "https://api.cased.com/workflows/workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", w.ApiURL)
	assert.Equal(t, "workflow", *w.Name)
	assert.Empty(t, w.Conditions)
	assert.Empty(t, w.Controls)
	assert.Equal(t, "2021-01-01T00:00:00Z", w.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", w.CreatedAt.Format(time.RFC3339Nano))
}

func TestWorkflow_Get(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	w, err := Get("workflow_id")

	assert.NoError(t, err)
	assert.Equal(t, "workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", w.ID)
	assert.Equal(t, "https://api.cased.com/workflows/workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", w.ApiURL)
	assert.Equal(t, "workflow", *w.Name)
	assert.Empty(t, w.Conditions)
	assert.Empty(t, w.Controls)
	assert.Equal(t, "2021-01-01T00:00:00Z", w.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", w.CreatedAt.Format(time.RFC3339Nano))
}

func TestWorkflow_Update(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	w, err := Update("workflow_id", &cased.WorkflowParams{})

	assert.NoError(t, err)
	assert.Equal(t, "workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", w.ID)
	assert.Equal(t, "https://api.cased.com/workflows/workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", w.ApiURL)
	assert.Equal(t, "workflow", *w.Name)
	assert.Empty(t, w.Conditions)
	assert.Empty(t, w.Controls)
	assert.Equal(t, "2021-01-01T00:00:00Z", w.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", w.CreatedAt.Format(time.RFC3339Nano))
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
