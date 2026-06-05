package service

import (
	"context"
	"log/slog"

	"goapp/internal/core"
	"goapp/modules/example/repository"
)

// Service holds the business logic for the Example module.
// Replace 'any' return types once the ent schema is generated.
type Service struct {
	core.BaseService
	repo  *repository.Repo
	store *core.Store
}

func NewService(repo *repository.Repo, logger *slog.Logger, store *core.Store) *Service {
	return &Service{
		BaseService: core.BaseService{Logger: logger},
		repo:        repo,
		store:       store,
	}
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (s *Service) List(ctx context.Context) ([]any, error) {
	return s.repo.List(ctx)
}

// ListPaginated returns a paginated result set.
func (s *Service) ListPaginated(ctx context.Context, params core.PaginationParams) (*core.PaginatedResult[any], error) {
	return core.Paginate(ctx, params, s.repo.ListPaginated)
}

func (s *Service) FindByID(ctx context.Context, id int) (any, error) {
	item, err := s.repo.FindByID(ctx, id)
	return item, s.WrapError("FindByID", err)
}

// ── Commands ──────────────────────────────────────────────────────────────────

func (s *Service) Create(ctx context.Context, name, description string) (any, error) {
	if err := s.ValidateRequired(map[string]string{
		"name": name,
	}); err != nil {
		return nil, err
	}
	if err := s.ValidateMinLength("name", name, 2); err != nil {
		return nil, err
	}
	if err := s.ValidateMaxLength("name", name, 255); err != nil {
		return nil, err
	}
	item, err := s.repo.Create(ctx, name, description)
	return item, s.WrapError("Create", err)
}

func (s *Service) Update(ctx context.Context, id int, name, description string) (any, error) {
	if err := s.ValidateRequired(map[string]string{
		"name": name,
	}); err != nil {
		return nil, err
	}
	if err := s.ValidateMinLength("name", name, 2); err != nil {
		return nil, err
	}
	if err := s.ValidateMaxLength("name", name, 255); err != nil {
		return nil, err
	}
	item, err := s.repo.Update(ctx, id, name, description)
	return item, s.WrapError("Update", err)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.WrapError("Delete", s.repo.Delete(ctx, id))
}
