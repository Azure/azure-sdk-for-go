//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package backup_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/security/keyvault/azadmin/testdata"
const fakeHsmURL = "https://fakehsm.managedhsm.azure.net/"
const fakeBlobURL = "https://fakestorageaccount.blob.core.windows.net/backup"
const fakeToken = "fakeSasToken"

var (
	credential azcore.TokenCredential
	hsmURL     string
	token      string
	blobURL    string
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
	blobURL = internal.GetEnvVar("BLOB_CONTAINER_URL", fakeBlobURL)
	token = internal.GetEnvVar("BLOB_STORAGE_SAS_TOKEN", fakeToken)

	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddBodyRegexSanitizer(fakeToken, `sv=[^"]*`, nil)
		if err != nil {
			panic(err)
		}
	}

	return m.Run()
}

func startBackupTest(t *testing.T) (*backup.Client, backup.SASTokenParameters) {
	internal.StartRecording(t, recordingDirectory)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &backup.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := backup.NewClient(hsmURL, credential, opts)
	require.NoError(t, err)
	sasToken := backup.SASTokenParameters{
		StorageResourceURI: &blobURL,
		Token:              &token,
	}

	return client, sasToken
}
