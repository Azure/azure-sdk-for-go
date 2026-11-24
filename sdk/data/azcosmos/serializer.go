// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"errors"
)

func marshalItem(item any) ([]byte, error) {
	body, err := json.Marshal(item)
	if err == nil {
		err = validateModelBody(body)
	}
	if err != nil {
		return nil, &Error{Code: CodeSerializationFailed, Message: "encoding item", cause: err}
	}
	if !json.Valid(body) || bytes.TrimSpace(body)[0] != '{' {
		return nil, &Error{
			Code:    CodeSerializationFailed,
			Message: "encoding item: serializer must produce a valid JSON object",
		}
	}
	return body, nil
}

func unmarshalModel[T any](body []byte) (T, error) {
	var model T
	if err := validateModelBody(body); err != nil {
		return model, err
	}
	if err := json.Unmarshal(body, &model); err != nil {
		var zero T
		return zero, err
	}
	return model, nil
}

func validateModelBody(body []byte) error {
	if len(body) == 0 {
		return errors.New("azcosmos: item body is empty")
	}
	if body[0] == 0x80 {
		return errors.New("azcosmos: Cosmos binary JSON is not supported by the model serializer")
	}
	return nil
}
