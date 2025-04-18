package utils

import (
	"errors"
	errorsutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a custom error handler for Fiber.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error.
	code := fiber.StatusInternalServerError
	var message string

	// Check if it's a Fiber error.
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	} else {
		// If it's a panic error, include the panic message.
		message = err.Error()
	}

	// Return the error response as JSON.
	return errorsutil.Response(
		c,
		code,
		errorsutil.InternalServerError,
		message,
	)
}
