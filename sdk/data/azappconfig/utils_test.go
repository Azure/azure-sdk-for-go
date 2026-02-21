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
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/data/azappconfig/testdata"

var (
	fakeConnStr  = fmt.Sprintf("Endpoint=%s;Id=fake;Secret=fake", fakeEndpoint)
	fakeEndpoint = fmt.Sprintf("https://%s.azconfig.io", recording.SanitizedValue)
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
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
			"AZSDK3447", // $.key
			"AZSDK3490", // $..etag
			"AZSDK3493", // $..name
		}, nil)
		if err != nil {
			panic(err)
		}

		if err := recording.AddHeaderRegexSanitizer("x-ms-content-sha256", "fake-content", "", nil); err != nil {
			panic(err)
		}

		if err := recording.AddHeaderRegexSanitizer("Operation-Location", fakeEndpoint, `https://\w+\.azconfig\.io`, nil); err != nil {
			panic(err)
		}
	}

	return m.Run()
}

func NewClientFromConnectionString(t *testing.T) *azappconfig.Client {
	connStr := recording.GetEnvVariable("APPCONFIGURATION_CONNECTION_STRING", fakeConnStr)
	if connStr == "" && recording.GetRecordMode() != recording.PlaybackMode {
		t.Skip("set APPCONFIGURATION_CONNECTION_STRING to run this test")
	}

	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	client, err := azappconfig.NewClientFromConnectionString(connStr, &azappconfig.ClientOptions{
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
