//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azwebpubsub_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
	"github.com/stretchr/testify/require"
)

func TestHealthClient_GetServiceStatus(t *testing.T) {
	client := newHealthClient(t)
	_, err := client.GetServiceStatus(context.Background(), &azwebpubsub.HealthClientGetServiceStatusOptions{})
	require.NoError(t, err)
}

func newHealthClient(t *testing.T) *azwebpubsub.HealthClient {
	tv, coreOptions := loadClientOptions(t)
	options := &azwebpubsub.ClientOptions{
		ClientOptions: *coreOptions,
	}
	println(tv.Endpoint)
	client, err := azwebpubsub.NewHealthClientWithNoCredential(tv.Endpoint, options)
	require.NoError(t, err)
	return client
}
