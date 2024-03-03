package realtime

import (
	"time"
)

//go:generate go run github.com/golang/mock/mockgen --source=realtime.go --destination=realtime_mock.go --package=realtime

type (
	// Time describe necessary methods for work with time.
	Time interface {
		Now() time.Time
	}
	// RealTime implements Time interface.
	RealTime struct {
		timeNowFunc func() time.Time
	}
)

// NewRealTime create new Time implementation.
func NewRealTime(
	timeNowFunc func() time.Time) Time {
	return &RealTime{
		timeNowFunc: timeNowFunc,
	}
}

// Now returns current time.
func (t *RealTime) Now() time.Time {
	return t.timeNowFunc()
}
