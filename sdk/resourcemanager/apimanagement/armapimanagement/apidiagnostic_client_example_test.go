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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListApiDiagnostics.json
func ExampleAPIDiagnosticClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAPIDiagnosticClient().NewListByServicePager("rg1", "apimService1", "echo-api", &armapimanagement.APIDiagnosticClientListByServiceOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
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
		// page.DiagnosticCollection = armapimanagement.DiagnosticCollection{
		// 	Count: to.Ptr[int64](1),
		// 	Value: []*armapimanagement.DiagnosticContract{
		// 		{
		// 			Name: to.Ptr("applicationinsights"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/apis/diagnostics"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/echo-api/diagnostics/applicationinsights"),
		// 			Properties: &armapimanagement.DiagnosticContractProperties{
		// 				AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
		// 				Backend: &armapimanagement.PipelineDiagnosticSettings{
		// 					Response: &armapimanagement.HTTPMessageDiagnostic{
		// 						Body: &armapimanagement.BodyDiagnosticSettings{
		// 							Bytes: to.Ptr[int32](100),
		// 						},
		// 						Headers: []*string{
		// 						},
		// 					},
		// 					Request: &armapimanagement.HTTPMessageDiagnostic{
		// 						Body: &armapimanagement.BodyDiagnosticSettings{
		// 							Bytes: to.Ptr[int32](100),
		// 						},
		// 						Headers: []*string{
		// 						},
		// 					},
		// 				},
		// 				Frontend: &armapimanagement.PipelineDiagnosticSettings{
		// 					Response: &armapimanagement.HTTPMessageDiagnostic{
		// 						Body: &armapimanagement.BodyDiagnosticSettings{
		// 							Bytes: to.Ptr[int32](100),
		// 						},
		// 						Headers: []*string{
		// 						},
		// 					},
		// 					Request: &armapimanagement.HTTPMessageDiagnostic{
		// 						Body: &armapimanagement.BodyDiagnosticSettings{
		// 							Bytes: to.Ptr[int32](100),
		// 						},
		// 						Headers: []*string{
		// 						},
		// 					},
		// 				},
		// 				HTTPCorrelationProtocol: to.Ptr(armapimanagement.HTTPCorrelationProtocolLegacy),
		// 				LogClientIP: to.Ptr(true),
		// 				LoggerID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/loggers/aisamplingtest"),
		// 				Sampling: &armapimanagement.SamplingSettings{
		// 					Percentage: to.Ptr[float64](100),
		// 					SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
		// 				},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementHeadApiDiagnostic.json
func ExampleAPIDiagnosticClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAPIDiagnosticClient().GetEntityTag(ctx, "rg1", "apimService1", "57d1f7558aa04f15146d9d8a", "applicationinsights", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetApiDiagnostic.json
func ExampleAPIDiagnosticClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIDiagnosticClient().Get(ctx, "rg1", "apimService1", "57d1f7558aa04f15146d9d8a", "applicationinsights", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DiagnosticContract = armapimanagement.DiagnosticContract{
	// 	Name: to.Ptr("applicationinsights"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/diagnostics"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/echo-api/diagnostics/applicationinsights"),
	// 	Properties: &armapimanagement.DiagnosticContractProperties{
	// 		AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
	// 		Backend: &armapimanagement.PipelineDiagnosticSettings{
	// 			Response: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 			Request: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 		},
	// 		Frontend: &armapimanagement.PipelineDiagnosticSettings{
	// 			Response: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 			Request: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 		},
	// 		HTTPCorrelationProtocol: to.Ptr(armapimanagement.HTTPCorrelationProtocolLegacy),
	// 		LogClientIP: to.Ptr(true),
	// 		LoggerID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/loggers/aisamplingtest"),
	// 		Sampling: &armapimanagement.SamplingSettings{
	// 			Percentage: to.Ptr[float64](100),
	// 			SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateApiDiagnostic.json
func ExampleAPIDiagnosticClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIDiagnosticClient().CreateOrUpdate(ctx, "rg1", "apimService1", "57d1f7558aa04f15146d9d8a", "applicationinsights", armapimanagement.DiagnosticContract{
		Properties: &armapimanagement.DiagnosticContractProperties{
			AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
			Backend: &armapimanagement.PipelineDiagnosticSettings{
				Response: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
				Request: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
			},
			Frontend: &armapimanagement.PipelineDiagnosticSettings{
				Response: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
				Request: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
			},
			LoggerID: to.Ptr("/loggers/applicationinsights"),
			Sampling: &armapimanagement.SamplingSettings{
				Percentage:   to.Ptr[float64](50),
				SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
			},
		},
	}, &armapimanagement.APIDiagnosticClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DiagnosticContract = armapimanagement.DiagnosticContract{
	// 	Name: to.Ptr("applicationinsights"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/diagnostics"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d1f7558aa04f15146d9d8a/diagnostics/applicationinsights"),
	// 	Properties: &armapimanagement.DiagnosticContractProperties{
	// 		AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
	// 		Backend: &armapimanagement.PipelineDiagnosticSettings{
	// 			Response: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](512),
	// 				},
	// 				Headers: []*string{
	// 					to.Ptr("Content-type")},
	// 				},
	// 				Request: &armapimanagement.HTTPMessageDiagnostic{
	// 					Body: &armapimanagement.BodyDiagnosticSettings{
	// 						Bytes: to.Ptr[int32](512),
	// 					},
	// 					Headers: []*string{
	// 						to.Ptr("Content-type")},
	// 					},
	// 				},
	// 				Frontend: &armapimanagement.PipelineDiagnosticSettings{
	// 					Response: &armapimanagement.HTTPMessageDiagnostic{
	// 						Body: &armapimanagement.BodyDiagnosticSettings{
	// 							Bytes: to.Ptr[int32](512),
	// 						},
	// 						Headers: []*string{
	// 							to.Ptr("Content-type")},
	// 						},
	// 						Request: &armapimanagement.HTTPMessageDiagnostic{
	// 							Body: &armapimanagement.BodyDiagnosticSettings{
	// 								Bytes: to.Ptr[int32](512),
	// 							},
	// 							Headers: []*string{
	// 								to.Ptr("Content-type")},
	// 							},
	// 						},
	// 						LoggerID: to.Ptr("/loggers/applicationinsights"),
	// 						Sampling: &armapimanagement.SamplingSettings{
	// 							Percentage: to.Ptr[float64](50),
	// 							SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
	// 						},
	// 					},
	// 				}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementUpdateApiDiagnostic.json
func ExampleAPIDiagnosticClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIDiagnosticClient().Update(ctx, "rg1", "apimService1", "echo-api", "applicationinsights", "*", armapimanagement.DiagnosticContract{
		Properties: &armapimanagement.DiagnosticContractProperties{
			AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
			Backend: &armapimanagement.PipelineDiagnosticSettings{
				Response: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
				Request: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
			},
			Frontend: &armapimanagement.PipelineDiagnosticSettings{
				Response: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
				Request: &armapimanagement.HTTPMessageDiagnostic{
					Body: &armapimanagement.BodyDiagnosticSettings{
						Bytes: to.Ptr[int32](512),
					},
					Headers: []*string{
						to.Ptr("Content-type")},
				},
			},
			LoggerID: to.Ptr("/loggers/applicationinsights"),
			Sampling: &armapimanagement.SamplingSettings{
				Percentage:   to.Ptr[float64](50),
				SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DiagnosticContract = armapimanagement.DiagnosticContract{
	// 	Name: to.Ptr("applicationinsights"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/diagnostics"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/echo-api/diagnostics/applicationinsights"),
	// 	Properties: &armapimanagement.DiagnosticContractProperties{
	// 		AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
	// 		Backend: &armapimanagement.PipelineDiagnosticSettings{
	// 			Response: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 			Request: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 		},
	// 		Frontend: &armapimanagement.PipelineDiagnosticSettings{
	// 			Response: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 			Request: &armapimanagement.HTTPMessageDiagnostic{
	// 				Body: &armapimanagement.BodyDiagnosticSettings{
	// 					Bytes: to.Ptr[int32](100),
	// 				},
	// 				Headers: []*string{
	// 				},
	// 			},
	// 		},
	// 		HTTPCorrelationProtocol: to.Ptr(armapimanagement.HTTPCorrelationProtocolLegacy),
	// 		LogClientIP: to.Ptr(true),
	// 		LoggerID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/loggers/aisamplingtest"),
	// 		Sampling: &armapimanagement.SamplingSettings{
	// 			Percentage: to.Ptr[float64](100),
	// 			SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeleteApiDiagnostic.json
func ExampleAPIDiagnosticClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAPIDiagnosticClient().Delete(ctx, "rg1", "apimService1", "57d1f7558aa04f15146d9d8a", "applicationinsights", "*", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
