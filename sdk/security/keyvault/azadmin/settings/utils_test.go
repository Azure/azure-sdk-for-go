//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package settings_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/settings"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/security/keyvault/azadmin/testdata"
const fakeHsmURL = "https://fakehsm.managedhsm.azure.net/"

var (
	credential azcore.TokenCredential
	hsmURL     string
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

	credential = internal.GetCredential("AZADMIN")

	hsmURL = internal.GetEnvVar("AZURE_MANAGEDHSM_URL", fakeHsmURL)

	return m.Run()
}

func startSettingsTest(t *testing.T) *settings.Client {
	internal.StartRecording(t, recordingDirectory)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &settings.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := settings.NewClient(hsmURL, credential, opts)
	require.NoError(t, err)
	return client
}
