package mcp

import (
	"strings"
	"testing"
)

func TestGetPrompts(t *testing.T) {
	prompts := GetPrompts()

	if len(prompts) != 2 {
		t.Fatalf("expected 2 prompts, got %d", len(prompts))
	}

	// Check codebase-tests-review prompt
	codebasePrompt := prompts[0]
	if codebasePrompt.Name != "codebase-tests-review" {
		t.Errorf("expected first prompt name %q, got %q", "codebase-tests-review", codebasePrompt.Name)
	}

	if codebasePrompt.Description == "" {
		t.Error("codebase-tests-review description is empty")
	}

	if len(codebasePrompt.Arguments) != 2 {
		t.Errorf("expected 2 arguments, got %d", len(codebasePrompt.Arguments))
	}

	foundFunctionName := false
	foundFilePath := false
	for _, arg := range codebasePrompt.Arguments {
		if arg.Name == "functionName" {
			foundFunctionName = true
			if arg.Description == "" {
				t.Error("functionName argument description is empty")
			}
			if !arg.Required {
				t.Error("functionName should be required")
			}
		}
		if arg.Name == "filePath" {
			foundFilePath = true
			if arg.Description == "" {
				t.Error("filePath argument description is empty")
			}
			if !arg.Required {
				t.Error("filePath should be required")
			}
		}
	}

	if !foundFunctionName {
		t.Error("functionName argument not found")
	}

	if !foundFilePath {
		t.Error("filePath argument not found")
	}

	// Check test-to-source-review prompt
	testToSourcePrompt := prompts[1]
	if testToSourcePrompt.Name != "test-to-source-review" {
		t.Errorf("expected second prompt name %q, got %q", "test-to-source-review", testToSourcePrompt.Name)
	}

	if testToSourcePrompt.Description == "" {
		t.Error("test-to-source-review description is empty")
	}

	if len(testToSourcePrompt.Arguments) != 2 {
		t.Errorf("expected 2 arguments, got %d", len(testToSourcePrompt.Arguments))
	}

	foundTestName := false
	foundTestFilePath := false
	for _, arg := range testToSourcePrompt.Arguments {
		if arg.Name == "testName" {
			foundTestName = true
			if arg.Description == "" {
				t.Error("testName argument description is empty")
			}
			if !arg.Required {
				t.Error("testName should be required")
			}
		}
		if arg.Name == "testFilePath" {
			foundTestFilePath = true
			if arg.Description == "" {
				t.Error("testFilePath argument description is empty")
			}
			if !arg.Required {
				t.Error("testFilePath should be required")
			}
		}
	}

	if !foundTestName {
		t.Error("testName argument not found")
	}

	if !foundTestFilePath {
		t.Error("testFilePath argument not found")
	}
}

