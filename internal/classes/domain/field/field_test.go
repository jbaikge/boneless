package field_test

import (
	"testing"

	"github.com/jbaikge/boneless/internal/classes/domain/field"
	"github.com/zeebo/assert"
)

func TestNewField(t *testing.T) {
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
		t.Run(test.Title, func(t *testing.T) {
			_, err := field.NewField(
				test.Type,
				test.Name,
				test.Label,
				test.Sort,
				test.Column,
				test.Min,
				test.Max,
				test.Step,
				test.Format,
				test.Options,
				test.ClassId,
				test.ClassField,
			)
			if test.Valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
