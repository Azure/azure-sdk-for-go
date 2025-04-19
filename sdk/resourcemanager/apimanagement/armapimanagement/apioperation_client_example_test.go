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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListApiOperations.json
func ExampleAPIOperationClient_NewListByAPIPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAPIOperationClient().NewListByAPIPager("rg1", "apimService1", "57d2ef278aa04f0888cba3f3", &armapimanagement.APIOperationClientListByAPIOptions{Filter: nil,
		Top:  nil,
		Skip: nil,
		Tags: nil,
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
		// page.OperationCollection = armapimanagement.OperationCollection{
		// 	Count: to.Ptr[int64](5),
		// 	Value: []*armapimanagement.OperationContract{
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cdc"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdc"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("CancelOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/CancelOrder"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cda"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cda"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("GetMostRecentOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/GetMostRecentOrder"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cd9"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cd9"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("GetOpenOrders"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/GetOpenOrders"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cdb"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdb"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("GetOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/GetOrder"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cd8"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cd8"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("submitOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/submitOrder"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementHeadApiOperation.json
func ExampleAPIOperationClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAPIOperationClient().GetEntityTag(ctx, "rg1", "apimService1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetApiOperation.json
func ExampleAPIOperationClient_Get_apiManagementGetApiOperation() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIOperationClient().Get(ctx, "rg1", "apimService1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OperationContract = armapimanagement.OperationContract{
	// 	Name: to.Ptr("57d2ef278aa04f0ad01d6cdc"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdc"),
	// 	Properties: &armapimanagement.OperationContractProperties{
	// 		TemplateParameters: []*armapimanagement.ParameterContract{
	// 		},
	// 		Request: &armapimanagement.RequestContract{
	// 			Description: to.Ptr("IFazioService_CancelOrder_InputMessage"),
	// 			Headers: []*armapimanagement.ParameterContract{
	// 			},
	// 			QueryParameters: []*armapimanagement.ParameterContract{
	// 			},
	// 			Representations: []*armapimanagement.RepresentationContract{
	// 				{
	// 					ContentType: to.Ptr("text/xml"),
	// 					SchemaID: to.Ptr("6980a395-f08b-4a59-8295-1440cbd909b8"),
	// 					TypeName: to.Ptr("CancelOrder"),
	// 			}},
	// 		},
	// 		Responses: []*armapimanagement.ResponseContract{
	// 			{
	// 				Description: to.Ptr("IFazioService_CancelOrder_OutputMessage"),
	// 				Headers: []*armapimanagement.ParameterContract{
	// 				},
	// 				Representations: []*armapimanagement.RepresentationContract{
	// 					{
	// 						ContentType: to.Ptr("text/xml"),
	// 						SchemaID: to.Ptr("6980a395-f08b-4a59-8295-1440cbd909b8"),
	// 						TypeName: to.Ptr("CancelOrderResponse"),
	// 				}},
	// 				StatusCode: to.Ptr[int32](200),
	// 		}},
	// 		Method: to.Ptr("POST"),
	// 		DisplayName: to.Ptr("CancelOrder"),
	// 		URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/CancelOrder"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetApiOperationPetStore.json
func ExampleAPIOperationClient_Get_apiManagementGetApiOperationPetStore() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIOperationClient().Get(ctx, "rg1", "apimService1", "swagger-petstore", "loginUser", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OperationContract = armapimanagement.OperationContract{
	// 	Name: to.Ptr("loginUser"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/swagger-petstore/operations/loginUser"),
	// 	Properties: &armapimanagement.OperationContractProperties{
	// 		Description: to.Ptr(""),
	// 		TemplateParameters: []*armapimanagement.ParameterContract{
	// 			{
	// 				Name: to.Ptr("username"),
	// 				Type: to.Ptr("string"),
	// 				Description: to.Ptr("The user name for login"),
	// 				Required: to.Ptr(true),
	// 				Values: []*string{
	// 				},
	// 			},
	// 			{
	// 				Name: to.Ptr("password"),
	// 				Type: to.Ptr("string"),
	// 				Description: to.Ptr("The password for login in clear text"),
	// 				Required: to.Ptr(true),
	// 				Values: []*string{
	// 				},
	// 		}},
	// 		Request: &armapimanagement.RequestContract{
	// 			Headers: []*armapimanagement.ParameterContract{
	// 			},
	// 			QueryParameters: []*armapimanagement.ParameterContract{
	// 			},
	// 			Representations: []*armapimanagement.RepresentationContract{
	// 			},
	// 		},
	// 		Responses: []*armapimanagement.ResponseContract{
	// 			{
	// 				Description: to.Ptr("successful operation"),
	// 				Headers: []*armapimanagement.ParameterContract{
	// 					{
	// 						Name: to.Ptr("X-Rate-Limit"),
	// 						Type: to.Ptr("integer"),
	// 						Description: to.Ptr("calls per hour allowed by the user"),
	// 						Values: []*string{
	// 						},
	// 					},
	// 					{
	// 						Name: to.Ptr("X-Expires-After"),
	// 						Type: to.Ptr("string"),
	// 						Description: to.Ptr("date in UTC when token expires"),
	// 						Values: []*string{
	// 						},
	// 				}},
	// 				Representations: []*armapimanagement.RepresentationContract{
	// 					{
	// 						ContentType: to.Ptr("application/xml"),
	// 						SchemaID: to.Ptr("5ba91a35f373b513a0bf31c6"),
	// 						TypeName: to.Ptr("UserLoginGet200ApplicationXmlResponse"),
	// 					},
	// 					{
	// 						ContentType: to.Ptr("application/json"),
	// 						SchemaID: to.Ptr("5ba91a35f373b513a0bf31c6"),
	// 						TypeName: to.Ptr("UserLoginGet200ApplicationJsonResponse"),
	// 				}},
	// 				StatusCode: to.Ptr[int32](200),
	// 			},
	// 			{
	// 				Description: to.Ptr("Invalid username/password supplied"),
	// 				Headers: []*armapimanagement.ParameterContract{
	// 				},
	// 				Representations: []*armapimanagement.RepresentationContract{
	// 					{
	// 						ContentType: to.Ptr("application/xml"),
	// 					},
	// 					{
	// 						ContentType: to.Ptr("application/json"),
	// 				}},
	// 				StatusCode: to.Ptr[int32](400),
	// 		}},
	// 		Method: to.Ptr("GET"),
	// 		DisplayName: to.Ptr("Logs user into the system"),
	// 		URLTemplate: to.Ptr("/user/login?username={username}&password={password}"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateApiOperation.json
func ExampleAPIOperationClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIOperationClient().CreateOrUpdate(ctx, "rg1", "apimService1", "PetStoreTemplate2", "newoperations", armapimanagement.OperationContract{
		Properties: &armapimanagement.OperationContractProperties{
			Description:        to.Ptr("This can only be done by the logged in user."),
			TemplateParameters: []*armapimanagement.ParameterContract{},
			Request: &armapimanagement.RequestContract{
				Description:     to.Ptr("Created user object"),
				Headers:         []*armapimanagement.ParameterContract{},
				QueryParameters: []*armapimanagement.ParameterContract{},
				Representations: []*armapimanagement.RepresentationContract{
					{
						ContentType: to.Ptr("application/json"),
						SchemaID:    to.Ptr("592f6c1d0af5840ca8897f0c"),
						TypeName:    to.Ptr("User"),
					}},
			},
			Responses: []*armapimanagement.ResponseContract{
				{
					Description: to.Ptr("successful operation"),
					Headers:     []*armapimanagement.ParameterContract{},
					Representations: []*armapimanagement.RepresentationContract{
						{
							ContentType: to.Ptr("application/xml"),
						},
						{
							ContentType: to.Ptr("application/json"),
						}},
					StatusCode: to.Ptr[int32](200),
				}},
			Method:      to.Ptr("POST"),
			DisplayName: to.Ptr("createUser2"),
			URLTemplate: to.Ptr("/user1"),
		},
	}, &armapimanagement.APIOperationClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OperationContract = armapimanagement.OperationContract{
	// 	Name: to.Ptr("newoperations"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/PetStoreTemplate2/operations/newoperations"),
	// 	Properties: &armapimanagement.OperationContractProperties{
	// 		Description: to.Ptr("This can only be done by the logged in user."),
	// 		TemplateParameters: []*armapimanagement.ParameterContract{
	// 		},
	// 		Request: &armapimanagement.RequestContract{
	// 			Description: to.Ptr("Created user object"),
	// 			Headers: []*armapimanagement.ParameterContract{
	// 			},
	// 			QueryParameters: []*armapimanagement.ParameterContract{
	// 			},
	// 			Representations: []*armapimanagement.RepresentationContract{
	// 				{
	// 					ContentType: to.Ptr("application/json"),
	// 					SchemaID: to.Ptr("592f6c1d0af5840ca8897f0c"),
	// 					TypeName: to.Ptr("User"),
	// 			}},
	// 		},
	// 		Responses: []*armapimanagement.ResponseContract{
	// 			{
	// 				Description: to.Ptr("successful operation"),
	// 				Headers: []*armapimanagement.ParameterContract{
	// 				},
	// 				Representations: []*armapimanagement.RepresentationContract{
	// 					{
	// 						ContentType: to.Ptr("application/xml"),
	// 					},
	// 					{
	// 						ContentType: to.Ptr("application/json"),
	// 				}},
	// 				StatusCode: to.Ptr[int32](200),
	// 		}},
	// 		Method: to.Ptr("POST"),
	// 		DisplayName: to.Ptr("createUser2"),
	// 		URLTemplate: to.Ptr("/user1"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementUpdateApiOperation.json
func ExampleAPIOperationClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIOperationClient().Update(ctx, "rg1", "apimService1", "echo-api", "operationId", "*", armapimanagement.OperationUpdateContract{
		Properties: &armapimanagement.OperationUpdateContractProperties{
			TemplateParameters: []*armapimanagement.ParameterContract{},
			Request: &armapimanagement.RequestContract{
				QueryParameters: []*armapimanagement.ParameterContract{
					{
						Name:         to.Ptr("param1"),
						Type:         to.Ptr("string"),
						Description:  to.Ptr("A sample parameter that is required and has a default value of \"sample\"."),
						DefaultValue: to.Ptr("sample"),
						Required:     to.Ptr(true),
						Values: []*string{
							to.Ptr("sample")},
					}},
			},
			Responses: []*armapimanagement.ResponseContract{
				{
					Description:     to.Ptr("Returned in all cases."),
					Headers:         []*armapimanagement.ParameterContract{},
					Representations: []*armapimanagement.RepresentationContract{},
					StatusCode:      to.Ptr[int32](200),
				},
				{
					Description:     to.Ptr("Server Error."),
					Headers:         []*armapimanagement.ParameterContract{},
					Representations: []*armapimanagement.RepresentationContract{},
					StatusCode:      to.Ptr[int32](500),
				}},
			Method:      to.Ptr("GET"),
			DisplayName: to.Ptr("Retrieve resource"),
			URLTemplate: to.Ptr("/resource"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OperationContract = armapimanagement.OperationContract{
	// 	Name: to.Ptr("57d2ef278aa04f0ad01d6cdc"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/operations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdc"),
	// 	Properties: &armapimanagement.OperationContractProperties{
	// 		TemplateParameters: []*armapimanagement.ParameterContract{
	// 		},
	// 		Request: &armapimanagement.RequestContract{
	// 			Description: to.Ptr("IFazioService_CancelOrder_InputMessage"),
	// 			Headers: []*armapimanagement.ParameterContract{
	// 			},
	// 			QueryParameters: []*armapimanagement.ParameterContract{
	// 			},
	// 			Representations: []*armapimanagement.RepresentationContract{
	// 				{
	// 					ContentType: to.Ptr("text/xml"),
	// 					SchemaID: to.Ptr("6980a395-f08b-4a59-8295-1440cbd909b8"),
	// 					TypeName: to.Ptr("CancelOrder"),
	// 			}},
	// 		},
	// 		Responses: []*armapimanagement.ResponseContract{
	// 			{
	// 				Description: to.Ptr("IFazioService_CancelOrder_OutputMessage"),
	// 				Headers: []*armapimanagement.ParameterContract{
	// 				},
	// 				Representations: []*armapimanagement.RepresentationContract{
	// 					{
	// 						ContentType: to.Ptr("text/xml"),
	// 						SchemaID: to.Ptr("6980a395-f08b-4a59-8295-1440cbd909b8"),
	// 						TypeName: to.Ptr("CancelOrderResponse"),
	// 				}},
	// 				StatusCode: to.Ptr[int32](200),
	// 		}},
	// 		Method: to.Ptr("POST"),
	// 		DisplayName: to.Ptr("CancelOrder"),
	// 		URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/CancelOrder"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeleteApiOperation.json
func ExampleAPIOperationClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAPIOperationClient().Delete(ctx, "rg1", "apimService1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", "*", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
