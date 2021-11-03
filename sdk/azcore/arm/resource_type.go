//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"fmt"
	"strings"
)

var (
	// SubscriptionResourceType is the ResourceType of a subscription
	SubscriptionResourceType = NewResourceType(builtInResourceNamespace, "subscriptions")
	// ResourceGroupResourceType is the ResourceType of a resource group
	ResourceGroupResourceType = NewResourceType(builtInResourceNamespace, "resourceGroups")
	// TenantResourceType is the ResourceType of a tenant
	TenantResourceType = NewResourceType(builtInResourceNamespace, "tenants")
	// ProviderResourceType is the ResourceType of a provider
	ProviderResourceType = NewResourceType(builtInResourceNamespace, "providers")
)

// ResourceType represents an Azure resource type, e.g. "Microsoft.Network/virtualNetworks/subnets"
type ResourceType struct {
	Namespace string
	Type      string
	Types     []string

	stringValue string
}

func (t *ResourceType) LastType() string {
	return t.Types[len(t.Types)-1]
}

func (t *ResourceType) String() string {
	return t.stringValue
}

// IsParentOf returns true when the receiver is the parent resource type of the child.
func (t *ResourceType) IsParentOf(child *ResourceType) bool {
	if child == nil {
		return false
	}
	if !strings.EqualFold(t.Namespace, child.Namespace) {
		return false
	}
	var types = splitStringAndOmitEmpty(t.Type, "/")
	var childTypes = splitStringAndOmitEmpty(child.Type, "/")
	if len(types) >= len(childTypes) {
		return false
	}
	for i := range types {
		if !strings.EqualFold(types[i], childTypes[i]) {
			return false
		}
	}

	return true
}

// NewResourceType initiate a simple instance of ResourceType using provider namespace such as "Microsoft.Network" and
// Type such as "virtualNetworks/subnets"
func NewResourceType(providerNamespace, typeName string) ResourceType {
	return ResourceType{
		Namespace:   providerNamespace,
		Type:        typeName,
		Types:       splitStringAndOmitEmpty(typeName, "/"),
		stringValue: fmt.Sprintf("%s/%s", providerNamespace, typeName),
	}
}

// AppendChild initiate an instance using the receiver ResourceType as parent and append childType to it.
func (t ResourceType) AppendChild(childType string) ResourceType {
	return NewResourceType(t.Namespace, fmt.Sprintf("%s/%s", t.Type, childType))
}

// ParseResourceType parses the ResourceType from a resource type string (e.g. Microsoft.Network/virtualNetworks/subsets)
// or a resource identifier string (e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySubnet)
func ParseResourceType(resourceIdOrType string) (ResourceType, error) {
	// split the path into segments
	parts := splitStringAndOmitEmpty(resourceIdOrType, "/")

	// There must be at least a namespace and type name
	if len(parts) < 1 {
		return ResourceType{}, fmt.Errorf("invalid resource id or type: %s", resourceIdOrType)
	}

	// if the type is just subscriptions, it is a built-in type in the Microsoft.Resources namespace
	if len(parts) == 1 {
		// Simple resource type
		return NewResourceType(builtInResourceNamespace, parts[0]), nil
	} else if strings.Contains(parts[0], ".") {
		// Handle resource types (Microsoft.Compute/virtualMachines, Microsoft.Network/virtualNetworks/subnets)
		// it is a full type name
		return NewResourceType(parts[0], strings.Join(parts[1:], "/")), nil
	} else {
		// Check if ResourceID
		id, err := ParseResourceID(resourceIdOrType)
		if err != nil {
			return ResourceType{}, fmt.Errorf("invalid resource id: %s", resourceIdOrType)
		}
		return NewResourceType(id.ResourceType.Namespace, id.ResourceType.Type), nil
	}
}
