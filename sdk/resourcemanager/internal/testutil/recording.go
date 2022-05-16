//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

// Host of the test proxy.
func (r *recordingPolicy) Host() string {
	if r.options.UseHTTPS {
		return "localhost:5001"
	}
	return "localhost:5000"
}

// Scheme of the test proxy.
func (r *recordingPolicy) Scheme() string {
	if r.options.UseHTTPS {
		return "https"
	}
	return "http"
}

// NewRecordingPolicy will create a recording policy which can be used in pipeline.
// The policy will change the destination of the request to the proxy server and add required header for the recording test.
func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{UseHTTPS: true}
	}
	p := &recordingPolicy{options: *o, t: t}
	return p
}

// Do with recording mode.
// When handling live request, the policy will do nothing.
// Otherwise, the policy will replace the URL of the request with the test proxy endpoint.
// After request, the policy will change back to the original URL for the request to prevent wrong polling URL for LRO.
func (r *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != "live" && !recording.IsLiveOnly(r.t) {
		oriSchema := req.Raw().URL.Scheme
		oriHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = r.Scheme()
		req.Raw().URL.Host = r.Host()
		req.Raw().Host = r.Host()

		// replace request target to use test proxy
		req.Raw().Header.Set(recording.UpstreamURIHeader, fmt.Sprintf("%v://%v", oriSchema, oriHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IDHeader, recording.GetRecordingId(r.t))

		resp, err = req.Next()
		// for any lro operation, need to change back to the original target to prevent
		if resp != nil {
			resp.Request.URL.Scheme = oriSchema
			resp.Request.URL.Host = oriHost
		}
		return resp, err
	} else {
		return req.Next()
	}
}

// StartRecording starts the recording with the path to store recording file.
// It will return a delegate function to stop recording.
func StartRecording(t *testing.T, pathToPackage string) func() {
	// sanitizer for any uuid string, e.g., subscriptionID
	err := recording.AddGeneralRegexSanitizer("00000000-0000-0000-0000-000000000000", `[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`, nil)
	if err != nil {
		t.Fatalf("Failed to add uuid sanitizer: %v", err)
	}
	// consolidate resource group name for recording and playback
	err = recording.AddGeneralRegexSanitizer("go-sdk-test-rg", `go-sdk-test-\d+`, nil)
	if err != nil {
		t.Fatalf("Failed to add resource group name sanitizer: %v", err)
	}
	err = recording.Start(t, pathToPackage, nil)
	if err != nil {
		t.Fatalf("Failed to start recording: %v", err)
	}
	return func() { StopRecording(t) }
}

// StopRecording stops the recording.
func StopRecording(t *testing.T) {
	err := recording.Stop(t, nil)
	if err != nil {
		t.Fatalf("Failed to stop recording: %v", err)
	}
}
