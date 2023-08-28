# Release History

## 0.2.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

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
