package service

import (
	"log/slog"

	"goapp/internal/core"
	"goapp/modules/accounting/repository"
)

// Services groups all services for the Accounting module.
// Add new service structs here as the module grows:
//
//	type Services struct {
//	    Accounting  *AccountingService
//	    Journal *JournalService
//	    Report  *ReportService
//	}
type Services struct {
	Accounting *AccountingService
}

// NewServices wires all services with their repository dependencies.
func NewServices(repos *repository.Repositories, logger *slog.Logger, store *core.Store) *Services {
	return &Services{
		Accounting: NewAccountingService(repos, logger),
	}
}
