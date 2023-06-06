// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package azeventgrid

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type AcknowledgeOptions.
func (a AcknowledgeOptions) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "lockTokens", a.LockTokens)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type AcknowledgeOptions.
func (a *AcknowledgeOptions) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "lockTokens":
			err = unpopulate(val, "LockTokens", &a.LockTokens)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type AcknowledgeResult.
func (a AcknowledgeResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "failedLockTokens", a.FailedLockTokens)
	populate(objectMap, "succeededLockTokens", a.SucceededLockTokens)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type AcknowledgeResult.
func (a *AcknowledgeResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "failedLockTokens":
			err = unpopulate(val, "FailedLockTokens", &a.FailedLockTokens)
			delete(rawMsg, key)
		case "succeededLockTokens":
			err = unpopulate(val, "SucceededLockTokens", &a.SucceededLockTokens)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type Error.
func (a Error) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "code", a.Code)
	populate(objectMap, "details", a.Details)
	populate(objectMap, "innererror", a.Innererror)
	populate(objectMap, "message", a.Message)
	populate(objectMap, "target", a.Target)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Error.
func (a *Error) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "code":
			err = unpopulate(val, "Code", &a.Code)
			delete(rawMsg, key)
		case "details":
			err = unpopulate(val, "Details", &a.Details)
			delete(rawMsg, key)
		case "innererror":
			err = unpopulate(val, "Innererror", &a.Innererror)
			delete(rawMsg, key)
		case "message":
			err = unpopulate(val, "Message", &a.Message)
			delete(rawMsg, key)
		case "target":
			err = unpopulate(val, "Target", &a.Target)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ErrorResponse.
func (a ErrorResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "error", a.Error)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ErrorResponse.
func (a *ErrorResponse) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "error":
			err = unpopulate(val, "Error", &a.Error)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type InnerError.
func (a InnerError) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "code", a.Code)
	populate(objectMap, "innererror", a.Innererror)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type InnerError.
func (a *InnerError) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "code":
			err = unpopulate(val, "Code", &a.Code)
			delete(rawMsg, key)
		case "innererror":
			err = unpopulate(val, "Innererror", &a.Innererror)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type BrokerProperties.
func (b BrokerProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "deliveryCount", b.DeliveryCount)
	populate(objectMap, "lockToken", b.LockToken)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type BrokerProperties.
func (b *BrokerProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", b, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "deliveryCount":
			err = unpopulate(val, "DeliveryCount", &b.DeliveryCount)
			delete(rawMsg, key)
		case "lockToken":
			err = unpopulate(val, "LockToken", &b.LockToken)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", b, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type CloudEvent.
func (c CloudEvent) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populateAny(objectMap, "data", c.Data)
	populateByteArray(objectMap, "data_base64", c.DataBase64, runtime.Base64StdFormat)
	populate(objectMap, "datacontenttype", c.DataContentType)
	populate(objectMap, "dataschema", c.DataSchema)
	populate(objectMap, "id", c.ID)
	populate(objectMap, "source", c.Source)
	populate(objectMap, "specversion", c.SpecVersion)
	populate(objectMap, "subject", c.Subject)
	populateTimeRFC3339(objectMap, "time", c.Time)
	populate(objectMap, "type", c.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type CloudEvent.
func (c *CloudEvent) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "data":
			err = unpopulate(val, "Data", &c.Data)
			delete(rawMsg, key)
		case "data_base64":
			err = runtime.DecodeByteArray(string(val), &c.DataBase64, runtime.Base64StdFormat)
			delete(rawMsg, key)
		case "datacontenttype":
			err = unpopulate(val, "DataContentType", &c.DataContentType)
			delete(rawMsg, key)
		case "dataschema":
			err = unpopulate(val, "DataSchema", &c.DataSchema)
			delete(rawMsg, key)
		case "id":
			err = unpopulate(val, "ID", &c.ID)
			delete(rawMsg, key)
		case "source":
			err = unpopulate(val, "Source", &c.Source)
			delete(rawMsg, key)
		case "specversion":
			err = unpopulate(val, "SpecVersion", &c.SpecVersion)
			delete(rawMsg, key)
		case "subject":
			err = unpopulate(val, "Subject", &c.Subject)
			delete(rawMsg, key)
		case "time":
			err = unpopulateTimeRFC3339(val, "Time", &c.Time)
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

// MarshalJSON implements the json.Marshaller interface for type FailedLockToken.
func (f FailedLockToken) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "errorCode", f.ErrorCode)
	populate(objectMap, "errorDescription", f.ErrorDescription)
	populate(objectMap, "lockToken", f.LockToken)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type FailedLockToken.
func (f *FailedLockToken) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", f, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "errorCode":
			err = unpopulate(val, "ErrorCode", &f.ErrorCode)
			delete(rawMsg, key)
		case "errorDescription":
			err = unpopulate(val, "ErrorDescription", &f.ErrorDescription)
			delete(rawMsg, key)
		case "lockToken":
			err = unpopulate(val, "LockToken", &f.LockToken)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", f, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ReceiveDetails.
func (r ReceiveDetails) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "brokerProperties", r.BrokerProperties)
	populate(objectMap, "event", r.Event)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ReceiveDetails.
func (r *ReceiveDetails) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "brokerProperties":
			err = unpopulate(val, "BrokerProperties", &r.BrokerProperties)
			delete(rawMsg, key)
		case "event":
			err = unpopulate(val, "Event", &r.Event)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", r, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ReceiveResult.
func (r ReceiveResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "value", r.Value)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ReceiveResult.
func (r *ReceiveResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "value":
			err = unpopulate(val, "Value", &r.Value)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", r, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type RejectOptions.
func (r RejectOptions) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "lockTokens", r.LockTokens)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type RejectOptions.
func (r *RejectOptions) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "lockTokens":
			err = unpopulate(val, "LockTokens", &r.LockTokens)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", r, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type RejectResult.
func (r RejectResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "failedLockTokens", r.FailedLockTokens)
	populate(objectMap, "succeededLockTokens", r.SucceededLockTokens)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type RejectResult.
func (r *RejectResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "failedLockTokens":
			err = unpopulate(val, "FailedLockTokens", &r.FailedLockTokens)
			delete(rawMsg, key)
		case "succeededLockTokens":
			err = unpopulate(val, "SucceededLockTokens", &r.SucceededLockTokens)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", r, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ReleaseOptions.
func (r ReleaseOptions) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "lockTokens", r.LockTokens)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ReleaseOptions.
func (r *ReleaseOptions) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "lockTokens":
			err = unpopulate(val, "LockTokens", &r.LockTokens)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", r, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ReleaseResult.
func (r ReleaseResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "failedLockTokens", r.FailedLockTokens)
	populate(objectMap, "succeededLockTokens", r.SucceededLockTokens)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ReleaseResult.
func (r *ReleaseResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", r, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "failedLockTokens":
			err = unpopulate(val, "FailedLockTokens", &r.FailedLockTokens)
			delete(rawMsg, key)
		case "succeededLockTokens":
			err = unpopulate(val, "SucceededLockTokens", &r.SucceededLockTokens)
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

func populateByteArray(m map[string]any, k string, b []byte, f runtime.Base64Encoding) {
	if azcore.IsNullValue(b) {
		m[k] = nil
	} else if len(b) == 0 {
		return
	} else {
		m[k] = runtime.EncodeByteArray(b, f)
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
