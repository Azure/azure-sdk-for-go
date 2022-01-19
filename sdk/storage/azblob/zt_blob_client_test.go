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
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
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
//	_assert.Nil(err)
//	correctURL := "https://" + accountName + "." + DefaultBlobEndpointSuffix + containerName + "/" + blobName
//	_assert.Equal(bbClient.URL(), correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClientWithSnapshotAndSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	_assert.Nil(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	blobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_assert.Nil(err)
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
//	_assert.Nil(err)
//	correctURL := "https://" + accountName + DefaultBlobEndpointSuffix + containerName + "/" + blobName +
//		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
//	_assert.Equal(blobURLParts, correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClientWithSnapshotAndSASUsingConnectionString() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClientFromConnectionString(nil, testAccountDefault, nil)
//	_assert.Nil(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	blobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_assert.Nil(err)
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
//	_assert.Nil(err)
//	correctURL := "https://" + accountName + DefaultBlobEndpointSuffix + containerName + "/" + blobName +
//		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
//	_assert.Equal(blobURLParts, correctURL)
//}

func waitForCopy(_assert *assert.Assertions, copyBlobClient BlockBlobClient, blobCopyResponse BlobStartCopyFromURLResponse) {
	status := *blobCopyResponse.CopyStatus
	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		props, _ := copyBlobClient.GetProperties(ctx, nil)
		status = *props.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_assert.Fail("If the copy takes longer than a minute, we will fail")
		}
	}
}

func (s *azblobTestSuite) TestBlobStartCopyDestEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	_assert.Nil(err)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	blobCopyResponse, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_assert.Nil(err)
	waitForCopy(_assert, copyBlobClient, blobCopyResponse)

	resp, err := copyBlobClient.Download(ctx, nil)
	_assert.Nil(err)

	// Read the blob data to verify the copy
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(*resp.ContentLength, int64(len(blockBlobDefaultData)))
	_assert.Equal(string(data), blockBlobDefaultData)
	_ = resp.Body(RetryReaderOptions{}).Close()
}

func (s *azblobTestSuite) TestBlobStartCopyMetadata() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	_assert.Nil(err)
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
	_assert.Nil(err)
	waitForCopy(_assert, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp2.Metadata, metadata)
}

func (s *azblobTestSuite) TestBlobStartCopyMetadataNil() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	anotherBlobName := "copy" + blockBlobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the nil metadata passed later takes effect
	_, err = copyBlobClient.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_assert.Nil(err)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_assert.Nil(err)

	waitForCopy(_assert, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp2.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobStartCopyMetadataEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the empty metadata passed later takes effect
	_, err = copyBlobClient.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_assert.Nil(err)

	metadata := make(map[string]string)
	options := StartCopyBlobOptions{
		Metadata: metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	waitForCopy(_assert, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp2.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobStartCopyMetadataInvalidField() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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
	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *azblobTestSuite) TestBlobStartCopySourceNonExistent() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), "not exist"), true)
}

