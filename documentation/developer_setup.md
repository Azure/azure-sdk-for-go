# Developer Set Up

* [Installing Go](#installing-go)
* [Create a Client](#create-a-client)
	* [Documenting Code](#documenting-code)
	* [Constructors](#constructors)
	* [Defining Methods](#defining-methods)
* [Write Tests](#write-tests)

## Installing Go

The Azure-sdk-for-go team supports Go versions latest and latest-1, to see the exact versions we support you can check the pipeline defintions [here][pipeline_definitions]. The CI pipelines test the latest and latest-1 versions on both Windows and Linux virtual machines. If you do not already have Go installed, refer to this [workspace setup][workspace_setup] article for a more in depth tutorial on setting up your Go environment (there is also an MSI if you are developing on Windows at the [go download page][go_download]). After installing Go and configuring your workspace, fork the `azure-sdk-for-go` repository and clone it to a directory that looks like: `<GO HOME>/src/github.com/Azure/azure-sdk-for-go`.


## Create a Client

After you have the generated code from Autorest, the next step is to wrap this generated code in a "convenience layer" that the customers will use directly to interact with the service. Go is not an object-oriented language like C#, Java, or Python. There is no type hierarchy in Go. Clients and models will be defined as `struct`s and methods will be defined on these structs to interact with the service.

In other languages, types can be specifically marked "public" or "private", in Go exported types and methods are defined by starting with a capital letter. The methods on structs also follow this rule, if it is for use outside of the model it must start with a capital letter.


### Documenting Code

Code is documented directly in line and can be created directly using the `doc` tool which is part of the Go toolchain. To document a type, variable, constant, function, or package write a regular comment directly preceding its declaration (with no intervening blank line). For an example, here is the documentation for the `fmt.Fprintf` function:
```golang
// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
```

Each package needs to include a `doc.go` file and not be a part of a service version. For more details about this file there is a detailed write-up in the [repo wiki][doc_go_template]. In the `doc.go` file you should include a short service overview, basic examples, and (if they exist) a link to samples in the [`azure-sdk-for-go-samples` repository][go_azsdk_samples]


### Constructors

All clients should be able to be initialized directly from the user and should begin with `New`. For example to define a constructor for a new client for the Tables service we start with defining the struct `ServiceClient`:
```golang
// A ServiceClient represents a client to the table service. It can be used to query the available tables, add/remove tables, and various other service level operations.
type ServiceClient struct {
	client  *tableClient
	service *serviceClient
	cred    SharedKeyCredential
}

```
Note that there are no exported fields on the `ServiceClient` struct, and as a rule of thumb, generated clients and credentials should be private.

Constructors for clients are separate methods that are not associated with the struct. The constructor for the ServiceClient is as follow:
```golang
// NewServiceClient creates a ServiceClient struct using the specified serviceURL, credential, and options.
func NewServiceClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (ServiceClient, error) {
	conOptions := getConnectionOptions(serviceURL, options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, runtime.NewBearerTokenPolicy(cred, []string{"https://storage.azure.com/.default"}, nil))
	con := generated.NewConnection(serviceURL, conOptions)
	return ServiceClient{
		client:  generated.NewTableClient(con, generated.Enum0TwoThousandNineteen0202),
		service: generated.NewServiceClient(con, generated.Enum0TwoThousandNineteen0202),
		con:     con,
	}, nil
}
```
In `Go`, the method parameters are enclosed with parenthesis immediately following the method name with the parameter name preceding the parameter type. The return arguments follow the parameters. If a method has more than one return parameter the types of the parameter must be enclosed in parenthesis. Note the `*` before a type indicates a pointer to that type. All methods that create a new client or interact with the service should return an `error` type as the last argument.

This client takes three parameters, the first is the service URL for the specific account. The second is an [`interface`][go_interfaces] which is a specific struct that has definitions for a certain set of methods. In the case of `azcore.TokenCredential` the `GetToken(context.Context, options policy.TokenRequestOptions)` method must be defined to be a valid interface. The final argument to methods that create clients or interact with the service should be a pointer to an `Options` parameter. This options struct should have `azcore.ClientOptions` embedded and any service specific options. Making this final parameter a pointer allows the customer to pass in `nil` if there are no specific options they want to change.


### Defining Methods

Defining a method follows the format:
```golang
// Create creates the table with the tableName specified when NewClient was called.
func (t *Client) Create(ctx context.Context, options *CreateTableOptions) (CreateTableResponse, error) {
	if options == nil {
		options = &CreateTableOptions{}
	}
	resp, err := t.client.Create(ctx, generated.Enum1Three0, generated.TableProperties{TableName: &t.name}, options.toGenerated(), &generated.QueryOptions{})
	return createTableResponseFromGen(&resp), err
}

```
The `(s *Client)` portion is the "receiver". Methods can be defined for either pointer (with a `*`) or receiver (without a `*`) types. Pointer receivers will not copy types on method calls and allows the method to mutate the receiving struct. Client methods should use a pointer receiver.

All methods that perform I/O of any kind, sleep, or perform a significant amount of CPU-bound work must have the first parameter be of type [`context.Context`][golang_context] which allows the customer to carry a deadline, cancellation signal, and other values across API boundaries. The remaining parameters should be parameters specific to that method. The return types for methods should be first a "Response" object and second an `error` object.


## Write Tests

Testing is built into the Go toolchain as well with the `testing` library. The testing infrastructure located in the `sdk/internal/recording` directory takes care of generating recordings, establishing the mode a test is being run in (options are "record" or "playback") and reading environment variables. The HTTP traffic is intercepted by a custom [test-proxy][test_proxy_docs] in both the "recording" and "playback" case to either persist or read HTTP interactions from a file. There is one small step that needs to be added to you client creation to route traffic to this test proxy. All three of these modes are specified in the `AZURE_RECORD_MODE` environment variable:

| Mode | Powershell Command | Usage |
| ---- | ------------------ | ----- |
| record | `$ENV:AZURE_RECORD_MODE="record"` | Running against a live service and recording HTTP interactions |
| playback | `$ENV:AZURE_RECORD_MODE="playback"` | Running tests against recording HTTP interactiosn |
| live | `$ENV:AZURE_RECORD_MODE="live"` | Bypassing test proxy, running against live service, and not recording HTTP interactions (used by live pipelines) |

To get started first install [`docker`][get_docker]. Then to start the proxy, from the root of the repository, run the command `./eng/common/testproxy/docker-start-proxy.ps1 start`. This command will take care of pulling the pinned docker image and running it in the background.

It is not required to run the test-proxy from within the docker container, but this is how the proxy is run in the Azure DevOps pipelines. If you would like to run the test-proxy in a different manner the [documentation][test_proxy_docs] has more information.


### Test Mode Options

There are three options for test modes: "recording", "playback", and "live, each with their own purpose.

Recording mode is for testing against a live service and 'recording' the HTTP interactions in a JSON file for use later. This is helpful for developers because not every request will have to run through the service and makes your tests run much quicker. This also allows us to run our tests in public pipelines without fear of leaking secrets to our developer subscriptions.

In playback mode the JSON file that the HTTP interactions are saved to is used in place of a real HTTP call. This is quicker and is used most often for quickly verifying you did not change the behavior of your library.

Live mode is used by the internal pipelines to test directly against a service (similar to how a customer would do so). This mode bypasses any interactions with the test proxy.

### Routing Requests to the Proxy

All clients contain an options struct as the last parameter of the constructor function. In this options struct you need to have a way to provide a custom HTTP transport object. In your tests, you will replace the default HTTP transport object with a custom one in the `internal/recording` library that takes care of routing requests. Here is an example:

```golang
package aztables

import (
	...

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

var pathToPackage = "sdk/data/aztables/testdata"

func createClientForRecording(t *testing.T, tableName string, serviceURL string, cred SharedKeyCredential) (*Client, error) {
	transport, err := recording.NewRecordingHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:       client,
		},
	}

	// Validate the URL ends with a "/"
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithSharedKey(serviceURL, &cred, options)
}
```

Including this in a file for test helper methods will ensure that before each test the developer simply has to add
```golang
func TestExample(t *testing.T) {
	err := recording.Start(t, "path/to/package", nil)
	defer recording.Stop(t, nil)

	client, err := createClientForRecording(t, "myTableName", "myServiceUrl", myCredential)
	require.NoError(t, err)
	...
	<test code>
}
```
The first two methods (`Start` and `Stop`) tell the proxy when an individual test is starting and stopping to communicate when to start recording HTTP interactions and when to persist it to disk. `Start` takes three parameters, the `t *testing.T` parameter of the test, the path to where the recordings live for a package (this should be the path to the package), and an optional options struct. `Stop` just takes the `t *testing.T` and an options struct as parameters.


### Writing Tests

A simple test for `aztables` is shown below:
```golang

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
	err := recording.Start(t, pathToPackage, nil)
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

```golang
func TestMain(m *testing.M) {
	// Initialize
	if recording.GetRecordMode() == "record" {
		// start all tests with a proxy using it's defaults.
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}

		vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
		err = recording.AddURISanitizer(fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}
	}

	exitVal := m.Run()

	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		// reset the proxy to it's defaults
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}
	}

```


Note that removing the names of accounts and other values in your recording can have side effects when running your tests in playback. To take care of this, there are additional methods in the `internal/recording` module for reading environment variables and defaulting to the processed recording value. For example, an `aztables` test for the client constructor and "requiring" the account name to be the same as provided could look like this:

```golang
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

```golang
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


## Create Pipelines

When you create the first PR for your library you will want to create this PR against a `track2-<package>` library. Submitting PRs to the `main` branch should only be done once your package is close to being released. Treating `track2-<package>` as your main development branch will allow nightly CI and live pipeline runs to pick up issues as soon as they are introduced. After creating this PR add a comment with the following:
```
/azp run prepare-pipelines
```
This creates the pipelines that will verify future PRs. The `azure-sdk-for-go` is tested against latest and latest-1 on Windows and Linux. All of your future PRs (regardless of whether they are made to `track2-<package>` or another branch) will be tested against these versions. For more information about the individual checks run by CI and troubleshooting common issues check out the `eng_sys.md` file.


<!-- LINKS -->
[doc_go_template]: https://github.com/Azure/azure-sdk-for-go/wiki/doc.go-template
[get_docker]: https://docs.docker.com/get-docker/
[go_azsdk_samples]: https://github.com/azure-samples/azure-sdk-for-go-samples
[go_download]: https://golang.org/dl/
[go_interfaces]: https://gobyexample.com/interfaces
[pipeline_definitions]: https://github.com/Azure/azure-sdk-for-go/blob/main/eng/pipelines/templates/jobs/archetype-sdk-client.yml
[require_package]: https://pkg.go.dev/github.com/stretchr/testify/require
[test_proxy_docs]: https://github.com/Azure/azure-sdk-tools/tree/main/tools/test-proxy
[workspace_setup]: https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-windows-10