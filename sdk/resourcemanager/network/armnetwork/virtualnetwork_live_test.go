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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"github.com/stretchr/testify/suite"
)

type VirtualNetworkTestSuite struct {
	suite.Suite

	ctx                       context.Context
	cred                      azcore.TokenCredential
	options                   *arm.ClientOptions
	virtualNetworkName        string
	virtualNetworkPeeringName string
	subnetName                string
	location                  string
	resourceGroupName         string
	subscriptionId            string
}

func (testsuite *VirtualNetworkTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.virtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualnet", 16, false)
	testsuite.virtualNetworkPeeringName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualnetpee", 19, false)
	testsuite.subnetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "subnetname", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *VirtualNetworkTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVirtualNetworkTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualNetworkTestSuite))
}

func (testsuite *VirtualNetworkTestSuite) Prepare() {
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
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/virtualNetworks/{virtualNetworkName}
func (testsuite *VirtualNetworkTestSuite) TestVirtualNetworks() {
	var err error
	// From step VirtualNetworks_ListAll
	fmt.Println("Call operation: VirtualNetworks_ListAll")
	virtualNetworksClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworksClientNewListAllPager := virtualNetworksClient.NewListAllPager(nil)
	for virtualNetworksClientNewListAllPager.More() {
		_, err := virtualNetworksClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworks_List
	fmt.Println("Call operation: VirtualNetworks_List")
	virtualNetworksClientNewListPager := virtualNetworksClient.NewListPager(testsuite.resourceGroupName, nil)
	for virtualNetworksClientNewListPager.More() {
		_, err := virtualNetworksClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworks_Get
	fmt.Println("Call operation: VirtualNetworks_Get")
	_, err = virtualNetworksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, &armnetwork.VirtualNetworksClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step VirtualNetworks_ListUsage
	fmt.Println("Call operation: VirtualNetworks_ListUsage")
	virtualNetworksClientNewListUsagePager := virtualNetworksClient.NewListUsagePager(testsuite.resourceGroupName, testsuite.virtualNetworkName, nil)
	for virtualNetworksClientNewListUsagePager.More() {
		_, err := virtualNetworksClientNewListUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworks_UpdateTags
	fmt.Println("Call operation: VirtualNetworks_UpdateTags")
	_, err = virtualNetworksClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworks_ListDdosProtectionStatus
	fmt.Println("Call operation: VirtualNetworks_ListDdosProtectionStatus")
	virtualNetworksClientListDdosProtectionStatusResponsePoller, err := virtualNetworksClient.BeginListDdosProtectionStatus(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, &armnetwork.VirtualNetworksClientBeginListDdosProtectionStatusOptions{Top: to.Ptr[int32](75),
		SkipToken: nil,
	})
	testsuite.Require().NoError(err)
	virtualNetworksClientListDdosProtectionStatusResponse, err := testutil.PollForTest(testsuite.ctx, virtualNetworksClientListDdosProtectionStatusResponsePoller)
	testsuite.Require().NoError(err)
	for (*virtualNetworksClientListDdosProtectionStatusResponse).More() {
		_, err := virtualNetworksClientNewListUsagePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Network/virtualNetworks/{virtualNetworkName}/subnets/{subnetName}
func (testsuite *VirtualNetworkTestSuite) TestSubnets() {
	var err error
	// From step Subnets_CreateOrUpdate
	fmt.Println("Call operation: Subnets_CreateOrUpdate")
	subnetsClient, err := armnetwork.NewSubnetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	subnetsClientCreateOrUpdateResponsePoller, err := subnetsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, testsuite.subnetName, armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.Ptr("10.0.0.0/16"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, subnetsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Subnets_List
	fmt.Println("Call operation: Subnets_List")
	subnetsClientNewListPager := subnetsClient.NewListPager(testsuite.resourceGroupName, testsuite.virtualNetworkName, nil)
	for subnetsClientNewListPager.More() {
		_, err := subnetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Subnets_Get
	fmt.Println("Call operation: Subnets_Get")
	_, err = subnetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, testsuite.subnetName, &armnetwork.SubnetsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step Subnets_Delete
	fmt.Println("Call operation: Subnets_Delete")
	subnetsClientDeleteResponsePoller, err := subnetsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, testsuite.subnetName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, subnetsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/virtualNetworks/{virtualNetworkName}/virtualNetworkPeerings/{virtualNetworkPeeringName}
func (testsuite *VirtualNetworkTestSuite) TestVirtualNetworkPeerings() {
	var err error
	// From step VirtualNetworks_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworks_CreateOrUpdate")
	virtualNetworksClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworksClientCreateOrUpdateResponsePoller, err := virtualNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkPeeringName, armnetwork.VirtualNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("11.0.0.0/16")},
			},
			FlowTimeoutInMinutes: to.Ptr[int32](10),
		},
	}, nil)
	testsuite.Require().NoError(err)
	virtualNetworksClientCreateOrUpdateResponse, err := testutil.PollForTest(testsuite.ctx, virtualNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	virtualNetworkSecondId := *virtualNetworksClientCreateOrUpdateResponse.ID

	// From step VirtualNetworkPeerings_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworkPeerings_CreateOrUpdate")
	virtualNetworkPeeringsClient, err := armnetwork.NewVirtualNetworkPeeringsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworkPeeringsClientCreateOrUpdateResponsePoller, err := virtualNetworkPeeringsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, testsuite.virtualNetworkPeeringName, armnetwork.VirtualNetworkPeering{
		Properties: &armnetwork.VirtualNetworkPeeringPropertiesFormat{
			AllowForwardedTraffic:     to.Ptr(true),
			AllowGatewayTransit:       to.Ptr(false),
			AllowVirtualNetworkAccess: to.Ptr(true),
			RemoteVirtualNetwork: &armnetwork.SubResource{
				ID: to.Ptr(virtualNetworkSecondId),
			},
			UseRemoteGateways: to.Ptr(false),
		},
	}, &armnetwork.VirtualNetworkPeeringsClientBeginCreateOrUpdateOptions{SyncRemoteAddressSpace: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkPeeringsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkPeerings_List
	fmt.Println("Call operation: VirtualNetworkPeerings_List")
	virtualNetworkPeeringsClientNewListPager := virtualNetworkPeeringsClient.NewListPager(testsuite.resourceGroupName, testsuite.virtualNetworkName, nil)
	for virtualNetworkPeeringsClientNewListPager.More() {
		_, err := virtualNetworkPeeringsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworkPeerings_Get
	fmt.Println("Call operation: VirtualNetworkPeerings_Get")
	_, err = virtualNetworkPeeringsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, testsuite.virtualNetworkPeeringName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkPeerings_Delete
	fmt.Println("Call operation: VirtualNetworkPeerings_Delete")
	virtualNetworkPeeringsClientDeleteResponsePoller, err := virtualNetworkPeeringsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, testsuite.virtualNetworkPeeringName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkPeeringsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *VirtualNetworkTestSuite) Cleanup() {
	var err error
	// From step VirtualNetworks_Delete
	fmt.Println("Call operation: VirtualNetworks_Delete")
	virtualNetworksClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworksClientDeleteResponsePoller, err := virtualNetworksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
