package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"time"
)

func (s *azfileLiveTestSuite) TestShareCreateRootDirectoryURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, shareClient, nil)
}

func (s *azfileLiveTestSuite) TestPutAndGetPermission() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)
	_assert.Nil(err)

	// Create a permission and check that it's not empty.
	createResp, err := srClient.CreatePermission(ctx, sampleSDDL, nil)
	_assert.Nil(err)
	_assert.NotEqual(*createResp.FilePermissionKey, "")

	getResp, err := srClient.GetPermission(ctx, *createResp.FilePermissionKey, nil)
	_assert.Nil(err)
	// Rather than checking against the original, we check for emptiness, as Azure Files has set a nil-ness flag on SACLs
	//        and converted our well-known SID.
	/*
		Expected :string = "O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)"
		Actual   :string = "O:AOG:S-1-5-21-397955417-626881126-188441444-512D:(A;;CCDCLCSWRPWPRCWDWOGA;;;S-1-0-0)S:NO_ACCESS_CONTROL"
	*/
	_assert.NotEqual(*getResp.Permission, "")
}

func (s *azfileLiveTestSuite) TestShareCreateDirectoryURL() {
	_assert := assert.New(s.T())
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient, err := svcClient.NewShareClient(sharePrefix)
	_assert.Nil(err)
	dirClient := srClient.NewDirectoryClient(directoryPrefix)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/" + sharePrefix + "/" + directoryPrefix
	_assert.Equal(dirClient.URL(), correctURL)
}

// Note: test share create with default parameter is covered with preparing phase for FileURL and etc.
//func (s *azfileLiveTestSuite) TestShareCreateDeleteNonDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient, err := svcClient.NewShareClient(shareName)
//	_assert.Nil(err)
//
//	md := map[string]string{
//		"foo": "FooValuE",
//		"bar": "bArvaLue",
//	}
//
//	quota := int32(1000)
//
//	cResp, err := srClient.Create(context.Background(), &CreateShareOptions{Quota: to.Int32Ptr(quota), Metadata: md})
//	_assert.Nil(err)
//	_assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
//	_assert(cResp.Date().IsZero(), chk.Equals, false)
//	_assert(cResp.ETag(), chk.Not(chk.Equals), ETagNone)
//	_assert(cResp.LastModified.IsZero(), chk.Equals, false)
//	_assert(cResp.RequestID(), chk.Not(chk.Equals), "")
//	_assert(cResp.Version(), chk.Not(chk.Equals), "")
//
//	shares, err := srClient.ListSharesSegment(context.Background(), Marker{}, ListSharesOptions{Prefix: shareName, Detail: ListSharesDetail{Metadata: true}})
//	_assert.Nil(err)
//	_assert(shares.ShareItems, chk.HasLen, 1)
//	_assert(shares.ShareItems[0].Name, chk.Equals, shareName)
//	_assert(shares.ShareItems[0].Metadata, chk.DeepEquals, md)
//	_assert(shares.ShareItems[0].Properties.Quota, chk.Equals, quota)
//
//	dResp, err := srClient.Delete(context.Background(), DeleteSnapshotsOptionNone)
//	_assert.Nil(err)
//	_assert(dResp.RawResponse.StatusCode, chk.Equals, 202)
//	_assert(dResp.Date().IsZero(), chk.Equals, false)
//	_assert(dResp.RequestID(), chk.Not(chk.Equals), "")
//	_assert(dResp.Version(), chk.Not(chk.Equals), "")
//
//	shares, err = srClient.ListSharesSegment(context.Background(), Marker{}, ListSharesOptions{Prefix: shareName})
//	_assert.Nil(err)
//	_assert(shares.ShareItems, chk.HasLen, 0)
//}

func (s *azfileLiveTestSuite) TestShareCreateNilMetadata() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	_, err = srClient.Create(ctx, nil)
	defer delShare(_assert, srClient, nil)
	_assert.Nil(err)

	response, err := srClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(response.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestShareCreateNegativeInvalidName() {
	_assert := assert.New(s.T())
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient, err := svcClient.NewShareClient("foo bar")
	_assert.Nil(err)

	_, err = srClient.Create(ctx, nil)

	validateStorageError(_assert, err, StorageErrorCodeInvalidResourceName)
}

func (s *azfileLiveTestSuite) TestShareCreateNegativeInvalidMetadata() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	_, err = srClient.Create(ctx, &CreateShareOptions{Metadata: map[string]string{"!@#$%^&*()": "!@#$%^&*()"}, Quota: to.Int32Ptr(0)})
	_assert.NotNil(err)
}

