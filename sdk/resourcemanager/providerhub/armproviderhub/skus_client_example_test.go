//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armproviderhub_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/providerhub/armproviderhub/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_Get.json
func ExampleSKUsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().Get(ctx, "Microsoft.Contoso", "testResourceType", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// 	Properties: &armproviderhub.SKUResourceProperties{
	// 		SKUSettings: []*armproviderhub.SKUSetting{
	// 			{
	// 				Name: to.Ptr("freeSku"),
	// 				Kind: to.Ptr("Standard"),
	// 				Tier: to.Ptr("Tier1"),
	// 			},
	// 			{
	// 				Name: to.Ptr("premiumSku"),
	// 				Costs: []*armproviderhub.SKUCost{
	// 					{
	// 						MeterID: to.Ptr("xxx"),
	// 				}},
	// 				Kind: to.Ptr("Premium"),
	// 				Tier: to.Ptr("Tier2"),
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_CreateOrUpdate.json
func ExampleSKUsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().CreateOrUpdate(ctx, "Microsoft.Contoso", "testResourceType", "testSku", armproviderhub.SKUResource{
		Properties: &armproviderhub.SKUResourceProperties{
			SKUSettings: []*armproviderhub.SKUSetting{
				{
					Name: to.Ptr("freeSku"),
					Kind: to.Ptr("Standard"),
					Tier: to.Ptr("Tier1"),
				},
				{
					Name: to.Ptr("premiumSku"),
					Costs: []*armproviderhub.SKUCost{
						{
							MeterID: to.Ptr("xxx"),
						}},
					Kind: to.Ptr("Premium"),
					Tier: to.Ptr("Tier2"),
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_Delete.json
func ExampleSKUsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewSKUsClient().Delete(ctx, "Microsoft.Contoso", "testResourceType", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_GetNestedResourceTypeFirst.json
func ExampleSKUsClient_GetNestedResourceTypeFirst() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().GetNestedResourceTypeFirst(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// 	Properties: &armproviderhub.SKUResourceProperties{
	// 		SKUSettings: []*armproviderhub.SKUSetting{
	// 			{
	// 				Name: to.Ptr("freeSku"),
	// 				Kind: to.Ptr("Standard"),
	// 				Tier: to.Ptr("Tier1"),
	// 			},
	// 			{
	// 				Name: to.Ptr("premiumSku"),
	// 				Costs: []*armproviderhub.SKUCost{
	// 					{
	// 						MeterID: to.Ptr("xxx"),
	// 				}},
	// 				Kind: to.Ptr("Premium"),
	// 				Tier: to.Ptr("Tier2"),
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_CreateOrUpdateNestedResourceTypeFirst.json
func ExampleSKUsClient_CreateOrUpdateNestedResourceTypeFirst() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().CreateOrUpdateNestedResourceTypeFirst(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "testSku", armproviderhub.SKUResource{
		Properties: &armproviderhub.SKUResourceProperties{
			SKUSettings: []*armproviderhub.SKUSetting{
				{
					Name: to.Ptr("freeSku"),
					Kind: to.Ptr("Standard"),
					Tier: to.Ptr("Tier1"),
				},
				{
					Name: to.Ptr("premiumSku"),
					Costs: []*armproviderhub.SKUCost{
						{
							MeterID: to.Ptr("xxx"),
						}},
					Kind: to.Ptr("Premium"),
					Tier: to.Ptr("Tier2"),
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_DeleteNestedResourceTypeFirst.json
func ExampleSKUsClient_DeleteNestedResourceTypeFirst() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewSKUsClient().DeleteNestedResourceTypeFirst(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_GetNestedResourceTypeSecond.json
func ExampleSKUsClient_GetNestedResourceTypeSecond() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().GetNestedResourceTypeSecond(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/nestedEmployee2/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/resourcetypeRegistrations/nestedEmployee2/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// 	Properties: &armproviderhub.SKUResourceProperties{
	// 		SKUSettings: []*armproviderhub.SKUSetting{
	// 			{
	// 				Name: to.Ptr("freeSku"),
	// 				Kind: to.Ptr("Standard"),
	// 				Tier: to.Ptr("Tier1"),
	// 			},
	// 			{
	// 				Name: to.Ptr("premiumSku"),
	// 				Costs: []*armproviderhub.SKUCost{
	// 					{
	// 						MeterID: to.Ptr("xxx"),
	// 				}},
	// 				Kind: to.Ptr("Premium"),
	// 				Tier: to.Ptr("Tier2"),
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_CreateOrUpdateNestedResourceTypeSecond.json
func ExampleSKUsClient_CreateOrUpdateNestedResourceTypeSecond() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().CreateOrUpdateNestedResourceTypeSecond(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", "testSku", armproviderhub.SKUResource{
		Properties: &armproviderhub.SKUResourceProperties{
			SKUSettings: []*armproviderhub.SKUSetting{
				{
					Name: to.Ptr("freeSku"),
					Kind: to.Ptr("Standard"),
					Tier: to.Ptr("Tier1"),
				},
				{
					Name: to.Ptr("premiumSku"),
					Costs: []*armproviderhub.SKUCost{
						{
							MeterID: to.Ptr("xxx"),
						}},
					Kind: to.Ptr("Premium"),
					Tier: to.Ptr("Tier2"),
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/nestedEmployee2/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/resourcetypeRegistrations/nestedEmployee2/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_DeleteNestedResourceTypeSecond.json
func ExampleSKUsClient_DeleteNestedResourceTypeSecond() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewSKUsClient().DeleteNestedResourceTypeSecond(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_GetNestedResourceTypeThird.json
func ExampleSKUsClient_GetNestedResourceTypeThird() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().GetNestedResourceTypeThird(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", "nestedResourceTypeThird", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/nestedEmployee2/nestedEmployee3/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/resourcetypeRegistrations/nestedEmployee2/resourcetypeRegistrations/nestedEmployee3/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// 	Properties: &armproviderhub.SKUResourceProperties{
	// 		SKUSettings: []*armproviderhub.SKUSetting{
	// 			{
	// 				Name: to.Ptr("freeSku"),
	// 				Kind: to.Ptr("Standard"),
	// 				Tier: to.Ptr("Tier1"),
	// 			},
	// 			{
	// 				Name: to.Ptr("premiumSku"),
	// 				Costs: []*armproviderhub.SKUCost{
	// 					{
	// 						MeterID: to.Ptr("xxx"),
	// 				}},
	// 				Kind: to.Ptr("Premium"),
	// 				Tier: to.Ptr("Tier2"),
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_CreateOrUpdateNestedResourceTypeThird.json
func ExampleSKUsClient_CreateOrUpdateNestedResourceTypeThird() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSKUsClient().CreateOrUpdateNestedResourceTypeThird(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", "nestedResourceTypeThird", "testSku", armproviderhub.SKUResource{
		Properties: &armproviderhub.SKUResourceProperties{
			SKUSettings: []*armproviderhub.SKUSetting{
				{
					Name: to.Ptr("freeSku"),
					Kind: to.Ptr("Standard"),
					Tier: to.Ptr("Tier1"),
				},
				{
					Name: to.Ptr("premiumSku"),
					Costs: []*armproviderhub.SKUCost{
						{
							MeterID: to.Ptr("xxx"),
						}},
					Kind: to.Ptr("Premium"),
					Tier: to.Ptr("Tier2"),
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SKUResource = armproviderhub.SKUResource{
	// 	Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/nestedEmployee2/nestedEmployee3/sku1"),
	// 	Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
	// 	ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/resourcetypeRegistrations/nestedEmployee2/resourcetypeRegistrations/nestedEmployee3/skus/sku1"),
	// 	SystemData: &armproviderhub.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_DeleteNestedResourceTypeThird.json
func ExampleSKUsClient_DeleteNestedResourceTypeThird() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewSKUsClient().DeleteNestedResourceTypeThird(ctx, "Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", "nestedResourceTypeThird", "testSku", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_ListByResourceTypeRegistrations.json
func ExampleSKUsClient_NewListByResourceTypeRegistrationsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSKUsClient().NewListByResourceTypeRegistrationsPager("Microsoft.Contoso", "testResourceType", nil)
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
		// page.SKUResourceArrayResponseWithContinuation = armproviderhub.SKUResourceArrayResponseWithContinuation{
		// 	Value: []*armproviderhub.SKUResource{
		// 		{
		// 			Name: to.Ptr("Microsoft.Contoso/employees/sku1"),
		// 			Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/skus"),
		// 			ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/skus/sku1"),
		// 			SystemData: &armproviderhub.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 			},
		// 			Properties: &armproviderhub.SKUResourceProperties{
		// 				SKUSettings: []*armproviderhub.SKUSetting{
		// 					{
		// 						Name: to.Ptr("freeSku"),
		// 						Kind: to.Ptr("Standard"),
		// 						Tier: to.Ptr("Tier1"),
		// 					},
		// 					{
		// 						Name: to.Ptr("premiumSku"),
		// 						Costs: []*armproviderhub.SKUCost{
		// 							{
		// 								MeterID: to.Ptr("xxx"),
		// 						}},
		// 						Kind: to.Ptr("Premium"),
		// 						Tier: to.Ptr("Tier2"),
		// 				}},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_ListByResourceTypeRegistrationsNestedResourceTypeFirst.json
func ExampleSKUsClient_NewListByResourceTypeRegistrationsNestedResourceTypeFirstPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSKUsClient().NewListByResourceTypeRegistrationsNestedResourceTypeFirstPager("Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", nil)
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
		// page.SKUResourceArrayResponseWithContinuation = armproviderhub.SKUResourceArrayResponseWithContinuation{
		// 	Value: []*armproviderhub.SKUResource{
		// 		{
		// 			Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/sku1"),
		// 			Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
		// 			ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/skus/sku1"),
		// 			SystemData: &armproviderhub.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 			},
		// 			Properties: &armproviderhub.SKUResourceProperties{
		// 				SKUSettings: []*armproviderhub.SKUSetting{
		// 					{
		// 						Name: to.Ptr("freeSku"),
		// 						Kind: to.Ptr("Standard"),
		// 						Tier: to.Ptr("Tier1"),
		// 					},
		// 					{
		// 						Name: to.Ptr("premiumSku"),
		// 						Costs: []*armproviderhub.SKUCost{
		// 							{
		// 								MeterID: to.Ptr("xxx"),
		// 						}},
		// 						Kind: to.Ptr("Premium"),
		// 						Tier: to.Ptr("Tier2"),
		// 				}},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_ListByResourceTypeRegistrationsNestedResourceTypeSecond.json
func ExampleSKUsClient_NewListByResourceTypeRegistrationsNestedResourceTypeSecondPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSKUsClient().NewListByResourceTypeRegistrationsNestedResourceTypeSecondPager("Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", nil)
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
		// page.SKUResourceArrayResponseWithContinuation = armproviderhub.SKUResourceArrayResponseWithContinuation{
		// 	Value: []*armproviderhub.SKUResource{
		// 		{
		// 			Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/nestedEmployee2/sku1"),
		// 			Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
		// 			ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/resourcetypeRegistrations/nestedEmployee2/skus/sku1"),
		// 			SystemData: &armproviderhub.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 			},
		// 			Properties: &armproviderhub.SKUResourceProperties{
		// 				SKUSettings: []*armproviderhub.SKUSetting{
		// 					{
		// 						Name: to.Ptr("freeSku"),
		// 						Kind: to.Ptr("Standard"),
		// 						Tier: to.Ptr("Tier1"),
		// 					},
		// 					{
		// 						Name: to.Ptr("premiumSku"),
		// 						Costs: []*armproviderhub.SKUCost{
		// 							{
		// 								MeterID: to.Ptr("xxx"),
		// 						}},
		// 						Kind: to.Ptr("Premium"),
		// 						Tier: to.Ptr("Tier2"),
		// 				}},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7611bb6c9bad11244f4351eecfc50b2c46a86fde/specification/providerhub/resource-manager/Microsoft.ProviderHub/stable/2024-09-01/examples/Skus_ListByResourceTypeRegistrationsNestedResourceTypeThird.json
func ExampleSKUsClient_NewListByResourceTypeRegistrationsNestedResourceTypeThirdPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armproviderhub.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSKUsClient().NewListByResourceTypeRegistrationsNestedResourceTypeThirdPager("Microsoft.Contoso", "testResourceType", "nestedResourceTypeFirst", "nestedResourceTypeSecond", "nestedResourceTypeThird", nil)
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
		// page.SKUResourceArrayResponseWithContinuation = armproviderhub.SKUResourceArrayResponseWithContinuation{
		// 	Value: []*armproviderhub.SKUResource{
		// 		{
		// 			Name: to.Ptr("Microsoft.Contoso/employees/nestedEmployee/nestedEmployee2/nestedEmployee3/sku1"),
		// 			Type: to.Ptr("Microsoft.ProviderHub/providerRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/resourcetypeRegistrations/skus"),
		// 			ID: to.Ptr("/subscriptions/ab7a8701-f7ef-471a-a2f4-d0ebbf494f77/providers/Microsoft.ProviderHub/providerRegistrations/Microsoft.Contoso/resourcetypeRegistrations/employees/resourcetypeRegistrations/nestedEmployee/resourcetypeRegistrations/nestedEmployee2/resourcetypeRegistrations/nestedEmployee3/skus/sku1"),
		// 			SystemData: &armproviderhub.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-01T01:01:01.107Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armproviderhub.CreatedByTypeUser),
		// 			},
		// 			Properties: &armproviderhub.SKUResourceProperties{
		// 				SKUSettings: []*armproviderhub.SKUSetting{
		// 					{
		// 						Name: to.Ptr("freeSku"),
		// 						Kind: to.Ptr("Standard"),
		// 						Tier: to.Ptr("Tier1"),
		// 					},
		// 					{
		// 						Name: to.Ptr("premiumSku"),
		// 						Costs: []*armproviderhub.SKUCost{
		// 							{
		// 								MeterID: to.Ptr("xxx"),
		// 						}},
		// 						Kind: to.Ptr("Premium"),
		// 						Tier: to.Ptr("Tier2"),
		// 				}},
		// 			},
		// 	}},
		// }
	}
}
