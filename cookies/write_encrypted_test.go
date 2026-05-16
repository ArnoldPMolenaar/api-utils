package cookies

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestWriteEncryptedReturnsErrorForInvalidSecretLength(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		err := WriteEncrypted(c, &fiber.Cookie{
			Name:  "session",
			Value: "payload",
			Path:  "/",
		}, "short-key")
		if err != nil {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		return c.SendStatus(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/set", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

func TestWriteEncryptedPropagatesErrValueTooLong(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		err := WriteEncrypted(c, &fiber.Cookie{
			Name:  "session",
			Value: strings.Repeat("x", 5000),
			Path:  "/",
		}, testAESKey)
		if err == ErrValueTooLong {
			return c.SendStatus(http.StatusRequestEntityTooLarge)
		}
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/set", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, resp.StatusCode)
	}
}

