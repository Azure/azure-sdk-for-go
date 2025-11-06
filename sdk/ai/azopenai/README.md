# Azure OpenAI extensions module for Go

This module provides models and convenience functions to make it simpler to use Azure OpenAI features, such as [Azure OpenAI On Your Data][openai_on_your_data], with the OpenAI Go client (https://pkg.go.dev/github.com/openai/openai-go/v3).

[Source code][repo] | [Package (pkg.go.dev)][pkggodev] | [REST API documentation][openai_rest_docs] | [Product documentation][openai_docs]

## Getting started

### Prerequisites

- Go, version 1.23 or higher - [Install Go](https://go.dev/doc/install)
- [Azure subscription][azure_sub]
- [Azure OpenAI access][azure_openai_access]

### Install the packages

Install the `azopenai` and `azidentity` modules with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai

# optional
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```

The [azidentity][azure_identity] module is used for Azure Active Directory authentication with Azure OpenAI.

## Key concepts

See [Key concepts][openai_key_concepts] in the product documentation for more details about general concepts.

# Examples

Examples for scenarios specific to Azure can be found on [pkg.go.dev](https://aka.ms/azsdk/go/azopenaiextensions/pkg#pkg-examples) or in the example\*\_test.go files in our GitHub repo for [azopenai](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/ai/azopenai).

For examples on using the openai-go client, see the examples in the [openai-go](https://github.com/openai/openai-go/v3/tree/main/examples) repository.

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a [Contributor License Agreement (CLA)][cla] declaring that you have the right to, and actually do, grant us the rights to use your contribution.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate
the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to
do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][coc]. For more information, see
the [Code of Conduct FAQ][coc_faq] or contact [opencode@microsoft.com][coc_contact] with any additional questions or
comments.

<!-- LINKS -->

[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_openai_access]: https://learn.microsoft.com/azure/cognitive-services/openai/overview#how-do-i-get-access-to-azure-openai
[azure_openai_quickstart]: https://learn.microsoft.com/azure/cognitive-services/openai/quickstart
[azure_sub]: https://azure.microsoft.com/free/
[cla]: https://cla.microsoft.com
[coc_contact]: mailto:opencode@microsoft.com
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc]: https://opensource.microsoft.com/codeofconduct/
[openai_docs]: https://learn.microsoft.com/azure/cognitive-services/openai
[openai_key_concepts]: https://learn.microsoft.com/azure/cognitive-services/openai/overview#key-concepts
[openai_on_your_data]: https://learn.microsoft.com/azure/ai-services/openai/concepts/use-your-data
[openai_rest_docs]: https://learn.microsoft.com/azure/cognitive-services/openai/reference
[pkggodev]: https://aka.ms/azsdk/go/azopenaiextensions/pkg
[repo]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/ai/azopenai
