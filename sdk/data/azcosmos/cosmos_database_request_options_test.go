// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestDatabaseRequestOptionsToHeaders(t *testing.T) {
	options := &DatabaseRequestOptions{}
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}

	etagValue := azcore.ETag("etag")
	noneEtagValue := azcore.ETag("noneetag")
	options.IfMatchEtag = &etagValue
	options.IfNoneMatchEtag = &noneEtagValue
	header := options.toHeaders()
	if header == nil {
		t.Error("toHeaders should return non-nil")
	}

	headers := *header
	if headers[headerIfMatch] != string(*options.IfMatchEtag) {
		t.Errorf("IfMatchEtag not set matching expected %v got %v", options.IfMatchEtag, headers[headerIfMatch])
	}
	if headers[headerIfNoneMatch] != string(*options.IfNoneMatchEtag) {
		t.Errorf("IfNoneMatchEtag not set matching expected %v got %v", options.IfNoneMatchEtag, headers[headerIfNoneMatch])
	}
}
