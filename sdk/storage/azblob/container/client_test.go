//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"net/url"
	"sort"
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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running container Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &ContainerRecordedTestsSuite{})
		suite.Run(t, &ContainerUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &ContainerRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &ContainerRecordedTestsSuite{})
	}
}

func (s *ContainerRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *ContainerRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *ContainerUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *ContainerUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type ContainerRecordedTestsSuite struct {
	suite.Suite
}

type ContainerUnrecordedTestsSuite struct {
	suite.Suite
}

//func (s *ContainerUnrecordedTestsSuite) TestNewContainerClientValidName() {
//	_require := require.New(s.T())
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	testURL := svcClient.NewContainerClient(containerPrefix)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_require.Nil(err)
//	correctURL := "https://" + accountName + "." + DefaultBlobEndpointSuffix + containerPrefix
//	_require.Equal(testURL.URL(), correctURL)
//}

//func (s *ContainerUnrecordedTestsSuite) TestCreateRootContainerURL() {
//	_require := require.New(s.T())
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	testURL := svcClient.NewContainerClient(ContainerNameRoot)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_require.Nil(err)
//	correctURL := "https://" + accountName + ".blob.core.windows.net/$root"
//	_require.Equal(testURL.URL(), correctURL)
//}

func (s *ContainerRecordedTestsSuite) TestContainerCreateInvalidName() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := svcClient.NewContainerClient("foo bar")

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access:   &access,
		Metadata: map[string]*string{},
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidResourceName)
}

func (s *ContainerRecordedTestsSuite) TestContainerCreateEmptyName() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := svcClient.NewContainerClient("")

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access:   &access,
		Metadata: map[string]*string{},
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidQueryParameterValue)
}

func (s *ContainerRecordedTestsSuite) TestContainerCreateNameCollision() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access:   &access,
		Metadata: map[string]*string{},
	}

	containerClient = svcClient.NewContainerClient(containerName)
	_, err = containerClient.Create(context.Background(), &createContainerOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ContainerAlreadyExists)
}

func (s *ContainerRecordedTestsSuite) TestContainerCreateInvalidMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access:   &access,
		Metadata: map[string]*string{"1 foo": to.Ptr("bar")},
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions)

	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring), true)
}

func (s *ContainerRecordedTestsSuite) TestContainerCreateNilMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access:   &access,
		Metadata: map[string]*string{},
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	_require.Nil(err)

	response, err := containerClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(response.Metadata)
}

func (s *ContainerRecordedTestsSuite) TestContainerCreateEmptyMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access:   &access,
		Metadata: map[string]*string{},
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	_require.Nil(err)

	response, err := containerClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(response.Metadata)
}

//func (s *ContainerRecordedTestsSuite) TestContainerCreateAccessContainer() {
//	// TOD0: NotWorking
//	_require := require.New(s.T())
//testName := s.T().Name()
////
//	svcClient := testcommon.GetServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	credential, err := getGenericCredential("")
//	_require.Nil(err)
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.GetContainerClient(containerName, svcClient)
//
//	access := container.PublicAccessTypeBlob
//	createContainerOptions := container.CreateOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(context.Background(), &createContainerOptions)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	_require.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(testcommon.BlobPrefix)
//	uploadBlockBlobOptions := blockblob.UploadOptions{
//		Metadata: testcommon.BasicMetadata,
//	}
//	_, err = bbClient.Upload(context.Background(), bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	_require.Nil(err)
//
//	// Anonymous enumeration should be valid with container access
//	containerClient2, _ := NewContainerClient(containerClient.URL(), credential, nil)
//	pager := containerClient2.NewListBlobsFlatPager(nil)
//
//	for pager.NextPage(context.Background()) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			_require.Equal(*blob.Name, testcommon.BlobPrefix)
//		}
//	}
//
//	_require.Nil(pager.Err())
//
//	// Getting blob data anonymously should still be valid with container access
//	blobURL2 := containerClient2.NewBlockBlobClient(testcommon.BlobPrefix)
//	resp, err := blobURL2.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
//}

//func (s *ContainerRecordedTestsSuite) TestContainerCreateAccessBlob() {
//	// TODO: Not Working
//	_require := require.New(s.T())
// testName := s.T().Name()
////	svcClient := testcommon.GetServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.GetContainerClient(containerName, svcClient)
//
//	access := container.PublicAccessTypeBlob
//	createContainerOptions := container.CreateOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(context.Background(), &createContainerOptions)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	_require.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(testcommon.BlobPrefix)
//	uploadBlockBlobOptions := blockblob.UploadOptions{
//		Metadata: testcommon.BasicMetadata,
//	}
//	_, err = bbClient.Upload(context.Background(), bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	_require.Nil(err)
//
//	// Reference the same container URL but with anonymous credentials
//	containerClient2, err := NewContainerClient(containerClient.URL(), azcore.AnonymousCredential(), nil)
//	_require.Nil(err)
//
//	pager := containerClient2.NewListBlobsFlatPager(nil)
//
//	_require.Equal(pager.NextPage(context.Background()), false)
//	_require.NotNil(pager.Err())
//
//	// Accessing blob specific data should be public
//	blobURL2 := containerClient2.NewBlockBlobClient(testcommon.BlobPrefix)
//	resp, err := blobURL2.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
//}

