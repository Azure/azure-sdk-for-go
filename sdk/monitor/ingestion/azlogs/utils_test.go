// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azlogs_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/monitor/ingestion/azlogs/testdata"
const fakeEndpoint = "https://test.eastus-1.ingest.monitor.azure.com"
const fakeRuleID = "Custom-TestTable_CL"
const fakeStreamName = "dcr-testing"

var (
	credential  azcore.TokenCredential
	endpoint    string
	ruleID      string
	streamName  string
	clientCloud cloud.Configuration
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
	if cloudEnv, ok := os.LookupEnv("AZLOGS_ENVIRONMENT"); ok {
		if strings.EqualFold(cloudEnv, "AzureUSGovernment") {
			clientCloud = cloud.AzureGovernment
		}
		if strings.EqualFold(cloudEnv, "AzureChinaCloud") {
			clientCloud = cloud.AzureChina
		}
	}
	endpoint = getEnvVar("AZURE_MONITOR_DCE", fakeEndpoint)
	ruleID = getEnvVar("AZURE_MONITOR_DCR_ID", fakeRuleID)
	streamName = getEnvVar("AZURE_MONITOR_STREAM_NAME", fakeStreamName)

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

func startTest(t *testing.T) *azlogs.Client {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azlogs.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport, Cloud: clientCloud}}

	client, err := azlogs.NewClient(endpoint, credential, opts)
	if err != nil {
		panic(err)
	}
	return client
}

func getEnvVar(lookupValue string, fakeValue string) string {
	// get value
	envVar := fakeValue
	if recording.GetRecordMode() == recording.LiveMode || recording.GetRecordMode() == recording.RecordingMode {
		envVar = os.Getenv(lookupValue)
		if envVar == "" {
			panic("no value for " + lookupValue)
		}
	}

	// sanitize value
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeValue, envVar, nil)
		if err != nil {
			panic(err)
		}
	}

	return envVar
}
