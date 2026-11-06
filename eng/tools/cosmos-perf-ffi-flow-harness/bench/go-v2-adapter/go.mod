module cosmos-perf-c0-control

go 1.27.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.22.0
	github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos v1.5.0
	github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2 v2.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/Azure/azure-cosmos-driver/darwin/arm64 v0.2.0 // indirect
	github.com/Azure/azure-cosmos-driver/linux/amd64 v0.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2 => /mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/src/go/control

replace github.com/Azure/azure-cosmos-driver/linux/amd64 => /mnt/cosmos-perf/cpu-isolation-synthetic-20261105-d/build/work/c0-control/native
