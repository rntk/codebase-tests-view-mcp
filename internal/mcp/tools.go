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
						"description": "Repo-relative path to the source file being tested, relative to the server's configured -dir root"
					},
					"tests": {
						"type": "array",
						"description": "Array of test references for this source file",
						"items": {
							"type": "object",
							"properties": {
								"testFile": {
									"type": "string",
									"description": "Repo-relative path to the test file, relative to the server's configured -dir root"
								},
								"functionName": {
									"type": "string",
									"description": "Name of the source function being tested"
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
							"required": ["testFile", "functionName", "testName", "comment", "lineRange", "coveredLines"]
						}
					}
				},
				"required": ["sourceFile", "tests"]
			}`),
		},
		{
			Name:        "get-function-metadata",
			Description: "Retrieve existing metadata for a specific function in a source file. Returns all tests associated with the given function. Use this before submitting new metadata to check what already exists.",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"sourceFile": {
						"type": "string",
						"description": "Repo-relative path to the source file, relative to the server's configured -dir root"
					},
					"functionName": {
						"type": "string",
						"description": "Name of the source function to retrieve metadata for"
					}
				},
				"required": ["sourceFile", "functionName"]
			}`),
		},
	}
}
