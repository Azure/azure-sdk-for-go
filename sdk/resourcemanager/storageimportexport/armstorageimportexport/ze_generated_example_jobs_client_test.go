//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstorageimportexport_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storageimportexport/armstorageimportexport"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storageimportexport/resource-manager/Microsoft.ImportExport/preview/2021-01-01/examples/ListJobsInSubscription.json
func ExampleJobsClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorageimportexport.NewJobsClient("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		nil, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListBySubscriptionPager(&armstorageimportexport.JobsClientListBySubscriptionOptions{Top: nil,
		Filter: nil,
	})
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storageimportexport/resource-manager/Microsoft.ImportExport/preview/2021-01-01/examples/ListJobsInResourceGroup.json
func ExampleJobsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorageimportexport.NewJobsClient("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		nil, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByResourceGroupPager("myResourceGroup",
		&armstorageimportexport.JobsClientListByResourceGroupOptions{Top: nil,
			Filter: nil,
		})
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storageimportexport/resource-manager/Microsoft.ImportExport/preview/2021-01-01/examples/GetExportJob.json
func ExampleJobsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorageimportexport.NewJobsClient("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		nil, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"myJob",
		"myResourceGroup",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storageimportexport/resource-manager/Microsoft.ImportExport/preview/2021-01-01/examples/UpdateExportJob.json
func ExampleJobsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorageimportexport.NewJobsClient("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		nil, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"myExportJob",
		"myResourceGroup",
		armstorageimportexport.UpdateJobParameters{
			Properties: &armstorageimportexport.UpdateJobParametersProperties{
				BackupDriveManifest: to.Ptr(true),
				LogLevel:            to.Ptr("Verbose"),
				State:               to.Ptr(""),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storageimportexport/resource-manager/Microsoft.ImportExport/preview/2021-01-01/examples/CreateExportJob.json
func ExampleJobsClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorageimportexport.NewJobsClient("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		nil, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Create(ctx,
		"myExportJob",
		"myResourceGroup",
		armstorageimportexport.PutJobParameters{
			Location: to.Ptr("West US"),
			Properties: &armstorageimportexport.JobDetails{
				BackupDriveManifest: to.Ptr(true),
				DiagnosticsPath:     to.Ptr("waimportexport"),
				Export: &armstorageimportexport.Export{
					BlobList: &armstorageimportexport.ExportBlobList{
						BlobPathPrefix: []*string{
							to.Ptr("/")},
					},
				},
				JobType:  to.Ptr("Export"),
				LogLevel: to.Ptr("Verbose"),
				ReturnAddress: &armstorageimportexport.ReturnAddress{
					City:            to.Ptr("Redmond"),
					CountryOrRegion: to.Ptr("USA"),
					Email:           to.Ptr("Test@contoso.com"),
					Phone:           to.Ptr("4250000000"),
					PostalCode:      to.Ptr("98007"),
					RecipientName:   to.Ptr("Test"),
					StateOrProvince: to.Ptr("wa"),
					StreetAddress1:  to.Ptr("Street1"),
					StreetAddress2:  to.Ptr("street2"),
				},
				ReturnShipping: &armstorageimportexport.ReturnShipping{
					CarrierAccountNumber: to.Ptr("989ffff"),
					CarrierName:          to.Ptr("FedEx"),
				},
				StorageAccountID: to.Ptr("/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/myResourceGroup/providers/Microsoft.ClassicStorage/storageAccounts/test"),
			},
		},
		&armstorageimportexport.JobsClientCreateOptions{ClientTenantID: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storageimportexport/resource-manager/Microsoft.ImportExport/preview/2021-01-01/examples/DeleteJob.json
func ExampleJobsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorageimportexport.NewJobsClient("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		nil, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"myJob",
		"myResourceGroup",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
