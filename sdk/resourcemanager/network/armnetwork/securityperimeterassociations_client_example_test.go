//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspAssociationGet.json
func ExampleSecurityPerimeterAssociationsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecurityPerimeterAssociationsClient().Get(ctx, "rg1", "nsp1", "association1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NspAssociation = armnetwork.NspAssociation{
	// 	Name: to.Ptr("association1"),
	// 	Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/resourceAssociations"),
	// 	ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/resourceAssociations/association1"),
	// 	SystemData: &armnetwork.SecurityPerimeterSystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		CreatedBy: to.Ptr("user"),
	// 		CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user"),
	// 		LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 	},
	// 	Properties: &armnetwork.NspAssociationProperties{
	// 		AccessMode: to.Ptr(armnetwork.AssociationAccessModeEnforced),
	// 		HasProvisioningIssues: to.Ptr("no"),
	// 		PrivateLinkResource: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/{paasSubscriptionId}/resourceGroups/{paasResourceGroupName}/providers/{providerName}/{resourceType}/{resourceName}"),
	// 		},
	// 		Profile: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/profiles/{profileName}"),
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.NspProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspAssociationPut.json
func ExampleSecurityPerimeterAssociationsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSecurityPerimeterAssociationsClient().BeginCreateOrUpdate(ctx, "rg1", "nsp1", "association1", armnetwork.NspAssociation{
		Properties: &armnetwork.NspAssociationProperties{
			AccessMode: to.Ptr(armnetwork.AssociationAccessModeEnforced),
			PrivateLinkResource: &armnetwork.SubResource{
				ID: to.Ptr("/subscriptions/{paasSubscriptionId}/resourceGroups/{paasResourceGroupName}/providers/{providerName}/{resourceType}/{resourceName}"),
			},
			Profile: &armnetwork.SubResource{
				ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/profiles/{profileName}"),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NspAssociation = armnetwork.NspAssociation{
	// 	Name: to.Ptr("association1"),
	// 	Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/resourceAssociations"),
	// 	ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/resourceAssociations/association1"),
	// 	SystemData: &armnetwork.SecurityPerimeterSystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		CreatedBy: to.Ptr("user"),
	// 		CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user"),
	// 		LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 	},
	// 	Properties: &armnetwork.NspAssociationProperties{
	// 		AccessMode: to.Ptr(armnetwork.AssociationAccessModeEnforced),
	// 		HasProvisioningIssues: to.Ptr("no"),
	// 		PrivateLinkResource: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/{paasSubscriptionId}/resourceGroups/{paasResourceGroupName}/providers/{providerName}/{resourceType}/{resourceName}"),
	// 		},
	// 		Profile: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/profiles/{profileName}"),
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.NspProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspAssociationDelete.json
func ExampleSecurityPerimeterAssociationsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSecurityPerimeterAssociationsClient().BeginDelete(ctx, "rg1", "nsp1", "association1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspAssociationList.json
func ExampleSecurityPerimeterAssociationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSecurityPerimeterAssociationsClient().NewListPager("rg1", "nsp1", &armnetwork.SecurityPerimeterAssociationsClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.NspAssociationsListResult = armnetwork.NspAssociationsListResult{
		// 	Value: []*armnetwork.NspAssociation{
		// 		{
		// 			Name: to.Ptr("association1"),
		// 			Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/resourceAssociations"),
		// 			ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/resourceAssociations/association1"),
		// 			SystemData: &armnetwork.SecurityPerimeterSystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				CreatedBy: to.Ptr("user"),
		// 				CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user"),
		// 				LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 			},
		// 			Properties: &armnetwork.NspAssociationProperties{
		// 				AccessMode: to.Ptr(armnetwork.AssociationAccessModeEnforced),
		// 				HasProvisioningIssues: to.Ptr("no"),
		// 				PrivateLinkResource: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/{paasSubscriptionId}/resourceGroups/{paasResourceGroupName}/providers/{providerName}/{resourceType}/{resourceName}"),
		// 				},
		// 				Profile: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/profiles/{profileName}"),
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.NspProvisioningStateSucceeded),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("association2"),
		// 			Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/resourceAssociations"),
		// 			ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/resourceAssociations/association2"),
		// 			SystemData: &armnetwork.SecurityPerimeterSystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				CreatedBy: to.Ptr("user"),
		// 				CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user"),
		// 				LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 			},
		// 			Properties: &armnetwork.NspAssociationProperties{
		// 				AccessMode: to.Ptr(armnetwork.AssociationAccessModeAudit),
		// 				HasProvisioningIssues: to.Ptr("no"),
		// 				PrivateLinkResource: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/{paasSubscriptionId}/resourceGroups/{paasResourceGroupName}/providers/{providerName}/{resourceType}/{resourceName}"),
		// 				},
		// 				Profile: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/profiles/{profileName}"),
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.NspProvisioningStateSucceeded),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("association3"),
		// 			Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/resourceAssociations"),
		// 			ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/resourceAssociations/association3"),
		// 			SystemData: &armnetwork.SecurityPerimeterSystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				CreatedBy: to.Ptr("user"),
		// 				CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user"),
		// 				LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 			},
		// 			Properties: &armnetwork.NspAssociationProperties{
		// 				AccessMode: to.Ptr(armnetwork.AssociationAccessModeLearning),
		// 				HasProvisioningIssues: to.Ptr("yes"),
		// 				PrivateLinkResource: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/{paasSubscriptionId}/resourceGroups/{paasResourceGroupName}/providers/{providerName}/{resourceType}/{resourceName}"),
		// 				},
		// 				Profile: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/profiles/{profileName}"),
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.NspProvisioningStateSucceeded),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspAssociationReconcile.json
func ExampleSecurityPerimeterAssociationsClient_Reconcile() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecurityPerimeterAssociationsClient().Reconcile(ctx, "rg1", "nsp1", "association1", map[string]any{
		"properties": map[string]any{},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Interface = map[string]any{
	// }
}
