# Performance Testing Framework
The `azgoperf` CLI tool provides a singular framework for writing and running performance tests.

## Default Command Options

| Flag | Short Flag | Default Value | Description |
| -----| ---------- | ------------- | ----------- |
| `--duration` | `-d` | 10 seconds| How long to run an individual performance test |
| `--iterations` | `-i` | 3 | How many iterations of a performance test to run |
| `--testproxy` | `-x` | N/A | Whether to run a test against a test proxy. If you want to run against `https` specify with `--testproxy https`, likewise for `http`. If you want to run normally omit this flag |
| `--timeout` | `-t` | 10 seconds| How long to allow an operation to block. |

## Running with the test proxy

To run with the test proxy, configure your client to route requests to the test proxy using the `NewProxyTransport()` method. For example in `azkeys`:

```golang
transport, err := recording.NewProxyTransport()
handle(err)

client, err := azkeys.NewClient("vault-url", cred, &azkeys.ClientOptions)
```

## Adding Performance Tests to an SDK

### Writing a test

### Testing with streams

### Writing a Batch Test

## Running Performance Tests

### Test Commands

Example test run commands
```
~/github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf> go build .
~/github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf> ./azgoperf.exe azkeys
```