func (s *azfileLiveTestSuite) TestShareDeleteNegativeNonExistent() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	_, err = srClient.Delete(ctx, nil)
	validateStorageError(_assert, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareGetSetPropertiesNonDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	newQuota := int32(1234)

	sResp, err := srClient.SetProperties(ctx, &SetSharePropertiesOptions{Quota: to.Int32Ptr(newQuota)})
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.Equal(sResp.LastModified.IsZero(), false)
	_assert.NotEqual(sResp.RequestID, "")
	_assert.NotEqual(sResp.Version, "")
	_assert.Equal(sResp.Date.IsZero(), false)

	props, err := srClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(props.RawResponse.StatusCode, 200)
	_assert.NotEqual(*props.ETag, "")
	_assert.Equal(props.LastModified.IsZero(), false)
	_assert.NotEqual(*props.RequestID, "")
	_assert.NotEqual(*props.Version, "")
	_assert.Equal(props.Date.IsZero(), false)
	_assert.Equal(*props.Quota, newQuota)
}

func (s *azfileLiveTestSuite) TestShareGetSetPropertiesDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	sResp, err := srClient.SetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.Equal(sResp.LastModified.IsZero(), false)
	_assert.NotEqual(sResp.RequestID, "")
	_assert.NotEqual(sResp.Version, "")
	_assert.Equal(sResp.Date.IsZero(), false)

	props, err := srClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(props.RawResponse.StatusCode, 200)
	_assert.NotEqual(*props.ETag, "")
	_assert.Equal(props.LastModified.IsZero(), false)
	_assert.NotEqual(*props.RequestID, "")
	_assert.NotEqual(*props.Version, "")
	_assert.Equal(props.Date.IsZero(), false)
	_assert.True(*props.Quota >= 0) // When using service default quota, it could be any value
}

func (s *azfileLiveTestSuite) TestShareSetQuotaNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	_, err = srClient.SetProperties(ctx, &SetSharePropertiesOptions{Quota: to.Int32Ptr(-1)})
	_assert.NotNil(err)
	_assert.Contains(err.Error(), "validation failed: share quote cannot be negative")
}

func (s *azfileLiveTestSuite) TestShareGetPropertiesNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	_, err = srClient.GetProperties(ctx, nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareGetSetPermissionsNonDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

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
	_assert.EqualValues(*pS2, pS)

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
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.Equal(sResp.LastModified.IsZero(), false)
	_assert.NotEqual(sResp.RequestID, "")
	_assert.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetPermissions(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.Len(gResp.SignedIdentifiers, 1)
	_assert.EqualValues(*(gResp.SignedIdentifiers[0]), *permissions[0])
}

func (s *azfileLiveTestSuite) TestShareGetSetPermissionsNonDefaultMultiple() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

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
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.Equal(sResp.LastModified.IsZero(), false)
	_assert.NotEqual(sResp.RequestID, "")
	_assert.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetPermissions(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.Len(gResp.SignedIdentifiers, 2)
	_assert.EqualValues(gResp.SignedIdentifiers[0], permissions[0])
}

func (s *azfileLiveTestSuite) TestShareGetSetPermissionsDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	sResp, err := srClient.SetPermissions(context.Background(), []*SignedIdentifier{}, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.Equal(sResp.LastModified.IsZero(), false)
	_assert.NotEqual(sResp.RequestID, "")
	_assert.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetPermissions(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.Len(gResp.SignedIdentifiers, 0)
}

func (s *azfileLiveTestSuite) TestShareGetPermissionNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	_, err = srClient.GetPermissions(ctx, nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareSetPermissionsNonDefaultDeleteAndModifyACL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)

	defer delShare(_assert, srClient, nil)

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
	_assert.Nil(err)

	resp, err := srClient.GetPermissions(ctx, nil)
	_assert.Nil(err)

	_assert.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	permissions[0].ID = to.StringPtr("0004") // Modify the remaining policy which is at index 0 in the new slice
	_, err = srClient.SetPermissions(ctx, permissions, nil)

	resp, err = srClient.GetPermissions(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.SignedIdentifiers, 1)

	_assert.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *azfileLiveTestSuite) TestShareSetPermissionsDeleteAllPolicies() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)

	defer delShare(_assert, srClient, nil)

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
	_assert.Nil(err)

	_, err = srClient.SetPermissions(ctx, []*SignedIdentifier{}, nil)
	_assert.Nil(err)

	resp, err := srClient.GetPermissions(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.SignedIdentifiers, 0)
}

