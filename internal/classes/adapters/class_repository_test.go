package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/jbaikge/boneless/internal/classes/adapters"
	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/zeebo/assert"
)

type Repository struct {
	Name       string
	Repository class.Repository
}

func TestRepository(t *testing.T) {
	t.Parallel()

	repositories := []Repository{
		{
			Name:       "Memory",
			Repository: adapters.NewClassMemoryRepository(),
		},
	}

	for i := range repositories {
		repo := repositories[i]

		t.Run(repo.Name, func(t *testing.T) {
			t.Parallel()

			t.Run("Store Retrieve", func(t *testing.T) {
				t.Parallel()
				testStoreRetrieve(t, repo.Repository)
			})

			t.Run("Add Same ID", func(t *testing.T) {
				t.Parallel()
				testAddSameID(t, repo.Repository)
			})
		})
	}
}

func testStoreRetrieve(t *testing.T, repo class.Repository) {
	t.Helper()

	ctx := context.Background()

	testClass := class.NewClass(t.Name(), "", nil)
	assert.NoError(t, repo.AddClass(ctx, testClass))

	check, err := repo.GetClass(ctx, testClass.ID())
	assert.NoError(t, err)
	assert.Equal(t, testClass.ID(), check.ID())
}

func testAddSameID(t *testing.T, repo class.Repository) {
	t.Helper()

	ctx := context.Background()

	initial := class.NewClass("initial", "", nil)
	overwrite, err := class.Unmarshal(initial.ID(), "", "overwrite", time.Now(), time.Now(), nil)
	assert.NoError(t, err)

	assert.NoError(t, repo.AddClass(ctx, initial))
	assert.Error(t, repo.AddClass(ctx, overwrite))

	check, err := repo.GetClass(ctx, initial.ID())
	assert.NoError(t, err)
	assert.Equal(t, initial.Name(), check.Name())
}