func (s *azblobTestSuite) TestBlobStartCopySourcePrivate() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	_assert.Nil(err)

	bbClient := createNewBlockBlob(_assert, generateBlobName(testName), containerClient)

	serviceClient2, err := getServiceClient(_context.recording, testAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
		return
	}

	copyContainerClient := createNewContainer(_assert, "cpyc"+containerName, serviceClient2)
	defer deleteContainer(_assert, copyContainerClient)
	copyBlobName := "copyb" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(copyBlobName, copyContainerClient)

	if svcClient.URL() == serviceClient2.URL() {
		s.T().Skip("Test not valid because primary and secondary accounts are the same")
	}
	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	validateStorageError(_assert, err, StorageErrorCodeCannotVerifyCopySource)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobStartCopyUsingSASSrc() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	_assert.Nil(err)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	// Create sas values for the source blob
	credential, err := getGenericCredential(nil, testAccountDefault)
	if err != nil {
		s.T().Fatal("Couldn't fetch credential because " + err.Error())
	}

	startTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)
	endTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_assert.Nil(err)
	serviceSASValues := BlobSASSignatureValues{
		StartTime:     startTime,
		ExpiryTime:    endTime,
		Permissions:   BlobSASPermissions{Read: true, Write: true}.String(),
		ContainerName: containerName,
		BlobName:      blockBlobName}
	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	// Create URLs to the destination blob with sas parameters
	sasURL := NewBlobURLParts(bbClient.URL())
	sasURL.SAS = queryParams

	// Create a new container for the destination
	serviceClient2, err := getServiceClient(nil, testAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
	}

	copyContainerName := "copy" + generateContainerName(testName)
	copyContainerClient := createNewContainer(_assert, copyContainerName, serviceClient2)
	defer deleteContainer(_assert, copyContainerClient)

	copyBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(copyBlobName, copyContainerClient)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, sasURL.URL(), nil)
	_assert.Nil(err)

	waitForCopy(_assert, copyBlobClient, resp)

	downloadBlobOptions := DownloadBlobOptions{
		Offset: to.Int64Ptr(0),
		Count:  to.Int64Ptr(int64(len(blockBlobDefaultData))),
	}
	resp2, err := copyBlobClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(*resp2.ContentLength, int64(len(blockBlobDefaultData)))
	_assert.Equal(string(data), blockBlobDefaultData)
	_ = resp2.Body(RetryReaderOptions{}).Close()
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobStartCopyUsingSASDest() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	var svcClient ServiceClient
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = getServiceClient(nil, testAccountDefault, nil)
		} else {
			svcClient, err = getServiceClientFromConnectionString(nil, testAccountDefault, nil)
		}
		_assert.Nil(err)

		containerClient := createNewContainer(_assert, generateContainerName(testName)+strconv.Itoa(i), svcClient)
		_, err := containerClient.SetAccessPolicy(ctx, nil)
		_assert.Nil(err)

		blobClient := createNewBlockBlob(_assert, generateBlobName(testName), containerClient)
		_, err = blobClient.Delete(ctx, nil)
		_assert.Nil(err)

		deleteContainer(_assert, containerClient)
	}
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient := getBlockBlobClient("dst"+generateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient := getBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}

	destBlobClient := getBlockBlobClient("dst"+generateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}
	destBlobClient := getBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeSourceConditionNotMet)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfNoneMatch: to.StringPtr("a"),
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfNoneMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeSourceConditionNotMet)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	destBlobClient := createNewBlockBlob(_assert, "dst"+bbName, containerClient) // The blob must exist to have a last-modified time
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

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

func (s *azblobTestSuite) TestBlobStartCopyDestIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	destBlobClient := createNewBlockBlob(_assert, "dst"+bbName, containerClient)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	destBlobClient := createNewBlockBlob(_assert, "dst"+bbName, containerClient)
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	metadata := make(map[string]string)
	metadata["bla"] = "bla"
	_, err = destBlobClient.SetMetadata(ctx, metadata, nil)
	_assert.Nil(err)

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeTargetConditionNotMet)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.SetMetadata(ctx, nil, nil) // SetMetadata chances the blob's etag
	_assert.Nil(err)

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.Nil(err)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_assert, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeTargetConditionNotMet)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobAbortCopyInProgress() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	// Create a large blob that takes time to copy
	blobSize := 8 * 1024 * 1024
	blobReader, _ := getRandomDataAndReader(blobSize)
	_, err = bbClient.Upload(ctx, internal.NopCloser(blobReader), nil)
	_assert.Nil(err)

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions) // So that we don't have to create a SAS
	_assert.Nil(err)

	// Must copy across accounts so it takes time to copy
	serviceClient2, err := getServiceClient(nil, testAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
	}

	copyContainerName := "copy" + generateContainerName(testName)
	copyContainerClient := createNewContainer(_assert, copyContainerName, serviceClient2)

	copyBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(copyBlobName, copyContainerClient)

	defer deleteContainer(_assert, copyContainerClient)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_assert.Nil(err)
	_assert.Equal(*resp.CopyStatus, CopyStatusTypePending)
	_assert.NotNil(resp.CopyID)

	_, err = copyBlobClient.AbortCopyFromURL(ctx, *resp.CopyID, nil)
	if err != nil {
		// If the error is nil, the test continues as normal.
		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
		validateStorageError(_assert, err, StorageErrorCodeNoPendingCopyOperation)
		_assert.Error(errors.New("the test failed because the copy completed because it was aborted"))
	}

	resp2, _ := copyBlobClient.GetProperties(ctx, nil)
	_assert.Equal(*resp2.CopyStatus, CopyStatusTypeAborted)
}

