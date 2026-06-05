package repository

import (
	"context"

	"goapp/ent"
	"goapp/internal/core"
)

// Repo handles all database operations for the Example module.
// TODO: replace 'any' types with generated ent types after running:
//
//	go generate ./ent/...
type Repo struct {
	db *ent.Client
}

func NewRepo(db *ent.Client) *Repo {
	return &Repo{db: db}
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (r *Repo) List(ctx context.Context) ([]any, error) {
	// TODO: return r.db.Example.Query().All(ctx)
	return []any{}, nil
}

func (r *Repo) ListPaginated(ctx context.Context, offset, limit int) ([]any, int, error) {
	// TODO: replace with real ent query
	// q := r.db.Example.Query()
	// total, err := q.Count(ctx)
	// items, err := q.Limit(limit).Offset(offset).All(ctx)
	return []any{}, 0, nil
}

func (r *Repo) FindByID(ctx context.Context, id int) (any, error) {
	// TODO: return r.db.Example.Get(ctx, id)
	return nil, nil
}

func (r *Repo) Create(ctx context.Context, name, description string) (any, error) {
	// TODO: return r.db.Example.Create().SetName(name).SetDescription(description).Save(ctx)
	return nil, nil
}

func (r *Repo) Update(ctx context.Context, id int, name, description string) (any, error) {
	// TODO: return r.db.Example.UpdateOneID(id).SetName(name).SetDescription(description).Save(ctx)
	return nil, nil
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	// TODO: return r.db.Example.DeleteOneID(id).Exec(ctx)
	return nil
}

// ── Module-specific queries (add your custom queries below) ───────────────────

// SearchByName filters records by name substring.
func (r *Repo) SearchByName(ctx context.Context, params core.PaginationParams) ([]any, int, error) {
	// TODO: implement search
	// q := r.db.Example.Query()
	// if params.Search != "" {
	//     q = q.Where(example.NameContainsFold(params.Search))
	// }
	return []any{}, 0, nil
}
