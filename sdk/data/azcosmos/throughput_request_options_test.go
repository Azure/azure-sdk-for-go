// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestThroughputRequestOptionsToHeaders(t *testing.T) {
	options := &ThroughputOptions{}
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}

	options.IfMatchEtag = "etag"
	options.IfNoneMatchEtag = "noneetag"
	header := options.toHeaders()
	if header == nil {
		t.Error("toHeaders should return non-nil")
	}

	headers := *header
	if headers[headerIfMatch] != options.IfMatchEtag {
		t.Errorf("IfMatchEtag not set matching expected %v got %v", options.IfMatchEtag, headers[headerIfMatch])
	}
	if headers[headerIfNoneMatch] != options.IfNoneMatchEtag {
		t.Errorf("IfNoneMatchEtag not set matching expected %v got %v", options.IfNoneMatchEtag, headers[headerIfNoneMatch])
	}
}
