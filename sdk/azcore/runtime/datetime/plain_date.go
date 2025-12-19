// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"fmt"
	"time"
)

// Time is a type constraint that represents the supported time types in the datetime package
type Time interface {
	PlainDate | PlainTime | RFC1123 | RFC3339 | Unix
}

const (
	fullDateJSON = `"2006-01-02"`
	jsonFormat   = `"%04d-%02d-%02d"`
)

// PlainDate represents a date value without time information in YYYY-MM-DD format.
// It wraps time.Time and can be marshaled to and unmarshaled from JSON.
type PlainDate time.Time

// MarshalJSON marshals the PlainDate to a JSON byte slice.
func (t PlainDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(jsonFormat, time.Time(t).Year(), time.Time(t).Month(), time.Time(t).Day())), nil
}

// UnmarshalJSON unmarshals a JSON byte slice into a PlainDate.
func (d *PlainDate) UnmarshalJSON(data []byte) (err error) {
	t, err := time.Parse(fullDateJSON, string(data))
	*d = (PlainDate)(t)
	return err
}
