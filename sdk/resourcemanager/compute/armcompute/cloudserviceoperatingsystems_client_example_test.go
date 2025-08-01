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

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4517f89a8ebd2f6a94e107e5ee60fff9886f3612/specification/compute/resource-manager/Microsoft.Compute/CloudserviceRP/stable/2024-11-04/examples/CloudServiceOSVersion_Get.json
func ExampleCloudServiceOperatingSystemsClient_GetOSVersion() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCloudServiceOperatingSystemsClient().GetOSVersion(ctx, "westus2", "WA-GUEST-OS-3.90_202010-02", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OSVersion = armcompute.OSVersion{
	// 	Name: to.Ptr("WA-GUEST-OS-3.90_202010-02"),
	// 	Type: to.Ptr("Microsoft.Compute/locations/cloudServiceOsVersions"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Compute/locations/westus2/cloudServiceOSVersions/WA-GUEST-OS-3.90_202010-02"),
	// 	Location: to.Ptr("westus2"),
	// 	Properties: &armcompute.OSVersionProperties{
	// 		Family: to.Ptr("3"),
	// 		FamilyLabel: to.Ptr("Windows Server 2012"),
	// 		IsActive: to.Ptr(true),
	// 		IsDefault: to.Ptr(true),
	// 		Label: to.Ptr("Windows Azure Guest OS 3.90 (Release 202010-02)"),
	// 		Version: to.Ptr("WA-GUEST-OS-3.90_202010-02"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4517f89a8ebd2f6a94e107e5ee60fff9886f3612/specification/compute/resource-manager/Microsoft.Compute/CloudserviceRP/stable/2024-11-04/examples/CloudServiceOSVersion_List.json
func ExampleCloudServiceOperatingSystemsClient_NewListOSVersionsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCloudServiceOperatingSystemsClient().NewListOSVersionsPager("westus2", nil)
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
		// page.OSVersionListResult = armcompute.OSVersionListResult{
		// 	Value: []*armcompute.OSVersion{
		// 		{
		// 			Name: to.Ptr("WA-GUEST-OS-3.90_202010-02"),
		// 			Type: to.Ptr("Microsoft.Compute/locations/cloudServiceOsVersions"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Compute/locations/westus2/cloudServiceOSVersions/WA-GUEST-OS-3.90_202010-02"),
		// 			Location: to.Ptr("westus2"),
		// 			Properties: &armcompute.OSVersionProperties{
		// 				Family: to.Ptr("3"),
		// 				FamilyLabel: to.Ptr("Windows Server 2012"),
		// 				IsActive: to.Ptr(true),
		// 				IsDefault: to.Ptr(true),
		// 				Label: to.Ptr("Windows Azure Guest OS 3.90 (Release 202010-02)"),
		// 				Version: to.Ptr("WA-GUEST-OS-3.90_202010-02"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("WA-GUEST-OS-4.83_202010-02"),
		// 			Type: to.Ptr("Microsoft.Compute/locations/cloudServiceOsVersions"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Compute/locations/westus2/cloudServiceOSVersions/WA-GUEST-OS-4.83_202010-02"),
		// 			Location: to.Ptr("westus2"),
		// 			Properties: &armcompute.OSVersionProperties{
		// 				Family: to.Ptr("4"),
		// 				FamilyLabel: to.Ptr("Windows Server 2012 R2"),
		// 				IsActive: to.Ptr(true),
		// 				IsDefault: to.Ptr(true),
		// 				Label: to.Ptr("Windows Azure Guest OS 4.83 (Release 202010-02)"),
		// 				Version: to.Ptr("WA-GUEST-OS-4.83_202010-02"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4517f89a8ebd2f6a94e107e5ee60fff9886f3612/specification/compute/resource-manager/Microsoft.Compute/CloudserviceRP/stable/2024-11-04/examples/CloudServiceOSFamily_Get.json
func ExampleCloudServiceOperatingSystemsClient_GetOSFamily() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCloudServiceOperatingSystemsClient().GetOSFamily(ctx, "westus2", "3", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OSFamily = armcompute.OSFamily{
	// 	Name: to.Ptr("3"),
	// 	Type: to.Ptr("Microsoft.Compute/locations/cloudServiceOsFamilies"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Compute/locations/westus2/cloudServiceOSFamilies/3"),
	// 	Location: to.Ptr("westus2"),
	// 	Properties: &armcompute.OSFamilyProperties{
	// 		Name: to.Ptr("3"),
	// 		Label: to.Ptr("Windows Server 2012"),
	// 		Versions: []*armcompute.OSVersionPropertiesBase{
	// 			{
	// 				IsActive: to.Ptr(true),
	// 				IsDefault: to.Ptr(true),
	// 				Label: to.Ptr("Windows Azure Guest OS 3.90 (Release 202010-02)"),
	// 				Version: to.Ptr("WA-GUEST-OS-3.90_202010-02"),
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4517f89a8ebd2f6a94e107e5ee60fff9886f3612/specification/compute/resource-manager/Microsoft.Compute/CloudserviceRP/stable/2024-11-04/examples/CloudServiceOSFamilies_List.json
func ExampleCloudServiceOperatingSystemsClient_NewListOSFamiliesPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCloudServiceOperatingSystemsClient().NewListOSFamiliesPager("westus2", nil)
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
		// page.OSFamilyListResult = armcompute.OSFamilyListResult{
		// 	Value: []*armcompute.OSFamily{
		// 		{
		// 			Name: to.Ptr("3"),
		// 			Type: to.Ptr("Microsoft.Compute/locations/cloudServiceOsFamilies"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Compute/locations/westus2/cloudServiceOSFamilies/3"),
		// 			Location: to.Ptr("westus2"),
		// 			Properties: &armcompute.OSFamilyProperties{
		// 				Name: to.Ptr("3"),
		// 				Label: to.Ptr("Windows Server 2012"),
		// 				Versions: []*armcompute.OSVersionPropertiesBase{
		// 					{
		// 						IsActive: to.Ptr(true),
		// 						IsDefault: to.Ptr(true),
		// 						Label: to.Ptr("Windows Azure Guest OS 3.90 (Release 202010-02)"),
		// 						Version: to.Ptr("WA-GUEST-OS-3.90_202010-02"),
		// 				}},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("4"),
		// 			Type: to.Ptr("Microsoft.Compute/locations/cloudServiceOsFamilies"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Compute/locations/westus2/cloudServiceOSFamilies/4"),
		// 			Location: to.Ptr("westus2"),
		// 			Properties: &armcompute.OSFamilyProperties{
		// 				Name: to.Ptr("4"),
		// 				Label: to.Ptr("Windows Server 2012 R2"),
		// 				Versions: []*armcompute.OSVersionPropertiesBase{
		// 					{
		// 						IsActive: to.Ptr(true),
		// 						IsDefault: to.Ptr(true),
		// 						Label: to.Ptr("Windows Azure Guest OS 4.83 (Release 202010-02)"),
		// 						Version: to.Ptr("WA-GUEST-OS-4.83_202010-02"),
		// 				}},
		// 			},
		// 	}},
		// }
	}
}
