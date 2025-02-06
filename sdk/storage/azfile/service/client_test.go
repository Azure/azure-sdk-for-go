//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strconv"
	"strings"
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

func (s *ServiceRecordedTestsSuite) TestAccountNewServiceURLValidName() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	correctURL := "https://" + accountName + "." + testcommon.DefaultFileEndpointSuffix
	_require.Equal(svcClient.URL(), correctURL)
}

func (s *ServiceRecordedTestsSuite) TestAccountNewShareURLValidName() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)
	_require.NoError(err)

	correctURL := "https://" + accountName + "." + testcommon.DefaultFileEndpointSuffix + shareName
	_require.Equal(shareClient.URL(), correctURL)
}

func (s *ServiceRecordedTestsSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())

	svcClient, err := testcommon.GetServiceClientFromConnectionString(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
}

func (s *ServiceRecordedTestsSuite) TestAccountProperties() {
	_require := require.New(s.T())

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	setPropertiesOptions := &service.SetPropertiesOptions{
		HourMetrics: &service.Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(true),
			RetentionPolicy: &service.RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(2)),
			},
		},
		MinuteMetrics: &service.Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(false),
			RetentionPolicy: &service.RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(2)),
			},
		},
		CORS: []*service.CORSRule{
			{
				AllowedOrigins:  to.Ptr("*"),
				AllowedMethods:  to.Ptr("PUT"),
				AllowedHeaders:  to.Ptr("x-ms-client-request-id"),
				ExposedHeaders:  to.Ptr("x-ms-*"),
				MaxAgeInSeconds: to.Ptr(int32(2)),
			},
		},
	}

	setPropsResp, err := svcClient.SetProperties(context.Background(), setPropertiesOptions)
	_require.NoError(err)
	_require.NotNil(setPropsResp.RequestID)

	time.Sleep(time.Second * 30)

	getPropsResp, err := svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getPropsResp.RequestID)
	_require.EqualValues(getPropsResp.HourMetrics.RetentionPolicy.Enabled, setPropertiesOptions.HourMetrics.RetentionPolicy.Enabled)
	_require.EqualValues(getPropsResp.HourMetrics.RetentionPolicy.Days, setPropertiesOptions.HourMetrics.RetentionPolicy.Days)
	_require.EqualValues(getPropsResp.MinuteMetrics.RetentionPolicy.Enabled, setPropertiesOptions.MinuteMetrics.RetentionPolicy.Enabled)
	_require.EqualValues(getPropsResp.MinuteMetrics.RetentionPolicy.Days, setPropertiesOptions.MinuteMetrics.RetentionPolicy.Days)
	_require.EqualValues(len(getPropsResp.CORS), len(setPropertiesOptions.CORS))
}

func (s *ServiceRecordedTestsSuite) TestAccountHourMetrics() {
	_require := require.New(s.T())

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	setPropertiesOptions := &service.SetPropertiesOptions{
		HourMetrics: &service.Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(true),
			RetentionPolicy: &service.RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(5)),
			},
		},
	}
	_, err = svcClient.SetProperties(context.Background(), setPropertiesOptions)
	_require.NoError(err)
}

type userAgentTest struct{}

func (u userAgentTest) Do(req *policy.Request) (*http.Response, error) {
	const userAgentHeader = "User-Agent"

	currentUserAgentHeader := req.Raw().Header.Get(userAgentHeader)
	if !strings.HasPrefix(currentUserAgentHeader, "azsdk-go-azfile/"+exported.ModuleVersion) {
		return nil, fmt.Errorf("%s user agent doesn't match expected agent: azsdk-go-azfile/vx.xx.x", currentUserAgentHeader)
	}

	return &http.Response{
		Request:    req.Raw(),
		Status:     "Created",
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       http.NoBody,
	}, nil
}

func newTelemetryTestPolicy() policy.Policy {
	return &userAgentTest{}
}