// Note: No error happened
func (s *azfileLiveTestSuite) TestShareSetPermissionsNegativeInvalidPolicyTimes() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

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
	_assert.Nil(err)
}

func (s *azfileLiveTestSuite) TestShareSetPermissionsNilPolicySlice() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)

	defer delShare(_assert, srClient, nil)

	_, err = srClient.SetPermissions(ctx, nil, nil)
	_assert.Nil(err)
}

// SignedIdentifier ID too long
func (s *azfileLiveTestSuite) TestShareSetPermissionsNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)

	defer delShare(_assert, srClient, nil)

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
	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}

func (s *azfileLiveTestSuite) TestShareGetSetMetadataDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	sResp, err := srClient.SetMetadata(context.Background(), map[string]string{}, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.Equal(sResp.LastModified.IsZero(), false)
	_assert.NotEqual(sResp.RequestID, "")
	_assert.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.Len(gResp.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestShareGetSetMetadataNonDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	md := map[string]string{
		"Foo": "FooValuE",
		"Bar": "bArvaLue", // Note: As testing result, currently only support case-insensitive keys(key will be saved in lower-case).
	}
	sResp, err := srClient.SetMetadata(context.Background(), md, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.Equal(sResp.LastModified.IsZero(), false)
	_assert.NotEqual(sResp.RequestID, "")
	_assert.NotEqual(sResp.Version, "")

	gResp, err := srClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.EqualValues(gResp.Metadata, md)
}

func (s *azfileLiveTestSuite) TestShareSetMetadataNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	md := map[string]string{
		"!@#$%^&*()": "!@#$%^&*()",
	}
	_, err = srClient.SetMetadata(context.Background(), md, nil)
	_assert.NotNil(err)
}

func (s *azfileLiveTestSuite) TestShareGetStats() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	newQuota := int32(300)

	// In order to test and get LastModified property.
	sResp, err := srClient.SetProperties(context.Background(), &SetSharePropertiesOptions{Quota: to.Int32Ptr(newQuota)})
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)

	gResp, err := srClient.GetStatistics(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	// _assert.NotEqual(*gResp.ETag, "") // TODO: The ETag would be ""
	// _assert.Equal(gResp.LastModified.IsZero(), false) // TODO: Even share is once updated, no LastModified would be returned.
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.Equal(*gResp.ShareUsageBytes, int64(0))
}

func (s *azfileLiveTestSuite) TestShareGetStatsNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	_, err = srClient.GetStatistics(ctx, nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestSetAndGetStatistics() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	cResp, err := srClient.Create(ctx, &CreateShareOptions{Quota: to.Int32Ptr(1024)})
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	defer delShare(_assert, srClient, nil)

	dirClient := srClient.NewDirectoryClient("testdir")
	_, err = dirClient.Create(ctx, nil)
	_assert.Nil(err)

	fCLient, err := dirClient.NewFileClient("testfile")
	_, err = fCLient.Create(ctx, &CreateFileOptions{FileContentLength: to.Int64Ptr(1024 * 1024 * 1024 * 1024)})
	_assert.Nil(err)

	getStats, err := srClient.GetStatistics(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getStats.ShareUsageBytes, int64(1024*1024*1024*1024))
}

