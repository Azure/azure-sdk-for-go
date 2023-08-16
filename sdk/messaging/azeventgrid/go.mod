module github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.8.0-beta.2
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.7.0
)

replace (
	// temporary until we officially release the next beta.
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.8.0-beta.2 => ../../azcore
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dnaeon/go-vcr v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
