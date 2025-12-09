# Developer Set Up

* [Install Dependencies](#install-dependencies)
	* [Configuring VSCode](#configuring-vscode)
* [New SDK Setup](#new-sdk-setup)
    * [Directory Structure](#directory-structure)
	* [Create Module Skeleton](#create-module-skeleton)
		* [Module Version Constant](#module-version-constant)
	* [Create Pipelines](#create-pipelines)

## Install Dependencies 

Install the following dependencies:
* [Go][go-install]
* [Node.js][node-install]
* [tsp-client][tsp-client-install]

The Azure SDK for Go supports the latest two versions of Go. When setting up a new development environment, we recommend installing the latest version of Go per the [Go download page][go_download].

### Configuring VSCode

If you're using VSCode, install the Go extension for VSCode. This usually happens automatically when opening a .go file for the first time.
See the [docs][vscode_go] for more information on using and configuring the extension.
After the extension is installed, you should be prompted to install the VSCode Go tools which are required for the extension to properly work.
To manually install or update the tools, open the VSCode command palette, select `Go: Install/Update Tools`, and select all boxes.

## New SDK Setup

### Directory Structure

Fork the `azure-sdk-for-go` repository and clone it to a directory that looks like: `<prefix-path>/Azure/azure-sdk-for-go`.
We use the `OneFlow` branching/workflow strategy with some minor variations.  See [repo branching][repo_branching] for further info.

After cloning the repository, create the appropriate directory structure for the module. It should look something like this.

`/sdk/<group>/<prefix><service>`

- `<group>` is the name of the technology, e.g. `messaging`.
- `<prefix>` is `az` for data-plane or `arm` for management-plane.
- `<service>` is the name of the service within the specified `<group>`, e.g. `servicebus`.

All directory structures **MUST** be approved by the Go SDK team (contact azsdkgo).

For more information, please consult the Azure Go SDK design guidelines on [directory structure][directory_structure].

### Create Module Skeleton

There are several files required to be in the root directory of your module.

- CHANGELOG.md for tracking released changes
- LICENSE.txt is the MIT license file
- README.md for getting started
- ci.yml for PR and release pipelines
- go.mod defines the Go module

These files can be copied from the [aztemplate][aztemplate] directory to jump-start the process. Be sure to update the contents as required, replacing all
occurrences of `template/aztemplate` with the correct values.

#### Module Version Constant

The release pipeline **requires** the presence of a constant named `moduleVersion` that contains the semantic version of the module.
The constant **must** be in a file named version.go or constants.go.  It does _not_ need to be in the root of the repo.

```go
const moduleVersion = "v1.2.3"
```

Or as part of a `const` block.

```go
const (
	moduleVersion = "v1.2.3"
	// other constants
)
```

### Create Pipelines

If your PR is the first for the module, you'll have to initialize the pipelines. After creating this PR add a comment with the following:

```
/azp run prepare-pipelines
```

This creates the pipelines that will verify future PRs. The `azure-sdk-for-go` is tested against latest and latest-1 on Windows and Linux. For more information about the individual checks run by CI and troubleshooting common issues check out the `eng_sys.md` file.


<!-- LINKS -->
[aztemplate]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/template/aztemplate
[directory_structure]: https://azure.github.io/azure-sdk/golang_introduction.html#azure-sdk-module-design
[go_download]: https://golang.org/dl/
[go-install]: https://go.dev/doc/install
[node-install]: https://nodejs.org/download/
[repo_branching]: https://github.com/Azure/azure-sdk/blob/main/docs/policies/repobranching.md
[tsp-client-install]: https://github.com/Azure/azure-sdk-tools/tree/main/tools/tsp-client#installation
[typespec-install]: https://typespec.io/docs/#install-typespec
[vscode_go]: https://code.visualstudio.com/docs/languages/go