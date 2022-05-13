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
	//"time"
	//"context"
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

	shareName := generateShareName(sharePrefix, testName)
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
		HourMetrics: &Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(true),
			Version:     to.Ptr("2020-02-10"),
			RetentionPolicy: &RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(1)),
			},
		},
		MinuteMetrics: &Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(false),
			Version:     to.Ptr("2020-02-10"),
			RetentionPolicy: &RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(2)),
			},
		},
		Cors: []*CorsRule{
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
	_require.Equal(setPropsResponse.RawResponse.StatusCode, 202)
	_require.NotEqual(setPropsResponse.RequestID, "")
	_require.NotEqual(setPropsResponse.Version, "")

	time.Sleep(time.Second * 30)

	props, err := svcClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(props.RawResponse.StatusCode, 200)
	_require.NotEqual(props.RequestID, "")
	_require.NotEqual(props.Version, "")
	_require.EqualValues(props.HourMetrics, setPropertiesOptions.HourMetrics)
	_require.EqualValues(props.MinuteMetrics, setPropertiesOptions.MinuteMetrics)
	_require.EqualValues(props.Cors, setPropertiesOptions.Cors)
}

func (s *azfileLiveTestSuite) TestAccountHourMetrics() {
	_require := require.New(s.T())
	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)

	setPropertiesOptions := &ServiceSetPropertiesOptions{
		HourMetrics: &Metrics{
			Enabled:     to.Ptr(true),
			IncludeAPIs: to.Ptr(true),
			RetentionPolicy: &RetentionPolicy{
				Enabled: to.Ptr(true),
				Days:    to.Ptr(int32(5)),
			},
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
//// 		HourMetrics: azfile.MetricProperties{
//// 			MetricEnabled: false,
//// 		},
//// 		MinuteMetrics: azfile.MetricProperties{
//// 			MetricEnabled: false,
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
//// 	_require(props.HourMetrics, chk.DeepEquals, azfile.MetricProperties{MetricEnabled: false})
//// 	_require(props.MinuteMetrics, chk.DeepEquals, azfile.MetricProperties{MetricEnabled: false})
//// 	_require(props.Cors, chk.IsNil)
//// }

func (s *azfileLiveTestSuite) TestAccountListSharesNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	mySharePrefix := generateEntityName(testName)
	pager := svcClient.ListShares(&ServiceListSharesOptions{Prefix: to.Ptr(mySharePrefix)})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)

		resp := pager.PageResponse()
		_require.NotNil(resp.Prefix)
		_require.Equal(*resp.Prefix, mySharePrefix)
		_require.NotNil(resp.ServiceEndpoint)
		_require.NotNil(resp.RequestID)
		_require.NotNil(resp.Version)
		_require.Equal(resp.RawResponse.StatusCode, 200)
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

	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)

		resp := pager.PageResponse()
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

	shareName := generateShareName(sharePrefix, testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	pager := svcClient.ListShares(&ServiceListSharesOptions{MaxResults: to.Ptr(int32(-2))})
	for pager.NextPage(ctx) {
		_require.Fail("Shouldn't have reached here")
	}
	err := pager.Err()
	_require.NotNil(err)
	_require.Contains(err.Error(), "OutOfRangeQueryParameterValue")

	pager = svcClient.ListShares(&ServiceListSharesOptions{MaxResults: to.Ptr(int32(0))})
	for pager.NextPage(ctx) {
		_require.Fail("Shouldn't have reached here")
	}
	err = pager.Err()
	_require.NotNil(err)
	_require.Contains(err.Error(), "OutOfRangeQueryParameterValue")

}

