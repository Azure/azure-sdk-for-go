//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDedicatedHostsClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"dhcreate","eastus")
	rgName := *rg.Name
	defer clean()

	// create dedicated host group
	dhgClient := armcompute.NewDedicatedHostGroupsClient(subscriptionID,cred,opt)
	dhgName, err := createRandomName(t, "createDH")
	require.NoError(t, err)
	dhgResp, err := dhgClient.CreateOrUpdate(
		context.Background(),
		rgName,
		dhgName,
		armcompute.DedicatedHostGroup{
			Resource: armcompute.Resource{
				Location: to.StringPtr("eastus"),
			},
			Properties: &armcompute.DedicatedHostGroupProperties{
				PlatformFaultDomainCount: to.Int32Ptr(3),
			},
			Zones: []*string{to.StringPtr("1")},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, *dhgResp.Name, dhgName)

	// create dedicated host
	dhClient := armcompute.NewDedicatedHostsClient(subscriptionID,cred,nil)
	dhName, err := createRandomName(t, "dh")
	require.NoError(t, err)
	dhPoller, err := dhClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		dhgName,
		dhName,
		armcompute.DedicatedHost{
			Resource: armcompute.Resource{
				Location: to.StringPtr("eastus"),
			},
			Properties: &armcompute.DedicatedHostProperties{
				PlatformFaultDomain: to.Int32Ptr(1),
			},
			SKU: &armcompute.SKU{
				Name: to.StringPtr("DSv3-Type1"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dhResp, err := dhPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *dhResp.Name, dhName)
}
