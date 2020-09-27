// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"net/url"
	"testing"
)

func TestAzurePublicCloudParse(t *testing.T) {
	_, err := url.Parse(AzurePublicCloud)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}

func TestAzureChinaParse(t *testing.T) {
	_, err := url.Parse(AzureChina)
	if err != nil {
		t.Fatalf("Failed to parse AzureChina authority host: %v", err)
	}
}

func TestAzureGermanyParse(t *testing.T) {
	_, err := url.Parse(AzureGermany)
	if err != nil {
		t.Fatalf("Failed to parse AzureGermany authority host: %v", err)
	}
}

func TestAzureGovernmentParse(t *testing.T) {
	_, err := url.Parse(AzureGovernment)
	if err != nil {
		t.Fatalf("Failed to parse AzureGovernment authority host: %v", err)
	}
}
