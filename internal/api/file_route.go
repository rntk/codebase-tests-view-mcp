package api

import (
	"net/http"
	"strings"
)

// GetFileOrTests routes file requests to GetFile, GetTests, GetSuggestions, or GetSources based on the path suffix.
func (h *Handler) GetFileOrTests(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	const testsSuffix = "/tests"
	const suggestionsSuffix = "/suggestions"
	const sourcesSuffix = "/sources"

	if len(path) > len(testsSuffix) && strings.HasSuffix(path, testsSuffix) {
		path = strings.TrimSuffix(path, testsSuffix)
		r.SetPathValue("path", path)
		h.GetTests(w, r)
		return
	}

	if len(path) > len(suggestionsSuffix) && strings.HasSuffix(path, suggestionsSuffix) {
		path = strings.TrimSuffix(path, suggestionsSuffix)
		r.SetPathValue("path", path)
		h.GetSuggestions(w, r)
		return
	}

	if len(path) > len(sourcesSuffix) && strings.HasSuffix(path, sourcesSuffix) {
		path = strings.TrimSuffix(path, sourcesSuffix)
		r.SetPathValue("path", path)
		h.GetSources(w, r)
		return
	}

	h.GetFile(w, r)
}
