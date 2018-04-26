package auth

import (
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// NewAuthorizerFromEnvironment creates a keyvault dataplane Authorizer configured from environment variables in the order:
// 1. Client credentials
// 2. Client certificate
// 3. Username password
// 4. MSI
func NewAuthorizerFromEnvironment() (autorest.Authorizer, error) {
	settings, err := auth.GetAuthSettings()
	if err != nil {
		return nil, err
	}

	settings.Resource = os.Getenv("AZURE_KEYVAULT_RESOURCE")
	if settings.Resource == "" {
		settings.Resource = strings.TrimSuffix(settings.Environment.KeyVaultEndpoint, "/")
	}

	return settings.GetAuth()
}
