// +build go1.13

package armcompute

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type availabilitySetMockTest struct {
	suite.Suite
	mode recording.RecordMode
}

// Hookup to the testing framework
func TestAvailabilitySetClient(t *testing.T) {
	mockTest := availabilitySetMockTest{mode: recording.Record}
	suite.Run(t, &mockTest)
}

func (s *availabilitySetMockTest) TestCreateTable() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	// generate a random recorded value for our table name.
	tableName, err := context.recording.GenerateAlphaNumericID(tableNamePrefix, 20, true)

	resp, err := context.client.Create(ctx, tableName)
	defer context.client.Delete(ctx, tableName)

	assert.Nil(err)
	assert.Equal(*resp.TableResponse.TableName, tableName)
}

func (s *availabilitySetMockTest) BeforeTest(suite string, test string) {
	// setup the test environment
	recordedTestSetup(s.T(), s.T().Name(), s.mode)
}

func (s *availabilitySetMockTest) AfterTest(suite string, test string) {
	// teardown the test context
	recordedTestTeardown(s.T().Name())
}

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
	accountName, err := recording.GetEnvVar(storageAccountNameEnvVar, testframework.Default)
	suffix := recording.GetOptionalEnvVar(storageEndpointSuffixEnvVar, DefaultStorageSuffix, testframework.Default)
	secret, err := recording.GetEnvVar(storageAccountKeyEnvVar, testframework.Secret_Base64String)
	cred, _ := NewSharedKeyCredential(accountName, secret)
	uri := storageURI(accountName, suffix)
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

// recordedTestTeardown fetches the context from our map based on test name and calls Stop on the Recording instance.
func recordedTestTeardown(key string) {
	context, ok := clientsMap[key]
	if ok && !(*context.context).IsFailed() {
		context.recording.Stop()
	}
}
