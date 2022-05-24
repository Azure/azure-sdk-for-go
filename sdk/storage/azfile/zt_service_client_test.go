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

func (s *azfileLiveTestSuite) TestAccountSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)

	dirClient := getDirectoryClientFromShare(_require, "dir1", srClient)
	fClient := getFileClientFromDirectory(_require, "file1", dirClient)

	svcURLWithSAS, err := svcClient.GetSASURL(
		AccountSASResourceTypes{Service: true, Container: true, Object: true},
		AccountSASPermissions{Read: true, List: true, Write: true, Delete: true, Add: true, Create: true, Update: true, Process: true},
		time.Now(),
		time.Now().Add(1*time.Hour),
	)

	//// Reverse validation all parse logics work as expect.
	//ap := &AccountSASPermissions{}
	//err = ap.Parse(sasQueryParams.Permissions())
	//_require.Nil(err)
	//_require.EqualValues(*ap, AccountSASPermissions{Read: true, List: true, Write: true, Delete: true, Add: true, Create: true, Update: true, Process: true})
	//
	//as := &AccountSASServices{}
	//err = as.Parse(sasQueryParams.Services())
	//_require.Nil(err)
	//_require.EqualValues(*as, AccountSASServices{File: true, Blob: true, Queue: true})
	//
	//ar := &AccountSASResourceTypes{}
	//err = ar.Parse(sasQueryParams.ResourceTypes())
	//_require.Nil(err)
	//_require.EqualValues(*ar, AccountSASResourceTypes{Service: true, Container: true, Object: true})

	svcClientWithSAS, err := NewServiceClientWithNoCredential(svcURLWithSAS, nil)

	// List
	pager := svcClientWithSAS.ListShares(nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.Nil(err)
		if err != nil {
			break
		}
	}

	// Write
	_, err = svcClientWithSAS.SetProperties(ctx, &ServiceSetPropertiesOptions{})
	_require.Nil(err)

	// Read
	_, err = svcClientWithSAS.GetProperties(ctx, &ServiceGetPropertiesOptions{})
	_require.Nil(err)

	// Test share URL
	srURLWithSAS, err := srClient.GetSASURL(
		ShareSASPermissions{Read: true, List: true, Write: true, Delete: true},
		time.Now(),
		time.Now().Add(1*time.Hour),
	)

	shareURLWithSAS, err := NewShareClientWithNoCredential(srURLWithSAS, nil)
	_require.Nil(err)

	// Create
	_, err = shareURLWithSAS.Create(ctx, &ShareCreateOptions{})
	_require.Nil(err)

	// Write
	_, err = shareURLWithSAS.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	// Read
	gResp, err := shareURLWithSAS.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(gResp.Metadata, basicMetadata)

	// Delete
	defer shareURLWithSAS.Delete(ctx, nil)

	// Test dir URL
	dirURLWithSAS, err := dirClient.GetSASURL(FileSASPermissions{Read: true, Create: true, Write: true, Delete: true}, time.Now(), time.Now().Add(1*time.Hour))
	_require.Nil(err)

	dirClientWithSAS, err := NewDirectoryClientWithNoCredential(dirURLWithSAS, nil)
	_require.Nil(err)

	// Create
	_, err = dirClientWithSAS.Create(ctx, nil)
	_require.Nil(err)

	// Write
	_, err = dirClientWithSAS.SetMetadata(ctx, basicMetadata, nil)
	_require.Nil(err)

	// Read
	gdResp, err := dirClientWithSAS.GetProperties(ctx, nil)
	_require.Nil(err)

	_require.EqualValues(gdResp.Metadata, basicMetadata)

	// List
	pager2 := dirClientWithSAS.ListFilesAndDirectories(nil)
	for pager2.More() {
		_, err := pager2.NextPage(ctx)
		_require.Nil(err)
		if err != nil {
			break
		}
	}

	// Delete
	defer dirClientWithSAS.Delete(ctx, nil)

	// Test file URL
	fileURLWithSAS, err := fClient.GetSASURL(FileSASPermissions{Read: true, Create: true, Write: true, Delete: true}, time.Now(), time.Now().Add(1*time.Hour))
	_require.Nil(err)

	fileClientWithSAS, err := NewFileClientWithNoCredential(fileURLWithSAS, nil)
	_require.Nil(err)

	// Create
	_, err = fileClientWithSAS.Create(ctx, nil)
	_require.Nil(err)

	// Write
	_, err = fileClientWithSAS.SetMetadata(ctx, basicMetadata, nil)
	// Read
	gfResp, err := fileClientWithSAS.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(gfResp.Metadata, basicMetadata)
	// Delete
	fileClientWithSAS.Delete(ctx, nil)
}
