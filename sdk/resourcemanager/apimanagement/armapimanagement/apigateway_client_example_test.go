//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armapimanagement_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateStandardGateway.json
func ExampleAPIGatewayClient_BeginCreateOrUpdate_apiManagementCreateStandardGateway() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAPIGatewayClient().BeginCreateOrUpdate(ctx, "rg1", "apimGateway1", armapimanagement.GatewayResource{
		Tags: map[string]*string{
			"Name": to.Ptr("Contoso"),
			"Test": to.Ptr("User"),
		},
		Location: to.Ptr("South Central US"),
		Properties: &armapimanagement.GatewayProperties{
			Backend: &armapimanagement.BackendConfiguration{
				Subnet: &armapimanagement.BackendSubnetConfiguration{
					ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
				},
			},
		},
		SKU: &armapimanagement.GatewaySKUProperties{
			Name:     to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
			Capacity: to.Ptr[int32](1),
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
	// res.GatewayResource = armapimanagement.GatewayResource{
	// 	Name: to.Ptr("apimGateway1"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/gateways"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/gateways/apimGateway1"),
	// 	Tags: map[string]*string{
	// 		"api-version": to.Ptr("2024-05-01"),
	// 	},
	// 	Etag: to.Ptr("AAAAAAAmREI="),
	// 	Location: to.Ptr("East US"),
	// 	Properties: &armapimanagement.GatewayProperties{
	// 		Backend: &armapimanagement.BackendConfiguration{
	// 			Subnet: &armapimanagement.BackendSubnetConfiguration{
	// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
	// 			},
	// 		},
	// 		ConfigurationAPI: &armapimanagement.GatewayConfigurationAPI{
	// 			Hostname: to.Ptr("apimGateway1.eastus.configuration.gateway.azure-api.net"),
	// 		},
	// 		CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:01.250Z"); return t}()),
	// 		Frontend: &armapimanagement.FrontendConfiguration{
	// 			DefaultHostname: to.Ptr("apimGateway1.eastus.gateway.azure-api.net"),
	// 		},
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		TargetProvisioningState: to.Ptr(""),
	// 	},
	// 	SKU: &armapimanagement.GatewaySKUProperties{
	// 		Name: to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
	// 		Capacity: to.Ptr[int32](1),
	// 	},
	// 	SystemData: &armapimanagement.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:00.939Z"); return t}()),
	// 		CreatedBy: to.Ptr("user@contoso.com"),
	// 		CreatedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:00.939Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user@contoso.com"),
	// 		LastModifiedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateWorkspacePremiumGateway.json
func ExampleAPIGatewayClient_BeginCreateOrUpdate_apiManagementCreateWorkspacePremiumGateway() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAPIGatewayClient().BeginCreateOrUpdate(ctx, "rg1", "apimGateway1", armapimanagement.GatewayResource{
		Tags: map[string]*string{
			"Name": to.Ptr("Contoso"),
			"Test": to.Ptr("User"),
		},
		Location: to.Ptr("South Central US"),
		Properties: &armapimanagement.GatewayProperties{
			Backend: &armapimanagement.BackendConfiguration{
				Subnet: &armapimanagement.BackendSubnetConfiguration{
					ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
				},
			},
			VirtualNetworkType: to.Ptr(armapimanagement.VirtualNetworkTypeExternal),
		},
		SKU: &armapimanagement.GatewaySKUProperties{
			Name:     to.Ptr(armapimanagement.APIGatewaySKUTypeWorkspaceGatewayPremium),
			Capacity: to.Ptr[int32](1),
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
	// res.GatewayResource = armapimanagement.GatewayResource{
	// 	Name: to.Ptr("apimGateway1"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/gateways"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/gateways/apimGateway1"),
	// 	Tags: map[string]*string{
	// 		"api-version": to.Ptr("2024-05-01"),
	// 	},
	// 	Etag: to.Ptr("AAAAAAAmREI="),
	// 	Location: to.Ptr("East US"),
	// 	Properties: &armapimanagement.GatewayProperties{
	// 		Backend: &armapimanagement.BackendConfiguration{
	// 			Subnet: &armapimanagement.BackendSubnetConfiguration{
	// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
	// 			},
	// 		},
	// 		CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:01.250Z"); return t}()),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		TargetProvisioningState: to.Ptr(""),
	// 		VirtualNetworkType: to.Ptr(armapimanagement.VirtualNetworkTypeExternal),
	// 	},
	// 	SKU: &armapimanagement.GatewaySKUProperties{
	// 		Name: to.Ptr(armapimanagement.APIGatewaySKUTypeWorkspaceGatewayPremium),
	// 		Capacity: to.Ptr[int32](1),
	// 	},
	// 	SystemData: &armapimanagement.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:00.939Z"); return t}()),
	// 		CreatedBy: to.Ptr("user@contoso.com"),
	// 		CreatedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:00.939Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user@contoso.com"),
	// 		LastModifiedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementUpdateStandardGateway.json
func ExampleAPIGatewayClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAPIGatewayClient().BeginUpdate(ctx, "rg1", "apimGateway1", armapimanagement.GatewayUpdateParameters{
		Tags: map[string]*string{
			"Name": to.Ptr("Contoso"),
			"Test": to.Ptr("User"),
		},
		Properties: &armapimanagement.GatewayUpdateProperties{},
		SKU: &armapimanagement.GatewaySKUPropertiesForPatch{
			Name:     to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
			Capacity: to.Ptr[int32](10),
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
	// res.GatewayResource = armapimanagement.GatewayResource{
	// 	Name: to.Ptr("apimGateway1"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/gateways"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/gateways/apimGateway1"),
	// 	Tags: map[string]*string{
	// 		"api-version": to.Ptr("2024-05-01"),
	// 	},
	// 	Etag: to.Ptr("AAAAAAAmREI="),
	// 	Location: to.Ptr("East US"),
	// 	Properties: &armapimanagement.GatewayProperties{
	// 		Backend: &armapimanagement.BackendConfiguration{
	// 			Subnet: &armapimanagement.BackendSubnetConfiguration{
	// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
	// 			},
	// 		},
	// 		ConfigurationAPI: &armapimanagement.GatewayConfigurationAPI{
	// 			Hostname: to.Ptr("apimGateway1.eastus.configuration.gateway.azure-api.net"),
	// 		},
	// 		CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:01.250Z"); return t}()),
	// 		Frontend: &armapimanagement.FrontendConfiguration{
	// 			DefaultHostname: to.Ptr("apimGateway1.eastus.gateway.azure-api.net"),
	// 		},
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		TargetProvisioningState: to.Ptr(""),
	// 	},
	// 	SKU: &armapimanagement.GatewaySKUProperties{
	// 		Name: to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
	// 		Capacity: to.Ptr[int32](1),
	// 	},
	// 	SystemData: &armapimanagement.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:00.939Z"); return t}()),
	// 		CreatedBy: to.Ptr("user@contoso.com"),
	// 		CreatedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-11T18:41:00.939Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user@contoso.com"),
	// 		LastModifiedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGatewayGetGateway.json
func ExampleAPIGatewayClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIGatewayClient().Get(ctx, "rg1", "apimService1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GatewayResource = armapimanagement.GatewayResource{
	// 	Name: to.Ptr("example-gateway"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/gateway"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/gateway/example-gateway"),
	// 	Tags: map[string]*string{
	// 		"ReleaseName": to.Ptr("Z3"),
	// 		"owner": to.Ptr("v-aswmoh"),
	// 	},
	// 	Etag: to.Ptr("AAAAAAAWN/4="),
	// 	Location: to.Ptr("East US"),
	// 	Properties: &armapimanagement.GatewayProperties{
	// 		Backend: &armapimanagement.BackendConfiguration{
	// 			Subnet: &armapimanagement.BackendSubnetConfiguration{
	// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
	// 			},
	// 		},
	// 		ConfigurationAPI: &armapimanagement.GatewayConfigurationAPI{
	// 			Hostname: to.Ptr("example-gateway.eastus.configuration.gateway.azure-api.net"),
	// 		},
	// 		CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.945Z"); return t}()),
	// 		Frontend: &armapimanagement.FrontendConfiguration{
	// 			DefaultHostname: to.Ptr("example-gateway.eastus.gateway.azure-api.net"),
	// 		},
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		TargetProvisioningState: to.Ptr(""),
	// 	},
	// 	SKU: &armapimanagement.GatewaySKUProperties{
	// 		Name: to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
	// 		Capacity: to.Ptr[int32](1),
	// 	},
	// 	SystemData: &armapimanagement.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.710Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-20T06:33:09.615Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("foo@contoso.com"),
	// 		LastModifiedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGatewayDeleteGateway.json
func ExampleAPIGatewayClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAPIGatewayClient().BeginDelete(ctx, "rg1", "example-gateway", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListGatewaysBySubscriptionAndResourceGroup.json
func ExampleAPIGatewayClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAPIGatewayClient().NewListByResourceGroupPager("rg1", nil)
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
		// page.GatewayListResult = armapimanagement.GatewayListResult{
		// 	Value: []*armapimanagement.GatewayResource{
		// 		{
		// 			Name: to.Ptr("standard-gw-1"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/gateways"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/gateways/standard-gw-1"),
		// 			Tags: map[string]*string{
		// 				"ReleaseName": to.Ptr("Z3"),
		// 				"owner": to.Ptr("v-aswmoh"),
		// 			},
		// 			Etag: to.Ptr("AAAAAAAWN/4="),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armapimanagement.GatewayProperties{
		// 				Backend: &armapimanagement.BackendConfiguration{
		// 					Subnet: &armapimanagement.BackendSubnetConfiguration{
		// 						ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
		// 					},
		// 				},
		// 				ConfigurationAPI: &armapimanagement.GatewayConfigurationAPI{
		// 					Hostname: to.Ptr("standard-gw-1.westus.configuration.gateway.azure-api.net"),
		// 				},
		// 				CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.945Z"); return t}()),
		// 				Frontend: &armapimanagement.FrontendConfiguration{
		// 					DefaultHostname: to.Ptr("standard-gw-1.westus.gateway.azure-api.net"),
		// 				},
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TargetProvisioningState: to.Ptr(""),
		// 			},
		// 			SKU: &armapimanagement.GatewaySKUProperties{
		// 				Name: to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
		// 				Capacity: to.Ptr[int32](1),
		// 			},
		// 			SystemData: &armapimanagement.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.710Z"); return t}()),
		// 				CreatedBy: to.Ptr("bar@contoso.com"),
		// 				CreatedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-20T06:33:09.615Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("foo@contoso.com"),
		// 				LastModifiedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("standard-gw-2"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/gateways"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/gateways/standard-gw-2"),
		// 			Tags: map[string]*string{
		// 				"Owner": to.Ptr("vitaliik"),
		// 			},
		// 			Etag: to.Ptr("AAAAAAAWKwo="),
		// 			Location: to.Ptr("East US"),
		// 			Properties: &armapimanagement.GatewayProperties{
		// 				Backend: &armapimanagement.BackendConfiguration{
		// 					Subnet: &armapimanagement.BackendSubnetConfiguration{
		// 						ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn2/subnets/sn2"),
		// 					},
		// 				},
		// 				ConfigurationAPI: &armapimanagement.GatewayConfigurationAPI{
		// 					Hostname: to.Ptr("standard-gw-2.eastus.configuration.gateway.azure-api.net"),
		// 				},
		// 				CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.945Z"); return t}()),
		// 				Frontend: &armapimanagement.FrontendConfiguration{
		// 					DefaultHostname: to.Ptr("standard-gw-2.eastus.gateway.azure-api.net"),
		// 				},
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TargetProvisioningState: to.Ptr(""),
		// 			},
		// 			SKU: &armapimanagement.GatewaySKUProperties{
		// 				Name: to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
		// 				Capacity: to.Ptr[int32](1),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListGatewaysBySubscription.json
func ExampleAPIGatewayClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAPIGatewayClient().NewListPager(nil)
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
		// page.GatewayListResult = armapimanagement.GatewayListResult{
		// 	Value: []*armapimanagement.GatewayResource{
		// 		{
		// 			Name: to.Ptr("standard-gw-1"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/gateways"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/gateways/standard-gw-1"),
		// 			Tags: map[string]*string{
		// 				"ReleaseName": to.Ptr("Z3"),
		// 				"owner": to.Ptr("v-aswmoh"),
		// 			},
		// 			Etag: to.Ptr("AAAAAAAWN/4="),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armapimanagement.GatewayProperties{
		// 				Backend: &armapimanagement.BackendConfiguration{
		// 					Subnet: &armapimanagement.BackendSubnetConfiguration{
		// 						ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vn1/subnets/sn1"),
		// 					},
		// 				},
		// 				ConfigurationAPI: &armapimanagement.GatewayConfigurationAPI{
		// 					Hostname: to.Ptr("standard-gw-1.westus.configuration.gateway.azure-api.net"),
		// 				},
		// 				CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.945Z"); return t}()),
		// 				Frontend: &armapimanagement.FrontendConfiguration{
		// 					DefaultHostname: to.Ptr("standard-gw-1.westus.gateway.azure-api.net"),
		// 				},
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TargetProvisioningState: to.Ptr(""),
		// 			},
		// 			SKU: &armapimanagement.GatewaySKUProperties{
		// 				Name: to.Ptr(armapimanagement.APIGatewaySKUTypeStandard),
		// 				Capacity: to.Ptr[int32](1),
		// 			},
		// 			SystemData: &armapimanagement.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.710Z"); return t}()),
		// 				CreatedBy: to.Ptr("bar@contoso.com"),
		// 				CreatedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-20T06:33:09.615Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("foo@contoso.com"),
		// 				LastModifiedByType: to.Ptr(armapimanagement.CreatedByTypeUser),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("workspace-gw-2"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/gateways"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg2/providers/Microsoft.ApiManagement/gateways/workspace-gw-2"),
		// 			Tags: map[string]*string{
		// 				"Owner": to.Ptr("foo"),
		// 			},
		// 			Etag: to.Ptr("AAAAAAAWKwo="),
		// 			Location: to.Ptr("East US"),
		// 			Properties: &armapimanagement.GatewayProperties{
		// 				CreatedAtUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-06-16T09:40:00.945Z"); return t}()),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TargetProvisioningState: to.Ptr(""),
		// 			},
		// 			SKU: &armapimanagement.GatewaySKUProperties{
		// 				Name: to.Ptr(armapimanagement.APIGatewaySKUTypeWorkspaceGatewayPremium),
		// 				Capacity: to.Ptr[int32](1),
		// 			},
		// 	}},
		// }
	}
}
