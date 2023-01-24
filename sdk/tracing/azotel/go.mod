module github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.3.0
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/sdk v1.11.2
	go.opentelemetry.io/otel/trace v1.11.2
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../../azcore
