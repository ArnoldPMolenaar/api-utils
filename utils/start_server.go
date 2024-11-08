package utils

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
)

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App) {
	// Create channel for idle connections.
	idleConnectionClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		// Catch OS signals.
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnectionClosed)
	}()

	// Build Fiber connection URL.
	fiberConnURL, _ := ConnectionURLBuilder("fiber")

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnectionClosed
}

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Build Fiber connection URL.
	fiberConnectionURL, _ := ConnectionURLBuilder("fiber")

	// Run server.
	if err := a.Listen(fiberConnectionURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
