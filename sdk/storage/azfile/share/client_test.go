// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running share Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &ShareRecordedTestsSuite{})
		suite.Run(t, &ShareUnrecordedTestsSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &ShareRecordedTestsSuite{})
	case recording.RecordingMode:
		suite.Run(t, &ShareRecordedTestsSuite{})
	}
}

func (s *ShareRecordedTestsSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *ShareRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
}

func (s *ShareRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *ShareRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *ShareUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *ShareUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type ShareRecordedTestsSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type ShareUnrecordedTestsSuite struct {
	suite.Suite
}

func (s *ShareRecordedTestsSuite) TestShareCreateRootDirectoryURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	rootDirClient := shareClient.NewRootDirectoryClient()
	_require.Equal(shareClient.URL(), rootDirClient.URL())
}

func (s *ShareRecordedTestsSuite) TestShareCreateDirectoryURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName, dirName := testcommon.GenerateShareName(testName), testcommon.GenerateDirectoryName(testName)
	shareClient := svcClient.NewShareClient(shareName)
	_require.NoError(err)
	dirClient := shareClient.NewDirectoryClient(dirName)
	_require.NoError(err)

	correctURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName
	_require.Equal(dirClient.URL(), correctURL)
}

func (s *ShareRecordedTestsSuite) TestShareCreateUsingSharedKey() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName
	options := &share.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	resp, err := shareClient.Create(context.Background(), nil)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
}

func (s *ShareRecordedTestsSuite) TestShareCreateUsingConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	connString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	options := &share.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClient, err := share.NewClientFromConnectionString(*connString, shareName, options)
	_require.NoError(err)

	resp, err := shareClient.Create(context.Background(), nil)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
}

func (s *ShareUnrecordedTestsSuite) TestShareClientUsingSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	permissions := sas.SharePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)

	shareSASURL, err := shareClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	shareSASClient, err := share.NewClientWithNoCredential(shareSASURL, nil)
	_require.NoError(err)

	_, err = shareSASClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.AuthorizationFailure)

	dirName1 := testcommon.GenerateDirectoryName(testName) + "1"
	_, err = shareSASClient.NewDirectoryClient(dirName1).Create(context.Background(), nil)
	_require.NoError(err)

	dirName2 := testcommon.GenerateDirectoryName(testName) + "2"
	_, err = shareSASClient.NewDirectoryClient(dirName2).Create(context.Background(), nil)
	_require.NoError(err)

	fileName1 := testcommon.GenerateFileName(testName) + "1"
	_, err = shareSASClient.NewRootDirectoryClient().NewFileClient(fileName1).Create(context.Background(), 1024, nil)
	_require.NoError(err)

	fileName2 := testcommon.GenerateFileName(testName) + "2"
	_, err = shareSASClient.NewDirectoryClient(dirName2).NewFileClient(fileName2).Create(context.Background(), 1024, nil)
	_require.NoError(err)

	dirCtr, fileCtr := 0, 0
	pager := shareSASClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		dirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
	}
	_require.Equal(dirCtr, 2)
	_require.Equal(fileCtr, 1)
}

func (s *ShareUnrecordedTestsSuite) TestShareCreateDeleteUsingOAuth() {
	_require := require.New(s.T())
	testName := s.T().Name()
	shareName := testcommon.GenerateShareName(testName)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	shareClientOAuth, err := share.NewClient("https://"+accountName+".file.core.windows.net/"+shareName, cred, nil)
	_require.NoError(err)

	_, err = shareClientOAuth.Delete(context.Background(), nil)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)

}

func (s *ShareRecordedTestsSuite) TestShareCreateDeleteNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	quota := int32(1000)

	cResp, err := shareClient.Create(context.Background(), &share.CreateOptions{
		AccessTier: to.Ptr(share.AccessTierCool),
		Quota:      to.Ptr(quota),
		Metadata:   testcommon.BasicMetadata})

	_require.NoError(err)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.ETag)
	_require.NotNil(cResp.LastModified)
	_require.NotNil(cResp.RequestID)
	_require.NotNil(cResp.Version)

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix:  to.Ptr(shareName),
		Include: service.ListSharesInclude{Metadata: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Shares, 1)
		_require.Equal(*resp.Shares[0].Name, shareName)
		_require.NotNil(resp.Shares[0].Metadata)
		_require.EqualValues(resp.Shares[0].Metadata, testcommon.BasicMetadata)
		_require.Equal(*resp.Shares[0].Properties.AccessTier, string(share.AccessTierCool))
		_require.Equal(*resp.Shares[0].Properties.Quota, quota)
	}

	dResp, err := shareClient.Delete(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(dResp.Date)
	_require.NotNil(dResp.RequestID)
	_require.NotNil(dResp.Version)

	pager1 := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix:  to.Ptr(shareName),
		Include: service.ListSharesInclude{Metadata: true},
	})
	for pager1.More() {
		resp, err := pager1.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Shares, 0)
	}
}

