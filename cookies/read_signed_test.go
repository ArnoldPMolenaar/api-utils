package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

var (
	signedKeyA = []byte("signed-secret-a")
	signedKeyB = []byte("signed-secret-b")
)

func TestReadSignedReturnsEmptyWhenCookieMissing(t *testing.T) {
	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		value, err := ReadSigned(c, "session", SignedConfig{SecretKey: signedKeyA})
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

func TestReadSignedReturnsErrInvalidValueWhenSignatureTooShort(t *testing.T) {
	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		_, err := ReadSigned(c, "session", SignedConfig{SecretKey: signedKeyA})
		if err == ErrInvalidValue {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: base64.URLEncoding.EncodeToString([]byte("short")),
	})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

func TestReadSignedReturnsErrInvalidValueWhenSignatureMismatch(t *testing.T) {
	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		_, err := ReadSigned(c, "session", SignedConfig{SecretKey: signedKeyA})
		if err == ErrInvalidValue {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		return err
	})

	tamperedValue := string(make([]byte, sha256.Size)) + "payload"
	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: base64.URLEncoding.EncodeToString([]byte(tamperedValue)),
	})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

func TestWriteSignedAndReadSignedRoundTrip(t *testing.T) {
	app := fiber.New()
	app.Get("/set", func(c fiber.Ctx) error {
		return WriteSigned(c, &fiber.Cookie{
			Name:  "session",
			Value: "payload",
			Path:  "/",
		}, SignedConfig{SecretKey: signedKeyA, Write: DefaultWriteConfig()})
	})
	app.Get("/get", func(c fiber.Ctx) error {
		value, err := ReadSigned(c, "session", SignedConfig{SecretKey: signedKeyA})
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

func TestReadSignedFailsWithDifferentSecret(t *testing.T) {
	mac := hmac.New(sha256.New, signedKeyA)
	mac.Write([]byte("session"))
	mac.Write([]byte("payload"))
	signedValue := string(mac.Sum(nil)) + "payload"

	app := fiber.New()
	app.Get("/get", func(c fiber.Ctx) error {
		_, err := ReadSigned(c, "session", SignedConfig{SecretKey: signedKeyB})
		if err == ErrInvalidValue {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: base64.URLEncoding.EncodeToString([]byte(signedValue)),
	})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("expected no test error, got %v", err)
	}

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}
}

