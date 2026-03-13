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
					Description: "Repo-relative path to the file containing the function, relative to the server's configured -dir root",
					Required:    true,
				},
			},
		},
		{
			Name:        "test-to-source-review",
			Description: "Analyze a test and trace it back to the source code it covers, then submit metadata using the submit-test-metadata tool",
			Arguments: []PromptArgument{
				{
					Name:        "testName",
					Description: "The name of the test function to analyze",
					Required:    true,
				},
				{
					Name:        "testFilePath",
					Description: "Repo-relative path to the test file containing the test, relative to the server's configured -dir root",
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

Use repo-relative paths only. Do not submit container-specific absolute paths such as /app/... or /workspace/....

**FIRST**: Read the source file and any related test files using `+"`cat -n %s`"+` (or equivalent) so that every line is prefixed with its line number. This ensures accurate line references in your analysis.

1. Examine the function's implementation.
2. If the function has associated tests, use the **submit-test-metadata** tool to submit metadata.
3. For each test, identify:
   * Which part of the source function the test exercises (line range in the source file).
   * The specific input data used in the test (line numbers in test file).
   * The expected result in the test (line numbers in test file).
   * A brief comment/description of what the test verifies.
     - Every test must include a non-empty comment; do not omit it.
   * The source function name (use the function name provided in this prompt, not the test name).

When reporting the analysis, include **only** the following information for each test:

- File name 
- Function name (source)
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
      "testFile": "relative/path/to/test_file.go",
      "functionName": "%s",
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
- "functionName" must be the source function name from this prompt (not the test name)
- "coveredLines" refers to the lines in the SOURCE file (%s) that this test covers
- "inputLines" and "outputLines" refer to lines in the TEST file
- "sourceFile" and "testFile" must stay repo-relative to the server root`, functionName, filePath, filePath, filePath, functionName, filePath)

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

	if name == "test-to-source-review" {
		testName := args["testName"]
		testFilePath := args["testFilePath"]

		if testName == "" {
			return nil, fmt.Errorf("testName argument is required")
		}
		if testFilePath == "" {
			return nil, fmt.Errorf("testFilePath argument is required")
		}

		promptText := fmt.Sprintf(`Please analyze the test **%s** in file **%s**.

Use repo-relative paths only. Do not submit container-specific absolute paths such as /app/... or /workspace/....

**FIRST**: Read the test file using `+"`cat -n %s`"+` (or equivalent) so that every line is prefixed with its line number. This ensures accurate line references.

1. Find the test **%s** and examine what it does.
2. Identify which source file and function it is testing — read that source file too, using `+"`cat -n`"+` with line numbers.
3. Use the **submit-test-metadata** tool to submit metadata about the relationship.

For the test, identify:
* Which source function it exercises (the function name in the source file).
* Which lines in the source file the test covers (coveredLines).
* The specific input data used in the test (inputLines — line numbers in the TEST file).
* The expected result in the test (outputLines — line numbers in the TEST file).
* A brief comment/description of what the test verifies (required, non-empty).
* The line range of the test itself (lineRange — line numbers in the TEST file).

Use the **submit-test-metadata** tool with the following structure:

{
  "sourceFile": "relative/path/to/source_file.go",
  "tests": [
    {
      "testFile": "%s",
      "functionName": "SourceFunctionName",
      "testName": "%s",
      "comment": "Brief description of what the test verifies",
      "lineRange": {"start": 10, "end": 25},
      "coveredLines": {"start": 45, "end": 60},
      "inputLines": {"start": 12, "end": 15},
      "outputLines": {"start": 20, "end": 22}
    }
  ]
}

**IMPORTANT**:
- "lineRange" refers to the lines in the TEST file (%s) where the test code is located
- "functionName" must be the source function name (not the test name)
- "coveredLines" refers to the lines in the SOURCE file that this test covers
- "inputLines" and "outputLines" refer to lines in the TEST file
- "sourceFile" and "testFile" must stay repo-relative to the server root`, testName, testFilePath, testFilePath, testName, testFilePath, testName, testFilePath)

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
