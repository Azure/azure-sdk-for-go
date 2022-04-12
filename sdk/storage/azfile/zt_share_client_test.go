//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.18 See License.txt in the project root for license information.

package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"os"
	"strconv"
	"time"
)

func (s *azfileLiveTestSuite) TestShareCreateRootDirectoryURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, shareClient, nil)
}

func (s *azfileLiveTestSuite) TestPutAndGetPermission() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)
	_require.Nil(err)

	// Create a permission and check that it's not empty.
	createResp, err := srClient.CreatePermission(ctx, sampleSDDL, nil)
	_require.Nil(err)
	_require.NotEqual(*createResp.FilePermissionKey, "")

	getResp, err := srClient.GetPermission(ctx, *createResp.FilePermissionKey, nil)
	_require.Nil(err)
	// Rather than checking against the original, we check for emptiness, as Azure Files has set a nil-ness flag on SACLs
	//        and converted our well-known SID.
	/*
		Expected :string = "O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)"
		Actual   :string = "O:AOG:S-1-5-21-397955417-626881126-188441444-512D:(A;;CCDCLCSWRPWPRCWDWOGA;;;S-1-0-0)S:NO_ACCESS_CONTROL"
	*/
	_require.NotEqual(*getResp.Permission, "")
}

func (s *azfileLiveTestSuite) TestShareCreateDirectoryURL() {
	_require := require.New(s.T())
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient, err := svcClient.NewShareClient(sharePrefix)
	_require.Nil(err)
	dirClient, err := srClient.NewDirectoryClient(directoryPrefix)
	_require.Nil(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/" + sharePrefix + "/" + directoryPrefix
	_require.Equal(dirClient.URL(), correctURL)
}

// Note: test share create with default parameter is covered with preparing phase for FileURL and etc.
//func (s *azfileLiveTestSuite) TestShareCreateDeleteNonDefault() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient, err := svcClient.NewShareClient(shareName)
//	_require.Nil(err)
//
//	md := map[string]string{
//		"foo": "FooValuE",
//		"bar": "bArvaLue",
//	}
//
//	quota := int32(1000)
//
//	cResp, err := srClient.Create(context.Background(), &ShareCreateOptions{Quota: to.Int32Ptr(quota), Metadata: md})
//	_require.Nil(err)
//	_assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
//	_assert(cResp.Date().IsZero(), chk.Equals, false)
//	_assert(cResp.ETag(), chk.Not(chk.Equals), ETagNone)
//	_assert(cResp.LastModified.IsZero(), chk.Equals, false)
//	_assert(cResp.RequestID(), chk.Not(chk.Equals), "")
//	_assert(cResp.Version(), chk.Not(chk.Equals), "")
//
//	shares, err := srClient.ListSharesSegment(context.Background(), Marker{}, ListSharesOptions{Prefix: shareName, Detail: ListSharesDetail{Metadata: true}})
//	_require.Nil(err)
//	_assert(shares.ShareItems, chk.HasLen, 1)
//	_assert(shares.ShareItems[0].Name, chk.Equals, shareName)
//	_assert(shares.ShareItems[0].Metadata, chk.DeepEquals, md)
//	_assert(shares.ShareItems[0].Properties.Quota, chk.Equals, quota)
//
//	dResp, err := srClient.Delete(context.Background(), DeleteSnapshotsOptionNone)
//	_require.Nil(err)
//	_assert(dResp.RawResponse.StatusCode, chk.Equals, 202)
//	_assert(dResp.Date().IsZero(), chk.Equals, false)
//	_assert(dResp.RequestID(), chk.Not(chk.Equals), "")
//	_assert(dResp.Version(), chk.Not(chk.Equals), "")
//
//	shares, err = srClient.ListSharesSegment(context.Background(), Marker{}, ListSharesOptions{Prefix: shareName})
//	_require.Nil(err)
//	_assert(shares.ShareItems, chk.HasLen, 0)
//}

func (s *azfileLiveTestSuite) TestShareCreateNilMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	_, err = srClient.Create(ctx, nil)
	defer delShare(_require, srClient, nil)
	_require.Nil(err)

	response, err := srClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_requireLen(response.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestShareCreateNegativeInvalidName() {
	_require := require.New(s.T())
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient, err := svcClient.NewShareClient("foo bar")
	_require.Nil(err)

	_, err = srClient.Create(ctx, nil)

	validateStorageError(_require, err, StorageErrorCodeInvalidResourceName)
}

func (s *azfileLiveTestSuite) TestShareCreateNegativeInvalidMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	_, err = srClient.Create(ctx, &ShareCreateOptions{Metadata: map[string]string{"!@#$%^&*()": "!@#$%^&*()"}, Quota: to.Int32Ptr(0)})
	_require.NotNil(err)
}

func (s *azfileLiveTestSuite) TestShareDeleteNegativeNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	_, err = srClient.Delete(ctx, nil)
	validateStorageError(_require, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareGetSetPropertiesNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	newQuota := int32(1234)

	sResp, err := srClient.SetProperties(ctx, &ShareSetPropertiesOptions{Quota: to.Int32Ptr(newQuota)})
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.NotEqual(*sResp.ETag, "")
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotEqual(sResp.RequestID, "")
	_require.NotEqual(sResp.Version, "")
	_require.Equal(sResp.Date.IsZero(), false)

	props, err := srClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(props.RawResponse.StatusCode, 200)
	_require.NotEqual(*props.ETag, "")
	_require.Equal(props.LastModified.IsZero(), false)
	_require.NotEqual(*props.RequestID, "")
	_require.NotEqual(*props.Version, "")
	_require.Equal(props.Date.IsZero(), false)
	_require.Equal(*props.Quota, newQuota)
}

func (s *azfileLiveTestSuite) TestShareGetSetPropertiesDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	sResp, err := srClient.SetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.NotEqual(*sResp.ETag, "")
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotEqual(sResp.RequestID, "")
	_require.NotEqual(sResp.Version, "")
	_require.Equal(sResp.Date.IsZero(), false)

	props, err := srClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(props.RawResponse.StatusCode, 200)
	_require.NotEqual(*props.ETag, "")
	_require.Equal(props.LastModified.IsZero(), false)
	_require.NotEqual(*props.RequestID, "")
	_require.NotEqual(*props.Version, "")
	_require.Equal(props.Date.IsZero(), false)
	_requireTrue(*props.Quota >= 0) // When using service default quota, it could be any value
}

