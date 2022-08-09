//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"crypto/md5"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
)

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClient() {
//	_require := require.New(s.T())
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
//	_require.Equal(blobURLParts.BlobName, blobName)
//	_require.Equal(blobURLParts.ContainerName, containerName)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_require.Nil(err)
//	correctURL := "https://" + accountName + "." + DefaultBlobEndpointSuffix + containerName + "/" + blobName
//	_require.Equal(bbClient.URL(), correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClientWithSnapshotAndSAS() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	_require.Nil(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	blobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_require.Nil(err)
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
//	_require.Nil(err)
//	correctURL := "https://" + accountName + DefaultBlobEndpointSuffix + containerName + "/" + blobName +
//		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
//	_require.Equal(blobURLParts, correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateBlobClientWithSnapshotAndSASUsingConnectionString() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClientFromConnectionString(nil, testAccountDefault, nil)
//	_require.Nil(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	blobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_require.Nil(err)
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
//	_require.Nil(err)
//	correctURL := "https://" + accountName + DefaultBlobEndpointSuffix + containerName + "/" + blobName +
//		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
//	_require.Equal(blobURLParts, correctURL)
//}

func waitForCopy(_require *require.Assertions, copyBlobClient *BlockBlobClient, blobCopyResponse BlobStartCopyFromURLResponse) {
	status := *blobCopyResponse.CopyStatus
	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		props, _ := copyBlobClient.GetProperties(ctx, nil)
		status = *props.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_require.Fail("If the copy takes longer than a minute, we will fail")
		}
	}
}

func (s *azblobTestSuite) TestBlobStartCopyDestEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	_require.Nil(err)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient, _ := getBlockBlobClient(anotherBlobName, containerClient)

	blobCopyResponse, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_require.Nil(err)
	waitForCopy(_require, copyBlobClient, blobCopyResponse)

	resp, err := copyBlobClient.Download(ctx, nil)
	_require.Nil(err)

	// Read the blob data to verify the copy
	data, err := io.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(len(blockBlobDefaultData)))
	_require.Equal(string(data), blockBlobDefaultData)
	_ = resp.Body(nil).Close()
}

func (s *azblobTestSuite) TestBlobStartCopyMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	_require.Nil(err)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient, _ := getBlockBlobClient(anotherBlobName, containerClient)

	metadata := make(map[string]string)
	metadata["Bla"] = "foo"
	options := BlobStartCopyOptions{
		Metadata: metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)
	waitForCopy(_require, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, metadata)
}

func (s *azblobTestSuite) TestBlobStartCopyMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	anotherBlobName := "copy" + blockBlobName
	copyBlobClient, _ := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the nil metadata passed later takes effect
	_, err = copyBlobClient.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_require.Nil(err)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_require.Nil(err)

	waitForCopy(_require, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Len(resp2.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobStartCopyMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient, _ := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the empty metadata passed later takes effect
	_, err = copyBlobClient.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_require.Nil(err)

	metadata := make(map[string]string)
	options := BlobStartCopyOptions{
		Metadata: metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	waitForCopy(_require, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Len(resp2.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobStartCopyMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blobName, containerClient)

	anotherBlobName := "copy" + generateBlobName(testName)
	copyBlobClient, _ := getBlockBlobClient(anotherBlobName, containerClient)

	metadata := make(map[string]string)
	metadata["I nvalid."] = "foo"
	options := BlobStartCopyOptions{
		Metadata: metadata,
	}
	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *azblobTestSuite) TestBlobStartCopySourceNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient, _ := getBlockBlobClient(anotherBlobName, containerClient)

	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), "not exist"), true)
}

