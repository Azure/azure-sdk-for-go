//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcompute_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/Skus/stable/2021-07-01/examples/skus/ListAvailableResourceSkus.json
func ExampleResourceSKUsClient_NewListPager_listsAllAvailableResourceSkUs() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewResourceSKUsClient().NewListPager(&armcompute.ResourceSKUsClientListOptions{Filter: nil,
		IncludeExtendedLocations: nil,
	})
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
		// page.ResourceSKUsResult = armcompute.ResourceSKUsResult{
		// 	Value: []*armcompute.ResourceSKU{
		// 		{
		// 			Name: to.Ptr("Standard_A0"),
		// 			Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 				{
		// 					Name: to.Ptr("MaxResourceVolumeMB"),
		// 					Value: to.Ptr("20480"),
		// 				},
		// 				{
		// 					Name: to.Ptr("OSVhdSizeMB"),
		// 					Value: to.Ptr("1047552"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUs"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("HyperVGenerations"),
		// 					Value: to.Ptr("V1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MemoryGB"),
		// 					Value: to.Ptr("0.75"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MaxDataDiskCount"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("LowPriorityCapable"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("PremiumIO"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUsAvailable"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("ACUs"),
		// 					Value: to.Ptr("50"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUsPerCore"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("EphemeralOSDiskSupported"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("AcceleratedNetworkingEnabled"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("RdmaEnabled"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MaxNetworkInterfaces"),
		// 					Value: to.Ptr("2"),
		// 			}},
		// 			Family: to.Ptr("standardA0_A7Family"),
		// 			LocationInfo: []*armcompute.ResourceSKULocationInfo{
		// 				{
		// 					Location: to.Ptr("westus"),
		// 					ZoneDetails: []*armcompute.ResourceSKUZoneDetails{
		// 						{
		// 							Name: []*string{
		// 								to.Ptr("2")},
		// 								Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 									{
		// 										Name: to.Ptr("UltraSSDAvailable"),
		// 										Value: to.Ptr("True"),
		// 								}},
		// 						}},
		// 						Zones: []*string{
		// 							to.Ptr("2"),
		// 							to.Ptr("1")},
		// 					}},
		// 					Locations: []*string{
		// 						to.Ptr("westus")},
		// 						ResourceType: to.Ptr("virtualMachines"),
		// 						Size: to.Ptr("A0"),
		// 						Tier: to.Ptr("Standard"),
		// 					},
		// 					{
		// 						Name: to.Ptr("Standard_A1"),
		// 						Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 							{
		// 								Name: to.Ptr("MaxResourceVolumeMB"),
		// 								Value: to.Ptr("71680"),
		// 							},
		// 							{
		// 								Name: to.Ptr("OSVhdSizeMB"),
		// 								Value: to.Ptr("1047552"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUs"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("HyperVGenerations"),
		// 								Value: to.Ptr("V1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MemoryGB"),
		// 								Value: to.Ptr("1.75"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MaxDataDiskCount"),
		// 								Value: to.Ptr("2"),
		// 							},
		// 							{
		// 								Name: to.Ptr("LowPriorityCapable"),
		// 								Value: to.Ptr("True"),
		// 							},
		// 							{
		// 								Name: to.Ptr("PremiumIO"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUsAvailable"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("ACUs"),
		// 								Value: to.Ptr("100"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUsPerCore"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("EphemeralOSDiskSupported"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("AcceleratedNetworkingEnabled"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("RdmaEnabled"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MaxNetworkInterfaces"),
		// 								Value: to.Ptr("2"),
		// 						}},
		// 						Family: to.Ptr("standardA0_A7Family"),
		// 						LocationInfo: []*armcompute.ResourceSKULocationInfo{
		// 							{
		// 								Location: to.Ptr("westus"),
		// 								Zones: []*string{
		// 									to.Ptr("1"),
		// 									to.Ptr("2"),
		// 									to.Ptr("3")},
		// 							}},
		// 							Locations: []*string{
		// 								to.Ptr("westus")},
		// 								ResourceType: to.Ptr("virtualMachines"),
		// 								Size: to.Ptr("A1"),
		// 								Tier: to.Ptr("Standard"),
		// 						}},
		// 					}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/Skus/stable/2021-07-01/examples/skus/ListAvailableResourceSkusForARegion.json
