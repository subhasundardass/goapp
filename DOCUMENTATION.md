# GoApp — Developer Documentation

> A modular, lifecycle-aware Go web framework built on Fiber, Ent, Templ, and Datastar.

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Tech Stack](#2-tech-stack)
3. [Project Structure](#3-project-structure)
4. [Installation & Setup](#4-installation--setup)
5. [Environment Configuration](#5-environment-configuration)
6. [Application Lifecycle](#6-application-lifecycle)
7. [Creating a New Module](#7-creating-a-new-module)
8. [Routing](#8-routing)
9. [EventBus — Inter-Module Communication](#9-eventbus--inter-module-communication)
10. [Render & Response](#10-render--response)
11. [Templ Components](#11-templ-components)
12. [Datastar — Reactivity](#12-datastar--reactivity)
13. [Database — Ent ORM](#13-database--ent-orm)
14. [Database Migration](#14-database-migration)
15. [Database Seeding](#15-database-seeding)
16. [Middleware](#16-middleware)
17. [Configuration System](#17-configuration-system)
18. [Logging](#18-logging)
19. [Error Handling](#19-error-handling)
20. [Makefile Commands](#20-makefile-commands)
21. [Adding a Full CRUD Module — End to End Example](#21-adding-a-full-crud-module--end-to-end-example)

---

## 1. Project Overview

GoApp is a **modular monolith** framework. Each feature lives in its own self-contained module under `modules/`. Modules register themselves automatically via Go's `init()` mechanism — the core never needs to know about a specific module by name.

Key design principles:

- **No magic** — every module explicitly declares its routes, lifecycle hooks, and dependencies.
- **Lifecycle-aware** — every module goes through `Init → Register → PostRegister → Shutdown`.
- **EventBus** — modules communicate via a hybrid sync/async event bus instead of direct imports.
- **Type-safe templates** — Templ compiles HTML components to Go at build time.
- **Zero JS framework** — Datastar handles all reactivity over SSE from the server.

---

## 2. Tech Stack

| Layer       | Technology                                   |
| ----------- | -------------------------------------------- |
| HTTP Server | [Fiber v2](https://github.com/gofiber/fiber) |
| ORM         | [Ent](https://entgo.io)                      |
| Templates   | [Templ](https://templ.guide)                 |
| Reactivity  | [Datastar](https://data-star.dev)            |
| Database    | SQLite (dev) / PostgreSQL (prod)             |
| Language    | Go 1.22+                                     |

---

## 3. Project Structure

```
goapp/
├── cmd/
│   └── web/
│       └── main.go              # Entry point
│
├── internal/
│   ├── core/
│   │   ├── app.go               # Application bootstrap + shutdown
│   │   ├── module.go            # Module interface (contract)
│   │   ├── base_module.go       # No-op embed helper
│   │   ├── registry.go          # Module factory registry + lifecycle runner
│   │   ├── container.go         # Dependency injection container
│   │   ├── config.go            # Config struct + env loader
│   │   ├── env.go               # .env file parser
│   │   ├── eventbus.go          # Hybrid sync/async event bus
│   │   ├── middleware.go        # Global middleware (CORS, logging, recovery)
│   │   ├── render.go            # Templ + SSE render helpers
│   │   └── response.go          # JSON response helpers
│   │
│   └── config/
│       └── db.go                # Database client factory
│
├── modules/
│   └── example/                 # Example feature module
│       ├── module.go            # Module wiring + lifecycle
│       ├── handler/
│       │   └── example_handler.go
│       ├── service/
│       │   └── example_service.go
│       ├── repository/
│       │   └── example_repo.go
│       ├── pages/               # Full-page templ templates
│       └── components/          # Reusable templ fragments
│
├── ent/
│   ├── schema/                  # Ent schema definitions
│   │   ├── example.go
│   │   └── user.go
│   ├── generate.go              # go:generate directive
│   └── ...                      # Generated ent code (do not edit)
│
├── web/
│   ├── static/                  # Static assets (datastar.js, css)
│   ├── app.css
│   └── app.js
│
├── .env                         # Local environment variables
├── Makefile
├── go.mod
└── go.sum
```

---

## 4. Installation & Setup

### Prerequisites

- Go 1.22+
- GCC (required for SQLite via `go-sqlite3`)
- [Templ CLI](https://templ.guide/quick-start/installation)
- [Air](https://github.com/air-verse/air) (optional, for live reload)

### Steps

```bash
# 1. Clone the repository
git clone https://github.com/your-org/goapp.git
cd goapp

# 2. Copy environment file
cp .env.example .env

# 3. Install Go dependencies
go mod download

# 4. Generate Ent ORM code (only needed after schema changes)
go generate ./ent/...

# 5. Generate Templ components (only needed after .templ file changes)
templ generate

# 6. Run the application
go run ./cmd/web

# Server starts at http://localhost:8000
```

### Live Reload (Development)

```bash
# Run templ watcher and air together
make watch
```

---

## 5. Environment Configuration

All configuration is driven by the `.env` file. Copy `.env.example` and edit as needed.

```dotenv
# =========================
# APPLICATION
# =========================
APP_NAME=GoERP
APP_ENV=development        # development | production
APP_PORT=8000
APP_URL=http://localhost:8000

# =========================
# DATABASE
# =========================
DB_DRIVER=sqlite
DB_DSN=file:goapp.db?_fk=1

# PostgreSQL (uncomment to switch)
# DB_DRIVER=postgres
# DB_DSN=postgres://user:password@localhost:5432/erp?sslmode=disable

# =========================
# LOGGING
# =========================
LOG_LEVEL=info             # debug | info | warn | error
LOG_FORMAT=json            # json | text

# =========================
# SECURITY
# =========================
JWT_SECRET=change_this_secret
JWT_EXPIRY=24h

# =========================
# MODULES
# =========================
MODULE_AUTO_LOAD=true

# =========================
# CORS
# =========================
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

# =========================
# DEVELOPMENT
# =========================
DEBUG=true
ENABLE_MIGRATION_AUTO=true
```

### Config fields reference

| Variable                | Config Field     | Default                 | Description                   |
| ----------------------- | ---------------- | ----------------------- | ----------------------------- |
| `APP_NAME`              | `AppName`        | `GoERP`                 | Application name              |
| `APP_ENV`               | `AppEnv`         | `development`           | Environment mode              |
| `APP_PORT`              | `AppPort`        | `8000`                  | HTTP listen port              |
| `APP_URL`               | `AppURL`         | `http://localhost:8000` | Public URL                    |
| `DB_DRIVER`             | `DBDriver`       | `sqlite`                | Database driver               |
| `DB_DSN`                | `DBDSN`          | `file:goapp.db?_fk=1`   | Connection string             |
| `LOG_LEVEL`             | `LogLevel`       | `info`                  | Minimum log level             |
| `LOG_FORMAT`            | `LogFormat`      | `json`                  | Log output format             |
| `JWT_SECRET`            | `JWTSecret`      | —                       | Signing key                   |
| `JWT_EXPIRY`            | `JWTExpiry`      | `24h`                   | Token lifetime                |
| `MODULE_AUTO_LOAD`      | `ModuleAutoLoad` | `true`                  | Master module switch          |
| `DEBUG`                 | `Debug`          | `false`                 | Enables debug logging         |
| `ENABLE_MIGRATION_AUTO` | `AutoMigration`  | `true`                  | Run schema migration on start |

---

## 6. Application Lifecycle

The application boots in strictly ordered phases:

```
main()
  │
  ├── LoadEnv(".env")               # Parse .env into os environment
  │
  ├── core.New()                    # Build Config → Container → Registry → Fiber
  │
  ├── NewMiddlewareRegistry()       # Attach global middleware to Fiber
  │
  └── app.Bootstrap(ctx)
        │
        ├── Phase 1 — Init()        # Each module validates config, opens files
        │                           # No cross-module dependencies here
        │
        ├── Phase 2 — Register()    # Each module wires DI, subscribes to EventBus
        │             + Routes      # Routes are mounted after Register returns
        │
        └── Phase 3 — PostRegister() # All modules ready; publish events, start workers

  [Server running]

  SIGINT / SIGTERM
        │
        └── app.Shutdown(ctx)
              │
              ├── Module Shutdown() # Called in reverse order
              ├── DB.Close()        # Close database connection
              └── Server.Shutdown() # Drain in-flight requests
```

### Rules per phase

| Phase          | Do                              | Don't                         |
| -------------- | ------------------------------- | ----------------------------- |
| `Init`         | Validate config, open files     | Call other modules            |
| `Register`     | Wire DI, subscribe to events    | Publish events                |
| `PostRegister` | Publish events, start workers   | Expect modules not yet loaded |
| `Shutdown`     | Stop workers, close connections | Start new work                |

---

## 7. Creating a New Module

### Step 1 — Create the folder structure

```bash
mkdir -p modules/customer/{handler,service,repository,pages,components}
```

### Step 2 — Create `module.go`

```go
// modules/customer/module.go
package customer

import (
    "context"

    "goapp/internal/core"
    "goapp/modules/customer/handler"
    "goapp/modules/customer/repository"
    "goapp/modules/customer/service"
)

// Compile-time check: Module must satisfy core.Module.
var _ core.Module = (*Module)(nil)

type Module struct {
    core.BaseModule          // provides no-op defaults for unused lifecycle methods
    deps    *core.Container
    handler *handler.Handler
    service *service.Service
    repo    *repository.Repo
}

// init() runs automatically when the package is imported.
// It registers the module factory with the core registry.
func init() {
    core.RegisterModule(func(c *core.Container) core.Module {
        return New(c)
    })
}

func New(c *core.Container) *Module {
    return &Module{deps: c}
}

func (m *Module) Name() string     { return "customer" }
func (m *Module) BasePath() string { return "/customer" }

func (m *Module) Init(_ context.Context) error {
    // Validate config, check required env vars, etc.
    return nil
}

func (m *Module) Register(bus *core.EventBus) error {
    db := m.deps.DB

    m.repo    = repository.NewRepo(db)
    m.service = service.NewService(m.repo, m.deps.Logger)
    m.handler = handler.NewHandler(m.service)

    // Subscribe to events from other modules
    bus.Subscribe("user.created", func(ctx context.Context, e core.Event) error {
        // react to user creation
        return nil
    })

    return nil
}

func (m *Module) PostRegister(ctx context.Context) error {
    // Announce readiness, start background jobs, etc.
    return m.deps.EventBus.Publish(ctx, core.Event{
        Topic:   "customer.ready",
        Source:  m.Name(),
        Payload: nil,
    })
}

func (m *Module) Shutdown(_ context.Context) error {
    // Clean up resources
    return nil
}

func (m *Module) Routes() []core.Route {
    return []core.Route{
        {Method: "GET",  Path: "/",       Handler: m.handler.ListPage},
        {Method: "GET",  Path: "/:id",    Handler: m.handler.DetailPage},
        {Method: "POST", Path: "/",       Handler: m.handler.Create},
        {Method: "PUT",  Path: "/:id",    Handler: m.handler.Update},
        {Method: "DELETE", Path: "/:id",  Handler: m.handler.Delete},
    }
}

func (m *Module) Service() any { return m.service }
```

### Step 3 — Blank-import the module in `main.go`

```go
// cmd/web/main.go
import (
    _ "goapp/modules/customer"   // ← add this line
    _ "goapp/modules/example"
)
```

This triggers `init()` which registers the module factory. Without this import the module never loads.

### Step 4 — Add module enable flag to `.env`

```dotenv
MODULE_CUSTOMER_ENABLED=true
```

Check it in `Init()` if you want the module to be disableable at runtime:

```go
func (m *Module) Init(_ context.Context) error {
    if !m.deps.Config.ModuleCustomerEnabled {
        return fmt.Errorf("customer module is disabled")
    }
    return nil
}
```

---

## 8. Routing

Routes are declared in each module's `Routes()` method and mounted automatically under the module's `BasePath()` during the `Register` phase.

### Route struct

```go
type Route struct {
    Method  string        // "GET", "POST", "PUT", "DELETE", "PATCH"
    Path    string        // relative to BasePath, e.g. "/list"
    Handler fiber.Handler
}
```

### Full URL = BasePath + Route.Path

```
BasePath:  /customer
Path:      /:id
Full URL:  GET /customer/:id
```

### Route parameters

```go
func (h *Handler) Detail(c *fiber.Ctx) error {
    id := c.Params("id")
    // ...
}
```

### Query strings

```go
func (h *Handler) List(c *fiber.Ctx) error {
    page  := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 20)
    // ...
}
```

### Request body

```go
type CreateRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (h *Handler) Create(c *fiber.Ctx) error {
    var req CreateRequest
    if err := c.BodyParser(&req); err != nil {
        return core.BadRequest(c, err)
    }
    // ...
}
```

---

## 9. EventBus — Inter-Module Communication

The EventBus enables modules to communicate without importing each other directly, keeping them loosely coupled.

### Subscribe (synchronous)

Blocks the publisher until this handler returns. Use for critical side-effects.

```go
bus.Subscribe("order.placed", func(ctx context.Context, e core.Event) error {
    order := e.Payload.(Order)
    return inventoryService.Reserve(ctx, order.Items)
})
```

### SubscribeAsync (asynchronous)

Runs in a new goroutine. Errors are logged, not returned to publisher. Use for notifications, cache busting, analytics.

```go
bus.SubscribeAsync("order.placed", func(ctx context.Context, e core.Event) error {
    order := e.Payload.(Order)
    return emailService.SendConfirmation(ctx, order)
})
```

### Publish

```go
err := m.deps.EventBus.Publish(ctx, core.Event{
    Topic:   "order.placed",
    Source:  m.Name(),
    Payload: order,
})
```

### Event struct

```go
type Event struct {
    Topic   string  // e.g. "user.created", "order.placed"
    Payload any     // cast to expected type in the handler
    Source  string  // name of the publishing module
}
```

### Topic naming convention

```
<module>.<action>

user.created
user.deleted
order.placed
order.cancelled
customer.ready
invoice.generated
```

### When to use EventBus vs direct service call

| Scenario                                       | Use                           |
| ---------------------------------------------- | ----------------------------- |
| Module A needs a result from Module B          | `Registry.Get("b").Service()` |
| Module A wants to notify Module B of something | EventBus                      |
| Multiple modules react to the same event       | EventBus                      |
| Real-time, non-blocking notification           | `SubscribeAsync`              |
| Must succeed before continuing                 | `Subscribe` (sync)            |

---

## 10. Render & Response

### Render — HTML (Templ)

```go
// Full page render — standard HTTP request
func (h *Handler) ListPage(c *fiber.Ctx) error {
    items, err := h.service.List()
    if err != nil {
        return core.InternalError(c, err)
    }
    return core.Render(c, pages.ListPage(items))
}

// With custom status code
return core.Render(c, pages.NotFoundPage(), core.WithStatus(404))
```

### RenderSSE — HTML fragment via Datastar

```go
// Push a fragment to the browser over SSE
func (h *Handler) Stream(c *fiber.Ctx) error {
    return core.StreamSSE(c, func() error {
        items, _ := h.service.List()
        return core.RenderSSE(c, components.ItemTable(items),
            core.WithSelector("#item-table"),
            core.WithMergeMode("morph"),
        )
    })
}
```

### Datastar merge modes

| Mode      | Effect                         |
| --------- | ------------------------------ |
| `morph`   | Smart diff and patch (default) |
| `inner`   | Replace innerHTML              |
| `outer`   | Replace outerHTML              |
| `prepend` | Insert before first child      |
| `append`  | Insert after last child        |
| `before`  | Insert before the element      |
| `after`   | Insert after the element       |

### JSON Responses

```go
core.OK(c, data)                          // 200 { success: true, data: ... }
core.Created(c, data)                     // 201 { success: true, data: ... }
core.NoContent(c)                         // 204

core.BadRequest(c, err)                   // 400 { success: false, error: "..." }
core.Unauthorized(c)                      // 401
core.Forbidden(c)                         // 403
core.NotFound(c, "customer")              // 404 { success: false, error: "customer not found" }
core.UnprocessableEntity(c, err)          // 422
core.InternalError(c, err)               // 500 (raw error never exposed)

core.ValidationErrors(c, map[string]string{
    "email": "must be a valid email",
    "name":  "required",
})                                        // 422 { success: false, errors: { ... } }

core.Redirect(c, "/customer")            // 303 See Other
```

---

## 11. Templ Components

Templ compiles `.templ` files into type-safe Go code at build time.

### File locations

```
modules/<name>/pages/         ← full page templates (include layout)
modules/<name>/components/    ← reusable fragments
web/layouts/                  ← shared base layouts
```

### Base layout example

```go
// web/layouts/base.templ
package layouts

templ Base(title string) {
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8"/>
        <title>{ title }</title>
        <link rel="stylesheet" href="/static/app.css"/>
        <script src="/static/datastar.js" defer></script>
    </head>
    <body>
        { children... }
    </body>
    </html>
}
```

### Page template

```go
// modules/customer/pages/list.templ
package pages

import (
    "goapp/modules/customer/components"
    "goapp/ent"
    "goapp/web/layouts"
)

templ ListPage(customers []*ent.Customer) {
    @layouts.Base("Customers") {
        <main>
            <h1>Customers</h1>
            @components.CustomerTable(customers)
        </main>
    }
}
```

### Component fragment

```go
// modules/customer/components/table.templ
package components

import "goapp/ent"

templ CustomerTable(customers []*ent.Customer) {
    <table id="customer-table">
        <tbody>
            for _, c := range customers {
                <tr>
                    <td>{ c.Name }</td>
                    <td>{ c.Email }</td>
                </tr>
            }
        </tbody>
    </table>
}
```

### Regenerate after changes

```bash
templ generate
# or watch mode
templ generate --watch
```

---

## 12. Datastar — Reactivity

Datastar replaces Alpine.js, htmx, and fetch calls. The server pushes HTML patches over SSE.

### Add to layout

```html
<script src="/static/datastar.js" defer></script>
```

Download the latest `datastar.js` from [data-star.dev](https://data-star.dev) into `web/static/`.

### Trigger a server action on click

```html
<!-- Sends GET /customer/list/fragment on click -->
<button data-on-click="$$get('/customer/list/fragment')">Refresh</button>
```

### Bind input to server state

```html
<input
  type="text"
  data-model="search"
  data-on-input__debounce.300ms="$$get('/customer/search?q=' + $search)"
/>
```

### Target for fragment merge

The server uses `core.WithSelector("#customer-table")` to target this element:

```html
<div id="customer-table">@components.CustomerTable(customers)</div>
```

### SSE stream route

```go
// module.go Routes()
{Method: "GET", Path: "/stream", Handler: m.handler.Stream},

// handler
func (h *Handler) Stream(c *fiber.Ctx) error {
    return core.StreamSSE(c, func() error {
        customers, _ := h.service.List()
        return core.RenderSSE(c, components.CustomerTable(customers),
            core.WithSelector("#customer-table"),
            core.WithMergeMode("morph"),
        )
    })
}
```

---

## 13. Database — Ent ORM

Ent generates a fully type-safe database client from schema definitions.

### Define a schema

```go
// ent/schema/customer.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
)

type Customer struct {
    ent.Schema
}

func (Customer) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").NotEmpty(),
        field.String("email").Unique(),
        field.Bool("active").Default(true),
        field.Time("created_at").Immutable().Default(time.Now),
    }
}

func (Customer) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("email").Unique(),
    }
}
```

### Regenerate after schema changes

```bash
go generate ./ent/...
# or
make generate
```

### Repository pattern

```go
// modules/customer/repository/customer_repo.go
package repository

import (
    "context"
    "goapp/ent"
    "goapp/ent/customer"
)

type Repo struct {
    db *ent.Client
}

func NewRepo(db *ent.Client) *Repo {
    return &Repo{db: db}
}

func (r *Repo) List(ctx context.Context) ([]*ent.Customer, error) {
    return r.db.Customer.
        Query().
        Where(customer.ActiveEQ(true)).
        Order(ent.Asc(customer.FieldName)).
        All(ctx)
}

func (r *Repo) FindByID(ctx context.Context, id int) (*ent.Customer, error) {
    return r.db.Customer.Get(ctx, id)
}

func (r *Repo) Create(ctx context.Context, name, email string) (*ent.Customer, error) {
    return r.db.Customer.
        Create().
        SetName(name).
        SetEmail(email).
        Save(ctx)
}

func (r *Repo) Update(ctx context.Context, id int, name string) (*ent.Customer, error) {
    return r.db.Customer.
        UpdateOneID(id).
        SetName(name).
        Save(ctx)
}

func (r *Repo) Delete(ctx context.Context, id int) error {
    return r.db.Customer.DeleteOneID(id).Exec(ctx)
}
```

### Transactions

```go
func (r *Repo) CreateWithOrder(ctx context.Context, ...) error {
    tx, err := r.db.Tx(ctx)
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    _, err = tx.Customer.Create().SetName(name).Save(ctx)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

---

## 14. Database Migration

Migrations are handled automatically by Ent using `Schema.Create()`.

### Auto migration (development)

Controlled by `ENABLE_MIGRATION_AUTO=true` in `.env`. Runs on every startup — safe for development, **disabled in production** by default.

```go
// internal/config/db.go
if cfg.AutoMigration {
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("migration failed: %v", err)
    }
}
```

### Manual migration (production)

```bash
# Generate migration SQL without applying
go run -mod=mod entgo.io/ent/cmd/ent migrate diff \
    --dialect sqlite \
    --dsn "file:goapp.db" \
    --dir migrations/

# Apply pending migrations
go run -mod=mod entgo.io/ent/cmd/ent migrate apply \
    --dir migrations/ \
    --dsn "file:goapp.db"
```

### Safety check

The app will refuse to start if `ENABLE_MIGRATION_AUTO=true` and `APP_ENV=production`:

```go
// internal/config/db.go
if cfg.AppEnv == "production" && cfg.AutoMigration {
    log.Fatal("ENABLE_MIGRATION_AUTO must be false in production")
}
```

---

## 15. Database Seeding

Create a seeder per module. Seeders are plain Go functions, not a framework concept.

### Create a seeder

```go
// modules/customer/seed.go
package customer

import (
    "context"
    "goapp/ent"
)

func Seed(ctx context.Context, db *ent.Client) error {
    exists, err := db.Customer.Query().Count(ctx)
    if err != nil {
        return err
    }
    if exists > 0 {
        return nil // already seeded
    }

    customers := []struct{ Name, Email string }{
        {"Alice", "alice@example.com"},
        {"Bob",   "bob@example.com"},
    }

    for _, c := range customers {
        _, err := db.Customer.Create().
            SetName(c.Name).
            SetEmail(c.Email).
            Save(ctx)
        if err != nil {
            return err
        }
    }

    return nil
}
```

### Wire it up in PostRegister

```go
func (m *Module) PostRegister(ctx context.Context) error {
    if m.deps.Config.AppEnv == "development" {
        if err := Seed(ctx, m.deps.DB); err != nil {
            return fmt.Errorf("seed failed: %w", err)
        }
    }
    return nil
}
```

### Or run via CLI

```bash
make seed
```

```makefile
# Makefile
seed:
	go run ./cmd/seed
```

---

## 16. Middleware

Global middleware is registered in `internal/core/middleware.go` and applies to all routes.

### Built-in middleware (applied in order)

| Middleware                | Purpose                             |
| ------------------------- | ----------------------------------- |
| `RecoveryMiddleware`      | Catches panics, returns 500         |
| `CORSMiddleware`          | Sets CORS headers from config       |
| `RequestTimingMiddleware` | Logs method, path, status, duration |

### Add a new global middleware

```go
// internal/core/middleware.go
func NewMiddlewareRegistry(app *fiber.App, cfg *Config, logger *slog.Logger) *MiddlewareRegistry {
    r := &MiddlewareRegistry{}

    r.Register(RecoveryMiddleware(logger))
    r.Register(CORSMiddleware(cfg))
    r.Register(RequestTimingMiddleware(logger))
    r.Register(YourNewMiddleware(cfg))    // ← add here

    for _, m := range r.middlewares {
        app.Use(m)
    }

    return r
}
```

### Add a module-scoped middleware

Apply middleware to specific routes only inside `Routes()`:

```go
func (m *Module) Routes() []core.Route {
    auth := middleware.RequireAuth(m.deps.Config.JWTSecret)

    return []core.Route{
        {Method: "GET",  Path: "/",    Handler: m.handler.ListPage},          // public
        {Method: "POST", Path: "/",    Handler: fiber.WrapHandler(auth, m.handler.Create)}, // protected
    }
}
```

---

## 17. Configuration System

### Reading a new environment variable

**Step 1 — Add to Config struct** (`internal/core/config.go`):

```go
type Config struct {
    // ...
    StripeKey string
}
```

**Step 2 — Read in NewConfig()**:

```go
func NewConfig() *Config {
    return &Config{
        // ...
        StripeKey: getEnv("STRIPE_KEY", ""),
    }
}
```

**Step 3 — Add to `.env`**:

```dotenv
STRIPE_KEY=sk_test_...
```

**Step 4 — Use anywhere via container**:

```go
key := m.deps.Config.StripeKey
```

### getEnv helper

```go
// Returns the env var value, or fallback if empty/unset
getEnv("KEY", "default")
```

---

## 18. Logging

The logger is `*slog.Logger` from Go stdlib, available on the container.

### Usage in a module

```go
m.deps.Logger.Info("customer created", "id", customer.ID, "name", customer.Name)
m.deps.Logger.Warn("rate limit approaching", "count", count)
m.deps.Logger.Error("payment failed", "error", err, "customerID", id)
```

### Module-scoped logger (recommended)

Add a module-specific logger with context in `New()`:

```go
func New(c *core.Container) *Module {
    return &Module{
        deps:   c,
        logger: c.Logger.With("module", "customer"),
    }
}
```

All log lines from this module will include `"module":"customer"` automatically.

### Log levels

Controlled by `LOG_LEVEL` in `.env`:

| Level   | When to use                          |
| ------- | ------------------------------------ |
| `debug` | Detailed flow tracing (dev only)     |
| `info`  | Normal operations, module start/stop |
| `warn`  | Recoverable issues, deprecations     |
| `error` | Failures that need attention         |

---

## 19. Error Handling

### Handler errors

Always return typed errors using `core/response.go` helpers — never write raw `c.Status(...).JSON(...)`:

```go
func (h *Handler) Detail(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return core.BadRequest(c, fmt.Errorf("invalid id"))
    }

    customer, err := h.service.FindByID(c.Context(), id)
    if ent.IsNotFound(err) {
        return core.NotFound(c, "customer")
    }
    if err != nil {
        return core.InternalError(c, err)
    }

    return core.Render(c, pages.DetailPage(customer))
}
```

### Ent error helpers

```go
ent.IsNotFound(err)        // record not found
ent.IsConstraintError(err) // unique constraint violation
ent.IsValidationError(err) // field validation failed
```

### Panics

The `RecoveryMiddleware` catches all panics automatically and returns a 500 response. The full stack trace is logged.

---

## 20. Makefile Commands

```makefile
make run        # go run ./cmd/web
make build      # templ generate + go build -o bin/goapp ./cmd/web
make watch      # templ generate --watch & air (live reload)
make generate   # go generate ./ent/... (regenerate Ent code)
make templ      # templ generate (regenerate templ components)
make seed       # go run ./cmd/seed
make test       # go test ./...
make lint       # golangci-lint run
make clean      # rm -rf bin/
```

---

## 21. Adding a Full CRUD Module — End to End Example

This walks through adding a complete `Product` module from scratch.

### 1. Ent schema

```go
// ent/schema/product.go
package schema

import (
    "time"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type Product struct{ ent.Schema }

func (Product) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").NotEmpty(),
        field.Float("price").Positive(),
        field.Int("stock").Default(0),
        field.Time("created_at").Immutable().Default(time.Now),
    }
}
```

```bash
go generate ./ent/...
```

### 2. Repository

```go
// modules/product/repository/product_repo.go
package repository

import (
    "context"
    "goapp/ent"
)

type Repo struct{ db *ent.Client }

func NewRepo(db *ent.Client) *Repo { return &Repo{db: db} }

func (r *Repo) List(ctx context.Context) ([]*ent.Product, error) {
    return r.db.Product.Query().All(ctx)
}

func (r *Repo) Create(ctx context.Context, name string, price float64) (*ent.Product, error) {
    return r.db.Product.Create().SetName(name).SetPrice(price).Save(ctx)
}
```

### 3. Service

```go
// modules/product/service/product_service.go
package service

import (
    "context"
    "fmt"
    "goapp/ent"
    "goapp/modules/product/repository"
)

type Service struct{ repo *repository.Repo }

func NewService(repo *repository.Repo) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context) ([]*ent.Product, error) {
    return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, name string, price float64) (*ent.Product, error) {
    if name == "" {
        return nil, fmt.Errorf("name is required")
    }
    if price <= 0 {
        return nil, fmt.Errorf("price must be positive")
    }
    return s.repo.Create(ctx, name, price)
}
```

### 4. Handler

```go
// modules/product/handler/product_handler.go
package handler

import (
    "goapp/internal/core"
    "goapp/modules/product/pages"
    "goapp/modules/product/service"

    "github.com/gofiber/fiber/v2"
)

type Handler struct{ service *service.Service }

func NewHandler(s *service.Service) *Handler { return &Handler{service: s} }

func (h *Handler) ListPage(c *fiber.Ctx) error {
    products, err := h.service.List(c.Context())
    if err != nil {
        return core.InternalError(c, err)
    }
    return core.Render(c, pages.ListPage(products))
}

func (h *Handler) Create(c *fiber.Ctx) error {
    type req struct {
        Name  string  `json:"name"`
        Price float64 `json:"price"`
    }
    var body req
    if err := c.BodyParser(&body); err != nil {
        return core.BadRequest(c, err)
    }
    _, err := h.service.Create(c.Context(), body.Name, body.Price)
    if err != nil {
        return core.UnprocessableEntity(c, err)
    }
    return core.Redirect(c, "/product")
}
```

### 5. Module

```go
// modules/product/module.go
package product

import (
    "context"
    "goapp/internal/core"
    "goapp/modules/product/handler"
    "goapp/modules/product/repository"
    "goapp/modules/product/service"
)

var _ core.Module = (*Module)(nil)

type Module struct {
    core.BaseModule
    deps    *core.Container
    handler *handler.Handler
    service *service.Service
    repo    *repository.Repo
}

func init() {
    core.RegisterModule(func(c *core.Container) core.Module { return New(c) })
}

func New(c *core.Container) *Module { return &Module{deps: c} }

func (m *Module) Name() string     { return "product" }
func (m *Module) BasePath() string { return "/product" }

func (m *Module) Register(bus *core.EventBus) error {
    m.repo    = repository.NewRepo(m.deps.DB)
    m.service = service.NewService(m.repo)
    m.handler = handler.NewHandler(m.service)
    return nil
}

func (m *Module) PostRegister(ctx context.Context) error {
    return m.deps.EventBus.Publish(ctx, core.Event{
        Topic: "product.ready", Source: m.Name(),
    })
}

func (m *Module) Routes() []core.Route {
    return []core.Route{
        {Method: "GET",  Path: "/",  Handler: m.handler.ListPage},
        {Method: "POST", Path: "/",  Handler: m.handler.Create},
    }
}

func (m *Module) Service() any { return m.service }
```

### 6. Register in main.go

```go
import (
    _ "goapp/modules/example"
    _ "goapp/modules/product"   // ← add this
)
```

### 7. Done

```
GET  /product      → renders product list page
POST /product      → creates product, redirects to list
```

---

_Last updated: May 2026_
