package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	errorsutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v3"
)

func TestNotFoundRoute(t *testing.T) {
	app := fiber.New()
	NotFoundRoute(app)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("expected valid JSON response, got %v", err)
	}

	if payload["code"] != errorsutil.NotFound {
		t.Fatalf("expected code %q, got %v", errorsutil.NotFound, payload["code"])
	}

	expectedMessage := "sorry, endpoint is not found"
	if payload["message"] != expectedMessage {
		t.Fatalf("expected message %q, got %v", expectedMessage, payload["message"])
	}
}

