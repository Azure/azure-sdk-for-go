//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem_test

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running datalake Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &RecordedTestSuite{})
		suite.Run(t, &UnrecordedTestSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &RecordedTestSuite{})
	case recording.RecordingMode:
		suite.Run(t, &RecordedTestSuite{})
	}
}

func (s *RecordedTestSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *RecordedTestSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
}

func (s *RecordedTestSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *RecordedTestSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *UnrecordedTestSuite) BeforeTest(suite string, test string) {

}

func (s *UnrecordedTestSuite) AfterTest(suite string, test string) {

}

type RecordedTestSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type UnrecordedTestSuite struct {
	suite.Suite
}

func validateFileSystemDeleted(_require *require.Assertions, filesystemClient *filesystem.Client) {
	_, err := filesystemClient.GetAccessPolicy(context.Background(), nil)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.FileSystemNotFound)
}

func (s *RecordedTestSuite) TestCreateFilesystem() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestCreateFilesystemWithOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()
	filesystemName := testcommon.GenerateFileSystemName(testName)
	testStr := "hello"
	metadata := map[string]*string{"foo": &testStr, "bar": &testStr}

	opts := filesystem.CreateOptions{
		Metadata:     metadata,
		CPKScopeInfo: &testcommon.TestCPKScopeInfo,
	}
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), &opts)
	_require.NoError(err)

	props, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(props.Metadata)
	_require.Equal(props.DefaultEncryptionScope, &testcommon.TestEncryptionScope)
}

func (s *RecordedTestSuite) TestCreateFilesystemWithFileAccess() {
	s.T().Skip("This test is not valid because public access is disabled")
	_require := require.New(s.T())
	testName := s.T().Name()
	filesystemName := testcommon.GenerateFileSystemName(testName)
	testStr := "hello"
	metadata := map[string]*string{"foo": &testStr, "bar": &testStr}
	access := filesystem.File
	opts := filesystem.CreateOptions{
		Metadata: metadata,
		Access:   &access,
	}
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), &opts)
	_require.NoError(err)
	props, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(props.Metadata)
	_require.Equal(*props.PublicAccess, filesystem.File)
}
func (s *RecordedTestSuite) TestCreateFilesystemEmptyMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	metadata := map[string]*string{"foo": nil, "bar": nil}
	opts := filesystem.CreateOptions{
		Metadata: metadata,
	}
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), &opts)
	_require.NoError(err)

	props, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(props.Metadata)

}

func (s *RecordedTestSuite) TestFilesystemCreateInvalidName() {
	_require := require.New(s.T())

	fsClient, err := testcommon.GetFileSystemClient("foo bar", s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.InvalidResourceName)
}

func (s *RecordedTestSuite) TestFilesystemCreateNameCollision() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.FileSystemAlreadyExists)
}

func (s *RecordedTestSuite) TestFilesystemGetProperties() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	resp, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.Nil(resp.Metadata)
}

func (s *RecordedTestSuite) TestFilesystemGetPropertiesWithEmptyOpts() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	opts := &filesystem.GetPropertiesOptions{}
	resp, err := fsClient.GetProperties(context.Background(), opts)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.Nil(resp.Metadata)
}

func (s *RecordedTestSuite) TestFilesystemGetPropertiesWithLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fsLeaseClient, err := lease.NewFileSystemClient(fsClient, &lease.FileSystemClientOptions{LeaseID: proposedLeaseIDs[0]})
	_require.NoError(err)
	_, err = fsLeaseClient.AcquireLease(context.Background(), int32(60), nil)
	_require.NoError(err)

	opts := &filesystem.GetPropertiesOptions{LeaseAccessConditions: &filesystem.LeaseAccessConditions{
		LeaseID: fsLeaseClient.LeaseID(),
	}}

	resp, err := fsClient.GetProperties(context.Background(), opts)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.Nil(resp.Metadata)

	_, err = fsLeaseClient.ReleaseLease(context.Background(), nil)
	_require.NoError(err)
}
func (s *RecordedTestSuite) TestFilesystemGetPropertiesDefaultEncryptionScopeAndOverride() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	testStr := "hello"
	metadata := map[string]*string{"foo": &testStr, "bar": &testStr}

	opts := filesystem.CreateOptions{
		Metadata:     metadata,
		CPKScopeInfo: &testcommon.TestCPKScopeInfo,
	}
	_, err = fsClient.Create(context.Background(), &opts)
	_require.NoError(err)

	resp, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(resp.DenyEncryptionScopeOverride, to.Ptr(false))
	_require.Equal(resp.DefaultEncryptionScope, &testcommon.TestEncryptionScope)

}

