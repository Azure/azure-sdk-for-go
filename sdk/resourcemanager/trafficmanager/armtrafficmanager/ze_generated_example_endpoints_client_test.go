//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armtrafficmanager_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/trafficmanager/resource-manager/Microsoft.Network/stable/2018-08-01/examples/Endpoint-PATCH-External-Target.json
func ExampleEndpointsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armtrafficmanager.NewEndpointsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"<resource-group-name>",
		"<profile-name>",
		armtrafficmanager.EndpointTypeExternalEndpoints,
		"<endpoint-name>",
		armtrafficmanager.Endpoint{
			Name: to.Ptr("<name>"),
			Type: to.Ptr("<type>"),
			ID:   to.Ptr("<id>"),
			Properties: &armtrafficmanager.EndpointProperties{
				Target: to.Ptr("<target>"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/trafficmanager/resource-manager/Microsoft.Network/stable/2018-08-01/examples/Endpoint-GET-External-WithGeoMapping.json
func ExampleEndpointsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armtrafficmanager.NewEndpointsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<profile-name>",
		armtrafficmanager.EndpointTypeExternalEndpoints,
		"<endpoint-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/trafficmanager/resource-manager/Microsoft.Network/stable/2018-08-01/examples/Endpoint-PUT-External-WithCustomHeaders.json
func ExampleEndpointsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armtrafficmanager.NewEndpointsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx,
		"<resource-group-name>",
		"<profile-name>",
		armtrafficmanager.EndpointTypeExternalEndpoints,
		"<endpoint-name>",
		armtrafficmanager.Endpoint{
			Name: to.Ptr("<name>"),
			Type: to.Ptr("<type>"),
			Properties: &armtrafficmanager.EndpointProperties{
				CustomHeaders: []*armtrafficmanager.EndpointPropertiesCustomHeadersItem{
					{
						Name:  to.Ptr("<name>"),
						Value: to.Ptr("<value>"),
					},
					{
						Name:  to.Ptr("<name>"),
						Value: to.Ptr("<value>"),
					}},
				EndpointLocation: to.Ptr("<endpoint-location>"),
				EndpointStatus:   to.Ptr(armtrafficmanager.EndpointStatusEnabled),
				Target:           to.Ptr("<target>"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/trafficmanager/resource-manager/Microsoft.Network/stable/2018-08-01/examples/Endpoint-DELETE-External.json
func ExampleEndpointsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armtrafficmanager.NewEndpointsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Delete(ctx,
		"<resource-group-name>",
		"<profile-name>",
		armtrafficmanager.EndpointTypeExternalEndpoints,
		"<endpoint-name>",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}
