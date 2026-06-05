package formui

// FieldProps carries presentation-level options for widget components.
// These are separate from the semantic FieldMeta so the form engine stays
// decoupled from rendering details.
//
// FIX: original had only Class, WrapperClass, Placeholder, TabIndex, Disabled.
// Added LabelClass, HintText, Rows (textarea), Step/Min/Max (number),
// AutoComplete, AutoFocus, and ID override so every HTML attribute a
// template author might need is available without forking the component.
type FieldProps struct {
	// Wrapper / layout
	WrapperClass string // class on the outer <div>
	LabelClass   string // class on the <label>

	// Input styling
	Class string // class on the <input> / <select> / <textarea>

	// Behaviour overrides
	Placeholder  string
	TabIndex     int
	Disabled     bool
	AutoComplete string // "on" | "off" | "email" | etc.
	AutoFocus    bool
	ID           string // overrides field.Key as the HTML id attribute

	// Hint / helper text shown below the input
	HintText string

	// Textarea-specific
	Rows int // defaults to 3 if zero

	// Datastart
	Endpoint string

	// Number-specific
	Step string // "1", "0.01", "any"
	Min  string // HTML min attribute
	Max  string // HTML max attribute
}

// DefaultProps returns a FieldProps with sensible defaults applied.
func DefaultProps() FieldProps {
	return FieldProps{
		Rows: 3,
		Step: "any",
	}
}

// Merge returns a new FieldProps where any non-zero value in override
// replaces the corresponding value in base. Useful for component defaults.
func (base FieldProps) Merge(override FieldProps) FieldProps {
	if override.WrapperClass != "" {
		base.WrapperClass = override.WrapperClass
	}
	if override.LabelClass != "" {
		base.LabelClass = override.LabelClass
	}
	if override.Class != "" {
		base.Class = override.Class
	}
	if override.Placeholder != "" {
		base.Placeholder = override.Placeholder
	}
	if override.TabIndex != 0 {
		base.TabIndex = override.TabIndex
	}
	if override.Disabled {
		base.Disabled = true
	}
	if override.AutoComplete != "" {
		base.AutoComplete = override.AutoComplete
	}
	if override.AutoFocus {
		base.AutoFocus = true
	}
	if override.ID != "" {
		base.ID = override.ID
	}
	if override.HintText != "" {
		base.HintText = override.HintText
	}
	if override.Rows != 0 {
		base.Rows = override.Rows
	}
	if override.Step != "" {
		base.Step = override.Step
	}
	if override.Min != "" {
		base.Min = override.Min
	}
	if override.Max != "" {
		base.Max = override.Max
	}
	return base
}
