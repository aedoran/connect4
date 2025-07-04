package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"mem0-go/internal/observability"

	"mem0-go/internal/config"
	"mem0-go/internal/docs"
	"mem0-go/internal/graphql"
	"mem0-go/internal/inmem"
	"mem0-go/internal/memory"
	"mem0-go/internal/rest"
)

func setupApp(logger *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error("unhandled error", "err", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		},
	})

	app.Use(recover.New())
	// tracing and metrics middleware
	app.Use(observability.Middleware("api"))
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		logger.Info("request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration", time.Since(start).String(),
		)
		return err
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	repo := inmem.NewRepo()
	vec := inmem.NewVector()
	g := inmem.NewGraph()
	svc := memory.NewService(repo, vec, g)
	graphql.Register(app, svc)
	rest.Register(app, svc)
	docs.Register(app)

	return app
}

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	shutdown := observability.Start(context.Background(), "api")

	app := setupApp(logger)
	addr := ":" + cfg.HTTPPort

	srvErr := make(chan error, 1)
	go func() {
		logger.Info("starting server", "addr", addr)
		srvErr <- app.Listen(addr)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-srvErr:
		if err != nil && err != fiber.ErrServerClosed {
			logger.Error("server error", "err", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
	}

	_ = shutdown(context.Background())
}
