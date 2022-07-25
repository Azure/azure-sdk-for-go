//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//
//import (
//	"github.com/stretchr/testify/require"
//	"sort"
//)
//
//// tests general functionality
//func (s *azblobTestSuite) TestBlobListWrapper() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	_, err = containerClient.Create(ctx, nil)
//	_require.Nil(err)
//	defer deleteContainer(_require, containerClient)
//
//	files := []string{"a123", "b234", "c345"}
//
//	createNewBlobs(_require, files, containerClient)
//
//	pager := containerClient.NewListBlobsFlatPager(nil)
//
//	found := make([]string, 0)
//
//	for pager.More() {
//		resp, err := pager.NextPage(ctx)
//		_require.Nil(err)
//
//		for _, blob := range resp.Segment.BlobItems {
//			found = append(found, *blob.Name)
//		}
//
//		if err != nil {
//			break
//		}
//	}
//
//	sort.Strings(files)
//	sort.Strings(found)
//
//	_require.EqualValues(found, files)
//}
//
//// tests that the buffer filling isn't a problem
//func (s *azblobTestSuite) TestBlobListWrapperFullBuffer() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerClient := getContainerClient(generateContainerName(testName), svcClient)
//
//	_, err = containerClient.Create(ctx, nil)
//	_require.Nil(err)
//	defer deleteContainer(_require, containerClient)
//
//	files := []string{"a123", "b234", "c345"}
//
//	createNewBlobs(_require, files, containerClient)
//
//	pager := containerClient.NewListBlobsFlatPager(nil)
//
//	found := make([]string, 0)
//
//	for pager.More() {
//		resp, err := pager.NextPage(ctx)
//		_require.Nil(err)
//
//		for _, blob := range resp.Segment.BlobItems {
//			found = append(found, *blob.Name)
//		}
//
//		if err != nil {
//			break
//		}
//	}
//
//	sort.Strings(files)
//	sort.Strings(found)
//
//	_require.EqualValues(files, found)
//}
//
//func (s *azblobTestSuite) TestBlobListWrapperListingError() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerClient := getContainerClient(generateContainerName(testName), svcClient)
//
//	pager := containerClient.NewListBlobsFlatPager(nil)
//	for pager.More() {
//		_, err := pager.NextPage(ctx)
//		_require.NotNil(err)
//		if err != nil {
//			break
//		}
//	}
//}
