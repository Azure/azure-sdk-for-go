// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClient() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	blobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//
//	blobURLParts := NewBlobURLParts(bbClient.URL())
//	_assert.Equal(blobURLParts.BlobName, blobName)
//	_assert.Equal(blobURLParts.ContainerName, containerName)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_assert.NoError(err)
//	correctURL := "https://" + accountName + "." + DefaultBlobEndpointSuffix + containerName + "/" + blobName
//	_assert.Equal(bbClient.URL(), correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClientWithSnapshotAndSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	_assert.NoError(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	blobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_assert.NoError(err)
//
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	if err != nil {
//		s.Fail(err.Error())
//	}
//	sasQueryParams, err := AccountSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    currentTime,
//		Permissions:   AccountSASPermissions{Read: true, List: true}.String(),
//		Services:      AccountSASServices{Blob: true}.String(),
//		ResourceTypes: AccountSASResourceTypes{Container: true, Object: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		s.Fail(err.Error())
//	}
//
//	parts := NewBlobURLParts(bbClient.URL())
//	parts.SAS = sasQueryParams
//	parts.Snapshot = currentTime.Format(SnapshotTimeFormat)
//	blobURLParts := parts.URL()
//
//	// The snapshot format string is taken from the snapshotTimeFormat value in parsing_urls.go. The field is not public, so
//	// it is copied here
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_assert.NoError(err)
//	correctURL := "https://" + accountName + DefaultBlobEndpointSuffix + containerName + "/" + blobName +
//		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
//	_assert.Equal(blobURLParts, correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClientWithSnapshotAndSASUsingConnectionString() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClientFromConnectionString(nil, testAccountDefault, nil)
//	_assert.NoError(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	blobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_assert.NoError(err)
//
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	if err != nil {
//		s.Fail(err.Error())
//	}
//	sasQueryParams, err := AccountSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    currentTime,
//		Permissions:   AccountSASPermissions{Read: true, List: true}.String(),
//		Services:      AccountSASServices{Blob: true}.String(),
//		ResourceTypes: AccountSASResourceTypes{Container: true, Object: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		s.Fail(err.Error())
//	}
//
//	parts := NewBlobURLParts(bbClient.URL())
//	parts.SAS = sasQueryParams
//	parts.Snapshot = currentTime.Format(SnapshotTimeFormat)
//	blobURLParts := parts.URL()
//
//	// The snapshot format string is taken from the snapshotTimeFormat value in parsing_urls.go. The field is not public, so
//	// it is copied here
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_assert.NoError(err)
//	correctURL := "https://" + accountName + DefaultBlobEndpointSuffix + containerName + "/" + blobName +
//		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
//	_assert.Equal(blobURLParts, correctURL)
//}

func waitForCopy(t *testing.T, copyBlobClient BlockBlobClient, blobCopyResponse BlobStartCopyFromURLResponse) {
	status := *blobCopyResponse.CopyStatus
	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		props, err := copyBlobClient.GetProperties(ctx, nil)
		require.NoError(t, err)
		status = *props.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			require.Fail(t, "If the copy takes longer than a minute, we will fail")
		}
	}
}

func TestBlobStartCopyDestEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	require.NoError(t, err)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	blobCopyResponse, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	require.NoError(t, err)
	waitForCopy(t, copyBlobClient, blobCopyResponse)

	resp, err := copyBlobClient.Download(ctx, nil)
	require.NoError(t, err)

	// Read the blob data to verify the copy
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentLength, int64(len(blockBlobDefaultData)))
	require.Equal(t, string(data), blockBlobDefaultData)
	err = resp.Body(nil).Close()
	require.NoError(t, err)
}

func TestBlobStartCopyMetadata(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	require.NoError(t, err)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	metadata := make(map[string]string)
	metadata["Bla"] = "foo"
	options := StartCopyBlobOptions{
		Metadata: metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)
	waitForCopy(t, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp2.Metadata, metadata)
}

func TestBlobStartCopyMetadataNil(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	anotherBlobName := "copy" + blockBlobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the nil metadata passed later takes effect
	_, err = copyBlobClient.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), nil)
	require.NoError(t, err)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	require.NoError(t, err)

	waitForCopy(t, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Len(t, resp2.Metadata, 0)
}

func TestBlobStartCopyMetadataEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the empty metadata passed later takes effect
	_, err = copyBlobClient.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), nil)
	require.NoError(t, err)

	metadata := make(map[string]string)
	options := StartCopyBlobOptions{
		Metadata: metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	waitForCopy(t, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Len(t, resp2.Metadata, 0)
}

func TestBlobStartCopyMetadataInvalidField(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)

	anotherBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	metadata := make(map[string]string)
	metadata["I nvalid."] = "foo"
	options := StartCopyBlobOptions{
		Metadata: metadata,
	}
	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
	require.Equal(t, strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func TestBlobStartCopySourceNonExistent(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	require.Error(t, err)
	require.Equal(t, strings.Contains(err.Error(), "not exist"), true)
}

func TestBlobStartCopySourcePrivate(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	require.NoError(t, err)

	bbClient := createNewBlockBlob(_assert, generateBlobName(testName), containerClient)

	serviceClient2, err := createServiceClient(t, testAccountSecondary)
	require.NoError(t, err)

	copyContainerClient := createNewContainer(t, "cpyc"+containerName, serviceClient2)
	defer deleteContainer(_assert, copyContainerClient)
	copyBlobName := "copyb" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(copyBlobName, copyContainerClient)

	require.NotEqual(t, svcClient.URL(), serviceClient2.URL(), "Test not valid because primary and secondary accounts are the same")
	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	validateStorageError(_assert, err, StorageErrorCodeCannotVerifyCopySource)
}

func TestBlobStartCopyUsingSASSrc(t *testing.T) {
	recording.LiveOnly(t) // SAS are live only
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	require.NoError(t, err)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	// Create sas values for the source blob
	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)

	startTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	require.NoError(t, err)
	endTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	require.NoError(t, err)
	serviceSASValues := BlobSASSignatureValues{
		StartTime:     startTime,
		ExpiryTime:    endTime,
		Permissions:   BlobSASPermissions{Read: true, Write: true}.String(),
		ContainerName: containerName,
		BlobName:      blockBlobName}
	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
	require.NoError(t, err)

	// Create URLs to the destination blob with sas parameters
	sasURL := NewBlobURLParts(bbClient.URL())
	sasURL.SAS = queryParams

	// Create a new container for the destination
	serviceClient2, err := createServiceClient(t, testAccountSecondary)
	require.NoError(t, err)

	copyContainerName := "copy" + generateContainerName(testName)
	copyContainerClient := createNewContainer(t, copyContainerName, serviceClient2)
	defer deleteContainer(_assert, copyContainerClient)

	copyBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(copyBlobName, copyContainerClient)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, sasURL.URL(), nil)
	require.NoError(t, err)

	waitForCopy(t, copyBlobClient, resp)

	downloadBlobOptions := DownloadBlobOptions{
		Offset: to.Int64Ptr(0),
		Count:  to.Int64Ptr(int64(len(blockBlobDefaultData))),
	}
	resp2, err := copyBlobClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)

	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, *resp2.ContentLength, int64(len(blockBlobDefaultData)))
	require.Equal(t, string(data), blockBlobDefaultData)
	err = resp2.Body(nil).Close()
	require.NoError(t, err)
}

func TestBlobStartCopyUsingSASDest(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()

	var svcClient ServiceClient
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = createServiceClient(t, testAccountDefault) // getServiceClient(nil, testAccountDefault, nil)
		} else {
			svcClient, err = createServiceClientFromConnectionString(t, testAccountDefault)
		}
		require.NoError(t, err)

		containerClient := createNewContainer(t, generateContainerName(testName)+strconv.Itoa(i), svcClient)
		_, err := containerClient.SetAccessPolicy(ctx, nil)
		require.NoError(t, err)

		blobClient := createNewBlockBlob(_assert, generateBlobName(testName), containerClient)
		_, err = blobClient.Delete(ctx, nil)
		require.NoError(t, err)

		deleteContainer(_assert, containerClient)
	}
}

func TestBlobStartCopySourceIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient := getBlockBlobClient("dst"+generateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopySourceIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient := getBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
}

func TestBlobStartCopySourceIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}

	destBlobClient := getBlockBlobClient("dst"+generateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopySourceIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}
	destBlobClient := getBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
}

func TestBlobStartCopySourceIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopySourceIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	randomEtag := "a"
	accessConditions := SourceModifiedAccessConditions{
		SourceIfMatch: &randomEtag,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeSourceConditionNotMet)
}

func TestBlobStartCopySourceIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfNoneMatch: to.StringPtr("a"),
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopySourceIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfNoneMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeSourceConditionNotMet)
}

func TestBlobStartCopyDestIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	destBlobClient := createNewBlockBlob(_assert, "dst"+bbName, containerClient) // The blob must exist to have a last-modified time
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopyDestIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	destBlobClient := createNewBlockBlob(_assert, "dst"+bbName, containerClient) // The blob must exist to have a last-modified time

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	validateStorageError(_assert, err, StorageErrorCodeTargetConditionNotMet)
}

func TestBlobStartCopyDestIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	destBlobClient := createNewBlockBlob(_assert, "dst"+bbName, containerClient)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopyDestIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	destBlobClient := createNewBlockBlob(_assert, "dst"+bbName, containerClient)
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
}

func TestBlobStartCopyDestIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopyDestIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	metadata := make(map[string]string)
	metadata["bla"] = "bla"
	_, err = destBlobClient.SetMetadata(ctx, metadata, nil)
	require.NoError(t, err)

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeTargetConditionNotMet)
}

func TestBlobStartCopyDestIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.SetMetadata(ctx, nil, nil) // SetMetadata chances the blob's etag
	require.NoError(t, err)

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.NoError(t, err)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)
}

func TestBlobStartCopyDestIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeTargetConditionNotMet)
}

func TestBlobAbortCopyInProgress(t *testing.T) {
	recording.LiveOnly(t) // Fails because of random data
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	// Create a large blob that takes time to copy
	blobSize := 8 * 1024 * 1024
	blobReader, _ := getRandomDataAndReader(blobSize)
	_, err = bbClient.Upload(ctx, internal.NopCloser(blobReader), nil)
	require.NoError(t, err)

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions) // So that we don't have to create a SAS
	require.NoError(t, err)

	// Must copy across accounts so it takes time to copy
	serviceClient2, err := createServiceClient(t, testAccountSecondary)
	require.NoError(t, err)

	copyContainerName := "copy" + generateContainerName(testName)
	copyContainerClient := createNewContainer(t, copyContainerName, serviceClient2)

	copyBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(copyBlobName, copyContainerClient)

	defer deleteContainer(_assert, copyContainerClient)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	require.NoError(t, err)
	require.Equal(t, *resp.CopyStatus, CopyStatusTypePending)
	require.NotNil(t, resp.CopyID)

	_, err = copyBlobClient.AbortCopyFromURL(ctx, *resp.CopyID, nil)
	if err != nil {
		// If the error is nil, the test continues as normal.
		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
		validateStorageError(_assert, err, StorageErrorCodeNoPendingCopyOperation)
		require.Error(t, errors.New("the test failed because the copy completed because it was aborted"))
	}

	resp2, _ := copyBlobClient.GetProperties(ctx, nil)
	require.Equal(t, *resp2.CopyStatus, CopyStatusTypeAborted)
}

func TestBlobAbortCopyNoCopyStarted(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(blockBlobName, containerClient)

	_, err = copyBlobClient.AbortCopyFromURL(ctx, "copynotstarted", nil)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidQueryParameterValue)
}

func TestBlobSnapshotMetadata(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		Metadata: basicMetadata,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	require.NoError(t, err)
	require.NotNil(t, resp.Snapshot)

	// Since metadata is specified on the snapshot, the snapshot should have its own metadata different from the (empty) metadata on the source
	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp2.Metadata, basicMetadata)
}

func TestBlobSnapshotMetadataEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.NoError(t, err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Snapshot)

	// In this case, because no metadata was specified, it should copy the basicMetadata from the source
	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp2.Metadata, basicMetadata)
}

func TestBlobSnapshotMetadataNil(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.NoError(t, err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Snapshot)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp2.Metadata, basicMetadata)
}

func TestBlobSnapshotMetadataInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		Metadata: map[string]string{"Invalid Field!": "value"},
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	require.Error(t, err)
	require.Contains(t, err.Error(), invalidHeaderErrorSubstring)
}

func TestBlobSnapshotBlobNotExist(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(ctx, nil)
	require.Error(t, err)
}

func TestBlobSnapshotOfSnapshot(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	snapshotString, err := time.Parse(SnapshotTimeFormat, "2021-01-01T01:01:01.0000000Z")
	require.NoError(t, err)
	snapshotURL := bbClient.WithSnapshot(snapshotString.String())
	// The library allows the server to handle the snapshot of snapshot error
	_, err = snapshotURL.CreateSnapshot(ctx, nil)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidQueryParameterValue)
}

func TestBlobSnapshotIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	require.NoError(t, err)
	require.NotEqual(t, *resp.Snapshot, "") // i.e. The snapshot time is not zero. If the service gives us back a snapshot time, it successfully created a snapshot
}

func TestBlobSnapshotIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	require.Error(t, err)
}

func TestBlobSnapshotIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	require.NoError(t, err)
	require.NotEqual(t, *resp.Snapshot, "")
}

func TestBlobSnapshotIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	require.Error(t, err)
}

func TestBlobSnapshotIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	resp2, err := bbClient.CreateSnapshot(ctx, &options)
	require.NoError(t, err)
	require.NotEqual(t, *resp2.Snapshot, "")
}

func TestBlobSnapshotIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	randomEtag := "garbage"
	access := ModifiedAccessConditions{
		IfMatch: &randomEtag,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	require.Error(t, err)
}

func TestBlobSnapshotIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	randomEtag := "garbage"
	access := ModifiedAccessConditions{
		IfNoneMatch: &randomEtag,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	require.NoError(t, err)
	require.NotEqual(t, *resp.Snapshot, "")
}

func TestBlobSnapshotIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	require.Error(t, err)
}

func TestBlobDownloadDataNonExistentBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := containerClient.NewBlobClient(blobName)

	_, err = bbClient.Download(ctx, nil)
	require.Error(t, err)
}

func TestBlobDownloadDataNegativeOffset(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Offset: to.Int64Ptr(-1),
	}
	_, err = bbClient.Download(ctx, &options)
	require.NoError(t, err)
}

func TestBlobDownloadDataOffsetOutOfRange(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Offset: to.Int64Ptr(int64(len(blockBlobDefaultData))),
	}
	_, err = bbClient.Download(ctx, &options)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidRange)
}

func TestBlobDownloadDataCountNegative(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count: to.Int64Ptr(-2),
	}
	_, err = bbClient.Download(ctx, &options)
	require.NoError(t, err)
}

func TestBlobDownloadDataCountZero(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count: to.Int64Ptr(0),
	}
	resp, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)

	// Specifying a count of 0 results in the value being ignored
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), blockBlobDefaultData)
}

func TestBlobDownloadDataCountExact(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	count := int64(len(blockBlobDefaultData))
	options := DownloadBlobOptions{
		Count: &count,
	}
	resp, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), blockBlobDefaultData)
}

func TestBlobDownloadDataCountOutOfRange(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count: to.Int64Ptr(int64((len(blockBlobDefaultData)) * 2)),
	}
	resp, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), blockBlobDefaultData)
}

func TestBlobDownloadDataEmptyRangeStruct(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count:  to.Int64Ptr(0),
		Offset: to.Int64Ptr(0),
	}
	resp, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), blockBlobDefaultData)
}

func TestBlobDownloadDataContentMD5(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count:              to.Int64Ptr(3),
		Offset:             to.Int64Ptr(10),
		RangeGetContentMD5: to.BoolPtr(true),
	}
	resp, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)
	mdf := md5.Sum([]byte(blockBlobDefaultData)[10:13])
	require.Equal(t, resp.ContentMD5, mdf[:])
}

func TestBlobDownloadDataIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}

	resp, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentLength, int64(len(blockBlobDefaultData)))
}

func TestBlobDownloadDataIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}
	_, err = bbClient.Download(ctx, &options)
	require.Error(t, err)
}

func TestBlobDownloadDataIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentLength, int64(len(blockBlobDefaultData)))
}

func TestBlobDownloadDataIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}
	_, err = bbClient.Download(ctx, &options)
	require.Error(t, err)
}

func TestBlobDownloadDataIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, *resp2.ContentLength, int64(len(blockBlobDefaultData)))
}

func TestBlobDownloadDataIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	require.NoError(t, err)

	_, err = bbClient.Download(ctx, &options)
	require.Error(t, err)
}

func TestBlobDownloadDataIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	access := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	require.NoError(t, err)

	resp2, err := bbClient.Download(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, *resp2.ContentLength, int64(len(blockBlobDefaultData)))
}

func TestBlobDownloadDataIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}

	_, err = bbClient.Download(ctx, &options)
	require.Error(t, err)
}

func TestBlobDeleteNonExistent(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blockBlobName)

	_, err = bbClient.Delete(ctx, nil)
	require.Error(t, err)
}

func TestBlobDeleteSnapshot(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	require.NoError(t, err)
	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)

	_, err = snapshotURL.Delete(ctx, nil)
	require.NoError(t, err)

	validateBlobDeleted(t, snapshotURL.BlobClient)
}

//
////func (s *azblobTestSuite) TestBlobDeleteSnapshotsInclude() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := bbClient.CreateSnapshot(ctx, nil)
////	_assert.NoError(err)
////
////	deleteSnapshots := DeleteSnapshotsOptionInclude
////	_, err = bbClient.Delete(ctx, &DeleteBlobOptions{
////		DeleteSnapshots: &deleteSnapshots,
////	})
////	_assert.NoError(err)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
////		Include: include,
////	}
////	blobs, errChan := containerClient.ListBlobsFlat(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(<- blobs, chk.HasLen, 0)
////}
//
////func (s *azblobTestSuite) TestBlobDeleteSnapshotsOnly() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := bbClient.CreateSnapshot(ctx, nil)
////	_assert.NoError(err)
////	deleteSnapshot := DeleteSnapshotsOptionOnly
////	deleteBlobOptions := DeleteBlobOptions{
////		DeleteSnapshots: &deleteSnapshot,
////	}
////	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
////	_assert.NoError(err)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
////		Include: include,
////	}
////	blobs, errChan := containerClient.ListBlobsFlat(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(blobs, chk.HasLen, 1)
////	_assert(*(<-blobs).Snapshot == "", chk.Equals, true)
////}

func TestBlobDeleteSnapshotsNoneWithSnapshots(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(ctx, nil)
	require.NoError(t, err)
	_, err = bbClient.Delete(ctx, nil)
	require.Error(t, err)
}

func validateBlobDeleted(t *testing.T, bbClient BlobClient) {
	_, err := bbClient.GetProperties(ctx, nil)
	require.Error(t, err)

	var storageError *StorageError
	require.Equal(t, true, errors.As(err, &storageError))
	require.Equal(t, storageError.ErrorCode, StorageErrorCodeBlobNotFound)
}

func TestBlobDeleteIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	require.NoError(t, err)

	validateBlobDeleted(t, bbClient.BlobClient)
}

func TestBlobDeleteIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobDeleteIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	require.NoError(t, err)

	validateBlobDeleted(t, bbClient.BlobClient)
}

func TestBlobDeleteIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobDeleteIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(ctx, nil)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	require.NoError(t, err)

	validateBlobDeleted(t, bbClient.BlobClient)
}

func TestBlobDeleteIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	etag := resp.ETag

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	require.NoError(t, err)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobDeleteIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(ctx, nil)
	etag := resp.ETag
	_, err = bbClient.SetMetadata(ctx, nil, nil)
	require.NoError(t, err)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	require.NoError(t, err)

	validateBlobDeleted(t, bbClient.BlobClient)
}

func TestBlobDeleteIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(ctx, nil)
	etag := resp.ETag

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobGetPropsAndMetadataIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.NoError(t, err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.NoError(t, err)
	require.EqualValues(t, resp.Metadata, basicMetadata)
}

func TestBlobGetPropsAndMetadataIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.NoError(t, err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.Error(t, err)
	var storageError *StorageError
	require.Equal(t, errors.As(err, &storageError), true)
	require.Equal(t, storageError.response.StatusCode, 304) // No service code returned for a HEAD
}

func TestBlobGetPropsAndMetadataIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.NoError(t, err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.NoError(t, err)
	require.EqualValues(t, resp.Metadata, basicMetadata)
}

//func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceFalse() {
//	// TODO: Not Working
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	blockBlobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blockBlobName, containerClient)
//
//	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
//
//	_assert.NoError(err)
//	_assert.Equal(cResp.RawResponse.StatusCode, 201)
//
//	currentTime := getRelativeTimeFromAnchor(cResp.Date,-10)
//
//	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert.NoError(err)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert.Error(err)
//}

func TestBlobGetPropsAndMetadataIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.NoError(t, err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.NoError(t, err)
	require.EqualValues(t, resp2.Metadata, basicMetadata)
}

func TestBlobGetPropsOnMissingBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := containerClient.NewBlobClient("MISSING")

	_, err = bbClient.GetProperties(ctx, nil)
	require.Error(t, err)
	var storageError *StorageError
	require.Equal(t, errors.As(err, &storageError), true)
	require.Equal(t, storageError.response.StatusCode, 404)
	require.Equal(t, storageError.response.Header.Get("x-ms-error-code"), string(StorageErrorCodeBlobNotFound))
}

func TestBlobGetPropsAndMetadataIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	eTag := "garbage"
	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.Error(t, err)
	var storageError *StorageError
	require.Equal(t, errors.As(err, &storageError), true)
	require.Equal(t, storageError.response.StatusCode, 412)
}

func TestBlobGetPropsAndMetadataIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.NoError(t, err)

	eTag := "garbage"
	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.NoError(t, err)
	require.EqualValues(t, resp.Metadata, basicMetadata)
}

func TestBlobGetPropsAndMetadataIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(ctx, nil, nil)
	require.NoError(t, err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.Error(t, err)
	var storageError *StorageError
	require.Equal(t, errors.As(err, &storageError), true)
	require.Equal(t, storageError.response.StatusCode, 304)
}

func TestBlobSetPropertiesBasic(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, basicHeaders, nil)
	require.NoError(t, err)

	resp, _ := bbClient.GetProperties(ctx, nil)
	h := resp.GetHTTPHeaders()
	require.EqualValues(t, h, basicHeaders)
}

func TestBlobSetPropertiesEmptyValue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	contentType := to.StringPtr("my_type")
	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentType: contentType}, nil)
	require.NoError(t, err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp.ContentType, contentType)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{}, nil)
	require.NoError(t, err)

	resp, err = bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.ContentType)
}

func validatePropertiesSet(t *testing.T, bbClient BlockBlobClient, disposition string) {
	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentDisposition, disposition)
}

func TestBlobSetPropertiesIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
	require.NoError(t, err)

	validatePropertiesSet(t, bbClient, "my_disposition")
}

func TestBlobSetPropertiesIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{
		BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	)
	require.Error(t, err)
}

func TestBlobSetPropertiesIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	require.NoError(t, err)

	validatePropertiesSet(t, bbClient, "my_disposition")
}

func TestBlobSetPropertiesIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	require.Error(t, err)
}

func TestBlobSetPropertiesIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}})
	require.NoError(t, err)

	validatePropertiesSet(t, bbClient, "my_disposition")
}

func TestBlobSetPropertiesIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: to.StringPtr("garbage")}})
	require.Error(t, err)
}

func TestBlobSetPropertiesIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: to.StringPtr("garbage")}})
	require.NoError(t, err)

	validatePropertiesSet(t, bbClient, "my_disposition")
}

func TestBlobSetPropertiesIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}})
	require.Error(t, err)
}

func TestBlobSetMetadataNil(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
	require.NoError(t, err)

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	require.NoError(t, err)

	blobGetResp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Len(t, blobGetResp.Metadata, 0)
}

func TestBlobSetMetadataEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
	require.NoError(t, err)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, nil)
	require.NoError(t, err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Len(t, resp.Metadata, 0)
}

func TestBlobSetMetadataInvalidField(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"Invalid field!": "value"}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), invalidHeaderErrorSubstring)
	//_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func validateMetadataSet(t *testing.T, bbClient BlockBlobClient) {
	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp.Metadata, basicMetadata)
}

func TestBlobSetMetadataIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	require.NoError(t, err)

	validateMetadataSet(t, bbClient)
}

func TestBlobSetMetadataIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobSetMetadataIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	require.NoError(t, err)

	validateMetadataSet(t, bbClient)
}

func TestBlobSetMetadataIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobSetMetadataIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	require.NoError(t, err)

	validateMetadataSet(t, bbClient)
}

func TestBlobSetMetadataIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	eTag := "garbage"
	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobSetMetadataIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	eTag := "garbage"
	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	require.NoError(t, err)

	validateMetadataSet(t, bbClient)
}

func TestBlobSetMetadataIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//nolint
func testBlobServiceClientDeleteImpl(_ *assert.Assertions, _ ServiceClient) error {
	//containerClient := createNewContainer(t, "gocblobserviceclientdeleteimpl", svcClient)
	//defer deleteContainer(_assert, containerClient)
	//bbClient := createNewBlockBlob(_assert, "goblobserviceclientdeleteimpl", containerClient)
	//
	//_, err := bbClient.Delete(ctx, nil)
	//_assert.NoError(err) // This call will not have errors related to slow update of service properties, so we assert.
	//
	//_, err = bbClient.Undelete(ctx)
	//if err != nil { // We want to give the wrapper method a chance to check if it was an error related to the service properties update.
	//	return err
	//}
	//
	//resp, err := bbClient.GetProperties(ctx, nil)
	//if err != nil {
	//	return errors.New(string(err.(*StorageError).ErrorCode))
	//}
	//_assert.Equal(resp.BlobType, BlobTypeBlockBlob) // We could check any property. This is just to double check it was undeleted.
	return nil
}

func TestBlobServiceClientDelete(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	code := 404
	runTestRequiringServiceProperties(_assert, svcClient, string(rune(code)), enableSoftDelete, testBlobServiceClientDeleteImpl, disableSoftDelete)
}

func setAndCheckBlobTier(t *testing.T, bbClient BlobClient, tier AccessTier) {
	_, err := bbClient.SetTier(ctx, tier, nil)
	require.NoError(t, err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.AccessTier, string(tier))
}

func TestBlobSetTierAllTiers(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	setAndCheckBlobTier(t, bbClient.BlobClient, AccessTierHot)
	setAndCheckBlobTier(t, bbClient.BlobClient, AccessTierCool)
	setAndCheckBlobTier(t, bbClient.BlobClient, AccessTierArchive)

	premiumServiceClient, err := createServiceClient(t, testAccountPremium)
	require.NoError(t, err)

	premContainerName := "prem" + generateContainerName(testName)
	premContainerClient := createNewContainer(t, premContainerName, premiumServiceClient)
	defer deleteContainer(_assert, premContainerClient)

	pbClient := createNewPageBlob(_assert, blockBlobName, premContainerClient)

	setAndCheckBlobTier(t, pbClient.BlobClient, AccessTierP4)
	setAndCheckBlobTier(t, pbClient.BlobClient, AccessTierP6)
	setAndCheckBlobTier(t, pbClient.BlobClient, AccessTierP10)
	setAndCheckBlobTier(t, pbClient.BlobClient, AccessTierP20)
	setAndCheckBlobTier(t, pbClient.BlobClient, AccessTierP30)
	setAndCheckBlobTier(t, pbClient.BlobClient, AccessTierP40)
	setAndCheckBlobTier(t, pbClient.BlobClient, AccessTierP50)
}