func (s *ShareRecordedTestsSuite) TestShareCreateNilMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Create(context.Background(), nil)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)
	_require.NoError(err)

	response, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(response.Metadata, 0)
}

func (s *ShareRecordedTestsSuite) TestShareCreateAccessTierPremium() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{
		AccessTier: to.Ptr(share.AccessTierPremium),
	})
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)
	_require.NoError(err)
	response, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*response.AccessTier, string(share.AccessTierPremium))
}

func (s *ShareRecordedTestsSuite) TestShareCreateWithSnapshotVirtualDirectoryAccess() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{EnabledProtocols: to.Ptr("NFS"), EnableSnapshotVirtualDirectoryAccess: to.Ptr(false)})
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)
	_require.NoError(err)

	response, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(response.EnableSnapshotVirtualDirectoryAccess, to.Ptr(false))
}

func (s *ShareRecordedTestsSuite) TestShareCreateWithSnapshotVirtualDirectoryAccessDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{EnabledProtocols: to.Ptr("NFS"), EnableSnapshotVirtualDirectoryAccess: to.Ptr(true)})
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)
	_require.NoError(err)

	response, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(response.EnableSnapshotVirtualDirectoryAccess, to.Ptr(true))
}

func (s *ShareRecordedTestsSuite) TestShareCreatePaidBursting() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName
	options := &share.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	createOptions := &share.CreateOptions{
		PaidBurstingEnabled:           to.Ptr(true),
		PaidBurstingMaxIops:           to.Ptr(int64(5000)),
		PaidBurstingMaxBandwidthMibps: to.Ptr(int64(1000)),
	}

	resp, err := shareClient.Create(context.Background(), createOptions)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)

	props, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(props.ETag)
	_require.Equal(props.LastModified.IsZero(), false)
	_require.NotNil(props.RequestID)
	_require.NotNil(props.Version)
	_require.Equal(props.PaidBurstingEnabled, to.Ptr(true))
	_require.Equal(*props.PaidBurstingMaxIops, int64(5000))
	_require.Equal(*props.PaidBurstingMaxBandwidthMibps, int64(1000))
}

func (s *ShareUnrecordedTestsSuite) TestAuthenticationErrorDetailError() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	id := "testAccessPolicy"
	ps := share.AccessPolicyPermission{
		Write: true,
	}
	_require.NoError(err)

	signedIdentifiers := make([]*share.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &share.SignedIdentifier{
		AccessPolicy: &share.AccessPolicy{
			Expiry:     to.Ptr(time.Now().Add(-1 * time.Hour)),
			Permission: to.Ptr(ps.String()),
		},
		ID: &id,
	})
	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: signedIdentifiers,
	})
	_require.NoError(err)

	sasQueryParams, err := sas.SignatureValues{
		Protocol:   sas.ProtocolHTTPS,
		Identifier: id,
		ShareName:  shareName,
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	shareSAS := shareClient.URL() + "?" + sasQueryParams.Encode()
	shareClientSAS, err := share.NewClientWithNoCredential(shareSAS, nil)
	_require.NoError(err)

	dirClient := testcommon.GetDirectoryClient(testcommon.GenerateDirectoryName(testName), shareClientSAS)
	_, err = dirClient.Create(context.Background(), nil)
	_require.Error(err)

	var responseErr *azcore.ResponseError
	_require.ErrorAs(err, &responseErr)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.AuthenticationFailed)
	_require.Contains(responseErr.Error(), "AuthenticationErrorDetail")
}

func (s *ShareRecordedTestsSuite) TestShareCreateNegativeInvalidName() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := svcClient.NewShareClient("foo bar")

	_, err = shareClient.Create(context.Background(), nil)

	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidResourceName)
}

func (s *ShareRecordedTestsSuite) TestShareCreateNegativeInvalidMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{
		Metadata: map[string]*string{"!@#$%^&*()": to.Ptr("!@#$%^&*()")},
		Quota:    to.Ptr(int32(0)),
	})
	_require.Error(err)
}

func (s *ShareRecordedTestsSuite) TestShareCreateWithSMBDirectoryLeaseDisabled() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName
	options := &share.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// Create share with directory leases disabled
	resp, err := shareClient.Create(context.Background(), &share.CreateOptions{
		EnableSMBDirectoryLease: to.Ptr(false),
	})
	_require.NoError(err)
	_require.NotNil(resp.ETag)

	// Verify with GetProperties
	getResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.LastModified)
	_require.Equal(false, *getResp.EnableSMBDirectoryLease)
}

func (s *ShareRecordedTestsSuite) TestShareCreateWithSMBDirectoryLeaseEnabled() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName
	options := &share.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// Create share with directory leases explicitly enabled
	resp, err := shareClient.Create(context.Background(), &share.CreateOptions{
		EnableSMBDirectoryLease: to.Ptr(true),
	})
	_require.NoError(err)
	_require.NotNil(resp.ETag)

	// Verify with GetProperties
	getResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(true, *getResp.EnableSMBDirectoryLease)
}

