//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import "testing"

var (
	testData = map[string]*ResourceIdentifier{
		// valid resource identifiers
		"/subscriptions/db1ab6f0-4769-4b27-930e-01e2ef9c123c": {
			parent:         RootResourceIdentifier,
			subscriptionId: "db1ab6f0-4769-4b27-930e-01e2ef9c123c",
			resourceType:   SubscriptionResourceType,
			name:           "db1ab6f0-4769-4b27-930e-01e2ef9c123c",
			isChild:        true,
		},
		"/providers/Microsoft.Billing/billingAccounts/3984c6f4-2d2a-4b04-93ce-43cf4824b698%3Ae2f1492a-a492-468d-909f-bf7fe6662c01_2019-05-31": {
			parent:       RootResourceIdentifier,
			resourceType: NewResourceType("Microsoft.Billing", "billingAccounts"),
			name:         "3984c6f4-2d2a-4b04-93ce-43cf4824b698%3Ae2f1492a-a492-468d-909f-bf7fe6662c01_2019-05-31",
		},
		"/subscriptions/db1ab6f0-4769-4b27-930e-01e2ef9c123c/providers/microsoft.insights": {
			parent: &ResourceIdentifier{
				parent:         RootResourceIdentifier,
				subscriptionId: "db1ab6f0-4769-4b27-930e-01e2ef9c123c",
				resourceType:   SubscriptionResourceType,
				name:           "db1ab6f0-4769-4b27-930e-01e2ef9c123c",
				isChild:        true,
			},
			subscriptionId: "db1ab6f0-4769-4b27-930e-01e2ef9c123c",
			provider:       "microsoft.insights",
			resourceType:   ProvidersResourceType,
			name:           "microsoft.insights",
			isChild:        true,
		},
		"/subscriptions/0c2f6471-1bf0-4dda-aec3-cb9272f09575/resourceGroups/myRg/providers/Microsoft.Compute/virtualMachines/myVm": {
			parent: &ResourceIdentifier{
				parent: &ResourceIdentifier{
					parent:         RootResourceIdentifier,
					subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
					resourceType:   SubscriptionResourceType,
					name:           "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
					isChild:        true,
				},
				subscriptionId:    "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
				resourceGroupName: "myRg",
				resourceType:      ResourceGroupResourceType,
				name:              "myRg",
				isChild:           true,
			},
			subscriptionId:    "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
			resourceGroupName: "myRg",
			resourceType:      NewResourceType("Microsoft.Compute", "virtualMachines"),
			name:              "myVm",
			isChild:           false,
		},
		"/subscriptions/0c2f6471-1bf0-4dda-aec3-cb9272f09575/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/myNet/subnets/mySubnet": {
			parent: &ResourceIdentifier{
				parent: &ResourceIdentifier{
					parent: &ResourceIdentifier{
						parent:         RootResourceIdentifier,
						subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
						resourceType:   SubscriptionResourceType,
						name:           "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
						isChild:        true,
					},
					subscriptionId:    "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
					resourceGroupName: "myRg",
					resourceType:      ResourceGroupResourceType,
					name:              "myRg",
					isChild:           true,
				},
				subscriptionId:    "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
				resourceGroupName: "myRg",
				resourceType:      NewResourceType("Microsoft.Network", "virtualNetworks"),
				name:              "myNet",
				isChild:           false,
			},
			subscriptionId:    "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
			resourceGroupName: "myRg",
			resourceType:      NewResourceType("Microsoft.Network", "virtualNetworks/subnets"),
			name:              "mySubnet",
			isChild:           true,
		},
		"/subscriptions/0c2f6471-1bf0-4dda-aec3-cb9272f09575/resourceGroups/myRg": {
			parent: &ResourceIdentifier{
				parent:         RootResourceIdentifier,
				subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
				resourceType:   SubscriptionResourceType,
				name:           "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
				isChild:        true,
			},
			subscriptionId:    "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
			resourceGroupName: "myRg",
			resourceType:      ResourceGroupResourceType,
			name:              "myRg",
			isChild:           true,
		},
		"/subscriptions/0c2f6471-1bf0-4dda-aec3-cb9272f09575/locations/MyLocation": {
			parent: &ResourceIdentifier{
				parent:         RootResourceIdentifier,
				subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
				resourceType:   SubscriptionResourceType,
				name:           "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
				isChild:        true,
			},
			subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
			resourceType:   NewResourceTypeFromParent(SubscriptionResourceType, locationsKey),
			name:           "MyLocation",
			location:       "MyLocation",
			isChild:        true,
		},
		"/subscriptions/0c2f6471-1bf0-4dda-aec3-cb9272f09575/locations/MyLocation/providers/Microsoft.Authorization/roleAssignments/myRa": {
			parent: &ResourceIdentifier{
				parent: &ResourceIdentifier{
					parent:         RootResourceIdentifier,
					subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
					resourceType:   SubscriptionResourceType,
					name:           "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
					isChild:        true,
				},
				subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
				resourceType:   NewResourceTypeFromParent(SubscriptionResourceType, locationsKey),
				name:           "MyLocation",
				location:       "MyLocation",
				isChild:        true,
			},
			subscriptionId: "0c2f6471-1bf0-4dda-aec3-cb9272f09575",
			resourceType:   NewResourceType("Microsoft.Authorization", "roleAssignments"),
			name:           "myRa",
			isChild:        false,
		},
		// invalid resource identifiers
		"/providers/MicrosoftSomething/billingAccounts/":             nil,
		"/MicrosoftSomething/billingAccounts/":                       nil,
		"providers/subscription/MicrosoftSomething/billingAccounts/": nil,
		"/subscription/providersSomething":                           nil,
		"/providers":                                                 nil,
		"":                                                           nil,
		" ":                                                          nil,
		"asdfghj":                                                    nil,
		"123456":                                                     nil,
		"!@#$%^&*/":                                                  nil,
		"/subscriptions/":                                            nil,
		"/0c2f6471-1bf0-4dda-aec3-cb9272f09575/myRg/":                                   nil,
		"/providers/Company.MyProvider/myResources/myResourceName/providers/incomplete": nil,
	}
)

func TestParseResourceIdentifier(t *testing.T) {
	for input, expected := range testData {
		t.Logf("testing %s...", input)
		id, err := ParseResourceIdentifier(input)
		if err != nil && expected != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if err != nil && expected == nil {
			continue
		}
		if !equals(id, expected) {
			t.Fatalf("resource id not identical, get %v, expected %v", *id, *expected)
		}
	}
}

func equals(left, right *ResourceIdentifier) bool {
	if left != nil && right != nil {
		if left.String() != right.String() {
			return false
		}
		fieldEquals := left.Name() == right.Name() &&
			left.Provider() == right.Provider() &&
			left.ResourceType().String() == right.ResourceType().String() &&
			left.SubscriptionId() == right.SubscriptionId() &&
			left.ResourceGroupName() == right.ResourceGroupName()
		if !fieldEquals {
			return false
		}
		return equals(left.parent, right.parent)
	}

	return left == right
}
