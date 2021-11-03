//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import "testing"

var (
	resourceTypeData = map[string]struct {
		namespace    string
		resourceType string
		typesLen     int
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
	}
)

func TestParseResourceType(t *testing.T) {
	for input, expected := range resourceTypeData {
		resourceType, err := ParseResourceType(input)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
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
