// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azlogs"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azlogs/fake"
	"github.com/stretchr/testify/require"
)

var fakeWorkspaceID = "fake-workspace-id"
var fakeResourceID = "testing"

func getServer() fake.Server {
	return fake.Server{
		QueryWorkspace: func(ctx context.Context, workspaceID string, body azlogs.QueryBody, opts *azlogs.QueryWorkspaceOptions) (resp azfake.Responder[azlogs.QueryWorkspaceResponse], errResp azfake.ErrorResponder) {
			if !*opts.Options.Statistics {
				resp.SetResponse(http.StatusNotFound, azlogs.QueryWorkspaceResponse{}, nil)
			}
			logsResp := azlogs.QueryWorkspaceResponse{
				QueryResults: azlogs.QueryResults{
					Tables: []azlogs.Table{
						{
							Name: to.Ptr(workspaceID),
						},
					},
				},
			}
			resp.SetResponse(http.StatusOK, logsResp, nil)
			return
		},
		QueryResource: func(ctx context.Context, resourceID string, body azlogs.QueryBody, opts *azlogs.QueryResourceOptions) (resp azfake.Responder[azlogs.QueryResourceResponse], errResp azfake.ErrorResponder) {
			logsResp := azlogs.QueryResourceResponse{
				QueryResults: azlogs.QueryResults{
					Tables: []azlogs.Table{
						{
							Name: to.Ptr(resourceID),
						},
					},
				},
			}
			resp.SetResponse(http.StatusOK, logsResp, nil)
			return
		},
		QueryBatch: func(ctx context.Context, body azlogs.BatchRequest, options *azlogs.QueryBatchOptions) (resp azfake.Responder[azlogs.QueryBatchResponse], errResp azfake.ErrorResponder) {
			logsResp := azlogs.QueryBatchResponse{
				BatchResponse: azlogs.BatchResponse{
					Responses: []azlogs.BatchQueryResponse{
						{
							ID: body.Requests[0].ID,
						},
						{
							ID: body.Requests[1].ID,
						},
					},
				},
			}
			resp.SetResponse(http.StatusOK, logsResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()
	client, err := azlogs.NewClient(&azfake.TokenCredential{}, &azlogs.ClientOptions{ClientOptions: azcore.ClientOptions{
		Transport: fake.NewServerTransport(&fakeServer),
	},
	})
	require.NoError(t, err)

	resp, err := client.QueryWorkspace(context.Background(), fakeWorkspaceID, azlogs.QueryBody{Query: to.Ptr("Fake query")},
		&azlogs.QueryWorkspaceOptions{
			Options: &azlogs.QueryOptions{Statistics: to.Ptr(true)},
		})
	require.NoError(t, err)
	require.Equal(t, *resp.Tables[0].Name, fakeWorkspaceID)

	resourceResp, err := client.QueryResource(context.Background(), fakeResourceID, azlogs.QueryBody{}, &azlogs.QueryResourceOptions{})
	require.NoError(t, err)
	require.Equal(t, *resourceResp.Tables[0].Name, fakeResourceID)

	batchResp, err := client.QueryBatch(context.Background(),
		azlogs.BatchRequest{
			Requests: []azlogs.BatchQueryRequest{
				{
					ID: to.Ptr("1"),
				},
				{
					ID: to.Ptr("2"),
				},
			},
		}, &azlogs.QueryBatchOptions{})
	require.NoError(t, err)
	require.Equal(t, *batchResp.Responses[0].ID, "1")
	require.Equal(t, *batchResp.Responses[1].ID, "2")
}
