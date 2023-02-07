//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

func (s *ServiceRecordedTestsSuite) TestGetAccountInfo() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sAccInfo, err := svcClient.GetAccountInfo(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(sAccInfo)
}

func (s *ServiceUnrecordedTestsSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)
	connectionString := testcommon.GetConnectionString(testcommon.TestAccountDefault)

	parsedConnStr, err := shared.ParseConnectionString(connectionString)
	_require.Nil(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".blob.core.windows.net/")

	sharedKeyCred, err := azblob.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.Nil(err)

	svcClient, err := service.NewClientWithSharedKeyCredential(parsedConnStr.ServiceURL, sharedKeyCred, nil)
	_require.Nil(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
}

func (s *ServiceUnrecordedTestsSuite) TestListContainersBasic() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)
	_, err = containerClient.Create(context.Background(), &container.CreateOptions{Metadata: md})
	defer func(containerClient *container.Client, ctx context.Context, options *container.DeleteOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			_require.Nil(err)
		}
	}(containerClient, context.Background(), nil)
	_require.Nil(err)
	prefix := testcommon.ContainerPrefix
	listOptions := service.ListContainersOptions{Prefix: &prefix, Include: service.ListContainersInclude{Metadata: true}}
	pager := svcClient.NewListContainersPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, ctnr := range resp.ContainerItems {
			_require.NotNil(ctnr.Name)

			if *ctnr.Name == containerName {
				_require.NotNil(ctnr.Properties)
				_require.NotNil(ctnr.Properties.LastModified)
				_require.NotNil(ctnr.Properties.ETag)
				_require.Equal(*ctnr.Properties.LeaseStatus, lease.StatusTypeUnlocked)
				_require.Equal(*ctnr.Properties.LeaseState, lease.StateTypeAvailable)
				_require.Nil(ctnr.Properties.LeaseDuration)
				_require.Nil(ctnr.Properties.PublicAccess)
				_require.NotNil(ctnr.Metadata)

				unwrappedMeta := map[string]*string{}
				for k, v := range ctnr.Metadata {
					if v != nil {
						unwrappedMeta[k] = v
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

func (s *ServiceRecordedTestsSuite) TestSetPropertiesLogging() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)

	loggingOpts := service.Logging{
		Read: enabled, Write: enabled, Delete: enabled,
		RetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}}
	opts := service.SetPropertiesOptions{Logging: &loggingOpts}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.Nil(err)
	resp1, err := svcClient.GetProperties(context.Background(), nil)

	_require.Nil(err)
	_require.Equal(resp1.Logging.Write, enabled)
	_require.Equal(resp1.Logging.Read, enabled)
	_require.Equal(resp1.Logging.Delete, enabled)
	_require.Equal(resp1.Logging.RetentionPolicy.Days, days)
	_require.Equal(resp1.Logging.RetentionPolicy.Enabled, enabled)
}

func (s *ServiceRecordedTestsSuite) TestSetPropertiesHourMetrics() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)

	metricsOpts := service.Metrics{
		Enabled: enabled, IncludeAPIs: enabled, RetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}}
	opts := service.SetPropertiesOptions{HourMetrics: &metricsOpts}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.Nil(err)
	resp1, err := svcClient.GetProperties(context.Background(), nil)

	_require.Nil(err)
	_require.Equal(resp1.HourMetrics.Enabled, enabled)
	_require.Equal(resp1.HourMetrics.IncludeAPIs, enabled)
	_require.Equal(resp1.HourMetrics.RetentionPolicy.Days, days)
	_require.Equal(resp1.HourMetrics.RetentionPolicy.Enabled, enabled)
}

