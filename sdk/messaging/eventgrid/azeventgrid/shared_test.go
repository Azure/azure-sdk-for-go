//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azeventgrid"
	"github.com/stretchr/testify/require"
)

type topicVars struct {
	Name     string
	Key      string
	Endpoint string
	SAS      string
}

type tokenCredentialVars struct {
	ClientID       string
	ClientSecret   string
	TenantID       string
	SubscriptionID string
}

type eventGridVars struct {
	// EG are connection variables for an EventGrid encoded topic.
	EG topicVars
	// CE are connection variables for a CloudEvent encoded topic.
	CE topicVars

	TokenCredVars tokenCredentialVars
}

var fakeVars = eventGridVars{
	EG: topicVars{
		Name:     "faketopic",
		Key:      base64.StdEncoding.EncodeToString([]byte("fakekey")),
		Endpoint: "https://localhost/fake-endpoint",
		SAS:      "fake-sas",
	},
	CE: topicVars{
		Name:     "faketopic",
		Key:      base64.StdEncoding.EncodeToString([]byte("fakekey")),
		Endpoint: "https://localhost/fake-endpoint",
		SAS:      "fake-sas",
	},
	TokenCredVars: tokenCredentialVars{
		TenantID:       "fake-tenant-id",
		ClientID:       "fake-client-id",
		ClientSecret:   "fake-client-secret",
		SubscriptionID: "fake-subscription-id",
	},
}

func newTestVars(t *testing.T) eventGridVars {
	if recording.GetRecordMode() == recording.PlaybackMode {
		err := recording.Start(t, recordingDirectory, nil)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := recording.Stop(t, nil)
			require.NoError(t, err)
		})

		// set these up so DefaultAzureCredential will pick these up when it auths.
		os.Setenv("AZURE_TENANT_ID", fakeVars.TokenCredVars.TenantID)
		os.Setenv("AZURE_CLIENT_ID", fakeVars.TokenCredVars.ClientID)
		os.Setenv("AZURE_CLIENT_SECRET", fakeVars.TokenCredVars.ClientSecret)

		sanitizeForPlayback(t)
		return fakeVars
	}

	egVars := eventGridVars{
		EG: topicVars{Name: os.Getenv("EVENTGRID_TOPIC_NAME"),
			Key:      os.Getenv("EVENTGRID_TOPIC_KEY"),
			Endpoint: os.Getenv("EVENTGRID_TOPIC_ENDPOINT"),
		},
		CE: topicVars{Name: os.Getenv("EVENTGRID_CE_TOPIC_NAME"),
			Key:      os.Getenv("EVENTGRID_CE_TOPIC_KEY"),
			Endpoint: os.Getenv("EVENTGRID_CE_TOPIC_ENDPOINT"),
		},
		TokenCredVars: tokenCredentialVars{
			ClientID:       os.Getenv("AZURE_CLIENT_ID"),
			ClientSecret:   os.Getenv("AZURE_CLIENT_SECRET"),
			TenantID:       os.Getenv("AZURE_TENANT_ID"),
			SubscriptionID: os.Getenv("AZURE_SUBSCRIPTION_ID"),
		},
	}

	for _, v := range []topicVars{egVars.EG, egVars.CE} {
		if v.Endpoint == "" || v.Key == "" || v.Name == "" {
			t.Logf("WARNING: not enabling azeventgrid integration tests, environment variables not set")
			t.Skip()
			break
		}
	}

	egVars.EG.SAS = generateSAS(egVars.EG.Endpoint, egVars.EG.Key, time.Now())
	egVars.CE.SAS = generateSAS(egVars.CE.Endpoint, egVars.CE.Key, time.Now())

	if egVars.TokenCredVars.ClientID == "" ||
		egVars.TokenCredVars.ClientSecret == "" ||
		egVars.TokenCredVars.TenantID == "" ||
		egVars.TokenCredVars.SubscriptionID == "" {
		t.Logf("WARNING: not enabling azeventgrid integration tests, environment variables for token credential auth not set")
		t.Skip()
	}

	if recording.GetRecordMode() == recording.LiveMode {
		return egVars
	}

	// we're recording then, let's setup the sanitizers.
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})

	sanitizeForRecording(t, egVars)
	return egVars
}

func newClientOptionsForTest(t *testing.T) struct {
	EG  *azeventgrid.ClientOptions
	DAC *azidentity.DefaultAzureCredentialOptions
} {
	var ret = struct {
		EG  *azeventgrid.ClientOptions
		DAC *azidentity.DefaultAzureCredentialOptions
	}{}

	if recording.GetRecordMode() != recording.LiveMode {
		recordingClient, err := recording.NewRecordingHTTPClient(t, nil)
		require.NoError(t, err)

		clientOptions := azcore.ClientOptions{
			Transport: recordingClient,
		}

		ret.DAC = &azidentity.DefaultAzureCredentialOptions{
			ClientOptions: clientOptions,
		}

		ret.EG = &azeventgrid.ClientOptions{
			ClientOptions: clientOptions,
		}
	}

	return ret
}

type dumpFullPolicy struct {
	Prefix string
}

func (p dumpFullPolicy) Do(req *policy.Request) (*http.Response, error) {
	fmt.Printf("\n\n===> BEGIN: REQUEST (%s) <===\n\n", p.Prefix)

	requestBytes, err := httputil.DumpRequestOut(req.Raw(), false)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(requestBytes))
	fmt.Printf("Body: %s\n", string(FormatRequestBytes(req.Raw())))
	fmt.Printf("\n\n===> END: REQUEST (%s)<===\n\n", p.Prefix)

	resp, err := req.Next()

	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n===> BEGIN: RESPONSE (%s) <===\n\n", p.Prefix)

	responseBytes, err := httputil.DumpResponse(resp, false)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(responseBytes))
	fmt.Printf("Body: %s\n", string(FormatResponseBytes(resp)))

	fmt.Printf("\n\n===> END: RESPONSE (%s) <===\n\n", p.Prefix)
	return resp, err
}

func FormatRequestBytes(req *http.Request) []byte {
	if req.Body == nil {
		return nil
	}

	requestBytes, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(requestBytes))
	return FormatBytes(requestBytes)
}

func FormatResponseBytes(resp *http.Response) []byte {
	requestBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(requestBytes))
	return FormatBytes(requestBytes)
}

func FormatBytes(body []byte) []byte {
	var m *map[string]any
	var l *[]any

	candidates := []any{&m, &l}

	for _, v := range candidates {
		err := json.Unmarshal(body, v)

		if err != nil {
			continue
		}

		if err == nil {
			formattedBytes, err := json.MarshalIndent(v, "  ", "  ")

			if err != nil {
				continue
			}

			return formattedBytes
		}
	}

	// if we can't format it we'll just give it back.
	return body
}