func (s *azblobTestSuite) TestBlobStartCopySourcePrivate() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	_require.Nil(err)

	bbClient := createNewBlockBlob(_require, generateBlobName(testName), containerClient)

	serviceClient2, err := getServiceClient(_context.recording, testAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
		return
	}

	copyContainerClient := createNewContainer(_require, "cpyc"+containerName, serviceClient2)
	defer deleteContainer(_require, copyContainerClient)
	copyBlobName := "copyb" + generateBlobName(testName)
	copyBlobClient, _ := getBlockBlobClient(copyBlobName, copyContainerClient)

	if svcClient.URL() == serviceClient2.URL() {
		s.T().Skip("Test not valid because primary and secondary accounts are the same")
	}
	_, err = copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	validateStorageError(_require, err, StorageErrorCodeCannotVerifyCopySource)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobStartCopyUsingSASSrc() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	_require.Nil(err)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	// Create sas values for the source blob
	credential, err := getGenericCredential(nil, testAccountDefault)
	if err != nil {
		s.T().Fatal("Couldn't fetch credential because " + err.Error())
	}

	startTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	endTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
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
	sasURL, _ := NewBlobURLParts(bbClient.URL())
	sasURL.SAS = queryParams

	// Create a new container for the destination
	serviceClient2, err := getServiceClient(nil, testAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
	}

	copyContainerName := "copy" + generateContainerName(testName)
	copyContainerClient := createNewContainer(_require, copyContainerName, serviceClient2)
	defer deleteContainer(_require, copyContainerClient)

	copyBlobName := "copy" + generateBlobName(testName)
	copyBlobClient, _ := getBlockBlobClient(copyBlobName, copyContainerClient)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, sasURL.URL(), nil)
	_require.Nil(err)

	waitForCopy(_require, copyBlobClient, resp)

	downloadBlobOptions := BlobDownloadOptions{
		Offset: to.Ptr[int64](0),
		Count:  to.Ptr(int64(len(blockBlobDefaultData))),
	}
	resp2, err := copyBlobClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)

	data, err := io.ReadAll(resp2.RawResponse.Body)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(len(blockBlobDefaultData)))
	_require.Equal(string(data), blockBlobDefaultData)
	_ = resp2.Body(nil).Close()
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobStartCopyUsingSASDest() {
	_require := require.New(s.T())
	testName := s.T().Name()
	var svcClient *ServiceClient
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = getServiceClient(nil, testAccountDefault, nil)
		} else {
			svcClient, err = getServiceClientFromConnectionString(nil, testAccountDefault, nil)
		}
		_require.Nil(err)

		containerClient := createNewContainer(_require, generateContainerName(testName)+strconv.Itoa(i), svcClient)
		_, err := containerClient.SetAccessPolicy(ctx, nil)
		_require.Nil(err)

		blobClient := createNewBlockBlob(_require, generateBlobName(testName), containerClient)
		_, err = blobClient.Delete(ctx, nil)
		_require.Nil(err)

		deleteContainer(_require, containerClient)
	}
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbClient, _ := getBlockBlobClient(generateBlobName(testName), containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient, _ := getBlockBlobClient("dst"+generateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient, _ := getBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}

	destBlobClient, _ := getBlockBlobClient("dst"+generateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}
	destBlobClient, _ := getBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient, _ := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	randomEtag := "a"
	accessConditions := SourceModifiedAccessConditions{
		SourceIfMatch: &randomEtag,
	}
	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient, _ := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeSourceConditionNotMet)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfNoneMatch: to.Ptr("a"),
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient, _ := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopySourceIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobStartCopyOptions{
		SourceModifiedAccessConditions: &SourceModifiedAccessConditions{
			SourceIfNoneMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient, _ := getBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeSourceConditionNotMet)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	destBlobClient := createNewBlockBlob(_require, "dst"+bbName, containerClient) // The blob must exist to have a last-modified time
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	destBlobClient := createNewBlockBlob(_require, "dst"+bbName, containerClient) // The blob must exist to have a last-modified time

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	validateStorageError(_require, err, StorageErrorCodeTargetConditionNotMet)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	destBlobClient := createNewBlockBlob(_require, "dst"+bbName, containerClient)

	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	destBlobClient := createNewBlockBlob(_require, "dst"+bbName, containerClient)
	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	metadata := make(map[string]string)
	metadata["bla"] = "bla"
	_, err = destBlobClient.SetMetadata(ctx, metadata, nil)
	_require.Nil(err)

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeTargetConditionNotMet)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.SetMetadata(ctx, nil, nil) // SetMetadata chances the blob's etag
	_require.Nil(err)

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.Nil(err)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobStartCopyDestIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	destBlobName := "dest" + generateBlobName(testName)
	destBlobClient := createNewBlockBlob(_require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobStartCopyOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}

	_, err = destBlobClient.StartCopyFromURL(ctx, bbClient.URL(), &options)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeTargetConditionNotMet)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobAbortCopyInProgress() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(blockBlobName, containerClient)

	// Create a large blob that takes time to copy
	blobSize := 8 * 1024 * 1024
	blobReader, _ := getRandomDataAndReader(blobSize)
	_, err = bbClient.Upload(ctx, internal.NopCloser(blobReader), nil)
	_require.Nil(err)

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := ContainerSetAccessPolicyOptions{
		Access: &access,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions) // So that we don't have to create a SAS
	_require.Nil(err)

	// Must copy across accounts so it takes time to copy
	serviceClient2, err := getServiceClient(nil, testAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
	}

	copyContainerName := "copy" + generateContainerName(testName)
	copyContainerClient := createNewContainer(_require, copyContainerName, serviceClient2)

	copyBlobName := "copy" + generateBlobName(testName)
	copyBlobClient, _ := getBlockBlobClient(copyBlobName, copyContainerClient)

	defer deleteContainer(_require, copyContainerClient)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, bbClient.URL(), nil)
	_require.Nil(err)
	_require.Equal(*resp.CopyStatus, CopyStatusTypePending)
	_require.NotNil(resp.CopyID)

	_, err = copyBlobClient.AbortCopyFromURL(ctx, *resp.CopyID, nil)
	if err != nil {
		// If the error is nil, the test continues as normal.
		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
		validateStorageError(_require, err, StorageErrorCodeNoPendingCopyOperation)
		_require.Error(errors.New("the test failed because the copy completed because it was aborted"))
	}

	resp2, _ := copyBlobClient.GetProperties(ctx, nil)
	_require.Equal(*resp2.CopyStatus, CopyStatusTypeAborted)
}

