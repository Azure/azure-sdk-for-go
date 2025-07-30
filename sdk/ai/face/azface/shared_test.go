// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/face/azface"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

func getEndpoint() string {
	return recording.GetEnvVariable("FACE_ENDPOINT", "https://fake.cognitiveservices.azure.com/")
}

func getAPIKey() string {
	return recording.GetEnvVariable("FACE_SUBSCRIPTION_KEY", "00000000000000000000000000000000")
}

func newClientForTest(t *testing.T) (*azface.Client, error) {
	endpoint := getEndpoint()
	
	if recording.GetRecordMode() == recording.PlaybackMode {
		// Use API key for recorded tests to avoid token dependencies
		return azface.NewClientWithKey(endpoint, getAPIKey(), nil)
	}
	
	// For live tests, prefer Azure AD authentication
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// Fall back to API key if Azure AD is not available
		t.Logf("Azure AD credential not available, using API key: %v", err)
		return azface.NewClientWithKey(endpoint, getAPIKey(), nil)
	}
	
	return azface.NewClient(endpoint, cred, nil)
}

func newClientForTestWithRecording(t *testing.T) *azface.Client {
	var transport policy.Transporter
	var err error
	
	// Only use recording transport if we're actually recording or have recordings
	if recording.GetRecordMode() == recording.RecordingMode || recording.GetRecordMode() == recording.PlaybackMode {
		transport, err = recording.NewRecordingHTTPClient(t, nil)
		if err != nil {
			t.Logf("Failed to create recording client, using regular client: %v", err)
			// Fall back to regular client
			client, err := newClientForTest(t)
			require.NoError(t, err)
			return client
		}
	}

	var options *azface.ClientOptions
	if transport != nil {
		options = &azface.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: transport,
			},
		}
	}

	endpoint := getEndpoint()
	
	if recording.GetRecordMode() == recording.PlaybackMode {
		// Use fake credential for playback
		client, err := azface.NewClient(endpoint, &credential.Fake{}, options)
		require.NoError(t, err)
		return client
	}
	
	// For recording, use API key authentication for simplicity
	client, err := azface.NewClientWithKey(endpoint, getAPIKey(), options)
	require.NoError(t, err)
	return client
}

func newAdministrationClientForTest(t *testing.T) *azface.AdministrationClient {
	var transport policy.Transporter
	var err error
	
	// Only use recording transport if we're actually recording or have recordings
	if recording.GetRecordMode() == recording.RecordingMode || recording.GetRecordMode() == recording.PlaybackMode {
		transport, err = recording.NewRecordingHTTPClient(t, nil)
		if err != nil {
			t.Logf("Failed to create recording client, using regular client: %v", err)
			// Fall back to regular client creation
			endpoint := getEndpoint()
			cred, err := azidentity.NewDefaultAzureCredential(nil)
			if err != nil {
				client, err := azface.NewAdministrationClientWithKey(endpoint, getAPIKey(), nil)
				require.NoError(t, err)
				return client
			}
			client, err := azface.NewAdministrationClient(endpoint, cred, nil)
			require.NoError(t, err)
			return client
		}
	}

	var options *azface.AdministrationClientOptions
	if transport != nil {
		options = &azface.AdministrationClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: transport,
			},
		}
	}

	endpoint := getEndpoint()
	
	if recording.GetRecordMode() == recording.PlaybackMode {
		// Use fake credential for playback
		client, err := azface.NewAdministrationClient(endpoint, &credential.Fake{}, options)
		require.NoError(t, err)
		return client
	}
	
	// For recording, use API key authentication for simplicity
	client, err := azface.NewAdministrationClientWithKey(endpoint, getAPIKey(), options)
	require.NoError(t, err)
	return client
}

func skipIfNotLive(t *testing.T) {
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skip("skipping live test")
	}
}

func requireLiveOnly(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("test requires live mode")
	}
}