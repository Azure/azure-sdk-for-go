//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testdata

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

// Snippet: PolicyDefinition
// This should be a 'testdata' directory in your module. `testdata` is ignored by the go tool, making it perfect for ancillary data
var pathToPackage = "sdk/data/aztables/testdata"

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{UseHTTPS: true}
	}
	p := &recordingPolicy{options: *o, t: t}
	return p
}

func (p *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != recording.LiveMode {
		originalURLHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = "https"
		req.Raw().URL.Host = p.options.Host
		req.Raw().Host = p.options.Host

		req.Raw().Header.Set(recording.UpstreamURIHeader, fmt.Sprintf("%s://%s", p.options.Scheme, originalURLHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IdHeader, recording.GetRecordingId(p.t))
	}
	return req.Next()
}
// EndSnippet
