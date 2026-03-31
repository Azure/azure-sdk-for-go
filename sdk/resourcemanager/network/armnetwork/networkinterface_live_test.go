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

type NetworkInterfaceTestSuite struct {
	suite.Suite

	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	networkInterfaceName string
	subnetId             string
	virtualNetworkName   string
	location             string
	resourceGroupName    string
	subscriptionId       string
}

func (testsuite *NetworkInterfaceTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.networkInterfaceName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkint", 16, false)
	testsuite.virtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vnetinterfacena", 21, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NetworkInterfaceTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestNetworkInterfaceTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkInterfaceTestSuite))
}

func (testsuite *NetworkInterfaceTestSuite) Prepare() {
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
					Name: to.Ptr("test-1"),
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
}

// Microsoft.Network/networkInterfaces/{networkInterfaceName}
func (testsuite *NetworkInterfaceTestSuite) TestNetworkInterfaces() {
	var err error
	// From step NetworkInterfaces_CreateOrUpdate
	fmt.Println("Call operation: NetworkInterfaces_CreateOrUpdate")
	interfacesClient, err := armnetwork.NewInterfacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	interfacesClientCreateOrUpdateResponsePoller, err := interfacesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkInterfaceName, armnetwork.Interface{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.InterfacePropertiesFormat{
			EnableAcceleratedNetworking: to.Ptr(true),
			IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
				{
					Name: to.Ptr("ipconfig1"),
					Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
						Subnet: &armnetwork.Subnet{
							ID: to.Ptr(testsuite.subnetId),
						},
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, interfacesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step NetworkInterfaces_Get
	fmt.Println("Call operation: NetworkInterfaces_Get")
	_, err = interfacesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkInterfaceName, &armnetwork.InterfacesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step NetworkInterfaces_UpdateTags
	fmt.Println("Call operation: NetworkInterfaces_UpdateTags")
	_, err = interfacesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkInterfaceName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NetworkInterfaces_ListAll
	fmt.Println("Call operation: NetworkInterfaces_ListAll")
	interfacesClientNewListAllPager := interfacesClient.NewListAllPager(nil)
	for interfacesClientNewListAllPager.More() {
		_, err := interfacesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkInterfaces_List
	fmt.Println("Call operation: NetworkInterfaces_List")
	interfacesClientNewListPager := interfacesClient.NewListPager(testsuite.resourceGroupName, nil)
	for interfacesClientNewListPager.More() {
		_, err := interfacesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkInterfaceIPConfigurations_List
	fmt.Println("Call operation: NetworkInterfaceIPConfigurations_List")
	interfaceIPConfigurationsClient, err := armnetwork.NewInterfaceIPConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	interfaceIPConfigurationsClientNewListPager := interfaceIPConfigurationsClient.NewListPager(testsuite.resourceGroupName, testsuite.networkInterfaceName, nil)
	for interfaceIPConfigurationsClientNewListPager.More() {
		_, err := interfaceIPConfigurationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkInterfaceIPConfigurations_Get
	fmt.Println("Call operation: NetworkInterfaceIPConfigurations_Get")
	_, err = interfaceIPConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkInterfaceName, "ipconfig1", nil)
	testsuite.Require().NoError(err)

	// From step NetworkInterfaceLoadBalancers_List
	fmt.Println("Call operation: NetworkInterfaceLoadBalancers_List")
	interfaceLoadBalancersClient, err := armnetwork.NewInterfaceLoadBalancersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	interfaceLoadBalancersClientNewListPager := interfaceLoadBalancersClient.NewListPager(testsuite.resourceGroupName, testsuite.networkInterfaceName, nil)
	for interfaceLoadBalancersClientNewListPager.More() {
		_, err := interfaceLoadBalancersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkInterfaces_Delete
	fmt.Println("Call operation: NetworkInterfaces_Delete")
	interfacesClientDeleteResponsePoller, err := interfacesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkInterfaceName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, interfacesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
