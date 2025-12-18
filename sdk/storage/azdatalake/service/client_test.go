//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
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

func (s *ServiceRecordedTestsSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *ServiceRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
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
	proxy *recording.TestProxyInstance
}

type ServiceUnrecordedTestsSuite struct {
	suite.Suite
}

func (s *ServiceUnrecordedTestsSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	connectionString, _ := testcommon.GetGenericConnectionString(testcommon.TestAccountDatalake)

	parsedConnStr, err := shared.ParseConnectionString(*connectionString)
	_require.NoError(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".blob.core.windows.net/")

	sharedKeyCred, err := azdatalake.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.NoError(err)

	svcClient, err := service.NewClientWithSharedKeyCredential(parsedConnStr.ServiceURL, sharedKeyCred, nil)
	_require.NoError(err)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, testcommon.GenerateFileSystemName(testName), svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)
}

func (s *ServiceRecordedTestsSuite) TestCreateFilesystemsWithOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}
	cpkScopeInfo := &testcommon.TestCPKScopeInfo

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.ServiceGetFileSystemClient(fsName, svcClient)

	_, err = fsClient.Create(context.Background(), &filesystem.CreateOptions{Metadata: md, CPKScopeInfo: cpkScopeInfo})
	defer func(fsClient *filesystem.Client, ctx context.Context, options *filesystem.DeleteOptions) {
		_, err := fsClient.Delete(ctx, options)
		if err != nil {
			_require.NoError(err)
		}
	}(fsClient, context.Background(), nil)

	_require.NoError(err)
	resp, err := fsClient.GetProperties(context.Background(), nil)

	_require.NoError(err)
	_require.Equal(resp.DefaultEncryptionScope, &testcommon.TestEncryptionScope)
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

	_require.NoError(err)
	resp1, err := svcClient.GetProperties(context.Background(), nil)

	_require.NoError(err)
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

	_require.NoError(err)
	resp1, err := svcClient.GetProperties(context.Background(), nil)

	_require.NoError(err)
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

	_require.NoError(err)
	resp1, err := svcClient.GetProperties(context.Background(), nil)

	_require.NoError(err)
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

	_require.NoError(err)
	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
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
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicy() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.NoError(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	disabled := false
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &disabled}})
	_require.NoError(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err = svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Empty retention policy causes an error, this is different from track 1.5
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{}})
	_require.Error(err)
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyNil() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled, Days: days}})
	_require.NoError(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{})
	_require.NoError(err)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	_require.EqualValues(*resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Disable for other tests
	enabled = to.Ptr(false)
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: enabled}})
	_require.NoError(err)
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
		_require.NoError(err)

		days := int32(366) // Max days is 365. Left to the service for validation.
		enabled := true
		_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled, Days: &days}})
		_require.Error(err)

		testcommon.ValidateErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
	}
}

func (s *ServiceRecordedTestsSuite) TestAccountDeleteRetentionPolicyDaysOmitted() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	// Days is required if enabled is true.
	enabled := true
	_, err = svcClient.SetProperties(context.Background(), &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: &enabled}})
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.InvalidXMLDocument)
}

func (s *ServiceUnrecordedTestsSuite) TestSASServiceClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	cred, _ := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDatalake)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.dfs.core.windows.net/", cred.AccountName()), cred, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)

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
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClientNoCredential(s.T(), sasUrl, nil)
	_require.NoError(err)

	// create fs using SAS
	_, err = svcClient.CreateFileSystem(context.Background(), fsName, nil)
	_require.NoError(err)

	_, err = svcClient.DeleteFileSystem(context.Background(), fsName, nil)
	_require.NoError(err)
}

func (s *ServiceRecordedTestsSuite) TestSASServiceClientNoKey() {
	_require := require.New(s.T())
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")

	serviceClient, err := service.NewClientWithNoCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)
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
	_require.NoError(err)

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

func (s *ServiceUnrecordedTestsSuite) TestGetFileSystemClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName + "1")
	fsClient := serviceClient.NewFileSystemClient(fsName)

	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fsClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *ServiceUnrecordedTestsSuite) TestSASFileSystemClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := serviceClient.NewFileSystemClient(fsName)

	permissions := sas.FileSystemPermissions{
		Read: true,
		Add:  true,
	}
	start := time.Now().Add(-5 * time.Minute).UTC()
	expiry := time.Now().Add(time.Hour)

	opts := filesystem.GetSASURLOptions{StartTime: &start}
	sasUrl, err := fsClient.GetSASURL(permissions, expiry, &opts)
	_require.NoError(err)

	fsClient2, err := filesystem.NewClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	_, err = fsClient2.Create(context.Background(), &filesystem.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.AuthorizationFailure)
}

