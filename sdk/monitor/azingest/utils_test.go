//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azingest_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azingest"
	"github.com/stretchr/testify/require"
)

const fakeEndpoint = "https://test.eastus-1.ingest.monitor.azure.com"
const fakeRuleID = "Custom-TestTable_CL"
const fakeStreamName = "dcr-testing"

var (
	credential azcore.TokenCredential
	endpoint   string
	ruleID     string
	streamName string
)

func TestMain(m *testing.M) {
	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("AZINGEST_TENANT_ID")
		clientID := lookupEnvVar("AZINGEST_CLIENT_ID")
		secret := lookupEnvVar("AZINGEST_CLIENT_SECRET")
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}
	endpoint = getEnvVar("AZURE_MONITOR_DCE", fakeEndpoint)
	ruleID = getEnvVar("AZURE_MONITOR_DCR_ID", fakeRuleID)
	streamName = getEnvVar("AZURE_MONITOR_STREAM_NAME", fakeStreamName)

	code := m.Run()
	os.Exit(code)
}

func startRecording(t *testing.T) {
	err := recording.Start(t, "sdk/monitor/azingest/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startTest(t *testing.T) *azingest.Client {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azingest.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}

	client, err := azingest.NewClient(endpoint, credential, opts)
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
