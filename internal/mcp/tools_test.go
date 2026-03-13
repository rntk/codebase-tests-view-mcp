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

	// Check submit-test-metadata tool
	submitTool := tools[0]
	if submitTool.Name != "submit-test-metadata" {
		t.Errorf("expected first tool name %q, got %q", "submit-test-metadata", submitTool.Name)
	}

	if submitTool.Description == "" {
		t.Error("submit-test-metadata description is empty")
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(submitTool.InputSchema, &schema); err != nil {
		t.Fatalf("failed to unmarshal submit-test-metadata schema: %v", err)
	}

	if schema["type"] != "object" {
		t.Errorf("expected schema type %q, got %q", "object", schema["type"])
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("schema properties is not a map")
	}

	if _, ok := props["sourceFile"]; !ok {
		t.Error("schema missing sourceFile property")
	}

	if _, ok := props["tests"]; !ok {
		t.Error("schema missing tests property")
	}

	required, ok := schema["required"].([]interface{})
	if !ok {
		t.Fatal("schema required is not an array")
	}

	foundSourceFile := false
	foundTests := false
	for _, req := range required {
		if req == "sourceFile" {
			foundSourceFile = true
		}
		if req == "tests" {
			foundTests = true
		}
	}

	if !foundSourceFile {
		t.Error("sourceFile not in required fields")
	}

	if !foundTests {
		t.Error("tests not in required fields")
	}

	// Check suggest-missing-tests tool
	suggestTool := tools[1]
	if suggestTool.Name != "suggest-missing-tests" {
		t.Errorf("expected second tool name %q, got %q", "suggest-missing-tests", suggestTool.Name)
	}

	if suggestTool.Description == "" {
		t.Error("suggest-missing-tests description is empty")
	}

	var schema2 map[string]interface{}
	if err := json.Unmarshal(suggestTool.InputSchema, &schema2); err != nil {
		t.Fatalf("failed to unmarshal suggest-missing-tests schema: %v", err)
	}

	if schema2["type"] != "object" {
		t.Errorf("expected schema type %q, got %q", "object", schema2["type"])
	}

	props2, ok := schema2["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("schema properties is not a map")
	}

	if _, ok := props2["sourceFile"]; !ok {
		t.Error("schema missing sourceFile property")
	}

	if _, ok := props2["suggestions"]; !ok {
		t.Error("schema missing suggestions property")
	}

	// Check that functionName is optional (not in required)
	required2, ok := schema2["required"].([]interface{})
	if !ok {
		t.Fatal("schema required is not an array")
	}

	foundFunctionName := false
	for _, req := range required2 {
		if req == "functionName" {
			foundFunctionName = true
		}
	}

	if foundFunctionName {
		t.Error("functionName should be optional")
	}
}

func TestGetToolsSchemaValidation(t *testing.T) {
	tools := GetTools()

	for _, tool := range tools {
		t.Run(tool.Name, func(t *testing.T) {
			var schema map[string]interface{}
			if err := json.Unmarshal(tool.InputSchema, &schema); err != nil {
				t.Fatalf("failed to unmarshal schema: %v", err)
			}

			// Verify schema structure
			if schema["type"] != "object" {
				t.Errorf("expected schema type %q, got %q", "object", schema["type"])
			}

			props, ok := schema["properties"].(map[string]interface{})
			if !ok {
				t.Fatal("schema properties is not a map")
			}

			// Check that required fields exist in properties
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
	tools := GetTools()

	for _, tool := range tools {
		if tool.Description == "" {
			t.Errorf("tool %s has empty description", tool.Name)
		}

		// Description should be meaningful (at least 20 chars)
		if len(tool.Description) < 20 {
			t.Errorf("tool %s description is too short: %q", tool.Name, tool.Description)
		}
	}
}
