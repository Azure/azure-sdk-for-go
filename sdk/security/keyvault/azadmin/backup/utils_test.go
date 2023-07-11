//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package backup_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup"
	"github.com/stretchr/testify/require"
)

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
	if recording.GetRecordMode() != recording.PlaybackMode {
		hsmURL = os.Getenv("AZURE_MANAGEDHSM_URL")
		blobURL = "https://" + os.Getenv("BLOB_STORAGE_ACCOUNT_NAME") + ".blob." + os.Getenv("KEYVAULT_STORAGE_ENDPOINT_SUFFIX") + "/" + os.Getenv("BLOB_CONTAINER_NAME")
		token = os.Getenv("BLOB_STORAGE_SAS_TOKEN")
	}
	if hsmURL == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for AZURE_MANAGEDHSM_URL")
		}
		hsmURL = fakeHsmURL
	}
	if blobURL == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for blob url")
		}
		blobURL = fakeBlobURL
	}
	if token == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for BLOB_STORAGE_SAS_TOKEN")
		}
		token = fakeToken
	}

	err = recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("AZADMIN_TENANT_ID")
		clientID := lookupEnvVar("AZADMIN_CLIENT_ID")
		secret := lookupEnvVar("AZADMIN_CLIENT_SECRET")
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeHsmURL, hsmURL, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddGeneralRegexSanitizer(fakeBlobURL, blobURL, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddBodyRegexSanitizer(fakeToken, `sv=[^"]*`, nil)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	code := m.Run()
	os.Exit(code)
}

func startRecording(t *testing.T) {
	err := recording.Start(t, "sdk/security/keyvault/azadmin/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startBackupTest(t *testing.T) (*backup.Client, backup.SASTokenParameters) {
	startRecording(t)
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

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

type FakeCredential struct{}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)

	// testing unmarshal error scenarios
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), ":false}")
	if modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}
