# Codebase Test Viewer

A test-viewer application for codebases that helps visualize test coverage, test data, and the relationship between code and tests.

## Features

- **File Browser**: Browse and navigate your codebase with a simple file explorer
- **File Viewer**: View file contents with metadata about associated tests
- **Test Visualization**: Mind-map visualization showing the relationship between code and tests
- **Test Details**: Detailed view of test information including input data and expected outputs
- **MCP Integration**: Model Context Protocol endpoint for LLM agents to submit test metadata

## Architecture

### Backend (Go)
- Pure `net/http` handlers (no external frameworks)
- JSON file persistence for test metadata
- MCP (Model Context Protocol) endpoint for LLM integration
- RESTful API endpoints

### Frontend (React + TypeScript)
- Three-panel layout with file explorer, file preview, and test panel
- SVG-based mind-map visualization
- No external UI dependencies

## Quick Start with Docker

The easiest way to run the application is with Docker Compose:

```bash
# Build and run
docker-compose up

# Access the application
open http://localhost:8080
```

The application will:
- Serve the frontend at `http://localhost:8080`
- Provide API endpoints at `http://localhost:8080/api/*`
- Store metadata in `./metadata/metadata.json`
- Browse files from the current directory

## Manual Build

### Prerequisites

- Go 1.22+ (for Go 1.25.4 compatibility)
- Node.js 18+
- npm

### Build Frontend

```bash
cd frontend
npm install
npm run build
```

### Build Backend

```bash
go build -o server ./cmd/server
```

### Run

```bash
./server -port 8080 -dir . -metadata metadata.json
```

## API Endpoints

### File Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/files?path=<path>` | List files/directories |
| GET | `/api/files/<path>` | Get file content + metadata |
| GET | `/api/files/<path>/tests` | Get related tests for a file |

### MCP Endpoint

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/mcp` | MCP JSON-RPC 2.0 endpoint |

#### Supported MCP Methods

- `initialize` - Initialize MCP session
- `tools/list` - List available tools
- `tools/call` - Execute a tool (e.g., `submit-test-metadata`)
- `prompts/list` - List available prompts
- `prompts/get` - Get a specific prompt

## Using the MCP Endpoint

### Submit Test Metadata

Use the `submit-test-metadata` tool to register test information:

```bash
curl -X POST http://localhost:8080/api/mcp \\
  -H "Content-Type: application/json" \\
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "submit-test-metadata",
      "arguments": {
        "sourceFile": "internal/files/service.go",
        "tests": [
          {
            "testFile": "internal/files/service_test.go",
            "testName": "TestListFiles",
            "lineRange": {"start": 10, "end": 25},
            "inputLines": {"start": 12, "end": 15},
            "outputLines": {"start": 20, "end": 22}
          }
        ]
      }
    }
  }'
```

### Get Prompt for LLM

```bash
curl -X POST http://localhost:8080/api/mcp \\
  -H "Content-Type: application/json" \\
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "prompts/get",
    "params": {
      "name": "codebase-tests-review",
      "arguments": {
        "functionName": "ListFiles",
        "filePath": "internal/files/service.go"
      }
    }
  }'
```

## Development

### Frontend Development

For frontend development with hot reload:

```bash
cd frontend
npm run dev
```

The dev server will proxy `/api/*` requests to `http://localhost:8080`.

### Running Backend Only

```bash
go run ./cmd/server -port 8080 -dir . -metadata metadata.json
```

### Linting

You can use the provided Dockerfiles to run linting checks for both Go and TypeScript:

```bash
# Run Golang lint
docker build -f Dockerfile.lint.golang -t lint-golang .
docker run --rm -v $(pwd):/workspace lint-golang

# Run TypeScript lint
docker build -f Dockerfile.lint.typescript -t lint-typescript .
docker run --rm -v $(pwd):/workspace lint-typescript
```

## Project Structure

```
.
├── cmd/
│   └── server/           # Server entry point
├── internal/
│   ├── api/              # HTTP handlers and routing
│   ├── files/            # File operations and models
│   ├── mcp/              # MCP protocol implementation
│   └── metadata/         # Test metadata storage
├── frontend/             # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── hooks/        # Custom React hooks
│   │   ├── api/          # API client
│   │   └── types/        # TypeScript types
│   └── public/
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yml    # Docker Compose configuration
└── go.mod               # Go module definition
```

## Configuration

### Command Line Flags

- `-port` - Server port (default: 8080)
- `-dir` - Base directory to serve files from (default: current directory)
- `-metadata` - Path to metadata JSON file (default: metadata.json)

### Environment Variables (Docker)

- `PORT` - Server port
- `METADATA_PATH` - Path to metadata file

## License

MIT
