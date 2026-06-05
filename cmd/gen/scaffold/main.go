package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// ─────────────────────────────────────────────────────────────────────────────
// Data model passed to every template
// ─────────────────────────────────────────────────────────────────────────────

type Module struct {
	Name      string // raw name as passed:          "sales_order"
	Pascal    string // PascalCase:                   "SalesOrder"
	Camel     string // camelCase:                    "salesOrder"
	Title     string // Title case with spaces:       "Sales Order"
	Snake     string // snake_case (same as Name):    "sales_order"
	BasePath  string // URL base path:                "/sales_order"
	ModulePkg string // Go import path of the module: "goapp/modules/sales_order"
	AppModule string // go.mod module name:           "goapp"
}

// ─────────────────────────────────────────────────────────────────────────────
// Templates
// ─────────────────────────────────────────────────────────────────────────────

var templates = map[string]string{

	// ── manifest.json ─────────────────────────────────────────────────────────
	"manifest.json": `{
    "name": "{{.Name}}",
    "base_path": "/{{.Name}}",
    "enabled": true,
    "version": "1.0.0",
    "description": "{{.Title}} module",
    "dependencies": [],
    "middlewares": [],
    "tags": ["api"]
}`,

	// ── module.go ─────────────────────────────────────────────────────────────
	"module.go": `package {{.Name}}

import (
	"context"

	"{{.AppModule}}/internal/core"
	"{{.ModulePkg}}/handler"
	"{{.ModulePkg}}/repository"
	"{{.ModulePkg}}/service"
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
		panic("{{.Name}} module: " + err.Error())
	}
	return m
}

// ── Lifecycle ─────────────────────────────────────────────────────────────────

func (m *Module) Init(_ context.Context) error {
	m.Logger.Info("init", "module", "{{.Name}}")
	return nil
}

func (m *Module) Register(bus *core.EventBus) error {
	m.repos    = repository.NewRepositories(m.deps.DB)
	m.services = service.NewServices(m.repos, m.Logger, m.deps.Store)
	m.handlers = handler.NewHandlers(m.services)

	// TODO: subscribe to events from other modules
	// bus.Subscribe("topic", func(ctx context.Context, e core.Event) error {
	// 	return nil
	// })

	return nil
}

func (m *Module) PostRegister(ctx context.Context, deps map[string]core.Module) error {
	m.Logger.Info("post-register", "module", "{{.Name}}")

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
	// if m.deps.Config.AppEnv == "development" {
	//     if err := Seed(ctx, m.deps.DB); err != nil {
	//         return fmt.Errorf("seed failed: %w", err)
	//     }
	// }

	// ── Background workers ────────────────────────────────────────────────────
	// go m.startWorker(ctx)

	return m.deps.EventBus.Publish(ctx, core.Event{
		Topic:   "{{.Name}}.ready",
		Payload: map[string]string{"module": "{{.Name}}"},
		Source:  "{{.Name}}",
	})
}

func (m *Module) Shutdown(ctx context.Context) error {
	m.Logger.Info("shutdown", "module", "{{.Name}}")
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
`,

	// ── routes.go ─────────────────────────────────────────────────────────────
	"routes.go": `package {{.Name}}

import (
	"{{.AppModule}}/internal/core"
	"{{.ModulePkg}}/handler"
)

// Routes declares all HTTP routes for the {{.Title}} module.
// Add new routes here as the module grows — no need to touch module.go.
func Routes(h *handler.Handlers) []core.Route {
	return []core.Route{
		// {{.Pascal}}
		{Method: "GET",    Path: "/",     Handler: h.{{.Pascal}}.Index},
		{Method: "GET",    Path: "/list", Handler: h.{{.Pascal}}.List},
		{Method: "POST",   Path: "/",     Handler: h.{{.Pascal}}.Create},
		{Method: "GET",    Path: "/:id",  Handler: h.{{.Pascal}}.Show},
		{Method: "PUT",    Path: "/:id",  Handler: h.{{.Pascal}}.Update},
		{Method: "DELETE", Path: "/:id",  Handler: h.{{.Pascal}}.Delete},
	}
}
`,

	// ── handler/handlers.go ───────────────────────────────────────────────────
	"handler/handlers.go": `package handler

import "{{.ModulePkg}}/service"

// Handlers groups all handlers for the {{.Title}} module.
// Add new handler structs here as the module grows:
//
//	type Handlers struct {
//	    {{.Pascal}}  *{{.Pascal}}Handler
//	    Journal *JournalHandler
//	    Report  *ReportHandler
//	}
type Handlers struct {
	{{.Pascal}} *{{.Pascal}}Handler
}

// NewHandlers wires all handlers with their service dependencies.
func NewHandlers(svc *service.Services) *Handlers {
	return &Handlers{
		{{.Pascal}}: New{{.Pascal}}Handler(svc),
	}
}
`,

	// ── handler/{{.Name}}_handler.go ──────────────────────────────────────────
	"handler/{{.Name}}_handler.go": `package handler

import (
	"strconv"

	"{{.AppModule}}/internal/core"
	"{{.ModulePkg}}/pages"
	"{{.ModulePkg}}/service"

	"github.com/gofiber/fiber/v2"
)

type {{.Pascal}}Handler struct {
	service *service.Services
}

func New{{.Pascal}}Handler(svc *service.Services) *{{.Pascal}}Handler {
	return &{{.Pascal}}Handler{service: svc}
}

// Index renders the landing page.
func (h *{{.Pascal}}Handler) Index(c *fiber.Ctx) error {
	return core.Render(c, pages.{{.Pascal}}Page("{{.Title}}"))
}

// List returns all records paginated.
func (h *{{.Pascal}}Handler) List(c *fiber.Ctx) error {
	params := core.PaginationParams{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
		Search:  c.Query("q"),
	}
	result, err := h.service.{{.Pascal}}.ListPaginated(c.Context(), params)
	if err != nil {
		return core.InternalError(c, err)
	}
	return core.OK(c, result)
}

// Show returns a single record by ID.
func (h *{{.Pascal}}Handler) Show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return core.BadRequest(c, err)
	}
	item, err := h.service.{{.Pascal}}.FindByID(c.Context(), id)
	if err != nil {
		return core.NotFound(c, "{{.Name}}")
	}
	return core.OK(c, item)
}

// Create creates a new record.
func (h *{{.Pascal}}Handler) Create(c *fiber.Ctx) error {
	var body struct {
		Name        string ` + "`json:\"name\"`" + `
		Description string ` + "`json:\"description\"`" + `
	}
	if err := c.BodyParser(&body); err != nil {
		return core.BadRequest(c, err)
	}
	item, err := h.service.{{.Pascal}}.Create(c.Context(), body.Name, body.Description)
	if err != nil {
		return core.UnprocessableEntity(c, err)
	}
	return core.Created(c, item)
}

// Update modifies an existing record.
func (h *{{.Pascal}}Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return core.BadRequest(c, err)
	}
	var body struct {
		Name        string ` + "`json:\"name\"`" + `
		Description string ` + "`json:\"description\"`" + `
	}
	if err := c.BodyParser(&body); err != nil {
		return core.BadRequest(c, err)
	}
	item, err := h.service.{{.Pascal}}.Update(c.Context(), id, body.Name, body.Description)
	if err != nil {
		return core.UnprocessableEntity(c, err)
	}
	return core.OK(c, item)
}

// Delete removes a record by ID.
func (h *{{.Pascal}}Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return core.BadRequest(c, err)
	}
	if err := h.service.{{.Pascal}}.Delete(c.Context(), id); err != nil {
		return core.InternalError(c, err)
	}
	return core.NoContent(c)
}
`,

	// ── service/services.go ───────────────────────────────────────────────────
	"service/services.go": `package service

import (
	"log/slog"

	"{{.AppModule}}/internal/core"
	"{{.ModulePkg}}/repository"
)

// Services groups all services for the {{.Title}} module.
// Add new service structs here as the module grows:
//
//	type Services struct {
//	    {{.Pascal}}  *{{.Pascal}}Service
//	    Journal *JournalService
//	    Report  *ReportService
//	}
type Services struct {
	{{.Pascal}} *{{.Pascal}}Service
}

// NewServices wires all services with their repository dependencies.
func NewServices(repos *repository.Repositories, logger *slog.Logger, store *core.Store) *Services {
	return &Services{
		{{.Pascal}}: New{{.Pascal}}Service(repos, logger),
	}
}
`,

	// ── service/{{.Name}}_service.go ──────────────────────────────────────────
	"service/{{.Name}}_service.go": `package service

import (
	"context"
	"log/slog"

	"{{.AppModule}}/internal/core"
	"{{.ModulePkg}}/repository"
)

// {{.Pascal}}Service holds business logic for {{.Title}}.
// Replace 'any' return types once the ent schema is generated.
type {{.Pascal}}Service struct {
	core.BaseService
	repos *repository.Repositories
}

func New{{.Pascal}}Service(repos *repository.Repositories, logger *slog.Logger) *{{.Pascal}}Service {
	return &{{.Pascal}}Service{
		BaseService: core.BaseService{Logger: logger},
		repos:       repos,
	}
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (s *{{.Pascal}}Service) List(ctx context.Context) ([]any, error) {
	return s.repos.{{.Pascal}}.List(ctx)
}

func (s *{{.Pascal}}Service) ListPaginated(ctx context.Context, params core.PaginationParams) (*core.PaginatedResult[any], error) {
	return core.Paginate(ctx, params, s.repos.{{.Pascal}}.ListPaginated)
}

func (s *{{.Pascal}}Service) FindByID(ctx context.Context, id int) (any, error) {
	item, err := s.repos.{{.Pascal}}.FindByID(ctx, id)
	return item, s.WrapError("FindByID", err)
}

// ── Commands ──────────────────────────────────────────────────────────────────

func (s *{{.Pascal}}Service) Create(ctx context.Context, name, description string) (any, error) {
	if err := s.ValidateRequired(map[string]string{"name": name}); err != nil {
		return nil, err
	}
	if err := s.ValidateMinLength("name", name, 2); err != nil {
		return nil, err
	}
	if err := s.ValidateMaxLength("name", name, 255); err != nil {
		return nil, err
	}
	item, err := s.repos.{{.Pascal}}.Create(ctx, name, description)
	return item, s.WrapError("Create", err)
}

func (s *{{.Pascal}}Service) Update(ctx context.Context, id int, name, description string) (any, error) {
	if err := s.ValidateRequired(map[string]string{"name": name}); err != nil {
		return nil, err
	}
	if err := s.ValidateMinLength("name", name, 2); err != nil {
		return nil, err
	}
	if err := s.ValidateMaxLength("name", name, 255); err != nil {
		return nil, err
	}
	item, err := s.repos.{{.Pascal}}.Update(ctx, id, name, description)
	return item, s.WrapError("Update", err)
}

func (s *{{.Pascal}}Service) Delete(ctx context.Context, id int) error {
	return s.WrapError("Delete", s.repos.{{.Pascal}}.Delete(ctx, id))
}
`,

	// ── service/{{.Name}}_service_test.go ─────────────────────────────────────
	"service/{{.Name}}_service_test.go": `package service_test

import (
	"context"
	"testing"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	_ = ctx

	t.Run("requires name", func(t *testing.T) {
		// TODO: wire up service with a mock repo and assert ValidateRequired fires
	})

	t.Run("name too short", func(t *testing.T) {
		// TODO: assert ValidateMinLength fires for single-char names
	})
}
`,

	// ── repository/repositories.go ────────────────────────────────────────────
	"repository/repositories.go": `package repository

import "{{.AppModule}}/ent"

// Repositories groups all repositories for the {{.Title}} module.
// Add new repository structs here as the module grows:
//
//	type Repositories struct {
//	    {{.Pascal}}  *{{.Pascal}}Repository
//	    Journal *JournalRepository
//	}
type Repositories struct {
	{{.Pascal}} *{{.Pascal}}Repository
}

// NewRepositories wires all repositories with the shared DB client.
func NewRepositories(db *ent.Client) *Repositories {
	return &Repositories{
		{{.Pascal}}: New{{.Pascal}}Repository(db),
	}
}
`,

	// ── repository/{{.Name}}_repository.go ───────────────────────────────────
	"repository/{{.Name}}_repository.go": `package repository

import (
	"context"

	"{{.AppModule}}/ent"
	"{{.AppModule}}/internal/core"
)

// {{.Pascal}}Repository handles all DB operations for {{.Title}}.
// TODO: replace 'any' with generated ent types after: go generate ./ent/...
type {{.Pascal}}Repository struct {
	db *ent.Client
}

func New{{.Pascal}}Repository(db *ent.Client) *{{.Pascal}}Repository {
	return &{{.Pascal}}Repository{db: db}
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (r *{{.Pascal}}Repository) List(ctx context.Context) ([]any, error) {
	// TODO: return r.db.{{.Pascal}}.Query().All(ctx)
	return []any{}, nil
}

func (r *{{.Pascal}}Repository) ListPaginated(ctx context.Context, offset, limit int) ([]any, int, error) {
	// TODO:
	// q     := r.db.{{.Pascal}}.Query()
	// total, _ := q.Count(ctx)
	// items, _ := q.Limit(limit).Offset(offset).All(ctx)
	return []any{}, 0, nil
}

func (r *{{.Pascal}}Repository) FindByID(ctx context.Context, id int) (any, error) {
	// TODO: return r.db.{{.Pascal}}.Get(ctx, id)
	return nil, nil
}

func (r *{{.Pascal}}Repository) Create(ctx context.Context, name, description string) (any, error) {
	// TODO: return r.db.{{.Pascal}}.Create().SetName(name).SetDescription(description).Save(ctx)
	return nil, nil
}

func (r *{{.Pascal}}Repository) Update(ctx context.Context, id int, name, description string) (any, error) {
	// TODO: return r.db.{{.Pascal}}.UpdateOneID(id).SetName(name).SetDescription(description).Save(ctx)
	return nil, nil
}

func (r *{{.Pascal}}Repository) Delete(ctx context.Context, id int) error {
	// TODO: return r.db.{{.Pascal}}.DeleteOneID(id).Exec(ctx)
	return nil
}

// ── Custom queries ────────────────────────────────────────────────────────────

func (r *{{.Pascal}}Repository) SearchByName(ctx context.Context, params core.PaginationParams) ([]any, int, error) {
	// TODO:
	// q := r.db.{{.Pascal}}.Query()
	// if params.Search != "" {
	//     q = q.Where({{.Name}}.NameContainsFold(params.Search))
	// }
	return []any{}, 0, nil
}
`,

	// ── ent schema ────────────────────────────────────────────────────────────
	"ent/schema/{{.Name}}.go": `package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	coremixin "{{.AppModule}}/internal/core/mixin"
)

// {{.Pascal}} holds the schema definition for the {{.Pascal}} entity.
type {{.Pascal}} struct {
	ent.Schema
}

func ({{.Pascal}}) Mixin() []ent.Mixin {
	return []ent.Mixin{
		coremixin.Base{},
	}
}

func ({{.Pascal}}) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(255),

		field.String("description").
			Optional().
			Nillable(),
	}
}

func ({{.Pascal}}) Edges() []ent.Edge {
	return nil
}
`,

	// ── templ page ────────────────────────────────────────────────────────────
	"pages/{{.Name}}_page.templ": `package pages

import "{{.AppModule}}/web/layouts"

templ {{.Pascal}}Page(title string) {
	@layouts.Base(title) {
		<div class="p-6">
			<h1 class="text-2xl font-bold mb-4">{ title }</h1>
			<p class="text-gray-500">{{.Title}} module — add your UI here.</p>
		</div>
	}
}
`,

	// ── seed ──────────────────────────────────────────────────────────────────
	"seed.go": `package {{.Name}}

import (
	"context"
	"{{.AppModule}}/ent"
)

// Seed inserts default records for the {{.Title}} module.
// It is a no-op if records already exist.
// Called from PostRegister when APP_ENV=development.
func Seed(ctx context.Context, db *ent.Client) error {
	// TODO: replace with real ent query once schema is generated
	// count, err := db.{{.Pascal}}.Query().Count(ctx)
	// if err != nil || count > 0 {
	//     return err
	// }
	// _, err = db.{{.Pascal}}.Create().SetName("Sample {{.Title}}").Save(ctx)
	// return err
	return nil
}
`,
}

