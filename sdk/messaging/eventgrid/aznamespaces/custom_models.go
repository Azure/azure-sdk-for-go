//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces

import (
	"encoding/json"
	"fmt"
)

// NOTE: most of this code doesn't need to be here but there's a bug
// that's trimming it out: https://github.com/Azure/typespec-azure/issues/1412

// Error implements the error interface for type Error.
// Note that the message contents are not contractual and can change over time.
func (e *Error) Error() string {
	if e.Message == nil {
		return ""
	}

	return *e.Message
}

// Error - The error object.
type Error struct {
	// REQUIRED; One of a server-defined set of error codes.
	Code *string

	// REQUIRED; A human-readable representation of the error.
	Message *string
}

// MarshalJSON implements the json.Marshaller interface for type Error.
func (e Error) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "code", e.Code)

	populate(objectMap, "message", e.Message)

	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Error.
func (e *Error) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err) //nolint:errorlint
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "code":
			err = unpopulate(val, "Code", &e.Code)
			delete(rawMsg, key)
		case "message":
			err = unpopulate(val, "Message", &e.Message)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err) //nolint:errorlint
		}
	}
	return nil
}
