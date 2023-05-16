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
	"os"
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

func (s *ContainerRecordedTestsSuite) TestContainerGetAccountInfo() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	cAccInfo, err := containerClient.GetAccountInfo(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(cAccInfo)
}

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

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsNonexistentPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	prefix := testcommon.BlobPrefix + testcommon.BlobPrefix
	for i := 0; i < 3; i++ {
		blobName := testcommon.GenerateBlobName(testName + strconv.Itoa(i))
		testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	}

	containerListBlobFlatSegmentOptions := container.ListBlobsFlatOptions{
		Prefix: &prefix,
	}
	pager := containerClient.NewListBlobsFlatPager(&containerListBlobFlatSegmentOptions)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(len(resp.Segment.BlobItems), 0)
	}
}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsSpecificValidPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	prefix := testcommon.BlobPrefix
	for i := 0; i < 3; i++ {
		blobName := prefix + testcommon.GenerateBlobName(testName+strconv.Itoa(i))
		testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	}

	containerListBlobFlatSegmentOptions := container.ListBlobsFlatOptions{
		Prefix: &prefix,
	}
	pager := containerClient.NewListBlobsFlatPager(&containerListBlobFlatSegmentOptions)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(len(resp.Segment.BlobItems), 3)
	}
}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsValidDelimiter() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	prefixes := []string{"a/1", "a/2", "b/1", "blob"}
	pre := []string{"a/", "b/"}
	prefixMap := testcommon.BlobListToMap(pre)
	for _, prefix := range prefixes {
		blobName := prefix + testcommon.GenerateBlobName(testName)
		testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	}

	pager := containerClient.NewListBlobsHierarchyPager("/", nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, prefix := range resp.Segment.BlobPrefixes {
			_require.Equal(prefixMap[*prefix.Name], true) // checks if prefix exists in prefixes
		}
		if err != nil {
			break
		}
	}
}

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

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsWithSnapshots() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// initialize a blob and create a snapshot of it
	snapBlobName := testcommon.GenerateBlobName(testName)
	snapBlob := testcommon.CreateNewBlockBlob(context.Background(), _require, snapBlobName, containerClient)
	snap, err := snapBlob.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	listBlobFlatSegmentOptions := container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	}
	pager := containerClient.NewListBlobsFlatPager(&listBlobFlatSegmentOptions)

	wasFound := false // hold the for loop accountable for finding the blob and it's snapshot
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		for _, blob := range resp.Segment.BlobItems {
			if *blob.Name == snapBlobName && blob.Snapshot != nil {
				wasFound = true
				_require.Equal(*blob.Snapshot, *snap.Snapshot)
			}
		}
		if err != nil {
			break
		}
	}
	_require.Equal(wasFound, true)
}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeCopy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobCopyName := testcommon.GenerateBlobName("copy" + testName)
	blobCopyClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobCopyName, containerClient)

	_, err = blobCopyClient.StartCopyFromURL(context.Background(), bbClient.URL(), nil)
	_require.Nil(err)

	opts := container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Copy: true},
	}
	pager := containerClient.NewListBlobsFlatPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		// These are sufficient to show that the blob copy was in fact included
		_require.Nil(err)
		_require.Equal(len(resp.Segment.BlobItems), 2)
		_require.EqualValues(*resp.Segment.BlobItems[1].Name, blobCopyName)
		_require.EqualValues(*resp.Segment.BlobItems[0].Name, blobName)
		_require.Equal(*resp.Segment.BlobItems[0].Properties.ContentLength, int64(len(testcommon.BlockBlobDefaultData)))
		_require.Equal(*resp.Segment.BlobItems[1].Properties.CopyStatus, container.CopyStatusTypeSuccess)
	}
}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeUncommitted() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blockID := testcommon.GenerateBlockIDsList(1)

	_, err = bbClient.StageBlock(context.Background(), blockID[0], streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	opts := container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{UncommittedBlobs: false},
	}
	pager := containerClient.NewListBlobsFlatPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Equal(len(resp.Segment.BlobItems), 1)
		_require.EqualValues(*resp.Segment.BlobItems[0].Name, blobName)
	}

}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeTypeDeletedWithVersion() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)

	opts := container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{DeletedWithVersions: true},
	}
	pager := containerClient.NewListBlobsFlatPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Equal(len(resp.Segment.BlobItems), 1)
		_require.EqualValues(*resp.Segment.BlobItems[0].Name, blobName)
	}

	deleteOpts := blob.DeleteOptions{
		DeleteSnapshots: to.Ptr(blob.DeleteSnapshotsOptionTypeInclude),
	}
	_, err = bbClient.Delete(context.Background(), &deleteOpts)
	_require.Nil(err)

	pager = containerClient.NewListBlobsFlatPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Equal(len(resp.Segment.BlobItems), 1)
		_require.EqualValues(*resp.Segment.BlobItems[0].Name, blobName)
	}
}

