// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestContainerRequestOptionsToHeaders(t *testing.T) {
	options := &ContainerRequestOptions{}
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}

	options.PopulateQuotaInfo = true
	header := options.toHeaders()
	if header == nil {
		t.Error("toHeaders should return non-nil")
	}

	headers := *header
	if headers[cosmosHeaderPopulateQuotaInfo] != "true" {
		t.Errorf("PopulateQuotaInfo not set matching expected %v got %v", true, headers[cosmosHeaderPopulateQuotaInfo])
	}

	options.PopulateQuotaInfo = false
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}
}
