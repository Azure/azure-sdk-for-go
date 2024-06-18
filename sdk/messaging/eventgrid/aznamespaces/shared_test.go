//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces_test

import (
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
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces"
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

var initTopic sync.Once
var tv testVars = fakeTestVars

func newClientOptions(t *testing.T) azcore.ClientOptions {
	if recording.GetRecordMode() != recording.PlaybackMode {
		initTopic.Do(func() {
			tmpTestVars, err := loadEnv()
			require.NoError(t, err)
			tv = tmpTestVars
		})
	}

	options := azcore.ClientOptions{}

	if recording.GetRecordMode() == recording.LiveMode {
		if tv.KeyLogPath != "" {
			keyLogWriter, err := os.OpenFile(tv.KeyLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
			require.NoError(t, err)

			t.Cleanup(func() { keyLogWriter.Close() })

			tp := http.DefaultTransport.(*http.Transport).Clone()
			tp.TLSClientConfig = &tls.Config{
				KeyLogWriter: keyLogWriter,
			}

			options = azcore.ClientOptions{
				Transport: &http.Client{Transport: tp},
			}
		}
	} else {
		options = azcore.ClientOptions{
			Transport: newRecordingTransporter(t, tv),
		}
	}

	options.Logging = policy.LogOptions{
		IncludeBody: true,
		AllowedQueryParams: []string{
			"maxEvents",
			"maxWaitTime",
		},
		AllowedHeaders: []string{
			// these are the standard headers for binary content mode with CloudEvents
			"ce-id",
			"ce-specversion",
			"ce-time",
			"ce-source",
			"ce-type",

			// these are the Extension fields that I use in testing binary content mode.
			"ce-extensiondatastring",
			"ce-extensiondatastring2",
			"ce-extensiondataint",
			"ce-extensiondataurl",
			"ce-extensiondatauint",
			"ce-extensiondatatime",
			"ce-extensiondatabytes",

			"Authorization",
		},
	}

	return options
}

func newClients(t *testing.T, useSASKey bool) (*aznamespaces.SenderClient, *aznamespaces.ReceiverClient) {
	if os.Getenv("FORCE_SASKEY") == "true" {
		t.Logf("Switching from TokenCredential -> SAS Key because FORCE_SASKEY is true. See https://github.com/Azure/azure-sdk-for-go/issues/22961 for more details")
		useSASKey = true
	}

	return newSenderClient(t, useSASKey), newReceiverClient(t, useSASKey)
}

func newSenderClient(t *testing.T, useSASKey bool) *aznamespaces.SenderClient {
	options := newClientOptions(t)

	if useSASKey {
		client, err := aznamespaces.NewSenderClientWithSharedKeyCredential(tv.Endpoint, tv.Topic, azcore.NewKeyCredential(tv.Key), &aznamespaces.SenderClientOptions{ClientOptions: options})
		require.NoError(t, err)
		return client
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	client, err := aznamespaces.NewSenderClient(tv.Endpoint, tv.Topic, cred, &aznamespaces.SenderClientOptions{ClientOptions: options})
	require.NoError(t, err)

	return client
}

func newReceiverClient(t *testing.T, useSASKey bool) *aznamespaces.ReceiverClient {
	options := newClientOptions(t)

	if useSASKey {
		client, err := aznamespaces.NewReceiverClientWithSharedKeyCredential(tv.Endpoint, tv.Topic, tv.Subscription, azcore.NewKeyCredential(tv.Key), &aznamespaces.ReceiverClientOptions{ClientOptions: options})
		require.NoError(t, err)
		return client
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	client, err := aznamespaces.NewReceiverClient(tv.Endpoint, tv.Topic, tv.Subscription, cred, &aznamespaces.ReceiverClientOptions{ClientOptions: options})
	require.NoError(t, err)

	return client
}

func newRecordingTransporter(t *testing.T, testVars testVars) policy.Transporter {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	err = recording.Start(t, recordingDirectory, nil)
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

	require.Equal(t, expected, actual)
}
