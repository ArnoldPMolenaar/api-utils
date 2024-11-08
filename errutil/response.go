package errutil

import "github.com/gofiber/fiber/v2"

// Response creates a JSON response with a message and code.
func Response(c *fiber.Ctx, status int, code, message interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"code":    code,
		"message": message,
	})
}
