//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azmetrics_test

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics"
	"github.com/stretchr/testify/require"
)

const (
	recordingDirectory = "sdk/monitor/query/azmetrics/testdata"
	fakeResourceURI    = "/subscriptions/faa080af-c1d8-40ad-9cce-e1a451va7b87/resourceGroups/rg-example/providers/Microsoft.AppConfiguration/configurationStores/example"
	fakeSubscriptionID = "faa080af-c1d8-40ad-9cce-e1a451va7b87"
	fakeRegion         = "westus"
)

var (
	credential     azcore.TokenCredential
	resourceURI    string
	subscriptionID string
	region         string
	endpoint       string
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

	resourceURI = getEnvVar("RESOURCE_URI", fakeResourceURI)
	subscriptionID = getEnvVar("AZMETRICS_SUBSCRIPTION_ID", fakeSubscriptionID)
	region = getEnvVar("AZMETRICS_LOCATION", fakeRegion)

	var err error
	credential, err = azcred.New(nil)
	if err != nil {
		panic(err)
	}
	endpoint = "https://" + region + ".metrics.monitor.azure.com"
	if cloudEnv, ok := os.LookupEnv("AZMETRICS_ENVIRONMENT"); ok {
		if strings.EqualFold(cloudEnv, "AzureUSGovernment") {
			clientCloud = cloud.AzureGovernment
			endpoint = "https://" + region + ".metrics.monitor.azure.us"
		}
		if strings.EqualFold(cloudEnv, "AzureChinaCloud") {
			clientCloud = cloud.AzureChina
			endpoint = "https://" + region + ".metrics.monitor.azure.cn"
		}
	}
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

func startTest(t *testing.T) *azmetrics.Client {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azmetrics.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport, Cloud: clientCloud}}
	client, err := azmetrics.NewClient(endpoint, credential, opts)
	require.NoError(t, err)
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
