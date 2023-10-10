//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpolicy_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/deleteVariable.json
func ExampleVariablesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewVariablesClient().Delete(ctx, "DemoTestVariable", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/createOrUpdateVariable.json
func ExampleVariablesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVariablesClient().CreateOrUpdate(ctx, "DemoTestVariable", armpolicy.Variable{
		Properties: &armpolicy.VariableProperties{
			Columns: []*armpolicy.VariableColumn{
				{
					ColumnName: to.Ptr("TestColumn"),
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Variable = armpolicy.Variable{
	// 	Name: to.Ptr("DemoTestVariable"),
	// 	Type: to.Ptr("Microsoft.Authorization/variables"),
	// 	ID: to.Ptr("/subscriptions/ae640e6b-ba3e-4256-9d62-2993eecfa6f2/providers/Microsoft.Authorization/variables/DemoTestVariable"),
	// 	Properties: &armpolicy.VariableProperties{
	// 		Columns: []*armpolicy.VariableColumn{
	// 			{
	// 				ColumnName: to.Ptr("TestColumn"),
	// 		}},
	// 	},
	// 	SystemData: &armpolicy.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T01:01:01.1075056Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T02:01:01.1075056Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/getVariable.json
func ExampleVariablesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVariablesClient().Get(ctx, "DemoTestVariable", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Variable = armpolicy.Variable{
	// 	Name: to.Ptr("DemoTestVariable"),
	// 	Type: to.Ptr("Microsoft.Authorization/variables"),
	// 	ID: to.Ptr("/subscriptions/ae640e6b-ba3e-4256-9d62-2993eecfa6f2/providers/Microsoft.Authorization/variables/DemoTestVariable"),
	// 	Properties: &armpolicy.VariableProperties{
	// 		Columns: []*armpolicy.VariableColumn{
	// 			{
	// 				ColumnName: to.Ptr("TestColumn"),
	// 		}},
	// 	},
	// 	SystemData: &armpolicy.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T01:01:01.1075056Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-08-01T01:01:01.1075056Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/deleteVariableAtManagementGroup.json
func ExampleVariablesClient_DeleteAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewVariablesClient().DeleteAtManagementGroup(ctx, "DevOrg", "DemoTestVariable", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/createOrUpdateVariableAtManagementGroup.json
func ExampleVariablesClient_CreateOrUpdateAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVariablesClient().CreateOrUpdateAtManagementGroup(ctx, "DevOrg", "DemoTestVariable", armpolicy.Variable{
		Properties: &armpolicy.VariableProperties{
			Columns: []*armpolicy.VariableColumn{
				{
					ColumnName: to.Ptr("TestColumn"),
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Variable = armpolicy.Variable{
	// 	Name: to.Ptr("DemoTestVariable"),
	// 	Type: to.Ptr("Microsoft.Authorization/variables"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/DevOrg/providers/Microsoft.Authorization/variables/DemoTestVariable"),
	// 	Properties: &armpolicy.VariableProperties{
	// 		Columns: []*armpolicy.VariableColumn{
	// 			{
	// 				ColumnName: to.Ptr("TestColumn"),
	// 		}},
	// 	},
	// 	SystemData: &armpolicy.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T01:01:01.1075056Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T02:01:01.1075056Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/getVariableAtManagementGroup.json
func ExampleVariablesClient_GetAtManagementGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVariablesClient().GetAtManagementGroup(ctx, "DevOrg", "DemoTestVariable", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Variable = armpolicy.Variable{
	// 	Name: to.Ptr("DemoTestVariable"),
	// 	Type: to.Ptr("Microsoft.Authorization/variables"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/DevOrg/providers/Microsoft.Authorization/variables/DemoTestVariable"),
	// 	Properties: &armpolicy.VariableProperties{
	// 		Columns: []*armpolicy.VariableColumn{
	// 			{
	// 				ColumnName: to.Ptr("TestColumn"),
	// 		}},
	// 	},
	// 	SystemData: &armpolicy.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T01:01:01.1075056Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-08-01T01:01:01.1075056Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/listVariablesForSubscription.json
func ExampleVariablesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVariablesClient().NewListPager(nil)
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
		// page.VariableListResult = armpolicy.VariableListResult{
		// 	Value: []*armpolicy.Variable{
		// 		{
		// 			Name: to.Ptr("DemoTestVariable"),
		// 			Type: to.Ptr("Microsoft.Authorization/variables"),
		// 			ID: to.Ptr("/subscriptions/ae640e6b-ba3e-4256-9d62-2993eecfa6f2/providers/Microsoft.Authorization/variables/DemoTestVariable"),
		// 			Properties: &armpolicy.VariableProperties{
		// 				Columns: []*armpolicy.VariableColumn{
		// 					{
		// 						ColumnName: to.Ptr("TestColumn"),
		// 				}},
		// 			},
		// 			SystemData: &armpolicy.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-08-01T01:01:01.1075056Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("NetworkingVariable"),
		// 			Type: to.Ptr("Microsoft.Authorization/variables"),
		// 			ID: to.Ptr("/subscriptions/ae640e6b-ba3e-4256-9d62-2993eecfa6f2/providers/Microsoft.Authorization/variables/NetworkingVariable"),
		// 			Properties: &armpolicy.VariableProperties{
		// 				Columns: []*armpolicy.VariableColumn{
		// 					{
		// 						ColumnName: to.Ptr("NetworkResourceName"),
		// 				}},
		// 			},
		// 			SystemData: &armpolicy.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-08-01T01:01:01.1075056Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4e7b46ac4e63391b4d436265622f7b029a707d4f/specification/resources/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/listVariablesForManagementGroup.json
func ExampleVariablesClient_NewListForManagementGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpolicy.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVariablesClient().NewListForManagementGroupPager("DevOrg", nil)
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
		// page.VariableListResult = armpolicy.VariableListResult{
		// 	Value: []*armpolicy.Variable{
		// 		{
		// 			Name: to.Ptr("DemoTestVariable"),
		// 			Type: to.Ptr("Microsoft.Authorization/variables"),
		// 			ID: to.Ptr("/providers/Microsoft.Management/managementGroups/DevOrg/providers/Microsoft.Authorization/variables/DemoTestVariable"),
		// 			Properties: &armpolicy.VariableProperties{
		// 				Columns: []*armpolicy.VariableColumn{
		// 					{
		// 						ColumnName: to.Ptr("TestColumn"),
		// 				}},
		// 			},
		// 			SystemData: &armpolicy.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-07-01T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-08-01T01:01:01.1075056Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("NetworkingVariable"),
		// 			Type: to.Ptr("Microsoft.Authorization/variables"),
		// 			ID: to.Ptr("/providers/Microsoft.Management/managementGroups/DevOrg/providers/Microsoft.Authorization/variables/NetworkingVariable"),
		// 			Properties: &armpolicy.VariableProperties{
		// 				Columns: []*armpolicy.VariableColumn{
		// 					{
		// 						ColumnName: to.Ptr("NetworkResourceName"),
		// 				}},
		// 			},
		// 			SystemData: &armpolicy.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-07-01T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-07-01T02:01:01.1075056Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armpolicy.CreatedByTypeUser),
		// 			},
		// 	}},
		// }
	}
}