func (s *azblobTestSuite) TestBlobAbortCopyNoCopyStarted() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(blockBlobName, containerClient)

	_, err = copyBlobClient.AbortCopyFromURL(ctx, "copynotstarted", nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidQueryParameterValue)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadata() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		Metadata: basicMetadata,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_assert.Nil(err)
	_assert.NotNil(resp.Snapshot)

	// Since metadata is specified on the snapshot, the snapshot should have its own metadata different from the (empty) metadata on the source
	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadataEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.Nil(err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Snapshot)

	// In this case, because no metadata was specified, it should copy the basicMetadata from the source
	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadataNil() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.Nil(err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Snapshot)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadataInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		Metadata: map[string]string{"Invalid Field!": "value"},
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_assert.NotNil(err)
	_assert.Contains(err.Error(), invalidHeaderErrorSubstring)
}

func (s *azblobTestSuite) TestBlobSnapshotBlobNotExist() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(ctx, nil)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotOfSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	snapshotString, err := time.Parse(SnapshotTimeFormat, "2021-01-01T01:01:01.0000000Z")
	_assert.Nil(err)
	snapshotURL := bbClient.WithSnapshot(snapshotString.String())
	// The library allows the server to handle the snapshot of snapshot error
	_, err = snapshotURL.CreateSnapshot(ctx, nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidQueryParameterValue)
}

func (s *azblobTestSuite) TestBlobSnapshotIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	_assert.Nil(err)
	_assert.NotEqual(*resp.Snapshot, "") // i.e. The snapshot time is not zero. If the service gives us back a snapshot time, it successfully created a snapshot
}

func (s *azblobTestSuite) TestBlobSnapshotIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	_assert.Nil(err)
	_assert.NotEqual(*resp.Snapshot, "")
}

func (s *azblobTestSuite) TestBlobSnapshotIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	resp2, err := bbClient.CreateSnapshot(ctx, &options)
	_assert.Nil(err)
	_assert.NotEqual(*resp2.Snapshot, "")
}

func (s *azblobTestSuite) TestBlobSnapshotIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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
	_assert.Nil(err)
	_assert.NotEqual(*resp.Snapshot, "")
}

func (s *azblobTestSuite) TestBlobSnapshotIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataNonExistentBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := containerClient.NewBlobClient(blobName)

	_, err = bbClient.Download(ctx, nil)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataNegativeOffset() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Offset: to.Int64Ptr(-1),
	}
	_, err = bbClient.Download(ctx, &options)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataOffsetOutOfRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Offset: to.Int64Ptr(int64(len(blockBlobDefaultData))),
	}
	_, err = bbClient.Download(ctx, &options)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidRange)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count: to.Int64Ptr(-2),
	}
	_, err = bbClient.Download(ctx, &options)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountZero() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count: to.Int64Ptr(0),
	}
	resp, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)

	// Specifying a count of 0 results in the value being ignored
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountExact() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	count := int64(len(blockBlobDefaultData))
	options := DownloadBlobOptions{
		Count: &count,
	}
	resp, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountOutOfRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count: to.Int64Ptr(int64((len(blockBlobDefaultData)) * 2)),
	}
	resp, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataEmptyRangeStruct() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count:  to.Int64Ptr(0),
		Offset: to.Int64Ptr(0),
	}
	resp, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataContentMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	options := DownloadBlobOptions{
		Count:              to.Int64Ptr(3),
		Offset:             to.Int64Ptr(10),
		RangeGetContentMD5: to.BoolPtr(true),
	}
	resp, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)
	mdf := md5.Sum([]byte(blockBlobDefaultData)[10:13])
	_assert.Equal(resp.ContentMD5, mdf[:])
}

