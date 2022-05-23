package azfile

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"strconv"
	"time"

	//"context"
	//"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	//"github.com/stretchr/testify/require"
	//chk "gopkg.in/check.v1"
	"github.com/stretchr/testify/require"
	"os"
)

func (s *azfileLiveTestSuite) TestAccountNewServiceURLValidName() {
	_require := require.New(s.T())
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/"
	_require.Equal(svcClient.URL(), correctURL)
}

func (s *azfileLiveTestSuite) TestAccountNewShareURLValidName() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)

	shareName := generateShareName(testName)
	srClient, err := svcClient.NewShareClient(shareName)
	_require.Nil(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/" + shareName
	_require.Equal(srClient.URL(), correctURL)
}

////type testPipeline struct{}
////
////const testPipelineMessage string = "Test factory invoked"
////
////func (tm testPipeline) Do(ctx context.Context, methodFactory pipeline.Factory, request pipeline.Request) (pipeline.Response, error) {
////	return nil, errors.New(testPipelineMessage)
////}
////
////func (s *azfileLiveTestSuite) TestAccountWithPipeline() {
////	_require := require.New(s.T())
////	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
////	if err != nil {
////		s.Fail("Unable to fetch service client because " + err.Error())
////	}
////	fsu = svcClient.WithPipeline(testPipeline{}) // testPipeline returns an identifying message as an error
////	shareURL := fsu.NewShareURL("name")
////
////	_, err := shareURL.Create(ctx, azfile.Metadata{}, 0)
////
////	_require(err.Error(), chk.Equals, testPipelineMessage)
////}
//
////// This case is not stable, as service side returns 202, if it previously has value,
////// it need unpredictable time to make updates take effect.
////func (s *azfileLiveTestSuite) TestAccountGetSetPropertiesDefault() {
////	_require := require.New(s.T())
////	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
////	if err != nil {
////		s.Fail("Unable to fetch service client because " + err.Error())
////	}
////
////	resp, err := svcClient.SetProperties(context.Background(), &ServiceSetPropertiesOptions{})
////	_require.Nil(err)
////	_require(resp.Response().StatusCode, chk.Equals, 202)
////	_require(resp.RequestID(), chk.Not(chk.Equals), "")
////	_require(resp.Version(), chk.Not(chk.Equals), "")
////
////	time.Sleep(time.Second * 15)
////
////	// Note: service side is 202, might depend on timing
////	props, err := sa.GetProperties(context.Background())
////	_require.Nil(err)
////	_require(props.Response().StatusCode, chk.Equals, 200)
////	_require(props.RequestID(), chk.Not(chk.Equals), "")
////	_require(props.Version(), chk.Not(chk.Equals), "")
////	_require(props.HourMetrics, chk.NotNil)
////	_require(props.MinuteMetrics, chk.NotNil)
////	//_require(props.Cors, chk.HasLen, 0) //Unstable evaluation
////}

