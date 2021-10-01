module github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus

go 1.16

require (
	github.com/Azure/azure-amqp-common-go/v3 v3.2.0
	github.com/Azure/azure-sdk-for-go v51.1.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.19.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.11.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.7.1 // indirect
	github.com/Azure/go-amqp v0.15.0
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/adal v0.9.13
	github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/devigned/tab v0.1.1
	github.com/joho/godotenv v1.3.0
	github.com/jpillora/backoff v1.0.0
	github.com/microsoft/ApplicationInsights-Go v0.4.4
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	nhooyr.io/websocket v1.8.6
)

// remove after go-amqp releases https://github.com/Azure/go-amqp/pull/72
replace github.com/Azure/go-amqp v0.15.0 => github.com/Azure/go-amqp v0.15.1-0.20210923181113-8f9a02b39d60
