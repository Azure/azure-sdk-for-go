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

// Hookup to the testing framework
func Test(t *testing.T) { chk.TestingT(t) }

type aztestsSuite struct{}

var _ = chk.Suite(&aztestsSuite{})

//func (s *aztestsSuite) TestRetryPolicyRetryReadsFromSecondaryHostField(c *chk.C) {
//	_, found := reflect.TypeOf(RetryOptions{}).FieldByName("RetryReadsFromSecondaryHost")
//	if !found {
//		// Make sure the RetryOption was not erroneously overwritten
//		c.Fatal("RetryOption's RetryReadsFromSecondaryHost field must exist in the Blob SDK - uncomment it and make sure the field is returned from the retryReadsFromSecondaryHost() method too!")
//	}
//}

const (
	containerPrefix             = "go"
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

func generateContainerName() string {
	return generateName(containerPrefix)
}

func generateBlobName() string {
	return generateName(blobPrefix)
}

func getContainerClient(c *chk.C, s ServiceClient) (container ContainerClient, name string) {
	name = generateContainerName()
	container = s.NewContainerClient(name)

	return container, name
}

func getBlockBlobClient(c *chk.C, container ContainerClient) (blob BlockBlobClient, name string) {
	name = generateBlobName()
	blob = container.NewBlockBlobClient(name)

	return blob, name
}

func getAppendBlobClient(c *chk.C, container ContainerClient) (blob AppendBlobClient, name string) {
	name = generateBlobName()
	blob = container.NewAppendBlobURL(name)

	return blob, name
}

func getPageBlobClient(c *chk.C, container ContainerClient) (blob PageBlobClient, name string) {
	name = generateBlobName()
	blob = container.NewPageBlobClient(name)

	return
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

func createNewContainer(c *chk.C, bsu ServiceClient) (container ContainerClient, name string) {
	container, name = getContainerClient(c, bsu)

	cResp, err := container.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	return container, name
}

func deleteContainer(c *chk.C, container ContainerClient) {
	resp, err := container.Delete(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
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

func createNewBlockBlob(c *chk.C, container ContainerClient) (blob BlockBlobClient, name string) {
	blob, name = getBlockBlobClient(c, container)

	cResp, err := blob.Upload(ctx, azcore.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)

	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)

	return
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

func createNewAppendBlob(c *chk.C, container ContainerClient) (blob AppendBlobClient, name string) {
	blob, name = getAppendBlobClient(c, container)

	resp, err := blob.Create(ctx, nil)

	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	return
}

func createNewPageBlob(c *chk.C, container ContainerClient) (blob PageBlobClient, name string) {
	blob, name = getPageBlobClient(c, container)

	resp, err := blob.Create(ctx, PageBlobPageBytes*10, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	return
}

func createNewPageBlobWithSize(c *chk.C, container ContainerClient, sizeInBytes int64) (blob PageBlobClient, name string) {
	blob, name = getPageBlobClient(c, container)

	resp, err := blob.Create(ctx, sizeInBytes, nil)

	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	return
}

func createNewBlockBlobWithPrefix(c *chk.C, container ContainerClient, prefix string) (blob BlockBlobClient, name string) {
	name = prefix + generateName(blobPrefix)
	blob = container.NewBlockBlobClient(name)

	cResp, err := blob.Upload(ctx, strings.NewReader(blockBlobDefaultData), nil)

	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	return
}

func getGenericCredential(accountType string) (*SharedKeyCredential, error) {
	accountNameEnvVar := accountType + "AZURE_STORAGE_ACCOUNT_NAME"
	accountKeyEnvVar := accountType + "AZURE_STORAGE_ACCOUNT_KEY"
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
	accountNameEnvVar := accountType + "AZURE_STORAGE_ACCOUNT_NAME"
	accountName := os.Getenv(accountNameEnvVar)
	if accountName == "" {
		return ServiceClient{}, errors.New(accountNameEnvVar + " environment variables not specified.")
	}

	blobPrimaryURL, _ := url.Parse("https://" + accountName + ".blob.core.windows.net/")
	return NewServiceClient(blobPrimaryURL.String(), getOAuthCredential(c), nil)
}

func getGenericBSU(accountType string) (ServiceClient, error) {
	credential, err := getGenericCredential(accountType)
	if err != nil {
		return ServiceClient{}, err
	}

	blobPrimaryURL, _ := url.Parse("https://" + credential.AccountName() + ".blob.core.windows.net/")
	return NewServiceClient(blobPrimaryURL.String(), credential, nil)
}

func getBSU() ServiceClient {
	bsu, _ := getGenericBSU("")
	return bsu
}

func getBSUFromConnectionString() ServiceClient {
	connectionString := getConnectionString()
	primaryURL, _, cred, err := ParseConnectionString(connectionString, "")
	if err != nil {
		return ServiceClient{}
	}

	svcClient, _ := NewServiceClient(primaryURL, cred, nil)
	return svcClient
}

func getAlternateBSU() (ServiceClient, error) {
	return getGenericBSU("SECONDARY_")
}

func getPremiumBSU() (ServiceClient, error) {
	return getGenericBSU("PREMIUM_")
}

func getBlobStorageBSU() (ServiceClient, error) {
	return getGenericBSU("BLOB_STORAGE_")
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
func validateStorageError(c *chk.C, err error, code StorageErrorCode) {
	c.Assert(err, chk.NotNil)
	var storageError *StorageError
	c.Assert(errors.As(err, &storageError), chk.Equals, true)

	c.Assert(storageError.ErrorCode, chk.Equals, code)
}

func blobListToMap(list []string) map[string]bool {
	out := make(map[string]bool)

	for _, v := range list {
		out[v] = true
	}

	return out
}

func getConnectionString() string {
	accountName, accountKey := accountInfo()
	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/",
		accountName, accountKey)
	return connectionString
}
