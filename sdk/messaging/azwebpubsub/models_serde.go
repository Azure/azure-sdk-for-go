//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azwebpubsub

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// MarshalJSON implements the json.Marshaller interface for type AddToGroupsRequest.
func (a AddToGroupsRequest) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "filter", a.Filter)
	populate(objectMap, "groups", a.Groups)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type AddToGroupsRequest.
func (a *AddToGroupsRequest) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "filter":
			err = unpopulate(val, "Filter", &a.Filter)
			delete(rawMsg, key)
		case "groups":
			err = unpopulate(val, "Groups", &a.Groups)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type TokenResponse.
func (c TokenResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "token", c.Token)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type TokenResponse.
func (c *TokenResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "token":
			err = unpopulate(val, "Token", &c.Token)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type RemoveFromGroupsRequest.
func (r RemoveFromGroupsRequest) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "filter", r.Filter)
	populate(objectMap, "groups", r.Groups)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type RemoveFromGroupsRequest.
func (r *RemoveFromGroupsRequest) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "filter":
			err = unpopulate(val, "Filter", &r.Filter)
			delete(rawMsg, key)
		case "groups":
			err = unpopulate(val, "Groups", &r.Groups)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", r, err)
		}
	}
	return nil
}

func populate(m map[string]any, k string, v any) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, fn string, v any) error {
	if data == nil {
		return nil
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}
	return nil
}
