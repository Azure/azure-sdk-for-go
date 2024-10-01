# Examples

## Conventions

If you're already using the `azopenai` package, but would like to switch to using `openai-go`, you'll need to adjust your code to accomodate the different conventions in that package.

- Fields for input types are wrapped in an `openai.Field` type, using the `openai.F()`, or helper functions like `openai.Int`:
   ```go
   	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.F(model),
		MaxTokens: openai.Int(512),
    }
   ```
- Model deployment names are passed in the `Model` input field, instead of `DeploymentName`.

## Using "Azure OpenAI On Your Data" with openai-go

["Azure OpenAI On Your Data"](https://learn.microsoft.com/azure/ai-services/openai/concepts/use-your-data) allows you to use external data sources, such as Azure AI Search, in combination with Azure OpenAI. This package provides a helper function to make it easy to include `DataSources` using `openai-go`:

For a full example see [example_azure_on_your_data_test.go](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiextensions#example-package-UsingAzureOnYourData).
