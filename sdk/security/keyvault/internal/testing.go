//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package internal

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func GetEnvVar(envVar string, fakeValue string) string {
	// get value
	value := fakeValue
	if recording.GetRecordMode() == recording.LiveMode || recording.GetRecordMode() == recording.RecordingMode {
		value = os.Getenv(envVar)
		if value == "" {
			panic("no value for " + envVar)
		}
	}

	// sanitize value
	if fakeValue != "" && recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeValue, value, nil)
		if err != nil {
			panic(err)
		}
	}

	return value
}

type FakeCredential struct{}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func TestSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)

	// testing unmarshal error scenarios
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)
}

func StartRecording(t *testing.T, recordingDirectory string) {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func GetCredential(module string) azcore.TokenCredential {
	var credential azcore.TokenCredential

	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := GetEnvVar(module+"_TENANT_ID", "")
		clientID := GetEnvVar(module+"_CLIENT_ID", "")
		secret := GetEnvVar(module+"_CLIENT_SECRET", "")
		var err error
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}

	return credential
}