func (s *RecordedTestSuite) TestFilesystemDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fsClient.Delete(context.Background(), nil)
	_require.NoError(err)

	validateFileSystemDeleted(_require, fsClient)
}

func (s *RecordedTestSuite) TestFilesystemDeleteNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Delete(context.Background(), nil)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.FileSystemNotFound)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteFileSystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFileSystemOptions)
	_require.NoError(err)
	validateFileSystemDeleted(_require, fsClient)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteFileSystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFileSystemOptions)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteFileSystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFileSystemOptions)
	_require.NoError(err)

	validateFileSystemDeleted(_require, fsClient)
}

func (s *RecordedTestSuite) TestFilesystemDeleteIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteFileSystemOptions := filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = fsClient.Delete(context.Background(), &deleteFileSystemOptions)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	opts := filesystem.SetMetadataOptions{
		Metadata: testcommon.BasicMetadata,
	}
	_, err = fsClient.SetMetadata(context.Background(), &opts)
	_require.NoError(err)

	resp1, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	for k, v := range testcommon.BasicMetadata {
		_require.Equal(v, resp1.Metadata[k])
	}
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	opts := filesystem.SetMetadataOptions{
		Metadata: map[string]*string{},
	}

	_, err = fsClient.SetMetadata(context.Background(), &opts)
	_require.NoError(err)

	resp1, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp1.Metadata)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetMetadata(context.Background(), nil)
	_require.NoError(err)

	resp1, err := fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp1.Metadata)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	opts := filesystem.SetMetadataOptions{
		Metadata: map[string]*string{"!nval!d Field!@#%": to.Ptr("value")},
	}
	_, err = fsClient.SetMetadata(context.Background(), &opts)
	_require.Error(err)
	_require.Equal(strings.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring), true)
}

func (s *RecordedTestSuite) TestFilesystemSetMetadataNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.SetMetadata(context.Background(), nil)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.FileSystemNotFound)
}

func (s *RecordedTestSuite) TestFilesystemSetEmptyAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetAccessPolicy(context.Background(), &filesystem.SetAccessPolicyOptions{})
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFilesystemSetNilAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFilesystemSetAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

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
	options := filesystem.SetAccessPolicyOptions{FileSystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFilesystemSetMultipleAccessPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

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
	options := filesystem.SetAccessPolicyOptions{FileSystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.NoError(err)

	// Make a Get to assert two access policies
	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.SignedIdentifiers, 3)
}

func (s *RecordedTestSuite) TestFilesystemSetNullAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	id := "null"

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id,
	})
	options := filesystem.SetAccessPolicyOptions{FileSystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.NoError(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(len(resp.SignedIdentifiers), 1)
}

func (s *RecordedTestSuite) TestFilesystemGetAccessPolicyWithEmptyOpts() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	id := "null"

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id,
	})
	options := filesystem.SetAccessPolicyOptions{FileSystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.NoError(err)

	opts := &filesystem.GetAccessPolicyOptions{}

	resp, err := fsClient.GetAccessPolicy(context.Background(), opts)
	_require.NoError(err)
	_require.Equal(len(resp.SignedIdentifiers), 1)
}

