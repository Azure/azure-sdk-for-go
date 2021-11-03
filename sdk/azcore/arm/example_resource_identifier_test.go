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
	fmt.Printf("Namespace: %s, Type: %s, LastType: %s\n", resourceType.NamespaceOld(), resourceType.Type(), resourceType.LastType())

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
	fmt.Printf("Namespace: %s, Type: %s, LastType: %s\n", resourceType.NamespaceOld(), resourceType.Type(), resourceType.LastType())

	// Output:
	// ResourceType: Microsoft.Network/virtualNetworks/subnets
	// Namespace: Microsoft.Network, Type: virtualNetworks/subnets, LastType: subnets
}

func ExampleParseResourceIdentifier() {
	rawResourceId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subsets/mySub"
	id, err := ParseResourceIdentifier(rawResourceId)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %s\n", id.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.NameOld(), id.ResourceTypeOld(), id.SubscriptionIdOld(), id.ResourceGroupNameOld())
	fmt.Printf("Parent: %s\n", id.ParentOld().String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.ParentOld().NameOld(), id.ParentOld().ResourceTypeOld(), id.ParentOld().SubscriptionIdOld(), id.ParentOld().ResourceGroupNameOld())
	fmt.Printf("Parent: %s\n", id.ParentOld().ParentOld().String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.ParentOld().ParentOld().NameOld(), id.ParentOld().ParentOld().ResourceTypeOld(), id.ParentOld().ParentOld().SubscriptionIdOld(), id.ParentOld().ParentOld().ResourceGroupNameOld())
	fmt.Printf("Parent: %s\n", id.ParentOld().ParentOld().ParentOld().String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.ParentOld().ParentOld().ParentOld().NameOld(), id.ParentOld().ParentOld().ParentOld().ResourceTypeOld(), id.ParentOld().ParentOld().ParentOld().SubscriptionIdOld(), id.ParentOld().ParentOld().ParentOld().ResourceGroupNameOld())

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
