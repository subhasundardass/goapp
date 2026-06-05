package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

// Route associates a URL path with a Fiber handler.
type Route struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

// Module is the contract every feature module must satisfy.
type Module interface {
	// Meta returns the parsed manifest.json config.
	// Core reads Name, BasePath, Enabled from here automatically.
	Meta() ModuleMeta

	// Init is called once before Register. Use it to set up internal state
	// that does not depend on other modules (e.g. validate config, open files).
	// Returning an error halts application startup.
	Init(ctx context.Context) error

	// Register initialises the module's services, repos, and wires dependencies.
	// The EventBus is injected here so modules can subscribe to topics published
	// by other modules (Hybrid EventBus approach).
	// An error here should halt application startup.
	Register(bus *EventBus) error

	// PostRegister is called after ALL modules have been registered.
	// Use it to publish startup events, start background workers, or perform
	// cross-module wiring that requires other modules to already be initialised.
	PostRegister(ctx context.Context, deps map[string]Module) error

	// Routes returns the URL routes this module exposes.
	Routes() []Route

	// Service optionally exposes an internal service for direct consumption by
	// other modules (alternative to event-based communication).
	// Other module usage (safe access):
	//   exampleModule := app.Registry.Get("example")
	Service() any

	// Shutdown is called during graceful shutdown. Close connections, flush
	// buffers, cancel background workers, etc.
	// ctx carries the shutdown deadline.
	Shutdown(ctx context.Context) error
}

type ModuleMeta struct {
	Name         string   `json:"name"`
	BasePath     string   `json:"base_path"`
	Enabled      bool     `json:"enabled"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"`
	Middlewares  []string `json:"middlewares"`
	Tags         []string `json:"tags"`
}

// BaseModule provides no-op implementations of every Module lifecycle method
// plus a shared logger initialised from the module name.
// Embed it in your module struct so you only override what you need.
//
//	type MyModule struct {
//	    *core.BaseModule          // ← pointer so LoadMeta can set fields
//	    deps *core.Container
//	}
type BaseModule struct {
	meta   ModuleMeta
	Logger *slog.Logger // FIX: exported so embedded modules can call m.Logger directly
}

// LoadMeta reads manifest.json from the module's own directory and initialises
// the logger with the module name as an attribute.
// Call this once at the end of New():
//
//	func New(c *core.Container) *Module {
//	    m := &Module{BaseModule: &core.BaseModule{}, deps: c}
//	    m.LoadMeta()        // sets m.meta AND m.Logger
//	    return m
//	}
func (b *BaseModule) LoadMeta() error {
	// Walk up one extra frame because LoadMeta is called from the module's New()
	// which itself is one frame above the module directory.
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("could not determine module path")
	}

	path := filepath.Join(filepath.Dir(file), "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("manifest.json not found at %s: %w", path, err)
	}

	if err := json.Unmarshal(data, &b.meta); err != nil {
		return fmt.Errorf("invalid manifest.json: %w", err)
	}

	// FIX: initialise the logger here so m.Logger is never nil
	b.Logger = slog.With("module", b.meta.Name)

	return nil
}

func (b *BaseModule) Meta() ModuleMeta { return b.meta }

// No-op lifecycle defaults — override only what your module needs.
func (b *BaseModule) Init(_ context.Context) error                              { return nil }
func (b *BaseModule) Register(_ *EventBus) error                                { return nil }
func (b *BaseModule) PostRegister(_ context.Context, _ map[string]Module) error { return nil }
func (b *BaseModule) Shutdown(_ context.Context) error                          { return nil }
func (b *BaseModule) Service() any                                              { return nil }