func (s *RecordedTestSuite) TestFilesystemGetAccessPolicyOnLeasedFilesystem() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

	fsLeaseClient, err := lease.NewFileSystemClient(fsClient, &lease.FileSystemClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	id := "null"

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id,
	})
	options := filesystem.SetAccessPolicyOptions{FileSystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.NoError(err)

	_, err = fsLeaseClient.AcquireLease(context.Background(), int32(15), nil)
	_require.NoError(err)

	opts1 := &filesystem.GetAccessPolicyOptions{
		LeaseAccessConditions: &filesystem.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		},
	}

	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.NoError(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), opts1)
	_require.NoError(err)
	_require.Equal(len(resp.SignedIdentifiers), 1)

	_, err = fsClient.Delete(context.Background(), &filesystem.DeleteOptions{AccessConditions: &filesystem.AccessConditions{
		LeaseAccessConditions: &filesystem.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		},
	}})
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFilesystemGetSetPermissionsMultiplePolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
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
	options := filesystem.SetAccessPolicyOptions{FileSystemACL: permissions}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)

	_require.NoError(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *RecordedTestSuite) TestFilesystemGetPermissionsPublicAccessNotNone() {
	s.T().Skip("this test is not needed as public access is disabled")
	_require := require.New(s.T())
	testName := s.T().Name()
	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	access := filesystem.File
	createContainerOptions := filesystem.CreateOptions{
		Access: &access,
	}
	_, err = fsClient.Create(context.Background(), &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)

	_require.NoError(err)
	_require.Equal(*resp.PublicAccess, filesystem.File)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsPublicAccessTypeFile() {
	s.T().Skip("this test is not needed as public access is disabled")
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)
	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.File),
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.PublicAccess, filesystem.File)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsPublicAccessFilesystem() {
	s.T().Skip("this test is not needed as public access is disabled")
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		Access: to.Ptr(filesystem.FileSystem),
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.PublicAccess, filesystem.FileSystem)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsACLMoreThanFive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		FileSystemACL: permissions,
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsDeleteAndModifyACL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		FileSystemACL: permissions,
	}

	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := filesystem.SetAccessPolicyOptions{
		FileSystemACL: permissions,
	}

	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions1)
	_require.NoError(err)

	resp, err = fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.SignedIdentifiers, 1)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsDeleteAllPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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
		FileSystemACL: permissions,
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = filesystem.SetAccessPolicyOptions{
		FileSystemACL: []*filesystem.SignedIdentifier{},
	}

	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp, err = fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.SignedIdentifiers)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsInvalidPolicyTimes() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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
		FileSystemACL: permissions,
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsNilPolicySlice() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.SetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsSignedIdentifierTooLong() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
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
		FileSystemACL: permissions,
	}

	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp1, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp1.PublicAccess)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp1, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp1.PublicAccess)
}

func (s *RecordedTestSuite) TestFilesystemSetPermissionsIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	setAccessPolicyOptions := filesystem.SetAccessPolicyOptions{
		AccessConditions: &filesystem.AccessConditions{
			ModifiedAccessConditions: &filesystem.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = fsClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *UnrecordedTestSuite) TestFilesystemSetAccessPoliciesInDifferentTimeFormats() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	id := "timeInEST"
	permission := "rw"
	loc, err := time.LoadLocation("EST")
	_require.NoError(err)
	start := time.Now().In(loc)
	expiry := start.Add(10 * time.Hour)

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id,
		AccessPolicy: &filesystem.AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission,
		},
	})

	id2 := "timeInIST"
	permission2 := "r"
	loc2, err := time.LoadLocation("Asia/Kolkata")
	_require.NoError(err)
	start2 := time.Now().In(loc2)
	expiry2 := start2.Add(5 * time.Hour)

	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id2,
		AccessPolicy: &filesystem.AccessPolicy{
			Start:      &start2,
			Expiry:     &expiry2,
			Permission: &permission2,
		},
	})

	id3 := "nilTime"
	permission3 := "r"

	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		ID: &id3,
		AccessPolicy: &filesystem.AccessPolicy{
			Permission: &permission3,
		},
	})
	options := filesystem.SetAccessPolicyOptions{FileSystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.NoError(err)

	// make a Get to assert three access policies
	resp1, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp1.SignedIdentifiers, 3)
	_require.EqualValues(resp1.SignedIdentifiers, signedIdentifiers)
}