func (s *azfileLiveTestSuite) TestShareSetQuotaNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	_, err = srClient.SetProperties(ctx, &ShareSetPropertiesOptions{Quota: to.Int32Ptr(-1)})
	_require.NotNil(err)
	_requireContains(err.Error(), "validation failed: share quote cannot be negative")
}

func (s *azfileLiveTestSuite) TestShareGetPropertiesNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	_, err = srClient.GetProperties(ctx, nil)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareGetSetPermissionsNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	now := time.Now().UTC().Truncate(10000 * time.Millisecond) // Enough resolution
	expiryTIme := now.Add(5 * time.Minute).UTC()
	pS := AccessPolicyPermission{
		Read:   true,
		Write:  true,
		Create: true,
		Delete: true,
		List:   true,
	}
	pS2 := &AccessPolicyPermission{}
	pS2.Parse("ldcwr")
	_require.EqualValues(*pS2, pS)

	permission := pS.String()

	permissions := []*SignedIdentifier{
		{
			ID: to.StringPtr("MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI="),
			AccessPolicy: &AccessPolicy{
				Start:      &now,
				Expiry:     &expiryTIme,
				Permission: &permission,
			},
		}}

	sResp, err := srClient.SetPermissions(context.Background(), permissions, nil)
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotEqual(sResp.RequestID, "")
	_require.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetPermissions(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_requireLen(gResp.SignedIdentifiers, 1)
	_require.EqualValues(*(gResp.SignedIdentifiers[0]), *permissions[0])
}

func (s *azfileLiveTestSuite) TestShareGetSetPermissionsNonDefaultMultiple() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	now := time.Now().UTC().Truncate(10000 * time.Millisecond) // Enough resolution
	expiryTIme := now.Add(5 * time.Minute).UTC()
	permission := AccessPolicyPermission{
		Read:  true,
		Write: true,
	}.String()

	permissions := []*SignedIdentifier{
		{
			ID: to.StringPtr("MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI="),
			AccessPolicy: &AccessPolicy{
				Start:      &now,
				Expiry:     &expiryTIme,
				Permission: &permission,
			},
		},
		{
			ID: to.StringPtr("2"),
			AccessPolicy: &AccessPolicy{
				Start:      &now,
				Expiry:     &expiryTIme,
				Permission: &permission,
			},
		}}

	sResp, err := srClient.SetPermissions(context.Background(), permissions, nil)
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotEqual(sResp.RequestID, "")
	_require.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetPermissions(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_requireLen(gResp.SignedIdentifiers, 2)
	_require.EqualValues(gResp.SignedIdentifiers[0], permissions[0])
}

