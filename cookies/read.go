package cookies

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v3"
)

// Read decodes the cookie value using base64 and returns it as a string.
func Read(c fiber.Ctx, name string) (string, error) {
	// Read the cookie as normal.
	cookie := c.Cookies(name)
	if cookie == "" {
		return "", nil
	}

	// Decode the base64-encoded cookie value. If the cookie didn't contain a
	// valid base64-encoded value, this operation will fail, and we return an
	// ErrInvalidValue error.
	value, err := base64.URLEncoding.DecodeString(cookie)
	if err != nil {
		return "", ErrInvalidValue
	}

	// Return the decoded cookie value.
	return string(value), nil
}
