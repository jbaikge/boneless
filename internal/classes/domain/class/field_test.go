package class_test

import (
	"testing"

	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/zeebo/assert"
)

func TestNewField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Title      string
		Valid      bool
		Type       string
		Name       string
		Label      string
		Sort       bool
		Column     int
		Min        string
		Max        string
		Step       string
		Format     string
		Options    string
		ClassId    string
		ClassField string
	}{
		{
			Title: "Valid Arguments",
			Valid: true,
			Type:  "test",
			Name:  "test",
			Label: "Test",
		},
		{
			Title: "Empty Name",
			Type:  "test",
			Label: "Test",
		},
		{
			Title: "Invalid Name",
			Type:  "test",
			Name:  "TEST-NAME",
			Label: "Test",
		},
		{
			Title: "Empty Label",
			Type:  "test",
			Name:  "test",
		},
		{
			Title:   "Class ID with No Field",
			Type:    "test",
			Name:    "test",
			Label:   "Test",
			ClassId: "1234",
		},
	}

	for _, test := range tests {
		data := test
		t.Run(data.Title, func(t *testing.T) {
			t.Parallel()

			_, err := class.NewField(
				data.Type,
				data.Name,
				data.Label,
				data.Sort,
				data.Column,
				data.Min,
				data.Max,
				data.Step,
				data.Format,
				data.Options,
				data.ClassId,
				data.ClassField,
			)
			if data.Valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