func (s *azfileLiveTestSuite) TestShareGetSetPermissionsDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	sResp, err := srClient.SetPermissions(context.Background(), []*SignedIdentifier{}, nil)
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotEqual(sResp.RequestID, "")
	_require.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetPermissions(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_requireLen(gResp.SignedIdentifiers, 0)
}

func (s *azfileLiveTestSuite) TestShareGetPermissionNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	_, err = srClient.GetPermissions(ctx, nil)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareSetPermissionsNonDefaultDeleteAndModifyACL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)

	defer delShare(_require, srClient, nil)

	start := time.Now().UTC().Truncate(10000 * time.Millisecond)
	expiry := start.Add(5 * time.Minute).UTC()
	accessPermission := AccessPolicyPermission{List: true}.String()
	permissions := make([]*SignedIdentifier, 2, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &SignedIdentifier{
			ID: to.StringPtr("000" + strconv.Itoa(i)),
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = srClient.SetPermissions(ctx, permissions, nil)
	_require.Nil(err)

	resp, err := srClient.GetPermissions(ctx, nil)
	_require.Nil(err)

	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	permissions[0].ID = to.StringPtr("0004") // Modify the remaining policy which is at index 0 in the new slice
	_, err = srClient.SetPermissions(ctx, permissions, nil)

	resp, err = srClient.GetPermissions(ctx, nil)
	_require.Nil(err)
	_requireLen(resp.SignedIdentifiers, 1)

	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *azfileLiveTestSuite) TestShareSetPermissionsDeleteAllPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)

	defer delShare(_require, srClient, nil)

	start := time.Now().UTC()
	expiry := start.Add(5 * time.Minute).UTC()
	accessPermission := AccessPolicyPermission{List: true}.String()
	permissions := make([]*SignedIdentifier, 2, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &SignedIdentifier{
			ID: to.StringPtr("000" + strconv.Itoa(i)),
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = srClient.SetPermissions(ctx, permissions, nil)
	_require.Nil(err)

	_, err = srClient.SetPermissions(ctx, []*SignedIdentifier{}, nil)
	_require.Nil(err)

	resp, err := srClient.GetPermissions(ctx, nil)
	_require.Nil(err)
	_requireLen(resp.SignedIdentifiers, 0)
}

// Note: No error happened
func (s *azfileLiveTestSuite) TestShareSetPermissionsNegativeInvalidPolicyTimes() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	// Swap start and expiry
	expiry := time.Now().UTC()
	start := expiry.Add(5 * time.Minute).UTC()
	accessPermission := AccessPolicyPermission{List: true}.String()
	permissions := make([]*SignedIdentifier, 2, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &SignedIdentifier{
			ID: to.StringPtr("000" + strconv.Itoa(i)),
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = srClient.SetPermissions(ctx, permissions, nil)
	_require.Nil(err)
}

func (s *azfileLiveTestSuite) TestShareSetPermissionsNilPolicySlice() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)

	defer delShare(_require, srClient, nil)

	_, err = srClient.SetPermissions(ctx, nil, nil)
	_require.Nil(err)
}

// SignedIdentifier ID too long
func (s *azfileLiveTestSuite) TestShareSetPermissionsNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)

	defer delShare(_require, srClient, nil)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry := time.Now().UTC()
	start := expiry.Add(5 * time.Minute).UTC()
	accessPermission := AccessPolicyPermission{List: true}.String()
	permissions := make([]*SignedIdentifier, 2, 2)
	for i := 0; i < 2; i++ {
		permissions[i] = &SignedIdentifier{
			ID: to.StringPtr(id),
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &accessPermission,
			},
		}
	}

	_, err = srClient.SetPermissions(ctx, permissions, nil)
	validateStorageError(_require, err, StorageErrorCodeInvalidXMLDocument)
}

func (s *azfileLiveTestSuite) TestShareGetSetMetadataDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	sResp, err := srClient.SetMetadata(context.Background(), map[string]string{}, nil)
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotEqual(sResp.RequestID, "")
	_require.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_requireLen(gResp.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestShareGetSetMetadataNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	md := map[string]string{
		"Foo": "FooValuE",
		"Bar": "bArvaLue", // Note: As testing result, currently only support case-insensitive keys(key will be saved in lower-case).
	}
	sResp, err := srClient.SetMetadata(context.Background(), md, nil)
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.Equal(sResp.LastModified.IsZero(), false)
	_require.NotEqual(sResp.RequestID, "")
	_require.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_require.EqualValues(gResp.Metadata, md)
}

