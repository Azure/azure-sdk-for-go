package azblob

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

func (s *aztestsSuite) TestGetAccountInfo(c *chk.C) {
	sa := getBSU()

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sAccInfo, err := sa.GetAccountInfo(context.Background())
	c.Assert(err, chk.IsNil)
	c.Assert(sAccInfo, chk.Not(chk.DeepEquals), ServiceGetAccountInfoResponse{})

	//Test on a container
	containerClient := sa.NewContainerClient(generateContainerName())
	_, err = containerClient.Create(ctx, nil)
	defer containerClient.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)
	cAccInfo, err := containerClient.GetAccountInfo(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(cAccInfo, chk.Not(chk.DeepEquals), ContainerGetAccountInfoResponse{})

	// test on a block blob URL. They all call the same thing on the base URL, so only one test is needed for that.
	blobClient := containerClient.NewBlockBlobClient(generateBlobName())
	_, err = blobClient.Upload(ctx, azcore.NopCloser(strings.NewReader("blah")), nil)
	c.Assert(err, chk.IsNil)
	bAccInfo, err := blobClient.GetAccountInfo(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(bAccInfo, chk.Not(chk.DeepEquals), BlobGetAccountInfoResponse{})
}

func (s *aztestsSuite) TestListContainersBasic(c *chk.C) {
	sa, err := getGenericBSU("")
	c.Assert(err, chk.IsNil)

	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	container, name := getContainerClient(c, sa)
	_, err = container.Create(ctx, &CreateContainerOptions{Metadata: &md})
	defer container.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)

	prefix := containerPrefix
	listOptions := ListContainersSegmentOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
	pager, errs := sa.ListContainersSegment(context.Background(), 100, 0, &listOptions)

	count := 0
	for container := range pager {
		count++

		c.Assert(container.Name, chk.NotNil)

		if *container.Name == name {
			c.Assert(container.Properties, chk.NotNil)
			c.Assert(container.Properties.LastModified, chk.NotNil)
			c.Assert(container.Properties.Etag, chk.NotNil)
			c.Assert(*container.Properties.LeaseStatus, chk.Equals, LeaseStatusTypeUnlocked)
			c.Assert(*container.Properties.LeaseState, chk.Equals, LeaseStateTypeAvailable)
			c.Assert(container.Properties.LeaseDuration, chk.IsNil)
			c.Assert(container.Properties.PublicAccess, chk.IsNil)
			c.Assert(container.Metadata, chk.DeepEquals, &md)
		}
	}
	err = <-errs

	c.Assert(err, chk.IsNil)
	c.Assert(count >= 0, chk.Equals, true)
}

func (s *aztestsSuite) TestListContainersPaged(c *chk.C) {
	sa := getBSU()

	const numContainers = 6
	maxResults := int32(2)
	const pagedContainersPrefix = "azcontainerpaged"

	containers := make([]ContainerClient, numContainers)
	expectedResults := make(map[string]bool)
	for i := 0; i < numContainers; i++ {
		containerClient, containerName := createNewContainerWithSuffix(c, sa, pagedContainersPrefix)
		containers[i] = containerClient
		expectedResults[containerName] = false
	}

	defer func() {
		for i := range containers {
			deleteContainer(c, containers[i])
		}
	}()

	// list for a first time
	prefix := containerPrefix + pagedContainersPrefix
	listOptions := ListContainersSegmentOptions{MaxResults: &maxResults, Prefix: &prefix}
	count := 0
	results := make([]ContainerItem, 0)

	pager, errs := sa.ListContainersSegment(context.Background(), 100, time.Hour, &listOptions)
	for container := range pager {
		// record the results
		results = append(results, container)
		count += 1
		c.Assert(container.Name, chk.NotNil)
	}
	err := <-errs
	c.Assert(err, chk.IsNil)
	c.Assert(count, chk.Equals, numContainers)
	c.Assert(len(results), chk.Equals, numContainers)

	// make sure each container we see is expected
	for _, container := range results {
		_, ok := expectedResults[*container.Name]
		c.Assert(ok, chk.Equals, true)

		expectedResults[*container.Name] = true
	}

	// make sure every expected container was seen
	for _, seen := range expectedResults {
		c.Assert(seen, chk.Equals, true)
	}
}

