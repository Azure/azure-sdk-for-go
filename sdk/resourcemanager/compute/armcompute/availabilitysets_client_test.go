//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
)

func cleanup(t *testing.T, client *armresources.ResourceGroupsClient, resourceGroupName string) {
	_, err := client.BeginDelete(context.Background(), resourceGroupName, nil)
	require.NoError(t, err)
}

func TestAvailabilitySetsClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	location := recording.GetEnvVariable("AZURE_LOCATION", "westus2")

	// create resource group
	rgName, err := createRandomName(t, "test-compute")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID, cred, opt)
	_, err = rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: &location,
	}, nil)
	defer cleanup(t, rgClient, rgName)
	require.NoError(t, err)
	client := armcompute.NewAvailabilitySetsClient(subscriptionID, cred, opt)
	name, err := createRandomName(t, "set")
	require.NoError(t, err)
	resp, err := client.CreateOrUpdate(
		context.Background(),
		rgName,
		name,
		armcompute.AvailabilitySet{
			Resource: armcompute.Resource{
				Location: &location,
			},
			SKU: &armcompute.SKU{
				Name: to.StringPtr(string(armcompute.AvailabilitySetSKUTypesAligned)),
			},
			Properties: &armcompute.AvailabilitySetProperties{
				PlatformFaultDomainCount:  to.Int32Ptr(1),
				PlatformUpdateDomainCount: to.Int32Ptr(1),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, *resp.Name, name)
}
