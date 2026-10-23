//go:build live

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test is separately selected because it requires a resource key in addition
// to the token-only live suite. It never prints or persists the key.
func TestResourceKeyAuthentication(t *testing.T) {
	endpoint := os.Getenv("CONTENTUNDERSTANDING_ENDPOINT")
	key := os.Getenv("CONTENTUNDERSTANDING_KEY")
	require.NotEmpty(t, endpoint, "CONTENTUNDERSTANDING_ENDPOINT is required")
	require.NotEmpty(t, key, "CONTENTUNDERSTANDING_KEY is required")
	credential := azcore.NewKeyCredential(key)
	client, err := sdk.NewClientWithKeyCredential(endpoint, credential, nil)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	defaults, err := client.GetDefaults(ctx, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, defaults.ModelDeployments)
	analyzer, err := client.GetAnalyzer(ctx, "prebuilt-layout", nil)
	require.NoError(t, err)
	require.NotNil(t, analyzer.AnalyzerID)
	assert.Equal(t, "prebuilt-layout", *analyzer.AnalyzerID)
}
