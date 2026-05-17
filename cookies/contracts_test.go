package cookies

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestWriteWithConfigHonorsCustomMaxSerializedSize(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		err := Write(c, &fiber.Cookie{Name: "s", Value: "payload", Path: "/"}, WriteConfig{MaxSerializedSize: 10})
		if errors.Is(err, ErrValueTooLong) {
			return c.SendStatus(http.StatusRequestEntityTooLarge)
		}
		return err
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/set", nil))
	if err != nil {
		t.Fatalf("expected no request error, got %v", err)
	}
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, resp.StatusCode)
	}
}

func TestWriteEncryptedWithDepsReturnsConfigErrorForNilEncrypt(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		err := WriteEncrypted(
			c,
			&fiber.Cookie{Name: "session", Value: "payload", Path: "/"},
			EncryptedConfig{SecretKey: testAESKey, Write: DefaultWriteConfig()},
			CryptoDeps{},
		)
		if errors.Is(err, ErrInvalidConfiguration) {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return err
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/set", nil))
	if err != nil {
		t.Fatalf("expected no request error, got %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestReadEncryptedWithDepsReturnsConfigErrorForNilDecrypt(t *testing.T) {
	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		_, err := ReadEncrypted(c, "session", EncryptedConfig{SecretKey: testAESKey}, CryptoDeps{})
		if errors.Is(err, ErrInvalidConfiguration) {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return err
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/get", nil))
	if err != nil {
		t.Fatalf("expected no request error, got %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestWriteSignedWithConfigAndReadSignedWithConfigRoundTrip(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		return WriteSigned(c, &fiber.Cookie{Name: "session", Value: "payload", Path: "/"}, SignedConfig{SecretKey: signedKeyA, Write: DefaultWriteConfig()})
	})
	app.Get("/get", func(c fiber.Ctx) error {
		value, err := ReadSigned(c, "session", SignedConfig{SecretKey: signedKeyA})
		if err != nil {
			return err
		}
		return c.SendString(value)
	})

	setResp, err := app.Test(httptest.NewRequest(http.MethodGet, "/set", nil))
	if err != nil {
		t.Fatalf("expected no request error, got %v", err)
	}
	cookies := setResp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}

	getReq := httptest.NewRequest(http.MethodGet, "/get", nil)
	getReq.AddCookie(cookies[0])
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("expected no request error, got %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResp.StatusCode)
	}
}

func TestNewCookieUsesPolicy(t *testing.T) {
	policy := CookiePolicy{
		Path:     "/api",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteStrictMode,
		MaxAge:   3600,
	}

	c := NewCookie("session", "payload", policy)
	if c.Name != "session" || c.Value != "payload" {
		t.Fatalf("unexpected name/value: %+v", c)
	}
	if c.Path != "/api" || !c.Secure || !c.HTTPOnly || c.SameSite != fiber.CookieSameSiteStrictMode || c.MaxAge != 3600 {
		t.Fatalf("unexpected cookie policy fields: %+v", c)
	}
}

func TestNewCookieDefaultsPath(t *testing.T) {
	c := NewCookie("session", "payload", CookiePolicy{})
	if c.Path != "/" {
		t.Fatalf("expected default path '/', got %q", c.Path)
	}
}


