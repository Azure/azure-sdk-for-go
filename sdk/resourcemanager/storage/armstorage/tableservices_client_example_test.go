//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstorage_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/TableServicesList.json
func ExampleTableServicesClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewTableServicesClient().List(ctx, "res9290", "sto1590", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ListTableServices = armstorage.ListTableServices{
	// 	Value: []*armstorage.TableServiceProperties{
	// 		{
	// 			Name: to.Ptr("default"),
	// 			Type: to.Ptr("Microsoft.Storage/storageAccounts/tableServices"),
	// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res4410/providers/Microsoft.Storage/storageAccounts/sto8607/tableServices/default"),
	// 			TableServiceProperties: &armstorage.TableServicePropertiesProperties{
	// 				Cors: &armstorage.CorsRules{
	// 					CorsRules: []*armstorage.CorsRule{
	// 						{
	// 							AllowedHeaders: []*string{
	// 								to.Ptr("x-ms-meta-abc"),
	// 								to.Ptr("x-ms-meta-data*"),
	// 								to.Ptr("x-ms-meta-target*")},
	// 								AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 									to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
	// 									to.Ptr(armstorage.CorsRuleAllowedMethodsItemHEAD),
	// 									to.Ptr(armstorage.CorsRuleAllowedMethodsItemPOST),
	// 									to.Ptr(armstorage.CorsRuleAllowedMethodsItemOPTIONS),
	// 									to.Ptr(armstorage.CorsRuleAllowedMethodsItemMERGE),
	// 									to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
	// 									AllowedOrigins: []*string{
	// 										to.Ptr("http://www.contoso.com"),
	// 										to.Ptr("http://www.fabrikam.com")},
	// 										ExposedHeaders: []*string{
	// 											to.Ptr("x-ms-meta-*")},
	// 											MaxAgeInSeconds: to.Ptr[int32](100),
	// 										},
	// 										{
	// 											AllowedHeaders: []*string{
	// 												to.Ptr("*")},
	// 												AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 													to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET)},
	// 													AllowedOrigins: []*string{
	// 														to.Ptr("*")},
	// 														ExposedHeaders: []*string{
	// 															to.Ptr("*")},
	// 															MaxAgeInSeconds: to.Ptr[int32](2),
	// 														},
	// 														{
	// 															AllowedHeaders: []*string{
	// 																to.Ptr("x-ms-meta-12345675754564*")},
	// 																AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 																	to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
	// 																	to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
	// 																	AllowedOrigins: []*string{
	// 																		to.Ptr("http://www.abc23.com"),
	// 																		to.Ptr("https://www.fabrikam.com/*")},
	// 																		ExposedHeaders: []*string{
	// 																			to.Ptr("x-ms-meta-abc"),
	// 																			to.Ptr("x-ms-meta-data*"),
	// 																			to.Ptr("x-ms-meta-target*")},
	// 																			MaxAgeInSeconds: to.Ptr[int32](2000),
	// 																	}},
	// 																},
	// 															},
	// 													}},
	// 												}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/TableServicesPut.json
func ExampleTableServicesClient_SetServiceProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewTableServicesClient().SetServiceProperties(ctx, "res4410", "sto8607", armstorage.TableServiceProperties{
		TableServiceProperties: &armstorage.TableServicePropertiesProperties{
			Cors: &armstorage.CorsRules{
				CorsRules: []*armstorage.CorsRule{
					{
						AllowedHeaders: []*string{
							to.Ptr("x-ms-meta-abc"),
							to.Ptr("x-ms-meta-data*"),
							to.Ptr("x-ms-meta-target*")},
						AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemHEAD),
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemPOST),
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemOPTIONS),
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemMERGE),
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
						AllowedOrigins: []*string{
							to.Ptr("http://www.contoso.com"),
							to.Ptr("http://www.fabrikam.com")},
						ExposedHeaders: []*string{
							to.Ptr("x-ms-meta-*")},
						MaxAgeInSeconds: to.Ptr[int32](100),
					},
					{
						AllowedHeaders: []*string{
							to.Ptr("*")},
						AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET)},
						AllowedOrigins: []*string{
							to.Ptr("*")},
						ExposedHeaders: []*string{
							to.Ptr("*")},
						MaxAgeInSeconds: to.Ptr[int32](2),
					},
					{
						AllowedHeaders: []*string{
							to.Ptr("x-ms-meta-12345675754564*")},
						AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
							to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
						AllowedOrigins: []*string{
							to.Ptr("http://www.abc23.com"),
							to.Ptr("https://www.fabrikam.com/*")},
						ExposedHeaders: []*string{
							to.Ptr("x-ms-meta-abc"),
							to.Ptr("x-ms-meta-data*"),
							to.Ptr("x-ms-meta-target*")},
						MaxAgeInSeconds: to.Ptr[int32](2000),
					}},
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.TableServiceProperties = armstorage.TableServiceProperties{
	// 	Name: to.Ptr("default"),
	// 	Type: to.Ptr("Microsoft.Storage/storageAccounts/tableServices"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res4410/providers/Microsoft.Storage/storageAccounts/sto8607/tableServices/default"),
	// 	TableServiceProperties: &armstorage.TableServicePropertiesProperties{
	// 		Cors: &armstorage.CorsRules{
	// 			CorsRules: []*armstorage.CorsRule{
	// 				{
	// 					AllowedHeaders: []*string{
	// 						to.Ptr("x-ms-meta-abc"),
	// 						to.Ptr("x-ms-meta-data*"),
	// 						to.Ptr("x-ms-meta-target*")},
	// 						AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemHEAD),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemPOST),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemOPTIONS),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemMERGE),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
	// 							AllowedOrigins: []*string{
	// 								to.Ptr("http://www.contoso.com"),
	// 								to.Ptr("http://www.fabrikam.com")},
	// 								ExposedHeaders: []*string{
	// 									to.Ptr("x-ms-meta-*")},
	// 									MaxAgeInSeconds: to.Ptr[int32](100),
	// 								},
	// 								{
	// 									AllowedHeaders: []*string{
	// 										to.Ptr("*")},
	// 										AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 											to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET)},
	// 											AllowedOrigins: []*string{
	// 												to.Ptr("*")},
	// 												ExposedHeaders: []*string{
	// 													to.Ptr("*")},
	// 													MaxAgeInSeconds: to.Ptr[int32](2),
	// 												},
	// 												{
	// 													AllowedHeaders: []*string{
	// 														to.Ptr("x-ms-meta-12345675754564*")},
	// 														AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 															to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
	// 															to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
	// 															AllowedOrigins: []*string{
	// 																to.Ptr("http://www.abc23.com"),
	// 																to.Ptr("https://www.fabrikam.com/*")},
	// 																ExposedHeaders: []*string{
	// 																	to.Ptr("x-ms-meta-abc"),
	// 																	to.Ptr("x-ms-meta-data*"),
	// 																	to.Ptr("x-ms-meta-target*")},
	// 																	MaxAgeInSeconds: to.Ptr[int32](2000),
	// 															}},
	// 														},
	// 													},
	// 												}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/97ee23a6db6078abcbec7b75bf9af8c503e9bb8b/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/TableServicesGet.json
func ExampleTableServicesClient_GetServiceProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewTableServicesClient().GetServiceProperties(ctx, "res4410", "sto8607", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.TableServiceProperties = armstorage.TableServiceProperties{
	// 	Name: to.Ptr("default"),
	// 	Type: to.Ptr("Microsoft.Storage/storageAccounts/tableServices"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res4410/providers/Microsoft.Storage/storageAccounts/sto8607/tableServices/default"),
	// 	TableServiceProperties: &armstorage.TableServicePropertiesProperties{
	// 		Cors: &armstorage.CorsRules{
	// 			CorsRules: []*armstorage.CorsRule{
	// 				{
	// 					AllowedHeaders: []*string{
	// 						to.Ptr("x-ms-meta-abc"),
	// 						to.Ptr("x-ms-meta-data*"),
	// 						to.Ptr("x-ms-meta-target*")},
	// 						AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemHEAD),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemPOST),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemOPTIONS),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemMERGE),
	// 							to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
	// 							AllowedOrigins: []*string{
	// 								to.Ptr("http://www.contoso.com"),
	// 								to.Ptr("http://www.fabrikam.com")},
	// 								ExposedHeaders: []*string{
	// 									to.Ptr("x-ms-meta-*")},
	// 									MaxAgeInSeconds: to.Ptr[int32](100),
	// 								},
	// 								{
	// 									AllowedHeaders: []*string{
	// 										to.Ptr("*")},
	// 										AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 											to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET)},
	// 											AllowedOrigins: []*string{
	// 												to.Ptr("*")},
	// 												ExposedHeaders: []*string{
	// 													to.Ptr("*")},
	// 													MaxAgeInSeconds: to.Ptr[int32](2),
	// 												},
	// 												{
	// 													AllowedHeaders: []*string{
	// 														to.Ptr("x-ms-meta-12345675754564*")},
	// 														AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
	// 															to.Ptr(armstorage.CorsRuleAllowedMethodsItemGET),
	// 															to.Ptr(armstorage.CorsRuleAllowedMethodsItemPUT)},
	// 															AllowedOrigins: []*string{
	// 																to.Ptr("http://www.abc23.com"),
	// 																to.Ptr("https://www.fabrikam.com/*")},
	// 																ExposedHeaders: []*string{
	// 																	to.Ptr("x-ms-meta-abc"),
	// 																	to.Ptr("x-ms-meta-data*"),
	// 																	to.Ptr("x-ms-meta-target*")},
	// 																	MaxAgeInSeconds: to.Ptr[int32](2000),
	// 															}},
	// 														},
	// 													},
	// 												}
}
