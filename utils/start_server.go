package utils

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"
)

var (
	connectionURLBuilderFn = ConnectionURLBuilder
	listenFn               = func(a *fiber.App, addr string) error { return a.Listen(addr) }
	shutdownFn             = func(a *fiber.App) error { return a.Shutdown() }
	notifySignalFn         = signal.Notify
)

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App) {
	// Build Fiber connection URL.
	fiberConnURL, err := connectionURLBuilderFn("fiber")
	if err != nil {
		log.Printf("Oops... Fiber connection URL could not be built! Reason: %v", err)
		return
	}

	// Create channel for idle connections.
	idleConnectionClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		// Catch OS signals.
		notifySignalFn(sigint, os.Interrupt)
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := shutdownFn(a); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnectionClosed)
	}()

	// Run server.
	if err := listenFn(a, fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnectionClosed
}

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Build Fiber connection URL.
	fiberConnectionURL, err := connectionURLBuilderFn("fiber")
	if err != nil {
		log.Printf("Oops... Fiber connection URL could not be built! Reason: %v", err)
		return
	}

	// Run server.
	if err := listenFn(a, fiberConnectionURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
