package service

import (
	"context"
	"fmt"
	"log/slog"

	"goapp/internal/core"
	"goapp/modules/settings/repository"
)

// Service holds the business logic for the Settings module.
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
// Load reads all settings from DB into Store. Call at startup.
func (s *Service) Load(ctx context.Context) error {
	rows, err := s.repo.All(ctx)
	if err != nil {
		return fmt.Errorf("settings: load failed: %w", err)
	}
	for _, row := range rows {
		s.store.Set(row.Key, row.Value)
	}
	return nil
}

// Update persists a new value to DB and refreshes Store immediately.
func (s *Service) Update(ctx context.Context, key, value string) error {
	if err := s.repo.Set(ctx, key, value); err != nil {
		return fmt.Errorf("settings: update %q failed: %w", key, err)
	}
	s.store.Set(key, value)
	return nil
}

// Get reads a setting from Store (no DB hit).
func (s *Service) Get(key string) (string, bool) {
	return s.store.GetString(key)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.WrapError("Delete", s.repo.Delete(ctx, id))
}
