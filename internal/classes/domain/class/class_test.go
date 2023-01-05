package class_test

import (
	"testing"
	"time"

	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/id"
	"github.com/zeebo/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Title    string
		Valid    bool
		Id       string
		ParentId string
		Name     string
		Created  time.Time
		Updated  time.Time
		Fields   []*class.Field
	}{
		{
			Title:    "Valid Arguments",
			Valid:    true,
			Id:       id.New(),
			ParentId: id.New(),
			Name:     "Test",
			Created:  time.Now(),
			Updated:  time.Now(),
			Fields:   make([]*class.Field, 0),
		},
		{
			Title: "Empty Everything",
		},
		{
			Title: "Invalid Class ID",
			Id:    "1234",
		},
		{
			Title:    "Invalid Parent ID",
			Id:       id.New(),
			ParentId: "1234",
		},
		{
			Title:    "Empty Name",
			Id:       id.New(),
			ParentId: "",
		},
		{
			Title:    "Empty Created Time",
			Id:       id.New(),
			ParentId: "",
			Name:     "Test",
		},
		{
			Title:    "Empty Created Time",
			Id:       id.New(),
			ParentId: "",
			Name:     "Test",
			Created:  time.Now(),
		},
	}

	for _, test := range tests {
		data := test
		t.Run(data.Title, func(t *testing.T) {
			t.Parallel()

			_, err := class.Unmarshal(
				data.Id,
				data.ParentId,
				data.Name,
				data.Created,
				data.Updated,
				data.Fields,
			)
			if data.Valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestClassID(t *testing.T) {
	t.Parallel()

	t.Run("Automatic ID", func(t *testing.T) {
		c := class.NewClass(t.Name(), "", nil)
		assert.True(t, id.IsValid(c.ID()))
	})

	t.Run("Explicit ID", func(t *testing.T) {
		classId := id.New()
		c, err := class.Unmarshal(classId, "", t.Name(), time.Now(), time.Now(), nil)
		assert.NoError(t, err)
		assert.Equal(t, classId, c.ID())
	})
}

func TestClassName(t *testing.T) {
	t.Parallel()

	c := class.NewClass(t.Name(), "", nil)
	assert.Equal(t, t.Name(), c.Name())
}
