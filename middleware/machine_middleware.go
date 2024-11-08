package middleware

import (
	"github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

// MachineProtected middleware checks if the machine key is valid.
// It reads the header x-machine-key and compares it with the value from the .env file.
// If the machine key is not valid, it returns an error response.
// Otherwise, it calls the next handler.
func MachineProtected() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		machineKey := os.Getenv("MACHINE_KEY")
		if machineKey == "" {
			log.Fatal("MACHINE_KEY is not configured in the .env file")
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
