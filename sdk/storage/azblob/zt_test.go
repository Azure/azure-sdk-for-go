// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

type azblobTestSuite struct {
	suite.Suite
	mode testframework.RecordMode
}

//nolint
type azblobUnrecordedTestSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &azblobTestSuite{mode: testframework.Playback})
	//suite.Run(t, &azblobUnrecordedTestSuite{})
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
func (s *azblobTestSuite) BeforeTest(suite string, test string) {
	// set up the test environment
	recordedTestSetup(s.T(), s.mode)
}

//nolint
func (s *azblobTestSuite) AfterTest(suite string, test string) {
	// teardown the test context
	recordedTestTeardown(s.T().Name())
}

//nolint
func (s *azblobUnrecordedTestSuite) BeforeTest(suite string, test string) {

}

//nolint
func (s *azblobUnrecordedTestSuite) AfterTest(suite string, test string) {

}

//// Note that this function is adding to the list of ignored headers not creating from scratch
//func ignoreHeaders(recording *testframework.Recording, headers []string) {
//	ignoredHeaders := recording.Matcher.IgnoredHeaders
//	for _, header := range headers {
//		// TODO: Mohit Come Here
//		//ignoredHeaders[header] = nil
//		_ = header
//	}
//	_ = ignoredHeaders
//}

// Vars for
const DefaultEndpointSuffix = "core.windows.net/"

//const DefaultBlobEndpointSuffix = "blob.core.windows.net/"
const AccountNameEnvVar = "AZURE_STORAGE_ACCOUNT_NAME"
const AccountKeyEnvVar = "AZURE_STORAGE_ACCOUNT_KEY"
const DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"

const (
	containerPrefix             = "goc"
	blobPrefix                  = "gotestblob"
	blockBlobDefaultData        = "GoBlockBlobData"
	invalidHeaderErrorSubstring = "invalid header field" // error thrown by the http client
)

var ctx = context.Background()

var (
	blobContentType        = "my_type"
	blobContentDisposition = "my_disposition"
	blobCacheControl       = "control"
	blobContentLanguage    = "my_language"
	blobContentEncoding    = "my_encoding"
)

var basicHeaders = BlobHTTPHeaders{
	BlobContentType:        &blobContentType,
	BlobContentDisposition: &blobContentDisposition,
	BlobCacheControl:       &blobCacheControl,
	BlobContentMD5:         nil,
	BlobContentLanguage:    &blobContentLanguage,
	BlobContentEncoding:    &blobContentEncoding,
}

var basicMetadata = map[string]string{"Foo": "bar"}

//nolint
var basicBlobTagsMap = map[string]string{
	"azure": "blob",
	"blob":  "sdk",
	"sdk":   "go",
}

