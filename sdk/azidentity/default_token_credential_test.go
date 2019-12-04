package azidentity

import (
	"errors"
	"testing"
)

func TestDefaultTokenCredential_ExcludeEnvCredential(t *testing.T) {
	cred, err := NewDefaultTokenCredential(&DefaultTokenCredentialOptions{ExcludeEnvironmentCredential: true})
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}

	if len(cred.sources) != 1 {
		t.Fatalf("Length of ChainedTokenCredential sources for DefaultTokenCredential. Expected: 1, Received: %d", len(cred.sources))
	}

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
	var authFailedErr *AuthenticationFailedError
	_, err = NewDefaultTokenCredential(&DefaultTokenCredentialOptions{ExcludeEnvironmentCredential: false, ExcludeMSICredential: true})
	if err == nil {
		t.Fatalf("Expected an error but received nil")
	}
	if !errors.As(err, &authFailedErr) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}

}

func TestDefaultTokenCredential_NilOptions(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unexpected error when initializing environment variables: %v", err)
	}
	cred, err := NewDefaultTokenCredential(nil)
	if err != nil {
		t.Fatalf("Did not expect to receive an error in creating the credential")
	}
	if len(cred.sources) != 1 {
		t.Fatalf("Length of ChainedTokenCredential sources for DefaultTokenCredential. Expected: 1, Received: %d", len(cred.sources))
	}
}
