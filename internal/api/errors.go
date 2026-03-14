package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorCode represents application-specific error codes
type ErrorCode string

const (
	ErrInvalidPath       ErrorCode = "ERR_INVALID_PATH"
	ErrFileNotFound      ErrorCode = "ERR_FILE_NOT_FOUND"
	ErrInvalidRequest    ErrorCode = "ERR_INVALID_REQUEST"
	ErrValidation        ErrorCode = "ERR_VALIDATION"
	ErrInternal          ErrorCode = "ERR_INTERNAL"
	ErrMethodNotAllowed  ErrorCode = "ERR_METHOD_NOT_ALLOWED"
	ErrEncodeFailed      ErrorCode = "ERR_ENCODE_FAILED"
	ErrNotFound          ErrorCode = "ERR_NOT_FOUND"
	ErrConflict          ErrorCode = "ERR_CONFLICT"
)

// ErrorResponse represents a JSON error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    ErrorCode              `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, statusCode int, code ErrorCode, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}
