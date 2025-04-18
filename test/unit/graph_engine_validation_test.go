package unit

import (
	"goiam/internal/auth/graph"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngine_ValidMinimalFlow(t *testing.T) {
	flow := &graph.FlowDefinition{
		Name:  "simple_init_only",
		Start: "init",
		Nodes: map[string]*graph.GraphNode{
			"init": {
				Name: "init",
				Use:  "init",
				Next: map[string]string{
					"start": "end",
				},
			},
			"end": {
				Name: "end",
				Use:  "successResult",
			},
		},
	}

	engine, err := graph.NewEngine(flow)
	assert.NoError(t, err)
	assert.NotNil(t, engine)
}

func TestEngine_MissingStartNode(t *testing.T) {
	flow := &graph.FlowDefinition{
		Name:  "no_start",
		Start: "init",
		Nodes: map[string]*graph.GraphNode{},
	}

	engine, err := graph.NewEngine(flow)
	assert.Nil(t, engine)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start node")
}

func TestEngine_StartNotInit(t *testing.T) {
	flow := &graph.FlowDefinition{
		Name:  "bad_start_type",
		Start: "askUsername",
		Nodes: map[string]*graph.GraphNode{
			"askUsername": {
				Name: "askUsername",
				Use:  "askUsername",
				Next: map[string]string{
					"submitted": "end",
				},
			},
			"end": {
				Name: "end",
				Use:  "successResult",
			},
		},
	}

	engine, err := graph.NewEngine(flow)
	assert.Nil(t, engine)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be of type 'init'")
}

func TestEngine_MissingNextOnLogicNode(t *testing.T) {
	flow := &graph.FlowDefinition{
		Name:  "missing_next",
		Start: "init",
		Nodes: map[string]*graph.GraphNode{
			"init": {
				Name: "init",
				Use:  "init",
				Next: map[string]string{
					"start": "logicStep",
				},
			},
			"logicStep": {
				Name: "logicStep",
				Use:  "validateUsernamePassword",
				// Next is missing here!
			},
		},
	}

	engine, err := graph.NewEngine(flow)
	assert.Nil(t, engine)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must define a 'Next' map")
}
