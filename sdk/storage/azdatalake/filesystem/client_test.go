//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running blob Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &RecordedTestSuite{})
		suite.Run(t, &UnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &RecordedTestSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &RecordedTestSuite{})
	}
}

func (s *RecordedTestSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *RecordedTestSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *UnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *UnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type RecordedTestSuite struct {
	suite.Suite
}

type UnrecordedTestsSuite struct {
	suite.Suite
}

func validateFilesystemDeleted(_require *require.Assertions, filesystemClient *filesystem.Client) {
	_, err := filesystemClient.GetAccessPolicy(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ContainerNotFound)
}

func (s *RecordedTestSuite) TestCreateFilesystem() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestCreateFilesystemWithOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	testStr := "hello"
	metadata := map[string]*string{"foo": &testStr, "bar": &testStr}
	access := filesystem.Filesystem
	opts := filesystem.CreateOptions{
		Metadata: metadata,
		Access:   &access,
	}
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), &opts)
	_require.Nil(err)

	props, err := fsClient.GetProperties(context.Background(), nil)
	_require.NotNil(props.Metadata)
	_require.Equal(*props.PublicAccess, filesystem.Filesystem)
}

func (s *RecordedTestSuite) TestCreateFilesystemWithFileAccess() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	testStr := "hello"
	metadata := map[string]*string{"foo": &testStr, "bar": &testStr}
	access := filesystem.File
	opts := filesystem.CreateOptions{
		Metadata: metadata,
		Access:   &access,
	}
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), &opts)
	_require.Nil(err)
	props, err := fsClient.GetProperties(context.Background(), nil)
	_require.NotNil(props.Metadata)
	_require.Equal(*props.PublicAccess, filesystem.File)
}

func (s *RecordedTestSuite) TestCreateFilesystemEmptyMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	metadata := map[string]*string{"foo": nil, "bar": nil}
	access := filesystem.Filesystem
	opts := filesystem.CreateOptions{
		Metadata: metadata,
		Access:   &access,
	}
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), &opts)
	_require.Nil(err)

	props, err := fsClient.GetProperties(context.Background(), nil)
	_require.Nil(props.Metadata)
	_require.Equal(*props.PublicAccess, filesystem.Filesystem)

}

func (s *RecordedTestSuite) TestFilesystemCreateInvalidName() {
	_require := require.New(s.T())

	fsClient, err := testcommon.GetFilesystemClient("foo bar", s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.InvalidResourceName)
}

func (s *RecordedTestSuite) TestFilesystemCreateNameCollision() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.FilesystemAlreadyExists)
}

func (s *RecordedTestSuite) TestFilesystemGetProperties() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	resp, err := fsClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.ETag)
	_require.Nil(resp.Metadata)
}

func (s *RecordedTestSuite) TestFilesystemDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	_, err = fsClient.Delete(context.Background(), nil)
	_require.Nil(err)

	validateFilesystemDeleted(_require, fsClient)
}

func (s *RecordedTestSuite) TestFilesystemDeleteNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Delete(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ContainerNotFound)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteFilesystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFilesystemOptions)
	_require.Nil(err)
	validateFilesystemDeleted(_require, fsClient)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteFilesystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFilesystemOptions)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteFilesystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFilesystemOptions)
	_require.Nil(err)

	validateFilesystemDeleted(_require, fsClient)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteFilesystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFilesystemOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFilesystemListPaths() {
	_require := require.New(s.T())
	//testName := s.T().Name()

	//filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient("cont1", s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	//defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	//_, err = fsClient.Create(context.Background(), nil)
	//_require.Nil(err)

	resp, err := fsClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.ETag)
	_require.Nil(resp.Metadata)

	pager := fsClient.NewListPathsPager(true, nil)

	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		if err != nil {
			break
		}
	}
}
