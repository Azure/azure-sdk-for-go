//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdnsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dnsresolver/armdnsresolver/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type DnsresolverTestSuite struct {
	suite.Suite

	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	armEndpoint          string
	dnsResolverName      string
	inboundEndpointName  string
	outboundEndpointName string
	subnetId             string
	virtualNetworkId     string
	virtualNetworkName   string
	location             string
	resourceGroupName    string
	subscriptionId       string
}

func (testsuite *DnsresolverTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.dnsResolverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "dnsresol", 14, false)
	testsuite.inboundEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "inbounde", 14, false)
	testsuite.outboundEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "outbound", 14, false)
	testsuite.virtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualn", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DnsresolverTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDnsresolverTestSuite(t *testing.T) {
	suite.Run(t, new(DnsresolverTestSuite))
}

func (testsuite *DnsresolverTestSuite) Prepare() {
	var err error
	// From step Create_VirtualNetwork
	template := map[string]any{
		"$schema":        "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"subnetId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks/subnets', parameters('virtualNetworkName'), 'default')]",
			},
			"virtualNetworkId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Network/virtualNetworks', parameters('virtualNetworkName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"virtualNetworkName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.virtualNetworkName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('virtualNetworkName')]",
				"type":       "Microsoft.Network/virtualNetworks",
				"apiVersion": "2019-02-01",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"addressSpace": map[string]any{
						"addressPrefixes": []any{
							"10.0.0.0/16",
						},
					},
					"subnets": []any{
						map[string]any{
							"name": "default",
							"properties": map[string]any{
								"addressPrefix": "10.0.0.0/24",
							},
						},
					},
				},
			},
		},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_VirtualNetwork", &deployment)
	testsuite.Require().NoError(err)
	testsuite.subnetId = deploymentExtend.Properties.Outputs.(map[string]interface{})["subnetId"].(map[string]interface{})["value"].(string)
	testsuite.virtualNetworkId = deploymentExtend.Properties.Outputs.(map[string]interface{})["virtualNetworkId"].(map[string]interface{})["value"].(string)

	// From step DnsResolvers_CreateOrUpdate
	fmt.Println("Call operation: DnsResolvers_CreateOrUpdate")
	dNSResolversClient, err := armdnsresolver.NewDNSResolversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dNSResolversClientCreateOrUpdateResponsePoller, err := dNSResolversClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, armdnsresolver.DNSResolver{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armdnsresolver.Properties{
			VirtualNetwork: &armdnsresolver.SubResource{
				ID: to.Ptr(testsuite.virtualNetworkId),
			},
		},
	}, &armdnsresolver.DNSResolversClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dNSResolversClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/dnsResolvers/{dnsResolverName}
