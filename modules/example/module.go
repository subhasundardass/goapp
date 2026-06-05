package example

import (
	"context"

	"goapp/internal/core"
	"goapp/modules/example/handler"
	"goapp/modules/example/repository"
	"goapp/modules/example/service"
)

// Compile-time interface check.
var _ core.Module = (*Module)(nil)

type Module struct {
	*core.BaseModule
	deps    *core.Container
	handler *handler.Handler
	service *service.Service
	repo    *repository.Repo
}

// build object
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
		panic("example module: " + err.Error())
	}
	return m
}

// ── Lifecycle ─────────────────────────────────────────────────────────────────

func (m *Module) Init(_ context.Context) error {
	m.Logger.Info("init", "module", "example")
	return nil
}

func (m *Module) Register(bus *core.EventBus) error {
	m.repo = repository.NewRepo(m.deps.DB)
	m.service = service.NewService(m.repo, m.Logger, m.deps.Store)
	m.handler = handler.NewHandler(m.service)

	// TODO: subscribe to events from other modules
	// bus.Subscribe("topic", func(ctx context.Context, e core.Event) error {
	// 	return nil
	// })

	return nil
}

func (m *Module) PostRegister(ctx context.Context, deps map[string]core.Module) error {
	m.Logger.Info("post-register", "module", "example")

	// ── Cross-module dependency injection ─────────────────────────────────────
	//
	// deps contains ONLY the modules declared in manifest.json "dependencies".
	// This is the ONLY place you should resolve other module services.
	//
	// Rules:
	//   ✅ DO    resolve other module services here
	//   ✅ DO    rewire handler here if it needs cross-module services
	//   ✅ DO    start background workers here
	//   ✅ DO    run seed data here (gate with AppEnv == "development")
	//   ❌ DON'T access Registry or Resolver directly
	//   ❌ DON'T subscribe to EventBus here (do that in Register)
	//   ❌ DON'T wire your own repo/service/handler here (do that in Register)
	//
	// Example — resolving a single dependency:
	//
	//   invoiceSvc, ok := deps["invoice"].Service().(*invoiceservice.Service)
	//   if !ok || invoiceSvc == nil {
	//       return fmt.Errorf("invoice service unavailable")
	//   }
	//
	// Example — resolving multiple dependencies:
	//
	//   invoiceSvc, ok := deps["invoice"].Service().(*invoiceservice.Service)
	//   if !ok || invoiceSvc == nil {
	//       return fmt.Errorf("invoice service unavailable")
	//   }
	//
	//   orderSvc, ok := deps["order"].Service().(*orderservice.Service)
	//   if !ok || orderSvc == nil {
	//       return fmt.Errorf("order service unavailable")
	//   }
	//
	// Example — rewire handler after resolving cross-module deps:
	//
	//   m.handler = handler.NewHandler(handler.Deps{
	//       Service:  m.service,      // own service — always required
	//       Invoice:  invoiceSvc,     // cross-module — only if needed
	//       Order:    orderSvc,       // cross-module — only if needed
	//   })
	//
	// If this module has NO cross-module dependencies:
	//   — leave this section empty
	//   — handler is already wired in Register(), no need to rewire here
	// ─────────────────────────────────────────────────────────────────────────

	// ── Seed data (development only) ─────────────────────────────────────────
	// if m.deps.Config.AppEnv == "development" {
	//     if err := Seed(ctx, m.deps.DB); err != nil {
	//         return fmt.Errorf("seed failed: %w", err)
	//     }
	// }

	// ── Background workers ────────────────────────────────────────────────────
	// Start long-running goroutines here. Use ctx for cancellation.
	// go m.startWorker(ctx)

	// ── Announce module is ready ──────────────────────────────────────────────
	// Always publish the ready event last — after all wiring is complete.
	// Other modules can subscribe to "example.ready" to react.
	return m.deps.EventBus.Publish(ctx, core.Event{
		Topic:   "example.ready",
		Payload: map[string]string{"module": "example"},
		Source:  "example",
	})
}

func (m *Module) Shutdown(ctx context.Context) error {
	m.Logger.Info("shutdown", "module", "example")

	// ── Cleanup ───────────────────────────────────────────────────────────────
	//
	// Called automatically on SIGINT / SIGTERM in reverse registration order.
	// ctx carries the shutdown deadline (default 15s) — respect it.
	//
	// Rules:
	//   ✅ DO    stop background workers / goroutines
	//   ✅ DO    flush pending queues or buffers
	//   ✅ DO    close module-specific connections (redis, external APIs)
	//   ✅ DO    check ctx.Done() in long-running cleanup
	//   ❌ DON'T close the shared DB — app.go handles that
	//   ❌ DON'T start new work here
	//   ❌ DON'T ignore ctx — shutdown has a deadline
	//
	// Example — stop a background worker:
	//
	//   if m.cancelWorker != nil {
	//       m.cancelWorker()
	//   }
	//
	// Example — flush a buffer with deadline awareness:
	//
	//   done := make(chan struct{})
	//   go func() {
	//       m.buffer.Flush()
	//       close(done)
	//   }()
	//   select {
	//   case <-done:
	//   case <-ctx.Done():
	//       m.Logger.Warn("shutdown timeout, buffer not fully flushed", "module", "example")
	//   }
	// ─────────────────────────────────────────────────────────────────────────

	return nil
}

func (m *Module) Service() any {
	if m.service == nil {
		return nil
	}
	return m.service
}

// ── Routes ────────────────────────────────────────────────────────────────────

func (m *Module) Routes() []core.Route {
	return []core.Route{
		{Method: "GET", Path: "/", Handler: m.handler.Index},
		{Method: "GET", Path: "/list", Handler: m.handler.List},
		{Method: "POST", Path: "/", Handler: m.handler.Create},
		{Method: "GET", Path: "/:id", Handler: m.handler.Show},
		{Method: "PUT", Path: "/:id", Handler: m.handler.Update},
		{Method: "DELETE", Path: "/:id", Handler: m.handler.Delete},

		//---
		{Method: "POST", Path: "/compute", Handler: m.handler.Compute},
	}
}
