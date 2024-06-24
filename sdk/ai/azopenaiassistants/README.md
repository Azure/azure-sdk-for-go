# Azure OpenAI assistants client module for Go

NOTE: this client can be used with Azure OpenAI and OpenAI.

OpenAI assistants makes it simpler to have a create, manage and use Assistant, where conversation state is stored and managed by the service.  These assistants are backed by the same powerful models you're used to with OpenAI, and also allows the use of the Code Interpreter, Retrieval and Function Calling tools.

Use this module to:

- Create and manage assistants, threads, messages, and runs.
- Configure and use tools with assistants.
- Upload and manage files for use with assistants.

[Source code][azopenaiassistants_repo] | [Package (pkg.go.dev)][azopenaiassistants_pkg_go] | [REST API documentation][openai_rest_docs] | [Product documentation][openai_docs]

## Getting started

### Prerequisites

* Go, version 1.18 or higher - [Install Go](https://go.dev/doc/install)
* [Azure subscription][azure_sub]
* [Azure OpenAI access][azure_openai_access]

### Install the packages

Install the `azopenaiassistants` and `azidentity` modules with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants

# optional
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

The [azidentity][azure_identity] module is used for Azure Active Directory authentication with Azure OpenAI.

### Authentication

#### Azure OpenAI

Azure OpenAI clients can authenticate using Azure Active Directory or with an API key:

* Using Azure Active Directory, with a TokenCredential: [example][azopenaiassistants_example_tokencredential]
* Using an API key: [example][azopenaiassistants_example_keycredential]

#### OpenAI

OpenAI supports connecting using an API key: [example][azopenaiassistants_example_openai].

## Key concepts

See [Key concepts][openai_key_concepts_assistants] in the product documentation for more details about general concepts.

# Examples

Examples for various scenarios can be found on [pkg.go.dev][azopenaiassistants_examples] or in the example*_test.go files in our GitHub repo for [azopenaiassistants][azopenaiassistants_github].

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
[azure_openai_access]: https://learn.microsoft.com/azure/cognitive-services/openai/overview#how-do-i-get-access-to-azure-openai
[azopenaiassistants_repo]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/ai/azopenaiassistants
[azopenaiassistants_pkg_go]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants
[azopenaiassistants_examples]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#pkg-examples
[azopenaiassistants_example_tokencredential]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#example-NewClient
[azopenaiassistants_example_keycredential]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#example-NewClientWithKeyCredential
[azopenaiassistants_example_openai]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#example-NewClientForOpenAI
[azopenaiassistants_github]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/ai/azopenaiassistants
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[openai_docs]: https://learn.microsoft.com/azure/cognitive-services/openai
[openai_key_concepts]: https://learn.microsoft.com/azure/cognitive-services/openai/overview#key-concepts
[openai_key_concepts_assistants]: https://platform.openai.com/docs/assistants/overview
[openai_rest_docs]: https://learn.microsoft.com/azure/cognitive-services/openai/reference
[cla]: https://cla.microsoft.com
[coc]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com
[azure_openai_quickstart]: https://learn.microsoft.com/azure/cognitive-services/openai/quickstart
