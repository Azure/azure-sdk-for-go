//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
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

func (s *ServiceRecordedTestsSuite) TestAccountNewServiceURLValidName() {
	_require := require.New(s.T())

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + "." + testcommon.DefaultFileEndpointSuffix
	_require.Equal(svcClient.URL(), correctURL)
}

func (s *ServiceRecordedTestsSuite) TestAccountNewShareURLValidName() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)
	_require.NoError(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + "." + testcommon.DefaultFileEndpointSuffix + shareName
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

func (s *ServiceRecordedTestsSuite) TestAccountListSharesNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	mySharePrefix := testcommon.GenerateEntityName(testName)
	pager := svcClient.NewListSharesPager(&service.ListSharesOptions{Prefix: to.Ptr(mySharePrefix)})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.NotNil(resp.Prefix)
		_require.Equal(*resp.Prefix, mySharePrefix)
		_require.NotNil(resp.ServiceEndpoint)
		_require.NotNil(resp.Version)
		_require.Len(resp.Shares, 0)
	}

	/*shareClients := map[string]*share.Client{}
	for i := 0; i < 4; i++ {
		shareName := mySharePrefix + "share" + strconv.Itoa(i)
		shareClients[shareName] = createNewShare(_require, shareName, svcClient)

		_, err := shareClients[shareName].SetMetadata(context.Background(), basicMetadata, nil)
		_require.NoError(err)

		_, err = shareClients[shareName].CreateSnapshot(context.Background(), nil)
		_require.NoError(err)

		defer delShare(_require, shareClients[shareName], &ShareDeleteOptions{
			DeleteSnapshots: to.Ptr(DeleteSnapshotsOptionTypeInclude),
		})
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
			_require.Len(shareItem.Metadata, len(basicMetadata))
			for key, val1 := range basicMetadata {
				if val2, ok := shareItem.Metadata[key]; !(ok && val1 == *val2) {
					_require.Fail("metadata mismatch")
				}
			}
			_require.NotNil(resp.Shares[0].Snapshot)
			_require.Nil(resp.Shares[1].Snapshot)
		}
	}*/
}
