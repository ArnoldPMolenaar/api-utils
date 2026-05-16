package cookies

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestReadReturnsEmptyWhenCookieIsMissing(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		value, err := Read(c, "session")
		if err != nil {
			return err
		}
		return c.SendString(value)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestReadDecodesCookieValue(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		value, err := Read(c, "session")
		if err != nil {
			return err
		}
		return c.SendString(value)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: base64.URLEncoding.EncodeToString([]byte("abc123")),
	})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestReadReturnsErrInvalidValueForInvalidBase64(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		_, err := Read(c, "session")
		if err != nil {
			if errors.Is(err, ErrInvalidValue) {
				return c.SendStatus(http.StatusUnprocessableEntity)
			}
			return err
		}
		return c.SendStatus(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "not-base64"})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