//func (s *azfileLiveTestSuite) TestShareCreateSnapshotNonDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_assert, shareName, svcClient)
//	defer delShare(_assert, srClient, &DeleteShareOptions{DeleteSnapshots: &deleteSnapshotsInclude})
//
//	ctx := context.Background()
//
//	md := map[string]string{
//		"foo": "FooValuE",
//		"bar": "bArvaLue",
//	}
//
//	cResp, err := srClient.CreateSnapshot(ctx, &CreateShareSnapshotOptions{Metadata: md})
//	_assert.Nil(err)
//	_assert.Equal(cResp.RawResponse.StatusCode, chk.Equals, 201)
//	_assert.Equal(cResp.Date().IsZero(), chk.Equals, false)
//	_assert.NotEqual(*cResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.Equal(cResp.LastModified.IsZero(), chk.Equals, false)
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
//	_assert.Nil(err)
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
//	_assert.Nil(err)
//
//	defer srClient.Delete(ctx, DeleteSnapshotsOptionTypeInclude)
//
//	// Let's create a file in the base share.
//	fileURL := srClient.NewRootDirectoryClient().NewFileURL("myfile")
//	_, err = fileURL.Create(ctx, 0, FileHTTPHeaders{}, map[string]string{})
//	_assert.Nil(err)
//
//	// Create share snapshot, the snapshot contains the create file.
//	snapshotShare, err := srClient.CreateSnapshot(ctx, map[string]string{})
//	_assert.Nil(err)
//
//	// Delete file in base share.
//	_, err = fileURL.Delete(ctx)
//	_assert.Nil(err)
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
//	_assert.Nil(err)
//
//	// Build a file snapshot URL.
//	fileParts := NewFileURLParts(fileURL.URL())
//	fileParts.ShareSnapshot = snapshotShare.Snapshot()
//	fileParts.SAS = sasQueryParams
//	sourceURL := fileParts.URL()
//
//	// Do restore.
//	_, err = fileURL.StartCopy(ctx, sourceURL, map[string]string{})
//	_assert.Nil(err)
//
//	_, err = srClient.WithSnapshot(snapshotShare.Snapshot()).Delete(ctx, DeleteSnapshotsOptionNone)
//	_assert.Nil(err)
//}

func (s *azfileLiveTestSuite) TestShareCreateSnapshotNegativeShareNotExist() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient, err := getShareClient(shareName, svcClient)
	_assert.Nil(err)

	_, err = srClient.CreateSnapshot(ctx, &CreateShareSnapshotOptions{Metadata: map[string]string{}})
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareCreateSnapshotNegativeMetadataInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	_, err = srClient.CreateSnapshot(ctx, &CreateShareSnapshotOptions{Metadata: map[string]string{"!@#$%^&*()": "!@#$%^&*()"}})
	_assert.NotNil(err)
}

// Note behavior is different from blob's snapshot.
func (s *azfileLiveTestSuite) TestShareCreateSnapshotNegativeSnapshotOfSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, &DeleteShareOptions{DeleteSnapshots: &deleteSnapshotsInclude})

	snapshotURL := srClient.WithSnapshot(time.Now().UTC().String())
	cResp, err := snapshotURL.CreateSnapshot(ctx, nil)
	_assert.Nil(err) //Note: this would not fail, snapshot would be ignored.

	snapshotRecursiveURL := srClient.WithSnapshot(*cResp.Snapshot)
	_, err = snapshotRecursiveURL.CreateSnapshot(ctx, nil)
	_assert.Nil(err) //Note: this would not fail, snapshot would be ignored.
}

func validateShareDeleted(_assert *assert.Assertions, srClient ShareClient) {
	_, err := srClient.GetProperties(ctx, nil)
	validateStorageError(_assert, err, StorageErrorCodeShareNotFound)
}

func (s *azfileLiveTestSuite) TestShareDeleteSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, nil)

	resp, err := srClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	snapshotURL := srClient.WithSnapshot(*resp.Snapshot)

	_, err = snapshotURL.Delete(ctx, nil)
	_assert.Nil(err)

	validateShareDeleted(_assert, snapshotURL)
}

//func (s *azfileLiveTestSuite) TestShareDeleteSnapshotsInclude() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_assert, shareName, svcClient)
//
//	_, err = srClient.CreateSnapshot(ctx, nil)
//	_assert.Nil(err)
//	_, err = srClient.Delete(ctx, &DeleteShareOptions{DeleteSnapshots: &deleteSnapshotsInclude})
//	_assert.Nil(err)
//
//	lResp, _ := svcClient.ListSharesSegment(ctx, Marker{}, ListSharesOptions{Detail: ListSharesDetail{Snapshots: true}, Prefix: shareName})
//	_assert(lResp.ShareItems, chk.HasLen, 0)
//}

func (s *azfileLiveTestSuite) TestShareDeleteSnapshotsNoneWithSnapshots() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient, &DeleteShareOptions{DeleteSnapshots: &deleteSnapshotsInclude})

	_, err = srClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	_, err = srClient.Delete(ctx, nil)
	validateStorageError(_assert, err, StorageErrorCodeShareHasSnapshots)
}
