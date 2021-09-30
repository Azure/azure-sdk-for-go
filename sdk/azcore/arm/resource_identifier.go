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
	// RootResourceIdentifier defines a ResourceIdentifier of a tenant as a root level parent of all other ResourceIdentifier
	RootResourceIdentifier = &ResourceIdentifier{
		parent:       nil,
		resourceType: TenantResourceType,
		name:         "",
	}
)

// ResourceIdentifier represents a resource ID such as `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg`
// Don't create this type directly, use ParseResourceID() instead
type ResourceIdentifier struct {
	parent            *ResourceIdentifier
	subscriptionId    string
	provider          string
	resourceGroupName string
	location          string
	resourceType      *ResourceType
	name              string
	isChild           bool

	stringValue string
}

func (id *ResourceIdentifier) Parent() *ResourceIdentifier {
	return id.parent
}

func (id *ResourceIdentifier) SubscriptionId() string {
	return id.subscriptionId
}

func (id *ResourceIdentifier) Provider() string {
	return id.provider
}

func (id *ResourceIdentifier) ResourceGroupName() string {
	return id.resourceGroupName
}

func (id *ResourceIdentifier) Location() string {
	return id.location
}

func (id *ResourceIdentifier) ResourceType() *ResourceType {
	return id.resourceType
}

func (id *ResourceIdentifier) Name() string {
	return id.name
}

func newResourceIdentifier(parent *ResourceIdentifier, resourceTypeName string, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	id.init(parent, chooseResourceType(resourceTypeName, parent), resourceName, true)
	return id
}

func newResourceIdentifierWithResourceType(parent *ResourceIdentifier, resourceType *ResourceType, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	id.init(parent, resourceType, resourceName, true)
	return id
}

func newResourceIdentifierWithProvider(parent *ResourceIdentifier, providerNamespace, resourceTypeName, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	id.init(parent, NewResourceType(providerNamespace, resourceTypeName), resourceName, false)
	return id
}

func chooseResourceType(resourceTypeName string, parent *ResourceIdentifier) *ResourceType {
	switch strings.ToLower(resourceTypeName) {
	case resourceGroupsLowerKey:
		return ResourceGroupResourceType
	case subscriptionsKey:
		return SubscriptionResourceType
	default:
		return parent.resourceType.AppendChild(resourceTypeName)
	}
}

func (id *ResourceIdentifier) init(parent *ResourceIdentifier, resourceType *ResourceType, name string, isChild bool) {
	if parent != nil {
		id.provider = parent.provider
		id.subscriptionId = parent.subscriptionId
		id.resourceGroupName = parent.resourceGroupName
		id.location = parent.location
	}

	if resourceType.String() == SubscriptionResourceType.String() {
		id.subscriptionId = name
	}

	if resourceType.LastType() == locationsKey {
		id.location = name
	}

	if resourceType.String() == ResourceGroupResourceType.String() {
		id.resourceGroupName = name
	}

	if resourceType.String() == ProviderResourceType.String() {
		id.provider = name
	}

	if parent == nil {
		id.parent = RootResourceIdentifier
	} else {
		id.parent = parent
	}
	id.isChild = isChild
	id.resourceType = resourceType
	id.name = name
}

// ParseResourceIdentifier parses a string to an instance of ResourceIdentifier
func ParseResourceIdentifier(id string) (*ResourceIdentifier, error) {
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

func appendNext(parent *ResourceIdentifier, parts []string, id string) (*ResourceIdentifier, error) {
	if len(parts) == 0 {
		return parent, nil
	}

	if len(parts) == 1 {
		// subscriptions and resourceGroups are not valid ids without their names
		if strings.EqualFold(parts[0], subscriptionsKey) || strings.EqualFold(parts[0], resourceGroupsLowerKey) {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		// resourceGroup must contain either child or provider resource type
		if parent.resourceType.String() == ResourceGroupResourceType.String() {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		return newResourceIdentifier(parent, parts[0], ""), nil
	}

	if strings.EqualFold(parts[0], providersKey) && (len(parts) == 2 || strings.EqualFold(parts[2], providersKey)) {
		//provider resource can only be on a tenant or a subscription parent
		if parent.resourceType.String() != SubscriptionResourceType.String() && parent.resourceType.String() != TenantResourceType.String() {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		return appendNext(newResourceIdentifierWithResourceType(parent, ProviderResourceType, parts[1]), parts[2:], id)
	}

	if len(parts) > 3 && strings.EqualFold(parts[0], providersKey) {
		return appendNext(newResourceIdentifierWithProvider(parent, parts[1], parts[2], parts[3]), parts[4:], id)
	}

	if len(parts) > 1 && !strings.EqualFold(parts[0], providersKey) {
		return appendNext(newResourceIdentifier(parent, parts[0], parts[1]), parts[2:], id)
	}

	return nil, fmt.Errorf("invalid resource id: %s", id)
}

func (id ResourceIdentifier) String() string {
	if len(id.stringValue) > 0 {
		return id.stringValue
	}

	id.stringValue = id.resourceString()
	return id.stringValue
}

func (id ResourceIdentifier) resourceString() string {
	if id.parent == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(id.parent.String())

	if id.isChild {
		builder.WriteString(fmt.Sprintf("/%s", id.resourceType.LastType()))
		if len(id.name) > 0 {
			builder.WriteString(fmt.Sprintf("/%s", id.name))
		}
	} else {
		builder.WriteString(fmt.Sprintf("/providers/%s/%s/%s", id.resourceType.namespace, id.resourceType.typeName, id.name))
	}

	return builder.String()
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
