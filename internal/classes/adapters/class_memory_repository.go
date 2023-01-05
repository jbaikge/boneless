package adapters

import (
	"context"
	"fmt"
	"sync"

	"github.com/jbaikge/boneless/internal/classes/domain/class"
)

type ClassMemoryRepository struct {
	classes map[string]class.Class
	lock    *sync.RWMutex
}

func NewClassMemoryRepository() *ClassMemoryRepository {
	return &ClassMemoryRepository{
		classes: make(map[string]class.Class),
		lock:    &sync.RWMutex{},
	}
}

func (r ClassMemoryRepository) AddClass(ctx context.Context, c *class.Class) (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.classes[c.ID()]; ok {
		return fmt.Errorf("class already exists with ID: %s", c.ID())
	}

	r.classes[c.ID()] = *c

	return
}

func (r ClassMemoryRepository) GetClass(ctx context.Context, id string) (c *class.Class, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	value, ok := r.classes[id]
	if !ok {
		return nil, fmt.Errorf("no class exists with ID: %s", id)
	}

	return &value, nil
}

func (r ClassMemoryRepository) UpdateClass(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, c *class.Class) (*class.Class, error),
) (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	existing, ok := r.classes[id]
	if !ok {
		return fmt.Errorf("no class exists wtih ID: %s", id)
	}

	updated, err := updateFn(ctx, &existing)
	if err != nil {
		return
	}

	r.classes[id] = *updated

	return
}
