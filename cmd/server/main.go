package main

import (
	"flag"
	"log/slog"
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
		slog.Error("Failed to resolve base directory", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting Codebase Test Viewer")
	slog.Info("Base directory", "dir", absBaseDir)
	slog.Info("Metadata file", "path", *metadataPath)
	slog.Info("Server port", "port", *port)

	// Check if base directory exists
	if _, err := os.Stat(absBaseDir); os.IsNotExist(err) {
		slog.Error("Base directory does not exist", "dir", absBaseDir)
		os.Exit(1)
	}

	// Initialize services
	fileService := files.NewService(absBaseDir)
	metaStore := metadata.NewStore(*metadataPath)
	mcpHandler := mcp.NewHandler(metaStore, fileService)

	// Initialize API handler
	apiHandler := api.NewHandler(fileService, metaStore, mcpHandler)

	// Setup routes
	router := api.SetupRoutes(apiHandler)

	// Apply middleware
	handler := api.Logging(api.CORS(router))

	// Start server
	addr := ":" + *port
	slog.Info("Server starting", "address", addr, "api", "/api", "mcp", "/api/mcp")

	if err := http.ListenAndServe(addr, handler); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
