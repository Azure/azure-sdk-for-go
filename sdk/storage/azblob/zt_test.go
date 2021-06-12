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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testframework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"github.com/Azure/azure-sdk-for-go/sdk/to"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	chk "gopkg.in/check.v1"
)

// For testing docs, see: https://labix.org/gocheck
// To test a specific test: go test -check.f MyTestSuite

type aztestsSuite struct {
	suite.Suite
	mode testframework.RecordMode
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &aztestsSuite{mode: testframework.Playback})
}

type testContext struct {
	recording *testframework.Recording
	context   *testframework.TestContext
}

// a map to store our created test contexts
var clientsMap = make(map[string]*testContext)

// recordedTestSetup is called before each test execution by the test suite's BeforeTest method
func recordedTestSetup(t *testing.T, accountType string, testName string, mode testframework.RecordMode) {
	_assert := assert.New(t)

	// init the test framework
	_testContext := testframework.NewTestContext(
		func(msg string) { _assert.FailNow(msg) },
		func(msg string) { t.Log(msg) },
		func() string { return testName })

	// mode should be testframework.Playback.
	//This will automatically record if no test recording is available and playback if it is.
	recording, err := testframework.NewRecording(_testContext, mode)
	_assert.Nil(err)

	_, err = recording.GetRecordedVariable(accountType+AccountNameEnvVar, testframework.Default)
	_, err = recording.GetRecordedVariable(accountType+AccountKeyEnvVar, testframework.Secret_Base64String)
	_ = recording.GetOptionalRecordedVariable(DefaultEndpointSuffixEnvVar, DefaultEndpointSuffix, testframework.Default)
	_, err = recording.GetRecordedVariable("SECONDARY_"+AccountNameEnvVar, testframework.Default)
	_, err = recording.GetRecordedVariable("SECONDARY_"+AccountKeyEnvVar, testframework.Secret_Base64String)

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

func (s *aztestsSuite) BeforeTest(suite string, test string) {
	// setup the test environment
	recordedTestSetup(s.T(), "", s.T().Name(), s.mode)
}

func (s *aztestsSuite) AfterTest(suite string, test string) {
	// teardown the test context
	recordedTestTeardown(s.T().Name())
}

// Vars for
const DefaultEndpointSuffix = "core.windows.net/"
const DefaultBlobEndpointSuffix = "blob.core.windows.net/"
const AccountNameEnvVar = "AZURE_STORAGE_ACCOUNT_NAME"
const AccountKeyEnvVar = "AZURE_STORAGE_ACCOUNT_KEY"
const DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"

const (
	containerPrefix             = "goc"
	blobPrefix                  = "gotestblob"
	blockBlobDefaultData        = "GoBlockBlobData"
	validationErrorSubstring    = "validation failed"
	invalidHeaderErrorSubstring = "invalid header field" // error thrown by the http client
)

var ctx = context.Background()

var (
	blobContentType        string = "my_type"
	blobContentDisposition string = "my_disposition"
	blobCacheControl       string = "control"
	blobContentLanguage    string = "my_language"
	blobContentEncoding    string = "my_encoding"
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

var basicBlobTagsMap = map[string]string{
	"azure": "blob",
	"blob":  "sdk",
	"sdk":   "go",
}

var specialCharBlobTagsMap = map[string]string{
	"+-./:=_ ":        "firsttag",
	"tag2":            "+-./:=_",
	"+-./:=_1":        "+-./:=_",
	"Microsoft Azure": "Azure Storage",
	"Storage+SDK":     "SDK/GO",
	"GO ":             ".Net",
}

const testPipelineMessage string = "test factory invoked"

func newTestPipeline() azcore.Pipeline {
	return azcore.NewPipeline(nil, newTestPolicy())
}

func newTestPolicy() azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		return nil, errors.New(testPipelineMessage)
	})
}

// This function generates an entity name by concatenating the passed prefix,
// the name of the test requesting the entity name, and the minute, second, and nanoseconds of the call.
// This should make it easy to associate the entities with their test, uniquely identify
// them, and determine the order in which they were created.
// Note that this imposes a restriction on the length of test names
func generateName(prefix string) string {
	// These next lines up through the for loop are obtaining and walking up the stack
	// trace to extrat the test name, which is stored in name
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
	return containerClient.NewAppendBlobURL(appendBlobName)
}

