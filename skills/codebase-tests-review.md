# Skill: codebase-tests-review

## Description
Analyze a function and submit metadata about its tests using the `submit-test-metadata` tool.

## MCP Tools Required
- `submit-test-metadata`: Specifically used to submit the gathered metadata about the source code, tests, and line coverage.

## Prompt
**Arguments:**
- `functionName` (Required): The name of the function to analyze.
- `filePath` (Required): Path to the file containing the function.

**Prompt Content & Instructions:**
This prompt directs the assistant to first read the provided source file with line numbers (e.g. `cat -n`). It then outlines clear instructions to:
1. Examine the implementation of `functionName`.
2. Find tests associated with this function.
3. For each test found, identify specific components by exact line ranges:
    - `coveredLines`: Which part of the source function the test exercises.
    - `inputLines`: The specific input data within the test.
    - `outputLines`: The expected result in the test.
    - `lineRange`: The full span of the test itself.
    - A brief comment/description of the test's purpose.
4. Finalize the analysis by structuring a precisely formatted JSON and invoking the `submit-test-metadata` tool.

## How to Use
1. Invoke the prompt `codebase-tests-review`, ensuring you supply the arguments `functionName` and `filePath`.
2. The LLM will execute the prompt by utilizing file reading tools to examine both the target source function and its test files.
3. The LLM then structures the metadata and executes the `submit-test-metadata` tool exactly as outlined by the prompt guidelines.
