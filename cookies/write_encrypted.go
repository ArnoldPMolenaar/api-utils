package cookies

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// WriteEncrypted encrypts the cookie value using the secret key from config
// and writes it to the response. The cookie name is included in the encrypted
// payload so the value cannot be replayed for other cookie names.
func WriteEncrypted(c fiber.Ctx, cookie *fiber.Cookie, cfg EncryptedConfig, deps CryptoDeps) error {
	if err := deps.validateEncrypt(); err != nil {
		return err
	}

	// Prepare the plaintext input for encryption. Because we want to
	// authenticate the cookie name as well as the value, we make this plaintext
	// in the format "{cookie name}:{cookie value}". We use the : character as a
	// separator because it is an invalid character for cookie names and
	// therefore shouldn't appear in them.
	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	// Encrypt the data using aesGCM.Seal(). By passing the nonce as the first
	// parameter, the encrypted data will be appended to the nonce — meaning
	// that the returned encryptedValue variable will be in the format
	// "{nonce}{encrypted plaintext data}".
	encryptedValue, err := deps.Encrypt(cfg.SecretKey, plaintext)
	if err != nil {
		return err
	}

	// Set the cookie value to the encryptedValue.
	cookie.Value = encryptedValue

	// Write the cookie as normal.
	return Write(c, cookie, cfg.writeConfig())
}
