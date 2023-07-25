//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
	"github.com/stretchr/testify/require"
)

var fakeTestVars = testVars{
	Key:          "key",
	Endpoint:     "https://fake.eastus-1.eventgrid.azure.net",
	Topic:        "topic",
	Subscription: "subscription",
}

type testVars struct {
	Key          string
	Endpoint     string
	Topic        string
	Subscription string

	// KeyLogPath is the value of environment "SSLKEYLOGFILE_TEST", which
	// points to a file on disk where we'll write the TLS pre-master-secret.
	// This is useful if you want to trace parts of this test using Wireshark.
	KeyLogPath string
}

func loadEnv() (testVars, error) {
	var missing []string

	get := func(n string) string {
		if v := os.Getenv(n); v == "" {
			missing = append(missing, n)
		}

		return os.Getenv(n)
	}

	tv := testVars{
		Key:          get("EVENTGRID_KEY"),
		Endpoint:     get("EVENTGRID_ENDPOINT"),
		Topic:        get("EVENTGRID_TOPIC"),
		Subscription: get("EVENTGRID_SUBSCRIPTION"),
	}

	if len(missing) > 0 {
		return testVars{}, fmt.Errorf("Missing env variables: %s", strings.Join(missing, ","))
	}

	// Setting this variable will cause the test clients to dump out the pre-master-key
	// for your HTTP connection. This allows you decrypt a packet capture from wireshark.
	//
	// If you want to do this just set SSLKEYLOGFILE env var to a path on disk and
	// Go will write out the key.
	tv.KeyLogPath = os.Getenv("SSLKEYLOGFILE")
	return tv, nil
}

type clientWrapper struct {
	*azeventgrid.Client
	TestVars testVars
}

type clientWrapperOptions struct {
	DontPurgeEvents bool
}

func newClientWrapper(t *testing.T, opts *clientWrapperOptions) clientWrapper {
	var client *azeventgrid.Client
	var tv testVars

	if recording.GetRecordMode() != recording.PlaybackMode {
		tmpTestVars, err := loadEnv()
		require.NoError(t, err)
		tv = tmpTestVars
	} else {
		tv = fakeTestVars
	}

	if recording.GetRecordMode() == recording.LiveMode {
		if tv.KeyLogPath != "" {
			keyLogWriter, err := os.OpenFile(tv.KeyLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
			require.NoError(t, err)

			t.Cleanup(func() { keyLogWriter.Close() })

			tp := http.DefaultTransport.(*http.Transport).Clone()
			tp.TLSClientConfig = &tls.Config{
				KeyLogWriter: keyLogWriter,
			}

			httpClient := &http.Client{Transport: tp}

			tmpClient, err := azeventgrid.NewClientWithSharedKeyCredential(tv.Endpoint, tv.Key, &azeventgrid.ClientOptions{
				ClientOptions: azcore.ClientOptions{
					Transport: httpClient,
				},
			})
			require.NoError(t, err)
			client = tmpClient
		} else {
			tmpClient, err := azeventgrid.NewClientWithSharedKeyCredential(tv.Endpoint, tv.Key, nil)
			require.NoError(t, err)
			client = tmpClient
		}

		purgePreviousEvents(t, client, tv)
	} else {
		tmpClient, err := azeventgrid.NewClientWithSharedKeyCredential(tv.Endpoint, tv.Key, &azeventgrid.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: newRecordingTransporter(t, tv),
			},
		})
		require.NoError(t, err)
		client = tmpClient
	}

	return clientWrapper{
		Client:   client,
		TestVars: tv,
	}
}

func newRecordingTransporter(t *testing.T, testVars testVars) policy.Transporter {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	err = recording.Start(t, "sdk/messaging/azeventgrid/testdata", nil)
	require.NoError(t, err)

	// err = recording.ResetProxy(nil)
	// require.NoError(t, err)

	err = recording.AddURISanitizer(fakeTestVars.Endpoint, testVars.Endpoint, nil)
	require.NoError(t, err)

	err = recording.AddURISanitizer(fakeTestVars.Topic, testVars.Topic, nil)
	require.NoError(t, err)

	err = recording.AddURISanitizer(fakeTestVars.Subscription, testVars.Subscription, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(`"time": "2023-06-17T00:33:32Z"`, `"time":".+?"`, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(
		`"id":"00000000-0000-0000-0000-000000000000"`,
		`"id":"[^"]+"`, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(
		`"lockToken":"fake-lock-token"`,
		`"lockToken":\s*"[^"]+"`, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(
		`"lockTokens": ["fake-lock-token"]`,
		`"lockTokens":\s*\[\s*"[^"]+"\s*\]`, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(
		`"succeededLockTokens": ["fake-lock-token"]`,
		`"succeededLockTokens":\s*\[\s*"[^"]+"\s*\]`, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(
		`"succeededLockTokens": ["fake-lock-token", "fake-lock-token", "fake-lock-token"]`,
		`"succeededLockTokens":\s*`+
			`\[`+
			`(\s*"[^"]+"\s*\,){2}`+
			`\s*"[^"]+"\s*`+
			`\]`, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(
		`"lockTokens": ["fake-lock-token", "fake-lock-token"]`,
		`"lockTokens":\s*\[\s*"[^"]+"\s*\,\s*"[^"]+"\s*\]`, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(
		`"lockTokens": ["fake-lock-token", "fake-lock-token", "fake-lock-token"]`,
		`"lockTokens":\s*`+
			`\[`+
			`(\s*"[^"]+"\s*\,){2}`+
			`\s*"[^"]+"\s*`+
			`\]`, nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})

	return transport
}

func requireEqualCloudEvent(t *testing.T, expected messaging.CloudEvent, actual messaging.CloudEvent) {
	t.Helper()

	require.NotEmpty(t, actual.ID, "ID is not empty")
	require.NotEmpty(t, actual.SpecVersion, "SpecVersion is not empty")

	expected.ID = actual.ID
	expected.Time = actual.Time

	require.Equal(t, actual, expected)
}

var purge sync.Once

func purgePreviousEvents(t *testing.T, c *azeventgrid.Client, testVars testVars) {
	purge.Do(func() {
		if recording.GetRecordMode() != recording.LiveMode {
			return
		}

		// we'll purge all the events in the queue just to ensure tests
		// run clean.
		events, err := c.ReceiveCloudEvents(context.Background(), testVars.Topic, testVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
			MaxEvents:   to.Ptr[int32](100),
			MaxWaitTime: to.Ptr[int32](10),
		})
		require.NoError(t, err)

		var lockTokens []string

		for _, e := range events.Value {
			lockTokens = append(lockTokens, *e.BrokerProperties.LockToken)
		}

		if len(lockTokens) > 0 {
			resp, err := c.AcknowledgeCloudEvents(context.Background(), testVars.Topic, testVars.Subscription, azeventgrid.AcknowledgeOptions{
				LockTokens: lockTokens,
			}, nil)
			require.NoError(t, err)
			require.Empty(t, resp.FailedLockTokens)
		}
	})

}
