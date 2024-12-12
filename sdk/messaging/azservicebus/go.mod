module github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus

go 1.18

retract v1.1.2 // Breaks customers in situations where close is slow/infinite.

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.7.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0
	github.com/Azure/go-amqp v1.1.0
)

require (
	// used in tests only
	github.com/joho/godotenv v1.5.1

	// used in stress tests
	// github.com/microsoft/ApplicationInsights-Go v0.4.4
	github.com/stretchr/testify v1.9.0

	// used in examples only
	nhooyr.io/websocket v1.8.11
)

require github.com/golang/mock v1.6.0

require (
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
