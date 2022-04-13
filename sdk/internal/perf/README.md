# Performance Testing Framework
The `perf` sub-module provides a framework for writing and running performance tests.

## Default Command Options

| Flag | Short Flag | Default Value | Description |
| -----| ---------- | ------------- | ----------- |
| `--duration` | `-d` | 10 seconds | How long to run an individual performance test |
| `--test-proxies` | `-x` | N/A | A semicolon separated list of proxy urls. |
| `--warmup` | `-w` | 3 seconds| How long to allow the connection to warm up. |


## Adding Performance Tests to an SDK

1. Create a performance test directory at `testdata/perf` within your module. For example, the storage performance tests live in `sdk/storage/azblob/testdata/perf`.

2. Run `go mod init` to create a new module. Add the `sdk/internal` module and the module you are testing for to your go.mod file and run `go mod tidy`

3. Create a `struct` that maintains all global values for the performance test (ie. account name, blob name, etc.).

4. Implement `GlobalPerfTest` interface by adding the `NewPerfTest(context.Context, *PerfTestOptions) (PerfTest, error)` and the `GlobalCleanup(context.Context)` functions on the global struct. The `NewPerfTest` method creates a new struct that is responsible for running a performance test in its own `goroutine`. `GlobalCleanup` takes care of any session level cleanup.

5. Implement the `PerfTest` interface on the struct returned from the `NewPerfTest` method. Add a `Run(context.Context)` and a `Cleanup(context.Context)` function. The `Run` method is the method you want to measure the performance of. The `Cleanup` method is for cleaning up within a single `goroutine`, most of the time this will be empty.

6. Create a main package that uses the `perf.Run` function to provide all the global constructor methods. `perf.Run` takes a map[string]PerfMethods

6. (Optional): Create a method that registers local flags.

```golang
package main

import (
    "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

// uploadTestRegister registers flags for the "UploadBlobTest"
// This is optional and does not to be included for every test.
func uploadTestRegister() {
	flag.IntVar(&uploadTestOpts.size, "size", 10240, "Size in bytes of data to be transferred in upload or download tests. Default is 10240.")
}

func main() {
	perf.Run(map[string]perf.PerfMethods{
		"UploadBlobTest":   {Register: uploadTestRegister, New: NewUploadTest},
		"ListBlobTest":     {Register: listTestRegister, New: NewListTest},
		"DownloadBlobTest": {Register: downloadTestRegister, New: NewDownloadTest},
	})
}
```

### Writing a test
The following walks through the `azblob` Blob Upload performance test:

#### `GlobalPerfTest` Interface
This struct handles global set up for your account and spawning structs for each `goroutine`. Make sure to embed the `perf.PerfTestOptions` struct in your global struct. The `NewUploadTest` will be run once per process and is used to spawn each parallel test.
```go
type uploadTestGlobal struct {
	perf.PerfTestOptions
	containerName         string
	blobName              string
	globalContainerClient azblob.ContainerClient
}

// NewUploadTest is called once per process
func NewUploadTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	u := &uploadTestGlobal{
		PerfTestOptions: options,
		containerName:   "uploadcontainer",
		blobName:        "uploadblob",
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, u.containerName, nil)
	if err != nil {
		return nil, err
	}
	u.globalContainerClient = containerClient
	_, err = u.globalContainerClient.Create(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return u, nil
}
```

The `GlobalCleanup` method cleans up the account, resetting it into the default.
```go
func (m *uploadTestGlobal) GlobalCleanup(ctx context.Context) error {
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

`NewPerfTest` is called once per `goroutine`, and creates a new `PerfTest` interface which will be used by each goroutine. This method should also include the setup for the eventual returned struct.
```go
// uploadPerfTest implements perf.PerfTest and is used to test performance of your event
type uploadPerfTest struct {
	*uploadTestGlobal
	perf.PerfTestOptions
	data       io.ReadSeekCloser
	blobClient azblob.BlockBlobClient
}

