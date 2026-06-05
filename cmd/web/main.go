package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goapp/internal/core"
)

func main() {

	// ── Load .env FIRST (before any config or logger init) ────────────
	if err := core.LoadEnv(".env"); err != nil {
		fmt.Println("warn: .env not found, using environment variables")
	}

	// ── Bootstrap application ──────────────────────────────────────────────────
	app := core.New()

	server := app.Server
	deps := app.Deps
	cfg := deps.Config
	logger := deps.Logger

	server.Static("/static", "./web/static")

	// Middleware (must run before Bootstrap so routes inherit middleware)
	core.NewMiddlewareRegistry(server, cfg, logger)

	// Run the three-phase lifecycle: Init → Register → PostRegister
	startCtx, startCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer startCancel()

	if err := app.Bootstrap(startCtx); err != nil {
		logger.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
	core.PrintBanner(app.Server, cfg, len(app.Registry.Modules()))

	// ── Graceful shutdown ──────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("shutdown signal received")

		shutCtx, shutCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutCancel()

		if err := app.Shutdown(shutCtx); err != nil {
			logger.Error("shutdown error", "error", err)
		}
	}()

	// ── Start HTTP server ──────────────────────────────────────────────────────
	addr := ":" + cfg.AppPort
	logger.Info("server starting", "addr", addr)

	if err := app.Run(addr); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)

	}
}
