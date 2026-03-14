package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		code       ErrorCode
		message    string
		details    map[string]interface{}
		wantStatus int
		wantCode   ErrorCode
		wantError  string
	}{
		{
			name:       "simple error without details",
			statusCode: http.StatusNotFound,
			code:       ErrFileNotFound,
			message:    "file not found",
			details:    nil,
			wantStatus: http.StatusNotFound,
			wantCode:   ErrFileNotFound,
			wantError:  "file not found",
		},
		{
			name:       "error with details",
			statusCode: http.StatusBadRequest,
			code:       ErrValidation,
			message:    "validation failed",
			details:    map[string]interface{}{"field": "line", "reason": "must be >= 1"},
			wantStatus: http.StatusBadRequest,
			wantCode:   ErrValidation,
			wantError:  "validation failed",
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			code:       ErrInternal,
			message:    "internal server error",
			details:    nil,
			wantStatus: http.StatusInternalServerError,
			wantCode:   ErrInternal,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeError(w, tt.statusCode, tt.code, tt.message, tt.details)

			if w.Code != tt.wantStatus {
				t.Errorf("status code = %d, want %d", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type = %s, want application/json", contentType)
			}

			var response ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Code != tt.wantCode {
				t.Errorf("error code = %s, want %s", response.Code, tt.wantCode)
			}

			if response.Error != tt.wantError {
				t.Errorf("error message = %s, want %s", response.Error, tt.wantError)
			}

			if tt.details != nil && response.Details == nil {
				t.Error("expected details but got nil")
			}

			if tt.details != nil {
				for key, expectedValue := range tt.details {
					if actualValue, ok := response.Details[key]; !ok {
						t.Errorf("missing detail key: %s", key)
					} else if actualValue != expectedValue {
						t.Errorf("detail[%s] = %v, want %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}
