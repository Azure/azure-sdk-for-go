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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/DscpConfigurationCreate.json
func ExampleDscpConfigurationClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDscpConfigurationClient().BeginCreateOrUpdate(ctx, "rg1", "mydscpconfig", armnetwork.DscpConfiguration{
		Location: to.Ptr("eastus"),
		Properties: &armnetwork.DscpConfigurationPropertiesFormat{
			QosDefinitionCollection: []*armnetwork.QosDefinition{
				{
					DestinationIPRanges: []*armnetwork.QosIPRange{
						{
							EndIP:   to.Ptr("127.0.10.2"),
							StartIP: to.Ptr("127.0.10.1"),
						}},
					DestinationPortRanges: []*armnetwork.QosPortRange{
						{
							End:   to.Ptr[int32](15),
							Start: to.Ptr[int32](15),
						}},
					Markings: []*int32{
						to.Ptr[int32](1)},
					SourceIPRanges: []*armnetwork.QosIPRange{
						{
							EndIP:   to.Ptr("127.0.0.2"),
							StartIP: to.Ptr("127.0.0.1"),
						}},
					SourcePortRanges: []*armnetwork.QosPortRange{
						{
							End:   to.Ptr[int32](11),
							Start: to.Ptr[int32](10),
						},
						{
							End:   to.Ptr[int32](21),
							Start: to.Ptr[int32](20),
						}},
					Protocol: to.Ptr(armnetwork.ProtocolTypeTCP),
				},
				{
					DestinationIPRanges: []*armnetwork.QosIPRange{
						{
							EndIP:   to.Ptr("12.0.10.2"),
							StartIP: to.Ptr("12.0.10.1"),
						}},
					DestinationPortRanges: []*armnetwork.QosPortRange{
						{
							End:   to.Ptr[int32](52),
							Start: to.Ptr[int32](51),
						}},
					Markings: []*int32{
						to.Ptr[int32](2)},
					SourceIPRanges: []*armnetwork.QosIPRange{
						{
							EndIP:   to.Ptr("12.0.0.2"),
							StartIP: to.Ptr("12.0.0.1"),
						}},
					SourcePortRanges: []*armnetwork.QosPortRange{
						{
							End:   to.Ptr[int32](12),
							Start: to.Ptr[int32](11),
						}},
					Protocol: to.Ptr(armnetwork.ProtocolTypeUDP),
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
	// res.DscpConfiguration = armnetwork.DscpConfiguration{
	// 	Name: to.Ptr("mydscpConfig"),
	// 	Type: to.Ptr("Microsoft.Network/dscpConfiguration"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/dscpConfiguration/mydscpConfig"),
	// 	Location: to.Ptr("eastus"),
	// 	Properties: &armnetwork.DscpConfigurationPropertiesFormat{
	// 		AssociatedNetworkInterfaces: []*armnetwork.Interface{
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		QosCollectionID: to.Ptr("0f8fad5b-d9cb-469f-a165-70867728950e"),
	// 		QosDefinitionCollection: []*armnetwork.QosDefinition{
	// 			{
	// 				DestinationIPRanges: []*armnetwork.QosIPRange{
	// 					{
	// 						EndIP: to.Ptr("127.0.10.2"),
	// 						StartIP: to.Ptr("127.0.10.1"),
	// 				}},
	// 				DestinationPortRanges: []*armnetwork.QosPortRange{
	// 					{
	// 						End: to.Ptr[int32](62),
	// 						Start: to.Ptr[int32](61),
	// 				}},
	// 				Markings: []*int32{
	// 					to.Ptr[int32](1)},
	// 					SourceIPRanges: []*armnetwork.QosIPRange{
	// 						{
	// 							EndIP: to.Ptr("127.0.0.2"),
	// 							StartIP: to.Ptr("127.0.0.1"),
	// 					}},
	// 					SourcePortRanges: []*armnetwork.QosPortRange{
	// 						{
	// 							End: to.Ptr[int32](12),
	// 							Start: to.Ptr[int32](11),
	// 					}},
	// 					Protocol: to.Ptr(armnetwork.ProtocolTypeTCP),
	// 				},
	// 				{
	// 					DestinationIPRanges: []*armnetwork.QosIPRange{
	// 						{
	// 							EndIP: to.Ptr("12.0.10.2"),
	// 							StartIP: to.Ptr("12.0.10.1"),
	// 					}},
	// 					DestinationPortRanges: []*armnetwork.QosPortRange{
	// 						{
	// 							End: to.Ptr[int32](52),
	// 							Start: to.Ptr[int32](51),
	// 					}},
	// 					Markings: []*int32{
	// 						to.Ptr[int32](2)},
	// 						SourceIPRanges: []*armnetwork.QosIPRange{
	// 							{
	// 								EndIP: to.Ptr("12.0.0.2"),
	// 								StartIP: to.Ptr("12.0.0.1"),
	// 						}},
	// 						SourcePortRanges: []*armnetwork.QosPortRange{
	// 							{
	// 								End: to.Ptr[int32](12),
	// 								Start: to.Ptr[int32](11),
	// 						}},
	// 						Protocol: to.Ptr(armnetwork.ProtocolTypeUDP),
	// 				}},
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/DscpConfigurationDelete.json
func ExampleDscpConfigurationClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDscpConfigurationClient().BeginDelete(ctx, "rg1", "mydscpConfig", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/DscpConfigurationGet.json
func ExampleDscpConfigurationClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDscpConfigurationClient().Get(ctx, "rg1", "mydscpConfig", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DscpConfiguration = armnetwork.DscpConfiguration{
	// 	Name: to.Ptr("mydscpConfig"),
	// 	Type: to.Ptr("Microsoft.Network/dscpConfiguration"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/dscpConfiguration/mydscpConfig"),
	// 	Location: to.Ptr("eastus"),
	// 	Properties: &armnetwork.DscpConfigurationPropertiesFormat{
	// 		AssociatedNetworkInterfaces: []*armnetwork.Interface{
	// 			{
	// 			},
	// 			{
	// 		}},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		QosCollectionID: to.Ptr("0f8fad5b-d9cb-469f-a165-70867728950e"),
	// 		QosDefinitionCollection: []*armnetwork.QosDefinition{
	// 			{
	// 				DestinationIPRanges: []*armnetwork.QosIPRange{
	// 					{
	// 						EndIP: to.Ptr("127.0.10.2"),
	// 						StartIP: to.Ptr("127.0.10.1"),
	// 				}},
	// 				DestinationPortRanges: []*armnetwork.QosPortRange{
	// 					{
	// 						End: to.Ptr[int32](62),
	// 						Start: to.Ptr[int32](61),
	// 				}},
	// 				Markings: []*int32{
	// 					to.Ptr[int32](1)},
	// 					SourceIPRanges: []*armnetwork.QosIPRange{
	// 						{
	// 							EndIP: to.Ptr("127.0.0.2"),
	// 							StartIP: to.Ptr("127.0.0.1"),
	// 					}},
	// 					SourcePortRanges: []*armnetwork.QosPortRange{
	// 						{
	// 							End: to.Ptr[int32](12),
	// 							Start: to.Ptr[int32](11),
	// 					}},
	// 					Protocol: to.Ptr(armnetwork.ProtocolTypeTCP),
	// 				},
	// 				{
	// 					DestinationIPRanges: []*armnetwork.QosIPRange{
	// 						{
	// 							EndIP: to.Ptr("12.0.10.2"),
	// 							StartIP: to.Ptr("12.0.10.1"),
	// 					}},
	// 					DestinationPortRanges: []*armnetwork.QosPortRange{
	// 						{
	// 							End: to.Ptr[int32](52),
	// 							Start: to.Ptr[int32](51),
	// 					}},
	// 					Markings: []*int32{
	// 						to.Ptr[int32](2)},
	// 						SourceIPRanges: []*armnetwork.QosIPRange{
	// 							{
	// 								EndIP: to.Ptr("12.0.0.2"),
	// 								StartIP: to.Ptr("12.0.0.1"),
	// 						}},
	// 						SourcePortRanges: []*armnetwork.QosPortRange{
	// 							{
	// 								End: to.Ptr[int32](12),
	// 								Start: to.Ptr[int32](11),
	// 						}},
	// 						Protocol: to.Ptr(armnetwork.ProtocolTypeUDP),
	// 				}},
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/DscpConfigurationList.json
func ExampleDscpConfigurationClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDscpConfigurationClient().NewListPager("rg1", nil)
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
		// page.DscpConfigurationListResult = armnetwork.DscpConfigurationListResult{
		// 	Value: []*armnetwork.DscpConfiguration{
		// 		{
		// 			Name: to.Ptr("mydscpConfig"),
		// 			Type: to.Ptr("Microsoft.Network/dscpConfiguration"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/dscpConfiguration/mydscpConfig"),
		// 			Location: to.Ptr("eastus"),
		// 			Properties: &armnetwork.DscpConfigurationPropertiesFormat{
		// 				AssociatedNetworkInterfaces: []*armnetwork.Interface{
		// 					{
		// 					},
		// 					{
		// 				}},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				QosCollectionID: to.Ptr("0f8fad5b-d9cb-469f-a165-70867728950e"),
		// 				QosDefinitionCollection: []*armnetwork.QosDefinition{
		// 					{
		// 						DestinationIPRanges: []*armnetwork.QosIPRange{
		// 							{
		// 								EndIP: to.Ptr("127.0.10.2"),
		// 								StartIP: to.Ptr("127.0.10.1"),
		// 						}},
		// 						DestinationPortRanges: []*armnetwork.QosPortRange{
		// 							{
		// 								End: to.Ptr[int32](62),
		// 								Start: to.Ptr[int32](61),
		// 						}},
		// 						Markings: []*int32{
		// 							to.Ptr[int32](1)},
		// 							SourceIPRanges: []*armnetwork.QosIPRange{
		// 								{
		// 									EndIP: to.Ptr("127.0.0.2"),
		// 									StartIP: to.Ptr("127.0.0.1"),
		// 							}},
		// 							SourcePortRanges: []*armnetwork.QosPortRange{
		// 								{
		// 									End: to.Ptr[int32](12),
		// 									Start: to.Ptr[int32](11),
		// 							}},
		// 							Protocol: to.Ptr(armnetwork.ProtocolTypeTCP),
		// 						},
		// 						{
		// 							DestinationIPRanges: []*armnetwork.QosIPRange{
		// 								{
		// 									EndIP: to.Ptr("12.0.10.2"),
		// 									StartIP: to.Ptr("12.0.10.1"),
		// 							}},
		// 							DestinationPortRanges: []*armnetwork.QosPortRange{
		// 								{
		// 									End: to.Ptr[int32](52),
		// 									Start: to.Ptr[int32](51),
		// 							}},
		// 							Markings: []*int32{
		// 								to.Ptr[int32](2)},
		// 								SourceIPRanges: []*armnetwork.QosIPRange{
		// 									{
		// 										EndIP: to.Ptr("12.0.0.2"),
		// 										StartIP: to.Ptr("12.0.0.1"),
		// 								}},
		// 								SourcePortRanges: []*armnetwork.QosPortRange{
		// 									{
		// 										End: to.Ptr[int32](12),
		// 										Start: to.Ptr[int32](11),
		// 								}},
		// 								Protocol: to.Ptr(armnetwork.ProtocolTypeUDP),
		// 						}},
		// 					},
		// 				},
		// 				{
		// 					Name: to.Ptr("mydscpConfig2"),
		// 					Type: to.Ptr("Microsoft.Network/dscpConfiguration"),
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/dscpConfiguration/mydscpConfig2"),
		// 					Location: to.Ptr("eastus"),
		// 					Properties: &armnetwork.DscpConfigurationPropertiesFormat{
		// 						AssociatedNetworkInterfaces: []*armnetwork.Interface{
		// 							{
		// 							},
		// 							{
		// 						}},
		// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 						QosCollectionID: to.Ptr("9as24mf6-d9cb-7a7f-a165-70867728950e"),
		// 						QosDefinitionCollection: []*armnetwork.QosDefinition{
		// 							{
		// 								DestinationIPRanges: []*armnetwork.QosIPRange{
		// 									{
		// 										EndIP: to.Ptr("127.0.10.2"),
		// 										StartIP: to.Ptr("127.0.10.1"),
		// 								}},
		// 								DestinationPortRanges: []*armnetwork.QosPortRange{
		// 									{
		// 										End: to.Ptr[int32](62),
		// 										Start: to.Ptr[int32](61),
		// 								}},
		// 								Markings: []*int32{
		// 									to.Ptr[int32](1)},
		// 									SourceIPRanges: []*armnetwork.QosIPRange{
		// 										{
		// 											EndIP: to.Ptr("127.0.0.2"),
		// 											StartIP: to.Ptr("127.0.0.1"),
		// 									}},
		// 									SourcePortRanges: []*armnetwork.QosPortRange{
		// 										{
		// 											End: to.Ptr[int32](12),
		// 											Start: to.Ptr[int32](11),
		// 									}},
		// 									Protocol: to.Ptr(armnetwork.ProtocolTypeTCP),
		// 								},
		// 								{
		// 									DestinationIPRanges: []*armnetwork.QosIPRange{
		// 										{
		// 											EndIP: to.Ptr("12.0.10.2"),
		// 											StartIP: to.Ptr("12.0.10.1"),
		// 									}},
		// 									DestinationPortRanges: []*armnetwork.QosPortRange{
		// 										{
		// 											End: to.Ptr[int32](52),
		// 											Start: to.Ptr[int32](51),
		// 									}},
		// 									Markings: []*int32{
		// 										to.Ptr[int32](2)},
		// 										SourceIPRanges: []*armnetwork.QosIPRange{
		// 											{
		// 												EndIP: to.Ptr("12.0.0.2"),
		// 												StartIP: to.Ptr("12.0.0.1"),
		// 										}},
		// 										SourcePortRanges: []*armnetwork.QosPortRange{
		// 											{
		// 												End: to.Ptr[int32](12),
		// 												Start: to.Ptr[int32](11),
		// 										}},
		// 										Protocol: to.Ptr(armnetwork.ProtocolTypeUDP),
		// 								}},
		// 							},
		// 					}},
		// 				}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/DscpConfigurationListAll.json
func ExampleDscpConfigurationClient_NewListAllPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDscpConfigurationClient().NewListAllPager(nil)
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
		// page.DscpConfigurationListResult = armnetwork.DscpConfigurationListResult{
		// 	Value: []*armnetwork.DscpConfiguration{
		// 		{
		// 			Name: to.Ptr("mydscpConfig"),
		// 			Type: to.Ptr("Microsoft.Network/dscpConfiguration"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/dscpConfiguration/mydscpConfig"),
		// 			Location: to.Ptr("eastus"),
		// 			Properties: &armnetwork.DscpConfigurationPropertiesFormat{
		// 				AssociatedNetworkInterfaces: []*armnetwork.Interface{
		// 					{
		// 					},
		// 					{
		// 				}},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				QosCollectionID: to.Ptr("0f8fad5b-d9cb-469f-a165-70867728950e"),
		// 				QosDefinitionCollection: []*armnetwork.QosDefinition{
		// 					{
		// 						DestinationIPRanges: []*armnetwork.QosIPRange{
		// 							{
		// 								EndIP: to.Ptr("127.0.10.2"),
		// 								StartIP: to.Ptr("127.0.10.1"),
		// 						}},
		// 						DestinationPortRanges: []*armnetwork.QosPortRange{
		// 							{
		// 								End: to.Ptr[int32](62),
		// 								Start: to.Ptr[int32](61),
		// 						}},
		// 						Markings: []*int32{
		// 							to.Ptr[int32](1)},
		// 							SourceIPRanges: []*armnetwork.QosIPRange{
		// 								{
		// 									EndIP: to.Ptr("127.0.0.2"),
		// 									StartIP: to.Ptr("127.0.0.1"),
		// 							}},
		// 							SourcePortRanges: []*armnetwork.QosPortRange{
		// 								{
		// 									End: to.Ptr[int32](12),
		// 									Start: to.Ptr[int32](11),
		// 							}},
		// 							Protocol: to.Ptr(armnetwork.ProtocolTypeTCP),
		// 						},
		// 						{
		// 							DestinationIPRanges: []*armnetwork.QosIPRange{
		// 								{
		// 									EndIP: to.Ptr("12.0.10.2"),
		// 									StartIP: to.Ptr("12.0.10.1"),
		// 							}},
		// 							DestinationPortRanges: []*armnetwork.QosPortRange{
		// 								{
		// 									End: to.Ptr[int32](52),
		// 									Start: to.Ptr[int32](51),
		// 							}},
		// 							Markings: []*int32{
		// 								to.Ptr[int32](2)},
		// 								SourceIPRanges: []*armnetwork.QosIPRange{
		// 									{
		// 										EndIP: to.Ptr("12.0.0.2"),
		// 										StartIP: to.Ptr("12.0.0.1"),
		// 								}},
		// 								SourcePortRanges: []*armnetwork.QosPortRange{
		// 									{
		// 										End: to.Ptr[int32](12),
		// 										Start: to.Ptr[int32](11),
		// 								}},
		// 								Protocol: to.Ptr(armnetwork.ProtocolTypeUDP),
		// 						}},
		// 					},
		// 				},
		// 				{
		// 					Name: to.Ptr("mydscpConfig2"),
		// 					Type: to.Ptr("Microsoft.Network/dscpConfiguration"),
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg2/providers/Microsoft.Network/dscpConfiguration/mydscpConfig2"),
		// 					Location: to.Ptr("eastus"),
		// 					Properties: &armnetwork.DscpConfigurationPropertiesFormat{
		// 						AssociatedNetworkInterfaces: []*armnetwork.Interface{
		// 							{
		// 							},
		// 							{
		// 						}},
		// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 						QosCollectionID: to.Ptr("9as24mf6-d9cb-7a7f-a165-70867728950e"),
		// 						QosDefinitionCollection: []*armnetwork.QosDefinition{
		// 							{
		// 								DestinationIPRanges: []*armnetwork.QosIPRange{
		// 									{
		// 										EndIP: to.Ptr("127.0.10.2"),
		// 										StartIP: to.Ptr("127.0.10.1"),
		// 								}},
		// 								DestinationPortRanges: []*armnetwork.QosPortRange{
		// 									{
		// 										End: to.Ptr[int32](62),
		// 										Start: to.Ptr[int32](61),
		// 								}},
		// 								Markings: []*int32{
		// 									to.Ptr[int32](1)},
		// 									SourceIPRanges: []*armnetwork.QosIPRange{
		// 										{
		// 											EndIP: to.Ptr("127.0.0.2"),
		// 											StartIP: to.Ptr("127.0.0.1"),
		// 									}},
		// 									SourcePortRanges: []*armnetwork.QosPortRange{
		// 										{
		// 											End: to.Ptr[int32](12),
		// 											Start: to.Ptr[int32](11),
		// 									}},
		// 									Protocol: to.Ptr(armnetwork.ProtocolTypeTCP),
		// 								},
		// 								{
		// 									DestinationIPRanges: []*armnetwork.QosIPRange{
		// 										{
		// 											EndIP: to.Ptr("12.0.10.2"),
		// 											StartIP: to.Ptr("12.0.10.1"),
		// 									}},
		// 									DestinationPortRanges: []*armnetwork.QosPortRange{
		// 										{
		// 											End: to.Ptr[int32](52),
		// 											Start: to.Ptr[int32](51),
		// 									}},
		// 									Markings: []*int32{
		// 										to.Ptr[int32](2)},
		// 										SourceIPRanges: []*armnetwork.QosIPRange{
		// 											{
		// 												EndIP: to.Ptr("12.0.0.2"),
		// 												StartIP: to.Ptr("12.0.0.1"),
		// 										}},
		// 										SourcePortRanges: []*armnetwork.QosPortRange{
		// 											{
		// 												End: to.Ptr[int32](12),
		// 												Start: to.Ptr[int32](11),
		// 										}},
		// 										Protocol: to.Ptr(armnetwork.ProtocolTypeUDP),
		// 								}},
		// 							},
		// 					}},
		// 				}
	}
}
