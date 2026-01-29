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

type LoadBalancerTestSuite struct {
	suite.Suite

	ctx                       context.Context
	cred                      azcore.TokenCredential
	options                   *arm.ClientOptions
	backendAddressPoolName    string
	frontendIPConfigurationId string
	inboundNatRuleName        string
	loadBalancerName          string
	subnetId                  string
	location                  string
	resourceGroupName         string
	subscriptionId            string
}

func (testsuite *LoadBalancerTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.backendAddressPoolName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "backendadd", 16, false)
	testsuite.inboundNatRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "inboundnat", 16, false)
	testsuite.loadBalancerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "loadbalanc", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *LoadBalancerTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestLoadBalancerTestSuite(t *testing.T) {
	suite.Run(t, new(LoadBalancerTestSuite))
}

func (testsuite *LoadBalancerTestSuite) Prepare() {
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

	// From step LoadBalancers_CreateOrUpdate
	fmt.Println("Call operation: LoadBalancers_CreateOrUpdate")
	loadBalancersClient, err := armnetwork.NewLoadBalancersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancersClientCreateOrUpdateResponsePoller, err := loadBalancersClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, armnetwork.LoadBalancer{
		Location: to.Ptr(testsuite.location),
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Name: to.Ptr("frontendipconf"),
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						Subnet: &armnetwork.Subnet{
							ID: to.Ptr(testsuite.subnetId),
						},
					},
				}},
		},
		SKU: &armnetwork.LoadBalancerSKU{
			Name: to.Ptr(armnetwork.LoadBalancerSKUNameStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var loadBalancersClientCreateOrUpdateResponse *armnetwork.LoadBalancersClientCreateOrUpdateResponse
	loadBalancersClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, loadBalancersClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.frontendIPConfigurationId = *loadBalancersClientCreateOrUpdateResponse.Properties.FrontendIPConfigurations[0].ID
}

// Microsoft.Network/loadBalancers/{loadBalancerName}
func (testsuite *LoadBalancerTestSuite) TestLoadBalancers() {
	var err error
	// From step LoadBalancers_ListAll
	fmt.Println("Call operation: LoadBalancers_ListAll")
	loadBalancersClient, err := armnetwork.NewLoadBalancersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancersClientNewListAllPager := loadBalancersClient.NewListAllPager(nil)
	for loadBalancersClientNewListAllPager.More() {
		_, err := loadBalancersClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LoadBalancers_List
	fmt.Println("Call operation: LoadBalancers_List")
	loadBalancersClientNewListPager := loadBalancersClient.NewListPager(testsuite.resourceGroupName, nil)
	for loadBalancersClientNewListPager.More() {
		_, err := loadBalancersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LoadBalancers_Get
	fmt.Println("Call operation: LoadBalancers_Get")
	_, err = loadBalancersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, &armnetwork.LoadBalancersClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step LoadBalancers_UpdateTags
	fmt.Println("Call operation: LoadBalancers_UpdateTags")
	_, err = loadBalancersClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/loadBalancers/{loadBalancerName}/backendAddressPools/{backendAddressPoolName}
func (testsuite *LoadBalancerTestSuite) TestLoadBalancerBackendAddressPools() {
	var err error
	// From step LoadBalancerBackendAddressPools_CreateOrUpdate
	fmt.Println("Call operation: LoadBalancerBackendAddressPools_CreateOrUpdate")
	loadBalancerBackendAddressPoolsClient, err := armnetwork.NewLoadBalancerBackendAddressPoolsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancerBackendAddressPoolsClientCreateOrUpdateResponsePoller, err := loadBalancerBackendAddressPoolsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, testsuite.backendAddressPoolName, armnetwork.BackendAddressPool{
		Properties: &armnetwork.BackendAddressPoolPropertiesFormat{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, loadBalancerBackendAddressPoolsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step LoadBalancerBackendAddressPools_List
	fmt.Println("Call operation: LoadBalancerBackendAddressPools_List")
	loadBalancerBackendAddressPoolsClientNewListPager := loadBalancerBackendAddressPoolsClient.NewListPager(testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	for loadBalancerBackendAddressPoolsClientNewListPager.More() {
		_, err := loadBalancerBackendAddressPoolsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LoadBalancerBackendAddressPools_Get
	fmt.Println("Call operation: LoadBalancerBackendAddressPools_Get")
	_, err = loadBalancerBackendAddressPoolsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, testsuite.backendAddressPoolName, nil)
	testsuite.Require().NoError(err)

	// From step LoadBalancerBackendAddressPools_Delete
	fmt.Println("Call operation: LoadBalancerBackendAddressPools_Delete")
	loadBalancerBackendAddressPoolsClientDeleteResponsePoller, err := loadBalancerBackendAddressPoolsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, testsuite.backendAddressPoolName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, loadBalancerBackendAddressPoolsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/loadBalancers/{loadBalancerName}/frontendIPConfigurations
func (testsuite *LoadBalancerTestSuite) TestLoadBalancerFrontendIpConfigurations() {
	var err error
	// From step LoadBalancerFrontendIPConfigurations_List
	fmt.Println("Call operation: LoadBalancerFrontendIPConfigurations_List")
	loadBalancerFrontendIPConfigurationsClient, err := armnetwork.NewLoadBalancerFrontendIPConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancerFrontendIPConfigurationsClientNewListPager := loadBalancerFrontendIPConfigurationsClient.NewListPager(testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	for loadBalancerFrontendIPConfigurationsClientNewListPager.More() {
		_, err := loadBalancerFrontendIPConfigurationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step LoadBalancerFrontendIPConfigurations_Get
	fmt.Println("Call operation: LoadBalancerFrontendIPConfigurations_Get")
	_, err = loadBalancerFrontendIPConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, "frontendipconf", nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/loadBalancers/{loadBalancerName}/inboundNatRules/{inboundNatRuleName}
func (testsuite *LoadBalancerTestSuite) TestInboundNatRules() {
	var err error
	// From step InboundNatRules_CreateOrUpdate
	fmt.Println("Call operation: InboundNatRules_CreateOrUpdate")
	inboundNatRulesClient, err := armnetwork.NewInboundNatRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	inboundNatRulesClientCreateOrUpdateResponsePoller, err := inboundNatRulesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, testsuite.inboundNatRuleName, armnetwork.InboundNatRule{
		Properties: &armnetwork.InboundNatRulePropertiesFormat{
			BackendPort:      to.Ptr[int32](3389),
			EnableFloatingIP: to.Ptr(false),
			EnableTCPReset:   to.Ptr(false),
			FrontendIPConfiguration: &armnetwork.SubResource{
				ID: to.Ptr(testsuite.frontendIPConfigurationId),
			},
			FrontendPort:         to.Ptr[int32](3390),
			IdleTimeoutInMinutes: to.Ptr[int32](4),
			Protocol:             to.Ptr(armnetwork.TransportProtocolTCP),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, inboundNatRulesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step InboundNatRules_List
	fmt.Println("Call operation: InboundNatRules_List")
	inboundNatRulesClientNewListPager := inboundNatRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	for inboundNatRulesClientNewListPager.More() {
		_, err := inboundNatRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step InboundNatRules_Get
	fmt.Println("Call operation: InboundNatRules_Get")
	_, err = inboundNatRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, testsuite.inboundNatRuleName, &armnetwork.InboundNatRulesClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step InboundNatRules_Delete
	fmt.Println("Call operation: InboundNatRules_Delete")
	inboundNatRulesClientDeleteResponsePoller, err := inboundNatRulesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, testsuite.inboundNatRuleName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, inboundNatRulesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/loadBalancers/{loadBalancerName}/loadBalancingRules
func (testsuite *LoadBalancerTestSuite) TestLoadBalancerLoadBalancingRules() {
	var err error
	// From step LoadBalancerLoadBalancingRules_List
	fmt.Println("Call operation: LoadBalancerLoadBalancingRules_List")
	loadBalancerLoadBalancingRulesClient, err := armnetwork.NewLoadBalancerLoadBalancingRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancerLoadBalancingRulesClientNewListPager := loadBalancerLoadBalancingRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	for loadBalancerLoadBalancingRulesClientNewListPager.More() {
		_, err := loadBalancerLoadBalancingRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Network/loadBalancers/{loadBalancerName}/outboundRules
func (testsuite *LoadBalancerTestSuite) TestLoadBalancerOutboundRules() {
	var err error
	// From step LoadBalancerOutboundRules_List
	fmt.Println("Call operation: LoadBalancerOutboundRules_List")
	loadBalancerOutboundRulesClient, err := armnetwork.NewLoadBalancerOutboundRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancerOutboundRulesClientNewListPager := loadBalancerOutboundRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	for loadBalancerOutboundRulesClientNewListPager.More() {
		_, err := loadBalancerOutboundRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Network/loadBalancers/{loadBalancerName}/networkInterfaces
func (testsuite *LoadBalancerTestSuite) TestLoadBalancerNetworkInterfaces() {
	var err error
	// From step LoadBalancerNetworkInterfaces_List
	fmt.Println("Call operation: LoadBalancerNetworkInterfaces_List")
	loadBalancerNetworkInterfacesClient, err := armnetwork.NewLoadBalancerNetworkInterfacesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancerNetworkInterfacesClientNewListPager := loadBalancerNetworkInterfacesClient.NewListPager(testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	for loadBalancerNetworkInterfacesClientNewListPager.More() {
		_, err := loadBalancerNetworkInterfacesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Network/loadBalancers/{loadBalancerName}/probes
func (testsuite *LoadBalancerTestSuite) TestLoadBalancerProbes() {
	var err error
	// From step LoadBalancerProbes_List
	fmt.Println("Call operation: LoadBalancerProbes_List")
	loadBalancerProbesClient, err := armnetwork.NewLoadBalancerProbesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancerProbesClientNewListPager := loadBalancerProbesClient.NewListPager(testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	for loadBalancerProbesClientNewListPager.More() {
		_, err := loadBalancerProbesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *LoadBalancerTestSuite) Cleanup() {
	var err error
	// From step LoadBalancers_Delete
	fmt.Println("Call operation: LoadBalancers_Delete")
	loadBalancersClient, err := armnetwork.NewLoadBalancersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	loadBalancersClientDeleteResponsePoller, err := loadBalancersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.loadBalancerName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, loadBalancersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
