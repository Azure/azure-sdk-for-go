//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybriddatamanager_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybriddatamanager/armhybriddatamanager"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/hybriddatamanager/resource-manager/Microsoft.HybridData/stable/2019-06-01/examples/JobDefinitions_ListByDataService-GET-example-71.json
func ExampleJobDefinitionsClient_NewListByDataServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybriddatamanager.NewJobDefinitionsClient("6e0219f5-327a-4365-904f-05eed4227ad7", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByDataServicePager("DataTransformation",
		"ResourceGroupForSDKTest",
		"TestAzureSDKOperations",
		&armhybriddatamanager.JobDefinitionsClientListByDataServiceOptions{Filter: nil})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/hybriddatamanager/resource-manager/Microsoft.HybridData/stable/2019-06-01/examples/JobDefinitions_Get-GET-example-81.json
func ExampleJobDefinitionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybriddatamanager.NewJobDefinitionsClient("6e0219f5-327a-4365-904f-05eed4227ad7", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"DataTransformation",
		"jobdeffromtestcode1",
		"ResourceGroupForSDKTest",
		"TestAzureSDKOperations",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/hybriddatamanager/resource-manager/Microsoft.HybridData/stable/2019-06-01/examples/JobDefinitions_CreateOrUpdate-PUT-example-83.json
func ExampleJobDefinitionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybriddatamanager.NewJobDefinitionsClient("6e0219f5-327a-4365-904f-05eed4227ad7", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreateOrUpdate(ctx,
		"DataTransformation",
		"jobdeffromtestcode1",
		"ResourceGroupForSDKTest",
		"TestAzureSDKOperations",
		armhybriddatamanager.JobDefinition{
			Properties: &armhybriddatamanager.JobDefinitionProperties{
				DataServiceInput: map[string]interface{}{
					"AzureStorageType": "Blob",
					"BackupChoice":     "UseExistingLatest",
					"ContainerName":    "containerfromtest",
					"DeviceName":       "8600-SHG0997877L71FC",
					"FileNameFilter":   "*",
					"IsDirectoryMode":  false,
					"RootDirectories": []interface{}{
						"\\",
					},
					"VolumeNames": []interface{}{
						"TestAutomation",
					},
				},
				DataSinkID:       to.Ptr("/subscriptions/6e0219f5-327a-4365-904f-05eed4227ad7/resourceGroups/ResourceGroupForSDKTest/providers/Microsoft.HybridData/dataManagers/TestAzureSDKOperations/dataStores/TestAzureStorage1"),
				DataSourceID:     to.Ptr("/subscriptions/6e0219f5-327a-4365-904f-05eed4227ad7/resourceGroups/ResourceGroupForSDKTest/providers/Microsoft.HybridData/dataManagers/TestAzureSDKOperations/dataStores/TestStorSimpleSource1"),
				RunLocation:      to.Ptr(armhybriddatamanager.RunLocationWestus),
				State:            to.Ptr(armhybriddatamanager.StateEnabled),
				UserConfirmation: to.Ptr(armhybriddatamanager.UserConfirmationRequired),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/hybriddatamanager/resource-manager/Microsoft.HybridData/stable/2019-06-01/examples/JobDefinitions_Delete-DELETE-example-81.json
func ExampleJobDefinitionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybriddatamanager.NewJobDefinitionsClient("6e0219f5-327a-4365-904f-05eed4227ad7", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"DataTransformation",
		"jobdeffromtestcode1",
		"ResourceGroupForSDKTest",
		"TestAzureSDKOperations",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/hybriddatamanager/resource-manager/Microsoft.HybridData/stable/2019-06-01/examples/JobDefinitions_Run-POST-example-132.json
func ExampleJobDefinitionsClient_BeginRun() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybriddatamanager.NewJobDefinitionsClient("6e0219f5-327a-4365-904f-05eed4227ad7", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginRun(ctx,
		"DataTransformation",
		"jobdeffromtestcode1",
		"ResourceGroupForSDKTest",
		"TestAzureSDKOperations",
		armhybriddatamanager.RunParameters{
			CustomerSecrets: []*armhybriddatamanager.CustomerSecret{},
			DataServiceInput: map[string]interface{}{
				"AzureStorageType": "Blob",
				"BackupChoice":     "UseExistingLatest",
				"ContainerName":    "containerfromtest",
				"DeviceName":       "8600-SHG0997877L71FC",
				"FileNameFilter":   "*",
				"IsDirectoryMode":  false,
				"RootDirectories": []interface{}{
					"\\",
				},
				"VolumeNames": []interface{}{
					"TestAutomation",
				},
			},
			UserConfirmation: to.Ptr(armhybriddatamanager.UserConfirmationNotRequired),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/hybriddatamanager/resource-manager/Microsoft.HybridData/stable/2019-06-01/examples/JobDefinitions_ListByDataManager-GET-example-191.json
func ExampleJobDefinitionsClient_NewListByDataManagerPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybriddatamanager.NewJobDefinitionsClient("6e0219f5-327a-4365-904f-05eed4227ad7", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByDataManagerPager("ResourceGroupForSDKTest",
		"TestAzureSDKOperations",
		&armhybriddatamanager.JobDefinitionsClientListByDataManagerOptions{Filter: nil})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}
