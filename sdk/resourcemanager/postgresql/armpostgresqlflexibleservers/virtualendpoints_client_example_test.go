//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpostgresqlflexibleservers_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/v4"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ecee919199a39cc0d864410f540aa105bf7cdb64/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/VirtualEndpointCreate.json
func ExampleVirtualEndpointsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualEndpointsClient().BeginCreate(ctx, "testrg", "pgtestsvc4", "pgVirtualEndpoint1", armpostgresqlflexibleservers.VirtualEndpointResource{
		Properties: &armpostgresqlflexibleservers.VirtualEndpointResourceProperties{
			EndpointType: to.Ptr(armpostgresqlflexibleservers.VirtualEndpointTypeReadWrite),
			Members: []*string{
				to.Ptr("testPrimary1")},
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
	// res.VirtualEndpointResource = armpostgresqlflexibleservers.VirtualEndpointResource{
	// 	Name: to.Ptr("pgVirtualEndpoint1"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/flexibleServers/pgtestsvc4/virtualEndpoints/pgVirtualEndpoint1"),
	// 	Properties: &armpostgresqlflexibleservers.VirtualEndpointResourceProperties{
	// 		EndpointType: to.Ptr(armpostgresqlflexibleservers.VirtualEndpointTypeReadWrite),
	// 		Members: []*string{
	// 			to.Ptr("testPrimary1")},
	// 			VirtualEndpoints: []*string{
	// 				to.Ptr("pgVirtualEndpoint1.reader.postgres.database.azure.com"),
	// 				to.Ptr("pgVirtualEndpoint1.writer.postgres.database.azure.com")},
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ecee919199a39cc0d864410f540aa105bf7cdb64/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/VirtualEndpointUpdate.json
func ExampleVirtualEndpointsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualEndpointsClient().BeginUpdate(ctx, "testrg", "pgtestsvc4", "pgVirtualEndpoint1", armpostgresqlflexibleservers.VirtualEndpointResourceForPatch{
		Properties: &armpostgresqlflexibleservers.VirtualEndpointResourceProperties{
			EndpointType: to.Ptr(armpostgresqlflexibleservers.VirtualEndpointTypeReadWrite),
			Members: []*string{
				to.Ptr("testReplica1")},
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
	// res.VirtualEndpointResource = armpostgresqlflexibleservers.VirtualEndpointResource{
	// 	Name: to.Ptr("pgVirtualEndpoint1"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/flexibleServers/pgtestsvc4/virtualEndpoints/pgVirtualEndpoint1"),
	// 	Properties: &armpostgresqlflexibleservers.VirtualEndpointResourceProperties{
	// 		EndpointType: to.Ptr(armpostgresqlflexibleservers.VirtualEndpointTypeReadWrite),
	// 		Members: []*string{
	// 			to.Ptr("testReplica1")},
	// 			VirtualEndpoints: []*string{
	// 				to.Ptr("pgVirtualEndpoint1.reader.postgres.database.azure.com"),
	// 				to.Ptr("pgVirtualEndpoint1.writer.postgres.database.azure.com")},
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ecee919199a39cc0d864410f540aa105bf7cdb64/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/VirtualEndpointDelete.json
func ExampleVirtualEndpointsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualEndpointsClient().BeginDelete(ctx, "testrg", "pgtestsvc4", "pgVirtualEndpoint1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ecee919199a39cc0d864410f540aa105bf7cdb64/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/VirtualEndpointsGet.json
func ExampleVirtualEndpointsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualEndpointsClient().Get(ctx, "testrg", "pgtestsvc4", "pgVirtualEndpoint1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualEndpointResource = armpostgresqlflexibleservers.VirtualEndpointResource{
	// 	Name: to.Ptr("pgVirtualEndpoint1"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/flexibleServers/pgtestsvc4/virtualEndpoints/pgVirtualEndpoint1"),
	// 	Properties: &armpostgresqlflexibleservers.VirtualEndpointResourceProperties{
	// 		EndpointType: to.Ptr(armpostgresqlflexibleservers.VirtualEndpointTypeReadWrite),
	// 		Members: []*string{
	// 			to.Ptr("testReplica1")},
	// 			VirtualEndpoints: []*string{
	// 				to.Ptr("pgVirtualEndpoint1.reader.postgres.database.azure.com"),
	// 				to.Ptr("pgVirtualEndpoint1.writer.postgres.database.azure.com")},
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ecee919199a39cc0d864410f540aa105bf7cdb64/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/VirtualEndpointsListByServer.json
func ExampleVirtualEndpointsClient_NewListByServerPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualEndpointsClient().NewListByServerPager("testrg", "pgtestsvc4", nil)
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
		// page.VirtualEndpointsListResult = armpostgresqlflexibleservers.VirtualEndpointsListResult{
		// 	Value: []*armpostgresqlflexibleservers.VirtualEndpointResource{
		// 		{
		// 			Name: to.Ptr("pgVirtualEndpoint1"),
		// 			Type: to.Ptr("Microsoft.DBforPostgreSQL/flexibleServers/virtualEndpoints"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforPostgreSQL/flexibleServers/pgtestsvc4/virtualEndpoints/pgVirtualEndpoint1"),
		// 			Properties: &armpostgresqlflexibleservers.VirtualEndpointResourceProperties{
		// 				EndpointType: to.Ptr(armpostgresqlflexibleservers.VirtualEndpointTypeReadWrite),
		// 				Members: []*string{
		// 					to.Ptr("testReplica1")},
		// 					VirtualEndpoints: []*string{
		// 						to.Ptr("pgVirtualEndpoint1.reader.postgres.database.azure.com"),
		// 						to.Ptr("pgVirtualEndpoint1.writer.postgres.database.azure.com")},
		// 					},
		// 			}},
		// 		}
	}
}
