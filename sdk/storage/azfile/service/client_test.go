//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running service Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &ServiceRecordedTestsSuite{})
		suite.Run(t, &ServiceUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &ServiceRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &ServiceRecordedTestsSuite{})
	}
}

func (s *ServiceRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *ServiceRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *ServiceUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *ServiceUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type ServiceRecordedTestsSuite struct {
	suite.Suite
}

type ServiceUnrecordedTestsSuite struct {
	suite.Suite
}
