// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"github.com/stretchr/testify/suite"
)

type WebapplicationfirewallTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	policyName        string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *WebapplicationfirewallTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.policyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "wafpolicyname", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *WebapplicationfirewallTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestWebapplicationfirewallTestSuite(t *testing.T) {
	suite.Run(t, new(WebapplicationfirewallTestSuite))
}

// Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies/{policyName}
func (testsuite *WebapplicationfirewallTestSuite) TestWebApplicationFirewallPolicies() {
	var err error
	// From step WebApplicationFirewallPolicies_CreateOrUpdate
	fmt.Println("Call operation: WebApplicationFirewallPolicies_CreateOrUpdate")
	webApplicationFirewallPoliciesClient, err := armnetwork.NewWebApplicationFirewallPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = webApplicationFirewallPoliciesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.policyName, armnetwork.WebApplicationFirewallPolicy{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.WebApplicationFirewallPolicyPropertiesFormat{
			ManagedRules: &armnetwork.ManagedRulesDefinition{
				ManagedRuleSets: []*armnetwork.ManagedRuleSet{
					{
						RuleGroupOverrides: []*armnetwork.ManagedRuleGroupOverride{
							{
								RuleGroupName: to.Ptr("REQUEST-931-APPLICATION-ATTACK-RFI"),
								Rules: []*armnetwork.ManagedRuleOverride{
									{
										Action: to.Ptr(armnetwork.ActionTypeLog),
										RuleID: to.Ptr("931120"),
										State:  to.Ptr(armnetwork.ManagedRuleEnabledStateEnabled),
									}},
							}},
						RuleSetType:    to.Ptr("OWASP"),
						RuleSetVersion: to.Ptr("3.2"),
					}},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step WebApplicationFirewallPolicies_List
	fmt.Println("Call operation: WebApplicationFirewallPolicies_List")
	webApplicationFirewallPoliciesClientNewListPager := webApplicationFirewallPoliciesClient.NewListPager(testsuite.resourceGroupName, nil)
	for webApplicationFirewallPoliciesClientNewListPager.More() {
		_, err := webApplicationFirewallPoliciesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WebApplicationFirewallPolicies_ListAll
	fmt.Println("Call operation: WebApplicationFirewallPolicies_ListAll")
	webApplicationFirewallPoliciesClientNewListAllPager := webApplicationFirewallPoliciesClient.NewListAllPager(nil)
	for webApplicationFirewallPoliciesClientNewListAllPager.More() {
		_, err := webApplicationFirewallPoliciesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step WebApplicationFirewallPolicies_Get
	fmt.Println("Call operation: WebApplicationFirewallPolicies_Get")
	_, err = webApplicationFirewallPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.policyName, nil)
	testsuite.Require().NoError(err)

	// From step WebApplicationFirewallPolicies_Delete
	fmt.Println("Call operation: WebApplicationFirewallPolicies_Delete")
	webApplicationFirewallPoliciesClientDeleteResponsePoller, err := webApplicationFirewallPoliciesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.policyName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, webApplicationFirewallPoliciesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