func (s *azblobTestSuite) TestBlobAbortCopyNoCopyStarted() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	copyBlobClient, _ := getBlockBlobClient(blockBlobName, containerClient)

	_, err = copyBlobClient.AbortCopyFromURL(ctx, "copynotstarted", nil)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeInvalidQueryParameterValue)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	createBlobSnapshotOptions := BlobCreateSnapshotOptions{
		Metadata: basicMetadata,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	// Since metadata is specified on the snapshot, the snapshot should have its own metadata different from the (empty) metadata on the source
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	// In this case, because no metadata was specified, it should copy the basicMetadata from the source
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSnapshotMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	createBlobSnapshotOptions := BlobCreateSnapshotOptions{
		Metadata: map[string]string{"Invalid Field!": "value"},
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_require.NotNil(err)
	_require.Contains(err.Error(), invalidHeaderErrorSubstring)
}

func (s *azblobTestSuite) TestBlobSnapshotBlobNotExist() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(ctx, nil)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotOfSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	snapshotString, err := time.Parse(SnapshotTimeFormat, "2021-01-01T01:01:01.0000000Z")
	_require.Nil(err)
	snapshotURL, _ := bbClient.WithSnapshot(snapshotString.String())
	// The library allows the server to handle the snapshot of snapshot error
	_, err = snapshotURL.CreateSnapshot(ctx, nil)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeInvalidQueryParameterValue)
}

