module examples

go 1.17

require (
	github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus v0.3.6
	github.com/joho/godotenv v1.4.0
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.8.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/messaging/internal v0.0.0-20211208010914-2b10e91d237e // indirect
	github.com/Azure/go-amqp v0.17.4 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	// temporary until https://github.com/nhooyr/websocket/pull/310 is merged and released.
	github.com/gin-gonic/gin v1.7.7 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus v0.3.6 => ../../
