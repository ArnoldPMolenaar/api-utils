package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	errorsutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v3"
)

func TestMachineProtected(t *testing.T) {
	t.Run("missing machine key config", func(t *testing.T) {
		t.Setenv("MACHINE_KEY", "")

		app := fiber.New()
		app.Use(MachineProtected())
		app.Get("/", func(c fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
		if err != nil {
			t.Fatalf("expected no request error, got %v", err)
		}
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", resp.StatusCode)
		}

		var body map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("expected json body, got %v", err)
		}
		if body["code"] != errorsutil.InternalServerError {
			t.Fatalf("expected code %q, got %v", errorsutil.InternalServerError, body["code"])
		}
	})

	t.Run("invalid header key", func(t *testing.T) {
		t.Setenv("MACHINE_KEY", "secret")

		app := fiber.New()
		app.Use(MachineProtected())
		app.Get("/", func(c fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("x-machine-key", "wrong")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("expected no request error, got %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("valid header key", func(t *testing.T) {
		t.Setenv("MACHINE_KEY", "secret")

		app := fiber.New()
		app.Use(MachineProtected())
		app.Get("/", func(c fiber.Ctx) error { return c.SendStatus(http.StatusNoContent) })

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("x-machine-key", "secret")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("expected no request error, got %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("expected status 204, got %d", resp.StatusCode)
		}
	})
}

