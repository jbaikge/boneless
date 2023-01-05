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

func TestClassParentID(t *testing.T) {
	t.Parallel()

	c := class.NewClass(t.Name(), "", nil)
	assert.Equal(t, "", c.ParentID())

	// Give the class a new Parent ID
	newId := id.New()
	assert.NoError(t, c.UpdateParentID(newId))
	assert.Equal(t, newId, c.ParentID())

	// Reset to an empty ID
	assert.NoError(t, c.UpdateParentID(""))
	assert.Equal(t, "", c.ParentID())

	// Attempt an invalid value
	assert.Error(t, c.UpdateParentID("0000"))
}

func TestClassName(t *testing.T) {
	t.Parallel()

	c := class.NewClass(t.Name(), "", nil)
	assert.Equal(t, t.Name(), c.Name())

	// Give the class a new name
	newName := "updated"
	assert.NoError(t, c.UpdateName(newName))
	assert.Equal(t, newName, c.Name())

	// Attempt to blank the name
	assert.Error(t, c.UpdateName(""))
}

func TestClassCreated(t *testing.T) {
	t.Parallel()

	c := class.NewClass(t.Name(), "", nil)
	created := c.Created()
	assert.False(t, created.IsZero())
}

func TestClassUpdated(t *testing.T) {
	t.Parallel()

	c := class.NewClass(t.Name(), "", nil)
	updated := c.Updated()
	assert.False(t, updated.IsZero())

	// Dramatic pause
	<-time.After(time.Millisecond)

	// Do something that triggers the modified call
	c.UpdateName("New Name")
	assert.True(t, c.Updated().After(updated))
}

func TestClassFields(t *testing.T) {
	t.Parallel()

	fields := make([]*class.Field, 0, 2)
	f1, err := class.NewField("text", "slug", "Slug", false, 0, "", "", "", "", "", "", "")
	assert.NoError(t, err)
	fields = append(fields, f1)

	f2, err := class.NewField("text", "title", "Title", true, 1, "", "", "", "", "", "", "")
	assert.NoError(t, err)
	fields = append(fields, f2)

	c := class.NewClass(t.Name(), "", fields)
	assert.NotNil(t, c.Fields())
	assert.Equal(t, 2, len(c.Fields()))

	c.UpdateFields(fields[0:1])
	assert.NotNil(t, c.Fields())
	assert.Equal(t, 1, len(c.Fields()))

	c.UpdateFields(nil)
	assert.Nil(t, c.Fields())
}
