//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmarketplaceordering

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type AgreementProperties.
func (a AgreementProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "accepted", a.Accepted)
	populate(objectMap, "licenseTextLink", a.LicenseTextLink)
	populate(objectMap, "marketplaceTermsLink", a.MarketplaceTermsLink)
	populate(objectMap, "plan", a.Plan)
	populate(objectMap, "privacyPolicyLink", a.PrivacyPolicyLink)
	populate(objectMap, "product", a.Product)
	populate(objectMap, "publisher", a.Publisher)
	populateTimeRFC3339(objectMap, "retrieveDatetime", a.RetrieveDatetime)
	populate(objectMap, "signature", a.Signature)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type AgreementProperties.
func (a *AgreementProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "accepted":
			err = unpopulate(val, &a.Accepted)
			delete(rawMsg, key)
		case "licenseTextLink":
			err = unpopulate(val, &a.LicenseTextLink)
			delete(rawMsg, key)
		case "marketplaceTermsLink":
			err = unpopulate(val, &a.MarketplaceTermsLink)
			delete(rawMsg, key)
		case "plan":
			err = unpopulate(val, &a.Plan)
			delete(rawMsg, key)
		case "privacyPolicyLink":
			err = unpopulate(val, &a.PrivacyPolicyLink)
			delete(rawMsg, key)
		case "product":
			err = unpopulate(val, &a.Product)
			delete(rawMsg, key)
		case "publisher":
			err = unpopulate(val, &a.Publisher)
			delete(rawMsg, key)
		case "retrieveDatetime":
			err = unpopulateTimeRFC3339(val, &a.RetrieveDatetime)
			delete(rawMsg, key)
		case "signature":
			err = unpopulate(val, &a.Signature)
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

// MarshalJSON implements the json.Marshaller interface for type SystemData.
func (s SystemData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "createdAt", s.CreatedAt)
	populate(objectMap, "createdBy", s.CreatedBy)
	populate(objectMap, "createdByType", s.CreatedByType)
	populateTimeRFC3339(objectMap, "lastModifiedAt", s.LastModifiedAt)
	populate(objectMap, "lastModifiedBy", s.LastModifiedBy)
	populate(objectMap, "lastModifiedByType", s.LastModifiedByType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type SystemData.
func (s *SystemData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
			err = unpopulateTimeRFC3339(val, &s.CreatedAt)
			delete(rawMsg, key)
		case "createdBy":
			err = unpopulate(val, &s.CreatedBy)
			delete(rawMsg, key)
		case "createdByType":
			err = unpopulate(val, &s.CreatedByType)
			delete(rawMsg, key)
		case "lastModifiedAt":
			err = unpopulateTimeRFC3339(val, &s.LastModifiedAt)
			delete(rawMsg, key)
		case "lastModifiedBy":
			err = unpopulate(val, &s.LastModifiedBy)
			delete(rawMsg, key)
		case "lastModifiedByType":
			err = unpopulate(val, &s.LastModifiedByType)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
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
