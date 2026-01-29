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

type NetworkManagerTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	networkManagerName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *NetworkManagerTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.networkManagerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkman", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *NetworkManagerTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestNetworkManagerTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkManagerTestSuite))
}

// Microsoft.Network/networkManagers/{networkManagerName}
func (testsuite *NetworkManagerTestSuite) TestNetworkManagers() {
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

	// From step NetworkManagers_ListBySubscription
	fmt.Println("Call operation: NetworkManagers_ListBySubscription")
	managersClientNewListBySubscriptionPager := managersClient.NewListBySubscriptionPager(&armnetwork.ManagersClientListBySubscriptionOptions{Top: nil,
		SkipToken: nil,
	})
	for managersClientNewListBySubscriptionPager.More() {
		_, err := managersClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkManagers_List
	fmt.Println("Call operation: NetworkManagers_List")
	managersClientNewListPager := managersClient.NewListPager(testsuite.resourceGroupName, &armnetwork.ManagersClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for managersClientNewListPager.More() {
		_, err := managersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkManagers_Get
	fmt.Println("Call operation: NetworkManagers_Get")
	_, err = managersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, nil)
	testsuite.Require().NoError(err)

	// From step NetworkManagers_Patch
	fmt.Println("Call operation: NetworkManagers_Patch")
	_, err = managersClient.Patch(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, armnetwork.PatchObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NetworkManagerCommits_Post
	fmt.Println("Call operation: NetworkManagerCommits_Post")
	managerCommitsClient, err := armnetwork.NewManagerCommitsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	managerCommitsClientPostResponsePoller, err := managerCommitsClient.BeginPost(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, armnetwork.ManagerCommit{
		CommitType: to.Ptr(armnetwork.ConfigurationTypeConnectivity),
		TargetLocations: []*string{
			to.Ptr("eastus")},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managerCommitsClientPostResponsePoller)
	testsuite.Require().NoError(err)

	// From step NetworkManagerDeploymentStatus_List
	fmt.Println("Call operation: NetworkManagerDeploymentStatus_List")
	managerDeploymentStatusClient, err := armnetwork.NewManagerDeploymentStatusClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managerDeploymentStatusClient.List(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, armnetwork.ManagerDeploymentStatusParameter{
		DeploymentTypes: []*armnetwork.ConfigurationType{
			to.Ptr(armnetwork.ConfigurationTypeConnectivity),
			to.Ptr(armnetwork.ConfigurationTypeSecurityAdmin)},
		Regions: []*string{
			to.Ptr("eastus"),
			to.Ptr("westus")},
	}, &armnetwork.ManagerDeploymentStatusClientListOptions{Top: nil})
	testsuite.Require().NoError(err)

	// From step ListActiveConnectivityConfigurations
	fmt.Println("Call operation: NetworkManagementClient_ListActiveConnectivityConfigurations")
	managementClient, err := armnetwork.NewManagementClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = managementClient.ListActiveConnectivityConfigurations(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, armnetwork.ActiveConfigurationParameter{
		Regions: []*string{
			to.Ptr("westus")},
	}, &armnetwork.ManagementClientListActiveConnectivityConfigurationsOptions{Top: nil})
	testsuite.Require().NoError(err)

	// From step ListActiveSecurityAdminRules
	fmt.Println("Call operation: NetworkManagementClient_ListActiveSecurityAdminRules")
	_, err = managementClient.ListActiveSecurityAdminRules(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, armnetwork.ActiveConfigurationParameter{
		Regions: []*string{
			to.Ptr("westus")},
	}, &armnetwork.ManagementClientListActiveSecurityAdminRulesOptions{Top: nil})
	testsuite.Require().NoError(err)

	// From step NetworkManagers_Delete
	fmt.Println("Call operation: NetworkManagers_Delete")
	managersClientDeleteResponsePoller, err := managersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, &armnetwork.ManagersClientBeginDeleteOptions{Force: to.Ptr(true)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, managersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
