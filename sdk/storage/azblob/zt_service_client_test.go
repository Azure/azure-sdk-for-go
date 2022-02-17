// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccountInfo(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sAccInfo, err := svcClient.GetAccountInfo(context.Background())
	require.NoError(t, err)
	require.NotEqualValues(t, sAccInfo, ServiceGetAccountInfoResponse{})
}

func TestServiceClientFromConnectionString(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()

	svcClient, err := createServiceClientFromConnectionString(t, testAccountDefault)
	require.NoError(t, err)
	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)
}

func TestListContainersBasic(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &CreateContainerOptions{Metadata: md})
	defer func(containerClient ContainerClient, ctx context.Context, options *DeleteContainerOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			require.NoError(t, err)
		}
	}(containerClient, ctx, nil)
	require.NoError(t, err)
	prefix := containerPrefix
	listOptions := ListContainersOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
	pager := svcClient.ListContainers(&listOptions)

	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			require.NotNil(t, container.Name)

			if *container.Name == containerName {
				require.NotNil(t, container.Properties)
				require.NotNil(t, container.Properties.LastModified)
				require.NotNil(t, container.Properties.Etag)
				require.Equal(t, *container.Properties.LeaseStatus, LeaseStatusTypeUnlocked)
				require.Equal(t, *container.Properties.LeaseState, LeaseStateTypeAvailable)
				require.Nil(t, container.Properties.LeaseDuration)
				require.Nil(t, container.Properties.PublicAccess)
				require.NotNil(t, container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				require.EqualValues(t, unwrappedMeta, md)
			}
		}
	}

	require.Nil(t, pager.Err())
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, 0)
}

func TestListContainersBasicUsingConnectionString(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClientFromConnectionString(t, testAccountDefault)
	require.NoError(t, err)

	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &CreateContainerOptions{Metadata: md})
	defer func(containerClient ContainerClient, ctx context.Context, options *DeleteContainerOptions) {
		_, err := containerClient.Delete(ctx, options)
		require.NoError(t, err)
	}(containerClient, ctx, nil)
	require.NoError(t, err)
	prefix := containerPrefix
	listOptions := ListContainersOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
	pager := svcClient.ListContainers(&listOptions)

	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			require.NotNil(t, container.Name)

			if *container.Name == containerName {
				require.NotNil(t, container.Properties)
				require.NotNil(t, container.Properties.LastModified)
				require.NotNil(t, container.Properties.Etag)
				require.Equal(t, *container.Properties.LeaseStatus, LeaseStatusTypeUnlocked)
				require.Equal(t, *container.Properties.LeaseState, LeaseStateTypeAvailable)
				require.Nil(t, container.Properties.LeaseDuration)
				require.Nil(t, container.Properties.PublicAccess)
				require.NotNil(t, container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				require.EqualValues(t, unwrappedMeta, md)
			}
		}
	}

	require.NoError(t, pager.Err())
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, 0)
}

//func (s *azblobTestSuite) TestListContainersPaged() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	const numContainers = 6
//	maxResults := int32(2)
//	const pagedContainersPrefix = "azcontainerpaged"
//
//	containers := make([]ContainerClient, numContainers)
//	expectedResults := make(map[string]bool)
//	for i := 0; i < numContainers; i++ {
//		containerName := pagedContainersPrefix + generateContainerName(testName) + string(i)
//		containerClient := createNewContainer(t, containerName, svcClient)
//		containers[i] = containerClient
//		expectedResults[containerName] = false
//	}
//
//	defer func() {
//		for i := range containers {
//			deleteContainer(_assert, containers[i])
//		}
//	}()
//
//	// list for a first time
//	prefix := containerPrefix + pagedContainersPrefix
//	listOptions := ListContainersOptions{MaxResults: &maxResults, Prefix: &prefix}
//	count := 0
//	results := make([]ContainerItem, 0)
//
//	pager := sa.ListContainers(&listOptions)
//
//	for pager.NextPage(ctx) {
//		for _, container := range *pager.PageResponse().EnumerationResults.ContainerItems {
//			if container == nil {
//				continue
//			}
//
//			results = append(results, *container)
//			count += 1
//			_assert.(container.Name, chk.NotNil)
//		}
//	}
//
//	_assert.(pager.Err(), chk.IsNil)
//	_assert.(count, chk.Equals, numContainers)
//	_assert.(len(results), chk.Equals, numContainers)
//
//	// make sure each container we see is expected
//	for _, container := range results {
//		_, ok := expectedResults[*container.Name]
//		_assert.(ok, chk.Equals, true)
//
//		expectedResults[*container.Name] = true
//	}
//
//	// make sure every expected container was seen
//	for _, seen := range expectedResults {
//		_assert.(seen, chk.Equals, true)
//	}
//}

