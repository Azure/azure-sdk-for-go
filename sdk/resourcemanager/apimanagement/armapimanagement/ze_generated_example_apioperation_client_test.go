//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementHeadApiOperation.json
func ExampleAPIOperationClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIOperationClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.GetEntityTag(ctx,
		"rg1",
		"apimService1",
		"57d2ef278aa04f0888cba3f3",
		"57d2ef278aa04f0ad01d6cdc",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementGetApiOperation.json
func ExampleAPIOperationClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIOperationClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"rg1",
		"apimService1",
		"57d2ef278aa04f0888cba3f3",
		"57d2ef278aa04f0ad01d6cdc",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementCreateApiOperation.json
func ExampleAPIOperationClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIOperationClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx,
		"rg1",
		"apimService1",
		"PetStoreTemplate2",
		"newoperations",
		armapimanagement.OperationContract{
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
		},
		&armapimanagement.APIOperationClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementUpdateApiOperation.json
func ExampleAPIOperationClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIOperationClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"rg1",
		"apimService1",
		"echo-api",
		"operationId",
		"*",
		armapimanagement.OperationUpdateContract{
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
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementDeleteApiOperation.json
func ExampleAPIOperationClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIOperationClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"rg1",
		"apimService1",
		"57d2ef278aa04f0888cba3f3",
		"57d2ef278aa04f0ad01d6cdc",
		"*",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
