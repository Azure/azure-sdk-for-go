//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOperationsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	ctx := context.Background()

	operationsClient := armstorage.NewOperationsClient(cred,opt)
	resp,err := operationsClient.List(ctx,nil)
	require.NoError(t, err)
	require.Greater(t, len(resp.Value),1)
}