// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package datetime

import (
	"encoding/json"
	"fmt"
	"time"
)

// Unix represents a Unix timestamp (seconds since January 1, 1970 UTC).
type Unix time.Time

// MarshalJSON marshals the Unix timestamp to a JSON byte slice.
func (t Unix) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix())
}

// MarshalText returns a textual representation of Unix.
func (t Unix) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalJSON unmarshals a JSON byte slice into a Unix timestamp.
func (t *Unix) UnmarshalJSON(data []byte) error {
	return t.parse(data)
}

// UnmarshalText decodes the textual representation of Unix.
func (t *Unix) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return t.parse(data)
}

// parses a Unix timestamp from a byte slice.
func (t *Unix) parse(data []byte) error {
	var seconds int64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return err
	}
	*t = Unix(time.Unix(seconds, 0))
	return nil
}

// String returns the string of Unix.
func (t Unix) String() string {
	return fmt.Sprintf("%d", time.Time(t).Unix())
}