func (s *azblobTestSuite) TestBlobSnapshotIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	_require.Nil(err)
	_require.NotEqual(*resp.Snapshot, "") // i.e. The snapshot time is not zero. If the service gives us back a snapshot time, it successfully created a snapshot
}

func (s *azblobTestSuite) TestBlobSnapshotIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	_require.Nil(err)
	_require.NotEqual(*resp.Snapshot, "")
}

func (s *azblobTestSuite) TestBlobSnapshotIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	resp2, err := bbClient.CreateSnapshot(ctx, &options)
	_require.Nil(err)
	_require.NotEqual(*resp2.Snapshot, "")
}

func (s *azblobTestSuite) TestBlobSnapshotIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	randomEtag := "garbage"
	access := ModifiedAccessConditions{
		IfMatch: &randomEtag,
	}
	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSnapshotIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	randomEtag := "garbage"
	access := ModifiedAccessConditions{
		IfNoneMatch: &randomEtag,
	}
	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &options)
	_require.Nil(err)
	_require.NotEqual(*resp.Snapshot, "")
}

func (s *azblobTestSuite) TestBlobSnapshotIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobCreateSnapshotOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}
	_, err = bbClient.CreateSnapshot(ctx, &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient, _ := containerClient.NewBlobClient(blobName)

	_, err = bbClient.Download(ctx, nil)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataNegativeOffset() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	options := BlobDownloadOptions{
		Offset: to.Ptr[int64](-1),
	}
	_, err = bbClient.Download(ctx, &options)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataOffsetOutOfRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	options := BlobDownloadOptions{
		Offset: to.Ptr(int64(len(blockBlobDefaultData))),
	}
	_, err = bbClient.Download(ctx, &options)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeInvalidRange)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	options := BlobDownloadOptions{
		Count: to.Ptr[int64](-2),
	}
	_, err = bbClient.Download(ctx, &options)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	options := BlobDownloadOptions{
		Count: to.Ptr[int64](0),
	}
	resp, err := bbClient.Download(ctx, &options)
	_require.Nil(err)

	// Specifying a count of 0 results in the value being ignored
	data, err := io.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountExact() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	count := int64(len(blockBlobDefaultData))
	options := BlobDownloadOptions{
		Count: &count,
	}
	resp, err := bbClient.Download(ctx, &options)
	_require.Nil(err)

	data, err := io.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataCountOutOfRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	options := BlobDownloadOptions{
		Count: to.Ptr(int64((len(blockBlobDefaultData)) * 2)),
	}
	resp, err := bbClient.Download(ctx, &options)
	_require.Nil(err)

	data, err := io.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataEmptyRangeStruct() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	options := BlobDownloadOptions{
		Count:  to.Ptr[int64](0),
		Offset: to.Ptr[int64](0),
	}
	resp, err := bbClient.Download(ctx, &options)
	_require.Nil(err)

	data, err := io.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobDownloadDataContentMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	options := BlobDownloadOptions{
		Count:              to.Ptr[int64](3),
		Offset:             to.Ptr[int64](10),
		RangeGetContentMD5: to.Ptr(true),
	}
	resp, err := bbClient.Download(ctx, &options)
	_require.Nil(err)
	mdf := md5.Sum([]byte(blockBlobDefaultData)[10:13])
	_require.Equal(resp.ContentMD5, mdf[:])
}

func (s *azblobTestSuite) TestBlobDownloadDataIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}

	resp, err := bbClient.Download(ctx, &options)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	bbClient, _ := getBlockBlobClient(generateBlobName(testName), containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	resp, err := bbClient.Download(ctx, &BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		}},
	})
	_require.Nil(err)
	_require.Equal(*resp.ErrorCode, string(StorageErrorCodeConditionNotMet))

}