func (s *azblobTestSuite) TestBlobDownloadDataIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}

	resp, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)
	_assert.Equal(*resp.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}
	_, err = bbClient.Download(ctx, &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)
	_assert.Equal(*resp.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}
	_, err = bbClient.Download(ctx, &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)
	_assert.Equal(*resp2.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_assert.Nil(err)

	_, err = bbClient.Download(ctx, &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	access := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_assert.Nil(err)

	resp2, err := bbClient.Download(ctx, &options)
	_assert.Nil(err)
	_assert.Equal(*resp2.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	options := DownloadBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}

	_, err = bbClient.Download(ctx, &options)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDeleteNonExistent() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blockBlobName)

	_, err = bbClient.Delete(ctx, nil)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDeleteSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)

	_, err = snapshotURL.Delete(ctx, nil)
	_assert.Nil(err)

	validateBlobDeleted(_assert, snapshotURL.BlobClient)
}

//
////func (s *azblobTestSuite) TestBlobDeleteSnapshotsInclude() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := bbClient.CreateSnapshot(ctx, nil)
////	_assert.Nil(err)
////
////	deleteSnapshots := DeleteSnapshotsOptionInclude
////	_, err = bbClient.Delete(ctx, &DeleteBlobOptions{
////		DeleteSnapshots: &deleteSnapshots,
////	})
////	_assert.Nil(err)
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
////	_assert.Nil(err)
////	deleteSnapshot := DeleteSnapshotsOptionOnly
////	deleteBlobOptions := DeleteBlobOptions{
////		DeleteSnapshots: &deleteSnapshot,
////	}
////	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
////	_assert.Nil(err)
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

func (s *azblobTestSuite) TestBlobDeleteSnapshotsNoneWithSnapshots() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	_, err = bbClient.Delete(ctx, nil)
	_assert.NotNil(err)
}

func validateBlobDeleted(_assert *assert.Assertions, bbClient BlobClient) {
	_, err := bbClient.GetProperties(ctx, nil)
	_assert.NotNil(err)

	var storageError *StorageError
	_assert.Equal(true, errors.As(err, &storageError))
	_assert.Equal(storageError.ErrorCode, StorageErrorCodeBlobNotFound)
}

