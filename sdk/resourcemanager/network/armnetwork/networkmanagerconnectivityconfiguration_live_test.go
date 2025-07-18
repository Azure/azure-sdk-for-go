//go:build go1.18
// +build go1.18

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

type NetworkManagerConnectivityConfigurationTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	configurationName  string
	networkGroupId     string
	networkGroupName   string
	networkManagerName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *NetworkManagerConnectivityConfigurationTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.configurationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "connectivityconf", 22, false)
	testsuite.networkGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkgroconeconfig", 26, false)
	testsuite.networkManagerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkmancc", 18, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NetworkManagerConnectivityConfigurationTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestNetworkManagerConnectivityConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkManagerConnectivityConfigurationTestSuite))
}

func (testsuite *NetworkManagerConnectivityConfigurationTestSuite) Prepare() {
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
				to.Ptr(armnetwork.ConfigurationTypeConnectivity)},
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
}

// Microsoft.Network/networkManagers/{networkManagerName}/connectivityConfigurations/{configurationName}
func (testsuite *NetworkManagerConnectivityConfigurationTestSuite) TestConnectivityConfigurations() {
	var err error
	// From step ConnectivityConfigurations_CreateOrUpdate
	fmt.Println("Call operation: ConnectivityConfigurations_CreateOrUpdate")
	connectivityConfigurationsClient, err := armnetwork.NewConnectivityConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = connectivityConfigurationsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, armnetwork.ConnectivityConfiguration{
		Properties: &armnetwork.ConnectivityConfigurationProperties{
			Description: to.Ptr("Sample Configuration"),
			AppliesToGroups: []*armnetwork.ConnectivityGroupItem{
				{
					GroupConnectivity: to.Ptr(armnetwork.GroupConnectivityNone),
					IsGlobal:          to.Ptr(armnetwork.IsGlobalFalse),
					NetworkGroupID:    to.Ptr(testsuite.networkGroupId),
					UseHubGateway:     to.Ptr(armnetwork.UseHubGatewayTrue),
				}},
			ConnectivityTopology:  to.Ptr(armnetwork.ConnectivityTopologyMesh),
			DeleteExistingPeering: to.Ptr(armnetwork.DeleteExistingPeeringTrue),
			IsGlobal:              to.Ptr(armnetwork.IsGlobalTrue),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ConnectivityConfigurations_List
	fmt.Println("Call operation: ConnectivityConfigurations_List")
	connectivityConfigurationsClientNewListPager := connectivityConfigurationsClient.NewListPager(testsuite.resourceGroupName, testsuite.networkManagerName, &armnetwork.ConnectivityConfigurationsClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for connectivityConfigurationsClientNewListPager.More() {
		_, err := connectivityConfigurationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ConnectivityConfigurations_Get
	fmt.Println("Call operation: ConnectivityConfigurations_Get")
	_, err = connectivityConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, nil)
	testsuite.Require().NoError(err)

	// From step ConnectivityConfigurations_Delete
	fmt.Println("Call operation: ConnectivityConfigurations_Delete")
	connectivityConfigurationsClientDeleteResponsePoller, err := connectivityConfigurationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.configurationName, &armnetwork.ConnectivityConfigurationsClientBeginDeleteOptions{Force: to.Ptr(false)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, connectivityConfigurationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