func (s *ContainerRecordedTestsSuite) TestContainerCreateAccessNone() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	// Public Access Type None
	_, err = containerClient.Create(context.Background(), nil)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	_require.Nil(err)

	bbClient := containerClient.NewBlockBlobClient(testcommon.BlobPrefix)
	uploadBlockBlobOptions := blockblob.UploadOptions{
		Metadata: testcommon.BasicMetadata,
	}
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader("Content")), &uploadBlockBlobOptions)
	_require.Nil(err)

	// Reference the same container URL but with anonymous credentials
	containerClient2, err := container.NewClientWithNoCredential(containerClient.URL(), nil)
	_require.Nil(err)

	pager := containerClient2.NewListBlobsFlatPager(nil)

	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		if err != nil {
			break
		}
	}

	// Blob data is not public
	// TODO: Fix Inheritance
	//blobURL2 := containerClient2.NewBlockBlobClient(testcommon.BlobPrefix)
	//_, err = blobURL2.GetProperties(context.Background(), nil)
	//_require.NotNil(err)

	//serr := err.(StorageError)
	//_assert(serr.Response().StatusCode, chk.Equals, 401) // HEAD request does not return a status code
}

//func (s *ContainerRecordedTestsSuite) TestContainerCreateIfExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	serviceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.GetContainerClient(containerName, serviceClient)
//
//	// Public Access Type None
//	_, err = containerClient.Create(context.Background(), nil)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	_require.Nil(err)
//
//	access := container.PublicAccessTypeBlob
//	createContainerOptions := container.CreateOptions{
//		Access:   &access,
//		Metadata: nil,
//	}
//	_, err = containerClient.CreateIfNotExists(context.Background(), &createContainerOptions)
//	_require.Nil(err)
//
//	// Ensure that next create call doesn't update the properties of already created container
//	getResp, err := containerClient.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//	_require.Nil(getResp.BlobPublicAccess)
//	_require.Nil(getResp.Metadata)
//}
//
//func (s *ContainerRecordedTestsSuite) TestContainerCreateIfNotExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	serviceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.GetContainerClient(containerName, serviceClient)
//
//	access := container.PublicAccessTypeBlob
//	createContainerOptions := container.CreateOptions{
//		Access:   &access,
//		Metadata: testcommon.BasicMetadata,
//	}
//	_, err = containerClient.CreateIfNotExists(context.Background(), &createContainerOptions)
//	_require.Nil(err)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	// Ensure that next create call doesn't update the properties of already created container
//	getResp, err := containerClient.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//	_require.EqualValues(*getResp.BlobPublicAccess, PublicAccessTypeBlob)
//	_require.EqualValues(getResp.Metadata, testcommon.BasicMetadata)
//}

func validateContainerDeleted(_require *require.Assertions, containerClient *container.Client) {
	_, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ContainerNotFound)
}

func (s *ContainerRecordedTestsSuite) TestContainerDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	_, err = containerClient.Delete(context.Background(), nil)
	_require.Nil(err)

	validateContainerDeleted(_require, containerClient)
}

//func (s *ContainerRecordedTestsSuite) TestContainerDeleteIfExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	serviceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.GetContainerClient(containerName, serviceClient)
//
//	// Public Access Type None
//	_, err = containerClient.Create(context.Background(), nil)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	_require.Nil(err)
//
//	_, err = containerClient.DeleteIfExists(context.Background(), nil)
//	_require.Nil(err)
//
//	validateContainerDeleted(_require, containerClient)
//}
//
//func (s *ContainerRecordedTestsSuite) TestContainerDeleteIfNotExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	serviceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.GetContainerClient(containerName, serviceClient)
//
//	_, err = containerClient.DeleteIfExists(context.Background(), nil)
//	_require.Nil(err)
//}

func (s *ContainerRecordedTestsSuite) TestContainerDeleteNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	_, err = containerClient.Delete(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ContainerNotFound)
}

func (s *ContainerRecordedTestsSuite) TestContainerDeleteIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	deleteContainerOptions := container.DeleteOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = containerClient.Delete(context.Background(), &deleteContainerOptions)
	_require.Nil(err)
	validateContainerDeleted(_require, containerClient)
}

func (s *ContainerRecordedTestsSuite) TestContainerDeleteIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	deleteContainerOptions := container.DeleteOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = containerClient.Delete(context.Background(), &deleteContainerOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *ContainerRecordedTestsSuite) TestContainerDeleteIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	deleteContainerOptions := container.DeleteOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = containerClient.Delete(context.Background(), &deleteContainerOptions)
	_require.Nil(err)

	validateContainerDeleted(_require, containerClient)
}

