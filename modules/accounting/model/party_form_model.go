package model

import (
	"goapp/ent/partymaster"
	"goapp/internal/form"
)

func PartyFormModel() form.FormModel {
	return form.FormModel{
		"display_name": {
			Name:   "display_name",
			Type:   form.String,
			Label:  "Party Name",
			Access: form.Writable,
			Order:  1,
		},
		"legal_name": {
			Name:   "legal_name",
			Type:   form.String,
			Label:  "Legal Name",
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
				{Value: "CUSTOMER", Label: "Customer"},
				{Value: "SUPPLIER", Label: "Supplier"},
				{Value: "BOTH", Label: "Both"},
			},
			Default: "CUSTOMER",
		},
		"ledger_id": {
			Name:    "ledgers",
			Type:    form.Select,
			Label:   "Ledger",
			Access:  form.Writable,
			Order:   4,
			Options: []form.Option{},
		},
		"gst_no": {
			Name:   "gst_no",
			Type:   form.String,
			Label:  "GST No",
			Access: form.Writable,
			Order:  5,
		},
		"pan_no": {
			Name:   "pan_no",
			Type:   form.String,
			Label:  "PAN No",
			Access: form.Writable,
			Order:  6,
		},
		"contact_person": {
			Name:   "contact_person",
			Type:   form.String,
			Label:  "Contact Person",
			Access: form.Writable,
			Order:  7,
		},
		"mobile": {
			Name:   "mobile",
			Type:   form.String,
			Label:  "Mobile",
			Access: form.Writable,
			Order:  8,
		},
		"phone": {
			Name:   "phone",
			Type:   form.String,
			Label:  "Phone",
			Access: form.Writable,
			Order:  9,
		},
		"email": {
			Name:   "email",
			Type:   form.Email,
			Label:  "Email",
			Access: form.Writable,
			Order:  10,
		},
		"website": {
			Name:   "website",
			Type:   form.String,
			Label:  "Website",
			Access: form.Writable,
			Order:  11,
		},
		"credit_limit": {
			Name:    "credit_limit",
			Type:    form.Number,
			Label:   "Credit Limit",
			Access:  form.Writable,
			Order:   12,
			Default: "0",
		},
		"credit_days": {
			Name:    "credit_days",
			Type:    form.Number,
			Label:   "Credit Days",
			Access:  form.Writable,
			Order:   13,
			Default: "0",
		},
		"opening_balance": {
			Name:    "opening_balance",
			Type:    form.Number,
			Label:   "Opening Balance",
			Access:  form.Writable,
			Order:   14,
			Default: "0",
		},
		"address": {
			Name:   "address",
			Type:   form.String,
			Label:  "Address",
			Access: form.Writable,
			Order:  15,
		},
		"city": {
			Name:   "city",
			Type:   form.String,
			Label:  "City",
			Access: form.Writable,
			Order:  16,
		},
		"states": {
			Name:    "states",
			Type:    form.Select,
			Label:   "State",
			Access:  form.Writable,
			Order:   17,
			Options: []form.Option{},
		},
		"countries": {
			Name:    "countries",
			Type:    form.Select,
			Label:   "Country",
			Access:  form.Writable,
			Order:   18,
			Options: []form.Option{},
		},
		"pincode": {
			Name:   "pincode",
			Type:   form.String,
			Label:  "Pincode",
			Access: form.Writable,
			Order:  19,
		},
	}
}

type PartyInput struct {
	LedgerID       int              `form:"ledger_id"`
	Type           partymaster.Type `form:"party_type"`
	DisplayName    string           `form:"display_name"`
	LegalName      *string          `form:"legal_name"`
	GSTNo          *string          `form:"gst_no"`
	PANNo          *string          `form:"pan_no"`
	ContactPerson  *string          `form:"contact_person"`
	Mobile         *string          `form:"mobile"`
	Phone          *string          `form:"phone"`
	Email          *string          `form:"email"`
	Website        *string          `form:"website"`
	CreditLimit    float64          `form:"credit_limit"`
	CreditDays     int              `form:"credit_days"`
	OpeningBalance float64          `form:"opening_balance"`
	Address        *string          `form:"address"`
	City           *string          `form:"city"`
	State          *string          `form:"state"`
	Country        string           `form:"country"`
	Pincode        *string          `form:"pincode"`
}
