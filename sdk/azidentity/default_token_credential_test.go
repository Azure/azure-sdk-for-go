// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestDefaultTokenCredential_ExcludeEnvCredential(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	_ = os.Setenv("MSI_ENDPOINT", "http://localhost:3000")
	cred, err := NewDefaultTokenCredential(&DefaultTokenCredentialOptions{ExcludeEnvironmentCredential: true})
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}

	if len(cred.sources) != 1 {
		t.Fatalf("Length of ChainedTokenCredential sources for DefaultTokenCredential. Expected: 1, Received: %d", len(cred.sources))
	}
	_ = os.Setenv("MSI_ENDPOINT", "")

}

func TestDefaultTokenCredential_ExcludeMSICredential(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewDefaultTokenCredential(&DefaultTokenCredentialOptions{ExcludeMSICredential: true})
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}
	if len(cred.sources) != 1 {
		t.Fatalf("Length of ChainedTokenCredential sources for DefaultTokenCredential. Expected: 1, Received: %d", len(cred.sources))
	}

}

func TestDefaultTokenCredential_ExcludeAllCredentials(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	var credUnavailable *CredentialUnavailableError
	_, err = NewDefaultTokenCredential(&DefaultTokenCredentialOptions{ExcludeEnvironmentCredential: false, ExcludeMSICredential: true})
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
	if !errors.As(err, &credUnavailable) {
		t.Fatalf("Expected: CredentialUnavailableError, Received: %T", err)
	}

}

func TestDefaultTokenCredential_NilOptions(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	err = initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewDefaultTokenCredential(nil)
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}
	c := newManagedIdentityClient(nil)
	// if the test is running in a MSI environment then the length of sources would be two since it will include environmnet credential and managed identity credential
	if msiType, err := c.getMSIType(context.Background()); msiType == msiTypeIMDS || msiType == msiTypeCloudShell || msiType == msiTypeAppService {
		if len(cred.sources) != 2 {
			t.Fatalf("Length of ChainedTokenCredential sources for DefaultTokenCredential. Expected: 2, Received: %d", len(cred.sources))
		}
		//if a credential unavailable error is received or msiType is unknown then only the environment credential will be added
	} else if unavailableErr := (*CredentialUnavailableError)(nil); errors.As(err, &unavailableErr) || msiType == msiTypeUnknown {
		if len(cred.sources) != 1 {
			t.Fatalf("Length of ChainedTokenCredential sources for DefaultTokenCredential. Expected: 1, Received: %d", len(cred.sources))
		}
		// if there is some other unexpected error then we fail here
	} else if err != nil {
		t.Fatalf("Received an error when trying to determine MSI type: %v", err)
	}
}
