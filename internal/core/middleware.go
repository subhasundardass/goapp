package core

import (
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareRegistry struct {
	middlewares []fiber.Handler
}

func NewMiddlewareRegistry(app *fiber.App, cfg *Config, logger *slog.Logger) *MiddlewareRegistry {

	r := &MiddlewareRegistry{
		middlewares: make([]fiber.Handler, 0),
	}

	r.Register(RecoveryMiddleware(logger))
	r.Register(LoggingMiddleware(logger))
	r.Register(CORSMiddleware(cfg))
	r.Register(RequestTimingMiddleware(logger))

	// IMPORTANT: attach to fiber app
	for _, m := range r.middlewares {
		app.Use(m)
	}

	return r
}

func (r *MiddlewareRegistry) Register(m fiber.Handler) {
	r.middlewares = append(r.middlewares, m)
}

// RecoveryMiddleware recovers from panics, logs the stack trace, and returns 500.
// response even after a panic.
func RecoveryMiddleware(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (retErr error) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered",
					"error", rec,
					"stack", string(debug.Stack()),
				)
				retErr = c.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{"error": "internal server error"})
			}
		}()

		return c.Next()
	}
}

// LoggingMiddleware logs each incoming request's method and path.
func LoggingMiddleware(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Info("request",
			"method", c.Method(),
			"path", c.Path(),
		)

		return c.Next()
	}
}

// CORSMiddleware sets CORS headers using the configured app URL as the allowed origin.
func CORSMiddleware(cfg *Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", cfg.CORSAllowedOrigins)
		c.Set("Access-Control-Allow-Methods", cfg.CORSAllowedMethods)
		c.Set("Access-Control-Allow-Headers", cfg.CORSAllowedHeaders)

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	}
}

// RequestTimingMiddleware logs the method, path, status code, and duration of each request.
func RequestTimingMiddleware(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		logger.Info("response",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration", time.Since(start).String(),
		)

		return err
	}
}
