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

All clients should be able to be initialized directly from the user and should begin with `New`. For example to define a constructor for a new client for the Tables service we start with defining the struct `TableServiceClient`:
```golang
// A TableServiceClient represents a client to the table service. It can be used to query the available tables, add/remove tables, and various other service level operations.
type TableServiceClient struct {
	client  *tableClient
	service *serviceClient
	cred    SharedKeyCredential
}

```
Note that there are no exported fields on the `TableServiceClient` struct, and as a rule of thumb, generated clients and credentials should be private.

Constructors for clients are separate methods that are not associated with the struct. The constructor for the TableServiceClient is as follow:
```golang
// NewTableServiceClient creates a TableServiceClient struct using the specified serviceURL, credential, and options.
func NewTableServiceClient(serviceURL string, credential azcore.Credential, options *TableServiceClientOptions) (*TableServiceClient, error) {
	conOptions := options.getConnectionOptions()
	if isCosmosEndpoint(serviceURL) {
		conOptions.PerCallPolicies = []azcore.Policy{CosmosPatchTransformPolicy{}}
	}
	con := newConnection(serviceURL, cred, conOptions)
	c, err := cred.(*SharedKeyCredential)
	return &TableServiceClient{client: &tableClient{con}, service: &serviceClient{con}, cred: *c}, err
}
```
In `Go`, the method parameters are enclosed with parenthesis immediately following the method name with the parameter name preceding the parameter type. The return arguments follow the parameters. If a method has more than one return parameter the types of the parameter must be enclosed in parenthesis. Note the `*` before a type indicates a pointer to that type. All methods that create a new client or interact with the service should return an `error` type as the last argument.

This client takes three parameters, the first is the service URL for the specific account. The second is an [`interface`][go_interfaces] which is a specific struct that has definitions for a certain set of methods. In the case of `azcore.Credential` the `AuthenticationPolicy(options AuthenticationPolicyOptions) Policy` method must be defined to be a valid interface. The final argument to methods that create clients or interact with the service should be a pointer to an `Options` parameter. Making this final parameter a pointer allows the customer to pass in `nil` if there are no specific options they want to change. The `Options` type should have a name that is intuitive to what the customer is trying to do, in this case `TableClientOptions`.


### Defining Methods

Defining a method follows the format:
```golang
// Create creates a table with the specified name.
func (t *TableServiceClient) Create(ctx context.Context, name string) (TableResponseResponse, error) {
	resp, err := t.client.Create(ctx, TableProperties{&name}, new(TableCreateOptions), new(QueryOptions))
	if err == nil {
		tableResp := resp.(TableResponseResponse)
		return tableResp, nil
	}
	return TableResponseResponse{}, err
}
```
The `(t *TableServiceClient)` portion is the "receiver". Methods can be defined for either pointer (with a `*`) or receiver (without a `*`) types. Pointer receivers will not copy types on method calls and allows the method to mutate the receiving struct. It is best practice to use pointer receivers wherever possible to limit memory copies.

All methods that perform I/O of any kind, sleep, or perform a significant amount of CPU-bound work must have the first parameter be of type [`context.Context`][golang_context] which allows the customer to carry a deadline, cancellation signal, and other values across API boundaries. The remaining parameters should be parameters specific to that method. The return types for methods should be first a "Response" object and second an `error` object.


## Write Tests

Testing is built into the Go toolchain as well with the `testing` library. The testing infrastructure located in the `sdk/internal/recording` directory takes care of generating recordings, establishing the mode a test is being run in (options are "recording", "playback", "live-no-playback"), and reading environment variables. The HTTP traffic is intercepted by a custom [test-proxy][test_proxy_docs] in both the "recording" and "playback" case to either persist or read HTTP interactions from a file. There is one small step that needs to be added to you client creation to route traffic to this test proxy. All three of these modes are specified in the `AZURE_RECORD_MODE` environment variable:

