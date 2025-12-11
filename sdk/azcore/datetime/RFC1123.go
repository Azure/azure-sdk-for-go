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
	DateTimeRFC1123JSON = `"` + time.RFC1123 + `"`
)

type DateTimeRFC1123 time.Time

func (t DateTimeRFC1123) MarshalJSON() ([]byte, error) {
	b := []byte(time.Time(t).Format(DateTimeRFC1123JSON))
	return b, nil
}

func (t DateTimeRFC1123) MarshalText() ([]byte, error) {
	b := []byte(time.Time(t).Format(time.RFC1123))
	return b, nil
}

func (t *DateTimeRFC1123) UnmarshalJSON(data []byte) error {
	p, err := time.Parse(DateTimeRFC1123JSON, strings.ToUpper(string(data)))
	*t = DateTimeRFC1123(p)
	return err
}

func (t *DateTimeRFC1123) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	p, err := time.Parse(time.RFC1123, string(data))
	*t = DateTimeRFC1123(p)
	return err
}

func (t DateTimeRFC1123) String() string {
	return time.Time(t).Format(time.RFC1123)
}

func PopulateDateTimeRFC1123(m map[string]any, k string, t *time.Time) {
	if t == nil {
		return
	} else if azcore.IsNullValue(t) {
		m[k] = nil
		return
	} else if reflect.ValueOf(t).IsNil() {
		return
	}
	m[k] = (*DateTimeRFC1123)(t)
}

func UnpopulateDateTimeRFC1123(data json.RawMessage, fn string, t **time.Time) error {
	if data == nil || string(data) == "null" {
		return nil
	}
	var aux DateTimeRFC1123
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}
	*t = (*time.Time)(&aux)
	return nil
}
