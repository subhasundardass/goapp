package accounting

import (
	"context"
	"fmt"

	"goapp/internal/core"
	"goapp/modules/accounting/handler"
	"goapp/modules/accounting/repository"
	"goapp/modules/accounting/service"
)

// Compile-time interface check.
var _ core.Module = (*Module)(nil)

type Module struct {
	*core.BaseModule
	deps     *core.Container
	handlers *handler.Handlers
	services *service.Services
	repos    *repository.Repositories
}

func init() {
	core.RegisterModule(func(c *core.Container) core.Module {
		return New(c)
	})
}

func New(c *core.Container) *Module {
	m := &Module{
		BaseModule: &core.BaseModule{},
		deps:       c,
	}
	if err := m.LoadMeta(); err != nil {
		panic("accounting module: " + err.Error())
	}
	return m
}

// ── Lifecycle ─────────────────────────────────────────────────────────────────

func (m *Module) Init(_ context.Context) error {
	m.Logger.Info("init", "module", "accounting")
	return nil
}

func (m *Module) Register(bus *core.EventBus) error {
	m.repos = repository.NewRepositories(m.deps.DB)
	m.services = service.NewServices(m.repos, m.Logger, m.deps.Store)
	m.handlers = handler.NewHandlers(m.services)

	// TODO: subscribe to events from other modules
	// bus.Subscribe("topic", func(ctx context.Context, e core.Event) error {
	// 	return nil
	// })

	return nil
}

func (m *Module) PostRegister(ctx context.Context, deps map[string]core.Module) error {
	m.Logger.Info("post-register", "module", "accounting")

	// ── Cross-module dependency injection ─────────────────────────────────────
	//
	// deps contains ONLY the modules declared in manifest.json "dependencies".
	// This is the ONLY place you should resolve other module services.
	//
	// Rules:
	//   ✅ DO    resolve other module services here
	//   ✅ DO    rewire handlers here if they need cross-module services
	//   ✅ DO    start background workers here
	//   ✅ DO    run seed data here (gate with AppEnv == "development")
	//   ❌ DON'T access Registry or Resolver directly
	//   ❌ DON'T subscribe to EventBus here (do that in Register)
	//   ❌ DON'T wire your own repos/services/handlers here (do that in Register)
	//
	// Example:
	//   invoiceSvc, ok := deps["invoice"].Service().(*invoiceservice.Services)
	//   if !ok || invoiceSvc == nil {
	//       return fmt.Errorf("invoice service unavailable")
	//   }
	//   m.handlers = handler.NewHandlers(handler.Deps{
	//       Services: m.services,
	//       Invoice:  invoiceSvc,
	//   })
	// ─────────────────────────────────────────────────────────────────────────

	// ── Seed data (development only) ─────────────────────────────────────────
	if m.deps.Config.AppEnv == "development" {
		if err := Seed(ctx, m.deps.DB); err != nil {
			return fmt.Errorf("seed failed: %w", err)
		}
	}

	// ── Background workers ────────────────────────────────────────────────────
	// go m.startWorker(ctx)

	return m.deps.EventBus.Publish(ctx, core.Event{
		Topic:   "accounting.ready",
		Payload: map[string]string{"module": "accounting"},
		Source:  "accounting",
	})
}

func (m *Module) Shutdown(ctx context.Context) error {
	m.Logger.Info("shutdown", "module", "accounting")
	return nil
}

func (m *Module) Service() any {
	if m.services == nil {
		return nil
	}
	return m.services
}

// ── Routes ────────────────────────────────────────────────────────────────────

func (m *Module) Routes() []core.Route {
	return Routes(m.handlers)
}
