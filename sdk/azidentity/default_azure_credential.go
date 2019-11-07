package azidentity

// -------------------- NOTES ------------------------------
/*
Currently all of the languages implement the DefaultAzureCredential as an abstraction over the ChainedTokenCredential
DAC only calls env cred and msi cred (sdks with msal include shared token cache)
There is no guarantee that the credential type used will not change

*/
import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	developerSignOnClientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
)

// TODO: check how we're going to implement this, delete probably
// DefaultAzureCredential provides a default ChainedTokenCredential configuration for applications that will be deployed to Azure.  The following credential
// types will be tried, in order:
// - EnvironmentCredential
// - ManagedIdentityCredential
// Consult the documentation of these credential types for more information on how they attempt authentication.
type DefaultAzureCredential struct {
	defaultCredentialChain *[]azcore.TokenCredential
}

// NewDefaultTokenCredential provides a default ChainedTokenCredential configuration for applications that will be deployed to Azure.  The following credential
// types will be tried, in order:
// - EnvironmentCredential
// - ManagedIdentityCredential
// Consult the documentation of these credential types for more information on how they attempt authentication.
func NewDefaultTokenCredential(o *IdentityClientOptions) (*ChainedTokenCredential, error) {
	// CP: This is fine because we are not calling GetToken we are simple creating the new EnvironmentClient
	envClient, err := NewEnvironmentCredential(o)
	if err != nil {
		// CP: Should we return here? OR allow the program to continue? There could be a check later on to see if both types of credentials return an error
		return nil, err
	}
	// TODO: check this implementation:
	// 1. params for constructor should be nilable
	// 2. Should this func ask for a client id? or get it from somewhere else? client id is optional anyways so it is also not necessary
	msiClient, err := NewManagedIdentityCredential("", o)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultTokenCredential: %w", err)
	}

	return NewChainedTokenCredential(
		envClient,
		msiClient,
		&credentialNotFoundGuard{})
}

type credentialNotFoundGuard struct {
	azcore.TokenCredential
}

func (c *credentialNotFoundGuard) GetToken(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	return nil, &CredentialUnavailableError{Message: "example"}
}