func (s *ServiceUnrecordedTestsSuite) TestSASFileSystem2() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := serviceClient.NewFileSystemClient(fsName)
	start := time.Now().Add(-5 * time.Minute).UTC()
	opts := filesystem.GetSASURLOptions{StartTime: &start}

	sasUrlReadAdd, err := fsClient.GetSASURL(sas.FileSystemPermissions{Read: true, Add: true}, time.Now().Add(time.Hour), &opts)
	_require.NoError(err)
	_, err = fsClient.Create(context.Background(), &filesystem.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fsClient1, err := filesystem.NewClientWithNoCredential(sasUrlReadAdd, nil)
	_require.NoError(err)

	// filesystem metadata and properties can't be read or written with SAS auth
	_, err = fsClient1.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.AuthorizationFailure)

	start = time.Now().Add(-5 * time.Minute).UTC()
	opts = filesystem.GetSASURLOptions{StartTime: &start}

	sasUrlRCWL, err := fsClient.GetSASURL(sas.FileSystemPermissions{Add: true, Create: true, Delete: true, List: true}, time.Now().Add(time.Hour), &opts)
	_require.NoError(err)

	fsClient2, err := filesystem.NewClientWithNoCredential(sasUrlRCWL, nil)
	_require.NoError(err)

	// filesystems can't be created, deleted, or listed with SAS auth
	_, err = fsClient2.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.AuthorizationFailure)
}

func (s *ServiceRecordedTestsSuite) TestListFilesystemsBasic() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := svcClient.NewFileSystemClient(fsName)
	_, err = fsClient.Create(context.Background(), &filesystem.CreateOptions{Metadata: md})
	defer func(fsClient *filesystem.Client, ctx context.Context, options *filesystem.DeleteOptions) {
		_, err := fsClient.Delete(ctx, options)
		if err != nil {
			_require.NoError(err)
		}
	}(fsClient, context.Background(), nil)
	_require.NoError(err)
	prefix := testcommon.FileSystemPrefix
	listOptions := service.ListFileSystemsOptions{Prefix: &prefix, Include: service.ListFileSystemsInclude{Metadata: to.Ptr(true)}}
	pager := svcClient.NewListFileSystemsPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, ctnr := range resp.FileSystemItems {
			_require.NotNil(ctnr.Name)

			if *ctnr.Name == fsName {
				_require.NotNil(ctnr.Properties)
				_require.NotNil(ctnr.Properties.LastModified)
				_require.NotNil(ctnr.Properties.ETag)
				_require.Equal(*ctnr.Properties.LeaseStatus, service.StatusTypeUnlocked)
				_require.Equal(*ctnr.Properties.LeaseState, service.StateTypeAvailable)
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

	_require.NoError(err)
	_require.GreaterOrEqual(count, 0)
}

func (s *ServiceRecordedTestsSuite) TestListFilesystemsBasicUsingConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.ServiceGetFileSystemClient(fsName, svcClient)
	_, err = fsClient.Create(context.Background(), &filesystem.CreateOptions{Metadata: md})
	defer func(fsClient *filesystem.Client, ctx context.Context, options *filesystem.DeleteOptions) {
		_, err := fsClient.Delete(ctx, options)
		if err != nil {
			_require.NoError(err)
		}
	}(fsClient, context.Background(), nil)
	_require.NoError(err)
	prefix := testcommon.FileSystemPrefix
	listOptions := service.ListFileSystemsOptions{Prefix: &prefix, Include: service.ListFileSystemsInclude{Metadata: to.Ptr(true)}}
	pager := svcClient.NewListFileSystemsPager(&listOptions)

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, ctnr := range resp.FileSystemItems {
			_require.NotNil(ctnr.Name)

			if *ctnr.Name == fsName {
				_require.NotNil(ctnr.Properties)
				_require.NotNil(ctnr.Properties.LastModified)
				_require.NotNil(ctnr.Properties.ETag)
				_require.Equal(*ctnr.Properties.LeaseStatus, service.StatusTypeUnlocked)
				_require.Equal(*ctnr.Properties.LeaseState, service.StateTypeAvailable)
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

	_require.NoError(err)
	_require.GreaterOrEqual(count, 0)
}

func (s *ServiceRecordedTestsSuite) TestListFilesystemsIncludeSystemFileSystems() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	listOptions := service.ListFileSystemsOptions{Include: service.ListFileSystemsInclude{System: to.Ptr(true)}}
	count := 0
	pager := svcClient.NewListFileSystemsPager(&listOptions)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, ctnr := range resp.FileSystemItems {
			_require.NotNil(ctnr.Name)
			_require.Equal("$logs", *ctnr.Name)
			count += 1
		}
	}
	_require.Equal(1, count)
}

