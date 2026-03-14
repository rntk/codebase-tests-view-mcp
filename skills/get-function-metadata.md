# Skill: get-function-metadata

## Description
Retrieve existing metadata for a specific function in a source file. Use this to check what tests, suggestions, and comments have already been submitted before running a full analysis.

## MCP Tools Required
- `get-function-metadata`: Queries stored metadata filtered to a single function.

## Prompt Arguments
- `sourceFile` (Required): Repo-relative path to the source file.
- `functionName` (Required): Name of the source function to look up.

## Returned Data
The tool returns a JSON object with:
- `sourceFile`: The canonicalized source file path.
- `functionName`: The queried function name.
- `tests`: Array of `TestReference` objects already submitted for this function (may be empty).
- `suggestions`: Array of `TestSuggestion` objects already submitted for this function (may be empty).

## How to Use
1. Call `get-function-metadata` with `sourceFile` and `functionName` before running `codebase-tests-review` or `test-to-source-review`.
2. If `tests` is non-empty, the function already has metadata — decide whether to update or skip.
3. If `tests` is empty, proceed with the full analysis and submit via `submit-test-metadata`.

## Example
```json
{
  "name": "get-function-metadata",
  "arguments": {
    "sourceFile": "internal/files/service.go",
    "functionName": "ReadFile"
  }
}
```
