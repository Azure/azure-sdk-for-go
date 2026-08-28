// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/data/azappconfig/testdata"

var (
	credential   azcore.TokenCredential
	endpoint     string
	fakeEndpoint = fmt.Sprintf("https://%s.azconfig.io", recording.SanitizedValue)
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	var err error
	credential, err = azcred.New(nil)
	if err != nil {
		panic(err)
	}

	endpoint = recording.GetEnvVariable("APPCONFIGURATION_ENDPOINT_STRING", fakeEndpoint)

	if recording.GetRecordMode() != recording.LiveMode {
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

		err = recording.RemoveRegisteredSanitizers([]string{
			"AZSDK2030", // operation-location header
			"AZSDK3424", // $..to (used by feature flag PercentileAllocation bounds)
			"AZSDK3425", // $..from (used by feature flag PercentileAllocation bounds)
			"AZSDK3447", // $.key
			"AZSDK3490", // $..etag
			"AZSDK3493", // $..name
		}, nil)
		if err != nil {
			panic(err)
		}

		if err := recording.AddHeaderRegexSanitizer("Operation-Location", fakeEndpoint, `https://\w+\.azconfig\.io`, nil); err != nil {
			panic(err)
		}
	}

	return m.Run()
}

// newTestClient constructs an [*azappconfig.Client] authenticated with Entra ID and wired to the test
// recording transport. In playback mode it uses a fake credential; in live/record mode it uses the
// shared credential resolved by [azcred.New].
func newTestClient(t *testing.T) *azappconfig.Client {
	if recording.GetRecordMode() != recording.PlaybackMode && os.Getenv("APPCONFIGURATION_ENDPOINT_STRING") == "" {
		t.Skip("set APPCONFIGURATION_ENDPOINT_STRING to run this test")
	}

	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)

	err = recording.SetDefaultMatcher(t, &recording.SetDefaultMatcherOptions{
		IgnoredQueryParameters: []string{"api-version"},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	client, err := azappconfig.NewClient(endpoint, credential, &azappconfig.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
			Logging: policy.LogOptions{
				IncludeBody: true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	return client
}
