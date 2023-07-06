//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
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

func (s *ServiceRecordedTestsSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	connectionString, _ := testcommon.GetGenericConnectionString(testcommon.TestAccountDatalake)

	parsedConnStr, err := shared.ParseConnectionString(*connectionString)
	_require.Nil(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".blob.core.windows.net/")

	sharedKeyCred, err := azdatalake.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.Nil(err)

	svcClient, err := service.NewClientWithSharedKeyCredential(parsedConnStr.ServiceURL, sharedKeyCred, nil)
	_require.Nil(err)
	fsClient := testcommon.CreateNewFilesystem(context.Background(), _require, testcommon.GenerateFilesystemName(testName), svcClient)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)
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
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
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
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
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
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
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

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicy() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
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
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
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
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
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

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyDaysTooLarge() {
	_require := require.New(s.T())
	var svcClient *service.Client
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
		} else {
			svcClient, err = testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDatalake, nil)
		}
		_require.Nil(err)

		days := int32(366) // Max days is 365. Left to the service for validation.
		enabled := true
		_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled, Days: &days}})
		_require.NotNil(err)

		testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
	}
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyDaysOmitted() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	// Days is required if enabled is true.
	enabled := true
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled}})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
}

func (s *ServiceRecordedTestsSuite) TestSASServiceClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	cred, _ := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDatalake)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.dfs.core.windows.net/", cred.AccountName()), cred, nil)
	_require.Nil(err)

	fsName := testcommon.GenerateFilesystemName(testName)

	// Note: Always set all permissions, services, types to true to ensure order of string formed is correct.
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:                  true,
		Write:                 true,
		Delete:                true,
		DeletePreviousVersion: true,
		List:                  true,
		Add:                   true,
		Create:                true,
		Update:                true,
		Process:               true,
		Tag:                   true,
		FilterByTags:          true,
		PermanentDelete:       true,
	}
	expiry := time.Now().Add(time.Hour)
	sasUrl, err := serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.Nil(err)

	svcClient, err := testcommon.GetServiceClientNoCredential(s.T(), sasUrl, nil)
	_require.Nil(err)

	// create fs using SAS
	_, err = svcClient.CreateFilesystem(context.Background(), fsName, nil)
	_require.Nil(err)

	_, err = svcClient.DeleteFilesystem(context.Background(), fsName, nil)
	_require.Nil(err)
}

func (s *ServiceRecordedTestsSuite) TestSASServiceClientNoKey() {
	_require := require.New(s.T())
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")

	serviceClient, err := service.NewClientWithNoCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), nil)
	_require.Nil(err)
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:                  true,
		Write:                 true,
		Delete:                true,
		DeletePreviousVersion: true,
		List:                  true,
		Add:                   true,
		Create:                true,
		Update:                true,
		Process:               true,
		Tag:                   true,
		FilterByTags:          true,
		PermanentDelete:       true,
	}

	expiry := time.Now().Add(time.Hour)
	_, err = serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.Equal(err.Error(), "SAS can only be signed with a SharedKeyCredential")
}

func (s *ServiceUnrecordedTestsSuite) TestSASServiceClientSignNegative() {
	_require := require.New(s.T())
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	_require.Nil(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.Nil(err)
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:                  true,
		Write:                 true,
		Delete:                true,
		DeletePreviousVersion: true,
		List:                  true,
		Add:                   true,
		Create:                true,
		Update:                true,
		Process:               true,
		Tag:                   true,
		FilterByTags:          true,
		PermanentDelete:       true,
	}
	expiry := time.Time{}
	_, err = serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.Equal(err.Error(), "account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType")
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
	_require.Equal(err, datalakeerror.MissingSharedKeyCredential)

}

func (s *ServiceRecordedTestsSuite) TestSASFilesystemClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	_require.Nil(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.Nil(err)

	fsName := testcommon.GenerateFilesystemName(testName)
	fsClient := serviceClient.NewFilesystemClient(fsName)

	permissions := sas.FilesystemPermissions{
		Read: true,
		Add:  true,
	}
	start := time.Now().Add(-5 * time.Minute).UTC()
	expiry := time.Now().Add(time.Hour)

	opts := filesystem.GetSASURLOptions{StartTime: &start}
	sasUrl, err := fsClient.GetSASURL(permissions, expiry, &opts)
	_require.Nil(err)

	fsClient2, err := filesystem.NewClientWithNoCredential(sasUrl, nil)
	_require.Nil(err)

	_, err = fsClient2.Create(context.Background(), &filesystem.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.AuthorizationFailure)
}

