//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rbac_test

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
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac"
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

	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("AZADMIN_TENANT_ID")
		clientID := lookupEnvVar("AZADMIN_CLIENT_ID")
		secret := lookupEnvVar("AZADMIN_CLIENT_SECRET")
		var err error
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}

	hsmURL = getEnvVar("AZURE_MANAGEDHSM_URL", fakeHsmURL)

	if recording.GetRecordMode() != recording.LiveMode {
		err := recording.RemoveRegisteredSanitizers([]string{
			"AZSDK3493", // name in body
			"AZSDK3430", // ID in body
		}, nil)
		if err != nil {
			panic(err)
		}
	}

	return m.Run()
}

func startRecording(t *testing.T) {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startAccessControlTest(t *testing.T) *rbac.Client {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &rbac.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := rbac.NewClient(hsmURL, credential, opts)
	require.NoError(t, err)
	return client
}

func getEnvVar(lookupValue string, fakeValue string) string {
	envVar := fakeValue
	if recording.GetRecordMode() != recording.PlaybackMode {
		envVar = lookupEnvVar(lookupValue)
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeValue, envVar, nil)
		if err != nil {
			panic(err)
		}
	}

	return envVar
}

func lookupEnvVar(s string) string {
	v := os.Getenv(s)
	if v == "" {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return v
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