// NewPerfTest is called once per goroutine
func (g *uploadTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	u := &uploadPerfTest{
		uploadTestGlobal: g,
		PerfTestOptions:  *options,
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(
		connStr,
		u.uploadTestGlobal.containerName,
		&azblob.ClientOptions{
			Transporter: u.PerfTestOptions.Transporter,
		},
	)
	if err != nil {
		return nil, err
	}
	bc := containerClient.NewBlockBlobClient(u.blobName)
	u.blobClient = bc

	data, err := perf.NewRandomStream(uploadTestOpts.size)
	if err != nil {
		return nil, err
	}
	u.data = data

	return u, nil
}
```

#### `PerfTest` Interface
The `PerfTest` interface is an interface that is responsible for running a single performance test. Each performance test in run within a single `goroutine`, these `goroutine`s are created by the `perf` framework. The `Run` method is the method that is being measured and the `Cleanup` method is responsible for cleanup work after each goroutine completes or errors out.

```go
func (m *uploadPerfTest) Run(ctx context.Context) error {
	_, err := m.data.Seek(0, io.SeekStart) // rewind to the beginning
	if err != nil {
		return err
	}
	_, err = m.blobClient.Upload(ctx, m.data, nil)
	return err
}

func (m *uploadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
```

#### Constructor
Each method needs to have a constructor that returns an instantialized `perf.GlobalPerfTest`.
```go
// NewUploadTest is called once per process
func NewUploadTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	u := &uploadTestGlobal{
		PerfTestOptions: options,
		containerName:   "uploadcontainer",
		blobName:        "uploadblob",
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, u.containerName, nil)
	if err != nil {
		return nil, err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return u, nil
}
```

#### `main.go` file
The `main` function must have the `perf.Run` method with a map of each performance test constructor. Flag parsing, test selection, reporting, proxy interfaces, and all other portions are taken care of for the user in the `perf.Run` method.

```go
package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

func main() {
	perf.Run(map[string]perf.PerfMethods{
		"UploadBlobTest":   {Register: uploadTestRegister, New: NewUploadTest},
		"ListBlobTest":     {Register: listTestRegister, New: NewListTest},
		"DownloadBlobTest": {Register: downloadTestRegister, New: NewDownloadTest},
	})
}
```

### Running with the test proxy

To run with the test proxy, configure your client to route requests to the test proxy using the `Transporter` parameters for the `Transport` in your client options. If the `--test-proxy` flag is not specified, `Transporter` is the same default HTTP client from `azcore`.

```go
func (k *keysPerfTest) GlobalSetup() error {
    ...
    options = &azkeys.ClientOptions{
        ClientOptions: azcore.ClientOptions{
            Transport: k.Transporter,
        },
    }

    client, err := azkeys.NewClient("<my-vault-url>", cred, options)
    if err != nil {return err}
    ...
}
```

### Registering Local Arguments
We use the `flag` library for argument parsing. Each perf test suite can have optional arguments and they are registered using the `perf.RegisterArguments` method, which takes a simple function that registers arguments using `flag.Int`, `flag.String`, etc. If you have no local arguments, you can skip this step.

```go
type uploadTestOptions struct {
	size int
}

var uploadTestOpts uploadTestOptions = uploadTestOptions{size: 10240}

// uploadTestRegister is called once per process
func uploadTestRegister() {
	flag.IntVar(&uploadTestOpts.size, "size", 10240, "Size in bytes of data to be transferred in upload or download tests.")
}
```

The `main.go` file changes to:
```go
package main

import "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"

func main() {
	perf.Run(map[string]perf.PerfMethods{
		"UploadBlobTest":   {Register: uploadTestRegister, New: NewUploadTest},
		"ListBlobTest":     {Register: listTestRegister, New: NewListTest},
		"DownloadBlobTest": {Register: downloadTestRegister, New: NewDownloadTest},
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

To run against multiple proxies and with multiple parallel tests running at the same time use the `-p/--parallel` flag
```pwsh
go run . CreateEntityTest --duration 7 --parallel 4 --test-proxies "https://localhost:5001;http://localhost:5000" --num-blobs 100
```

This runs four goroutines of the same test and splits the traffic to two different proxy addresses. The first and third routine will target `https://localhost:5001`, the second and fourth will use `http://localhost:5000`.

For help run `go run . --help`. A list of registered performance tests, global command flags, and local command flags will print out.
