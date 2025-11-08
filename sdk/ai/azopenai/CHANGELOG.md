# Release History

## 0.9.0 (2025-11-10)

### Features Added

- Updating to /v3 of the OpenAI SDK (github.com/openai/openai-go/v3).

### Other Changes

- Added examples demonstrating support for Managed Identity.
- Added examples demonstrating support for deepseek-r1 reasoning.
- Migrated examples to using the openai/v1 endpoint.

## 0.8.0 (2025-06-03)

### Breaking Changes

This library has been updated to function as a companion to the [official OpenAI Go client library](https://github.com/openai/openai-go). It provides types and functions that allow interaction with Azure-specific extensions available in the Azure OpenAI service.

See the [migration guide](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/ai/azopenai/MIGRATION.md) for details on how to update your code to use this library alongside the official OpenAI Go client.

### Other Changes

- Updating to `v1.2.1` of the [OpenAI go module](https://github.com/openai/openai-go).
- Added samples for responses API.

## 0.7.2 (2025-02-05)

### Features Added

- Updating to support Azure OpenAI API version `2025-01-01-preview`.
- Updated `ChatCompletionsOptions` and `ChatCompletionsStreamOptions`:
  - Added `Audio` parameter.
  - Added `Metadata` parameter.
  - Added `Modalities` parameter.
  - Added `Prediction` parameter.
  - Added `ReasoningEffort` parameter.
  - Added `Store` parameter.
  - Added `UserSecurityContext` parameter.
- Added `Audio` field to `ChatResponseMessage`
- Added `AudioOutputParameters` type.
- Added `AudioResponseData` type.
- Updated `CompletionsUsageCompletionTokensDetails`:
  - Added `AcceptedPredictionTokens` field.
  - Added `AudioTokens` field.
  - Added `RejectedPredictionTokens` field.
- Updated `CompletionsUsagePromptTokensDetails`:
  - Added `AudioTokens` field.
- Added `InputAudioContent` type.
- Added `ChatRequestDeveloperMessage` type.
- Added `PredictionContent` type.
- Added `UserSecurityContext` type.
- Added `ChatMessageAudioContentItem` type.
- Added `ChatCompletionModality` enum.
- Added `ChatRoleDeveloper` to the `ChatRole` enum.
- Added `InputAudioFormat` enum.
- Added `OutputAudioFormat` enum.
- Added `ReasoningEffortValue` enum.

## 0.7.1 (2024-11-13)

### Features Added

- `StreamOptions` parameter added to `ChatCompletionsOptions` and `CompletionsOptions`.
- `MaxCompletionTokens` parameter added to `ChatCompletionsOptions`.
- `ParallelToolCalls` parameter added to `ChatCompletionsOptions`.

### Breaking Changes

- `MongoDBChatExtensionParameters.Authentication`'s type has been changed to a `OnYourDataUsernameAndPasswordAuthenticationOptions`. (PR#23620)
- `GetCompletions` and `GetCompletionsStream` now receive different options (`CompletionsOptions` and `CompletionsStreamOptions` respectively)
- `GetChatCompletions` and `GetChatCompletionsStream` now receive different options (`ChatCompletionsOptions` and `ChatCompletionsStreamOptions` respectively)

## 0.7.0 (2024-10-14)

### Features Added

- MongoDBChatExtensionConfiguration has been added as an "On Your Data" data source.
- Several types now have union types for their content or dependency information:
  - ChatRequestAssistantMessage.Content is now a ChatRequestAssistantMessageContent.
  - ChatRequestSystemMessage.Content is now a ChatRequestSystemMessageContent.
  - ChatRequestToolMessage.Content is now a ChatRequestToolMessageContent.
  - MongoDBChatExtensionParameters.EmbeddingDependency is now a MongoDBChatExtensionParametersEmbeddingDependency

### Breaking Changes

- FunctionDefinition has been renamed to ChatCompletionsFunctionToolDefinitionFunction.
- AzureCosmosDBChatExtensionParameters.RoleInformation has been removed.
- AzureMachineLearningIndexChatExtension and related types have been removed.
- Several types now have union types for their content or dependency information:
  - ChatRequestAssistantMessage.Content is now a ChatRequestAssistantMessageContent.
  - ChatRequestSystemMessage.Content is now a ChatRequestSystemMessageContent.
  - ChatRequestToolMessage.Content is now a ChatRequestToolMessageContent.

## 0.6.2 (2024-09-10)

### Features Added

- Added Batch and File APIs.

### Breaking Changes

- FunctionDefinition.Parameters has been changed to take JSON instead of an object/map. You can set it using code
  similar to this:

  ```go
    parametersJSON, err := json.Marshal(map[string]any{
      "required": []string{"location"},
      "type":     "object",
      "properties": map[string]any{
        "location": map[string]any{
          "type":        "string",
          "description": "The city and state, e.g. San Francisco, CA",
        },
      },
    })

    if err != nil {
      // TODO: Update the following line with your application specific error handling logic
      log.Printf("ERROR: %s", err)
      return
    }

    // and then, in ChatCompletionsOptions
    opts := azopenai.ChatCompletionsOptions{
      Functions: []azopenai.FunctionDefinition{
        {
          Name:        to.Ptr("get_current_weather"),
          Description: to.Ptr("Get the current weather in a given location"),
          Parameters: parametersJSON,
        },
      },
    }
  ```

## 0.6.1 (2024-08-14)

### Bugs Fixed

- Client now respects the `InsecureAllowCredentialWithHTTP` flag for allowing non-HTTPS connections. Thank you @ukrocks007! (PR#23188)

## 0.6.0 (2024-06-11)

### Features Added

- Updating to the `2024-05-01-preview` API version for Azure OpenAI. (PR#22967)

### Breaking Changes

- ContentFilterResultDetailsForPrompt.CustomBlocklists has been changed from a []ContentFilterBlocklistIDResult to a struct,
  containing the slice of []ContentFilterBlocklistIDResult.
- OnYourDataEndpointVectorizationSource.Authentication's type has changed to OnYourDataVectorSearchAuthenticationOptionsClassification
- Casing has been corrected for fields:
  - Filepath -> FilePath
  - FilepathField -> FilePathField
  - CustomBlocklists -> CustomBlockLists

### Bugs Fixed

- EventReader can now handle chunks of text larger than 64k. Thank you @ChrisTrenkamp for finding the issue and suggesting a fix. (PR#22703)

## 0.5.1 (2024-04-02)

### Features Added

- Updating to the `2024-03-01-preview` API version. This adds support for using Dimensions with Embeddings as well as the ability to choose the embeddings format.
  This update also adds in the `Model` field for ChatCompletions responses. PR(#22603)

## 0.5.0 (2024-03-05)

### Features Added

- Updating to the `2024-02-15-preview` API version.
- `GetAudioSpeech` enables translating text to speech.

### Breaking Changes

- Citations, previously returned as an unparsed JSON blob, are now deserialized into a real type in `ChatResponseMessage.Citations`.
- `AzureCognitiveSearchChatExtensionConfiguration` has been renamed to `AzureSearchChatExtensionConfiguration`.
- `AzureCognitiveSearchChatExtensionParameters` has been renamed to `AzureSearchChatExtensionParameters`.

## 0.4.1 (2024-01-16)

### Bugs Fixed

- `AudioTranscriptionOptions.Filename` and `AudioTranslationOptions.Filename` fields are now properly propagated, allowing
  for disambiguating the format of an audio file when OpenAI can't detect it. (PR#22210)

## 0.4.0 (2023-12-11)

Support for many of the features mentioned in OpenAI's November Dev Day and Microsoft's 2023 Ignite conference

### Features Added

- Chat completions has been extended to accomodate new features:
  - Parallel function calling via Tools. See the function `ExampleClient_GetChatCompletions_functions` in `example_client_getchatcompletions_extensions_test.go` for an example of specifying a Tool.
  - "JSON mode", via `ChatCompletionOptions.ResponseFormat` for guaranteed function outputs.
- ChatCompletions can now be used with both text and images using `gpt-4-vision-preview`.
  - Azure enhancements to `gpt-4-vision-preview` results that include grounding and OCR features
- GetImageGenerations now works with DallE-3.
- `-1106` model feature support for `gpt-35-turbo` and `gpt-4-turbo`, including use of a seed via `ChatCompletionsOptions.Seed` and system fingerprints returned in `ChatCompletions.SystemFingerprint`.
- `dall-e-3` image generation capabilities via `GetImageGenerations`, featuring higher model quality, automatic prompt revisions by `gpt-4`, and customizable quality/style settings

### Breaking Changes

- `azopenai.KeyCredential` has been replaced by [azcore.KeyCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore#KeyCredential).
- `Deployment` has been renamed to `DeploymentName` throughout all APIs.
- `CreateImage` has been replaced with `GetImageGenerations`.
- `ChatMessage` has been split into per-role types. The function `ExampleClient_GetChatCompletions` in `example_client_getcompletions_test.go` shows an example of this.

## 0.3.0 (2023-09-26)

### Features Added

- Support for Whisper audio APIs for transcription and translation using `GetAudioTranscription` and `GetAudioTranslation`.

### Breaking Changes

- ChatChoiceContentFilterResults content filtering fields are now all typed as ContentFilterResult, instead of unique types for each field.
- `PromptAnnotations` renamed to `PromptFilterResults` in `ChatCompletions` and `Completions`.

## 0.2.0 (2023-08-28)

### Features Added

- ChatCompletions supports Azure OpenAI's newest feature to use Azure OpenAI with your own data. See `example_client_getchatcompletions_extensions_test.go`
  for a working example. (PR#21426)

### Breaking Changes

- ChatCompletionsOptions, CompletionsOptions, EmbeddingsOptions `DeploymentID` field renamed to `Deployment`.
- Method `Close()` on `EventReader[T]` now returns an error.

### Bugs Fixed

- EventReader, used by GetChatCompletionsStream and GetCompletionsStream for streaming results, would not return an
  error if the underlying Body reader was closed or EOF'd before the actual DONE: token arrived. This could result in an
  infinite loop for callers. (PR#21323)

## 0.1.1 (2023-07-26)

### Breaking Changes

- Moved from `sdk/cognitiveservices/azopenai` to `sdk/ai/azopenai`.

## 0.1.0 (2023-07-20)

- Initial release of the `azopenai` library
