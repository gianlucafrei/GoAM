package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"goiam/internal/auth/graph"
	"html/template"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

// Configuration
var (
	baseTemplates      *template.Template
	LayoutTemplatePath = "templates/layout.html"
	NodeTemplatesPath  = "templates/nodes"
)

// ViewData is passed to all templates for dynamic rendering
type ViewData struct {
	Title        string
	NodeName     string
	Prompts      map[string]string
	Debug        bool
	Error        string
	StateJSON    string
	FlowName     string
	StylePath    string
	ScriptPath   string
	Message      string
	CustomConfig map[string]string
	Tenant       string
	Realm        string
	FlowPath     string
}

//go:embed templates/*
var templatesFS embed.FS

// InitTemplates loads and parses the base layout template
func InitTemplates() error {

	// Parse the base layout template
	tmpl, err := template.New("layout").Funcs(template.FuncMap{
		"title": title,
	}).ParseFS(templatesFS, LayoutTemplatePath)

	if err != nil {
		return fmt.Errorf("failed to parse base template: %w", err)
	}

	baseTemplates = tmpl
	return nil
}

// Render is the single public entry point
func Render(ctx *fasthttp.RequestCtx, flow *graph.FlowDefinition, state *graph.FlowState, resultNode *graph.GraphNode, prompts map[string]string) {
	var templateFile string
	var customMessage string

	switch {
	case resultNode != nil:
		templateFile = "result.html"
		customMessage = resultNode.CustomConfig["message"]

	case prompts != nil:
		templateFile = fmt.Sprintf("%s.html", state.Current)

	default:
		RenderError(ctx, "Unknown flow state")
		return
	}

	debug := isDebugMode(ctx)
	var stateJSON string
	if debug {
		if js, err := json.MarshalIndent(state, "", "  "); err == nil {
			stateJSON = string(js)
		}
	}

	currentGraphNode, ok := flow.Nodes[state.Current]
	if !ok {
		RenderError(ctx, "Did not find current graph node: "+state.Current)
		return
	}

	CustomConfig := currentGraphNode.CustomConfig
	if CustomConfig == nil {
		CustomConfig = make(map[string]string)
	}

	view := &ViewData{
		Title:        state.Current,
		NodeName:     state.Current,
		Prompts:      prompts,
		Debug:        debug,
		Error:        resolveErrorMessage(state),
		StateJSON:    stateJSON,
		FlowName:     flow.Name,
		Message:      customMessage,
		StylePath:    "../static/style.css",
		ScriptPath:   "../static/script.js",
		CustomConfig: CustomConfig,
		Tenant:       ctx.UserValue("tenant").(string),
		Realm:        ctx.UserValue("realm").(string),
		FlowPath:     ctx.UserValue("path").(string),
	}

	tmpl, err := baseTemplates.Clone()
	if err != nil {
		RenderError(ctx, "Template clone error: "+err.Error())
		return
	}

	filepath := filepath.Join(NodeTemplatesPath, templateFile)
	_, err = tmpl.ParseFS(templatesFS, filepath)
	if err != nil {
		RenderError(ctx, "Parse error: "+err.Error())
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", view); err != nil {
		RenderError(ctx, "Render error: "+err.Error())
		return
	}

	ctx.SetContentType("text/html")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(buf.Bytes())
}

func RenderError(ctx *fasthttp.RequestCtx, msg string) {
	ctx.SetContentType("text/html")
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetBodyString(fmt.Sprintf("<html><body><h2>Error</h2><p>%s</p></body></html>", msg))
}

func isDebugMode(ctx *fasthttp.RequestCtx) bool {

	debugParam := ctx.URI().QueryArgs().Peek("debug")
	return debugParam != nil
}

func resolveErrorMessage(state *graph.FlowState) string {
	if state.Error != nil {
		return *state.Error
	}
	return ""
}

func title(s string) string {
	if len(s) == 0 {
		return ""
	}
	return string(s[0]-32) + s[1:] // capitalize first letter
}
