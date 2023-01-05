package command_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jbaikge/boneless/internal/classes/app/command"
	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/logger"
	"github.com/zeebo/assert"
)

func TestAddClass(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &mockAddRepository{
		Classes: make(map[string]class.Class),
	}
	cmd := command.AddClass{
		Class: class.NewClass(t.Name(), "", nil),
	}
	handler := command.NewAddClassHandler(repo, logger.Default)

	// Ensure class goes in
	assert.NoError(t, handler.Handle(ctx, cmd))

	// Make sure the repository error bubbles up through the handler
	assert.Error(t, handler.Handle(ctx, cmd))
}

type mockAddRepository struct {
	Classes map[string]class.Class
}

func (r *mockAddRepository) AddClass(ctx context.Context, class *class.Class) error {
	_, ok := r.Classes[class.ID()]
	if ok {
		return fmt.Errorf("class already exists with ID: %s", class.ID())
	}
	r.Classes[class.ID()] = *class
	return nil
}

func (r *mockAddRepository) GetClass(ctx context.Context, id string) (*class.Class, error) {
	return nil, nil
}

func (r *mockAddRepository) UpdateClass(ctx context.Context, class *class.Class) error {
	return nil
}
