// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"

	"github.com/stretchr/testify/require"
)

func TestCreateDeleteResourceGroup(t *testing.T) {
	ctx := context.Background()
	cred, options := GetCredAndClientOptions(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	stop := StartRecording(t, pathToPackage)
	defer stop()
	resourceGroup, _, err := CreateResourceGroup(ctx, subscriptionID, cred, options, "eastus")
	require.NoError(t, err)
	_, err = DeleteResourceGroup(ctx, subscriptionID, cred, options, *resourceGroup.Name)
	require.NoError(t, err)
}
