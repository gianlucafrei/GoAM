package model

import (
	"errors"
	"time"
)

// Enum of node types
type NodeType string

const (
	NodeTypeInit           NodeType = "init"
	NodeTypeQuery          NodeType = "query"
	NodeTypeQueryWithLogic NodeType = "queryWithLogic"
	NodeTypeLogic          NodeType = "logic"
	NodeTypeResult         NodeType = "result"
)

type NodeResult struct {
	Prompts   map[string]string // Prompts to be shown to the user, if applicable
	Condition string            // The next state, if applicable
}

type GraphNode struct {
	Name         string            `json:"name"`                    // unique in graph
	Use          string            `json:"use"`                     // reference to NodeDefinition.Name
	Next         map[string]string `json:"next"`                    // condition -> next GraphNode.Name
	CustomConfig map[string]string `json:"custom_config,omitempty"` // for overrides (optional)
}

// This is the flow definition, usually stored as a yaml file
type FlowDefinition struct {
	Description string                `json:"description"`
	Start       string                `json:"start"` // e.g., "init"
	Nodes       map[string]*GraphNode `json:"nodes"`
}

// This is a flow together with meta information such as route, realm and tenant.
// It contains the flow defition which is a yaml file
type Flow struct {
	Tenant             string          `json:"tenant" yaml:"tenant"`                                               // e.g. "acme"
	Realm              string          `json:"realm" yaml:"realm"`                                                 // e.g. "customers"
	Id                 string          `json:"id" yaml:"id"`                                                       // e.g. "login"
	Route              string          `json:"route" yaml:"route"`                                                 // e.g. "/login"
	Active             bool            `json:"active" yaml:"active"`                                               // whether the flow is active
	Definition         *FlowDefinition `json:"-" yaml:"-"`                                                         // pre-loaded flow definition
	DefinitionYaml     string          `json:"-" yaml:"-"`                                                         // original yaml content, we keep that in order to perserve the exactly same yaml
	DefinitionLocation string          `json:"definition_location,omitempty" yaml:"definition_location,omitempty"` // path to the yaml file
	CreatedAt          time.Time       `json:"created_at" yaml:"-"`
	UpdatedAt          time.Time       `json:"updated_at" yaml:"-"`
}

type AuthLevel string

const (
	AuthLevelUnauthenticated AuthLevel = "0"
	AuthLevel1FA             AuthLevel = "1"
	AuthLevel2FA             AuthLevel = "2"
)

type FlowResult struct {
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	Authenticated bool      `json:"authenticated"`
	AuthLevel     AuthLevel `json:"auth_level"`
	FlowName      string    `json:"flow_name,omitempty"`
}

// Create NodeResult with state
func NewNodeResultWithCondition(condition string) (*NodeResult, error) {
	return &NodeResult{
		Prompts:   nil,
		Condition: condition,
	}, nil
}

// Create NodeResult with prompts
func NewNodeResultWithPrompts(prompts map[string]string) (*NodeResult, error) {
	return &NodeResult{
		Prompts:   prompts,
		Condition: "",
	}, nil
}

// Create NodeResult with error
func NewNodeResultWithError(err error) (*NodeResult, error) {
	return nil, err
}

// Create NodeResult with text error
func NewNodeResultWithTextError(text string) (*NodeResult, error) {
	return nil, errors.New(text)
}
