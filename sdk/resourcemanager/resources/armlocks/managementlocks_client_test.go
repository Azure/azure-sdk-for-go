//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armlocks_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armlocks"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestManagementLocksClient_CreateOrUpdateAtSubscriptionLevel(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	locksClient := armlocks.NewManagementLocksClient(subscriptionID, cred, opt)
	ctx := context.Background()
	lockName, err := createRandomName(t, "lock")
	require.NoError(t, err)
	subLevel, err := locksClient.CreateOrUpdateAtSubscriptionLevel(
		ctx,
		lockName,
		armlocks.ManagementLockObject{
			Properties: &armlocks.ManagementLockProperties{
				Level: armlocks.LockLevelCanNotDelete.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, lockName, *subLevel.Name)
}

func TestManagementLocksClient_CreateOrUpdateByScope(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()

	resourceName, _ := createRandomName(t, "compute")
	resourceID := "/subscriptions/{guid}/resourceGroups/{resourcegroupname}/providers/{resourceprovidernamespace}/{resourcetype}/{resourcename}"
	resourceID = strings.ReplaceAll(resourceID, "{guid}", subscriptionID)
	resourceID = strings.ReplaceAll(resourceID, "{resourcegroupname}", *rg.Name)
	resourceID = strings.ReplaceAll(resourceID, "{resourceprovidernamespace}", "Microsoft.Compute")
	resourceID = strings.ReplaceAll(resourceID, "{resourcetype}", "availabilitySets")
	resourceID = strings.ReplaceAll(resourceID, "{resourcename}", resourceName)

	resourcesClient := armresources.NewClient(subscriptionID, cred, opt)
	resourcePollerResp, err := resourcesClient.BeginCreateOrUpdateByID(
		ctx,
		resourceID,
		"2019-07-01",
		armresources.GenericResource{
			Location: to.StringPtr(location),
		},
		nil,
	)
	require.NoError(t, err)
	resourceResp, err := resourcePollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, resourceName, *resourceResp.Name)

	locksClient := armlocks.NewManagementLocksClient(subscriptionID, cred, opt)

	lockName, err := createRandomName(t, "lock")
	require.NoError(t, err)
	subLevel, err := locksClient.CreateOrUpdateByScope(
		ctx,
		resourceID,
		lockName,
		armlocks.ManagementLockObject{
			Properties: &armlocks.ManagementLockProperties{
				Level: armlocks.LockLevelCanNotDelete.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, lockName, *subLevel.Name)
}