func TestUserAgentForAzFile(t *testing.T) {
	client, err := service.NewClientWithNoCredential("https://fake/blob/testpath", &service.ClientOptions{
		ClientOptions: policy.ClientOptions{
			PerCallPolicies: []policy.Policy{newTelemetryTestPolicy()},
		},
	})
	require.NoError(t, err)

	_, err = client.GetProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func (s *ServiceRecordedTestsSuite) TestAccountListSharesNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	mySharePrefix := testcommon.GenerateEntityName(testName)
	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix: to.Ptr(mySharePrefix),
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.NotNil(resp.Prefix)
		_require.Equal(*resp.Prefix, mySharePrefix)
		_require.NotNil(resp.ServiceEndpoint)
		_require.NotNil(resp.Version)
		_require.Len(resp.Shares, 0)
	}

	shareClients := map[string]*share.Client{}
	for i := 0; i < 4; i++ {
		shareName := mySharePrefix + "share" + strconv.Itoa(i)
		shareClients[shareName] = testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
		defer testcommon.DeleteShare(context.Background(), _require, shareClients[shareName])

		_, err := shareClients[shareName].SetMetadata(context.Background(), &share.SetMetadataOptions{
			Metadata: testcommon.BasicMetadata,
		})
		_require.NoError(err)
	}

	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include:    service.ListSharesInclude{Metadata: true, Snapshots: true},
		Prefix:     to.Ptr(mySharePrefix),
		MaxResults: to.Ptr(int32(2)),
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		if len(resp.Shares) > 0 {
			_require.Len(resp.Shares, 2)
		}
		for _, shareItem := range resp.Shares {
			_require.NotNil(shareItem.Properties)
			_require.NotNil(shareItem.Properties.LastModified)
			_require.NotNil(shareItem.Properties.ETag)
			_require.EqualValues(shareItem.Metadata, testcommon.BasicMetadata)
		}
	}
}

func (s *ServiceRecordedTestsSuite) TestListSharesEnableSnapshotVirtualDirectoryAccess() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	mySharePrefix := testcommon.GenerateEntityName(testName)
	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix: to.Ptr(mySharePrefix),
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(*resp.Prefix, mySharePrefix)
		_require.Len(resp.Shares, 0)
	}

	shareClients := map[string]*share.Client{}
	for i := 0; i < 4; i++ {
		shareName := mySharePrefix + "share" + strconv.Itoa(i)
		shareClients[shareName] = svcClient.NewShareClient(shareName)
		_, err = shareClients[shareName].Create(context.Background(), &share.CreateOptions{EnabledProtocols: to.Ptr("NFS")})
		defer testcommon.DeleteShare(context.Background(), _require, shareClients[shareName])
		_require.NoError(err)

		_, err := shareClients[shareName].SetMetadata(context.Background(), &share.SetMetadataOptions{
			Metadata: testcommon.BasicMetadata,
		})
		_require.NoError(err)

		_, err = shareClients[shareName].SetProperties(context.Background(), &share.SetPropertiesOptions{
			EnableSnapshotVirtualDirectoryAccess: to.Ptr(true),
		})
		_require.NoError(err)

	}

	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include:    service.ListSharesInclude{Metadata: true, Snapshots: true},
		Prefix:     to.Ptr(mySharePrefix),
		MaxResults: to.Ptr(int32(2)),
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		if len(resp.Shares) > 0 {
			_require.Len(resp.Shares, 2)
		}
		for _, shareItem := range resp.Shares {
			_require.NotNil(shareItem.Properties)
			_require.Equal(shareItem.Properties.EnableSnapshotVirtualDirectoryAccess, to.Ptr(true))
		}
	}
}

