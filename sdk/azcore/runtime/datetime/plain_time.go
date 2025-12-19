// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"fmt"
	"strings"
	"time"
)

const (
	plainTimeUTC     = "15:04:05.999999999"
	plainTimeUTCJSON = `"` + plainTimeUTC + `"`
	plainTime        = "15:04:05.999999999Z07:00"
	plainTimeJSON    = `"` + plainTime + `"`
)

// PlainTime represents a time value without date information. It supports HH:MM:SS format
// with optional nanosecond precision and timezone information.
type PlainTime time.Time

// MarshalJSON marshals the PlainTime to a JSON byte slice.
func (t PlainTime) MarshalJSON() ([]byte, error) {
	s, _ := t.MarshalText()
	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}

// MarshalText returns a textual representation of PlainTime
func (t PlainTime) MarshalText() ([]byte, error) {
	tt := time.Time(t)
	return []byte(tt.Format(plainTime)), nil
}

// UnmarshalJSON unmarshals a JSON byte slice into PlainTime.
func (t *PlainTime) UnmarshalJSON(data []byte) error {
	layout := plainTimeUTCJSON
	if tzOffsetRegex.Match(data) {
		layout = plainTimeJSON
	}
	return t.Parse(layout, string(data))
}

// UnmarshalText decodes the textual representation of PlainTime
func (t *PlainTime) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	layout := plainTimeUTC
	if tzOffsetRegex.Match(data) {
		layout = plainTime
	}
	return t.Parse(layout, string(data))
}

// Parse parses a time string using the specified layout
func (t *PlainTime) Parse(layout, value string) error {
	p, err := time.Parse(layout, strings.ToUpper(value))
	*t = PlainTime(p)
	return err
}

// String returns the string of PlainTime
func (t PlainTime) String() string {
	tt := time.Time(t)
	return tt.Format(plainTime)
}