func (s *azblobTestSuite) TestBlobDownloadDataIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	options := BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.Download(ctx, &options)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}
	_, err = bbClient.Download(ctx, &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.Download(ctx, &options)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	options := BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_require.Nil(err)

	_, err = bbClient.Download(ctx, &options)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDownloadDataIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	access := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &access},
	}

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_require.Nil(err)

	resp2, err := bbClient.Download(ctx, &options)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(len(blockBlobDefaultData)))
}

func (s *azblobTestSuite) TestBlobDownloadDataIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	options := BlobDownloadOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}

	resp2, err := bbClient.Download(ctx, &options)
	_require.Nil(err)
	_require.Equal(*resp2.ErrorCode, string(StorageErrorCodeConditionNotMet))
}

func (s *azblobTestSuite) TestBlobDeleteNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient, _ := containerClient.NewBlockBlobClient(blockBlobName)

	_, err = bbClient.Delete(ctx, nil)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobDeleteSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)

	_, err = snapshotURL.Delete(ctx, nil)
	_require.Nil(err)

	validateBlobDeleted(_require, snapshotURL.BlobClient)
}

//
////func (s *azblobTestSuite) TestBlobDeleteSnapshotsInclude() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := bbClient.CreateSnapshot(ctx, nil)
////	_require.Nil(err)
////
////	deleteSnapshots := DeleteSnapshotsOptionInclude
////	_, err = bbClient.Delete(ctx, &BlobDeleteOptions{
////		DeleteSnapshots: &deleteSnapshots,
////	})
////	_require.Nil(err)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
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
////	defer deleteContainer(_require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := bbClient.CreateSnapshot(ctx, nil)
////	_require.Nil(err)
////	deleteSnapshot := DeleteSnapshotsOptionOnly
////	deleteBlobOptions := BlobDeleteOptions{
////		DeleteSnapshots: &deleteSnapshot,
////	}
////	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
////	_require.Nil(err)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
////		Include: include,
////	}
////	blobs, errChan := containerClient.ListBlobsFlat(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(blobs, chk.HasLen, 1)
////	_assert(*(<-blobs).Snapshot == "", chk.Equals, true)
////}

func (s *azblobTestSuite) TestBlobDeleteSnapshotsNoneWithSnapshots() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	_, err = bbClient.Delete(ctx, nil)
	_require.NotNil(err)
}

func validateBlobDeleted(_require *require.Assertions, bbClient BlobClient) {
	_, err := bbClient.GetProperties(ctx, nil)
	_require.NotNil(err)

	var storageError *StorageError
	_require.Equal(true, errors.As(err, &storageError))
	_require.Equal(storageError.ErrorCode, StorageErrorCodeBlobNotFound)
}

func (s *azblobTestSuite) TestBlobDeleteIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobDeleteIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobDeleteIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(ctx, nil)

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	etag := resp.ETag

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_require.Nil(err)

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobDeleteIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(ctx, nil)
	etag := resp.ETag
	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_require.Nil(err)

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobDeleteIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(ctx, nil)
	etag := resp.ETag

	deleteBlobOptions := BlobDeleteOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: etag},
		},
	}
	_, err = bbClient.Delete(ctx, &deleteBlobOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.NotNil(err)
	var storageError *StorageError
	_require.Equal(errors.As(err, &storageError), true)
	_require.Equal(storageError.response.StatusCode, 304) // No service code returned for a HEAD
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, basicMetadata)
}

