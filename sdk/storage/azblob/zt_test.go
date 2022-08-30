//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/suite"
)

var ctx = context.Background()

type AZBlobRecordedTestsSuite struct {
	suite.Suite
}

type AZBlobUnrecordedTestsSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &AZBlobRecordedTestsSuite{})
	//suite.Run(t, &AZBlobUnrecordedTestsSuite{})
}

// nolint
func (s *AZBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

// nolint
func (s *AZBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

// nolint
func (s *AZBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}
