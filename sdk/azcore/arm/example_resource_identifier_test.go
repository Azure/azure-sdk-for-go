//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import "fmt"

func ExampleParseResourceType() {
	rawResourceType := "Microsoft.Network/virtualNetworks/subnets"
	resourceType, err := ParseResourceType(rawResourceType)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ResourceType: %s\n", resourceType.String())
	fmt.Printf("Namespace: %s, Type: %s, LastType: %s\n", resourceType.Namespace, resourceType.Type, resourceType.LastType())

	// Output:
	// ResourceType: Microsoft.Network/virtualNetworks/subnets
	// Namespace: Microsoft.Network, Type: virtualNetworks/subnets, LastType: subnets
}

func ExampleParseResourceType_fromResourceIdentifier() {
	rawResourceId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySub"
	resourceType, err := ParseResourceType(rawResourceId)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ResourceType: %s\n", resourceType.String())
	fmt.Printf("Namespace: %s, Type: %s, LastType: %s\n", resourceType.Namespace, resourceType.Type, resourceType.LastType())

	// Output:
	// ResourceType: Microsoft.Network/virtualNetworks/subnets
	// Namespace: Microsoft.Network, Type: virtualNetworks/subnets, LastType: subnets
}

func ExampleParseResourceIdentifier() {
	rawResourceId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subsets/mySub"
	id, err := ParseResourceID(rawResourceId)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %s\n", id.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Name, id.ResourceType, id.SubscriptionId, id.ResourceGroupName)
	fmt.Printf("Parent: %s\n", id.Parent.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Parent.Name, id.Parent.ResourceType, id.Parent.SubscriptionId, id.Parent.ResourceGroupName)
	fmt.Printf("Parent: %s\n", id.Parent.Parent.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Parent.Parent.Name, id.Parent.Parent.ResourceType, id.Parent.Parent.SubscriptionId, id.Parent.Parent.ResourceGroupName)
	fmt.Printf("Parent: %s\n", id.Parent.Parent.Parent.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Parent.Parent.Parent.Name, id.Parent.Parent.Parent.ResourceType, id.Parent.Parent.Parent.SubscriptionId, id.Parent.Parent.Parent.ResourceGroupName)

	// Output:
	// ID: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subsets/mySub
	// Name: mySub, ResourceType: Microsoft.Network/virtualNetworks/subsets, SubscriptionId: 00000000-0000-0000-0000-000000000000, ResourceGroupName: myRg
	// Parent: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet
	// Name: vnet, ResourceType: Microsoft.Network/virtualNetworks, SubscriptionId: 00000000-0000-0000-0000-000000000000, ResourceGroupName: myRg
	// Parent: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg
	// Name: myRg, ResourceType: Microsoft.Resources/resourceGroups, SubscriptionId: 00000000-0000-0000-0000-000000000000, ResourceGroupName: myRg
	// Parent: /subscriptions/00000000-0000-0000-0000-000000000000
	// Name: 00000000-0000-0000-0000-000000000000, ResourceType: Microsoft.Resources/subscriptions, SubscriptionId: 00000000-0000-0000-0000-000000000000, ResourceGroupName:
}
