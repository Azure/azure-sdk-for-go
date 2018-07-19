# Microsoft Azure Service Bus Client for Golang
[![Go Report Card](https://goreportcard.com/badge/github.com/Azure/azure-service-bus-go)](https://goreportcard.com/report/github.com/Azure/azure-service-bus-go)
[![godoc](https://godoc.org/github.com/Azure/azure-service-bus-go?status.svg)](https://godoc.org/github.com/Azure/azure-service-bus-go)
[![Build Status](https://travis-ci.org/Azure/azure-service-bus-go.svg?branch=master)](https://travis-ci.org/Azure/azure-service-bus-go)
[![Coverage Status](https://coveralls.io/repos/github/Azure/azure-service-bus-go/badge.svg?branch=master)](https://coveralls.io/github/Azure/azure-service-bus-go?branch=master)

Microsoft Azure Service Bus is a reliable cloud messaging service (MaaS) which simplifies enterprise cloud messaging. It
enables developers to build scalable cloud solutions and implement complex messaging workflows over an efficient binary
protocol called AMQP.

This library provides a simple interface for sending, receiving and managing Service Bus entities such as Queues, Topics
and Subscriptions.

For more information about Service Bus, check out the [Azure documentation](https://azure.microsoft.com/en-us/services/service-bus/).

This library is a pure Golang implementation of Azure Service Bus over AMQP.

## Preview of Service Bus for Golang
This library is currently a preview. There may be breaking interface changes until it reaches semantic version `v1.0.0`. 
If you run into an issue, please don't hesitate to log a 
[new issue](https://github.com/Azure/azure-service-bus-go/issues/new) or open a pull request.

## Installing the library
To more reliably manage dependencies in your application we recommend [golang/dep](https://github.com/golang/dep).

With dep:
```
dep ensure -add github.com/Azure/azure-service-bus-go
```

With go get:
```
go get -u github.com/Azure/azure-service-bus-go/...
```

If you need to install Go, follow [the official instructions](https://golang.org/dl/)

## Using Service Bus
In this section we'll cover some basics of the library to help you get started.

### Quick start
Let's send and receive `"hello, world!"`.
```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Azure/azure-service-bus-go"
)

func main() {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		// handle error
	}

	// Initialize and create a Service Bus Queue named helloworld if it doesn't exist
	queueName := "helloworld"
	q, err := ns.NewQueue(queueName)
	if err != nil {
		// handle queue creation error
	}

	// Send message to the Queue named helloworld
	err = q.Send(context.Background(), servicebus.NewMessageFromString("Hello World!"))
	if err != nil {
		// handle message send error
	}

	// Receive message from queue named helloworld
	listenHandle, err := q.Receive(context.Background(),
		func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
			fmt.Println(string(msg.Data))
			return msg.Accept()
		})
	if err != nil {
		// handle queue listener creation err
	}
	// Close the listener handle for the Service Bus Queue
	defer listenHandle.Close(context.Background())

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
}
```

## Examples
- [HelloWorld: Producer and Consumer](./_examples/helloworld): an example of sending and receiving messages from a 
Service Bus Queue.
- [OpenTracing Support](./_examples/opentracing): an example of how to enable distributed tracing using [OpenTracing](http://opentracing.io)
- [Scheduled Messages](./_examples/scheduledMessage): an example of how to send a message to be delivered in the future

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