func ExampleResourceSKUsClient_NewListPager_listsAllAvailableResourceSkUsForTheSpecifiedRegion() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewResourceSKUsClient().NewListPager(&armcompute.ResourceSKUsClientListOptions{Filter: to.Ptr("location eq 'westus'"),
		IncludeExtendedLocations: nil,
	})
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
		// page.ResourceSKUsResult = armcompute.ResourceSKUsResult{
		// 	Value: []*armcompute.ResourceSKU{
		// 		{
		// 			Name: to.Ptr("Standard_A0"),
		// 			Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 				{
		// 					Name: to.Ptr("MaxResourceVolumeMB"),
		// 					Value: to.Ptr("20480"),
		// 				},
		// 				{
		// 					Name: to.Ptr("OSVhdSizeMB"),
		// 					Value: to.Ptr("1047552"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUs"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("HyperVGenerations"),
		// 					Value: to.Ptr("V1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MemoryGB"),
		// 					Value: to.Ptr("0.75"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MaxDataDiskCount"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("LowPriorityCapable"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("PremiumIO"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUsAvailable"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("ACUs"),
		// 					Value: to.Ptr("50"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUsPerCore"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("EphemeralOSDiskSupported"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("AcceleratedNetworkingEnabled"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("RdmaEnabled"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MaxNetworkInterfaces"),
		// 					Value: to.Ptr("2"),
		// 			}},
		// 			Family: to.Ptr("standardA0_A7Family"),
		// 			LocationInfo: []*armcompute.ResourceSKULocationInfo{
		// 				{
		// 					Location: to.Ptr("westus"),
		// 					ZoneDetails: []*armcompute.ResourceSKUZoneDetails{
		// 						{
		// 							Name: []*string{
		// 								to.Ptr("2")},
		// 								Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 									{
		// 										Name: to.Ptr("UltraSSDAvailable"),
		// 										Value: to.Ptr("True"),
		// 								}},
		// 						}},
		// 						Zones: []*string{
		// 							to.Ptr("2"),
		// 							to.Ptr("1")},
		// 					}},
		// 					Locations: []*string{
		// 						to.Ptr("westus")},
		// 						ResourceType: to.Ptr("virtualMachines"),
		// 						Size: to.Ptr("A0"),
		// 						Tier: to.Ptr("Standard"),
		// 					},
		// 					{
		// 						Name: to.Ptr("Standard_A1"),
		// 						Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 							{
		// 								Name: to.Ptr("MaxResourceVolumeMB"),
		// 								Value: to.Ptr("71680"),
		// 							},
		// 							{
		// 								Name: to.Ptr("OSVhdSizeMB"),
		// 								Value: to.Ptr("1047552"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUs"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("HyperVGenerations"),
		// 								Value: to.Ptr("V1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MemoryGB"),
		// 								Value: to.Ptr("1.75"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MaxDataDiskCount"),
		// 								Value: to.Ptr("2"),
		// 							},
		// 							{
		// 								Name: to.Ptr("LowPriorityCapable"),
		// 								Value: to.Ptr("True"),
		// 							},
		// 							{
		// 								Name: to.Ptr("PremiumIO"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUsAvailable"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("ACUs"),
		// 								Value: to.Ptr("100"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUsPerCore"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("EphemeralOSDiskSupported"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("AcceleratedNetworkingEnabled"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("RdmaEnabled"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MaxNetworkInterfaces"),
		// 								Value: to.Ptr("2"),
		// 						}},
		// 						Family: to.Ptr("standardA0_A7Family"),
		// 						LocationInfo: []*armcompute.ResourceSKULocationInfo{
		// 							{
		// 								Location: to.Ptr("westus"),
		// 								Zones: []*string{
		// 									to.Ptr("1"),
		// 									to.Ptr("2"),
		// 									to.Ptr("3")},
		// 							}},
		// 							Locations: []*string{
		// 								to.Ptr("westus")},
		// 								ResourceType: to.Ptr("virtualMachines"),
		// 								Size: to.Ptr("A1"),
		// 								Tier: to.Ptr("Standard"),
		// 						}},
		// 					}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d6d0798c6f5eb196fba7bd1924db2b145a94f58c/specification/compute/resource-manager/Microsoft.Compute/Skus/stable/2021-07-01/examples/skus/ListAvailableResourceSkusWithExtendedLocations.json
