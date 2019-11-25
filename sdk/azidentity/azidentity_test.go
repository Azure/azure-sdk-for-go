// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"net/url"
	"testing"
)

func TestAuthorityHostParse(t *testing.T) {
	_, err := url.Parse(defaultAuthorityHost)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}

func TestNilAuthFailedParam(t *testing.T) {
	var err *AuthenticationFailedError
	authFailed := newAuthenticationFailedError(nil)

	if !errors.As(authFailed, &err) {
		t.Fatalf("Failed to create the right error type")
	}
}
