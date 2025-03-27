module github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys/testdata/perf

go 1.23.0

replace github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys => ../..

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.17.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.2
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys v1.3.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal v1.2.0-beta.1 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.4.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