func (s *azfileLiveTestSuite) TestAccountProperties() {
	_require := require.New(s.T())
	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)

	setPropertiesOptions := &ServiceSetPropertiesOptions{
		HourMetrics: &ShareMetricProperties{
			Enabled:                to.Ptr(true),
			IncludeAPIs:            to.Ptr(true),
			RetentionPolicyEnabled: to.Ptr(true),
			RetentionDays:          to.Ptr(int32(2)),
		},
		MinuteMetrics: &ShareMetricProperties{
			Enabled:                to.Ptr(true),
			IncludeAPIs:            to.Ptr(false),
			RetentionPolicyEnabled: to.Ptr(true),
			RetentionDays:          to.Ptr(int32(2)),
		},
		Cors: []*ShareCorsRule{
			{
				AllowedOrigins:  to.Ptr("*"),
				AllowedMethods:  to.Ptr("PUT"),
				AllowedHeaders:  to.Ptr("x-ms-client-request-id"),
				ExposedHeaders:  to.Ptr("x-ms-*"),
				MaxAgeInSeconds: to.Ptr(int32(2)),
			},
		},
	}
	setPropsResponse, err := svcClient.SetProperties(ctx, setPropertiesOptions)
	_require.Nil(err)
	// _require.Equal(setPropsResponse.RawResponse.StatusCode, 202)
	_require.NotEqual(setPropsResponse.RequestID, "")
	_require.NotEqual(setPropsResponse.Version, "")

	time.Sleep(time.Second * 30)

	props, err := svcClient.GetProperties(ctx, nil)
	_require.Nil(err)
	//_require.Equal(props.RawResponse.StatusCode, 200)
	_require.NotEqual(props.RequestID, "")
	_require.NotEqual(props.Version, "")
	_require.EqualValues(props.HourMetrics.RetentionPolicy.Enabled, setPropertiesOptions.HourMetrics.RetentionPolicyEnabled)
	_require.EqualValues(props.HourMetrics.RetentionPolicy.Days, setPropertiesOptions.HourMetrics.RetentionDays)
	_require.EqualValues(props.MinuteMetrics.RetentionPolicy.Enabled, setPropertiesOptions.MinuteMetrics.RetentionPolicyEnabled)
	_require.EqualValues(props.MinuteMetrics.RetentionPolicy.Days, setPropertiesOptions.MinuteMetrics.RetentionDays)
	_require.Len(props.Cors, len(setPropertiesOptions.Cors))
}

func (s *azfileLiveTestSuite) TestAccountHourMetrics() {
	_require := require.New(s.T())
	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)

	setPropertiesOptions := &ServiceSetPropertiesOptions{
		HourMetrics: &ShareMetricProperties{
			Enabled:                to.Ptr(true),
			IncludeAPIs:            to.Ptr(true),
			RetentionPolicyEnabled: to.Ptr(true),
			RetentionDays:          to.Ptr(int32(5)),
		},
	}
	setPropertiesResponse, err := svcClient.SetProperties(ctx, setPropertiesOptions)
	_require.Nil(err)
	_ = setPropertiesResponse
}

//// TODO: This case is not stable... As SetProperties returns 202 Accepted, it depends on server side how fast properties would be set.
//// func (s *azfileLiveTestSuite) TestAccountGetSetPropertiesNonDefaultWithDisable() {
//// 	sa := getFSU()
//
//// 	setProps := azfile.FileServiceProperties{
//// 		HourMetrics: azfile.ShareMetricProperties{
//// 			Enabled: false,
//// 		},
//// 		MinuteMetrics: azfile.ShareMetricProperties{
//// 			Enabled: false,
//// 		},
//// 	}
//// 	resp, err := sa.SetProperties(context.Background(), setProps)
//// 	_require.Nil(err)
//// 	_require(resp.Response().StatusCode, chk.Equals, 202)
//// 	_require(resp.RequestID(), chk.Not(chk.Equals), "")
//// 	_require(resp.Version(), chk.Not(chk.Equals), "")
//
//// 	time.Sleep(time.Second * 5)
//
//// 	props, err := sa.GetProperties(context.Background())
//// 	_require.Nil(err)
//// 	_require(props.Response().StatusCode, chk.Equals, 200)
//// 	_require(props.RequestID(), chk.Not(chk.Equals), "")
//// 	_require(props.Version(), chk.Not(chk.Equals), "")
//// 	_require(props.HourMetrics, chk.DeepEquals, azfile.ShareMetricProperties{Enabled: false})
//// 	_require(props.MinuteMetrics, chk.DeepEquals, azfile.ShareMetricProperties{Enabled: false})
//// 	_require(props.Cors, chk.IsNil)
//// }

