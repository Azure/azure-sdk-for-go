// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	utcTimeJSON = "15:04:05.999999999"
	utcTime     = "15:04:05.999999999"
	timeFormat  = "15:04:05.999999999Z07:00"
)

type PlainTime time.Time

func (t PlainTime) MarshalJSON() ([]byte, error) {
	s, _ := t.MarshalText()
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

func (t PlainTime) MarshalText() ([]byte, error) {
	tt := time.Time(t)
	return []byte(tt.Format(timeFormat)), nil
}

func (t *PlainTime) UnmarshalJSON(data []byte) error {
	layout := utcTimeJSON
	if tzOffsetRegex.Match(data) {
		layout = timeFormat
	}
	return t.Parse(layout, string(data))
}

func (t *PlainTime) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	layout := utcTime
	if tzOffsetRegex.Match(data) {
		layout = timeFormat
	}
	return t.Parse(layout, string(data))
}

func (t *PlainTime) Parse(layout, value string) error {
	p, err := time.Parse(layout, strings.ToUpper(value))
	*t = PlainTime(p)
	return err
}

func (t PlainTime) String() string {
	tt := time.Time(t)
	return tt.Format(timeFormat)
}

func PopulatePlainTime(m map[string]any, k string, t *time.Time) {
	if t == nil {
		return
	} else if azcore.IsNullValue(t) {
		m[k] = nil
		return
	} else if reflect.ValueOf(t).IsNil() {
		return
	}
	m[k] = (*PlainTime)(t)
}

func UnpopulatePlainTime(data json.RawMessage, fn string, t **time.Time) error {
	if data == nil || string(data) == "null" {
		return nil
	}
	var aux PlainTime
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}
	*t = (*time.Time)(&aux)
	return nil
}
