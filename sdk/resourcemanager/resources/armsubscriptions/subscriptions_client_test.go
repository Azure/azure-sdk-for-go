//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armsubscriptions_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubscriptionsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	subscriptionsClient := armsubscriptions.NewClient(cred, opt)
	resp, err := subscriptionsClient.Get(context.Background(), subscriptionID, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.SubscriptionID, subscriptionID)
}

func TestSubscriptionsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)

	subscriptionsClient := armsubscriptions.NewClient(cred, opt)
	resp := subscriptionsClient.List(nil)
	require.NoError(t, resp.Err())
}

func TestSubscriptionsClient_ListLocations(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	subscriptionsClient := armsubscriptions.NewClient(cred, opt)
	resp, err := subscriptionsClient.ListLocations(context.Background(), subscriptionID, nil)
	require.NoError(t, err)
	require.Greater(t, len(resp.Value), 0)
}
