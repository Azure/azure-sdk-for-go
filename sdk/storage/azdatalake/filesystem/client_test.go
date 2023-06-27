//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
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

func (s *RecordedTestSuite) TestCreateFilesystem() {
	_require := require.New(s.T())
	testName := s.T().Name()
	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

}
