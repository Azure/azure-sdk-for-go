module github.com/Azure/azure-sdk-for-go/sdk/azcore/testdata/perf

go 1.24.0

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.20.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.opentelemetry.io/otel v1.41.0 // indirect
	go.opentelemetry.io/otel/trace v1.41.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../../
