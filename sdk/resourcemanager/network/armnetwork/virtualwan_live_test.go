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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v8"
	"github.com/stretchr/testify/suite"
)

type VirtualWanTestSuite struct {
	suite.Suite

	ctx                        context.Context
	cred                       azcore.TokenCredential
	options                    *arm.ClientOptions
	gatewayName                string
	virtualHubId               string
	virtualHubName             string
	virtualWANName             string
	virtualWanId               string
	vpnServerConfigurationName string
	vpnSiteName                string
	routeMapName               string
	location                   string
	resourceGroupName          string
	subscriptionId             string
}

func (testsuite *VirtualWanTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.gatewayName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "gatewaynam", 16, false)
	testsuite.virtualHubName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualhub", 16, false)
	testsuite.virtualWANName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "virtualwan", 16, false)
	testsuite.vpnServerConfigurationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vpnserverc", 16, false)
	testsuite.vpnSiteName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "vpnsitenam", 16, false)
	testsuite.routeMapName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "routemapna", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *VirtualWanTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestVirtualWanTestSuite(t *testing.T) {
	suite.Run(t, new(VirtualWanTestSuite))
}

func (testsuite *VirtualWanTestSuite) Prepare() {
	var err error
	// From step VirtualWans_CreateOrUpdate
	fmt.Println("Call operation: VirtualWans_CreateOrUpdate")
	virtualWansClient, err := armnetwork.NewVirtualWansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualWansClientCreateOrUpdateResponsePoller, err := virtualWansClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, "wan1", armnetwork.VirtualWAN{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armnetwork.VirtualWanProperties{
			Type:                 to.Ptr("Standard"),
			DisableVPNEncryption: to.Ptr(false),
		},
	}, nil)
	testsuite.Require().NoError(err)
	var virtualWansClientCreateOrUpdateResponse *armnetwork.VirtualWansClientCreateOrUpdateResponse
	virtualWansClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, virtualWansClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.virtualWanId = *virtualWansClientCreateOrUpdateResponse.ID

	// From step VpnSites_CreateOrUpdate
	fmt.Println("Call operation: VPNSites_CreateOrUpdate")
	vPNSitesClient, err := armnetwork.NewVPNSitesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vPNSitesClientCreateOrUpdateResponsePoller, err := vPNSitesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.vpnSiteName, armnetwork.VPNSite{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armnetwork.VPNSiteProperties{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.0.0.0/16")},
			},
			IsSecuritySite: to.Ptr(false),
			O365Policy: &armnetwork.O365PolicyProperties{
				BreakOutCategories: &armnetwork.O365BreakOutCategoryPolicies{
					Default:  to.Ptr(false),
					Allow:    to.Ptr(true),
					Optimize: to.Ptr(true),
				},
			},
			VirtualWan: &armnetwork.SubResource{
				ID: to.Ptr(testsuite.virtualWanId),
			},
			VPNSiteLinks: []*armnetwork.VPNSiteLink{
				{
					Name: to.Ptr("vpnSiteLink1"),
					Properties: &armnetwork.VPNSiteLinkProperties{
						BgpProperties: &armnetwork.VPNLinkBgpSettings{
							Asn:               to.Ptr[int64](1234),
							BgpPeeringAddress: to.Ptr("192.168.0.0"),
						},
						IPAddress: to.Ptr("50.50.50.56"),
						LinkProperties: &armnetwork.VPNLinkProviderProperties{
							LinkProviderName: to.Ptr("vendor1"),
							LinkSpeedInMbps:  to.Ptr[int32](0),
						},
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, vPNSitesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualHubs_CreateOrUpdate
	fmt.Println("Call operation: VirtualHubs_CreateOrUpdate")
	virtualHubsClient, err := armnetwork.NewVirtualHubsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualHubsClientCreateOrUpdateResponsePoller, err := virtualHubsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, armnetwork.VirtualHub{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armnetwork.VirtualHubProperties{
			AddressPrefix: to.Ptr("10.168.0.0/24"),
			SKU:           to.Ptr("Standard"),
			VirtualWan: &armnetwork.SubResource{
				ID: to.Ptr(testsuite.virtualWanId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var virtualHubsClientCreateOrUpdateResponse *armnetwork.VirtualHubsClientCreateOrUpdateResponse
	virtualHubsClientCreateOrUpdateResponse, err = testutil.PollForTest(testsuite.ctx, virtualHubsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.virtualHubId = *virtualHubsClientCreateOrUpdateResponse.ID
}

// Microsoft.Network/virtualWans/{VirtualWANName}
func (testsuite *VirtualWanTestSuite) TestVirtualWans() {
	var err error
	// From step VirtualWans_List
	fmt.Println("Call operation: VirtualWans_List")
	virtualWansClient, err := armnetwork.NewVirtualWansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualWansClientNewListPager := virtualWansClient.NewListPager(nil)
	for virtualWansClientNewListPager.More() {
		_, err := virtualWansClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualWans_ListByResourceGroup
	fmt.Println("Call operation: VirtualWans_ListByResourceGroup")
	virtualWansClientNewListByResourceGroupPager := virtualWansClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for virtualWansClientNewListByResourceGroupPager.More() {
		_, err := virtualWansClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualWans_Get
	fmt.Println("Call operation: VirtualWans_Get")
	_, err = virtualWansClient.Get(testsuite.ctx, testsuite.resourceGroupName, "wan1", nil)
	testsuite.Require().NoError(err)

	// From step VirtualWans_UpdateTags
	fmt.Println("Call operation: VirtualWans_UpdateTags")
	_, err = virtualWansClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, "wan1", armnetwork.TagsObject{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/vpnSites/{vpnSiteName}
func (testsuite *VirtualWanTestSuite) TestVpnSites() {
	var vpnSiteLinkName string
	var err error
	// From step VpnSites_List
	fmt.Println("Call operation: VPNSites_List")
	vPNSitesClient, err := armnetwork.NewVPNSitesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vPNSitesClientNewListPager := vPNSitesClient.NewListPager(nil)
	for vPNSitesClientNewListPager.More() {
		_, err := vPNSitesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VpnSites_ListByResourceGroup
	fmt.Println("Call operation: VPNSites_ListByResourceGroup")
	vPNSitesClientNewListByResourceGroupPager := vPNSitesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for vPNSitesClientNewListByResourceGroupPager.More() {
		_, err := vPNSitesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VpnSites_Get
	fmt.Println("Call operation: VPNSites_Get")
	_, err = vPNSitesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vpnSiteName, nil)
	testsuite.Require().NoError(err)

	// From step VpnSites_UpdateTags
	fmt.Println("Call operation: VPNSites_UpdateTags")
	_, err = vPNSitesClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.vpnSiteName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step VpnSiteLinks_ListByVpnSite
	fmt.Println("Call operation: VPNSiteLinks_ListByVPNSite")
	vPNSiteLinksClient, err := armnetwork.NewVPNSiteLinksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vPNSiteLinksClientNewListByVPNSitePager := vPNSiteLinksClient.NewListByVPNSitePager(testsuite.resourceGroupName, testsuite.vpnSiteName, nil)
	for vPNSiteLinksClientNewListByVPNSitePager.More() {
		nextResult, err := vPNSiteLinksClientNewListByVPNSitePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)

		vpnSiteLinkName = *nextResult.Value[0].Name
		break
	}

	// From step VpnSiteLinks_Get
	fmt.Println("Call operation: VPNSiteLinks_Get")
	_, err = vPNSiteLinksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.vpnSiteName, vpnSiteLinkName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/virtualHubs/{virtualHubName}
func (testsuite *VirtualWanTestSuite) TestVirtualHubs() {
	var err error
	// From step VirtualHubs_List
	fmt.Println("Call operation: VirtualHubs_List")
	virtualHubsClient, err := armnetwork.NewVirtualHubsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualHubsClientNewListPager := virtualHubsClient.NewListPager(nil)
	for virtualHubsClientNewListPager.More() {
		_, err := virtualHubsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualHubs_ListByResourceGroup
	fmt.Println("Call operation: VirtualHubs_ListByResourceGroup")
	virtualHubsClientNewListByResourceGroupPager := virtualHubsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for virtualHubsClientNewListByResourceGroupPager.More() {
		_, err := virtualHubsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VirtualHubs_Get
	fmt.Println("Call operation: VirtualHubs_Get")
	_, err = virtualHubsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, nil)
	testsuite.Require().NoError(err)

	// From step VirtualHubs_UpdateTags
	fmt.Println("Call operation: VirtualHubs_UpdateTags")
	_, err = virtualHubsClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/virtualHubs/{virtualHubName}/routeMaps/{routeMapName}
func (testsuite *VirtualWanTestSuite) TestRouteMaps() {
	var err error
	// From step RouteMaps_CreateOrUpdate
	fmt.Println("Call operation: RouteMaps_CreateOrUpdate")
	routeMapsClient, err := armnetwork.NewRouteMapsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	routeMapsClientCreateOrUpdateResponsePoller, err := routeMapsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, testsuite.routeMapName, armnetwork.RouteMap{
		Properties: &armnetwork.RouteMapProperties{
			AssociatedInboundConnections: []*string{
				to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Network/expressRouteGateways/exrGateway1/expressRouteConnections/exrConn1")},
			AssociatedOutboundConnections: []*string{},
			Rules: []*armnetwork.RouteMapRule{
				{
					Name: to.Ptr("rule1"),
					Actions: []*armnetwork.Action{
						{
							Type: to.Ptr(armnetwork.RouteMapActionTypeAdd),
							Parameters: []*armnetwork.Parameter{
								{
									AsPath: []*string{
										to.Ptr("22334")},
									Community:   []*string{},
									RoutePrefix: []*string{},
								}},
						}},
					MatchCriteria: []*armnetwork.Criterion{
						{
							AsPath:         []*string{},
							Community:      []*string{},
							MatchCondition: to.Ptr(armnetwork.RouteMapMatchConditionContains),
							RoutePrefix: []*string{
								to.Ptr("10.0.0.0/8")},
						}},
					NextStepIfMatched: to.Ptr(armnetwork.NextStepContinue),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routeMapsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step RouteMaps_List
	fmt.Println("Call operation: RouteMaps_List")
	routeMapsClientNewListPager := routeMapsClient.NewListPager(testsuite.resourceGroupName, testsuite.virtualHubName, nil)
	for routeMapsClientNewListPager.More() {
		_, err := routeMapsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RouteMaps_Get
	fmt.Println("Call operation: RouteMaps_Get")
	_, err = routeMapsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, testsuite.routeMapName, nil)
	testsuite.Require().NoError(err)

	// From step RouteMaps_Delete
	fmt.Println("Call operation: RouteMaps_Delete")
	routeMapsClientDeleteResponsePoller, err := routeMapsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, testsuite.routeMapName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routeMapsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/virtualHubs/{virtualHubName}/hubRouteTables/{routeTableName}
func (testsuite *VirtualWanTestSuite) TestHubRouteTables() {
	var err error
	// From step HubRouteTables_CreateOrUpdate
	fmt.Println("Call operation: HubRouteTables_CreateOrUpdate")
	hubRouteTablesClient, err := armnetwork.NewHubRouteTablesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	hubRouteTablesClientCreateOrUpdateResponsePoller, err := hubRouteTablesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, "hubRouteTable1", armnetwork.HubRouteTable{
		Properties: &armnetwork.HubRouteTableProperties{
			Labels: []*string{
				to.Ptr("label1"),
				to.Ptr("label2"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, hubRouteTablesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step HubRouteTables_List
	fmt.Println("Call operation: HubRouteTables_List")
	hubRouteTablesClientNewListPager := hubRouteTablesClient.NewListPager(testsuite.resourceGroupName, testsuite.virtualHubName, nil)
	for hubRouteTablesClientNewListPager.More() {
		_, err := hubRouteTablesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step HubRouteTables_Get
	fmt.Println("Call operation: HubRouteTables_Get")
	_, err = hubRouteTablesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, "hubRouteTable1", nil)
	testsuite.Require().NoError(err)

	// From step HubRouteTables_Delete
	fmt.Println("Call operation: HubRouteTables_Delete")
	hubRouteTablesClientDeleteResponsePoller, err := hubRouteTablesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, "hubRouteTable1", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, hubRouteTablesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Network/virtualHubs/{virtualHubName}/routingIntent/{routingIntentName}
func (testsuite *VirtualWanTestSuite) TestRoutingIntent() {
	var err error
	// From step RoutingIntent_CreateOrUpdate
	fmt.Println("Call operation: RoutingIntent_CreateOrUpdate")
	routingIntentClient, err := armnetwork.NewRoutingIntentClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	routingIntentClientCreateOrUpdateResponsePoller, err := routingIntentClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, "Intent1", armnetwork.RoutingIntent{
		Properties: &armnetwork.RoutingIntentProperties{},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routingIntentClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step RoutingIntent_List
	fmt.Println("Call operation: RoutingIntent_List")
	routingIntentClientNewListPager := routingIntentClient.NewListPager(testsuite.resourceGroupName, testsuite.virtualHubName, nil)
	for routingIntentClientNewListPager.More() {
		_, err := routingIntentClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RoutingIntent_Get
	fmt.Println("Call operation: RoutingIntent_Get")
	_, err = routingIntentClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, "Intent1", nil)
	testsuite.Require().NoError(err)

	// From step RoutingIntent_Delete
	fmt.Println("Call operation: RoutingIntent_Delete")
	routingIntentClientDeleteResponsePoller, err := routingIntentClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.virtualHubName, "Intent1", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, routingIntentClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *VirtualWanTestSuite) Cleanup() {
	var err error
	// From step VpnSites_Delete
	fmt.Println("Call operation: VPNSites_Delete")
	vPNSitesClient, err := armnetwork.NewVPNSitesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vPNSitesClientDeleteResponsePoller, err := vPNSitesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.vpnSiteName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, vPNSitesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step VirtualWans_Delete
	fmt.Println("Call operation: VirtualWans_Delete")
	virtualWansClient, err := armnetwork.NewVirtualWansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	virtualWansClientDeleteResponsePoller, err := virtualWansClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, "virtualWan1", nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, virtualWansClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
