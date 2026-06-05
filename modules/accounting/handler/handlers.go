package handler

import "goapp/modules/accounting/service"

// Handlers groups all handlers for the Accounting module.
// Add new handler structs here as the module grows:
//
//	type Handlers struct {
//	    Accounting  *AccountingHandler
//	    Journal *JournalHandler
//	    Report  *ReportHandler
//	}
type Handlers struct {
	Accounting *AccountingHandler
}

// NewHandlers wires all handlers with their service dependencies.
func NewHandlers(svc *service.Services) *Handlers {
	return &Handlers{
		Accounting: NewAccountingHandler(svc),
	}
}
