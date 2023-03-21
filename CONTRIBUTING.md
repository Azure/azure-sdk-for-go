# Azure SDK for Go Contributing Guide

Thank you for your interest in contributing to Azure SDK for Go.

- For reporting bugs, requesting features, or asking for support, please file an issue in the [issues](https://github.com/Azure/azure-sdk-for-go/issues) section of the project.

- If you would like to become an active contributor to this project please follow the instructions provided in [Microsoft Azure Projects Contribution Guidelines](https://azure.github.io/azure-sdk/policies_opensource.html).

- To make code changes, or contribute something new, please follow the [GitHub Forks / Pull requests model](https://help.github.com/articles/fork-a-repo/): Fork the repo, make the change and propose it back by submitting a pull request.

## Pull Requests

- **DO** follow the API design and implementation [Go Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html).
  - When submitting large changes or features, **DO** have an issue or spec doc that describes the design, usage, and motivating scenario.
- **DO** submit all code changes via pull requests (PRs) rather than through a direct commit. PRs will be reviewed and potentially merged by the repo maintainers after a peer review that includes at least one maintainer.
- **DO** review your own PR to make sure there are no unintended changes or commits before submitting it.
- **DO NOT** submit "work in progress" PRs. A PR should only be submitted when it is considered ready for review and subsequent merging by the contributor.
  - If the change is work-in-progress or an experiment, **DO** start off as a temporary draft PR.
- **DO** give PRs short-but-descriptive names (e.g. "Improve code coverage for Azure.Core by 10%", not "Fix #1234") and add a description which explains why the change is being made.
- **DO** refer to any relevant issues, and include [keywords](https://help.github.com/articles/closing-issues-via-commit-messages/) that automatically close issues when the PR is merged.
- **DO** tag any users that should know about and/or review the change.
- **DO** ensure each commit successfully builds. The entire PR must pass all tests in the Continuous Integration (CI) system before it'll be merged.
- **DO** address PR feedback in an additional commit(s) rather than amending the existing commits, and only rebase/squash them when necessary. This makes it easier for reviewers to track changes.
- **DO** assume that ["Squash and Merge"](https://github.com/blog/2141-squash-your-commits) will be used to merge your commit unless you request otherwise in the PR.
- **DO NOT** mix independent, unrelated changes in one PR. Separate real product/test code changes from larger code formatting/dead code removal changes. Separate unrelated fixes into separate PRs, especially if they are in different modules or files that otherwise wouldn't be changed.
- **DO** comment your code focusing on "why", where necessary. Otherwise, aim to keep it self-documenting with appropriate names and style.
- **DO** add [GoDoc style comments](https://azure.github.io/azure-sdk/golang_introduction.html#documentation-style) when adding new APIs or modifying header files.
- **DO** make sure there are no typos or spelling errors, especially in user-facing documentation.
- **DO** verify if your changes have impact elsewhere. For instance, do you need to update other docs or exiting markdown files that might be impacted?
- **DO** add relevant unit tests to ensure CI will catch future regressions.

## Merging Pull Requests (for project contributors with write access)

- **DO** use ["Squash and Merge"](https://github.com/blog/2141-squash-your-commits) by default for individual contributions unless requested by the PR author.
  Do so, even if the PR contains only one commit. It creates a simpler history than "Create a Merge Commit".
  Reasons that PR authors may request "Merge and Commit" may include (but are not limited to):

  - The change is easier to understand as a series of focused commits. Each commit in the series must be buildable so as not to break `git bisect`.
  - Contributor is using an e-mail address other than the primary GitHub address and wants that preserved in the history. Contributor must be willing to squash
    the commits manually before acceptance.

## Developer Guide

### Repo structure

Most packages under the `services` directory in the SDK are generated from [Azure API specs][azure_rest_specs]
using [Azure/autorest.go][] and [Azure/autorest][]. These generated packages depend on the HTTP client implemented at [Azure/go-autorest][]. Therefore when contributing, please make sure you do not change anything under the `services` directory.

[azure_rest_specs]: https://github.com/Azure/azure-rest-api-specs
[azure/autorest]: https://github.com/Azure/autorest
[azure/autorest.go]: https://github.com/Azure/autorest.go
[azure/go-autorest]: https://github.com/Azure/go-autorest

For bugs or feature requests you can submit them using the [Github issues page][issues] and filling the appropriate template.

### Codespaces

Codespaces is new technology that allows you to use a container as your development environment. This repo provides a Codespaces container which is supported by both GitHub Codespaces and VS Code Codespaces.

#### GitHub Codespaces

1. From the Azure SDK GitHub repo, click on the "Code -> Open with Codespaces" button.
1. Open a Terminal. The development environment will be ready for you. Continue to [Building and Testing](https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#building-and-testing).

#### VS Code Codespaces

1. Install the [VS Code Remote Extension Pack](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)
1. When you open the Azure SDK for Go repo in VS Code, it will prompt you to open the project in the Dev Container. If it does not prompt you, then hit CTRL+P, and select "Remote-Containers: Open Folder in Container..."
1. Open a Terminal. The development environment will be ready for you. Continue to [Building and Testing](https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#building-and-testing).

### Building and Testing

#### Building

SDKs are either old (track 1) or new (track 2):

- Old (Track 1) SDKs are found in the services/ and profiles/ top level folders. 
    - CI is in /azure-pipelines.yml
- New (Track 2) SDKs are found in the sdk/ top level folder.
    - CI is in /eng/pipelines/templates/steps/build.yml

To build, run `go build` from the respective SDK directory.

There currently is not a repository wide way to build or regenerate code.

#### Testing

To test, run 'go test' from the respective directory.

## Samples

### Third-party dependencies

Third party libraries should only be included in samples when necessary to demonstrate usage of an Azure SDK package; they should not be suggested or endorsed as alternatives to the Azure SDK.

When code samples take dependencies, readers should be able to use the material without significant license burden or research on terms. This goal requires restricting dependencies to certain types of open source or commercial licenses.

Samples may take the following categories of dependencies:

- **Open-source** : Open source offerings that use an [Open Source Initiative (OSI) approved license](https://opensource.org/licenses). Any component whose license isn't OSI-approved is considered a commercial offering. Prefer OSS projects that are members of any of the [OSS foundations that Microsoft is part of](https://opensource.microsoft.com/ecosystem/). Prefer permissive licenses for libraries, like [MIT](https://opensource.org/licenses/MIT) and [Apache 2](https://opensource.org/licenses/Apache-2.0). Copy-left licenses like [GPL](https://opensource.org/licenses/gpl-license) are acceptable for tools, and OSs. [Kubernetes](https://github.com/kubernetes/kubernetes), [Linux](https://github.com/torvalds/linux), and [Newtonsoft.Json](https://github.com/JamesNK/Newtonsoft.Json) are examples of this license type. Links to open source components should be to where the source is hosted, including any applicable license, such as a GitHub repository (or similar).

- **Commercial**: Commercial offerings that enable readers to learn from our content without unnecessary extra costs. Typically, the offering has some form of a community edition, or a free trial sufficient for its use in content. A commercial license may be a form of dual-license, or tiered license. Links to commercial components should be to the commercial site for the software, even if the source software is hosted publicly on GitHub (or similar).

- **Dual licensed**: Commercial offerings that enable readers to choose either license based on their needs. For example, if the offering has an OSS and commercial license, readers can  choose between them. [MySql](https://github.com/mysql/mysql-server) is an example of this license type.

- **Tiered licensed**: Offerings that enable readers to use the license tier that corresponds to their characteristics. For example, tiers may be available for students, hobbyists, or companies with defined revenue  thresholds. For offerings with tiered licenses, strive to limit our use in tutorials to the features available in the lowest tier. This policy enables the widest audience for the article. [Docker](https://www.docker.com/), [IdentityServer](https://duendesoftware.com/products/identityserver), [ImageSharp](https://sixlabors.com/products/imagesharp/), and [Visual Studio](https://visualstudio.com) are examples of this license type.

In general, we prefer taking dependencies on licensed components in the order of the listed categories. In cases where the category may not be well known, we'll document the category so that readers understand the choice that they're making by using that dependency.
