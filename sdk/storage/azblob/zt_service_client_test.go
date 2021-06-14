// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

//
//import (
//	"context"
//	"strings"
//	"time"
//
//	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
//
//	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
//)
//
//func (s *azblobTestSuite) TestGetAccountInfo() {
//	sa := getServiceClient(nil)
//
//	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
//	sAccInfo, err := sa.GetAccountInfo(context.Background())
//	_assert.Nil(err)
//	c.Assert(sAccInfo, chk.Not(chk.DeepEquals), ServiceGetAccountInfoResponse{})
//
//	//Test on a container
//	containerClient := sa.NewContainerClient(generateContainerName())
//	_, err = containerClient.Create(ctx, nil)
//	defer containerClient.Delete(ctx, nil)
//	_assert.Nil(err)
//	cAccInfo, err := containerClient.GetAccountInfo(ctx)
//	_assert.Nil(err)
//	c.Assert(cAccInfo, chk.Not(chk.DeepEquals), ContainerGetAccountInfoResponse{})
//
//	// test on a block blob URL. They all call the same thing on the base URL, so only one test is needed for that.
//	blobClient := containerClient.NewBlockBlobClient(generateBlobName())
//	_, err = blobClient.Upload(ctx, azcore.NopCloser(strings.NewReader("blah")), nil)
//	_assert.Nil(err)
//	bAccInfo, err := blobClient.GetAccountInfo(ctx)
//	_assert.Nil(err)
//	c.Assert(bAccInfo, chk.Not(chk.DeepEquals), BlobGetAccountInfoResponse{})
//}
//
//func (s *azblobTestSuite) TestListContainersBasic() {
//	sa, err := getGenericBSU("", nil)
//	_assert.Nil(err)
//
//	md := map[string]string{
//		"foo": "foovalue",
//		"bar": "barvalue",
//	}
//
//	container, name := getContainerClient(sa)
//	_, err = container.Create(ctx, &CreateContainerOptions{Metadata: &md})
//	defer container.Delete(ctx, nil)
//	_assert.Nil(err)
//
//	prefix := containerPrefix
//	listOptions := ListContainersSegmentOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
//	pager := sa.ListContainersSegment(&listOptions)
//
//	count := 0
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, container := range *resp.EnumerationResults.ContainerItems {
//			c.Assert(container.Name, chk.NotNil)
//
//			if *container.Name == name {
//				c.Assert(container.Properties, chk.NotNil)
//				c.Assert(container.Properties.LastModified, chk.NotNil)
//				c.Assert(container.Properties.Etag, chk.NotNil)
//				c.Assert(*container.Properties.LeaseStatus, chk.Equals, LeaseStatusUnlocked)
//				c.Assert(*container.Properties.LeaseState, chk.Equals, LeaseStateAvailable)
//				c.Assert(container.Properties.LeaseDuration, chk.IsNil)
//				c.Assert(container.Properties.PublicAccess, chk.IsNil)
//				c.Assert(container.Metadata, chk.NotNil)
//
//				unwrappedMeta := map[string]string{}
//				for k, v := range *container.Metadata {
//					if v != nil {
//						unwrappedMeta[k] = *v
//					}
//				}
//
//				c.Assert(unwrappedMeta, chk.DeepEquals, md)
//			}
//		}
//	}
//
//	c.Assert(pager.Err(), chk.IsNil)
//
//	// for container := range pager {
//	// 	count++
//	//
//	// 	c.Assert(container.Name, chk.NotNil)
//	//
//	// 	if *container.Name == name {
//	//
//	// 	}
//	// }
//	// err = <-errs
//
//	_assert.Nil(err)
//	c.Assert(count >= 0, chk.Equals, true)
//}
//
//func (s *azblobTestSuite) TestListContainersPaged() {
//	sa := getServiceClient(nil)
//
//	const numContainers = 6
//	maxResults := int32(2)
//	const pagedContainersPrefix = "azcontainerpaged"
//
//	containers := make([]ContainerClient, numContainers)
//	expectedResults := make(map[string]bool)
//	for i := 0; i < numContainers; i++ {
//		containerClient, containerName := createNewContainerWithSuffix(c, sa, pagedContainersPrefix)
//		containers[i] = containerClient
//		expectedResults[containerName] = false
//	}
//
//	defer func() {
//		for i := range containers {
//			deleteContainer(containers[i])
//		}
//	}()
//
//	// list for a first time
//	prefix := containerPrefix + pagedContainersPrefix
//	listOptions := ListContainersSegmentOptions{MaxResults: &maxResults, Prefix: &prefix}
//	count := 0
//	results := make([]ContainerItem, 0)
//
//	pager := sa.ListContainersSegment(&listOptions)
//
//	for pager.NextPage(ctx) {
//		for _, container := range *pager.PageResponse().EnumerationResults.ContainerItems {
//			if container == nil {
//				continue
//			}
//
//			results = append(results, *container)
//			count += 1
//			c.Assert(container.Name, chk.NotNil)
//		}
//	}
//
//	c.Assert(pager.Err(), chk.IsNil)
//	c.Assert(count, chk.Equals, numContainers)
//	c.Assert(len(results), chk.Equals, numContainers)
//
//	// make sure each container we see is expected
//	for _, container := range results {
//		_, ok := expectedResults[*container.Name]
//		c.Assert(ok, chk.Equals, true)
//
//		expectedResults[*container.Name] = true
//	}
//
//	// make sure every expected container was seen
//	for _, seen := range expectedResults {
//		c.Assert(seen, chk.Equals, true)
//	}
//}
//
//func (s *azblobTestSuite) TestAccountListContainersEmptyPrefix() {
//	bsu := getServiceClient(nil)
//	containerClient1, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient1)
//	containerClient2, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient2)
//
//	count := 0
//	pager := bsu.ListContainersSegment(nil)
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, container := range *resp.EnumerationResults.ContainerItems {
//			count++
//			c.Assert(container.Name, chk.NotNil)
//		}
//	}
//	c.Assert(pager.Err(), chk.IsNil)
//
//	c.Assert(count >= 2, chk.Equals, true)
//}
//
//// TODO re-enable after fixing error handling
////func (s *azblobTestSuite) TestAccountListContainersMaxResultsNegative() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////
////	illegalMaxResults := []int32{-2, 0}
////	for _, num := range illegalMaxResults {
////		options := ListContainersSegmentOptions{MaxResults: &num}
////
////		// getting the pager should still work
////		pager, err := bsu.ListContainersSegment(context.Background(), 100, time.Hour, &options)
////		_assert.Nil(err)
////
////		// getting the next page should fail
////
////	}
////}
//
////func (s *azblobTestSuite) TestAccountListContainersMaxResultsExact() {
////	// If this test fails, ensure there are no extra containers prefixed with go in the account. These may be left over if a test is interrupted.
////	bsu := getServiceClient()
////	containerClient1, containerName1 := createNewContainerWithSuffix(c, bsu, "abc")
////	defer deleteContainer(containerClient1)
////	containerClient2, containerName2 := createNewContainerWithSuffix(c, bsu, "abcde")
////	defer deleteContainer(containerClient2)
////
////	prefix := containerPrefix + "abc"
////	maxResults := int32(2)
////	options := ListContainersSegmentOptions{Prefix: &prefix, MaxResults: &maxResults}
////	pager, err := bsu.ListContainersSegment(&options)
////	_assert.Nil(err)
////
////	// getting the next page should work
////	hasPage := pager.NextPage(context.Background())
////	c.Assert(hasPage, chk.Equals, true)
////
////	page := pager.PageResponse()
////	_assert.Nil(err)
////	c.Assert(*page.EnumerationResults.ContainerItems, chk.HasLen, 2)
////	c.Assert(*(*page.EnumerationResults.ContainerItems)[0].Name, chk.DeepEquals, containerName1)
////	c.Assert(*(*page.EnumerationResults.ContainerItems)[1].Name, chk.DeepEquals, containerName2)
////}
//
//func (s *azblobTestSuite) TestAccountDeleteRetentionPolicy() {
//	bsu := getServiceClient(nil)
//
//	days := int32(5)
//	enabled := true
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
//	_assert.Nil(err)
//
//	// From FE, 30 seconds is guaranteed to be enough.
//	time.Sleep(time.Second * 30)
//
//	resp, err := bsu.GetProperties(ctx)
//	_assert.Nil(err)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))
//
//	disabled := false
//	_, err = bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &disabled}})
//	_assert.Nil(err)
//
//	// From FE, 30 seconds is guaranteed to be enough.
//	time.Sleep(time.Second * 30)
//
//	resp, err = bsu.GetProperties(ctx)
//	_assert.Nil(err)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, false)
//	c.Assert(resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyEmpty() {
//	bsu := getServiceClient(nil)
//
//	days := int32(5)
//	enabled := true
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
//	_assert.Nil(err)
//
//	// From FE, 30 seconds is guaranteed to be enough.
//	time.Sleep(time.Second * 30)
//
//	resp, err := bsu.GetProperties(ctx)
//	_assert.Nil(err)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))
//
//	// Empty retention policy causes an error, this is different from track 1.5
//	_, err = bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{}})
//	_assert.NotNil(err)
//}
//
//func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyNil() {
//	bsu := getServiceClient(nil)
//
//	days := int32(5)
//	enabled := true
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
//	_assert.Nil(err)
//
//	// From FE, 30 seconds is guaranteed to be enough.
//	time.Sleep(time.Second * 30)
//
//	resp, err := bsu.GetProperties(ctx)
//	_assert.Nil(err)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))
//
//	_, err = bsu.SetProperties(ctx, StorageServiceProperties{})
//	_assert.Nil(err)
//
//	// From FE, 30 seconds is guaranteed to be enough.
//	time.Sleep(time.Second * 30)
//
//	// If an element of service properties is not passed, the service keeps the current settings.
//	resp, err = bsu.GetProperties(ctx)
//	_assert.Nil(err)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
//	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))
//
//	// Disable for other tests
//	enabled = false
//	bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled}})
//}
//
//func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyDaysTooSmall() {
//	bsu := getServiceClient(nil)
//
//	days := int32(0) // Minimum days is 1. Validated on the client.
//	enabled := true
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
//	_assert.NotNil(err)
//}
//
//func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyDaysTooLarge() {
//	bsu := getServiceClient(nil)
//
//	days := int32(366) // Max days is 365. Left to the service for validation.
//	enabled := true
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidXMLDocument)
//}
//
//func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyDaysOmitted() {
//	bsu := getServiceClient(nil)
//
//	// Days is required if enabled is true.
//	enabled := true
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled}})
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidXMLDocument)
//}
