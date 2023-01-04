package class_test

import (
	"testing"
	"time"

	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/classes/domain/field"
	"github.com/jbaikge/boneless/internal/common/id"
	"github.com/zeebo/assert"
)

func TestNewClass(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Title    string
		Valid    bool
		Id       string
		ParentId string
		Name     string
		Created  time.Time
		Updated  time.Time
		Fields   []*field.Field
	}{
		{
			Title:    "Valid Arguments",
			Valid:    true,
			Id:       id.New(),
			ParentId: id.New(),
			Name:     "Test",
			Created:  time.Now(),
			Updated:  time.Now(),
			Fields:   make([]*field.Field, 0),
		},
		{
			Title: "Automatic Everything",
			Valid: true,
			Name:  "Test",
		},
		{
			Title: "Empty Name",
		},
		{
			Title: "Invalid Class ID",
			Id:    "1234",
		},
		{
			Title:    "Invalid Parent ID",
			ParentId: "1234",
		},
	}

	for _, test := range tests {
		data := test
		t.Run(data.Title, func(t *testing.T) {
			t.Parallel()

			_, err := class.NewClass(
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
