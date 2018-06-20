# Schedule Message Example

This example illustrates how to send a scheduled message using Azure Service Bus. The key to the example is to set the
ScheduledEnqueuedTime on the message before it is sent to the broker. The message will be delivered after the specified
UTC time.

## Setup
- Create an Azure Service Bus namespace with a queue named "helloworld" in the [Azure Portal](https://protal.azure.com).
- After creation copy the Service Bus connection string and use it as shown below.


## To Run
From the root directory run 
```
SERVICEBUS_CONNECTION_STRING='your-SB-conn-string' go run main.go
```