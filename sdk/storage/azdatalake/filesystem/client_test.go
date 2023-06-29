//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running datalake Tests in %s mode\n", recordMode)
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

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.ContainerNotFound)
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

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.ContainerNotFound)
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
	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.ConditionNotMet)
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

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	opts := filesystem.SetMetadataOptions{
		Metadata: testcommon.BasicMetadata,
	}
	_, err = fsClient.SetMetadata(context.Background(), &opts)
	_require.Nil(err)

	resp1, err := fsClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	for k, v := range testcommon.BasicMetadata {
		_require.Equal(v, resp1.Metadata[k])
	}
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	opts := filesystem.SetMetadataOptions{
		Metadata: map[string]*string{},
	}

	_, err = fsClient.SetMetadata(context.Background(), &opts)
	_require.Nil(err)

	resp1, err := fsClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp1.Metadata)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetMetadata(context.Background(), nil)
	_require.Nil(err)

	resp1, err := fsClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp1.Metadata)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	opts := filesystem.SetMetadataOptions{
		Metadata: map[string]*string{"!nval!d Field!@#%": to.Ptr("value")},
	}
	_, err = fsClient.SetMetadata(context.Background(), &opts)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring), true)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.SetMetadata(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.ContainerNotFound)
}

func (s *RecordedTestSuite) TestSetEmptyAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetAccessPolicy(context.Background(), &filesystem.SetAccessPolicyOptions{})
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetNilAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)

	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		AccessPolicy: &filesystem.AccessPolicy{
			Expiry:     &expiration,
			Start:      &start,
			Permission: &permission,
		},
		ID: &id,
	})
	options := filesystem.SetAccessPolicyOptions{FilesystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetMultipleAccessPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	id := "empty"

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id,
	})

	permission2 := "r"
	id2 := "partial"

	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id2,
		AccessPolicy: &filesystem.AccessPolicy{
			Permission: &permission2,
		},
	})

	id3 := "full"
	permission3 := "r"
	start := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)
	expiry := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)

	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id3,
		AccessPolicy: &filesystem.AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission3,
		},
	})
	options := filesystem.SetAccessPolicyOptions{FilesystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)

	// Make a Get to assert two access policies
	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 3)
}

func (s *RecordedTestSuite) TestSetNullAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	id := "null"

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id,
	})
	options := filesystem.SetAccessPolicyOptions{FilesystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(len(resp.SignedIdentifiers), 1)
}

func (s *RecordedTestSuite) TestFilesystemGetSetPermissionsMultiplePolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry := start.Add(5 * time.Minute)
	expiry2 := start.Add(time.Minute)
	readWrite := to.Ptr(filesystem.AccessPolicyPermission{Read: true, Write: true}).String()
	readOnly := to.Ptr(filesystem.AccessPolicyPermission{Read: true}).String()
	id1, id2 := "0000", "0001"
	permissions := []*filesystem.SignedIdentifier{
		{ID: &id1,
			AccessPolicy: &filesystem.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &readWrite,
			},
		},
		{ID: &id2,
			AccessPolicy: &filesystem.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry2,
				Permission: &readOnly,
			},
		},
	}
	options := filesystem.SetAccessPolicyOptions{FilesystemACL: permissions}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)

	_require.Nil(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *RecordedTestSuite) TestFilesystemGetPermissionsPublicAccessNotNone() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	access := filesystem.File
	createContainerOptions := filesystem.CreateOptions{
		Access: &access,
	}
	_, err = fsClient.Create(context.Background(), &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)

	_require.Nil(err)
	_require.Equal(*resp.PublicAccess, filesystem.File)
}

// TODO: TestFilesystemSetPermissionsPublicAccessNone()

func (s *RecordedTestSuite) TestFilesystemSetPermissionsPublicAccessTypeFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)
	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.File),
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.PublicAccess, filesystem.File)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsPublicAccessFilesystem() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.Filesystem),
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.PublicAccess, filesystem.Filesystem)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsACLMoreThanFive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*filesystem.SignedIdentifier, 6)
	listOnly := to.Ptr(filesystem.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 6; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &filesystem.SignedIdentifier{
			ID: &id,
			AccessPolicy: &filesystem.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := filesystem.File
	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: &access,
	}
	setAccessPolicyOptions.FilesystemACL = permissions
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsDeleteAndModifyACL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	listOnly := to.Ptr(filesystem.AccessPolicyPermission{Read: true}).String()
	permissions := make([]*filesystem.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &filesystem.SignedIdentifier{
			ID: &id,
			AccessPolicy: &filesystem.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := filesystem.File
	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: &access,
	}
	setAccessPolicyOptions.FilesystemACL = permissions
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := filesystem.SetAccessPolicyOptions{
		Access: &access,
	}
	setAccessPolicyOptions1.FilesystemACL = permissions
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions1)
	_require.Nil(err)

	resp, err = fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 1)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsDeleteAllPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*filesystem.SignedIdentifier, 2)
	listOnly := to.Ptr(filesystem.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &filesystem.SignedIdentifier{
			ID: &id,
			AccessPolicy: &filesystem.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.File),
	}
	setAccessPolicyOptions.FilesystemACL = permissions
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.File),
	}
	setAccessPolicyOptions.FilesystemACL = []*filesystem.SignedIdentifier{}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err = fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.SignedIdentifiers)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsInvalidPolicyTimes() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*filesystem.SignedIdentifier, 2)
	listOnly := to.Ptr(filesystem.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &filesystem.SignedIdentifier{
			ID: &id,
			AccessPolicy: &filesystem.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.File),
	}
	setAccessPolicyOptions.FilesystemACL = permissions
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsNilPolicySlice() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsSignedIdentifierTooLong() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]*filesystem.SignedIdentifier, 2)
	listOnly := to.Ptr(filesystem.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		permissions[i] = &filesystem.SignedIdentifier{
			ID: &id,
			AccessPolicy: &filesystem.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.File),
	}
	setAccessPolicyOptions.FilesystemACL = permissions
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp1, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp1.PublicAccess)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp1, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp1.PublicAccess)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

//func (s *RecordedTestSuite) TestFilesystemListPaths() {
//	_require := require.New(s.T())
//	//testName := s.T().Name()
//
//	//filesystemName := testcommon.GenerateFilesystemName(testName)
//	fsClient, err := testcommon.GetFilesystemClient("cont1", s.T(), testcommon.TestAccountDatalake, nil)
//	_require.NoError(err)
//	//defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)
//
//	//_, err = fsClient.Create(context.Background(), nil)
//	//_require.Nil(err)
//
//	resp, err := fsClient.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//	_require.NotNil(resp.ETag)
//	_require.Nil(resp.Metadata)
//
//	pager := fsClient.NewListPathsPager(true, nil)
//
//	for pager.More() {
//		_, err := pager.NextPage(context.Background())
//		_require.NotNil(err)
//		if err != nil {
//			break
//		}
//	}
//}

// TODO: Lease tests
