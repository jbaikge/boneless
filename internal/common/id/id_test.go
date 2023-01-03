package id_test

import (
	"testing"
	"time"

	"github.com/jbaikge/boneless/internal/common/id"
	"github.com/zeebo/assert"
)

func TestID(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	stamp := time.Date(2023, time.January, 3, 12, 0, 0, 0, loc)

	randomId := id.New()
	assert.True(t, id.IsValid(randomId))

	stampedId := id.NewWithTime(stamp)
	assert.True(t, id.IsValid(stampedId))

	assert.False(t, id.IsValid("0000"))
}
