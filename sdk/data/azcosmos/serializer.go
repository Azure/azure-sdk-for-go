// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"reflect"
)

// Serializer converts application models to and from JSON. It does not encode protocol metadata.
// Implementations must support concurrent calls and remain usable by responses after Client.Close.
type Serializer interface {
	// Marshal returns one UTF-8 JSON value. The returned buffer must not be reused or mutated
	// until the SDK operation returns. The SDK does not infer or modify item fields.
	Marshal(value any) ([]byte, error)

	// Unmarshal decodes a JSON value into a non-nil pointer. data is read-only and borrowed
	// for this call; implementations retaining it must copy it.
	Unmarshal(data []byte, value any) error
}

type jsonSerializer struct{}

func (jsonSerializer) Marshal(value any) ([]byte, error) { return json.Marshal(value) }
func (jsonSerializer) Unmarshal(data []byte, value any) error {
	return json.Unmarshal(data, value)
}

func resolveSerializer(serializer Serializer) Serializer {
	if serializer == nil {
		return jsonSerializer{}
	}
	return serializer
}

func validateSerializer(serializer Serializer) error {
	if serializer == nil {
		return nil
	}
	value := reflect.ValueOf(serializer)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if value.IsNil() {
			return errors.New("azcosmos: Serializer must not contain a nil implementation")
		}
	}
	return nil
}

func marshalItem(serializer Serializer, item any) ([]byte, error) {
	body, err := resolveSerializer(serializer).Marshal(item)
	if err == nil {
		err = validateModelBody(body)
	}
	if err != nil {
		return nil, &Error{Code: CodeSerializationFailed, Message: "encoding item", cause: err}
	}
	return body, nil
}

func unmarshalModel[T any](serializer Serializer, body []byte) (T, error) {
	var model T
	if err := validateModelBody(body); err != nil {
		return model, err
	}
	if err := resolveSerializer(serializer).Unmarshal(body, &model); err != nil {
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
