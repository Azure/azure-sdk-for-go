//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"
)

const fakeWorkspaceID = "32d1e136-gg81-4b0a-b647-260cdc471f68"
const fakeWorkspaceID2 = "asdjkfj8k20-gg81-4b0a-9fu2-260c09fn1f68"
const fakeResourceURI = "/subscriptions/faa080af-c1d8-40ad-9cce-e1a451va7b87/resourceGroups/rg-example/providers/Microsoft.AppConfiguration/configurationStores/example"

var (
	credential   azcore.TokenCredential
	workspaceID  string
	workspaceID2 string
	resourceURI  string
	clientCloud  cloud.Configuration = cloud.AzurePublic
)

func TestMain(m *testing.M) {
	if recording.GetRecordMode() == recording.LiveMode || recording.GetRecordMode() == recording.RecordingMode {
		workspaceID, workspaceID2, resourceURI = os.Getenv("WORKSPACE_ID"), os.Getenv("WORKSPACE_ID2"), os.Getenv("RESOURCE_URI")
	}
	if workspaceID == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for WORKSPACE_ID")
		}
		workspaceID = fakeWorkspaceID
	}
	if workspaceID2 == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for WORKSPACE_ID2")
		}
		workspaceID2 = fakeWorkspaceID2
	}
	if resourceURI == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for RESOURCE_URI")
		}
		resourceURI = fakeResourceURI
	}
	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("AZQUERY_TENANT_ID")
		clientID := lookupEnvVar("AZQUERY_CLIENT_ID")
		secret := lookupEnvVar("AZQUERY_CLIENT_SECRET")
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		//credOptions := azidentity.DefaultAzureCredentialOptions{ClientOptions: azcore.ClientOptions{Cloud: cloud.AzureGovernment}}
		//credential, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			panic(err)
		}

		if cloudEnv, ok := os.LookupEnv("AZQUERY_ENVIRONMENT"); ok {
			if cloudEnv == "AzureUSGovernment" {
				clientCloud = cloud.AzureGovernment
			}
			if cloudEnv == "AzureChinaCloud" {
				clientCloud = cloud.AzureChina
			}
		}

	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeWorkspaceID, workspaceID, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddGeneralRegexSanitizer(fakeWorkspaceID2, workspaceID2, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddGeneralRegexSanitizer(fakeResourceURI, resourceURI, nil)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	code := m.Run()
	os.Exit(code)
}

func startRecording(t *testing.T) {
	err := recording.Start(t, "sdk/monitor/azquery/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startLogsTest(t *testing.T) *azquery.LogsClient {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azquery.LogsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport, Cloud: clientCloud}}
	return azquery.NewLogsClient(credential, opts)
}

func startMetricsTest(t *testing.T) *azquery.MetricsClient {
	/*startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azquery.MetricsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	return azquery.NewMetricsClient(credential, opts)*/
	//transport = azcore.defaultTransport{TLSClientConfig: }
	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:    tls.VersionTLS12,
			Renegotiation: tls.RenegotiateFreelyAsClient,
		},
	}
	transport := &http.Client{
		Transport: defaultTransport,
	}

	opts := &azquery.MetricsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}

	return azquery.NewMetricsClient(credential, opts)
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

type FakeCredential struct{}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)
}
