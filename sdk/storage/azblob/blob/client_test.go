//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running blob Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &BlobRecordedTestsSuite{})
		suite.Run(t, &BlobUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &BlobRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &BlobRecordedTestsSuite{})
	}
}

func (s *BlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *BlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *BlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *BlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type BlobRecordedTestsSuite struct {
	suite.Suite
}

type BlobUnrecordedTestsSuite struct {
	suite.Suite
}

func (s *BlobUnrecordedTestsSuite) TestCreateBlobClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	blobURLParts, err := blob.ParseURL(bbClient.URL())
	_require.Nil(err)
	_require.Equal(blobURLParts.BlobName, blobName)
	_require.Equal(blobURLParts.ContainerName, containerName)

	accountName, err := testcommon.GetRequiredEnv(testcommon.AccountNameEnvVar)
	_require.Nil(err)
	correctURL := "https://" + accountName + "." + testcommon.DefaultBlobEndpointSuffix + containerName + "/" + blobName
	_require.Equal(bbClient.URL(), correctURL)
}

func (s *BlobUnrecordedTestsSuite) TestCreateBlobClientWithSnapshotAndSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    currentTime,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	parts, err := blob.ParseURL(bbClient.URL())
	_require.Nil(err)
	parts.SAS = sasQueryParams
	parts.Snapshot = currentTime.Format(blob.SnapshotTimeFormat)
	blobURLParts := parts.String()

	// The snapshot format string is taken from the snapshotTimeFormat value in parsing_urls.go. The field is not public, so
	// it is copied here
	accountName, err := testcommon.GetRequiredEnv(testcommon.AccountNameEnvVar)
	_require.Nil(err)
	correctURL := "https://" + accountName + "." + testcommon.DefaultBlobEndpointSuffix + containerName + "/" + blobName +
		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
	_require.Equal(blobURLParts, correctURL)
}

func (s *BlobUnrecordedTestsSuite) TestCreateBlobClientWithSnapshotAndSASUsingConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)
	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    currentTime,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	parts, err := blob.ParseURL(bbClient.URL())
	_require.Nil(err)
	parts.SAS = sasQueryParams
	parts.Snapshot = currentTime.Format(blob.SnapshotTimeFormat)
	blobURLParts := parts.String()

	// The snapshot format string is taken from the snapshotTimeFormat value in parsing_urls.go. The field is not public, so
	// it is copied here
	accountName, err := testcommon.GetRequiredEnv(testcommon.AccountNameEnvVar)
	_require.Nil(err)
	correctURL := "https://" + accountName + "." + testcommon.DefaultBlobEndpointSuffix + containerName + "/" + blobName +
		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
	_require.Equal(blobURLParts, correctURL)
}

func waitForCopy(_require *require.Assertions, copyBlobClient *blockblob.Client, blobCopyResponse blob.StartCopyFromURLResponse) {
	status := *blobCopyResponse.CopyStatus
	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
	start := time.Now()
	for status != blob.CopyStatusTypeSuccess {
		props, _ := copyBlobClient.GetProperties(context.Background(), nil)
		status = *props.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_require.Fail("If the copy takes longer than a minute, we will fail")
		}
	}
}