func (s *azfileLiveTestSuite) TestAccountListSharesNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	mySharePrefix := generateEntityName(testName)
	pager := svcClient.ListShares(&ServiceListSharesOptions{Prefix: to.Ptr(mySharePrefix)})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.NoError(err)
		_require.NotNil(resp.Prefix)
		_require.Equal(*resp.Prefix, mySharePrefix)
		_require.NotNil(resp.ServiceEndpoint)
		_require.NotNil(resp.RequestID)
		_require.NotNil(resp.Version)
		//_require.Equal(resp.RawResponse.StatusCode, 200)
		_require.Len(resp.ShareItems, 0)
	}

	shareClients := map[string]*ShareClient{}
	for i := 0; i < 4; i++ {
		shareName := mySharePrefix + "share" + strconv.Itoa(i)
		shareClients[shareName] = createNewShare(_require, shareName, svcClient)

		_, err := shareClients[shareName].SetMetadata(ctx, basicMetadata, nil)
		_require.Nil(err)

		_, err = shareClients[shareName].CreateSnapshot(ctx, nil)
		_require.Nil(err)

		defer delShare(_require, shareClients[shareName], &ShareDeleteOptions{
			DeleteSnapshots: to.Ptr(DeleteSnapshotsOptionTypeInclude),
		})
	}

	pager = svcClient.ListShares(&ServiceListSharesOptions{
		Include:    []ListSharesIncludeType{ListSharesIncludeTypeMetadata, ListSharesIncludeTypeSnapshots},
		Prefix:     to.Ptr(mySharePrefix),
		MaxResults: to.Ptr(int32(2)),
	})

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		if len(resp.ShareItems) > 0 {
			_require.Len(resp.ShareItems, 2)
		}
		for _, shareItem := range resp.ShareItems {
			_require.NotNil(shareItem.Properties)
			_require.NotNil(shareItem.Properties.LastModified)
			_require.NotNil(shareItem.Properties.Etag)
			_require.Len(shareItem.Metadata, len(basicMetadata))
			for key, val1 := range basicMetadata {
				if val2, ok := shareItem.Metadata[key]; !(ok && val1 == *val2) {
					_require.Fail("metadata mismatch")
				}
			}
			_require.NotNil(resp.ShareItems[0].Snapshot)
			_require.Nil(resp.ShareItems[1].Snapshot)
		}
	}
}

func (s *azfileLiveTestSuite) TestAccountListSharesInvalidMaxResults() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	pager := svcClient.ListShares(&ServiceListSharesOptions{MaxResults: to.Ptr(int32(-2))})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		_require.Contains(err.Error(), "OutOfRangeQueryParameterValue")
		break
	}

	pager2 := svcClient.ListShares(&ServiceListSharesOptions{MaxResults: to.Ptr(int32(0))})
	for pager2.More() {
		_, err := pager2.NextPage(ctx)
		_require.NotNil(err)
		_require.Contains(err.Error(), "OutOfRangeQueryParameterValue")
		break
	}
}

