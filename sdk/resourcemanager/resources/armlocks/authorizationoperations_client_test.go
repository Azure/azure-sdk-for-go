//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armlocks_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armlocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthorizationOperationsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	operationsClient := armlocks.NewAuthorizationOperationsClient(cred, opt)
	pager := operationsClient.List(nil)
	require.NoError(t, pager.Err())
}
