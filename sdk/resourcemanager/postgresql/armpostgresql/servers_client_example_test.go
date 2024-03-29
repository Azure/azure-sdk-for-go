//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpostgresql_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerCreatePointInTimeRestore.json
func ExampleServersClient_BeginCreate_createADatabaseAsAPointInTimeRestore() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServersClient().BeginCreate(ctx, "TargetResourceGroup", "targetserver", armpostgresql.ServerForCreate{
		Location: to.Ptr("brazilsouth"),
		Properties: &armpostgresql.ServerPropertiesForRestore{
			CreateMode:         to.Ptr(armpostgresql.CreateModePointInTimeRestore),
			RestorePointInTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-14T00:00:37.467Z"); return t }()),
			SourceServerID:     to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/SourceResourceGroup/providers/Microsoft.DBforPostgreSQL/servers/sourceserver"),
		},
		SKU: &armpostgresql.SKU{
			Name:     to.Ptr("B_Gen5_2"),
			Capacity: to.Ptr[int32](2),
			Family:   to.Ptr("Gen5"),
			Tier:     to.Ptr(armpostgresql.SKUTierBasic),
		},
		Tags: map[string]*string{
			"ElasticServer": to.Ptr("1"),
		},
	}, nil)
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
	// res.Server = armpostgresql.Server{
	// 	Name: to.Ptr("targetserver"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/targetserver"),
	// 	Location: to.Ptr("brazilsouth"),
	// 	Tags: map[string]*string{
	// 		"ElasticServer": to.Ptr("1"),
	// 	},
	// 	Properties: &armpostgresql.ServerProperties{
	// 		AdministratorLogin: to.Ptr("cloudsa"),
	// 		EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-14T21:08:24.637Z"); return t}()),
	// 		FullyQualifiedDomainName: to.Ptr("targetserver.postgres.database.azure.com"),
	// 		SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
	// 		StorageProfile: &armpostgresql.StorageProfile{
	// 			BackupRetentionDays: to.Ptr[int32](7),
	// 			GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
	// 			StorageMB: to.Ptr[int32](128000),
	// 		},
	// 		UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
	// 		Version: to.Ptr(armpostgresql.ServerVersionNine6),
	// 	},
	// 	SKU: &armpostgresql.SKU{
	// 		Name: to.Ptr("B_Gen5_2"),
	// 		Capacity: to.Ptr[int32](2),
	// 		Family: to.Ptr("Gen5"),
	// 		Tier: to.Ptr(armpostgresql.SKUTierBasic),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerCreate.json
func ExampleServersClient_BeginCreate_createANewServer() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServersClient().BeginCreate(ctx, "TestGroup", "pgtestsvc4", armpostgresql.ServerForCreate{
		Location: to.Ptr("westus"),
		Properties: &armpostgresql.ServerPropertiesForDefaultCreate{
			CreateMode:        to.Ptr(armpostgresql.CreateModeDefault),
			MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
			SSLEnforcement:    to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
			StorageProfile: &armpostgresql.StorageProfile{
				BackupRetentionDays: to.Ptr[int32](7),
				GeoRedundantBackup:  to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
				StorageMB:           to.Ptr[int32](128000),
			},
			AdministratorLogin:         to.Ptr("cloudsa"),
			AdministratorLoginPassword: to.Ptr("<administratorLoginPassword>"),
		},
		SKU: &armpostgresql.SKU{
			Name:     to.Ptr("B_Gen5_2"),
			Capacity: to.Ptr[int32](2),
			Family:   to.Ptr("Gen5"),
			Tier:     to.Ptr(armpostgresql.SKUTierBasic),
		},
		Tags: map[string]*string{
			"ElasticServer": to.Ptr("1"),
		},
	}, nil)
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
	// res.Server = armpostgresql.Server{
	// 	Name: to.Ptr("pgtestsvc4"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc4"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"ElasticServer": to.Ptr("1"),
	// 	},
	// 	Properties: &armpostgresql.ServerProperties{
	// 		AdministratorLogin: to.Ptr("cloudsa"),
	// 		EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-14T21:08:24.637Z"); return t}()),
	// 		FullyQualifiedDomainName: to.Ptr("pgtestsvc4.postgres.database.azure.com"),
	// 		SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
	// 		StorageProfile: &armpostgresql.StorageProfile{
	// 			BackupRetentionDays: to.Ptr[int32](7),
	// 			GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
	// 			StorageMB: to.Ptr[int32](128000),
	// 		},
	// 		UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
	// 		Version: to.Ptr(armpostgresql.ServerVersionNine6),
	// 	},
	// 	SKU: &armpostgresql.SKU{
	// 		Name: to.Ptr("B_Gen5_2"),
	// 		Capacity: to.Ptr[int32](2),
	// 		Family: to.Ptr("Gen5"),
	// 		Tier: to.Ptr(armpostgresql.SKUTierBasic),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerCreateReplicaMode.json