func (s *BlobUnrecordedTestsSuite) TestCopyBlockBlobFromUrlSourceContentMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	contentMD5 := md5.Sum(content)
	body := bytes.NewReader(content)

	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")

	// Prepare source bbClient for copy.
	_, err = srcBlob.Upload(context.Background(), streaming.NopCloser(body), nil)
	_require.Nil(err)

	expiryTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	if err != nil {
		s.T().Fatal("Couldn't fetch credential because " + err.Error())
	}

	// Get source blob url with SAS for StageFromURL.
	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	// Invoke CopyFromURL.
	sourceContentMD5 := contentMD5[:]
	resp, err := destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &blob.CopyFromURLOptions{
		SourceContentMD5: sourceContentMD5,
	})
	_require.Nil(err)
	_require.EqualValues(resp.ContentMD5, sourceContentMD5)

	// Provide bad MD5 and make sure the copy fails
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	resp, err = destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &blob.CopyFromURLOptions{
		SourceContentMD5: badMD5,
	})
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	_, err = containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := testcommon.GetBlockBlobClient(anotherBlobName, containerClient)

	blobCopyResponse, err := copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), nil)
	_require.Nil(err)
	waitForCopy(_require, copyBlobClient, blobCopyResponse)

	resp, err := copyBlobClient.DownloadStream(context.Background(), nil)
	_require.Nil(err)

	// Read the blob data to verify the copy
	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(len(testcommon.BlockBlobDefaultData)))
	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
	_ = resp.Body.Close()
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	_, err = containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := testcommon.GetBlockBlobClient(anotherBlobName, containerClient)

	resp, err := copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &blob.StartCopyFromURLOptions{
		Metadata: testcommon.BasicMetadata,
	})
	_require.Nil(err)
	waitForCopy(_require, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	anotherBlobName := "copy" + blockBlobName
	copyBlobClient := testcommon.GetBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the nil metadata passed later takes effect
	_, err = copyBlobClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_require.Nil(err)

	resp, err := copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), nil)
	_require.Nil(err)

	waitForCopy(_require, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp2.Metadata, 0)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := testcommon.GetBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata, so we ensure the empty metadata passed later takes effect
	_, err = copyBlobClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_require.Nil(err)

	metadata := make(map[string]*string)
	options := blob.StartCopyFromURLOptions{
		Metadata: metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	waitForCopy(_require, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp2.Metadata, 0)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)

	anotherBlobName := "copy" + testcommon.GenerateBlobName(testName)
	copyBlobClient := testcommon.GetBlockBlobClient(anotherBlobName, containerClient)

	metadata := make(map[string]*string)
	metadata["I nvalid."] = to.Ptr("foo")
	options := blob.StartCopyFromURLOptions{
		Metadata: metadata,
	}
	_, err = copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring), true)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := testcommon.GetBlockBlobClient(anotherBlobName, containerClient)

	_, err = copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), nil)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), "not exist"), true)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourcePrivate() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	serviceClient2, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
		return
	}

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err = containerClient.SetAccessPolicy(context.Background(), nil)
	_require.Nil(err)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	copyContainerClient := testcommon.CreateNewContainer(context.Background(), _require, "cpyc"+containerName, serviceClient2)
	defer testcommon.DeleteContainer(context.Background(), _require, copyContainerClient)
	copyBlobName := "copyb" + testcommon.GenerateBlobName(testName)
	copyBlobClient := testcommon.GetBlockBlobClient(copyBlobName, copyContainerClient)

	if svcClient.URL() == serviceClient2.URL() {
		s.T().Skip("Test not valid because primary and secondary accounts are the same")
	}
	_, err = copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CannotVerifyCopySource)
}

//func (s *BlobUnrecordedTestsSuite) TestBlobStartCopyUsingSASSrc() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(),testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	_, err = containerClient.SetAccessPolicy(context.Background(), nil)
//	_require.Nil(err)
//
//	blockBlobName := testcommon.GenerateBlobName(testName)
//	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)
//
//	// Create sas values for the source blob
//	credential, err := testcommon.GetGenericSharedKeyCredential(nil, testcommon.TestAccountDefault)
//	if err != nil {
//		s.T().Fatal("Couldn't fetch credential because " + err.Error())
//	}
//
//	startTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//	_require.Nil(err)
//	endTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_require.Nil(err)
//	serviceSASValues := azblob.SASSignatureValues{
//		StartTime:     startTime,
//		ExpiryTime:    endTime,
//		Permissions:   BlobSASPermissions{Read: true, Write: true}.String(),
//		ContainerName: containerName,
//		BlobName:      blockBlobName}
//	queryParams, err := serviceSASValues.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	// Create URLs to the destination blob with sas parameters
//	sasURL, _ := NewBlobURLParts(bbClient.URL())
//	sasURL.SAS = queryParams
//
//	// Create a new container for the destination
//	serviceClient2, err := testcommon.GetServiceClient(s.T(),testcommon.TestAccountSecondary, nil)
//	if err != nil {
//		s.T().Skip(err.Error())
//	}
//
//	copyContainerName := "copy" + testcommon.GenerateContainerName(testName)
//	copyContainerClient := testcommon.CreateNewContainer(context.Background(), _require, copyContainerName, serviceClient2)
//	defer testcommon.DeleteContainer(context.Background(), _require, copyContainerClient)
//
//	copyBlobName := "copy" + testcommon.GenerateBlobName(testName)
//	copyBlobClient := testcommon.GetBlockBlobClient(copyBlobName, copyContainerClient)
//
//	resp, err := copyBlobClient.StartCopyFromURL(context.Background(), sasURL.URL(), nil)
//	_require.Nil(err)
//
//	waitForCopy(_require, copyBlobClient, resp)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		Offset: to.Ptr[int64](0),
//		Count:  to.Ptr(int64(len(testcommon.BlockBlobDefaultData))),
//	}
//	resp2, err := copyBlobClient.DownloadStream(context.Background(), &downloadBlobOptions)
//	_require.Nil(err)
//
//	data, err := io.ReadAll(resp2.Body(nil))
//	_require.Nil(err)
//	_require.Equal(*resp2.ContentLength, int64(len(testcommon.BlockBlobDefaultData)))
//	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
//	_ = resp2.Body(nil).Close()
//}

