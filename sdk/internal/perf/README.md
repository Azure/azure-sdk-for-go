# Performance Testing Framework
The `perf` sub-module provides a singular framework for writing performance tests.

## Default Command Options

| Flag | Short Flag | Default Value | Variable Name | Description |
| -----| ---------- | ------------- | ------------- | ----------- |
| `--duration` | `-d` | 10 seconds | internal.Duration (`int`) | How long to run an individual performance test |
| `--test-proxies` | `-x` | N/A | internal.TestProxy (`string`) | Whether to run a test against a test proxy. If you want to run against `https` specify with `--test-proxies https`, likewise for `http`. If you want to run normally omit this flag |
| `--warmup` | `-w` | 3 seconds| internal.WarmUp (`int`) | How long to allow the connection to warm up. |


## Adding Performance Tests to an SDK

1. Create a performance test directory at `testdata/perf` within your module. For example, the storage performance tests live in `sdk/storage/azblob/testdata/perf`.
2. Run `go mod init` to create a new module.
3. Create a `struct` that implements the `perf.PerfTest` interface by implementing the `GlobalSetup`, `Setup`, `Run`, `Cleanup`, `GlobalCleanup`, and `GetMetadata` functions. The first five functions take the `context.Context` type to prevent an erroneous test from blocking. `GetMetadata` should returns the `perf.PerfTestOptions` that is embedded in the struct.
    `GlobalSetup` and `GlobalCleanup` are called once each, at the beginning and end of the performance test respectively.
    `Setup` and `Cleanup` are called once for each test instance. If you're running four performance tests in parallel, this method is run four times, once for each parallel instance. If there is nothing to do in any of these steps, you can `return nil`.
4. Create a `New<Name>Test` method that returns an instantialized struct.

5. Create a main package that uses the `perf.Run` function to provide all the constructor methods.

6. (Optional): Create a method that registers local flags.

```golang
package main

import (
    ...
    "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

func main() {
    perf.Run([]perf.NewPerfTest{
        FirstTestConstructor,
        SecondTestConstructor,
        ...
    })
}
```

### Writing a test
The following walks through the `azblob` Blob Upload performance test:

#### GlobalSetup
Setup for the test, this is run one time per test. This should not have any struct initializations. In the example, we are creating the container in the storage account.
```go
func (u *uploadPerfTest) GlobalSetup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, m.containerName, nil)
	if err != nil {
		return err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}
```

#### Setup
In the `Setup` method, initialize your clients for the Run stage. Always use the `ProxyInstance` from the embedded `perf.PerfTestOptions` struct as your ClientOptions.`Transport`. The `ProxyInstance` will handle all routing for performance tests run against a proxy. You should create any clients you will use in the `Run` and `Cleanup` methods.
```go
func (u *uploadPerfTest) Setup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, u.containerName, &azblob.ClientOptions{Transporter: u.ProxyInstance})
	if err != nil {
		return err
	}

	m.blobClient = containerClient.NewBlockBlobClient(u.blobName)
	return nil
}
```

In this example we create a `BlobClient` and save it in the struct.

#### Run
This method is the function you are testing for performance. Use the context from the function signature if your method uses one. This will prevent test runs from hanging or blocking. Don't handle errors in this method, rather return them and let the performance framework handle erro handling.

```go
func (u *uploadPerfTest) Run(ctx context.Context) error {
	_, err := u.data.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = u.blobClient.Upload(ctx, u.data, nil)
	return err
}
```

#### Cleanup
This method is for cleaning up after a single test run. This method will most likely be empty, just returning `nil`. Don't do any deletions or modifications to the account in this method.

```go
func (m *uploadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
```

#### GlobalCleanup
Use the method to delete any provisioned resources, and restore your account to the defaults.

```go
func (m *uploadPerfTest) GlobalCleanup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, m.containerName, nil)
	if err != nil {
		return err
	}

	_, err = containerClient.Delete(context.Background(), nil)
	return err
}
```
#### GetMetadata
This is the simplest method, just return the `PerfTestOptions` struct from your defined struct.

```go
func (m *uploadPerfTest) GetMetadata() perf.PerfTestOptions {
	return m.PerfTestOptions
}
```

#### Constructor
Each method needs to have a constructor that returns an instantialized `perf.PerfTest`. This storage example does a few things:

* sets the `Name` parameter on the options struct
* sets the default value for the size parameter used by this test.
* Creates a data source from the `perf.RandomStream` (which implements the `io.ReadSeekCloser`).
```go
func NewUploadTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "BlobUploadTest"

	if size == nil {
		size = to.Int64Ptr(10240)
	}

	data, err := perf.NewRandomStream(int(*size))
	if err != nil {
		panic(err)
	}
	return &uploadPerfTest{
		PerfTestOptions: *options,
		blobName:        "uploadtest",
		containerName:   "uploadcontainer",
		data:            data,
	}
}
```

#### `main.go` file
The main function must have the `perf.Run` method with a slice of each performance test constructor. Flag parsing, test selection, reporting, proxy interfaces, and all other portions are taken care of for the user in the `perf.Run` method.

```go
package main

import "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"

func main() {
	perf.RegisterArguments(RegisterArguments)
	perf.Run([]perf.NewPerfTest{
		NewDownloadTest,
		NewListTest,
		NewUploadTest,
	})
}
```

### Running with the test proxy

To run with the test proxy, configure your client to route requests to the test proxy using the `ProxyInstance` parameters for the `Transport` in your client options. For example in `azkeys`:

```go
func (k *keysPerfTest) GlobalSetup() error {
    ...
    options = &azkeys.ClientOptions{
        ClientOptions: azcore.ClientOptions{
            Transport: k.ProxyInstance,
        },
    }

    client, err := azkeys.NewClient("<my-vault-url>", cred, options)
    if err != nil {return err}
    ...
}
```

### Registering Local Arguments
We use the `pflag` library for argument parsing (the standard library does not support double dashed arguments). Each perf test suite can have optional arguments and they are registered using the `perf.RegisterArguments` method, which takes a simple function that registers arguments using `pflag.Int32`, `pflag.String`, etc. If you have no local arguments, you can skip this step.

```go
var size *int64
var count *int32

func RegisterArguments() {
	count = pflag.Int32("num-blobs", 100, "Number of blobs to list. Defaults to 100.")
	size = pflag.Int64("size", 10240, "Size in bytes of data to be transferred in upload or download tests. Default is 10240.")
}
```

The `main.go` file changes to:
```go
package main

import "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"

func main() {
	perf.RegisterArguments(RegisterArguments)
	perf.Run([]perf.NewPerfTest{
		NewDownloadTest,
		NewListTest,
		NewUploadTest,
	})
}
```

## Running Performance Tests

First, make sure you have go version 1.17 or later installed. You can install go from the [go.dev](https://go.dev/doc/install) site.

Navigate to your SDK performance test, (ie. `sdk/keyvault/azkeys/testdata/perf` folder and run `go mod download` to download all the requirements for the performance framework.

To run a single performance test specify the test as the first argument:
```pwsh
go run . CreateEntityTest
```

To specify flags for a performance test, add them after the first argument:
```pwsh
go run . CreateEntityTest --duration 7 --test-proxies https://localhost:5001 --num-blobs 100
```

For help run `go run . --help`. A list of registered performance tests, global command flags, and local command flags will print out.
