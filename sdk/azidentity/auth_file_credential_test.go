// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"testing"
)

const (
	scopeToResource = "https://storage.azure.com/.default"
)

func TestAuthFileCredential_BadSdkAuthFilePathThrowsDuringGetToken(t *testing.T) {
	_, err := NewAuthFileCredential("Bougs*File*Path", nil)
	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}