func (s *ShareRecordedTestsSuite) TestShareCreateWithSMBDirectoryLeaseDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName
	options := &share.ClientOptions{}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// Create share with no explicit option (should default to true)
	resp, err := shareClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)

	// Verify with GetProperties
	getResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(true, *getResp.EnableSMBDirectoryLease)
}

func (s *ShareRecordedTestsSuite) TestShareDeleteNegativeNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Delete(context.Background(), nil)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareRecordedTestsSuite) TestShareGetSetPropertiesNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	newQuota := int32(1234)

	sResp, err := shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		Quota:      to.Ptr(newQuota),
		AccessTier: to.Ptr(share.AccessTierHot),
	})
	_require.NoError(err)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)
	_require.Equal(sResp.Date.IsZero(), false)

	props, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(props.ETag)
	_require.Equal(props.LastModified.IsZero(), false)
	_require.NotNil(props.RequestID)
	_require.NotNil(props.Version)
	_require.Equal(props.Date.IsZero(), false)
	_require.Equal(*props.Quota, newQuota)
	_require.Equal(*props.AccessTier, string(share.AccessTierHot))
}

func (s *ShareRecordedTestsSuite) TestShareGetSetPropertiesDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	sResp, err := shareClient.SetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)
	_require.Equal(sResp.Date.IsZero(), false)

	props, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(props.ETag)
	_require.Equal(props.LastModified.IsZero(), false)
	_require.NotNil(props.RequestID)
	_require.NotNil(props.Version)
	_require.Equal(props.Date.IsZero(), false)
	_require.Greater(*props.Quota, int32(0)) // When using service default quota, it could be any value
}

func (s *ShareRecordedTestsSuite) TestShareGetSetPropertiesAccessTierPremium() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		AccessTier: to.Ptr(share.AccessTierPremium),
	})
	_require.NoError(err)

	props, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*props.AccessTier, string(share.AccessTierPremium))
}

func (s *ShareUnrecordedTestsSuite) TestShareGetSetPropertiesOAuth() {
	_require := require.New(s.T())
	testName := s.T().Name()
	shareName := testcommon.GenerateShareName(testName)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	shareClientOAuth, err := share.NewClient("https://"+accountName+".file.core.windows.net/"+shareName, cred, nil)
	_require.NoError(err)

	_, err = shareClientOAuth.Create(context.Background(), nil)
	_require.NoError(err)

	sResp, err := shareClientOAuth.SetProperties(context.Background(), &share.SetPropertiesOptions{
		AccessTier: to.Ptr(share.AccessTierCool),
	})
	_require.NoError(err)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)
	_require.Equal(sResp.Date.IsZero(), false)

	properties, err := shareClientOAuth.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(properties.ETag)
	_require.Equal(*properties.AccessTier, string(share.AccessTierCool))

}

func (s *ShareRecordedTestsSuite) TestShareGetSetPropertiesWithSnapshotVirtualDirectory() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{EnabledProtocols: to.Ptr("NFS")})
	_require.NoError(err)

	defer testcommon.DeleteShare(context.Background(), _require, shareClient)
	_require.NoError(err)

	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		EnableSnapshotVirtualDirectoryAccess: to.Ptr(true),
	})
	_require.NoError(err)

	props, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*props.EnableSnapshotVirtualDirectoryAccess, true)
}

func (s *ShareRecordedTestsSuite) TestShareSetPropertiesPaidBursting() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	sResp, err := shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{
		PaidBurstingEnabled:           to.Ptr(true),
		PaidBurstingMaxIops:           to.Ptr(int64(5000)),
		PaidBurstingMaxBandwidthMibps: to.Ptr(int64(1000)),
	})
	_require.NoError(err)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)
	_require.Equal(sResp.Date.IsZero(), false)

	props, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(props.ETag)
	_require.Equal(props.LastModified.IsZero(), false)
	_require.NotNil(props.RequestID)
	_require.NotNil(props.Version)
	_require.Equal(props.Date.IsZero(), false)
	_require.Equal(props.PaidBurstingEnabled, to.Ptr(true))
	_require.Equal(*props.PaidBurstingMaxIops, int64(5000))
	_require.Equal(*props.PaidBurstingMaxBandwidthMibps, int64(1000))
}

func (s *ShareRecordedTestsSuite) TestShareSetQuotaNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{Quota: to.Ptr(int32(-1))})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidHeaderValue)
}