func (s *ContainerRecordedTestsSuite) TestContainerDeleteIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	deleteContainerOptions := container.DeleteOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = containerClient.Delete(context.Background(), &deleteContainerOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

////func (s *ContainerRecordedTestsSuite) TestContainerAccessConditionsUnsupportedConditions() {
////	// This test defines that the library will panic if the user specifies conditional headers
////	// that will be ignored by the service
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////
////	invalidEtag := "invalid"
////	deleteContainerOptions := ContainerSetMetadataOptions{
////		Metadata: testcommon.BasicMetadata,
////		ModifiedAccessConditions: &ModifiedAccessConditions{
////			IfMatch: &invalidEtag,
////		},
////	}
////	_, err := containerClient.SetMetadata(context.Background(), &deleteContainerOptions)
////	_require.NotNil(err)
////}
//
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsNonexistentPrefix() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	prefix := testcommon.BlobPrefix + testcommon.BlobPrefix
////	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
////		Prefix: &prefix,
////	}
////	listResponse, errChan := containerClient.NewListBlobsFlatPager(context.Background(), 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(listResponse, chk.IsNil)
////}
//
//func (s *ContainerRecordedTestsSuite) TestContainerListBlobsSpecificValidPrefix() {
//	svcClient := testcommon.GetServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	prefix := testcommon.BlobPrefix
//	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
//		Prefix: &prefix,
//	}
//	pager := containerClient.NewListBlobsFlatPager(&containerListBlobFlatSegmentOptions)
//
//	count := 0
//
//	for pager.NextPage(context.Background()) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			count++
//			_assert(*blob.Name, chk.Equals, blobName)
//		}
//	}
//
//	_assert(pager.Err(), chk.IsNil)
//
//	_assert(count, chk.Equals, 1)
//}
//
//func (s *ContainerRecordedTestsSuite) TestContainerListBlobsValidDelimiter() {
//	svcClient := testcommon.GetServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	prefixes := []string{"a/1", "a/2", "b/2", "blob"}
//	blobNames := make([]string, 4)
//	for idx, prefix := range prefixes {
//		_, blobNames[idx] = createNewBlockBlobWithPrefix(c, containerClient, prefix)
//	}
//
//	pager := containerClient.NewListBlobsHierarchyPager("/", nil)
//
//	count := 0
//
//	for pager.NextPage(context.Background()) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			count++
//			_assert(*blob.Name, chk.Equals, blobNames[3])
//		}
//	}
//
//	_assert(pager.Err(), chk.IsNil)
//	_assert(count, chk.Equals, 1)
//
//	// TODO: Ask why the output is BlobItemInternal and why other fields are not there for ex: prefix array
//	//_require.Nil(err)
//	//_assert(len(resp.Segment.BlobItems), chk.Equals, 1)
//	//_assert(len(resp.Segment.BlobPrefixes), chk.Equals, 2)
//	//_assert(resp.Segment.BlobPrefixes[0].Name, chk.Equals, "a/")
//	//_assert(resp.Segment.BlobPrefixes[1].Name, chk.Equals, "b/")
//	//_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

//func (s *ContainerRecordedTestsSuite) TestContainerListBlobsWithSnapshots() {
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
//	// initialize a blob and create a snapshot of it
//	snapBlobName := testcommon.GenerateBlobName(testName)
//	snapBlob := testcommon.CreateNewBlockBlob(context.Background(), _require, snapBlobName, containerClient)
//	snap, err := snapBlob.CreateSnapshot(context.Background(), nil)
//	// snap.
//	_require.Nil(err)
//
//	listBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
//		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots},
//	}
//	pager := containerClient.NewListBlobsFlatPager(&listBlobFlatSegmentOptions)
//
//	wasFound := false // hold the for loop accountable for finding the blob and it's snapshot
//	for pager.More() {
//		resp, err := pager.NextPage(context.Background())
//		_require.Nil(err)
//
//		for _, blob := range resp.Segment.BlobItems {
//			if *blob.Name == snapBlobName && blob.Snapshot != nil {
//				wasFound = true
//				_require.Equal(*blob.Snapshot, *snap.Snapshot)
//			}
//		}
//		if err != nil {
//			break
//		}
//	}
//	_require.Equal(wasFound, true)
//}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsInvalidDelimiter() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	prefixes := []string{"a/1", "a/2", "b/1", "blob"}
	for _, prefix := range prefixes {
		blobName := prefix + testcommon.GenerateBlobName(testName)
		testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	}

	pager := containerClient.NewListBlobsHierarchyPager("^", nil)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Nil(resp.Segment.BlobPrefixes)
		if err != nil {
			break
		}
	}
}

