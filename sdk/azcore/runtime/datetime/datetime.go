// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"encoding"
	"fmt"
	"time"
)

// Constraints is a type constraint that represents the supported time types in the datetime package.
type Constraints interface {
	PlainDate | PlainTime | RFC1123 | RFC3339 | Unix
}

// Parse parses the provided string into one of the supported datetime types.
func Parse[T Constraints](val string) (time.Time, error) {
	var result T

	// Each datetime type exposes UnmarshalText on its pointer receiver.
	um, ok := any(&result).(encoding.TextUnmarshaler)
	if !ok {
		return time.Time{}, fmt.Errorf("type %T does not implement encoding.TextUnmarshaler", result)
	}

	if err := um.UnmarshalText([]byte(val)); err != nil {
		return time.Time{}, err
	}

	return time.Time(result), nil
}
