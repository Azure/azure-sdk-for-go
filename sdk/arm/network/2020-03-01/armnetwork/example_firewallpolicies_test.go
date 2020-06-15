// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armnetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2019-05-01/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

const (
	policyName    = "samplepolicy"
	ruleGroupName = "sampleRuleGroup"
)

func getResourceGroupsOperations() armresources.ResourceGroupsOperations {
	client, err := armresources.NewDefaultClient(getCredential(), nil)
	if err != nil {
		panic(err)
	}
	return client.ResourceGroupsOperations(subscriptionID)
}

func getFirewallPolicyRuleGroupsOperations() FirewallPolicyRuleGroupsOperations {
	client, err := NewDefaultClient(getCredential(), nil)
	if err != nil {
		panic(err)
	}
	return client.FirewallPolicyRuleGroupsOperations(subscriptionID)
}

func getFirewallPoliciesOperations() FirewallPoliciesOperations {
	client, err := NewDefaultClient(getCredential(), nil)
	if err != nil {
		panic(err)
	}
	return client.FirewallPoliciesOperations(subscriptionID)
}

func ExampleFirewallPolicyRuleGroupsOperations_BeginCreateOrUpdate() {
	// create a new resource group to create resource in
	rg := createResourceGroup(resourceGroupName, location)
	rgName := *rg.Name
	// get FirewallPoliciesOperations and create a new FirewallPolicy to use in the FirewallPolicyRuleGroup
	fwPolicy := createFirewallPolicy(rgName, location, policyName)
	fwPolicyName := *fwPolicy.Name
	// get FirewallPolicyRuleGroupsOperations and create a new FirewallPolicyRuleGroup using the FirewallPolicy that was previously created
	fwClient := getFirewallPolicyRuleGroupsOperations()
	fwResp, err := fwClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		fwPolicyName,
		ruleGroupName,
		FirewallPolicyRuleGroup{
			Name: to.StringPtr(ruleGroupName),
			Properties: &FirewallPolicyRuleGroupProperties{
				Priority: to.Int32Ptr(110),
				Rules: &[]FirewallPolicyRuleClassification{
					&FirewallPolicyFilterRule{
						FirewallPolicyRule: FirewallPolicyRule{
							Priority: to.Int32Ptr(110),
							Name:     to.StringPtr("rule1"),
							RuleType: FirewallPolicyRuleTypeFirewallPolicyFilterRule.ToPtr(),
						},
						Action: &FirewallPolicyFilterRuleAction{
							Type: FirewallPolicyFilterRuleActionTypeDeny.ToPtr(),
						},
					},
				},
			},
		})
	if err != nil {
		panic(err)
	}
	res, err := fwResp.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(*res)
}

func ExampleFirewallPolicyRuleGroupsOperations_BeginDelete() {
	fwClient := getFirewallPolicyRuleGroupsOperations()
	fwResp, err := fwClient.BeginDelete(context.Background(), resourceGroupName, policyName, ruleGroupName)
	if err != nil {
		panic(err)
	}
	res, err := fwResp.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.StatusCode)
}

func ExampleResourceGroupsOperations_BeginDelete() {
	rgClient := getResourceGroupsOperations()
	rgResp, err := rgClient.BeginDelete(context.Background(), resourceGroupName)
	if err != nil {
		panic(err)
	}
	// the following demonstrates the recommended way to manually handle polling
	poller := rgResp.Poller
	for {
		resp, err := poller.Poll(context.Background())
		if err != nil {
			panic(err)
		}
		if poller.Done() {
			break
		}
		if delay := azcore.RetryAfter(resp); delay > 0 {
			time.Sleep(delay)
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	res := poller.FinalResponse()
	fmt.Println(res.StatusCode)
}

func createResourceGroup(rgName, loc string) *armresources.ResourceGroup {
	rgClient := getResourceGroupsOperations()
	rgResp, err := rgClient.CreateOrUpdate(
		context.Background(),
		rgName,
		armresources.ResourceGroup{
			Name:     to.StringPtr(rgName),
			Location: to.StringPtr(loc),
		})
	if err != nil {
		panic(err)
	}
	return rgResp.ResourceGroup
}

func createFirewallPolicy(rgName, loc, policyName string) *FirewallPolicy {
	fpClient := getFirewallPoliciesOperations()
	fpResp, err := fpClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		policyName,
		FirewallPolicy{
			Resource: Resource{
				Name:     to.StringPtr(policyName),
				Location: to.StringPtr(loc),
			},
		})
	if err != nil {
		panic(err)
	}
	return fpResp.FirewallPolicy
}