func (s *BlobUnrecordedTestsSuite) TestBlobStartCopyUsingSASDest() {
	_require := require.New(s.T())
	testName := s.T().Name()
	var svcClient *service.Client
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
		} else {
			svcClient, err = testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
		}
		_require.Nil(err)

		containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+strconv.Itoa(i), svcClient)
		_, err := containerClient.SetAccessPolicy(context.Background(), nil)
		_require.Nil(err)

		blobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
		_, err = blobClient.Delete(context.Background(), nil)
		_require.Nil(err)

		testcommon.DeleteContainer(context.Background(), _require, containerClient)
	}
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient := testcommon.GetBlockBlobClient("dst"+testcommon.GenerateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)
	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfModifiedSince: &currentTime,
		},
	}

	destBlobClient := testcommon.GetBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)
	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}

	destBlobClient := testcommon.GetBlockBlobClient("dst"+testcommon.GenerateBlobName(testName), containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)
	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfUnmodifiedSince: &currentTime,
		},
	}
	destBlobClient := testcommon.GetBlockBlobClient("dst"+blobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.GetBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	randomEtag := azcore.ETag("a")
	accessConditions := blob.SourceModifiedAccessConditions{
		SourceIfMatch: &randomEtag,
	}
	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.GetBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SourceConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfNoneMatch: to.Ptr(azcore.ETag("a")),
		},
	}

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.GetBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopySourceIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.StartCopyFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfNoneMatch: resp.ETag,
		},
	}

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.GetBlockBlobClient(destBlobName, containerClient)
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SourceConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "dst"+bbName, containerClient) // The blob must exist to have a last-modified time
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)

	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "dst"+bbName, containerClient) // The blob must exist to have a last-modified time

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)
	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.TargetConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "dst"+bbName, containerClient)

	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	_, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "dst"+bbName, containerClient)
	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}

	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	resp, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	metadata := make(map[string]*string)
	metadata["bla"] = to.Ptr("bla")
	_, err = destBlobClient.SetMetadata(context.Background(), metadata, nil)
	_require.Nil(err)

	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.TargetConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}

	_, err = destBlobClient.SetMetadata(context.Background(), nil, nil) // SetMetadata chances the blob's etag
	_require.Nil(err)

	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.Nil(err)

	resp, err = destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobStartCopyDestIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	destBlobName := "dest" + testcommon.GenerateBlobName(testName)
	destBlobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, destBlobName, containerClient)
	resp, err := destBlobClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.StartCopyFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}

	_, err = destBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), &options)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.TargetConditionNotMet)
}

func (s *BlobUnrecordedTestsSuite) TestBlobAbortCopyInProgress() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	// Create a large blob that takes time to copy
	blobSize := 8 * 1024 * 1024
	blobReader, _ := testcommon.GetDataAndReader(testName, blobSize)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(blobReader), nil)
	_require.Nil(err)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: to.Ptr(container.PublicAccessTypeBlob),
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions) // So that we don't have to create a SAS
	_require.Nil(err)

	// Must copy across accounts so it takes time to copy
	serviceClient2, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSecondary, nil)
	if err != nil {
		s.T().Skip(err.Error())
	}

	copyContainerName := "copy" + testcommon.GenerateContainerName(testName)
	copyContainerClient := testcommon.CreateNewContainer(context.Background(), _require, copyContainerName, serviceClient2)

	copyBlobName := "copy" + testcommon.GenerateBlobName(testName)
	copyBlobClient := testcommon.GetBlockBlobClient(copyBlobName, copyContainerClient)

	defer testcommon.DeleteContainer(context.Background(), _require, copyContainerClient)

	resp, err := copyBlobClient.StartCopyFromURL(context.Background(), bbClient.URL(), nil)
	_require.Nil(err)
	_require.Equal(*resp.CopyStatus, blob.CopyStatusTypePending)
	_require.NotNil(resp.CopyID)

	_, err = copyBlobClient.AbortCopyFromURL(context.Background(), *resp.CopyID, nil)
	if err != nil {
		// If the error is nil, the test continues as normal.
		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.NoPendingCopyOperation)
		_require.Error(errors.New("the test failed because the copy completed because it was aborted"))
	}

	resp2, _ := copyBlobClient.GetProperties(context.Background(), nil)
	_require.Equal(*resp2.CopyStatus, blob.CopyStatusTypeAborted)
}

