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

type DnsforwardingrulesetTestSuite struct {
	suite.Suite

	ctx                      context.Context
	cred                     azcore.TokenCredential
	options                  *arm.ClientOptions
	armEndpoint              string
	dnsForwardingRulesetName string
	dnsResolverName          string
	forwardingRuleName       string
	outboundEndpointId       string
	outboundEndpointName     string
	subnetId                 string
	virtualNetworkId         string
	virtualNetworkLinkName   string
	virtualNetworkName       string
	location                 string
	resourceGroupName        string
	subscriptionId           string
}

func (testsuite *DnsforwardingrulesetTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.dnsForwardingRulesetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "dnsforwa", 14, false)
	testsuite.dnsResolverName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "dnsresol", 14, false)
	testsuite.forwardingRuleName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "forwardi", 14, false)
	testsuite.outboundEndpointName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "outbound", 14, false)
	testsuite.virtualNetworkLinkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualn", 14, false)
	testsuite.virtualNetworkName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualn", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DnsforwardingrulesetTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDnsforwardingrulesetTestSuite(t *testing.T) {
	suite.Run(t, new(DnsforwardingrulesetTestSuite))
}

func (testsuite *DnsforwardingrulesetTestSuite) Prepare() {
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
	fmt.Println("Call operation: DNSResolvers_CreateOrUpdate")
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
	var outboundEndpointsClientCreateOrUpdateResponse *armdnsresolver.OutboundEndpointsClientCreateOrUpdateResponse
	outboundEndpointsClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, outboundEndpointsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.outboundEndpointId = *outboundEndpointsClientCreateOrUpdateResponse.ID

	// From step DnsForwardingRulesets_CreateOrUpdate
	fmt.Println("Call operation: DnsForwardingRulesets_CreateOrUpdate")
	dNSForwardingRulesetsClient, err := armdnsresolver.NewDNSForwardingRulesetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dNSForwardingRulesetsClientCreateOrUpdateResponsePoller, err := dNSForwardingRulesetsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, armdnsresolver.DNSForwardingRuleset{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armdnsresolver.DNSForwardingRulesetProperties{
			DNSResolverOutboundEndpoints: []*armdnsresolver.SubResource{
				{
					ID: to.Ptr(testsuite.outboundEndpointId),
				}},
		},
	}, &armdnsresolver.DNSForwardingRulesetsClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dNSForwardingRulesetsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/dnsForwardingRulesets/{dnsForwardingRulesetName}
func (testsuite *DnsforwardingrulesetTestSuite) TestDnsForwardingRulesets() {
	var err error
	// From step DnsForwardingRulesets_List
	fmt.Println("Call operation: DnsForwardingRulesets_List")
	dNSForwardingRulesetsClient, err := armdnsresolver.NewDNSForwardingRulesetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dNSForwardingRulesetsClientNewListPager := dNSForwardingRulesetsClient.NewListPager(&armdnsresolver.DNSForwardingRulesetsClientListOptions{Top: nil})
	for dNSForwardingRulesetsClientNewListPager.More() {
		_, err := dNSForwardingRulesetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DnsForwardingRulesets_Get
	fmt.Println("Call operation: DnsForwardingRulesets_Get")
	_, err = dNSForwardingRulesetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, nil)
	testsuite.Require().NoError(err)

	// From step DnsForwardingRulesets_ListByResourceGroup
	fmt.Println("Call operation: DnsForwardingRulesets_ListByResourceGroup")
	dNSForwardingRulesetsClientNewListByResourceGroupPager := dNSForwardingRulesetsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armdnsresolver.DNSForwardingRulesetsClientListByResourceGroupOptions{Top: nil})
	for dNSForwardingRulesetsClientNewListByResourceGroupPager.More() {
		_, err := dNSForwardingRulesetsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DnsForwardingRulesets_Update
	fmt.Println("Call operation: DnsForwardingRulesets_Update")
	dNSForwardingRulesetsClientUpdateResponsePoller, err := dNSForwardingRulesetsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, armdnsresolver.DNSForwardingRulesetPatch{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, &armdnsresolver.DNSForwardingRulesetsClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dNSForwardingRulesetsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DnsForwardingRulesets_ListByVirtualNetwork
	fmt.Println("Call operation: DnsForwardingRulesets_ListByVirtualNetwork")
	dNSForwardingRulesetsClientNewListByVirtualNetworkPager := dNSForwardingRulesetsClient.NewListByVirtualNetworkPager(testsuite.resourceGroupName, testsuite.virtualNetworkName, &armdnsresolver.DNSForwardingRulesetsClientListByVirtualNetworkOptions{Top: nil})
	for dNSForwardingRulesetsClientNewListByVirtualNetworkPager.More() {
		_, err := dNSForwardingRulesetsClientNewListByVirtualNetworkPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.Network/dnsForwardingRulesets/{dnsForwardingRulesetName}/forwardingRules/{forwardingRuleName}
func (testsuite *DnsforwardingrulesetTestSuite) TestForwardingRules() {
	var err error
	// From step ForwardingRules_CreateOrUpdate
	fmt.Println("Call operation: ForwardingRules_CreateOrUpdate")
	forwardingRulesClient, err := armdnsresolver.NewForwardingRulesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = forwardingRulesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.forwardingRuleName, armdnsresolver.ForwardingRule{
		Properties: &armdnsresolver.ForwardingRuleProperties{
			DomainName:          to.Ptr("contoso.com."),
			ForwardingRuleState: to.Ptr(armdnsresolver.ForwardingRuleStateEnabled),
			Metadata: map[string]*string{
				"additionalProp1": to.Ptr("value1"),
			},
			TargetDNSServers: []*armdnsresolver.TargetDNSServer{
				{
					IPAddress: to.Ptr("10.0.0.1"),
					Port:      to.Ptr[int32](53),
				},
				{
					IPAddress: to.Ptr("10.0.0.2"),
					Port:      to.Ptr[int32](53),
				}},
		},
	}, &armdnsresolver.ForwardingRulesClientCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step ForwardingRules_List
	fmt.Println("Call operation: ForwardingRules_List")
	forwardingRulesClientNewListPager := forwardingRulesClient.NewListPager(testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, &armdnsresolver.ForwardingRulesClientListOptions{Top: nil})
	for forwardingRulesClientNewListPager.More() {
		_, err := forwardingRulesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ForwardingRules_Get
	fmt.Println("Call operation: ForwardingRules_Get")
	_, err = forwardingRulesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.forwardingRuleName, nil)
	testsuite.Require().NoError(err)

	// From step ForwardingRules_Update
	fmt.Println("Call operation: ForwardingRules_Update")
	_, err = forwardingRulesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.forwardingRuleName, armdnsresolver.ForwardingRulePatch{
		Properties: &armdnsresolver.ForwardingRulePatchProperties{
			ForwardingRuleState: to.Ptr(armdnsresolver.ForwardingRuleStateDisabled),
			Metadata: map[string]*string{
				"additionalProp2": to.Ptr("value2"),
			},
		},
	}, &armdnsresolver.ForwardingRulesClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step ForwardingRules_Delete
	fmt.Println("Call operation: ForwardingRules_Delete")
	_, err = forwardingRulesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.forwardingRuleName, &armdnsresolver.ForwardingRulesClientDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.Network/dnsForwardingRulesets/{dnsForwardingRulesetName}/virtualNetworkLinks/{virtualNetworkLinkName}
func (testsuite *DnsforwardingrulesetTestSuite) TestVirtualNetworkLinks() {
	var err error
	// From step VirtualNetworkLinks_CreateOrUpdate
	fmt.Println("Call operation: VirtualNetworkLinks_CreateOrUpdate")
	virtualNetworkLinksClient, err := armdnsresolver.NewVirtualNetworkLinksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualNetworkLinksClientCreateOrUpdateResponsePoller, err := virtualNetworkLinksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.virtualNetworkLinkName, armdnsresolver.VirtualNetworkLink{
		Properties: &armdnsresolver.VirtualNetworkLinkProperties{
			Metadata: map[string]*string{
				"additionalProp1": to.Ptr("value1"),
			},
			VirtualNetwork: &armdnsresolver.SubResource{
				ID: to.Ptr(testsuite.virtualNetworkId),
			},
		},
	}, &armdnsresolver.VirtualNetworkLinksClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkLinksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkLinks_List
	fmt.Println("Call operation: VirtualNetworkLinks_List")
	virtualNetworkLinksClientNewListPager := virtualNetworkLinksClient.NewListPager(testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, &armdnsresolver.VirtualNetworkLinksClientListOptions{Top: nil})
	for virtualNetworkLinksClientNewListPager.More() {
		_, err := virtualNetworkLinksClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualNetworkLinks_Get
	fmt.Println("Call operation: VirtualNetworkLinks_Get")
	_, err = virtualNetworkLinksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.virtualNetworkLinkName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkLinks_Update
	fmt.Println("Call operation: VirtualNetworkLinks_Update")
	virtualNetworkLinksClientUpdateResponsePoller, err := virtualNetworkLinksClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.virtualNetworkLinkName, armdnsresolver.VirtualNetworkLinkPatch{
		Properties: &armdnsresolver.VirtualNetworkLinkPatchProperties{
			Metadata: map[string]*string{
				"additionalProp1": to.Ptr("value1"),
			},
		},
	}, &armdnsresolver.VirtualNetworkLinksClientBeginUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkLinksClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualNetworkLinks_Delete
	fmt.Println("Call operation: VirtualNetworkLinks_Delete")
	virtualNetworkLinksClientDeleteResponsePoller, err := virtualNetworkLinksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, testsuite.virtualNetworkLinkName, &armdnsresolver.VirtualNetworkLinksClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualNetworkLinksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *DnsforwardingrulesetTestSuite) Cleanup() {
	var err error
	// From step DnsForwardingRulesets_Delete
	fmt.Println("Call operation: DnsForwardingRulesets_Delete")
	dNSForwardingRulesetsClient, err := armdnsresolver.NewDNSForwardingRulesetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dNSForwardingRulesetsClientDeleteResponsePoller, err := dNSForwardingRulesetsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.dnsForwardingRulesetName, &armdnsresolver.DNSForwardingRulesetsClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dNSForwardingRulesetsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
