# Service Onboarding

This guide describes how to take an OpenAPI (or Swagger) spec located in the azure-rest-api-specs repository and generate the autorest code, create a basic client, and write basic tests for your client.

* [Install Go](#install-go)
* [Install AutoRest](#install-autorest)
    * [Generating Code](#generating-code)
* [Create a Client](#create-a-client)
	* [Documenting Code](#documenting-code)
    * [Constructors](#constructors)
	* [Defining Methods](#defining-methods)
* [Write Tests](#write-tests)
* [Create Pipelines](#create-pipelines)

## Install Go

The Azure-sdk-for-go team supports Go versions 1.14 and greater, with our CI pipelines testing against versions 1.14 and 1.15 on both Windows and Linux virtual machines. If you do not already have Go installed, refer to this [workspace setup][workspace_setup] article for a more in depth tutorial on setting up your Go environment. After installing Go and configuring your workspace, fork the `azure-sdk-for-go` repository and clone it to a directory that looks like: `<GO HOME>/src/github.com/Azure/azure-sdk-for-go`.

## Install Autorest

* swagger.md file
### Generating Code


## Create a Client

After you have the generated code from Autorest, the next step is to wrap this generated code in a "convenience layer" that the customers will use directly to interact with the service. Go is not an object-oriented language like C#, Java, or Python. There is no type hierarchy in Go. Clients and models will be defined as `struct`s and methods will be defined on these structs to interact with the service.

In other languages, types can be specifically marked "public" or "private", in Go exported types and methods are defined by starting with a capital letter. The methods on structs also follow this rule, if it is for use outside the model it must start with a capital letter.

### Documenting Code
Code is documented directly in line and can be created directly using the Go toolchain.

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
func NewTableServiceClient(serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableServiceClient, error) {
	conOptions := options.getConnectionOptions()
	if isCosmosEndpoint(serviceURL) {
		conOptions.PerCallPolicies = []azcore.Policy{CosmosPatchTransformPolicy{}}
	}
	con := newConnection(serviceURL, cred, conOptions)
	c, err := cred.(*SharedKeyCredential)
	return &TableServiceClient{client: &tableClient{con}, service: &serviceClient{con}, cred: *c}, err
}
```
In `Go`, the parameters are surrounded in parenthesis immediately following the method name with the parameter name preceding the type of the parameter. Following the parameters and a closing parentheses is the return arguments. If a method has more than one return parameter the types of the parameter must be enclosed in parenthesis. Note the `*` before a type indicates a pointer to that type. All methods that create a new client or interact with the service should return an `error` type as the last argument.

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
func (t *TableServiceClient) Create(ctx context.Context, name string) (TableResponseResponse, error) {
	resp, err := t.client.Create(ctx, TableProperties{&name}, new(TableCreateOptions), new(QueryOptions))
	if err == nil {
		tableResp := resp.(TableResponseResponse)
		return tableResp, nil
	}
	return TableResponseResponse{}, err
}
```

All methods that make a call to service must have the first parameter be of type [`context.Context`][golang_context] which allows the customer to do SOMETHING TO BE FILLED IN LATER. The remaining parameters should be parameters specific to that method. The return types for methods should be first a "Response" object and second an `error` object.



## Write Tests

## Create Pipelines

<!-- LINKS -->
[workspace_setup]: https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-windows-10
[golang_context]: https://golang.org/pkg/context/#Context