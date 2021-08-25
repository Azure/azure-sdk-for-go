// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
)

var (
	nonePartitionKey      = &partitionKeyInternal{components: nil}
	emptyPartitionKey     = &partitionKeyInternal{components: []interface{}{}}
	undefinedPartitionKey = &partitionKeyInternal{components: []interface{}{partitionKeyUndefinedComponent{}}}
)

type partitionKeyInternal struct {
	components []interface{}
}

func newPartitionKeyInternal(values []interface{}) (*partitionKeyInternal, error) {
	components := make([]interface{}, len(values))
	for i, v := range values {
		var component interface{}
		switch val := v.(type) {
		case nil:
			component = partitionKeyNullComponent{}
		case bool:
			component = partitionKeyBoolComponent{val}
		case string:
			component = partitionKeyStringComponent{val}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			component = partitionKeyNumberComponent{v.(float64)}
		default:
			return nil, fmt.Errorf("PartitionKey can only be a string, bool, or a number: '%T'", v)
		}

		components[i] = component
	}

	return &partitionKeyInternal{components: components}, nil
}

func (p *partitionKeyInternal) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.components)
}

type partitionKeyNumberComponent struct {
	value float64
}

func (p *partitionKeyNumberComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

type partitionKeyBoolComponent struct {
	value bool
}

func (p *partitionKeyBoolComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

type partitionKeyStringComponent struct {
	value string
}

func (p *partitionKeyStringComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

type partitionKeyNullComponent struct{}

func (p *partitionKeyNullComponent) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

type partitionKeyUndefinedComponent struct{}

func (p *partitionKeyUndefinedComponent) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}
