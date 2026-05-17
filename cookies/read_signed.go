package cookies

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/gofiber/fiber/v3"
)

// ReadSigned reads a signed cookie value and validates its HMAC signature.
func ReadSigned(c fiber.Ctx, name string, cfg SignedConfig) (string, error) {
	// Read in the signed value from the cookie. This should be in the format
	// "{signature}{original value}".
	signedValue, err := Read(c, name)
	if err != nil {
		return "", err
	}
	if signedValue == "" {
		return "", nil
	}

	// A SHA256 HMAC signature has a fixed length of 32 bytes. To avoid an
	// out-of-range panic, check that the signed cookie value is at least this long.
	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	// Split apart the signature and original cookie value.
	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	// Recalculate the HMAC signature of the cookie name and original value.
	mac := hmac.New(sha256.New, cfg.SecretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	// Check that the recalculated signature matches the signature we received
	// in the cookie. If they match, we can be confident that the cookie name
	// and value haven't been edited by the client.
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	// Return the original cookie value.
	return value, nil
}
