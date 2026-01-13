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

type NetworkProfileTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	networkProfileName string
	subnetId           string
	virtualNetworkName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *NetworkProfileTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.networkProfileName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkpro", 16, false)
	testsuite.virtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vnetprofilena", 19, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *NetworkProfileTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestNetworkProfileTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkProfileTestSuite))
}

func (testsuite *NetworkProfileTestSuite) Prepare() {
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

// Microsoft.Network/networkProfiles/{networkProfileName}
func (testsuite *NetworkProfileTestSuite) TestNetworkProfiles() {
	var err error
	// From step NetworkProfiles_CreateOrUpdate
	fmt.Println("Call operation: NetworkProfiles_CreateOrUpdate")
	profilesClient, err := armnetwork.NewProfilesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = profilesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkProfileName, armnetwork.Profile{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.ProfilePropertiesFormat{
			ContainerNetworkInterfaceConfigurations: []*armnetwork.ContainerNetworkInterfaceConfiguration{
				{
					Name: to.Ptr("eth1"),
					Properties: &armnetwork.ContainerNetworkInterfaceConfigurationPropertiesFormat{
						IPConfigurations: []*armnetwork.IPConfigurationProfile{
							{
								Name: to.Ptr("ipconfig1"),
								Properties: &armnetwork.IPConfigurationProfilePropertiesFormat{
									Subnet: &armnetwork.Subnet{
										ID: to.Ptr(testsuite.subnetId),
									},
								},
							}},
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NetworkProfiles_ListAll
	fmt.Println("Call operation: NetworkProfiles_ListAll")
	profilesClientNewListAllPager := profilesClient.NewListAllPager(nil)
	for profilesClientNewListAllPager.More() {
		_, err := profilesClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkProfiles_List
	fmt.Println("Call operation: NetworkProfiles_List")
	profilesClientNewListPager := profilesClient.NewListPager(testsuite.resourceGroupName, nil)
	for profilesClientNewListPager.More() {
		_, err := profilesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkProfiles_Get
	fmt.Println("Call operation: NetworkProfiles_Get")
	_, err = profilesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkProfileName, &armnetwork.ProfilesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step NetworkProfiles_UpdateTags
	fmt.Println("Call operation: NetworkProfiles_UpdateTags")
	_, err = profilesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkProfileName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NetworkProfiles_Delete
	fmt.Println("Call operation: NetworkProfiles_Delete")
	profilesClientDeleteResponsePoller, err := profilesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkProfileName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, profilesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
