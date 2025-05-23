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

type VirtualNetworkGatewayTestSuite struct {
	suite.Suite

	ctx                       context.Context
	cred                      azcore.TokenCredential
	options                   *arm.ClientOptions
	publicIpAddressId         string
	publicIpAddressName       string
	subnetId                  string
	virtualNetworkGatewayName string
	virtualNetworkName        string
	localNetworkGatewayName   string
	location                  string
	resourceGroupName         string
	subscriptionId            string
}

func (testsuite *VirtualNetworkGatewayTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.publicIpAddressName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "publicipadgateway", 23, false)
	testsuite.virtualNetworkGatewayName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualnetgateway", 23, false)
	testsuite.virtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualnetgateway", 23, false)
	testsuite.localNetworkGatewayName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "localnetwo", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *VirtualNetworkGatewayTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestVirtualNetworkGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualNetworkGatewayTestSuite))
}

func (testsuite *VirtualNetworkGatewayTestSuite) Prepare() {
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
			Subnets: []*armnetwork.Subnet{
				{
					Name: to.Ptr("GatewaySubnet"),
					Properties: &armnetwork.SubnetPropertiesFormat{
						AddressPrefix: to.Ptr("10.0.0.0/24"),
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var virtualNetworksClientCreateOrUpdateResponse *armnetwork.VirtualNetworksClientCreateOrUpdateResponse
	virtualNetworksClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, virtualNetworksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.subnetId = *virtualNetworksClientCreateOrUpdateResponse.Properties.Subnets[0].ID

	// From step PublicIPAddresses_CreateOrUpdate
	fmt.Println("Call operation: PublicIPAddresses_CreateOrUpdate")
	publicIPAddressesClient, err := armnetwork.NewPublicIPAddressesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	publicIPAddressesClientCreateOrUpdateResponsePoller, err := publicIPAddressesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.publicIpAddressName, armnetwork.PublicIPAddress{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	var publicIPAddressesClientCreateOrUpdateResponse *armnetwork.PublicIPAddressesClientCreateOrUpdateResponse
	publicIPAddressesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, publicIPAddressesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.publicIpAddressId = *publicIPAddressesClientCreateOrUpdateResponse.ID

	// From step VirtualNetworkGateways_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworkGateways_CreateOrUpdate")
	virtualNetworkGatewaysClient, err := armnetwork.NewVirtualNetworkGatewaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworkGatewaysClientCreateOrUpdateResponsePoller, err := virtualNetworkGatewaysClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkGatewayName, armnetwork.VirtualNetworkGateway{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
			GatewayType: to.Ptr(armnetwork.VirtualNetworkGatewayTypeVPN),
			IPConfigurations: []*armnetwork.VirtualNetworkGatewayIPConfiguration{
				{
					Name: to.Ptr("gwipconfig1"),
					Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
						PublicIPAddress: &armnetwork.SubResource{
							ID: to.Ptr(testsuite.publicIpAddressId),
						},
						Subnet: &armnetwork.SubResource{
							ID: to.Ptr(testsuite.subnetId),
						},
					},
				}},
			VPNType: to.Ptr(armnetwork.VPNTypeRouteBased),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkGatewaysClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/virtualNetworkGateways/{virtualNetworkGatewayName}
func (testsuite *VirtualNetworkGatewayTestSuite) TestVirtualNetworkGateways() {
	var err error
	// From step VirtualNetworkGateways_List
	fmt.Println("Call operation: VirtualNetworkGateways_List")
	virtualNetworkGatewaysClient, err := armnetwork.NewVirtualNetworkGatewaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworkGatewaysClientNewListPager := virtualNetworkGatewaysClient.NewListPager(testsuite.resourceGroupName, nil)
	for virtualNetworkGatewaysClientNewListPager.More() {
		_, err := virtualNetworkGatewaysClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworkGateways_Get
	fmt.Println("Call operation: VirtualNetworkGateways_Get")
	_, err = virtualNetworkGatewaysClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkGatewayName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkGateways_ListConnections
	fmt.Println("Call operation: VirtualNetworkGateways_ListConnections")
	virtualNetworkGatewaysClientNewListConnectionsPager := virtualNetworkGatewaysClient.NewListConnectionsPager(testsuite.resourceGroupName, testsuite.virtualNetworkGatewayName, nil)
	for virtualNetworkGatewaysClientNewListConnectionsPager.More() {
		_, err := virtualNetworkGatewaysClientNewListConnectionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworkGateways_UpdateTags
	fmt.Println("Call operation: VirtualNetworkGateways_UpdateTags")
	virtualNetworkGatewaysClientUpdateTagsResponsePoller, err := virtualNetworkGatewaysClient.BeginUpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkGatewayName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkGatewaysClientUpdateTagsResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkGateways_Delete
	fmt.Println("Call operation: VirtualNetworkGateways_Delete")
	virtualNetworkGatewaysClientDeleteResponsePoller, err := virtualNetworkGatewaysClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualNetworkGatewayName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkGatewaysClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/localNetworkGateways/{localNetworkGatewayName}
func (testsuite *VirtualNetworkGatewayTestSuite) TestLocalNetworkGateways() {
	var err error
	// From step LocalNetworkGateways_CreateOrUpdate
	fmt.Println("Call operation: LocalNetworkGateways_CreateOrUpdate")
	localNetworkGatewaysClient, err := armnetwork.NewLocalNetworkGatewaysClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	localNetworkGatewaysClientCreateOrUpdateResponsePoller, err := localNetworkGatewaysClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.localNetworkGatewayName, armnetwork.LocalNetworkGateway{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
			GatewayIPAddress: to.Ptr("11.12.13.14"),
			LocalNetworkAddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.1.0.0/16")},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, localNetworkGatewaysClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LocalNetworkGateways_List
	fmt.Println("Call operation: LocalNetworkGateways_List")
	localNetworkGatewaysClientNewListPager := localNetworkGatewaysClient.NewListPager(testsuite.resourceGroupName, nil)
	for localNetworkGatewaysClientNewListPager.More() {
		_, err := localNetworkGatewaysClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LocalNetworkGateways_Get
	fmt.Println("Call operation: LocalNetworkGateways_Get")
	_, err = localNetworkGatewaysClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.localNetworkGatewayName, nil)
	testsuite.Require().NoError(err)

	// From step LocalNetworkGateways_UpdateTags
	fmt.Println("Call operation: LocalNetworkGateways_UpdateTags")
	_, err = localNetworkGatewaysClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.localNetworkGatewayName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step LocalNetworkGateways_Delete
	fmt.Println("Call operation: LocalNetworkGateways_Delete")
	localNetworkGatewaysClientDeleteResponsePoller, err := localNetworkGatewaysClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.localNetworkGatewayName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, localNetworkGatewaysClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
