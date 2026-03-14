package mcp

import (
	"encoding/json"
	"testing"
)

func TestGetTools(t *testing.T) {
	tools := GetTools()

	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}

	submitTool := tools[0]
	if submitTool.Name != "submit-test-metadata" {
		t.Errorf("expected first tool name %q, got %q", "submit-test-metadata", submitTool.Name)
	}
	if submitTool.Description == "" {
		t.Error("submit-test-metadata description is empty")
	}

	var submitSchema map[string]interface{}
	if err := json.Unmarshal(submitTool.InputSchema, &submitSchema); err != nil {
		t.Fatalf("failed to unmarshal submit-test-metadata schema: %v", err)
	}
	assertSchemaHasRequiredFields(t, submitSchema, "sourceFile", "tests")

	getFuncTool := tools[1]
	if getFuncTool.Name != "get-function-metadata" {
		t.Errorf("expected second tool name %q, got %q", "get-function-metadata", getFuncTool.Name)
	}
	if getFuncTool.Description == "" {
		t.Error("get-function-metadata description is empty")
	}

	var getSchema map[string]interface{}
	if err := json.Unmarshal(getFuncTool.InputSchema, &getSchema); err != nil {
		t.Fatalf("failed to unmarshal get-function-metadata schema: %v", err)
	}
	assertSchemaHasRequiredFields(t, getSchema, "sourceFile", "functionName")
}

func TestGetToolsSchemaValidation(t *testing.T) {
	for _, tool := range GetTools() {
		t.Run(tool.Name, func(t *testing.T) {
			var schema map[string]interface{}
			if err := json.Unmarshal(tool.InputSchema, &schema); err != nil {
				t.Fatalf("failed to unmarshal schema: %v", err)
			}

			if schema["type"] != "object" {
				t.Errorf("expected schema type %q, got %q", "object", schema["type"])
			}

			props, ok := schema["properties"].(map[string]interface{})
			if !ok {
				t.Fatal("schema properties is not a map")
			}

			required, ok := schema["required"].([]interface{})
			if !ok {
				t.Fatal("schema required is not an array")
			}

			for _, req := range required {
				reqStr, ok := req.(string)
				if !ok {
					continue
				}
				if _, ok := props[reqStr]; !ok {
					t.Errorf("required field %q not found in properties", reqStr)
				}
			}
		})
	}
}

func TestGetToolsDescriptions(t *testing.T) {
	for _, tool := range GetTools() {
		if tool.Description == "" {
			t.Errorf("tool %s has empty description", tool.Name)
		}
		if len(tool.Description) < 20 {
			t.Errorf("tool %s description is too short: %q", tool.Name, tool.Description)
		}
	}
}

func assertSchemaHasRequiredFields(t *testing.T, schema map[string]interface{}, fields ...string) {
	t.Helper()

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("schema properties is not a map")
	}

	required, ok := schema["required"].([]interface{})
	if !ok {
		t.Fatal("schema required is not an array")
	}

	requiredSet := make(map[string]bool, len(required))
	for _, req := range required {
		if reqStr, ok := req.(string); ok {
			requiredSet[reqStr] = true
		}
	}

	for _, field := range fields {
		if _, ok := props[field]; !ok {
			t.Errorf("schema missing %s property", field)
		}
		if !requiredSet[field] {
			t.Errorf("%s not in required fields", field)
		}
	}
}
