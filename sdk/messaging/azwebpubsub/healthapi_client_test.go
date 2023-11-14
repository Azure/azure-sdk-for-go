//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azwebpubsub_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub/internal"
	"github.com/stretchr/testify/require"
)

func TestHealthAPIClient_GetServiceStatus(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode || testing.Short() {
		t.Skip()
	}

	client := getHealthAPIClient(t)
	_, err := client.GetServiceStatus(context.Background(), &azwebpubsub.HealthAPIClientGetServiceStatusOptions{})
	require.NoError(t, err)
}

func loadEndpointFromEnv() (string, error) {
	if v := os.Getenv("WEBPUBSUB_ENDPOINT"); v != "" {
		return v, nil
	}
	if v := os.Getenv("WEBPUBSUB_CONNECTIONSTRING"); v != "" {
		props, err := internal.ParseConnectionString(v)
		if err != nil {
			return "", err
		}
		return props.Endpoint, nil
	}
	return "", nil
}

func getHealthAPIClient(t *testing.T) *azwebpubsub.HealthAPIClient {
	var options *azwebpubsub.HealthAPIClientOptions
	var endpoint string
	if recording.GetRecordMode() != recording.PlaybackMode {
		tmpEndpoint, err := loadEndpointFromEnv()
		require.NoError(t, err)
		require.NotEmpty(t, tmpEndpoint)
		endpoint = tmpEndpoint
	} else {
		endpoint = "https://fake.eastus-1.webpubsub.azure.com"
	}

	if recording.GetRecordMode() == recording.LiveMode {
		keyLogPath := os.Getenv("SSLKEYLOGFILE")
		if keyLogPath != "" {
			keyLogWriter, err := os.OpenFile(keyLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
			require.NoError(t, err)

			t.Cleanup(func() { keyLogWriter.Close() })

			tp := http.DefaultTransport.(*http.Transport).Clone()
			tp.TLSClientConfig = &tls.Config{
				KeyLogWriter: keyLogWriter,
			}

			httpClient := &http.Client{Transport: tp}
			options = &azwebpubsub.HealthAPIClientOptions{
				ClientOptions: azcore.ClientOptions{
					Transport: httpClient,
				},
			}
		} else {
			options = nil
		}
	} else {
		options = &azwebpubsub.HealthAPIClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: newRecordingTransporter(t, testVars{Endpoint: endpoint}),
			},
		}
	}

	client, err := azwebpubsub.NewHealthAPIClient("https://lianwei-test-1.webpubsub.azure.com", options)
	require.NoError(t, err)
	return client
}
