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

// Generated from example definition: 2024-12-01/SolutionTypes_Get.json
func ExampleSolutionTypesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSolutionTypesClient("5ACC4579-DB34-4C2F-8F8C-25061168F342").Get(ctx, "rgpublicCloud", "lulzqllpu", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.SolutionTypesClientGetResponse{
	// 	SolutionTypeResource: &armhybridconnectivity.SolutionTypeResource{
	// 		Properties: &armhybridconnectivity.SolutionTypeProperties{
	// 			SolutionType: to.Ptr("tjtoeycxhyqxtgd"),
	// 			Description: to.Ptr("wxyxcvtzuxgodtlanjevedwfdwnznc"),
	// 			SupportedAzureRegions: []*string{
	// 				to.Ptr("cimocdh"),
	// 			},
	// 			SolutionSettings: []*armhybridconnectivity.SolutionTypeSettingsProperties{
	// 				{
	// 					Name: to.Ptr("tepghdgbefujhnnue"),
	// 					DisplayName: to.Ptr("mwlzepoin"),
	// 					Type: to.Ptr("je"),
	// 					Description: to.Ptr("soq"),
	// 					AllowedValues: []*string{
	// 						to.Ptr("pwizyngpkpxsllpluffjspx"),
	// 					},
	// 					DefaultValue: to.Ptr("laekyetgapdpxyqervqaqfscfwagek"),
	// 				},
	// 			},
	// 		},
	// 		ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/solutionTypes/j"),
	// 		Name: to.Ptr("xczyyxuphhacyyj"),
	// 		Type: to.Ptr("mf"),
	// 		SystemData: &armhybridconnectivity.SystemData{
	// 			CreatedBy: to.Ptr("rpxzkcrobprrdvuoqxz"),
	// 			CreatedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("jidegyskxi"),
	// 			LastModifiedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-12-01/SolutionTypes_ListByResourceGroup.json
func ExampleSolutionTypesClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSolutionTypesClient("5ACC4579-DB34-4C2F-8F8C-25061168F342").NewListByResourceGroupPager("rgpublicCloud", nil)
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
		// page = armhybridconnectivity.SolutionTypesClientListByResourceGroupResponse{
		// 	SolutionTypeResourceListResult: armhybridconnectivity.SolutionTypeResourceListResult{
		// 		Value: []*armhybridconnectivity.SolutionTypeResource{
		// 			{
		// 				Properties: &armhybridconnectivity.SolutionTypeProperties{
		// 					SolutionType: to.Ptr("j"),
		// 					Description: to.Ptr("mhasmuazxsr"),
		// 					SupportedAzureRegions: []*string{
		// 						to.Ptr("jfvkgljymtuzfwbumgabpdpjjnxit"),
		// 					},
		// 					SolutionSettings: []*armhybridconnectivity.SolutionTypeSettingsProperties{
		// 						{
		// 							Name: to.Ptr("eepvybtmsjwgcpf"),
		// 							DisplayName: to.Ptr("npxunbwkjrklbjsvdryzsjtecm"),
		// 							Type: to.Ptr("fngmzlffmwmrglepeyce"),
		// 							Description: to.Ptr("vdtvoysdagvae"),
		// 							AllowedValues: []*string{
		// 								to.Ptr("cgbkgbmsgsfofmcrjerg"),
		// 							},
		// 							DefaultValue: to.Ptr("knshmo"),
		// 						},
		// 					},
		// 				},
		// 				ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/solutionTypes/j"),
		// 				Name: to.Ptr("yzgpvbtzwvijawjgfvnhgbqefqq"),
		// 				Type: to.Ptr("vamwfnqqpjosjnomwbkwnlnrg"),
		// 				SystemData: &armhybridconnectivity.SystemData{
		// 					CreatedBy: to.Ptr("rpxzkcrobprrdvuoqxz"),
		// 					CreatedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("jidegyskxi"),
		// 					LastModifiedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
		// 				},
		// 			},
		// 		},
		// 		NextLink: to.Ptr("https://microsoft.com/a"),
		// 	},
		// }
	}
}

// Generated from example definition: 2024-12-01/SolutionTypes_ListBySubscription.json
func ExampleSolutionTypesClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSolutionTypesClient("5ACC4579-DB34-4C2F-8F8C-25061168F342").NewListBySubscriptionPager(nil)
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
		// page = armhybridconnectivity.SolutionTypesClientListBySubscriptionResponse{
		// 	SolutionTypeResourceListResult: armhybridconnectivity.SolutionTypeResourceListResult{
		// 		Value: []*armhybridconnectivity.SolutionTypeResource{
		// 			{
		// 				Properties: &armhybridconnectivity.SolutionTypeProperties{
		// 					SolutionType: to.Ptr("dembhpcydwoiyszmjtniletpy"),
		// 					Description: to.Ptr("fkegiumpjdwgkde"),
		// 					SupportedAzureRegions: []*string{
		// 						to.Ptr("ujawbfint"),
		// 					},
		// 					SolutionSettings: []*armhybridconnectivity.SolutionTypeSettingsProperties{
		// 						{
		// 							Name: to.Ptr("i"),
		// 							DisplayName: to.Ptr("forzmqskffaub"),
		// 							Type: to.Ptr("d"),
		// 							Description: to.Ptr("cbyxsxfsaye"),
		// 							AllowedValues: []*string{
		// 								to.Ptr("uecqnmmssdeusxejcxrtkskfugvl"),
		// 							},
		// 							DefaultValue: to.Ptr("uzwiymoxrummkoowwvzjhyazeavzr"),
		// 						},
		// 					},
		// 				},
		// 				ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/solutionTypes/i"),
		// 				Name: to.Ptr("jtlxwihbuftmaobxfmfjojalhpwrv"),
		// 				Type: to.Ptr("zditfautattfhnffvjw"),
		// 				SystemData: &armhybridconnectivity.SystemData{
		// 					CreatedBy: to.Ptr("rpxzkcrobprrdvuoqxz"),
		// 					CreatedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("jidegyskxi"),
		// 					LastModifiedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
		// 				},
		// 			},
		// 		},
		// 		NextLink: to.Ptr("https://microsoft.com/a"),
		// 	},
		// }
	}
}
