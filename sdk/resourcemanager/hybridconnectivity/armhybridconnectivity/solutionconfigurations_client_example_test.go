// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armhybridconnectivity_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridconnectivity/armhybridconnectivity"
	"log"
)

// Generated from example definition: 2024-12-01/SolutionConfigurations_CreateOrUpdate.json
func ExampleSolutionConfigurationsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSolutionConfigurationsClient().CreateOrUpdate(ctx, "ymuj", "keebwujt", armhybridconnectivity.SolutionConfiguration{
		Properties: &armhybridconnectivity.SolutionConfigurationProperties{
			SolutionType:     to.Ptr("nmtqllkyohwtsthxaimsye"),
			SolutionSettings: &armhybridconnectivity.SolutionSettings{},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.SolutionConfigurationsClientCreateOrUpdateResponse{
	// 	SolutionConfiguration: &armhybridconnectivity.SolutionConfiguration{
	// 		Properties: &armhybridconnectivity.SolutionConfigurationProperties{
	// 			SolutionType: to.Ptr("nmtqllkyohwtsthxaimsye"),
	// 			SolutionSettings: &armhybridconnectivity.SolutionSettings{
	// 			},
	// 			ProvisioningState: to.Ptr(armhybridconnectivity.ResourceProvisioningStateSucceeded),
	// 			Status: to.Ptr(armhybridconnectivity.SolutionConfigurationStatusNew),
	// 			StatusDetails: to.Ptr("rqbrzildwecankrpukkbjjqrczxboz"),
	// 			LastSyncTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-20T03:24:15.820Z"); return t}()),
	// 		},
	// 		ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/SolutionConfigurations/qpwubemzmootxmtlxaerir"),
	// 		Name: to.Ptr("qpwubemzmootxmtlxaerir"),
	// 		Type: to.Ptr("uknrk"),
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

// Generated from example definition: 2024-12-01/SolutionConfigurations_Delete.json
func ExampleSolutionConfigurationsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSolutionConfigurationsClient().Delete(ctx, "ymuj", "stu", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.SolutionConfigurationsClientDeleteResponse{
	// }
}

// Generated from example definition: 2024-12-01/SolutionConfigurations_Get.json
func ExampleSolutionConfigurationsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSolutionConfigurationsClient().Get(ctx, "ymuj", "tks", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.SolutionConfigurationsClientGetResponse{
	// 	SolutionConfiguration: &armhybridconnectivity.SolutionConfiguration{
	// 		Properties: &armhybridconnectivity.SolutionConfigurationProperties{
	// 			SolutionType: to.Ptr("nmtqllkyohwtsthxaimsye"),
	// 			SolutionSettings: &armhybridconnectivity.SolutionSettings{
	// 			},
	// 			ProvisioningState: to.Ptr(armhybridconnectivity.ResourceProvisioningStateSucceeded),
	// 			Status: to.Ptr(armhybridconnectivity.SolutionConfigurationStatusNew),
	// 			StatusDetails: to.Ptr("rqbrzildwecankrpukkbjjqrczxboz"),
	// 			LastSyncTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-20T03:24:15.820Z"); return t}()),
	// 		},
	// 		ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/SolutionConfigurations/qpwubemzmootxmtlxaerir"),
	// 		Name: to.Ptr("qpwubemzmootxmtlxaerir"),
	// 		Type: to.Ptr("uknrk"),
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

// Generated from example definition: 2024-12-01/SolutionConfigurations_List.json
func ExampleSolutionConfigurationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSolutionConfigurationsClient().NewListPager("ymuj", nil)
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
		// page = armhybridconnectivity.SolutionConfigurationsClientListResponse{
		// 	SolutionConfigurationListResult: armhybridconnectivity.SolutionConfigurationListResult{
		// 		Value: []*armhybridconnectivity.SolutionConfiguration{
		// 			{
		// 				Properties: &armhybridconnectivity.SolutionConfigurationProperties{
		// 					SolutionType: to.Ptr("Microsoft.AssetManagement"),
		// 					SolutionSettings: &armhybridconnectivity.SolutionSettings{
		// 					},
		// 					ProvisioningState: to.Ptr(armhybridconnectivity.ResourceProvisioningStateSucceeded),
		// 					StatusDetails: to.Ptr("Aws authorization validation pending in Aws account"),
		// 					LastSyncTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-20T03:24:15.820Z"); return t}()),
		// 				},
		// 				ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/SolutionConfigurations/solutionconfigurationtest"),
		// 				Name: to.Ptr("solutionconfigurationtest"),
		// 				Type: to.Ptr("microsoft.hybridconnectivity/solutionconfigurations"),
		// 				SystemData: &armhybridconnectivity.SystemData{
		// 					CreatedBy: to.Ptr("rpxzkcrobprrdvuoqxz"),
		// 					CreatedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("jidegyskxi"),
		// 					LastModifiedByType: to.Ptr(armhybridconnectivity.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-18T22:52:07.890Z"); return t}()),
		// 				},
		// 			},
		// 			{
		// 				Properties: &armhybridconnectivity.SolutionConfigurationProperties{
		// 					SolutionType: to.Ptr("Microsoft.HybridCompute"),
		// 					SolutionSettings: &armhybridconnectivity.SolutionSettings{
		// 					},
		// 					ProvisioningState: to.Ptr(armhybridconnectivity.ResourceProvisioningStateSucceeded),
		// 					StatusDetails: to.Ptr("Aws authorization validation succeeded in Aws account"),
		// 					LastSyncTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-20T03:24:15.820Z"); return t}()),
		// 				},
		// 				ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/SolutionConfigurations/solutionconfigurationtest2"),
		// 				Name: to.Ptr("solutionconfigurationtest2"),
		// 				Type: to.Ptr("microsoft.hybridconnectivity/solutionconfigurations"),
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

// Generated from example definition: 2024-12-01/SolutionConfigurations_SyncNow.json
func ExampleSolutionConfigurationsClient_BeginSyncNow() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSolutionConfigurationsClient().BeginSyncNow(ctx, "ymuj", "tks", nil)
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
	// res = armhybridconnectivity.SolutionConfigurationsClientSyncNowResponse{
	// 	OperationStatusResult: &armhybridconnectivity.OperationStatusResult{
	// 		ID: to.Ptr("/subscriptions/5ACC4579-DB34-4C2F-8F8C-25061168F342/providers/Microsoft.HybridConnectivity/PublicCloudConnectors/esixipkbydb"),
	// 		ResourceID: to.Ptr("/subscriptions/5ACC4579-DB34-4C2F-8F8C-25061168F342/providers/Microsoft.HybridConnectivity/PublicCloudConnectors/esixipkbydb"),
	// 		Name: to.Ptr("svqtraeuwvyvblujlvqilypwpdrt"),
	// 		Status: to.Ptr("bevmrejij"),
	// 		PercentComplete: to.Ptr[float64](15),
	// 		StartTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-10-02T18:38:19.143Z"); return t}()),
	// 		EndTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-10-02T18:38:19.143Z"); return t}()),
	// 		Operations: []*armhybridconnectivity.OperationStatusResult{
	// 		},
	// 		Error: &armhybridconnectivity.ErrorDetail{
	// 			Code: to.Ptr("ykzvluyqiftfsumgvwzdh"),
	// 			Message: to.Ptr("krbjgtqkjgiux"),
	// 			Target: to.Ptr("nsaucxt"),
	// 			Details: []*armhybridconnectivity.ErrorDetail{
	// 			},
	// 			AdditionalInfo: []*armhybridconnectivity.ErrorAdditionalInfo{
	// 				{
	// 					Type: to.Ptr("qivvrewsjvcildjgwwytgimwklh"),
	// 					Info: &armhybridconnectivity.ErrorAdditionalInfoInfo{
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-12-01/SolutionConfigurations_Update.json
func ExampleSolutionConfigurationsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridconnectivity.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSolutionConfigurationsClient().Update(ctx, "ymuj", "dxt", armhybridconnectivity.SolutionConfiguration{
		Properties: &armhybridconnectivity.SolutionConfigurationProperties{
			SolutionType:     to.Ptr("myzljlstvmgkp"),
			SolutionSettings: &armhybridconnectivity.SolutionSettings{},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armhybridconnectivity.SolutionConfigurationsClientUpdateResponse{
	// 	SolutionConfiguration: &armhybridconnectivity.SolutionConfiguration{
	// 		Properties: &armhybridconnectivity.SolutionConfigurationProperties{
	// 			SolutionType: to.Ptr("nmtqllkyohwtsthxaimsye"),
	// 			SolutionSettings: &armhybridconnectivity.SolutionSettings{
	// 			},
	// 			ProvisioningState: to.Ptr(armhybridconnectivity.ResourceProvisioningStateSucceeded),
	// 			Status: to.Ptr(armhybridconnectivity.SolutionConfigurationStatusNew),
	// 			StatusDetails: to.Ptr("rqbrzildwecankrpukkbjjqrczxboz"),
	// 			LastSyncTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-20T03:24:15.820Z"); return t}()),
	// 		},
	// 		ID: to.Ptr("/subscriptions/testSubcrptions/resourceGroups/testResourceGroup/providers/Microsoft.HybridConnectivity/SolutionConfigurations/qpwubemzmootxmtlxaerir"),
	// 		Name: to.Ptr("qpwubemzmootxmtlxaerir"),
	// 		Type: to.Ptr("uknrk"),
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