func (s *BlobRecordedTestsSuite) TestBlobAbortCopyNoCopyStarted() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	copyBlobClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	_, err = copyBlobClient.AbortCopyFromURL(context.Background(), "copynotstarted", nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidQueryParameterValue)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	createBlobSnapshotOptions := blob.CreateSnapshotOptions{Metadata: testcommon.BasicMetadata}
	resp, err := bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	// Since metadata is specified on the snapshot, the snapshot should have its own metadata different from the (empty) metadata on the source
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Nil(err)

	resp, err := bbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	// In this case, because no metadata was specified, it should copy the testcommon.BasicMetadata from the source
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Nil(err)

	resp, err := bbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	createBlobSnapshotOptions := blob.CreateSnapshotOptions{
		Metadata: map[string]*string{"Invalid Field!": to.Ptr("value")},
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	_require.NotNil(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotBlobNotExist() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(context.Background(), nil)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotOfSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	snapshotString, err := time.Parse(blob.SnapshotTimeFormat, "2021-01-01T01:01:01.0000000Z")
	_require.Nil(err)
	snapshotURL, _ := bbClient.WithSnapshot(snapshotString.String())
	// The library allows the server to handle the snapshot of snapshot error
	_, err = snapshotURL.CreateSnapshot(context.Background(), nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidQueryParameterValue)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := bbClient.CreateSnapshot(context.Background(), &options)
	_require.Nil(err)
	_require.NotEqual(*resp.Snapshot, "") // i.e. The snapshot time is not zero. If the service gives us back a snapshot time, it successfully created a snapshot
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)
	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := bbClient.CreateSnapshot(context.Background(), &options)
	_require.Nil(err)
	_require.NotEqual(*resp.Snapshot, "")
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		}},
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := bbClient.CreateSnapshot(context.Background(), &options)
	_require.Nil(err)
	_require.NotEqual(*resp2.Snapshot, "")
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr(azcore.ETag("garbage")),
			},
		},
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	randomEtag := azcore.ETag("garbage")
	access := blob.ModifiedAccessConditions{
		IfNoneMatch: &randomEtag,
	}
	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &access,
		},
	}
	resp, err := bbClient.CreateSnapshot(context.Background(), &options)
	_require.Nil(err)
	_require.NotEqual(*resp.Snapshot, "")
}

func (s *BlobRecordedTestsSuite) TestBlobSnapshotIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.CreateSnapshotOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := containerClient.NewBlobClient(blobName)

	_, err = bbClient.DownloadStream(context.Background(), nil)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataNegativeOffset() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{
		Range: blob.HTTPRange{
			Offset: -1,
		},
	}
	_, err = bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataOffsetOutOfRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{
		Range: blob.HTTPRange{
			Offset: int64(len(testcommon.BlockBlobDefaultData)),
		},
	}
	_, err = bbClient.DownloadStream(context.Background(), &options)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidRange)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataCountNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{
		Range: blob.HTTPRange{
			Count: -2,
		},
	}
	_, err = bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataCountZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{}
	resp, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)

	// Specifying a count of 0 results in the value being ignored
	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataCountExact() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{
		Range: blob.HTTPRange{
			Count: int64(len(testcommon.BlockBlobDefaultData)),
		},
	}
	resp, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)

	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataCountOutOfRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{
		Range: blob.HTTPRange{
			Count: int64((len(testcommon.BlockBlobDefaultData)) * 2),
		},
	}
	resp, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)

	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataEmptyRangeStruct() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{}
	resp, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)

	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataContentMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	options := blob.DownloadStreamOptions{
		Range: blob.HTTPRange{
			Count:  3,
			Offset: 10,
		},
		RangeGetContentMD5: to.Ptr(true),
	}
	resp, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
	mdf := md5.Sum([]byte(testcommon.BlockBlobDefaultData)[10:13])
	_require.Equal(resp.ContentMD5, mdf[:])
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	options := blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(len(testcommon.BlockBlobDefaultData)))
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	resp, err := bbClient.DownloadStream(context.Background(), &blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		}},
	})
	_require.Nil(err)
	//testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
	_require.Equal(*resp.ErrorCode, string(bloberror.ConditionNotMet))
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	options := blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(len(testcommon.BlockBlobDefaultData)))
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	access := blob.ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &access},
	}
	_, err = bbClient.DownloadStream(context.Background(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(len(testcommon.BlockBlobDefaultData)))
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	options := blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}

	_, err = bbClient.SetMetadata(context.Background(), nil, nil)
	_require.Nil(err)

	_, err = bbClient.DownloadStream(context.Background(), &options)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	options := blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		}},
	}

	_, err = bbClient.SetMetadata(context.Background(), nil, nil)
	_require.Nil(err)

	resp2, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(len(testcommon.BlockBlobDefaultData)))
}

