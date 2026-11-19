// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceIDValidateResourceTypeChain(t *testing.T) {
	const subscription = "/subscriptions/00000000-0000-0000-0000-000000000000"
	const group = subscription + "/resourceGroups/myRg"
	const identity = "/providers/Microsoft.ManagedIdentity/userAssignedIdentities/"
	identityType := NewResourceType("Microsoft.ManagedIdentity", "userAssignedIdentities")
	groupIdentity := []ResourceType{TenantResourceType, SubscriptionResourceType, ResourceGroupResourceType, identityType}
	tests := []struct {
		name    string
		input   string
		types   []ResourceType
		wantErr string
	}{
		{
			name:  "resource group resource",
			input: group + identity + "bar",
			types: groupIdentity,
		},
		{
			name:    "reject nested identity with same leaf type",
			input:   group + identity + "foo" + identity + "bar",
			types:   groupIdentity,
			wantErr: `expected resource type "Microsoft.Resources/resourceGroups" at ancestry index 2, got "Microsoft.ManagedIdentity/userAssignedIdentities"`,
		},
		{
			name:  "explicit extension ancestry",
			input: group + identity + "foo" + identity + "bar",
			types: []ResourceType{TenantResourceType, SubscriptionResourceType, ResourceGroupResourceType, identityType, identityType},
		},
		{
			name:    "wrong leaf type",
			input:   group + "/providers/Microsoft.Compute/virtualMachines/vm",
			types:   groupIdentity,
			wantErr: `expected resource type "Microsoft.ManagedIdentity/userAssignedIdentities" at ancestry index 3, got "Microsoft.Compute/virtualMachines"`,
		},
		{
			name:    "wrong scope",
			input:   subscription + identity + "bar",
			types:   groupIdentity,
			wantErr: `expected resource type "Microsoft.Resources/resourceGroups" at ancestry index 2, got "Microsoft.Resources/subscriptions"`,
		},
		{
			name:  "case insensitive",
			input: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/myRg/providers/MICROSOFT.MANAGEDIDENTITY/USERASSIGNEDIDENTITIES/bar",
			types: []ResourceType{
				NewResourceType("microsoft.resources", "TENANTS"),
				NewResourceType("microsoft.resources", "SUBSCRIPTIONS"),
				NewResourceType("microsoft.resources", "RESOURCEGROUPS"),
				NewResourceType("microsoft.managedidentity", "userassignedidentities"),
			},
		},
		{
			name:  "ordinary child resource",
			input: group + "/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet",
			types: []ResourceType{
				TenantResourceType, SubscriptionResourceType, ResourceGroupResourceType,
				NewResourceType("Microsoft.Network", "virtualNetworks"),
				NewResourceType("Microsoft.Network", "virtualNetworks/subnets"),
			},
		},
		{
			name:  "tenant scoped resource",
			input: "/providers/Microsoft.Management/managementGroups/myGroup",
			types: []ResourceType{TenantResourceType, NewResourceType("Microsoft.Management", "managementGroups")},
		},
		{
			name:  "subscription scoped resource",
			input: subscription + "/providers/Microsoft.Authorization/roleAssignments/role",
			types: []ResourceType{TenantResourceType, SubscriptionResourceType, NewResourceType("Microsoft.Authorization", "roleAssignments")},
		},
		{
			name:  "subscription",
			input: subscription,
			types: []ResourceType{TenantResourceType, SubscriptionResourceType},
		},
		{
			name:  "resource group",
			input: group,
			types: []ResourceType{TenantResourceType, SubscriptionResourceType, ResourceGroupResourceType},
		},
		{
			name:    "extra ancestors",
			input:   group + identity + "bar",
			types:   []ResourceType{identityType},
			wantErr: `unexpected ancestor of resource type "Microsoft.Resources/resourceGroups"`,
		},
		{
			name:    "missing ancestor",
			input:   group + identity + "bar",
			types:   []ResourceType{TenantResourceType, TenantResourceType, SubscriptionResourceType, ResourceGroupResourceType, identityType},
			wantErr: `missing resource type "Microsoft.Resources/tenants" at ancestry index 0`,
		},
		{
			name:    "empty chain",
			input:   group + identity + "bar",
			wantErr: "resource types cannot be empty",
		},
		{
			name:    "zero resource type",
			input:   group + identity + "bar",
			types:   []ResourceType{{}},
			wantErr: "resource type at ancestry index 0 must have a namespace and type",
		},
		{
			name:    "empty namespace",
			input:   group + identity + "bar",
			types:   []ResourceType{NewResourceType("", "userAssignedIdentities")},
			wantErr: "resource type at ancestry index 0 must have a namespace and type",
		},
		{
			name:    "empty type",
			input:   group + identity + "bar",
			types:   []ResourceType{NewResourceType("Microsoft.ManagedIdentity", "")},
			wantErr: "resource type at ancestry index 0 must have a namespace and type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ParseResourceID(tt.input)
			require.NoError(t, err)
			before := id.String()
			typesBefore := append([]ResourceType(nil), tt.types...)
			err = id.ValidateResourceTypeChain(tt.types)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
			assert.Equal(t, before, id.String())
			assert.Equal(t, typesBefore, tt.types)
		})
	}
}

func TestResourceIDValidateResourceTypeChainRootAndNil(t *testing.T) {
	expected := []ResourceType{TenantResourceType}
	assert.NoError(t, RootResourceID.ValidateResourceTypeChain(expected))
	var id *ResourceID
	assert.ErrorContains(t, id.ValidateResourceTypeChain(expected), `missing resource type "Microsoft.Resources/tenants"`)
	assert.ErrorContains(t, new(ResourceID).ValidateResourceTypeChain(expected), `expected resource type "Microsoft.Resources/tenants"`)
}
