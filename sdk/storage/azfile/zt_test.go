package azfile

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

// Constants & Variables declaration -----------------------------------------------------------------------------------

var ctx = context.Background()

var (
	basicHeaders = FileHTTPHeaders{
		FileContentType:        to.StringPtr("my_type"),
		FileContentDisposition: to.StringPtr("my_disposition"),
		FileCacheControl:       to.StringPtr("control"),
		FileContentMD5:         nil,
		FileContentLanguage:    to.StringPtr("my_language"),
		FileContentEncoding:    to.StringPtr("my_encoding"),
	}

	basicMetadata = map[string]string{"foo": "bar"}

	sampleSDDL = `O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)`

	deleteSnapshotsInclude = DeleteSnapshotsOptionTypeInclude
)

const (
	DefaultEndpointSuffix       = "core.windows.net/"
	AccountNameEnvVar           = "AZURE_STORAGE_ACCOUNT_NAME"
	AccountKeyEnvVar            = "AZURE_STORAGE_ACCOUNT_KEY"
	DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"
)

const (
	sharePrefix                 = "azftestshare"
	directoryPrefix             = "azftestdir"
	filePrefix                  = "azftestfile"
	fileDefaultData             = "GoFileDefaultData"
	invalidHeaderErrorSubstring = "invalid header field"
)

// Testsuite -----------------------------------------------------------------------------------------------------------

type azfileTestSuite struct {
	suite.Suite
	mode testframework.RecordMode
}

//nolint
type azfileLiveTestSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &azfileTestSuite{mode: testframework.Playback})
	suite.Run(t, &azfileLiveTestSuite{})
}

type testContext struct {
	recording *testframework.Recording
	context   *testframework.TestContext
}

// a map to store our created test contexts
var clientsMap = make(map[string]*testContext)

// recordedTestSetup is called before each test execution by the test suite's BeforeTest method
func recordedTestSetup(t *testing.T, mode testframework.RecordMode) {
	testName := t.Name()
	_assert := assert.New(t)

	// init the test framework
	_testContext := testframework.NewTestContext(
		func(msg string) { _assert.FailNow(msg) },
		func(msg string) { t.Log(msg) },
		func() string { return testName })

	// mode should be test_framework.Playback.
	// This will automatically record if no test recording is available and playback if it is.
	recording, err := testframework.NewRecording(_testContext, mode)
	_assert.Nil(err)

	_, err = recording.GetEnvVar(AccountNameEnvVar, testframework.NoSanitization)
	if err != nil {
		log.Fatal(err)
	}
	_, err = recording.GetEnvVar(AccountKeyEnvVar, testframework.Secret_Base64String)
	if err != nil {
		log.Fatal(err)
	}
	_ = recording.GetOptionalEnvVar(DefaultEndpointSuffixEnvVar, DefaultEndpointSuffix, testframework.NoSanitization)

	clientsMap[testName] = &testContext{recording: recording, context: &_testContext}
}

func getTestContext(key string) *testContext {
	return clientsMap[key]
}

func recordedTestTeardown(key string) {
	_context, ok := clientsMap[key]
	if ok && !(*_context.context).IsFailed() {
		_ = _context.recording.Stop()
	}
}

//nolint
func (s *azfileTestSuite) BeforeTest(suite string, test string) {
	// set up the test environment
	recordedTestSetup(s.T(), s.mode)
}

//nolint
func (s *azfileTestSuite) AfterTest(suite string, test string) {
	// teardown the test context
	recordedTestTeardown(s.T().Name())
}

//nolint
func (s *azfileLiveTestSuite) BeforeTest(suite string, test string) {
}

//nolint
func (s *azfileLiveTestSuite) AfterTest(suite string, test string) {

}

//----------------------------------------------------------------------------------------------------------------------

func validateStorageError(_assert *assert.Assertions, err error, code StorageErrorCode) {
	_assert.NotNil(err)
	var storageError *StorageError
	_assert.Equal(errors.As(err, &storageError), true)

	_assert.Equal(storageError.ErrorCode, code)
}
