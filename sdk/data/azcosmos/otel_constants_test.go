// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/url"
	"testing"
)

func TestSpanForClient(t *testing.T) {
	endpoint, _ := url.Parse("https://localhost:8081/")
	aSpan, err := getSpanNameForClient(endpoint, operationTypeQuery, resourceTypeDatabase, "test")
	if err != nil {
		t.Fatalf("Failed to get span name: %v", err)
	}
	if aSpan.name != "query_databases test" {
		t.Fatalf("Expected span name to be 'query_databases test', but got %s", aSpan.name)
	}
	if len(aSpan.options.Attributes) == 0 {
		t.Fatalf("Expected span options to have attributes, but got none")
	}
}
