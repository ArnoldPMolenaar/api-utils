package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	errorsutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v3"
)

func TestNewValidatorAndValidatorErrors(t *testing.T) {
	type payload struct {
		ID string `validate:"required,uuid"`
	}

	validate := NewValidator()

	if err := validate.Struct(payload{ID: "550e8400-e29b-41d4-a716-446655440000"}); err != nil {
		t.Fatalf("expected valid uuid to pass validation, got %v", err)
	}

	err := validate.Struct(payload{ID: "not-a-uuid"})
	if err == nil {
		t.Fatal("expected invalid uuid validation error")
	}

	fields := ValidatorErrors(err)
	if _, ok := fields["ID"]; !ok {
		t.Fatalf("expected ID field validation error, got %v", fields)
	}

	if got := ValidatorErrors(nil); len(got) != 0 {
		t.Fatalf("expected empty map for nil error, got %v", got)
	}
	if got := ValidatorErrors(errors.New("other")); len(got) != 0 {
		t.Fatalf("expected empty map for non-validation error, got %v", got)
	}
}

func TestErrorHandler(t *testing.T) {
	app := fiber.New()

	app.Get("/fiber", func(c fiber.Ctx) error {
		return ErrorHandler(c, fiber.NewError(fiber.StatusNotFound, "missing"))
	})
	app.Get("/generic", func(c fiber.Ctx) error {
		return ErrorHandler(c, errors.New("boom"))
	})

	respFiber, err := app.Test(httptest.NewRequest(http.MethodGet, "/fiber", nil))
	if err != nil {
		t.Fatalf("expected no request error, got %v", err)
	}
	if respFiber.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected 404, got %d", respFiber.StatusCode)
	}

	var bodyFiber map[string]any
	if err := json.NewDecoder(respFiber.Body).Decode(&bodyFiber); err != nil {
		t.Fatalf("expected valid json body, got %v", err)
	}
	if bodyFiber["code"] != errorsutil.InternalServerError {
		t.Fatalf("expected code %q, got %v", errorsutil.InternalServerError, bodyFiber["code"])
	}
	if bodyFiber["message"] != "missing" {
		t.Fatalf("expected message %q, got %v", "missing", bodyFiber["message"])
	}

	respGeneric, err := app.Test(httptest.NewRequest(http.MethodGet, "/generic", nil))
	if err != nil {
		t.Fatalf("expected no request error, got %v", err)
	}
	if respGeneric.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", respGeneric.StatusCode)
	}

	var bodyGeneric map[string]any
	if err := json.NewDecoder(respGeneric.Body).Decode(&bodyGeneric); err != nil {
		t.Fatalf("expected valid json body, got %v", err)
	}
	if bodyGeneric["message"] != "boom" {
		t.Fatalf("expected message %q, got %v", "boom", bodyGeneric["message"])
	}
}

func TestStartServer(t *testing.T) {
	origBuilder := connectionURLBuilderFn
	origListen := listenFn
	defer func() {
		connectionURLBuilderFn = origBuilder
		listenFn = origListen
	}()

	called := false
	var gotAddr string
	connectionURLBuilderFn = func(name string) (string, error) {
		if name != "fiber" {
			t.Fatalf("expected builder name fiber, got %q", name)
		}
		return ":9999", nil
	}
	listenFn = func(a *fiber.App, addr string) error {
		called = true
		gotAddr = addr
		return nil
	}

	StartServer(fiber.New())

	if !called {
		t.Fatal("expected listen to be called")
	}
	if gotAddr != ":9999" {
		t.Fatalf("expected listen addr :9999, got %q", gotAddr)
	}
}

func TestStartServerBuilderError(t *testing.T) {
	origBuilder := connectionURLBuilderFn
	origListen := listenFn
	defer func() {
		connectionURLBuilderFn = origBuilder
		listenFn = origListen
	}()

	connectionURLBuilderFn = func(name string) (string, error) {
		return "", errors.New("builder failed")
	}
	listenCalled := false
	listenFn = func(a *fiber.App, addr string) error {
		listenCalled = true
		return nil
	}

	StartServer(fiber.New())

	if listenCalled {
		t.Fatal("expected listen not to be called when URL builder fails")
	}
}

func TestStartServerWithGracefulShutdown(t *testing.T) {
	origBuilder := connectionURLBuilderFn
	origListen := listenFn
	origShutdown := shutdownFn
	origNotify := notifySignalFn
	defer func() {
		connectionURLBuilderFn = origBuilder
		listenFn = origListen
		shutdownFn = origShutdown
		notifySignalFn = origNotify
	}()

	connectionURLBuilderFn = func(name string) (string, error) {
		return ":9999", nil
	}
	listenCalled := false
	shutdownCalled := false
	listenFn = func(a *fiber.App, addr string) error {
		listenCalled = true
		return nil
	}
	shutdownFn = func(a *fiber.App) error {
		shutdownCalled = true
		return nil
	}
	notifySignalFn = func(c chan<- os.Signal, sig ...os.Signal) {
		go func() {
			c <- os.Interrupt
		}()
	}

	StartServerWithGracefulShutdown(fiber.New())

	if !listenCalled {
		t.Fatal("expected listen to be called")
	}
	if !shutdownCalled {
		t.Fatal("expected shutdown to be called")
	}
}