func (s *ShareRecordedTestsSuite) TestShareGetPropertiesNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareRecordedTestsSuite) TestSharePutAndGetPermission() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// Create a permission and check that it's not empty.
	createResp, err := shareClient.CreatePermission(context.Background(), testcommon.SampleSDDL, nil)
	_require.NoError(err)
	_require.NotEqual(*createResp.FilePermissionKey, "")

	getResp, err := shareClient.GetPermission(context.Background(), *createResp.FilePermissionKey, nil)
	_require.NoError(err)
	// Rather than checking against the original, we check for emptiness, as Azure Files has set a nil-ness flag on SACLs
	//        and converted our well-known SID.
	/*
		Expected :string = "O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)"
		Actual   :string = "O:AOG:S-1-5-21-397955417-626881126-188441444-512D:(A;;CCDCLCSWRPWPRCWDWOGA;;;S-1-0-0)S:NO_ACCESS_CONTROL"
	*/
	_require.NotNil(getResp.Permission)
	_require.NotEmpty(*getResp.Permission)
}

func (s *ShareRecordedTestsSuite) TestShareGetSetAccessPolicyNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)
	now := currTime.UTC().Truncate(10000 * time.Millisecond) // Enough resolution
	expiryTIme := now.Add(5 * time.Minute).UTC()
	pS := share.AccessPolicyPermission{
		Read:   true,
		Write:  true,
		Create: true,
		Delete: true,
		List:   true,
	}
	pS2 := &share.AccessPolicyPermission{}
	err = pS2.Parse("ldcwr")
	_require.NoError(err)
	_require.EqualValues(*pS2, pS)

	permission := pS.String()
	permissions := []*share.SignedIdentifier{
		{
			ID: to.Ptr("MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI="),
			AccessPolicy: &share.AccessPolicy{
				Start:      &now,
				Expiry:     &expiryTIme,
				Permission: &permission,
			},
		}}

	sResp, err := shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.NoError(err)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)

	gResp, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.Len(gResp.SignedIdentifiers, 1)
	_require.EqualValues(*(gResp.SignedIdentifiers[0]), *permissions[0])
}

func (s *ShareRecordedTestsSuite) TestShareGetSetAccessPolicyNonDefaultMultiple() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)
	now := currTime.UTC().Truncate(10000 * time.Millisecond) // Enough resolution
	expiryTIme := now.Add(5 * time.Minute).UTC()
	permission := share.AccessPolicyPermission{
		Read:  true,
		Write: true,
	}.String()

	permissions := []*share.SignedIdentifier{
		{
			ID: to.Ptr("MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI="),
			AccessPolicy: &share.AccessPolicy{
				Start:      &now,
				Expiry:     &expiryTIme,
				Permission: &permission,
			},
		},
		{
			ID: to.Ptr("2"),
			AccessPolicy: &share.AccessPolicy{
				Start:      &now,
				Expiry:     &expiryTIme,
				Permission: &permission,
			},
		}}

	sResp, err := shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.NoError(err)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)

	gResp, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.Len(gResp.SignedIdentifiers, 2)
	_require.EqualValues(gResp.SignedIdentifiers[0], permissions[0])
	_require.EqualValues(gResp.SignedIdentifiers[1], permissions[1])
}

func (s *ShareRecordedTestsSuite) TestShareSetAccessPolicyMoreThanFive() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)
	now := currTime.UTC().Truncate(10000 * time.Millisecond) // Enough resolution
	expiryTIme := now.Add(5 * time.Minute).UTC()
	permission := share.AccessPolicyPermission{
		Read:   true,
		Create: true,
		Write:  true,
		Delete: true,
		List:   true,
	}.String()

	var permissions []*share.SignedIdentifier
	for i := 0; i <= len(permission); i++ {
		p := permission
		if i < len(permission) {
			p = string(permission[i])
		}
		permissions = append(permissions, &share.SignedIdentifier{
			ID: to.Ptr(fmt.Sprintf("%v", i)),
			AccessPolicy: &share.AccessPolicy{
				Start:      &now,
				Expiry:     &expiryTIme,
				Permission: &p,
			},
		})
	}
	_require.Len(permissions, 6)

	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidXMLDocument)
}

func (s *ShareRecordedTestsSuite) TestShareGetSetAccessPolicyDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	sResp, err := shareClient.SetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)

	gResp, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.Len(gResp.SignedIdentifiers, 0)
}

func (s *ShareRecordedTestsSuite) TestShareGetAccessPolicyNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.GetAccessPolicy(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareRecordedTestsSuite) TestShareSetAccessPolicyNonDefaultDeleteAndModifyACL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	currTime, err := time.Parse(time.UnixDate, "Thu Mar 30 20:00:00 GMT 2023")
	_require.NoError(err)
	start := currTime.UTC().Truncate(10000 * time.Millisecond)
	expiry := start.Add(5 * time.Minute).UTC()
	accessPermission := share.AccessPolicyPermission{List: true}.String()
	permissions := make([]*share.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &share.SignedIdentifier{
			ID: to.Ptr("000" + strconv.Itoa(i)),
			AccessPolicy: &share.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.NoError(err)

	resp, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the second policy by removing it from the slice
	permissions[0].ID = to.Ptr("0004")       // Modify the remaining policy which is at index 0 in the new slice
	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.NoError(err)

	resp, err = shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.SignedIdentifiers, 1)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *ShareRecordedTestsSuite) TestShareSetAccessPolicyDeleteAllPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)
	start := currTime.UTC()
	expiry := start.Add(5 * time.Minute).UTC()
	accessPermission := share.AccessPolicyPermission{List: true}.String()
	permissions := make([]*share.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &share.SignedIdentifier{
			ID: to.Ptr("000" + strconv.Itoa(i)),
			AccessPolicy: &share.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.NoError(err)

	resp1, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp1.SignedIdentifiers, 2)

	_, err = shareClient.SetAccessPolicy(context.Background(), nil)
	_require.NoError(err)

	resp2, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp2.SignedIdentifiers, 0)
}

