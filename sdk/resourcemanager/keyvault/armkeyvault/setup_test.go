//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var pathToPackage = "sdk/resourcemanager/keyvault/armkeyvault/testdata"

func TestMain(m *testing.M) {
	// Initialize
	if err := recording.AddGeneralRegexSanitizer("00000000-0000-0000-0000-000000000000", `/subscriptions/(?<subsId>[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`, &recording.RecordingOptions{
		UseHTTPS:        true,
		GroupForReplace: "subsId",
	}); err != nil {
		panic(err)
	}
	tenantID := os.Getenv("AZURE_TENANT_ID")
	err := recording.AddBodyKeySanitizer("$.properties.tenantId", "11111111-1111-1111-1111-111111111111", tenantID, nil)
	if err != nil {
		panic(err)
	}
	err = recording.AddBodyKeySanitizer("$.properties.accessPolicies[0].tenantId", "11111111-1111-1111-1111-111111111111", tenantID, nil)
	if err != nil {
		panic(err)
	}
	objectID := os.Getenv("AZURE_OBJECT_ID")
	err = recording.AddBodyKeySanitizer("$.properties.accessPolicies[0].objectId", "22222222-2222-2222-2222-222222222222", objectID, nil)
	if err != nil {
		panic(err)
	}

	// Run
	exitVal := m.Run()

	// cleanup
	os.Exit(exitVal)
}

func startTest(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}
}
