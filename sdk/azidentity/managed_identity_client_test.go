// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"net/url"
	"testing"
)

func TestIMDSEndpointParse(t *testing.T) {
	_, err := url.Parse(imdsEndpoint)
	if err != nil {
		t.Fatalf("Failed to parse the IMDS endpoint: %v", err)
	}
}

// func TestNewDefaultMSIPipeline(t *testing.T) {
// 	p := newDefaultMSIPipeline(ManagedIdentityCredentialOptions{})
// }