func (s *ShareRecordedTestsSuite) TestShareSetPermissionsNegativeInvalidPolicyTimes() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// Swap start and expiry
	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)
	expiry := currTime.UTC()
	start := expiry.Add(5 * time.Minute).UTC()
	accessPermission := share.AccessPolicyPermission{List: true}.String()
	permissions := make([]*share.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &share.SignedIdentifier{
			ID: to.Ptr("000" + strconv.Itoa(i)),
			AccessPolicy: &share.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.NoError(err)

	resp, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

// SignedIdentifier ID too long
func (s *ShareRecordedTestsSuite) TestShareSetPermissionsNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	currTime, err := time.Parse(time.UnixDate, "Wed Mar 29 20:00:00 GMT 2023")
	_require.NoError(err)
	expiry := currTime.UTC()
	start := expiry.Add(5 * time.Minute).UTC()
	accessPermission := share.AccessPolicyPermission{List: true}.String()
	permissions := make([]*share.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &share.SignedIdentifier{
			ID: to.Ptr(id),
			AccessPolicy: &share.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: permissions,
	})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidXMLDocument)
}

func (s *ShareRecordedTestsSuite) TestShareGetSetMetadataDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	sResp, err := shareClient.SetMetadata(context.Background(), &share.SetMetadataOptions{
		Metadata: map[string]*string{},
	})
	_require.NoError(err)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)

	gResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.Len(gResp.Metadata, 0)
}

func (s *ShareRecordedTestsSuite) TestShareGetSetMetadataNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}
	sResp, err := shareClient.SetMetadata(context.Background(), &share.SetMetadataOptions{
		Metadata: md,
	})
	_require.NoError(err)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotNil(sResp.ETag)
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)

	gResp, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.EqualValues(gResp.Metadata, md)
}

func (s *ShareRecordedTestsSuite) TestShareSetMetadataNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	md := map[string]*string{
		"!@#$%^&*()": to.Ptr("!@#$%^&*()"),
	}
	_, err = shareClient.SetMetadata(context.Background(), &share.SetMetadataOptions{
		Metadata: md,
	})
	_require.Error(err)
}

func (s *ShareRecordedTestsSuite) TestShareGetStats() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	newQuota := int32(300)

	// In order to test and get LastModified property.
	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{Quota: to.Ptr(newQuota)})
	_require.NoError(err)

	gResp, err := shareClient.GetStatistics(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	// _require.NotEqual(*gResp.ETag, "") // TODO: The ETag would be ""
	// _require.Equal(gResp.LastModified.IsZero(), false) // TODO: Even share is once updated, no LastModified would be returned.
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.Equal(*gResp.ShareUsageBytes, int64(0))
}

func (s *ShareRecordedTestsSuite) TestShareGetStatsNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.GetStatistics(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareRecordedTestsSuite) TestSetAndGetStatistics() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.Create(context.Background(), &share.CreateOptions{Quota: to.Ptr(int32(1024))})
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := shareClient.NewDirectoryClient("testdir")
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileClient := dirClient.NewFileClient("testfile")
	_, err = fileClient.Create(context.Background(), int64(1024*1024*1024*1024), nil)
	_require.NoError(err)

	getStats, err := shareClient.GetStatistics(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getStats.ShareUsageBytes, int64(1024*1024*1024*1024))
}

func deleteShare(ctx context.Context, _require *require.Assertions, shareClient *share.Client, o *share.DeleteOptions) {
	_, err := shareClient.Delete(ctx, o)
	_require.NoError(err)
}

func (s *ShareRecordedTestsSuite) TestShareCreateSnapshotNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer deleteShare(context.Background(), _require, shareClient, &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})

	cResp, err := shareClient.CreateSnapshot(context.Background(), &share.CreateSnapshotOptions{Metadata: testcommon.BasicMetadata})
	_require.NoError(err)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.ETag)
	_require.NotEqual(*cResp.ETag, "")
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotNil(cResp.RequestID)
	_require.NotNil(cResp.Version)
	_require.NotNil(cResp.Snapshot)
	_require.NotEqual(*cResp.Snapshot, "")

	cSnapshot := *cResp.Snapshot

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Metadata: true, Snapshots: true},
		Prefix:  &shareName,
	})

	foundSnapshot := false
	for pager.More() {
		lResp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(lResp.Shares, 2)

		for _, s := range lResp.Shares {
			if s.Snapshot != nil {
				foundSnapshot = true
				_require.Equal(*s.Snapshot, cSnapshot)
				_require.NotNil(s.Metadata)
				_require.EqualValues(s.Metadata, testcommon.BasicMetadata)
			} else {
				_require.Len(s.Metadata, 0)
			}
		}
	}
	_require.True(foundSnapshot)
}

