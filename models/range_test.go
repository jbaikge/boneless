package models

import (
	"testing"

	"github.com/zeebo/assert"
)

func TestContentRangeHeader(t *testing.T) {
	var r Range
	unit := "test"

	assert.Equal(t, unit+" 0-0/0", r.ContentRangeHeader(unit))

	r.Size = 100
	assert.Equal(t, unit+" 0-0/100", r.ContentRangeHeader(unit))

	r.End = 49
	assert.Equal(t, unit+" 0-49/100", r.ContentRangeHeader(unit))

	r.Start = 40
	assert.Equal(t, unit+" 40-49/100", r.ContentRangeHeader(unit))
}

func TestParseHeader(t *testing.T) {
	unit := "test"

	t.Run("Normal", func(t *testing.T) {
		var r Range
		assert.NoError(t, r.ParseHeader(unit+"=0-9", unit))
		assert.Equal(t, 0, r.Start)
		assert.Equal(t, 9, r.End)
	})

	t.Run("InvalidUnit", func(t *testing.T) {
		var r Range
		assert.Error(t, r.ParseHeader("invalid=0-9", unit))
	})

	t.Run("Multiple", func(t *testing.T) {
		var r Range
		assert.Error(t, r.ParseHeader(unit+"=0-9, 10-14", unit))
	})

	t.Run("Negative", func(t *testing.T) {
		var r Range
		assert.Error(t, r.ParseHeader(unit+"=-10", unit))
	})

	t.Run("Malformed", func(t *testing.T) {
		var r Range
		assert.Error(t, r.ParseHeader(unit+"=0~9", unit))
		assert.Error(t, r.ParseHeader(unit+"=9", unit))
	})

	t.Run("MalformedStart", func(t *testing.T) {
		var r Range
		assert.Error(t, r.ParseHeader(unit+"=a-9", unit))
	})

	t.Run("MalformedEnd", func(t *testing.T) {
		var r Range
		assert.Error(t, r.ParseHeader(unit+"=0-b", unit))
	})

	t.Run("EndBeforeStart", func(t *testing.T) {
		var r Range
		assert.Error(t, r.ParseHeader(unit+"=9-0", unit))
	})
}

func TestParseParams(t *testing.T) {
	tests := []struct {
		Name   string
		Params map[string]string
		Start  int
		End    int
		Error  bool
	}{
		{
			Name:   "Empty",
			Params: map[string]string{},
			Start:  0,
			End:    9,
		},
		// React Admin data provider: ra-data-json-server
		{
			Name:   "StartOnly",
			Params: map[string]string{"_start": "10"},
			Error:  true,
		},
		{
			Name:   "EndOnly",
			Params: map[string]string{"_end": "19"},
			Error:  true,
		},
		{
			Name:   "StartAndEnd",
			Params: map[string]string{"_start": "10", "_end": "19"},
			Start:  10,
			End:    19,
		},
		{
			Name:   "StartGreaterEnd",
			Params: map[string]string{"_start": "19", "_end": "10"},
			Error:  true,
		},
		{
			Name:   "NegativeStart",
			Params: map[string]string{"_start": "-5", "_end": "5"},
			Error:  true,
		},
		{
			Name:   "NegativeEnd",
			Params: map[string]string{"_start": "5", "_end": "-5"},
			Error:  true,
		},
		{
			Name:   "InvalidStartInt",
			Params: map[string]string{"_start": "zero", "_end": "4"},
			Error:  true,
		},
		{
			Name:   "InvalidEndInt",
			Params: map[string]string{"_start": "0", "_end": "four"},
			Error:  true,
		},
		// React Admin data provider: ra-data-simple-rest
		{
			Name:   "Range",
			Params: map[string]string{"range": "[5,9]"},
			Start:  5,
			End:    9,
		},
		{
			Name:   "TooManyRangeElements",
			Params: map[string]string{"range": "[5,9,2]"},
			Error:  true,
		},
		{
			Name:   "TooFewRangeElements",
			Params: map[string]string{"range": "[5]"},
			Error:  true,
		},
		{
			Name:   "BadRangeJSON",
			Params: map[string]string{"range": "5,9"},
			Error:  true,
		},
		{
			Name:   "NegativeRangeStart",
			Params: map[string]string{"range": "[-5,9]"},
			Error:  true,
		},
		{
			Name:   "NegativeRangeEnd",
			Params: map[string]string{"range": "[5,-9]"},
			Error:  true,
		},
		// page/per-page
		{
			Name:   "PageOnly",
			Params: map[string]string{"_page": "4"},
			Start:  30,
			End:    39,
		},
		{
			Name:   "PerPageOnly",
			Params: map[string]string{"_per_page": "5"},
			Start:  0,
			End:    4,
		},
		{
			Name:   "PageAndPerPage",
			Params: map[string]string{"_page": "5", "_per_page": "5"},
			Start:  20,
			End:    24,
		},
		{
			Name:   "NegativePage",
			Params: map[string]string{"_page": "-1", "_per_page": "10"},
			Error:  true,
		},
		{
			Name:   "NegativePerPage",
			Params: map[string]string{"_page": "5", "_per_page": "-1"},
			Error:  true,
		},
		{
			Name:   "ZeroPerPage",
			Params: map[string]string{"_page": "5", "_per_page": "0"},
			Error:  true,
		},
		{
			Name:   "InvalidPageInt",
			Params: map[string]string{"_page": "one", "_per_page": "5"},
			Error:  true,
		},
		{
			Name:   "InvalidPerPageInt",
			Params: map[string]string{"_page": "1", "_per_page": "five"},
			Error:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var r Range
			err := r.ParseParams(test.Params)
			if test.Error {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.Start, r.Start)
			assert.Equal(t, test.End, r.End)
		})
	}
}
