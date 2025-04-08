// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armhybridconnectivity_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridconnectivity/armhybridconnectivity"
	"log"
)

// Generated from example definition: 2024-12-01/ServiceConfigurationsPutSSH.json
func ExampleServiceConfigurationsClient_CreateOrupdate_serviceConfigurationsPutSsh() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceConfigurationsClient().CreateOrupdate(ctx, "subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default", "default", "SSH", armhybridconnectivity.ServiceConfigurationResource{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.ServiceConfigurationsClientCreateOrupdateResponse{
	// 	ServiceConfigurationResource: &armhybridconnectivity.ServiceConfigurationResource{
	// 		Type: to.Ptr("Microsoft.HybridConnectivity/endpoints/serviceConfigurations"),
	// 		ID: to.Ptr("/subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default/serviceconfigurations/SSH"),
	// 		Properties: &armhybridconnectivity.ServiceConfigurationProperties{
	// 			Port: to.Ptr[int64](22),
	// 			ProvisioningState: to.Ptr(armhybridconnectivity.ProvisioningStateSucceeded),
	// 			ServiceName: to.Ptr(armhybridconnectivity.ServiceNameSSH),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-12-01/ServiceConfigurationsPutWAC.json
func ExampleServiceConfigurationsClient_CreateOrupdate_serviceConfigurationsPutWac() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceConfigurationsClient().CreateOrupdate(ctx, "subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default", "default", "WAC", armhybridconnectivity.ServiceConfigurationResource{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.ServiceConfigurationsClientCreateOrupdateResponse{
	// 	ServiceConfigurationResource: &armhybridconnectivity.ServiceConfigurationResource{
	// 		Type: to.Ptr("Microsoft.HybridConnectivity/endpoints/serviceConfigurations"),
	// 		ID: to.Ptr("/subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default/serviceconfigurations/WAC"),
	// 		Properties: &armhybridconnectivity.ServiceConfigurationProperties{
	// 			Port: to.Ptr[int64](6516),
	// 			ProvisioningState: to.Ptr(armhybridconnectivity.ProvisioningStateSucceeded),
	// 			ServiceName: to.Ptr(armhybridconnectivity.ServiceNameWAC),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-12-01/ServiceConfigurationsDeleteSSH.json
func ExampleServiceConfigurationsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceConfigurationsClient().Delete(ctx, "subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default", "default", "SSH", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.ServiceConfigurationsClientDeleteResponse{
	// }
}

// Generated from example definition: 2024-12-01/ServiceConfigurationsGetSSH.json
func ExampleServiceConfigurationsClient_Get_hybridConnectivityEndpointsServiceconfigurationsGetSsh() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceConfigurationsClient().Get(ctx, "subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default", "default", "SSH", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.ServiceConfigurationsClientGetResponse{
	// 	ServiceConfigurationResource: &armhybridconnectivity.ServiceConfigurationResource{
	// 		Type: to.Ptr("Microsoft.HybridConnectivity/endpoints/serviceConfigurations"),
	// 		ID: to.Ptr("/subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default/serviceconfigurations/SSH"),
	// 		Properties: &armhybridconnectivity.ServiceConfigurationProperties{
	// 			Port: to.Ptr[int64](22),
	// 			ServiceName: to.Ptr(armhybridconnectivity.ServiceNameSSH),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-12-01/ServiceConfigurationsGetWAC.json
func ExampleServiceConfigurationsClient_Get_hybridConnectivityEndpointsServiceconfigurationsGetWac() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceConfigurationsClient().Get(ctx, "subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default", "default", "WAC", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.ServiceConfigurationsClientGetResponse{
	// 	ServiceConfigurationResource: &armhybridconnectivity.ServiceConfigurationResource{
	// 		Type: to.Ptr("Microsoft.HybridConnectivity/endpoints/serviceConfigurations"),
	// 		ID: to.Ptr("/subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default/serviceconfigurations/WAC"),
	// 		Properties: &armhybridconnectivity.ServiceConfigurationProperties{
	// 			Port: to.Ptr[int64](6516),
	// 			ServiceName: to.Ptr(armhybridconnectivity.ServiceNameWAC),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-12-01/ServiceConfigurationsList.json
func ExampleServiceConfigurationsClient_NewListByEndpointResourcePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewServiceConfigurationsClient().NewListByEndpointResourcePager("subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default", "default", nil)
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
		// page = armhybridconnectivity.ServiceConfigurationsClientListByEndpointResourceResponse{
		// 	ServiceConfigurationList: armhybridconnectivity.ServiceConfigurationList{
		// 		Value: []*armhybridconnectivity.ServiceConfigurationResource{
		// 			{
		// 				Type: to.Ptr("Microsoft.HybridConnectivity/endpoints/serviceConfigurations"),
		// 				ID: to.Ptr("/subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default/serviceconfigurations/SSH"),
		// 				Properties: &armhybridconnectivity.ServiceConfigurationProperties{
		// 					Port: to.Ptr[int64](22),
		// 					ServiceName: to.Ptr(armhybridconnectivity.ServiceNameSSH),
		// 				},
		// 			},
		// 			{
		// 				Type: to.Ptr("Microsoft.HybridConnectivity/endpoints/serviceConfigurations"),
		// 				ID: to.Ptr("/subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default/serviceconfigurations/WAC"),
		// 				Properties: &armhybridconnectivity.ServiceConfigurationProperties{
		// 					Port: to.Ptr[int64](6516),
		// 					ServiceName: to.Ptr(armhybridconnectivity.ServiceNameWAC),
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}

// Generated from example definition: 2024-12-01/ServiceConfigurationsPatchSSH.json
func ExampleServiceConfigurationsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceConfigurationsClient().Update(ctx, "subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default", "default", "SSH", armhybridconnectivity.ServiceConfigurationResourcePatch{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.ServiceConfigurationsClientUpdateResponse{
	// 	ServiceConfigurationResource: &armhybridconnectivity.ServiceConfigurationResource{
	// 		ID: to.Ptr("/subscriptions/f5bcc1d9-23af-4ae9-aca1-041d0f593a63/resourceGroups/hybridRG/providers/Microsoft.HybridCompute/machines/testMachine/providers/Microsoft.HybridConnectivity/endpoints/default/serviceConfigurations/SSH"),
	// 		Properties: &armhybridconnectivity.ServiceConfigurationProperties{
	// 			Port: to.Ptr[int64](22),
	// 			ServiceName: to.Ptr(armhybridconnectivity.ServiceNameSSH),
	// 		},
	// 	},
	// }
}
