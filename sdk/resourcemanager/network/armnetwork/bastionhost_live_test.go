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

type BastionHostTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	bastionHostName   string
	publicIpAddressId string
	subnetId          string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *BastionHostTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.bastionHostName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "bastionhos", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *BastionHostTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestBastionHostTestSuite(t *testing.T) {
	suite.Run(t, new(BastionHostTestSuite))
}

func (testsuite *BastionHostTestSuite) Prepare() {
	var err error
	// From step VirtualNetworks_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworks_CreateOrUpdate")
	virtualNetworksClient, err := armnetwork.NewVirtualNetworksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworksClientCreateOrUpdateResponsePoller, err := virtualNetworksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, "test-vnet", armnetwork.VirtualNetwork{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.0.0.0/16")},
			},
			Subnets: []*armnetwork.Subnet{
				{
					Name: to.Ptr("AzureBastionSubnet"),
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
	publicIPAddressesClientCreateOrUpdateResponsePoller, err := publicIPAddressesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, "test-ip", armnetwork.PublicIPAddress{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			IdleTimeoutInMinutes:     to.Ptr[int32](10),
			PublicIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv4),
			PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
		},
		SKU: &armnetwork.PublicIPAddressSKU{
			Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
			Tier: to.Ptr(armnetwork.PublicIPAddressSKUTierRegional),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var publicIPAddressesClientCreateOrUpdateResponse *armnetwork.PublicIPAddressesClientCreateOrUpdateResponse
	publicIPAddressesClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, publicIPAddressesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.publicIpAddressId = *publicIPAddressesClientCreateOrUpdateResponse.ID
}

// Microsoft.Network/bastionHosts/{bastionHostName}
func (testsuite *BastionHostTestSuite) TestBastionHosts() {
	var err error
	// From step BastionHosts_CreateOrUpdate
	fmt.Println("Call operation: BastionHosts_CreateOrUpdate")
	bastionHostsClient, err := armnetwork.NewBastionHostsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	bastionHostsClientCreateOrUpdateResponsePoller, err := bastionHostsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.bastionHostName, armnetwork.BastionHost{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.BastionHostPropertiesFormat{
			IPConfigurations: []*armnetwork.BastionHostIPConfiguration{
				{
					Name: to.Ptr("bastionHostIpConfiguration"),
					Properties: &armnetwork.BastionHostIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.SubResource{
							ID: to.Ptr(testsuite.publicIpAddressId),
						},
						Subnet: &armnetwork.SubResource{
							ID: to.Ptr(testsuite.subnetId),
						},
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, bastionHostsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step BastionHosts_List
	fmt.Println("Call operation: BastionHosts_List")
	bastionHostsClientNewListPager := bastionHostsClient.NewListPager(nil)
	for bastionHostsClientNewListPager.More() {
		_, err := bastionHostsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BastionHosts_ListByResourceGroup
	fmt.Println("Call operation: BastionHosts_ListByResourceGroup")
	bastionHostsClientNewListByResourceGroupPager := bastionHostsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for bastionHostsClientNewListByResourceGroupPager.More() {
		_, err := bastionHostsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BastionHosts_Get
	fmt.Println("Call operation: BastionHosts_Get")
	_, err = bastionHostsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.bastionHostName, nil)
	testsuite.Require().NoError(err)

	// From step BastionHosts_UpdateTags
	fmt.Println("Call operation: BastionHosts_UpdateTags")
	bastionHostsClientUpdateTagsResponsePoller, err := bastionHostsClient.BeginUpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.bastionHostName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, bastionHostsClientUpdateTagsResponsePoller)
	testsuite.Require().NoError(err)

	// From step BastionHosts_Delete
	fmt.Println("Call operation: BastionHosts_Delete")
	bastionHostsClientDeleteResponsePoller, err := bastionHostsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.bastionHostName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, bastionHostsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
