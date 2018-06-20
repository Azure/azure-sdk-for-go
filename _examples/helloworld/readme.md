# Hello World Producer / Consumer

This example illustrates a producer sending messages into a Service Bus FIFO Queue. The consumer
receives from each message in FIFO order from the queue, and outputs the message it receives. Upon entering 'exit' into the
producer, the producer will send, then exit, and the receiver will receive the message and close.

## Setup
- Create an Azure Service Bus namespace with a queue named "helloworld" in the [Azure Portal](https://protal.azure.com).
- After creation copy the Service Bus connection string and use it as shown below.

## To Run
- from this directory execute `make`
- open two terminal windows
  - in the first terminal, execute `SERVICEBUS_CONNECTION_STRING='your-connstring' ./bin/consumer`
  - in the second terminal, execute `SERVICEBUS_CONNECTION_STRING='your-connstring' ./bin/producer`
  - in the second terminal, type some works and press enter
- see the words you typed in the second terminal in the first
- type 'exit' in the second terminal when you'd like to end your session

## [Producer](./producer/main.go)
Send messages to the consumer process.

## [Consumer](./consumer/main.go)
Receive messages from the producer process.