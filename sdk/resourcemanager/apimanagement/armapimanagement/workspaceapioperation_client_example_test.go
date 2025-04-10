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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListWorkspaceApiOperations.json
func ExampleWorkspaceAPIOperationClient_NewListByAPIPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewWorkspaceAPIOperationClient().NewListByAPIPager("rg1", "apimService1", "wks1", "57d2ef278aa04f0888cba3f3", &armapimanagement.WorkspaceAPIOperationClientListByAPIOptions{Filter: nil,
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
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdc"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("CancelOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/CancelOrder"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cda"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cda"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("GetMostRecentOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/GetMostRecentOrder"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cd9"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cd9"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("GetOpenOrders"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/GetOpenOrders"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cdb"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdb"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("GetOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/GetOrder"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cd8"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cd8"),
		// 			Properties: &armapimanagement.OperationContractProperties{
		// 				Method: to.Ptr("POST"),
		// 				DisplayName: to.Ptr("submitOrder"),
		// 				URLTemplate: to.Ptr("/?soapAction=http://tempuri.org/IFazioService/submitOrder"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementHeadWorkspaceApiOperation.json
func ExampleWorkspaceAPIOperationClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewWorkspaceAPIOperationClient().GetEntityTag(ctx, "rg1", "apimService1", "wks1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetWorkspaceApiOperation.json
func ExampleWorkspaceAPIOperationClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewWorkspaceAPIOperationClient().Get(ctx, "rg1", "apimService1", "wks1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OperationContract = armapimanagement.OperationContract{
	// 	Name: to.Ptr("57d2ef278aa04f0ad01d6cdc"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdc"),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateWorkspaceApiOperation.json
func ExampleWorkspaceAPIOperationClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewWorkspaceAPIOperationClient().CreateOrUpdate(ctx, "rg1", "apimService1", "wks1", "PetStoreTemplate2", "newoperations", armapimanagement.OperationContract{
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
	}, &armapimanagement.WorkspaceAPIOperationClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OperationContract = armapimanagement.OperationContract{
	// 	Name: to.Ptr("newoperations"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/PetStoreTemplate2/operations/newoperations"),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementUpdateWorkspaceApiOperation.json
func ExampleWorkspaceAPIOperationClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewWorkspaceAPIOperationClient().Update(ctx, "rg1", "apimService1", "wks1", "echo-api", "operationId", "*", armapimanagement.OperationUpdateContract{
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
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/workspaces/apis/operations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/workspaces/wks1/apis/57d2ef278aa04f0888cba3f3/operations/57d2ef278aa04f0ad01d6cdc"),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeleteWorkspaceApiOperation.json
func ExampleWorkspaceAPIOperationClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewWorkspaceAPIOperationClient().Delete(ctx, "rg1", "apimService1", "wks1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", "*", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
