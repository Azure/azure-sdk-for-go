//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armwindowsiot

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type DeviceService.
func (d DeviceService) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "etag", d.Etag)
	populate(objectMap, "id", d.ID)
	populate(objectMap, "location", d.Location)
	populate(objectMap, "name", d.Name)
	populate(objectMap, "properties", d.Properties)
	populate(objectMap, "tags", d.Tags)
	populate(objectMap, "type", d.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DeviceService.
func (d *DeviceService) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", d, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "etag":
				err = unpopulate(val, "Etag", &d.Etag)
			delete(rawMsg, key)
		case "id":
				err = unpopulate(val, "ID", &d.ID)
			delete(rawMsg, key)
		case "location":
				err = unpopulate(val, "Location", &d.Location)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &d.Name)
			delete(rawMsg, key)
		case "properties":
				err = unpopulate(val, "Properties", &d.Properties)
			delete(rawMsg, key)
		case "tags":
				err = unpopulate(val, "Tags", &d.Tags)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &d.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", d, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type DeviceServiceCheckNameAvailabilityParameters.
func (d DeviceServiceCheckNameAvailabilityParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "name", d.Name)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DeviceServiceCheckNameAvailabilityParameters.
func (d *DeviceServiceCheckNameAvailabilityParameters) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", d, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "name":
				err = unpopulate(val, "Name", &d.Name)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", d, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type DeviceServiceDescriptionListResult.
func (d DeviceServiceDescriptionListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "nextLink", d.NextLink)
	populate(objectMap, "value", d.Value)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DeviceServiceDescriptionListResult.
func (d *DeviceServiceDescriptionListResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", d, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "nextLink":
				err = unpopulate(val, "NextLink", &d.NextLink)
			delete(rawMsg, key)
		case "value":
				err = unpopulate(val, "Value", &d.Value)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", d, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type DeviceServiceNameAvailabilityInfo.
func (d DeviceServiceNameAvailabilityInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "message", d.Message)
	populate(objectMap, "nameAvailable", d.NameAvailable)
	populate(objectMap, "reason", d.Reason)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DeviceServiceNameAvailabilityInfo.
func (d *DeviceServiceNameAvailabilityInfo) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", d, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "message":
				err = unpopulate(val, "Message", &d.Message)
			delete(rawMsg, key)
		case "nameAvailable":
				err = unpopulate(val, "NameAvailable", &d.NameAvailable)
			delete(rawMsg, key)
		case "reason":
				err = unpopulate(val, "Reason", &d.Reason)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", d, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type DeviceServiceProperties.
func (d DeviceServiceProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "adminDomainName", d.AdminDomainName)
	populate(objectMap, "billingDomainName", d.BillingDomainName)
	populate(objectMap, "notes", d.Notes)
	populate(objectMap, "quantity", d.Quantity)
	populateTimeRFC3339(objectMap, "startDate", d.StartDate)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type DeviceServiceProperties.
func (d *DeviceServiceProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", d, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "adminDomainName":
				err = unpopulate(val, "AdminDomainName", &d.AdminDomainName)
			delete(rawMsg, key)
		case "billingDomainName":
				err = unpopulate(val, "BillingDomainName", &d.BillingDomainName)
			delete(rawMsg, key)
		case "notes":
				err = unpopulate(val, "Notes", &d.Notes)
			delete(rawMsg, key)
		case "quantity":
				err = unpopulate(val, "Quantity", &d.Quantity)
			delete(rawMsg, key)
		case "startDate":
				err = unpopulateTimeRFC3339(val, "StartDate", &d.StartDate)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", d, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ErrorDetails.
func (e ErrorDetails) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "error", e.Error)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorDetails.
func (e *ErrorDetails) UnmarshalJSON(data []byte) error {
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

// MarshalJSON implements the json.Marshaller interface for type ErrorDetailsError.
func (e ErrorDetailsError) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorDetailsError.
func (e *ErrorDetailsError) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
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

// MarshalJSON implements the json.Marshaller interface for type OperationDisplayInfo.
func (o OperationDisplayInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "description", o.Description)
	populate(objectMap, "operation", o.Operation)
	populate(objectMap, "provider", o.Provider)
	populate(objectMap, "resource", o.Resource)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type OperationDisplayInfo.
func (o *OperationDisplayInfo) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", o, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "description":
				err = unpopulate(val, "Description", &o.Description)
			delete(rawMsg, key)
		case "operation":
				err = unpopulate(val, "Operation", &o.Operation)
			delete(rawMsg, key)
		case "provider":
				err = unpopulate(val, "Provider", &o.Provider)
			delete(rawMsg, key)
		case "resource":
				err = unpopulate(val, "Resource", &o.Resource)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", o, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type OperationEntity.
func (o OperationEntity) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "display", o.Display)
	populate(objectMap, "isDataAction", o.IsDataAction)
	populate(objectMap, "name", o.Name)
	populate(objectMap, "origin", o.Origin)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type OperationEntity.
func (o *OperationEntity) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", o, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "display":
				err = unpopulate(val, "Display", &o.Display)
			delete(rawMsg, key)
		case "isDataAction":
				err = unpopulate(val, "IsDataAction", &o.IsDataAction)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &o.Name)
			delete(rawMsg, key)
		case "origin":
				err = unpopulate(val, "Origin", &o.Origin)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", o, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type OperationListResult.
func (o OperationListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "nextLink", o.NextLink)
	populate(objectMap, "value", o.Value)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type OperationListResult.
func (o *OperationListResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", o, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "nextLink":
				err = unpopulate(val, "NextLink", &o.NextLink)
			delete(rawMsg, key)
		case "value":
				err = unpopulate(val, "Value", &o.Value)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", o, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ProxyResource.
func (p ProxyResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", p.ID)
	populate(objectMap, "name", p.Name)
	populate(objectMap, "type", p.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ProxyResource.
func (p *ProxyResource) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", p, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &p.ID)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &p.Name)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &p.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", p, err)
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

// MarshalJSON implements the json.Marshaller interface for type TrackedResource.
func (t TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", t.ID)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "name", t.Name)
	populate(objectMap, "tags", t.Tags)
	populate(objectMap, "type", t.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type TrackedResource.
func (t *TrackedResource) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", t, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &t.ID)
			delete(rawMsg, key)
		case "location":
				err = unpopulate(val, "Location", &t.Location)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &t.Name)
			delete(rawMsg, key)
		case "tags":
				err = unpopulate(val, "Tags", &t.Tags)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &t.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", t, err)
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

