module github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus

go 1.18

retract v1.1.2 // Breaks customers in situations where close is slow/infinite.

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.17.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.2
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.1-0.20250301010503-a03cfeb12d70
	github.com/Azure/go-amqp v1.4.0
)

require (
	// used in tests only
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/coder/websocket v1.8.12
	github.com/golang/mock v1.6.0
)

require (
	github.com/AzureAD/microsoft-authentication-library-for-go v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => github.com/Azure/azure-sdk-for-go/sdk/azcore v1.17.1-0.20250224170238-070e44629503