func (s *BlobRecordedTestsSuite) TestBlobDownloadDataIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	options := blob.DownloadStreamOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}

	resp2, err := bbClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
	_require.Equal(*resp2.ErrorCode, string(bloberror.ConditionNotMet))
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blockBlobName)

	_, err = bbClient.Delete(context.Background(), nil)
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)

	_, err = snapshotURL.Delete(context.Background(), nil)
	_require.Nil(err)

	validateBlobDeleted(_require, snapshotURL)
}

//
////func (s *BlobRecordedTestsSuite) TestBlobDeleteSnapshotsInclude() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := testcommon.CreateNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := bbClient.CreateSnapshot(context.Background(), nil)
////	_require.Nil(err)
////
////	deleteSnapshots := DeleteSnapshotsOptionInclude
////	_, err = bbClient.Delete(context.Background(), &BlobDeleteOptions{
////		DeleteSnapshots: &deleteSnapshots,
////	})
////	_require.Nil(err)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
////		Include: include,
////	}
////	blobs, errChan := containerClient.NewListBlobsFlatPager(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(<- blobs, chk.HasLen, 0)
////}
//
////func (s *BlobRecordedTestsSuite) TestBlobDeleteSnapshotsOnly() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := testcommon.CreateNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := bbClient.CreateSnapshot(context.Background(), nil)
////	_require.Nil(err)
////	deleteSnapshot := DeleteSnapshotsOptionOnly
////	deleteBlobOptions := blob.DeleteOptions{
////		DeleteSnapshots: &deleteSnapshot,
////	}
////	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
////	_require.Nil(err)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
////		Include: include,
////	}
////	blobs, errChan := containerClient.NewListBlobsFlatPager(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(blobs, chk.HasLen, 1)
////	_assert(*(<-blobs).Snapshot == "", chk.Equals, true)
////}

func (s *BlobRecordedTestsSuite) TestBlobDeleteSnapshotsNoneWithSnapshots() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	_, err = bbClient.Delete(context.Background(), nil)
	_require.NotNil(err)
}

func validateBlobDeleted(_require *require.Assertions, bbClient *blockblob.Client) {
	_, err := bbClient.GetProperties(context.Background(), nil)
	_require.NotNil(err)
	_require.Contains(err.Error(), bloberror.BlobNotFound)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(context.Background(), nil)

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	etag := resp.ETag

	_, err = bbClient.SetMetadata(context.Background(), nil, nil)
	_require.Nil(err)

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: etag},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(context.Background(), nil)
	etag := resp.ETag
	_, err = bbClient.SetMetadata(context.Background(), nil, nil)
	_require.Nil(err)

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: etag},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	_require.Nil(err)

	validateBlobDeleted(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobDeleteIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, _ := bbClient.GetProperties(context.Background(), nil)
	etag := resp.ETag

	deleteBlobOptions := blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: etag},
		},
	}
	_, err = bbClient.Delete(context.Background(), &deleteBlobOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.NotNil(err)
	testcommon.ValidateHTTPErrorCode(_require, err, 304)
}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	resp, err := bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

//func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceFalse() {
//	// TODO: Not Working
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	blockBlobName := testcommon.GenerateBlobName(testName)
//	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)
//
//	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
//
//	_require.Nil(err)
//	_require.Equal(cResp.RawResponse.StatusCode, 201)
//
//	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date,-10)
//
//	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
//	_require.Nil(err)
//
//	getBlobPropertiesOptions := blob.GetPropertiesOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err = bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
//	_require.NotNil(err)
//}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	resp2, err := bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp2.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsOnMissingBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := containerClient.NewBlobClient("MISSING")

	_, err = bbClient.GetProperties(context.Background(), nil)
	testcommon.ValidateHTTPErrorCode(_require, err, 404)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	eTag := azcore.ETag("garbage")
	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &eTag},
		},
	}
	_, err = bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.NotNil(err)
	testcommon.ValidateHTTPErrorCode(_require, err, 412)
}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Nil(err)

	eTag := azcore.ETag("garbage")
	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: &eTag},
		},
	}
	resp, err := bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobGetPropsAndMetadataIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.SetMetadata(context.Background(), nil, nil)
	_require.Nil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}
	_, err = bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.NotNil(err)
	testcommon.ValidateHTTPErrorCode(_require, err, 304)
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesBasic() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, nil)
	_require.Nil(err)

	resp, _ := bbClient.GetProperties(context.Background(), nil)
	h := blob.ParseHTTPHeaders(resp)
	_require.EqualValues(h, testcommon.BasicHeaders)
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesEmptyValue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	contentType := to.Ptr("my_type")
	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentType: contentType}, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.ContentType, contentType)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{}, nil)
	_require.Nil(err)

	resp, err = bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.ContentType)
}