func (s *ServiceRecordedTestsSuite) TestSetPropertiesMinuteMetrics() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)

	metricsOpts := service.Metrics{
		Enabled: enabled, IncludeAPIs: enabled, RetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}}
	opts := service.SetPropertiesOptions{MinuteMetrics: &metricsOpts}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.Nil(err)
	resp1, err := svcClient.GetProperties(context.Background(), nil)

	_require.Nil(err)
	_require.Equal(resp1.MinuteMetrics.Enabled, enabled)
	_require.Equal(resp1.MinuteMetrics.IncludeAPIs, enabled)
	_require.Equal(resp1.MinuteMetrics.RetentionPolicy.Days, days)
	_require.Equal(resp1.MinuteMetrics.RetentionPolicy.Enabled, enabled)
}

func (s *ServiceRecordedTestsSuite) TestSetPropertiesSetCORSMultiple() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	defaultAge := to.Ptr[int32](500)
	defaultStr := to.Ptr[string]("")

	allowedOrigins1 := "www.xyz.com"
	allowedMethods1 := "GET"
	CORSOpts1 := &service.CORSRule{AllowedOrigins: &allowedOrigins1, AllowedMethods: &allowedMethods1}

	allowedOrigins2 := "www.xyz.com,www.ab.com,www.bc.com"
	allowedMethods2 := "GET, PUT"
	maxAge2 := to.Ptr[int32](500)
	exposedHeaders2 := "x-ms-meta-data*,x-ms-meta-source*,x-ms-meta-abc,x-ms-meta-bcd"
	allowedHeaders2 := "x-ms-meta-data*,x-ms-meta-target*,x-ms-meta-xyz,x-ms-meta-foo"

	CORSOpts2 := &service.CORSRule{
		AllowedOrigins: &allowedOrigins2, AllowedMethods: &allowedMethods2,
		MaxAgeInSeconds: maxAge2, ExposedHeaders: &exposedHeaders2, AllowedHeaders: &allowedHeaders2}

	CORSRules := []*service.CORSRule{CORSOpts1, CORSOpts2}

	opts := service.SetPropertiesOptions{CORS: CORSRules}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.Nil(err)
	resp, err := svcClient.GetProperties(context.Background(), nil)
	for i := 0; i < len(resp.CORS); i++ {
		if resp.CORS[i].AllowedOrigins == &allowedOrigins1 {
			_require.Equal(resp.CORS[i].AllowedMethods, &allowedMethods1)
			_require.Equal(resp.CORS[i].MaxAgeInSeconds, defaultAge)
			_require.Equal(resp.CORS[i].ExposedHeaders, defaultStr)
			_require.Equal(resp.CORS[i].AllowedHeaders, defaultStr)

		} else if resp.CORS[i].AllowedOrigins == &allowedOrigins2 {
			_require.Equal(resp.CORS[i].AllowedMethods, &allowedMethods2)
			_require.Equal(resp.CORS[i].MaxAgeInSeconds, &maxAge2)
			_require.Equal(resp.CORS[i].ExposedHeaders, &exposedHeaders2)
			_require.Equal(resp.CORS[i].AllowedHeaders, &allowedHeaders2)
		}
	}
	_require.Nil(err)
}

