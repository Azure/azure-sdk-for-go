package azfile

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

var ctx = context.Background()

type azfileTestSuite struct {
	suite.Suite
	mode testframework.RecordMode
}

const (
	DefaultEndpointSuffix       = "core.windows.net/"
	AccountNameEnvVar           = "AZURE_STORAGE_ACCOUNT_NAME"
	AccountKeyEnvVar            = "AZURE_STORAGE_ACCOUNT_KEY"
	DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"
)

const (
	sharePrefix                 = "gos"
	filePrefix                  = "gotestfile"
	fileDefaultData             = "GoFileDefaultData"
	invalidHeaderErrorSubstring = "invalid header field"
)

//nolint
type azfileLiveTestSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &azfileTestSuite{mode: testframework.Playback})
	//suite.Run(t, &azfileLiveTestSuite{})
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

// This function generates an entity name by concatenating the passed prefix,
// the name of the test requesting the entity name, and the minute, second, and nanoseconds of the call.
// This should make it easy to associate the entities with their test, uniquely identify
// them, and determine the order in which they were created.
// Note that this imposes a restriction on the length of test names
//nolint
func generateName(prefix string) string {
	// These next lines up through the for loop are obtaining and walking up the stack
	// trace to extract the test name, which is stored in name
	pc := make([]uintptr, 10)
	runtime.Callers(0, pc)
	frames := runtime.CallersFrames(pc)
	name := ""
	for f, next := frames.Next(); next; f, next = frames.Next() {
		name = f.Function
		if strings.Contains(name, "Suite") {
			break
		}
	}
	funcNameStart := strings.Index(name, "Test")
	name = name[funcNameStart+len("Test"):] // Just get the name of the test and not any of the garbage at the beginning
	name = strings.ToLower(name)            // Ensure it is a valid resource name
	currentTime := time.Now()
	name = fmt.Sprintf("%s%s%d%d%d", prefix, strings.ToLower(name), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond())
	return name
}

func generateEntityName(testName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(testName), "/", ""), "test", "")
}
func generateShareName(testName string) string {
	return sharePrefix + generateEntityName(testName)
}

func generateFileName(testName string) string {
	return filePrefix + generateEntityName(testName)
}

// getRequiredEnv gets an environment variable by name and returns an error if it is not found
func getRequiredEnv(name string) (string, error) {
	env, ok := os.LookupEnv(name)
	if ok {
		return env, nil
	} else {
		return "", errors.New("Required environment variable not set: " + name)
	}
}
func getAccountInfo(recording *testframework.Recording, accountType testAccountType) (string, string, error) {
	accountNameEnvVar := string(accountType) + AccountNameEnvVar
	accountKeyEnvVar := string(accountType) + AccountKeyEnvVar
	accountName, accountKey := "", ""
	var err error
	if recording == nil {
		accountName, err = getRequiredEnv(accountNameEnvVar)
		if err != nil {
			return "", "", err
		}
		accountKey, err = getRequiredEnv(accountKeyEnvVar)
		if err != nil {
			return "", "", err
		}

	} else {
		accountName, err = recording.GetEnvVar(accountNameEnvVar, testframework.NoSanitization)
		if err != nil {
			return "", "", err
		}
		accountKey, err = recording.GetEnvVar(accountKeyEnvVar, testframework.Secret_Base64String)
		if err != nil {
			return "", "", err
		}
	}
	return accountName, accountKey, nil
}
func getGenericCredential(recording *testframework.Recording, accountType testAccountType) (*SharedKeyCredential, error) {
	accountName, accountKey, err := getAccountInfo(recording, accountType)
	if err != nil {
		return nil, err
	} else if accountName == "" || accountKey == "" {
		return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
	}
	return NewSharedKeyCredential(accountName, accountKey)
}

//nolint
func getConnectionString(recording *testframework.Recording, accountType testAccountType) (string, error) {
	accountName, accountKey, err := getAccountInfo(recording, accountType)
	if err != nil {
		return "", err
	}
	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/",
		accountName, accountKey)
	return connectionString, nil
}

//nolint
func getServiceClientFromConnectionString(recording *testframework.Recording, accountType testAccountType, options *ClientOptions) (ServiceClient, error) {
	if recording != nil {
		if options == nil {
			options = &ClientOptions{
				Transporter: recording,
				Retry:       policy.RetryOptions{MaxRetries: -1}}
		}
	}

	connectionString, err := getConnectionString(recording, accountType)
	if err != nil {
		return ServiceClient{}, nil
	}
	primaryURL, cred, err := parseConnectionString(connectionString)
	if err != nil {
		return ServiceClient{}, nil
	}

	svcClient, err := NewServiceClient(primaryURL, cred, options)
	return svcClient, err
}

type testAccountType string

const (
	testAccountDefault   testAccountType = ""
	testAccountSecondary testAccountType = "SECONDARY_"
	testAccountPremium   testAccountType = "PREMIUM_"
)

func getServiceClient(recording *testframework.Recording, accountType testAccountType, options *ClientOptions) (ServiceClient, error) {
	if recording != nil {
		if options == nil {
			options = &ClientOptions{
				Transporter: recording,
				Retry:       policy.RetryOptions{MaxRetries: -1}}
		}
	}

	cred, err := getGenericCredential(recording, accountType)
	if err != nil {
		return ServiceClient{}, err
	}

	serviceURL, _ := url.Parse("https://" + cred.AccountName() + ".file.core.windows.net/")
	serviceClient, err := NewServiceClient(serviceURL.String(), cred, options)

	return serviceClient, err
}

//nolint
func getRelativeTimeGMT(amount time.Duration) time.Time {
	currentTime := time.Now().In(time.FixedZone("GMT", 0))
	currentTime = currentTime.Add(amount * time.Second)
	return currentTime
}

func getRelativeTimeFromAnchor(anchorTime *time.Time, amount time.Duration) time.Time {
	return anchorTime.Add(amount * time.Second)
}

func getReaderToGeneratedBytes(n int) io.ReadSeekCloser {
	r, _ := generateData(n)
	return internal.NopCloser(r)
}

//nolint
func getRandomDataAndReader(n int) (*bytes.Reader, []byte) {
	data := make([]byte, n)
	rand.Read(data)
	return bytes.NewReader(data), data
}

const random64BString string = "2SDgZj6RkKYzJpu04sweQek4uWHO8ndPnYlZ0tnFS61hjnFZ5IkvIGGY44eKABov"

func generateData(sizeInBytes int) (io.ReadSeekCloser, []byte) {
	data := make([]byte, sizeInBytes)
	_len := len(random64BString)
	if sizeInBytes > _len {
		count := sizeInBytes / _len
		if sizeInBytes%_len != 0 {
			count = count + 1
		}
		copy(data[:], strings.Repeat(random64BString, count))
	} else {
		copy(data[:], random64BString)
	}
	return internal.NopCloser(bytes.NewReader(data)), data
}

func getShareClient(containerName string, s ServiceClient) ShareClient {
	return s.NewShareClient(containerName)
}

func createNewShare(_assert *assert.Assertions, shareName string, serviceClient ServiceClient) ShareClient {
	shareClient := getShareClient(shareName, serviceClient)

	cResp, err := shareClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	return shareClient
}

func deleteContainer(_assert *assert.Assertions, shareClient ShareClient) {
	deleteContainerResp, err := shareClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(deleteContainerResp.RawResponse.StatusCode, 202)
}
