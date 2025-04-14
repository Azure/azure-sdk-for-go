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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ecee919199a39cc0d864410f540aa105bf7cdb64/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/PrivateEndpointConnectionUpdate.json
func ExamplePrivateEndpointConnectionClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPrivateEndpointConnectionClient().BeginUpdate(ctx, "Default", "test-svr", "private-endpoint-connection-name.1fa229cd-bf3f-47f0-8c49-afb36723997e", armpostgresqlflexibleservers.PrivateEndpointConnection{
		Properties: &armpostgresqlflexibleservers.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armpostgresqlflexibleservers.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Approved by johndoe@contoso.com"),
				Status:      to.Ptr(armpostgresqlflexibleservers.PrivateEndpointServiceConnectionStatusApproved),
			},
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
	// res.PrivateEndpointConnection = armpostgresqlflexibleservers.PrivateEndpointConnection{
	// 	Name: to.Ptr("private-endpoint-connection-name.1fa229cd-bf3f-47f0-8c49-afb36723997e"),
	// 	Type: to.Ptr("Microsoft.DBforPostgreSQL/flexibleServers/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/Default/providers/Microsoft.DBforPostgreSQL/flexibleServers/test-svr/privateEndpointConnections/private-endpoint-connection-name.1fa229cd-bf3f-47f0-8c49-afb36723997e"),
	// 	Properties: &armpostgresqlflexibleservers.PrivateEndpointConnectionProperties{
	// 		GroupIDs: []*string{
	// 			to.Ptr("postgresqlServer")},
	// 			PrivateEndpoint: &armpostgresqlflexibleservers.PrivateEndpoint{
	// 				ID: to.Ptr("/subscriptions/55555555-6666-7777-8888-999999999999/resourceGroups/Default-Network/providers/Microsoft.Network/privateEndpoints/private-endpoint-name"),
	// 			},
	// 			PrivateLinkServiceConnectionState: &armpostgresqlflexibleservers.PrivateLinkServiceConnectionState{
	// 				Description: to.Ptr("Approved by johndoe@contoso.com"),
	// 				ActionsRequired: to.Ptr("None"),
	// 				Status: to.Ptr(armpostgresqlflexibleservers.PrivateEndpointServiceConnectionStatusApproved),
	// 			},
	// 			ProvisioningState: to.Ptr(armpostgresqlflexibleservers.PrivateEndpointConnectionProvisioningStateSucceeded),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ecee919199a39cc0d864410f540aa105bf7cdb64/specification/postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2024-08-01/examples/PrivateEndpointConnectionDelete.json
func ExamplePrivateEndpointConnectionClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armpostgresqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPrivateEndpointConnectionClient().BeginDelete(ctx, "Default", "test-svr", "private-endpoint-connection-name.1fa229cd-bf3f-47f0-8c49-afb36723997e", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
