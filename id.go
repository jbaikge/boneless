package gocms

import (
	"time"

	"github.com/rs/xid"
)

var idProvider IdProvider = new(XidProvider)

func SetIdProvider(provider IdProvider) {
	idProvider = provider
}

type IdProvider interface {
	New() string
	NewWithTime(time.Time) string
	IsValid(string) bool
}

type XidProvider struct {
}

func (p XidProvider) New() string {
	return xid.New().String()
}

func (p XidProvider) NewWithTime(t time.Time) string {
	return xid.NewWithTime(t).String()
}

func (p XidProvider) IsValid(id string) bool {
	_, err := xid.FromString(id)
	return err == nil
}