func (s *ServiceRecordedTestsSuite) TestListFilesystemsPaged() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	const numFileSystems = 6
	maxResults := int32(2)
	const pagedFileSystemsPrefix = "azfilesystempaged"

	filesystems := make([]*filesystem.Client, numFileSystems)
	expectedResults := make(map[string]bool)
	for i := 0; i < numFileSystems; i++ {
		fsName := pagedFileSystemsPrefix + testcommon.GenerateFileSystemName(testName) + fmt.Sprintf("%d", i)
		fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
		filesystems[i] = fsClient
		expectedResults[fsName] = false
	}

	defer func() {
		for i := range filesystems {
			testcommon.DeleteFileSystem(context.Background(), _require, filesystems[i])
		}
	}()

	prefix := pagedFileSystemsPrefix + testcommon.FileSystemPrefix
	listOptions := service.ListFileSystemsOptions{MaxResults: &maxResults, Prefix: &prefix, Include: service.ListFileSystemsInclude{Metadata: to.Ptr(true)}}
	count := 0
	results := make([]service.FileSystemItem, 0)
	pager := svcClient.NewListFileSystemsPager(&listOptions)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, ctnr := range resp.FileSystemItems {
			_require.NotNil(ctnr.Name)
			results = append(results, *ctnr)
			count += 1
		}
	}

	_require.Equal(count, numFileSystems)
	_require.Equal(len(results), numFileSystems)

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

	fsClient1 := testcommon.CreateNewFileSystem(context.Background(), _require, testcommon.GenerateFileSystemName(testName)+"1", svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient1)
	fsClient2 := testcommon.CreateNewFileSystem(context.Background(), _require, testcommon.GenerateFileSystemName(testName)+"2", svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient2)

	count := 0
	pager := svcClient.NewListFileSystemsPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, container := range resp.FileSystemItems {
			count++
			_require.NotNil(container.Name)
		}
		if err != nil {
			break
		}
	}
	_require.GreaterOrEqual(count, 2)
}

func (s *ServiceUnrecordedTestsSuite) TestServiceClientWithHTTP() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDatalake)
	_require.NoError(err)

	svcClient, err := service.NewClientWithSharedKeyCredential("http://"+cred.AccountName()+".dfs.core.windows.net/", cred, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fileName := testcommon.GenerateFileName(testName)
	fileClient := svcClient.NewFileSystemClient(fsName).NewFileClient(fileName)
	_require.Equal(fileClient.DFSURL(), "http://"+cred.AccountName()+".dfs.core.windows.net/"+fsName+"/"+fileName)

	_, err = fileClient.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, "AccountRequiresHttps")
}

func (s *ServiceRecordedTestsSuite) TestServiceClientWithNilSharedKey() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := service.NewClientWithSharedKeyCredential("http://"+accountName+".dfs.core.windows.net/", nil, nil)
	_require.Error(err)
	_require.Nil(svcClient)
}

func (s *ServiceRecordedTestsSuite) TestServiceClientUsingOauth() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	serviceUrl := "https://" + accountName + ".dfs.core.windows.net/"

	svcClient, err := service.NewClient(serviceUrl, cred, nil)
	_require.NoError(err)
	_require.NotNil(svcClient)

	fs, _ := svcClient.CreateFileSystem(context.Background(), "test", nil)
	_require.NotNil(fs)
	_require.NoError(err)
}

func (s *ServiceRecordedTestsSuite) TestServiceClientUsingOauthWithDefaultAudience() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	serviceUrl := "https://" + accountName + ".dfs.core.windows.net/"

	options := service.ClientOptions{
		Audience: "https://storage.azure.com/",
	}

	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	svcClient, err := service.NewClient(serviceUrl, cred, &options)
	_require.NoError(err)
	_require.NotNil(svcClient)

	fs, _ := svcClient.CreateFileSystem(context.Background(), "test", nil)
	_require.NotNil(fs)
	_require.NoError(err)

}

func (s *ServiceRecordedTestsSuite) TestServiceClientUsingOauthWithCustomAudience() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	serviceUrl := "https://" + accountName + ".dfs.core.windows.net/"

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := service.ClientOptions{
		Audience: "https://" + accountName + ".blob.core.windows.net",
	}

	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	svcClient, err := service.NewClient(serviceUrl, cred, &options)
	_require.NoError(err)
	_require.NotNil(svcClient)

	fs, _ := svcClient.CreateFileSystem(context.Background(), "test", nil)
	_require.NotNil(fs)
	_require.NoError(err)

}
