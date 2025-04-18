# Migration Guide from Azure OpenAI SDK v0.7.x to v0.8.0

## Table of Contents
- [Overview](#overview)
- [Summary of Major Changes](#summary-of-major-changes)
- [Key Changes](#key-changes)
- [Authentication and Client Creation](#authentication-and-client-creation)
- [API Changes](#api-changes)
- [Common Migration Scenarios](#common-migration-scenarios)
- [Additional Resources](#additional-resources)

## Overview

Starting with version 0.8.0, the `azopenai.Client` provided by this package has been retired in favor of the [official OpenAI Go client library](https://github.com/openai/openai-go). That package contains all that is needed to connect to both the Azure OpenAI and OpenAI services. This library is now a companion, enabling Azure-specific extensions (such as Azure OpenAI On Your Data). The `azopenaiassistants` package has also been deprecated in favor of the official client.

Migrating to the official client offers:
- Consistent API experience between Azure OpenAI and OpenAI services.
- Direct access to the latest OpenAI features.

> [!NOTE]
> This document is a work-in-progress and may change to reflect updates to the package. We value your feedback—please [create an issue](https://github.com/Azure/azure-sdk-for-go/issues/new/choose) to suggest improvements or report problems with this guide or the package.

## Summary of Major Changes

| Area                | v0.7.x Approach                | v0.8.0+ Approach (Recommended)         |
|---------------------|--------------------------------|----------------------------------------|
| Client              | `azopenai.Client`              | `openai.Client`                        |
| Assistants          | `azopenaiassistants`           | `openai.Client.Beta.Assistants`        |
| Azure Extensions    | Built-in                       | Use `azopenai` as a companion          |
| API Structure       | Flat methods                   | Subclients per service category        |
| Authentication      | Azure-specific                 | Use `azure.With...` options            |

## Key Changes

### New Dependency

Your projects must now include the official OpenAI Go client:

```go
import (
    "github.com/openai/openai-go"
)
```

If you need Azure-specific extensions, also include the `azopenai` package:

```go
import (
    "github.com/openai/openai-go"
    "github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
)
```

> [!NOTE]
> **Azure extensions** refer to features unique to Azure OpenAI, such as new request/response structures (e.g., Azure OpenAI On Your Data), not basic authentication or connection differences.

## Authentication and Client Creation

Instead of using the Azure OpenAI client directly for all operations, you'll now:
- Create an OpenAI client configured for Azure.
- Use the Azure OpenAI companion library for Azure-specific extensions.

### Azure OpenAI with API Key

**Before:**
```go
endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
key := os.Getenv("AZURE_OPENAI_API_KEY")
client, err := azopenai.NewClientWithKeyCredential(endpoint, azcore.NewKeyCredential(key), nil)
if err != nil {
    panic(err)
}
```

**After:**
```go
endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
api_version := os.Getenv("AZURE_OPENAI_API_VERSION")
key := os.Getenv("AZURE_OPENAI_API_KEY")

client := openai.NewClient(
    azure.WithEndpoint(endpoint, api_version),
    azure.WithAPIKey(key),
)
```

### Azure OpenAI with Token Credentials

**Before:**
```go
endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

credential, err := azidentity.NewDefaultAzureCredential(nil)
if err != nil {
    panic(err)
}
client, err := azopenai.NewClient(endpoint, credential, nil)
if err != nil {
    panic(err)
}
```

**After:**
```go
endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
api_version := os.Getenv("AZURE_OPENAI_API_VERSION")

credential, err := azidentity.NewDefaultAzureCredential(nil)
if err != nil {
    panic(err)
}
client := openai.NewClient(
    azure.WithEndpoint(endpoint, api_version),
    azure.WithTokenCredential(credential),
)
```

### OpenAI (not Azure)

**Before:**
```go
key := os.Getenv("OPENAI_API_KEY")

client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", azcore.NewKeyCredential(key), nil)
if err != nil {
    panic(err)
}
```

**After:**
```go
key := os.Getenv("OPENAI_API_KEY")
client := openai.NewClient(
    option.WithAPIKey(key),
)
```

## API Changes

The official OpenAI Go client organizes operations into subclients for each service category, rather than providing all operations as flat methods.

| Service               | Description |
|-----------------------|-------------|
| `client.Completions`  | [Completions API](https://platform.openai.com/docs/api-reference/completions) |
| `client.Chat`         | [Chat Completions API](https://platform.openai.com/docs/api-reference/chat) |
| `client.Embeddings`   | [Embeddings API](https://platform.openai.com/docs/api-reference/embeddings) |
| `client.Files`        | [Files API](https://platform.openai.com/docs/api-reference/files) |
| `client.Images`       | [Images API](https://platform.openai.com/docs/api-reference/images) |
| `client.Audio`        | [Audio API](https://platform.openai.com/docs/api-reference/audio) |
| `client.Moderations`  | [Moderations API](https://platform.openai.com/docs/api-reference/moderations) |
| `client.Models`       | [Models API](https://platform.openai.com/docs/api-reference/models) |
| `client.FineTuning`   | [Fine-tuning API](https://platform.openai.com/docs/api-reference/fine-tuning) |
| `client.VectorStores` | [Vector Stores API](https://platform.openai.com/docs/api-reference/vector-stores) |
| `client.Beta`         | Beta features (e.g., [Assistants](https://platform.openai.com/docs/api-reference/assistants), [Threads](https://platform.openai.com/docs/api-reference/threads)) |
| `client.Batches`      | [Batch API](https://platform.openai.com/docs/api-reference/batch) |
| `client.Uploads`      | [Uploads API](https://platform.openai.com/docs/api-reference/uploads) |
| `client.Responses`    | [Responses API](https://platform.openai.com/docs/api-reference/responses) |

Refer to the [official OpenAI Go client documentation](https://github.com/openai/openai-go) for details.

For Azure-specific extensions, see the reference documentation and examples in this companion library.

## Common Migration Scenarios

### Chat Completions

**Before:**
```go
resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
    // DeploymentName: "gpt-4o", // This only applies for the OpenAI service.
    Messages: []azopenai.ChatRequestMessageClassification{
        &azopenai.ChatRequestUserMessage{
            Content: azopenai.NewChatRequestUserMessageContent("What is OpenAI, in 20 words or less?"),
        },
    },
}, nil)
if err != nil {
    return err
}
for _, choice := range resp.Choices {
    // Process the response content from each choice
    // choice.Message.Content contains the message text
}
```

**After:**
```go
deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME")
resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
    Model: openai.ChatModel(deployment), // For Azure OpenAI, deployment name is used as the model.
    // Model: openai.ChatModelGPT4o, // For OpenAI, use the model name.
    Messages: []openai.ChatCompletionMessageParamUnion{
        {
            OfUser: &openai.ChatCompletionUserMessageParam{
                Content: openai.ChatCompletionUserMessageParamContentUnion{
                    OfString: openai.String("What is OpenAI, in 20 words or less?"),
                },
            },
        },
    },
})

if err != nil {
    return err
}

for _, choice := range resp.Choices {
    // Process the response content from each choice
    // choice.Message.Content contains the message text
}
```

#### Streaming Chat Completions

**Before:**
```go
resp, err := client.GetChatCompletionsStream(context.TODO(), azopenai.ChatCompletionsStreamOptions{
    // DeploymentName: "gpt-4o", // This only applies for the OpenAI service.
    Messages: []azopenai.ChatRequestMessageClassification{
        &azopenai.ChatRequestUserMessage{
            Content: azopenai.NewChatRequestUserMessageContent("What is OpenAI, in 20 words or less?"),
        },
    },
}, nil)
if err != nil {
    return err
}
defer resp.ChatCompletionsStream.Close()

for {
    entry, err := resp.ChatCompletionsStream.Read()

    if errors.Is(err, io.EOF) {
        break
    }

    if err != nil {
        return err
    }

    for _, choice := range entry.Choices {
        // Process each chunk of streaming content
        // choice.Message.Content contains the partial message
    }
}
```

**After:**
```go
deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME")
stream := client.Chat.Completions.NewStreaming(context.TODO(), openai.ChatCompletionNewParams{
    Model: openai.ChatModel(deployment), // For Azure OpenAI, deployment name is used as the model.
    // Model: openai.ChatModelGPT4o, // For OpenAI, use the model name.
    Messages: []openai.ChatCompletionMessageParamUnion{
        {
            OfUser: &openai.ChatCompletionUserMessageParam{
                Content: openai.ChatCompletionUserMessageParamContentUnion{
                    OfString: openai.String("What is OpenAI, in 20 words or less?"),
                },
            },
        },
    },
})

for stream.Next() {
    chunk := stream.Current()

    for _, choice := range chunk.Choices {
        // Process each chunk of streaming content
        // choice.Delta.Content contains the partial message
    }
}
```

### Chat Completions (On Your Data)

**Before:**
```go
resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
    Messages: []azopenai.ChatRequestMessageClassification{
        &azopenai.ChatRequestUserMessage{
            Content: azopenai.NewChatRequestUserMessageContent("Your message here"),
        },
    },
    AzureExtensionsOptions: []azopenai.AzureChatExtensionConfigurationClassification{
        &azopenai.AzureSearchChatExtensionConfiguration{
            Parameters: &azopenai.AzureSearchChatExtensionParameters{
                Endpoint:       &search_endpoint,
                IndexName:      &search_index,
                Authentication: &azopenai.OnYourDataSystemAssignedManagedIdentityAuthenticationOptions{},
            },
        },
    },
}, nil)

// Access citations from the response
for _, choice := range resp.Choices {
    // Get the response content from the message
    // choice.Message.Content contains the message text

    // Access citations if available
    if context := choice.Message.Context; context != nil {
        for _, citation := range context.Citations {
            // Process each citation
            // citation.Content contains the citation text
        }
    }
}
```

**After:**
```go
// Create Azure Search data source configuration
azureSearchDataSource := &azopenai.AzureSearchChatExtensionConfiguration{
    Parameters: &azopenai.AzureSearchChatExtensionParameters{
        Endpoint:       &search_endpoint,
        IndexName:      &search_index,
        Authentication: &azopenai.OnYourDataSystemAssignedManagedIdentityAuthenticationOptions{},
    },
}

// Use the standard OpenAI client with Azure data source extension
resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
    Model: openai.ChatModel("my-deployment"), // Azure deployment name goes here
    Messages: []openai.ChatCompletionMessageParamUnion{
        {
            OfUser: &openai.ChatCompletionUserMessageParam{
                Content: openai.ChatCompletionUserMessageParamContentUnion{
                    OfString: openai.String("Your message here"),
                },
            },
        },
    },
}, azopenai.WithDataSource(azureSearchDataSource))

// Access citations from the response
for _, choice := range resp.Choices {
    // Get the response content from the message
    // choice.Message.Content contains the message text

    // Access citations using helper method from azopenai
    azureChatCompletionMessage := azopenai.ChatCompletionMessage(choice.Message)
    context, err := azureChatCompletionMessage.Context()
    if err == nil {
        for _, citation := range context.Citations {
            if citation.Content != nil {
                // Process each citation
                // citation.Content contains the citation text
            }
        }
    }
}
```

### Embeddings

**Before:**
```go
resp, err := client.GetEmbeddings(context.TODO(), azopenai.EmbeddingsOptions{
    // DeploymentName: to.Ptr("text-embedding-3-large"), // This only applies for the OpenAI service.
    Input: []string{"Text to embed here"},
}, nil)
if err != nil {
    // Handle error
}
for _, embedding := range resp.Data {
    // Use the embedding vector here
    // embedding.Embedding contains the vector data
}
```

**After:**
```go
resp, err := client.Embeddings.New(context.TODO(), openai.EmbeddingNewParams{
    Model: openai.EmbeddingModel("my-deployment"), // Azure deployment name here
    // Model: openai.EmbeddingModelTextEmbedding3Large, // For OpenAI, use the model name
    Input: openai.EmbeddingNewParamsInputUnion{
        OfString: openai.String("Text to embed here"),
    },
})

if err != nil {
    // Handle error
}

for _, embedding := range resp.Data {
    // Use the embedding vector here
    // embedding.Embedding contains the vector data
}
```

## Additional Resources

- [OpenAI Go Client Documentation](https://github.com/openai/openai-go)
- [Azure OpenAI Service Documentation](https://learn.microsoft.com/azure/ai-services/openai/)