////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeMetadata() {
////	svcClient := testcommon.GetServiceClient()
////	container, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(container)
////	_, blobNameNoMetadata := createNewBlockBlobWithPrefix(c, container, "a")
////	blobMetadata, blobNameMetadata := createNewBlockBlobWithPrefix(c, container, "b")
////	_, err := blobMetadata.SetMetadata(context.Background(), Metadata{"field": "value"}, LeaseAccessConditions{}, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////
////	resp, err := container.NewListBlobsFlatPager(context.Background(), Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Metadata: true}})
////
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobNameNoMetadata)
////	_assert(resp.Segment.BlobItems[0].Metadata, chk.HasLen, 0)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobNameMetadata)
////	_assert(resp.Segment.BlobItems[1].Metadata["field"], chk.Equals, "value")
////}
//
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeSnapshots() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	blob, blobName := createNewBlockBlob(c, containerClient)
////	_, err := blob.CreateSnapshot(context.Background(), Metadata{}, LeaseAccessConditions{}, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////
////	resp, err := containerClient.NewListBlobsFlatPager(context.Background(), Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
////
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 2)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[0].Snapshot, chk.NotNil)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[1].Snapshot, chk.Equals, "")
////}
////
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeCopy() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, blobName := createNewBlockBlob(c, containerClient)
////	blobCopyURL, blobCopyName := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	_, err := blobCopyURL.StartCopyFromURL(context.Background(), bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, LeaseAccessConditions{}, DefaultAccessTier, nil)
////	_require.Nil(err)
////
////	resp, err := containerClient.NewListBlobsFlatPager(context.Background(), Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Copy: true}})
////
////	// These are sufficient to show that the blob copy was in fact included
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 2)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobCopyName)
////	_assert(*resp.Segment.BlobItems[0].Properties.ContentLength, chk.Equals, int64(len(testcommon.BlockBlobDefaultData)))
////	temp := bbClient.URL()
////	_assert(*resp.Segment.BlobItems[0].Properties.CopySource, chk.Equals, temp.String())
////	_assert(resp.Segment.BlobItems[0].Properties.CopyStatus, chk.Equals, CopyStatusTypeSuccess)
////}
////
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeUncommitted() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, blobName := getBlockBlobURL(c, containerClient)
////	_, err := bbClient.StageBlock(context.Background(), blockID, strings.NewReader(testcommon.BlockBlobDefaultData), LeaseAccessConditions{}, nil, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////
////	resp, err := containerClient.NewListBlobsFlatPager(context.Background(), Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{UncommittedBlobs: true}})
////
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////}
//
////func testContainerListBlobsIncludeTypeDeletedImpl(, svcClient ServiceURL) error {
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	resp, err := containerClient.NewListBlobsFlatPager(context.Background(), Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////
////	_, err = bbClient.Delete(context.Background(), DeleteSnapshotsOptionInclude, LeaseAccessConditions{})
////	_require.Nil(err)
////
////	resp, err = containerClient.NewListBlobsFlatPager(context.Background(), Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_require.Nil(err)
////	if len(resp.Segment.BlobItems) != 1 {
////		return errors.New("DeletedBlobNotFound")
////	}
////
////	// resp.Segment.BlobItems[0].Deleted == true/false if versioning is disabled/enabled.
////	_assert(resp.Segment.BlobItems[0].Deleted, chk.Equals, false)
////	return nil
////}
////
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeDeleted() {
////	svcClient := testcommon.GetServiceClient()
////
////	runTestRequiringServiceProperties(c, svcClient, "DeletedBlobNotFound", enableSoftDelete,
////		testContainerListBlobsIncludeTypeDeletedImpl, disableSoftDelete)
////}
////
////func testContainerListBlobsIncludeMultipleImpl(, svcClient ServiceURL) error {
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////
////	bbClient, _ := createNewBlockBlobWithPrefix(c, containerClient, "z")
////	_, err := bbClient.CreateSnapshot(context.Background(), Metadata{}, LeaseAccessConditions{}, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////	blobURL2, _ := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	resp2, err := blobURL2.StartCopyFromURL(context.Background(), bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, LeaseAccessConditions{}, DefaultAccessTier, nil)
////	_require.Nil(err)
////	waitForCopy(c, blobURL2, resp2)
////	blobURL3, _ := createNewBlockBlobWithPrefix(c, containerClient, "deleted")
////
////	_, err = blobURL3.Delete(context.Background(), DeleteSnapshotsOptionNone, LeaseAccessConditions{})
////
////	resp, err := containerClient.NewListBlobsFlatPager(context.Background(), Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Copy: true, Deleted: true, Versions: true}})
////
////	_require.Nil(err)
////	if len(resp.Segment.BlobItems) != 6 {
////		// If there are fewer blobs in the container than there should be, it will be because one was permanently deleted.
////		return errors.New("DeletedBlobNotFound")
////	}
////
////	//_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName2)
////	//_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName) // With soft delete, the overwritten blob will have a backup snapshot
////	//_assert(resp.Segment.BlobItems[2].Name, chk.Equals, blobName)
////	return nil
////}
////
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeMultiple() {
////	svcClient := testcommon.GetServiceClient()
////
////	runTestRequiringServiceProperties(c, svcClient, "DeletedBlobNotFound", enableSoftDelete,
////		testContainerListBlobsIncludeMultipleImpl, disableSoftDelete)
////}
////
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsMaxResultsNegative() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	_, err := containerClient.NewListBlobsFlatPager(context.Background(), Marker{}, ListBlobsSegmentOptions{MaxResults: -2})
////	_assert(err, chk.Not(chk.IsNil))
////}
//
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsMaxResultsZero() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	maxResults := int32(0)
////	resp, errChan := containerClient.NewListBlobsFlatPager(context.Background(), 1, 0, &ContainerListBlobsFlatOptions{MaxResults: &maxResults})
////
////	_assert(<-errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////}
//
//// TODO: Adele: Case failing
////func (s *ContainerRecordedTestsSuite) TestContainerListBlobsMaxResultsInsufficient() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
////	createNewBlockBlobWithPrefix(c, containerClient, "b")
////
////	maxResults := int32(1)
////	resp, errChan := containerClient.NewListBlobsFlatPager(context.Background(), 3, 0, &ContainerListBlobsFlatOptions{MaxResults: &maxResults})
////	_assert(<- errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////	_assert((<- resp).Name, chk.Equals, blobName)
////}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsMaxResultsExact() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	blobNames := make([]string, 2)
	blobName := testcommon.GenerateBlobName(testName)
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	testcommon.CreateNewBlockBlob(context.Background(), _require, blobNames[0], containerClient)
	testcommon.CreateNewBlockBlob(context.Background(), _require, blobNames[1], containerClient)

	maxResult := int32(2)
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		MaxResults: &maxResult,
	})

	nameMap := testcommon.BlobListToMap(blobNames)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		for _, blob := range resp.Segment.BlobItems {
			_require.Equal(nameMap[*blob.Name], true)
		}
		if err != nil {
			break
		}
	}
}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsMaxResultsSufficient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobNames := make([]string, 2)
	blobName := testcommon.GenerateBlobName(testName)
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	testcommon.CreateNewBlockBlob(context.Background(), _require, blobNames[0], containerClient)
	testcommon.CreateNewBlockBlob(context.Background(), _require, blobNames[1], containerClient)

	maxResult := int32(3)
	containerListBlobFlatSegmentOptions := container.ListBlobsFlatOptions{
		MaxResults: &maxResult,
	}
	pager := containerClient.NewListBlobsFlatPager(&containerListBlobFlatSegmentOptions)

	nameMap := testcommon.BlobListToMap(blobNames)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		for _, blob := range resp.Segment.BlobItems {
			_require.Equal(nameMap[*blob.Name], true)
		}
		if err != nil {
			break
		}
	}
}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsNonExistentContainer() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	pager := containerClient.NewListBlobsFlatPager(nil)
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		if err != nil {
			break
		}
	}
}