func (s *RecordedTestSuite) TestFilesystemSetAccessPolicyWithNullId() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	signedIdentifiers := make([]*filesystem.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &filesystem.SignedIdentifier{
		AccessPolicy: &filesystem.AccessPolicy{
			Permission: to.Ptr("rw"),
		},
	})

	options := filesystem.SetAccessPolicyOptions{FileSystemACL: signedIdentifiers}
	_, err = fsClient.SetAccessPolicy(context.Background(), &options)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.InvalidXMLDocument)

	resp1, err := fsClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp1.SignedIdentifiers, 0)
}

func (s *UnrecordedTestSuite) TestSASFileSystemClient() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	// Adding SAS and options
	permissions := sas.FileSystemPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	// filesystemSASURL is created with GetSASURL
	sasUrl, err := fsClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	// Create filesystem client with sasUrl
	_, err = filesystem.NewClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestCreateFileSystemDirectoryClientWithSpecialDirName() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirClient := fsClient.NewDirectoryClient("#,%,?/%,#")
	_require.NoError(err)

	response, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(response)

	// Perform an operation on dfs endpoint
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	_, err = dirClient.SetAccessControl(context.Background(), &file.SetAccessControlOptions{Owner: &owner})
	_require.NoError(err)

	// Perform an operation on blob endpoint
	_, err = dirClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestCreateFileSystemFileClientWithSpecialFileName() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileClient := fsClient.NewFileClient("#,%,?/#")
	_require.NoError(err)

	response, err := fileClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(response)

	// Perform an operation on dfs endpoint
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	_, err = fileClient.SetAccessControl(context.Background(), &file.SetAccessControlOptions{Owner: &owner})
	_require.NoError(err)

	// Perform an operation on blob endpoint
	_, err = fileClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestCreateFileClientFromDirectoryClientWithSpecialFileName() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirClient := fsClient.NewDirectoryClient("#,%,?/%,#")
	_require.NoError(err)

	fileClient, err := dirClient.NewFileClient("%,?/##")
	_require.NoError(err)

	response, err := fileClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(response)

	_, err = fileClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFilesystemListPathsWithRecursive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	pager := fsClient.NewListPathsPager(true, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(5, len(resp.Paths))
		_require.NotNil(resp.PathList.Paths[0].IsDirectory)

		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListPathsRecursiveWithEtagCheck() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	pager := fsClient.NewListPathsPager(true, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, p := range resp.Paths {
			_require.NotNil(p.ETag)
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListPathsWithRecursiveNoPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient("file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient("file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient("dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient("dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	pager := fsClient.NewListPathsPager(true, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(4, len(resp.Paths))
		_require.NotNil(resp.PathList.Paths[0].IsDirectory)
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListPathsWithoutRecursive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	pager := fsClient.NewListPathsPager(false, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(1, len(resp.Paths))
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListPathsWithMaxResults() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	opts := filesystem.ListPathsOptions{
		MaxResults: to.Ptr(int32(2)),
	}
	pages := 3
	count := 0
	pager := fsClient.NewListPathsPager(true, &opts)
	for pager.More() {
		_, err = pager.NextPage(context.Background())
		_require.NoError(err)
		count += 1
		if err != nil {
			break
		}
	}
	_require.Equal(pages, count)
}

func (s *RecordedTestSuite) TestFilesystemListPathsWithPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	opts := filesystem.ListPathsOptions{
		Prefix: to.Ptr("Test"),
	}
	pager := fsClient.NewListPathsPager(true, &opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(4, len(resp.Paths))
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListPathsWithContinuation() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	opts := filesystem.ListPathsOptions{
		MaxResults: to.Ptr(int32(3)),
	}
	pager := fsClient.NewListPathsPager(true, &opts)

	resp, err := pager.NextPage(context.Background())
	_require.NoError(err)
	_require.Equal(3, len(resp.Paths))
	_require.NotNil(resp.Continuation)

	token := resp.Continuation
	pager = fsClient.NewListPathsPager(true, &filesystem.ListPathsOptions{
		Marker: token,
	})
	resp, err = pager.NextPage(context.Background())
	_require.NoError(err)
	_require.Equal(2, len(resp.Paths))
	_require.Nil(resp.Continuation)
}

func (s *RecordedTestSuite) TestFilesystemListPathsWithEncryptionContext() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), &file.CreateOptions{EncryptionContext: &testcommon.TestEncryptionContext})
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), &file.CreateOptions{EncryptionContext: &testcommon.TestEncryptionContext})
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	pager := fsClient.NewListPathsPager(true, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(5, len(resp.Paths))
		_require.Equal(resp.PathList.Paths[2].IsDirectory, to.Ptr(true))
		_require.Nil(resp.PathList.Paths[3].IsDirectory)
		_require.Nil(resp.PathList.Paths[2].EncryptionContext)
		// Encryption context is only applicable on files, not directories.
		_require.Equal(resp.PathList.Paths[3].EncryptionContext, &testcommon.TestEncryptionContext)

		if err != nil {
			break
		}
	}
}

func (s *UnrecordedTestSuite) TestFilesystemListDirectoryPaths() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	pager := fsClient.NewListDirectoryPathsPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(3, len(resp.Segment.PathItems))
		if err != nil {
			break
		}
	}
}

