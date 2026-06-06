package form

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

// ─────────────────────────────────────────────────────────────────────────────
// Engine
// ─────────────────────────────────────────────────────────────────────────────

// Engine holds the live state of one form instance.
// One Engine per form instance — never share across requests.
type Engine struct {
	model  FormModel
	state  FormState
	graph  *depGraph
	schema FormSchema // nil when constructed from FormModel directly
}

// NewEngine builds an Engine from a FormModel and initialises all fields.
// Computed fields are evaluated immediately if their dependencies have defaults.
func NewEngine(model FormModel) (*Engine, error) {
	if err := detectCycles(model); err != nil {
		return nil, fmt.Errorf("form: invalid model: %w", err)
	}

	e := &Engine{
		model: model,
		state: make(FormState, len(model)),
		graph: buildDepGraph(model),
	}

	// Apply defaults for source fields first
	for key, field := range model {
		if field.Compute == nil && field.Default != nil {
			e.state[key] = field.Default
		} else {
			e.state[key] = nil
		}
	}

	// Recompute all computed fields in dependency order
	e.recomputeAll()

	return e, nil
}

// NewEngineFromSchema builds an Engine from a FormSchema (ordered).
func NewEngineFromSchema(schema FormSchema) (*Engine, error) {
	e, err := NewEngine(schema.ToModel())
	if err != nil {
		return nil, err
	}
	e.schema = schema
	return e, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Read / Write
// ─────────────────────────────────────────────────────────────────────────────

// SetField writes a value to a source field and recomputes all dependents.
// Returns an error if the key is unknown or the field is computed.
func (e *Engine) SetField(key string, value any) error {

	field, ok := e.model[key]
	if !ok {
		return fmt.Errorf("form: unknown field %q", key)
	}
	if field.Compute != nil {
		return fmt.Errorf("form: field %q is computed and cannot be set directly", key)
	}
	e.state[key] = value
	e.recompute(key)
	return nil
}

// GetField returns the current value of any field (source or computed).
func (e *Engine) GetField(key string) (any, error) {
	if _, ok := e.model[key]; !ok {
		return nil, fmt.Errorf("form: unknown field %q", key)
	}
	return e.state[key], nil
}

// GetState returns a shallow copy of the full form state.
func (e *Engine) GetState() FormState {
	snap := make(FormState, len(e.state))
	for k, v := range e.state {
		snap[k] = v
	}
	return snap
}

// Apply writes multiple source-field values from a map (e.g. HTTP form POST).
// Unknown keys and computed keys are silently skipped — safe for untrusted input.
// Returns a map of fieldKey → error for any field whose value failed SetField.
//

func (e *Engine) Apply(values map[string]any) map[string]error {
	errs := map[string]error{}
	for key, value := range values {
		field, ok := e.model[key]
		if !ok || field.Compute != nil {
			continue // skip unknown and computed
		}
		if err := e.SetField(key, value); err != nil {
			errs[key] = err
		}
	}
	return errs
}

// Binding to Fiber Context
func (e *Engine) BindRequest(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var values map[string]any

		if err := c.BodyParser(&values); err != nil {
			return err
		}

		e.Apply(values)
		return nil
	}

	values := map[string]any{}

	for key := range e.model {
		values[key] = c.FormValue(key)
	}

	e.Apply(values)

	return nil
}

func (f FormState) Bind(out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("out must be pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("form")
		if key == "" {
			key = toSnake(field.Name)
			// continue // force explicit mapping only
		}

		val, ok := f[key]
		if !ok {
			continue
		}

		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		setValue(fv, val)
	}

	return nil
}

