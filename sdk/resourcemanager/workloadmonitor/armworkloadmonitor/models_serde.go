//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armworkloadmonitor

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type ErrorDetails.
func (e ErrorDetails) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "message", e.Message)
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
		case "code":
				err = unpopulate(val, "Code", &e.Code)
			delete(rawMsg, key)
		case "message":
				err = unpopulate(val, "Message", &e.Message)
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

// MarshalJSON implements the json.Marshaller interface for type ErrorResponseError.
func (e ErrorResponseError) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorResponseError.
func (e *ErrorResponseError) UnmarshalJSON(data []byte) error {
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
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type HealthMonitor.
func (h HealthMonitor) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", h.ID)
	populate(objectMap, "name", h.Name)
	populate(objectMap, "properties", h.Properties)
	populate(objectMap, "type", h.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type HealthMonitor.
func (h *HealthMonitor) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", h, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &h.ID)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &h.Name)
			delete(rawMsg, key)
		case "properties":
				err = unpopulate(val, "Properties", &h.Properties)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &h.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", h, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type HealthMonitorList.
func (h HealthMonitorList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "nextLink", h.NextLink)
	populate(objectMap, "value", h.Value)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type HealthMonitorList.
func (h *HealthMonitorList) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", h, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "nextLink":
				err = unpopulate(val, "NextLink", &h.NextLink)
			delete(rawMsg, key)
		case "value":
				err = unpopulate(val, "Value", &h.Value)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", h, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type HealthMonitorProperties.
func (h HealthMonitorProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "currentMonitorState", h.CurrentMonitorState)
	populate(objectMap, "currentStateFirstObservedTimestamp", h.CurrentStateFirstObservedTimestamp)
	populate(objectMap, "evaluationTimestamp", h.EvaluationTimestamp)
	populateAny(objectMap, "evidence", h.Evidence)
	populate(objectMap, "lastReportedTimestamp", h.LastReportedTimestamp)
	populateAny(objectMap, "monitorConfiguration", h.MonitorConfiguration)
	populate(objectMap, "monitorName", h.MonitorName)
	populate(objectMap, "monitorType", h.MonitorType)
	populate(objectMap, "monitoredObject", h.MonitoredObject)
	populate(objectMap, "parentMonitorName", h.ParentMonitorName)
	populate(objectMap, "previousMonitorState", h.PreviousMonitorState)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type HealthMonitorProperties.
func (h *HealthMonitorProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", h, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "currentMonitorState":
				err = unpopulate(val, "CurrentMonitorState", &h.CurrentMonitorState)
			delete(rawMsg, key)
		case "currentStateFirstObservedTimestamp":
				err = unpopulate(val, "CurrentStateFirstObservedTimestamp", &h.CurrentStateFirstObservedTimestamp)
			delete(rawMsg, key)
		case "evaluationTimestamp":
				err = unpopulate(val, "EvaluationTimestamp", &h.EvaluationTimestamp)
			delete(rawMsg, key)
		case "evidence":
				err = unpopulate(val, "Evidence", &h.Evidence)
			delete(rawMsg, key)
		case "lastReportedTimestamp":
				err = unpopulate(val, "LastReportedTimestamp", &h.LastReportedTimestamp)
			delete(rawMsg, key)
		case "monitorConfiguration":
				err = unpopulate(val, "MonitorConfiguration", &h.MonitorConfiguration)
			delete(rawMsg, key)
		case "monitorName":
				err = unpopulate(val, "MonitorName", &h.MonitorName)
			delete(rawMsg, key)
		case "monitorType":
				err = unpopulate(val, "MonitorType", &h.MonitorType)
			delete(rawMsg, key)
		case "monitoredObject":
				err = unpopulate(val, "MonitoredObject", &h.MonitoredObject)
			delete(rawMsg, key)
		case "parentMonitorName":
				err = unpopulate(val, "ParentMonitorName", &h.ParentMonitorName)
			delete(rawMsg, key)
		case "previousMonitorState":
				err = unpopulate(val, "PreviousMonitorState", &h.PreviousMonitorState)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", h, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type HealthMonitorStateChange.
func (h HealthMonitorStateChange) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "id", h.ID)
	populate(objectMap, "name", h.Name)
	populate(objectMap, "properties", h.Properties)
	populate(objectMap, "type", h.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type HealthMonitorStateChange.
func (h *HealthMonitorStateChange) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", h, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "id":
				err = unpopulate(val, "ID", &h.ID)
			delete(rawMsg, key)
		case "name":
				err = unpopulate(val, "Name", &h.Name)
			delete(rawMsg, key)
		case "properties":
				err = unpopulate(val, "Properties", &h.Properties)
			delete(rawMsg, key)
		case "type":
				err = unpopulate(val, "Type", &h.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", h, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type HealthMonitorStateChangeList.
func (h HealthMonitorStateChangeList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "nextLink", h.NextLink)
	populate(objectMap, "value", h.Value)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type HealthMonitorStateChangeList.
func (h *HealthMonitorStateChangeList) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", h, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "nextLink":
				err = unpopulate(val, "NextLink", &h.NextLink)
			delete(rawMsg, key)
		case "value":
				err = unpopulate(val, "Value", &h.Value)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", h, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type HealthMonitorStateChangeProperties.
func (h HealthMonitorStateChangeProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "currentMonitorState", h.CurrentMonitorState)
	populate(objectMap, "currentStateFirstObservedTimestamp", h.CurrentStateFirstObservedTimestamp)
	populate(objectMap, "evaluationTimestamp", h.EvaluationTimestamp)
	populateAny(objectMap, "evidence", h.Evidence)
	populateAny(objectMap, "monitorConfiguration", h.MonitorConfiguration)
	populate(objectMap, "monitorName", h.MonitorName)
	populate(objectMap, "monitorType", h.MonitorType)
	populate(objectMap, "monitoredObject", h.MonitoredObject)
	populate(objectMap, "previousMonitorState", h.PreviousMonitorState)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type HealthMonitorStateChangeProperties.
func (h *HealthMonitorStateChangeProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", h, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "currentMonitorState":
				err = unpopulate(val, "CurrentMonitorState", &h.CurrentMonitorState)
			delete(rawMsg, key)
		case "currentStateFirstObservedTimestamp":
				err = unpopulate(val, "CurrentStateFirstObservedTimestamp", &h.CurrentStateFirstObservedTimestamp)
			delete(rawMsg, key)
		case "evaluationTimestamp":
				err = unpopulate(val, "EvaluationTimestamp", &h.EvaluationTimestamp)
			delete(rawMsg, key)
		case "evidence":
				err = unpopulate(val, "Evidence", &h.Evidence)
			delete(rawMsg, key)
		case "monitorConfiguration":
				err = unpopulate(val, "MonitorConfiguration", &h.MonitorConfiguration)
			delete(rawMsg, key)
		case "monitorName":
				err = unpopulate(val, "MonitorName", &h.MonitorName)
			delete(rawMsg, key)
		case "monitorType":
				err = unpopulate(val, "MonitorType", &h.MonitorType)
			delete(rawMsg, key)
		case "monitoredObject":
				err = unpopulate(val, "MonitoredObject", &h.MonitoredObject)
			delete(rawMsg, key)
		case "previousMonitorState":
				err = unpopulate(val, "PreviousMonitorState", &h.PreviousMonitorState)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", h, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type Operation.
func (o Operation) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "display", o.Display)
	populate(objectMap, "name", o.Name)
	populate(objectMap, "origin", o.Origin)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Operation.
func (o *Operation) UnmarshalJSON(data []byte) error {
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

// MarshalJSON implements the json.Marshaller interface for type OperationDisplay.
func (o OperationDisplay) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "description", o.Description)
	populate(objectMap, "operation", o.Operation)
	populate(objectMap, "provider", o.Provider)
	populate(objectMap, "resource", o.Resource)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type OperationDisplay.
func (o *OperationDisplay) UnmarshalJSON(data []byte) error {
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

// MarshalJSON implements the json.Marshaller interface for type OperationList.
func (o OperationList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "nextLink", o.NextLink)
	populate(objectMap, "value", o.Value)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type OperationList.
func (o *OperationList) UnmarshalJSON(data []byte) error {
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