func (s *azblobTestSuite) TestBlobDeleteIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	_assert.Nil(err)

	validateBlobDeleted(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobDeleteIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	_assert.Nil(err)

	validateBlobDeleted(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobDeleteIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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
	_assert.Nil(err)

	validateBlobDeleted(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	etag := resp.ETag

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_assert.Nil(err)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobDeleteIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(ctx, nil)
	etag := resp.ETag
	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_assert.Nil(err)

	deleteBlobOptions := DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	_assert.Nil(err)

	validateBlobDeleted(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.Nil(err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.Nil(err)
	_assert.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.Nil(err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.NotNil(err)
	var storageError *StorageError
	_assert.Equal(errors.As(err, &storageError), true)
	_assert.Equal(storageError.response.StatusCode, 304) // No service code returned for a HEAD
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.Nil(err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.Nil(err)
	_assert.EqualValues(resp.Metadata, basicMetadata)
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
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	blockBlobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blockBlobName, containerClient)
//
//	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
//
//	_assert.Nil(err)
//	_assert.Equal(cResp.RawResponse.StatusCode, 201)
//
//	currentTime := getRelativeTimeFromAnchor(cResp.Date,-10)
//
//	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert.Nil(err)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert.NotNil(err)
//}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.Nil(err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.Nil(err)
	_assert.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobGetPropsOnMissingBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := containerClient.NewBlobClient("MISSING")

	_, err = bbClient.GetProperties(ctx, nil)
	_assert.NotNil(err)
	var storageError *StorageError
	_assert.Equal(errors.As(err, &storageError), true)
	_assert.Equal(storageError.response.StatusCode, 404)
	_assert.Equal(storageError.response.Header.Get("x-ms-error-code"), string(StorageErrorCodeBlobNotFound))
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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
	_assert.NotNil(err)
	var storageError *StorageError
	_assert.Equal(errors.As(err, &storageError), true)
	_assert.Equal(storageError.response.StatusCode, 412)
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.Nil(err)

	eTag := "garbage"
	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.Nil(err)
	_assert.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(ctx, nil, nil)
	_assert.Nil(err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.NotNil(err)
	var storageError *StorageError
	_assert.Equal(errors.As(err, &storageError), true)
	_assert.Equal(storageError.response.StatusCode, 304)
}

func (s *azblobTestSuite) TestBlobSetPropertiesBasic() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, basicHeaders, nil)
	_assert.Nil(err)

	resp, _ := bbClient.GetProperties(ctx, nil)
	h := resp.GetHTTPHeaders()
	_assert.EqualValues(h, basicHeaders)
}

func (s *azblobTestSuite) TestBlobSetPropertiesEmptyValue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	contentType := to.StringPtr("my_type")
	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentType: contentType}, nil)
	_assert.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp.ContentType, contentType)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{}, nil)
	_assert.Nil(err)

	resp, err = bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.ContentType)
}