func (s *aztestsSuite) TestAccountListContainersEmptyPrefix(c *chk.C) {
	bsu := getBSU()
	containerURL1, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL1)
	containerURL2, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL2)

	count := 0
	pager, errChan := bsu.ListContainersSegment(context.Background(), 100, time.Hour, nil)
	c.Assert(<-errChan, chk.IsNil)

	for container := range pager {
		// record the results
		count += 1
		c.Assert(container.Name, chk.NotNil)
	}
	c.Assert(count >= 2, chk.Equals, true)
}

// TODO re-enable after fixing error handling
//func (s *aztestsSuite) TestAccountListContainersMaxResultsNegative(c *chk.C) {
//	bsu := getBSU()
//	containerURL, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerURL)
//
//	illegalMaxResults := []int32{-2, 0}
//	for _, num := range illegalMaxResults {
//		options := ListContainersSegmentOptions{MaxResults: &num}
//
//		// getting the pager should still work
//		pager, err := bsu.ListContainersSegment(context.Background(), 100, time.Hour, &options)
//		c.Assert(err, chk.IsNil)
//
//		// getting the next page should fail
//
//	}
//}

//func (s *aztestsSuite) TestAccountListContainersMaxResultsExact(c *chk.C) {
//	// If this test fails, ensure there are no extra containers prefixed with go in the account. These may be left over if a test is interrupted.
//	bsu := getBSU()
//	containerURL1, containerName1 := createNewContainerWithSuffix(c, bsu, "abc")
//	defer deleteContainer(c, containerURL1)
//	containerURL2, containerName2 := createNewContainerWithSuffix(c, bsu, "abcde")
//	defer deleteContainer(c, containerURL2)
//
//	prefix := containerPrefix + "abc"
//	maxResults := int32(2)
//	options := ListContainersSegmentOptions{Prefix: &prefix, MaxResults: &maxResults}
//	pager, err := bsu.ListContainersSegment(&options)
//	c.Assert(err, chk.IsNil)
//
//	// getting the next page should work
//	hasPage := pager.NextPage(context.Background())
//	c.Assert(hasPage, chk.Equals, true)
//
//	page := pager.PageResponse()
//	c.Assert(err, chk.IsNil)
//	c.Assert(*page.EnumerationResults.ContainerItems, chk.HasLen, 2)
//	c.Assert(*(*page.EnumerationResults.ContainerItems)[0].Name, chk.DeepEquals, containerName1)
//	c.Assert(*(*page.EnumerationResults.ContainerItems)[1].Name, chk.DeepEquals, containerName2)
//}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicy(c *chk.C) {
	bsu := getBSU()

	days := int32(5)
	enabled := true
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))

	disabled := false
	_, err = bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &disabled}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err = bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, false)
	c.Assert(resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.IsNil)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyEmpty(c *chk.C) {
	bsu := getBSU()

	days := int32(5)
	enabled := true
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))

	// Empty retention policy causes an error, this is different from track 1.5
	_, err = bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{}})
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyNil(c *chk.C) {
	bsu := getBSU()

	days := int32(5)
	enabled := true
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))

	_, err = bsu.SetProperties(ctx, StorageServiceProperties{})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, chk.DeepEquals, true)
	c.Assert(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, chk.DeepEquals, int32(5))

	// Disable for other tests
	enabled = false
	bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled}})
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyDaysTooSmall(c *chk.C) {
	bsu := getBSU()

	days := int32(0) // Minimum days is 1. Validated on the client.
	enabled := true
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyDaysTooLarge(c *chk.C) {
	bsu := getBSU()

	days := int32(366) // Max days is 365. Left to the service for validation.
	enabled := true
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	c.Assert(err, chk.NotNil)

	// TODO the error should have more details, follow up with Joel
	//validateStorageError(c, err, ServiceCodeInvalidXMLDocument)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyDaysOmitted(c *chk.C) {
	bsu := getBSU()

	// Days is required if enabled is true.
	enabled := true
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled}})
	c.Assert(err, chk.NotNil)

	// TODO the error should have more details, follow up with Joel
	//validateStorageError(c, err, ServiceCodeInvalidXMLDocument)
}