func validatePropertiesSet(_require *require.Assertions, bbClient *blockblob.Client, disposition string) {
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentDisposition, disposition)
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
			},
		})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
			}})
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		}})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		}})
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		}})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: to.Ptr(azcore.ETag("garbage"))},
		}})
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: to.Ptr(azcore.ETag("garbage"))},
		}})
	_require.Nil(err)

	validatePropertiesSet(_require, bbClient, "my_disposition")
}

func (s *BlobRecordedTestsSuite) TestBlobSetPropertiesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	_, err = bbClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		}})
	_require.NotNil(err)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(context.Background(), map[string]*string{"not": to.Ptr("nil")}, nil)
	_require.Nil(err)

	_, err = bbClient.SetMetadata(context.Background(), nil, nil)
	_require.Nil(err)

	blobGetResp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Len(blobGetResp.Metadata, 0)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(context.Background(), map[string]*string{"not": to.Ptr("nil")}, nil)
	_require.Nil(err)

	_, err = bbClient.SetMetadata(context.Background(), map[string]*string{}, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.Metadata, 0)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.SetMetadata(context.Background(), map[string]*string{"Invalid field!": to.Ptr("value")}, nil)
	_require.NotNil(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)
	//_require.Equal(strings.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring), true)
}

func validateMetadataSet(_require *require.Assertions, bbClient *blockblob.Client) {
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(bbName, containerClient)

	cResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: to.Ptr(currentTime)},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: to.Ptr(azcore.ETag("garbage"))},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: to.Ptr(azcore.ETag("garbage"))},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)

	validateMetadataSet(_require, bbClient)
}

func (s *BlobRecordedTestsSuite) TestBlobSetMetadataIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlobRecordedTestsSuite) TestPermanentDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.Nil(err)

	// Create container and blob, upload blob to container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)

	parts, err := sas.ParseURL(bbClient.URL()) // Get parts for BlobURL
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	// Set Account SAS and set Permanent Delete to true
	parts.SAS, err = sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true, PermanentDelete: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	// Create snapshot of Blob and get snapshot URL
	resp, err := bbClient.CreateSnapshot(context.Background(), &blob.CreateSnapshotOptions{})
	_require.Nil(err)
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)

	// Check that there are two items in the container: one snapshot, one blob
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{Include: container.ListBlobsInclude{Snapshots: true}})
	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	// Delete snapshot (snapshot will be soft deleted)
	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = bbClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.Nil(err)

	// Check that only blob exists (snapshot is soft-deleted)
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 1)

	// Check that soft-deleted snapshot exists by including deleted items
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Deleted: true},
	})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	// Options for PermanentDeleteOptions
	perm := blob.DeleteTypePermanent
	deleteBlobOptions := blob.DeleteOptions{
		BlobDeleteType: &perm,
	}
	// Execute Delete with DeleteTypePermanent
	pdResp, err := snapshotURL.Delete(context.Background(), &deleteBlobOptions)
	_require.Nil(err)
	_require.NotNil(pdResp)

	// Check that only blob exists even after including snapshots and deleted items
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Deleted: true}})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 1)
}

func (s *BlobRecordedTestsSuite) TestPermanentDeleteWithoutPermission() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)

	// Create container and blob, upload blob to container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)

	parts, err := sas.ParseURL(bbClient.URL()) // Get parts for BlobURL
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	// Set Account SAS
	parts.SAS, err = sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	// Create snapshot of Blob and get snapshot URL
	resp, err := bbClient.CreateSnapshot(context.Background(), &blob.CreateSnapshotOptions{})
	_require.Nil(err)
	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)

	// Check that there are two items in the container: one snapshot, one blob
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{Include: container.ListBlobsInclude{Snapshots: true}})
	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	// Delete snapshot
	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = bbClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.Nil(err)

	// Check that only blob exists
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 1)

	// Options for PermanentDeleteOptions
	perm := blob.DeleteTypePermanent
	deleteBlobOptions := blob.DeleteOptions{
		BlobDeleteType: &perm,
	}
	// Execute Delete with DeleteTypePermanent,should fail because permissions are not set and snapshot is not soft-deleted
	_, err = snapshotURL.Delete(context.Background(), &deleteBlobOptions)
	_require.NotNil(err)
}

