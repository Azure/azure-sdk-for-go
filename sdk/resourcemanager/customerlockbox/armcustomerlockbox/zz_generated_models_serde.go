//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcustomerlockbox

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type ErrorBody.
func (e ErrorBody) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type LockboxRequestResponseProperties.
func (l LockboxRequestResponseProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "accessLevel", l.AccessLevel)
	populateTimeRFC3339(objectMap, "createdDateTime", l.CreatedDateTime)
	populate(objectMap, "duration", l.Duration)
	populateTimeRFC3339(objectMap, "expirationDateTime", l.ExpirationDateTime)
	populate(objectMap, "justification", l.Justification)
	populate(objectMap, "requestId", l.RequestID)
	populate(objectMap, "resourceIds", l.ResourceIDs)
	populate(objectMap, "resourceType", l.ResourceType)
	populate(objectMap, "status", l.Status)
	populate(objectMap, "subscriptionId", l.SubscriptionID)
	populate(objectMap, "supportCaseUrl", l.SupportCaseURL)
	populate(objectMap, "supportRequest", l.SupportRequest)
	populate(objectMap, "workitemsource", l.Workitemsource)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type LockboxRequestResponseProperties.
func (l *LockboxRequestResponseProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "accessLevel":
			err = unpopulate(val, &l.AccessLevel)
			delete(rawMsg, key)
		case "createdDateTime":
			err = unpopulateTimeRFC3339(val, &l.CreatedDateTime)
			delete(rawMsg, key)
		case "duration":
			err = unpopulate(val, &l.Duration)
			delete(rawMsg, key)
		case "expirationDateTime":
			err = unpopulateTimeRFC3339(val, &l.ExpirationDateTime)
			delete(rawMsg, key)
		case "justification":
			err = unpopulate(val, &l.Justification)
			delete(rawMsg, key)
		case "requestId":
			err = unpopulate(val, &l.RequestID)
			delete(rawMsg, key)
		case "resourceIds":
			err = unpopulate(val, &l.ResourceIDs)
			delete(rawMsg, key)
		case "resourceType":
			err = unpopulate(val, &l.ResourceType)
			delete(rawMsg, key)
		case "status":
			err = unpopulate(val, &l.Status)
			delete(rawMsg, key)
		case "subscriptionId":
			err = unpopulate(val, &l.SubscriptionID)
			delete(rawMsg, key)
		case "supportCaseUrl":
			err = unpopulate(val, &l.SupportCaseURL)
			delete(rawMsg, key)
		case "supportRequest":
			err = unpopulate(val, &l.SupportRequest)
			delete(rawMsg, key)
		case "workitemsource":
			err = unpopulate(val, &l.Workitemsource)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type OperationListResult.
func (o OperationListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", o.NextLink)
	populate(objectMap, "value", o.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type RequestListResult.
func (r RequestListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", r.NextLink)
	populate(objectMap, "value", r.Value)
	return json.Marshal(objectMap)
}

func populate(m map[string]interface{}, k string, v interface{}) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}