func (s *ContainerRecordedTestsSuite) TestContainerGetPropertiesAndMetadataNoMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := containerClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *ContainerRecordedTestsSuite) TestContainerGetPropsAndMetaNonExistentContainer() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	_, err = containerClient.GetProperties(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ContainerNotFound)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Metadata: testcommon.BasicMetadata,
		Access:   &access,
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	_require.Nil(err)

	setMetadataContainerOptions := container.SetMetadataOptions{
		Metadata: map[string]*string{},
	}
	_, err = containerClient.SetMetadata(context.Background(), &setMetadataContainerOptions)
	_require.Nil(err)

	resp, err := containerClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)
	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access:   &access,
		Metadata: testcommon.BasicMetadata,
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions)
	_require.Nil(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err = containerClient.SetMetadata(context.Background(), nil)
	_require.Nil(err)

	resp, err := containerClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	setMetadataContainerOptions := container.SetMetadataOptions{
		Metadata: map[string]*string{"!nval!d Field!@#%": to.Ptr("value")},
	}
	_, err = containerClient.SetMetadata(context.Background(), &setMetadataContainerOptions)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring), true)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetMetadataNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	_, err = containerClient.SetMetadata(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ContainerNotFound)
}

//
//func (s *ContainerRecordedTestsSuite) TestContainerSetMetadataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	svcClient := testcommon.GetServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	setMetadataContainerOptions := ContainerSetMetadataOptions{
//		Metadata: testcommon.BasicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := containerClient.SetMetadata(context.Background(), &setMetadataContainerOptions)
//	_require.Nil(err)
//
//	resp, err := containerClient.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//	_assert(resp.Metadata, chk.NotNil)
//	_assert(resp.Metadata, chk.DeepEquals, testcommon.BasicMetadata)
//
//}

//func (s *ContainerRecordedTestsSuite) TestContainerSetMetadataIfModifiedSinceFalse() {
//	// TODO: NotWorking
//	_require := require.New(s.T())
// testName := s.T().Name()
////	svcClient := testcommon.GetServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	containerClient, _ := testcommon.CreateNewContainer(context.Background(), _require, testName, svcClient)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	//currentTime := getRelativeTimeGMT(10)
//	//currentTime, err := time.Parse(time.UnixDate, "Wed Jan 07 11:11:11 PST 2099")
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//	_require.Nil(err)
//	setMetadataContainerOptions := ContainerSetMetadataOptions{
//		Metadata: testcommon.BasicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err = containerClient.SetMetadata(context.Background(), &setMetadataContainerOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
//}

func (s *ContainerRecordedTestsSuite) TestContainerNewBlobURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	bbClient := containerClient.NewBlobClient(testcommon.BlobPrefix)

	_require.Equal(bbClient.URL(), containerClient.URL()+"/"+testcommon.BlobPrefix)
	_require.IsTypef(bbClient, &blob.Client{}, fmt.Sprintf("%T should be of type %T", bbClient, blob.Client{}))
}

