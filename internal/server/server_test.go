package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nimishgj/aws-radar/internal/config"
)

func TestHealthEndpoint(t *testing.T) {
	s := New(&config.ServerConfig{
		Port:        0,
		MetricsPath: "/metrics",
		HealthPath:  "/health",
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.httpServer.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != "OK" {
		t.Fatalf("expected OK body, got %q", rec.Body.String())
	}
}

func TestRootEndpoint(t *testing.T) {
	s := New(&config.ServerConfig{
		Port:        0,
		MetricsPath: "/metrics",
		HealthPath:  "/health",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.httpServer.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "/metrics") || !strings.Contains(body, "/health") {
		t.Fatalf("expected root page to include metrics and health links")
	}
}
