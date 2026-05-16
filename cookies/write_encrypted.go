package cookies

import (
	"fmt"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"
)

// WriteEncrypted encrypts the cookie value using the provided secret key and
// writes it to the response. The cookie name is also included in the
// encryption, so that the cookie value cannot be used for other cookies.
func WriteEncrypted(c fiber.Ctx, cookie *fiber.Cookie, secretKey string) error {
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
	encryptedValue, err := utils.Encrypt(secretKey, plaintext)
	if err != nil {
		return err
	}

	// Set the cookie value to the encryptedValue.
	cookie.Value = encryptedValue

	// Write the cookie as normal.
	return Write(c, cookie)
}