| Mode | Powershell Command |
| ---- | ------------------ |
| record | $ENV:AZURE_RECORD_MODE="record" |
| playback | $ENV:AZURE_RECORD_MODE="playback" |
| live-no-playback | $ENV:AZURE_RECORD_MODE="live-no-playback" |

To get started first install [`docker`][get_docker]. Then to start the proxy, from the directory your tests live in, run the command `../path-to-root/eng/scripts proxy-server.ps1 start`. This command will take care of pulling the pinned docker image and running it in the background. Note that wherever this command is run from is where the recordings will be persisted or (if you are running in playback) searched for.

It is not required to run the test-proxy from within the docker container, but this is how the proxy is run in the Azure DevOps pipelines. If you would like to run the test-proxy in a different way check out the test-proxy [documentation][test_proxy_docs] for more information.


### Test Mode Options

There are three options for test modes: "recording", "playback", and "live-no-playback" each with their own purpose.

Recording mode is for testing against a live service and 'recording' the HTTP interactions in a JSON file for use later. This is helpful for developers because not every request will have to run through the service and makes your tests run much quicker. This also allows us to run our tests in public pipelines without fear of leaking secrets to our developer subscriptions.

In playback mode the JSON file that the HTTP interactions are saved to is used in place of a real HTTP call. This is quicker and is used most often for quickly verifying you did not change the behavior of your library.

In live-no-playback the tests run against the live service and no recordings are generated, this is most commonly done in our live pipelines that run nightly and (at a developers request) to verify PRs against a live service in addition to the normal check against recordings.


### Routing Requests to the Proxy

All clients should contain an options struct as the last parameter on the constructor. In this options struct you need to have a way to provide a custom HTTP transport object. In your tests, you will replace the default HTTP transport object with a custom one in the `internal/recording` library that takes care of all the routing for you. For example, here is that code snippet in the `aztable` package:

```golang
package aztable

import (
	...

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

type recordingPolicy struct {
	options recording.RecordingOptions
}

func NewRecordingPolicy(o *recording.RecordingOptions) azcore.Policy {
	if o == nil {
		o = &recording.RecordingOptions{}
	}
	p := &recordingPolicy{options: *o}
	p.options.Init()
	return p
}

func (p *recordingPolicy) Do(req *azcore.Request) (resp *azcore.Response, err error) {
	originalURLHost := req.URL.Host
	req.URL.Scheme = "https"
	req.URL.Host = p.options.Host
	req.Host = p.options.Host

	req.Header.Set(recording.UpstreamUriHeader, fmt.Sprintf("%v://%v", p.options.Scheme, originalURLHost))
	req.Header.Set(recording.ModeHeader, recording.GetRecordMode())
	req.Header.Set(recording.IdHeader, recording.GetRecordingId())

	return req.Next()
}

func createTableClientForRecording(t *testing.T, tableName string, serviceURL string, cred azcore.Credential) (*TableClient, error) {
	policy := NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient()
	require.NoError(t, err)
	options := &TableClientOptions{
		PerCallOptions: []azcore.Policy{policy},
		HTTPClient: client,
	}
	return NewTableClient(tableName, serviceURL, cred, options)
}
```

