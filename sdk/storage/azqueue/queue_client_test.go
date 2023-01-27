//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/testcommon"
	"github.com/stretchr/testify/require"
)

func (s *RecordedTestSuite) TestQueueCreateQueue() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	resp, err := queueClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(resp)
}

func (s *RecordedTestSuite) TestQueueCreateQueueWithMetadata() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	opts := azqueue.CreateOptions{Metadata: testcommon.BasicMetadata}
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	resp, err := queueClient.Create(context.Background(), &opts)
	_require.Nil(err)
	_require.NotZero(resp)
}

func (s *RecordedTestSuite) TestQueueDeleteQueue() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	resp, err := queueClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(resp)

	delResp, err := queueClient.Delete(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(delResp)
}

func (s *RecordedTestSuite) TestQueueSetMetadata() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.SetMetadataOptions{Metadata: testcommon.BasicMetadata}
	_, err = queueClient.SetMetadata(context.Background(), &opts)
	_require.Nil(err)

	resp, err := queueClient.GetProperties(context.Background(), nil)
	_require.Equal(resp.Metadata, testcommon.BasicMetadata)
}

//TODO: TestCreateQueue
//TODO: TestDeleteQueue
//TODO: TestSetQueueMetadata
//TODO: TestGetQueueMetadata
//TODO: TestSetQueueACL
//TODO: TestGetQueueACL
//TODO: TestPutMessage
//TODO: TestGetMessages
//TODO: TestPeekMessages
//TODO: TestDeleteMessage
//TODO: TestClearMessages
//TODO: TestUpdateMessage
