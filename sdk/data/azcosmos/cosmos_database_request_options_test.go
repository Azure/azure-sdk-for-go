// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestDatabaseRequestOptionsToHeaders(t *testing.T) {
	options := &CosmosDatabaseRequestOptions{}
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
	if headers[azcore.HeaderIfMatch] != options.IfMatchEtag {
		t.Errorf("IfMatchEtag not set matching expected %v got %v", options.IfMatchEtag, headers[azcore.HeaderIfMatch])
	}
	if headers[azcore.HeaderIfNoneMatch] != options.IfNoneMatchEtag {
		t.Errorf("IfNoneMatchEtag not set matching expected %v got %v", options.IfNoneMatchEtag, headers[azcore.HeaderIfNoneMatch])
	}
}
