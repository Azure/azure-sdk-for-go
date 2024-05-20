// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestContainerRequestOptionsToHeaders(t *testing.T) {
	options := &ReadContainerOptions{}
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}

	options.PopulateQuotaInfo = true
	header := options.toHeaders()
	if header == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	headers := *header
	if headers[cosmosHeaderPopulateQuotaInfo] != "true" {
		t.Errorf("PopulateQuotaInfo not set matching expected %v got %v", true, headers[cosmosHeaderPopulateQuotaInfo])
	}

	options.PopulateQuotaInfo = false
	if options.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}

	replaceOptions := &ReplaceContainerOptions{}
	if replaceOptions.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}

	deleteOptions := &DeleteContainerOptions{}
	if deleteOptions.toHeaders() != nil {
		t.Error("toHeaders should return nil")
	}
}

func TestQueryContainersRequestOptionsToHeaders(t *testing.T) {
	options := &QueryContainersOptions{}
	continuation := "continuationToken"
	options.ContinuationToken = &continuation
	header := options.toHeaders()
	if header == nil {
		t.Fatal("toHeaders should return non-nil")
	}

	headers := *header
	if headers[cosmosHeaderContinuationToken] != "continuationToken" {
		t.Errorf("ContinuationToken should be continuationToken but got %v", headers[cosmosHeaderContinuationToken])
	}
}
