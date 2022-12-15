//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azadmin_test

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin"
	"github.com/stretchr/testify/require"
)

const fakeHsmURL = "https://fakehsm.managedhsm.azure.net/"

var (
	credential azcore.TokenCredential
	hsmURL     string
)

func TestMain(m *testing.M) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		//TODO - doublecheck env var name
		hsmURL = os.Getenv("AZURE_MANAGEDHSM_URL")
	}
	if hsmURL == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for AZURE_MANAGEDHSM_URL")
		}
		hsmURL = fakeHsmURL
	}

	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		credential, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			panic(err)
		}
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeHsmURL, hsmURL, nil)
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
	err := recording.Start(t, "sdk/keyvault/azadmin/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startAccessControlTest(t *testing.T) *azadmin.AccessControlClient {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azadmin.AccessControlClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := azadmin.NewAccessControlClient(hsmURL, credential, opts)
	require.NoError(t, err)
	return client
}

func startSettingsTest(t *testing.T) *azadmin.SettingsClient {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azadmin.SettingsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := azadmin.NewSettingsClient(hsmURL, credential, opts)
	require.NoError(t, err)
	return client
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
