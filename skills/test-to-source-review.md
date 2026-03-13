# Skill: test-to-source-review

## Description
Analyze a test and trace it back to the source code it covers, then submit metadata using the `submit-test-metadata` tool.

## MCP Tools Required
- `submit-test-metadata`: Specifically used to submit the gathered metadata mapping the given test back to the targeted source function.

## Prompt
**Arguments:**
- `testName` (Required): The name of the test function to analyze.
- `testFilePath` (Required): Path to the test file containing the test.

**Prompt Content & Instructions:**
This prompt instructs the assistant to start by reading the specified test file (`testFilePath`) with line numbers prefixed. Following that, it must:
1. Find and examine the behavior of `testName`.
2. Identify the source file and function being tested and read that source file with line numbers.
3. Map the test behavior to the implementation and identify line ranges for:
    - The source function being exercised.
    - `coveredLines`: The lines in the actual source code the test is covering.
    - `inputLines`: Where test data/inputs are declared.
    - `outputLines`: Where the expected result/assertions are.
    - `lineRange`: The full span of the test.
    - A descriptive comment on what the test validates.
4. Lastly, format all retrieved line numbers and metadata into JSON and pass it to the `submit-test-metadata` tool.

## How to Use
1. Invoke the prompt `test-to-source-review`, providing the inputs for `testName` and `testFilePath`.
2. The LLM will read the test file to understand what it's testing, trace it over to the appropriate source code, and read that source file.
3. The LLM correlates the test data/assertions to the exact source lines and then finalizes the workflow by calling the `submit-test-metadata` tool with the specific JSON payload.
