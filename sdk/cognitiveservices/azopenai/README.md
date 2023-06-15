# Azure OpenAI client module for Go

Azure OpenAI is a managed service that allows developers to deploy, tune, and generate content from OpenAI models on Azure resources.

The Azure OpenAI client library for GO is an adaptation of OpenAI's REST APIs that provides an idiomatic interface and rich integration with the rest of the Azure SDK ecosystem.

[Source code][azopenai_repo] | [Package (pkg.go.dev)][azopenai_pkg_go] | [REST API documentation][openai_rest_docs] | [Product documentation][openai_docs]

## Getting started

### Prerequisites

* Go, version 1.18 or higher - [Install Go](https://go.dev/doc/install)
* [Azure subscription][azure_sub]
* [Azure OpenAI access][azure_openai_access]

### Install the packages

Install the `azopenai` and `azidentity` modules with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

The [azidentity][azure_identity] module is used for authentication during client construction.

### Authentication

<!-- TODO: Add api-key authentication instructions -->

#### Create a client

Constructing the client requires your vault's URL, which you can get from the Azure CLI or the Azure Portal.

```go
import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

func main() {
	endpoint     := "https://<TODO: OpenAI endpoint>"
	apiKey       := "<TODO: OpenAI API key>"

	var err error
	cred := azopenai.KeyCredential{APIKey: apiKey}
	client, err := azopenai.NewClientWithKeyCredential(endpoint, cred, &options)
	if err != nil {
		// TODO: handle error
	}
}
```

## Key concepts

See [Key concepts][openai_key_concepts] in the product documentation for more details about general concepts.

## Troubleshooting

### Error Handling

All methods that send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from the service.

### Logging

This module uses the logging implementation in `azcore`. To turn on logging for all Azure SDK modules, set `AZURE_SDK_GO_LOGGING` to `all`. By default, the logger writes to stderr. Use the `azcore/log` package to control log output. For example, logging only HTTP request and response events, and printing them to stdout:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

// Print log events to stdout
azlog.SetListener(func(cls azlog.Event, msg string) {
	fmt.Println(msg)
})

// Includes only requests and responses in credential logs
azlog.SetEvents(azlog.EventRequest, azlog.EventResponse)
```

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a [Contributor License Agreement (CLA)][cla] declaring that you have the right to, and actually do, grant us the rights to use your contribution.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate
the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to
do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][coc]. For more information, see
the [Code of Conduct FAQ][coc_faq] or contact [opencode@microsoft.com][coc_contact] with any additional questions or
comments.

<!-- LINKS -->
<!-- https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/cognitiveservices/azopenai ->
<!-- https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai -->
[azopenai_repo]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk
[azopenai_pkg_go]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[openai_docs]: https://learn.microsoft.com/azure/cognitive-services/openai
[openai_key_concepts]: https://learn.microsoft.com/azure/cognitive-services/openai/overview#key-concepts
[openai_rest_docs]: https://learn.microsoft.com/azure/cognitive-services/openai/reference
[cla]: https://cla.microsoft.com
[coc]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com