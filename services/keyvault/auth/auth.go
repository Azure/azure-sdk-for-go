package auth

import (
	"os"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// NewAuthorizerFromEnvironment creates a keyvault dataplane Authorizer configured from environment variables in the order:
// 1. Client credentials
// 2. Client certificate
// 3. Username password
// 4. MSI
func NewAuthorizerFromEnvironment() (autorest.Authorizer, error) {
	envName := os.Getenv("AZURE_ENVIRONMENT")
	var env azure.Environment
	var err error

	if envName == "" {
		env = azure.PublicCloud
	} else {
		env, err = azure.EnvironmentFromName(envName)
		if err != nil {
			return nil, err
		}
	}

	resource := os.Getenv("AZURE_KEYVAULT_RESOURCE")
	if resource == "" {
		resource = strings.TrimSuffix(env.KeyVaultEndpoint, "/")
	}

	return auth.NewAuthorizerFromEnvironmentWithResource(resource)
}