func (s *ServiceUnrecordedTestsSuite) TestListContainersBasicUsingConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.GetContainerClient(containerName, svcClient)
	_, err = containerClient.Create(context.Background(), &container.CreateOptions{Metadata: md})
	defer func(containerClient *container.Client, ctx context.Context, options *container.DeleteOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			_require.Nil(err)
		}
	}(containerClient, context.Background(), nil)
	_require.Nil(err)
	prefix := testcommon.ContainerPrefix
	listOptions := service.ListContainersOptions{Prefix: &prefix, Include: service.ListContainersInclude{Metadata: true}}
	pager := svcClient.NewListContainersPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		for _, ctnr := range resp.ContainerItems {
			_require.NotNil(ctnr.Name)

			if *ctnr.Name == containerName {
				_require.NotNil(ctnr.Properties)
				_require.NotNil(ctnr.Properties.LastModified)
				_require.NotNil(ctnr.Properties.ETag)
				_require.Equal(*ctnr.Properties.LeaseStatus, lease.StatusTypeUnlocked)
				_require.Equal(*ctnr.Properties.LeaseState, lease.StateTypeAvailable)
				_require.Nil(ctnr.Properties.LeaseDuration)
				_require.Nil(ctnr.Properties.PublicAccess)
				_require.NotNil(ctnr.Metadata)

				unwrappedMeta := map[string]*string{}
				for k, v := range ctnr.Metadata {
					if v != nil {
						unwrappedMeta[k] = v
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

func (s *ServiceUnrecordedTestsSuite) TestListContainersPaged() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)
	const numContainers = 6
	maxResults := int32(2)
	const pagedContainersPrefix = "azcontainerpaged"

	containers := make([]*container.Client, numContainers)
	expectedResults := make(map[string]bool)
	for i := 0; i < numContainers; i++ {
		containerName := pagedContainersPrefix + testcommon.GenerateContainerName(testName) + fmt.Sprintf("%d", i)
		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
		containers[i] = containerClient
		expectedResults[containerName] = false
	}

	defer func() {
		for i := range containers {
			testcommon.DeleteContainer(context.Background(), _require, containers[i])
		}
	}()

	prefix := pagedContainersPrefix + testcommon.ContainerPrefix
	listOptions := service.ListContainersOptions{MaxResults: &maxResults, Prefix: &prefix, Include: service.ListContainersInclude{Metadata: true}}
	count := 0
	results := make([]service.ContainerItem, 0)
	pager := svcClient.NewListContainersPager(&listOptions)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, ctnr := range resp.ContainerItems {
			_require.NotNil(ctnr.Name)
			results = append(results, *ctnr)
			count += 1
		}
	}

	_require.Equal(count, numContainers)
	_require.Equal(len(results), numContainers)

	// make sure each container we see is expected
	for _, ctnr := range results {
		_, ok := expectedResults[*ctnr.Name]
		_require.Equal(ok, true)
		expectedResults[*ctnr.Name] = true
	}

	// make sure every expected container was seen
	for _, seen := range expectedResults {
		_require.Equal(seen, true)
	}

}

//func (s *ServiceRecordedTestsSuite) TestListContainersPaged() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := testcommon.GetServiceClient(_context.recording, testcommon.TestAccountDefault, nil)
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
//		containerName := pagedContainersPrefix + testcommon.GenerateContainerName(testName) + string(i)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		containers[i] = containerClient
//		expectedResults[containerName] = false
//	}
//
//	defer func() {
//		for i := range containers {
//			testcommon.DeleteContainer(context.Background(), _require, containers[i])
//		}
//	}()
//
//	// list for a first time
//	prefix := testcommon.ContainerPrefix + pagedContainersPrefix
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

func (s *ServiceRecordedTestsSuite) TestAccountListContainersEmptyPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient1 := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"1", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient1)
	containerClient2 := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"2", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient2)

	count := 0
	pager := svcClient.NewListContainersPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
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
////func (s *ServiceRecordedTestsSuite) TestAccountListContainersMaxResultsNegative() {
////	svcClient := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
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
////func (s *ServiceRecordedTestsSuite) TestAccountListContainersMaxResultsExact() {
////	// If this test fails, ensure there are no extra containers prefixed with go in the account. These may be left over if a test is interrupted.
////	svcClient := testcommon.GetServiceClient()
////	containerClient1, containerName1 := createNewContainerWithSuffix(c, svcClient, "abc")
////	defer deleteContainer(containerClient1)
////	containerClient2, containerName2 := createNewContainerWithSuffix(c, svcClient, "abcde")
////	defer deleteContainer(containerClient2)
////
////	prefix := testcommon.ContainerPrefix + "abc"
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

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicy() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	disabled := false
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &disabled}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err = svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, false)
	_require.Nil(resp.StorageServiceProperties.DeleteRetentionPolicy.Days)
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyEmpty() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Empty retention policy causes an error, this is different from track 1.5
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{}})
	_require.NotNil(err)
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyNil() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{})
	_require.Nil(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Disable for other tests
	enabled = to.Ptr(false)
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled}})
	_require.Nil(err)
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyDaysTooSmall() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := int32(0) // Minimum days is 1. Validated on the client.
	enabled := true
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled, Days: &days}})
	_require.NotNil(err)
}

