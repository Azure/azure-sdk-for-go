# Azure SDK for Go Recorded Test Framework

[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/Azure.azure-sdk-for-go?branchName=master)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=1842&branchName=master)

The `recording` package makes it easy to add recorded tests to your track-2 client package.
Below are some examples that walk through setting up a recorded test end to end.

## Set Up
The `recording` package supports three different test modes. These modes are set by setting the `AZURE_RECORD_MODE` environment variable to one of the three values:
1. `record`: Used when making requests against real resources. In this mode new recording will be generated and saved to file.
2. `playback`: This mode is for running tests against local recordings.
3. `live`: This mode should not be used locally, it is used by the nightly live pipelines to run against real resources and skip any routing to the proxy. This mode closely mimics how our customers will use the libraries.

After you've set the `AZURE_RECORD_MODE`, set the `PROXY_CERT` environment variable to:
```pwsh
$ENV:PROXY_CERT="C:/ <path-to-repo> /azure-sdk-for-go/eng/common/testproxy/dotnet-devcert.crt"
```

## Running the test proxy

Recording and playing back tests relies on the [Test Proxy](https://github.com/Azure/azure-sdk-tools/blob/main/tools/test-proxy/Azure.Sdk.Tools.TestProxy/README.md) to intercept traffic. The recording package can automatically install and run an instance of the test-proxy server per package. The below code needs to be added to test setup and teardown in order to achieve this. The service directory value should correspond to the `testdata` directory within the directory containing the [assets.json](https://github.com/Azure/azure-sdk-tools/blob/main/tools/assets-automation/asset-sync/README.md#assetsjson-discovery) config, see [examples](https://github.com/search?q=repo%3AAzure%2Fazure-sdk-for-go+assets.json+language%3AJSON&type=code&l=JSON).

```golang
const recordingDirectory = "<path to service directory with assets.json file>/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
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

    ... all other test code, including proxy recording setup ...
	return m.Run()
}
```

## Routing Traffic

The first step in instrumenting a client to interact with recorded tests is to direct traffic to the proxy through a custom `policy`. In these examples we'll use testify's [`require`](https://pkg.go.dev/github.com/stretchr/testify/require) library but you can use the framework of your choice. Each test has to call `recording.Start` and `recording.Stop`, the rest is taken care of by the `recording` library and the [`test-proxy`](https://github.com/Azure/azure-sdk-tools/tree/main/tools/test-proxy).

The snippet below demonstrates an example test policy:

```go
type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

func (r recordingPolicy) Host() string {
	if r.options.UseHTTPS {
		return "localhost:5001"
	}
	return "localhost:5000"
}

func (r recordingPolicy) Scheme() string {
	if r.options.UseHTTPS {
		return "https"
	}
	return "http"
}

func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{UseHTTPS: true}
	}
	p := &recordingPolicy{options: *o, t: t}
	return p
}

func (p *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != "live" {
		p.options.ReplaceAuthority(t, req.Raw())
	}
	return req.Next()
}
```

After creating a recording policy, it has to be added to the client on the `ClientOptions.PerCallPolicies` option:
```go
func TestSomething(t *testing.T) {
    p := NewRecordingPolicy(t)
    httpClient, err := recording.GetHTTPClient(t)
    require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			PerCallPolicies: []policy.Policy{p},
			Transport:       client,
		},
	}

    client, err := NewClient("https://mystorageaccount.table.core.windows.net", myCred, options)
    require.NoError(t, err)
    // Continue test
}
```

## Starting and Stopping a Test
To start and stop your tests use the `recording.Start` and `recording.Stop` (make sure to use `defer recording.Stop` to ensure the proxy cleans up your test on failure) methods:
```go
func TestSomething(t *testing.T) {
    err := recording.Start(t, recordingDirectory, nil)
    defer recording.Stop(t, recordingDirectory, nil)

    // Continue test
}
```

## Using Sanitizers
The recording files generated by the test-proxy are committed along with the code to the public repository. We have to keep our recording files free of secrets that can be used by bad actors to infilitrate services. To do so, the `recording` package has several sanitizers for taking care of this. Sanitizers are added at the session level (ie. for an entire test run) to apply to all recordings generated during a test run. For example, to replace the account name from a storage url use the `recording.AddURISanitizer` method:
```go
func TestSomething(t *testing.T) {
    err := recording.AddURISanitizer("fakeaccountname", "my-real-account-name", nil)
    require.NoError(t, err)

    // To remove the sanitizer after this test use the following:
    defer recording.ResetSanitizers(nil)

    err := recording.Start(t, recordingDirectory, nil)
    defer recording.Stop(t, recordingDirectory, nil)

    // Continue test
}
```

In addition to URI sanitizers, there are sanitizers for headers, response bodies, OAuth responses, continuation tokens, and more. For more information about all the sanitizers check out the [source code](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/internal/recording/sanitizer.go)


## Reading Environment Variables
The CI pipelines for PRs do not run against live resources, you will need to make sure that the values that are replaced in the recording files are also replaced in your requests when running in playback. The best way to do this is to use the `recording.GetEnvVariable` and use the replaced value as the `recordedValue` argument:

```go
func TestSomething(t *testing.T) {
    accountName := recording.GetEnvVariable(t, "TABLES_PRIMARY_ACCOUNT_NAME", "fakeaccountname")
    if recording.GetRecordMode() = recording.RecordMode {
        err := recording.AddURISanitizer("fakeaccountname", accountName, nil)
        require.NoError(t, err)
    }

    // Continue test
}
```
In this snippet, if the test is running in live mode and we have the real account name, we want to add a URI sanitizer for the account name to ensure the value does not appear in any recordings.

