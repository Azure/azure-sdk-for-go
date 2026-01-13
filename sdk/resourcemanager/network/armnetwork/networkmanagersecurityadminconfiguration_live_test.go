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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v8"
	"github.com/stretchr/testify/suite"
)

type NetworkManagerSecurityAdminConfigurationTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	configurationName  string
	networkGroupId     string
	networkGroupName   string
	networkManagerName string
	ruleCollectionName string
	ruleName           string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *NetworkManagerSecurityAdminConfigurationTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.configurationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "configurationsecurity", 27, false)
	testsuite.networkGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkgrosecurity", 24, false)
	testsuite.networkManagerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkmanagersecurity", 28, false)
	testsuite.ruleCollectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rulecollec", 16, false)
	testsuite.ruleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "rulename", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NetworkManagerSecurityAdminConfigurationTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestNetworkManagerSecurityAdminConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkManagerSecurityAdminConfigurationTestSuite))
}

func (testsuite *NetworkManagerSecurityAdminConfigurationTestSuite) Prepare() {
	var err error
	// From step NetworkManagers_CreateOrUpdate
	fmt.Println("Call operation: NetworkManagers_CreateOrUpdate")
	managersClient, err := armnetwork.NewManagersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, armnetwork.Manager{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.ManagerProperties{
			Description: to.Ptr("My Test Network Manager"),
			NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
				to.Ptr(armnetwork.ConfigurationTypeSecurityAdmin)},
			NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
				Subscriptions: []*string{
					to.Ptr("/subscriptions/" + testsuite.subscriptionId)},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NetworkGroups_CreateOrUpdate
	fmt.Println("Call operation: NetworkGroups_CreateOrUpdate")
	groupsClient, err := armnetwork.NewGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	groupsClientCreateOrUpdateResponse, err := groupsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, armnetwork.Group{
		Properties: &armnetwork.GroupProperties{
			Description: to.Ptr("A sample group"),
		},
	}, &armnetwork.GroupsClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	testsuite.networkGroupId = *groupsClientCreateOrUpdateResponse.ID

	// From step SecurityAdminConfigurations_CreateOrUpdate
	fmt.Println("Call operation: SecurityAdminConfigurations_CreateOrUpdate")
	securityAdminConfigurationsClient, err := armnetwork.NewSecurityAdminConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = securityAdminConfigurationsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, armnetwork.SecurityAdminConfiguration{
		Properties: &armnetwork.SecurityAdminConfigurationPropertiesFormat{
			Description: to.Ptr("A sample policy"),
			ApplyOnNetworkIntentPolicyBasedServices: []*armnetwork.NetworkIntentPolicyBasedService{
				to.Ptr(armnetwork.NetworkIntentPolicyBasedServiceNone)},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AdminRuleCollections_CreateOrUpdate
	fmt.Println("Call operation: AdminRuleCollections_CreateOrUpdate")
	adminRuleCollectionsClient, err := armnetwork.NewAdminRuleCollectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = adminRuleCollectionsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, testsuite.ruleCollectionName, armnetwork.AdminRuleCollection{
		Properties: &armnetwork.AdminRuleCollectionPropertiesFormat{
			Description: to.Ptr("A sample policy"),
			AppliesToGroups: []*armnetwork.ManagerSecurityGroupItem{
				{
					NetworkGroupID: to.Ptr(testsuite.networkGroupId),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/networkManagers/{networkManagerName}/securityAdminConfigurations/{configurationName}
func (testsuite *NetworkManagerSecurityAdminConfigurationTestSuite) TestSecurityAdminConfigurations() {
	var err error
	// From step SecurityAdminConfigurations_List
	fmt.Println("Call operation: SecurityAdminConfigurations_List")
	securityAdminConfigurationsClient, err := armnetwork.NewSecurityAdminConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	securityAdminConfigurationsClientNewListPager := securityAdminConfigurationsClient.NewListPager(testsuite.resourceGroupName, testsuite.networkManagerName, &armnetwork.SecurityAdminConfigurationsClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for securityAdminConfigurationsClientNewListPager.More() {
		_, err := securityAdminConfigurationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SecurityAdminConfigurations_Get
	fmt.Println("Call operation: SecurityAdminConfigurations_Get")
	_, err = securityAdminConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/networkManagers/{networkManagerName}/securityAdminConfigurations/{configurationName}/ruleCollections/{ruleCollectionName}
func (testsuite *NetworkManagerSecurityAdminConfigurationTestSuite) TestAdminRuleCollections() {
	var err error
	// From step AdminRuleCollections_List
	fmt.Println("Call operation: AdminRuleCollections_List")
	adminRuleCollectionsClient, err := armnetwork.NewAdminRuleCollectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	adminRuleCollectionsClientNewListPager := adminRuleCollectionsClient.NewListPager(testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, &armnetwork.AdminRuleCollectionsClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for adminRuleCollectionsClientNewListPager.More() {
		_, err := adminRuleCollectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AdminRuleCollections_Get
	fmt.Println("Call operation: AdminRuleCollections_Get")
	_, err = adminRuleCollectionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, testsuite.ruleCollectionName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/networkManagers/{networkManagerName}/securityAdminConfigurations/{configurationName}/ruleCollections/{ruleCollectionName}/rules/{ruleName}
func (testsuite *NetworkManagerSecurityAdminConfigurationTestSuite) TestAdminRules() {
	var err error
	// From step AdminRules_CreateOrUpdate
	fmt.Println("Call operation: AdminRules_CreateOrUpdate")
	adminRulesClient, err := armnetwork.NewAdminRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = adminRulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, testsuite.ruleCollectionName, testsuite.ruleName, &armnetwork.AdminRule{
		Kind: to.Ptr(armnetwork.AdminRuleKindCustom),
		Properties: &armnetwork.AdminPropertiesFormat{
			Description: to.Ptr("This is Sample Admin Rule"),
			Access:      to.Ptr(armnetwork.SecurityConfigurationRuleAccessDeny),
			DestinationPortRanges: []*string{
				to.Ptr("22")},
			Destinations: []*armnetwork.AddressPrefixItem{
				{
					AddressPrefix:     to.Ptr("*"),
					AddressPrefixType: to.Ptr(armnetwork.AddressPrefixTypeIPPrefix),
				}},
			Direction: to.Ptr(armnetwork.SecurityConfigurationRuleDirectionInbound),
			Priority:  to.Ptr[int32](1),
			SourcePortRanges: []*string{
				to.Ptr("0-65535")},
			Sources: []*armnetwork.AddressPrefixItem{
				{
					AddressPrefix:     to.Ptr("Internet"),
					AddressPrefixType: to.Ptr(armnetwork.AddressPrefixTypeServiceTag),
				}},
			Protocol: to.Ptr(armnetwork.SecurityConfigurationRuleProtocolTCP),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AdminRules_List
	fmt.Println("Call operation: AdminRules_List")
	adminRulesClientNewListPager := adminRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, testsuite.ruleCollectionName, &armnetwork.AdminRulesClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for adminRulesClientNewListPager.More() {
		_, err := adminRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AdminRules_Get
	fmt.Println("Call operation: AdminRules_Get")
	_, err = adminRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, testsuite.ruleCollectionName, testsuite.ruleName, nil)
	testsuite.Require().NoError(err)

	// From step AdminRules_Delete
	fmt.Println("Call operation: AdminRules_Delete")
	adminRulesClientDeleteResponsePoller, err := adminRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, testsuite.ruleCollectionName, testsuite.ruleName, &armnetwork.AdminRulesClientBeginDeleteOptions{Force: to.Ptr(false)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, adminRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *NetworkManagerSecurityAdminConfigurationTestSuite) Cleanup() {
	var err error
	// From step AdminRuleCollections_Delete
	fmt.Println("Call operation: AdminRuleCollections_Delete")
	adminRuleCollectionsClient, err := armnetwork.NewAdminRuleCollectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	adminRuleCollectionsClientDeleteResponsePoller, err := adminRuleCollectionsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, testsuite.ruleCollectionName, &armnetwork.AdminRuleCollectionsClientBeginDeleteOptions{Force: to.Ptr(false)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, adminRuleCollectionsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step SecurityAdminConfigurations_Delete
	fmt.Println("Call operation: SecurityAdminConfigurations_Delete")
	securityAdminConfigurationsClient, err := armnetwork.NewSecurityAdminConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	securityAdminConfigurationsClientDeleteResponsePoller, err := securityAdminConfigurationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, &armnetwork.SecurityAdminConfigurationsClientBeginDeleteOptions{Force: to.Ptr(false)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, securityAdminConfigurationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
