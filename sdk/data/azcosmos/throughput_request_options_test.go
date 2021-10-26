// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestThroughputRequestOptionsToHeaders(t *testing.T) {
	options := &ThroughputOptions{}
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}

	etag := azcore.ETag("etag")
	noneetag := azcore.ETag("noneetag")
	options.IfMatchEtag = &etag
	options.IfNoneMatchEtag = &noneetag

	header := options.toHeaders()
	if header == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	headers := *header
	if headers[headerIfMatch] != string(*options.IfMatchEtag) {
		t.Errorf("IfMatchEtag not set matching expected %v got %v", options.IfMatchEtag, headers[headerIfMatch])
	}
	if headers[headerIfNoneMatch] != string(*options.IfNoneMatchEtag) {
		t.Errorf("IfNoneMatchEtag not set matching expected %v got %v", options.IfNoneMatchEtag, headers[headerIfNoneMatch])
	}
}