//func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceFalse() {
//	// TODO: Not Working
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	blockBlobName := generateBlobName(testName)
//	bbClient := getBlockBlobClient(blockBlobName, containerClient)
//
//	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
//
//	_require.Nil(err)
//	_require.Equal(cResp.RawResponse.StatusCode, 201)
//
//	currentTime := getRelativeTimeFromAnchor(cResp.Date,-10)
//
//	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
//	_require.Nil(err)
//
//	getBlobPropertiesOptions := BlobGetPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_require.NotNil(err)
//}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobGetPropsOnMissingBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbClient, _ := containerClient.NewBlobClient("MISSING")

	_, err = bbClient.GetProperties(ctx, nil)
	_require.NotNil(err)
	var storageError *StorageError
	_require.Equal(errors.As(err, &storageError), true)
	_require.Equal(storageError.response.StatusCode, 404)
	_require.Equal(storageError.response.Header.Get("x-ms-error-code"), string(StorageErrorCodeBlobNotFound))
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	eTag := "garbage"
	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.NotNil(err)
	var storageError *StorageError
	_require.Equal(errors.As(err, &storageError), true)
	_require.Equal(storageError.response.StatusCode, 412)
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	eTag := "garbage"
	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
		},
	}
	resp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobGetPropsAndMetadataIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(ctx, nil, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}
	_, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.NotNil(err)
	var storageError *StorageError
	_require.Equal(errors.As(err, &storageError), true)
	_require.Equal(storageError.response.StatusCode, 304)
}

func (s *azblobTestSuite) TestBlobSetPropertiesBasic() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, basicHeaders, nil)
	_require.Nil(err)

	resp, _ := bbClient.GetProperties(ctx, nil)
	h := resp.GetHTTPHeaders()
	_require.EqualValues(h, basicHeaders)
}

func (s *azblobTestSuite) TestBlobSetPropertiesEmptyValue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	contentType := to.Ptr("my_type")
	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentType: contentType}, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp.ContentType, contentType)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{}, nil)
	_require.Nil(err)

	resp, err = bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.ContentType)
}

func validatePropertiesSet(_require *require.Assertions, bbClient *BlockBlobClient, disposition string) {
	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentDisposition, disposition)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: to.Ptr("garbage")}})
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: to.Ptr("garbage")}})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *azblobTestSuite) TestBlobSetPropertiesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	_, err = bbClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&BlobSetHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}})
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
	_require.Nil(err)

	_, err = bbClient.SetMetadata(ctx, nil, nil)
	_require.Nil(err)

	blobGetResp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Len(blobGetResp.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobSetMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
	_require.Nil(err)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Len(resp.Metadata, 0)
}

func (s *azblobTestSuite) TestBlobSetMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(ctx, map[string]string{"Invalid field!": "value"}, nil)
	_require.NotNil(err)
	_require.Contains(err.Error(), invalidHeaderErrorSubstring)
	//_require.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func validateMetadataSet(_require *require.Assertions, bbClient *BlockBlobClient) {
	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient, _ := getBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	eTag := "garbage"
	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	eTag := "garbage"
	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobSetMetadataIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
	}
	_, err = bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

// nolint
func testBlobServiceClientDeleteImpl(_ *require.Assertions, _ *ServiceClient) error {
	//containerClient := createNewContainer(_require, "gocblobserviceclientdeleteimpl", svcClient)
	//defer deleteContainer(_require, containerClient)
	//bbClient := createNewBlockBlob(_require, "goblobserviceclientdeleteimpl", containerClient)
	//
	//_, err := bbClient.Delete(ctx, nil)
	//_require.Nil(err) // This call will not have errors related to slow update of service properties, so we assert.
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
	//_require.Equal(resp.BlobType, BlobTypeBlockBlob) // We could check any property. This is just to double check it was undeleted.
	return nil
}

func (s *azblobTestSuite) TestBlobServiceClientDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	code := 404
	runTestRequiringServiceProperties(_require, svcClient, string(rune(code)), enableSoftDelete, testBlobServiceClientDeleteImpl, disableSoftDelete)
}

func setAndCheckBlobTier(_require *require.Assertions, bbClient BlobClient, tier AccessTier) {
	_, err := bbClient.SetTier(ctx, tier, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.AccessTier, string(tier))
}

