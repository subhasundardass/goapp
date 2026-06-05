package form

// ─────────────────────────────────────────────────────────────────────────────
// Types
// ─────────────────────────────────────────────────────────────────────────────

// FieldType is the semantic data type of a field.
type FieldType string

// Access controls whether a field is visible and editable in the UI.
type Access int

const (
	Hidden   Access = iota // not sent to the UI at all
	Readable               // sent to UI, read-only
	Writable               // sent to UI, user can edit
)

const (
	String   FieldType = "string"
	Number   FieldType = "number"
	Bool     FieldType = "bool" // use Bool — Checkbox is an HTML concern, not a schema concern
	Select   FieldType = "select"
	Radio    FieldType = "radio"
	Date     FieldType = "date"
	DateTime FieldType = "datetime-local" // matches HTML input type="datetime-local"
	Email    FieldType = "email"
	TextArea FieldType = "textarea"
	File     FieldType = "file"
)

// Option is a single choice in a Select or Radio field.
type Option struct {
	Value    string
	Label    string
	Disabled bool // renders as <option disabled> or a greyed-out radio
}

// FormState is the live map of field key → value held by the engine.
type FormState map[string]any

// ComputeFn derives a field value from the current form state.
// Return nil if the value cannot be computed yet (missing dependencies).
type ComputeFn func(state FormState) any

// ValidateFn validates a single field value against business rules.
// Return a non-nil error to signal a validation failure.
type ValidateFn func(value any) error

// FormatFn converts a raw any value to its display string.
// Used by the UI projection layer to format numbers, dates, etc.
type FormatFn func(value any) string

// ─────────────────────────────────────────────────────────────────────────────
// Field
// ─────────────────────────────────────────────────────────────────────────────

// Field describes one field in a form schema.
//
//   - Source fields: Compute is nil — values are written by the user/caller.
//   - Computed fields: Compute is set — values are owned by the engine.
//   - Both kinds may have a Validate function.
//   - Format is optional; defaults to fmt.Sprintf("%v", value).
//
// Note: Placeholder is intentionally omitted here — it is a rendering concern
// and belongs in FieldProps (formui package), not the schema.
type Field struct {
	// Core identity
	Name  string
	Label string
	Type  FieldType

	// Value behaviour
	Default   any        // applied when the field has no value
	Required  bool       // shorthand; engine also checks ValidateFn
	DependsOn []string   // keys this field reads in Compute
	Compute   ComputeFn  // nil = source field (user-editable)
	Validate  ValidateFn // nil = no validation
	Format    FormatFn   // nil = fmt.Sprintf("%v", value)

	// UI behaviour
	Access      Access
	Options     []Option             // for Select / Radio fields
	VisibleWhen func(FormState) bool // nil = always visible
	Group       string               // visual grouping in rendered form
	Order       int                  // render order (lower = earlier)
}

// ─────────────────────────────────────────────────────────────────────────────
// FormModel  (map — O(1) lookup, non-deterministic order)
// FormSchema (slice — preserves declaration order for rendering)
// ─────────────────────────────────────────────────────────────────────────────

// FormModel is the schema used internally by the engine.
type FormModel map[string]Field

// FormSchema is an ordered field list.
// Prefer this when declaring forms so field order is preserved.
type FormSchema []Field

// ToModel converts a FormSchema into a FormModel for engine use.
func (s FormSchema) ToModel() FormModel {
	m := make(FormModel, len(s))
	for _, f := range s {
		m[f.Name] = f
	}
	return m
}

// OrderedKeys returns field keys sorted by their Order value,
// preserving declaration order for ties.
func (m FormModel) OrderedKeys() []string {
	type entry struct {
		key   string
		order int
	}
	entries := make([]entry, 0, len(m))
	for k, f := range m {
		entries = append(entries, entry{k, f.Order})
	}
	// stable insertion sort
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].order < entries[j-1].order; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.key
	}
	return keys
}
