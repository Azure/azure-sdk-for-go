// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	scopeToResource = "https://storage.azure.com/.default"
)

func TestAuthFileCredential_BadSdkAuthFilePathThrowsDuringGetToken(t *testing.T) {
	cred, err := NewAuthFileCredential("Bougs*File*Path", nil)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scopeToResource}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}
