package unit

import (
	"goiam/internal/auth/graph"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunPasskeyRegisterNode_HappyPath(t *testing.T) {
	// Setup flow state with a username
	state := &graph.FlowState{
		Context: map[string]string{
			"username": "alice",
		},
	}

	// Minimal graph node (not used in logic here, but required)
	node := &graph.GraphNode{
		Name: "registerPasskey",
		Use:  "registerPasskey",
	}

	// Run the node logic
	prompts, err := graph.GeneratePasskeysOptions(state, node)

	// Validate result
	assert.NoError(t, err)
	assert.Contains(t, prompts, "passkeysOptions")

	// Check if JSON strings are stored in context
	sessionStr, sessionExists := state.Context["passkeysSession"]
	optionsStr, optionsExists := state.Context["passkeysOptions"]

	assert.True(t, sessionExists, "passkeysSession should be set in context")
	assert.True(t, optionsExists, "passkeysOptions should be set in context")
	assert.NotEmpty(t, sessionStr)
	assert.NotEmpty(t, optionsStr)
}
