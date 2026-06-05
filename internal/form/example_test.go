package form

// example.go — a complete, working example of the form engine.
//
// Run with: go test -v -run TestWalkthrough ./internal/form/
//
// This shows every engine feature:
//   - source fields with validation
//   - computed fields with dependencies
//   - hidden fields (internal use only)
//   - default values
//   - Apply (batch set from HTTP POST)
//   - Reset
//   - Validate + ValidationError
//   - UI projection (ordered, with errors)
//   - VisibleWhen (conditional field visibility)
//   - FormatFn (money formatting)

import (
	"errors"
	"fmt"
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// Invoice line-item schema
// ─────────────────────────────────────────────────────────────────────────────

func invoiceSchema() FormSchema {
	return FormSchema{
		{
			Name:     "qty",
			Label:    "Quantity",
			Type:     Number,
			Required: true,
			Access:   Writable,
			Order:    1,
			Validate: func(v any) error {
				if v == nil {
					return errors.New("qty is required")
				}
				q, ok := v.(float64)
				if !ok || q <= 0 {
					return errors.New("qty must be a positive number")
				}
				return nil
			},
		},
		{
			Name:     "price",
			Label:    "Unit Price",
			Type:     Number,
			Required: true,
			Access:   Writable,
			Order:    2,
			Format: func(v any) string {
				if v == nil {
					return "0.00"
				}
				return fmt.Sprintf("%.2f", v.(float64))
			},
			Validate: func(v any) error {
				if v == nil {
					return errors.New("price is required")
				}
				p, ok := v.(float64)
				if !ok || p < 0 {
					return errors.New("price must be zero or positive")
				}
				return nil
			},
		},
		{
			Name:      "discount",
			Label:     "Discount %",
			Type:      Number,
			Access:    Writable,
			Default:   0.0,
			Order:     3,
			// VisibleWhen: only show discount field when price > 100
			VisibleWhen: func(s FormState) bool {
				p, _ := s["price"].(float64)
				return p > 100
			},
			Validate: func(v any) error {
				if v == nil {
					return nil
				}
				d, ok := v.(float64)
				if !ok || d < 0 || d > 100 {
					return errors.New("discount must be between 0 and 100")
				}
				return nil
			},
		},
		{
			Name:      "subtotal",
			Label:     "Subtotal",
			Type:      Number,
			Access:    Readable, // computed — UI cannot write
			Order:     4,
			DependsOn: []string{"qty", "price", "discount"},
			Compute: func(s FormState) any {
				q, ok1 := s["qty"].(float64)
				p, ok2 := s["price"].(float64)
				if !ok1 || !ok2 {
					return nil
				}
				d, _ := s["discount"].(float64)
				return q * p * (1 - d/100)
			},
			Format: func(v any) string {
				if v == nil {
					return "0.00"
				}
				return fmt.Sprintf("%.2f", v.(float64))
			},
		},
		{
			Name:      "tax",
			Label:     "Tax (18%)",
			Type:      Number,
			Access:    Readable,
			Order:     5,
			DependsOn: []string{"subtotal"},
			Compute: func(s FormState) any {
				sub, ok := s["subtotal"].(float64)
				if !ok {
					return nil
				}
				return sub * 0.18
			},
			Format: func(v any) string {
				if v == nil {
					return "0.00"
				}
				return fmt.Sprintf("%.2f", v.(float64))
			},
		},
		{
			Name:      "total",
			Label:     "Total",
			Type:      Number,
			Access:    Readable,
			Order:     6,
			DependsOn: []string{"subtotal", "tax"},
			Compute: func(s FormState) any {
				sub, ok1 := s["subtotal"].(float64)
				tax, ok2 := s["tax"].(float64)
				if !ok1 || !ok2 {
					return nil
				}
				return sub + tax
			},
			Format: func(v any) string {
				if v == nil {
					return "0.00"
				}
				return fmt.Sprintf("₹ %.2f", v.(float64))
			},
		},
		{
			Name:    "note",
			Label:   "Note",
			Type:    TextArea,
			Access:  Writable,
			Order:   7,
		},
		{
			Name:   "_internal_ref",
			Label:  "",
			Type:   String,
			Access: Hidden, // never sent to UI
			Order:  99,
		},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Tests
// ─────────────────────────────────────────────────────────────────────────────

func TestComputedChain(t *testing.T) {
	e, err := NewEngineFromSchema(invoiceSchema())
	if err != nil {
		t.Fatal(err)
	}

	_ = e.SetField("qty", 4.0)
	_ = e.SetField("price", 250.0)

	subtotal, _ := e.GetField("subtotal")
	tax, _ := e.GetField("tax")
	total, _ := e.GetField("total")

	if subtotal != 1000.0 {
		t.Errorf("subtotal: want 1000.0 got %v", subtotal)
	}
	if tax != 180.0 {
		t.Errorf("tax: want 180.0 got %v", tax)
	}
	if total != 1180.0 {
		t.Errorf("total: want 1180.0 got %v", total)
	}
}

func TestDiscount(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	_ = e.SetField("qty", 2.0)
	_ = e.SetField("price", 200.0)
	_ = e.SetField("discount", 10.0) // 10% off

	subtotal, _ := e.GetField("subtotal")
	if subtotal != 360.0 { // 2*200*(1-0.10)
		t.Errorf("subtotal with discount: want 360.0 got %v", subtotal)
	}
}

func TestDefaultApplied(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	discount, _ := e.GetField("discount")
	if discount != 0.0 {
		t.Errorf("default discount: want 0.0 got %v", discount)
	}
}

func TestValidationRequiredMissing(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	errs := e.Validate()
	if _, ok := errs["qty"]; !ok {
		t.Error("expected validation error for missing qty")
	}
	if _, ok := errs["price"]; !ok {
		t.Error("expected validation error for missing price")
	}
}

func TestValidationPasses(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	_ = e.SetField("qty", 1.0)
	_ = e.SetField("price", 100.0)
	if !e.IsValid() {
		t.Errorf("expected valid form, got errors: %v", e.Validate())
	}
}

func TestApply(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	applyErrs := e.Apply(map[string]any{
		"qty":      3.0,
		"price":    50.0,
		"subtotal": 999.0, // computed — should be ignored
		"unknown":  "x",   // unknown — should be ignored
	})
	if len(applyErrs) > 0 {
		t.Errorf("unexpected apply errors: %v", applyErrs)
	}
	sub, _ := e.GetField("subtotal")
	if sub != 150.0 { // 3 * 50
		t.Errorf("subtotal after Apply: want 150.0 got %v", sub)
	}
}

func TestReset(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	_ = e.SetField("qty", 5.0)
	_ = e.SetField("price", 100.0)
	e.Reset()
	qty, _ := e.GetField("qty")
	if qty != nil {
		t.Errorf("qty after reset: want nil got %v", qty)
	}
}

func TestProjectOrdered(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	_ = e.SetField("qty", 1.0)
	_ = e.SetField("price", 100.0)

	fields := Project(e.model, e.state)

	// _internal_ref has Access=Hidden — must not appear
	for _, f := range fields {
		if f.Key == "_internal_ref" {
			t.Error("hidden field _internal_ref should not appear in projection")
		}
	}

	// Fields should be in Order sequence: qty=1, price=2, discount=3 ...
	if len(fields) < 2 {
		t.Fatal("expected at least 2 projected fields")
	}
	if fields[0].Key != "qty" {
		t.Errorf("first field should be qty, got %q", fields[0].Key)
	}
	if fields[1].Key != "price" {
		t.Errorf("second field should be price, got %q", fields[1].Key)
	}
}

func TestVisibleWhen(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())

	// price ≤ 100 → discount should not be visible
	_ = e.SetField("price", 50.0)
	fields := ProjectMap(e.model, e.state)
	if fields["discount"].Visible {
		t.Error("discount should be hidden when price <= 100")
	}

	// price > 100 → discount should be visible
	_ = e.SetField("price", 200.0)
	fields = ProjectMap(e.model, e.state)
	if !fields["discount"].Visible {
		t.Error("discount should be visible when price > 100")
	}
}

func TestProjectWithErrors(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	errs := e.Validate()
	fields := ProjectWithErrors(e.model, e.state, errs)

	var qtyField UIField
	for _, f := range fields {
		if f.Key == "qty" {
			qtyField = f
		}
	}
	if qtyField.Error == "" {
		t.Error("expected error on qty field in projection")
	}
}

func TestCycleDetection(t *testing.T) {
	cyclic := FormModel{
		"a": {Name: "a", Type: Number, DependsOn: []string{"b"}, Compute: func(s FormState) any { return s["b"] }},
		"b": {Name: "b", Type: Number, DependsOn: []string{"a"}, Compute: func(s FormState) any { return s["a"] }},
	}
	_, err := NewEngine(cyclic)
	if err == nil {
		t.Error("expected cycle detection error, got nil")
	}
}

func TestErrValidation(t *testing.T) {
	e, _ := NewEngineFromSchema(invoiceSchema())
	errs := e.Validate()
	err := ErrValidation(errs)

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Error("expected *ValidationError")
	}
	if ve.Fields["qty"] == nil {
		t.Error("expected qty error in ValidationError.Fields")
	}
}

func TestWalkthrough(t *testing.T) {
	e, err := NewEngineFromSchema(invoiceSchema())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("\n── Step 1: empty form ─────────────────────────────")
	for _, f := range Project(e.model, e.state) {
		fmt.Printf("  %-15s type=%-10s writable=%-5v visible=%-5v value=%q\n",
			f.Key, f.Meta.Type, f.Writable, f.Visible, f.Value)
	}

	fmt.Println("\n── Step 2: apply user input ───────────────────────")
	e.Apply(map[string]any{
		"qty":   4.0,
		"price": 250.0,
		"note":  "rush order",
	})

	fmt.Println("\n── Step 3: validate ───────────────────────────────")
	errs := e.Validate()
	fmt.Printf("  errors: %d\n", len(errs))

	fmt.Println("\n── Step 4: projection with computed values ────────")
	for _, f := range ProjectWithErrors(e.model, e.state, errs) {
		fmt.Printf("  %-15s writable=%-5v visible=%-5v value=%-12q error=%q\n",
			f.Key, f.Writable, f.Visible, f.Value, f.Error)
	}

	total, _ := e.GetField("total")
	fmt.Printf("\n  TOTAL: %v\n\n", total)
}
