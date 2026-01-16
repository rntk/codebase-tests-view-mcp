package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codebase-view-mcp/internal/metadata"
)

// JSON-RPC 2.0 types

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP-specific types

// InitializeParams for initialize method
type InitializeParams struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ClientInfo      ClientInfo   `json:"clientInfo"`
}

// Capabilities represents client/server capabilities
type Capabilities struct {
	Tools   *ToolsCapability   `json:"tools,omitempty"`
	Prompts *PromptsCapability `json:"prompts,omitempty"`
}

// ToolsCapability indicates tools support
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// PromptsCapability indicates prompts support
type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ClientInfo contains client information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult for initialize response
type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
}

// ServerInfo contains server information
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ToolsListResult for tools/list response
type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// ToolsCallParams for tools/call request
type ToolsCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolsCallResult for tools/call response
type ToolsCallResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ContentItem represents a piece of content
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// PromptsListResult for prompts/list response
type PromptsListResult struct {
	Prompts []Prompt `json:"prompts"`
}

// Prompt represents an MCP prompt
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument defines a prompt argument
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptsGetParams for prompts/get request
type PromptsGetParams struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments,omitempty"`
}

// PromptsGetResult for prompts/get response
type PromptsGetResult struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}

// PromptMessage represents a message in a prompt
type PromptMessage struct {
	Role    string      `json:"role"`
	Content TextContent `json:"content"`
}

// TextContent represents text content
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Handler handles MCP protocol requests
type Handler struct {
	metaStore *metadata.Store
}

// NewHandler creates a new MCP handler
func NewHandler(metaStore *metadata.Store) *Handler {
	return &Handler{
		metaStore: metaStore,
	}
}

// Handle processes an MCP request
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, nil, -32700, "Parse error")
		return
	}

	log.Printf("MCP Request: %s", req.Method)

	var result interface{}
	var err error

	switch req.Method {
	case "initialize":
		result, err = h.handleInitialize(req.Params)
	case "tools/list":
		result, err = h.handleToolsList()
	case "tools/call":
		result, err = h.handleToolsCall(req.Params)
	case "prompts/list":
		result, err = h.handlePromptsList()
	case "prompts/get":
		result, err = h.handlePromptsGet(req.Params)
	default:
		h.sendError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
		return
	}

	if err != nil {
		h.sendError(w, req.ID, -32603, err.Error())
		return
	}

	h.sendSuccess(w, req.ID, result)
}

// handleInitialize handles the initialize request
func (h *Handler) handleInitialize(params json.RawMessage) (interface{}, error) {
	return InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: Capabilities{
			Tools:   &ToolsCapability{},
			Prompts: &PromptsCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "codebase-view-mcp",
			Version: "1.0.0",
		},
	}, nil
}

// handleToolsList handles the tools/list request
func (h *Handler) handleToolsList() (interface{}, error) {
	return ToolsListResult{
		Tools: GetTools(),
	}, nil
}

// handleToolsCall handles the tools/call request
func (h *Handler) handleToolsCall(params json.RawMessage) (interface{}, error) {
	var callParams ToolsCallParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	switch callParams.Name {
	case "submit-test-metadata":
		return h.executeSubmitTestMetadata(callParams.Arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", callParams.Name)
	}
}

// executeSubmitTestMetadata executes the submit-test-metadata tool
func (h *Handler) executeSubmitTestMetadata(args map[string]interface{}) (interface{}, error) {
	// Extract sourceFile
	sourceFile, ok := args["sourceFile"].(string)
	if !ok {
		return nil, fmt.Errorf("sourceFile is required and must be a string")
	}

	// Extract tests array
	testsRaw, ok := args["tests"]
	if !ok {
		return nil, fmt.Errorf("tests is required")
	}

	testsJSON, err := json.Marshal(testsRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tests: %w", err)
	}

	var tests []metadata.TestReference
	if err := json.Unmarshal(testsJSON, &tests); err != nil {
		return nil, fmt.Errorf("invalid tests format: %w", err)
	}

	// Store metadata (merge with existing tests)
	if err := h.metaStore.AddTestMetadata(sourceFile, tests); err != nil {
		return nil, fmt.Errorf("failed to store metadata: %w", err)
	}

	return ToolsCallResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: fmt.Sprintf("Successfully stored test metadata for %s (%d tests)", sourceFile, len(tests)),
			},
		},
	}, nil
}

// handlePromptsList handles the prompts/list request
func (h *Handler) handlePromptsList() (interface{}, error) {
	return PromptsListResult{
		Prompts: GetPrompts(),
	}, nil
}

// handlePromptsGet handles the prompts/get request
func (h *Handler) handlePromptsGet(params json.RawMessage) (interface{}, error) {
	var getParams PromptsGetParams
	if err := json.Unmarshal(params, &getParams); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	messages, err := GetPromptContent(getParams.Name, getParams.Arguments)
	if err != nil {
		return nil, err
	}

	return PromptsGetResult{
		Messages: messages,
	}, nil
}

// sendSuccess sends a successful JSON-RPC response
func (h *Handler) sendSuccess(w http.ResponseWriter, id interface{}, result interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendError sends an error JSON-RPC response
func (h *Handler) sendError(w http.ResponseWriter, id interface{}, code int, message string) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