//
////func (s *azblobTestSuite) TestBlobTierInferred() {
////	svcClient, err := getPremiumserviceClient()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	bbClient, _ := createNewPageBlob(c, containerClient)
////
////	resp, err := bbClient.GetProperties(ctx, nil)
////	_assert.NoError(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "true")
////
////	resp2, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.NoError(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.NotNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTier, chk.Not(chk.Equals), "")
////
////	_, err = bbClient.SetTier(ctx, AccessTierP4, LeaseAccessConditions{})
////	_assert.NoError(err)
////
////	resp, err = bbClient.GetProperties(ctx, nil)
////	_assert.NoError(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "")
////
////	resp2, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.NoError(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.IsNil) // AccessTierInferred never returned if false
////}
////
////func (s *azblobTestSuite) TestBlobArchiveStatus() {
////	svcClient, err := getBlobStorageserviceClient()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_assert.NoError(err)
////	_, err = bbClient.SetTier(ctx, AccessTierCool, LeaseAccessConditions{})
////	_assert.NoError(err)
////
////	resp, err := bbClient.GetProperties(ctx, nil)
////	_assert.NoError(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToCool))
////
////	resp2, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.NoError(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToCool)
////
////	// delete first blob
////	_, err = bbClient.Delete(ctx, DeleteSnapshotsOptionNone, nil)
////	_assert.NoError(err)
////
////	bbClient, _ = createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_assert.NoError(err)
////	_, err = bbClient.SetTier(ctx, AccessTierHot, LeaseAccessConditions{})
////	_assert.NoError(err)
////
////	resp, err = bbClient.GetProperties(ctx, nil)
////	_assert.NoError(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToHot))
////
////	resp2, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.NoError(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToHot)
////}
////
////func (s *azblobTestSuite) TestBlobTierInvalidValue() {
////	svcClient, err := getBlobStorageserviceClient()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierType("garbage"), LeaseAccessConditions{})
////	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
////}
////

func TestBlobClientPartsSASQueryTimes(t *testing.T) {
	StartTimesInputs := []string{
		"2020-04-20",
		"2020-04-20T07:00Z",
		"2020-04-20T07:15:00Z",
		"2020-04-20T07:30:00.1234567Z",
	}
	StartTimesExpected := []time.Time{
		time.Date(2020, time.April, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 7, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 7, 15, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 7, 30, 0, 123456700, time.UTC),
	}
	ExpiryTimesInputs := []string{
		"2020-04-21",
		"2020-04-20T08:00Z",
		"2020-04-20T08:15:00Z",
		"2020-04-20T08:30:00.2345678Z",
	}
	ExpiryTimesExpected := []time.Time{
		time.Date(2020, time.April, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 8, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 8, 15, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 8, 30, 0, 234567800, time.UTC),
	}

	for i := 0; i < len(StartTimesInputs); i++ {
		urlString :=
			"https://myaccount.blob.core.windows.net/mycontainer/mydirectory/myfile.txt?" +
				"se=" + url.QueryEscape(ExpiryTimesInputs[i]) + "&" +
				"sig=NotASignature&" +
				"sp=r&" +
				"spr=https&" +
				"sr=b&" +
				"st=" + url.QueryEscape(StartTimesInputs[i]) + "&" +
				"sv=2019-10-10"

		parts := NewBlobURLParts(urlString)
		require.Equal(t, parts.Scheme, "https")
		require.Equal(t, parts.Host, "myaccount.blob.core.windows.net")
		require.Equal(t, parts.ContainerName, "mycontainer")
		require.Equal(t, parts.BlobName, "mydirectory/myfile.txt")

		sas := parts.SAS
		require.Equal(t, sas.StartTime(), StartTimesExpected[i])
		require.Equal(t, sas.ExpiryTime(), ExpiryTimesExpected[i])

		require.Equal(t, parts.URL(), urlString)
	}
}

func TestDownloadBlockBlobUnexpectedEOF(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.Download(ctx, nil)
	require.NoError(t, err)

	// Verify that we can inject errors first.
	reader := resp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))

	_, err = ioutil.ReadAll(reader)
	require.Error(t, err)
	require.Equal(t, err.Error(), "unrecoverable error")

	// Then inject the retryable error.
	reader = resp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))

	t.Skip("The range specified is invalid for the current size of the resource")
	buf, err := ioutil.ReadAll(reader)
	require.NoError(t, err)
	require.EqualValues(t, buf, []byte(blockBlobDefaultData))
}

//nolint
func InjectErrorInRetryReaderOptions(err error) *RetryReaderOptions {
	return &RetryReaderOptions{
		MaxRetryRequests:       1,
		doInjectError:          true,
		doInjectErrorRound:     0,
		injectedError:          err,
		NotifyFailedRead:       nil,
		TreatEarlyCloseAsError: false,
		CpkInfo:                nil,
		CpkScopeInfo:           nil,
	}
}