func (s *ContainerRecordedTestsSuite) TestContainerNewBlockBlobClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	bbClient := containerClient.NewBlockBlobClient(testcommon.BlobPrefix)

	_require.Equal(bbClient.URL(), containerClient.URL()+"/"+testcommon.BlobPrefix)
	_require.IsTypef(bbClient, &blockblob.Client{}, fmt.Sprintf("%T should be of type %T", bbClient, blockblob.Client{}))
}

func (s *ContainerRecordedTestsSuite) TestListBlobIncludeMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	for i := 0; i < 6; i++ {
		bbClient := testcommon.GetBlockBlobClient(blobName+strconv.Itoa(i), containerClient)
		_, err = bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &blockblob.UploadOptions{Metadata: testcommon.BasicMetadata})
		_require.Nil(err)
		// _require.Equal(cResp.RawResponse.StatusCode, 201)
	}

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Metadata: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		_require.Len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems, 6)
		for _, blob := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			_require.NotNil(blob.Metadata)
			_require.Len(blob.Metadata, len(testcommon.BasicMetadata))
		}
		if err != nil {
			break
		}
	}

	//----------------------------------------------------------

	pager1 := containerClient.NewListBlobsHierarchyPager("/", &container.ListBlobsHierarchyOptions{
		Include: container.ListBlobsInclude{Metadata: true, Tags: true},
	})

	for pager1.More() {
		resp, err := pager1.NextPage(context.Background())
		_require.Nil(err)
		if err != nil {
			break
		}
		_require.Len(resp.Segment.BlobItems, 6)
		for _, blob := range resp.Segment.BlobItems {
			_require.NotNil(blob.Metadata)
			_require.Len(blob.Metadata, len(testcommon.BasicMetadata))
		}
	}
}

func (s *ContainerRecordedTestsSuite) TestBlobListWrapper() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	_, err = containerClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	files := []string{"a123", "b234", "c345"}

	testcommon.CreateNewBlobs(context.Background(), _require, files, containerClient)

	found := make([]string, 0)

	pager := containerClient.NewListBlobsFlatPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, blob := range resp.Segment.BlobItems {
			found = append(found, *blob.Name)
		}
	}

	sort.Strings(files)
	sort.Strings(found)

	_require.EqualValues(files, found)
}

func (s *ContainerRecordedTestsSuite) TestBlobListWrapperListingError() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.GetContainerClient(testcommon.GenerateContainerName(testName), svcClient)

	pager := containerClient.NewListBlobsFlatPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		var respErr *azcore.ResponseError
		_require.ErrorAs(err, &respErr)
		_require.Equal(bloberror.ContainerNotFound, bloberror.Code(respErr.ErrorCode))
		_require.Empty(resp)
		break
	}
}

func (s *ContainerUnrecordedTestsSuite) TestSetEmptyAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err = containerClient.SetAccessPolicy(context.Background(), &container.SetAccessPolicyOptions{})
	_require.Nil(err)
}

func (s *ContainerUnrecordedTestsSuite) TestSetAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"

	signedIdentifiers := make([]*container.SignedIdentifier, 0)

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		AccessPolicy: &container.AccessPolicy{
			Expiry:     &expiration,
			Start:      &start,
			Permission: &permission,
		},
		ID: &id,
	})
	options := container.SetAccessPolicyOptions{ContainerACL: signedIdentifiers}
	_, err = containerClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)
}

func (s *ContainerUnrecordedTestsSuite) TestSetMultipleAccessPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	id := "empty"

	signedIdentifiers := make([]*container.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id,
	})

	permission2 := "r"
	id2 := "partial"

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id2,
		AccessPolicy: &container.AccessPolicy{
			Permission: &permission2,
		},
	})

	id3 := "full"
	permission3 := "r"
	start := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)
	expiry := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id3,
		AccessPolicy: &container.AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission3,
		},
	})
	options := container.SetAccessPolicyOptions{ContainerACL: signedIdentifiers}
	_, err = containerClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)

	// Make a Get to assert two access policies
	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 3)
}

func (s *ContainerUnrecordedTestsSuite) TestSetNullAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	id := "null"

	signedIdentifiers := make([]*container.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id,
	})
	options := container.SetAccessPolicyOptions{ContainerACL: signedIdentifiers}
	_, err = containerClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(len(resp.SignedIdentifiers), 1)
}