func (s *ServiceUnrecordedTestsSuite) TestAccountDeleteRetentionPolicyDaysTooLarge() {
	_require := require.New(s.T())
	var svcClient *service.Client
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
		} else {
			svcClient, err = testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
		}
		_require.Nil(err)

		days := int32(366) // Max days is 365. Left to the service for validation.
		enabled := true
		_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled, Days: &days}})
		_require.NotNil(err)

		testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
	}
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyDaysOmitted() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Days is required if enabled is true.
	enabled := true
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled}})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
}

func (s *ServiceUnrecordedTestsSuite) TestSASServiceClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	_require.Nil(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.Nil(err)

	containerName := testcommon.GenerateContainerName(testName)

	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Update: true,
		Delete: true,
	}

	expiry := time.Now().Add(time.Hour)

	sasUrl, err := serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.Nil(err)

	svcClient, err := service.NewClientWithNoCredential(sasUrl, nil)
	_require.Nil(err)

	// create container using SAS
	_, err = svcClient.CreateContainer(context.Background(), containerName, nil)
	_require.Nil(err)

	_, err = svcClient.DeleteContainer(context.Background(), containerName, nil)
	_require.Nil(err)
}

func (s *ServiceUnrecordedTestsSuite) TestNoSharedKeyCredError() {
	_require := require.New(s.T())
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")

	// Creating service client without credentials
	serviceClient, err := service.NewClientWithNoCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), nil)
	_require.Nil(err)

	// Adding SAS and options
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Update: true,
		Delete: true,
	}
	start := time.Now().Add(-time.Hour)
	expiry := start.Add(time.Hour)
	opts := service.GetSASURLOptions{StartTime: &start}

	// GetSASURL fails (with MissingSharedKeyCredential) because service client is created without credentials
	_, err = serviceClient.GetSASURL(resources, permissions, expiry, &opts)
	_require.Equal(err, bloberror.MissingSharedKeyCredential)

}

