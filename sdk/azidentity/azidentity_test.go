// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func Test_AuthorityHost_Parse(t *testing.T) {
	_, err := url.Parse(defaultAuthorityHost)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}

func Test_NonNilTokenCredentialOptsNilAuthorityHost(t *testing.T) {
	opts := &TokenCredentialOptions{Retry: azcore.RetryOptions{MaxTries: 6}}
	opts, err := opts.setDefaultValues()
	if err != nil {
		t.Fatalf("Received an error: %v", err)
	}
	if opts.AuthorityHost == nil {
		t.Fatalf("Did not set default authority host")
	}
}
