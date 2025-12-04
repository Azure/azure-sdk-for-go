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

type NetworkSecurityGroupTestSuite struct {
	suite.Suite

	ctx                      context.Context
	cred                     azcore.TokenCredential
	options                  *arm.ClientOptions
	networkSecurityGroupName string
	securityRuleName         string
	location                 string
	resourceGroupName        string
	subscriptionId           string
}

func (testsuite *NetworkSecurityGroupTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.networkSecurityGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networksec", 16, false)
	testsuite.securityRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "securityru", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NetworkSecurityGroupTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestNetworkSecurityGroupTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkSecurityGroupTestSuite))
}

func (testsuite *NetworkSecurityGroupTestSuite) Prepare() {
	var err error
	// From step NetworkSecurityGroups_CreateOrUpdate
	fmt.Println("Call operation: NetworkSecurityGroups_CreateOrUpdate")
	securityGroupsClient, err := armnetwork.NewSecurityGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	securityGroupsClientCreateOrUpdateResponsePoller, err := securityGroupsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, armnetwork.SecurityGroup{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, securityGroupsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/networkSecurityGroups/{networkSecurityGroupName}
func (testsuite *NetworkSecurityGroupTestSuite) TestNetworkSecurityGroups() {
	var err error
	// From step NetworkSecurityGroups_ListAll
	fmt.Println("Call operation: NetworkSecurityGroups_ListAll")
	securityGroupsClient, err := armnetwork.NewSecurityGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	securityGroupsClientNewListAllPager := securityGroupsClient.NewListAllPager(nil)
	for securityGroupsClientNewListAllPager.More() {
		_, err := securityGroupsClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkSecurityGroups_List
	fmt.Println("Call operation: NetworkSecurityGroups_List")
	securityGroupsClientNewListPager := securityGroupsClient.NewListPager(testsuite.resourceGroupName, nil)
	for securityGroupsClientNewListPager.More() {
		_, err := securityGroupsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkSecurityGroups_Get
	fmt.Println("Call operation: NetworkSecurityGroups_Get")
	_, err = securityGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, &armnetwork.SecurityGroupsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step NetworkSecurityGroups_UpdateTags
	fmt.Println("Call operation: NetworkSecurityGroups_UpdateTags")
	_, err = securityGroupsClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/networkSecurityGroups/{networkSecurityGroupName}/securityRules/{securityRuleName}
func (testsuite *NetworkSecurityGroupTestSuite) TestSecurityRules() {
	var defaultSecurityRuleName string
	var err error
	// From step SecurityRules_CreateOrUpdate
	fmt.Println("Call operation: SecurityRules_CreateOrUpdate")
	securityRulesClient, err := armnetwork.NewSecurityRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	securityRulesClientCreateOrUpdateResponsePoller, err := securityRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, testsuite.securityRuleName, armnetwork.SecurityRule{
		Properties: &armnetwork.SecurityRulePropertiesFormat{
			Access:                   to.Ptr(armnetwork.SecurityRuleAccessDeny),
			DestinationAddressPrefix: to.Ptr("11.0.0.0/8"),
			DestinationPortRange:     to.Ptr("8080"),
			Direction:                to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
			Priority:                 to.Ptr[int32](100),
			SourceAddressPrefix:      to.Ptr("10.0.0.0/8"),
			SourcePortRange:          to.Ptr("*"),
			Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, securityRulesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step SecurityRules_List
	fmt.Println("Call operation: SecurityRules_List")
	securityRulesClientNewListPager := securityRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.networkSecurityGroupName, nil)
	for securityRulesClientNewListPager.More() {
		_, err := securityRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SecurityRules_Get
	fmt.Println("Call operation: SecurityRules_Get")
	_, err = securityRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, testsuite.securityRuleName, nil)
	testsuite.Require().NoError(err)

	// From step DefaultSecurityRules_List
	fmt.Println("Call operation: DefaultSecurityRules_List")
	defaultSecurityRulesClient, err := armnetwork.NewDefaultSecurityRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	defaultSecurityRulesClientNewListPager := defaultSecurityRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.networkSecurityGroupName, nil)
	for defaultSecurityRulesClientNewListPager.More() {
		nextResult, err := defaultSecurityRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		defaultSecurityRuleName = *nextResult.Value[0].Name
		break
	}

	// From step DefaultSecurityRules_Get
	fmt.Println("Call operation: DefaultSecurityRules_Get")
	_, err = defaultSecurityRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, defaultSecurityRuleName, nil)
	testsuite.Require().NoError(err)

	// From step SecurityRules_Delete
	fmt.Println("Call operation: SecurityRules_Delete")
	securityRulesClientDeleteResponsePoller, err := securityRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, testsuite.securityRuleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, securityRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *NetworkSecurityGroupTestSuite) Cleanup() {
	var err error
	// From step NetworkSecurityGroups_Delete
	fmt.Println("Call operation: NetworkSecurityGroups_Delete")
	securityGroupsClient, err := armnetwork.NewSecurityGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	securityGroupsClientDeleteResponsePoller, err := securityGroupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkSecurityGroupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, securityGroupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