func (s *ShareUnrecordedTestsSuite) TestShareCreateSnapshotDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer deleteShare(context.Background(), _require, shareClient, &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})

	// create a file in the base share.
	dirClient := shareClient.NewRootDirectoryClient()
	_require.NoError(err)

	fClient := dirClient.NewFileClient("myfile")
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	// Create share snapshot, the snapshot contains the create file.
	snapshotShare, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)

	// Delete file in base share.
	_, err = fClient.Delete(context.Background(), nil)
	_require.NoError(err)

	// To produce a share SAS (as opposed to a file SAS), assign to FilePermissions using
	// ShareSASPermissions and make sure the DirectoryAndFilePath field is "" (the default).
	perms := sas.SharePermissions{Read: true, Write: true}

	// Restore file from share snapshot.
	// Create a SAS.
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ShareName:   shareName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	// Build a file snapshot URL.
	fileParts, err := sas.ParseURL(fClient.URL())
	_require.NoError(err)
	fileParts.ShareSnapshot = *snapshotShare.Snapshot
	fileParts.SAS = sasQueryParams
	sourceURL := fileParts.String()

	// Before restore
	_, err = fClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)

	// Do restore.
	_, err = fClient.StartCopyFromURL(context.Background(), sourceURL, nil)
	_require.NoError(err)

	time.Sleep(2 * time.Second)

	// After restore
	_, err = fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = shareClient.Delete(context.Background(), &share.DeleteOptions{
		ShareSnapshot: snapshotShare.Snapshot,
	})
	_require.NoError(err)
}

func (s *ShareRecordedTestsSuite) TestShareCreateSnapshotNegativeShareNotExist() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.CreateSnapshot(context.Background(), &share.CreateSnapshotOptions{Metadata: map[string]*string{}})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareRecordedTestsSuite) TestShareDeleteSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer deleteShare(context.Background(), _require, shareClient, &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})

	resp1, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp1.Snapshot)
	_require.NotEmpty(*resp1.Snapshot)

	resp2, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp2.Snapshot)
	_require.NotEmpty(*resp2.Snapshot)

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Snapshots: true},
		Prefix:  &shareName,
	})

	snapshotsCtr := 0
	for pager.More() {
		lResp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(lResp.Shares, 3) // 2 snapshots and 1 share

		for _, s := range lResp.Shares {
			if s.Snapshot != nil {
				snapshotsCtr++
			}
		}
	}
	_require.Equal(snapshotsCtr, 2)

	snapClient, err := shareClient.WithSnapshot(*resp1.Snapshot)
	_require.NoError(err)

	_, err = snapClient.Delete(context.Background(), nil)
	_require.NoError(err)

	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Snapshots: true},
		Prefix:  &shareName,
	})

	snapshotsCtr = 0
	for pager.More() {
		lResp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(lResp.Shares, 2)

		for _, s := range lResp.Shares {
			if s.Snapshot != nil {
				snapshotsCtr++
				_require.Equal(*s.Snapshot, *resp2.Snapshot)
			}
		}
	}
	_require.Equal(snapshotsCtr, 1)
}

func (s *ShareRecordedTestsSuite) TestShareCreateSnapshotNegativeMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.CreateSnapshot(context.Background(), &share.CreateSnapshotOptions{Metadata: map[string]*string{"!@#$%^&*()": to.Ptr("!@#$%^&*()")}})
	_require.Error(err)
}

func (s *ShareRecordedTestsSuite) TestShareCreateSnapshotNegativeSnapshotOfSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer deleteShare(context.Background(), _require, shareClient, &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})

	snapTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)

	snapshotClient, err := shareClient.WithSnapshot(snapTime.UTC().String())
	_require.NoError(err)

	cResp, err := snapshotClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err) //Note: this would not fail, snapshot would be ignored.
	_require.NotNil(cResp)
	_require.NotEmpty(*cResp.Snapshot)

	snapshotRecursiveClient, err := shareClient.WithSnapshot(*cResp.Snapshot)
	_require.NoError(err)
	_, err = snapshotRecursiveClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err) //Note: this would not fail, snapshot would be ignored.
}

func (s *ShareRecordedTestsSuite) TestShareDeleteSnapshotsInclude() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)

	_, err = shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Snapshots: true},
		Prefix:  &shareName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Shares, 2)
	}

	_, err = shareClient.Delete(context.Background(), &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})
	_require.NoError(err)

	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Snapshots: true},
		Prefix:  &shareName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Shares, 0)
	}
}

