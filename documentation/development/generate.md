# Generate a SDK

* [Prerequisites](#prerequisites)
* [Generate a SDK](#generate-a-sdk)
    * [Using TypeSpec and tsp-client](#using-typespec-and-tsp-client)
    * [Using Autorest **DEPRECATED**](#using-autorest-deprecated)
* [Testing](#testing)


## Prerequisites

Follow the environment set up instructions in [setup.md]. If generating a brand-new module, follow the prerequisite instructions for new SDKs, also in [setup.md].

If you are re-generating an existing module, fork the `azure-sdk-for-go` repository and clone it to a directory that looks like: `<prefix-path>/Azure/azure-sdk-for-go`.
We use the `OneFlow` branching/workflow strategy with some minor variations.  See [repo branching][repo_branching] for further info.

Please refer to the Azure [Go SDK API design guidelines][api_design] for detailed information on how to structure clients, their APIs, and more.

## Generate a SDK 

### Using TypeSpec and tsp-client

Once the prerequisites are complete, you can generate using the tsp-client tool. tsp-client bundles up the client generation process into a single program, making it simple to generate your code once you've onboarded your TypeSpec files into the azure-rest-api-specs repository.

Setting up your project involves a few steps:

1. Add the go configuration to your TypeSpec project in the `tspconfig.yaml` file in azure-rest-api-specs: ([example](https://github.com/Azure/azure-rest-api-specs/blob/bd235f0c4ef6b3887dae6658a0a3a766a6fa4887/specification/eventgrid/Azure.Messaging.EventGrid/tspconfig.yaml#L57)).

    ```yaml
    # other YAML elided
    options:
      # other emitters elided
      "@azure-tools/typespec-go":
        module: "github.com/Azure/azure-sdk-for-go/{service-dir}/aznamespaces"
        service-dir: "sdk/messaging/eventgrid"
        emitter-output-dir: "{output-dir}/{service-dir}/aznamespaces"
    ```
2. Create a tsp-location.yaml file at the root of your module directory. This file gives the location and commit that should be used to generate your code: ([example](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/eventgrid/aznamespaces/tsp-location.yaml)).
    ``` yaml
    directory: specification/eventgrid/Azure.Messaging.EventGrid
    commit: 8d6deb81acb126a071f6f7dbf18d87a49a82e7e2
    repo: Azure/azure-rest-api-specs
    ```

3. Install [`tsp-client`](https://github.com/Azure/azure-sdk-tools/blob/main/tools/tsp-client/README.md). If already installed, be sure to update to the latest version.
4. In a terminal, `cd` into your package folder and type `tsp-client update`. This should run without error and will create a client, along with code needed to serialize and deserialize models.

    ```shell
    azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces$ tsp-client update
    ```

    To generate using a local TypeSpec project,
    ```shell
    azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces$ tsp-client update --local-spec-repo <path to TypeSpec project>
    ```

Generated code **must not** be edited, as any edits would be lost on future regeneration of content. To make customizations, update the TypeSpec project's `client.tsp` file, then regenerate.

### Using Autorest **DEPRECATED**

If your SDK doesn't require any Autorest-generated content, please skip this section. All new SDKs should be created using TypeSpec.

When using [Autorest][autorest_intro] to generate code, it's best to create a configuration file that contains all of the parameters.
This ensures that the build is repeatable and any changes are documented.
The convention is to place the parameters in a file named `autorest.md`.
Below is a template to get you started (you **must** include the yaml delimiters).

```yaml
clear-output-folder: false
export-clients: true
go: true
input-file: <URI to OpenAPI spec file>
license-header: MICROSOFT_MIT_NO_VERSION
module: <full module name> (e.g. github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys)
openapi-type: "data-plane"
output-folder: <output directory>
use: "@autorest/go@4.0.0-preview.44"
```

For the `use` section, the value should always be the latest version of the `@autorest/go` package.
The latest version can be found at the NPM [page][autorest_go] for `@autorest/go`.

For services that authenticate with Microsoft Entra ID, you **must** include the `security-scopes` parameter with the appropriate values (example below).

```yaml
security-scopes: "https://vault.azure.net/.default"
```

Generated code **must not** be edited, as any edits would be lost on future regeneration of content.
That said, if there is a need to customize the generated code, you can add one or more [Autorest directives][autorest_directives] to your autorest.md file.
This way, the changes are documented and preserved across regenerations.

## Testing

Once your SDK code is generated, you'll need to write tests and examples for the module. See [testing.md] for more information.

<!-- LINKS -->
[api_design]: https://azure.github.io/azure-sdk/golang_introduction.html#azure-sdk-module-design
[autorest_directives]: https://github.com/Azure/autorest/blob/main/docs/generate/directives.md
[autorest_go]: https://www.npmjs.com/package/@autorest/go
[autorest_intro]: https://github.com/Azure/autorest/blob/main/docs/readme.md
[repo_branching]: https://github.com/Azure/azure-sdk/blob/main/docs/policies/repobranching.md
