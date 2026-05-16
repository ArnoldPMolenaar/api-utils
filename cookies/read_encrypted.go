package cookies

import (
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"

	"strings"
)

// ReadEncrypted reads an encrypted cookie value from the request and
// decrypts it using the provided secret key. The cookie name is also
// included in the encryption, so that the cookie value cannot be used for
// other cookies. If the cookie name has been changed, or if the
// decryption fails, an ErrInvalidValue error is returned.
func ReadEncrypted(c fiber.Ctx, name, secretKey string) (string, error) {
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
	plaintext, err := utils.Decrypt(secretKey, encryptedValue)
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
