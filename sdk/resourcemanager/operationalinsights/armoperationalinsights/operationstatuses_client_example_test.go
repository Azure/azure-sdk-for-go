//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armoperationalinsights_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/operationalinsights/resource-manager/Microsoft.OperationalInsights/stable/2020-08-01/examples/OperationStatusesGet.json
func ExampleOperationStatusesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armoperationalinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewOperationStatusesClient().Get(ctx, "West US", "713192d7-503f-477a-9cfe-4efc3ee2bd11", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OperationStatus = armoperationalinsights.OperationStatus{
	// 	Name: to.Ptr("713192d7-503f-477a-9cfe-4efc3ee2bd11"),
	// 	EndTime: to.Ptr("2017-01-01T16:13:13.933Z"),
	// 	ID: to.Ptr("/subscriptions/613192d7-503f-477a-9cfe-4efc3ee2bd60/locations/westus/operationStatuses/713192d7-503f-477a-9cfe-4efc3ee2bd11"),
	// 	StartTime: to.Ptr("2017-01-01T13:13:13.933Z"),
	// 	Status: to.Ptr("Succeeded"),
	// }
}