func (s *UnrecordedTestSuite) TestFilesystemListDirectoryPathsMaxResults() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	pages := 3
	count := 0
	opts := filesystem.ListDirectoryPathsOptions{
		MaxResults: to.Ptr(int32(1)),
	}

	pager := fsClient.NewListDirectoryPathsPager(&opts)
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NoError(err)
		count += 1
		if err != nil {
			break
		}
	}
	_require.Equal(pages, count)
}

func (s *UnrecordedTestSuite) TestFilesystemListDirectoryPathsWithPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	opts := filesystem.ListDirectoryPathsOptions{
		Prefix: to.Ptr("Test"),
	}

	pager := fsClient.NewListDirectoryPathsPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(3, len(resp.Segment.PathItems))
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListDeletedPaths() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)

	pager := fsClient.NewListDeletedPathsPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(1, len(resp.Segment.PathItems))
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListDeletedPathsWithMaxResults() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = client.Delete(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = client.Delete(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)

	opts := filesystem.ListDeletedPathsOptions{
		MaxResults: to.Ptr(int32(2)),
	}
	pages := 2
	count := 0
	pager := fsClient.NewListDeletedPathsPager(&opts)
	for pager.More() {
		_, err = pager.NextPage(context.Background())
		_require.NoError(err)
		count += 1
		if err != nil {
			break
		}
	}
	_require.Equal(pages, count)
}

func (s *RecordedTestSuite) TestFilesystemListDeletedPathsWithPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = client.Delete(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = client.Delete(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)

	opts := filesystem.ListDeletedPathsOptions{
		Prefix: to.Ptr("Test"),
	}
	pager := fsClient.NewListDeletedPathsPager(&opts)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(4, len(resp.Segment.PathItems))
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestFilesystemListDeletedPathsWithContinuation() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	client := fsClient.NewFileClient(testName + "file1")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = client.Delete(context.Background(), nil)
	_require.NoError(err)
	client = fsClient.NewFileClient(testName + "file2")
	_, err = client.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = client.Delete(context.Background(), nil)
	_require.NoError(err)
	dirClient := fsClient.NewDirectoryClient(testName + "dir1")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)
	dirClient = fsClient.NewDirectoryClient(testName + "dir2")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)

	opts := filesystem.ListDeletedPathsOptions{
		MaxResults: to.Ptr(int32(3)),
	}
	pager := fsClient.NewListDeletedPathsPager(&opts)

	resp, err := pager.NextPage(context.Background())
	_require.NoError(err)
	_require.Equal(3, len(resp.Segment.PathItems))
	_require.NotNil(resp.NextMarker)

	token := resp.NextMarker
	pager = fsClient.NewListDeletedPathsPager(&filesystem.ListDeletedPathsOptions{
		Marker: token,
	})
	resp, err = pager.NextPage(context.Background())
	_require.NoError(err)
	_require.Equal(1, len(resp.Segment.PathItems))
	_require.Equal("", *resp.NextMarker)
}

