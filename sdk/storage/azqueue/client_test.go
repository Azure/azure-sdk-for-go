//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/testcommon"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running service Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &RecordedTestSuite{})
		suite.Run(t, &UnrecordedTestSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &RecordedTestSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &RecordedTestSuite{})
	}
}

func (s *RecordedTestSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *RecordedTestSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *UnrecordedTestSuite) BeforeTest(suite string, test string) {

}

func (s *UnrecordedTestSuite) AfterTest(suite string, test string) {

}

type RecordedTestSuite struct {
	suite.Suite
}

type UnrecordedTestSuite struct {
	suite.Suite
}

func (s *UnrecordedTestSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())
	//testName := s.T().Name()

	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)
	connectionString := testcommon.GetConnectionString(testcommon.TestAccountDefault)

	parsedConnStr, err := shared.ParseConnectionString(connectionString)
	_require.Nil(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".queue.core.windows.net/")

	sharedKeyCred, err := azqueue.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.Nil(err)

	svcClient, err := azqueue.NewServiceClientWithSharedKeyCredential(parsedConnStr.ServiceURL, sharedKeyCred, nil)
	_require.Nil(err)

	sProps, err := svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(sProps)
}

func (s *RecordedTestSuite) TestGetProperties() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sProps, err := svcClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(sProps)
}

func (s *RecordedTestSuite) TestCreateQueue() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := svcClient.NewQueueClient(queueName)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	resp, err := svcClient.CreateQueue(context.Background(), queueName, nil)
	_require.Nil(err)
	_require.NotZero(resp)
}

func (s *RecordedTestSuite) TestCreateQueueWithMetadata() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := svcClient.NewQueueClient(queueName)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)
	opts := azqueue.CreateOptions{Metadata: testcommon.BasicMetadata}

	resp, err := svcClient.CreateQueue(context.Background(), queueName, &opts)
	_require.Nil(err)
	_require.NotZero(resp)
}

func (s *RecordedTestSuite) TestDeleteQueue() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	createResp, err := svcClient.CreateQueue(context.Background(), queueName, nil)
	_require.Nil(err)
	_require.NotZero(createResp)

	delResp, err := svcClient.DeleteQueue(context.Background(), queueName, nil)
	_require.Nil(err)
	_require.NotZero(delResp)
}

func (s *RecordedTestSuite) TestListQueuesWithMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)
	md := map[string]*string{
		"foo": to.Ptr("foovalue"),
		"bar": to.Ptr("barvalue"),
	}

	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	_, err = queueClient.Create(context.Background(), &azqueue.CreateOptions{Metadata: md})
	defer func(queueClient *azqueue.QueueClient, ctx context.Context, options *azqueue.DeleteOptions) {
		_, err := queueClient.Delete(ctx, options)
		if err != nil {
			_require.Nil(err)
		}
	}(queueClient, context.Background(), nil)
	_require.Nil(err)
	listOptions := azqueue.ListQueuesOptions{Include: azqueue.ListQueuesInclude{Metadata: true}}
	pager := svcClient.NewListQueuesPager(&listOptions)

	exists := false
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, queue := range resp.QueueItems {
			_require.NotNil(queue.Name)
			if *queue.Name == queueName {
				_require.NotNil(queue.Metadata)
				unwrappedMeta := map[string]*string{}
				for k, v := range queue.Metadata {
					if v != nil {
						unwrappedMeta[k] = v
					}
				}
				_require.EqualValues(unwrappedMeta, md)
				exists = true
			}
		}
		if err != nil {
			break
		}
	}

	_require.Nil(err)
	_require.True(exists)
}

func (s *RecordedTestSuite) TestListQueuesPagination() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)

	queueName := testcommon.GenerateQueueName(testName)
	queueNames := []string{queueName + "1", queueName + "2", queueName + "3"}

	for i := 0; i < len(queueNames); i++ {
		_, err := svcClient.CreateQueue(context.Background(), queueNames[i], nil)
		if err != nil {
			_require.Nil(err)
		}
	}
	// cleanup created queues
	defer func(queueNames []string, ctx context.Context, options *azqueue.DeleteOptions) {
		for i := 0; i < len(queueNames); i++ {
			_, err := svcClient.DeleteQueue(ctx, queueNames[i], nil)
			if err != nil {
				_require.Nil(err)
			}
		}
	}(queueNames, context.Background(), nil)
	_require.Nil(err)
	pager := svcClient.NewListQueuesPager(&azqueue.ListQueuesOptions{MaxResults: to.Ptr(int32(1))})

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		for _, queue := range resp.QueueItems {
			_require.NotNil(queue.Name)
			count += 1
		}
		if err != nil {
			break
		}
	}
	_require.GreaterOrEqual(count, 3)
}

//TODO: TestSAS
