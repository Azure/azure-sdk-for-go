//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package resource

import (
	"testing"
)

func TestParseResourceType(t *testing.T) {
	resourceTypeData := map[string]struct {
		namespace    string
		resourceType string
		typesLen     int
		err          bool
	}{
		"/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines/myVmName": {
			namespace:    "Microsoft.Compute",
			resourceType: "virtualMachines",
			typesLen:     1,
		},
		"Microsoft.Compute/virtualMachines": {
			namespace:    "Microsoft.Compute",
			resourceType: "virtualMachines",
			typesLen:     1,
		},
		"/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines/myVmName/fooType/fooName": {
			namespace:    "Microsoft.Compute",
			resourceType: "virtualMachines/fooType",
			typesLen:     2,
		},
		"Microsoft.Compute/virtualMachines/fooType": {
			namespace:    "Microsoft.Compute",
			resourceType: "virtualMachines/fooType",
			typesLen:     2,
		},
		"/providers/Microsoft.Insights/providers/Microsoft.Compute/virtualMachines/myVmName": {
			namespace:    "Microsoft.Compute",
			resourceType: "virtualMachines",
			typesLen:     1,
		},
		"/providers/Microsoft.Insights/providers/Microsoft.Network/virtualNetworks/testvnet/subnets/testsubnet": {
			namespace:    "Microsoft.Network",
			resourceType: "virtualNetworks/subnets",
			typesLen:     2,
		},
		"/providers/Microsoft.Compute/virtualMachines/myVmName/fooType/fooName": {
			namespace:    "Microsoft.Compute",
			resourceType: "virtualMachines/fooType",
			typesLen:     2,
		},
		"": {
			err: true,
		},
		" ": {
			err: true,
		},
		"/": {
			err: true,
		},
	}
	for input, expected := range resourceTypeData {
		resourceType, err := ParseResourceType(input)
		if err != nil && !expected.err {
			t.Fatalf("unexpected error: %+v", err)
		}
		if expected.err {
			continue
		}
		if resourceType.Namespace != expected.namespace {
			t.Fatalf("expecting %s, but got %s", expected.namespace, resourceType.Namespace)
		}
		if resourceType.Type != expected.resourceType {
			t.Fatalf("expecting %s, but got %s", expected.resourceType, resourceType.Type)
		}
		if len(resourceType.Types) != expected.typesLen {
			t.Fatalf("expecting %d, but got %d", expected.typesLen, len(resourceType.Types))
		}
	}
}

func TestResourceType_IsParentOf(t *testing.T) {
	resourceTypes := []struct {
		left     ResourceType
		right    ResourceType
		expected bool
	}{
		{
			left:     NewResourceType("Microsoft.Compute", "virtualMachines"),
			right:    NewResourceType("Microsoft.Compute", "virtualMachines"),
			expected: false,
		},
		{
			left:     NewResourceType("Microsoft.Compute", "virtualMachines"),
			right:    NewResourceType("Microsoft.Compute", "virtualMachines/extensions"),
			expected: true,
		},
		{
			left:     NewResourceType("Microsoft.Compute", "virtualMachines"),
			right:    NewResourceType("Microsoft.Compute", "virtualMachineScaleSets/someScaleset"),
			expected: false,
		},
		{
			left:     NewResourceType("Microsoft.Network", "virtualMachines"),
			right:    NewResourceType("Microsoft.Compute", "virtualMachines"),
			expected: false,
		},
		{
			left:     NewResourceType("Microsoft.Network", "virtualNetworks"),
			right:    NewResourceType("Microsoft.Network", "virtualNetworks/subnets"),
			expected: true,
		},
		{
			left:     NewResourceType("Microsoft.Network", "virtualNetworks"),
			right:    NewResourceType("Microsoft.Network", "virtualNetworks/subnets/ipConfigurations"),
			expected: true,
		},
		{
			left:     NewResourceType("Microsoft.Network", "virtualNetworks/subnets"),
			right:    NewResourceType("Microsoft.Network", "virtualNetworks"),
			expected: false,
		},
	}

	for _, c := range resourceTypes {
		result := c.left.IsParentOf(c.right)
		if result != c.expected {
			t.Fatalf("expected %v but got %v", c.expected, result)
		}
	}
}
