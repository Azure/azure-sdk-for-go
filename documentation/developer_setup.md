# Developer Set Up

* [Installing Go](#installing-go)
* [Directory Structure](#directory-structure)
* [Module Skeleton](#create-module-skeleton)
* [Create SDK](#create-your-sdk)
* [Write Tests](#write-tests)
* [Write Examples](#write-examples)

## Installing Go

The Azure SDK for Go supports the latest two versions of Go. When setting up a new development environment, we recommend installing the latest version of Go per the [Go download page][go_download].

### Configuring VSCode

If you're using VSCode, install the Go extension for VSCode. This usually happens automatically when opening a .go file for the first time.
See the [docs][vscode_go] for more information on using and configuring the extension.
After the extension is installed, you should be prompted to install the VSCode Go tools which are required for the extension to properly work.
To manually install or update the tools, open the VSCode command palette, select `Go: Install/Update Tools`, and select all boxes.

## Directory Structure

Fork the `azure-sdk-for-go` repository and clone it to a directory that looks like: `<prefix-path>/Azure/azure-sdk-for-go`.
We use the `OneFlow` branching/workflow strategy with some minor variations.  See [repo branching][repo_branching] for further info.

After cloning the repository, create the appropriate directory structure for the module. It should look something like this.

`/sdk/<group>/<prefix><service>`

- `<group>` is the name of the technology, e.g. `messaging`.
- `<prefix>` is `az` for data-plane or `arm` for management plane.
- `<service>` is the name of the service within the specified `<group>`, e.g. `servicebus`.

All directory structures **MUST** be approved by the Go SDK team (contact azsdkgo).

For more information, please consult the Azure Go SDK design guidelines on [directory structure][directory_structure].

## Create Module Skeleton

There are several files required to be in the root directory of your module.

- CHANGELOG.md for tracking released changes
- LICENSE.txt is the MIT license file
- README.md for getting started
- ci.yml for PR and release pipelines
- go.mod defines the Go module

These files can be copied from the [aztemplate][aztemplate] directory to jump-start the process. Be sure to update the contents as required, replacing all
occurrences of `template/aztemplate` with the correct values.

### Module Version Constant

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

## Create Your SDK

Once the skeleton for your SDK has been created, you can start writing your SDK.
Please refer to the Azure [Go SDK API design guidelines][api_design] for detailed information on how to structure clients, their APIs, and more.

From here you will choose which code generation stack you'll be using: [TypeSpec](https://typespec.io/), which should be used for all new stacks or autorest, which is only used for legacy stacks that use Swagger/OpenAPI specifications.

### Using TypeSpec and `tsp-client`

tsp-client bundles up the client generation process into a single program, making it simple to generate your code once you've onboarded your TypeSpec files into the azure-rest-api-specs repository.

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

Generated code **must not** be edited, as any edits would be lost on future regeneration of content. To make customizations, update the TypeSpec project's `client.tsp` file then regenerate.

### Using Autorest

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

## Write Tests

Testing is built into the Go toolchain as well with the `testing` library. The testing infrastructure located in the `sdk/internal/recording` directory takes care of generating recordings, establishing the mode a test is being run in (options are "record" or "playback") and reading environment variables. The HTTP traffic is intercepted by a custom [test-proxy][test_proxy_docs] in both the "recording" and "playback" case to either persist or read HTTP interactions from a file. There is one small step that needs to be added to you client creation to route traffic to this test proxy. All three of these modes are specified in the `AZURE_RECORD_MODE` environment variable:

| Mode | Powershell Command | Usage |
| ---- | ------------------ | ----- |
| record | `$ENV:AZURE_RECORD_MODE="record"` | Running against a live service and recording HTTP interactions |
| playback | `$ENV:AZURE_RECORD_MODE="playback"` | Running tests against recording HTTP interactiosn |
| live | `$ENV:AZURE_RECORD_MODE="live"` | Bypassing test proxy, running against live service, and not recording HTTP interactions (used by live pipelines) |

By default the [recording](recording_package) package will automatically install and run the test proxy server. If there are issues with auto-install or the proxy needs to be run standalone, it can be run manually instead. To get started first [install test-proxy][test_proxy_install] via the standalone executable, then to start the proxy, from the root of the repository, run the command `test-proxy start`. When invoking tests, set the environment variable `PROXY_MANUAL_START` to `true`.

### Test Mode Options

There are three options for test modes: `recording`, `playback`, and `live`, each with their own purpose.

Recording mode is for testing against a live service and 'recording' the HTTP interactions in a JSON file for use later. This is helpful for developers because not every request will have to run through the service and makes your tests run much quicker. This also allows us to run our tests in public pipelines without fear of leaking secrets to our developer subscriptions.

In playback mode the JSON file that the HTTP interactions are saved to is used in place of a real HTTP call. This is quicker and is used most often for quickly verifying you did not change the behavior of your library.

Live mode is used by the internal pipelines to test directly against a service (similar to how a customer would do so). This mode bypasses any interactions with the test proxy.

### Routing Requests to the Proxy

All clients contain an options struct as the last parameter of the constructor function. In this options struct you need to have a way to provide a custom HTTP transport object. In your tests, you will replace the default HTTP transport object with a custom one in the `internal/recording` library that takes care of routing requests. Here is an example:

```go
package aztables

import (
	...

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
)

func createClientForRecording(t *testing.T, tableName string, serviceURL string) (*Client, error) {
	transport, err := recording.NewRecordingHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
		},
	}

	// credential.New returns a credential for Entra ID authentication. It works in CI
	// and local development, in all recording modes. To authenticate live tests on
	// your machine, sign in to the Azure CLI or Azure Developer CLI.
	cred, err := credential.New(nil)
	require.NoError(t, err)

	return NewClient(runtime.JoinPaths(serviceURL, tableName), &cred, options)
}

func startRecording(t *testing.T) {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}
```

Including this in a file for test helper methods will ensure that before each test the developer simply has to add

```go
func TestExample(t *testing.T) {
	startRecording(t)

	client, err := createClientForRecording(t, "myTableName", "myServiceUrl")
	require.NoError(t, err)
	...
	<test code>
}
```

The first two methods (`Start` and `Stop`) tell the proxy when an individual test is starting and stopping to communicate when to start recording HTTP interactions and when to persist it to disk. `Start` takes three parameters, the `t *testing.T` parameter of the test, the path to where the recordings live for a package, and an optional options struct. `Stop` just takes the `t *testing.T` and an options struct as parameters.

NOTE: the path to the recordings **must** be in or under a directory named `testdata`; this will prevent the recordings from being included in the module disk-footprint.

### Writing Tests

A simple test for `aztables` is shown below:

```go
import (
	"fmt"
	"os"

	"github.com/stretchr/testify/require"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const accountName = os.GetEnv("TABLES_PRIMARY_ACCOUNT_NAME")

// Test creating a single table
func TestCreateTable(t *testing.T) {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	defer func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}()

	serviceURL := fmt.Sprintf("https://%v.table.core.windows.net", accountName)
	client, err := createClientForRecording(t, "tableName", serviceURL)
	require.NoError(t, err)

	resp, err := client.Create()
	require.NoError(t, err)
	require.Equal(t, resp.TableResponse.TableName, "tableName")
	defer client.Delete()  // Clean up resources
	..
	.. More test functionality
	..
}
```

The first part of the test above is for getting resource configuration from your environment.

The rest of the snippet shows a test that creates a single table and requirements (similar to assertions in other languages) that the response from the service has the same table name as the supplied parameter. Every test in Go has to have exactly one parameter, the `t *testing.T` object, and it must begin with `Test`. After making a service call or creating an object you can make assertions on that object by using the external `testify/require` library. In the example above, we "require" that the error returned is `nil`, meaning the call was successful and then we require that the response object has the same table name as supplied.

Check out the docs for more information about the methods available in the [`require`][require_package] libraries.

If you set the environment variable `AZURE_RECORD_MODE` to "record" and run `go test` with this code and the proper environment variables this test would pass and you would be left with a new directory and file. Test recordings are saved to a `recording` directory in the same directory that your test code lives. Running the above test would also create a file `recording/TestCreateTable.json` with the HTTP interactions persisted on disk. Now you can set `AZURE_RECORD_MODE` to "playback" and run `go test` again, the test will have the same output but without reaching the service.

### Scrubbing Secrets

Recording files live in the assets repository (`github.com/Azure/azure-sdk-assets`) and must not contain secrets. We use sanitizers with regular expression replacements to prevent recording secrets. The test proxy has many built-in sanitizers enabled by default. However, you may need to add your own by calling functions from the `recording` package. These functions generally take three parameters: the test instance (`t *testing.T`), the value to be removed (ie. an account name or key), and the value to use in replacement.

| Sanitizer Type | Method |
| -------------- | ------ |
| Body Key Sanitizer | `AddBodyKeySanitizer(jsonPath, value, regex string, options *RecordingOptions)` |
| Body Regex Sanitizer | `AddBodyRegexSanitizer(value, regex string, options *RecordingOptions)` |
| Continuation Sanitizer | `AddContinuationSanitizer(key, method string, resetAfterFirst bool, options *RecordingOptions)` |
| General Regex Sanitizer | `AddGeneralRegexSanitizer(value, regex string, options *RecordingOptions)` |
| Header Regex Sanitizer | `AddHeaderRegexSanitizer(key, value, regex string, options *RecordingOptions)` |
| OAuth Response Sanitizer | `AddOAuthResponseSanitizer(options *RecordingOptions)` |
| Remove Header Sanitizer | `AddRemoveHeaderSanitizer(headersForRemoval []string, options *RecordingOptions)` |
| URI Sanitizer | `AddURISanitizer(value, regex string, options *RecordingOptions)` |
| URI Subscription ID Sanitizer | `AddURISubscriptionIDSanitizer(value string, options *RecordingOptions)` |

You may also need to remove a built-in sanitizer overwriting a non-secret value needed by a test. To do this, call `RemoveRegisteredSanitizers` with a list of sanitizer IDs such as "AZSDK3430". See the [test proxy documentation](https://github.com/Azure/azure-sdk-tools/blob/main/tools/test-proxy/Azure.Sdk.Tools.TestProxy/Common/SanitizerDictionary.cs) for a complete list of all the built-in sanitizers and their IDs.

Configure sanitizers before running tests:

```go
const recordingDirectory = "<path to service directory with assets.json file>/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	// Initialize
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		// NOTE: defer should not be used directly within TestMain as it will not be executed due to os.Exit()
		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

	// sanitizers run in playback and recording modes
	if recording.GetRecordMode() != recording.LiveMode {
		// configure sanitizers here. For example:
		err := recording.RemoveRegisteredSanitizers([]string{
			"AZSDK3430", // body key $..id
		}, nil)
		if err != nil {
			panic(err)
		}
	}

	// run test cases
	return m.Run()
}
```

Note that removing the names of accounts and other values in your recording can have side effects when running your tests in playback. To take care of this, there are additional methods in the `internal/recording` module for reading environment variables and defaulting to the processed recording value. For example, an `aztables` test for the client constructor and "requiring" the account name to be the same as provided could look like this:

```go
func TestClient(t *testing.T) {
	accountName := recording.GetEnvVariable(t, "TABLES_PRIMARY_ACCOUNT_NAME", "fakeAccountName")
	// If running in playback, the value is "fakeAccountName". If running in "record" the value is the environment variable
	accountKey := recording.GetEnvVariable(t, "TABLES_PRIMARY_ACCOUNT_KEY", "fakeAccountKey")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	client, err := NewClient("someTableName", someServiceURL, cred, nil)
	require.NoError(t, err)
	require.Equal(t, accountName, client.AccountName())
}
```

### Live Test Resource Management

If you have live tests that require Azure resources, you'll need to create a test resources config file for deployment during CI.
Please see the [test resource][test_resources] documentation for more info.

### Committing/Updating Recordings

Always inspect recordings for secrets before pushing them. The `assets.json` file located in your module directory is used by the Test Framework to figure out how to retrieve session records from the assets repo. In order to push new session records, you need to invoke:

```PowerShell
test-proxy push -a <path-to-assets.json>
```

On completion of the push, a newly created tag will be stamped into the `assets.json` file. This new tag must be committed and pushed to your package directory along with any other changes.

## Write Examples

Examples are built into the Go toolchain by way of [testable examples][testable_examples]. By convention, examples are placed in a file named `example_test.go` and
may be spread across multiple files, grouped by feature (e.g. `example_<feature>_test.go`). Since testable examples are by definition tests, the file(s) must have the `_test.go` suffix.

Examples **should** be succinct allowing for copy/paste usage and **must** be clearly commented so they're easy to understand.

Examples **must** be provided as testable examples, not as markdown blocks in README files (code snippets are ok but should be used sparingly as they tend to rot over time).
This ensures that examples actually compile (and work!) and remain current as a SDK evolves. It also allows the doc tooling to automatically link API docs to their examples.

Please consult the canonical documentation on [testable examples][testable_examples] for instructions on how to create/name testable examples and enabling testable example execution.

All SDKs **must** include, at minimum, examples for their champion scenarios.

## Create Pipelines

When you create the first PR for your library you will want to create this PR against a `track2-<package>` library. Submitting PRs to the `main` branch should only be done once your package is close to being released. Treating `track2-<package>` as your main development branch will allow nightly CI and live pipeline runs to pick up issues as soon as they are introduced. After creating this PR add a comment with the following:

```
/azp run prepare-pipelines
```

This creates the pipelines that will verify future PRs. The `azure-sdk-for-go` is tested against latest and latest-1 on Windows and Linux. All of your future PRs (regardless of whether they are made to `track2-<package>` or another branch) will be tested against these versions. For more information about the individual checks run by CI and troubleshooting common issues check out the `eng_sys.md` file.


<!-- LINKS -->
[get_docker]: https://docs.docker.com/get-docker/
[go_download]: https://golang.org/dl/
[require_package]: https://pkg.go.dev/github.com/stretchr/testify/require
[test_proxy_docs]: https://github.com/Azure/azure-sdk-tools/tree/main/tools/test-proxy
[test_proxy_install]: https://github.com/Azure/azure-sdk-tools/blob/main/tools/test-proxy/Azure.Sdk.Tools.TestProxy/README.md#installation
[workspace_setup]: https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-windows-10
[directory_structure]: https://azure.github.io/azure-sdk/golang_introduction.html#azure-sdk-module-design
[module_design]: https://azure.github.io/azure-sdk/golang_introduction.html#azure-sdk-module-design
[type_declarations]: https://go.dev/ref/spec#Type_declarations
[azkeys_directory]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/security/keyvault/azkeys
[aztables_directory]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/data/aztables
[aztemplate]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/template/aztemplate
[api_design]: https://azure.github.io/azure-sdk/golang_introduction.html#azure-sdk-module-design
[vscode_go]: https://code.visualstudio.com/docs/languages/go
[repo_branching]: https://github.com/Azure/azure-sdk/blob/main/docs/policies/repobranching.md
[autorest_go]: https://www.npmjs.com/package/@autorest/go
[autorest_intro]: https://github.com/Azure/autorest/blob/main/docs/readme.md
[autorest_directives]: https://github.com/Azure/autorest/blob/main/docs/generate/directives.md
[test_resources]: https://github.com/Azure/azure-sdk-tools/tree/main/eng/common/TestResources
[recording_package]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/internal/recording
[testable_examples]: https://go.dev/blog/examples
