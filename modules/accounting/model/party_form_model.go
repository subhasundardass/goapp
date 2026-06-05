package model

import "goapp/internal/form"

func PartyFormModel() form.FormModel {
	return form.FormModel{
		"display_name": {
			Name:   "display_name",
			Type:   form.String,
			Label:  "Party Name",
			Access: form.Writable,
			Order:  1,
		},

		"city": {
			Name:   "city",
			Type:   form.String,
			Label:  "City",
			Access: form.Writable,
			Order:  2,
		},

		"party_type": {
			Name:   "party_type",
			Type:   form.Select,
			Label:  "Party Type",
			Access: form.Writable,
			Order:  3,
			Options: []form.Option{
				{Value: "buyer", Label: "Buyer"},
				{Value: "supplier", Label: "Supplier"},
				{Value: "both", Label: "Both"},
			},
			Default: "buyer",
		},
		"ledgers": {
			Name:    "ledgers",
			Type:    form.Select,
			Label:   "Ledger",
			Access:  form.Writable,
			Order:   4,
			Options: []form.Option{},
		},
		"countries": {
			Name:    "countries",
			Type:    form.Select,
			Label:   "Country",
			Access:  form.Writable,
			Order:   4,
			Options: []form.Option{},
		},
	}
}