func (s *ContainerRecordedTestsSuite) TestContainerListBlobsIncludeMultipleImpl() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "z"+blobName, containerClient)
	_, err = bbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	bbClient2 := testcommon.CreateNewBlockBlob(context.Background(), _require, "copy"+blobName, containerClient)
	_, err = bbClient2.StartCopyFromURL(context.Background(), bbClient.URL(), nil)
	_require.Nil(err)

	// Copy should finish within one minute
	time.Sleep(60 * time.Second)

	bbClient3 := testcommon.CreateNewBlockBlob(context.Background(), _require, "deleted"+blobName, containerClient)
	_, err = bbClient3.Delete(context.Background(), nil)
	_require.NoError(err)

	opts := container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Copy: true, Deleted: true, Versions: true},
	}
	pager := containerClient.NewListBlobsFlatPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())

		// These are sufficient to show that the blob copy was in fact included
		_require.Nil(err)
		_require.Equal(len(resp.Segment.BlobItems), 6)
		_require.Equal(*resp.Segment.BlobItems[1].Properties.CopyStatus, container.CopyStatusTypeSuccess)
	}
}

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
			_require.Equal(nameMap[*blob.Name], true) // checks if name exists in blob names
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
			_require.Equal(nameMap[*blob.Name], true) // checks if name exists in blob names
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

