//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share_test

import (
	"context"
	"fmt"
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
	if recordMode == recording.LiveMode {
		suite.Run(t, &ShareRecordedTestsSuite{})
		suite.Run(t, &ShareUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &ShareRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &ShareRecordedTestsSuite{})
	}
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
}

type ShareUnrecordedTestsSuite struct {
	suite.Suite
}

func (s *ShareUnrecordedTestsSuite) TestShareCreateRootDirectoryURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	rootDirClient := shareClient.NewRootDirectoryClient()
	_require.Equal(shareClient.URL(), rootDirClient.URL())
}

func (s *ShareUnrecordedTestsSuite) TestShareCreateDirectoryURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, err := testcommon.GetRequiredEnv(testcommon.AccountNameEnvVar)
	_require.NoError(err)

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

func (s *ShareUnrecordedTestsSuite) TestShareCreateUsingSharedKey() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName
	shareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	_require.NoError(err)

	resp, err := shareClient.Create(context.Background(), nil)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
}

func (s *ShareUnrecordedTestsSuite) TestShareCreateUsingConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	connString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient, err := share.NewClientFromConnectionString(*connString, shareName, nil)
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

	// TODO: create directories and files and uncomment this
	//dirCtr, fileCtr := 0, 0
	//pager := shareSASClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(nil)
	//for pager.More() {
	//	resp, err := pager.NextPage(context.Background())
	//	_require.NoError(err)
	//	dirCtr += len(resp.Segment.Directories)
	//	fileCtr += len(resp.Segment.Files)
	//}
	//_require.Equal(dirCtr, 0)
	//_require.Equal(fileCtr, 0)
}

func (s *ShareUnrecordedTestsSuite) TestShareCreateDeleteNonDefault() {
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

func (s *ShareUnrecordedTestsSuite) TestShareCreateNilMetadata() {
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

func (s *ShareUnrecordedTestsSuite) TestShareCreateNegativeInvalidName() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := svcClient.NewShareClient("foo bar")

	_, err = shareClient.Create(context.Background(), nil)

	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidResourceName)
}

func (s *ShareUnrecordedTestsSuite) TestShareCreateNegativeInvalidMetadata() {
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
	_require.NotNil(err)
}

func (s *ShareUnrecordedTestsSuite) TestShareDeleteNegativeNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	_, err = shareClient.Delete(context.Background(), nil)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareUnrecordedTestsSuite) TestShareGetSetPropertiesNonDefault() {
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

func (s *ShareUnrecordedTestsSuite) TestShareGetSetPropertiesDefault() {
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

func (s *ShareUnrecordedTestsSuite) TestShareSetQuotaNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.SetProperties(context.Background(), &share.SetPropertiesOptions{Quota: to.Ptr(int32(-1))})
	_require.NotNil(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidHeaderValue)
}

func (s *ShareUnrecordedTestsSuite) TestShareGetPropertiesNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.GetProperties(context.Background(), nil)
	_require.NotNil(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareUnrecordedTestsSuite) TestSharePutAndGetPermission() {
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

func (s *ShareUnrecordedTestsSuite) TestShareGetSetAccessPolicyNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	now := time.Now().UTC().Truncate(10000 * time.Millisecond) // Enough resolution
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

func (s *ShareUnrecordedTestsSuite) TestShareGetSetAccessPolicyNonDefaultMultiple() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	now := time.Now().UTC().Truncate(10000 * time.Millisecond) // Enough resolution
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

func (s *ShareUnrecordedTestsSuite) TestShareSetAccessPolicyMoreThanFive() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	now := time.Now().UTC().Truncate(10000 * time.Millisecond) // Enough resolution
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

func (s *ShareUnrecordedTestsSuite) TestShareGetSetAccessPolicyDefault() {
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

func (s *ShareUnrecordedTestsSuite) TestShareGetAccessPolicyNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.GetShareClient(shareName, svcClient)

	_, err = shareClient.GetAccessPolicy(context.Background(), nil)
	_require.NotNil(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ShareNotFound)
}

func (s *ShareUnrecordedTestsSuite) TestShareSetAccessPolicyNonDefaultDeleteAndModifyACL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	start := time.Now().UTC().Truncate(10000 * time.Millisecond)
	expiry := start.Add(5 * time.Minute).UTC()
	accessPermission := share.AccessPolicyPermission{List: true}.String()
	permissions := make([]*share.SignedIdentifier, 2, 2)
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

func (s *ShareUnrecordedTestsSuite) TestShareSetAccessPolicyDeleteAllPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	start := time.Now().UTC()
	expiry := start.Add(5 * time.Minute).UTC()
	accessPermission := share.AccessPolicyPermission{List: true}.String()
	permissions := make([]*share.SignedIdentifier, 2, 2)
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
