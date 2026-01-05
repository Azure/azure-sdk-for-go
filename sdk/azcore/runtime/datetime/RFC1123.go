// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"strings"
	"time"
)

const (
	rfc1123JSON = `"` + time.RFC1123 + `"`
)

// RFC1123 represents a date and time value in RFC 1123 format as defined in
// https://tools.ietf.org/html/rfc1123.
type RFC1123 time.Time

// MarshalJSON marshals the RFC1123 timestamp to a JSON byte slice.
func (t RFC1123) MarshalJSON() ([]byte, error) {
	b := []byte(time.Time(t).Format(rfc1123JSON))
	return b, nil
}

// MarshalText returns a textual representation of RFC1123.
func (t RFC1123) MarshalText() ([]byte, error) {
	b := []byte(time.Time(t).Format(time.RFC1123))
	return b, nil
}

// UnmarshalJSON unmarshals a JSON byte slice into an RFC1123 timestamp.
func (t *RFC1123) UnmarshalJSON(data []byte) error {
	if string(data) == jsonNull {
		return nil
	}
	p, err := time.Parse(rfc1123JSON, strings.ToUpper(string(data)))
	*t = RFC1123(p)
	return err
}

// UnmarshalText decodes the textual representation of RFC1123.
func (t *RFC1123) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	p, err := time.Parse(time.RFC1123, string(data))
	*t = RFC1123(p)
	return err
}

// String returns the string of RFC1123.
func (t RFC1123) String() string {
	return time.Time(t).Format(time.RFC1123)
}
