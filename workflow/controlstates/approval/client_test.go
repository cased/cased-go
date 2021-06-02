package approval

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
	APIURL:    "https://api.cased.com/workflows/workflow_1sY67fxGiiHp4f7dlcY26pL4eY5",
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

func TestWorkflow_Cancel(t *testing.T) {
	e := cased.GetEndpoint(cased.WorkflowsEndpoint)
	cased.SetEndpoint(cased.WorkflowsEndpoint, mockedEndpoint)
	defer func(e cased.Endpoint) {
		cased.SetEndpoint(cased.WorkflowsEndpoint, e)
	}(e)

	c, err := Cancel("approval_control_state")

	assert.NoError(t, err)
	assert.Equal(t, "workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", c.ID)
	assert.Equal(t, "https://api.cased.com/workflows/workflow_1sY67fxGiiHp4f7dlcY26pL4eY5", c.APIURL)
	assert.Equal(t, "workflow", *c.Name)
	assert.Empty(t, c.Conditions)
	assert.Empty(t, c.Controls)
	assert.Equal(t, "2021-01-01T00:00:00Z", c.UpdatedAt.Format(time.RFC3339Nano))
	assert.Equal(t, "2021-01-01T00:00:00Z", c.CreatedAt.Format(time.RFC3339Nano))
}
