//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	pathToPackage = "sdk/resourcemanager/internal/testdata"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("test", "test")
	require.Equal(t, GetEnv("test", ""), "test")
	require.Equal(t, GetEnv("testfail", "fail"), "fail")
}

func TestCreateDeleteResourceGroup(t *testing.T) {
	ctx := context.Background()
	cred, options := GetCredAndClientOptions(t)
	subscriptionID := GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	StartRecording(t, pathToPackage)
	resourceGroup, _, err := CreateResourceGroup(ctx, subscriptionID, cred, options, "eastus")
	require.NoError(t, err)
	require.Equal(t, strings.HasPrefix(*resourceGroup.Name, "go-sdk-test-"), true)
	_, err = DeleteResourceGroup(ctx, subscriptionID, cred, options, *resourceGroup.Name)
	require.NoError(t, err)
	StopRecording(t)
}
