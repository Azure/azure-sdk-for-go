//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhealthcareapis_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/healthcareapis/armhealthcareapis"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/healthcareapis/resource-manager/Microsoft.HealthcareApis/preview/2022-01-31-preview/examples/fhirservices/FhirServices_List.json
func ExampleFhirServicesClient_NewListByWorkspacePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhealthcareapis.NewFhirServicesClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByWorkspacePager("testRG",
		"workspace1",
		nil)
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/healthcareapis/resource-manager/Microsoft.HealthcareApis/preview/2022-01-31-preview/examples/fhirservices/FhirServices_Get.json
func ExampleFhirServicesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhealthcareapis.NewFhirServicesClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"testRG",
		"workspace1",
		"fhirservices1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/healthcareapis/resource-manager/Microsoft.HealthcareApis/preview/2022-01-31-preview/examples/fhirservices/FhirServices_Create.json
func ExampleFhirServicesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhealthcareapis.NewFhirServicesClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreateOrUpdate(ctx,
		"testRG",
		"workspace1",
		"fhirservice1",
		armhealthcareapis.FhirService{
			Identity: &armhealthcareapis.ServiceManagedIdentityIdentity{
				Type: to.Ptr(armhealthcareapis.ServiceManagedIdentityTypeSystemAssigned),
			},
			Location: to.Ptr("westus"),
			Tags: map[string]*string{
				"additionalProp1": to.Ptr("string"),
				"additionalProp2": to.Ptr("string"),
				"additionalProp3": to.Ptr("string"),
			},
			Kind: to.Ptr(armhealthcareapis.FhirServiceKindFhirR4),
			Properties: &armhealthcareapis.FhirServiceProperties{
				AccessPolicies: []*armhealthcareapis.FhirServiceAccessPolicyEntry{
					{
						ObjectID: to.Ptr("c487e7d1-3210-41a3-8ccc-e9372b78da47"),
					},
					{
						ObjectID: to.Ptr("5b307da8-43d4-492b-8b66-b0294ade872f"),
					}},
				AcrConfiguration: &armhealthcareapis.FhirServiceAcrConfiguration{
					LoginServers: []*string{
						to.Ptr("test1.azurecr.io")},
				},
				AuthenticationConfiguration: &armhealthcareapis.FhirServiceAuthenticationConfiguration{
					Audience:          to.Ptr("https://azurehealthcareapis.com"),
					Authority:         to.Ptr("https://login.microsoftonline.com/abfde7b2-df0f-47e6-aabf-2462b07508dc"),
					SmartProxyEnabled: to.Ptr(true),
				},
				CorsConfiguration: &armhealthcareapis.FhirServiceCorsConfiguration{
					AllowCredentials: to.Ptr(false),
					Headers: []*string{
						to.Ptr("*")},
					MaxAge: to.Ptr[int32](1440),
					Methods: []*string{
						to.Ptr("DELETE"),
						to.Ptr("GET"),
						to.Ptr("OPTIONS"),
						to.Ptr("PATCH"),
						to.Ptr("POST"),
						to.Ptr("PUT")},
					Origins: []*string{
						to.Ptr("*")},
				},
				ExportConfiguration: &armhealthcareapis.FhirServiceExportConfiguration{
					StorageAccountName: to.Ptr("existingStorageAccount"),
				},
				ImportConfiguration: &armhealthcareapis.FhirServiceImportConfiguration{
					Enabled:              to.Ptr(false),
					InitialImportMode:    to.Ptr(false),
					IntegrationDataStore: to.Ptr("existingStorageAccount"),
				},
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/healthcareapis/resource-manager/Microsoft.HealthcareApis/preview/2022-01-31-preview/examples/fhirservices/FhirServices_Patch.json
func ExampleFhirServicesClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhealthcareapis.NewFhirServicesClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginUpdate(ctx,
		"testRG",
		"fhirservice1",
		"workspace1",
		armhealthcareapis.FhirServicePatchResource{
			Tags: map[string]*string{
				"tagKey": to.Ptr("tagValue"),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/healthcareapis/resource-manager/Microsoft.HealthcareApis/preview/2022-01-31-preview/examples/fhirservices/FhirServices_Delete.json
func ExampleFhirServicesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhealthcareapis.NewFhirServicesClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"testRG",
		"fhirservice1",
		"workspace1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
