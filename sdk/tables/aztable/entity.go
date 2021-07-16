// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type EdmInt64 int64

func (e EdmInt64) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("\"%d\"", e))
	return buffer.Bytes(), nil
}

func (e *EdmInt64) UnmarshalJSON(b []byte) error {
	stringValue := string(b)
	i, err := strconv.ParseInt(stringValue[1:len(stringValue)-1], 10, 64) // Have to peel off the quotations
	*e = EdmInt64(i)
	return err
}

type EdmGuid string

func (e EdmGuid) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("\"%v\"", e))
	return buffer.Bytes(), nil
}

type EdmDateTime time.Time

func (e EdmDateTime) MarshalJSON() ([]byte, error) {
	u := time.Time(e).UTC()
	formatted := fmt.Sprintf("\"%d-%d-%dT%d:%d:%d.%dZ\"", u.Year(), u.Month(), u.Day(), u.Hour(), u.Minute(), u.Second(), u.Nanosecond())
	buffer := bytes.NewBufferString(formatted)
	return buffer.Bytes(), nil
}

func (e *EdmDateTime) UnmarshalJSON(b []byte) error {
	stringValue := string(b)
	splitValue := strings.Split(stringValue, "T")
	date := splitValue[0]
	d := strings.Split(date, "-")
	year, err := strconv.Atoi(d[0][1:len(d[0])])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(d[1])
	if err != nil {
		return err
	}
	day, err := strconv.Atoi(d[2])
	if err != nil {
		return err
	}

	timeVal := splitValue[1]
	t := strings.Split(timeVal, ":")
	hours, err := strconv.Atoi(t[0])
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(t[1])
	if err != nil {
		return err
	}
	s := strings.Split(string(t[2]), ".")
	seconds, err := strconv.Atoi(s[0])
	if err != nil {
		return err
	}
	nano, err := strconv.Atoi(s[1][0 : len(s[1])-2]) // Peel off the 'Z' for UTC
	if err != nil {
		return err
	}

	*e = EdmDateTime(time.Date(year, time.Month(month), day, hours, minutes, seconds, nano, time.UTC))
	return nil
}
