// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

//
//import (
//	"sort"
//
//	chk "gopkg.in/check.v1"
//)
//
//// tests general functionality
//func (s *azblobTestSuite) TestBlobListWrapper() {
//	bsu, err := getGenericBSU("", nil)
//
//	_assert.Nil(err)
//
//	container, _ := getContainerClient(bsu)
//
//	_, err = container.Create(ctx, nil)
//	_assert.Nil(err)
//	defer container.Delete(ctx, nil)
//
//	files := []string{"a", "b", "c"}
//
//	createNewBlobs(c, container, files)
//
//	pager := container.ListBlobsFlatSegment(nil)
//
//	found := make([]string, 0)
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
//			found = append(found, *blob.Name)
//		}
//	}
//	c.Assert(pager.Err(), chk.IsNil)
//
//	sort.Strings(files)
//	sort.Strings(found)
//
//	c.Assert(found, chk.DeepEquals, files)
//}
//
//// tests that the buffer filling isn't a problem
//func (s *azblobTestSuite) TestBlobListWrapperFullBuffer() {
//	bsu, err := getGenericBSU("", nil)
//
//	_assert.Nil(err)
//
//	container, _ := getContainerClient(bsu)
//
//	_, err = container.Create(ctx, nil)
//	_assert.Nil(err)
//	defer container.Delete(ctx, nil)
//
//	files := []string{"a", "b", "c"}
//
//	createNewBlobs(c, container, files)
//
//	pager := container.ListBlobsFlatSegment(nil)
//
//	found := make([]string, 0)
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
//			found = append(found, *blob.Name)
//		}
//	}
//	c.Assert(pager.Err(), chk.IsNil)
//
//	sort.Strings(files)
//	sort.Strings(found)
//
//	c.Assert(files, chk.DeepEquals, found)
//}
//
//func (s *azblobTestSuite) TestBlobListWrapperListingError() {
//	bsu, err := getGenericBSU("", nil)
//
//	_assert.Nil(err)
//
//	container, _ := getContainerClient(bsu)
//
//	pager := container.ListBlobsFlatSegment(nil)
//
//	c.Assert(pager.NextPage(ctx), chk.Equals, false)
//	c.Assert(pager.Err(), chk.NotNil)
//}
