// Deprecated: use github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai instead
module github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai

go 1.18

retract [v0.1.0, v0.1.2] // use github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai instead.

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.6.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.4.0
	github.com/joho/godotenv v1.3.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/AzureAD/microsoft-authentication-library-for-go v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dnaeon/go-vcr v1.2.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
