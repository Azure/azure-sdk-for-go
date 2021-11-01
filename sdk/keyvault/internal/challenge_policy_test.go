//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import "testing"

func TestParseTenantID(t *testing.T) {
	sampleURL := "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000"
	tenant := parseTenant(sampleURL)
	expectedTenant := "00000000-0000-0000-0000-000000000000"
	if *tenant != expectedTenant {
		t.Fatalf("tenant was not properly parsed, got %s, expected %s", *tenant, expectedTenant)
	}
}