func (s *UnrecordedTestSuite) TestSASFileSystemCreateAndDeleteFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	// Adding SAS and options
	permissions := sas.FileSystemPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	// filesystemSASURL is created with GetSASURL
	sasUrl, err := fsClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	// Create filesystem client with sasUrl
	client, err := filesystem.NewClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	fClient := client.NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestSASFileSystemCreateAndDeleteDirectory() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	// Adding SAS and options
	permissions := sas.FileSystemPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	// filesystemSASURL is created with GetSASURL
	sasUrl, err := fsClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	// Create filesystem client with sasUrl
	client, err := filesystem.NewClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	dClient := client.NewDirectoryClient(testcommon.GenerateDirName(testName))
	_, err = dClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestFSCreateDeleteUsingOAuth() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName

	fsClient, err := filesystem.NewClient(fsURL, cred, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

}

func (s *RecordedTestSuite) TestCreateFileInFileSystemSetOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	umask := "0000"
	user := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := user
	acl := "user::rwx,group::r-x,other::rwx"
	leaseDuration := to.Ptr(int64(15))

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	createFileOptions := &filesystem.CreateFileOptions{
		Umask: &umask,
		Owner: &user,
		Group: &group,
		ACL:   &acl,
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeNeverExpire,
		},
		LeaseDuration:   leaseDuration,
		ProposedLeaseID: proposedLeaseIDs[0],
	}
	resp, err := fsClient.CreateFile(context.Background(), testName, createFileOptions)
	_require.NoError(err)
	_require.NotNil(resp)

	fClient := fsClient.NewFileClient(testName)

	response, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal("4cf4e284-f6a8-4540-b53e-c3469af032dc", *response.Owner)
	_require.Equal("rwxr-xrwx", *response.Permissions)
	_require.Equal(filesystem.StateTypeLeased, *response.LeaseState)

}

func (s *RecordedTestSuite) TestCreateDirectoryInFileSystemSetOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	perms := "0777"
	umask := "0000"
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	leaseDuration := to.Ptr(int64(-1))

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	createDirOptions := &filesystem.CreateDirectoryOptions{
		Permissions:     &perms,
		Umask:           &umask,
		Owner:           &owner,
		Group:           &group,
		LeaseDuration:   leaseDuration,
		ProposedLeaseID: proposedLeaseIDs[0],
	}

	resp, err := fsClient.CreateDirectory(context.Background(), testName, createDirOptions)
	_require.NoError(err)
	_require.NotNil(resp)

	dirClient := fsClient.NewDirectoryClient(testName)

	response, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*response.Owner, "4cf4e284-f6a8-4540-b53e-c3469af032dc")
	_require.Equal("rwxrwxrwx", *response.Permissions)
	_require.Equal(filesystem.StateTypeLeased, *response.LeaseState)

}

func (s *RecordedTestSuite) TestFSCreateDefaultAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName

	options := &filesystem.ClientOptions{Audience: "https://storage.azure.com/"}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	fsClient, err := filesystem.NewClient(fsURL, cred, options)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

}

func (s *RecordedTestSuite) TestFSCreateCustomAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName

	options := &filesystem.ClientOptions{Audience: "https://" + accountName + ".blob.core.windows.net"}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	fsClient, err := filesystem.NewClient(fsURL, cred, options)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

}
