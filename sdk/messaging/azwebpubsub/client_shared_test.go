// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azwebpubsub_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub/internal"
	"github.com/stretchr/testify/require"
)

type clientWrapper struct {
	*azwebpubsub.Client
	TestVars testVars
}

var fakeTestVars = testVars{
	ConnectionString: "Endpoint=https://fake.eastus-1.webpubsub.azure.com;AccessKey=ABCDE;",
	Hub:              "chat",
	Endpoint:         "https://fake.eastus-1.webpubsub.azure.com",
}

type testVars struct {
	// NewClientFromConnectionString when ConnectionString is set
	ConnectionString string
	Hub              string
	Endpoint         string
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
		ConnectionString: get("WEBPUBSUB_CONNECTIONSTRING"),
		Hub:              get("WEBPUBSUB_HUB"),
		Endpoint:         get("WEBPUBSUB_ENDPOINT"),
	}

	if len(missing) > 1 {
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

func newClientWrapper(t *testing.T) clientWrapper {
	var client *azwebpubsub.Client
	var tv testVars
	var options *azwebpubsub.ClientOptions
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
			options = &azwebpubsub.ClientOptions{
				ClientOptions: azcore.ClientOptions{
					Transport: httpClient,
				},
			}
		} else {
			options = nil
		}
	} else {
		options = &azwebpubsub.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: newRecordingTransporter(t, tv),
			},
		}
	}

	if tv.ConnectionString != "" {
		tmpClient, err := azwebpubsub.NewClientFromConnectionString(tv.ConnectionString, tv.Hub, options)
		require.NoError(t, err)
		client = tmpClient
	} else {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		require.NoError(t, err)

		tmpClient, err := azwebpubsub.NewClient(tv.Endpoint, tv.Hub, cred, options)
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

	err = recording.Start(t, "sdk/messaging/azwebpubsub/testdata", nil)
	require.NoError(t, err)

	// err = recording.ResetProxy(nil)
	// require.NoError(t, err)
	props, _ := internal.ParseConnectionString(testVars.ConnectionString)
	err = recording.AddURISanitizer("https://fake_endpoint.com", props.Endpoint, nil)
	require.NoError(t, err)

	err = recording.AddGeneralRegexSanitizer(`"time": "2023-11-31T00:33:32Z"`, `"time":".+?"`, nil)
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
