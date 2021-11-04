//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"fmt"
	"strings"
)

const (
	providersKey             = "providers"
	subscriptionsKey         = "subscriptions"
	resourceGroupsLowerKey   = "resourcegroups"
	locationsKey             = "locations"
	builtInResourceNamespace = "Microsoft.Resources"
)

var (
	// RootResourceIdentifier defines a ResourceID of a tenant as a root level parent of all other ResourceID
	RootResourceIdentifier = &ResourceID{
		Parent:       nil,
		ResourceType: TenantResourceType,
		Name:         "",
	}
)

// ResourceID represents a resource ID such as `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg`
// Don't create this type directly, use ParseResourceID() instead
type ResourceID struct {
	Parent         *ResourceID
	SubscriptionID string
	Provider       string
	ResourceGroupName string
	Location          string
	ResourceType      ResourceType
	Name              string

	isChild     bool
	stringValue string
}

// ParseResourceID parses a string to an instance of ResourceID
func ParseResourceID(id string) (*ResourceID, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("invalid resource id: id cannot be empty")
	}

	if !strings.HasPrefix(id, "/") {
		return nil, fmt.Errorf("invalid resource id: resource id '%s' must start with '/'", id)
	}

	parts := splitStringAndOmitEmpty(id, "/")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid resource id: %s", id)
	}

	if !strings.EqualFold(parts[0], subscriptionsKey) && !strings.EqualFold(parts[0], providersKey) {
		return nil, fmt.Errorf("invalid resource id: %s", id)
	}

	return appendNext(RootResourceIdentifier, parts, id)
}

func newResourceID(parent *ResourceID, resourceTypeName string, resourceName string) *ResourceID {
	id := &ResourceID{}
	id.init(parent, chooseResourceType(resourceTypeName, parent), resourceName, true)
	return id
}

func newResourceIDWithResourceType(parent *ResourceID, resourceType ResourceType, resourceName string) *ResourceID {
	id := &ResourceID{}
	id.init(parent, resourceType, resourceName, true)
	return id
}

func newResourceIDWithProvider(parent *ResourceID, providerNamespace, resourceTypeName, resourceName string) *ResourceID {
	id := &ResourceID{}
	id.init(parent, NewResourceType(providerNamespace, resourceTypeName), resourceName, false)
	return id
}

func chooseResourceType(resourceTypeName string, parent *ResourceID) ResourceType {
	switch strings.ToLower(resourceTypeName) {
	case resourceGroupsLowerKey:
		return ResourceGroupResourceType
	case subscriptionsKey:
		return SubscriptionResourceType
	default:
		return parent.ResourceType.AppendChild(resourceTypeName)
	}
}

func (id *ResourceID) init(parent *ResourceID, resourceType ResourceType, name string, isChild bool) {
	if parent != nil {
		id.Provider = parent.Provider
		id.SubscriptionID = parent.SubscriptionID
		id.ResourceGroupName = parent.ResourceGroupName
		id.Location = parent.Location
	}

	if resourceType.String() == SubscriptionResourceType.String() {
		id.SubscriptionID = name
	}

	if resourceType.lastType() == locationsKey {
		id.Location = name
	}

	if resourceType.String() == ResourceGroupResourceType.String() {
		id.ResourceGroupName = name
	}

	if resourceType.String() == ProviderResourceType.String() {
		id.Provider = name
	}

	if parent == nil {
		id.Parent = RootResourceIdentifier
	} else {
		id.Parent = parent
	}
	id.isChild = isChild
	id.ResourceType = resourceType
	id.Name = name
}

func appendNext(parent *ResourceID, parts []string, id string) (*ResourceID, error) {
	if len(parts) == 0 {
		return parent, nil
	}

	if len(parts) == 1 {
		// subscriptions and resourceGroups are not valid ids without their names
		if strings.EqualFold(parts[0], subscriptionsKey) || strings.EqualFold(parts[0], resourceGroupsLowerKey) {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		// resourceGroup must contain either child or provider resource type
		if parent.ResourceType.String() == ResourceGroupResourceType.String() {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		return newResourceID(parent, parts[0], ""), nil
	}

	if strings.EqualFold(parts[0], providersKey) && (len(parts) == 2 || strings.EqualFold(parts[2], providersKey)) {
		//provider resource can only be on a tenant or a subscription parent
		if parent.ResourceType.String() != SubscriptionResourceType.String() && parent.ResourceType.String() != TenantResourceType.String() {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		return appendNext(newResourceIDWithResourceType(parent, ProviderResourceType, parts[1]), parts[2:], id)
	}

	if len(parts) > 3 && strings.EqualFold(parts[0], providersKey) {
		return appendNext(newResourceIDWithProvider(parent, parts[1], parts[2], parts[3]), parts[4:], id)
	}

	if len(parts) > 1 && !strings.EqualFold(parts[0], providersKey) {
		return appendNext(newResourceID(parent, parts[0], parts[1]), parts[2:], id)
	}

	return nil, fmt.Errorf("invalid resource id: %s", id)
}

func (id ResourceID) String() string {
	if len(id.stringValue) > 0 {
		return id.stringValue
	}

	if id.Parent == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(id.Parent.String())

	if id.isChild {
		builder.WriteString(fmt.Sprintf("/%s", id.ResourceType.lastType()))
		if len(id.Name) > 0 {
			builder.WriteString(fmt.Sprintf("/%s", id.Name))
		}
	} else {
		builder.WriteString(fmt.Sprintf("/providers/%s/%s/%s", id.ResourceType.Namespace, id.ResourceType.Type, id.Name))
	}

	id.stringValue = builder.String()

	return id.stringValue
}

func splitStringAndOmitEmpty(v, sep string) []string {
	r := make([]string, 0)
	for _, s := range strings.Split(v, sep) {
		if len(s) == 0 {
			continue
		}
		r = append(r, s)
	}

	return r
}
