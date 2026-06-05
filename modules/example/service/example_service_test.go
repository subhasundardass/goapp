package service_test

import (
	"context"
	"testing"
)

// TestCreate validates the Create business rules.
func TestCreate(t *testing.T) {
	ctx := context.Background()
	_ = ctx

	t.Run("requires name", func(t *testing.T) {
		// TODO: wire up service with a mock repo and assert ValidateRequired fires
	})

	t.Run("name too short", func(t *testing.T) {
		// TODO: assert ValidateMinLength fires for single-char names
	})
}