func (s *ContainerRecordedTestsSuite) TestContainerGetSetPermissionsMultiplePolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry := start.Add(5 * time.Minute)
	expiry2 := start.Add(time.Minute)
	readWrite := to.Ptr(container.AccessPolicyPermission{Read: true, Write: true}).String()
	readOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	id1, id2 := "0000", "0001"
	permissions := []*container.SignedIdentifier{
		{ID: &id1,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &readWrite,
			},
		},
		{ID: &id2,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry2,
				Permission: &readOnly,
			},
		},
	}
	options := container.SetAccessPolicyOptions{ContainerACL: permissions}
	_, err = containerClient.SetAccessPolicy(context.Background(), &options)

	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *ContainerRecordedTestsSuite) TestContainerGetPermissionsPublicAccessNotNone() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access: &access,
	}
	_, err = containerClient.Create(context.Background(), &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	_require.Nil(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)

	_require.Nil(err)
	_require.Equal(*resp.BlobPublicAccess, container.PublicAccessTypeBlob)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsPublicAccessNone() {
	// Test the basic one by making an anonymous request to ensure it's actually doing it and also with GetPermissions
	// For all the others, can just use GetPermissions since we've validated that it at least registers on the server correctly
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	_ = testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)

	// Container is created with PublicAccessTypeBlob, so setting it to None will actually test that it is changed through this method
	_, err = containerClient.SetAccessPolicy(context.Background(), nil)
	_require.Nil(err)

	bsu2, err := service.NewClientWithNoCredential(svcClient.URL(), nil)
	_require.Nil(err)

	containerClient2 := bsu2.NewContainerClient(containerName)

	// Get permissions via the original container URL so the request succeeds
	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(resp.BlobPublicAccess)
	_require.Nil(err)

	// If we cannot access a blob's data, we will also not be able to enumerate blobs
	pager := containerClient2.NewListBlobsFlatPager(nil)
	for pager.More() {
		_, err = pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.NoAuthenticationInformation)
		break
	}

	blobClient2 := containerClient2.NewBlockBlobClient(blobName)
	_, err = blobClient2.DownloadStream(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.NoAuthenticationInformation)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsPublicAccessTypeBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: to.Ptr(container.PublicAccessTypeBlob),
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.BlobPublicAccess, container.PublicAccessTypeBlob)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsPublicAccessContainer() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: to.Ptr(container.PublicAccessTypeContainer),
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.BlobPublicAccess, container.PublicAccessTypeContainer)
}

//
//func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsACLSinglePolicy() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	_ = testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
//
//	start := time.Now().UTC().Add(-15 * time.Second)
//	expiry := start.Add(5 * time.Minute).UTC()
//	listOnly := AccessPolicyPermission{List: true}.String()
//	id := "0000"
//	permissions := []*container.SignedIdentifier{{
//		ID: &id,
//		AccessPolicy: &container.AccessPolicy{
//			Start:      &start,
//			Expiry:     &expiry,
//			Permission: &listOnly,
//		},
//	}}
//
//	setAccessPolicyOptions := container.SetAccessPolicyOptions{
//		ContainerACL: permissions,
//	}
//	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	_require.Nil(err)
//
//	serviceSASValues := service.SASSignatureValues{Identifier: "0000", ContainerName: containerName}
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	sasURL := svcClient.URL()
//	sasURL.RawQuery = queryParams.Encode()
//	sasPipeline := (NewAnonymousCredential(), PipelineOptions{})
//	sasBlobServiceURL := service.NewServiceURL(sasURL, sasPipeline)
//
//	// Verifies that the SAS can access the resource
//	sasContainer := sasBlobServiceURL.NewContainerClient(containerName)
//	resp, err := sasContainer.NewListBlobsFlatPager(context.Background(), Marker{}, ListBlobsSegmentOptions{})
//	_require.Nil(err)
//	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//
//	// Verifies that successful sas access is not just because it's public
//	anonymousBlobService := NewServiceURL(svcClient.URL(), sasPipeline)
//	anonymousContainer := anonymousBlobService.NewContainerClient(containerName)
//	_, err = anonymousContainer.NewListBlobsFlatPager(context.Background(), Marker{}, ListBlobsSegmentOptions{})
//	testcommon.ValidateBlobErrorCode(c, err, StorageErrorCodeNoAuthenticationInformation)
//}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsACLMoreThanFive() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*container.SignedIdentifier, 6)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 6; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: &access,
	}
	setAccessPolicyOptions.ContainerACL = permissions
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsDeleteAndModifyACL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	permissions := make([]*container.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: &access,
	}
	setAccessPolicyOptions.ContainerACL = permissions
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := container.SetAccessPolicyOptions{
		Access: &access,
	}
	setAccessPolicyOptions1.ContainerACL = permissions
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions1)
	_require.Nil(err)

	resp, err = containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 1)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsDeleteAllPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*container.SignedIdentifier, 2)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: to.Ptr(container.PublicAccessTypeBlob),
	}
	setAccessPolicyOptions.ContainerACL = permissions
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = container.SetAccessPolicyOptions{
		Access: to.Ptr(container.PublicAccessTypeBlob),
	}
	setAccessPolicyOptions.ContainerACL = []*container.SignedIdentifier{}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err = containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.SignedIdentifiers)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsInvalidPolicyTimes() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*container.SignedIdentifier, 2)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: to.Ptr(container.PublicAccessTypeBlob),
	}
	setAccessPolicyOptions.ContainerACL = permissions
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsNilPolicySlice() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err = containerClient.SetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsSignedIdentifierTooLong() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]*container.SignedIdentifier, 2)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: to.Ptr(container.PublicAccessTypeBlob),
	}
	setAccessPolicyOptions.ContainerACL = permissions
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.BlobPublicAccess)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.BlobPublicAccess)
}

