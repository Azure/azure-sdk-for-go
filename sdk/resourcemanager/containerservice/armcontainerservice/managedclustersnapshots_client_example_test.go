//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcontainerservice_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/28a163a3b110ac1d7002449ee2d04553c725ec0d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/ManagedClusterSnapshotsList.json
func ExampleManagedClusterSnapshotsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewManagedClusterSnapshotsClient().NewListPager(nil)
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
		// page.ManagedClusterSnapshotListResult = armcontainerservice.ManagedClusterSnapshotListResult{
		// 	Value: []*armcontainerservice.ManagedClusterSnapshot{
		// 		{
		// 			Name: to.Ptr("snapshot1"),
		// 			Type: to.Ptr("Microsoft.ContainerService/ManagedClusterSnapshots"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedclustersnapshots/snapshot1"),
		// 			SystemData: &armcontainerservice.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-08-09T20:13:23.298Z"); return t}()),
		// 				CreatedBy: to.Ptr("user1"),
		// 				CreatedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
		// 			},
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"key1": to.Ptr("val1"),
		// 				"key2": to.Ptr("val2"),
		// 			},
		// 			Properties: &armcontainerservice.ManagedClusterSnapshotProperties{
		// 				CreationData: &armcontainerservice.CreationData{
		// 					SourceResourceID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
		// 				},
		// 				ManagedClusterPropertiesReadOnly: &armcontainerservice.ManagedClusterPropertiesForSnapshot{
		// 					EnableRbac: to.Ptr(true),
		// 					KubernetesVersion: to.Ptr("1.20.5"),
		// 					NetworkProfile: &armcontainerservice.NetworkProfileForSnapshot{
		// 						LoadBalancerSKU: to.Ptr(armcontainerservice.LoadBalancerSKUStandard),
		// 						NetworkMode: to.Ptr(armcontainerservice.NetworkModeBridge),
		// 						NetworkPlugin: to.Ptr(armcontainerservice.NetworkPluginKubenet),
		// 						NetworkPolicy: to.Ptr(armcontainerservice.NetworkPolicyCalico),
		// 					},
		// 					SKU: &armcontainerservice.ManagedClusterSKU{
		// 						Name: to.Ptr(armcontainerservice.ManagedClusterSKUName("Basic")),
		// 						Tier: to.Ptr(armcontainerservice.ManagedClusterSKUTierFree),
		// 					},
		// 				},
		// 				SnapshotType: to.Ptr(armcontainerservice.SnapshotTypeManagedCluster),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/28a163a3b110ac1d7002449ee2d04553c725ec0d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/ManagedClusterSnapshotsListByResourceGroup.json
func ExampleManagedClusterSnapshotsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewManagedClusterSnapshotsClient().NewListByResourceGroupPager("rg1", nil)
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
		// page.ManagedClusterSnapshotListResult = armcontainerservice.ManagedClusterSnapshotListResult{
		// 	Value: []*armcontainerservice.ManagedClusterSnapshot{
		// 		{
		// 			Name: to.Ptr("snapshot1"),
		// 			Type: to.Ptr("Microsoft.ContainerService/ManagedClusterSnapshots"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedclustersnapshots/snapshot1"),
		// 			SystemData: &armcontainerservice.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-08-09T20:13:23.298Z"); return t}()),
		// 				CreatedBy: to.Ptr("user1"),
		// 				CreatedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
		// 			},
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"key1": to.Ptr("val1"),
		// 				"key2": to.Ptr("val2"),
		// 			},
		// 			Properties: &armcontainerservice.ManagedClusterSnapshotProperties{
		// 				CreationData: &armcontainerservice.CreationData{
		// 					SourceResourceID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
		// 				},
		// 				ManagedClusterPropertiesReadOnly: &armcontainerservice.ManagedClusterPropertiesForSnapshot{
		// 					EnableRbac: to.Ptr(true),
		// 					KubernetesVersion: to.Ptr("1.20.5"),
		// 					NetworkProfile: &armcontainerservice.NetworkProfileForSnapshot{
		// 						LoadBalancerSKU: to.Ptr(armcontainerservice.LoadBalancerSKUStandard),
		// 						NetworkMode: to.Ptr(armcontainerservice.NetworkModeBridge),
		// 						NetworkPlugin: to.Ptr(armcontainerservice.NetworkPluginKubenet),
		// 						NetworkPolicy: to.Ptr(armcontainerservice.NetworkPolicyCalico),
		// 					},
		// 					SKU: &armcontainerservice.ManagedClusterSKU{
		// 						Name: to.Ptr(armcontainerservice.ManagedClusterSKUName("Basic")),
		// 						Tier: to.Ptr(armcontainerservice.ManagedClusterSKUTierFree),
		// 					},
		// 				},
		// 				SnapshotType: to.Ptr(armcontainerservice.SnapshotTypeManagedCluster),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/28a163a3b110ac1d7002449ee2d04553c725ec0d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/ManagedClusterSnapshotsGet.json
func ExampleManagedClusterSnapshotsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagedClusterSnapshotsClient().Get(ctx, "rg1", "snapshot1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ManagedClusterSnapshot = armcontainerservice.ManagedClusterSnapshot{
	// 	Name: to.Ptr("snapshot1"),
	// 	Type: to.Ptr("Microsoft.ContainerService/ManagedClusterSnapshots"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedclustersnapshots/snapshot1"),
	// 	SystemData: &armcontainerservice.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-08-09T20:13:23.298Z"); return t}()),
	// 		CreatedBy: to.Ptr("user1"),
	// 		CreatedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"key1": to.Ptr("val1"),
	// 		"key2": to.Ptr("val2"),
	// 	},
	// 	Properties: &armcontainerservice.ManagedClusterSnapshotProperties{
	// 		CreationData: &armcontainerservice.CreationData{
	// 			SourceResourceID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
	// 		},
	// 		ManagedClusterPropertiesReadOnly: &armcontainerservice.ManagedClusterPropertiesForSnapshot{
	// 			EnableRbac: to.Ptr(true),
	// 			KubernetesVersion: to.Ptr("1.20.5"),
	// 			NetworkProfile: &armcontainerservice.NetworkProfileForSnapshot{
	// 				LoadBalancerSKU: to.Ptr(armcontainerservice.LoadBalancerSKUStandard),
	// 				NetworkMode: to.Ptr(armcontainerservice.NetworkModeBridge),
	// 				NetworkPlugin: to.Ptr(armcontainerservice.NetworkPluginKubenet),
	// 				NetworkPolicy: to.Ptr(armcontainerservice.NetworkPolicyCalico),
	// 			},
	// 			SKU: &armcontainerservice.ManagedClusterSKU{
	// 				Name: to.Ptr(armcontainerservice.ManagedClusterSKUName("Basic")),
	// 				Tier: to.Ptr(armcontainerservice.ManagedClusterSKUTierFree),
	// 			},
	// 		},
	// 		SnapshotType: to.Ptr(armcontainerservice.SnapshotTypeManagedCluster),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/28a163a3b110ac1d7002449ee2d04553c725ec0d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/ManagedClusterSnapshotsCreate.json
func ExampleManagedClusterSnapshotsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagedClusterSnapshotsClient().CreateOrUpdate(ctx, "rg1", "snapshot1", armcontainerservice.ManagedClusterSnapshot{
		Location: to.Ptr("westus"),
		Tags: map[string]*string{
			"key1": to.Ptr("val1"),
			"key2": to.Ptr("val2"),
		},
		Properties: &armcontainerservice.ManagedClusterSnapshotProperties{
			CreationData: &armcontainerservice.CreationData{
				SourceResourceID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ManagedClusterSnapshot = armcontainerservice.ManagedClusterSnapshot{
	// 	Name: to.Ptr("snapshot1"),
	// 	Type: to.Ptr("Microsoft.ContainerService/ManagedClusterSnapshots"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedclustersnapshots/snapshot1"),
	// 	SystemData: &armcontainerservice.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-08-09T20:13:23.298Z"); return t}()),
	// 		CreatedBy: to.Ptr("user1"),
	// 		CreatedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"key1": to.Ptr("val1"),
	// 		"key2": to.Ptr("val2"),
	// 	},
	// 	Properties: &armcontainerservice.ManagedClusterSnapshotProperties{
	// 		CreationData: &armcontainerservice.CreationData{
	// 			SourceResourceID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
	// 		},
	// 		ManagedClusterPropertiesReadOnly: &armcontainerservice.ManagedClusterPropertiesForSnapshot{
	// 			EnableRbac: to.Ptr(true),
	// 			KubernetesVersion: to.Ptr("1.20.5"),
	// 			NetworkProfile: &armcontainerservice.NetworkProfileForSnapshot{
	// 				LoadBalancerSKU: to.Ptr(armcontainerservice.LoadBalancerSKUStandard),
	// 				NetworkMode: to.Ptr(armcontainerservice.NetworkModeBridge),
	// 				NetworkPlugin: to.Ptr(armcontainerservice.NetworkPluginKubenet),
	// 				NetworkPolicy: to.Ptr(armcontainerservice.NetworkPolicyCalico),
	// 			},
	// 			SKU: &armcontainerservice.ManagedClusterSKU{
	// 				Name: to.Ptr(armcontainerservice.ManagedClusterSKUName("Basic")),
	// 				Tier: to.Ptr(armcontainerservice.ManagedClusterSKUTierFree),
	// 			},
	// 		},
	// 		SnapshotType: to.Ptr(armcontainerservice.SnapshotTypeManagedCluster),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/28a163a3b110ac1d7002449ee2d04553c725ec0d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/ManagedClusterSnapshotsUpdateTags.json
func ExampleManagedClusterSnapshotsClient_UpdateTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagedClusterSnapshotsClient().UpdateTags(ctx, "rg1", "snapshot1", armcontainerservice.TagsObject{
		Tags: map[string]*string{
			"key2": to.Ptr("new-val2"),
			"key3": to.Ptr("val3"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ManagedClusterSnapshot = armcontainerservice.ManagedClusterSnapshot{
	// 	Name: to.Ptr("snapshot1"),
	// 	Type: to.Ptr("Microsoft.ContainerService/ManagedClusterSnapshots"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedclustersnapshots/snapshot1"),
	// 	SystemData: &armcontainerservice.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-08-09T20:13:23.298Z"); return t}()),
	// 		CreatedBy: to.Ptr("user1"),
	// 		CreatedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"key1": to.Ptr("val1"),
	// 		"key2": to.Ptr("val2"),
	// 	},
	// 	Properties: &armcontainerservice.ManagedClusterSnapshotProperties{
	// 		CreationData: &armcontainerservice.CreationData{
	// 			SourceResourceID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster1"),
	// 		},
	// 		ManagedClusterPropertiesReadOnly: &armcontainerservice.ManagedClusterPropertiesForSnapshot{
	// 			EnableRbac: to.Ptr(true),
	// 			KubernetesVersion: to.Ptr("1.20.5"),
	// 			NetworkProfile: &armcontainerservice.NetworkProfileForSnapshot{
	// 				LoadBalancerSKU: to.Ptr(armcontainerservice.LoadBalancerSKUStandard),
	// 				NetworkMode: to.Ptr(armcontainerservice.NetworkModeBridge),
	// 				NetworkPlugin: to.Ptr(armcontainerservice.NetworkPluginKubenet),
	// 				NetworkPolicy: to.Ptr(armcontainerservice.NetworkPolicyCalico),
	// 			},
	// 			SKU: &armcontainerservice.ManagedClusterSKU{
	// 				Name: to.Ptr(armcontainerservice.ManagedClusterSKUName("Basic")),
	// 				Tier: to.Ptr(armcontainerservice.ManagedClusterSKUTierFree),
	// 			},
	// 		},
	// 		SnapshotType: to.Ptr(armcontainerservice.SnapshotTypeManagedCluster),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/28a163a3b110ac1d7002449ee2d04553c725ec0d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/ManagedClusterSnapshotsDelete.json
func ExampleManagedClusterSnapshotsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewManagedClusterSnapshotsClient().Delete(ctx, "rg1", "snapshot1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