func validatePropertiesSet(_assert *assert.Assertions, bbClient BlockBlobClient, disposition string) {
	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.ContentDisposition, disposition)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
	_assert.Nil(err)

	validatePropertiesSet(_assert, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	_assert.Nil(err)

	validatePropertiesSet(_assert, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}})
	_assert.Nil(err)

	validatePropertiesSet(_assert, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: to.StringPtr("garbage")}})
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: to.StringPtr("garbage")}})
	_assert.Nil(err)

	validatePropertiesSet(_assert, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}})
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetMetadataNil() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
	_assert.Nil(err)

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_assert.Nil(err)

	blobGetResp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(blobGetResp.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobSetMetadataEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
	_assert.Nil(err)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, nil)
	_assert.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobSetMetadataInvalidField() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"Invalid field!": "value"}, nil)
	_assert.NotNil(err)
	_assert.Contains(err.Error(), invalidHeaderErrorSubstring)
	//_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func validateMetadataSet(_assert *assert.Assertions, bbClient BlockBlobClient) {
	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_assert.Nil(err)

	validateMetadataSet(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_assert.Nil(err)

	validateMetadataSet(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_assert.Nil(err)

	validateMetadataSet(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
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

func (s *azblobTestSuite) TestBlobSetMetadataIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	eTag := "garbage"
	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_assert.Nil(err)

	validateMetadataSet(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//nolint
func testBlobServiceClientDeleteImpl(_ *assert.Assertions, _ ServiceClient) error {
	//containerClient := createNewContainer(_assert, "gocblobserviceclientdeleteimpl", svcClient)
	//defer deleteContainer(_assert, containerClient)
	//bbClient := createNewBlockBlob(_assert, "goblobserviceclientdeleteimpl", containerClient)
	//
	//_, err := bbClient.Delete(ctx, nil)
	//_assert.Nil(err) // This call will not have errors related to slow update of service properties, so we assert.
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

func (s *azblobTestSuite) TestBlobServiceClientDelete() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	code := 404
	runTestRequiringServiceProperties(_assert, svcClient, string(rune(code)), enableSoftDelete, testBlobServiceClientDeleteImpl, disableSoftDelete)
}

func setAndCheckBlobTier(_assert *assert.Assertions, bbClient BlobClient, tier AccessTier) {
	_, err := bbClient.SetTier(ctx, tier, nil)
	_assert.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.AccessTier, string(tier))
}

func (s *azblobTestSuite) TestBlobSetTierAllTiers() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	setAndCheckBlobTier(_assert, bbClient.BlobClient, AccessTierHot)
	setAndCheckBlobTier(_assert, bbClient.BlobClient, AccessTierCool)
	setAndCheckBlobTier(_assert, bbClient.BlobClient, AccessTierArchive)

	premiumServiceClient, err := getServiceClient(_context.recording, testAccountPremium, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	premContainerName := "prem" + generateContainerName(testName)
	premContainerClient := createNewContainer(_assert, premContainerName, premiumServiceClient)
	defer deleteContainer(_assert, premContainerClient)

	pbClient := createNewPageBlob(_assert, blockBlobName, premContainerClient)

	setAndCheckBlobTier(_assert, pbClient.BlobClient, AccessTierP4)
	setAndCheckBlobTier(_assert, pbClient.BlobClient, AccessTierP6)
	setAndCheckBlobTier(_assert, pbClient.BlobClient, AccessTierP10)
	setAndCheckBlobTier(_assert, pbClient.BlobClient, AccessTierP20)
	setAndCheckBlobTier(_assert, pbClient.BlobClient, AccessTierP30)
	setAndCheckBlobTier(_assert, pbClient.BlobClient, AccessTierP40)
	setAndCheckBlobTier(_assert, pbClient.BlobClient, AccessTierP50)
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
////	_assert.Nil(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "true")
////
////	resp2, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.NotNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTier, chk.Not(chk.Equals), "")
////
////	_, err = bbClient.SetTier(ctx, AccessTierP4, LeaseAccessConditions{})
////	_assert.Nil(err)
////
////	resp, err = bbClient.GetProperties(ctx, nil)
////	_assert.Nil(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "")
////
////	resp2, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.Nil(err)
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
////	_assert.Nil(err)
////	_, err = bbClient.SetTier(ctx, AccessTierCool, LeaseAccessConditions{})
////	_assert.Nil(err)
////
////	resp, err := bbClient.GetProperties(ctx, nil)
////	_assert.Nil(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToCool))
////
////	resp2, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToCool)
////
////	// delete first blob
////	_, err = bbClient.Delete(ctx, DeleteSnapshotsOptionNone, nil)
////	_assert.Nil(err)
////
////	bbClient, _ = createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_assert.Nil(err)
////	_, err = bbClient.SetTier(ctx, AccessTierHot, LeaseAccessConditions{})
////	_assert.Nil(err)
////
////	resp, err = bbClient.GetProperties(ctx, nil)
////	_assert.Nil(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToHot))
////
////	resp2, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.Nil(err)
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

func (s *azblobTestSuite) TestBlobClientPartsSASQueryTimes() {
	_assert := assert.New(s.T())
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
		_assert.Equal(parts.Scheme, "https")
		_assert.Equal(parts.Host, "myaccount.blob.core.windows.net")
		_assert.Equal(parts.ContainerName, "mycontainer")
		_assert.Equal(parts.BlobName, "mydirectory/myfile.txt")

		sas := parts.SAS
		_assert.Equal(sas.StartTime(), StartTimesExpected[i])
		_assert.Equal(sas.ExpiryTime(), ExpiryTimesExpected[i])

		_assert.Equal(parts.URL(), urlString)
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestDownloadBlockBlobUnexpectedEOF() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.Download(ctx, nil)
	_assert.Nil(err)

	// Verify that we can inject errors first.
	reader := resp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))

	_, err = ioutil.ReadAll(reader)
	_assert.NotNil(err)
	_assert.Equal(err.Error(), "unrecoverable error")

	// Then inject the retryable error.
	reader = resp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))

	buf, err := ioutil.ReadAll(reader)
	_assert.Nil(err)
	_assert.EqualValues(buf, []byte(blockBlobDefaultData))
}

//nolint
func InjectErrorInRetryReaderOptions(err error) RetryReaderOptions {
	return RetryReaderOptions{
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
