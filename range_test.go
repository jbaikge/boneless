package gocms

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

func TestRangeIsZero(t *testing.T) {
	var r Range
	assert.True(t, r.IsZero())

	r.Size = 100
	assert.False(t, r.IsZero())

	r.Size, r.End = 0, 99
	assert.False(t, r.IsZero())

	r.End, r.Start = 0, 10
	assert.False(t, r.IsZero())
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
