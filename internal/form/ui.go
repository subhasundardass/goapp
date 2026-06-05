package form

import "fmt"

// ─────────────────────────────────────────────────────────────────────────────
// UI projection types
// ─────────────────────────────────────────────────────────────────────────────

// UIField is the read-only view of a field sent to the template layer.
// It contains everything a templ component needs to render the field.
type UIField struct {
	Key      string
	Value    string    // formatted display value
	RawValue any       // original unformatted value (for computed display etc.)
	Writable bool      // false = render as disabled / readonly
	Visible  bool      // false = field should be hidden (VisibleWhen returned false)
	Error    string    // non-empty = show error message
	Meta     FieldMeta // static metadata from the schema
}

// FieldMeta carries the static schema metadata needed by UI components.
// Placeholder is intentionally absent — it is a rendering concern and
// belongs in FieldProps (formui package), not the schema.
type FieldMeta struct {
	Name     string
	Label    string
	Type     FieldType
	Required bool
	Options  []Option
	Group    string
	Order    int
}

// ─────────────────────────────────────────────────────────────────────────────
// Projection
// ─────────────────────────────────────────────────────────────────────────────

// Project converts a FormModel + FormState into an ordered slice of UIFields.
// Templates that range over a map render fields in random order every
// request — use the slice form for stable rendering.
func Project(model FormModel, state FormState) []UIField {
	keys := model.OrderedKeys()
	out := make([]UIField, 0, len(keys))

	for _, key := range keys {
		field := model[key]
		if field.Access == Hidden {
			continue
		}

		visible := true
		if field.VisibleWhen != nil {
			visible = field.VisibleWhen(state)
		}

		out = append(out, UIField{
			Key:      key,
			Value:    formatVal(field, state[key]),
			RawValue: state[key],
			Writable: field.Access == Writable && field.Compute == nil,
			Visible:  visible,
			Meta: FieldMeta{
				Name:     field.Name,
				Label:    field.Label,
				Type:     field.Type,
				Required: field.Required,
				Options:  field.Options,
				Group:    field.Group,
				Order:    field.Order,
			},
		})
	}

	return out
}

// ProjectMap returns the same data as Project but keyed by field name.
// Useful for template lookups: ui["email"].Value.
func ProjectMap(model FormModel, state FormState) map[string]UIField {
	fields := Project(model, state)
	m := make(map[string]UIField, len(fields))
	for _, f := range fields {
		m[f.Key] = f
	}
	return m
}

// ProjectWithErrors attaches validation errors to projected fields.
func ProjectWithErrors(model FormModel, state FormState, errs map[string]error) []UIField {
	fields := Project(model, state)
	for i, f := range fields {
		if err, ok := errs[f.Key]; ok && err != nil {
			fields[i].Error = err.Error()
		}
	}
	return fields
}

// ProjectMapWithErrors is the map form of ProjectWithErrors.
func ProjectMapWithErrors(model FormModel, state FormState, errs map[string]error) map[string]UIField {
	fields := ProjectWithErrors(model, state, errs)
	m := make(map[string]UIField, len(fields))
	for _, f := range fields {
		m[f.Key] = f
	}
	return m
}

// ProjectByGroup returns fields grouped by their Group value.
// Fields with no Group are placed under "default".
// Useful for multi-section form layouts.
func ProjectByGroup(model FormModel, state FormState) map[string][]UIField {
	fields := Project(model, state)
	groups := map[string][]UIField{}
	for _, f := range fields {
		g := f.Meta.Group
		if g == "" {
			g = "default"
		}
		groups[g] = append(groups[g], f)
	}
	return groups
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func formatVal(field Field, value any) string {
	if value == nil {
		return ""
	}
	if field.Format != nil {
		return field.Format(value)
	}
	return fmt.Sprintf("%v", value)
}
