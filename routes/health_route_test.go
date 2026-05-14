package routes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestHealthRoute(t *testing.T) {
	app := fiber.New()
	HealthRoute(app)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("expected response body to be readable, got %v", err)
	}

	if len(body) != 0 {
		t.Fatalf("expected empty body, got %q", string(body))
	}
}

