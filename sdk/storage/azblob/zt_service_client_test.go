// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *azblobTestSuite) TestGetAccountInfo() {
	// // _assert := assert.New(s.T())
	// testName := s.T().Name()
	_context := getTestContext(s.T().Name())
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sAccInfo, err := svcClient.GetAccountInfo(context.Background())
	require.NoError(s.T(), err)
	require.NotEqualValues(s.T(), sAccInfo, ServiceGetAccountInfoResponse{})
}

//nolint
func (s *azblobUnrecordedTestSuite) TestServiceClientFromConnectionString() {
	// // _assert := assert.New(s.T())
	// testName := s.T().Name()

	accountName, _ := getAccountInfo(nil, testAccountDefault)
	connectionString := getConnectionString(nil, testAccountDefault)

	serviceURL, cred, err := parseConnectionString(connectionString)
	require.NoError(s.T(), err)
	require.Equal(s.T(), serviceURL, "https://"+accountName+".blob.core.windows.net/")

	svcClient, err := NewServiceClientWithSharedKey(serviceURL, cred, nil)
	require.NoError(s.T(), err)
	containerClient := createNewContainer(assert.New(s.T()), generateContainerName(s.T().Name()), svcClient)
	defer deleteContainer(assert.New(s.T()), containerClient)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestListContainersBasic() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	require.NoError(s.T(), err)
	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &CreateContainerOptions{Metadata: md})
	defer func(containerClient ContainerClient, ctx context.Context, options *DeleteContainerOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			require.NoError(s.T(), err)
		}
	}(containerClient, ctx, nil)
	require.NoError(s.T(), err)
	prefix := containerPrefix
	listOptions := ListContainersOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
	pager := svcClient.ListContainers(&listOptions)

	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			assert.NotNil(s.T(), container.Name)

			if *container.Name == containerName {
				assert.NotNil(s.T(), container.Properties)
				assert.NotNil(s.T(), container.Properties.LastModified)
				assert.NotNil(s.T(), container.Properties.Etag)
				assert.Equal(s.T(), *container.Properties.LeaseStatus, LeaseStatusTypeUnlocked)
				assert.Equal(s.T(), *container.Properties.LeaseState, LeaseStateTypeAvailable)
				assert.Nil(s.T(), container.Properties.LeaseDuration)
				assert.Nil(s.T(), container.Properties.PublicAccess)
				assert.NotNil(s.T(), container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				assert.EqualValues(s.T(), unwrappedMeta, md)
			}
		}
	}

	assert.Nil(s.T(), pager.Err())
	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), count, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestListContainersBasicUsingConnectionString() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	svcClient, err := getServiceClientFromConnectionString(nil, testAccountDefault, nil)
	require.NoError(s.T(), err)
	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &CreateContainerOptions{Metadata: md})
	defer func(containerClient ContainerClient, ctx context.Context, options *DeleteContainerOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			require.NoError(s.T(), err)
		}
	}(containerClient, ctx, nil)
	require.NoError(s.T(), err)
	prefix := containerPrefix
	listOptions := ListContainersOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
	pager := svcClient.ListContainers(&listOptions)

	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			assert.NotNil(s.T(), container.Name)

			if *container.Name == containerName {
				assert.NotNil(s.T(), container.Properties)
				assert.NotNil(s.T(), container.Properties.LastModified)
				assert.NotNil(s.T(), container.Properties.Etag)
				assert.Equal(s.T(), *container.Properties.LeaseStatus, LeaseStatusTypeUnlocked)
				assert.Equal(s.T(), *container.Properties.LeaseState, LeaseStateTypeAvailable)
				assert.Nil(s.T(), container.Properties.LeaseDuration)
				assert.Nil(s.T(), container.Properties.PublicAccess)
				assert.NotNil(s.T(), container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				assert.EqualValues(s.T(), unwrappedMeta, md)
			}
		}
	}

	assert.Nil(s.T(), pager.Err())
	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), count, 0)
}

