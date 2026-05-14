package routes

import "github.com/gofiber/fiber/v3"

// HealthRoute registers a lightweight health check endpoint.
func HealthRoute(a *fiber.App) {
	a.Get("/health", func(c fiber.Ctx) error {
		// 204 keeps the response empty and fast for health probes.
		return c.SendStatus(fiber.StatusNoContent)
	})
}

