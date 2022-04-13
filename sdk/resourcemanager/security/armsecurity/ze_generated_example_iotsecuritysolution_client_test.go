//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2019-08-01/examples/IoTSecuritySolutions/GetIoTSecuritySolutionsListByIotHub.json
func ExampleIotSecuritySolutionClient_ListBySubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armsecurity.NewIotSecuritySolutionClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	pager := client.ListBySubscription(&armsecurity.IotSecuritySolutionClientListBySubscriptionOptions{Filter: to.Ptr("<filter>")})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
			return
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2019-08-01/examples/IoTSecuritySolutions/GetIoTSecuritySolutionsListByRg.json
func ExampleIotSecuritySolutionClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armsecurity.NewIotSecuritySolutionClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	pager := client.ListByResourceGroup("<resource-group-name>",
		&armsecurity.IotSecuritySolutionClientListByResourceGroupOptions{Filter: nil})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
			return
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2019-08-01/examples/IoTSecuritySolutions/GetIoTSecuritySolution.json
func ExampleIotSecuritySolutionClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armsecurity.NewIotSecuritySolutionClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<solution-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2019-08-01/examples/IoTSecuritySolutions/CreateIoTSecuritySolution.json
func ExampleIotSecuritySolutionClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armsecurity.NewIotSecuritySolutionClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.CreateOrUpdate(ctx,
		"<resource-group-name>",
		"<solution-name>",
		armsecurity.IoTSecuritySolutionModel{
			Tags:     map[string]*string{},
			Location: to.Ptr("<location>"),
			Properties: &armsecurity.IoTSecuritySolutionProperties{
				DisabledDataSources: []*armsecurity.DataSource{},
				DisplayName:         to.Ptr("<display-name>"),
				Export:              []*armsecurity.ExportData{},
				IotHubs: []*string{
					to.Ptr("/subscriptions/075423e9-7d33-4166-8bdf-3920b04e3735/resourceGroups/myRg/providers/Microsoft.Devices/IotHubs/FirstIotHub")},
				RecommendationsConfiguration: []*armsecurity.RecommendationConfigurationProperties{
					{
						RecommendationType: to.Ptr(armsecurity.RecommendationTypeIoTOpenPorts),
						Status:             to.Ptr(armsecurity.RecommendationConfigStatusDisabled),
					},
					{
						RecommendationType: to.Ptr(armsecurity.RecommendationTypeIoTSharedCredentials),
						Status:             to.Ptr(armsecurity.RecommendationConfigStatusDisabled),
					}},
				Status:                  to.Ptr(armsecurity.SecuritySolutionStatusEnabled),
				UnmaskedIPLoggingStatus: to.Ptr(armsecurity.UnmaskedIPLoggingStatusEnabled),
				UserDefinedResources: &armsecurity.UserDefinedResourcesProperties{
					Query: to.Ptr("<query>"),
					QuerySubscriptions: []*string{
						to.Ptr("075423e9-7d33-4166-8bdf-3920b04e3735")},
				},
				Workspace: to.Ptr("<workspace>"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2019-08-01/examples/IoTSecuritySolutions/UpdateIoTSecuritySolution.json
func ExampleIotSecuritySolutionClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armsecurity.NewIotSecuritySolutionClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	res, err := client.Update(ctx,
		"<resource-group-name>",
		"<solution-name>",
		armsecurity.UpdateIotSecuritySolutionData{
			Tags: map[string]*string{
				"foo": to.Ptr("bar"),
			},
			Properties: &armsecurity.UpdateIoTSecuritySolutionProperties{
				RecommendationsConfiguration: []*armsecurity.RecommendationConfigurationProperties{
					{
						RecommendationType: to.Ptr(armsecurity.RecommendationTypeIoTOpenPorts),
						Status:             to.Ptr(armsecurity.RecommendationConfigStatusDisabled),
					},
					{
						RecommendationType: to.Ptr(armsecurity.RecommendationTypeIoTSharedCredentials),
						Status:             to.Ptr(armsecurity.RecommendationConfigStatusDisabled),
					}},
				UserDefinedResources: &armsecurity.UserDefinedResourcesProperties{
					Query: to.Ptr("<query>"),
					QuerySubscriptions: []*string{
						to.Ptr("075423e9-7d33-4166-8bdf-3920b04e3735")},
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/security/resource-manager/Microsoft.Security/stable/2019-08-01/examples/IoTSecuritySolutions/DeleteIoTSecuritySolution.json
func ExampleIotSecuritySolutionClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
		return
	}
	ctx := context.Background()
	client, err := armsecurity.NewIotSecuritySolutionClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		return
	}
	_, err = client.Delete(ctx,
		"<resource-group-name>",
		"<solution-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
		return
	}
}
