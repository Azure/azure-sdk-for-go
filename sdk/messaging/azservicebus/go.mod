module github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus

go 1.16

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.13.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.8.3
	github.com/Azure/azure-sdk-for-go/sdk/messaging/internal v0.0.0-20211208010914-2b10e91d237e
	github.com/Azure/go-amqp v0.17.0
	github.com/devigned/tab v0.1.1
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

require (

	// temporary until https://github.com/nhooyr/websocket/pull/310 is merged and released.
	github.com/gin-gonic/gin v1.7.7 // indirect
	// used in tests only
	github.com/joho/godotenv v1.3.0

	// used in stress tests
	github.com/microsoft/ApplicationInsights-Go v0.4.4
	github.com/stretchr/testify v1.7.0

	// used in examples only
	nhooyr.io/websocket v1.8.6
)