func ExampleResourceSKUsClient_NewListPager_listsAllAvailableResourceSkUsWithExtendedLocationInformation() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewResourceSKUsClient().NewListPager(&armcompute.ResourceSKUsClientListOptions{Filter: nil,
		IncludeExtendedLocations: to.Ptr("true"),
	})
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
		// page.ResourceSKUsResult = armcompute.ResourceSKUsResult{
		// 	Value: []*armcompute.ResourceSKU{
		// 		{
		// 			Name: to.Ptr("Standard_A0"),
		// 			Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 				{
		// 					Name: to.Ptr("MaxResourceVolumeMB"),
		// 					Value: to.Ptr("20480"),
		// 				},
		// 				{
		// 					Name: to.Ptr("OSVhdSizeMB"),
		// 					Value: to.Ptr("1047552"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUs"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("HyperVGenerations"),
		// 					Value: to.Ptr("V1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MemoryGB"),
		// 					Value: to.Ptr("0.75"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MaxDataDiskCount"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("LowPriorityCapable"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("PremiumIO"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUsAvailable"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("ACUs"),
		// 					Value: to.Ptr("50"),
		// 				},
		// 				{
		// 					Name: to.Ptr("vCPUsPerCore"),
		// 					Value: to.Ptr("1"),
		// 				},
		// 				{
		// 					Name: to.Ptr("EphemeralOSDiskSupported"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("AcceleratedNetworkingEnabled"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("RdmaEnabled"),
		// 					Value: to.Ptr("False"),
		// 				},
		// 				{
		// 					Name: to.Ptr("MaxNetworkInterfaces"),
		// 					Value: to.Ptr("2"),
		// 			}},
		// 			Family: to.Ptr("standardA0_A7Family"),
		// 			LocationInfo: []*armcompute.ResourceSKULocationInfo{
		// 				{
		// 					Location: to.Ptr("westus"),
		// 					ZoneDetails: []*armcompute.ResourceSKUZoneDetails{
		// 						{
		// 							Name: []*string{
		// 								to.Ptr("2")},
		// 								Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 									{
		// 										Name: to.Ptr("UltraSSDAvailable"),
		// 										Value: to.Ptr("True"),
		// 								}},
		// 						}},
		// 						Zones: []*string{
		// 							to.Ptr("2"),
		// 							to.Ptr("1")},
		// 					}},
		// 					Locations: []*string{
		// 						to.Ptr("westus")},
		// 						ResourceType: to.Ptr("virtualMachines"),
		// 						Size: to.Ptr("A0"),
		// 						Tier: to.Ptr("Standard"),
		// 					},
		// 					{
		// 						Name: to.Ptr("Standard_A1"),
		// 						Capabilities: []*armcompute.ResourceSKUCapabilities{
		// 							{
		// 								Name: to.Ptr("MaxResourceVolumeMB"),
		// 								Value: to.Ptr("71680"),
		// 							},
		// 							{
		// 								Name: to.Ptr("OSVhdSizeMB"),
		// 								Value: to.Ptr("1047552"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUs"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("HyperVGenerations"),
		// 								Value: to.Ptr("V1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MemoryGB"),
		// 								Value: to.Ptr("1.75"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MaxDataDiskCount"),
		// 								Value: to.Ptr("2"),
		// 							},
		// 							{
		// 								Name: to.Ptr("LowPriorityCapable"),
		// 								Value: to.Ptr("True"),
		// 							},
		// 							{
		// 								Name: to.Ptr("PremiumIO"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUsAvailable"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("ACUs"),
		// 								Value: to.Ptr("100"),
		// 							},
		// 							{
		// 								Name: to.Ptr("vCPUsPerCore"),
		// 								Value: to.Ptr("1"),
		// 							},
		// 							{
		// 								Name: to.Ptr("EphemeralOSDiskSupported"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("AcceleratedNetworkingEnabled"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("RdmaEnabled"),
		// 								Value: to.Ptr("False"),
		// 							},
		// 							{
		// 								Name: to.Ptr("MaxNetworkInterfaces"),
		// 								Value: to.Ptr("2"),
		// 						}},
		// 						Family: to.Ptr("standardA0_A7Family"),
		// 						LocationInfo: []*armcompute.ResourceSKULocationInfo{
		// 							{
		// 								Location: to.Ptr("westus"),
		// 								Zones: []*string{
		// 									to.Ptr("1"),
		// 									to.Ptr("2"),
		// 									to.Ptr("3")},
		// 								},
		// 								{
		// 									Type: to.Ptr(armcompute.ExtendedLocationTypeEdgeZone),
		// 									ExtendedLocations: []*string{
		// 										to.Ptr("Las Vegas"),
		// 										to.Ptr("Seattle"),
		// 										to.Ptr("Portland")},
		// 										Location: to.Ptr("westus"),
		// 								}},
		// 								Locations: []*string{
		// 									to.Ptr("westus")},
		// 									ResourceType: to.Ptr("virtualMachines"),
		// 									Size: to.Ptr("A1"),
		// 									Tier: to.Ptr("Standard"),
		// 							}},
		// 						}
	}
}
