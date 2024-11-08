package routes

import (
	"github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
)

// NotFoundRoute func for describe 404 Error route.
func NotFoundRoute(a *fiber.App) {
	// Register new special route.
	a.Use(
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return errors.Response(
				c,
				fiber.StatusNotFound,
				errors.NotFound,
				"sorry, endpoint is not found",
			)
		},
	)
}