func ExampleServersClient_BeginCreate_createAReplicaServer() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServersClient().BeginCreate(ctx, "TestGroup_WestCentralUS", "testserver-replica1", armpostgresql.ServerForCreate{
		Location: to.Ptr("westcentralus"),
		Properties: &armpostgresql.ServerPropertiesForReplica{
			CreateMode:     to.Ptr(armpostgresql.CreateModeReplica),
			SourceServerID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/TestGroup_WestCentralUS/providers/Microsoft.DBforPostgreSQL/servers/testserver-master"),
		},
		SKU: &armpostgresql.SKU{
			Name:     to.Ptr("GP_Gen5_2"),
			Capacity: to.Ptr[int32](2),
			Family:   to.Ptr("Gen5"),
			Tier:     to.Ptr(armpostgresql.SKUTierGeneralPurpose),
		},
	}, nil)
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
	// res.Server = armpostgresql.Server{
	// 	Name: to.Ptr("testserver-replica1"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/TestGroup_WestCentralUS/providers/Microsoft.DBforPostgreSQL/servers/testserver-replica1"),
	// 	Location: to.Ptr("westcentralus"),
	// 	Properties: &armpostgresql.ServerProperties{
	// 		AdministratorLogin: to.Ptr("postgres"),
	// 		EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-06-20T00:17:56.677Z"); return t}()),
	// 		FullyQualifiedDomainName: to.Ptr("testserver-replica1.postgres.database.azure.com"),
	// 		MasterServerID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/TestGroup_WestCentralUS/providers/Microsoft.DBforPostgreSQL/servers/testserver-master"),
	// 		ReplicaCapacity: to.Ptr[int32](0),
	// 		ReplicationRole: to.Ptr("Replica"),
	// 		SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumDisabled),
	// 		StorageProfile: &armpostgresql.StorageProfile{
	// 			BackupRetentionDays: to.Ptr[int32](7),
	// 			GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
	// 			StorageMB: to.Ptr[int32](2048000),
	// 		},
	// 		UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
	// 		Version: to.Ptr(armpostgresql.ServerVersionNine6),
	// 	},
	// 	SKU: &armpostgresql.SKU{
	// 		Name: to.Ptr("GP_Gen5_2"),
	// 		Capacity: to.Ptr[int32](2),
	// 		Family: to.Ptr("Gen4"),
	// 		Tier: to.Ptr(armpostgresql.SKUTierGeneralPurpose),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerCreateGeoRestoreMode.json
func ExampleServersClient_BeginCreate_createAServerAsAGeoRestore() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServersClient().BeginCreate(ctx, "TargetResourceGroup", "targetserver", armpostgresql.ServerForCreate{
		Location: to.Ptr("westus"),
		Properties: &armpostgresql.ServerPropertiesForGeoRestore{
			CreateMode:     to.Ptr(armpostgresql.CreateModeGeoRestore),
			SourceServerID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/SourceResourceGroup/providers/Microsoft.DBforPostgreSQL/servers/sourceserver"),
		},
		SKU: &armpostgresql.SKU{
			Name:     to.Ptr("GP_Gen5_2"),
			Capacity: to.Ptr[int32](2),
			Family:   to.Ptr("Gen5"),
			Tier:     to.Ptr(armpostgresql.SKUTierGeneralPurpose),
		},
		Tags: map[string]*string{
			"ElasticServer": to.Ptr("1"),
		},
	}, nil)
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
	// res.Server = armpostgresql.Server{
	// 	Name: to.Ptr("targetserver"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/targetserver"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"ElasticServer": to.Ptr("1"),
	// 	},
	// 	Properties: &armpostgresql.ServerProperties{
	// 		AdministratorLogin: to.Ptr("cloudsa"),
	// 		EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-14T21:08:24.637Z"); return t}()),
	// 		FullyQualifiedDomainName: to.Ptr("targetserver.postgres.database.azure.com"),
	// 		SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
	// 		StorageProfile: &armpostgresql.StorageProfile{
	// 			BackupRetentionDays: to.Ptr[int32](7),
	// 			GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
	// 			StorageMB: to.Ptr[int32](128000),
	// 		},
	// 		UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
	// 		Version: to.Ptr(armpostgresql.ServerVersionNine6),
	// 	},
	// 	SKU: &armpostgresql.SKU{
	// 		Name: to.Ptr("GP_Gen5_2"),
	// 		Capacity: to.Ptr[int32](2),
	// 		Family: to.Ptr("Gen5"),
	// 		Tier: to.Ptr(armpostgresql.SKUTierGeneralPurpose),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerUpdate.json
func ExampleServersClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServersClient().BeginUpdate(ctx, "testrg", "pgtestsvc4", armpostgresql.ServerUpdateParameters{
		Properties: &armpostgresql.ServerUpdateParametersProperties{
			AdministratorLoginPassword: to.Ptr("<administratorLoginPassword>"),
			MinimalTLSVersion:          to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
			SSLEnforcement:             to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		},
	}, nil)
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
	// res.Server = armpostgresql.Server{
	// 	Name: to.Ptr("pgtestsvc4"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc4"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"ElasticServer": to.Ptr("1"),
	// 	},
	// 	Properties: &armpostgresql.ServerProperties{
	// 		AdministratorLogin: to.Ptr("cloudsa"),
	// 		EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-14T21:08:24.637Z"); return t}()),
	// 		FullyQualifiedDomainName: to.Ptr("pgtestsvc4.postgres.database.azure.com"),
	// 		MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
	// 		SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
	// 		StorageProfile: &armpostgresql.StorageProfile{
	// 			BackupRetentionDays: to.Ptr[int32](7),
	// 			GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
	// 			StorageMB: to.Ptr[int32](128000),
	// 		},
	// 		UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
	// 		Version: to.Ptr(armpostgresql.ServerVersionNine6),
	// 	},
	// 	SKU: &armpostgresql.SKU{
	// 		Name: to.Ptr("B_Gen4_2"),
	// 		Capacity: to.Ptr[int32](2),
	// 		Family: to.Ptr("Gen4"),
	// 		Tier: to.Ptr(armpostgresql.SKUTierBasic),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerDelete.json
func ExampleServersClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServersClient().BeginDelete(ctx, "TestGroup", "testserver", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerGet.json
func ExampleServersClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServersClient().Get(ctx, "testrg", "pgtestsvc1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Server = armpostgresql.Server{
	// 	Name: to.Ptr("pgtestsvc1"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc1"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armpostgresql.ServerProperties{
	// 		AdministratorLogin: to.Ptr("testuser"),
	// 		EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-04T21:00:58.924Z"); return t}()),
	// 		FullyQualifiedDomainName: to.Ptr("pgtestsvc1.postgres.database.azure.com"),
	// 		MasterServerID: to.Ptr(""),
	// 		MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
	// 		PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
	// 			{
	// 				ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc1/privateEndpointConnections/private-endpoint-name-00000000-1111-2222-3333-444444444444"),
	// 				Properties: &armpostgresql.ServerPrivateEndpointConnectionProperties{
	// 					PrivateEndpoint: &armpostgresql.PrivateEndpointProperty{
	// 						ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/Default-Network/providers/Microsoft.Network/privateEndpoints/private-endpoint-name"),
	// 					},
	// 					PrivateLinkServiceConnectionState: &armpostgresql.ServerPrivateLinkServiceConnectionStateProperty{
	// 						Description: to.Ptr("Auto-approved"),
	// 						ActionsRequired: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateActionsRequireNone),
	// 						Status: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateStatusApproved),
	// 					},
	// 					ProvisioningState: to.Ptr(armpostgresql.PrivateEndpointProvisioningState("Succeeded")),
	// 				},
	// 		}},
	// 		PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
	// 		ReplicationRole: to.Ptr(""),
	// 		SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
	// 		StorageProfile: &armpostgresql.StorageProfile{
	// 			BackupRetentionDays: to.Ptr[int32](10),
	// 			GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
	// 			StorageMB: to.Ptr[int32](5120),
	// 		},
	// 		UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
	// 		Version: to.Ptr(armpostgresql.ServerVersionNine5),
	// 	},
	// 	SKU: &armpostgresql.SKU{
	// 		Name: to.Ptr("B_Gen4_1"),
	// 		Capacity: to.Ptr[int32](1),
	// 		Family: to.Ptr("Gen4"),
	// 		Tier: to.Ptr(armpostgresql.SKUTierBasic),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerListByResourceGroup.json
func ExampleServersClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewServersClient().NewListByResourceGroupPager("TestGroup", nil)
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
		// page.ServerListResult = armpostgresql.ServerListResult{
		// 	Value: []*armpostgresql.Server{
		// 		{
		// 			Name: to.Ptr("pgtestsvc1"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc1"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armpostgresql.ServerProperties{
		// 				AdministratorLogin: to.Ptr("testuser"),
		// 				EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-04T21:01:55.149Z"); return t}()),
		// 				FullyQualifiedDomainName: to.Ptr("pgtestsvc1.postgres.database.azure.com"),
		// 				PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
		// 				},
		// 				PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
		// 				SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		// 				StorageProfile: &armpostgresql.StorageProfile{
		// 					BackupRetentionDays: to.Ptr[int32](10),
		// 					GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
		// 					StorageMB: to.Ptr[int32](5120),
		// 				},
		// 				UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
		// 				Version: to.Ptr(armpostgresql.ServerVersionNine5),
		// 			},
		// 			SKU: &armpostgresql.SKU{
		// 				Name: to.Ptr("B_Gen4_1"),
		// 				Capacity: to.Ptr[int32](1),
		// 				Family: to.Ptr("Gen4"),
		// 				Tier: to.Ptr(armpostgresql.SKUTierBasic),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("pgtestsvc2"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc2"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armpostgresql.ServerProperties{
		// 				AdministratorLogin: to.Ptr("testuser"),
		// 				EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-07T21:01:55.149Z"); return t}()),
		// 				FullyQualifiedDomainName: to.Ptr("pgtestsvc2.postgres.database.azure.com"),
		// 				PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
		// 					{
		// 						ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc2/privateEndpointConnections/private-endpoint-name-00000000-1111-2222-3333-444444444444"),
		// 						Properties: &armpostgresql.ServerPrivateEndpointConnectionProperties{
		// 							PrivateEndpoint: &armpostgresql.PrivateEndpointProperty{
		// 								ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/Default-Network/providers/Microsoft.Network/privateEndpoints/private-endpoint-name"),
		// 							},
		// 							PrivateLinkServiceConnectionState: &armpostgresql.ServerPrivateLinkServiceConnectionStateProperty{
		// 								Description: to.Ptr("Auto-approved"),
		// 								ActionsRequired: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateActionsRequireNone),
		// 								Status: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateStatusApproved),
		// 							},
		// 							ProvisioningState: to.Ptr(armpostgresql.PrivateEndpointProvisioningState("Succeeded")),
		// 						},
		// 				}},
		// 				PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
		// 				SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		// 				StorageProfile: &armpostgresql.StorageProfile{
		// 					BackupRetentionDays: to.Ptr[int32](7),
		// 					GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
		// 					StorageMB: to.Ptr[int32](5120),
		// 				},
		// 				UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
		// 				Version: to.Ptr(armpostgresql.ServerVersionNine6),
		// 			},
		// 			SKU: &armpostgresql.SKU{
		// 				Name: to.Ptr("GP_Gen4_2"),
		// 				Capacity: to.Ptr[int32](2),
		// 				Family: to.Ptr("Gen4"),
		// 				Tier: to.Ptr(armpostgresql.SKUTierGeneralPurpose),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("pgtestsvc4"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc4"),
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"ElasticServer": to.Ptr("1"),
		// 			},
		// 			Properties: &armpostgresql.ServerProperties{
		// 				AdministratorLogin: to.Ptr("cloudsa"),
		// 				EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-14T21:08:24.637Z"); return t}()),
		// 				FullyQualifiedDomainName: to.Ptr("pgtestsvc4.postgres.database.azure.com"),
		// 				PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
		// 				},
		// 				PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
		// 				SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		// 				StorageProfile: &armpostgresql.StorageProfile{
		// 					BackupRetentionDays: to.Ptr[int32](7),
		// 					GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
		// 					StorageMB: to.Ptr[int32](128000),
		// 				},
		// 				UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
		// 				Version: to.Ptr(armpostgresql.ServerVersionNine6),
		// 			},
		// 			SKU: &armpostgresql.SKU{
		// 				Name: to.Ptr("B_Gen4_2"),
		// 				Capacity: to.Ptr[int32](2),
		// 				Family: to.Ptr("Gen4"),
		// 				Tier: to.Ptr(armpostgresql.SKUTierBasic),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerList.json
func ExampleServersClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewServersClient().NewListPager(nil)
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
		// page.ServerListResult = armpostgresql.ServerListResult{
		// 	Value: []*armpostgresql.Server{
		// 		{
		// 			Name: to.Ptr("pgtestsvc1"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc1"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armpostgresql.ServerProperties{
		// 				AdministratorLogin: to.Ptr("testuser"),
		// 				EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-04T21:01:55.149Z"); return t}()),
		// 				FullyQualifiedDomainName: to.Ptr("pgtestsvc1.postgres.database.azure.com"),
		// 				MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS11),
		// 				PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
		// 				},
		// 				PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
		// 				SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		// 				StorageProfile: &armpostgresql.StorageProfile{
		// 					BackupRetentionDays: to.Ptr[int32](10),
		// 					GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
		// 					StorageMB: to.Ptr[int32](5120),
		// 				},
		// 				UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
		// 				Version: to.Ptr(armpostgresql.ServerVersionNine5),
		// 			},
		// 			SKU: &armpostgresql.SKU{
		// 				Name: to.Ptr("B_Gen4_1"),
		// 				Capacity: to.Ptr[int32](1),
		// 				Family: to.Ptr("Gen4"),
		// 				Tier: to.Ptr(armpostgresql.SKUTierBasic),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("pgtestsvc2"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc2"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armpostgresql.ServerProperties{
		// 				AdministratorLogin: to.Ptr("testuser"),
		// 				EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-07T21:01:55.149Z"); return t}()),
		// 				FullyQualifiedDomainName: to.Ptr("pgtestsvc2.postgres.database.azure.com"),
		// 				MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS12),
		// 				PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
		// 					{
		// 						ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc2/privateEndpointConnections/private-endpoint-name-00000000-1111-2222-3333-444444444444"),
		// 						Properties: &armpostgresql.ServerPrivateEndpointConnectionProperties{
		// 							PrivateEndpoint: &armpostgresql.PrivateEndpointProperty{
		// 								ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/Default-Network/providers/Microsoft.Network/privateEndpoints/private-endpoint-name"),
		// 							},
		// 							PrivateLinkServiceConnectionState: &armpostgresql.ServerPrivateLinkServiceConnectionStateProperty{
		// 								Description: to.Ptr("Auto-approved"),
		// 								ActionsRequired: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateActionsRequireNone),
		// 								Status: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateStatusApproved),
		// 							},
		// 							ProvisioningState: to.Ptr(armpostgresql.PrivateEndpointProvisioningState("Succeeded")),
		// 						},
		// 				}},
		// 				PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
		// 				SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		// 				StorageProfile: &armpostgresql.StorageProfile{
		// 					BackupRetentionDays: to.Ptr[int32](7),
		// 					GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
		// 					StorageMB: to.Ptr[int32](5120),
		// 				},
		// 				UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
		// 				Version: to.Ptr(armpostgresql.ServerVersionNine6),
		// 			},
		// 			SKU: &armpostgresql.SKU{
		// 				Name: to.Ptr("GP_Gen4_2"),
		// 				Capacity: to.Ptr[int32](2),
		// 				Family: to.Ptr("Gen4"),
		// 				Tier: to.Ptr(armpostgresql.SKUTierGeneralPurpose),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("pgtestsvc3"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg1/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc3"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armpostgresql.ServerProperties{
		// 				AdministratorLogin: to.Ptr("testuser"),
		// 				EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-01T00:11:08.550Z"); return t}()),
		// 				FullyQualifiedDomainName: to.Ptr("pgtestsvc3.postgres.database.azure.com"),
		// 				MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS10),
		// 				PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
		// 					{
		// 						ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc3/privateEndpointConnections/private-endpoint-name-00000000-1111-2222-3333-444444444444"),
		// 						Properties: &armpostgresql.ServerPrivateEndpointConnectionProperties{
		// 							PrivateEndpoint: &armpostgresql.PrivateEndpointProperty{
		// 								ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/Default-Network/providers/Microsoft.Network/privateEndpoints/private-endpoint-name"),
		// 							},
		// 							PrivateLinkServiceConnectionState: &armpostgresql.ServerPrivateLinkServiceConnectionStateProperty{
		// 								Description: to.Ptr("Auto-approved"),
		// 								ActionsRequired: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateActionsRequireNone),
		// 								Status: to.Ptr(armpostgresql.PrivateLinkServiceConnectionStateStatusApproved),
		// 							},
		// 							ProvisioningState: to.Ptr(armpostgresql.PrivateEndpointProvisioningState("Succeeded")),
		// 						},
		// 				}},
		// 				PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
		// 				SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		// 				StorageProfile: &armpostgresql.StorageProfile{
		// 					BackupRetentionDays: to.Ptr[int32](35),
		// 					GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupEnabled),
		// 					StorageMB: to.Ptr[int32](204800),
		// 				},
		// 				UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
		// 				Version: to.Ptr(armpostgresql.ServerVersionNine6),
		// 			},
		// 			SKU: &armpostgresql.SKU{
		// 				Name: to.Ptr("GP_Gen4_4"),
		// 				Capacity: to.Ptr[int32](4),
		// 				Family: to.Ptr("Gen4"),
		// 				Tier: to.Ptr(armpostgresql.SKUTierGeneralPurpose),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("pgtestsvc4"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/servers"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/servers/pgtestsvc4"),
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"ElasticServer": to.Ptr("1"),
		// 			},
		// 			Properties: &armpostgresql.ServerProperties{
		// 				AdministratorLogin: to.Ptr("cloudsa"),
		// 				EarliestRestoreDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-03-14T21:08:24.637Z"); return t}()),
		// 				FullyQualifiedDomainName: to.Ptr("pgtestsvc4.postgres.database.azure.com"),
		// 				MinimalTLSVersion: to.Ptr(armpostgresql.MinimalTLSVersionEnumTLS10),
		// 				PrivateEndpointConnections: []*armpostgresql.ServerPrivateEndpointConnection{
		// 				},
		// 				PublicNetworkAccess: to.Ptr(armpostgresql.PublicNetworkAccessEnumEnabled),
		// 				SSLEnforcement: to.Ptr(armpostgresql.SSLEnforcementEnumEnabled),
		// 				StorageProfile: &armpostgresql.StorageProfile{
		// 					BackupRetentionDays: to.Ptr[int32](7),
		// 					GeoRedundantBackup: to.Ptr(armpostgresql.GeoRedundantBackupDisabled),
		// 					StorageMB: to.Ptr[int32](128000),
		// 				},
		// 				UserVisibleState: to.Ptr(armpostgresql.ServerStateReady),
		// 				Version: to.Ptr(armpostgresql.ServerVersionNine6),
		// 			},
		// 			SKU: &armpostgresql.SKU{
		// 				Name: to.Ptr("B_Gen4_2"),
		// 				Capacity: to.Ptr[int32](2),
		// 				Family: to.Ptr("Gen4"),
		// 				Tier: to.Ptr(armpostgresql.SKUTierBasic),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c767823fdfd9d5e96bad245e3ea4d14d94a716bb/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2017-12-01/examples/ServerRestart.json
func ExampleServersClient_BeginRestart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServersClient().BeginRestart(ctx, "TestGroup", "testserver", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
