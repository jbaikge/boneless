package services

import (
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/zeebo/assert"
)

func TestXidProvider(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	stamp := time.Date(2022, time.August, 9, 12, 0, 0, 0, loc)

	provider := new(XidProvider)

	id := provider.NewWithTime(stamp)

	check, err := xid.FromString(id)
	assert.NoError(t, err)
	assert.True(t, check.Time().Equal(stamp))

	assert.True(t, idProvider.IsValid(id))
	assert.False(t, idProvider.IsValid("0000"))
}
