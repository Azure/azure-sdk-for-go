//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armaad

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type ErrorDefinition.
func (e ErrorDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateEndpointConnectionListResult.
func (p PrivateEndpointConnectionListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", p.NextLink)
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateLinkPolicy.
func (p PrivateLinkPolicy) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "allTenants", p.AllTenants)
	populate(objectMap, "id", p.ID)
	populate(objectMap, "name", p.Name)
	populate(objectMap, "ownerTenantId", p.OwnerTenantID)
	populate(objectMap, "resourceGroup", p.ResourceGroup)
	populate(objectMap, "resourceName", p.ResourceName)
	populate(objectMap, "subscriptionId", p.SubscriptionID)
	populate(objectMap, "tags", p.Tags)
	populate(objectMap, "tenants", p.Tenants)
	populate(objectMap, "type", p.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateLinkPolicyListResult.
func (p PrivateLinkPolicyListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", p.NextLink)
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateLinkPolicyUpdateParameter.
func (p PrivateLinkPolicyUpdateParameter) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "tags", p.Tags)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateLinkResourceListResult.
func (p PrivateLinkResourceListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", p.NextLink)
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateLinkResourceProperties.
func (p PrivateLinkResourceProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "groupId", p.GroupID)
	populate(objectMap, "requiredMembers", p.RequiredMembers)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type TagsResource.
func (t TagsResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "tags", t.Tags)
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