func (s *ShareRecordedTestsSuite) TestShareDeleteSnapshotsNoneWithSnapshots() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer deleteShare(context.Background(), _require, shareClient, &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})

	_, err = shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)

	_, err = shareClient.Delete(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareHasSnapshots)
}

func (s *ShareRecordedTestsSuite) TestShareRestoreSuccess() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.Delete(context.Background(), nil)
	_require.NoError(err)

	// wait for share deletion
	time.Sleep(60 * time.Second)

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Deleted: true},
		Prefix:  &shareName,
	})

	shareVersion := ""
	shareCtr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, s := range resp.Shares {
			if s.Deleted != nil && *s.Deleted {
				shareVersion = *s.Version
			} else {
				shareCtr++
			}
		}
	}
	_require.NotEmpty(shareVersion)
	_require.Equal(shareCtr, 0)

	rResp, err := shareClient.Restore(context.Background(), shareVersion, nil)
	_require.NoError(err)
	_require.NotNil(rResp.ETag)
	_require.NotNil(rResp.RequestID)
	_require.NotNil(rResp.Version)

	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Prefix: &shareName,
	})

	shareCtr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		shareCtr += len(resp.Shares)
	}
	_require.Equal(shareCtr, 1)
}

func (s *ShareRecordedTestsSuite) TestShareRestoreFailures() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.Restore(context.Background(), "", nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.MissingRequiredHeader)
}

func (s *ShareRecordedTestsSuite) TestShareRestoreWithSnapshotsAgain() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer deleteShare(context.Background(), _require, shareClient, &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})

	cResp, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(cResp.Snapshot)

	_, err = shareClient.Delete(context.Background(), &share.DeleteOptions{
		DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude),
	})
	_require.NoError(err)

	// wait for share deletion
	time.Sleep(60 * time.Second)

	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Deleted: true},
		Prefix:  &shareName,
	})

	shareVersion := ""
	shareCtr := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, s := range resp.Shares {
			if s.Deleted != nil && *s.Deleted {
				shareVersion = *s.Version
			} else {
				shareCtr++
			}
		}
	}
	_require.NotEmpty(shareVersion)
	_require.Equal(shareCtr, 0)

	rResp, err := shareClient.Restore(context.Background(), shareVersion, nil)
	_require.NoError(err)
	_require.NotNil(rResp.ETag)
	_require.NotNil(rResp.RequestID)
	_require.NotNil(rResp.Version)

	pager = svcClient.NewListSharesPager(&service.ListSharesOptions{
		Include: service.ListSharesInclude{Snapshots: true},
		Prefix:  &shareName,
	})

	shareCtr = 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		shareCtr += len(resp.Shares)
		for _, s := range resp.Shares {
			if s.Snapshot != nil {
				_require.Equal(*s.Snapshot, *cResp.Snapshot)
			}
		}
	}
	_require.Equal(shareCtr, 2) // 1 share and 1 snapshot
}

func (s *ShareRecordedTestsSuite) TestSASShareClientNoKey() {
	_require := require.New(s.T())
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	testName := s.T().Name()
	shareName := testcommon.GenerateShareName(testName)
	shareClient, err := share.NewClientWithNoCredential(fmt.Sprintf("https://%s.file.core.windows.net/%v", accountName, shareName), nil)
	_require.NoError(err)

	permissions := sas.SharePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)

	_, err = shareClient.GetSASURL(permissions, expiry, nil)
	_require.Equal(err, fileerror.MissingSharedKeyCredential)
}

func (s *ShareRecordedTestsSuite) TestSASShareClientSignNegative() {
	_require := require.New(s.T())
	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)
	_require.Greater(len(accountKey), 0)

	cred, err := share.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	testName := s.T().Name()
	shareName := testcommon.GenerateShareName(testName)
	shareClient, err := share.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.file.core.windows.net/%v", accountName, shareName), cred, nil)
	_require.NoError(err)

	permissions := sas.SharePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		List:   true,
		Create: true,
	}
	expiry := time.Time{}

	// zero expiry time
	_, err = shareClient.GetSASURL(permissions, expiry, &share.GetSASURLOptions{StartTime: to.Ptr(time.Now())})
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")

	// zero start and expiry time
	_, err = shareClient.GetSASURL(permissions, expiry, &share.GetSASURLOptions{})
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")

	// empty permissions
	_, err = shareClient.GetSASURL(sas.SharePermissions{}, expiry, nil)
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")
}

