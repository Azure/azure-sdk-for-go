//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/GetApplicationGatewayWafDynamicManifestsDefault.json
func ExampleApplicationGatewayWafDynamicManifestsDefaultClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewApplicationGatewayWafDynamicManifestsDefaultClient().Get(ctx, "westus", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ApplicationGatewayWafDynamicManifestResult = armnetwork.ApplicationGatewayWafDynamicManifestResult{
	// 	Name: to.Ptr("default"),
	// 	Type: to.Ptr("Microsoft.Network/applicationGatewayWafDynamicManifest"),
	// 	ID: to.Ptr("/subscriptions/subid/providers/Microsoft.Network/applicationGatewayWafDynamicManifests/default"),
	// 	Properties: &armnetwork.ApplicationGatewayWafDynamicManifestPropertiesResult{
	// 		AvailableRuleSets: []*armnetwork.ApplicationGatewayFirewallManifestRuleSet{
	// 			{
	// 				RuleGroups: []*armnetwork.ApplicationGatewayFirewallRuleGroup{
	// 					{
	// 						Description: to.Ptr(""),
	// 						RuleGroupName: to.Ptr("General"),
	// 						Rules: []*armnetwork.ApplicationGatewayFirewallRule{
	// 							{
	// 								Description: to.Ptr("Failed to Parse Request Body."),
	// 								Action: to.Ptr(armnetwork.ApplicationGatewayWafRuleActionTypesAnomalyScoring),
	// 								RuleID: to.Ptr[int32](200002),
	// 								RuleIDString: to.Ptr("200002"),
	// 								State: to.Ptr(armnetwork.ApplicationGatewayWafRuleStateTypesEnabled),
	// 							},
	// 							{
	// 								Description: to.Ptr("Multipart Request Body Strict Validation."),
	// 								Action: to.Ptr(armnetwork.ApplicationGatewayWafRuleActionTypesAnomalyScoring),
	// 								RuleID: to.Ptr[int32](200003),
	// 								RuleIDString: to.Ptr("200003"),
	// 								State: to.Ptr(armnetwork.ApplicationGatewayWafRuleStateTypesEnabled),
	// 							},
	// 							{
	// 								Description: to.Ptr("Possible Multipart Unmatched Boundary."),
	// 								Action: to.Ptr(armnetwork.ApplicationGatewayWafRuleActionTypesAnomalyScoring),
	// 								RuleID: to.Ptr[int32](200004),
	// 								RuleIDString: to.Ptr("200004"),
	// 								State: to.Ptr(armnetwork.ApplicationGatewayWafRuleStateTypesEnabled),
	// 						}},
	// 				}},
	// 				RuleSetType: to.Ptr("OWASP"),
	// 				RuleSetVersion: to.Ptr("3.2"),
	// 				Status: to.Ptr(armnetwork.ApplicationGatewayRuleSetStatusOptions("0")),
	// 				Tiers: []*armnetwork.ApplicationGatewayTierTypes{
	// 					to.Ptr(armnetwork.ApplicationGatewayTierTypesWAFV2)},
	// 			}},
	// 			DefaultRuleSet: &armnetwork.DefaultRuleSetPropertyFormat{
	// 				RuleSetType: to.Ptr("OWASP"),
	// 				RuleSetVersion: to.Ptr("3.2"),
	// 			},
	// 		},
	// 	}
}
