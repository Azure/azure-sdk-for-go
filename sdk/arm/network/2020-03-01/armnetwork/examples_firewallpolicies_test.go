// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armnetwork_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-03-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

// This sample requires the following environment variables to be set correctly in order to run:
// AZURE_LOCATION: Azure location where the resource will be created.
// AZURE_RESOURCE_GROUP: Azure resource group to retrieve and create resources.
// AZURE_SUBSCRIPTION_ID: The subscription ID where the resource group exists.
// AZURE_TENANT_ID: The Azure Active Directory tenant (directory) ID of the service principal.
// AZURE_CLIENT_ID: The client (application) ID of the service principal.
// AZURE_CLIENT_SECRET: A client secret that was generated for the App Registration used to authenticate the client.

const (
	policyName    = "samplepolicy"
	ruleGroupName = "sampleRuleGroup"
)

func getFirewallPolicyRuleGroupsOperations() armnetwork.FirewallPolicyRuleGroupsOperations {
	return armnetwork.NewFirewallPolicyRuleGroupsClient(
		armnetwork.NewDefaultClient(getCredential(), nil),
		subscriptionID)
}

func getFirewallPoliciesOperations() armnetwork.FirewallPoliciesOperations {
	return armnetwork.NewFirewallPoliciesClient(
		armnetwork.NewDefaultClient(getCredential(), nil),
		subscriptionID)
}

// Environment variables that are required to run this example:
// AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_RESOURCE_GROUP, AZURE_LOCATION
func ExampleFirewallPolicyRuleGroupsOperations_BeginCreateOrUpdate() {
	// get FirewallPoliciesOperations and create a new FirewallPolicy to use in the FirewallPolicyRuleGroup
	fwPolicy := createFirewallPolicy(resourceGroupName, location, policyName)
	fwPolicyName := *fwPolicy.Name
	// get FirewallPolicyRuleGroupsOperations and create a new FirewallPolicyRuleGroup using the FirewallPolicy that was previously created
	fwClient := getFirewallPolicyRuleGroupsOperations()
	fwResp, err := fwClient.BeginCreateOrUpdate(
		context.Background(),
		resourceGroupName,
		fwPolicyName,
		ruleGroupName,
		armnetwork.FirewallPolicyRuleGroup{
			Properties: &armnetwork.FirewallPolicyRuleGroupProperties{
				Priority: to.Int32Ptr(110),
				Rules: &[]armnetwork.FirewallPolicyRuleClassification{
					&armnetwork.FirewallPolicyFilterRule{
						FirewallPolicyRule: armnetwork.FirewallPolicyRule{
							Priority: to.Int32Ptr(110),
							Name:     to.StringPtr("rule1"),
							RuleType: armnetwork.FirewallPolicyRuleTypeFirewallPolicyFilterRule.ToPtr(),
						},
						Action: &armnetwork.FirewallPolicyFilterRuleAction{
							Type: armnetwork.FirewallPolicyFilterRuleActionTypeDeny.ToPtr(),
						},
					},
				},
			},
		},
		nil)
	if err != nil {
		panic(err)
	}
	res, err := fwResp.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.RawResponse.StatusCode)
	// Output:
	// 200
}

// Environment variables that are required to run this example:
// AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_RESOURCE_GROUP
func ExampleFirewallPolicyRuleGroupsOperations_BeginDelete() {
	fwClient := getFirewallPolicyRuleGroupsOperations()
	fwResp, err := fwClient.BeginDelete(context.Background(), resourceGroupName, policyName, ruleGroupName, nil)
	if err != nil {
		panic(err)
	}
	res, err := fwResp.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.StatusCode)
	// Output:
	// 200
}

func createFirewallPolicy(rgName, loc, policyName string) *armnetwork.FirewallPolicy {
	fpClient := getFirewallPoliciesOperations()
	fpResp, err := fpClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		policyName,
		armnetwork.FirewallPolicy{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(policyName),
				Location: to.StringPtr(loc),
			},
		},
		nil)
	if err != nil {
		panic(err)
	}
	res, err := fpResp.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	return res.FirewallPolicy
}
