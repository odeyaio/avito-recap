package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

type pingerFunc func(context.Context) error

func (f pingerFunc) Ping(ctx context.Context) error {
	return f(ctx)
}

func TestRegisteredHandlers(t *testing.T) {
	t.Parallel()

	router := echo.New()
	RegisterHandlers(router, pingerFunc(func(context.Context) error { return nil }))

	tests := []struct {
		name        string
		method      string
		target      string
		body        string
		wantStatus  int
		contentType string
	}{
		{
			name:        "health",
			method:      http.MethodGet,
			target:      "/api/v1/health",
			wantStatus:  http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "profiles stub",
			method:      http.MethodGet,
			target:      "/api/v1/profiles",
			wantStatus:  http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "generate recap stub",
			method:      http.MethodPost,
			target:      "/api/v1/profiles/7a9e06c2-42a2-4cb9-ae7d-3f187a51735c/recaps",
			body:        `{"year":2025}`,
			wantStatus:  http.StatusServiceUnavailable,
			contentType: "application/problem+json",
		},
		{
			name:        "get recap stub",
			method:      http.MethodGet,
			target:      "/api/v1/recaps/410942ba-9544-49db-99fb-a02cc19e5b84",
			wantStatus:  http.StatusNotFound,
			contentType: "application/problem+json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequestWithContext(t.Context(), tt.method, tt.target, strings.NewReader(tt.body))
			if tt.body != "" {
				request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.wantStatus, response.Body.String())
			}
			if got := response.Header().Get(echo.HeaderContentType); got != tt.contentType {
				t.Fatalf("Content-Type = %q, want %q", got, tt.contentType)
			}
		})
	}
}

func TestCheckHealthDatabaseUnavailable(t *testing.T) {
	t.Parallel()

	router := echo.New()
	RegisterHandlers(router, pingerFunc(func(context.Context) error {
		return errors.New("database is down")
	}))

	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/health", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusServiceUnavailable, response.Body.String())
	}
}
