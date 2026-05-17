package utils

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"
)

// ServerConfig defines runtime settings for Fiber server startup.
type ServerConfig struct {
	Address string
}

// ServerDeps defines startup dependencies to keep behavior explicit and testable.
type ServerDeps struct {
	Listen       func(a *fiber.App, addr string) error
	Shutdown     func(a *fiber.App) error
	NotifySignal func(c chan<- os.Signal, sig ...os.Signal)
}

// DefaultServerDeps returns production defaults for server startup dependencies.
func DefaultServerDeps() ServerDeps {
	return ServerDeps{
		Listen:       func(a *fiber.App, addr string) error { return a.Listen(addr) },
		Shutdown:     func(a *fiber.App) error { return a.Shutdown() },
		NotifySignal: signal.Notify,
	}
}

func (d ServerDeps) validate() error {
	if d.Listen == nil {
		return fmt.Errorf("server dependency Listen cannot be nil")
	}
	if d.Shutdown == nil {
		return fmt.Errorf("server dependency Shutdown cannot be nil")
	}
	if d.NotifySignal == nil {
		return fmt.Errorf("server dependency NotifySignal cannot be nil")
	}
	return nil
}

// ServerConfigFromEnv loads Fiber server config from environment variables.
func ServerConfigFromEnv() (ServerConfig, error) {
	addr, err := ConnectionURLBuilder("fiber")
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{Address: addr}, nil
}

func validateServerConfig(cfg ServerConfig) error {
	if cfg.Address == "" {
		return fmt.Errorf("server address cannot be empty")
	}

	return nil
}

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App) {
	cfg, err := ServerConfigFromEnv()
	if err != nil {
		log.Printf("Oops... Fiber connection URL could not be built! Reason: %v", err)
		return
	}

	StartServerWithGracefulShutdownConfig(a, cfg, DefaultServerDeps())
}

// StartServerWithGracefulShutdownConfig starts Fiber with graceful shutdown using explicit config.
func StartServerWithGracefulShutdownConfig(a *fiber.App, cfg ServerConfig, deps ServerDeps) {
	if err := validateServerConfig(cfg); err != nil {
		log.Printf("Oops... Server config is not valid! Reason: %v", err)
		return
	}
	if err := deps.validate(); err != nil {
		log.Printf("Oops... Server dependencies are not valid! Reason: %v", err)
		return
	}

	// Create channel for idle connections.
	idleConnectionClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		// Catch OS signals.
		deps.NotifySignal(sigint, os.Interrupt)
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := deps.Shutdown(a); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnectionClosed)
	}()

	// Run server.
	if err := deps.Listen(a, cfg.Address); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnectionClosed
}

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	cfg, err := ServerConfigFromEnv()
	if err != nil {
		log.Printf("Oops... Fiber connection URL could not be built! Reason: %v", err)
		return
	}

	StartServerWithConfig(a, cfg, DefaultServerDeps())
}

// StartServerWithConfig starts Fiber with explicit config.
func StartServerWithConfig(a *fiber.App, cfg ServerConfig, deps ServerDeps) {
	if err := validateServerConfig(cfg); err != nil {
		log.Printf("Oops... Server config is not valid! Reason: %v", err)
		return
	}
	if err := deps.validate(); err != nil {
		log.Printf("Oops... Server dependencies are not valid! Reason: %v", err)
		return
	}

	// Run server.
	if err := deps.Listen(a, cfg.Address); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
