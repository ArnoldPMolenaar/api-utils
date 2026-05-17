package cookies

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestWriteEncodesCookieValue(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		return Write(c, &fiber.Cookie{
			Name:  "session",
			Value: "plain-value",
			Path:  "/",
		}, DefaultWriteConfig())
	})

	req := httptest.NewRequest(http.MethodGet, "/set", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one set-cookie header, got %d", len(cookies))
	}

	expected := base64.URLEncoding.EncodeToString([]byte("plain-value"))
	if cookies[0].Value != expected {
		t.Fatalf("expected encoded value %q, got %q", expected, cookies[0].Value)
	}
}

func TestWriteReturnsErrValueTooLong(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		err := Write(c, &fiber.Cookie{
			Name:  "session",
			Value: strings.Repeat("a", 5000),
			Path:  "/",
		}, DefaultWriteConfig())
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

func TestWriteBoundaryAtSerializedCookieLimit(t *testing.T) {
	baseCookie := fiber.Cookie{
		Name:     "session",
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	}

	maxRawLen := -1
	for i := 0; i < 6000; i++ {
		cookieAtLen := baseCookie
		cookieAtLen.Value = base64.URLEncoding.EncodeToString([]byte(strings.Repeat("a", i)))
		if serializedCookieSize(&cookieAtLen) <= 4096 {
			maxRawLen = i
			continue
		}
		break
	}

	if maxRawLen < 0 {
		t.Fatal("expected to find a valid value length within the cookie limit")
	}

	fits := baseCookie
	fits.Value = base64.URLEncoding.EncodeToString([]byte(strings.Repeat("a", maxRawLen)))
	if serializedCookieSize(&fits) > 4096 {
		t.Fatalf("expected boundary value to fit limit, got %d", serializedCookieSize(&fits))
	}

	over := baseCookie
	over.Value = base64.URLEncoding.EncodeToString([]byte(strings.Repeat("a", maxRawLen+1)))
	if serializedCookieSize(&over) <= 4096 {
		t.Fatalf("expected boundary+1 value to overflow limit, got %d", serializedCookieSize(&over))
	}

	app := fiber.New()
	app.Get("/fit", func(c fiber.Ctx) error {
		return Write(c, &fiber.Cookie{
			Name:     "session",
			Value:    strings.Repeat("a", maxRawLen),
			Path:     "/",
			Secure:   true,
			HTTPOnly: true,
			SameSite: fiber.CookieSameSiteLaxMode,
		}, DefaultWriteConfig())
	})
	app.Get("/overflow", func(c fiber.Ctx) error {
		err := Write(c, &fiber.Cookie{
			Name:     "session",
			Value:    strings.Repeat("a", maxRawLen+1),
			Path:     "/",
			Secure:   true,
			HTTPOnly: true,
			SameSite: fiber.CookieSameSiteLaxMode,
		}, DefaultWriteConfig())
		if err == ErrValueTooLong {
			return c.SendStatus(http.StatusRequestEntityTooLarge)
		}
		return err
	})

	fitReq := httptest.NewRequest(http.MethodGet, "/fit", nil)
	fitResp, err := app.Test(fitReq)
	if err != nil {
		t.Fatalf("expected no test error for fit route, got %v", err)
	}
	if fitResp.StatusCode != http.StatusOK {
		t.Fatalf("expected fit status %d, got %d", http.StatusOK, fitResp.StatusCode)
	}

	overflowReq := httptest.NewRequest(http.MethodGet, "/overflow", nil)
	overflowResp, err := app.Test(overflowReq)
	if err != nil {
		t.Fatalf("expected no test error for overflow route, got %v", err)
	}
	if overflowResp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected overflow status %d, got %d", http.StatusRequestEntityTooLarge, overflowResp.StatusCode)
	}
}

