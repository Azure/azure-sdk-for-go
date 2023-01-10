package service_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
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

//TODO: TestListQueues
//TODO: TestCreateQueue
//TODO: TestDeleteQueue
//TODO: TestSAS...

func (s *ServiceUnrecordedTestsSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())
	//testName := s.T().Name()

	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)
	connectionString := testcommon.GetConnectionString(testcommon.TestAccountDefault)

	parsedConnStr, err := shared.ParseConnectionString(connectionString)
	_require.Nil(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".queue.core.windows.net/")

	sharedKeyCred, err := azqueue.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.Nil(err)

	svcClient, err := service.NewClientWithSharedKeyCredential(parsedConnStr.ServiceURL, sharedKeyCred, nil)
	_require.Nil(err)
	//containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	//defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	sProps, err := svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(sProps)
}

func (s *ServiceRecordedTestsSuite) TestGetProperties() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sProps, err := svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(sProps)
}