func (s *ServiceRecordedTestsSuite) TestSASFilesystem2() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	_require.Nil(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.Nil(err)

	fsName := testcommon.GenerateFilesystemName(testName)
	fsClient := serviceClient.NewFilesystemClient(fsName)
	start := time.Now().Add(-5 * time.Minute).UTC()
	opts := filesystem.GetSASURLOptions{StartTime: &start}

	sasUrlReadAdd, err := fsClient.GetSASURL(sas.FilesystemPermissions{Read: true, Add: true}, time.Now().Add(time.Hour), &opts)
	_require.Nil(err)
	_, err = fsClient.Create(context.Background(), &filesystem.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.Nil(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	fsClient1, err := filesystem.NewClientWithNoCredential(sasUrlReadAdd, nil)
	_require.Nil(err)

	// filesystem metadata and properties can't be read or written with SAS auth
	_, err = fsClient1.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.AuthorizationFailure)

	start = time.Now().Add(-5 * time.Minute).UTC()
	opts = filesystem.GetSASURLOptions{StartTime: &start}

	sasUrlRCWL, err := fsClient.GetSASURL(sas.FilesystemPermissions{Add: true, Create: true, Delete: true, List: true}, time.Now().Add(time.Hour), &opts)
	_require.Nil(err)

	fsClient2, err := filesystem.NewClientWithNoCredential(sasUrlRCWL, nil)
	_require.Nil(err)

	// filesystems can't be created, deleted, or listed with SAS auth
	_, err = fsClient2.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, datalakeerror.AuthorizationFailure)
}

func (s *ServiceRecordedTestsSuite) TestListFilesystemsBasic() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.Nil(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}

	fsName := testcommon.GenerateFilesystemName(testName)
	fsClient := testcommon.ServiceGetFilesystemClient(fsName, svcClient)
	_, err = fsClient.Create(context.Background(), &filesystem.CreateOptions{Metadata: md})
	defer func(fsClient *filesystem.Client, ctx context.Context, options *filesystem.DeleteOptions) {
		_, err := fsClient.Delete(ctx, options)
		if err != nil {
			_require.Nil(err)
		}
	}(fsClient, context.Background(), nil)
	_require.Nil(err)
	prefix := testcommon.FilesystemPrefix
	listOptions := service.ListFilesystemsOptions{Prefix: &prefix, Include: service.ListFilesystemsInclude{Metadata: true}}
	pager := svcClient.NewListFilesystemsPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, ctnr := range resp.Filesystems {
			_require.NotNil(ctnr.Name)

			if *ctnr.Name == fsName {
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

func (s *ServiceRecordedTestsSuite) TestListFilesystemsBasicUsingConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}

	fsName := testcommon.GenerateFilesystemName(testName)
	fsClient := testcommon.ServiceGetFilesystemClient(fsName, svcClient)
	_, err = fsClient.Create(context.Background(), &filesystem.CreateOptions{Metadata: md})
	defer func(fsClient *filesystem.Client, ctx context.Context, options *filesystem.DeleteOptions) {
		_, err := fsClient.Delete(ctx, options)
		if err != nil {
			_require.Nil(err)
		}
	}(fsClient, context.Background(), nil)
	_require.Nil(err)
	prefix := testcommon.FilesystemPrefix
	listOptions := service.ListFilesystemsOptions{Prefix: &prefix, Include: service.ListFilesystemsInclude{Metadata: true}}
	pager := svcClient.NewListFilesystemsPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		for _, ctnr := range resp.Filesystems {
			_require.NotNil(ctnr.Name)

			if *ctnr.Name == fsName {
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

func (s *ServiceRecordedTestsSuite) TestListFilesystemsPaged() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)
	const numFilesystems = 6
	maxResults := int32(2)
	const pagedFilesystemsPrefix = "azfilesystempaged"

	filesystems := make([]*filesystem.Client, numFilesystems)
	expectedResults := make(map[string]bool)
	for i := 0; i < numFilesystems; i++ {
		fsName := pagedFilesystemsPrefix + testcommon.GenerateFilesystemName(testName) + fmt.Sprintf("%d", i)
		fsClient := testcommon.CreateNewFilesystem(context.Background(), _require, fsName, svcClient)
		filesystems[i] = fsClient
		expectedResults[fsName] = false
	}

	defer func() {
		for i := range filesystems {
			testcommon.DeleteFilesystem(context.Background(), _require, filesystems[i])
		}
	}()

	prefix := pagedFilesystemsPrefix + testcommon.FilesystemPrefix
	listOptions := service.ListFilesystemsOptions{MaxResults: &maxResults, Prefix: &prefix, Include: service.ListFilesystemsInclude{Metadata: true}}
	count := 0
	results := make([]service.FilesystemItem, 0)
	pager := svcClient.NewListFilesystemsPager(&listOptions)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, ctnr := range resp.Filesystems {
			_require.NotNil(ctnr.Name)
			results = append(results, *ctnr)
			count += 1
		}
	}

	_require.Equal(count, numFilesystems)
	_require.Equal(len(results), numFilesystems)

	// make sure each fs we see is expected
	for _, ctnr := range results {
		_, ok := expectedResults[*ctnr.Name]
		_require.Equal(ok, true)
		expectedResults[*ctnr.Name] = true
	}

	// make sure every expected fs was seen
	for _, seen := range expectedResults {
		_require.Equal(seen, true)
	}

}

func (s *ServiceRecordedTestsSuite) TestAccountListFilesystemsEmptyPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	fsClient1 := testcommon.CreateNewFilesystem(context.Background(), _require, testcommon.GenerateFilesystemName(testName)+"1", svcClient)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient1)
	fsClient2 := testcommon.CreateNewFilesystem(context.Background(), _require, testcommon.GenerateFilesystemName(testName)+"2", svcClient)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient2)

	count := 0
	pager := svcClient.NewListFilesystemsPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		for _, container := range resp.Filesystems {
			count++
			_require.NotNil(container.Name)
		}
		if err != nil {
			break
		}
	}
	_require.GreaterOrEqual(count, 2)
}
