package api

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed ../../frontend/build/*
var staticFiles embed.FS

// GetStaticFS returns the embedded static file system
func GetStaticFS() (fs.FS, error) {
	return fs.Sub(staticFiles, "frontend/build")
}

// ServeStaticFiles returns an http.Handler for serving static files
func ServeStaticFiles() http.Handler {
	staticFS, err := GetStaticFS()
	if err != nil {
		// Return a fallback handler if static files are not embedded
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Codebase Test Viewer</title>
</head>
<body>
    <h1>Codebase Test Viewer</h1>
    <p>Frontend build not found. Please build the frontend first:</p>
    <pre>cd frontend && npm install && npm run build</pre>
    <p>API is running at <code>/api/*</code></p>
</body>
</html>`))
		})
	}

	return http.FileServer(http.FS(staticFS))
}
