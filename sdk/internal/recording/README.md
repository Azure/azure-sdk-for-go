# Azure SDK for Go Recorded Test Framework

[![Build Status](https://dev.azure.com/azure-sdk/public/_apis/build/status/go/Azure.azure-sdk-for-go?branchName=master)](https://dev.azure.com/azure-sdk/public/_build/latest?definitionId=1842&branchName=master)

The `testframework` package makes it easy to add recorded tests to your track-2 client package.
Below are some examples that walk through setting up a recorded test end to end.

## Examples

### Initializing a Recording instance for a test

The first step in instrumenting a client to interact with recorded tests is to create a `TestContext`.
This acts as the interface between the recorded test framework and your chosen test package.
In these examples we'll use testify's [assert](https://pkg.go.dev/github.com/stretchr/testify/assert),
but you can use the framework of your choice.

In the snippet below, demonstrates an example test setup func in which we are initializing the `TestContext`
with the methods that will be invoked when your recorded test needs to Log, Fail, get the Name of the test,
or indicate that the test IsFailed.

***Note**: an instance of TestContext should be initialized for each test.*

```go
type testState struct {
    recording *recording.Recording
    client    *TableServiceClient
    context   *recording.TestContext
}
// a map to store our created test contexts
var clientsMap map[string]*testState = make(map[string]*testState)

// recordedTestSetup is called before each test execution by the test suite's BeforeTest method
func recordedTestSetup(t *testing.T, testName string, mode recording.RecordMode) {
    var accountName string
    var suffix string
    var cred *SharedKeyCredential
    var secret string
    var uri string
    assert := assert.New(t)

    // init the test framework
    context := recording.NewTestContext(func(msg string) { assert.FailNow(msg) }, func(msg string) { t.Log(msg) }, func() string { return testName })
    //mode should be recording.Playback. This will automatically record if no test recording is available and playback if it is.
    recording, err := recording.NewRecording(context, mode) 
    assert.Nil(err)
```

After creating the TestContext, it must be passed to a new instance of `Recording` along with the current test mode.
`Recording` is the main component of the testframework package.

```go
//func recordedTestSetup(t *testing.T, testName string, mode recording.RecordMode) {
//  <...>
    record, err := recording.NewRecording(context, mode)
    assert.Nil(err)
```

### Initializing recorded variables

A key component to recorded tests is recorded variables.
They allow creation of values that stay with the test recording so that playback of service operations is consistent.

In the snippet below we are calling `GetRecordedVariable` to acquire details such as the service account name and
client secret to configure the client.

```go
//func recordedTestSetup(t *testing.T, testName string, mode recording.RecordMode) {
//  <...>
    accountName, err := record.GetEnvVar(storageAccountNameEnvVar, recording.NoSanitization)
    suffix := record.GetOptionalEnvVar(storageEndpointSuffixEnvVar, DefaultStorageSuffix, recording.NoSanitization)
    secret, err := record.GetEnvVar(storageAccountKeyEnvVar, recording.Secret_Base64String)
    cred, _ := NewSharedKeyCredential(accountName, secret)
    uri := storageURI(accountName, suffix)
```

The last step is to instrument your client by replacing its transport with your `Recording` instance.
`Recording` satisfies the `azcore.Transport` interface.

```go
//func recordedTestSetup(t *testing.T, testName string, mode recording.RecordMode) {
//  <...>
    // Set our client's HTTPClient to our recording instance.
    // Optionally, we can also configure MaxRetries to -1 to avoid the default retry behavior.
    client, err := NewTableServiceClient(uri, cred, &TableClientOptions{HTTPClient: recording, Retry: azcore.RetryOptions{MaxRetries: -1}})
    assert.Nil(err)

    // either return your client instance, or store it somewhere that your test can use it for test execution.
    clientsMap[testName] = &testState{client: client, recording: recording, context: &context}
}


func getTestState(key string) *testState {
    return clientsMap[key]
}
```

### Completing the recorded test session

After the test run completes we need to signal the `Recording` instance to save the recording.

```go
// recordedTestTeardown fetches the context from our map based on test name and calls Stop on the Recording instance.
func recordedTestTeardown(key string) {
    context, ok := clientsMap[key]
    if ok && !(*context.context).IsFailed() {
        context.recording.Stop()
    }
}
```

### Setting up a test to use our Recording instance

Test frameworks like testify suite allow for configuration of a `BeforeTest` method to be executed before each test.
We can use this to call our `recordedTestSetup` method

Below is an example test setup which executes a single test.

```go
package aztable

import (
    "errors"
    "fmt"
    "net/http"
    "testing"

    "github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
    "github.com/Azure/azure-sdk-for-go/sdk/internal/testframework"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type tableServiceClientLiveTests struct {
    suite.Suite
    mode         recording.RecordMode
}

// Hookup to the testing framework
func TestServiceClient_Storage(t *testing.T) {
    storage := tableServiceClientLiveTests{mode: recording.Playback /* change to Record to re-record tests */}
    suite.Run(t, &storage)
}

func (s *tableServiceClientLiveTests) TestCreateTable() {
    assert := assert.New(s.T())
    context := getTestState(s.T().Name())
    // generate a random recorded value for our table name.
    tableName, err := context.recording.GenerateAlphaNumericID(tableNamePrefix, 20, true)

    resp, err := context.client.Create(ctx, tableName)
    defer context.client.Delete(ctx, tableName)

    assert.Nil(err)
    assert.Equal(*resp.TableResponse.TableName, tableName)
}

func (s *tableServiceClientLiveTests) BeforeTest(suite string, test string) {
    // setup the test environment
    recordedTestSetup(s.T(), s.T().Name(), s.mode)
}

func (s *tableServiceClientLiveTests) AfterTest(suite string, test string) {
    // teardown the test context
    recordedTestTeardown(s.T().Name())
}
```