//func (s *azfileLiveTestSuite) TestAccountSAS() {
//	fsu := getFSU()
//	shareURL, shareName := getShareURL(c, fsu)
//	dirURL, _ := getDirectoryURLFromShare(c, shareURL)
//	fileURL, _ := getFileURLFromDirectory(c, dirURL)
//
//	credential, _ := getCredential()
//	sasQueryParams, err := azfile.AccountSASSignatureValues{
//		Protocol:      azfile.SASProtocolHTTPS,
//		ExpiryTime:    time.Now().Add(48 * time.Hour),
//		Permissions:   azfile.AccountSASPermissions{Read: true, List: true, Write: true, Delete: true, Add: true, Create: true, Update: true, Process: true}.String(),
//		Services:      azfile.AccountSASServices{File: true, Blob: true, Queue: true}.String(),
//		ResourceTypes: azfile.AccountSASResourceTypes{Service: true, Container: true, Object: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_require.Nil(err)
//
//	// Reverse valiadation all parse logics work as expect.
//	ap := &azfile.AccountSASPermissions{}
//	err = ap.Parse(sasQueryParams.Permissions())
//	_require.Nil(err)
//	_require(*ap, chk.DeepEquals, azfile.AccountSASPermissions{Read: true, List: true, Write: true, Delete: true, Add: true, Create: true, Update: true, Process: true})
//
//	as := &azfile.AccountSASServices{}
//	err = as.Parse(sasQueryParams.Services())
//	_require.Nil(err)
//	_require(*as, chk.DeepEquals, azfile.AccountSASServices{File: true, Blob: true, Queue: true})
//
//	ar := &azfile.AccountSASResourceTypes{}
//	err = ar.Parse(sasQueryParams.ResourceTypes())
//	_require.Nil(err)
//	_require(*ar, chk.DeepEquals, azfile.AccountSASResourceTypes{Service: true, Container: true, Object: true})
//
//	// Test service URL
//	svcParts := azfile.NewFileURLParts(fsu.URL())
//	svcParts.SAS = sasQueryParams
//	testSvcURL := svcParts.URL()
//	svcURLWithSAS := azfile.NewServiceURL(testSvcURL, azfile.NewPipeline(azfile.NewAnonymousCredential(), azfile.PipelineOptions{}))
//	// List
//	_, err = svcURLWithSAS.ListSharesSegment(ctx, azfile.Marker{}, azfile.ListSharesOptions{})
//	_require.Nil(err)
//	// Write
//	_, err = svcURLWithSAS.SetProperties(ctx, azfile.FileServiceProperties{})
//	_require.Nil(err)
//	// Read
//	_, err = svcURLWithSAS.GetProperties(ctx)
//	_require.Nil(err)
//
//	// Test share URL
//	sParts := azfile.NewFileURLParts(shareURL.URL())
//	_require(sParts.ShareName, chk.Equals, shareName)
//	sParts.SAS = sasQueryParams
//	testShareURL := sParts.URL()
//	shareURLWithSAS := azfile.NewShareURL(testShareURL, azfile.NewPipeline(azfile.NewAnonymousCredential(), azfile.PipelineOptions{}))
//	// Create
//	_, err = shareURLWithSAS.Create(ctx, azfile.Metadata{}, 0)
//	_require.Nil(err)
//	// Write
//	metadata := azfile.Metadata{"foo": "bar"}
//	_, err = shareURLWithSAS.SetMetadata(ctx, metadata)
//	// Read
//	gResp, err := shareURLWithSAS.GetProperties(ctx)
//	_require.Nil(err)
//	_require(gResp.NewMetadata(), chk.DeepEquals, metadata)
//	// Delete
//	defer shareURLWithSAS.Delete(ctx, azfile.DeleteSnapshotsOptionNone)
//
//	// Test dir URL
//	dParts := azfile.NewFileURLParts(dirURL.URL())
//	dParts.SAS = sasQueryParams
//	testDirURL := dParts.URL()
//	dirURLWithSAS := azfile.NewDirectoryURL(testDirURL, azfile.NewPipeline(azfile.NewAnonymousCredential(), azfile.PipelineOptions{}))
//	// Create
//	_, err = dirURLWithSAS.Create(ctx, azfile.Metadata{}, azfile.SMBProperties{})
//	_require.Nil(err)
//	// Write
//	_, err = dirURLWithSAS.SetMetadata(ctx, metadata)
//	// Read
//	gdResp, err := dirURLWithSAS.GetProperties(ctx)
//	_require.Nil(err)
//	_require(gdResp.NewMetadata(), chk.DeepEquals, metadata)
//	// List
//	_, err = dirURLWithSAS.ListFilesAndDirectoriesSegment(ctx, azfile.Marker{}, azfile.ListFilesAndDirectoriesOptions{})
//	_require.Nil(err)
//	// Delete
//	defer dirURLWithSAS.Delete(ctx)
//
//	// Test file URL
//	fParts := azfile.NewFileURLParts(fileURL.URL())
//	fParts.SAS = sasQueryParams
//	testFileURL := fParts.URL()
//	fileURLWithSAS := azfile.NewFileURL(testFileURL, azfile.NewPipeline(azfile.NewAnonymousCredential(), azfile.PipelineOptions{}))
//	// Create
//	_, err = fileURLWithSAS.Create(ctx, 0, azfile.FileHTTPHeaders{}, azfile.Metadata{})
//	_require.Nil(err)
//	// Write
//	_, err = fileURLWithSAS.SetMetadata(ctx, metadata)
//	// Read
//	gfResp, err := fileURLWithSAS.GetProperties(ctx)
//	_require.Nil(err)
//	_require(gfResp.NewMetadata(), chk.DeepEquals, metadata)
//	// Delete
//	defer fileURLWithSAS.Delete(ctx)
//}
