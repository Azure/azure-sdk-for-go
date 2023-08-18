//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armauthorization_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/310a0100f5b020c1900c527a6aa70d21992f078a/specification/authorization/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/GetAlertOperationById.json
func ExampleAlertOperationClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armauthorization.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAlertOperationClient().Get(ctx, "subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f", "{operationId}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AlertOperationResult = armauthorization.AlertOperationResult{
	// 	CreatedDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-04-05T03:04:06.467+00:00"); return t}()),
	// 	ID: to.Ptr("{operationId}"),
	// 	LastActionDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-04-05T03:04:06.467+00:00"); return t}()),
	// 	ResourceLocation: to.Ptr("https://management.azure.com/providers/Microsoft.Subscription/subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f/providers/Microsoft.Authorization/roleManagementAlerts/DuplicateRoleCreated?api-version=2022-08-01-preview"),
	// 	Status: to.Ptr("Running"),
	// 	StatusDetail: to.Ptr("{\"result\":[{\"name\":\"DuplicateRoleCreated\",\"statusDetail\":\"Alert scan is not yet started.\"}]}"),
	// }
}
