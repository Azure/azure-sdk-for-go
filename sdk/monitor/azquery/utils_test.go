//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/monitor/azquery/testdata"
const fakeWorkspaceID = "32d1e136-gg81-4b0a-b647-260cdc471f68"
const fakeWorkspaceID2 = "asdjkfj8k20-gg81-4b0a-9fu2-260c09fn1f68"
const fakeResourceURI = "/subscriptions/faa080af-c1d8-40ad-9cce-e1a451va7b87/resourceGroups/rg-example/providers/Microsoft.AppConfiguration/configurationStores/example"
const fakeSubscrtiptionID = "faa080af-c1d8-40ad-9cce-e1a451va7b87"
const fakeRegion = "westus2"

var (
	credential     azcore.TokenCredential
	workspaceID    string
	workspaceID2   string
	resourceURI    string
	subscriptionID string
	region         string
	clientCloud    cloud.Configuration
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
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

	var err error
	credential, err = azcred.New(nil)
	if err != nil {
		panic(err)
	}
	if cloudEnv, ok := os.LookupEnv("AZQUERY_ENVIRONMENT"); ok {
		if strings.EqualFold(cloudEnv, "AzureUSGovernment") {
			clientCloud = cloud.AzureGovernment
		}
		if strings.EqualFold(cloudEnv, "AzureChinaCloud") {
			clientCloud = cloud.AzureChina
		}
	}
	workspaceID = getEnvVar("WORKSPACE_ID", fakeWorkspaceID)
	workspaceID2 = getEnvVar("WORKSPACE_ID2", fakeWorkspaceID2)
	resourceURI = getEnvVar("RESOURCE_URI", fakeResourceURI)
	subscriptionID = getEnvVar("AZQUERY_SUBSCRIPTION_ID", fakeSubscrtiptionID)
	region = getEnvVar("AZQUERY_LOCATION", fakeRegion)

	return m.Run()
}

func startRecording(t *testing.T) {
	err := recording.Start(t, recordingDirectory, nil)
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
	client, err := azquery.NewLogsClient(credential, opts)
	if err != nil {
		panic(err)
	}
	return client
}

func startMetricsTest(t *testing.T) *azquery.MetricsClient {
	var opts *azquery.MetricsClientOptions
	if recording.GetRecordMode() == recording.LiveMode {
		transport := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Renegotiation: tls.RenegotiateOnceAsClient,
				},
			},
		}
		opts = &azquery.MetricsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport, Cloud: clientCloud}}
	} else {
		startRecording(t)
		transport, err := recording.NewRecordingHTTPClient(t, nil)
		require.NoError(t, err)
		opts = &azquery.MetricsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	}

	client, err := azquery.NewMetricsClient(credential, opts)
	if err != nil {
		panic(err)
	}
	return client
}

func startMetricsBatchTest(t *testing.T) *azquery.MetricsBatchClient {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azquery.MetricsBatchClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	endpoint := "https://" + region + ".metrics.monitor.azure.com"
	client, err := azquery.NewMetricsBatchClient(endpoint, credential, opts)
	if err != nil {
		panic(err)
	}
	return client
}

func getEnvVar(envVar string, fakeValue string) string {
	// get value
	value := fakeValue
	if recording.GetRecordMode() == recording.LiveMode || recording.GetRecordMode() == recording.RecordingMode {
		value = os.Getenv(envVar)
		if value == "" {
			panic("no value for " + envVar)
		}
	}

	// sanitize value
	if fakeValue != "" && recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeValue, value, nil)
		if err != nil {
			panic(err)
		}
	}

	return value
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

	// testing unmarshal error scenarios
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), ":false}")
	if !strings.Contains(modifiedData, "render") && modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}
