//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VpnConnectionGet.json
func ExampleVPNConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVPNConnectionsClient().Get(ctx, "rg1", "gateway1", "vpnConnection1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VPNConnection = armnetwork.VPNConnection{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/vpnConnection1"),
	// 	Name: to.Ptr("vpnConnection1"),
	// 	Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 	Properties: &armnetwork.VPNConnectionProperties{
	// 		EgressBytesTransferred: to.Ptr[int64](0),
	// 		EnableInternetSecurity: to.Ptr(false),
	// 		IngressBytesTransferred: to.Ptr[int64](0),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		RemoteVPNSite: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnSites/vpnSite1"),
	// 		},
	// 		RoutingConfiguration: &armnetwork.RoutingConfiguration{
	// 			AssociatedRouteTable: &armnetwork.SubResource{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable1"),
	// 			},
	// 			InboundRouteMap: &armnetwork.SubResource{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1"),
	// 			},
	// 			OutboundRouteMap: &armnetwork.SubResource{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap2"),
	// 			},
	// 			PropagatedRouteTables: &armnetwork.PropagatedRouteTable{
	// 				IDs: []*armnetwork.SubResource{
	// 					{
	// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable1"),
	// 					},
	// 					{
	// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable2"),
	// 					},
	// 					{
	// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable3"),
	// 				}},
	// 				Labels: []*string{
	// 					to.Ptr("label1"),
	// 					to.Ptr("label2")},
	// 				},
	// 			},
	// 			TrafficSelectorPolicies: []*armnetwork.TrafficSelectorPolicy{
	// 			},
	// 			VPNLinkConnections: []*armnetwork.VPNSiteLinkConnection{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/vpnConnection1/VpnSiteLinkConnections/Connection-Link1"),
	// 					Name: to.Ptr("Connection-Link1"),
	// 					Type: to.Ptr("Microsoft.Network/vpnGateways/vpnConnections/VpnSiteLinkConnections"),
	// 					Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 					Properties: &armnetwork.VPNSiteLinkConnectionProperties{
	// 						ConnectionBandwidth: to.Ptr[int32](200),
	// 						EgressBytesTransferred: to.Ptr[int64](0),
	// 						EnableBgp: to.Ptr(false),
	// 						EnableRateLimiting: to.Ptr(false),
	// 						IngressBytesTransferred: to.Ptr[int64](0),
	// 						IPSecPolicies: []*armnetwork.IPSecPolicy{
	// 						},
	// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 						RoutingWeight: to.Ptr[int32](0),
	// 						SharedKey: to.Ptr("key"),
	// 						UseLocalAzureIPAddress: to.Ptr(false),
	// 						UsePolicyBasedTrafficSelectors: to.Ptr(false),
	// 						VPNConnectionProtocolType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
	// 						VPNLinkConnectionMode: to.Ptr(armnetwork.VPNLinkConnectionModeResponderOnly),
	// 						VPNSiteLink: &armnetwork.SubResource{
	// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnSites/vpnSite1/vpnSiteLinks/siteLink1"),
	// 						},
	// 					},
	// 				},
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/vpnConnection1/VpnSiteLinkConnections/Connection-Link2"),
	// 					Name: to.Ptr("Connection-Link2"),
	// 					Type: to.Ptr("Microsoft.Network/vpnGateways/vpnConnections/VpnSiteLinkConnections"),
	// 					Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 					Properties: &armnetwork.VPNSiteLinkConnectionProperties{
	// 						ConnectionBandwidth: to.Ptr[int32](200),
	// 						EgressBytesTransferred: to.Ptr[int64](0),
	// 						EnableBgp: to.Ptr(false),
	// 						EnableRateLimiting: to.Ptr(false),
	// 						IngressBytesTransferred: to.Ptr[int64](0),
	// 						IPSecPolicies: []*armnetwork.IPSecPolicy{
	// 						},
	// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 						RoutingWeight: to.Ptr[int32](0),
	// 						SharedKey: to.Ptr("key"),
	// 						UseLocalAzureIPAddress: to.Ptr(false),
	// 						UsePolicyBasedTrafficSelectors: to.Ptr(false),
	// 						VPNConnectionProtocolType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
	// 						VPNLinkConnectionMode: to.Ptr(armnetwork.VPNLinkConnectionModeInitiatorOnly),
	// 						VPNSiteLink: &armnetwork.SubResource{
	// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnSites/vpnSite1/vpnSiteLinks/siteLink2"),
	// 						},
	// 					},
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VpnConnectionPut.json
func ExampleVPNConnectionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVPNConnectionsClient().BeginCreateOrUpdate(ctx, "rg1", "gateway1", "vpnConnection1", armnetwork.VPNConnection{
		Properties: &armnetwork.VPNConnectionProperties{
			RemoteVPNSite: &armnetwork.SubResource{
				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnSites/vpnSite1"),
			},
			RoutingConfiguration: &armnetwork.RoutingConfiguration{
				AssociatedRouteTable: &armnetwork.SubResource{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable1"),
				},
				InboundRouteMap: &armnetwork.SubResource{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1"),
				},
				OutboundRouteMap: &armnetwork.SubResource{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap2"),
				},
				PropagatedRouteTables: &armnetwork.PropagatedRouteTable{
					IDs: []*armnetwork.SubResource{
						{
							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable1"),
						},
						{
							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable2"),
						},
						{
							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable3"),
						}},
					Labels: []*string{
						to.Ptr("label1"),
						to.Ptr("label2")},
				},
			},
			TrafficSelectorPolicies: []*armnetwork.TrafficSelectorPolicy{},
			VPNLinkConnections: []*armnetwork.VPNSiteLinkConnection{
				{
					Name: to.Ptr("Connection-Link1"),
					Properties: &armnetwork.VPNSiteLinkConnectionProperties{
						ConnectionBandwidth:            to.Ptr[int32](200),
						SharedKey:                      to.Ptr("key"),
						UsePolicyBasedTrafficSelectors: to.Ptr(false),
						VPNConnectionProtocolType:      to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
						VPNLinkConnectionMode:          to.Ptr(armnetwork.VPNLinkConnectionModeDefault),
						VPNSiteLink: &armnetwork.SubResource{
							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnSites/vpnSite1/vpnSiteLinks/siteLink1"),
						},
					},
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VPNConnection = armnetwork.VPNConnection{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/vpnConnection1"),
	// 	Name: to.Ptr("vpnConnection1"),
	// 	Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 	Properties: &armnetwork.VPNConnectionProperties{
	// 		EgressBytesTransferred: to.Ptr[int64](0),
	// 		EnableInternetSecurity: to.Ptr(false),
	// 		IngressBytesTransferred: to.Ptr[int64](0),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		RemoteVPNSite: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnSites/vpnSite1"),
	// 		},
	// 		RoutingConfiguration: &armnetwork.RoutingConfiguration{
	// 			AssociatedRouteTable: &armnetwork.SubResource{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable1"),
	// 			},
	// 			InboundRouteMap: &armnetwork.SubResource{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1"),
	// 			},
	// 			OutboundRouteMap: &armnetwork.SubResource{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap2"),
	// 			},
	// 			PropagatedRouteTables: &armnetwork.PropagatedRouteTable{
	// 				IDs: []*armnetwork.SubResource{
	// 					{
	// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable1"),
	// 					},
	// 					{
	// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable2"),
	// 					},
	// 					{
	// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubRouteTables/hubRouteTable3"),
	// 				}},
	// 				Labels: []*string{
	// 					to.Ptr("label1"),
	// 					to.Ptr("label2")},
	// 				},
	// 			},
	// 			TrafficSelectorPolicies: []*armnetwork.TrafficSelectorPolicy{
	// 			},
	// 			VPNLinkConnections: []*armnetwork.VPNSiteLinkConnection{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnGateways/gateway1/vpnConnections/vpnConnection1/VpnSiteLinkConnections/Connection-Link1"),
	// 					Name: to.Ptr("Connection-Link1"),
	// 					Type: to.Ptr("Microsoft.Network/vpnGateways/vpnConnections/VpnSiteLinkConnections"),
	// 					Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 					Properties: &armnetwork.VPNSiteLinkConnectionProperties{
	// 						ConnectionBandwidth: to.Ptr[int32](200),
	// 						EgressBytesTransferred: to.Ptr[int64](0),
	// 						EnableBgp: to.Ptr(false),
	// 						EnableRateLimiting: to.Ptr(false),
	// 						IngressBytesTransferred: to.Ptr[int64](0),
	// 						IPSecPolicies: []*armnetwork.IPSecPolicy{
	// 						},
	// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 						RoutingWeight: to.Ptr[int32](0),
	// 						SharedKey: to.Ptr("key"),
	// 						UseLocalAzureIPAddress: to.Ptr(false),
	// 						UsePolicyBasedTrafficSelectors: to.Ptr(false),
	// 						VPNConnectionProtocolType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
	// 						VPNLinkConnectionMode: to.Ptr(armnetwork.VPNLinkConnectionModeDefault),
	// 						VPNSiteLink: &armnetwork.SubResource{
	// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/vpnSites/vpnSite1/vpnSiteLinks/siteLink1"),
	// 						},
	// 					},
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VpnConnectionDelete.json
func ExampleVPNConnectionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVPNConnectionsClient().BeginDelete(ctx, "rg1", "gateway1", "vpnConnection1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VpnConnectionStartPacketCaptureFilterData.json
func ExampleVPNConnectionsClient_BeginStartPacketCapture_startPacketCaptureOnVpnConnectionWithFilter() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVPNConnectionsClient().BeginStartPacketCapture(ctx, "rg1", "gateway1", "vpnConnection1", &armnetwork.VPNConnectionsClientBeginStartPacketCaptureOptions{Parameters: &armnetwork.VPNConnectionPacketCaptureStartParameters{
		FilterData: to.Ptr("{'TracingFlags': 11,'MaxPacketBufferSize': 120,'MaxFileSize': 200,'Filters': [{'SourceSubnets': ['20.1.1.0/24'],'DestinationSubnets': ['10.1.1.0/24'],'SourcePort': [500],'DestinationPort': [4500],'Protocol': 6,'TcpFlags': 16,'CaptureSingleDirectionTrafficOnly': true}]}"),
		LinkConnectionNames: []*string{
			to.Ptr("siteLink1"),
			to.Ptr("siteLink2")},
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Value = "\"{\"Status\":\"Successful\",\"Data\":null}\""
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VpnConnectionStartPacketCapture.json
func ExampleVPNConnectionsClient_BeginStartPacketCapture_startPacketCaptureOnVpnConnectionWithoutFilter() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVPNConnectionsClient().BeginStartPacketCapture(ctx, "rg1", "gateway1", "vpnConnection1", &armnetwork.VPNConnectionsClientBeginStartPacketCaptureOptions{Parameters: &armnetwork.VPNConnectionPacketCaptureStartParameters{
		LinkConnectionNames: []*string{
			to.Ptr("siteLink1"),
			to.Ptr("siteLink2")},
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Value = "\"{\"Status\":\"Successful\",\"Data\":null}\""
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VpnConnectionStopPacketCapture.json
func ExampleVPNConnectionsClient_BeginStopPacketCapture() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVPNConnectionsClient().BeginStopPacketCapture(ctx, "rg1", "gateway1", "vpnConnection1", &armnetwork.VPNConnectionsClientBeginStopPacketCaptureOptions{Parameters: &armnetwork.VPNConnectionPacketCaptureStopParameters{
		LinkConnectionNames: []*string{
			to.Ptr("vpnSiteLink1"),
			to.Ptr("vpnSiteLink2")},
		SasURL: to.Ptr("https://teststorage.blob.core.windows.net/?sv=2018-03-28&ss=bfqt&srt=sco&sp=rwdlacup&se=2019-09-13T07:44:05Z&st=2019-09-06T23:44:05Z&spr=https&sig=V1h9D1riltvZMI69d6ihENnFo%2FrCvTqGgjO2lf%2FVBhE%3D"),
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Value = "\"{\"Status\":\"Successful\",\"Data\":null}\""
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VpnConnectionList.json
func ExampleVPNConnectionsClient_NewListByVPNGatewayPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVPNConnectionsClient().NewListByVPNGatewayPager("rg1", "gateway1", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.ListVPNConnectionsResult = armnetwork.ListVPNConnectionsResult{
		// }
	}
}
