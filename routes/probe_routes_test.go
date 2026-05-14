package routes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestLivezRoute(t *testing.T) {
	app := fiber.New()
	LivezRoute(app)

	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestReadyzRouteWithoutChecks(t *testing.T) {
	app := fiber.New()
	ReadyzRoute(app)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestReadyzRouteWithFailingCheck(t *testing.T) {
	app := fiber.New()
	ReadyzRoute(app,
		func() error { return nil },
		func() error { return errors.New("postgres unavailable") },
	)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
}

func TestReadyzRouteWithPassingChecks(t *testing.T) {
	app := fiber.New()
	ReadyzRoute(app,
		func() error { return nil },
		func() error { return nil },
	)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

