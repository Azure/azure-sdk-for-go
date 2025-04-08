//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcompute_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/250861bb6a886b75255edfa0aa5ee2dd0d6e7a11/specification/compute/resource-manager/Microsoft.Compute/GalleryRP/stable/2024-03-03/examples/galleryResourceProfileExamples/GalleryInVMAccessControlProfile_Create.json
func ExampleGalleryInVMAccessControlProfilesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGalleryInVMAccessControlProfilesClient().BeginCreateOrUpdate(ctx, "myResourceGroup", "myGalleryName", "myInVMAccessControlProfileName", armcompute.GalleryInVMAccessControlProfile{
		Location: to.Ptr("West US"),
		Properties: &armcompute.GalleryInVMAccessControlProfileProperties{
			ApplicableHostEndpoint: to.Ptr(armcompute.EndpointTypesWireServer),
			OSType:                 to.Ptr(armcompute.OperatingSystemTypesLinux),
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
	// res.GalleryInVMAccessControlProfile = armcompute.GalleryInVMAccessControlProfile{
	// 	Name: to.Ptr("myInVMAccessControlProfileName"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/galleries/myGallery/inVMAccessControlProfiles/myInVMAccessControlProfileName"),
	// 	Location: to.Ptr("West US"),
	// 	Properties: &armcompute.GalleryInVMAccessControlProfileProperties{
	// 		ProvisioningState: to.Ptr(armcompute.GalleryProvisioningStateSucceeded),
	// 		ApplicableHostEndpoint: to.Ptr(armcompute.EndpointTypesWireServer),
	// 		OSType: to.Ptr(armcompute.OperatingSystemTypesLinux),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/250861bb6a886b75255edfa0aa5ee2dd0d6e7a11/specification/compute/resource-manager/Microsoft.Compute/GalleryRP/stable/2024-03-03/examples/galleryResourceProfileExamples/GalleryInVMAccessControlProfile_Update.json
func ExampleGalleryInVMAccessControlProfilesClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGalleryInVMAccessControlProfilesClient().BeginUpdate(ctx, "myResourceGroup", "myGalleryName", "myInVMAccessControlProfileName", armcompute.GalleryInVMAccessControlProfileUpdate{
		Properties: &armcompute.GalleryInVMAccessControlProfileProperties{
			ApplicableHostEndpoint: to.Ptr(armcompute.EndpointTypesWireServer),
			OSType:                 to.Ptr(armcompute.OperatingSystemTypesLinux),
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
	// res.GalleryInVMAccessControlProfile = armcompute.GalleryInVMAccessControlProfile{
	// 	Name: to.Ptr("myInVMAccessControlProfileName"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/galleries/myGallery/inVMAccessControlProfiles/myInVMAccessControlProfileName"),
	// 	Location: to.Ptr("West US"),
	// 	Properties: &armcompute.GalleryInVMAccessControlProfileProperties{
	// 		ProvisioningState: to.Ptr(armcompute.GalleryProvisioningStateSucceeded),
	// 		ApplicableHostEndpoint: to.Ptr(armcompute.EndpointTypesWireServer),
	// 		OSType: to.Ptr(armcompute.OperatingSystemTypesLinux),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/250861bb6a886b75255edfa0aa5ee2dd0d6e7a11/specification/compute/resource-manager/Microsoft.Compute/GalleryRP/stable/2024-03-03/examples/galleryResourceProfileExamples/GalleryInVMAccessControlProfile_Get.json
func ExampleGalleryInVMAccessControlProfilesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGalleryInVMAccessControlProfilesClient().Get(ctx, "myResourceGroup", "myGalleryName", "myInVMAccessControlProfileName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GalleryInVMAccessControlProfile = armcompute.GalleryInVMAccessControlProfile{
	// 	Name: to.Ptr("myInVMAccessControlProfileName"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/galleries/myGallery/inVMAccessControlProfiles/myInVMAccessControlProfileName"),
	// 	Location: to.Ptr("West US"),
	// 	Properties: &armcompute.GalleryInVMAccessControlProfileProperties{
	// 		ProvisioningState: to.Ptr(armcompute.GalleryProvisioningStateSucceeded),
	// 		ApplicableHostEndpoint: to.Ptr(armcompute.EndpointTypesWireServer),
	// 		OSType: to.Ptr(armcompute.OperatingSystemTypesLinux),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/250861bb6a886b75255edfa0aa5ee2dd0d6e7a11/specification/compute/resource-manager/Microsoft.Compute/GalleryRP/stable/2024-03-03/examples/galleryResourceProfileExamples/GalleryInVMAccessControlProfile_Delete.json
func ExampleGalleryInVMAccessControlProfilesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGalleryInVMAccessControlProfilesClient().BeginDelete(ctx, "myResourceGroup", "myGalleryName", "myInVMAccessControlProfileName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/250861bb6a886b75255edfa0aa5ee2dd0d6e7a11/specification/compute/resource-manager/Microsoft.Compute/GalleryRP/stable/2024-03-03/examples/galleryResourceProfileExamples/GalleryInVMAccessControlProfile_ListByGallery.json
func ExampleGalleryInVMAccessControlProfilesClient_NewListByGalleryPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGalleryInVMAccessControlProfilesClient().NewListByGalleryPager("myResourceGroup", "myGalleryName", nil)
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
		// page.GalleryInVMAccessControlProfileList = armcompute.GalleryInVMAccessControlProfileList{
		// 	Value: []*armcompute.GalleryInVMAccessControlProfile{
		// 		{
		// 			Name: to.Ptr("myInVMAccessControlProfileName"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/galleries/myGallery/inVMAccessControlProfiles/myInVMAccessControlProfileName"),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armcompute.GalleryInVMAccessControlProfileProperties{
		// 				ProvisioningState: to.Ptr(armcompute.GalleryProvisioningStateSucceeded),
		// 				ApplicableHostEndpoint: to.Ptr(armcompute.EndpointTypesWireServer),
		// 				OSType: to.Ptr(armcompute.OperatingSystemTypesLinux),
		// 			},
		// 	}},
		// }
	}
}
