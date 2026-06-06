package repository

import (
	"context"
	"fmt"

	"goapp/ent"
	"goapp/ent/ledger"
	"goapp/ent/state"
	"goapp/internal/core"
	"goapp/modules/accounting/model"
)

// AccountingRepository handles all DB operations for Accounting.
// TODO: replace 'any' with generated ent types after: go generate ./ent/...
type AccountingRepository struct {
	db *ent.Client
}

func NewAccountingRepository(db *ent.Client) *AccountingRepository {
	return &AccountingRepository{db: db}
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (r *AccountingRepository) List(ctx context.Context) ([]any, error) {
	// TODO: return r.db.Accounting.Query().All(ctx)
	return []any{}, nil
}

func (r *AccountingRepository) ListPaginated(ctx context.Context, offset, limit int) ([]any, int, error) {
	// TODO:
	// q     := r.db.Accounting.Query()
	// total, _ := q.Count(ctx)
	// items, _ := q.Limit(limit).Offset(offset).All(ctx)
	return []any{}, 0, nil
}

// Country List
func (r *AccountingRepository) CountryList(ctx context.Context) ([]*ent.Country, error) {
	return r.db.Country.
		Query().
		All(ctx)
}

// State List
func (r *AccountingRepository) StateList(ctx context.Context) ([]*ent.State, error) {
	return r.db.State.
		Query().
		Order(ent.Asc(state.FieldName)).
		All(ctx)
}

func (r *AccountingRepository) LedgerList(ctx context.Context) ([]*ent.Ledger, error) {
	return r.db.Ledger.
		Query().
		Where(ledger.StatusEQ(1)).
		All(ctx)
}

func (r *AccountingRepository) LedgerListPaginated(ctx context.Context, offset, limit int) ([]*ent.Ledger, int, error) {
	query := r.db.Ledger.Query()

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := r.db.Ledger.
		Query().
		WithGroup().
		Order(ent.Desc(ledger.FieldID)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *AccountingRepository) GroupListPaginated(ctx context.Context, offset, limit int) ([]*ent.Acct_Group, int, error) {
	query := r.db.Acct_Group.Query()

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := r.db.Acct_Group.
		Query().
		Order(ent.Desc(ledger.FieldID)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *AccountingRepository) PartyListPaginated(ctx context.Context, offset, limit int) ([]*ent.PartyMaster, int, error) {
	query := r.db.PartyMaster.Query()

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	items, err := r.db.PartyMaster.
		Query().
		// Order(ent.Desc(ledger.FieldID)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *AccountingRepository) FindByID(ctx context.Context, id int) (any, error) {
	// TODO: return r.db.Accounting.Get(ctx, id)
	return nil, nil
}

func (r *AccountingRepository) FindLedgerByID(
	ctx context.Context,
	id int,
) (*ent.Ledger, error) {

	ledger, err := r.db.Ledger.
		Query().
		Where(ledger.IDEQ(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("ledger not found: %d", id)
		}
		return nil, err
	}

	return ledger, nil
}

func (r *AccountingRepository) CreateParty(ctx context.Context, input model.PartyInput) (*ent.PartyMaster, error) {

	q := r.db.PartyMaster.Create()

	q.SetDisplayName(input.DisplayName)
	q.SetType(input.Type)

	q.SetLedgerID(core.ToInt(input.LedgerID))

	if input.LegalName != nil {
		q.SetLegalName(*input.LegalName)
	}

	if input.GSTNo != nil {
		q.SetGstNo(*input.GSTNo)
	}

	if input.PANNo != nil {
		q.SetPanNo(*input.PANNo)
	}

	if input.ContactPerson != nil {
		q.SetContactPerson(*input.ContactPerson)
	}

	if input.Mobile != nil {
		q.SetMobile(*input.Mobile)
	}

	if input.Phone != nil {
		q.SetPhone(*input.Phone)
	}

	if input.Email != nil {
		q.SetEmail(*input.Email)
	}

	if input.Website != nil {
		q.SetWebsite(*input.Website)
	}

	q.SetCreditLimit(core.ToFloat(input.CreditLimit))
	q.SetCreditDays(input.CreditDays)
	q.SetOpeningBalance(core.ToFloat(input.OpeningBalance))

	if input.Address != nil {
		q.SetAddress(*input.Address)
	}

	if input.City != nil {
		q.SetCity(*input.City)
	}

	if input.State != nil {
		q.SetState(*input.State)
	}

	q.SetCountry(input.Country)

	if input.Pincode != nil {
		q.SetPincode(*input.Pincode)
	}

	return q.Save(ctx)
}

func (r *AccountingRepository) Update(ctx context.Context, id int, name, description string) (any, error) {
	// TODO: return r.db.Accounting.UpdateOneID(id).SetName(name).SetDescription(description).Save(ctx)
	return nil, nil
}

func (r *AccountingRepository) Delete(ctx context.Context, id int) error {
	// TODO: return r.db.Accounting.DeleteOneID(id).Exec(ctx)
	return nil
}

// ── Custom queries ────────────────────────────────────────────────────────────

func (r *AccountingRepository) SearchByName(ctx context.Context, params core.PaginationParams) ([]any, int, error) {
	// TODO:
	// q := r.db.Accounting.Query()
	// if params.Search != "" {
	//     q = q.Where(accounting.NameContainsFold(params.Search))
	// }
	return []any{}, 0, nil
}
