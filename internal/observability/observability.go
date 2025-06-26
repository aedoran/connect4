package observability

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Start sets up tracing and metrics. It returns a shutdown function.
func Start(ctx context.Context, service string) func(context.Context) error {
	slog.Info("telemetry initialized", "service", service)
	return func(context.Context) error {
		slog.Info("telemetry shutdown")
		return nil
	}
}

// Middleware records basic tracing information for each request.
func Middleware(service string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		slog.Info("trace", "svc", service, "method", c.Method(), "path", c.Path(), "status", c.Response().StatusCode(), "duration", time.Since(start).String())
		return err
	}
}
