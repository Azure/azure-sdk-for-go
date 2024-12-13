//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcognitiveservices_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/cognitiveservices/resource-manager/Microsoft.CognitiveServices/stable/2024-10-01/examples/ListRaiPolicies.json
func ExampleRaiPoliciesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcognitiveservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewRaiPoliciesClient().NewListPager("resourceGroupName", "accountName", nil)
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
		// page.RaiPolicyListResult = armcognitiveservices.RaiPolicyListResult{
		// 	Value: []*armcognitiveservices.RaiPolicy{
		// 		{
		// 			Name: to.Ptr("raiPolicyName"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.CognitiveServices/accounts/accountName/raiPolicies/raiPolicyName"),
		// 			Properties: &armcognitiveservices.RaiPolicyProperties{
		// 				BasePolicyName: to.Ptr("Microsoft.Default"),
		// 				ContentFilters: []*armcognitiveservices.RaiPolicyContentFilter{
		// 					{
		// 						Name: to.Ptr("Hate"),
		// 						Blocking: to.Ptr(false),
		// 						Enabled: to.Ptr(false),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
		// 					},
		// 					{
		// 						Name: to.Ptr("Hate"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
		// 					},
		// 					{
		// 						Name: to.Ptr("Sexual"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
		// 					},
		// 					{
		// 						Name: to.Ptr("Sexual"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
		// 					},
		// 					{
		// 						Name: to.Ptr("Selfharm"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
		// 					},
		// 					{
		// 						Name: to.Ptr("Selfharm"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
		// 					},
		// 					{
		// 						Name: to.Ptr("Violence"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
		// 					},
		// 					{
		// 						Name: to.Ptr("Violence"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
		// 					},
		// 					{
		// 						Name: to.Ptr("Jailbreak"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
		// 					},
		// 					{
		// 						Name: to.Ptr("Protected Material Text"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
		// 					},
		// 					{
		// 						Name: to.Ptr("Protected Material Code"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
		// 					},
		// 					{
		// 						Name: to.Ptr("Profanity"),
		// 						Blocking: to.Ptr(true),
		// 						Enabled: to.Ptr(true),
		// 						Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
		// 				}},
		// 				Mode: to.Ptr(armcognitiveservices.RaiPolicyModeAsynchronousFilter),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/cognitiveservices/resource-manager/Microsoft.CognitiveServices/stable/2024-10-01/examples/GetRaiPolicy.json
func ExampleRaiPoliciesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcognitiveservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewRaiPoliciesClient().Get(ctx, "resourceGroupName", "accountName", "raiPolicyName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.RaiPolicy = armcognitiveservices.RaiPolicy{
	// 	Name: to.Ptr("raiPolicyName"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.CognitiveServices/accounts/accountName/raiPolicies/raiPolicyName"),
	// 	Properties: &armcognitiveservices.RaiPolicyProperties{
	// 		BasePolicyName: to.Ptr("Microsoft.Default"),
	// 		ContentFilters: []*armcognitiveservices.RaiPolicyContentFilter{
	// 			{
	// 				Name: to.Ptr("Hate"),
	// 				Blocking: to.Ptr(false),
	// 				Enabled: to.Ptr(false),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Hate"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Sexual"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Sexual"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Selfharm"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Selfharm"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Violence"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Violence"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Jailbreak"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Protected Material Text"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Protected Material Code"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Profanity"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 		}},
	// 		Mode: to.Ptr(armcognitiveservices.RaiPolicyModeAsynchronousFilter),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/cognitiveservices/resource-manager/Microsoft.CognitiveServices/stable/2024-10-01/examples/PutRaiPolicy.json
func ExampleRaiPoliciesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcognitiveservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewRaiPoliciesClient().CreateOrUpdate(ctx, "resourceGroupName", "accountName", "raiPolicyName", armcognitiveservices.RaiPolicy{
		Properties: &armcognitiveservices.RaiPolicyProperties{
			BasePolicyName: to.Ptr("Microsoft.Default"),
			ContentFilters: []*armcognitiveservices.RaiPolicyContentFilter{
				{
					Name:              to.Ptr("Hate"),
					Blocking:          to.Ptr(false),
					Enabled:           to.Ptr(false),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
				},
				{
					Name:              to.Ptr("Hate"),
					Blocking:          to.Ptr(true),
					Enabled:           to.Ptr(true),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
				},
				{
					Name:              to.Ptr("Sexual"),
					Blocking:          to.Ptr(true),
					Enabled:           to.Ptr(true),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
				},
				{
					Name:              to.Ptr("Sexual"),
					Blocking:          to.Ptr(true),
					Enabled:           to.Ptr(true),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
				},
				{
					Name:              to.Ptr("Selfharm"),
					Blocking:          to.Ptr(true),
					Enabled:           to.Ptr(true),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
				},
				{
					Name:              to.Ptr("Selfharm"),
					Blocking:          to.Ptr(true),
					Enabled:           to.Ptr(true),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
				},
				{
					Name:              to.Ptr("Violence"),
					Blocking:          to.Ptr(true),
					Enabled:           to.Ptr(true),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
				},
				{
					Name:              to.Ptr("Violence"),
					Blocking:          to.Ptr(true),
					Enabled:           to.Ptr(true),
					SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
					Source:            to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
				},
				{
					Name:     to.Ptr("Jailbreak"),
					Blocking: to.Ptr(true),
					Enabled:  to.Ptr(true),
					Source:   to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
				},
				{
					Name:     to.Ptr("Protected Material Text"),
					Blocking: to.Ptr(true),
					Enabled:  to.Ptr(true),
					Source:   to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
				},
				{
					Name:     to.Ptr("Protected Material Code"),
					Blocking: to.Ptr(true),
					Enabled:  to.Ptr(true),
					Source:   to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
				},
				{
					Name:     to.Ptr("Profanity"),
					Blocking: to.Ptr(true),
					Enabled:  to.Ptr(true),
					Source:   to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
				}},
			Mode: to.Ptr(armcognitiveservices.RaiPolicyModeAsynchronousFilter),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.RaiPolicy = armcognitiveservices.RaiPolicy{
	// 	Name: to.Ptr("raiPolicyName"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.CognitiveServices/accounts/accountName/raiPolicies/raiPolicyName"),
	// 	Properties: &armcognitiveservices.RaiPolicyProperties{
	// 		BasePolicyName: to.Ptr("Microsoft.Default"),
	// 		ContentFilters: []*armcognitiveservices.RaiPolicyContentFilter{
	// 			{
	// 				Name: to.Ptr("Hate"),
	// 				Blocking: to.Ptr(false),
	// 				Enabled: to.Ptr(false),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Hate"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Sexual"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Sexual"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Selfharm"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelHigh),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Selfharm"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Violence"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Violence"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				SeverityThreshold: to.Ptr(armcognitiveservices.ContentLevelMedium),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Jailbreak"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 			},
	// 			{
	// 				Name: to.Ptr("Protected Material Text"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Protected Material Code"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourceCompletion),
	// 			},
	// 			{
	// 				Name: to.Ptr("Profanity"),
	// 				Blocking: to.Ptr(true),
	// 				Enabled: to.Ptr(true),
	// 				Source: to.Ptr(armcognitiveservices.RaiPolicyContentSourcePrompt),
	// 		}},
	// 		Mode: to.Ptr(armcognitiveservices.RaiPolicyModeAsynchronousFilter),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/cognitiveservices/resource-manager/Microsoft.CognitiveServices/stable/2024-10-01/examples/DeleteRaiPolicy.json
func ExampleRaiPoliciesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcognitiveservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewRaiPoliciesClient().BeginDelete(ctx, "resourceGroupName", "accountName", "raiPolicyName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}