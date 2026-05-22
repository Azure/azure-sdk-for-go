// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"net/http"
	"strings"
	"testing"

	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

type testBody struct {
	body *strings.Reader
}

func (r *testBody) Close() error {
	return nil
}

func (r *testBody) Read(b []byte) (int, error) {
	return r.body.Read(b)
}

func (r *testBody) Seek(offset int64, whence int) (int64, error) {
	return r.body.Seek(offset, whence)
}

func TestGetCredAndClientOptions(t *testing.T) {
	testEndpoint := "https://test"
	cred, options := GetCredAndClientOptions(t)
	pl, err := armruntime.NewPipeline("testmodule", "v0.1.0", cred, runtime.PipelineOptions{}, options)
	require.NoError(t, err)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, testEndpoint)
	require.NoError(t, err)
	err = req.SetBody(&testBody{body: strings.NewReader("test")}, "text/plain")
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	if recording.GetRecordMode() == recording.PlaybackMode {
		require.Equal(t, "Bearer "+recording.SanitizedValue, resp.Request.Header.Get("Authorization"))
	}
	require.Equal(t, testEndpoint, resp.Request.URL.String())

}
