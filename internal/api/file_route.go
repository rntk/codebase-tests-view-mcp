package api

import (
	"net/http"
	"strings"
)

// GetFileOrTests routes file requests to GetFile or GetTests based on the path suffix.
func (h *Handler) GetFileOrTests(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	const testsSuffix = "/tests"

	if len(path) > len(testsSuffix) && strings.HasSuffix(path, testsSuffix) {
		path = strings.TrimSuffix(path, testsSuffix)
		r.SetPathValue("path", path)
		h.GetTests(w, r)
		return
	}

	h.GetFile(w, r)
}
