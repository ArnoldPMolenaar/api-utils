package cookies

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestWriteSignedPropagatesErrValueTooLong(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		err := WriteSigned(c, &fiber.Cookie{
			Name:  "session",
			Value: strings.Repeat("x", 5000),
			Path:  "/",
		}, signedKeyA)
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