func (s *azblobTestSuite) TestBlobSetTierAllTiers() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	setAndCheckBlobTier(_require, bbClient.BlobClient, AccessTierHot)
	setAndCheckBlobTier(_require, bbClient.BlobClient, AccessTierCool)
	setAndCheckBlobTier(_require, bbClient.BlobClient, AccessTierArchive)

	premiumServiceClient, err := getServiceClient(_context.recording, testAccountPremium, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	premContainerName := "prem" + generateContainerName(testName)
	premContainerClient := createNewContainer(_require, premContainerName, premiumServiceClient)
	defer deleteContainer(_require, premContainerClient)

	pbClient := createNewPageBlob(_require, blockBlobName, premContainerClient)

	setAndCheckBlobTier(_require, pbClient.BlobClient, AccessTierP4)
	setAndCheckBlobTier(_require, pbClient.BlobClient, AccessTierP6)
	setAndCheckBlobTier(_require, pbClient.BlobClient, AccessTierP10)
	setAndCheckBlobTier(_require, pbClient.BlobClient, AccessTierP20)
	setAndCheckBlobTier(_require, pbClient.BlobClient, AccessTierP30)
	setAndCheckBlobTier(_require, pbClient.BlobClient, AccessTierP40)
	setAndCheckBlobTier(_require, pbClient.BlobClient, AccessTierP50)
}

//
////func (s *azblobTestSuite) TestBlobTierInferred() {
////	svcClient, err := getPremiumserviceClient()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	bbClient, _ := createNewPageBlob(c, containerClient)
////
////	resp, err := bbClient.GetProperties(ctx, nil)
////	_require.Nil(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "true")
////
////	resp2, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.NotNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTier, chk.Not(chk.Equals), "")
////
////	_, err = bbClient.SetTier(ctx, AccessTierP4, LeaseAccessConditions{})
////	_require.Nil(err)
////
////	resp, err = bbClient.GetProperties(ctx, nil)
////	_require.Nil(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "")
////
////	resp2, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
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
////	defer deleteContainer(_require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_require.Nil(err)
////	_, err = bbClient.SetTier(ctx, AccessTierCool, LeaseAccessConditions{})
////	_require.Nil(err)
////
////	resp, err := bbClient.GetProperties(ctx, nil)
////	_require.Nil(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToCool))
////
////	resp2, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToCool)
////
////	// delete first blob
////	_, err = bbClient.Delete(ctx, DeleteSnapshotsOptionNone, nil)
////	_require.Nil(err)
////
////	bbClient, _ = createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_require.Nil(err)
////	_, err = bbClient.SetTier(ctx, AccessTierHot, LeaseAccessConditions{})
////	_require.Nil(err)
////
////	resp, err = bbClient.GetProperties(ctx, nil)
////	_require.Nil(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToHot))
////
////	resp2, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
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
////	defer deleteContainer(_require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierType("garbage"), LeaseAccessConditions{})
////	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
////}
////

func (s *azblobTestSuite) TestBlobClientPartsSASQueryTimes() {
	_require := require.New(s.T())
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

		parts, _ := NewBlobURLParts(urlString)
		_require.Equal(parts.Scheme, "https")
		_require.Equal(parts.Host, "myaccount.blob.core.windows.net")
		_require.Equal(parts.ContainerName, "mycontainer")
		_require.Equal(parts.BlobName, "mydirectory/myfile.txt")

		sas := parts.SAS
		_require.Equal(sas.StartTime(), StartTimesExpected[i])
		_require.Equal(sas.ExpiryTime(), ExpiryTimesExpected[i])

		_require.Equal(parts.URL(), urlString)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestDownloadBlockBlobUnexpectedEOF() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.Download(ctx, nil)
	_require.Nil(err)

	// Verify that we can inject errors first.
	reader := resp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))

	_, err = io.ReadAll(reader)
	_require.NotNil(err)
	_require.Equal(err.Error(), "unrecoverable error")

	// Then inject the retryable error.
	reader = resp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))

	buf, err := io.ReadAll(reader)
	_require.Nil(err)
	_require.EqualValues(buf, []byte(blockBlobDefaultData))
}

// nolint
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