//func (s *azblobTestSuite) TestListContainersPaged() {
//	// _assert := assert.New(s.T())
//	// testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
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
//		containerName := pagedContainersPrefix + generateContainerName(s.T().Name()) + string(i)
//		containerClient := createNewContainer(_assert, containerName, svcClient)
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

func (s *azblobTestSuite) TestAccountListContainersEmptyPrefix() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	_context := getTestContext(s.T().Name())
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient1 := createNewContainer(assert.New(s.T()), generateContainerName(s.T().Name())+"1", svcClient)
	defer deleteContainer(assert.New(s.T()), containerClient1)
	containerClient2 := createNewContainer(assert.New(s.T()), generateContainerName(s.T().Name())+"2", svcClient)
	defer deleteContainer(assert.New(s.T()), containerClient2)

	count := 0
	pager := svcClient.ListContainers(nil)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			count++
			assert.NotNil(s.T(), container.Name)
		}
	}
	assert.Nil(s.T(), pager.Err())
	assert.GreaterOrEqual(s.T(), count, 2)
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
////		require.NoError(s.T(), err)
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
////	require.NoError(s.T(), err)
////
////	// getting the next page should work
////	hasPage := pager.NextPage(context.Background())
////	_assert.(hasPage, chk.Equals, true)
////
////	page := pager.PageResponse()
////	require.NoError(s.T(), err)
////	_assert.(*page.EnumerationResults.ContainerItems, chk.HasLen, 2)
////	_assert.(*(*page.EnumerationResults.ContainerItems)[0].Name, chk.DeepEquals, containerName1)
////	_assert.(*(*page.EnumerationResults.ContainerItems)[1].Name, chk.DeepEquals, containerName2)
////}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicy() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	_context := getTestContext(s.T().Name())
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(s.T(), err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(s.T(), err)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	disabled := false
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &disabled}})
	require.NoError(s.T(), err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err = svcClient.GetProperties(ctx)
	require.NoError(s.T(), err)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, false)
	assert.Nil(s.T(), resp.StorageServiceProperties.DeleteRetentionPolicy.Days)
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyEmpty() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	_context := getTestContext(s.T().Name())
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(s.T(), err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(s.T(), err)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Empty retention policy causes an error, this is different from track 1.5
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{}})
	assert.Error(s.T(), err)
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyNil() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	_context := getTestContext(s.T().Name())
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(s.T(), err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(s.T(), err)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{})
	require.NoError(s.T(), err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = svcClient.GetProperties(ctx)
	require.NoError(s.T(), err)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	assert.EqualValues(s.T(), *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Disable for other tests
	enabled = to.BoolPtr(false)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled}})
	require.NoError(s.T(), err)
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyDaysTooSmall() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	_context := getTestContext(s.T().Name())
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := int32(0) // Minimum days is 1. Validated on the client.
	enabled := true
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	assert.Error(s.T(), err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestAccountDeleteRetentionPolicyDaysTooLarge() {
	// _assert := assert.New(s.T())
	var svcClient ServiceClient
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = getServiceClient(nil, testAccountDefault, nil)
		} else {
			svcClient, err = getServiceClientFromConnectionString(nil, testAccountDefault, nil)
		}
		require.NoError(s.T(), err)

		days := int32(366) // Max days is 365. Left to the service for validation.
		enabled := true
		_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
		assert.Error(s.T(), err)

		validateStorageError(assert.New(s.T()), err, StorageErrorCodeInvalidXMLDocument)
	}
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyDaysOmitted() {
	// _assert := assert.New(s.T())
	// testName := s.T().Name()
	_context := getTestContext(s.T().Name())
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	// Days is required if enabled is true.
	enabled := true
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled}})
	assert.Error(s.T(), err)

	validateStorageError(assert.New(s.T()), err, StorageErrorCodeInvalidXMLDocument)
}
