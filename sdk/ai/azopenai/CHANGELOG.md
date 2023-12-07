# Release History

## 0.4.0 (2023-12-07)

### Features Added

### Breaking Changes

- `azopenai.KeyCredential` has been replaced by [azcore.KeyCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore#KeyCredential).
- `Deployment` has been renamed to `DeploymentName` throughout all APIs.
- `CreateImage` has been replaced with `GetImageGenerations`.

### Bugs Fixed

### Other Changes

- Chat completions with functions have been updated to allow for the newer tools-style. See the example function `ExampleClient_GetChatCompletions_functions` 
  in `example_client_getchatcompletions_extensions_test.go`. The legacy style using `ChatCompletionsOptions.Functions` continues to be supported for older models.

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

* Initial release of the `azopenai` library
