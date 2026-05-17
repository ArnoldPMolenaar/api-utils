package cookies

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/gofiber/fiber/v3"
)

// WriteSigned encodes the cookie value using base64, signs it with a HMAC.
func WriteSigned(c fiber.Ctx, cookie *fiber.Cookie, cfg SignedConfig) error {
	// Calculate a HMAC signature of the cookie name and value using SHA256.
	mac := hmac.New(sha256.New, cfg.SecretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	// Prepend the cookie value with the HMAC signature.
	cookie.Value = string(signature) + cookie.Value

	// Call our Write() helper to base64-encode the new cookie value and write
	// the cookie.
	return Write(c, cookie, cfg.writeConfig())
}