func (s *azfileLiveTestSuite) TestShareSetMetadataNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	md := map[string]string{
		"!@#$%^&*()": "!@#$%^&*()",
	}
	_, err = srClient.SetMetadata(context.Background(), md, nil)
	_require.NotNil(err)
}

func (s *azfileLiveTestSuite) TestShareGetStats() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	newQuota := int32(300)

	// In order to test and get LastModified property.
	sResp, err := srClient.SetProperties(context.Background(), &ShareSetPropertiesOptions{Quota: to.Int32Ptr(newQuota)})
	_require.Nil(err)
	_require.Equal(sResp.RawResponse.StatusCode, 200)

	gResp, err := srClient.GetStatistics(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	// _require.NotEqual(*gResp.ETag, "") // TODO: The ETag would be ""
	// _require.Equal(gResp.LastModified.IsZero(), false) // TODO: Even share is once updated, no LastModified would be returned.
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_require.Equal(*gResp.ShareUsageBytes, int64(0))
}

func (s *azfileLiveTestSuite) TestShareGetStatsNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	_, err = srClient.GetStatistics(ctx, nil)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestSetAndGetStatistics() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	cResp, err := srClient.Create(ctx, &ShareCreateOptions{Quota: to.Int32Ptr(1024)})
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer delShare(_require, srClient, nil)

	dirClient, err := srClient.NewDirectoryClient("testdir")
	_require.Nil(err)
	_, err = dirClient.Create(ctx, nil)
	_require.Nil(err)

	fCLient, err := dirClient.NewFileClient("testfile")
	_, err = fCLient.Create(ctx, &CreateFileOptions{FileContentLength: to.Int64Ptr(1024 * 1024 * 1024 * 1024)})
	_require.Nil(err)

	getStats, err := srClient.GetStatistics(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getStats.ShareUsageBytes, int64(1024*1024*1024*1024))
}

//func (s *azfileLiveTestSuite) TestShareCreateSnapshotNonDefault() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_require, shareName, svcClient)
//	defer delShare(_require, srClient, &ShareDeleteOptions{DeleteSnapshots: &deleteSnapshotsInclude})
//
//	ctx := context.Background()
//
//	md := map[string]string{
//		"foo": "FooValuE",
//		"bar": "bArvaLue",
//	}
//
//	cResp, err := srClient.CreateSnapshot(ctx, &ShareCreateSnapshotOptions{Metadata: md})
//	_require.Nil(err)
//	_require.Equal(cResp.RawResponse.StatusCode, chk.Equals, 201)
//	_require.Equal(cResp.Date().IsZero(), chk.Equals, false)
//	_require.NotEqual(*cResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_require.Equal(cResp.LastModified.IsZero(), chk.Equals, false)
//	_assert(cResp.RequestID(), chk.Not(chk.Equals), "")
//	_assert(cResp.Version(), chk.Not(chk.Equals), "")
//	_assert(cResp.Snapshot(), chk.Not(chk.Equals), nil)
//
//	cSnapshot := cResp.Snapshot()
//
//	lResp, err := svcClient.ListSharesSegment(
//		ctx, Marker{},
//		ListSharesOptions{
//			Detail: ListSharesDetail{
//				Metadata:  true,
//				Snapshots: true,
//			},
//			Prefix: shareName,
//		})
//
//	_require.Nil(err)
//	_assert(lResp.Response().StatusCode, chk.Equals, 200)
//	_assert(lResp.ShareItems, chk.HasLen, 2)
//
//	if lResp.ShareItems[0].Snapshot != nil {
//		_assert(*(lResp.ShareItems[0].Snapshot), chk.Equals, cSnapshot)
//		_assert(lResp.ShareItems[0].Metadata, chk.DeepEquals, md)
//		_assert(len(lResp.ShareItems[1].Metadata), chk.Equals, 0)
//	} else {
//		_assert(*(lResp.ShareItems[1].Snapshot), chk.Equals, cSnapshot)
//		_assert(lResp.ShareItems[1].Metadata, chk.DeepEquals, md)
//		_assert(len(lResp.ShareItems[0].Metadata), chk.Equals, 0)
//	}
//
//}

