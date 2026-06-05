package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// ModuleFactory is a constructor function registered by each module's init().
type ModuleFactory func(*Container) Module

// moduleFactories is the global slice populated by init() in every module package.
var moduleFactories []ModuleFactory

// RegisterModule is called from each module's init() to register its factory.
// This is the only way modules self-register — no manual wiring in main.go.
func RegisterModule(f ModuleFactory) {
	moduleFactories = append(moduleFactories, f)
}

// ─────────────────────────────────────────────────────────────────────────────
// Registry
// ─────────────────────────────────────────────────────────────────────────────

type Registry struct {
	modules   []Module
	container *Container
	index     map[string]Module // fast name-based lookup
}

func NewRegistry(container *Container) *Registry {
	return &Registry{
		container: container,
		index:     make(map[string]Module),
	}
}

// Load instantiates all registered modules and validates dependencies.
//
// FIX — two-pass loading:
//
//	Pass 1: instantiate every enabled module and build the full index.
//	Pass 2: verify dependencies against the complete index.
//
// The old single-pass check only saw modules registered BEFORE the current one,
// so load order determined success — fragile and wrong.
// Now ALL modules are indexed first; dependency validation is order-independent.
func (r *Registry) Load() error {
	logger := r.container.Logger

	// ── Pass 1: instantiate all enabled modules ──────────────────────────────
	for _, factory := range moduleFactories {
		mod := factory(r.container)
		meta := mod.Meta()

		if !meta.Enabled {
			logger.Warn("module disabled — skipping",
				"module", meta.Name,
			)
			continue
		}

		// Duplicate name guard
		if _, exists := r.index[meta.Name]; exists {
			return fmt.Errorf(
				"duplicate module name %q — each module must have a unique name in manifest.json",
				meta.Name,
			)
		}

		r.modules = append(r.modules, mod)
		r.index[meta.Name] = mod

		logger.Info("module loaded", "module", meta.Name, "version", meta.Version)
	}

	// ── Pass 2: validate dependencies against the FULL index ────────────────
	// FIX: only now that all modules are indexed can we reliably check deps.
	for _, mod := range r.modules {
		meta := mod.Meta()
		for _, dep := range meta.Dependencies {
			if _, ok := r.index[dep]; !ok {

				return fmt.Errorf(
					"❌ module %q requires dependency %q — add it to modules/ and modules_gen.go, or remove it from manifest.json",
					meta.Name, dep,
				)
			}
		}
	}

	logger.Info("all modules loaded", "count", len(r.modules))
	return nil
}

// Modules returns all loaded module instances.
func (r *Registry) Modules() []Module { return r.modules }

// Get returns the module with the given name, or nil if not found.
func (r *Registry) Get(name string) Module { return r.index[name] }

// ─────────────────────────────────────────────────────────────────────────────
// Bootstrap: Init → Register → PostRegister → mount routes
// ─────────────────────────────────────────────────────────────────────────────

// internal/core/registry.go
func (r *Registry) Bootstrap(ctx context.Context, server *fiber.App, bus *EventBus, logger *slog.Logger) error {

	// ── Phase 1: Init ──────────────────────────────────────────────────
	for _, mod := range r.modules {
		meta := mod.Meta()
		if err := mod.Init(ctx); err != nil {
			return fmt.Errorf("module %q Init failed: %w", meta.Name, err)
		}
		logger.Info("module initialised", "module", meta.Name)
	}

	// ── Phase 2: Register + Routes ─────────────────────────────────────
	for _, mod := range r.modules {
		meta := mod.Meta()
		if err := mod.Register(bus); err != nil {
			return fmt.Errorf("module %q Register failed: %w", meta.Name, err)
		}

		group := server.Group(meta.BasePath)
		for _, route := range mod.Routes() {
			group.Add(route.Method, route.Path, route.Handler)
		}
		logger.Info("module registered", "module", meta.Name, "path", meta.BasePath)
	}

	// ── Phase 3: PostRegister — inject deps from manifest.json ───────────
	for _, mod := range r.modules {
		meta := mod.Meta()

		// Build resolved dep map from manifest.json dependencies
		resolved := make(map[string]Module)
		for _, depName := range meta.Dependencies {
			depMod, ok := r.index[depName]
			if !ok {
				return fmt.Errorf("module %q dependency %q not found", meta.Name, depName)
			}
			resolved[depName] = depMod
		}

		if err := mod.PostRegister(ctx, resolved); err != nil {
			return fmt.Errorf("module %q PostRegister failed: %w", meta.Name, err)
		}
		logger.Info("module post-registered", "module", meta.Name)
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Shutdown — reverse order, errors logged but non-fatal
// ─────────────────────────────────────────────────────────────────────────────

func (r *Registry) Shutdown(ctx context.Context, logger *slog.Logger) {
	for i := len(r.modules) - 1; i >= 0; i-- {
		mod := r.modules[i]
		meta := mod.Meta()
		if err := mod.Shutdown(ctx); err != nil {
			logger.Error("module shutdown error", "module", meta.Name, "error", err)
		} else {
			logger.Info("module stopped", "module", meta.Name)
		}
	}
}
