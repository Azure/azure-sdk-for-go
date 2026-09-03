// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"context"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

func newSubRequest(t *testing.T, headers map[string]string) *policy.Request {
	t.Helper()
	req, err := runtime.NewRequest(context.Background(), "DELETE", "https://account.blob.core.windows.net/container/blob")
	require.NoError(t, err)
	for k, v := range headers {
		req.Raw().Header.Set(k, v)
	}
	return req
}

func TestBuildSubRequest_Clean(t *testing.T) {
	req := newSubRequest(t, map[string]string{
		"x-ms-if-tags": "\"tenant\"='A'",
	})
	result, err := buildSubRequest(req)
	require.NoError(t, err)
	require.Contains(t, string(result), "X-Ms-If-Tags: \"tenant\"='A'")
}

func TestBuildSubRequest_CRLFInHeaderValue(t *testing.T) {
	req := newSubRequest(t, map[string]string{
		"x-ms-if-tags": "\"tenant\"='A'\r\nx-ms-delete-snapshots: include",
	})
	_, err := buildSubRequest(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid CR/LF")
}

func TestBuildSubRequest_LFInHeaderValue(t *testing.T) {
	req := newSubRequest(t, map[string]string{
		"x-ms-if-tags": "\"tenant\"='A'\nx-ms-delete-snapshots: include",
	})
	_, err := buildSubRequest(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid CR/LF")
}

func TestBuildSubRequest_CRInHeaderValue(t *testing.T) {
	req := newSubRequest(t, map[string]string{
		"x-ms-if-tags": "value\rwith-cr",
	})
	_, err := buildSubRequest(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid CR/LF")
}

func TestBuildSubRequest_CRLFInHeaderName(t *testing.T) {
	req := newSubRequest(t, nil)
	req.Raw().Header["bad\r\nheader"] = []string{"value"}
	_, err := buildSubRequest(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "header name")
}

func TestBuildSubRequest_MultipleInjectedHeaders(t *testing.T) {
	req := newSubRequest(t, map[string]string{
		"x-ms-if-tags": "\"tenant\"='A'\r\nx-ms-delete-snapshots: include\r\nx-ms-extra: injected",
	})
	_, err := buildSubRequest(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid CR/LF")
}

func TestBuildSubRequest_SkipsXmsVersion(t *testing.T) {
	req := newSubRequest(t, map[string]string{
		"x-ms-version": "2023-11-03",
		"x-ms-if-tags": "\"tenant\"='A'",
	})
	result, err := buildSubRequest(req)
	require.NoError(t, err)
	require.False(t, strings.Contains(string(result), "x-ms-version"))
	require.False(t, strings.Contains(string(result), "X-Ms-Version"))
}
