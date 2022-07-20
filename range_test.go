package gocms

import (
	"testing"

	"github.com/zeebo/assert"
)

func TestRange(t *testing.T) {
	unit := "classes"

	r := Range{}
	assert.NoError(t, r.ParseHeader(unit+"=0-9", unit))
}
