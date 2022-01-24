//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armsubscriptions_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTenantsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	tenantsClient := armsubscriptions.NewTenantsClient(cred, opt)

	pager := tenantsClient.List(nil)
	require.NoError(t, pager.Err())
}
