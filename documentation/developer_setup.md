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

If your SDK won't be generated from OpenAPI (aka swagger) files, skip to the [next section](#create-module-skeleton).

Once the directory structure has been created, you must decide if your SDK will directly export generated types (commonly referred to as a code generated client).
The alternative is to make the generated content internal and export hand-written types, possibly along with generated types via type aliasing.

### Code Generated Clients

An SDK that uses code generated clients directly exposes the Autorest-generated code to consumers of the module and is the preferred approach.
The [azkeys][azkeys_directory] module is an example of a code generated client (CGC).

Note that for data-plane CGCs, client constructors must be hand-written as there's no consistent form of authentication across data-plane services.

### Internally Generated Clients

Internally generated clients are used when the Autorest-generated code isn't fit for direct, public consumption.
In this design, the generated code is placed under an `/internal` directory, prohibiting it from being directly imported by module consumers,
and all publicly exposed content is either hand-written or a type alias of internal types (see `Alias declarations` in the [Go language specification][type_declarations] for more info on creating type aliases).

The [aztables][aztables_directory] modules is an example that uses internally generated clients.

## Create Module Skeleton

There are several files required to be in the root directory of your module.

- CHANGELOG.md for tracking released changes
- LICENSE.txt is the MIT license file
- NOTICE.txt for legal attributions
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

### Using Autorest

If your SDK doesn't require any Autorest-generated content, please skip this section.

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

For services that authenticate with Azure Active Directory, you **must** include the `security-scopes` parameter with the appropriate values (example below).

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
)

func createClientForRecording(t *testing.T, tableName string, serviceURL string, cred SharedKeyCredential) (*Client, error) {
	transport, err := recording.NewRecordingHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:       client,
		},
	}

	return NewClientWithSharedKey(runtime.JoinPaths(serviceURL, tableName), &cred, options)
}
```

Including this in a file for test helper methods will ensure that before each test the developer simply has to add

```go
func TestExample(t *testing.T) {
	err := recording.Start(t, "path/to/package", nil)
	defer recording.Stop(t, nil)

	client, err := createClientForRecording(t, "myTableName", "myServiceUrl", myCredential)
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
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const (
	accountName := os.GetEnv("TABLES_PRIMARY_ACCOUNT_NAME")
	accountKey := os.GetEnv("TABLES_PRIMARY_ACCOUNT_KEY")
)

// Test creating a single table
func TestCreateTable(t *testing.T) {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	defer func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}()

	serviceUrl := fmt.Sprintf("https://%v.table.core.windows.net", accountName)
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	client, err := createClientForRecording(t, "tableName", serviceUrl, cred)
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

The first part of the test above is for getting the secrets needed for authentication from your environment, best practice is to store your test secrets in environment variables.

The rest of the snippet shows a test that creates a single table and requirements (similar to assertions in other languages) that the response from the service has the same table name as the supplied parameter. Every test in Go has to have exactly one parameter, the `t *testing.T` object, and it must begin with `Test`. After making a service call or creating an object you can make assertions on that object by using the external `testify/require` library. In the example above, we "require" that the error returned is `nil`, meaning the call was successful and then we require that the response object has the same table name as supplied.

Check out the docs for more information about the methods available in the [`require`][require_package] libraries.

If you set the environment variable `AZURE_RECORD_MODE` to "record" and run `go test` with this code and the proper environment variables this test would pass and you would be left with a new directory and file. Test recordings are saved to a `recording` directory in the same directory that your test code lives. Running the above test would also create a file `recording/TestCreateTable.json` with the HTTP interactions persisted on disk. Now you can set `AZURE_RECORD_MODE` to "playback" and run `go test` again, the test will have the same output but without reaching the service.

### Scrubbing Secrets

The recording files eventually live in the main repository (`github.com/Azure/azure-sdk-for-go`) and we need to make sure that all of these recordings are free from secrets. To do this we use Sanitizers with regular expressions for replacements. All of the available sanitizers are available as methods from the `recording` package. The recording methods generally take three parameters: the test instance (`t *testing.T`), the value to be removed (ie. an account name or key), and the value to use in replacement.

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

To add a scrubber that replaces the URL of your account use the `TestMain()` function to set sanitizers before you begin running tests.

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

    // Set sanitizers in record mode
	if recording.GetRecordMode() == "record" {
		vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
		err = recording.AddURISanitizer(fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}
	}

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

### Using `azidentity` Credentials In Tests

The credentials in `azidentity` are not automatically configured to run in playback mode. To make sure your tests run in playback mode even with `azidentity` credentials the best practice is to use a simple `FakeCredential` type that inserts a fake Authorization header to mock a credential. An example for swapping the `DefaultAzureCredential` using a helper function is shown below in the context of `aztables`

```go
type FakeCredential struct {}

func NewFakeCredential() *FakeCredential {
	return &FakeCredential{}
}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func getAADCredential() (azcore.TokenCredential, error) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return NewFakeCredential(), nil
	}
	return azidentity.NewDefaultCredential(nil)
}

func TestClientWithAAD(t *testing.T) {
	accountName := recording.GetEnvVariable(t, "TABLES_PRIMARY_ACCOUNT_NAME", "fakeAccountName")
	cred, err := getAADCredential()
	require.NoError(t, err)
	...
	...run tests...
}
```

The `FakeCredential` show here implements the `azcore.TokenCredential` interface and can be used anywhere the `azcore.TokenCredential` is used.

### Live Test Resource Management

If you have live tests that require Azure resources, you'll need to create a test resources config file for deployment during CI.
Please see the [test resource][test_resources] documentation for more info.

### Committing/Updating Recordings

The `assets.json` file located in your module directory is used by the Test Framework to figure out how to retrieve session records from the assets repo. In order to push new session records, you need to invoke:

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
[directory_structure]: https://azure.github.io/azure-sdk/golang_introduction.html
[module_design]: https://azure.github.io/azure-sdk/golang_introduction.html#azure-sdk-module-design
[type_declarations]: https://go.dev/ref/spec#Type_declarations
[azkeys_directory]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/security/keyvault/azkeys
[aztables_directory]: https://github.com/Azure/azure-sdk-for-go/tree/sdk/data/aztables/v1.0.1/sdk/data/aztables
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
