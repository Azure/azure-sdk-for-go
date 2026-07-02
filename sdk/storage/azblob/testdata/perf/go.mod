module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf

go 1.25.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.21.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.4
	github.com/stretchr/testify v1.11.1
	golang.org/x/sys v0.46.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// TODO: This will be removed once the internal module changes are released
replace github.com/Azure/azure-sdk-for-go/sdk/internal => ../../../../internal