/*func testBlobServiceClientDeleteImpl(_ *require.Assertions, _ *service.Client) error {
	//containerClient := testcommon.CreateNewContainer(context.Background(), _require, "gocblobserviceclientdeleteimpl", svcClient)
	//defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	//bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "goblobserviceclientdeleteimpl", containerClient)
	//
	//_, err := bbClient.Delete(context.Background(), nil)
	//_require.Nil(err) // This call will not have errors related to slow update of service properties, so we assert.
	//
	//_, err = bbClient.Undelete(ctx)
	//if err != nil { // We want to give the wrapper method a chance to check if it was an error related to the service properties update.
	//	return err
	//}
	//
	//resp, err := bbClient.GetProperties(context.Background(), nil)
	//if err != nil {
	//	return errors.New(string(err.(*StorageError).ErrorCode))
	//}
	//_require.Equal(resp.BlobType, BlobTypeBlockBlob) // We could check any property. This is just to double check it was undeleted.
	return nil
}*/

//
////func (s *BlobRecordedTestsSuite) TestBlobTierInferred() {
////	svcClient, err := getPremiumserviceClient()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := testcommon.CreateNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, _ := createNewPageBlob(c, containerClient)
////
////	resp, err := bbClient.GetProperties(context.Background(), nil)
////	_require.Nil(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "true")
////
////	resp2, err := containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.NotNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTier, chk.Not(chk.Equals), "")
////
////	_, err = bbClient.SetTier(ctx, AccessTierP4, LeaseAccessConditions{})
////	_require.Nil(err)
////
////	resp, err = bbClient.GetProperties(context.Background(), nil)
////	_require.Nil(err)
////	_assert(resp.AccessTierInferred(), chk.Equals, "")
////
////	resp2, err = containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.IsNil) // AccessTierInferred never returned if false
////}
////
////func (s *BlobRecordedTestsSuite) TestBlobArchiveStatus() {
////	svcClient, err := getBlobStorageserviceClient()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := testcommon.CreateNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_require.Nil(err)
////	_, err = bbClient.SetTier(ctx, AccessTierCool, LeaseAccessConditions{})
////	_require.Nil(err)
////
////	resp, err := bbClient.GetProperties(context.Background(), nil)
////	_require.Nil(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToCool))
////
////	resp2, err := containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToCool)
////
////	// delete first blob
////	_, err = bbClient.Delete(context.Background(), DeleteSnapshotsOptionNone, nil)
////	_require.Nil(err)
////
////	bbClient, _ = createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_require.Nil(err)
////	_, err = bbClient.SetTier(ctx, AccessTierHot, LeaseAccessConditions{})
////	_require.Nil(err)
////
////	resp, err = bbClient.GetProperties(context.Background(), nil)
////	_require.Nil(err)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToHot))
////
////	resp2, err = containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToHot)
////}
////
////func (s *BlobRecordedTestsSuite) TestBlobTierInvalidValue() {
////	svcClient, err := getBlobStorageserviceClient()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := testcommon.CreateNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = bbClient.SetTier(ctx, AccessTierType("garbage"), LeaseAccessConditions{})
////	testcommon.ValidateBlobErrorCode(c, err, bloberror.InvalidHeaderValue)
////}
////

func (s *BlobRecordedTestsSuite) TestBlobClientPartsSASQueryTimes() {
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

		parts, _ := blob.ParseURL(urlString)
		_require.Equal(parts.Scheme, "https")
		_require.Equal(parts.Host, "myaccount.blob.core.windows.net")
		_require.Equal(parts.ContainerName, "mycontainer")
		_require.Equal(parts.BlobName, "mydirectory/myfile.txt")

		sas := parts.SAS
		_require.Equal(sas.StartTime(), StartTimesExpected[i])
		_require.Equal(sas.ExpiryTime(), ExpiryTimesExpected[i])

		_require.Equal(parts.String(), urlString)
	}
}

//func (s *BlobUnrecordedTestsSuite) TestDownloadBlockBlobUnexpectedEOF() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(),testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	blockBlobName := testcommon.GenerateBlobName(testName)
//	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)
//
//	resp, err := bbClient.DownloadStream(context.Background(), nil)
//	_require.Nil(err)
//
//	// Verify that we can inject errors first.
//	reader := resp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))
//
//	_, err = io.ReadAll(reader)
//	_require.NotNil(err)
//	_require.Equal(err.Error(), "unrecoverable error")
//
//	// Then inject the retryable error.
//	reader = resp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))
//
//	buf, err := io.ReadAll(reader)
//	_require.Nil(err)
//	_require.EqualValues(buf, []byte(testcommon.BlockBlobDefaultData))
//}

