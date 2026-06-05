package core

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Application struct {
	Config   *Config
	Deps     *Container
	Registry *Registry
	Server   *fiber.App
}

func New() *Application {

	cfg := NewConfig()
	container := NewContainer(cfg)
	registry := NewRegistry(container)

	app := &Application{
		Config:   cfg,
		Deps:     container,
		Registry: registry,
		Server: fiber.New(fiber.Config{
			AppName:               cfg.AppName,
			DisableStartupMessage: true,
		}),
	}

	app.registerRoutes()
	return app
}

// central route registration
func (a *Application) registerRoutes() {
	app := a.Server

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Prototype Running 🚀")
	})
}

// Bootstrap runs Init → Register → PostRegister for all modules.
func (a *Application) Bootstrap(ctx context.Context) error {
	// Load halts immediately if a dependency is missing
	if err := a.Registry.Load(); err != nil {
		return fmt.Errorf("registry load failed: %w", err)
	}

	return a.Registry.Bootstrap(ctx, a.Server, a.Deps.EventBus, a.Deps.Logger)
}

// Shutdown gracefully shuts down all modules then stops the HTTP server.
func (a *Application) Shutdown(ctx context.Context) error {
	a.Registry.Shutdown(ctx, a.Deps.Logger)

	// Close DB connection to avoid leaking file handles / connections.
	if a.Deps.DB != nil {
		if err := a.Deps.DB.Close(); err != nil {
			a.Deps.Logger.Error("db close error", "error", err)
		}
	}

	return a.Server.ShutdownWithContext(ctx)
}

func (a *Application) Run(addr string) error {
	return a.Server.Listen(addr)
}
