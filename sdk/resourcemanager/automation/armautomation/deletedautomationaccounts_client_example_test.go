//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armautomation_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automation/armautomation"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/main/specification/automation/resource-manager/Microsoft.Automation/stable/2022-01-31/examples/getDeletedAutomationAccount.json
func ExampleDeletedAutomationAccountsClient_ListBySubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armautomation.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDeletedAutomationAccountsClient().ListBySubscription(ctx, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DeletedAutomationAccountListResult = armautomation.DeletedAutomationAccountListResult{
	// 	Value: []*armautomation.DeletedAutomationAccount{
	// 		{
	// 			Name: to.Ptr("myAutomationAccount"),
	// 			Type: to.Ptr("Microsoft.Automation/deletedAutomationAccounts"),
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/msitest/providers/Microsoft.Automation/deletedAutomationAccounts/myAutomationAccount"),
	// 			Location: to.Ptr("westus"),
	// 			Properties: &armautomation.DeletedAutomationAccountProperties{
	// 				AutomationAccountID: to.Ptr("cb855f13-0223-4fe4-8260-9e6583dfef24"),
	// 				AutomationAccountResourceID: to.Ptr("/subscriptions/subid/resourceGroups/msitest/providers/Microsoft.Automation/automationAccounts/myAutomationAccount"),
	// 				DeletionTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-04-24T16:30:55+00:00"); return t}()),
	// 				Location: to.Ptr("westus"),
	// 			},
	// 	}},
	// }
}
