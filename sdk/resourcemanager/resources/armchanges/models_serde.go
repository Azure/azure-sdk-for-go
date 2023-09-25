//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armchanges

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type ChangeAttributes.
func (c ChangeAttributes) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "changesCount", c.ChangesCount)
	populate(objectMap, "correlationId", c.CorrelationID)
	populate(objectMap, "newResourceSnapshotId", c.NewResourceSnapshotID)
	populate(objectMap, "previousResourceSnapshotId", c.PreviousResourceSnapshotID)
	populate(objectMap, "timestamp", c.Timestamp)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ChangeAttributes.
func (c *ChangeAttributes) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "changesCount":
				err = unpopulate(val, "ChangesCount", &c.ChangesCount)
			delete(rawMsg, key)
		case "correlationId":
				err = unpopulate(val, "CorrelationID", &c.CorrelationID)
			delete(rawMsg, key)
		case "newResourceSnapshotId":
				err = unpopulate(val, "NewResourceSnapshotID", &c.NewResourceSnapshotID)
			delete(rawMsg, key)
		case "previousResourceSnapshotId":
				err = unpopulate(val, "PreviousResourceSnapshotID", &c.PreviousResourceSnapshotID)
			delete(rawMsg, key)
		case "timestamp":
				err = unpopulate(val, "Timestamp", &c.Timestamp)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ChangeBase.
func (c ChangeBase) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "changeCategory", c.ChangeCategory)
	populate(objectMap, "newValue", c.NewValue)
	populate(objectMap, "previousValue", c.PreviousValue)
	populate(objectMap, "propertyChangeType", c.PropertyChangeType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ChangeBase.
func (c *ChangeBase) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "changeCategory":
				err = unpopulate(val, "ChangeCategory", &c.ChangeCategory)
			delete(rawMsg, key)
		case "newValue":
				err = unpopulate(val, "NewValue", &c.NewValue)
			delete(rawMsg, key)
		case "previousValue":
				err = unpopulate(val, "PreviousValue", &c.PreviousValue)
			delete(rawMsg, key)
		case "propertyChangeType":
				err = unpopulate(val, "PropertyChangeType", &c.PropertyChangeType)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ChangeProperties.
func (c ChangeProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "changeAttributes", c.ChangeAttributes)
	populate(objectMap, "changeType", c.ChangeType)
	populate(objectMap, "changes", c.Changes)
	populate(objectMap, "targetResourceId", c.TargetResourceID)
	populate(objectMap, "targetResourceType", c.TargetResourceType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ChangeProperties.
func (c *ChangeProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "changeAttributes":
				err = unpopulate(val, "ChangeAttributes", &c.ChangeAttributes)
			delete(rawMsg, key)
		case "changeType":
				err = unpopulate(val, "ChangeType", &c.ChangeType)
			delete(rawMsg, key)
		case "changes":
				err = unpopulate(val, "Changes", &c.Changes)
			delete(rawMsg, key)
		case "targetResourceId":
				err = unpopulate(val, "TargetResourceID", &c.TargetResourceID)
			delete(rawMsg, key)
		case "targetResourceType":
				err = unpopulate(val, "TargetResourceType", &c.TargetResourceType)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ChangeResourceListResult.
func (c ChangeResourceListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "nextLink", c.NextLink)
	populate(objectMap, "value", c.Value)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ChangeResourceListResult.
func (c *ChangeResourceListResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "nextLink":
				err = unpopulate(val, "NextLink", &c.NextLink)
			delete(rawMsg, key)
		case "value":
				err = unpopulate(val, "Value", &c.Value)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ChangeResourceResult.
func (c ChangeResourceResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", c.ID)
	populate(objectMap, "name", c.Name)
	populate(objectMap, "properties", c.Properties)
	populate(objectMap, "type", c.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ChangeResourceResult.
func (c *ChangeResourceResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &c.ID)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &c.Name)
			delete(rawMsg, key)
		case "properties":
				err = unpopulate(val, "Properties", &c.Properties)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &c.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ErrorAdditionalInfo.
func (e ErrorAdditionalInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populateAny(objectMap, "info", e.Info)
	populate(objectMap, "type", e.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorAdditionalInfo.
func (e *ErrorAdditionalInfo) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "info":
				err = unpopulate(val, "Info", &e.Info)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &e.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ErrorDetail.
func (e ErrorDetail) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorDetail.
func (e *ErrorDetail) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "additionalInfo":
				err = unpopulate(val, "AdditionalInfo", &e.AdditionalInfo)
			delete(rawMsg, key)
		case "code":
				err = unpopulate(val, "Code", &e.Code)
			delete(rawMsg, key)
		case "details":
				err = unpopulate(val, "Details", &e.Details)
			delete(rawMsg, key)
		case "message":
				err = unpopulate(val, "Message", &e.Message)
			delete(rawMsg, key)
		case "target":
				err = unpopulate(val, "Target", &e.Target)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ErrorResponse.
func (e ErrorResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "error", e.Error)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorResponse.
func (e *ErrorResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "error":
				err = unpopulate(val, "Error", &e.Error)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", r.ID)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "type", r.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Resource.
func (r *Resource) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &r.ID)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &r.Name)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &r.Type)
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

func populateAny(m map[string]any, k string, v any) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else {
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

