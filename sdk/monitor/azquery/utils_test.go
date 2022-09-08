//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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
		credential, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			panic(err)
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
	opts := &azquery.LogsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	return azquery.NewLogsClient(credential, opts)
}

func startMetricsTest(t *testing.T) *azquery.MetricsClient {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azquery.MetricsClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	return azquery.NewMetricsClient(credential, opts)
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
