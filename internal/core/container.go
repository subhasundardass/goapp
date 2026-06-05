package core

import (
	"log/slog"
	"os"

	"goapp/ent"
	core "goapp/internal/config"
)

type Container struct {
	Config   *Config
	DB       *ent.Client
	Logger   *slog.Logger
	EventBus *EventBus
	Registry *Registry
	Store    *Store
}

func NewContainer(cfg *Config) *Container {
	// Production → structured JSON for log aggregators (Datadog, CloudWatch, etc.)
	// Development → coloured pretty output for the terminal
	var logger *slog.Logger
	if cfg.AppEnv == "production" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	} else {
		logger = slog.New(NewPrettyHandler(slog.LevelDebug))
	}

	// DB INIT
	db := core.NewDBClient(cfg.DBDSN, cfg.AutoMigration)

	// EventBus — shared across all modules
	bus := NewEventBus(logger)

	return &Container{
		Config:   cfg,
		DB:       db,
		Logger:   logger,
		EventBus: bus,
		Store:    NewStore(),
	}
}