//func InjectErrorInRetryReaderOptions(err error) *blob.RetryReaderOptions {
//	return &blob.RetryReaderOptions{
//		MaxRetryRequests:       1,
//		doInjectError:          true,
//		doInjectErrorRound:     0,
//		injectedError:          err,
//		NotifyFailedRead:       nil,
//		TreatEarlyCloseAsError: false,
//		CPKInfo:                nil,
//		CPKScopeInfo:           nil,
//	}
//}

func (s *BlobRecordedTestsSuite) TestBlobSetExpiry() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.ExpiresOn)

	_, err = bbClient.SetExpiry(context.Background(), blockblob.ExpiryTypeRelativeToNow(8*time.Second), nil)
	_require.Nil(err)

	resp, err = bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.ExpiresOn)

	time.Sleep(time.Second * 10)

	_, err = bbClient.GetProperties(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *BlobRecordedTestsSuite) TestSetImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.Nil(err)
	policy := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.Nil(err)

	setImmutabilityPolicyOptions := &blob.SetImmutabilityPolicyOptions{
		Mode:                     &policy,
		ModifiedAccessConditions: nil,
	}
	_, err = bbClient.SetImmutabilityPolicy(context.Background(), currentTime, setImmutabilityPolicyOptions)
	_require.Nil(err)

	_, err = bbClient.SetLegalHold(context.Background(), false, nil)
	_require.Nil(err)

	_, err = bbClient.Delete(context.Background(), nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobImmutableDueToPolicy)

	_, err = bbClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.Nil(err)

	_, err = bbClient.Delete(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestDeleteImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.Nil(err)

	policy := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.Nil(err)

	setImmutabilityPolicyOptions := &blob.SetImmutabilityPolicyOptions{
		Mode:                     &policy,
		ModifiedAccessConditions: nil,
	}
	_, err = bbClient.SetImmutabilityPolicy(context.Background(), currentTime, setImmutabilityPolicyOptions)
	_require.Nil(err)

	_, err = bbClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.Nil(err)

	_, err = bbClient.Delete(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlobRecordedTestsSuite) TestSetLegalHold() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	_, err = bbClient.SetLegalHold(context.Background(), true, nil)
	_require.Nil(err)

	// should fail since time has not passed yet
	_, err = bbClient.Delete(context.Background(), nil)
	_require.NotNil(err)

	_, err = bbClient.SetLegalHold(context.Background(), false, nil)
	_require.Nil(err)

	_, err = bbClient.Delete(context.Background(), nil)
	_require.Nil(err)

}

func (s *BlobUnrecordedTestsSuite) TestSASURLBlobClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	_require.Nil(err)

	// Creating service client with credentials
	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.Nil(err)

	// Creating container client
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, serviceClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Creating blob client with credentials
	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)
	blobClient, err := blob.NewClientWithSharedKeyCredential(bbClient.URL(), cred, nil)
	_require.NoError(err)

	// Adding SAS and options
	permissions := sas.BlobPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(5 * time.Minute)

	// BlobSASURL is created with GetSASURL
	sasUrl, err := blobClient.GetSASURL(permissions, expiry, nil)
	_require.Nil(err)

	// Get new blob client with sasUrl and attempt GetProperties
	_, err = blob.NewClientWithNoCredential(sasUrl, nil)
	_require.Nil(err)
}

func (s *BlobUnrecordedTestsSuite) TestNoSharedKeyCredError() {
	_require := require.New(s.T())
	testName := s.T().Name()

	// Creating service client
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Creating container client
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Creating blob client without credentials
	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	// Adding SAS and options
	permissions := sas.BlobPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	start := time.Now().Add(-time.Hour)
	expiry := start.Add(time.Hour)
	opts := blob.GetSASURLOptions{StartTime: &start}

	// GetSASURL fails (with MissingSharedKeyCredential) because blob client is created without credentials
	_, err = bbClient.BlobClient().GetSASURL(permissions, expiry, &opts)
	_require.Equal(err, bloberror.MissingSharedKeyCredential)
}

func (s *BlobRecordedTestsSuite) TestBlobGetAccountInfo() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	bAccInfo, err := bbClient.GetAccountInfo(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(bAccInfo)
}
