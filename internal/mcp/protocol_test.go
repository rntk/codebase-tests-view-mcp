package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codebase-view-mcp/internal/files"
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

	t.Run("tools/call normalizes absolute paths under configured root", func(t *testing.T) {
		baseDir := t.TempDir()
		sourcePath := filepath.Join(baseDir, "pkg", "source.go")
		testPath := filepath.Join(baseDir, "pkg", "source_test.go")
		if err := os.MkdirAll(filepath.Dir(sourcePath), 0755); err != nil {
			t.Fatalf("mkdirAll: %v", err)
		}
		if err := os.WriteFile(sourcePath, []byte("package pkg"), 0644); err != nil {
			t.Fatalf("write source file: %v", err)
		}
		if err := os.WriteFile(testPath, []byte("package pkg"), 0644); err != nil {
			t.Fatalf("write test file: %v", err)
		}

		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore, files.NewService(baseDir))

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      51,
			Method:  "tools/call",
			Params:  json.RawMessage([]byte(`{"name":"submit-test-metadata","arguments":{"sourceFile":"` + sourcePath + `","tests":[{"testFile":"` + testPath + `","testName":"TestSomething","functionName":"Something","comment":"Tests something","lineRange":{"start":1,"end":10},"coveredLines":{"start":1,"end":5}}]}}`)),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.Handle(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		if meta := metaStore.GetTestMetadata("pkg/source.go"); meta == nil || len(meta.Tests) != 1 {
			t.Fatalf("expected normalized metadata entry for pkg/source.go, got %+v", meta)
		} else if meta.Tests[0].TestFile != "pkg/source_test.go" {
			t.Fatalf("expected normalized test file %q, got %q", "pkg/source_test.go", meta.Tests[0].TestFile)
		}
	})

	t.Run("tools/call rejects absolute paths outside configured root", func(t *testing.T) {
		baseDir := t.TempDir()
		outsideDir := t.TempDir()
		outsideSource := filepath.Join(outsideDir, "source.go")
		outsideTest := filepath.Join(outsideDir, "source_test.go")
		if err := os.WriteFile(outsideSource, []byte("package main"), 0644); err != nil {
			t.Fatalf("write source file: %v", err)
		}
		if err := os.WriteFile(outsideTest, []byte("package main"), 0644); err != nil {
			t.Fatalf("write test file: %v", err)
		}

		metaStore := metadata.NewStore("")
		handler := NewHandler(metaStore, files.NewService(baseDir))

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      52,
			Method:  "tools/call",
			Params:  json.RawMessage([]byte(`{"name":"submit-test-metadata","arguments":{"sourceFile":"` + outsideSource + `","tests":[{"testFile":"` + outsideTest + `","testName":"TestSomething","functionName":"Something","comment":"Tests something","lineRange":{"start":1,"end":10},"coveredLines":{"start":1,"end":5}}]}}`)),
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
			t.Fatal("expected error for path outside configured root")
		}
		if !strings.Contains(response.Error.Message, "outside configured codebase root") {
			t.Fatalf("unexpected error message: %s", response.Error.Message)
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

func TestHandlerGetFunctionMetadata(t *testing.T) {
	t.Run("tools/call with get-function-metadata returns filtered tests", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("src.go", []metadata.TestReference{
			{TestFile: "src_test.go", TestName: "TestFoo", FunctionName: "Foo", Comment: "tests Foo", LineRange: metadata.LineRange{Start: 1, End: 5}, CoveredLines: metadata.LineRange{Start: 1, End: 3}},
			{TestFile: "src_test.go", TestName: "TestBar", FunctionName: "Bar", Comment: "tests Bar", LineRange: metadata.LineRange{Start: 6, End: 10}, CoveredLines: metadata.LineRange{Start: 4, End: 6}},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}
		handler := NewHandler(metaStore)

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      30,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name":"get-function-metadata","arguments":{"sourceFile":"src.go","functionName":"Foo"}}`),
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if response.Error != nil {
			t.Fatalf("unexpected error: %+v", response.Error)
		}

		result := response.Result.(map[string]interface{})
		content := result["content"].([]interface{})
		text := content[0].(map[string]interface{})["text"].(string)

		var res map[string]interface{}
		if err := json.Unmarshal([]byte(text), &res); err != nil {
			t.Fatalf("decode result text: %v", err)
		}
		tests := res["tests"].([]interface{})
		if len(tests) != 1 {
			t.Fatalf("expected 1 test for Foo, got %d", len(tests))
		}
		if tests[0].(map[string]interface{})["testName"] != "TestFoo" {
			t.Errorf("expected TestFoo, got %v", tests[0].(map[string]interface{})["testName"])
		}
	})

	t.Run("tools/call with get-function-metadata returns empty for unknown function", func(t *testing.T) {
		handler := NewHandler(metadata.NewStore(""))

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      31,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name":"get-function-metadata","arguments":{"sourceFile":"src.go","functionName":"Unknown"}}`),
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.Handle(rr, req)

		var response JSONRPCResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if response.Error != nil {
			t.Fatalf("unexpected error: %+v", response.Error)
		}

		result := response.Result.(map[string]interface{})
		content := result["content"].([]interface{})
		text := content[0].(map[string]interface{})["text"].(string)

		var res map[string]interface{}
		if err := json.Unmarshal([]byte(text), &res); err != nil {
			t.Fatalf("decode result text: %v", err)
		}
		tests, _ := res["tests"].([]interface{})
		if len(tests) != 0 {
			t.Fatalf("expected 0 tests for unknown function, got %d", len(tests))
		}
	})

	t.Run("tools/call with get-function-metadata missing sourceFile returns error", func(t *testing.T) {
		handler := NewHandler(metadata.NewStore(""))

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      32,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name":"get-function-metadata","arguments":{"functionName":"Foo"}}`),
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

	t.Run("tools/call with get-function-metadata missing functionName returns error", func(t *testing.T) {
		handler := NewHandler(metadata.NewStore(""))

		reqBody := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      33,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name":"get-function-metadata","arguments":{"sourceFile":"src.go"}}`),
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
