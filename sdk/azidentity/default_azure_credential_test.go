// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"os"
	"testing"
)

const (
	lengthOfChainOneExcluded = 2
	lengthOfChainFull        = 3
)

func TestDefaultAzureCredential_ExcludeEnvCredential(t *testing.T) {
	resetEnvironmentVarsForTest()
	_ = os.Setenv("MSI_ENDPOINT", "http://localhost:3000")
	cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{ExcludeEnvironmentCredential: true})
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}

	if len(cred.sources) != lengthOfChainOneExcluded {
		t.Fatalf("Length of ChainedTokenCredential sources for DefaultAzureCredential. Expected: %d, Received: %d", lengthOfChainOneExcluded, len(cred.sources))
	}
	_ = os.Setenv("MSI_ENDPOINT", "")
}

func TestDefaultAzureCredential_ExcludeMSICredential(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{ExcludeMSICredential: true})
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}
	if len(cred.sources) != lengthOfChainOneExcluded {
		t.Fatalf("Length of ChainedTokenCredential sources for DefaultAzureCredential. Expected: %d, Received: %d", lengthOfChainOneExcluded, len(cred.sources))
	}
}

func TestDefaultAzureCredential_ExcludeAzureCLICredential(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	_ = os.Setenv("MSI_ENDPOINT", "http://localhost:3000")
	cred, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{ExcludeAzureCLICredential: true})
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}

	if len(cred.sources) != lengthOfChainOneExcluded {
		t.Fatalf("Length of ChainedTokenCredential sources for DefaultAzureCredential. Expected: %d, Received: %d", lengthOfChainOneExcluded, len(cred.sources))
	}
	clearEnvVars("MSI_ENDPOINT")
	resetEnvironmentVarsForTest()
}

func TestDefaultAzureCredential_ExcludeAllCredentials(t *testing.T) {
	resetEnvironmentVarsForTest()
	var credUnavailable *CredentialUnavailableError
	_, err := NewDefaultAzureCredential(&DefaultAzureCredentialOptions{
		ExcludeEnvironmentCredential: true,
		ExcludeMSICredential:         true,
		ExcludeAzureCLICredential:    true,
	})
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
	if !errors.As(err, &credUnavailable) {
		t.Fatalf("Expected: CredentialUnavailableError, Received: %T", err)
	}

}

func TestDefaultAzureCredential_NilOptions(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}
	options := DefaultManagedIdentityCredentialOptions()
	c := newManagedIdentityClient(&options)
	// if the test is running in a MSI environment then the length of sources would be two since it will include environment credential and managed identity credential
	if msiType, err := c.getMSIType(); !(msiType == msiTypeUnavailable || msiType == msiTypeUnknown) {
		if len(cred.sources) != lengthOfChainFull {
			t.Fatalf("Length of ChainedTokenCredential sources for DefaultAzureCredential. Expected: %d, Received: %d", lengthOfChainFull, len(cred.sources))
		}
		//if a credential unavailable error is received or msiType is unknown then only the environment credential will be added
	} else if unavailableErr := (*CredentialUnavailableError)(nil); errors.As(err, &unavailableErr) || msiType == msiTypeUnknown {
		if len(cred.sources) != lengthOfChainOneExcluded {
			t.Fatalf("Length of ChainedTokenCredential sources for DefaultAzureCredential. Expected: %d, Received: %d", lengthOfChainOneExcluded, len(cred.sources))
		}
		// if there is some other unexpected error then we fail here
	} else if err != nil {
		t.Fatalf("Received an error when trying to determine MSI type: %v", err)
	}
}
