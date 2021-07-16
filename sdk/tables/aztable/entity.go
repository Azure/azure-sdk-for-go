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
	fmt.Println(e)
	stringValue := string(b)
	i, err := strconv.ParseInt(stringValue[1:len(stringValue)-1], 10, 64) // Have to peel off the quotations
	*e = EdmInt64(i)
	fmt.Println("e: ", e)
	fmt.Println(e)
	return err
}

type EdmGuid string

func (e EdmGuid) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("\"%v\"", e))
	return buffer.Bytes(), nil
}

type EdmDateTime struct {
	time.Time
}

func (e EdmDateTime) MarshalJSON() ([]byte, error) {
	u := e.UTC()
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
	// fmt.Println("Year: ", year)
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(d[1])
	// fmt.Println("Month: ", month)
	if err != nil {
		return err
	}
	day, err := strconv.Atoi(d[2])
	// fmt.Println("Day: ", day)
	if err != nil {
		return err
	}

	timeVal := splitValue[1]
	// fmt.Println("timeVal: ", timeVal)
	t := strings.Split(timeVal, ":")
	// fmt.Println("t: ", t)
	hours, err := strconv.Atoi(t[0])
	// fmt.Println("Hours: ", hours)
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(t[1])
	// fmt.Println("Minutes: ", minutes)
	if err != nil {
		return err
	}
	s := strings.Split(string(t[2]), ".")
	// fmt.Println("s: ", s)
	seconds, err := strconv.Atoi(s[0])
	// fmt.Println("Seconds: ", seconds)
	if err != nil {
		return err
	}
	// fmt.Println(s[1][0 : len(s[1])-2])
	nano, err := strconv.Atoi(s[1][0 : len(s[1])-2]) // Peel off the 'Z' for UTC
	// fmt.Println("Nano: ", nano)
	if err != nil {
		return err
	}

	*e = EdmDateTime{time.Date(year, time.Month(month), day, hours, minutes, seconds, nano, time.UTC)}
	return nil
}