func (s *ContainerUnrecordedTestsSuite) TestSASContainerClientTags() {
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
		Tag:    true,
	}
	expiry := time.Now().Add(time.Hour)

	// ContainerSASURL is created with GetSASURL
	sasUrl, err := containerClient.GetSASURL(permissions, expiry, nil)
	_require.Nil(err)

	// Create container client with sasUrl
	containerSasClient, err := container.NewClientWithNoCredential(sasUrl, nil)
	_require.Nil(err)

	// Get Blob Client
	blobSasClient := containerSasClient.NewAppendBlobClient(testName)
	_, err = blobSasClient.Create(context.Background(), nil)
	_require.Nil(err)

	// Try getting tags with container SAS
	_, err = blobSasClient.GetTags(context.Background(), nil)
	_require.Nil(err)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchDeleteSuccessUsingSharedKey() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
		err = bb.Delete(bbName, nil)
		_require.NoError(err)
	}

	pager := containerClient.NewListBlobsFlatPager(nil)
	ctr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 10)

	resp, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	for _, subResp := range resp.Responses {
		_require.NotNil(subResp.ContentID)
		_require.NotNil(subResp.ContainerName)
		_require.NotNil(subResp.BlobName)
		_require.NotNil(subResp.RequestID)
		_require.NotNil(subResp.Version)
		_require.NoError(subResp.Error)
	}

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 0)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchSetTierPartialFailureUsingSharedKey() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	// add 5 blobs to BatchBuilder which does not exist
	for i := 0; i < 15; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		if i < 10 {
			_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
		}
		err = bb.SetTier(bbName, blob.AccessTierCool, nil)
		_require.NoError(err)
	}

	pager := containerClient.NewListBlobsFlatPager(nil)
	var ctrHot, ctrCool = 0, 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 10)
	_require.Equal(ctrCool, 0)

	resp, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	var ctrSuccess, ctrFailure = 0, 0
	for _, subResp := range resp.Responses {
		_require.NotNil(subResp.ContentID)
		_require.NotNil(subResp.ContainerName)
		_require.NotNil(subResp.BlobName)
		_require.NotNil(subResp.RequestID)
		_require.NotNil(subResp.Version)
		if subResp.Error == nil {
			ctrSuccess++
		} else {
			ctrFailure++
			_require.NotEmpty(subResp.Error.Error())
			testcommon.ValidateBlobErrorCode(_require, subResp.Error, bloberror.BlobNotFound)
		}
	}
	_require.Equal(ctrSuccess, 10)
	_require.Equal(ctrFailure, 5)

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctrHot = 0
	ctrCool = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 0)
	_require.Equal(ctrCool, 10)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchDeleteUsingTokenCredential() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient, err := container.NewClient("https://"+accountName+".blob.core.windows.net/"+containerName, cred, nil)
	_require.NoError(err)

	_, err = containerClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
		err = bb.Delete(bbName, nil)
		_require.NoError(err)
	}

	pager := containerClient.NewListBlobsFlatPager(nil)
	ctr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 10)

	resp, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 0)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchSetTierUsingTokenCredential() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient, err := container.NewClient("https://"+accountName+".blob.core.windows.net/"+containerName, cred, nil)
	_require.NoError(err)

	_, err = containerClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
		err = bb.SetTier(bbName, blob.AccessTierCool, nil)
		_require.NoError(err)
	}

	pager := containerClient.NewListBlobsFlatPager(nil)
	var ctrHot, ctrCool = 0, 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 10)
	_require.Equal(ctrCool, 0)

	resp, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctrHot = 0
	ctrCool = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 0)
	_require.Equal(ctrCool, 10)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchDeleteUsingAccountSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountSAS, err := testcommon.GetAccountSAS(sas.AccountPermissions{Read: true, Create: true, Write: true, List: true, Add: true, Delete: true},
		sas.AccountResourceTypes{Service: true, Container: true, Object: true})
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient, err := container.NewClientWithNoCredential("https://"+accountName+".blob.core.windows.net/"+containerName+"?"+accountSAS, nil)
	_require.NoError(err)

	_, err = containerClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
		err = bb.Delete(bbName, nil)
		_require.NoError(err)
	}

	pager := containerClient.NewListBlobsFlatPager(nil)
	ctr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 10)

	resp, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 0)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchSetTierUsingAccountSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountSAS, err := testcommon.GetAccountSAS(sas.AccountPermissions{Read: true, Create: true, Write: true, List: true, Add: true, Delete: true},
		sas.AccountResourceTypes{Service: true, Container: true, Object: true})
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient, err := container.NewClientWithNoCredential("https://"+accountName+".blob.core.windows.net/"+containerName+"?"+accountSAS, nil)
	_require.NoError(err)

	_, err = containerClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
		err = bb.SetTier(bbName, blob.AccessTierCool, nil)
		_require.NoError(err)
	}

	pager := containerClient.NewListBlobsFlatPager(nil)
	var ctrHot, ctrCool = 0, 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 10)
	_require.Equal(ctrCool, 0)

	resp, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctrHot = 0
	ctrCool = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 0)
	_require.Equal(ctrCool, 10)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchDeleteUsingServiceSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	cntClientSharedKey := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, cntClientSharedKey)

	serviceSAS, err := testcommon.GetServiceSAS(containerName, sas.BlobPermissions{Read: true, Create: true, Write: true, List: true, Add: true, Delete: true})
	_require.NoError(err)

	cntClientSAS, err := container.NewClientWithNoCredential(cntClientSharedKey.URL()+"?"+serviceSAS, nil)
	_require.NoError(err)

	bb, err := cntClientSAS.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, cntClientSAS)
		err = bb.Delete(bbName, nil)
		_require.NoError(err)
	}

	pager := cntClientSAS.NewListBlobsFlatPager(nil)
	ctr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 10)

	resp, err := cntClientSAS.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = cntClientSAS.NewListBlobsFlatPager(nil)
	ctr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 0)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchSetTierUsingServiceSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	cntClientSharedKey := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, cntClientSharedKey)

	serviceSAS, err := testcommon.GetServiceSAS(containerName, sas.BlobPermissions{Read: true, Create: true, Write: true, List: true})
	_require.NoError(err)

	cntClientSAS, err := container.NewClientWithNoCredential(cntClientSharedKey.URL()+"?"+serviceSAS, nil)
	_require.NoError(err)

	bb, err := cntClientSAS.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, cntClientSAS)
		err = bb.SetTier(bbName, blob.AccessTierCool, nil)
		_require.NoError(err)
	}

	pager := cntClientSAS.NewListBlobsFlatPager(nil)
	var ctrHot, ctrCool = 0, 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 10)
	_require.Equal(ctrCool, 0)

	resp, err := cntClientSAS.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = cntClientSAS.NewListBlobsFlatPager(nil)
	ctrHot = 0
	ctrCool = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 0)
	_require.Equal(ctrCool, 10)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchDeleteUsingUserDelegationSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := service.NewClient("https://"+accountName+".blob.core.windows.net/", cred, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	cntClientTokenCred := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, cntClientTokenCred)

	udSAS, err := testcommon.GetUserDelegationSAS(svcClient, containerName, sas.BlobPermissions{Read: true, Create: true, Write: true, List: true, Add: true, Delete: true})
	_require.NoError(err)

	cntClientSAS, err := container.NewClientWithNoCredential(cntClientTokenCred.URL()+"?"+udSAS, nil)
	_require.NoError(err)

	bb, err := cntClientSAS.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, cntClientSAS)
		err = bb.Delete(bbName, nil)
		_require.NoError(err)
	}

	pager := cntClientSAS.NewListBlobsFlatPager(nil)
	ctr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 10)

	resp, err := cntClientSAS.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = cntClientSAS.NewListBlobsFlatPager(nil)
	ctr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 0)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchSetTierUsingUserDelegationSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := service.NewClient("https://"+accountName+".blob.core.windows.net/", cred, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	cntClientTokenCred := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, cntClientTokenCred)

	udSAS, err := testcommon.GetUserDelegationSAS(svcClient, containerName, sas.BlobPermissions{Read: true, Create: true, Write: true, List: true})
	_require.NoError(err)

	cntClientSAS, err := container.NewClientWithNoCredential(cntClientTokenCred.URL()+"?"+udSAS, nil)
	_require.NoError(err)

	bb, err := cntClientSAS.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, cntClientSAS)
		err = bb.SetTier(bbName, blob.AccessTierCool, nil)
		_require.NoError(err)
	}

	pager := cntClientSAS.NewListBlobsFlatPager(nil)
	var ctrHot, ctrCool = 0, 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 10)
	_require.Equal(ctrCool, 0)

	resp, err := cntClientSAS.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	pager = cntClientSAS.NewListBlobsFlatPager(nil)
	ctrHot = 0
	ctrCool = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		for _, blobItem := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			if *blobItem.Properties.AccessTier == container.AccessTierHot {
				ctrHot++
			} else if *blobItem.Properties.AccessTier == container.AccessTierCool {
				ctrCool++
			}
		}
	}
	_require.Equal(ctrHot, 0)
	_require.Equal(ctrCool, 10)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchDeleteMoreThan256() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	for i := 0; i < 256; i++ {
		bbName := fmt.Sprintf("blockblob%v", i)
		_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
		err = bb.Delete(bbName, nil)
		_require.NoError(err)
	}

	pager := containerClient.NewListBlobsFlatPager(nil)
	ctr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 256)

	resp, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
	for _, subResp := range resp.Responses {
		_require.Nil(subResp.Error)
	}

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 0)

	// add more items to make batch size more than 256
	for i := 0; i < 10; i++ {
		bbName := fmt.Sprintf("fakeblob%v", i)
		err = bb.Delete(bbName, nil)
		_require.NoError(err)
	}

	resp2, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.Error(err)
	_require.Nil(resp2.RequestID)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchDeleteForOneBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bb, err := containerClient.NewBatchBuilder()
	_require.NoError(err)

	bbName := "blockblob1"
	_ = testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)
	err = bb.Delete(bbName, nil)
	_require.NoError(err)

	pager := containerClient.NewListBlobsFlatPager(nil)
	ctr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 1)

	resp1, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp1.RequestID)
	_require.Equal(len(resp1.Responses), 1)
	_require.NoError(resp1.Responses[0].Error)

	pager = containerClient.NewListBlobsFlatPager(nil)
	ctr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err)
		ctr += len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems)
	}
	_require.Equal(ctr, 0)

	resp2, err := containerClient.SubmitBatch(context.Background(), bb, nil)
	_require.NoError(err)
	_require.NotNil(resp2.RequestID)
	_require.Equal(len(resp2.Responses), 1)
	_require.Error(resp2.Responses[0].Error)
	testcommon.ValidateBlobErrorCode(_require, resp2.Responses[0].Error, bloberror.BlobNotFound)
}

func (s *ContainerUnrecordedTestsSuite) TestContainerBlobBatchErrors() {
	_require := require.New(s.T())

	svcClient, err := service.NewClientWithNoCredential("https://fakestorageaccount.blob.core.windows.net/", nil)
	_require.NoError(err)

	cntClient := svcClient.NewContainerClient("fakecontainer")

	bb1, err := cntClient.NewBatchBuilder()
	_require.NoError(err)

	// adding multiple operations to BatchBuilder
	err = bb1.Delete("blob1", nil)
	_require.NoError(err)

	err = bb1.SetTier("blob2", blob.AccessTierCool, nil)
	_require.Error(err)

	bb2, err := cntClient.NewBatchBuilder()
	_require.NoError(err)

	// submitting empty batch
	_, err = cntClient.SubmitBatch(context.Background(), bb2, nil)
	_require.Error(err)

	// submitting nil BatchBuilder
	_, err = cntClient.SubmitBatch(context.Background(), nil, nil)
	_require.Error(err)
}
