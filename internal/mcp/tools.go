package mcp

import "encoding/json"

// GetTools returns all available MCP tools
func GetTools() []Tool {
	return []Tool{
		{
			Name:        "submit-test-metadata",
			Description: "Submit metadata about tests for a source file. This tool allows LLM agents to register information about which tests cover which parts of a source file, including the line numbers for test code, input data, and expected output.",
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
							"required": ["testFile", "testName", "lineRange", "coveredLines"]
						}
					}
				},
				"required": ["sourceFile", "tests"]
			}`),
		},
	}
}
