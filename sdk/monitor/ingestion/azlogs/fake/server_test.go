// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs/fake"
	"github.com/stretchr/testify/require"
)

func getServer() fake.Server {
	return fake.Server{
		Upload: func(ctx context.Context, ruleID string, streamName string, logs []byte, options *azlogs.UploadOptions) (resp azfake.Responder[azlogs.UploadResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusNoContent, azlogs.UploadResponse{}, nil)
			return
		},
	}
}
func TestServer(t *testing.T) {
	fakeServer := getServer()
	client, err := azlogs.NewClient("https://fake.region.ingest.monitor.azure.com", &azfake.TokenCredential{}, &azlogs.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	_, err = client.Upload(context.Background(), "fake-rule-id", "fake-stream-name", []byte("fake-logs"), nil)
	require.NoError(t, err)
}
