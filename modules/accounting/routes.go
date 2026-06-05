package accounting

import (
	"goapp/internal/core"
	"goapp/modules/accounting/handler"
)

// Routes declares all HTTP routes for the Accounting module.
// Add new routes here as the module grows — no need to touch module.go.
func Routes(h *handler.Handlers) []core.Route {
	return []core.Route{
		{Method: "GET", Path: "/", Handler: h.Accounting.LedgerList},

		// ── Group ── specific paths BEFORE /:id
		{Method: "GET", Path: "/group", Handler: h.Accounting.GroupList},

		// ── Ledger ── specific paths BEFORE /:id
		{Method: "GET", Path: "/ledger", Handler: h.Accounting.LedgerList},
		// {Method: "GET", Path: "/ledger/:id/edit", Handler: h.Accounting.LedgerEdit},

		// ── Party ── specific paths BEFORE /:id
		{Method: "GET", Path: "/party", Handler: h.Accounting.PartyList},
		{Method: "GET", Path: "/party/new", Handler: h.Accounting.PartyNew},
		// {Method: "POST", Path: "/party/new", Handler: h.Accounting.PartyCreate},

		// ── Generic (keep at bottom) ──
		// {Method: "POST", Path: "/", Handler: h.Accounting.Create},
		// {Method: "GET", Path: "/:id", Handler: h.Accounting.Show},
		// {Method: "PUT", Path: "/:id", Handler: h.Accounting.Update},
		// {Method: "DELETE", Path: "/:id", Handler: h.Accounting.Delete},
	}
}
