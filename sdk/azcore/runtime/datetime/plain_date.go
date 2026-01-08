// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"time"
)

// Constraints is a type constraint that represents the supported time types in the datetime package.
type Constraints interface {
	PlainDate | PlainTime | RFC1123 | RFC3339 | Unix
}

const (
	plainDate     = "2006-01-02"
	plainDateJSON = `"` + plainDate + `"`
)

// PlainDate represents a date value without time information in YYYY-MM-DD format.
// It wraps time.Time and can be marshaled to and unmarshaled from JSON.
type PlainDate time.Time

// MarshalJSON marshals the PlainDate to a JSON byte slice.
func (t PlainDate) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(plainDateJSON)), nil
}

// MarshalText returns a textual representation of PlainDate.
func (t PlainDate) MarshalText() ([]byte, error) {
	tt := time.Time(t)
	return []byte(tt.Format(plainDate)), nil
}

// UnmarshalJSON unmarshals a JSON byte slice into a PlainDate.
func (d *PlainDate) UnmarshalJSON(data []byte) (err error) {
	t, err := time.Parse(plainDateJSON, string(data))
	*d = (PlainDate)(t)
	return err
}

// UnmarshalText decodes the textual representation of PlainDate.
func (t *PlainDate) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	p, err := time.Parse(plainDate, string(data))
	*t = PlainDate(p)
	return err
}

// String returns the string representation of PlainDate.
func (t PlainDate) String() string {
	return time.Time(t).Format(plainDate)
}
