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

type NetworkManagerGroupTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	networkGroupName   string
	networkManagerName string
	staticMemberName   string
	virtualNetworkName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *NetworkManagerGroupTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.networkGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkgro", 16, false)
	testsuite.networkManagerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkmanagergp", 22, false)
	testsuite.staticMemberName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "staticmemberna", 20, false)
	testsuite.virtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkgrovet", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NetworkManagerGroupTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestNetworkManagerGroupTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkManagerGroupTestSuite))
}

func (testsuite *NetworkManagerGroupTestSuite) Prepare() {
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
	_, err = groupsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, armnetwork.Group{
		Properties: &armnetwork.GroupProperties{
			Description: to.Ptr("A sample group"),
		},
	}, &armnetwork.GroupsClientCreateOrUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.Network/networkManagers/{networkManagerName}/networkGroups/{networkGroupName}
func (testsuite *NetworkManagerGroupTestSuite) TestNetworkGroups() {
	var err error
	// From step NetworkGroups_List
	fmt.Println("Call operation: NetworkGroups_List")
	groupsClient, err := armnetwork.NewGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	groupsClientNewListPager := groupsClient.NewListPager(testsuite.resourceGroupName, testsuite.networkManagerName, &armnetwork.GroupsClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for groupsClientNewListPager.More() {
		_, err := groupsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkGroups_Get
	fmt.Println("Call operation: NetworkGroups_Get")
	_, err = groupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/networkManagers/{networkManagerName}/networkGroups/{networkGroupName}/staticMembers/{staticMemberName}
func (testsuite *NetworkManagerGroupTestSuite) TestStaticMembers() {
	var virutalNetworkId string
	var err error
	// From step VirtualNetworks_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworks_CreateOrUpdate")
	virtualNetworksClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworksClientCreateOrUpdateResponsePoller, err := virtualNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, armnetwork.VirtualNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.0.0.0/16")},
			},
			FlowTimeoutInMinutes: to.Ptr[int32](10),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var virtualNetworksClientCreateOrUpdateResponse *armnetwork.VirtualNetworksClientCreateOrUpdateResponse
	virtualNetworksClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, virtualNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	virutalNetworkId = *virtualNetworksClientCreateOrUpdateResponse.ID

	// From step StaticMembers_CreateOrUpdate
	fmt.Println("Call operation: StaticMembers_CreateOrUpdate")
	staticMembersClient, err := armnetwork.NewStaticMembersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = staticMembersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, testsuite.staticMemberName, armnetwork.StaticMember{
		Properties: &armnetwork.StaticMemberProperties{
			ResourceID: to.Ptr(virutalNetworkId),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step StaticMembers_List
	fmt.Println("Call operation: StaticMembers_List")
	staticMembersClientNewListPager := staticMembersClient.NewListPager(testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, &armnetwork.StaticMembersClientListOptions{Top: nil,
		SkipToken: nil,
	})
	for staticMembersClientNewListPager.More() {
		_, err := staticMembersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StaticMembers_Get
	fmt.Println("Call operation: StaticMembers_Get")
	_, err = staticMembersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, testsuite.staticMemberName, nil)
	testsuite.Require().NoError(err)

	// From step StaticMembers_Delete
	fmt.Println("Call operation: StaticMembers_Delete")
	_, err = staticMembersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, testsuite.staticMemberName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *NetworkManagerGroupTestSuite) Cleanup() {
	var err error
	// From step NetworkGroups_Delete
	fmt.Println("Call operation: NetworkGroups_Delete")
	groupsClient, err := armnetwork.NewGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	groupsClientDeleteResponsePoller, err := groupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkManagerName, testsuite.networkGroupName, &armnetwork.GroupsClientBeginDeleteOptions{Force: to.Ptr(false)})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, groupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
