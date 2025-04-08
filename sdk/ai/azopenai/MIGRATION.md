# Migration Guide from Azure OpenAI SDK v0.7.x to v0.8.0

## Overview

Starting with version 0.8.0, this library has been updated to function as a companion to the [official OpenAI Go client library](https://github.com/openai/openai-go). This new approach offers several benefits:

- Consistent API experience between Azure OpenAI and OpenAI services
- Direct access to the latest OpenAI features through the official library
- Azure-specific extensions available through this companion library

## Key Changes

### New Dependency

Your projects will now need to include the official OpenAI Go client:

```go
import (
    "github.com/openai/openai-go"
    "github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
)
```

### Authentication and Client Creation

Instead of using the Azure OpenAI client directly for all operations, you'll now:

1. Create an OpenAI client configured for Azure
2. Use the Azure OpenAI companion library for Azure-specific extensions

Example:

```go
import (
    "github.com/openai/openai-go"
    "github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
)

// Create the OpenAI client configured for Azure
azureOpenAIKey := os.Getenv("AZURE_OPENAI_KEY")
azureOpenAIEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

// Create the main client using the OpenAI Go SDK
client, err := openai.New(
    openai.WithAzureOptions(azureOpenAIEndpoint, azureOpenAIKey),
)
if err != nil {
    // Handle error
}

// For Azure-specific features, use the companion library
// Example of using an Azure-specific feature
```

### API Changes

Please refer to the [official OpenAI Go client documentation](https://github.com/openai/openai-go) for details on the standard API operations.

For Azure-specific extensions provided by this companion library, see the reference documentation and examples.

## Common Migration Scenarios

### Chat Completions

Before:
```go
client, err := azopenai.NewClient(endpoint, azcore.NewKeyCredential(key), nil)
resp, err := client.GetChatCompletions(context.TODO(), "my-deployment", azopenai.ChatCompletionsOptions{
    // options
})
```

After:
```go
// Create the main OpenAI client configured for Azure
client, err := openai.New(
    openai.WithAzureOptions(endpoint, key),
)

// Use the standard OpenAI client for chat completions
resp, err := client.CreateChatCompletion(context.TODO(), &openai.ChatCompletionRequest{
    Model: "my-deployment", // Azure deployment name goes here
    // other options
})

// For Azure-specific extensions, use the companion library
```

### Chat Completions (On Your Data)

### Embeddings

## Additional Resources

- [OpenAI Go Client Documentation](https://github.com/openai/openai-go)
- [Azure OpenAI Service Documentation](https://learn.microsoft.com/azure/ai-services/openai/)
