package mcp

import "fmt"

// GetPrompts returns all available MCP prompts
func GetPrompts() []Prompt {
	return []Prompt{
		{
			Name:        "codebase-tests-review",
			Description: "Analyze a function and submit metadata about its tests using the submit-test-metadata tool",
			Arguments: []PromptArgument{
				{
					Name:        "functionName",
					Description: "The name of the function to analyze",
					Required:    true,
				},
				{
					Name:        "filePath",
					Description: "Path to the file containing the function",
					Required:    true,
				},
			},
		},
	}
}

// GetPromptContent returns the prompt messages with arguments filled in
func GetPromptContent(name string, args map[string]string) ([]PromptMessage, error) {
	if name == "codebase-tests-review" {
		functionName := args["functionName"]
		filePath := args["filePath"]

		if functionName == "" {
			return nil, fmt.Errorf("functionName argument is required")
		}
		if filePath == "" {
			return nil, fmt.Errorf("filePath argument is required")
		}

		promptText := fmt.Sprintf(`Please analyze the **%s** function in file **%s**.

1. Examine the function's implementation.
2. If the function has associated tests, use the **submit-test-metadata** tool to submit metadata.
3. For each test, identify:
   * Which part of the source function the test exercises (line range in the source file).
   * The specific input data used in the test (line numbers in test file).
   * The expected result in the test (line numbers in test file).
   * A brief comment/description of what the test verifies.

When reporting the analysis, include **only** the following information for each test:

- File name 
- Test name
- Comment
- Covered lines (source)
- Line numbers (test range)
- Line numbers (input data)
- Line numbers (expected result)

After identifying all tests, use the **submit-test-metadata** tool with the following structure:

{
  "sourceFile": "%s",
  "tests": [
    {
      "testFile": "path/to/test_file.go",
      "testName": "TestFunctionName",
      "comment": "Brief description of what the test verifies",
      "lineRange": {"start": 10, "end": 25},
      "coveredLines": {"start": 45, "end": 60},
      "inputLines": {"start": 12, "end": 15},
      "outputLines": {"start": 20, "end": 22}
    }
  ]
}

**IMPORTANT**:
- "lineRange" refers to the lines in the TEST file where the test code is located
- "coveredLines" refers to the lines in the SOURCE file (%s) that this test covers
- "inputLines" and "outputLines" refer to lines in the TEST file`, functionName, filePath, filePath, filePath)

		return []PromptMessage{
			{
				Role: "user",
				Content: TextContent{
					Type: "text",
					Text: promptText,
				},
			},
		}, nil
	}

	return nil, fmt.Errorf("prompt not found: %s", name)
}
