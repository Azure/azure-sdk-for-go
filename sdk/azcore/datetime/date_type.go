// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	fullDateJSON = `2006-01-02`
	jsonFormat   = `%04d-%02d-%02d`
)

type DateType time.Time

func (t DateType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(jsonFormat, time.Time(t).Year(), time.Time(t).Month(), time.Time(t).Day())), nil
}

func (d *DateType) UnmarshalJSON(data []byte) (err error) {
	t, err := time.Parse(fullDateJSON, string(data))
	*d = (DateType)(t)
	return err
}

func PopulateDateType(m map[string]any, k string, t *time.Time) {
	if t == nil {
		return
	} else if azcore.IsNullValue(t) {
		m[k] = nil
		return
	} else if reflect.ValueOf(t).IsNil() {
		return
	}
	m[k] = (*DateType)(t)
}

func UnpopulateDateType(data json.RawMessage, fn string, t **time.Time) error {
	if data == nil || string(data) == "null" {
		return nil
	}
	var aux DateType
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}
	*t = (*time.Time)(&aux)
	return nil
}
