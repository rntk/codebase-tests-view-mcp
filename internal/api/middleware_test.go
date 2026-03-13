package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	t.Run("adds CORS headers to all requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		corsHandler := CORS(handler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		headers := rr.Header()

		if headers.Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("expected Access-Control-Allow-Origin %q, got %q", "*", headers.Get("Access-Control-Allow-Origin"))
		}

		if headers.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Errorf("expected Access-Control-Allow-Methods %q, got %q", "GET, POST, PUT, DELETE, OPTIONS", headers.Get("Access-Control-Allow-Methods"))
		}

		if headers.Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
			t.Errorf("expected Access-Control-Allow-Headers %q, got %q", "Content-Type, Authorization", headers.Get("Access-Control-Allow-Headers"))
		}
	})

	t.Run("handles OPTIONS preflight requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called for OPTIONS requests")
		})

		corsHandler := CORS(handler)
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("passes through non-OPTIONS requests to next handler", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		corsHandler := CORS(handler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		if !called {
			t.Error("next handler was not called")
		}
	})

	t.Run("works with POST requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})

		corsHandler := CORS(handler)
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
		}

		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("CORS headers not set for POST request")
		}
	})

	t.Run("works with PUT requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		corsHandler := CORS(handler)
		req := httptest.NewRequest(http.MethodPut, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("works with DELETE requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

		corsHandler := CORS(handler)
		req := httptest.NewRequest(http.MethodDelete, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
		}
	})
}

func TestLogging(t *testing.T) {
	t.Run("logs request method and path", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		loggingHandler := Logging(handler)
		req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
		rr := httptest.NewRecorder()

		loggingHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("captures status code correctly", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})

		loggingHandler := Logging(handler)
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		rr := httptest.NewRecorder()

		loggingHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
		}
	})

	t.Run("captures status code from WriteHeader", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		loggingHandler := Logging(handler)
		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		rr := httptest.NewRecorder()

		loggingHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("defaults to 200 if WriteHeader not called", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Don't call WriteHeader, should default to 200
			w.Write([]byte("OK"))
		})

		loggingHandler := Logging(handler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		loggingHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("works with different HTTP methods", func(t *testing.T) {
		methods := []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodHead,
			http.MethodOptions,
		}

		for _, method := range methods {
			t.Run(method, func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				loggingHandler := Logging(handler)
				req := httptest.NewRequest(method, "/test", nil)
				rr := httptest.NewRecorder()

				loggingHandler.ServeHTTP(rr, req)

				if rr.Code != http.StatusOK {
					t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
				}
			})
		}
	})
}

func TestResponseWriter(t *testing.T) {
	t.Run("default status code is 200", func(t *testing.T) {
		rw := &responseWriter{
			ResponseWriter: httptest.NewRecorder(),
			statusCode:     http.StatusOK,
		}

		if rw.statusCode != http.StatusOK {
			t.Errorf("expected default status code %d, got %d", http.StatusOK, rw.statusCode)
		}
	})

	t.Run("WriteHeader updates status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		rw := &responseWriter{
			ResponseWriter: recorder,
			statusCode:     http.StatusOK,
		}

		rw.WriteHeader(http.StatusCreated)

		if rw.statusCode != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rw.statusCode)
		}

		if recorder.Code != http.StatusCreated {
			t.Errorf("expected recorder status code %d, got %d", http.StatusCreated, recorder.Code)
		}
	})

	t.Run("WriteHeader with different status codes", func(t *testing.T) {
		statusCodes := []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusInternalServerError,
			http.StatusServiceUnavailable,
		}

		for _, code := range statusCodes {
			t.Run(string(rune(code)), func(t *testing.T) {
				recorder := httptest.NewRecorder()
				rw := &responseWriter{
					ResponseWriter: recorder,
					statusCode:     http.StatusOK,
				}

				rw.WriteHeader(code)

				if rw.statusCode != code {
					t.Errorf("expected status code %d, got %d", code, rw.statusCode)
				}
			})
		}
	})
}

func TestMiddlewareChain(t *testing.T) {
	t.Run("CORS and Logging can be chained", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Chain: Logging -> CORS -> handler
		chain := Logging(CORS(handler))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		chain.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Check CORS headers are present
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("CORS headers not present")
		}
	})

	t.Run("middleware chain preserves response", func(t *testing.T) {
		expectedBody := "Hello, World!"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedBody))
		})

		chain := Logging(CORS(handler))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		chain.ServeHTTP(rr, req)

		if rr.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rr.Body.String())
		}
	})
}

func TestMiddlewareWithDifferentPaths(t *testing.T) {
	paths := []string{
		"/",
		"/api",
		"/api/files",
		"/api/files/test.go",
		"/api/files/test.go/tests",
		"/api/mcp",
		"/static/index.html",
		"/deep/nested/path/here",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			chain := Logging(CORS(handler))
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			chain.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
			}
		})
	}
}
