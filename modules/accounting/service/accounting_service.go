package service

import (
	"context"
	"fmt"
	"log/slog"

	"goapp/ent"
	"goapp/internal/core"
	"goapp/modules/accounting/model"
	"goapp/modules/accounting/repository"
)

// AccountingService holds business logic for Accounting.
// Replace 'any' return types once the ent schema is generated.
type AccountingService struct {
	core.BaseService
	repos *repository.Repositories
}

func NewAccountingService(repos *repository.Repositories, logger *slog.Logger) *AccountingService {
	return &AccountingService{
		BaseService: core.BaseService{Logger: logger},
		repos:       repos,
	}
}

// ── Queries ───────────────────────────────────────────────────────────────────
func (s *AccountingService) CountryList(ctx context.Context) ([]*ent.Country, error) {
	return s.repos.Accounting.CountryList(ctx)
}
func (s *AccountingService) StateList(ctx context.Context) ([]*ent.State, error) {
	return s.repos.Accounting.StateList(ctx)
}

// func (s *AccountingService) List(ctx context.Context) ([]any, error) {
// 	return s.repos.Accounting.List(ctx)
// }

// func (s *AccountingService) ListPaginated(ctx context.Context, params core.PaginationParams) (*core.PaginatedResult[any], error) {
// 	return core.Paginate(ctx, params, s.repos.Accounting.ListPaginated)
// }

// func (s *AccountingService) FindByID(ctx context.Context, id int) (any, error) {
// 	item, err := s.repos.Accounting.FindByID(ctx, id)
// 	return item, s.WrapError("FindByID", err)
// }

func (s *AccountingService) LedgerList(ctx context.Context) ([]*ent.Ledger, error) {
	return s.repos.Accounting.LedgerList(ctx)
}

func (s *AccountingService) LedgerListPaginated(ctx context.Context, params core.PaginationParams) (*core.PaginatedResult[*ent.Ledger], error) {
	return core.Paginate(ctx, params, s.repos.Accounting.LedgerListPaginated)
}
func (s *AccountingService) GroupListPaginated(ctx context.Context, params core.PaginationParams) (*core.PaginatedResult[*ent.Acct_Group], error) {
	return core.Paginate(ctx, params, s.repos.Accounting.GroupListPaginated)
}
func (s *AccountingService) PartyListPaginated(ctx context.Context, params core.PaginationParams) (*core.PaginatedResult[*ent.PartyMaster], error) {
	return core.Paginate(ctx, params, s.repos.Accounting.PartyListPaginated)
}

func (s *AccountingService) FindLedgerByID(ctx context.Context, id int) (any, error) {
	item, err := s.repos.Accounting.FindLedgerByID(ctx, id)
	return item, s.WrapError("FindByID", err)
}

// ── Commands ──────────────────────────────────────────────────────────────────

func (s *AccountingService) CreateParty(ctx context.Context, input model.PartyInput) (any, error) {

	name := fmt.Sprint(input.DisplayName)

	if err := s.ValidateRequired(map[string]string{
		"display_name": name,
	}); err != nil {
		return nil, err
	}

	if err := s.ValidateMinLength("display_name", name, 2); err != nil {
		return nil, err
	}

	if err := s.ValidateMaxLength("display_name", name, 255); err != nil {
		return nil, err
	}

	// call repo (pass full input or extracted values)
	item, err := s.repos.Accounting.CreateParty(ctx, input)

	return item, s.WrapError("CreateParty", err)
}

func (s *AccountingService) Update(ctx context.Context, id int, name, description string) (any, error) {
	if err := s.ValidateRequired(map[string]string{"name": name}); err != nil {
		return nil, err
	}
	if err := s.ValidateMinLength("name", name, 2); err != nil {
		return nil, err
	}
	if err := s.ValidateMaxLength("name", name, 255); err != nil {
		return nil, err
	}
	item, err := s.repos.Accounting.Update(ctx, id, name, description)
	return item, s.WrapError("Update", err)
}

func (s *AccountingService) Delete(ctx context.Context, id int) error {
	return s.WrapError("Delete", s.repos.Accounting.Delete(ctx, id))
}
