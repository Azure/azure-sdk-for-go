//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const (
	providersKey             = "providers"
	subscriptionsKey         = "subscriptions"
	resourceGroupsLowerKey   = "resourcegroups"
	builtInResourceNamespace = "Microsoft.Resources"

	SubscriptionResourceType  = "Microsoft.Resources/subscriptions"
	ResourceGroupResourceType = "Microsoft.Resources/resourceGroups"
	TenantResourceType        = "Microsoft.Resources/tenants"
	ProvidersResourceType     = "Microsoft.Resources/providers"
)

var (
	RootResourceIdentifier = &ResourceIdentifier{
		Parent:       nil,
		ResourceType: TenantResourceType,
		Name:         "",
	}
)

// ResourceIdentifier represents a resource ID
type ResourceIdentifier struct {
	Parent            *ResourceIdentifier
	SubscriptionId    string
	Provider          string
	ResourceGroupName string
	ResourceType      string
	Name              string
	IsChild           bool
}

func newResourceIdentifier(parent *ResourceIdentifier, resourceTypeName string, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	_ = id.init(parent, chooseResourceType(resourceTypeName, parent), resourceName, true)
	return id
}

func newResourceIdentifierWithProvider(parent *ResourceIdentifier, providerNamespace, resourceTypeName, resourceName string) *ResourceIdentifier {
	id := &ResourceIdentifier{}
	_ = id.init(parent, fmt.Sprintf("%s/%s", providerNamespace, resourceTypeName), resourceName, false)
	return id
}

func chooseResourceType(resourceTypeName string, parent *ResourceIdentifier) string {
	switch strings.ToLower(resourceTypeName) {
	case resourceGroupsLowerKey:
		return ResourceGroupResourceType
	case subscriptionsKey:
		return SubscriptionResourceType
	default:
		return fmt.Sprintf("%s/%s", parent.ResourceType, resourceTypeName)
	}
}

func (id *ResourceIdentifier) init(parent *ResourceIdentifier, resourceType string, name string, isChild bool) error {
	if parent != nil {
		id.Provider = parent.Provider
		id.SubscriptionId = parent.SubscriptionId
		id.ResourceGroupName = parent.ResourceGroupName
	}

	if resourceType == SubscriptionResourceType {
		_, err := uuid.FromString(name)
		if err != nil {
			return err
		}
		id.SubscriptionId = name
	}

	if resourceType == ResourceGroupResourceType {
		id.ResourceGroupName = name
	}

	if resourceType == ProvidersResourceType {
		id.Provider = name
	}

	if parent == nil {
		id.Parent = RootResourceIdentifier
	} else {
		id.Parent = parent
	}
	id.IsChild = isChild
	id.ResourceType = resourceType
	id.Name = name

	return nil
}

func ParseResourceIdentifier(id string) (*ResourceIdentifier, error) {
	if len(id) == 0 {
		return nil, errors.New("id cannot be empty")
	}

	if !strings.HasPrefix(id, "/") {
		return nil, fmt.Errorf("resource id '%s' must start with '/'", id)
	}

	parts := make([]string, 0)
	for _, part := range strings.Split(id, "/") {
		if len(part) == 0 {
			continue
		}
		parts = append(parts, part)
	}

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
		if parent.ResourceType == ResourceGroupResourceType {
			return nil, fmt.Errorf("invalid resource id: %s", id)
		}

		return newResourceIdentifier(parent, parts[0], ""), nil
	}

	if len(parts) > 3 && strings.ToLower(parts[0]) == providersKey {
		return appendNext(newResourceIdentifierWithProvider(parent, parts[1], parts[2], parts[3]), parts[4:], id)
	}

	if len(parts) > 1 && strings.ToLower(parts[0]) == providersKey {
		return appendNext(newResourceIdentifier(parent, parts[0], parts[1]), parts[2:], id)
	}

	return nil, fmt.Errorf("invalid resource id: %s", id)
}
