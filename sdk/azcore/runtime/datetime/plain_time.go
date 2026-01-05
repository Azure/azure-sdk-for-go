// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"strings"
	"time"
)

const (
	timeOnlyJSON = `"` + time.TimeOnly + `"`
)

// PlainTime represents a time value without date information. It supports HH:MM:SS format
// with optional nanosecond precision and timezone information.
type PlainTime time.Time

// MarshalJSON marshals the PlainTime to a JSON byte slice.
func (t PlainTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(timeOnlyJSON)), nil
}

// MarshalText returns a textual representation of PlainTime.
func (t PlainTime) MarshalText() ([]byte, error) {
	tt := time.Time(t)
	return []byte(tt.Format(time.TimeOnly)), nil
}

// UnmarshalJSON unmarshals a JSON byte slice into PlainTime.
func (t *PlainTime) UnmarshalJSON(data []byte) error {
	if string(data) == jsonNull {
		return nil
	}
	return t.parse(timeOnlyJSON, string(data))
}

// UnmarshalText decodes the textual representation of PlainTime.
func (t *PlainTime) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return t.parse(time.TimeOnly, string(data))
}

// parse parses a time string using the specified layout
func (t *PlainTime) parse(layout, value string) error {
	p, err := time.Parse(layout, strings.ToUpper(value))
	*t = PlainTime(p)
	return err
}

// String returns the string of PlainTime.
func (t PlainTime) String() string {
	tt := time.Time(t)
	return tt.Format(time.TimeOnly)
}
