//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcontainerregistry_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/dc4c1eaef16e0bc8b1e96c3d1e014deb96259b35/specification/containerregistry/resource-manager/Microsoft.ContainerRegistry/preview/2025-03-01-preview/examples/ExportPipelineList.json
func ExampleExportPipelinesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerregistry.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewExportPipelinesClient().NewListPager("myResourceGroup", "myRegistry", nil)
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
		// page.ExportPipelineListResult = armcontainerregistry.ExportPipelineListResult{
		// 	Value: []*armcontainerregistry.ExportPipeline{
		// 		{
		// 			Name: to.Ptr("myExportPipeline"),
		// 			Type: to.Ptr("Microsoft.ContainerRegistry/registries/exportPipelines"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.ContainerRegistry/registries/myRegistry/exportPipelines/myExportPipeline"),
		// 			Identity: &armcontainerregistry.IdentityProperties{
		// 				Type: to.Ptr(armcontainerregistry.ResourceIdentityTypeSystemAssigned),
		// 				PrincipalID: to.Ptr("fa153151-b9fd-46f4-9088-5e6600f2689v"),
		// 				TenantID: to.Ptr("f686d426-8d16-42db-81b7-abu4gm510ccd"),
		// 			},
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armcontainerregistry.ExportPipelineProperties{
		// 				Options: []*armcontainerregistry.PipelineOptions{
		// 					to.Ptr(armcontainerregistry.PipelineOptionsOverwriteBlobs)},
		// 					ProvisioningState: to.Ptr(armcontainerregistry.ProvisioningStateSucceeded),
		// 					Target: &armcontainerregistry.ExportPipelineTargetProperties{
		// 						Type: to.Ptr("AzureStorageBlobContainer"),
		// 						KeyVaultURI: to.Ptr("https://myvault.vault.azure.net/secrets/acrexportsas"),
		// 						URI: to.Ptr("https://accountname.blob.core.windows.net/containername"),
		// 					},
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/dc4c1eaef16e0bc8b1e96c3d1e014deb96259b35/specification/containerregistry/resource-manager/Microsoft.ContainerRegistry/preview/2025-03-01-preview/examples/ExportPipelineGet.json
func ExampleExportPipelinesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerregistry.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewExportPipelinesClient().Get(ctx, "myResourceGroup", "myRegistry", "myExportPipeline", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ExportPipeline = armcontainerregistry.ExportPipeline{
	// 	Name: to.Ptr("myExportPipeline"),
	// 	Type: to.Ptr("Microsoft.ContainerRegistry/registries/exportPipelines"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.ContainerRegistry/registries/myRegistry/exportPipelines/myExportPipeline"),
	// 	Identity: &armcontainerregistry.IdentityProperties{
	// 		Type: to.Ptr(armcontainerregistry.ResourceIdentityTypeSystemAssigned),
	// 		PrincipalID: to.Ptr("fa153151-b9fd-46f4-9088-5e6600f2689v"),
	// 		TenantID: to.Ptr("f686d426-8d16-42db-81b7-abu4gm510ccd"),
	// 	},
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armcontainerregistry.ExportPipelineProperties{
	// 		Options: []*armcontainerregistry.PipelineOptions{
	// 			to.Ptr(armcontainerregistry.PipelineOptionsOverwriteBlobs)},
	// 			ProvisioningState: to.Ptr(armcontainerregistry.ProvisioningStateSucceeded),
	// 			Target: &armcontainerregistry.ExportPipelineTargetProperties{
	// 				Type: to.Ptr("AzureStorageBlobContainer"),
	// 				KeyVaultURI: to.Ptr("https://myvault.vault.azure.net/secrets/acrexportsas"),
	// 				URI: to.Ptr("https://accountname.blob.core.windows.net/containername"),
	// 			},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/dc4c1eaef16e0bc8b1e96c3d1e014deb96259b35/specification/containerregistry/resource-manager/Microsoft.ContainerRegistry/preview/2025-03-01-preview/examples/ExportPipelineCreate.json
func ExampleExportPipelinesClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerregistry.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExportPipelinesClient().BeginCreate(ctx, "myResourceGroup", "myRegistry", "myExportPipeline", armcontainerregistry.ExportPipeline{
		Identity: &armcontainerregistry.IdentityProperties{
			Type: to.Ptr(armcontainerregistry.ResourceIdentityTypeSystemAssigned),
		},
		Location: to.Ptr("westus"),
		Properties: &armcontainerregistry.ExportPipelineProperties{
			Options: []*armcontainerregistry.PipelineOptions{
				to.Ptr(armcontainerregistry.PipelineOptionsOverwriteBlobs)},
			Target: &armcontainerregistry.ExportPipelineTargetProperties{
				Type:        to.Ptr("AzureStorageBlobContainer"),
				KeyVaultURI: to.Ptr("https://myvault.vault.azure.net/secrets/acrexportsas"),
				URI:         to.Ptr("https://accountname.blob.core.windows.net/containername"),
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
	// res.ExportPipeline = armcontainerregistry.ExportPipeline{
	// 	Name: to.Ptr("myExportPipeline"),
	// 	Type: to.Ptr("Microsoft.ContainerRegistry/registries/exportPipelines"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.ContainerRegistry/registries/myRegistry/exportPipelines/myExportPipeline"),
	// 	Identity: &armcontainerregistry.IdentityProperties{
	// 		Type: to.Ptr(armcontainerregistry.ResourceIdentityTypeSystemAssigned),
	// 		PrincipalID: to.Ptr("fa153151-b9fd-46f4-9088-5e6600f2689v"),
	// 		TenantID: to.Ptr("f686d426-8d16-42db-81b7-abu4gm510ccd"),
	// 	},
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armcontainerregistry.ExportPipelineProperties{
	// 		Options: []*armcontainerregistry.PipelineOptions{
	// 			to.Ptr(armcontainerregistry.PipelineOptionsOverwriteBlobs)},
	// 			ProvisioningState: to.Ptr(armcontainerregistry.ProvisioningStateSucceeded),
	// 			Target: &armcontainerregistry.ExportPipelineTargetProperties{
	// 				Type: to.Ptr("AzureStorageBlobContainer"),
	// 				KeyVaultURI: to.Ptr("https://myvault.vault.azure.net/secrets/acrexportsas"),
	// 				URI: to.Ptr("https://accountname.blob.core.windows.net/containername"),
	// 			},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/dc4c1eaef16e0bc8b1e96c3d1e014deb96259b35/specification/containerregistry/resource-manager/Microsoft.ContainerRegistry/preview/2025-03-01-preview/examples/ExportPipelineDelete.json
func ExampleExportPipelinesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerregistry.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExportPipelinesClient().BeginDelete(ctx, "myResourceGroup", "myRegistry", "myExportPipeline", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
