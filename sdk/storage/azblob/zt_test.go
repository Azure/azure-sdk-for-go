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
	"io"
	"io/ioutil"
	"math/rand"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
)

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

// This function generates an entity name by concatenating the passed prefix,
// the name of the test requesting the entity name, and the minute, second, and nanoseconds of the call.
// This should make it easy to associate the entities with their test, uniquely identify
// them, and determine the order in which they were created.
// Note that this imposes a restriction on the length of test names
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

func createNewContainer(t *testing.T, containerName string, serviceClient ServiceClient) ContainerClient {
	containerClient := getContainerClient(containerName, serviceClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	return containerClient
}

func deleteContainer(t *testing.T, containerClient ContainerClient) {
	deleteContainerResp, err := containerClient.Delete(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, deleteContainerResp.RawResponse.StatusCode, 202)
}

func createNewBlockBlob(t *testing.T, blockBlobName string, containerClient ContainerClient) BlockBlobClient {
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)

	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	return bbClient
}

func createNewBlobs(t *testing.T, blobNames []string, containerClient ContainerClient) {
	for _, blobName := range blobNames {
		createNewBlockBlob(t, blobName, containerClient)
	}
}

func createNewAppendBlob(t *testing.T, appendBlobName string, containerClient ContainerClient) AppendBlobClient {
	abClient := getAppendBlobClient(appendBlobName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	return abClient
}

func createNewPageBlob(t *testing.T, pageBlobName string, containerClient ContainerClient) PageBlobClient {
	return createNewPageBlobWithSize(t, pageBlobName, containerClient, PageBlobPageBytes*10)
}

func createNewPageBlobWithSize(t *testing.T, pageBlobName string,
	containerClient ContainerClient, sizeInBytes int64) PageBlobClient {
	pbClient := getPageBlobClient(pageBlobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, sizeInBytes, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	return pbClient
}

func createNewBlockBlobWithCPK(t *testing.T, blockBlobName string, containerClient ContainerClient, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (bbClient BlockBlobClient) {
	bbClient = getBlockBlobClient(blockBlobName, containerClient)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkInfo:      cpkInfo,
		CpkScopeInfo: cpkScopeInfo,
	}
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &uploadBlockBlobOptions)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	require.Equal(t, *cResp.IsServerEncrypted, true)
	if cpkInfo != nil {
		require.EqualValues(t, cResp.EncryptionKeySHA256, cpkInfo.EncryptionKeySHA256)
	}
	if cpkScopeInfo != nil {
		require.EqualValues(t, cResp.EncryptionScope, cpkScopeInfo.EncryptionScope)
	}
	return
}

func createNewPageBlobWithCPK(t *testing.T, pageBlobName string, container ContainerClient,
	sizeInBytes int64, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (pbClient PageBlobClient) {
	pbClient = getPageBlobClient(pageBlobName, container)

	resp, err := pbClient.Create(ctx, sizeInBytes, &CreatePageBlobOptions{
		CpkInfo:      cpkInfo,
		CpkScopeInfo: cpkScopeInfo,
	})
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)
	return
}

type testAccountType string

const (
	testAccountDefault   testAccountType = ""
	testAccountSecondary testAccountType = "SECONDARY_"
	testAccountPremium   testAccountType = "PREMIUM_"
)

func getRelativeTimeGMT(amount time.Duration) time.Time {
	currentTime := time.Now().In(time.FixedZone("GMT", 0))
	currentTime = currentTime.Add(amount * time.Second)
	return currentTime
}

func getRelativeTimeFromAnchor(anchorTime *time.Time, amount time.Duration) time.Time {
	return anchorTime.Add(amount * time.Second)
}

// Some tests require setting service properties. It can take up to 30 seconds for the new properties to be reflected across all FEs.
// We will enable the necessary property and try to run the test implementation. If it fails with an error that should be due to
// those changes not being reflected yet, we will wait 30 seconds and try the test again. If it fails this time for any reason,
// we fail the test. It is the responsibility of the the testImplFunc to determine which error string indicates the test should be retried.
// There can only be one such string. All errors that cannot be due to this detail should be asserted and not returned as an error string.
func runTestRequiringServiceProperties(t *testing.T, bsu ServiceClient, code string,
	enableServicePropertyFunc func(*testing.T, ServiceClient),
	testImplFunc func(*testing.T, ServiceClient) error,
	disableServicePropertyFunc func(*testing.T, ServiceClient)) {

	enableServicePropertyFunc(t, bsu)
	defer disableServicePropertyFunc(t, bsu)

	err := testImplFunc(t, bsu)
	// We cannot assume that the error indicative of slow update will necessarily be a StorageError. As in ListBlobs.
	if err != nil && err.Error() == code {
		recording.Sleep(time.Second * 30)
		err = testImplFunc(t, bsu)
		require.NoError(t, err)
	}
}

func enableSoftDelete(t *testing.T, serviceClient ServiceClient) {
	_, err := serviceClient.SetProperties(ctx, StorageServiceProperties{
		DeleteRetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(1),
		},
	})
	require.NoError(t, err)
}

func disableSoftDelete(t *testing.T, bsu ServiceClient) {
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: to.BoolPtr(false)}})
	require.NoError(t, err)
}

func validateUpload(t *testing.T, blobClient BlobClient) {
	resp, err := blobClient.Download(ctx, nil)
	require.NoError(t, err)
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(t, err)
	require.Len(t, data, 0)
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
func validateStorageError(t *testing.T, err error, code StorageErrorCode) {
	require.Error(t, err)
	var storageError *StorageError

	require.Equal(t, true, errors.As(err, &storageError))
	require.Equal(t, storageError.ErrorCode, code)
}

func blobListToMap(list []string) map[string]bool {
	out := make(map[string]bool)

	for _, v := range list {
		out[v] = true
	}

	return out
}
