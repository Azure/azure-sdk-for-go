// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/stretchr/testify/require"
)

type sequenceTransport struct {
	responses []*http.Response
}

func (s *sequenceTransport) Do(r *http.Request) (*http.Response, error) {
	resp := s.responses[0]
	s.responses = s.responses[1:]
	resp.Request = r
	return resp, nil
}

func newTestResponse(status int, body string, headers ...string) *http.Response {
	h := http.Header{}
	for i := 0; i+1 < len(headers); i += 2 {
		h.Set(headers[i], headers[i+1])
	}
	return &http.Response{StatusCode: status, Header: h, Body: io.NopCloser(strings.NewReader(body))}
}

// TestBeginCreateSnapshotResultFilters verifies the poller's final result keeps each
// filter's key and label.
func TestBeginCreateSnapshotResultFilters(t *testing.T) {
	transport := &sequenceTransport{responses: []*http.Response{
		newTestResponse(201, `{"name":"s1","status":"provisioning"}`,
			"Operation-Location", "https://store.azconfig.example.test/operations?snapshot=s1&api-version=2026-04-01",
			"Sync-Token", "a=b;sn=1"),
		newTestResponse(200, `{"id":"op","status":"Succeeded"}`),
		newTestResponse(200, `{"name":"s1","status":"ready","filters":[{"key":"app/*","label":"prod"}]}`),
	}}

	client, err := azappconfig.NewClientFromConnectionString(
		"Endpoint=https://store.azconfig.example.test;Id=id;Secret=c2VjcmV0",
		&azappconfig.ClientOptions{ClientOptions: azcore.ClientOptions{
			Transport: transport,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		}})
	require.NoError(t, err)

	ctx := context.Background()
	poller, err := client.BeginCreateSnapshot(ctx, "s1",
		[]azappconfig.SettingFilter{{KeyFilter: to.Ptr("app/*"), LabelFilter: to.Ptr("prod")}}, nil)
	require.NoError(t, err)
	_, err = poller.Poll(ctx)
	require.NoError(t, err)

	result, err := poller.Result(ctx)
	require.NoError(t, err)
	require.Len(t, result.Filters, 1)
	require.NotNil(t, result.Filters[0].KeyFilter)
	require.Equal(t, "app/*", *result.Filters[0].KeyFilter)
	require.NotNil(t, result.Filters[0].LabelFilter)
	require.Equal(t, "prod", *result.Filters[0].LabelFilter)
}
