# Troubleshooting Azure Event Hubs client library issues

This troubleshooting guide contains instructions to diagnose frequently encountered issues while using the Azure Event Hubs client library for Go.

## Table of contents

* [General Troubleshooting](#general-troubleshooting)
  * [Error Handling](#error-handling)
  * [Logging](#logging)
  * [Authentication Issues](#authentication-issues)
* [Common Error Scenarios](#common-error-scenarios)
  * [Unauthorized Access Errors](#unauthorized-access-errors)
  * [Connection Lost Errors](#connection-lost-errors)
  * [Ownership Lost Errors](#ownership-lost-errors)
* [Event Processor Troubleshooting](#event-processor-troubleshooting)
  * [Load Balancing Issues](#load-balancing-issues)
  * [Checkpoint Store Problems](#checkpoint-store-problems)
  * [Performance Considerations](#performance-considerations)
* [Connectivity Issues](#connectivity-issues)
  * [Enterprise Environments and Firewalls](#enterprise-environments-and-firewalls)
  * [Using WebSockets Transport](#using-websockets-transport)
  * [Working with Proxies](#working-with-proxies)
* [Advanced Troubleshooting](#advanced-troubleshooting)
  * [Logs to Collect](#logs-to-collect)
  * [Interpreting Logs](#interpreting-logs)
* [Additional Resources](#additional-resources)
  * [Filing GitHub Issues](#filing-github-issues)

## General Troubleshooting

### Error Handling

The Event Hubs client library provides strongly-typed error handling through the `azeventhubs.Error` type with specific error codes that can be checked programmatically. This allows you to handle different error scenarios in your code.

```go
if err != nil {
    var azError *azeventhubs.Error
    if errors.As(err, &azError) {
        switch azError.Code {
        case azeventhubs.ErrorCodeUnauthorizedAccess:
            // Handle authentication errors
        case azeventhubs.ErrorCodeConnectionLost:
            // Handle connection problems
        case azeventhubs.ErrorCodeOwnershipLost:
            // Handle partition ownership changes
        }
    }
    // Handle other error types
}
```

### Logging

Event Hubs uses the classification-based logging implementation in `azcore`. You can enable logging for all Azure SDK modules by setting the environment variable `AZURE_SDK_GO_LOGGING` to `all`.

For more fine-grained control, use the `azcore/log` package to enable specific log events:

```go
import (
    "fmt"
    azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
    "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
)

// Print log output to stdout
azlog.SetListener(func(event azlog.Event, s string) {
    fmt.Printf("[%s] %s\n", event, s)
})

// Enable specific event types
azlog.SetEvents(
    azeventhubs.EventConn,    // Connection-related events
    azeventhubs.EventAuth,    // Authentication events
    azeventhubs.EventProducer, // Producer operations
    azeventhubs.EventConsumer, // Consumer operations
)
```

When troubleshooting, it's recommended to enable all Event Hubs log events to get comprehensive information.

### Authentication Issues

Authentication errors typically manifest as `azeventhubs.ErrorCodeUnauthorizedAccess`. Common causes include:

1. **Expired credentials**: Ensure your credentials haven't expired, especially when using SAS tokens
2. **Insufficient permissions**: Verify that the service principal or managed identity has appropriate permissions
3. **Invalid connection string**: Double-check your connection string for typos or incorrect values
4. **Network restrictions**: Verify that network security rules allow connections to Event Hubs

For more help with troubleshooting authentication errors when using Azure Identity, see the Azure Identity client library [troubleshooting guide][azidentity_troubleshooting].

## Common Error Scenarios

### Unauthorized Access Errors

If you receive an `ErrorCodeUnauthorizedAccess` error, it means the credentials provided are not valid for use with a particular entity, or they have expired.

**Common causes and solutions:**

- **Expired credentials**: If using SAS tokens, they expire after a certain duration. Generate a new token or use a credential that automatically refreshes.
- **Missing permissions**: Ensure the identity you're using has the correct RBAC roles assigned (e.g., "Azure Event Hubs Data Sender", "Azure Event Hubs Data Receiver").
- **Incorrect entity name**: Verify that the Event Hub name, consumer group, or namespace name is spelled correctly.

**Example error handling:**

```go
if err != nil {
    var ehError *azeventhubs.Error
    if errors.As(err, &ehError) && ehError.Code == azeventhubs.ErrorCodeUnauthorizedAccess {
        // Handle unauthorized access - most likely credentials issue
        // Log error and possibly attempt to refresh credentials
    }
}
```

### Connection Lost Errors

An `ErrorCodeConnectionLost` error indicates that the connection was lost and all retry attempts failed. This typically reflects an extended outage or connection disruption.

**Common causes and solutions:**

- **Network instability**: Check your network connection and try again after ensuring stability.
- **Service outage**: Check the [Azure status page](https://status.azure.com) for any ongoing Event Hubs outages.
- **Firewall or proxy issues**: Ensure firewall rules aren't blocking the connection.
- **Service throttling**: If sending too many events at once, try increasing your batch intervals.

**Recovery strategy:**

The client will automatically attempt to recover from transient issues. For persistent issues:

1. Implement a backoff policy in your application
2. Consider using the WebSockets transport mode in restrictive network environments (see [Using WebSockets Transport](#using-websockets-transport))

### Ownership Lost Errors

An `ErrorCodeOwnershipLost` error occurs when a partition that you were reading from was opened by another link with a higher epoch/owner level.

**Common causes:**

- Another consumer instance with a higher epoch value claimed the partition
- When using the `Processor`, this is a normal part of partition rebalancing
- Manually creating multiple PartitionClients for the same partition with different epochs

**Solutions:**

- If not using `Processor`, ensure only one client reads from a given partition at a time
- If using `Processor`, this is normal behavior during rebalancing and doesn't require action
- Allow the `Processor` to handle partition ownership management automatically

## Event Processor Troubleshooting

### Load Balancing Issues

The `Processor` uses a load balancer to distribute partition ownership across multiple instances. Issues can occur when:

**Partitions aren't distributed evenly:**

- Check your `ProcessorStrategy` setting
  - `ProcessorStrategyBalanced` aims for equal distribution (default)
  - `ProcessorStrategyGreedy` allows a single processor to claim all partitions if possible

**Excessive rebalancing:**

- If partitions constantly switch owners, consider increasing `ProcessorOptions.LoadBalancingOptions.PartitionOwnershipExpirationInterval`
- Add a delay between processor startup when launching multiple instances

**Code example for processor configuration:**

```go
processorOptions := &azeventhubs.ProcessorOptions{
    LoadBalancingOptions: &azeventhubs.ProcessorLoadBalancingOptions{
        PartitionOwnershipExpirationInterval: 30 * time.Second,
    },
    Strategy: azeventhubs.ProcessorStrategyBalanced,
}
```

### Checkpoint Store Problems

The checkpoint store is crucial for maintaining processing state across `Processor` instances.

**Common issues:**

1. **Storage permission errors**: Ensure your credential has write access to the storage container
2. **Container doesn't exist**: Create the container before running the processor
3. **High storage latency**: If checkpoint operations take too long, consider upgrading your storage account tier

**Addressing checkpoint issues:**

```go
// Create a container client with appropriate permissions
containerClient, err := azblob.NewContainerClient(containerURL, credential, nil)
if err != nil {
    // Handle error
}

// Create the container if it doesn't exist
_, err = containerClient.Create(context.Background(), nil)
if err != nil {
    var storageErr *azblob.StorageError
    if errors.As(err, &storageErr) && storageErr.ErrorCode == "ContainerAlreadyExists" {
        // Container already exists, which is fine
    } else {
        // Handle error
    }
}

// Create the checkpoint store
checkpointStore, err := checkpoints.NewBlobStore(containerClient, nil)
if err != nil {
    // Handle error
}
```

### Performance Considerations

**If the processor can't keep up with event flow:**

1. **Increase processor instances**: Add more processor instances across machines to distribute the load
2. **Optimize handler code**: Make sure your event handler function is efficient
3. **Batch checkpoints**: Don't checkpoint after every event; batch them periodically
4. **Scale up hardware**: If CPU or memory is a bottleneck, increase machine resources
5. **Increase Event Hubs partitions**: Consider creating an Event Hub with more partitions (requires a new Event Hub)

**Example of batched checkpointing:**

```go
var processedCount int
var lastCheckpoint time.Time

// In your handler:
handler := func(ctx context.Context, event *azeventhubs.ReceivedEventData) error {
    // Process the event
    processedCount++
    
    // Checkpoint every 100 messages or every 30 seconds
    if processedCount >= 100 || time.Since(lastCheckpoint) > 30*time.Second {
        err := pc.UpdateCheckpoint(ctx, event)
        if err == nil {
            processedCount = 0
            lastCheckpoint = time.Now()
        }
    }
    
    return nil
}
```

## Connectivity Issues

### Enterprise Environments and Firewalls

In corporate networks with strict firewall rules, you may encounter connectivity issues when connecting to Event Hubs.

**Common solutions:**

1. **Allow the necessary endpoints**: Ensure your firewall allows connectivity to `*.servicebus.windows.net` on ports 443 (HTTPS) and 5671/5672 (AMQP)
2. **Use WebSockets**: If AMQP ports are blocked, configure the client to use WebSockets which operate over standard HTTPS ports
3. **Configure network security rules**: If using Azure VNet integration, configure service endpoints or private endpoints

### Using WebSockets Transport

WebSockets allows connections through restrictive firewalls that permit only HTTPS traffic (port 443). To enable WebSockets:

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
)

// For a producer client
producerOptions := &azeventhubs.ProducerClientOptions{
    WebSocketOptions: &azeventhubs.WebSocketOptions{},  // Enable WebSockets
}
producerClient, err := azeventhubs.NewProducerClient("myhub.servicebus.windows.net", "myeventhub", credential, producerOptions)

// For a consumer client
consumerOptions := &azeventhubs.ConsumerClientOptions{
    WebSocketOptions: &azeventhubs.WebSocketOptions{},  // Enable WebSockets
}
consumerClient, err := azeventhubs.NewConsumerClient("myhub.servicebus.windows.net", "myeventhub", "myconsumergroup", credential, consumerOptions)
```

### Working with Proxies

To use Event Hubs through a proxy server:

```go
import (
    "net/http"
    "net/url"
    "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
)

// Create proxy URL
proxyURL, _ := url.Parse("http://your-proxy-address:port")

// Configure HTTP transport with the proxy
httpClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

// Configure options with the custom HTTP client
options := &azeventhubs.ProducerClientOptions{
    HTTPClient: httpClient,
    WebSocketOptions: &azeventhubs.WebSocketOptions{}, // Recommended when using proxies
}

producerClient, err := azeventhubs.NewProducerClient("myhub.servicebus.windows.net", "myeventhub", credential, options)
```

## Advanced Troubleshooting

### Logs to Collect

When troubleshooting issues with Event Hubs that you need to escalate to support or report in GitHub issues, collect the following logs:

1. **Enable DEBUG logging**: Set `AZURE_SDK_GO_LOGGING=all` and capture all output
2. **Specific Event Hubs logs**: Focus on the Event Hubs specific log events:
   - `EventConn`: Connection-related events
   - `EventAuth`: Authentication events
   - `EventProducer`: Producer operations
   - `EventConsumer`: Consumer operations
3. **Timeframe**: Capture logs from at least 5 minutes before until 5 minutes after the issue occurs
4. **Include timestamps**: Ensure your logging setup includes timestamps

### Interpreting Logs

When analyzing Event Hubs logs:

1. **Connection errors**: Look for AMQP connection and link errors in `EventConn` logs
2. **Authentication failures**: Check `EventAuth` logs for credential or authorization failures
3. **Producer errors**: `EventProducer` logs show message send operations and errors
4. **Consumer errors**: `EventConsumer` logs show message receive operations and partition ownership changes
5. **Load balancing**: Look for ownership claims and changes in `EventConsumer` logs

## Additional Resources

- [Event Hubs Documentation](https://learn.microsoft.com/azure/event-hubs/)
- [Event Hubs Pricing](https://azure.microsoft.com/pricing/details/event-hubs/)
- [Event Hubs Quotas](https://learn.microsoft.com/azure/event-hubs/event-hubs-quotas)
- [Event Hubs FAQ](https://learn.microsoft.com/azure/event-hubs/event-hubs-faq)

### Filing GitHub Issues

When filing GitHub issues for Event Hubs, please include:

1. **Event Hub details**:
   - How many partitions?
   - What tier (Standard/Premium/Dedicated)?

2. **Client environment**:
   - Machine specifications
   - Number of client instances running
   - Go version

3. **Message patterns**:
   - Average message size
   - Throughput (messages per second)
   - Whether traffic is consistent or bursty

4. **Reproduction steps**:
   - A minimal code example that reproduces the issue
   - Steps to reproduce the problem

5. **Logs**:
   - DEBUG level logs if possible (at minimum INFO)
   - Logs from before, during, and after the issue occurs

Having this information ready will greatly help in diagnosing and resolving your issue quickly.

<!-- LINKS -->
[azidentity_troubleshooting]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/TROUBLESHOOTING.md