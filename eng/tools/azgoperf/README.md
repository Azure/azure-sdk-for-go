# Performance Testing Framework
The `azgoperf` CLI tool provides a singular framework for writing and running performance tests.

## Default Command Options

| Flag | Short Flag | Default Value | Variable Name | Description |
| -----| ---------- | ------------- | ------------- | ----------- |
| `--duration` | `-d` | 10 seconds | duration (`int`) | How long to run an individual performance test |
| `--iterations` | `-i` | 3 | iterations (`int`) | How many iterations of a performance test to run |
| `--testproxy` | `-x` | N/A | testProxy (`string`) | Whether to run a test against a test proxy. If you want to run against `https` specify with `--testproxy https`, likewise for `http`. If you want to run normally omit this flag |
| `--timeout` | `-t` | 10 seconds| timeoutSeconds (`int`) | How long to allow an operation to block. |

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

1. Copy the `template.go` file:
```
Copy-Item template.go aztables.go
```

2. Change the name of `templateCmd`. Best practices are to use a `<packageName>Cmd` as the name. Fill in `Use`, `Short`, `Long`, and `RunE` for your specific test.

3. Implement the `GlobalSetup`, `Setup`, `Run`, `TearDown`, `GlobalTearDown`, and `GetMetadata` functions. All of these functions will take a `context.Context` type to prevent an erroneous test from blocking. `GetMetadata` should return a `string` with the name of the specific test. This is only used by the test proxy for generating and reading recordings. `GlobalSetup` and `GlobalTearDown` are called once each, at the beginning and end of the performance test respectively. These are good places to do client instantiation, create instances, etc. `Setup` and `TearDown` are called once before each test iteration. If there is nothing to do in any of these steps, you can use `return nil` for the implementation.

### Writing a test

### Testing with streams
TODO

### Writing a Batch Test
TODO

## Running Performance Tests

First, compile the cli into a single executable:
```pwsh
go build .
```

To run a single performance test specify the test as the second argument:
```pwsh
./azgoperf.exe azkeys
```

To specify flags for a performance test, add them after the second argument:
```pwsh
./azgoperf.exe azkeys --duration 7 --testproxy https
```

### Available Performance Tests

| Name | Options | Package Testing | Description |
| ---- | ------- | ----------- |
| CreateEntityTest | None | `aztables` | Creates a single entity |
| CreateKeyTest | None | `azkeys` | Creates a single RSA key |