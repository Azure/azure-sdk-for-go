module github.com/Azure/azure-sdk-for-python/sdk/samples/azidentity/manual-tests/managed-identity/functions

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets v0.13.0
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.6.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal v0.8.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.0.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)

// intent is to test main:HEAD
replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../../../../../azidentity
