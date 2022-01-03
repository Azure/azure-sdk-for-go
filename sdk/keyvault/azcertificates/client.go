package azcertificates

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client is the struct for interacting with a Key Vault Certificates instance.
type Client struct {
	genClient *generated.KeyVaultClient
	vaultURL  string
}

// ClientOptions are the optional parameters for the NewClient function
type ClientOptions struct {
	azcore.ClientOptions
}

func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	genOptions := &policy.ClientOptions{
		Logging:         options.Logging,
		PerRetryPolicies: options.PerRetryPolicies,
	}
	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)

	conn := generated.NewConnection(genOptions)

	return Client{
		genClient: generated.NewKeyVaultClient(conn),
		vaultURL:  vaultURL,
	}, nil
}