Including this in a file for test helper methods will ensure that before each test the developer simply has to add
```golang
func TestStartTests(t *testing.T) {
	recording.StartRecording(t, nil)
	defer recording.StopRecording(t, nil)
	client, err := createTableClientForRecording(t, "myTableName", "myServiceUrl", myCredential)
	require.NoError(t, err)
	...
	<test stuff>
}
```
and nearly all of the test proxy interactions will be handled for them. In a later section ([scrubbing secrets](#scrubbing-secrets)) there is more information about making sure the recording files are free of secrets. The first two methods (`StartRecording` and `StopRecording`) tell the proxy when an individual test is starting and stopping to communicate when to start creating the recording file and when to persist it to disk.


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

// Test creating a single table
func TestCreateTable(t *testing.T) {
	recording.StartRecording(t, nil)
	defer recording.StopRecording(t, nil)

	serviceUrl := fmt.Sprintf("https://%v.table.core.windows.net", accountName)
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	client, err := createTableClientForRecording(t, "tableName", serviceUrl, cred)
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

| Sanitizer Type | Method | Parameters | Notes |
| -------------- | ------ | ---------- | ----- |
| BodyKeySanitizer | `recording.AddBodyKeySanitizer(t, ...)` | ... | ... |
| BodyRegexSanitizer | `recording.BodyRegexSanitizer(t, ...)` | ... | ... |
| ContinuationSanitizer | `recording.ContinuationSanitizer(t, ...)` | ... | ... |
| GeneralRegexSanitizer | `recording.GeneralRegexSanitizer(t, ...)` | ... | ... |
| HeaderRegexSanitizer | `recording.HeaderRegexSanitizer(t, ...)` | ... | ... |
| OAuthResponseSanitizer | `recording.OAuthResponseSanitizer(t, ...)` | ... | ... |
| RemoveHeaderSanitizer | `recording.RemoveHeaderSanitizer(t, ...)` | ... | ... |
| ReplaceRequestSubscriptionId | `recording.ReplaceRequestSubscriptionId(t, ...)` | ... | ... |
| UriRegexSanitizer | `recording.UriRegexSanitizer(t, ...)` | ... | ... |

Note that removing the names of accounts and other values in your recording can have side effects when running your tests in playback. To take care of this, there are additional methods in the `internal/recording` module for reading environment variables and defaulting to the processed recording value. For example, if the `aztable` library had a test for creating a client and "requiring" the account name to be the same as provided it could potentially look similar to this:

```golang
func TestTableClient(t *testing.T) {
	accountName := recording.GetEnvVariable(t, "TABLES_PRIMARY_ACCOUNT_NAME", "fakeAccountName")
	// If running in playback, the value is "fakeAccountName". If running in "record" the value is whatever is stored in the environment variable
	accountKey := recording.GetEnvVariable(t, "TABLES_PRIMARY_ACCOUNT_KEY", "fakeAccountKey")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	client, err := NewTableClient("someTableName", someServiceURL, cred, nil)
	require.NoError(t, err)
	require.Equal(t, accountName, client.AccountName())
}
```

### Using `azidentity` Credentials In Tests

The credentials in `azidentity` are not automatically configured to run in playback mode. To make sure your tests run in playback mode even with `azidentity` credentials the best practice is to use a simple `FakeCredential` type that inserts a fake Authorization header to mock a credential. An example for swapping the `DefaultAzureCredential` using a helper function is shown below in the context of `aztable`

```golang
func getAADCredential() (azcore.Credential, error) {
	if recording.InPlayback() {
		return recording.NewFakeCredential("fakeAccountName", "fakeAccountKey"), nil
	}
	return azidentity.NewDefaultCredential(nil)
}

func TestTableClientWithAAD(t *testing.T) {
	accountName := recording.GetEnvVariable(t, "TABLES_PRIMARY_ACCOUNT_NAME", "fakeAccountName")
	cred, err := getAADCredential()
	require.NoError(t, err)
	...
	...run tests...
}
```

The `FakeCredential` type in `internal/recording` implements the `azcore.Credential` interface and can be used anywhere that the `azcore.Credential` is used.

<!-- LINKS -->
[doc_go_template]: https://github.com/Azure/azure-sdk-for-go/wiki/doc.go-template
[get_docker]: https://docs.docker.com/get-docker/
[go_azsdk_samples]: https://github.com/azure-samples/azure-sdk-for-go-samples
[go_download]: https://golang.org/dl/
[go_interfaces]: (https://gobyexample.com/interfaces)
[pipeline_definitions]: https://github.com/Azure/azure-sdk-for-go/blob/main/eng/pipelines/templates/jobs/archetype-sdk-client.yml
[require_package]: https://pkg.go.dev/github.com/stretchr/testify/require
[test_proxy_docs]: https://github.com/Azure/azure-sdk-tools/tree/main/tools/test-proxy
[workspace_setup]: https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-windows-10