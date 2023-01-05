package command_test

import (
	"context"
	"testing"

	"github.com/jbaikge/boneless/internal/classes/app/command"
	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/logger"
	"github.com/zeebo/assert"
)

func TestAddClass(t *testing.T) {
	t.Parallel()

	repo := &mockAddRepository{}

	cmd := command.AddClass{
		Class: class.NewClass(t.Name(), "", nil),
	}
	handler := command.NewAddClassHandler(repo, logger.Default)
	err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

type mockAddRepository struct {
	Classes map[string]class.Class
}

func (r *mockAddRepository) AddClass(ctx context.Context, class *class.Class) error {
	return nil
}

func (r *mockAddRepository) GetClass(ctx context.Context, id string) (*class.Class, error) {
	return nil, nil
}

func (r *mockAddRepository) UpdateClass(ctx context.Context, class *class.Class) error {
	return nil
}
