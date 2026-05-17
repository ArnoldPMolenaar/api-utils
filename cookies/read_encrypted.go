package cookies

import (
	"github.com/gofiber/fiber/v3"

	"strings"
)

// ReadEncrypted reads an encrypted cookie value from the request and decrypts
// it using the secret key from config. The cookie name is authenticated in the
// encrypted payload so values cannot be replayed across different cookie names.
// If decryption fails or the authenticated name does not match, it returns
// ErrInvalidValue.
func ReadEncrypted(c fiber.Ctx, name string, cfg EncryptedConfig, deps CryptoDeps) (string, error) {
	if err := deps.validateDecrypt(); err != nil {
		return "", err
	}

	// Read the encrypted value from the cookie as normal.
	encryptedValue, err := Read(c, name)
	if err != nil {
		return "", err
	}
	if encryptedValue == "" {
		return "", nil
	}

	// Use aesGCM.Open() to decrypt and authenticate the data. If this fails,
	// return a ErrInvalidValue error.
	plaintext, err := deps.Decrypt(cfg.SecretKey, encryptedValue)
	if err != nil {
		return "", ErrInvalidValue
	}

	// The plaintext value is in the format "{cookie name}:{cookie value}". We
	// use strings.Cut() to split it on the first ":" character.
	expectedName, value, ok := strings.Cut(plaintext, ":")
	if !ok {
		return "", ErrInvalidValue
	}

	// Check that the cookie name is the expected one and hasn't been changed.
	if expectedName != name {
		return "", ErrInvalidValue
	}

	// Return the plaintext cookie value.
	return value, nil
}