func (s *ServiceUnrecordedTestsSuite) TestSASContainerClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	_require.Nil(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.Nil(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := serviceClient.NewContainerClient(containerName)

	permissions := sas.ContainerPermissions{
		Read: true,
		Add:  true,
	}
	start := time.Now().Add(-5 * time.Minute).UTC()
	expiry := time.Now().Add(time.Hour)

	opts := container.GetSASURLOptions{StartTime: &start}
	sasUrl, err := containerClient.GetSASURL(permissions, expiry, &opts)
	_require.Nil(err)

	containerClient2, err := container.NewClientWithNoCredential(sasUrl, nil)
	_require.Nil(err)

	_, err = containerClient2.Create(context.Background(), &container.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.AuthorizationFailure)
}

func (s *ServiceUnrecordedTestsSuite) TestSASContainerClient2() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	_require.Nil(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.Nil(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := serviceClient.NewContainerClient(containerName)
	start := time.Now().Add(-5 * time.Minute).UTC()
	opts := container.GetSASURLOptions{StartTime: &start}

	sasUrlReadAdd, err := containerClient.GetSASURL(sas.ContainerPermissions{Read: true, Add: true}, time.Now().Add(time.Hour), &opts)
	_require.Nil(err)
	_, err = containerClient.Create(context.Background(), &container.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.Nil(err)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	containerClient1, err := container.NewClientWithNoCredential(sasUrlReadAdd, nil)
	_require.Nil(err)

	// container metadata and properties can't be read or written with SAS auth
	_, err = containerClient1.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.AuthorizationFailure)

	start = time.Now().Add(-5 * time.Minute).UTC()
	opts = container.GetSASURLOptions{StartTime: &start}

	sasUrlRCWL, err := containerClient.GetSASURL(sas.ContainerPermissions{Add: true, Create: true, Delete: true, List: true}, time.Now().Add(time.Hour), &opts)
	_require.Nil(err)

	containerClient2, err := container.NewClientWithNoCredential(sasUrlRCWL, nil)
	_require.Nil(err)

	// containers can't be created, deleted, or listed with SAS auth
	_, err = containerClient2.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.AuthorizationFailure)
}

// make sure that container soft delete is enabled
func (s *ServiceRecordedTestsSuite) TestContainerRestore() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	testName := s.T().Name()
	containerName := testcommon.GenerateContainerName(testName)

	_, err = svcClient.CreateContainer(context.Background(), containerName, nil)
	_require.Nil(err)

	_, err = svcClient.DeleteContainer(context.Background(), containerName, nil)
	_require.Nil(err)

	// it appears that deleting the container involves acquiring a lease.
	// since leases can only be 15-60s or infinite, we just wait for 60 seconds.
	time.Sleep(60 * time.Second)
	prefix := testcommon.ContainerPrefix
	listOptions := service.ListContainersOptions{Prefix: &prefix, Include: service.ListContainersInclude{Metadata: true, Deleted: true}}
	pager := svcClient.NewListContainersPager(&listOptions)

	contRestored := false
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, cont := range resp.ContainerItems {
			_require.NotNil(cont.Name)

			if *cont.Deleted && *cont.Name == containerName {
				_, err = svcClient.RestoreContainer(context.Background(), containerName, *cont.Version, nil)
				_require.NoError(err)
				contRestored = true
				break
			}
		}
		if contRestored {
			break
		}
	}

	_require.Equal(contRestored, true)

	for i := 0; i < 5; i++ {
		_, err = svcClient.DeleteContainer(context.Background(), containerName, nil)
		if err == nil {
			// container was deleted
			break
		} else if bloberror.HasCode(err, bloberror.Code("ConcurrentContainerOperationInProgress")) {
			// the container is still being restored, sleep a bit then try again
			time.Sleep(10 * time.Second)
		} else {
			// some other error
			break
		}
	}
	_require.NoError(err)
}

func (s *ServiceRecordedTestsSuite) TestContainerRestoreFailures() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	containerName := testcommon.GenerateContainerName(testName)

	_, err = svcClient.RestoreContainer(context.Background(), containerName, "", nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MissingRequiredHeader)

	_, err = svcClient.RestoreContainer(context.Background(), "", "", &service.RestoreContainerOptions{})
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MissingRequiredHeader)
}

func (s *ServiceUnrecordedTestsSuite) TestServiceSASUploadDownload() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	containerName := testcommon.GenerateContainerName(testName)

	_, err = svcClient.CreateContainer(context.Background(), containerName, nil)
	_require.Nil(err)

	credential, err := testcommon.GetGenericCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		Permissions:   to.Ptr(sas.BlobPermissions{Read: true, Create: true, Write: true, Tag: true}).String(),
		ContainerName: containerName,
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	sasURL := svcClient.URL()
	if len(sasURL) > 0 && sasURL[len(sasURL)-1:] != "/" {
		sasURL += "/"
	}
	sasURL += "?" + sasQueryParams.Encode()

	azClient, err := azblob.NewClientWithNoCredential(sasURL, nil)
	_require.Nil(err)

	const blobData = "test data"
	blobName := testcommon.GenerateBlobName(testName)
	_, err = azClient.UploadStream(context.TODO(),
		containerName,
		blobName,
		strings.NewReader(blobData),
		&azblob.UploadStreamOptions{
			Metadata: testcommon.BasicMetadata,
			Tags:     map[string]string{"Year": "2022"},
		})
	_require.Nil(err)

	blobDownloadResponse, err := azClient.DownloadStream(context.TODO(), containerName, blobName, nil)
	_require.Nil(err)

	reader := blobDownloadResponse.Body
	downloadData, err := io.ReadAll(reader)
	_require.Nil(err)
	_require.Equal(string(downloadData), blobData)

	_, err = svcClient.DeleteContainer(context.Background(), containerName, nil)
	_require.Nil(err)

	err = reader.Close()
	_require.Nil(err)
}
