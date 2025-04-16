# Migration Guide from Azure OpenAI SDK v0.7.x to v0.8.0

## Overview

Starting with version 0.8.0, the `azopenai.Client` provided by this package has been retired in favor of the  [official OpenAI Go client library](https://github.com/openai/openai-go). That package contains all that is needed to connect to both the Azure OpenAI and OpenAI services. In that context, this library became a companion meant to enable Azure-specific extensions (e.g Azure OpenAI On Your Data). Similarly, the `azopenaiassistants` package has been completely deprecated in favor of the before mentioned official client.

Although it is understood that there is cost associated to migrating from using this package to using the official library, it is also acknowledged that in the  long term the benefits outweight these costs:

- Consistent API experience between Azure OpenAI and OpenAI services
- Direct access to the latest OpenAI features through the official library

This document is provided as a way to make this transition as smooth as possible.

> [!NOTE]
> This document is a work-in-progress and may change to reflect updates to the package. We value your feedback, please [create an issue](https://github.com/Azure/azure-sdk-for-go/issues/new/choose) to suggest any improvements or report any problems with this guide or with the package itself.

## Key Changes

### New Dependency

Your projects will now need to include the official OpenAI Go client:

```go
import (
    "github.com/openai/openai-go"
)
```

If you need to use any Azure extension, you will also need to include the `azopenai` package.

```go
import (
    "github.com/openai/openai-go"
    "github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
)
```

> [!NOTE]
> When we speak of Azure extensions we do not mean authentication or any other basic connection difference with the OpenAI service, but rather differences that introduce new models and modify the structure of the requests or responses (e.g Azure OpenAI On Your Data)

## Authentication and Client Creation

Instead of using the Azure OpenAI client directly for all operations, you'll now:

1. Create an OpenAI client configured for Azure
2. Use the Azure OpenAI companion library for Azure-specific extensions

### Azure OpenAI with API Key

Before:
```go
endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
key := os.Getenv("AZURE_OPENAI_API_KEY")
client, err := azopenai.NewClientWithKeyCredential(endpoint, azcore.NewKeyCredential(key), nil)
if err != nil {
    panic(err)
}
```

After:
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

Before:
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

After:
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

Before:
```go
key := os.Getenv("OPENAI_API_KEY")

client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", azcore.NewKeyCredential(key), nil)
if err != nil {
    panic(err)
}
```

After:
```go
key := os.Getenv("OPENAI_API_KEY")
client := openai.NewClient(
    option.WithAPIKey(key),
)
```

### API Changes

Please refer to the [official OpenAI Go client documentation](https://github.com/openai/openai-go) for details on the standard API operations.

For Azure-specific extensions provided by this companion library, see the reference documentation and examples.

## Common Migration Scenarios

### Chat Completions

Before:
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

After:
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

Before:
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

After:
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

Before:
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
        for _, citation := context.Citations {
            // Process each citation
            // citation.Content contains the citation text
        }
    }
}
```

After:
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
        for _, citation := context.Citations {
            if citation.Content != nil {
                // Process each citation
                // citation.Content contains the citation text
            }
        }
    }
}
```

### Embeddings

Before:
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

After:
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