//
//func (s *azfileLiveTestSuite) TestAccountSAS() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)
//
//	shareName := generateShareName(sharePrefix, testName)
//	srClient := createNewShare(_require, shareName, svcClient)
//
//	dirClient := getDirectoryClientFromShare(_require, "dir1", srClient)
//	fClient := getFileClientFromDirectory(_require, "file1", dirClient)
//
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	_require.Nil(err)
//
//	sasQueryParams, err := AccountSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().Add(1 * time.Hour),
//		Permissions:   AccountSASPermissions{Read: true, List: true, Write: true, Delete: true, Add: true, Create: true, Update: true, Process: true}.String(),
//		Services:      AccountSASServices{File: true, Blob: true, Queue: true}.String(),
//		ResourceTypes: AccountSASResourceTypes{Service: true, Container: true, Object: true}.String(),
//	}.Sign(credential)
//	_require.Nil(err)
//
//	// Reverse validation all parse logics work as expect.
//	ap := &AccountSASPermissions{}
//	err = ap.Parse(sasQueryParams.Permissions())
//	_require.Nil(err)
//	_require.EqualValues(*ap, AccountSASPermissions{Read: true, List: true, Write: true, Delete: true, Add: true, Create: true, Update: true, Process: true})
//
//	as := &AccountSASServices{}
//	err = as.Parse(sasQueryParams.Services())
//	_require.Nil(err)
//	_require.EqualValues(*as, AccountSASServices{File: true, Blob: true, Queue: true})
//
//	ar := &AccountSASResourceTypes{}
//	err = ar.Parse(sasQueryParams.ResourceTypes())
//	_require.Nil(err)
//	_require.EqualValues(*ar, AccountSASResourceTypes{Service: true, Container: true, Object: true})
//
//	// Test service URL
//	svcParts, err := NewFileURLParts(svcClient.URL())
//	svcParts.SAS = sasQueryParams
//	testSvcURL := svcParts.URL()
//	svcURLWithSAS, err := NewServiceClient(testSvcURL, azcore.TokenCredential(nil), nil)
//
//	// List
//	pager := svcURLWithSAS.ListShares(nil)
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//		_ = resp
//	}
//	err = pager.Err()
//	_require.Nil(err)
//
//	// Write
//	_, err = svcURLWithSAS.SetProperties(ctx, &ServiceSetPropertiesOptions{})
//	_require.Nil(err)
//
//	// Read
//	_, err = svcURLWithSAS.GetProperties(ctx, &ServiceGetPropertiesOptions{})
//	_require.Nil(err)
//
//	// Test share URL
//	sParts, err := NewFileURLParts(srClient.URL())
//	_require.Nil(err)
//	_require.Equal(sParts.ShareName, shareName)
//	sParts.SAS = sasQueryParams
//	testShareURL := sParts.URL()
//
//	shareURLWithSAS, err := NewShareClient(testShareURL, azcore.TokenCredential(nil), nil)
//	_require.Nil(err)
//
//	// Create
//	_, err = shareURLWithSAS.Create(ctx, &ShareCreateOptions{})
//	_require.Nil(err)
//
//	// Write
//	_, err = shareURLWithSAS.SetMetadata(ctx, basicMetadata, nil)
//	_require.Nil(err)
//
//	// Read
//	gResp, err := shareURLWithSAS.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.EqualValues(gResp.Metadata, basicMetadata)
//
//	// Delete
//	defer shareURLWithSAS.Delete(ctx, nil)
//
//	// Test dir URL
//	dParts, err := NewFileURLParts(dirClient.URL())
//	_require.Nil(err)
//
//	dParts.SAS = sasQueryParams
//	testDirURL := dParts.URL()
//
//	dirURLWithSAS, err := NewDirectoryClient(testDirURL, azcore.TokenCredential(nil), nil)
//	_require.Nil(err)
//
//	// Create
//	_, err = dirURLWithSAS.Create(ctx, nil)
//	_require.Nil(err)
//
//	// Write
//	_, err = dirURLWithSAS.SetMetadata(ctx, basicMetadata, nil)
//	_require.Nil(err)
//
//	// Read
//	gdResp, err := dirURLWithSAS.GetProperties(ctx, nil)
//	_require.Nil(err)
//
//	_require.EqualValues(gdResp.Metadata, basicMetadata)
//
//	// List
//	pager2 := dirURLWithSAS.ListFilesAndDirectories(nil)
//	for pager2.NextPage(ctx) {
//		resp := pager2.PageResponse()
//		_ = resp
//	}
//	err = pager2.Err()
//	_require.Nil(err)
//
//	// Delete
//	defer dirURLWithSAS.Delete(ctx, nil)
//
//	// Test file URL
//	fParts, err := NewFileURLParts(fClient.URL())
//	_require.Nil(err)
//
//	fParts.SAS = sasQueryParams
//	testFileURL := fParts.URL()
//	fileURLWithSAS, err := NewFileClient(testFileURL, azcore.TokenCredential(nil), nil)
//	_require.Nil(err)
//
//	// Create
//	_, err = fileURLWithSAS.Create(ctx, nil)
//	_require.Nil(err)
//
//	// Write
//	_, err = fileURLWithSAS.SetMetadata(ctx, basicMetadata, nil)
//	// Read
//	gfResp, err := fileURLWithSAS.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.EqualValues(gfResp.Metadata, basicMetadata)
//	// Delete
//	defer fileURLWithSAS.Delete(ctx, nil)
//}
