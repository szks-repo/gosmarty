package object

import (
	"time"
)

type Time struct {
	Value time.Time
}

func NewTime(t time.Time) *Time {
	return &Time{Value: t}
}

func (s *Time) Type() ObjectType {
	return StringType
}

func (s *Time) Inspect() string {
	return s.Value.Format(time.RFC3339)
}
