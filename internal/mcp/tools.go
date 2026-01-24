package mcp

import "encoding/json"

// GetTools returns all available MCP tools
func GetTools() []Tool {
	return []Tool{
		{
			Name:        "submit-test-metadata",
			Description: "Submit metadata about tests for a source file. This tool allows LLM agents to register information about which tests cover which parts of a source file, including the line numbers for test code, input data, expected output, and a brief comment. Multiple submissions for the same file will be merged (tests with the same testFile+testName will be updated, new tests will be added).",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"sourceFile": {
						"type": "string",
						"description": "Path to the source file being tested"
					},
					"tests": {
						"type": "array",
						"description": "Array of test references for this source file",
						"items": {
							"type": "object",
							"properties": {
								"testFile": {
									"type": "string",
									"description": "Path to the test file"
								},
								"testName": {
									"type": "string",
									"description": "Name of the test function/method"
								},
								"comment": {
									"type": "string",
									"description": "Brief description of what the test verifies",
									"minLength": 1
								},
								"lineRange": {
									"type": "object",
									"description": "Line range of the test code",
									"properties": {
										"start": {"type": "integer", "description": "Starting line number (1-indexed)"},
										"end": {"type": "integer", "description": "Ending line number (1-indexed, inclusive)"}
									},
									"required": ["start", "end"]
								},
								"coveredLines": {
									"type": "object",
									"description": "Line range in the source file that this test covers",
									"properties": {
										"start": {"type": "integer", "description": "Starting line number (1-indexed)"},
										"end": {"type": "integer", "description": "Ending line number (1-indexed, inclusive)"}
									},
									"required": ["start", "end"]
								},
								"inputLines": {
									"type": "object",
									"description": "Line range containing the input/test data",
									"properties": {
										"start": {"type": "integer", "description": "Starting line number (1-indexed)"},
										"end": {"type": "integer", "description": "Ending line number (1-indexed, inclusive)"}
									}
								},
								"outputLines": {
									"type": "object",
									"description": "Line range containing the expected output/assertions",
									"properties": {
										"start": {"type": "integer", "description": "Starting line number (1-indexed)"},
										"end": {"type": "integer", "description": "Ending line number (1-indexed, inclusive)"}
									}
								}
							},
							"required": ["testFile", "testName", "comment", "lineRange", "coveredLines"]
						}
					}
				},
				"required": ["sourceFile", "tests"]
			}`),
		},
		{
			Name:        "suggest-missing-tests",
			Description: "Submit suggestions for missing tests for uncovered code. This tool allows LLM agents to suggest tests that should be written for functions or code sections that lack test coverage. Suggestions include a test skeleton, priority level, and reasoning.",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"sourceFile": {
						"type": "string",
						"description": "Path to the source file that needs tests"
					},
					"functionName": {
						"type": "string",
						"description": "Optional name of the function to suggest tests for"
					},
					"suggestions": {
						"type": "array",
						"description": "Array of test suggestions",
						"items": {
							"type": "object",
							"properties": {
								"targetLines": {
									"type": "object",
									"description": "Line range in the source file that needs test coverage",
									"properties": {
										"start": {"type": "integer", "description": "Starting line number (1-indexed)"},
										"end": {"type": "integer", "description": "Ending line number (1-indexed, inclusive)"}
									},
									"required": ["start", "end"]
								},
								"reason": {
									"type": "string",
									"description": "Explanation of why this test is needed"
								},
								"suggestedName": {
									"type": "string",
									"description": "Suggested name for the test function"
								},
								"testSkeleton": {
									"type": "string",
									"description": "Code skeleton for the suggested test"
								},
								"priority": {
									"type": "string",
									"enum": ["high", "medium", "low"],
									"description": "Priority level of the suggested test"
								}
							},
							"required": ["targetLines", "reason", "suggestedName", "testSkeleton", "priority"]
						}
					}
				},
				"required": ["sourceFile", "suggestions"]
			}`),
		},
	}
}
