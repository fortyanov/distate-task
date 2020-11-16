package utils

import (
	"strings"
	"time"
)

type TimeISO time.Time

func (mt *TimeISO) UnmarshalJSON(b []byte) (err error) {
	var t time.Time
	
	s := strings.Trim(string(b), "\"")
	
	if s == "" {
		*mt = TimeISO(t)
		return
	}

	if t, err = time.Parse(time.RFC3339, s); err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z07:00", s)
	}
	if err == nil {
		*mt = TimeISO(t)
	}
	return
}