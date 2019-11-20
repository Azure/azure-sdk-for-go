// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"fmt"
	"testing"
)

func TestNilCredentialInChain(t *testing.T) {
	var unavailableError *CredentialUnavailableError
	cred := NewClientSecretCredential("expected_tenant", "client", "secret", nil)

	_, err := NewChainedTokenCredential(cred, nil, cred)
	if err != nil {
		switch i := err.(type) {
		case *CredentialUnavailableError:
		default:
			t.Errorf("Actual error: %v, Expected error: %v, wrong type %t", err, unavailableError, i)
		}
	}
}

func TestNilChain(t *testing.T) {
	var unavailableError *CredentialUnavailableError

	_, err := NewChainedTokenCredential()
	if err != nil {
		switch i := err.(type) {
		case *CredentialUnavailableError:
			fmt.Println("Received: ", err.Error())
		default:
			t.Errorf("Actual error: %v, Expected error: %v, wrong type %t", err, unavailableError, i)
		}
	}
}