//func (s *azfileLiveTestSuite) TestShareCreateSnapshotDefault() {
//	credential, accountName := getCredential()
//
//	ctx := context.Background()
//
//	u, _ := url.Parse(fmt.Sprintf("https://%s.file.core.windows.net", accountName))
//	serviceURL := NewServiceURL(*u, NewPipeline(credential, PipelineOptions{}))
//
//	shareName := generateShareName(test)
//	shareURL := serviceURL.NewShareClient(shareName)
//
//	_, err := srClient.Create(ctx, map[string]string{}, 0)
//	_require.Nil(err)
//
//	defer srClient.Delete(ctx, DeleteSnapshotsOptionTypeInclude)
//
//	// Let's create a file in the base share.
//	fileURL := srClient.NewRootDirectoryClient().NewFileURL("myfile")
//	_, err = fileURL.Create(ctx, 0, FileHTTPHeaders{}, map[string]string{})
//	_require.Nil(err)
//
//	// Create share snapshot, the snapshot contains the create file.
//	snapshotShare, err := srClient.CreateSnapshot(ctx, map[string]string{})
//	_require.Nil(err)
//
//	// Delete file in base share.
//	_, err = fileURL.Delete(ctx)
//	_require.Nil(err)
//
//	// Restore file from share snapshot.
//	// Create a SAS.
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:   SASProtocolHTTPS,              // Users MUST use HTTPS (not HTTP)
//		ExpiryTime: time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ShareName:  shareName,
//
//		// To produce a share SAS (as opposed to a file SAS), assign to Permissions using
//		// ShareSASPermissions and make sure the DirectoryAndFilePath field is "" (the default).
//		Permissions: ShareSASPermissions{Read: true, Write: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_require.Nil(err)
//
//	// Build a file snapshot URL.
//	fileParts := NewFileURLParts(fileURL.URL())
//	fileParts.ShareSnapshot = snapshotShare.Snapshot()
//	fileParts.SAS = sasQueryParams
//	sourceURL := fileParts.URL()
//
//	// Do restore.
//	_, err = fileURL.StartCopy(ctx, sourceURL, map[string]string{})
//	_require.Nil(err)
//
//	_, err = srClient.WithSnapshot(snapshotShare.Snapshot()).Delete(ctx, DeleteSnapshotsOptionNone)
//	_require.Nil(err)
//}

func (s *azfileLiveTestSuite) TestShareCreateSnapshotNegativeShareNotExist() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_require.Nil(err)

	_, err = srClient.CreateSnapshot(ctx, &ShareCreateSnapshotOptions{Metadata: map[string]string{}})
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareCreateSnapshotNegativeMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	_, err = srClient.CreateSnapshot(ctx, &ShareCreateSnapshotOptions{Metadata: map[string]string{"!@#$%^&*()": "!@#$%^&*()"}})
	_require.NotNil(err)
}

// Note behavior is different from blob's snapshot.
func (s *azfileLiveTestSuite) TestShareCreateSnapshotNegativeSnapshotOfSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, &ShareDeleteOptions{DeleteSnapshots: &deleteSnapshotsInclude})

	snapshotURL := srClient.WithSnapshot(time.Now().UTC().String())
	cResp, err := snapshotURL.CreateSnapshot(ctx, nil)
	_require.Nil(err) //Note: this would not fail, snapshot would be ignored.

	snapshotRecursiveURL := srClient.WithSnapshot(*cResp.Snapshot)
	_, err = snapshotRecursiveURL.CreateSnapshot(ctx, nil)
	_require.Nil(err) //Note: this would not fail, snapshot would be ignored.
}

func validateShareDeleted(_assert *assert.Assertions, srClient ShareClient) {
	_, err := srClient.GetProperties(ctx, nil)
	validateStorageError(_require, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareDeleteSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	resp, err := srClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	snapshotURL := srClient.WithSnapshot(*resp.Snapshot)

	_, err = snapshotURL.Delete(ctx, nil)
	_require.Nil(err)

	validateShareDeleted(_require, snapshotURL)
}

//func (s *azfileLiveTestSuite) TestShareDeleteSnapshotsInclude() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_require, shareName, svcClient)
//
//	_, err = srClient.CreateSnapshot(ctx, nil)
//	_require.Nil(err)
//	_, err = srClient.Delete(ctx, &ShareDeleteOptions{DeleteSnapshots: &deleteSnapshotsInclude})
//	_require.Nil(err)
//
//	lResp, _ := svcClient.ListSharesSegment(ctx, Marker{}, ListSharesOptions{Detail: ListSharesDetail{Snapshots: true}, Prefix: shareName})
//	_assert(lResp.ShareItems, chk.HasLen, 0)
//}

func (s *azfileLiveTestSuite) TestShareDeleteSnapshotsNoneWithSnapshots() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, &ShareDeleteOptions{DeleteSnapshots: &deleteSnapshotsInclude})

	_, err = srClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	_, err = srClient.Delete(ctx, nil)
	validateStorageError(_require, err, StorageErrorCodeShareHasSnapshots)
}
