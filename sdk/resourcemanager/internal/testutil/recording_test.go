//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestRecordingPolicy(t *testing.T) {
	testEndpoint := "http://test"
	pl := runtime.NewPipeline("testmodule", "v0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{PerCallPolicies: []policy.Policy{NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: false})}})
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, testEndpoint)
	require.NoError(t, err)
	err = req.SetBody(&testBody{body: strings.NewReader("test")}, "text/plain")
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, testEndpoint, resp.Request.URL.String())
}

func TestStartStopRecording(t *testing.T) {
	stop := StartRecording(t, pathToPackage)
	defer stop()
}

func TestGenerateAlphaNumericID(t *testing.T) {
	stop := StartRecording(t, pathToPackage)
	rnd := GenerateAlphaNumericID(t, "test", 6)
	require.Equal(t, "testNlDAa8", rnd)
	defer stop()
}