func TestAccountListContainersEmptyPrefix(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient1 := createNewContainer(t, generateContainerName(testName)+"1", svcClient)
	defer deleteContainer(_assert, containerClient1)
	containerClient2 := createNewContainer(t, generateContainerName(testName)+"2", svcClient)
	defer deleteContainer(_assert, containerClient2)

	count := 0
	pager := svcClient.ListContainers(nil)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			count++
			require.NotNil(t, container.Name)
		}
	}
	require.NoError(t, pager.Err())
	require.GreaterOrEqual(t, count, 2)
}

//// TODO re-enable after fixing error handling
////func (s *azblobTestSuite) TestAccountListContainersMaxResultsNegative() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////
////	illegalMaxResults := []int32{-2, 0}
////	for _, num := range illegalMaxResults {
////		options := ListContainersOptions{MaxResults: &num}
////
////		// getting the pager should still work
////		pager, err := svcClient.ListContainers(context.Background(), 100, time.Hour, &options)
////		_assert.NoError(err)
////
////		// getting the next page should fail
////
////	}
////}
//
////func (s *azblobTestSuite) TestAccountListContainersMaxResultsExact() {
////	// If this test fails, ensure there are no extra containers prefixed with go in the account. These may be left over if a test is interrupted.
////	svcClient := getServiceClient()
////	containerClient1, containerName1 := createNewContainerWithSuffix(c, svcClient, "abc")
////	defer deleteContainer(containerClient1)
////	containerClient2, containerName2 := createNewContainerWithSuffix(c, svcClient, "abcde")
////	defer deleteContainer(containerClient2)
////
////	prefix := containerPrefix + "abc"
////	maxResults := int32(2)
////	options := ListContainersOptions{Prefix: &prefix, MaxResults: &maxResults}
////	pager, err := svcClient.ListContainers(&options)
////	_assert.NoError(err)
////
////	// getting the next page should work
////	hasPage := pager.NextPage(context.Background())
////	_assert.(hasPage, chk.Equals, true)
////
////	page := pager.PageResponse()
////	_assert.NoError(err)
////	_assert.(*page.EnumerationResults.ContainerItems, chk.HasLen, 2)
////	_assert.(*(*page.EnumerationResults.ContainerItems)[0].Name, chk.DeepEquals, containerName1)
////	_assert.(*(*page.EnumerationResults.ContainerItems)[1].Name, chk.DeepEquals, containerName2)
////}

func TestAccountDeleteRetentionPolicy(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	disabled := false
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &disabled}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err = svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, false)
	require.Nil(t, resp.StorageServiceProperties.DeleteRetentionPolicy.Days)
}

func TestAccountDeleteRetentionPolicyEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Empty retention policy causes an error, this is different from track 1.5
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{}})
	require.Error(t, err)
}

func TestAccountDeleteRetentionPolicyNil(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Disable for other tests
	enabled = to.BoolPtr(false)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled}})
	require.NoError(t, err)
}

func TestAccountDeleteRetentionPolicyDaysTooSmall(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := int32(0) // Minimum days is 1. Validated on the client.
	enabled := true
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	require.Error(t, err)
}

func TestAccountDeleteRetentionPolicyDaysTooLarge(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)

	var svcClient ServiceClient
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = createServiceClient(t, testAccountDefault)
		} else {
			svcClient, err = createServiceClientFromConnectionString(t, testAccountDefault)
		}
		require.NoError(t, err)

		days := int32(366) // Max days is 365. Left to the service for validation.
		enabled := true
		_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
		require.Error(t, err)

		validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
	}
}

func TestAccountDeleteRetentionPolicyDaysOmitted(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	// Days is required if enabled is true.
	enabled := true
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled}})
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}