//nolint
var specialCharBlobTagsMap = map[string]string{
	"+-./:=_ ":        "firsttag",
	"tag2":            "+-./:=_",
	"+-./:=_1":        "+-./:=_",
	"Microsoft Azure": "Azure Storage",
	"Storage+SDK":     "SDK/GO",
	"GO ":             ".Net",
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
func generateContainerName(testName string) string {
	return containerPrefix + generateEntityName(testName)
}

func generateBlobName(testName string) string {
	return blobPrefix + generateEntityName(testName)
}

func getContainerClient(containerName string, s ServiceClient) ContainerClient {
	return s.NewContainerClient(containerName)
}

func getBlockBlobClient(blockBlobName string, containerClient ContainerClient) BlockBlobClient {
	return containerClient.NewBlockBlobClient(blockBlobName)
}

func getAppendBlobClient(appendBlobName string, containerClient ContainerClient) AppendBlobClient {
	return containerClient.NewAppendBlobClient(appendBlobName)
}

func getPageBlobClient(pageBlobName string, containerClient ContainerClient) PageBlobClient {
	return containerClient.NewPageBlobClient(pageBlobName)
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

func createNewContainer(_assert *assert.Assertions, containerName string, serviceClient ServiceClient) ContainerClient {
	containerClient := getContainerClient(containerName, serviceClient)

	cResp, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	return containerClient
}

func deleteContainer(_assert *assert.Assertions, containerClient ContainerClient) {
	deleteContainerResp, err := containerClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(deleteContainerResp.RawResponse.StatusCode, 202)
}

func createNewBlockBlob(_assert *assert.Assertions, blockBlobName string, containerClient ContainerClient) BlockBlobClient {
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)

	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	return bbClient
}

func createNewBlobs(_assert *assert.Assertions, blobNames []string, containerClient ContainerClient) {
	for _, blobName := range blobNames {
		createNewBlockBlob(_assert, blobName, containerClient)
	}
}

func createNewAppendBlob(_assert *assert.Assertions, appendBlobName string, containerClient ContainerClient) AppendBlobClient {
	abClient := getAppendBlobClient(appendBlobName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	return abClient
}

func createNewPageBlob(_assert *assert.Assertions, pageBlobName string, containerClient ContainerClient) PageBlobClient {
	return createNewPageBlobWithSize(_assert, pageBlobName, containerClient, PageBlobPageBytes*10)
}

func createNewPageBlobWithSize(_assert *assert.Assertions, pageBlobName string,
	containerClient ContainerClient, sizeInBytes int64) PageBlobClient {
	pbClient := getPageBlobClient(pageBlobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, sizeInBytes, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	return pbClient
}

func createNewBlockBlobWithCPK(_assert *assert.Assertions, blockBlobName string, containerClient ContainerClient, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (bbClient BlockBlobClient) {
	bbClient = getBlockBlobClient(blockBlobName, containerClient)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkInfo:      cpkInfo,
		CpkScopeInfo: cpkScopeInfo,
	}
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &uploadBlockBlobOptions)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	_assert.Equal(*cResp.IsServerEncrypted, true)
	if cpkInfo != nil {
		_assert.EqualValues(cResp.EncryptionKeySHA256, cpkInfo.EncryptionKeySHA256)
	}
	if cpkScopeInfo != nil {
		_assert.EqualValues(cResp.EncryptionScope, cpkScopeInfo.EncryptionScope)
	}
	return
}

func createNewPageBlobWithCPK(_assert *assert.Assertions, pageBlobName string, container ContainerClient,
	sizeInBytes int64, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (pbClient PageBlobClient) {
	pbClient = getPageBlobClient(pageBlobName, container)

	resp, err := pbClient.Create(ctx, sizeInBytes, &CreatePageBlobOptions{
		CpkInfo:      cpkInfo,
		CpkScopeInfo: cpkScopeInfo,
	})
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)
	return
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
func getAccountInfo(recording *testframework.Recording, accountType testAccountType) (string, string) {
	accountNameEnvVar := string(accountType) + AccountNameEnvVar
	accountKeyEnvVar := string(accountType) + AccountKeyEnvVar
	var err error
	accountName, accountKey := "", ""
	if recording == nil {
		accountName, err = getRequiredEnv(accountNameEnvVar)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		_ = err
		accountKey, err = getRequiredEnv(accountKeyEnvVar)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		_ = err
	} else {
		accountName, err = recording.GetEnvVar(accountNameEnvVar, testframework.NoSanitization)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		_ = err
		accountKey, err = recording.GetEnvVar(accountKeyEnvVar, testframework.Secret_Base64String)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		_ = err
	}
	return accountName, accountKey
}
func getGenericCredential(recording *testframework.Recording, accountType testAccountType) (*SharedKeyCredential, error) {
	accountName, accountKey := getAccountInfo(recording, accountType)
	if accountName == "" || accountKey == "" {
		return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
	}
	return NewSharedKeyCredential(accountName, accountKey)
}

//nolint
func getConnectionString(recording *testframework.Recording, accountType testAccountType) string {
	accountName, accountKey := getAccountInfo(recording, accountType)
	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/",
		accountName, accountKey)
	return connectionString
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

	connectionString := getConnectionString(recording, accountType)
	primaryURL, cred, err := parseConnectionString(connectionString)
	if err != nil {
		return ServiceClient{}, nil
	}

	svcClient, err := NewServiceClientWithSharedKey(primaryURL, cred, options)
	return svcClient, err
}

type testAccountType string

const (
	testAccountDefault   testAccountType = ""
	testAccountSecondary testAccountType = "SECONDARY_"
	testAccountPremium   testAccountType = "PREMIUM_"
	//testAccountBlobStorage testAccountType = "BLOB_"
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

	serviceURL, _ := url.Parse("https://" + cred.AccountName() + ".blob.core.windows.net/")
	serviceClient, err := NewServiceClientWithSharedKey(serviceURL.String(), cred, options)

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

//func generateCurrentTimeWithModerateResolution() time.Time {
//	highResolutionTime := time.Now().UTC()
//	return time.Date(highResolutionTime.Year(), highResolutionTime.Month(), highResolutionTime.Day(), highResolutionTime.Hour(), highResolutionTime.Minute(),
//		highResolutionTime.Second(), 0, highResolutionTime.Location())
//}

// Some tests require setting service properties. It can take up to 30 seconds for the new properties to be reflected across all FEs.
// We will enable the necessary property and try to run the test implementation. If it fails with an error that should be due to
// those changes not being reflected yet, we will wait 30 seconds and try the test again. If it fails this time for any reason,
// we fail the test. It is the responsibility of the the testImplFunc to determine which error string indicates the test should be retried.
// There can only be one such string. All errors that cannot be due to this detail should be asserted and not returned as an error string.
func runTestRequiringServiceProperties(_assert *assert.Assertions, bsu ServiceClient, code string,
	enableServicePropertyFunc func(*assert.Assertions, ServiceClient),
	testImplFunc func(*assert.Assertions, ServiceClient) error,
	disableServicePropertyFunc func(*assert.Assertions, ServiceClient)) {

	enableServicePropertyFunc(_assert, bsu)
	defer disableServicePropertyFunc(_assert, bsu)

	err := testImplFunc(_assert, bsu)
	// We cannot assume that the error indicative of slow update will necessarily be a StorageError. As in ListBlobs.
	if err != nil && err.Error() == code {
		time.Sleep(time.Second * 30)
		err = testImplFunc(_assert, bsu)
		_assert.Nil(err)
	}
}

func enableSoftDelete(_assert *assert.Assertions, serviceClient ServiceClient) {
	days := int32(1)
	_, err := serviceClient.SetProperties(ctx, StorageServiceProperties{
		DeleteRetentionPolicy: &RetentionPolicy{Enabled: to.BoolPtr(true), Days: &days}})
	_assert.Nil(err)
}

func disableSoftDelete(_assert *assert.Assertions, bsu ServiceClient) {
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: to.BoolPtr(false)}})
	_assert.Nil(err)
}

func validateUpload(_assert *assert.Assertions, blobClient BlobClient) {
	resp, err := blobClient.Download(ctx, nil)
	_assert.Nil(err)
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Len(data, 0)
}

func generateBlockIDsList(count int) []string {
	blockIDs := make([]string, count)
	for i := 0; i < count; i++ {
		blockIDs[i] = blockIDIntToBase64(i)
	}
	return blockIDs
}

// blockIDIntToBase64 functions convert an int block ID to a base-64 string and vice versa
func blockIDIntToBase64(blockID int) string {
	binaryBlockID := (&[4]byte{})[:]
	binary.LittleEndian.PutUint32(binaryBlockID, uint32(blockID))
	return base64.StdEncoding.EncodeToString(binaryBlockID)
}

// TODO: Figure out in which scenario, the parsing will fail.
func validateStorageError(_assert *assert.Assertions, err error, code StorageErrorCode) {
	_assert.NotNil(err)
	var storageError *StorageError
	_assert.Equal(errors.As(err, &storageError), true)

	_assert.Equal(storageError.ErrorCode, code)
}

func blobListToMap(list []string) map[string]bool {
	out := make(map[string]bool)

	for _, v := range list {
		out[v] = true
	}

	return out
}
