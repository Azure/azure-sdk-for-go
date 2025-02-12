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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskRestorePointExamples/DiskRestorePoint_Get.json
func ExampleDiskRestorePointClient_Get_getAnIncrementalDiskRestorePointResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDiskRestorePointClient().Get(ctx, "myResourceGroup", "rpc", "vmrp", "TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DiskRestorePoint = armcompute.DiskRestorePoint{
	// 	Name: to.Ptr("TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/restorePointCollections/rpc/restorePoints/vmrp/diskRestorePoints/TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745"),
	// 	Properties: &armcompute.DiskRestorePointProperties{
	// 		FamilyID: to.Ptr("996bf3ce-b6ff-4e86-9db6-dc27ea06cea5"),
	// 		HyperVGeneration: to.Ptr(armcompute.HyperVGenerationV1),
	// 		LogicalSectorSize: to.Ptr[int32](4096),
	// 		NetworkAccessPolicy: to.Ptr(armcompute.NetworkAccessPolicyAllowAll),
	// 		OSType: to.Ptr(armcompute.OperatingSystemTypesWindows),
	// 		PublicNetworkAccess: to.Ptr(armcompute.PublicNetworkAccessDisabled),
	// 		SourceResourceID: to.Ptr("/subscriptions/d2260d06-e00d-422f-8b63-93df551a59ae/resourceGroups/rg0680fb0c-89f1-41b4-96c0-35733a181558/providers/Microsoft.Compute/disks/TestDisk45ceb03433006d1baee0"),
	// 		SourceUniqueID: to.Ptr("48e058b1-7eea-4968-b532-10a8a1130c13"),
	// 		TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-09-16T04:41:35.079Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskRestorePointExamples/DiskRestorePoint_Get_WhenSourceResourceIsFromDifferentRegion.json
func ExampleDiskRestorePointClient_Get_getAnIncrementalDiskRestorePointWhenSourceResourceIsFromADifferentRegion() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDiskRestorePointClient().Get(ctx, "myResourceGroup", "rpc", "vmrp", "TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DiskRestorePoint = armcompute.DiskRestorePoint{
	// 	Name: to.Ptr("TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/restorePointCollections/rpc/restorePoints/vmrp/diskRestorePoints/TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745"),
	// 	Properties: &armcompute.DiskRestorePointProperties{
	// 		CompletionPercent: to.Ptr[float32](100),
	// 		FamilyID: to.Ptr("996bf3ce-b6ff-4e86-9db6-dc27ea06cea5"),
	// 		HyperVGeneration: to.Ptr(armcompute.HyperVGenerationV1),
	// 		LogicalSectorSize: to.Ptr[int32](4096),
	// 		NetworkAccessPolicy: to.Ptr(armcompute.NetworkAccessPolicyAllowAll),
	// 		OSType: to.Ptr(armcompute.OperatingSystemTypesWindows),
	// 		PublicNetworkAccess: to.Ptr(armcompute.PublicNetworkAccessDisabled),
	// 		ReplicationState: to.Ptr("Succeeded"),
	// 		SourceResourceID: to.Ptr("/subscriptions/d2260d06-e00d-422f-8b63-93df551a59ae/resourceGroups/rg0680fb0c-89f1-41b4-96c0-35733a181558/providers/Microsoft.Compute/disks/TestDisk45ceb03433006d1baee0"),
	// 		SourceResourceLocation: to.Ptr("eastus2"),
	// 		SourceUniqueID: to.Ptr("48e058b1-7eea-4968-b532-10a8a1130c13"),
	// 		TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-09-16T04:41:35.079Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskRestorePointExamples/DiskRestorePoint_ListByVmRestorePoint.json
func ExampleDiskRestorePointClient_NewListByRestorePointPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDiskRestorePointClient().NewListByRestorePointPager("myResourceGroup", "rpc", "vmrp", nil)
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
		// page.DiskRestorePointList = armcompute.DiskRestorePointList{
		// 	Value: []*armcompute.DiskRestorePoint{
		// 		{
		// 			Name: to.Ptr("TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/restorePointCollections/rpc/restorePoints/vmrp/diskRestorePoints/TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745"),
		// 			Properties: &armcompute.DiskRestorePointProperties{
		// 				FamilyID: to.Ptr("996bf3ce-b6ff-4e86-9db6-dc27ea06cea5"),
		// 				HyperVGeneration: to.Ptr(armcompute.HyperVGenerationV1),
		// 				LogicalSectorSize: to.Ptr[int32](4096),
		// 				NetworkAccessPolicy: to.Ptr(armcompute.NetworkAccessPolicyAllowAll),
		// 				OSType: to.Ptr(armcompute.OperatingSystemTypesWindows),
		// 				PublicNetworkAccess: to.Ptr(armcompute.PublicNetworkAccessDisabled),
		// 				SourceResourceID: to.Ptr("/subscriptions/d2260d06-e00d-422f-8b63-93df551a59ae/resourceGroups/rg0680fb0c-89f1-41b4-96c0-35733a181558/providers/Microsoft.Compute/disks/TestDisk45ceb03433006d1baee0"),
		// 				SourceUniqueID: to.Ptr("48e058b1-7eea-4968-b532-10a8a1130c13"),
		// 				TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-09-16T04:41:35.079Z"); return t}()),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskRestorePointExamples/DiskRestorePoint_BeginGetAccess.json
func ExampleDiskRestorePointClient_BeginGrantAccess() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDiskRestorePointClient().BeginGrantAccess(ctx, "myResourceGroup", "rpc", "vmrp", "TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745", armcompute.GrantAccessData{
		Access:            to.Ptr(armcompute.AccessLevelRead),
		DurationInSeconds: to.Ptr[int32](300),
		FileFormat:        to.Ptr(armcompute.FileFormatVHDX),
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
	// res.AccessURI = armcompute.AccessURI{
	// 	AccessSAS: to.Ptr("https://md-gpvmcxzlzxgd.partition.blob.storage.azure.net/xx3cqcx53f0v/abcd?sv=2014-02-14&sr=b&sk=key1&sig=XXX&st=2021-05-24T18:02:34Z&se=2021-05-24T18:19:14Z&sp=r"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskRestorePointExamples/DiskRestorePoint_EndGetAccess.json
func ExampleDiskRestorePointClient_BeginRevokeAccess() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDiskRestorePointClient().BeginRevokeAccess(ctx, "myResourceGroup", "rpc", "vmrp", "TestDisk45ceb03433006d1baee0_b70cd924-3362-4a80-93c2-9415eaa12745", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
