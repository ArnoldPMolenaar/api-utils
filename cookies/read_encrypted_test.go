package cookies

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"
)

const testAESKey = "0123456789abcdef0123456789abcdef"

func TestReadEncryptedReturnsEmptyWhenCookieMissing(t *testing.T) {
	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		value, err := ReadEncrypted(c, "session", EncryptedConfig{SecretKey: testAESKey}, DefaultCryptoDeps())
		if err != nil {
			return err
		}
		return c.SendString(value)
	})

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestReadEncryptedReturnsErrInvalidValueOnDecryptFailure(t *testing.T) {
	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		_, err := ReadEncrypted(c, "session", EncryptedConfig{SecretKey: testAESKey}, DefaultCryptoDeps())
		if err == ErrInvalidValue {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: base64.URLEncoding.EncodeToString([]byte("not-encrypted")),
	})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

func TestReadEncryptedReturnsErrInvalidValueOnMalformedPlaintext(t *testing.T) {
	encrypted, err := utils.Encrypt(testAESKey, "missing-separator")
	if err != nil {
		t.Fatalf("expected encryption to succeed, got %v", err)
	}

	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		_, readErr := ReadEncrypted(c, "session", EncryptedConfig{SecretKey: testAESKey}, DefaultCryptoDeps())
		if readErr == ErrInvalidValue {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		return readErr
	})

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: base64.URLEncoding.EncodeToString([]byte(encrypted)),
	})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

func TestReadEncryptedReturnsErrInvalidValueOnCookieNameMismatch(t *testing.T) {
	encrypted, err := utils.Encrypt(testAESKey, "other:payload")
	if err != nil {
		t.Fatalf("expected encryption to succeed, got %v", err)
	}

	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		_, readErr := ReadEncrypted(c, "session", EncryptedConfig{SecretKey: testAESKey}, DefaultCryptoDeps())
		if readErr == ErrInvalidValue {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		return readErr
	})

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: base64.URLEncoding.EncodeToString([]byte(encrypted)),
	})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

func TestWriteEncryptedAndReadEncryptedRoundTrip(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		return WriteEncrypted(c, &fiber.Cookie{
			Name:  "session",
			Value: "payload",
			Path:  "/",
		}, EncryptedConfig{SecretKey: testAESKey, Write: DefaultWriteConfig()}, DefaultCryptoDeps())
	})
	app.Get("/get", func(c fiber.Ctx) error {
		value, err := ReadEncrypted(c, "session", EncryptedConfig{SecretKey: testAESKey}, DefaultCryptoDeps())
		if err != nil {
			return err
		}
		return c.SendString(value)
	})

	setReq := httptest.NewRequest(http.MethodGet, "/set", nil)
	setResp, err := app.Test(setReq)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	cookies := setResp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one set-cookie header, got %d", len(cookies))
	}

	getReq := httptest.NewRequest(http.MethodGet, "/get", nil)
	getReq.AddCookie(cookies[0])
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResp.StatusCode)
	}
}

