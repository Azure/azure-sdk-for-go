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

type NetworkManagerConnectionTestSuite struct {
	suite.Suite

	ctx                          context.Context
	cred                         azcore.TokenCredential
	options                      *arm.ClientOptions
	networkManagerConnectionName string
	networkManagerId             string
	networkManagerName           string
	location                     string
	resourceGroupName            string
	subscriptionId               string
}

func (testsuite *NetworkManagerConnectionTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.networkManagerConnectionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "netmanagerconnsub", 23, false)
	testsuite.networkManagerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkmanagerconn", 24, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NetworkManagerConnectionTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestNetworkManagerConnectionTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkManagerConnectionTestSuite))
}

func (testsuite *NetworkManagerConnectionTestSuite) Prepare() {
	var err error
	// From step NetworkManagers_CreateOrUpdate
	fmt.Println("Call operation: NetworkManagers_CreateOrUpdate")
	managersClient, err := armnetwork.NewManagersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	managersClientCreateOrUpdateResponse, err := managersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, armnetwork.Manager{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.ManagerProperties{
			Description: to.Ptr("My Test Network Manager"),
			NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
				to.Ptr(armnetwork.ConfigurationTypeConnectivity)},
			NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
				Subscriptions: []*string{
					to.Ptr("/subscriptions/" + testsuite.subscriptionId),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	testsuite.networkManagerId = *managersClientCreateOrUpdateResponse.ID
}

// Microsoft.Network/networkManagerConnections/{networkManagerConnectionName}
func (testsuite *NetworkManagerConnectionTestSuite) TestSubscriptionNetworkManagerConnections() {
	var err error
	// From step SubscriptionNetworkManagerConnections_CreateOrUpdate
	fmt.Println("Call operation: SubscriptionNetworkManagerConnections_CreateOrUpdate")
	subscriptionNetworkManagerConnectionsClient, err := armnetwork.NewSubscriptionNetworkManagerConnectionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = subscriptionNetworkManagerConnectionsClient.CreateOrUpdate(testsuite.ctx, testsuite.networkManagerConnectionName, armnetwork.ManagerConnection{
		Properties: &armnetwork.ManagerConnectionProperties{
			NetworkManagerID: to.Ptr(testsuite.networkManagerId),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step SubscriptionNetworkManagerConnections_List
	fmt.Println("Call operation: SubscriptionNetworkManagerConnections_List")
	subscriptionNetworkManagerConnectionsClientNewListPager := subscriptionNetworkManagerConnectionsClient.NewListPager(&armnetwork.SubscriptionNetworkManagerConnectionsClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for subscriptionNetworkManagerConnectionsClientNewListPager.More() {
		_, err := subscriptionNetworkManagerConnectionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step SubscriptionNetworkManagerConnections_Get
	fmt.Println("Call operation: SubscriptionNetworkManagerConnections_Get")
	_, err = subscriptionNetworkManagerConnectionsClient.Get(testsuite.ctx, testsuite.networkManagerConnectionName, nil)
	testsuite.Require().NoError(err)

	// From step SubscriptionNetworkManagerConnections_Delete
	fmt.Println("Call operation: SubscriptionNetworkManagerConnections_Delete")
	_, err = subscriptionNetworkManagerConnectionsClient.Delete(testsuite.ctx, testsuite.networkManagerConnectionName, nil)
	testsuite.Require().NoError(err)
}