func (testsuite *DnsresolverTestSuite) TestDnsResolvers() {
	var err error
	// From step DnsResolvers_List
	fmt.Println("Call operation: DnsResolvers_List")
	dNSResolversClient, err := armdnsresolver.NewDNSResolversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dNSResolversClientNewListPager := dNSResolversClient.NewListPager(&armdnsresolver.DNSResolversClientListOptions{Top: nil})
	for dNSResolversClientNewListPager.More() {
		_, err := dNSResolversClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DnsResolvers_ListByResourceGroup
	fmt.Println("Call operation: DnsResolvers_ListByResourceGroup")
	dNSResolversClientNewListByResourceGroupPager := dNSResolversClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armdnsresolver.DNSResolversClientListByResourceGroupOptions{Top: nil})
	for dNSResolversClientNewListByResourceGroupPager.More() {
		_, err := dNSResolversClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DnsResolvers_Get
	fmt.Println("Call operation: DnsResolvers_Get")
	_, err = dNSResolversClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, nil)
	testsuite.Require().NoError(err)

	// From step DnsResolvers_Update
	fmt.Println("Call operation: DnsResolvers_Update")
	dNSResolversClientUpdateResponsePoller, err := dNSResolversClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, armdnsresolver.Patch{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, &armdnsresolver.DNSResolversClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dNSResolversClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DnsResolvers_ListByVirtualNetwork
	fmt.Println("Call operation: DnsResolvers_ListByVirtualNetwork")
	dNSResolversClientNewListByVirtualNetworkPager := dNSResolversClient.NewListByVirtualNetworkPager(testsuite.resourceGroupName, testsuite.virtualNetworkName, &armdnsresolver.DNSResolversClientListByVirtualNetworkOptions{Top: nil})
	for dNSResolversClientNewListByVirtualNetworkPager.More() {
		_, err := dNSResolversClientNewListByVirtualNetworkPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Network/dnsResolvers/{dnsResolverName}/inboundEndpoints/{inboundEndpointName}
func (testsuite *DnsresolverTestSuite) TestInboundEndpoints() {
	var err error
	// From step InboundEndpoints_CreateOrUpdate
	fmt.Println("Call operation: InboundEndpoints_CreateOrUpdate")
	inboundEndpointsClient, err := armdnsresolver.NewInboundEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	inboundEndpointsClientCreateOrUpdateResponsePoller, err := inboundEndpointsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.inboundEndpointName, armdnsresolver.InboundEndpoint{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armdnsresolver.InboundEndpointProperties{
			IPConfigurations: []*armdnsresolver.IPConfiguration{
				{
					PrivateIPAllocationMethod: to.Ptr(armdnsresolver.IPAllocationMethodDynamic),
					Subnet: &armdnsresolver.SubResource{
						ID: to.Ptr(testsuite.subnetId),
					},
				}},
		},
	}, &armdnsresolver.InboundEndpointsClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, inboundEndpointsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step InboundEndpoints_List
	fmt.Println("Call operation: InboundEndpoints_List")
	inboundEndpointsClientNewListPager := inboundEndpointsClient.NewListPager(testsuite.resourceGroupName, testsuite.dnsResolverName, &armdnsresolver.InboundEndpointsClientListOptions{Top: nil})
	for inboundEndpointsClientNewListPager.More() {
		_, err := inboundEndpointsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step InboundEndpoints_Get
	fmt.Println("Call operation: InboundEndpoints_Get")
	_, err = inboundEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.inboundEndpointName, nil)
	testsuite.Require().NoError(err)

	// From step InboundEndpoints_Update
	fmt.Println("Call operation: InboundEndpoints_Update")
	inboundEndpointsClientUpdateResponsePoller, err := inboundEndpointsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.inboundEndpointName, armdnsresolver.InboundEndpointPatch{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, &armdnsresolver.InboundEndpointsClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, inboundEndpointsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step InboundEndpoints_Delete
	fmt.Println("Call operation: InboundEndpoints_Delete")
	inboundEndpointsClientDeleteResponsePoller, err := inboundEndpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.inboundEndpointName, &armdnsresolver.InboundEndpointsClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, inboundEndpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/dnsResolvers/{dnsResolverName}/outboundEndpoints/{outboundEndpointName}
func (testsuite *DnsresolverTestSuite) TestOutboundEndpoints() {
	var err error
	// From step OutboundEndpoints_CreateOrUpdate
	fmt.Println("Call operation: OutboundEndpoints_CreateOrUpdate")
	outboundEndpointsClient, err := armdnsresolver.NewOutboundEndpointsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	outboundEndpointsClientCreateOrUpdateResponsePoller, err := outboundEndpointsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.outboundEndpointName, armdnsresolver.OutboundEndpoint{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armdnsresolver.OutboundEndpointProperties{
			Subnet: &armdnsresolver.SubResource{
				ID: to.Ptr(testsuite.subnetId),
			},
		},
	}, &armdnsresolver.OutboundEndpointsClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, outboundEndpointsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step OutboundEndpoints_List
	fmt.Println("Call operation: OutboundEndpoints_List")
	outboundEndpointsClientNewListPager := outboundEndpointsClient.NewListPager(testsuite.resourceGroupName, testsuite.dnsResolverName, &armdnsresolver.OutboundEndpointsClientListOptions{Top: nil})
	for outboundEndpointsClientNewListPager.More() {
		_, err := outboundEndpointsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step OutboundEndpoints_Get
	fmt.Println("Call operation: OutboundEndpoints_Get")
	_, err = outboundEndpointsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.outboundEndpointName, nil)
	testsuite.Require().NoError(err)

	// From step OutboundEndpoints_Update
	fmt.Println("Call operation: OutboundEndpoints_Update")
	outboundEndpointsClientUpdateResponsePoller, err := outboundEndpointsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.outboundEndpointName, armdnsresolver.OutboundEndpointPatch{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, &armdnsresolver.OutboundEndpointsClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, outboundEndpointsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step OutboundEndpoints_Delete
	fmt.Println("Call operation: OutboundEndpoints_Delete")
	outboundEndpointsClientDeleteResponsePoller, err := outboundEndpointsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, testsuite.outboundEndpointName, &armdnsresolver.OutboundEndpointsClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, outboundEndpointsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *DnsresolverTestSuite) Cleanup() {
	var err error
	// From step DnsResolvers_Delete
	fmt.Println("Call operation: DnsResolvers_Delete")
	dNSResolversClient, err := armdnsresolver.NewDNSResolversClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dNSResolversClientDeleteResponsePoller, err := dNSResolversClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsResolverName, &armdnsresolver.DNSResolversClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dNSResolversClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
