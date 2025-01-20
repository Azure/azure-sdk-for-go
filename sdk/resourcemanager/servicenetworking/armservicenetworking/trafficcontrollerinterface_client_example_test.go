// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armservicenetworking_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicenetworking/armservicenetworking/v2"
	"log"
)

// Generated from example definition: 2025-01-01/TrafficControllerPut.json
func ExampleTrafficControllerInterfaceClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armservicenetworking.NewClientFactory("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewTrafficControllerInterfaceClient().BeginCreateOrUpdate(ctx, "rg1", "tc1", armservicenetworking.TrafficController{
		Location: to.Ptr("NorthCentralUS"),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
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
	// res = armservicenetworking.TrafficControllerInterfaceClientCreateOrUpdateResponse{
	// 	TrafficController: &armservicenetworking.TrafficController{
	// 		Name: to.Ptr("tc1"),
	// 		ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1"),
	// 		Type: to.Ptr("Microsoft.ServiceNetworking/trafficControllers"),
	// 		Location: to.Ptr("NorthCentralUS"),
	// 		Tags: map[string]*string{
	// 			"key1": to.Ptr("value1"),
	// 		},
	// 		Properties: &armservicenetworking.TrafficControllerProperties{
	// 			ConfigurationEndpoints: []*string{
	// 				to.Ptr("abc.trafficcontroller.azure.net"),
	// 			},
	// 			Frontends: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/frontends/fe1"),
	// 				},
	// 			},
	// 			Associations: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/association/as1"),
	// 				},
	// 			},
	// 			SecurityPolicies: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/sp1"),
	// 				},
	// 			},
	// 			SecurityPolicyConfigurations: &armservicenetworking.SecurityPolicyConfigurations{
	// 				WafSecurityPolicy: &armservicenetworking.WafSecurityPolicy{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/waf-0"),
	// 				},
	// 			},
	// 			ProvisioningState: to.Ptr(armservicenetworking.ProvisioningStateSucceeded),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-01-01/TrafficControllerDelete.json
func ExampleTrafficControllerInterfaceClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armservicenetworking.NewClientFactory("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewTrafficControllerInterfaceClient().BeginDelete(ctx, "rg1", "tc1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: 2025-01-01/TrafficControllerGet.json
func ExampleTrafficControllerInterfaceClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armservicenetworking.NewClientFactory("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewTrafficControllerInterfaceClient().Get(ctx, "rg1", "tc1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armservicenetworking.TrafficControllerInterfaceClientGetResponse{
	// 	TrafficController: &armservicenetworking.TrafficController{
	// 		Name: to.Ptr("tc1"),
	// 		ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1"),
	// 		Type: to.Ptr("Microsoft.ServiceNetworking/trafficControllers"),
	// 		Location: to.Ptr("NorthCentralUS"),
	// 		Tags: map[string]*string{
	// 			"key1": to.Ptr("value1"),
	// 		},
	// 		Properties: &armservicenetworking.TrafficControllerProperties{
	// 			ConfigurationEndpoints: []*string{
	// 				to.Ptr("abc.trafficcontroller.azure.net"),
	// 			},
	// 			Frontends: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/frontends/fe1"),
	// 				},
	// 			},
	// 			Associations: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/association/as1"),
	// 				},
	// 			},
	// 			ProvisioningState: to.Ptr(armservicenetworking.ProvisioningStateSucceeded),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-01-01/TrafficControllersGet.json
func ExampleTrafficControllerInterfaceClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armservicenetworking.NewClientFactory("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewTrafficControllerInterfaceClient().NewListByResourceGroupPager("rg1", nil)
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
		// page = armservicenetworking.TrafficControllerInterfaceClientListByResourceGroupResponse{
		// 	TrafficControllerListResult: armservicenetworking.TrafficControllerListResult{
		// 		Value: []*armservicenetworking.TrafficController{
		// 			{
		// 				Name: to.Ptr("tc1"),
		// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1"),
		// 				Type: to.Ptr("Microsoft.ServiceNetworking/trafficControllers"),
		// 				Location: to.Ptr("NorthCentralUS"),
		// 				Tags: map[string]*string{
		// 					"key1": to.Ptr("value1"),
		// 				},
		// 				Properties: &armservicenetworking.TrafficControllerProperties{
		// 					ConfigurationEndpoints: []*string{
		// 						to.Ptr("abc.trafficcontroller.azure.net"),
		// 					},
		// 					Frontends: []*armservicenetworking.ResourceID{
		// 						{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/frontends/fe1"),
		// 						},
		// 					},
		// 					Associations: []*armservicenetworking.ResourceID{
		// 						{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/association/as1"),
		// 						},
		// 					},
		// 					SecurityPolicies: []*armservicenetworking.ResourceID{
		// 						{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/sp1"),
		// 						},
		// 					},
		// 					SecurityPolicyConfigurations: &armservicenetworking.SecurityPolicyConfigurations{
		// 						WafSecurityPolicy: &armservicenetworking.WafSecurityPolicy{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/waf-0"),
		// 						},
		// 					},
		// 					ProvisioningState: to.Ptr(armservicenetworking.ProvisioningStateSucceeded),
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}

// Generated from example definition: 2025-01-01/TrafficControllersGetList.json
func ExampleTrafficControllerInterfaceClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armservicenetworking.NewClientFactory("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewTrafficControllerInterfaceClient().NewListBySubscriptionPager(nil)
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
		// page = armservicenetworking.TrafficControllerInterfaceClientListBySubscriptionResponse{
		// 	TrafficControllerListResult: armservicenetworking.TrafficControllerListResult{
		// 		Value: []*armservicenetworking.TrafficController{
		// 			{
		// 				Name: to.Ptr("tc1"),
		// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1"),
		// 				Type: to.Ptr("Microsoft.ServiceNetworking/trafficControllers"),
		// 				Location: to.Ptr("NorthCentralUS"),
		// 				Tags: map[string]*string{
		// 					"key1": to.Ptr("value1"),
		// 				},
		// 				Properties: &armservicenetworking.TrafficControllerProperties{
		// 					ConfigurationEndpoints: []*string{
		// 						to.Ptr("abc.trafficcontroller.azure.net"),
		// 					},
		// 					Frontends: []*armservicenetworking.ResourceID{
		// 						{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/frontends/fe1"),
		// 						},
		// 					},
		// 					Associations: []*armservicenetworking.ResourceID{
		// 						{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/association/as1"),
		// 						},
		// 					},
		// 					SecurityPolicies: []*armservicenetworking.ResourceID{
		// 						{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/sp1"),
		// 						},
		// 					},
		// 					SecurityPolicyConfigurations: &armservicenetworking.SecurityPolicyConfigurations{
		// 						WafSecurityPolicy: &armservicenetworking.WafSecurityPolicy{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/waf-0"),
		// 						},
		// 					},
		// 					ProvisioningState: to.Ptr(armservicenetworking.ProvisioningStateSucceeded),
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}

// Generated from example definition: 2025-01-01/TrafficControllerPatch.json
func ExampleTrafficControllerInterfaceClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armservicenetworking.NewClientFactory("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewTrafficControllerInterfaceClient().Update(ctx, "rg1", "tc1", armservicenetworking.TrafficControllerUpdate{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armservicenetworking.TrafficControllerInterfaceClientUpdateResponse{
	// 	TrafficController: &armservicenetworking.TrafficController{
	// 		Name: to.Ptr("tc1"),
	// 		ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1"),
	// 		Type: to.Ptr("Microsoft.ServiceNetworking/trafficControllers"),
	// 		Location: to.Ptr("NorthCentralUS"),
	// 		Tags: map[string]*string{
	// 			"key1": to.Ptr("value1"),
	// 		},
	// 		Properties: &armservicenetworking.TrafficControllerProperties{
	// 			ConfigurationEndpoints: []*string{
	// 				to.Ptr("abc.trafficcontroller.azure.net"),
	// 			},
	// 			Frontends: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/frontends/fe1"),
	// 				},
	// 			},
	// 			Associations: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/association/as1"),
	// 				},
	// 			},
	// 			SecurityPolicies: []*armservicenetworking.ResourceID{
	// 				{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/sp1"),
	// 				},
	// 			},
	// 			SecurityPolicyConfigurations: &armservicenetworking.SecurityPolicyConfigurations{
	// 				WafSecurityPolicy: &armservicenetworking.WafSecurityPolicy{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/tc1/securityPolicies/waf-0"),
	// 				},
	// 			},
	// 			ProvisioningState: to.Ptr(armservicenetworking.ProvisioningStateSucceeded),
	// 		},
	// 	},
	// }
}
