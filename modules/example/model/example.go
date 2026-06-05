package model

import (
	"fmt"
	"goapp/internal/form"
)

func ExampleFormModel() form.FormModel {

	return form.FormModel{

		"itemName": {
			Name:     "itemName",
			Type:     form.String,
			Label:    "Item Name",
			Access:   form.Writable,
			Order:    1,
			Required: true,
		},

		"qty": {
			Name:   "qty",
			Type:   form.Number,
			Label:  "Quantity",
			Access: form.Writable,
			Order:  2,
			Validate: func(v any) error {
				q, ok := v.(float64)
				if !ok || q <= 0 {
					return fmt.Errorf("quantity must be greater than zero")
				}
				return nil
			},
		},

		"rate": {
			Name:   "rate",
			Type:   form.Number,
			Label:  "Rate",
			Access: form.Writable,
			Order:  3,
			Validate: func(v any) error {
				r, ok := v.(float64)
				if !ok || r < 0 {
					return fmt.Errorf("rate cannot be negative")
				}
				return nil
			},
		},

		"amount": {
			Name:      "amount",
			Type:      form.Number,
			Label:     "Amount",
			Access:    form.Readable,
			Order:     4,
			DependsOn: []string{"qty", "rate"},
			Compute: func(s form.FormState) any {
				qty, _ := s["qty"].(float64)
				rate, _ := s["rate"].(float64)
				return qty * rate
			},
		},

		"taxPercent": {
			Name:   "taxPercent",
			Type:   form.Number,
			Label:  "Tax %",
			Access: form.Writable,
			Order:  5,
		},

		"taxAmount": {
			Name:      "taxAmount",
			Type:      form.Number,
			Label:     "Tax Amount",
			Access:    form.Readable,
			Order:     6,
			DependsOn: []string{"amount", "taxPercent"},
			Compute: func(s form.FormState) any {
				amount, _ := s["amount"].(float64)
				tax, _ := s["tax_percent"].(float64)
				return amount * tax / 100
			},
		},

		"grandTotal": {
			Name:      "grandTotal",
			Type:      form.Number,
			Label:     "Grand Total",
			Access:    form.Readable,
			Order:     7,
			DependsOn: []string{"amount", "taxAmount"},
			Compute: func(s form.FormState) any {
				amount, _ := s["amount"].(float64)
				tax, _ := s["taxAmount"].(float64)
				return amount + tax
			},
		},

		"paymentMode": {
			Name:   "paymentMode",
			Type:   form.Select,
			Label:  "Payment Mode",
			Access: form.Writable,
			Order:  8,
			Options: []form.Option{
				{Value: "cash", Label: "Cash"},
				{Value: "upi", Label: "UPI"},
				{Value: "card", Label: "Card"},
			},
		},
	}
}
