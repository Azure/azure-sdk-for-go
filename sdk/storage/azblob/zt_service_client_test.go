//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

func (s *azblobTestSuite) TestGetAccountInfo() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sAccInfo, err := svcClient.GetAccountInfo(context.Background(), nil)
	_require.Nil(err)
	_require.NotEqualValues(sAccInfo, service.GetAccountInfoResponse{})
}

//nolint
func (s *azblobUnrecordedTestSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := getAccountInfo(nil, testAccountDefault)
	connectionString := getConnectionString(nil, testAccountDefault)

	parsedConnStr, err := azblob.ParseConnectionString(connectionString)
	_require.Nil(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".blob.core.windows.net/")

	sharedKeyCred, err := azblob.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.Nil(err)

	svcClient, err := service.NewClientWithSharedKey(parsedConnStr.ServiceURL, sharedKeyCred, nil)
	_require.Nil(err)
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestListContainersBasic() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	_require.Nil(err)
	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &container.CreateOptions{Metadata: md})
	defer func(containerClient *container.Client, ctx context.Context, options *container.DeleteOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			_require.Nil(err)
		}
	}(containerClient, ctx, nil)
	_require.Nil(err)
	prefix := containerPrefix
	listOptions := service.ListContainersOptions{Prefix: &prefix, Include: service.ListContainersDetail{Metadata: true}}
	pager := svcClient.NewListContainersPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		for _, container := range resp.ContainerItems {
			_require.NotNil(container.Name)

			if *container.Name == containerName {
				_require.NotNil(container.Properties)
				_require.NotNil(container.Properties.LastModified)
				_require.NotNil(container.Properties.Etag)
				_require.Equal(*container.Properties.LeaseStatus, azblob.LeaseStatusTypeUnlocked)
				_require.Equal(*container.Properties.LeaseState, azblob.LeaseStateTypeAvailable)
				_require.Nil(container.Properties.LeaseDuration)
				_require.Nil(container.Properties.PublicAccess)
				_require.NotNil(container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				_require.EqualValues(unwrappedMeta, md)
			}
		}
		if err != nil {
			break
		}
	}

	_require.Nil(err)
	_require.GreaterOrEqual(count, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestListContainersBasicUsingConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClientFromConnectionString(nil, testAccountDefault, nil)
	_require.Nil(err)
	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &container.CreateOptions{Metadata: md})
	defer func(containerClient *container.Client, ctx context.Context, options *container.DeleteOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			_require.Nil(err)
		}
	}(containerClient, ctx, nil)
	_require.Nil(err)
	prefix := containerPrefix
	listOptions := service.ListContainersOptions{Prefix: &prefix, Include: service.ListContainersDetail{Metadata: true}}
	pager := svcClient.NewListContainersPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)

		for _, container := range resp.ContainerItems {
			_require.NotNil(container.Name)

			if *container.Name == containerName {
				_require.NotNil(container.Properties)
				_require.NotNil(container.Properties.LastModified)
				_require.NotNil(container.Properties.Etag)
				_require.Equal(*container.Properties.LeaseStatus, azblob.LeaseStatusTypeUnlocked)
				_require.Equal(*container.Properties.LeaseState, azblob.LeaseStateTypeAvailable)
				_require.Nil(container.Properties.LeaseDuration)
				_require.Nil(container.Properties.PublicAccess)
				_require.NotNil(container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				_require.EqualValues(unwrappedMeta, md)
			}
		}
		if err != nil {
			break
		}
	}

	_require.Nil(err)
	_require.GreaterOrEqual(count, 0)
}