func (s *ShareRecordedTestsSuite) TestShareCreateAndGetPermissionOAuth() {
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

	options := &share.ClientOptions{FileRequestIntent: to.Ptr(share.TokenIntentBackup)}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClientOAuth, err := share.NewClient("https://"+accountName+".file.core.windows.net/"+shareName, cred, options)
	_require.NoError(err)

	// Create a permission and check that it's not empty.
	createResp, err := shareClientOAuth.CreatePermission(context.Background(), testcommon.SampleSDDL, nil)
	_require.NoError(err)
	_require.NotNil(createResp.FilePermissionKey)
	_require.NotEmpty(*createResp.FilePermissionKey)

	getResp, err := shareClientOAuth.GetPermission(context.Background(), *createResp.FilePermissionKey, nil)
	_require.NoError(err)
	_require.NotNil(getResp.Permission)
	_require.NotEmpty(*getResp.Permission)
}

func (s *ShareUnrecordedTestsSuite) TestShareSASUsingAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	id := "testAccessPolicy"
	ps := share.AccessPolicyPermission{
		Read:   true,
		Write:  true,
		Create: true,
		Delete: true,
		List:   true,
	}
	signedIdentifiers := make([]*share.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &share.SignedIdentifier{
		AccessPolicy: &share.AccessPolicy{
			Expiry:     to.Ptr(time.Now().Add(1 * time.Hour)),
			Start:      to.Ptr(time.Now()),
			Permission: to.Ptr(ps.String()),
		},
		ID: &id,
	})

	_, err = shareClient.SetAccessPolicy(context.Background(), &share.SetAccessPolicyOptions{
		ShareACL: signedIdentifiers,
	})
	_require.NoError(err)

	gResp, err := shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(gResp.SignedIdentifiers, 1)

	time.Sleep(30 * time.Second)

	sasQueryParams, err := sas.SignatureValues{
		Protocol:   sas.ProtocolHTTPS,
		Identifier: id,
		ShareName:  shareName,
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	shareSAS := shareClient.URL() + "?" + sasQueryParams.Encode()
	shareClientSAS, err := share.NewClientWithNoCredential(shareSAS, nil)
	_require.NoError(err)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClientSAS)
	fileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClientSAS)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)

	_, err = fileClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *ShareRecordedTestsSuite) TestPremiumShareBandwidth() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Create(context.Background(), nil)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)
	_require.NoError(err)

	response, err := shareClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(response.ProvisionedBandwidthMiBps)
	_require.NotNil(response.ProvisionedIngressMBps)
	_require.NotNil(response.ProvisionedEgressMBps)
	_require.NotNil(response.ProvisionedIops)
	_require.NotNil(response.NextAllowedQuotaDowngradeTime)
	_require.Greater(*response.ProvisionedBandwidthMiBps, (int32)(0))
}

func (s *ShareRecordedTestsSuite) TestShareClientDefaultAudience() {
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

	options := &share.ClientOptions{
		FileRequestIntent: to.Ptr(share.TokenIntentBackup),
		Audience:          "https://storage.azure.com/",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClientAudience, err := share.NewClient("https://"+accountName+".file.core.windows.net/"+shareName, cred, options)
	_require.NoError(err)

	// Create a permission and check that it's not empty.
	createResp, err := shareClientAudience.CreatePermission(context.Background(), testcommon.SampleSDDL, nil)
	_require.NoError(err)
	_require.NotNil(createResp.FilePermissionKey)
	_require.NotEmpty(*createResp.FilePermissionKey)

	getResp, err := shareClientAudience.GetPermission(context.Background(), *createResp.FilePermissionKey, nil)
	_require.NoError(err)
	_require.NotNil(getResp.Permission)
	_require.NotEmpty(*getResp.Permission)
}

func (s *ShareRecordedTestsSuite) TestShareClientCustomAudience() {
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

	options := &share.ClientOptions{
		FileRequestIntent: to.Ptr(share.TokenIntentBackup),
		Audience:          "https://" + accountName + ".file.core.windows.net",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClientAudience, err := share.NewClient("https://"+accountName+".file.core.windows.net/"+shareName, cred, options)
	_require.NoError(err)

	// Create a permission and check that it's not empty.
	createResp, err := shareClientAudience.CreatePermission(context.Background(), testcommon.SampleSDDL, nil)
	_require.NoError(err)
	_require.NotNil(createResp.FilePermissionKey)
	_require.NotEmpty(*createResp.FilePermissionKey)

	getResp, err := shareClientAudience.GetPermission(context.Background(), *createResp.FilePermissionKey, nil)
	_require.NoError(err)
	_require.NotNil(getResp.Permission)
	_require.NotEmpty(*getResp.Permission)
}

func (s *ShareUnrecordedTestsSuite) TestShareClientAudienceNegative() {
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

	options := &share.ClientOptions{
		FileRequestIntent: to.Ptr(share.TokenIntentBackup),
		Audience:          "https://badaudience.file.core.windows.net",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	shareClientAudience, err := share.NewClient("https://"+accountName+".file.core.windows.net/"+shareName, cred, options)
	_require.NoError(err)

	// Create a permission and check that it's not empty.
	_, err = shareClientAudience.CreatePermission(context.Background(), testcommon.SampleSDDL, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidAuthenticationInfo)
}
