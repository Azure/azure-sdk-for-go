# Developer Set Up

* [Installing Go](#installing-go)
* [Create a Client](#create-a-client)
	* [Documenting Code](#documenting-code)
    * [Constructors](#constructors)
	* [Defining Methods](#defining-methods)
* [Write Tests](#write-tests)

## Installing Go

The Azure-sdk-for-go team supports Go versions 1.14 and greater, with our CI pipelines testing against versions 1.14 and 1.15 on both Windows and Linux virtual machines. If you do not already have Go installed, refer to this [workspace setup][workspace_setup] article for a more in depth tutorial on setting up your Go environment. After installing Go and configuring your workspace, fork the `azure-sdk-for-go` repository and clone it to a directory that looks like: `<GO HOME>/src/github.com/Azure/azure-sdk-for-go`.


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

Each package needs to include a `doc.go` file and not be a part of a service version. For more details about this file there is a detailed write-up in the [repo wiki](https://github.com/Azure/azure-sdk-for-go/wiki/doc.go-template). In the `doc.go` file you should include a short service overview, basic examples, and (if they exist) a link to samples in the [`azure-sdk-for-go-samples` repository](https://github.com/azure-samples/azure-sdk-for-go-samples)

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

This client takes three parameters, the first is the service URL for the specific account. The second is an [`interface`](https://gobyexample.com/interfaces) which is a specific struct that has definitions for a certain set of methods. In the case of `azcore.Credential` the `AuthenticationPolicy(options AuthenticationPolicyOptions) Policy` method must be defined to be a valid interface. The final argument to methods that create clients or interact with the service should be a pointer to an `Options` parameter. Making this final parameter a pointer allows the customer to pass in `nil` if there are no specific options they want to change. The `Options` type should have a name that is intuitive to what the customer is trying to do, in this case `TableClientOptions`.

### Defining Methods
Defining a method follows the format:
```golang
func (m *<MyStruct>) MethodName(param1 param1Type, param2 param2Type) (ReturnType, ReturnType2) {

}
```
The `(m *<MyStruct>)` portion is the "receiver". Methods can be defined for either pointer (with a `*`) or receiver (without a `*`) types. Pointer receivers will avoid copying types on method calls and allow the method to mutate the receiving struct. You should use pointer receivers wherever possible to limit memory copies.


Both public and private methods can be declared on clients. Below is an example in the `aztables` package for a `Create` method on the `TableServiceClient`:
```golang
// Create creates a table with the specified name.
func (t *TableServiceClient) Create(ctx context.Context, name string) (TableCreateResponse, error) {
	resp, err := t.client.Create(ctx, TableProperties{&name}, new(TableCreateOptions), new(QueryOptions))
	if err == nil {
		tableResp := resp.(TableCreateResponse)
		return tableResp, nil
	}
	return TableCreateResponse{}, err
}
```

All methods that perform I/O of any kind, sleep, or perform a significant amount of CPU-bound work must have the first parameter be of type [`context.Context`][golang_context] which allows the customer to carry a deadline, cancellation signal, and other values across API boundaries. The remaining parameters should be parameters specific to that method. The return types for methods should be first a "Response" object and second an `error` object.

## Write Tests

Testing is built into the Go toolchain as well with the `testing` library. The testing infrastructure located in the `sdk/internal` directory takes care of generating recordings, establishing the mode a test is being run in (options are "recording", "playback", "live-no-playback"), and reading environment variables.

A simple test for `aztables` is shown below:
```golang

import (
	"os"

	"github.com/testify/assert"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const (
	accountName := os.GetEnv("TABLES_PRIMARY_ACCOUNT_NAME")
	accountKey := os.GetEnv("TABLES_PRIMARY_ACCOUNT_KEY")
	mode := recording.Recording
)

// Test creating a single table
func TestCreateTable(t *testing.T) {
	client := NewTableClient(accountName, accountKey, "tableName")
	resp, err := client.Create()
	assert.Nil(t, err)
	assert.Equal(t, respo.TableResponse.TableName, "tableName")
}
```

The first part of the test above is for getting the secrets needed for authentication from your environment, the current practice is to store your test secrets in environment variables.

The rest of the snippet shows a test that creates a single table and asserts that the response from the service has the same table name as the supplied parameter. Every test in Go has to have exactly one parameter, the `t *testing.T` object, and it must begin with `Test`. After making a service call or creating an object you can make assertions on that object by using the external `testify/assert` library. In the example above, we assert that the error returned is `nil`, meaning the call was successful and then we assert that the response object has the same table name as supplied.

You can also use the `testify/require` library instead of `testify/assert` if you want your test to fail as soon as you have an unexpected result.

Check out the docs for more information about the methods available in the [`assert`](https://pkg.go.dev/github.com/stretchr/testify/assert) and [`require`](https://pkg.go.dev/github.com/stretchr/testify/require) libraries.


<!-- LINKS -->
[workspace_setup]: https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-windows-10