//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdatafactory_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_ListByFactory.json
func ExampleIntegrationRuntimesClient_NewListByFactoryPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByFactoryPager("exampleResourceGroup", "exampleFactoryName", nil)
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
		// page.IntegrationRuntimeListResponse = armdatafactory.IntegrationRuntimeListResponse{
		// 	Value: []*armdatafactory.IntegrationRuntimeResource{
		// 		{
		// 			Name: to.Ptr("exampleIntegrationRuntime"),
		// 			Type: to.Ptr("Microsoft.DataFactory/factories/integrationruntimes"),
		// 			Etag: to.Ptr("0400f1a1-0000-0000-0000-5b2188640000"),
		// 			ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/integrationruntimes/exampleIntegrationRuntime"),
		// 			Properties: &armdatafactory.SelfHostedIntegrationRuntime{
		// 				Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
		// 				Description: to.Ptr("A selfhosted integration runtime"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_Create.json
func ExampleIntegrationRuntimesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", armdatafactory.IntegrationRuntimeResource{
		Properties: &armdatafactory.SelfHostedIntegrationRuntime{
			Type:        to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
			Description: to.Ptr("A selfhosted integration runtime"),
		},
	}, &armdatafactory.IntegrationRuntimesClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeResource = armdatafactory.IntegrationRuntimeResource{
	// 	Name: to.Ptr("exampleIntegrationRuntime"),
	// 	Type: to.Ptr("Microsoft.DataFactory/factories/integrationruntimes"),
	// 	Etag: to.Ptr("000046c4-0000-0000-0000-5b2198bf0000"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/integrationruntimes/exampleIntegrationRuntime"),
	// 	Properties: &armdatafactory.SelfHostedIntegrationRuntime{
	// 		Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
	// 		Description: to.Ptr("A selfhosted integration runtime"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_Get.json
func ExampleIntegrationRuntimesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", &armdatafactory.IntegrationRuntimesClientGetOptions{IfNoneMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeResource = armdatafactory.IntegrationRuntimeResource{
	// 	Name: to.Ptr("exampleIntegrationRuntime"),
	// 	Type: to.Ptr("Microsoft.DataFactory/factories/integrationruntimes"),
	// 	Etag: to.Ptr("15003c4f-0000-0200-0000-5cbe090b0000"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/integrationruntimes/exampleIntegrationRuntime"),
	// 	Properties: &armdatafactory.SelfHostedIntegrationRuntime{
	// 		Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
	// 		Description: to.Ptr("A selfhosted integration runtime"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_Update.json
func ExampleIntegrationRuntimesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", armdatafactory.UpdateIntegrationRuntimeRequest{
		AutoUpdate:        to.Ptr(armdatafactory.IntegrationRuntimeAutoUpdateOff),
		UpdateDelayOffset: to.Ptr("\"PT3H\""),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeResource = armdatafactory.IntegrationRuntimeResource{
	// 	Name: to.Ptr("exampleIntegrationRuntime"),
	// 	Type: to.Ptr("Microsoft.DataFactory/factories/integrationruntimes"),
	// 	Etag: to.Ptr("0400f1a1-0000-0000-0000-5b2188640000"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/integrationruntimes/exampleIntegrationRuntime"),
	// 	Properties: &armdatafactory.SelfHostedIntegrationRuntime{
	// 		Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
	// 		Description: to.Ptr("A selfhosted integration runtime"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_Delete.json
func ExampleIntegrationRuntimesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_GetStatus.json
func ExampleIntegrationRuntimesClient_GetStatus() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetStatus(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeStatusResponse = armdatafactory.IntegrationRuntimeStatusResponse{
	// 	Name: to.Ptr("exampleIntegrationRuntime"),
	// 	Properties: &armdatafactory.SelfHostedIntegrationRuntimeStatus{
	// 		Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
	// 		State: to.Ptr(armdatafactory.IntegrationRuntimeStateOnline),
	// 		TypeProperties: &armdatafactory.SelfHostedIntegrationRuntimeStatusTypeProperties{
	// 			AutoUpdate: to.Ptr(armdatafactory.IntegrationRuntimeAutoUpdateOff),
	// 			Capabilities: map[string]*string{
	// 				"connectedToResourceManager": to.Ptr("True"),
	// 				"credentialInSync": to.Ptr("True"),
	// 				"httpsPortEnabled": to.Ptr("True"),
	// 				"nodeEnabled": to.Ptr("True"),
	// 				"serviceBusConnected": to.Ptr("True"),
	// 			},
	// 			CreateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-14T09:17:45.1839685Z"); return t}()),
	// 			LatestVersion: to.Ptr("3.7.6711.1"),
	// 			LocalTimeZoneOffset: to.Ptr("PT8H"),
	// 			Nodes: []*armdatafactory.SelfHostedIntegrationRuntimeNode{
	// 				{
	// 					Capabilities: map[string]*string{
	// 						"connectedToResourceManager": to.Ptr("True"),
	// 						"credentialInSync": to.Ptr("True"),
	// 						"httpsPortEnabled": to.Ptr("True"),
	// 						"nodeEnabled": to.Ptr("True"),
	// 						"serviceBusConnected": to.Ptr("True"),
	// 					},
	// 					HostServiceURI: to.Ptr("https://yanzhang-dt.fareast.corp.microsoft.com:8050/HostServiceRemote.svc/"),
	// 					IsActiveDispatcher: to.Ptr(true),
	// 					LastConnectTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-14T14:52:59.8933313Z"); return t}()),
	// 					LastStartTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-14T14:52:59.8933313Z"); return t}()),
	// 					LastUpdateResult: to.Ptr(armdatafactory.IntegrationRuntimeUpdateResultNone),
	// 					MachineName: to.Ptr("YANZHANG-DT"),
	// 					MaxConcurrentJobs: to.Ptr[int32](56),
	// 					NodeName: to.Ptr("Node_1"),
	// 					RegisterTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-14T14:51:44.9237069Z"); return t}()),
	// 					Status: to.Ptr(armdatafactory.SelfHostedIntegrationRuntimeNodeStatusOnline),
	// 					Version: to.Ptr("3.8.6730.2"),
	// 					VersionStatus: to.Ptr("UpToDate"),
	// 			}},
	// 			ServiceUrls: []*string{
	// 				to.Ptr("wu.frontend.int.clouddatahub-int.net"),
	// 				to.Ptr("*.servicebus.windows.net")},
	// 				TaskQueueID: to.Ptr("1a6296ab-423c-4346-9bcc-85a78c2c0582"),
	// 				UpdateDelayOffset: to.Ptr("PT3H"),
	// 				Version: to.Ptr("3.8.6730.2"),
	// 				VersionStatus: to.Ptr("UpToDate"),
	// 			},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_ListOutboundNetworkDependenciesEndpoints.json
func ExampleIntegrationRuntimesClient_ListOutboundNetworkDependenciesEndpoints() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("7ad7c73b-38b8-4df3-84ee-52ff91092f61", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListOutboundNetworkDependenciesEndpoints(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse = armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointsResponse{
	// 	Value: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesCategoryEndpoint{
	// 		{
	// 			Category: to.Ptr("Azure Data Factory (Management)"),
	// 			Endpoints: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpoint{
	// 				{
	// 					DomainName: to.Ptr("wu.frontend.int.clouddatahub-int.net"),
	// 					EndpointDetails: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails{
	// 						{
	// 							Port: to.Ptr[int32](443),
	// 					}},
	// 			}},
	// 		},
	// 		{
	// 			Category: to.Ptr("Azure Storage (Management)"),
	// 			Endpoints: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpoint{
	// 				{
	// 					DomainName: to.Ptr("*.blob.core.windows.net"),
	// 					EndpointDetails: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails{
	// 						{
	// 							Port: to.Ptr[int32](443),
	// 					}},
	// 				},
	// 				{
	// 					DomainName: to.Ptr("*.table.core.windows.net"),
	// 					EndpointDetails: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails{
	// 						{
	// 							Port: to.Ptr[int32](443),
	// 					}},
	// 			}},
	// 		},
	// 		{
	// 			Category: to.Ptr("Event Hub (Logging)"),
	// 			Endpoints: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpoint{
	// 				{
	// 					DomainName: to.Ptr("*.servicebus.windows.net"),
	// 					EndpointDetails: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails{
	// 						{
	// 							Port: to.Ptr[int32](443),
	// 					}},
	// 			}},
	// 		},
	// 		{
	// 			Category: to.Ptr("Microsoft Logging service (Internal Use)"),
	// 			Endpoints: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpoint{
	// 				{
	// 					DomainName: to.Ptr("gcs.prod.monitoring.core.windows.net"),
	// 					EndpointDetails: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails{
	// 						{
	// 							Port: to.Ptr[int32](443),
	// 					}},
	// 				},
	// 				{
	// 					DomainName: to.Ptr("prod.warmpath.msftcloudes.com"),
	// 					EndpointDetails: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails{
	// 						{
	// 							Port: to.Ptr[int32](443),
	// 					}},
	// 				},
	// 				{
	// 					DomainName: to.Ptr("azurewatsonanalysis-prod.core.windows.net"),
	// 					EndpointDetails: []*armdatafactory.IntegrationRuntimeOutboundNetworkDependenciesEndpointDetails{
	// 						{
	// 							Port: to.Ptr[int32](443),
	// 					}},
	// 			}},
	// 	}},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_GetConnectionInfo.json
func ExampleIntegrationRuntimesClient_GetConnectionInfo() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetConnectionInfo(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeConnectionInfo = armdatafactory.IntegrationRuntimeConnectionInfo{
	// 	HostServiceURI: to.Ptr("https://yanzhang-dt.fareast.corp.microsoft.com:8050/HostServiceRemote.svc/"),
	// 	IdentityCertThumbprint: to.Ptr("**********"),
	// 	IsIdentityCertExprired: to.Ptr(false),
	// 	PublicKey: to.Ptr("**********"),
	// 	ServiceToken: to.Ptr("**********"),
	// 	Version: to.Ptr("3.8.6730.2"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_RegenerateAuthKey.json
func ExampleIntegrationRuntimesClient_RegenerateAuthKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.RegenerateAuthKey(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", armdatafactory.IntegrationRuntimeRegenerateKeyParameters{
		KeyName: to.Ptr(armdatafactory.IntegrationRuntimeAuthKeyNameAuthKey2),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeAuthKeys = armdatafactory.IntegrationRuntimeAuthKeys{
	// 	AuthKey2: to.Ptr("**********"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_ListAuthKeys.json
func ExampleIntegrationRuntimesClient_ListAuthKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListAuthKeys(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeAuthKeys = armdatafactory.IntegrationRuntimeAuthKeys{
	// 	AuthKey1: to.Ptr("**********"),
	// 	AuthKey2: to.Ptr("**********"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_Start.json
func ExampleIntegrationRuntimesClient_BeginStart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginStart(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleManagedIntegrationRuntime", nil)
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
	// res.IntegrationRuntimeStatusResponse = armdatafactory.IntegrationRuntimeStatusResponse{
	// 	Name: to.Ptr("exampleManagedIntegrationRuntime"),
	// 	Properties: &armdatafactory.ManagedIntegrationRuntimeStatus{
	// 		Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeManaged),
	// 		DataFactoryName: to.Ptr("exampleFactoryName"),
	// 		State: to.Ptr(armdatafactory.IntegrationRuntimeStateStarted),
	// 		TypeProperties: &armdatafactory.ManagedIntegrationRuntimeStatusTypeProperties{
	// 			CreateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-13T21:11:01.8695494Z"); return t}()),
	// 			Nodes: []*armdatafactory.ManagedIntegrationRuntimeNode{
	// 			},
	// 			OtherErrors: []*armdatafactory.ManagedIntegrationRuntimeError{
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_Stop.json
func ExampleIntegrationRuntimesClient_BeginStop() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginStop(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleManagedIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_SyncCredentials.json
func ExampleIntegrationRuntimesClient_SyncCredentials() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.SyncCredentials(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_GetMonitoringData.json
func ExampleIntegrationRuntimesClient_GetMonitoringData() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetMonitoringData(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeMonitoringData = armdatafactory.IntegrationRuntimeMonitoringData{
	// 	Name: to.Ptr("exampleIntegrationRuntime"),
	// 	Nodes: []*armdatafactory.IntegrationRuntimeNodeMonitoringData{
	// 		{
	// 			AvailableMemoryInMB: to.Ptr[int32](16740),
	// 			ConcurrentJobsLimit: to.Ptr[int32](28),
	// 			ConcurrentJobsRunning: to.Ptr[int32](0),
	// 			CPUUtilization: to.Ptr[int32](15),
	// 			NodeName: to.Ptr("Node_1"),
	// 			ReceivedBytes: to.Ptr[float32](6.731423377990723),
	// 			SentBytes: to.Ptr[float32](2.647491693496704),
	// 	}},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_Upgrade.json
func ExampleIntegrationRuntimesClient_Upgrade() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Upgrade(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_RemoveLinks.json
func ExampleIntegrationRuntimesClient_RemoveLinks() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.RemoveLinks(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", armdatafactory.LinkedIntegrationRuntimeRequest{
		LinkedFactoryName: to.Ptr("exampleFactoryName-linked"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/IntegrationRuntimes_CreateLinkedIntegrationRuntime.json
func ExampleIntegrationRuntimesClient_CreateLinkedIntegrationRuntime() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewIntegrationRuntimesClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateLinkedIntegrationRuntime(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleIntegrationRuntime", armdatafactory.CreateLinkedIntegrationRuntimeRequest{
		Name:                to.Ptr("bfa92911-9fb6-4fbe-8f23-beae87bc1c83"),
		DataFactoryLocation: to.Ptr("West US"),
		DataFactoryName:     to.Ptr("e9955d6d-56ea-4be3-841c-52a12c1a9981"),
		SubscriptionID:      to.Ptr("061774c7-4b5a-4159-a55b-365581830283"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IntegrationRuntimeStatusResponse = armdatafactory.IntegrationRuntimeStatusResponse{
	// 	Name: to.Ptr("exampleIntegrationRuntime"),
	// 	Properties: &armdatafactory.SelfHostedIntegrationRuntimeStatus{
	// 		Type: to.Ptr(armdatafactory.IntegrationRuntimeTypeSelfHosted),
	// 		DataFactoryName: to.Ptr("exampleFactoryName"),
	// 		State: to.Ptr(armdatafactory.IntegrationRuntimeStateOnline),
	// 		TypeProperties: &armdatafactory.SelfHostedIntegrationRuntimeStatusTypeProperties{
	// 			AutoUpdate: to.Ptr(armdatafactory.IntegrationRuntimeAutoUpdateOn),
	// 			AutoUpdateETA: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-08-20T19:00:00Z"); return t}()),
	// 			Capabilities: map[string]*string{
	// 				"connectedToResourceManager": to.Ptr("True"),
	// 				"credentialInSync": to.Ptr("True"),
	// 				"httpsPortEnabled": to.Ptr("True"),
	// 				"nodeEnabled": to.Ptr("True"),
	// 				"serviceBusConnected": to.Ptr("True"),
	// 			},
	// 			CreateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-08-17T03:43:25.7055573Z"); return t}()),
	// 			LatestVersion: to.Ptr("3.9.6774.1"),
	// 			Links: []*armdatafactory.LinkedIntegrationRuntime{
	// 				{
	// 					Name: to.Ptr("bfa92911-9fb6-4fbe-8f23-beae87bc1c83"),
	// 					CreateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-08-17T06:31:04.0617928Z"); return t}()),
	// 					DataFactoryLocation: to.Ptr("West US"),
	// 					DataFactoryName: to.Ptr("e9955d6d-56ea-4be3-841c-52a12c1a9981"),
	// 					SubscriptionID: to.Ptr("061774c7-4b5a-4159-a55b-365581830283"),
	// 			}},
	// 			LocalTimeZoneOffset: to.Ptr("PT8H"),
	// 			Nodes: []*armdatafactory.SelfHostedIntegrationRuntimeNode{
	// 				{
	// 					Capabilities: map[string]*string{
	// 						"connectedToResourceManager": to.Ptr("True"),
	// 						"credentialInSync": to.Ptr("True"),
	// 						"httpsPortEnabled": to.Ptr("True"),
	// 						"nodeEnabled": to.Ptr("True"),
	// 						"serviceBusConnected": to.Ptr("True"),
	// 					},
	// 					HostServiceURI: to.Ptr("https://yanzhang-dt.fareast.corp.microsoft.com:8050/HostServiceRemote.svc/"),
	// 					IsActiveDispatcher: to.Ptr(true),
	// 					LastConnectTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-08-17T06:30:46.6262976Z"); return t}()),
	// 					LastStartTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-08-17T03:45:30.8499851Z"); return t}()),
	// 					LastUpdateResult: to.Ptr(armdatafactory.IntegrationRuntimeUpdateResultNone),
	// 					MachineName: to.Ptr("YANZHANG-DT"),
	// 					MaxConcurrentJobs: to.Ptr[int32](20),
	// 					NodeName: to.Ptr("Node_1"),
	// 					RegisterTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-08-17T03:44:55.8012825Z"); return t}()),
	// 					Status: to.Ptr(armdatafactory.SelfHostedIntegrationRuntimeNodeStatusOnline),
	// 					Version: to.Ptr("3.8.6743.6"),
	// 					VersionStatus: to.Ptr("UpToDate"),
	// 			}},
	// 			PushedVersion: to.Ptr("3.9.6774.1"),
	// 			ScheduledUpdateDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-08-20T00:00:00Z"); return t}()),
	// 			ServiceUrls: []*string{
	// 				to.Ptr("wu.frontend.int.clouddatahub-int.net"),
	// 				to.Ptr("*.servicebus.windows.net")},
	// 				TaskQueueID: to.Ptr("823da112-f2d9-426b-a0d8-5f361b94f72a"),
	// 				UpdateDelayOffset: to.Ptr("PT19H"),
	// 				Version: to.Ptr("3.8.6743.6"),
	// 				VersionStatus: to.Ptr("UpdateAvailable"),
	// 			},
	// 		},
	// 	}
}
