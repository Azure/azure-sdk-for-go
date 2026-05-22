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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v9"
	"github.com/stretchr/testify/suite"
)

type FirewallPolicyTestSuite struct {
	suite.Suite

	ctx                     context.Context
	cred                    azcore.TokenCredential
	options                 *arm.ClientOptions
	firewallPolicyName      string
	ruleCollectionGroupName string
	location                string
	resourceGroupName       string
	subscriptionId          string
}

func (testsuite *FirewallPolicyTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.firewallPolicyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "firewallpo", 16, false)
	testsuite.ruleCollectionGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rulecollec", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *FirewallPolicyTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestFirewallPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(FirewallPolicyTestSuite))
}

func (testsuite *FirewallPolicyTestSuite) Prepare() {
	var err error
	// From step FirewallPolicies_CreateOrUpdate
	fmt.Println("Call operation: FirewallPolicies_CreateOrUpdate")
	firewallPoliciesClient, err := armnetwork.NewFirewallPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	firewallPoliciesClientCreateOrUpdateResponsePoller, err := firewallPoliciesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.firewallPolicyName, armnetwork.FirewallPolicy{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.FirewallPolicyPropertiesFormat{
			SKU: &armnetwork.FirewallPolicySKU{
				Tier: to.Ptr(armnetwork.FirewallPolicySKUTierPremium),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallPoliciesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/firewallPolicies/{firewallPolicyName}
func (testsuite *FirewallPolicyTestSuite) TestFirewallPolicies() {
	var err error
	// From step FirewallPolicies_ListAll
	fmt.Println("Call operation: FirewallPolicies_ListAll")
	firewallPoliciesClient, err := armnetwork.NewFirewallPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	firewallPoliciesClientNewListAllPager := firewallPoliciesClient.NewListAllPager(nil)
	for firewallPoliciesClientNewListAllPager.More() {
		_, err := firewallPoliciesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallPolicies_List
	fmt.Println("Call operation: FirewallPolicies_List")
	firewallPoliciesClientNewListPager := firewallPoliciesClient.NewListPager(testsuite.resourceGroupName, nil)
	for firewallPoliciesClientNewListPager.More() {
		_, err := firewallPoliciesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallPolicies_Get
	fmt.Println("Call operation: FirewallPolicies_Get")
	_, err = firewallPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.firewallPolicyName, &armnetwork.FirewallPoliciesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step FirewallPolicies_UpdateTags
	fmt.Println("Call operation: FirewallPolicies_UpdateTags")
	_, err = firewallPoliciesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.firewallPolicyName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/firewallPolicies/{firewallPolicyName}/ruleCollectionGroups/{ruleCollectionGroupName}
func (testsuite *FirewallPolicyTestSuite) TestFirewallPolicyRuleCollectionGroups() {
	var err error
	// From step FirewallPolicyRuleCollectionGroups_CreateOrUpdate
	fmt.Println("Call operation: FirewallPolicyRuleCollectionGroups_CreateOrUpdate")
	firewallPolicyRuleCollectionGroupsClient, err := armnetwork.NewFirewallPolicyRuleCollectionGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	firewallPolicyRuleCollectionGroupsClientCreateOrUpdateResponsePoller, err := firewallPolicyRuleCollectionGroupsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.firewallPolicyName, testsuite.ruleCollectionGroupName, armnetwork.FirewallPolicyRuleCollectionGroup{
		Properties: &armnetwork.FirewallPolicyRuleCollectionGroupProperties{
			Priority: to.Ptr[int32](100),
			RuleCollections: []armnetwork.FirewallPolicyRuleCollectionClassification{
				&armnetwork.FirewallPolicyFilterRuleCollection{
					Name:               to.Ptr("Example-Filter-Rule-Collection"),
					Priority:           to.Ptr[int32](100),
					RuleCollectionType: to.Ptr(armnetwork.FirewallPolicyRuleCollectionTypeFirewallPolicyFilterRuleCollection),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallPolicyRuleCollectionGroupsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step FirewallPolicyRuleCollectionGroups_List
	fmt.Println("Call operation: FirewallPolicyRuleCollectionGroups_List")
	firewallPolicyRuleCollectionGroupsClientNewListPager := firewallPolicyRuleCollectionGroupsClient.NewListPager(testsuite.resourceGroupName, testsuite.firewallPolicyName, nil)
	for firewallPolicyRuleCollectionGroupsClientNewListPager.More() {
		_, err := firewallPolicyRuleCollectionGroupsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FirewallPolicyRuleCollectionGroups_Get
	fmt.Println("Call operation: FirewallPolicyRuleCollectionGroups_Get")
	_, err = firewallPolicyRuleCollectionGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.firewallPolicyName, testsuite.ruleCollectionGroupName, nil)
	testsuite.Require().NoError(err)

	// From step FirewallPolicyRuleCollectionGroups_Delete
	fmt.Println("Call operation: FirewallPolicyRuleCollectionGroups_Delete")
	firewallPolicyRuleCollectionGroupsClientDeleteResponsePoller, err := firewallPolicyRuleCollectionGroupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.firewallPolicyName, testsuite.ruleCollectionGroupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallPolicyRuleCollectionGroupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *FirewallPolicyTestSuite) Cleanup() {
	var err error
	// From step FirewallPolicies_Delete
	fmt.Println("Call operation: FirewallPolicies_Delete")
	firewallPoliciesClient, err := armnetwork.NewFirewallPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	firewallPoliciesClientDeleteResponsePoller, err := firewallPoliciesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.firewallPolicyName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, firewallPoliciesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
