package routes

import "github.com/gofiber/fiber/v3"

// ReadinessCheck is a custom function that verifies whether an API dependency is ready.
// Return nil when healthy, or an error when the dependency is unavailable.
type ReadinessCheck func() error

// LivezRoute registers a lightweight liveness endpoint.
func LivezRoute(a *fiber.App) {
	a.Get("/livez", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
}

// ReadyzRoute registers a readiness endpoint with pluggable checks.
func ReadyzRoute(a *fiber.App, checks ...ReadinessCheck) {
	a.Get("/readyz", func(c fiber.Ctx) error {
		for _, check := range checks {
			if check == nil {
				continue
			}

			if err := check(); err != nil {
				return c.SendStatus(fiber.StatusServiceUnavailable)
			}
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}

// KubernetesProbeRoutes registers liveness and readiness endpoints for Kubernetes probes.
func KubernetesProbeRoutes(a *fiber.App, checks ...ReadinessCheck) {
	LivezRoute(a)
	ReadyzRoute(a, checks...)
}

