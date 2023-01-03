package id

import (
	"time"

	"github.com/rs/xid"
)

func New() string {
	return xid.New().String()
}

func NewWithTime(t time.Time) string {
	return xid.NewWithTime(t).String()
}

func IsValid(id string) bool {
	_, err := xid.FromString(id)
	return err == nil
}
