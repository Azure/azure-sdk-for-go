// +build go1.13

package armcompute

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	mockSubscriptionID = "00000000-0000-0000-0000-000000000000"
	location           = "westus"
	groupNamePrefix    = "rgMock"
	setNamePrefix      = "setMock"
)

var (
	ctx = context.Background()
)

type availabilitySetMockTest struct {
	suite.Suite
	mode recording.RecordMode
}

// Hookup to the testing framework
func TestAvailabilitySetClient(t *testing.T) {
	mockTest := availabilitySetMockTest{mode: recording.Playback}
	suite.Run(t, &mockTest)
}

func (s *availabilitySetMockTest) TestAvailabilitySetCRUD() {
	assertion := assert.New(s.T())
	state := getTestState(s.T().Name())
	// generate a random recorded value for our table name.
	groupName, _ := state.recording.GenerateAlphaNumericID(groupNamePrefix, 10, true)
	accountName, _ := state.recording.GenerateAlphaNumericID(setNamePrefix, 10, true)

	// test create
	parameters := AvailabilitySet{
		Resource: Resource{
			Location: to.StringPtr(location),
		},
	}
	createResp, err := state.client.CreateOrUpdate(ctx, groupName, accountName, parameters, nil)
	assertion.NoError(err)
	assertion.NotNil(createResp.AvailabilitySet)
	assertion.Equal(*createResp.AvailabilitySet.Name, accountName)

	// test get
	getResp, err := state.client.Get(ctx, groupName, accountName, nil)
	assertion.NoError(err)
	assertion.NotNil(getResp.AvailabilitySet)
	assertion.Equal(*getResp.AvailabilitySet.Name, accountName)
	//assertion.Equal(*getResp.AvailabilitySet.Location, location) // mock server is making this up, therefore this check cannot pass

	// test delete
	deleteResp, err := state.client.Delete(ctx, groupName, accountName, nil)
	assertion.NoError(err)
	assertion.NotNil(deleteResp)

	// test get (and get an error)
	getFailResp, err := state.client.Get(ctx, groupName, accountName, nil)
	assertion.Error(err)
	rawResponse := err.(azcore.HTTPResponse).RawResponse()
	assertion.Equal(http.StatusNotFound, rawResponse.StatusCode)
	assertion.Nil(getFailResp.AvailabilitySet)
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
	client    *AvailabilitySetsClient
	context   *recording.TestContext
}

// a map to store our created test contexts
var clientsMap = make(map[string]*testState)

// recordedTestSetup is called before each test execution by the test suite's BeforeTest method
func recordedTestSetup(t *testing.T, testName string, mode recording.RecordMode) {
	assertion := assert.New(t)

	// init the test framework
	testContext := recording.NewTestContext(func(msg string) { assertion.FailNow(msg) }, func(msg string) { t.Log(msg) }, func() string { return testName })
	//mode should be recording.Playback. This will automatically record if no test recording is available and playback if it is.
	record, err := recording.NewRecording(testContext, mode)
	assertion.Nil(err)
	// Set our client's HTTPClient to our recording instance.
	// Optionally, we can also configure MaxRetries to -1 to avoid the default retry behavior.
	conn := armcore.NewConnection("https://localhost:8441/", &mockTokenCred{}, &armcore.ConnectionOptions{
		HTTPClient: record,
		Retry: azcore.RetryOptions{MaxRetries: -1},
	})
	client := NewAvailabilitySetsClient(conn, mockSubscriptionID)
	assertion.Nil(err)

	// either return your client instance, or store it somewhere that your test can use it for test execution.
	clientsMap[testName] = &testState{client: client, recording: record, context: &testContext}
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

// The below part might need to migrate to a more common place to avoid to be duplicated everywhere
type mockTokenCred struct{}

func (mockTokenCred) AuthenticationPolicy(azcore.AuthenticationPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		return req.Next()
	})
}

func (mockTokenCred) GetToken(context.Context, azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		Token:     "abc123",
		ExpiresOn: time.Now().Add(1 * time.Hour),
	}, nil
}
