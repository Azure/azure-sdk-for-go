# Testing SDKs

* [Write Tests](#write-tests)
  * [Test Mode Options](#test-mode-options)
  * [Routing Request to the Proxy](#routing-requests-to-the-proxy)
  * [Writing Tests](#writing-tests)
    * [Scrubbing Secrets](#scrubbing-secrets)
  * [Live Test Resource Management](#live-test-resource-management)
  * [Committing/Updating Recordings](#committingupdating-recordings)
* [Write Examples](#write-examples)

## Write Tests

Testing is built into the Go toolchain as well with the `testing` library. The testing infrastructure located in the `sdk/internal/recording` directory takes care of generating recordings, establishing the mode a test is being run in (options are "record" or "playback") and reading environment variables. The HTTP traffic is intercepted by a custom [test-proxy][test_proxy_docs] in both the "recording" and "playback" case to either persist or read HTTP interactions from a file. There is one small step that needs to be added to you client creation to route traffic to this test proxy. All three of these modes are specified in the `AZURE_RECORD_MODE` environment variable:

| Mode | Powershell Command | Usage |
| ---- | ------------------ | ----- |
| record | `$ENV:AZURE_RECORD_MODE="record"` | Running against a live service and recording HTTP interactions |
| playback | `$ENV:AZURE_RECORD_MODE="playback"` | Running tests against recording HTTP interactiosn |
| live | `$ENV:AZURE_RECORD_MODE="live"` | Bypassing test proxy, running against live service, and not recording HTTP interactions (used by live pipelines) |

By default the recording package will automatically install and run the test proxy server. If there are issues with auto-install or the proxy needs to be run standalone, it can be run manually instead. To get started first [install test-proxy][test_proxy_install] via the standalone executable, then to start the proxy, from the root of the repository, run the command `test-proxy start`. When invoking tests, set the environment variable `PROXY_MANUAL_START` to `true`.

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

#### Example: Data Plane

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

#### Example: Management Plane

A simple test for `armchaos` is shown below:
##### The first step is to download prepared scripts to generated assets.json in the path and create file utils_test.go

1. Run the following PowerShell commands to download necessary scripts:

 ```powershell
 Invoke-WebRequest -OutFile "generate-assets-json.ps1" https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/generate-assets-json.ps1
 Invoke-WebRequest -OutFile "common-asset-functions.ps1" https://raw.githubusercontent.com/Azure/azure-sdk-tools/main/eng/common/testproxy/onboarding/common-asset-functions.ps1
```

2. Run the script in the service path:
`.\generate-assets-json.ps1 -InitialPush`
This will create a config file `assets.json` and push recordings to the Azure SDK Assets repo.

`assets.json`
```json
{
  "AssetsRepo": "Azure/azure-sdk-assets",
  "AssetsRepoPrefixPath": "go",
  "TagPrefix": "go/resourcemanager/chaos/armchaos",
  "Tag": "go/resourcemanager/chaos/armchaos_fd50b88100"
}
```

3. Before testing, create a utils_test.go file as the entry point for live tests. Modify "package" and pathToPackage to match your service.

`utils_test.go`
```go
//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armchaos_test

import (
    "os"
    "testing"
    
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
)

const (
    pathToPackage = "sdk/resourcemanager/chaos/armchaos/testdata"
)

func TestMain(m *testing.M) {
    code := run(m)
    os.Exit(code)
}

func run(m *testing.M) int {
    f := testutil.StartProxy(pathToPackage)
    defer f()
    return m.Run()
}

```

##### Then you can add the test file

`operation_live_test.go`

```go
//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armchaos_test

import (
    "context"
    "testing"
    
    "github.com/Azure/azure-sdk-for-go/sdk/azcore"
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
    "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/chaos/armchaos/v2"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
    "github.com/stretchr/testify/suite"
)

type OperationsTestSuite struct {
    suite.Suite
    
    ctx               context.Context
    cred              azcore.TokenCredential
    options           *arm.ClientOptions
    armEndpoint       string
    location          string
    resourceGroupName string
    subscriptionId    string
}

func (testsuite *OperationsTestSuite) SetupSuite() {
    testutil.StartRecording(testsuite.T(), pathToPackage)
    
    testsuite.ctx = context.Background()
    testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
    testsuite.armEndpoint = "https://management.azure.com"
    testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
    testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
    testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
    resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
    testsuite.Require().NoError(err)
    testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *OperationsTestSuite) TearDownSuite() {
    _, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
    testsuite.Require().NoError(err)
    testutil.StopRecording(testsuite.T())
}

func TestOperationsTestSuite(t *testing.T) {
    suite.Run(t, new(OperationsTestSuite))
}

// Microsoft.Chaos/operations
func (testsuite *OperationsTestSuite) TestOperation() {
    var err error
    // From step Operations_ListAll
    operationsClient, err := armchaos.NewOperationStatusesClient(testsuite.subscriptionId, testsuite.cred, nil)
    testsuite.Require().NoError(err)
    _, err = operationsClient.Get(testsuite.ctx, testsuite.location, testsuite.subscriptionId, nil)
    testsuite.Require().NoError(err)
}

```

##### At last, you can run test via command `go test -run TestOperationsTestSuite`

1. Set the test mode to "live" using:
`$ENV:AZURE_RECORD_MODE="live"`
2. Once tests pass, switch to "playback" mode and ensure all tests pass in both modes.
3. Push the final assets with `test-proxy push --assets-json-path assets.json`.

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

<!-- LINKS -->
[test_proxy_docs]: https://github.com/Azure/azure-sdk-tools/tree/main/tools/test-proxy
[test_proxy_install]: https://github.com/Azure/azure-sdk-tools/blob/main/tools/test-proxy/Azure.Sdk.Tools.TestProxy/README.md#installation
[testable_examples]: https://go.dev/blog/examples
[require_package]: https://pkg.go.dev/github.com/stretchr/testify/require
[test_resources]: https://github.com/Azure/azure-sdk-tools/tree/main/eng/common/TestResources
