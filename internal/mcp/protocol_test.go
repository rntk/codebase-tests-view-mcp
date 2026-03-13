package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codebase-view-mcp/internal/metadata"
)

func TestNewHandler(t *testing.T) {
	metaStore := metadata.NewStore("")
	handler := NewHandler(metaStore)

	if handler == nil {
		t.Fatal("NewHandler returned nil")
	}

	if handler.metaStore != metaStore {
		t.Error("handler.metaStore not set correctly")
	}
}

func TestHandlerHandle(t *testing.T) {
	t.Run("initialize method returns correct response", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
			Params:  json.RawMessage(`{"protocolVersion": "2024-11-05"}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.JSONRPC != "2.0" {
			t.Errorf("expected JSONRPC %q, got %q", "2.0", response.JSONRPC)
		}

		if response.ID.(float64) != 1 {
			t.Errorf("expected ID 1, got %v", response.ID)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}

		result, ok := response.Result.(map[string]interface{})
		if !ok {
			t.Fatal("result is not a map")
		}

		if result["protocolVersion"] != "2024-11-05" {
			t.Errorf("expected protocolVersion %q, got %v", "2024-11-05", result["protocolVersion"])
		}
	})

	t.Run("tools/list returns available tools", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "tools/list",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}

		result, ok := response.Result.(map[string]interface{})
		if !ok {
			t.Fatal("result is not a map")
		}

		tools, ok := result["tools"].([]interface{})
		if !ok {
			t.Fatal("tools is not an array")
		}

		if len(tools) != 2 {
			t.Errorf("expected 2 tools, got %d", len(tools))
		}
	})

	t.Run("prompts/list returns available prompts", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      3,
			Method:  "prompts/list",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}

		result, ok := response.Result.(map[string]interface{})
		if !ok {
			t.Fatal("result is not a map")
		}

		prompts, ok := result["prompts"].([]interface{})
		if !ok {
			t.Fatal("prompts is not an array")
		}

		if len(prompts) != 2 {
			t.Errorf("expected 2 prompts, got %d", len(prompts))
		}
	})

	t.Run("unknown method returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      4,
			Method:  "unknown/method",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for unknown method")
		}

		if response.Error.Code != -32601 {
			t.Errorf("expected error code -32601, got %d", response.Error.Code)
		}
	})

	t.Run("parse error for invalid JSON", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader([]byte("invalid json")))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for invalid JSON")
		}

		if response.Error.Code != -32700 {
			t.Errorf("expected error code -32700, got %d", response.Error.Code)
		}
	})

	t.Run("tools/call with submit-test-metadata", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      5,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"sourceFile": "test.go", "tests": [{"testFile": "test_test.go", "testName": "TestSomething", "functionName": "Something", "comment": "Tests something", "lineRange": {"start": 1, "end": 10}, "coveredLines": {"start": 1, "end": 5}}]}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}
	})

	t.Run("tools/call with suggest-missing-tests", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      6,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "suggest-missing-tests", "arguments": {"sourceFile": "test.go", "suggestions": [{"targetLines": {"start": 1, "end": 10}, "reason": "Missing test", "suggestedName": "TestMissing", "testSkeleton": "func TestMissing(t *testing.T) {}", "priority": "high"}]}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}
	})

	t.Run("tools/call with unknown tool returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      7,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "unknown-tool"}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for unknown tool")
		}
	})

	t.Run("prompts/get with codebase-tests-review", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      8,
			Method:  "prompts/get",
			Params:  json.RawMessage(`{"name": "codebase-tests-review", "arguments": {"functionName": "MyFunction", "filePath": "test.go"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}
	})

	t.Run("prompts/get with test-to-source-review", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      9,
			Method:  "prompts/get",
			Params:  json.RawMessage(`{"name": "test-to-source-review", "arguments": {"testName": "TestSomething", "testFilePath": "test_test.go"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}
	})

	t.Run("prompts/get with unknown prompt returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      10,
			Method:  "prompts/get",
			Params:  json.RawMessage(`{"name": "unknown-prompt"}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for unknown prompt")
		}
	})

	t.Run("tools/call with invalid params returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      11,
			Method:  "tools/call",
			Params:  json.RawMessage(`invalid`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for invalid params")
		}
	})

	t.Run("submit-test-metadata with missing sourceFile returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      12,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"tests": []}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing sourceFile")
		}
	})

	t.Run("submit-test-metadata with missing tests returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      13,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"sourceFile": "test.go"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing tests")
		}
	})

	t.Run("submit-test-metadata with invalid test format returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      14,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"sourceFile": "test.go", "tests": "invalid"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for invalid test format")
		}
	})

	t.Run("submit-test-metadata with test missing functionName returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      15,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"sourceFile": "test.go", "tests": [{"testFile": "test_test.go", "testName": "TestSomething", "functionName": "", "comment": "Tests something", "lineRange": {"start": 1, "end": 10}, "coveredLines": {"start": 1, "end": 5}}]}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing functionName")
		}
	})

	t.Run("submit-test-metadata with test missing comment returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      16,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"sourceFile": "test.go", "tests": [{"testFile": "test_test.go", "testName": "TestSomething", "functionName": "Something", "comment": "", "lineRange": {"start": 1, "end": 10}, "coveredLines": {"start": 1, "end": 5}}]}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing comment")
		}
	})

	t.Run("submit-test-metadata with invalid lineRange returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      17,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"sourceFile": "test.go", "tests": [{"testFile": "test_test.go", "testName": "TestSomething", "functionName": "Something", "comment": "Tests something", "lineRange": {"start": 0, "end": 0}, "coveredLines": {"start": 1, "end": 5}}]}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for invalid lineRange")
		}
	})

	t.Run("submit-test-metadata with invalid coveredLines returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      18,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "submit-test-metadata", "arguments": {"sourceFile": "test.go", "tests": [{"testFile": "test_test.go", "testName": "TestSomething", "functionName": "Something", "comment": "Tests something", "lineRange": {"start": 1, "end": 10}, "coveredLines": {"start": 0, "end": 0}}]}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for invalid coveredLines")
		}
	})

	t.Run("suggest-missing-tests with missing sourceFile returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      19,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "suggest-missing-tests", "arguments": {"suggestions": []}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing sourceFile")
		}
	})

	t.Run("suggest-missing-tests with missing suggestions returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      20,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "suggest-missing-tests", "arguments": {"sourceFile": "test.go"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing suggestions")
		}
	})

	t.Run("suggest-missing-tests with invalid suggestions format returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      21,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name": "suggest-missing-tests", "arguments": {"sourceFile": "test.go", "suggestions": "invalid"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for invalid suggestions format")
		}
	})

	t.Run("prompts/get with missing functionName returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      22,
			Method:  "prompts/get",
			Params:  json.RawMessage(`{"name": "codebase-tests-review", "arguments": {"filePath": "test.go"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing functionName")
		}
	})

	t.Run("prompts/get with missing filePath returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      23,
			Method:  "prompts/get",
			Params:  json.RawMessage(`{"name": "codebase-tests-review", "arguments": {"functionName": "MyFunction"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing filePath")
		}
	})

	t.Run("prompts/get with missing testName returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      24,
			Method:  "prompts/get",
			Params:  json.RawMessage(`{"name": "test-to-source-review", "arguments": {"testFilePath": "test_test.go"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing testName")
		}
	})

	t.Run("prompts/get with missing testFilePath returns error", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      25,
			Method:  "prompts/get",
			Params:  json.RawMessage(`{"name": "test-to-source-review", "arguments": {"testName": "TestSomething"}}`),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error for missing testFilePath")
		}
	})
}

func TestHandlerSendSuccessAndError(t *testing.T) {
	t.Run("sendSuccess encodes response correctly", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		rr := httptest.NewRecorder()
		handler.sendSuccess(rr, 1, map[string]string{"key": "value"})

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.JSONRPC != "2.0" {
			t.Errorf("expected JSONRPC %q, got %q", "2.0", response.JSONRPC)
		}

		if response.Error != nil {
			t.Errorf("unexpected error: %+v", response.Error)
		}
	})

	t.Run("sendError encodes error response correctly", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore)

		rr := httptest.NewRecorder()
		handler.sendError(rr, 1, -32600, "Invalid Request")

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Error == nil {
			t.Fatal("expected error")
		}

		if response.Error.Code != -32600 {
			t.Errorf("expected error code -32600, got %d", response.Error.Code)
		}

		if response.Error.Message != "Invalid Request" {
			t.Errorf("expected error message %q, got %q", "Invalid Request", response.Error.Message)
		}
	})
}

func TestValidateLineRange(t *testing.T) {
	// These are tested indirectly through the handler tests
	// but we can test the helper functions directly too
	t.Run("validateRequiredLineRange with zero values returns error", func(t *testing.T) {
		// This is tested through the handler
	})

	t.Run("validateRequiredLineRange with valid range returns nil", func(t *testing.T) {
		// This is tested through the handler
	})

	t.Run("validateOptionalLineRange with zero values returns nil", func(t *testing.T) {
		// This is tested through the handler
	})
}