//func (s *azblobTestSuite) TestListContainersPaged() {
//	_require := require.New(s.T())
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
//	containers := make([]containerClient, numContainers)
//	expectedResults := make(map[string]bool)
//	for i := 0; i < numContainers; i++ {
//		containerName := pagedContainersPrefix + generateContainerName(testName) + string(i)
//		containerClient := createNewContainer(_require, containerName, svcClient)
//		containers[i] = containerClient
//		expectedResults[containerName] = false
//	}
//
//	defer func() {
//		for i := range containers {
//			deleteContainer(_require, containers[i])
//		}
//	}()
//
//	// list for a first time
//	prefix := containerPrefix + pagedContainersPrefix
//	listOptions := ServiceListContainersOptions{MaxResults: &maxResults, Prefix: &prefix}
//	count := 0
//	results := make([]ContainerItem, 0)
//
//	pager := sa.NewListContainersPager(&listOptions)
//
//	for pager.NextPage(ctx) {
//		for _, container := range *pager.PageResponse().EnumerationResults.ContainerItems {
//			if container == nil {
//				continue
//			}
//
//			results = append(results, *container)
//			count += 1
//			_require.(container.Name, chk.NotNil)
//		}
//	}
//
//	_require.(pager.Err(), chk.IsNil)
//	_require.(count, chk.Equals, numContainers)
//	_require.(len(results), chk.Equals, numContainers)
//
//	// make sure each container we see is expected
//	for _, container := range results {
//		_, ok := expectedResults[*container.Name]
//		_require.(ok, chk.Equals, true)
//
//		expectedResults[*container.Name] = true
//	}
//
//	// make sure every expected container was seen
//	for _, seen := range expectedResults {
//		_require.(seen, chk.Equals, true)
//	}
//}

func (s *azblobTestSuite) TestAccountListContainersEmptyPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient1 := createNewContainer(_require, generateContainerName(testName)+"1", svcClient)
	defer deleteContainer(_require, containerClient1)
	containerClient2 := createNewContainer(_require, generateContainerName(testName)+"2", svcClient)
	defer deleteContainer(_require, containerClient2)

	count := 0
	pager := svcClient.NewListContainersPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)

		for _, container := range resp.ContainerItems {
			count++
			_require.NotNil(container.Name)
		}
		if err != nil {
			break
		}
	}
	_require.GreaterOrEqual(count, 2)
}

//// TODO re-enable after fixing error handling
////func (s *azblobTestSuite) TestAccountListContainersMaxResultsNegative() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////
////	illegalMaxResults := []int32{-2, 0}
////	for _, num := range illegalMaxResults {
////		options := ServiceListContainersOptions{MaxResults: &num}
////
////		// getting the pager should still work
////		pager, err := svcClient.NewListContainersPager(context.Background(), 100, time.Hour, &options)
////		_require.Nil(err)
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
////	options := ServiceListContainersOptions{Prefix: &prefix, MaxResults: &maxResults}
////	pager, err := svcClient.NewListContainersPager(&options)
////	_require.Nil(err)
////
////	// getting the next page should work
////	hasPage := pager.NextPage(context.Background())
////	_require.(hasPage, chk.Equals, true)
////
////	page := pager.PageResponse()
////	_require.Nil(err)
////	_require.(*page.EnumerationResults.ContainerItems, chk.HasLen, 2)
////	_require.(*(*page.EnumerationResults.ContainerItems)[0].Name, chk.DeepEquals, containerName1)
////	_require.(*(*page.EnumerationResults.ContainerItems)[1].Name, chk.DeepEquals, containerName2)
////}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	disabled := false
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &disabled}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err = svcClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, false)
	_require.Nil(resp.StorageServiceProperties.DeleteRetentionPolicy.Days)
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Empty retention policy causes an error, this is different from track 1.5
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{}})
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = svcClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Disable for other tests
	enabled = to.Ptr(false)
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled}})
	_require.Nil(err)
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyDaysTooSmall() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	days := int32(0) // Minimum days is 1. Validated on the client.
	enabled := true
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled, Days: &days}})
	_require.NotNil(err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestAccountDeleteRetentionPolicyDaysTooLarge() {
	_require := require.New(s.T())
	var svcClient *service.Client
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = getServiceClient(nil, testAccountDefault, nil)
		} else {
			svcClient, err = getServiceClientFromConnectionString(nil, testAccountDefault, nil)
		}
		_require.Nil(err)

		days := int32(366) // Max days is 365. Left to the service for validation.
		enabled := true
		_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled, Days: &days}})
		_require.NotNil(err)

		validateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
	}
}

func (s *azblobTestSuite) TestAccountDeleteRetentionPolicyDaysOmitted() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	// Days is required if enabled is true.
	enabled := true
	_, err = svcClient.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled}})
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
}
