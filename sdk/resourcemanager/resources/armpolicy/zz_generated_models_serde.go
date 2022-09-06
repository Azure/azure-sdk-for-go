//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpolicy

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type AssignmentProperties.
func (a AssignmentProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "description", a.Description)
	populate(objectMap, "displayName", a.DisplayName)
	populate(objectMap, "enforcementMode", a.EnforcementMode)
	populate(objectMap, "metadata", &a.Metadata)
	populate(objectMap, "nonComplianceMessages", a.NonComplianceMessages)
	populate(objectMap, "notScopes", a.NotScopes)
	populate(objectMap, "parameters", a.Parameters)
	populate(objectMap, "policyDefinitionId", a.PolicyDefinitionID)
	populate(objectMap, "scope", a.Scope)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type AssignmentUpdate.
func (a AssignmentUpdate) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "identity", a.Identity)
	populate(objectMap, "location", a.Location)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type DefinitionProperties.
func (d DefinitionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "description", d.Description)
	populate(objectMap, "displayName", d.DisplayName)
	populate(objectMap, "metadata", &d.Metadata)
	populate(objectMap, "mode", d.Mode)
	populate(objectMap, "parameters", d.Parameters)
	populate(objectMap, "policyRule", &d.PolicyRule)
	populate(objectMap, "policyType", d.PolicyType)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type DefinitionReference.
func (d DefinitionReference) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "groupNames", d.GroupNames)
	populate(objectMap, "parameters", d.Parameters)
	populate(objectMap, "policyDefinitionId", d.PolicyDefinitionID)
	populate(objectMap, "policyDefinitionReferenceId", d.PolicyDefinitionReferenceID)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ExemptionProperties.
func (e ExemptionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "description", e.Description)
	populate(objectMap, "displayName", e.DisplayName)
	populate(objectMap, "exemptionCategory", e.ExemptionCategory)
	populateTimeRFC3339(objectMap, "expiresOn", e.ExpiresOn)
	populate(objectMap, "metadata", &e.Metadata)
	populate(objectMap, "policyAssignmentId", e.PolicyAssignmentID)
	populate(objectMap, "policyDefinitionReferenceIds", e.PolicyDefinitionReferenceIDs)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ExemptionProperties.
func (e *ExemptionProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", e, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "description":
			err = unpopulate(val, "Description", &e.Description)
			delete(rawMsg, key)
		case "displayName":
			err = unpopulate(val, "DisplayName", &e.DisplayName)
			delete(rawMsg, key)
		case "exemptionCategory":
			err = unpopulate(val, "ExemptionCategory", &e.ExemptionCategory)
			delete(rawMsg, key)
		case "expiresOn":
			err = unpopulateTimeRFC3339(val, "ExpiresOn", &e.ExpiresOn)
			delete(rawMsg, key)
		case "metadata":
			err = unpopulate(val, "Metadata", &e.Metadata)
			delete(rawMsg, key)
		case "policyAssignmentId":
			err = unpopulate(val, "PolicyAssignmentID", &e.PolicyAssignmentID)
			delete(rawMsg, key)
		case "policyDefinitionReferenceIds":
			err = unpopulate(val, "PolicyDefinitionReferenceIDs", &e.PolicyDefinitionReferenceIDs)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", e, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type Identity.
func (i Identity) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "principalId", i.PrincipalID)
	populate(objectMap, "tenantId", i.TenantID)
	populate(objectMap, "type", i.Type)
	populate(objectMap, "userAssignedIdentities", i.UserAssignedIdentities)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ParameterDefinitionsValue.
func (p ParameterDefinitionsValue) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "allowedValues", p.AllowedValues)
	populate(objectMap, "defaultValue", &p.DefaultValue)
	populate(objectMap, "metadata", p.Metadata)
	populate(objectMap, "type", p.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ParameterDefinitionsValueMetadata.
func (p ParameterDefinitionsValueMetadata) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "assignPermissions", p.AssignPermissions)
	populate(objectMap, "description", p.Description)
	populate(objectMap, "displayName", p.DisplayName)
	populate(objectMap, "strongType", p.StrongType)
	if p.AdditionalProperties != nil {
		for key, val := range p.AdditionalProperties {
			objectMap[key] = val
		}
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ParameterDefinitionsValueMetadata.
func (p *ParameterDefinitionsValueMetadata) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", p, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "assignPermissions":
			err = unpopulate(val, "AssignPermissions", &p.AssignPermissions)
			delete(rawMsg, key)
		case "description":
			err = unpopulate(val, "Description", &p.Description)
			delete(rawMsg, key)
		case "displayName":
			err = unpopulate(val, "DisplayName", &p.DisplayName)
			delete(rawMsg, key)
		case "strongType":
			err = unpopulate(val, "StrongType", &p.StrongType)
			delete(rawMsg, key)
		default:
			if p.AdditionalProperties == nil {
				p.AdditionalProperties = map[string]interface{}{}
			}
			if val != nil {
				var aux interface{}
				err = json.Unmarshal(val, &aux)
				p.AdditionalProperties[key] = aux
			}
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", p, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type SetDefinitionProperties.
func (s SetDefinitionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "description", s.Description)
	populate(objectMap, "displayName", s.DisplayName)
	populate(objectMap, "metadata", &s.Metadata)
	populate(objectMap, "parameters", s.Parameters)
	populate(objectMap, "policyDefinitionGroups", s.PolicyDefinitionGroups)
	populate(objectMap, "policyDefinitions", s.PolicyDefinitions)
	populate(objectMap, "policyType", s.PolicyType)
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
		return fmt.Errorf("unmarshalling type %T: %v", s, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
			err = unpopulateTimeRFC3339(val, "CreatedAt", &s.CreatedAt)
			delete(rawMsg, key)
		case "createdBy":
			err = unpopulate(val, "CreatedBy", &s.CreatedBy)
			delete(rawMsg, key)
		case "createdByType":
			err = unpopulate(val, "CreatedByType", &s.CreatedByType)
			delete(rawMsg, key)
		case "lastModifiedAt":
			err = unpopulateTimeRFC3339(val, "LastModifiedAt", &s.LastModifiedAt)
			delete(rawMsg, key)
		case "lastModifiedBy":
			err = unpopulate(val, "LastModifiedBy", &s.LastModifiedBy)
			delete(rawMsg, key)
		case "lastModifiedByType":
			err = unpopulate(val, "LastModifiedByType", &s.LastModifiedByType)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", s, err)
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

func unpopulate(data json.RawMessage, fn string, v interface{}) error {
	if data == nil {
		return nil
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}
	return nil
}
