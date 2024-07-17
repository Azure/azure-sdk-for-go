//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package datetime contains helpers for different date/time formats.
// The content is intended for SDK authors.
package datetime

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Format defines the supported date/time formats.
type Format int

const (
	FormatRFC3339  Format = 0
	FormatRFC1123  Format = 1
	FormatRFC7231  Format = 2
	FormatDateOnly Format = 3
	FormatTimeOnly Format = 4
	FormatUnix     Format = 5
)

// Parse parses a formatted date/time string and returns a [time.Time] or an error.
//   - format is the expected format of the string value
//   - value is the string value to parse
func Parse(format Format, value string) (time.Time, error) {
	switch format {
	case FormatDateOnly:
		return parse(time.DateOnly, value)
	case FormatRFC1123:
		return parse(time.RFC1123, value)
	case FormatRFC3339:
		return parse(rfc3339Format([]byte(value)), value)
	case FormatRFC7231:
		return parse(http.TimeFormat, value)
	case FormatTimeOnly:
		return parse(timeOnly, value)
	case FormatUnix:
		return parseUnixTime(value)
	default:
		return time.Time{}, fmt.Errorf("invalid format %d", format)
	}
}

// UTC returns the non-nil t with the location set to UTC.
func UTC(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	tt := t.UTC()
	return &tt
}

// DateOnly is a [time.Time] where only the date components have a value.
type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return asJSON(d.String()), nil
}

func (d DateOnly) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(FormatDateOnly, data, d)
}

func (d *DateOnly) UnmarshalText(data []byte) error {
	return unmarshalText(FormatDateOnly, data, d)
}

func (d DateOnly) String() string {
	return time.Time(d).Format(dateOnly)
}

// RFC1123 is a [time.Time] in RFC1123 format.
type RFC1123 time.Time

func (r RFC1123) MarshalJSON() ([]byte, error) {
	return asJSON(r.String()), nil
}

func (r RFC1123) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *RFC1123) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(FormatRFC1123, data, r)
}

func (r *RFC1123) UnmarshalText(data []byte) error {
	return unmarshalText(FormatRFC1123, data, r)
}

func (r RFC1123) String() string {
	return time.Time(r).Format(time.RFC1123)
}

// RFC3339 is a [time.Time] in RFC3339 format.
type RFC3339 time.Time

func (r RFC3339) MarshalJSON() ([]byte, error) {
	return asJSON(r.String()), nil
}

func (r RFC3339) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *RFC3339) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(FormatRFC3339, data, r)
}

func (r *RFC3339) UnmarshalText(data []byte) error {
	return unmarshalText(FormatRFC3339, data, r)
}

func (r RFC3339) String() string {
	return time.Time(r).Format(time.RFC3339Nano)
}

// RFC7231 is a [time.Time] in RFC7231 format.
type RFC7231 time.Time

func (r RFC7231) MarshalJSON() ([]byte, error) {
	return asJSON(r.String()), nil
}

func (r RFC7231) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *RFC7231) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(FormatRFC7231, data, r)
}

func (r *RFC7231) UnmarshalText(data []byte) error {
	return unmarshalText(FormatRFC7231, data, r)
}

func (r RFC7231) String() string {
	return time.Time(r).Format(http.TimeFormat)
}

// TimeOnly is a [time.Time] where only the time components have a value.
type TimeOnly time.Time

func (t TimeOnly) MarshalJSON() ([]byte, error) {
	return asJSON(t.String()), nil
}

func (t TimeOnly) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(FormatTimeOnly, data, t)
}

func (t *TimeOnly) UnmarshalText(data []byte) error {
	return unmarshalText(FormatTimeOnly, data, t)
}

func (t TimeOnly) String() string {
	return time.Time(t).Format(timeOnly)
}

// Unix is a [time.Time] represented as a Unix time stamp.
type Unix time.Time

func (u Unix) MarshalJSON() ([]byte, error) {
	// unix time stamps are sent as a JSON number, so don't quote them
	return []byte(u.String()), nil
}

func (u Unix) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u *Unix) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == "null" {
		return nil
	}
	// Unix time is a JSON number, no quotes to strip off
	return u.UnmarshalText(data)
}

func (u *Unix) UnmarshalText(data []byte) error {
	// this is to handle XML with an empty value, e.g. <SomeTime />
	if len(data) == 0 {
		return nil
	}

	tt, err := parseUnixTime(string(data))
	if err != nil {
		return err
	}
	*u = Unix(tt)
	return nil
}

func (u Unix) String() string {
	return fmt.Sprintf("%d", time.Time(u).Unix())
}

// TypeConstraint is a generic type constraint for the supported date/time formats.
type TypeConstraint interface {
	DateOnly | RFC1123 | RFC3339 | RFC7231 | TimeOnly | Unix
}

func asJSON(v string) []byte {
	return []byte(`"` + v + `"`)
}

func parse(format, value string) (time.Time, error) {
	return time.Parse(format, strings.ToUpper(value))
}

func unmarshalJSON[T TypeConstraint](format Format, data []byte, dest *T) error {
	if data == nil || string(data) == "null" {
		return nil
	}
	// strip off quotes
	data = data[1 : len(data)-1]
	return unmarshal(format, data, dest)
}

func unmarshalText[T TypeConstraint](format Format, data []byte, dest *T) error {
	// this is to handle XML with an empty value, e.g. <SomeTime />
	if len(data) == 0 {
		return nil
	}
	return unmarshal(format, data, dest)
}

func unmarshal[T TypeConstraint](format Format, data []byte, dest *T) error {
	p, err := Parse(format, string(data))
	if err != nil {
		return err
	}
	*dest = T(p)
	return nil
}

const dateOnly = "2006-01-02"
const timeOnly = "15:04:05.999999999"

// Azure reports time in UTC but it doesn't include the 'Z' time zone suffix in some cases.
var tzOffsetRegex = regexp.MustCompile(`(?:Z|z|\+|-)(?:\d+:\d+)*"*$`)

const (
	utcDateTime    = "2006-01-02T15:04:05.999999999"
	utcDateTimeNoT = "2006-01-02 15:04:05.999999999"
	dateTimeNoT    = `2006-01-02 15:04:05.999999999Z07:00`
)

func rfc3339Format(data []byte) string {
	// for RFC3339 there are several corner-cases we need to handle
	tzOffset := tzOffsetRegex.Match(data)
	hasT := strings.Contains(string(data), "T") || strings.Contains(string(data), "t")
	var layout string
	if tzOffset && hasT {
		layout = time.RFC3339Nano
	} else if tzOffset {
		layout = dateTimeNoT
	} else if hasT {
		layout = utcDateTime
	} else {
		layout = utcDateTimeNoT
	}
	return layout
}

func parseUnixTime(value string) (time.Time, error) {
	sec, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}
