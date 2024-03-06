//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/data/azappconfig/testdata"
const fakeConnStr = "Endpoint=https://contoso.azconfig.io;Id=fake-id:fake-value;Secret=ZmFrZS1zZWNyZXQ="

func TestMain(m *testing.M) {
	os.Exit(run(m))
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

	switch recording.GetRecordMode() {
	case recording.PlaybackMode:
		if err := recording.SetDefaultMatcher(nil, &recording.SetDefaultMatcherOptions{
			ExcludedHeaders: []string{"Date", "X-Ms-Content-Sha256"},
		}); err != nil {
			panic(err)
		}

	case recording.RecordingMode:
		if err := recording.AddURISanitizer("https://contoso.azconfig.io", `https://\w+\.azconfig\.io`, nil); err != nil {
			panic(err)
		}

		if err := recording.AddHeaderRegexSanitizer("x-ms-content-sha256", "fake-content", "", nil); err != nil {
			panic(err)
		}

		if err := recording.AddHeaderRegexSanitizer("Operation-Location", "https://contoso.azconfig.io", `https://\w+\.azconfig\.io`, nil); err != nil {
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