func (s *ServiceUnrecordedTestsSuite) TestSASServiceClientRestoreShare() {
	_require := require.New(s.T())
	testName := s.T().Name()
	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.file.core.windows.net/", cred.AccountName()), cred, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)

	// Note: Always set all permissions, services, types to true to ensure order of string formed is correct.
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)
	sasUrl, err := serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClientNoCredential(s.T(), sasUrl, nil)
	_require.NoError(err)

	resp, err := svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)

	// create share using account SAS
	_, err = svcClient.CreateShare(context.Background(), shareName, nil)
	_require.NoError(err)

	defer func() {
		_, err := svcClient.DeleteShare(context.Background(), shareName, nil)
		_require.NoError(err)
	}()

	_, err = svcClient.DeleteShare(context.Background(), shareName, nil)
	_require.NoError(err)

	// wait for share deletion
	time.Sleep(60 * time.Second)

	sharesCnt := 0
	shareVersion := ""

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Deleted: true},
		Prefix:  &shareName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, s := range resp.Shares {
			if s.Deleted != nil && *s.Deleted {
				_require.NotNil(s.Version)
				shareVersion = *s.Version
			} else {
				sharesCnt++
			}
		}
	}

	_require.Equal(sharesCnt, 0)
	_require.NotEmpty(shareVersion)

	restoreResp, err := svcClient.RestoreShare(context.Background(), shareName, shareVersion, nil)
	_require.NoError(err)
	_require.NotNil(restoreResp.RequestID)

	sharesCnt = 0
	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix: &shareName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		sharesCnt += len(resp.Shares)
	}
	_require.Equal(sharesCnt, 1)
}

func (s *ServiceRecordedTestsSuite) TestSASServiceClientNoKey() {
	_require := require.New(s.T())
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	serviceClient, err := service.NewClientWithNoCredential(fmt.Sprintf("https://%s.file.core.windows.net/", accountName), nil)
	_require.NoError(err)
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}

	expiry := time.Now().Add(time.Hour)
	_, err = serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.Equal(err, fileerror.MissingSharedKeyCredential)
}

func (s *ServiceRecordedTestsSuite) TestSASServiceClientSignNegative() {
	_require := require.New(s.T())
	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)
	_require.Greater(len(accountKey), 0)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.file.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Time{}

	// zero expiry time
	_, err = serviceClient.GetSASURL(resources, permissions, expiry, &service.GetSASURLOptions{StartTime: to.Ptr(time.Now())})
	_require.Equal(err.Error(), "account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType")

	// zero start and expiry time
	_, err = serviceClient.GetSASURL(resources, permissions, expiry, &service.GetSASURLOptions{})
	_require.Equal(err.Error(), "account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType")

	// empty permissions
	_, err = serviceClient.GetSASURL(sas.AccountResourceTypes{}, sas.AccountPermissions{}, expiry, nil)
	_require.Equal(err.Error(), "account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType")
}

func (s *ServiceRecordedTestsSuite) TestServiceSetPropertiesDefault() {
	_require := require.New(s.T())

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	_, err = svcClient.SetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *ServiceRecordedTestsSuite) TestServiceCreateDeleteRestoreShare() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)

	_, err = svcClient.CreateShare(context.Background(), shareName, nil)
	_require.NoError(err)

	defer func() {
		_, err := svcClient.DeleteShare(context.Background(), shareName, nil)
		_require.NoError(err)
	}()

	_, err = svcClient.DeleteShare(context.Background(), shareName, nil)
	_require.NoError(err)

	// wait for share deletion
	time.Sleep(60 * time.Second)

	sharesCnt := 0
	shareVersion := ""

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Deleted: true},
		Prefix:  &shareName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, s := range resp.Shares {
			if s.Deleted != nil && *s.Deleted {
				_require.NotNil(s.Version)
				shareVersion = *s.Version
			} else {
				sharesCnt++
			}
		}
	}

	_require.Equal(sharesCnt, 0)
	_require.NotEmpty(shareVersion)

	restoreResp, err := svcClient.RestoreShare(context.Background(), shareName, shareVersion, nil)
	_require.NoError(err)
	_require.NotNil(restoreResp.RequestID)

	sharesCnt = 0
	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix: &shareName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		sharesCnt += len(resp.Shares)
	}
	_require.Equal(sharesCnt, 1)
}

