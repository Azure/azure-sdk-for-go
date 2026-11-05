// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"fmt"
	"log"
)

func ExampleResourceID_ValidateResourceTypeChain() {
	id, err := ParseResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myIdentity")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	err = id.ValidateResourceTypeChain([]ResourceType{
		TenantResourceType,
		SubscriptionResourceType,
		ResourceGroupResourceType,
		NewResourceType("Microsoft.ManagedIdentity", "userAssignedIdentities"),
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Validate the full chain before extracting fields for an operation that assumes
	// the identity belongs directly to a resource group.
	fmt.Printf("SubscriptionID: %s, ResourceGroupName: %s, Name: %s\n",
		id.SubscriptionID, id.ResourceGroupName, id.Name)
}

func ExampleResourceID_ValidateResourceTypeChain_rejectUnexpectedNesting() {
	id, err := ParseResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/parentIdentity/providers/Microsoft.ManagedIdentity/userAssignedIdentities/extensionIdentity")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Parsing succeeds, and the final type is userAssignedIdentities. However, the
	// immediate parent is another identity, not a resource group. This check rejects
	// the ID before its fields can be used to address a different, top-level identity.
	err = id.ValidateResourceTypeChain([]ResourceType{
		TenantResourceType,
		SubscriptionResourceType,
		ResourceGroupResourceType,
		NewResourceType("Microsoft.ManagedIdentity", "userAssignedIdentities"),
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func ExampleResourceID_ValidateResourceTypeChain_childResource() {
	id, err := ParseResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/myVNet/subnets/mySubnet")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Include the virtual network parent and the subnet's full resource type.
	err = id.ValidateResourceTypeChain([]ResourceType{
		TenantResourceType,
		SubscriptionResourceType,
		ResourceGroupResourceType,
		NewResourceType("Microsoft.Network", "virtualNetworks"),
		NewResourceType("Microsoft.Network", "virtualNetworks/subnets"),
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Printf("Virtual network: %s, Subnet: %s\n", id.Parent.Name, id.Name)
}

func ExampleParseResourceType() {
	rawResourceType := "Microsoft.Network/virtualNetworks/subnets"
	resourceType, err := ParseResourceType(rawResourceType)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ResourceType: %s\n", resourceType.String())
	fmt.Printf("Namespace: %s, Type: %s\n", resourceType.Namespace, resourceType.Type)

	// Output:
	// ResourceType: Microsoft.Network/virtualNetworks/subnets
	// Namespace: Microsoft.Network, Type: virtualNetworks/subnets
}

func ExampleParseResourceType_fromResourceID() {
	rawResourceID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySub"
	resourceType, err := ParseResourceType(rawResourceID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ResourceType: %s\n", resourceType.String())
	fmt.Printf("Namespace: %s, Type: %s\n", resourceType.Namespace, resourceType.Type)

	// Output:
	// ResourceType: Microsoft.Network/virtualNetworks/subnets
	// Namespace: Microsoft.Network, Type: virtualNetworks/subnets
}

func ExampleParseResourceID() {
	rawResourceID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subsets/mySub"
	id, err := ParseResourceID(rawResourceID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %s\n", id.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Name, id.ResourceType, id.SubscriptionID, id.ResourceGroupName)
	fmt.Printf("Parent: %s\n", id.Parent.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Parent.Name, id.Parent.ResourceType, id.Parent.SubscriptionID, id.Parent.ResourceGroupName)
	fmt.Printf("Parent: %s\n", id.Parent.Parent.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Parent.Parent.Name, id.Parent.Parent.ResourceType, id.Parent.Parent.SubscriptionID, id.Parent.Parent.ResourceGroupName)
	fmt.Printf("Parent: %s\n", id.Parent.Parent.Parent.String())
	fmt.Printf("Name: %s, ResourceType: %s, SubscriptionId: %s, ResourceGroupName: %s\n",
		id.Parent.Parent.Parent.Name, id.Parent.Parent.Parent.ResourceType, id.Parent.Parent.Parent.SubscriptionID, id.Parent.Parent.Parent.ResourceGroupName)

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