func toSnake(s string) string {
	var result strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			// add underscore before uppercase (except first char)
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func setValue(field reflect.Value, val any) {
	if val == nil {
		return
	}

	switch field.Kind() {

	case reflect.String:
		field.SetString(fmt.Sprint(val))

	case reflect.Int, reflect.Int64:
		i, err := strconv.ParseInt(fmt.Sprint(val), 10, 64)
		if err != nil {
			panic(fmt.Errorf("invalid int value: %v", val))
		}
		field.SetInt(i)

	case reflect.Float64:
		f, err := strconv.ParseFloat(fmt.Sprint(val), 64)
		if err != nil {
			return
		}
		field.SetFloat(f)

	case reflect.Pointer:
		elemType := field.Type().Elem()
		elem := reflect.New(elemType).Elem()

		setValue(elem, val)

		if !elem.IsZero() {
			field.Set(elem.Addr())
		}
	}
}

// Reset clears all fields to nil (or their Default) and recomputes.
func (e *Engine) Reset() {
	for key, field := range e.model {
		if field.Compute == nil {
			if field.Default != nil {
				e.state[key] = field.Default
			} else {
				e.state[key] = nil
			}
		}
	}
	e.recomputeAll()
}

// ─────────────────────────────────────────────────────────────────────────────
// Validation
// ─────────────────────────────────────────────────────────────────────────────

// Validate runs every field's ValidateFn (and built-in Required check)
// against the current state.
// Returns a map of fieldKey → error. Empty map = form is valid.
func (e *Engine) Validate() map[string]error {
	errs := map[string]error{}

	for key, field := range e.model {
		// FIX: built-in required check so callers don't have to repeat it
		if field.Required && isEmpty(e.state[key]) {
			errs[key] = fmt.Errorf("%s is required", labelOf(field))
			continue
		}
		if field.Validate == nil {
			continue
		}
		if err := field.Validate(e.state[key]); err != nil {
			errs[key] = err
		}
	}

	return errs
}

// ValidateField validates a single field and returns its error (or nil).
func (e *Engine) ValidateField(key string) error {
	field, ok := e.model[key]
	if !ok {
		return fmt.Errorf("form: unknown field %q", key)
	}
	if field.Required && isEmpty(e.state[key]) {
		return fmt.Errorf("%s is required", labelOf(field))
	}
	if field.Validate == nil {
		return nil
	}
	return field.Validate(e.state[key])
}

// IsValid returns true when Validate() produces no errors.
func (e *Engine) IsValid() bool {
	return len(e.Validate()) == 0
}

// ─────────────────────────────────────────────────────────────────────────────
// Dependency graph
// ─────────────────────────────────────────────────────────────────────────────

type depGraph struct {
	// dependents[key] = fields whose Compute reads key
	dependents map[string][]string
	// all computed field keys in any topological order (for recomputeAll)
	computedKeys []string
}

func buildDepGraph(model FormModel) *depGraph {
	g := &depGraph{
		dependents: make(map[string][]string),
	}
	for key, field := range model {
		for _, dep := range field.DependsOn {
			g.dependents[dep] = append(g.dependents[dep], key)
		}
		if field.Compute != nil {
			g.computedKeys = append(g.computedKeys, key)
		}
	}
	return g
}

// topoOrder returns computed fields that need refreshing after changedKey,
// in breadth-first dependency order.
func (g *depGraph) topoOrder(changedKey string) []string {
	visited := map[string]bool{}
	var order []string
	queue := []string{changedKey}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, dep := range g.dependents[cur] {
			if !visited[dep] {
				visited[dep] = true
				order = append(order, dep)
				queue = append(queue, dep)
			}
		}
	}
	return order
}

func (e *Engine) recompute(changedKey string) {
	for _, key := range e.graph.topoOrder(changedKey) {
		field := e.model[key]
		if field.Compute != nil {
			e.state[key] = field.Compute(e.state)
		}
	}
}

// recomputeAll re-evaluates all computed fields in a safe order.
// Used after Reset() and during engine initialisation.
func (e *Engine) recomputeAll() {
	// simple multi-pass: repeat until stable (handles chains of depth >1)
	const maxPasses = 10
	for pass := 0; pass < maxPasses; pass++ {
		changed := false
		for _, key := range e.graph.computedKeys {
			field := e.model[key]
			prev := e.state[key]
			next := field.Compute(e.state)
			e.state[key] = next
			if fmt.Sprintf("%v", prev) != fmt.Sprintf("%v", next) {
				changed = true
			}
		}
		if !changed {
			break
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Cycle detection
// ─────────────────────────────────────────────────────────────────────────────

// detectCycles returns an error if the model's dependency graph contains a cycle.
// A cycle would cause recompute() to loop forever.
func detectCycles(model FormModel) error {
	// build adjacency: field → its dependents
	adj := make(map[string][]string)
	for key, field := range model {
		for _, dep := range field.DependsOn {
			adj[dep] = append(adj[dep], key)
		}
	}

	// DFS-based cycle detection
	const (
		white = 0 // unvisited
		grey  = 1 // in current path
		black = 2 // done
	)
	color := make(map[string]int)
	var path []string

	var dfs func(node string) error
	dfs = func(node string) error {
		color[node] = grey
		path = append(path, node)
		for _, next := range adj[node] {
			if color[next] == grey {
				// found cycle — find the loop in path
				start := 0
				for start < len(path) && path[start] != next {
					start++
				}
				cycle := append(path[start:], next)
				return fmt.Errorf("dependency cycle detected: %s", strings.Join(cycle, " → "))
			}
			if color[next] == white {
				if err := dfs(next); err != nil {
					return err
				}
			}
		}
		path = path[:len(path)-1]
		color[node] = black
		return nil
	}

	for key := range model {
		if color[key] == white {
			if err := dfs(key); err != nil {
				return err
			}
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x) == ""
	case int:
		return x == 0
	case float64:
		return x == 0
	case bool:
		return false // false is a valid checkbox value
	default:
		return false
	}
}

func labelOf(f Field) string {
	if f.Label != "" {
		return f.Label
	}
	return f.Name
}

// ValidationError is a structured multi-field validation result.
// Implements the error interface so it can be returned from Submit handlers.
type ValidationError struct {
	Fields map[string]error
}

func (e *ValidationError) Error() string {
	msgs := make([]string, 0, len(e.Fields))
	for k, err := range e.Fields {
		msgs = append(msgs, fmt.Sprintf("%s: %s", k, err.Error()))
	}
	return "validation failed: " + strings.Join(msgs, "; ")
}

// ErrValidation wraps a field error map into a *ValidationError.
// Returns nil if the map is empty (no errors).
func ErrValidation(errs map[string]error) error {
	if len(errs) == 0 {
		return nil
	}
	return &ValidationError{Fields: errs}
}

// AsValidationError extracts a *ValidationError from an error.
func AsValidationError(err error) (*ValidationError, bool) {
	var ve *ValidationError
	return ve, errors.As(err, &ve)
}