func (s *ServiceRecordedTestsSuite) TestServiceCreateDeleteDirOAuth() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// create service client using token credential
	options := &service.ClientOptions{FileRequestIntent: to.Ptr(service.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	svcClientOAuth, err := service.NewClient("https://"+accountName+".file.core.windows.net/", cred, options)
	_require.NoError(err)

	dirClient := svcClientOAuth.NewShareClient(shareName).NewDirectoryClient(testcommon.GenerateDirectoryName(testName))

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *ServiceUnrecordedTestsSuite) TestAccountSASEncryptionScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.EncryptionScopeEnvVar)
	_require.NoError(err)

	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Now().Add(1 * time.Hour)

	qps, err := sas.AccountSignatureValues{
		Permissions:     permissions.String(),
		ResourceTypes:   resources.String(),
		ExpiryTime:      expiry.UTC(),
		EncryptionScope: encryptionScope,
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	svcSAS := fmt.Sprintf("https://%s.file.core.windows.net/?%s", cred.AccountName(), qps.Encode())
	svcClient, err := service.NewClientWithNoCredential(svcSAS, nil)
	_require.NoError(err)

	_, err = svcClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	_, err = svcClient.CreateShare(context.Background(), shareName, nil)
	_require.NoError(err)
	defer func() {
		_, err = svcClient.DeleteShare(context.Background(), shareName, nil)
		_require.NoError(err)
	}()

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix: &shareName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Shares, 1)
		_require.NotNil(resp.Shares[0].Name)
		_require.Equal(*resp.Shares[0].Name, shareName)
	}

	shareClient := svcClient.NewShareClient(shareName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))

	_, err = fileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	_, err = fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *ServiceRecordedTestsSuite) TestPremiumAccountListShares() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	mySharePrefix := testcommon.GenerateEntityName(testName)
	shareClients := map[string]*share.Client{}
	for i := 0; i < 4; i++ {
		shareName := mySharePrefix + "share" + strconv.Itoa(i)
		shareClients[shareName] = testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
		defer testcommon.DeleteShare(context.Background(), _require, shareClients[shareName])
	}

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix: to.Ptr(mySharePrefix),
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Shares, 4)
		for _, shareItem := range resp.Shares {
			_require.NotNil(shareItem.Properties)
			_require.NotNil(shareItem.Properties.ProvisionedBandwidthMiBps)
			_require.NotNil(shareItem.Properties.ProvisionedIngressMBps)
			_require.NotNil(shareItem.Properties.ProvisionedEgressMBps)
			_require.NotNil(shareItem.Properties.ProvisionedIops)
			_require.NotNil(shareItem.Properties.NextAllowedQuotaDowngradeTime)
			_require.Greater(*shareItem.Properties.ProvisionedBandwidthMiBps, (int32)(0))
		}
	}
}

func (s *ServiceUnrecordedTestsSuite) TestServiceClientWithHTTP() {
	_require := require.New(s.T())

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := service.NewClientWithSharedKeyCredential("http://"+cred.AccountName()+".file.core.windows.net/", cred, nil)
	_require.NoError(err)

	_, err = svcClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, "AccountRequiresHttps")
}

func (s *ServiceRecordedTestsSuite) TestServiceClientWithNilSharedKey() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	options := &service.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	svcClient, err := service.NewClientWithSharedKeyCredential("https://"+accountName+".file.core.windows.net/", nil, options)
	_require.NoError(err)

	_, err = svcClient.GetProperties(context.Background(), nil)
	_require.Error(err)
}

func (s *ServiceRecordedTestsSuite) TestServiceClientCustomAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// create service client using token credential
	options := &service.ClientOptions{
		FileRequestIntent: to.Ptr(service.ShareTokenIntentBackup),
		Audience:          "https://" + accountName + ".file.core.windows.net",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	svcClientAudience, err := service.NewClient("https://"+accountName+".file.core.windows.net/", cred, options)
	_require.NoError(err)

	dirClientAudience := svcClientAudience.NewShareClient(shareName).NewDirectoryClient(testcommon.GenerateDirectoryName(testName))

	_, err = dirClientAudience.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *ServiceRecordedTestsSuite) TestAccountPropertiesWithNFS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	resp, err := premiumShareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(resp.EnabledProtocols, to.Ptr("NFS"))
}