func getPageBlobClient(pageBlobName string, containerClient ContainerClient) PageBlobClient {
	return containerClient.NewPageBlobClient(pageBlobName)
}

func getReaderToRandomBytes(n int) *bytes.Reader {
	r, _ := getRandomDataAndReader(n)
	return r
}

func getRandomDataAndReader(n int) (*bytes.Reader, []byte) {
	data := make([]byte, n, n)
	rand.Read(data)
	return bytes.NewReader(data), data
}

func createNewContainer(_assert *assert.Assertions, testName string, bsu ServiceClient) (ContainerClient, string) {
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, bsu)

	cResp, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	return containerClient, containerName
}

func deleteContainer(container ContainerClient) {
	_, _ = container.Delete(context.Background(), nil)
	//c.Assert(err, chk.IsNil)
	//c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
}

func createNewContainerWithSuffix(c *chk.C, bsu ServiceClient, suffix string) (container ContainerClient, name string) {
	// The goal of adding the suffix is to be able to predetermine what order the containers will be in when listed.
	// We still need the container prefix to come first, though, to ensure only containers as a part of this test
	// are listed at all.
	name = generateName(containerPrefix + suffix)
	container = bsu.NewContainerClient(name)

	cResp, err := container.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	return container, name
}

func createNewBlockBlob(_assert *assert.Assertions, testName string, containerClient ContainerClient) (BlockBlobClient, string) {
	blockBlobName := generateBlobName(testName)
	blockBlobClient := getBlockBlobClient(blockBlobName, containerClient)

	cResp, err := blockBlobClient.Upload(ctx, azcore.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)

	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	return blockBlobClient, blockBlobName
}

func createNewBlobs(c *chk.C, container ContainerClient, blobNames []string) {
	for _, name := range blobNames {
		createNewBlockBlobWithName(c, container, name)
	}
}

func createNewBlockBlobWithName(c *chk.C, container ContainerClient, name string) (blob BlockBlobClient) {
	blob = container.NewBlockBlobClient(name)

	cResp, err := blob.Upload(ctx, strings.NewReader(blockBlobDefaultData), nil)

	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	return
}

func createNewAppendBlob(_assert *assert.Assertions, testName string, containerClient ContainerClient) (AppendBlobClient, string) {
	appendBlobName := generateBlobName(testName)
	appendBlobClient := getAppendBlobClient(appendBlobName, containerClient)

	appendBlobCreateResp, err := appendBlobClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	return appendBlobClient, appendBlobName
}

func createNewPageBlob(_assert *assert.Assertions, testName string, containerClient ContainerClient) (PageBlobClient, string) {
	pageBlobName := generateBlobName(testName)
	pageBlobClient := getPageBlobClient(pageBlobName, containerClient)

	pageBlobCreateResponse, err := pageBlobClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	return pageBlobClient, pageBlobName
}

func createNewPageBlobWithSize(_assert *assert.Assertions, testName string, containerClient ContainerClient, sizeInBytes int64) (PageBlobClient, string) {
	pageBlobName := generateBlobName(testName)
	pageBlobClient := getPageBlobClient(pageBlobName, containerClient)

	pageBlobCreateResponse, err := pageBlobClient.Create(ctx, sizeInBytes, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	return pageBlobClient, pageBlobName
}

func createNewBlockBlobWithPrefix(_assert *assert.Assertions, testName string, container ContainerClient, prefix string) (blob BlockBlobClient, name string) {
	name = prefix + generateBlobName(testName)
	blob = container.NewBlockBlobClient(name)

	blockBlobUploadResp, err := blob.Upload(ctx, strings.NewReader(blockBlobDefaultData), nil)

	_assert.Nil(err)
	_assert.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)
	return
}

func getGenericCredential(accountType string) (*SharedKeyCredential, error) {
	accountNameEnvVar := accountType + AccountNameEnvVar
	accountKeyEnvVar := accountType + AccountKeyEnvVar
	accountName, accountKey := os.Getenv(accountNameEnvVar), os.Getenv(accountKeyEnvVar)
	if accountName == "" || accountKey == "" {
		return nil, errors.New(accountNameEnvVar + " and/or " + accountKeyEnvVar + " environment variables not specified.")
	}
	return NewSharedKeyCredential(accountName, accountKey)
}

