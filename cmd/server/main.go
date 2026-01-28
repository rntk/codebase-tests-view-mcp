package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"codebase-view-mcp/internal/api"
	"codebase-view-mcp/internal/files"
	"codebase-view-mcp/internal/mcp"
	"codebase-view-mcp/internal/metadata"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to run the server on")
	baseDir := flag.String("dir", ".", "Base directory to serve files from")
	metadataPath := flag.String("metadata", "metadata.json", "Path to metadata JSON file")
	flag.Parse()

	// Resolve absolute path for base directory
	absBaseDir, err := filepath.Abs(*baseDir)
	if err != nil {
		log.Fatalf("Failed to resolve base directory: %v", err)
	}

	log.Printf("Starting Codebase Test Viewer")
	log.Printf("Base directory: %s", absBaseDir)
	log.Printf("Metadata file: %s", *metadataPath)
	log.Printf("Server port: %s", *port)

	// Check if base directory exists
	if _, err := os.Stat(absBaseDir); os.IsNotExist(err) {
		log.Fatalf("Base directory does not exist: %s", absBaseDir)
	}

	// Initialize services
	fileService := files.NewService(absBaseDir)
	metaStore := metadata.NewStore(*metadataPath)
	mcpHandler := mcp.NewHandler(metaStore)

	// Initialize API handler
	apiHandler := api.NewHandler(fileService, metaStore, mcpHandler)

	// Setup routes
	router := api.SetupRoutes(apiHandler)

	// Apply middleware
	handler := api.Logging(api.CORS(router))

	// Start server
	addr := ":" + *port
	log.Printf("Server listening on http://localhost%s", addr)
	log.Printf("API available at http://localhost%s/api", addr)
	log.Printf("MCP endpoint at http://localhost%s/api/mcp", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
