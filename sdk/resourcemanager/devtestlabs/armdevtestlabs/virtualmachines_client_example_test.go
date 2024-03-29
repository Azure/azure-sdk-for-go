//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdevtestlabs_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devtestlabs/armdevtestlabs"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_List.json
func ExampleVirtualMachinesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualMachinesClient().NewListPager("resourceGroupName", "{labName}", &armdevtestlabs.VirtualMachinesClientListOptions{Expand: nil,
		Filter:  nil,
		Top:     nil,
		Orderby: nil,
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
		// page.LabVirtualMachineList = armdevtestlabs.LabVirtualMachineList{
		// 	Value: []*armdevtestlabs.LabVirtualMachine{
		// 		{
		// 			Name: to.Ptr("{vmName}"),
		// 			Type: to.Ptr("Microsoft.DevTestLab/labs/virtualMachines"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualmachines/{vmName}"),
		// 			Location: to.Ptr("{location}"),
		// 			Tags: map[string]*string{
		// 				"tagName1": to.Ptr("tagValue1"),
		// 			},
		// 			Properties: &armdevtestlabs.LabVirtualMachineProperties{
		// 				AllowClaim: to.Ptr(true),
		// 				ArtifactDeploymentStatus: &armdevtestlabs.ArtifactDeploymentStatusProperties{
		// 					ArtifactsApplied: to.Ptr[int32](0),
		// 					TotalArtifacts: to.Ptr[int32](0),
		// 				},
		// 				ComputeID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/{labName}-{vmName}-{randomSuffix}/providers/Microsoft.Compute/virtualMachines/{vmName}"),
		// 				CreatedByUser: to.Ptr(""),
		// 				CreatedByUserID: to.Ptr(""),
		// 				CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-10-01T23:53:02.483Z"); return t}()),
		// 				DataDiskParameters: []*armdevtestlabs.DataDiskProperties{
		// 				},
		// 				DisallowPublicIPAddress: to.Ptr(true),
		// 				GalleryImageReference: &armdevtestlabs.GalleryImageReference{
		// 					Offer: to.Ptr("UbuntuServer"),
		// 					OSType: to.Ptr("Linux"),
		// 					Publisher: to.Ptr("Canonical"),
		// 					SKU: to.Ptr("16.04-LTS"),
		// 					Version: to.Ptr("Latest"),
		// 				},
		// 				LabSubnetName: to.Ptr("{virtualNetworkName}Subnet"),
		// 				LabVirtualNetworkID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualnetworks/{virtualNetworkName}"),
		// 				NetworkInterface: &armdevtestlabs.NetworkInterfaceProperties{
		// 				},
		// 				OSType: to.Ptr("Linux"),
		// 				OwnerObjectID: to.Ptr(""),
		// 				OwnerUserPrincipalName: to.Ptr(""),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				Size: to.Ptr("Standard_A2_v2"),
		// 				StorageType: to.Ptr("Standard"),
		// 				UniqueIdentifier: to.Ptr("{uniqueIdentifier}"),
		// 				UserName: to.Ptr("{userName}"),
		// 				VirtualMachineCreationSource: to.Ptr(armdevtestlabs.VirtualMachineCreationSourceFromGalleryImage),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Get.json
func ExampleVirtualMachinesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachinesClient().Get(ctx, "resourceGroupName", "{labName}", "{vmName}", &armdevtestlabs.VirtualMachinesClientGetOptions{Expand: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.LabVirtualMachine = armdevtestlabs.LabVirtualMachine{
	// 	Name: to.Ptr("{vmName}"),
	// 	Type: to.Ptr("Microsoft.DevTestLab/labs/virtualMachines"),
	// 	ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualmachines/{vmName}"),
	// 	Location: to.Ptr("{location}"),
	// 	Tags: map[string]*string{
	// 		"tagName1": to.Ptr("tagValue1"),
	// 	},
	// 	Properties: &armdevtestlabs.LabVirtualMachineProperties{
	// 		AllowClaim: to.Ptr(true),
	// 		ArtifactDeploymentStatus: &armdevtestlabs.ArtifactDeploymentStatusProperties{
	// 			ArtifactsApplied: to.Ptr[int32](0),
	// 			TotalArtifacts: to.Ptr[int32](0),
	// 		},
	// 		ComputeID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/{labName}-{vmName}-{randomSuffix}/providers/Microsoft.Compute/virtualMachines/{vmName}"),
	// 		CreatedByUser: to.Ptr(""),
	// 		CreatedByUserID: to.Ptr(""),
	// 		CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-10-01T23:53:02.483Z"); return t}()),
	// 		DataDiskParameters: []*armdevtestlabs.DataDiskProperties{
	// 		},
	// 		DisallowPublicIPAddress: to.Ptr(true),
	// 		GalleryImageReference: &armdevtestlabs.GalleryImageReference{
	// 			Offer: to.Ptr("UbuntuServer"),
	// 			OSType: to.Ptr("Linux"),
	// 			Publisher: to.Ptr("Canonical"),
	// 			SKU: to.Ptr("16.04-LTS"),
	// 			Version: to.Ptr("Latest"),
	// 		},
	// 		LabSubnetName: to.Ptr("{virtualNetworkName}Subnet"),
	// 		LabVirtualNetworkID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualnetworks/{virtualNetworkName}"),
	// 		NetworkInterface: &armdevtestlabs.NetworkInterfaceProperties{
	// 		},
	// 		OSType: to.Ptr("Linux"),
	// 		OwnerObjectID: to.Ptr(""),
	// 		OwnerUserPrincipalName: to.Ptr(""),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		Size: to.Ptr("Standard_A2_v2"),
	// 		StorageType: to.Ptr("Standard"),
	// 		UniqueIdentifier: to.Ptr("{uniqueIdentifier}"),
	// 		UserName: to.Ptr("{userName}"),
	// 		VirtualMachineCreationSource: to.Ptr(armdevtestlabs.VirtualMachineCreationSourceFromGalleryImage),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_CreateOrUpdate.json
func ExampleVirtualMachinesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginCreateOrUpdate(ctx, "resourceGroupName", "{labName}", "{vmName}", armdevtestlabs.LabVirtualMachine{
		Location: to.Ptr("{location}"),
		Tags: map[string]*string{
			"tagName1": to.Ptr("tagValue1"),
		},
		Properties: &armdevtestlabs.LabVirtualMachineProperties{
			AllowClaim:              to.Ptr(true),
			DisallowPublicIPAddress: to.Ptr(true),
			GalleryImageReference: &armdevtestlabs.GalleryImageReference{
				Offer:     to.Ptr("UbuntuServer"),
				OSType:    to.Ptr("Linux"),
				Publisher: to.Ptr("Canonical"),
				SKU:       to.Ptr("16.04-LTS"),
				Version:   to.Ptr("Latest"),
			},
			LabSubnetName:       to.Ptr("{virtualNetworkName}Subnet"),
			LabVirtualNetworkID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualnetworks/{virtualNetworkName}"),
			Password:            to.Ptr("{userPassword}"),
			Size:                to.Ptr("Standard_A2_v2"),
			StorageType:         to.Ptr("Standard"),
			UserName:            to.Ptr("{userName}"),
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
	// res.LabVirtualMachine = armdevtestlabs.LabVirtualMachine{
	// 	Name: to.Ptr("{vmName}"),
	// 	Type: to.Ptr("Microsoft.DevTestLab/labs/virtualMachines"),
	// 	ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualmachines/{vmName}"),
	// 	Location: to.Ptr("{location}"),
	// 	Tags: map[string]*string{
	// 		"tagName1": to.Ptr("tagValue1"),
	// 	},
	// 	Properties: &armdevtestlabs.LabVirtualMachineProperties{
	// 		AllowClaim: to.Ptr(true),
	// 		ArtifactDeploymentStatus: &armdevtestlabs.ArtifactDeploymentStatusProperties{
	// 			ArtifactsApplied: to.Ptr[int32](0),
	// 			TotalArtifacts: to.Ptr[int32](0),
	// 		},
	// 		CreatedByUser: to.Ptr(""),
	// 		CreatedByUserID: to.Ptr(""),
	// 		CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-10-01T23:53:02.483Z"); return t}()),
	// 		DataDiskParameters: []*armdevtestlabs.DataDiskProperties{
	// 		},
	// 		DisallowPublicIPAddress: to.Ptr(true),
	// 		GalleryImageReference: &armdevtestlabs.GalleryImageReference{
	// 			Offer: to.Ptr("UbuntuServer"),
	// 			OSType: to.Ptr("Linux"),
	// 			Publisher: to.Ptr("Canonical"),
	// 			SKU: to.Ptr("16.04-LTS"),
	// 			Version: to.Ptr("Latest"),
	// 		},
	// 		LabSubnetName: to.Ptr("{virtualNetworkName}Subnet"),
	// 		LabVirtualNetworkID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualnetworks/{virtualNetworkName}"),
	// 		NetworkInterface: &armdevtestlabs.NetworkInterfaceProperties{
	// 		},
	// 		OwnerObjectID: to.Ptr(""),
	// 		OwnerUserPrincipalName: to.Ptr(""),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		Size: to.Ptr("Standard_A2_v2"),
	// 		StorageType: to.Ptr("Standard"),
	// 		UniqueIdentifier: to.Ptr("{uniqueIdentifier}"),
	// 		UserName: to.Ptr("{userName}"),
	// 		VirtualMachineCreationSource: to.Ptr(armdevtestlabs.VirtualMachineCreationSourceFromGalleryImage),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Delete.json
func ExampleVirtualMachinesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginDelete(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Update.json
func ExampleVirtualMachinesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachinesClient().Update(ctx, "resourceGroupName", "{labName}", "{vmName}", armdevtestlabs.LabVirtualMachineFragment{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.LabVirtualMachine = armdevtestlabs.LabVirtualMachine{
	// 	Name: to.Ptr("{vmName}"),
	// 	Type: to.Ptr("Microsoft.DevTestLab/labs/virtualMachines"),
	// 	ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualmachines/{vmName}"),
	// 	Location: to.Ptr("{location}"),
	// 	Tags: map[string]*string{
	// 		"tagName1": to.Ptr("tagValue1"),
	// 	},
	// 	Properties: &armdevtestlabs.LabVirtualMachineProperties{
	// 		AllowClaim: to.Ptr(true),
	// 		ArtifactDeploymentStatus: &armdevtestlabs.ArtifactDeploymentStatusProperties{
	// 			ArtifactsApplied: to.Ptr[int32](0),
	// 			TotalArtifacts: to.Ptr[int32](0),
	// 		},
	// 		ComputeID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/{labName}-{vmName}-{randomSuffix}/providers/Microsoft.Compute/virtualMachines/{vmName}"),
	// 		CreatedByUser: to.Ptr(""),
	// 		CreatedByUserID: to.Ptr(""),
	// 		CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-10-01T23:53:02.483Z"); return t}()),
	// 		DataDiskParameters: []*armdevtestlabs.DataDiskProperties{
	// 		},
	// 		DisallowPublicIPAddress: to.Ptr(true),
	// 		GalleryImageReference: &armdevtestlabs.GalleryImageReference{
	// 			Offer: to.Ptr("UbuntuServer"),
	// 			OSType: to.Ptr("Linux"),
	// 			Publisher: to.Ptr("Canonical"),
	// 			SKU: to.Ptr("16.04-LTS"),
	// 			Version: to.Ptr("Latest"),
	// 		},
	// 		LabSubnetName: to.Ptr("{virtualNetworkName}Subnet"),
	// 		LabVirtualNetworkID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualnetworks/{virtualNetworkName}"),
	// 		NetworkInterface: &armdevtestlabs.NetworkInterfaceProperties{
	// 		},
	// 		OSType: to.Ptr("Linux"),
	// 		OwnerObjectID: to.Ptr(""),
	// 		OwnerUserPrincipalName: to.Ptr(""),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		Size: to.Ptr("Standard_A2_v2"),
	// 		StorageType: to.Ptr("Standard"),
	// 		UniqueIdentifier: to.Ptr("{uniqueIdentifier}"),
	// 		UserName: to.Ptr("{userName}"),
	// 		VirtualMachineCreationSource: to.Ptr(armdevtestlabs.VirtualMachineCreationSourceFromGalleryImage),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_AddDataDisk.json
func ExampleVirtualMachinesClient_BeginAddDataDisk() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginAddDataDisk(ctx, "resourceGroupName", "{labName}", "{virtualMachineName}", armdevtestlabs.DataDiskProperties{
		AttachNewDataDiskOptions: &armdevtestlabs.AttachNewDataDiskOptions{
			DiskName:    to.Ptr("{diskName}"),
			DiskSizeGiB: to.Ptr[int32](127),
			DiskType:    to.Ptr(armdevtestlabs.StorageType("{diskType}")),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_ApplyArtifacts.json
func ExampleVirtualMachinesClient_BeginApplyArtifacts() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginApplyArtifacts(ctx, "resourceGroupName", "{labName}", "{vmName}", armdevtestlabs.ApplyArtifactsRequest{
		Artifacts: []*armdevtestlabs.ArtifactInstallProperties{
			{
				ArtifactID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/resourceGroupName/providers/Microsoft.DevTestLab/labs/{labName}/artifactSources/public repo/artifacts/windows-restart"),
			}},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Claim.json
func ExampleVirtualMachinesClient_BeginClaim() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginClaim(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_DetachDataDisk.json
func ExampleVirtualMachinesClient_BeginDetachDataDisk() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginDetachDataDisk(ctx, "resourceGroupName", "{labName}", "{virtualMachineName}", armdevtestlabs.DetachDataDiskProperties{
		ExistingLabDiskID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/microsoft.devtestlab/labs/{labName}/virtualmachines/{virtualMachineName}"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_GetRdpFileContents.json
func ExampleVirtualMachinesClient_GetRdpFileContents() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachinesClient().GetRdpFileContents(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.RdpConnection = armdevtestlabs.RdpConnection{
	// 	Contents: to.Ptr("full address:s:10.0.0.4\r\nprompt for credentials:i:1\r\nusername:s:{vmName}\\{userName}\r\n"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_ListApplicableSchedules.json
func ExampleVirtualMachinesClient_ListApplicableSchedules() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachinesClient().ListApplicableSchedules(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ApplicableSchedule = armdevtestlabs.ApplicableSchedule{
	// 	Properties: &armdevtestlabs.ApplicableScheduleProperties{
	// 		LabVMsShutdown: &armdevtestlabs.Schedule{
	// 			Name: to.Ptr("LabVmsShutdown"),
	// 			Type: to.Ptr("Microsoft.DevTestLab/labs/virtualMachines/schedules"),
	// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/Microsoft.DevTestLab/labs/{labName}/schedules/myAutoShutdownSchedule"),
	// 			Location: to.Ptr("{location}"),
	// 			Properties: &armdevtestlabs.ScheduleProperties{
	// 				CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-29T21:48:14.136Z"); return t}()),
	// 				DailyRecurrence: &armdevtestlabs.DayDetails{
	// 					Time: to.Ptr("1900"),
	// 				},
	// 				HourlyRecurrence: &armdevtestlabs.HourDetails{
	// 					Minute: to.Ptr[int32](30),
	// 				},
	// 				NotificationSettings: &armdevtestlabs.NotificationSettings{
	// 					EmailRecipient: to.Ptr("{email}"),
	// 					NotificationLocale: to.Ptr("EN"),
	// 					Status: to.Ptr(armdevtestlabs.EnableStatusEnabled),
	// 					TimeInMinutes: to.Ptr[int32](30),
	// 					WebhookURL: to.Ptr("{webhookUrl}"),
	// 				},
	// 				ProvisioningState: to.Ptr("Succeeded"),
	// 				Status: to.Ptr(armdevtestlabs.EnableStatusEnabled),
	// 				TargetResourceID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{vmName}"),
	// 				TaskType: to.Ptr("LabVmsShutdownTask"),
	// 				TimeZoneID: to.Ptr("Pacific Standard Time"),
	// 				UniqueIdentifier: to.Ptr("4acf0408-1c10-49cb-96b7-28ce655c8320"),
	// 				WeeklyRecurrence: &armdevtestlabs.WeekDetails{
	// 					Time: to.Ptr("1700"),
	// 					Weekdays: []*string{
	// 						to.Ptr("Friday"),
	// 						to.Ptr("Saturday"),
	// 						to.Ptr("Sunday")},
	// 					},
	// 				},
	// 			},
	// 			LabVMsStartup: &armdevtestlabs.Schedule{
	// 				Name: to.Ptr("LabVmAutoStart"),
	// 				Type: to.Ptr("Microsoft.DevTestLab/labs/virtualMachines/schedules"),
	// 				ID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/Microsoft.DevTestLab/labs/{labName}/schedules/myAutoStartSchedule"),
	// 				Location: to.Ptr("{location}"),
	// 				Properties: &armdevtestlabs.ScheduleProperties{
	// 					CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-29T21:46:37.047Z"); return t}()),
	// 					DailyRecurrence: &armdevtestlabs.DayDetails{
	// 						Time: to.Ptr("0900"),
	// 					},
	// 					HourlyRecurrence: &armdevtestlabs.HourDetails{
	// 						Minute: to.Ptr[int32](30),
	// 					},
	// 					NotificationSettings: &armdevtestlabs.NotificationSettings{
	// 						EmailRecipient: to.Ptr("{email}"),
	// 						NotificationLocale: to.Ptr("EN"),
	// 						Status: to.Ptr(armdevtestlabs.EnableStatusEnabled),
	// 						TimeInMinutes: to.Ptr[int32](30),
	// 						WebhookURL: to.Ptr("{webhookUrl}"),
	// 					},
	// 					ProvisioningState: to.Ptr("Succeeded"),
	// 					Status: to.Ptr(armdevtestlabs.EnableStatusEnabled),
	// 					TargetResourceID: to.Ptr("/subscriptions/{subscriptionId}/resourcegroups/resourceGroupName/providers/Microsoft.DevTestLab/labs/{labName}/virtualmachines/{vmName}"),
	// 					TaskType: to.Ptr("LabVmsStartupTask"),
	// 					TimeZoneID: to.Ptr("Pacific Standard Time"),
	// 					WeeklyRecurrence: &armdevtestlabs.WeekDetails{
	// 						Time: to.Ptr("1000"),
	// 						Weekdays: []*string{
	// 							to.Ptr("Friday"),
	// 							to.Ptr("Saturday"),
	// 							to.Ptr("Sunday")},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Redeploy.json
func ExampleVirtualMachinesClient_BeginRedeploy() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginRedeploy(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Resize.json
func ExampleVirtualMachinesClient_BeginResize() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginResize(ctx, "resourceGroupName", "{labName}", "{vmName}", armdevtestlabs.ResizeLabVirtualMachineProperties{
		Size: to.Ptr("Standard_A4_v2"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Restart.json
func ExampleVirtualMachinesClient_BeginRestart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginRestart(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Start.json
func ExampleVirtualMachinesClient_BeginStart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginStart(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_Stop.json
func ExampleVirtualMachinesClient_BeginStop() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginStop(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_TransferDisks.json
func ExampleVirtualMachinesClient_BeginTransferDisks() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginTransferDisks(ctx, "resourceGroupName", "{labName}", "{virtualmachineName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/VirtualMachines_UnClaim.json
func ExampleVirtualMachinesClient_BeginUnClaim() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginUnClaim(ctx, "resourceGroupName", "{labName}", "{vmName}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