func getOAuthCredential(c *chk.C) azcore.Credential {
	cred, err := azidentity.NewEnvironmentCredential(nil)
	c.Assert(err, chk.IsNil)
	return cred
}

func getGenericServiceClientWithOAuth(c *chk.C, accountType string) (ServiceClient, error) {
	accountNameEnvVar := accountType + AccountNameEnvVar
	accountName := os.Getenv(accountNameEnvVar)
	if accountName == "" {
		return ServiceClient{}, errors.New(accountNameEnvVar + " environment variables not specified.")
	}

	blobPrimaryURL, _ := url.Parse("https://" + accountName + ".blob.core.windows.net/")
	return NewServiceClient(blobPrimaryURL.String(), getOAuthCredential(c), nil)
}

func getGenericBSU(accountType string, options *ClientOptions) (ServiceClient, error) {
	credential, err := getGenericCredential(accountType)
	if err != nil {
		return ServiceClient{}, err
	}

	blobPrimaryURL, _ := url.Parse("https://" + credential.AccountName() + ".blob.core.windows.net/")
	return NewServiceClient(blobPrimaryURL.String(), credential, options)
}

func getBSU(options *ClientOptions) ServiceClient {
	bsu, _ := getGenericBSU("", options)
	return bsu
}

func getAlternateBSU() (ServiceClient, error) {
	return getGenericBSU("SECONDARY_", nil)
}

func getPremiumBSU() (ServiceClient, error) {
	return getGenericBSU("PREMIUM_", nil)
}

func getBlobStorageBSU() (ServiceClient, error) {
	return getGenericBSU("BLOB_STORAGE_", nil)
}

func getRelativeTimeGMT(amount time.Duration) time.Time {
	currentTime := time.Now().In(time.FixedZone("GMT", 0))
	currentTime = currentTime.Add(amount * time.Second)
	return currentTime
}

func generateCurrentTimeWithModerateResolution() time.Time {
	highResolutionTime := time.Now().UTC()
	return time.Date(highResolutionTime.Year(), highResolutionTime.Month(), highResolutionTime.Day(), highResolutionTime.Hour(), highResolutionTime.Minute(),
		highResolutionTime.Second(), 0, highResolutionTime.Location())
}

// Some tests require setting service properties. It can take up to 30 seconds for the new properties to be reflected across all FEs.
// We will enable the necessary property and try to run the test implementation. If it fails with an error that should be due to
// those changes not being reflected yet, we will wait 30 seconds and try the test again. If it fails this time for any reason,
// we fail the test. It is the responsibility of the the testImplFunc to determine which error string indicates the test should be retried.
// There can only be one such string. All errors that cannot be due to this detail should be asserted and not returned as an error string.
func runTestRequiringServiceProperties(c *chk.C, bsu ServiceClient, code string,
	enableServicePropertyFunc func(*chk.C, ServiceClient),
	testImplFunc func(*chk.C, ServiceClient) error,
	disableServicePropertyFunc func(*chk.C, ServiceClient)) {
	enableServicePropertyFunc(c, bsu)
	defer disableServicePropertyFunc(c, bsu)
	err := testImplFunc(c, bsu)
	// We cannot assume that the error indicative of slow update will necessarily be a StorageError. As in ListBlobs.
	if err != nil && err.Error() == code {
		time.Sleep(time.Second * 30)
		err = testImplFunc(c, bsu)
		c.Assert(err, chk.IsNil)
	}
}

func enableSoftDelete(c *chk.C, bsu ServiceClient) {
	days := int32(1)
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: to.BoolPtr(true), Days: &days}})
	c.Assert(err, chk.IsNil)
}

func disableSoftDelete(c *chk.C, bsu ServiceClient) {
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: to.BoolPtr(false)}})
	c.Assert(err, chk.IsNil)
}

func validateUpload(c *chk.C, blobClient BlobClient) {
	resp, err := blobClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	data, _ := ioutil.ReadAll(resp.RawResponse.Body)
	c.Assert(data, chk.HasLen, 0)
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