func TestGetPromptContent(t *testing.T) {
	t.Run("codebase-tests-review returns prompt with functionName and filePath", func(t *testing.T) {
		args := map[string]string{
			"functionName": "MyFunction",
			"filePath":     "test.go",
		}

		messages, err := GetPromptContent("codebase-tests-review", args)
		if err != nil {
			t.Fatalf("GetPromptContent failed: %v", err)
		}

		if len(messages) != 1 {
			t.Fatalf("expected 1 message, got %d", len(messages))
		}

		msg := messages[0]
		if msg.Role != "user" {
			t.Errorf("expected role %q, got %q", "user", msg.Role)
		}

		if msg.Content.Type != "text" {
			t.Errorf("expected content type %q, got %q", "text", msg.Content.Type)
		}

		if msg.Content.Text == "" {
			t.Error("prompt text is empty")
		}

		// Check that the prompt contains the function name and file path
		if !strings.Contains(msg.Content.Text, "MyFunction") {
			t.Error("prompt text does not contain function name")
		}

		if !strings.Contains(msg.Content.Text, "test.go") {
			t.Error("prompt text does not contain file path")
		}

		// Check that the prompt mentions submit-test-metadata
		if !strings.Contains(msg.Content.Text, "submit-test-metadata") {
			t.Error("prompt text does not mention submit-test-metadata")
		}
	})

	t.Run("test-to-source-review returns prompt with testName and testFilePath", func(t *testing.T) {
		args := map[string]string{
			"testName":     "TestMyFunction",
			"testFilePath": "test_test.go",
		}

		messages, err := GetPromptContent("test-to-source-review", args)
		if err != nil {
			t.Fatalf("GetPromptContent failed: %v", err)
		}

		if len(messages) != 1 {
			t.Fatalf("expected 1 message, got %d", len(messages))
		}

		msg := messages[0]
		if msg.Role != "user" {
			t.Errorf("expected role %q, got %q", "user", msg.Role)
		}

		if msg.Content.Type != "text" {
			t.Errorf("expected content type %q, got %q", "text", msg.Content.Type)
		}

		if msg.Content.Text == "" {
			t.Error("prompt text is empty")
		}

		// Check that the prompt contains the test name and file path
		if !strings.Contains(msg.Content.Text, "TestMyFunction") {
			t.Error("prompt text does not contain test name")
		}

		if !strings.Contains(msg.Content.Text, "test_test.go") {
			t.Error("prompt text does not contain test file path")
		}
	})

	t.Run("codebase-tests-review with missing functionName returns error", func(t *testing.T) {
		args := map[string]string{
			"filePath": "test.go",
		}

		_, err := GetPromptContent("codebase-tests-review", args)
		if err == nil {
			t.Fatal("expected error for missing functionName")
		}

		if err.Error() != "functionName argument is required" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("codebase-tests-review with missing filePath returns error", func(t *testing.T) {
		args := map[string]string{
			"functionName": "MyFunction",
		}

		_, err := GetPromptContent("codebase-tests-review", args)
		if err == nil {
			t.Fatal("expected error for missing filePath")
		}

		if err.Error() != "filePath argument is required" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("test-to-source-review with missing testName returns error", func(t *testing.T) {
		args := map[string]string{
			"testFilePath": "test_test.go",
		}

		_, err := GetPromptContent("test-to-source-review", args)
		if err == nil {
			t.Fatal("expected error for missing testName")
		}

		if err.Error() != "testName argument is required" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("test-to-source-review with missing testFilePath returns error", func(t *testing.T) {
		args := map[string]string{
			"testName": "TestMyFunction",
		}

		_, err := GetPromptContent("test-to-source-review", args)
		if err == nil {
			t.Fatal("expected error for missing testFilePath")
		}

		if err.Error() != "testFilePath argument is required" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("unknown prompt returns error", func(t *testing.T) {
		args := map[string]string{}

		_, err := GetPromptContent("unknown-prompt", args)
		if err == nil {
			t.Fatal("expected error for unknown prompt")
		}

		if !strings.Contains(err.Error(), "prompt not found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("codebase-tests-review prompt includes line number instructions", func(t *testing.T) {
		args := map[string]string{
			"functionName": "MyFunction",
			"filePath":     "test.go",
		}

		messages, err := GetPromptContent("codebase-tests-review", args)
		if err != nil {
			t.Fatalf("GetPromptContent failed: %v", err)
		}

		msg := messages[0]
		// Check that the prompt mentions line numbers
		if !strings.Contains(msg.Content.Text, "line number") && !strings.Contains(msg.Content.Text, "cat -n") {
			t.Error("prompt text does not mention line numbers")
		}
	})

	t.Run("test-to-source-review prompt includes line number instructions", func(t *testing.T) {
		args := map[string]string{
			"testName":     "TestMyFunction",
			"testFilePath": "test_test.go",
		}

		messages, err := GetPromptContent("test-to-source-review", args)
		if err != nil {
			t.Fatalf("GetPromptContent failed: %v", err)
		}

		msg := messages[0]
		// Check that the prompt mentions line numbers
		if !strings.Contains(msg.Content.Text, "line number") && !strings.Contains(msg.Content.Text, "cat -n") {
			t.Error("prompt text does not mention line numbers")
		}
	})
}

func TestGetPromptsDescriptions(t *testing.T) {
	prompts := GetPrompts()

	for _, prompt := range prompts {
		if prompt.Description == "" {
			t.Errorf("prompt %s has empty description", prompt.Name)
		}

		// Description should be meaningful (at least 20 chars)
		if len(prompt.Description) < 20 {
			t.Errorf("prompt %s description is too short: %q", prompt.Name, prompt.Description)
		}
	}
}

func TestGetPromptContentEmptyArgs(t *testing.T) {
	// Test with empty args map
	_, err := GetPromptContent("codebase-tests-review", map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty args")
	}

	_, err = GetPromptContent("test-to-source-review", map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty args")
	}
}
