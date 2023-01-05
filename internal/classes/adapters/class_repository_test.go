package adapters_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jbaikge/boneless/internal/classes/adapters"
	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/id"
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

			t.Run("Unknown ID", func(t *testing.T) {
				t.Parallel()
				testRetrieveUnknownID(t, repo.Repository)
			})

			t.Run("Update", func(t *testing.T) {
				t.Parallel()
				testUpdate(t, repo.Repository)
			})

			t.Run("Update Function Fail", func(t *testing.T) {
				t.Parallel()
				testUpdateFnFail(t, repo.Repository)
			})

			t.Run("Update Non-Existent", func(t *testing.T) {
				t.Parallel()
				testBadUpdate(t, repo.Repository)
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

func testRetrieveUnknownID(t *testing.T, repo class.Repository) {
	t.Helper()

	// Nothing exists with this ID
	_, err := repo.GetClass(context.Background(), id.New())
	assert.Error(t, err)
}

func testUpdate(t *testing.T, repo class.Repository) {
	t.Helper()

	ctx := context.Background()

	initial := class.NewClass("initial", "", nil)
	assert.NoError(t, repo.AddClass(ctx, initial))

	newName := "update"
	assert.NoError(t, repo.UpdateClass(
		ctx,
		initial.ID(),
		func(ctx context.Context, c *class.Class) (*class.Class, error) {
			c.UpdateName(newName)
			return c, nil
		},
	))

	check, err := repo.GetClass(ctx, initial.ID())
	assert.NoError(t, err)
	assert.Equal(t, newName, check.Name())
}

func testUpdateFnFail(t *testing.T, repo class.Repository) {
	t.Helper()

	ctx := context.Background()

	initial := class.NewClass("initial", "", nil)
	assert.NoError(t, repo.AddClass(ctx, initial))

	assert.Error(t, repo.UpdateClass(
		ctx,
		initial.ID(),
		func(ctx context.Context, c *class.Class) (*class.Class, error) {
			return nil, errors.New("bad update")
		},
	))
}

// This should silently do nothing, though maybe it should return an error if
// the record is not found?
func testBadUpdate(t *testing.T, repo class.Repository) {
	t.Helper()

	ctx := context.Background()
	update := class.NewClass(t.Name(), "", nil)
	assert.Error(t, repo.UpdateClass(
		ctx,
		update.ID(),
		func(ctx context.Context, c *class.Class) (*class.Class, error) {
			t.Fatal("should never reach the update function")
			return nil, nil
		},
	))
}