func (s *ContainerRecordedTestsSuite) TestContainerSetPermissionsIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

// make sure that container soft delete is enabled
func (s *ContainerRecordedTestsSuite) TestContainerUndelete() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	testName := s.T().Name()
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := svcClient.NewContainerClient(containerName)

	_, err = containerClient.Create(context.Background(), nil)
	_require.Nil(err)

	_, err = containerClient.Delete(context.Background(), nil)
	_require.Nil(err)

	// it appears that deleting the container involves acquiring a lease.
	// since leases can only be 15-60s or infinite, we just wait for 60 seconds.
	time.Sleep(60 * time.Second)

	prefix := testcommon.ContainerPrefix
	listOptions := service.ListContainersOptions{Prefix: &prefix, Include: service.ListContainersInclude{Metadata: true, Deleted: true}}
	pager := svcClient.NewListContainersPager(&listOptions)

	contRestored := false
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, cont := range resp.ContainerItems {
			_require.NotNil(cont.Name)

			if *cont.Deleted && *cont.Name == containerName {
				_, err = containerClient.Restore(context.Background(), *cont.Version, nil)
				_require.NoError(err)
				contRestored = true
				break
			}
		}
		if contRestored {
			break
		}
	}

	_require.Equal(contRestored, true)

	for i := 0; i < 5; i++ {
		_, err = containerClient.Delete(context.Background(), nil)
		if err == nil {
			// container was deleted
			break
		} else if bloberror.HasCode(err, bloberror.Code("ConcurrentContainerOperationInProgress")) {
			// the container is still being restored, sleep a bit then try again
			time.Sleep(10 * time.Second)
		} else {
			// some other error
			break
		}
	}
	_require.Nil(err)
}

func (s *ContainerUnrecordedTestsSuite) TestSetAccessPoliciesInDifferentTimeFormats() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	id := "timeInEST"
	permission := "rw"
	loc, err := time.LoadLocation("EST")
	_require.Nil(err)
	start := time.Now().In(loc)
	expiry := start.Add(10 * time.Hour)

	signedIdentifiers := make([]*container.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id,
		AccessPolicy: &container.AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission,
		},
	})

	id2 := "timeInIST"
	permission2 := "r"
	loc2, err := time.LoadLocation("Asia/Kolkata")
	_require.Nil(err)
	start2 := time.Now().In(loc2)
	expiry2 := start2.Add(5 * time.Hour)

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id2,
		AccessPolicy: &container.AccessPolicy{
			Start:      &start2,
			Expiry:     &expiry2,
			Permission: &permission2,
		},
	})

	id3 := "nilTime"
	permission3 := "r"

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id3,
		AccessPolicy: &container.AccessPolicy{
			Permission: &permission3,
		},
	})
	options := container.SetAccessPolicyOptions{ContainerACL: signedIdentifiers}
	_, err = containerClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)

	// make a Get to assert three access policies
	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 3)
	_require.EqualValues(resp.SignedIdentifiers, signedIdentifiers)
}

func (s *ContainerRecordedTestsSuite) TestSetAccessPolicyWithNullId() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	signedIdentifiers := make([]*container.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		AccessPolicy: &container.AccessPolicy{
			Permission: to.Ptr("rw"),
		},
	})

	options := container.SetAccessPolicyOptions{ContainerACL: signedIdentifiers}
	_, err = containerClient.SetAccessPolicy(context.Background(), &options)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)

	resp, err := containerClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 0)
}

func (s *ContainerUnrecordedTestsSuite) TestBlobNameSpecialCharacters() {
	_require := require.New(s.T())

	const containerURL = testcommon.FakeStorageURL + "/fakecontainer"
	client, err := container.NewClientWithNoCredential(containerURL, nil)
	_require.NoError(err)
	_require.NotNil(client)

	blobNames := []string{"foo%5Cbar", "hello? sausage/Hello.txt", "世界.txt"}
	for _, blobName := range blobNames {
		expected := containerURL + "/" + url.PathEscape(blobName)

		abc := client.NewAppendBlobClient(blobName)
		_require.Equal(expected, abc.URL())

		bbc := client.NewBlockBlobClient(blobName)
		_require.Equal(expected, bbc.URL())

		pbc := client.NewPageBlobClient(blobName)
		_require.Equal(expected, pbc.URL())

		bc := client.NewBlobClient(blobName)
		_require.Equal(expected, bc.URL())
	}
}

func (s *ContainerUnrecordedTestsSuite) TestSASContainerClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Creating container client
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Adding SAS and options
	permissions := sas.ContainerPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	// ContainerSASURL is created with GetSASURL
	sasUrl, err := containerClient.GetSASURL(permissions, expiry, nil)
	_require.Nil(err)

	// Create container client with sasUrl
	_, err = container.NewClientWithNoCredential(sasUrl, nil)
	_require.Nil(err)
}
