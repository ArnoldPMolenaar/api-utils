package middleware

import (
	"os"

	"github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

// MachineProtected middleware checks if the machine key is valid.
// It reads the header x-machine-key and compares it with the value from the .env file.
// If the machine key is not valid, it returns an error response.
// Otherwise, it calls the next handler.
func MachineProtected() func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		machineKey := os.Getenv("MACHINE_KEY")
		if machineKey == "" {
			log.Error("MACHINE_KEY is not configured in the .env file")
			return errors.Response(
				c,
				fiber.StatusInternalServerError,
				errors.InternalServerError,
				"Machine key is not configured.",
			)
		}

		headerKey := c.Get("x-machine-key")
		if headerKey == "" || headerKey != machineKey {
			return errors.Response(
				c,
				fiber.StatusUnauthorized,
				errors.Unauthorized,
				"Machine key is invalid.",
			)
		}

		return c.Next()
	}
}
