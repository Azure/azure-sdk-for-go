//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armsubscriptions_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubscriptionClient_CheckResourceName(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionClient := armsubscriptions.NewSubscriptionClient(cred, opt)

	randomName, err := createRandomName(t, "sub")
	require.NoError(t, err)

	resourceName, err := subscriptionClient.CheckResourceName(context.Background(), &armsubscriptions.SubscriptionClientCheckResourceNameOptions{
		ResourceNameDefinition: &armsubscriptions.ResourceName{
			Name: to.StringPtr(randomName),
			Type: to.StringPtr("Microsoft.Compute"),
		},
	})
	require.NoError(t, err)
	require.Equal(t, *resourceName.Status, armsubscriptions.ResourceNameStatusAllowed)
}
