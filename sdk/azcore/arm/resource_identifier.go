//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

const (
	providersKey             = "providers"
	subscriptionsKey         = "subscriptions"
	resourceGroupsLowerKey   = "resourcegroups"
	locationsKey             = "locations"
	builtInResourceNamespace = "Microsoft.Resources"
)

var (
	RootResourceIdentifier = &ResourceIdentifier{
		parent:       nil,
		resourceType: TenantResourceType,
		name:         "",
	}

	SubscriptionResourceType  = NewResourceType("Microsoft.Resources", "subscriptions")
	ResourceGroupResourceType = NewResourceType("Microsoft.Resources", "resourceGroups")
	TenantResourceType        = NewResourceType("Microsoft.Resources", "tenants")
	ProvidersResourceType     = NewResourceType("Microsoft.Resources", "providers")
)

// ResourceIdentifier represents a resource ID
type ResourceIdentifier struct {
	parent            *ResourceIdentifier
	subscriptionId    string
	provider          string
	resourceGroupName string
	location          string
	resourceType      ResourceType
	name              string
	isChild           bool

	stringValue *string
}

func (id ResourceIdentifier) Parent() *ResourceIdentifier {
	return id.parent
}

func (id ResourceIdentifier) SubscriptionId() string {
	return id.subscriptionId
}

func (id ResourceIdentifier) Provider() string {
	return id.provider
}

func (id ResourceIdentifier) ResourceGroupName() string {
	return id.resourceGroupName
}

func (id ResourceIdentifier) Location() string {
	return id.location
}

func (id ResourceIdentifier) ResourceType() ResourceType {
	return id.resourceType
}

func (id ResourceIdentifier) Name() string {
	return id.name
}

func newResourceIdentifier(parent *ResourceIdentifier, resourceTypeName string, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	id.init(parent, chooseResourceType(resourceTypeName, parent), resourceName, true)
	return id
}

func newResourceIdentifierWithResourceType(parent *ResourceIdentifier, resourceType ResourceType, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	id.init(parent, resourceType, resourceName, true)
	return id
}

func newResourceIdentifierWithProvider(parent *ResourceIdentifier, providerNamespace, resourceTypeName, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	id.init(parent, NewResourceType(providerNamespace, resourceTypeName), resourceName, false)
	return id
}

func chooseResourceType(resourceTypeName string, parent *ResourceIdentifier) ResourceType {
	switch strings.ToLower(resourceTypeName) {
	case resourceGroupsLowerKey:
		return ResourceGroupResourceType
	case subscriptionsKey:
		return SubscriptionResourceType
	default:
		return NewResourceTypeFromParent(parent.resourceType, resourceTypeName)
	}
}

func (id *ResourceIdentifier) init(parent *ResourceIdentifier, resourceType ResourceType, name string, isChild bool) {
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

	if resourceType.String() == ProvidersResourceType.String() {
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

	firstToLower := strings.ToLower(parts[0])
	if firstToLower != subscriptionsKey && firstToLower != providersKey {
		return nil, fmt.Errorf("invalid resource id: %s", id)
	}

	return appendNext(RootResourceIdentifier, parts, id)
}

func appendNext(parent *ResourceIdentifier, parts []string, id string) (*ResourceIdentifier, error) {
	if len(parts) == 0 {
		return parent, nil
	}

	lowerFirstPart := strings.ToLower(parts[0])

	if len(parts) == 1 {
		// subscriptions and resourceGroups are not valid ids without their names
		if lowerFirstPart == subscriptionsKey || lowerFirstPart == resourceGroupsLowerKey {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		// resourceGroup must contain either child or provider resource type
		if parent.resourceType.String() == ResourceGroupResourceType.String() {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		return newResourceIdentifier(parent, parts[0], ""), nil
	}

	if lowerFirstPart == providersKey && (len(parts) == 2 || strings.ToLower(parts[2]) == providersKey) {
		//provider resource can only be on a tenant or a subscription parent
		if parent.resourceType.String() != SubscriptionResourceType.String() && parent.resourceType.String() != TenantResourceType.String() {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		return appendNext(newResourceIdentifierWithResourceType(parent, ProvidersResourceType, parts[1]), parts[2:], id)
	}

	if len(parts) > 3 && strings.ToLower(parts[0]) == providersKey {
		return appendNext(newResourceIdentifierWithProvider(parent, parts[1], parts[2], parts[3]), parts[4:], id)
	}

	if len(parts) > 1 && strings.ToLower(parts[0]) != providersKey {
		return appendNext(newResourceIdentifier(parent, parts[0], parts[1]), parts[2:], id)
	}

	return nil, fmt.Errorf("invalid resource id: %s", id)
}

func (id ResourceIdentifier) String() string {
	if id.stringValue != nil {
		return *id.stringValue
	}

	id.stringValue = to.StringPtr(id.toResourceString())
	return *id.stringValue
}

func (id ResourceIdentifier) toResourceString() string {
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
		builder.WriteString(fmt.Sprintf("/providers/%s/%s/%s", id.resourceType.namespace, id.resourceType.t, id.name))
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

type ResourceType struct {
	namespace string
	t         string
	types     []string

	stringValue string
}

func (t ResourceType) Namespace() string {
	return t.namespace
}

func (t ResourceType) Type() string {
	return t.t
}

func (t ResourceType) LastType() string {
	return t.types[len(t.types)-1]
}

func (t ResourceType) String() string {
	return t.stringValue
}

func NewResourceType(providerNamespace, name string) ResourceType {
	return ResourceType{
		namespace:   providerNamespace,
		t:           name,
		types:       splitStringAndOmitEmpty(name, "/"),
		stringValue: fmt.Sprintf("%s/%s", providerNamespace, name),
	}
}

func NewResourceTypeFromParent(parent ResourceType, childType string) ResourceType {
	return NewResourceType(parent.namespace, fmt.Sprintf("%s/%s", parent.t, childType))
}

func ParseResourceType(resourceIdOrType string) (*ResourceType, error) {
	// split the path into segments
	parts := splitStringAndOmitEmpty(resourceIdOrType, "/")

	// There must be at least a namespace and type name
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid resource id or type: %s", resourceIdOrType)
	}

	resourceType := ResourceType{}
	// if the type is just subscriptions, it is a built-in type in the Microsoft.Resources namespace
	if len(parts) == 1 {
		// Simple resource type
		resourceType = NewResourceType(builtInResourceNamespace, parts[0])
		return &resourceType, nil
	} else if strings.Contains(parts[0], ".") {
		// Handle resource types (Microsoft.Compute/virtualMachines, Microsoft.Network/virtualNetworks/subnets)
		// Type
		// it is a full type name
		resourceType = NewResourceType(parts[0], strings.Join(parts[1:], "/"))
		return &resourceType, nil
	} else {
		// Check if ResourceIdentifier
		id, err := ParseResourceIdentifier(resourceIdOrType)
		if err != nil {
			return nil, fmt.Errorf("invalid resource id: %s", resourceIdOrType)
		}
		resourceType = NewResourceType(id.resourceType.namespace, id.resourceType.t)
		return &resourceType, nil
	}
}
