module github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0
	github.com/golang/mock v1.6.0
	github.com/joho/godotenv v1.4.0
	github.com/stretchr/testify v1.7.1
	nhooyr.io/websocket v1.8.7
)

require (
	code.cloudfoundry.org/clock v0.0.0-20180518195852-02e53af36e6c // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v0.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220511200225-c6db032c6c88 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// used in stress tests
require github.com/microsoft/ApplicationInsights-Go v0.4.4
