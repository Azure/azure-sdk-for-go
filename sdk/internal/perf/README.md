# Performance Testing Framework
The `perf` sub-module provides a singular framework for writing performance tests.

## Default Command Options

| Flag | Short Flag | Default Value | Variable Name | Description |
| -----| ---------- | ------------- | ------------- | ----------- |
| `--duration` | `-d` | 10 seconds | internal.Duration (`int`) | How long to run an individual performance test |
| `--proxy` | `-x` | N/A | internal.TestProxy (`string`) | Whether to run a test against a test proxy. If you want to run against `https` specify with `--proxy https`, likewise for `http`. If you want to run normally omit this flag |
| `--warmup` | `-w` | 3 seconds| internal.WarmUp (`int`) | How long to allow the connection to warm up. |
| `--timeout` | `-t` | 10 seconds| internal.TimeoutSeconds (`int`) | How long to allow an operation to block. |

## Running with the test proxy

To run with the test proxy, configure your client to route requests to the test proxy using the `NewProxyTransport()` method. For example in `azkeys`:

```go
func (k *keysPerfTest) GlobalSetup() error {
    ...
    t, err := recording.NewProxyTransport(&recording.TransportOptions{UseHTTPS: true, TestName: a.GetMetadata()})
    if err != nil {return err}
    options = &azkeys.ClientOptions{
        ClientOptions: azcore.ClientOptions{
            Transport: t,
        },
    }

    client, err := azkeys.NewClient("<my-vault-url>", cred, options)
    if err != nil {return err}
    k.client = client
}
```

Before you can use the proxy in playback mode, you have to generate live recordings. To do this, call `recording.Start` before the test run and `recording.Stop` after the test run.

```go
func (k *keysPerfTest) GlobalSetup() error {
    ...
    err := recording.Start(k.GetMetadata(), nil)
    ...
}

func (k *keysPerfTest) GlobalTeardown() error {
    ...
    err := recording.Stop(k.GetMetadata(), nil)
    return err
}
```

## Adding Performance Tests to an SDK

1. Copy the `cmd/template` directory into a directory for your package:
```pwsh
Copy-Item cmd/template cmd/<mypackage> -Recurse
```

2. Change the name of `TemplateCmd`. Best practices are to use a `<packageName>Cmd` as the name. Fill in `Use`, `Short`, `Long`, and `RunE` for your specific test.

3. Implement the `GlobalSetup`, `Setup`, `Run`, `TearDown`, `GlobalTearDown`, and `GetMetadata` functions. All of these functions will take a `context.Context` type to prevent an erroneous test from blocking. `GetMetadata` should return a `string` with the name of the specific test. This is only used by the test proxy for generating and reading recordings. `GlobalSetup` and `GlobalTearDown` are called once each, at the beginning and end of the performance test respectively. These are good places to do client instantiation, create instances, etc. `Setup` and `TearDown` are called once before each test iteration. If there is nothing to do in any of these steps, you can use `return nil` for the implementation.

4. Register the test in `cmd/root.go` by adding the command to the `init` function and importing your package:

```golang
import (
    ...
    mypackage "github.com/Azure/azure-sdk-for-go/eng/tools/azperf/cmd/mypackage"
)


func init() {
    ...
    rootCmd.AddCommand(mypackage.MyPackageCmd)
}
```

You can create multiple performance tests all in one directory. Each performance test should be it's own unique `*cobra.Command` type and each needs to be registered with the `rootCmd.AddCommand` function.

### Writing a test

### Testing with streams
TODO

### Writing a Batch Test
TODO

## Running Performance Tests

First, make sure you have go version 1.17 or later installed. You can install go from the [go.dev](https://go.dev/doc/install) site.

Navigate to the `eng/tools/azperf` folder and run `go mod download` to download all the requirements for the performance framework.

Build the CLI with the command:
```pwsh
go build .
```

Now you will have a single executable `azperf.exe` on Windows, or `azperf` on Mac/Linux, which can be used to run all performance tests.

To list all the performance tests available run `./azperf.exe --help` and a list will show up along with optional parameters.

To run a single performance test specify the test as the second argument:
```pwsh
./azperf.exe CreateEntityTest
```

To specify flags for a performance test, add them after the second argument:
```pwsh
./azperf.exe CreateEntityTest --duration 7 --proxy https
```

For more information about running individual tests, refer to the READMEs of each test linked below.

### Available Performance Tests

| Name | Options | Package Testing | Description | Additional Flags | More Information |
| ---- | ------- | --------------- | ----------- | ---------------- | ---------------- |
| CreateEntityTest | None | `aztables` | Creates a single entity | None | [README](https://github.com/Azure/azure-sdk-for-go) |
| CreateKeyTest | None | `azkeys` | Creates a single RSA key | None | [README](https://github.com/Azure/azure-sdk-for-go) |
| UploadBlobTest | None | `azblobs` | Uploads random data to a blob | None | [README](https://github.com/Azure/azure-sdk-for-go) |