// ─────────────────────────────────────────────────────────────────────────────
// String helpers
// ─────────────────────────────────────────────────────────────────────────────

func toPascal(s string) string {
	parts := splitWords(s)
	var b strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		runes := []rune(p)
		b.WriteRune(unicode.ToUpper(runes[0]))
		b.WriteString(strings.ToLower(string(runes[1:])))
	}
	return b.String()
}

func toCamel(s string) string {
	pascal := toPascal(s)
	if pascal == "" {
		return ""
	}
	runes := []rune(pascal)
	return string(unicode.ToLower(runes[0])) + string(runes[1:])
}

func toTitle(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		runes := []rune(p)
		parts[i] = string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
	}
	return strings.Join(parts, " ")
}

func splitWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
}

func readAppModule() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "goapp"
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "goapp"
}

// ─────────────────────────────────────────────────────────────────────────────
// Validation
// ─────────────────────────────────────────────────────────────────────────────

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("module name too long (max 64 characters)")
	}
	runes := []rune(name)
	if !unicode.IsLetter(runes[0]) {
		return fmt.Errorf("module name must start with a letter, got %q", string(runes[0]))
	}
	for _, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("invalid character %q — use letters, digits, and underscores only", string(r))
		}
	}
	reserved := map[string]bool{
		"main": true, "core": true, "ent": true, "web": true,
		"cmd": true, "internal": true, "pkg": true,
	}
	if reserved[name] {
		return fmt.Errorf("%q is a reserved name, choose a different module name", name)
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// File writing
// ─────────────────────────────────────────────────────────────────────────────

func writeFile(path, tmplSrc string, data Module, dryRun bool) (skipped bool, err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return false, err
	}
	if _, err := os.Stat(path); err == nil {
		return true, nil // never overwrite
	}
	if dryRun {
		return false, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	tmpl, err := template.New("").Parse(tmplSrc)
	if err != nil {
		return false, fmt.Errorf("template parse error: %w", err)
	}
	return false, tmpl.Execute(f, data)
}

// ─────────────────────────────────────────────────────────────────────────────
// Main
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	nameFlag := flag.String("name", "", "module name in snake_case (e.g. example, sales_order)")
	dryRun := flag.Bool("dry-run", false, "print files that would be created without writing them")
	flag.Parse()

	name := *nameFlag
	if name == "" && flag.NArg() > 0 {
		name = flag.Arg(0)
	}

	if name == "" {
		fmt.Fprintln(os.Stderr, "usage: scaffold --name=<module_name>")
		fmt.Fprintln(os.Stderr, "       scaffold example")
		fmt.Fprintln(os.Stderr, "       scaffold --dry-run example")
		os.Exit(1)
	}

	name = strings.ToLower(strings.TrimSpace(name))

	if err := validateName(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "error: run scaffold from the project root (where go.mod is)")
		os.Exit(1)
	}

	appModule := readAppModule()

	data := Module{
		Name:      name,
		Pascal:    toPascal(name),
		Camel:     toCamel(name),
		Title:     toTitle(name),
		Snake:     name,
		BasePath:  "/" + name,
		ModulePkg: appModule + "/modules/" + name,
		AppModule: appModule,
	}

	moduleDir := filepath.Join("modules", name)

	if _, err := os.Stat(moduleDir); err == nil {
		fmt.Fprintf(os.Stderr, "error: module %q already exists at %s\n", name, moduleDir)
		os.Exit(1)
	}

	if *dryRun {
		fmt.Printf("\n🔍 Dry run — files that would be created for module: %s\n\n", data.Pascal)
	} else {
		fmt.Printf("\n🔧 Scaffolding module: %s\n\n", data.Pascal)
	}

	created := 0
	skipped := 0

	// Deterministic output order
	order := []string{
		"manifest.json",
		"module.go",
		"routes.go",
		"handler/handlers.go",
		"handler/{{.Name}}_handler.go",
		"service/services.go",
		"service/{{.Name}}_service.go",
		"service/{{.Name}}_service_test.go",
		"repository/repositories.go",
		"repository/{{.Name}}_repository.go",
		"pages/{{.Name}}_page.templ",
		"seed.go",
		"ent/schema/{{.Name}}.go",
	}

	for _, pathTmpl := range order {
		content, ok := templates[pathTmpl]
		if !ok {
			continue
		}

		pt, err := template.New("").Parse(pathTmpl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "path template error: %v\n", err)
			os.Exit(1)
		}
		var pathBuf strings.Builder
		if err := pt.Execute(&pathBuf, data); err != nil {
			fmt.Fprintf(os.Stderr, "path render error: %v\n", err)
			os.Exit(1)
		}

		var fullPath string
		if strings.HasPrefix(pathBuf.String(), "ent/") {
			fullPath = pathBuf.String()
		} else {
			fullPath = filepath.Join(moduleDir, pathBuf.String())
		}

		wasSkipped, err := writeFile(fullPath, content, data, *dryRun)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", fullPath, err)
			os.Exit(1)
		}

		if wasSkipped {
			fmt.Printf("  skip    %s\n", fullPath)
			skipped++
		} else {
			fmt.Printf("  create  %s\n", fullPath)
			created++
		}
	}

	fmt.Println()
	if *dryRun {
		fmt.Printf("Dry run complete — %d files would be created, %d skipped.\n\n", created, skipped)
		return
	}

	fmt.Printf("✅  Module %q scaffolded — %d files created, %d skipped.\n\n", name, created, skipped)
	fmt.Println("Next steps:")
	fmt.Printf("  1. Review  ent/schema/%s.go\n", name)
	fmt.Printf("  2. Run     go generate ./ent/...\n")
	fmt.Printf("  3. Update  modules/%s/repository/%s_repository.go  (replace any → *ent.%s)\n", name, name, data.Pascal)
	fmt.Printf("  4. Update  modules/%s/service/%s_service.go        (replace any → *ent.%s)\n", name, name, data.Pascal)
	fmt.Printf("  5. Run     make modules\n")
	fmt.Printf("  6. Run     templ generate\n")
	fmt.Printf("  7. Run     make run\n\n")
	fmt.Printf("  Routes registered at: %s\n\n", data.BasePath)
}
