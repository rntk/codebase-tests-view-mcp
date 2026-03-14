package mcp

import (
	"encoding/json"
	"testing"
)

func TestGetTools(t *testing.T) {
	tools := GetTools()

	if len(tools) != 3 {
		t.Fatalf("expected 3 tools, got %d", len(tools))
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

	// Check get-function-metadata tool
	var getFuncTool Tool
	for _, tool := range tools {
		if tool.Name == "get-function-metadata" {
			getFuncTool = tool
			break
		}
	}
	if getFuncTool.Name != "get-function-metadata" {
		t.Errorf("expected third tool name %q, got %q", "get-function-metadata", getFuncTool.Name)
	}
	if getFuncTool.Description == "" {
		t.Error("get-function-metadata description is empty")
	}
	var schema3 map[string]interface{}
	if err := json.Unmarshal(getFuncTool.InputSchema, &schema3); err != nil {
		t.Fatalf("failed to unmarshal get-function-metadata schema: %v", err)
	}
	props3, ok := schema3["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("get-function-metadata schema properties is not a map")
	}
	if _, ok := props3["sourceFile"]; !ok {
		t.Error("get-function-metadata schema missing sourceFile property")
	}
	if _, ok := props3["functionName"]; !ok {
		t.Error("get-function-metadata schema missing functionName property")
	}
	required3, ok := schema3["required"].([]interface{})
	if !ok {
		t.Fatal("get-function-metadata schema required is not an array")
	}
	foundFuncName := false
	foundSrcFile := false
	for _, req := range required3 {
		if req == "functionName" {
			foundFuncName = true
		}
		if req == "sourceFile" {
			foundSrcFile = true
		}
	}
	if !foundFuncName {
		t.Error("functionName not in get-function-metadata required fields")
	}
	if !foundSrcFile {
		t.Error("sourceFile not in get-function-metadata required fields")
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
