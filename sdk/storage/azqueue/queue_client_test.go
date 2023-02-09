//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/queueerror"
	"github.com/stretchr/testify/require"
	"strconv"
	"time"
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
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestQueueSetMetadataNilOptions() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	_, err = queueClient.SetMetadata(context.Background(), nil)
	_require.Nil(err)

	_, err = queueClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestQueueSetEmptyACL() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.SetAccessPolicyOptions{QueueACL: nil}
	_, err = queueClient.SetAccessPolicy(context.Background(), &opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestQueueSetACLNil() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	_, err = queueClient.SetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestQueueSetEmptyACL2() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	sI := make([]*azqueue.SignedIdentifier, 0)
	opts := azqueue.SetAccessPolicyOptions{QueueACL: sI}
	_, err = queueClient.SetAccessPolicy(context.Background(), &opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestQueueSetBasicACL() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"

	signedIdentifiers := make([]*azqueue.SignedIdentifier, 0)

	signedIdentifiers = append(signedIdentifiers, &azqueue.SignedIdentifier{
		AccessPolicy: &azqueue.AccessPolicy{
			Expiry:     &expiration,
			Start:      &start,
			Permission: &permission,
		},
		ID: &id,
	})
	options := azqueue.SetAccessPolicyOptions{QueueACL: signedIdentifiers}
	_, err = queueClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestQueueSetMultipleACL() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	id := "empty"

	signedIdentifiers := make([]*azqueue.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &azqueue.SignedIdentifier{
		ID: &id,
	})

	permission2 := "r"
	id2 := "partial"

	signedIdentifiers = append(signedIdentifiers, &azqueue.SignedIdentifier{
		ID: &id2,
		AccessPolicy: &azqueue.AccessPolicy{
			Permission: &permission2,
		},
	})

	id3 := "full"
	permission3 := "r"
	start := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)
	expiry := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)

	signedIdentifiers = append(signedIdentifiers, &azqueue.SignedIdentifier{
		ID: &id3,
		AccessPolicy: &azqueue.AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission3,
		},
	})
	options := azqueue.SetAccessPolicyOptions{QueueACL: signedIdentifiers}
	_, err = queueClient.SetAccessPolicy(context.Background(), &options)
	_require.Nil(err)

	// Make a Get to assert two access policies
	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 3)
}

func (s *RecordedTestSuite) TestQueueGetMultipleACL() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry := start.Add(5 * time.Minute)
	expiry2 := start.Add(time.Minute)
	readWrite := to.Ptr(azqueue.AccessPolicyPermission{Read: true, Update: true}).String()
	readOnly := to.Ptr(azqueue.AccessPolicyPermission{Read: true}).String()
	id1, id2 := "0000", "0001"
	permissions := []*azqueue.SignedIdentifier{
		{ID: &id1,
			AccessPolicy: &azqueue.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &readWrite,
			},
		},
		{ID: &id2,
			AccessPolicy: &azqueue.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry2,
				Permission: &readOnly,
			},
		},
	}
	options := azqueue.SetAccessPolicyOptions{QueueACL: permissions}
	_, err = queueClient.SetAccessPolicy(context.Background(), &options)

	_require.Nil(err)

	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *RecordedTestSuite) TestQueueSetACLMoreThanFive() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*azqueue.SignedIdentifier, 6)
	listOnly := to.Ptr(azqueue.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 6; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &azqueue.SignedIdentifier{
			ID: &id,
			AccessPolicy: &azqueue.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := azqueue.SetAccessPolicyOptions{
		QueueACL: permissions,
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.InvalidXMLDocument)
}

func (s *RecordedTestSuite) TestQueueSetPermissionsDeleteAndModifyACL() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	listOnly := to.Ptr(azqueue.AccessPolicyPermission{Read: true}).String()
	permissions := make([]*azqueue.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &azqueue.SignedIdentifier{
			ID: &id,
			AccessPolicy: &azqueue.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := azqueue.SetAccessPolicyOptions{
		QueueACL: permissions,
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := azqueue.SetAccessPolicyOptions{
		QueueACL: permissions,
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions1)
	_require.Nil(err)

	resp, err = queueClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 1)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *RecordedTestSuite) TestQueueSetPermissionsDeleteAllPolicies() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*azqueue.SignedIdentifier, 2)
	listOnly := to.Ptr(azqueue.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &azqueue.SignedIdentifier{
			ID: &id,
			AccessPolicy: &azqueue.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := azqueue.SetAccessPolicyOptions{
		QueueACL: permissions,
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = azqueue.SetAccessPolicyOptions{
		QueueACL: []*azqueue.SignedIdentifier{},
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err = queueClient.GetAccessPolicy(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.SignedIdentifiers)
}

func (s *RecordedTestSuite) TestQueueSetPermissionsInvalidPolicyTimes() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*azqueue.SignedIdentifier, 2)
	listOnly := to.Ptr(azqueue.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &azqueue.SignedIdentifier{
			ID: &id,
			AccessPolicy: &azqueue.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := azqueue.SetAccessPolicyOptions{
		QueueACL: permissions,
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestQueueSetPermissionsSignedIdentifierTooLong() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]*azqueue.SignedIdentifier, 2)
	listOnly := to.Ptr(azqueue.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		permissions[i] = &azqueue.SignedIdentifier{
			ID: &id,
			AccessPolicy: &azqueue.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	setAccessPolicyOptions := azqueue.SetAccessPolicyOptions{
		QueueACL: permissions,
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NotNil(err)

	testcommon.ValidateQueueErrorCode(_require, err, queueerror.InvalidXMLDocument)
}

func (s *RecordedTestSuite) TestEnqueueMessageBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)

	}
}

func (s *RecordedTestSuite) TestEnqueueMessageWithTimeToLive() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestEnqueueMessageWithTimeToLiveExpired() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.Nil(err)

	time.Sleep(time.Second * 2)
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(0, len(resp.QueueMessagesList))
}

func (s *RecordedTestSuite) TestEnqueueMessageWithInfiniteTimeToLive() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(-1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestEnqueueMessageWithVisibilityTimeout() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.EnqueueMessageOptions{VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestEnqueueMessageWithVisibilityTimeoutSmallerThanTTL() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(2)), VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestEnqueueMessageWithVisibilityTimeoutLargerThanTTL() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(1)), VisibilityTimeout: to.Ptr(int32(2))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	// cannot have visibility timeout be greater than ttl
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.InvalidQueryParameterValue)
}

func (s *RecordedTestSuite) TestDequeueMessageBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// dequeue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.DequeueMessage(context.Background(), nil)
		_require.Nil(err)
		_require.Equal(1, len(resp.QueueMessagesList))
		_require.NotNil(resp.QueueMessagesList[0].MessageID)
	}
	// should be 0 now
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Equal(0, len(resp.QueueMessagesList))
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestDequeueMessageWithVisibilityTimeout() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	opts := azqueue.DequeueMessageOptions{VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.Nil(err)

	resp, err := queueClient.DequeueMessage(context.Background(), &opts)
	_require.Nil(err)
	_require.NotNil(resp.QueueMessagesList[0].TimeNextVisible)
}

func (s *RecordedTestSuite) TestDequeueMessagesBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// dequeue 4 messages
	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	resp, err := queueClient.DequeueMessages(context.Background(), &opts)
	_require.Nil(err)
	_require.Equal(4, len(resp.QueueMessagesList))
}

func (s *RecordedTestSuite) TestDequeueMessagesDefault() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// should dequeue only 1 message (since default num of messages is 1 when not specified)
	resp, err := queueClient.DequeueMessages(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(1, len(resp.QueueMessagesList))
}

func (s *RecordedTestSuite) TestDequeueMessagesWithVisibilityTimeout() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// dequeue 4 messages
	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4)), VisibilityTimeout: to.Ptr(int32(2))}
	_, err = queueClient.DequeueMessages(context.Background(), &opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestDequeueMessagesWithNumMessagesLargerThan32() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 33 messages
	for i := 0; i < 33; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(33))}
	_, err = queueClient.DequeueMessages(context.Background(), &opts)
	// should fail
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.OutOfRangeQueryParameterValue)
}

func (s *RecordedTestSuite) TestDequeueMessagesWithLeftovers() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 10 messages
	for i := 0; i < 10; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// dequeue 5 messages
	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(5))}
	resp, err := queueClient.DequeueMessages(context.Background(), &opts)
	_require.Nil(err)
	_require.Equal(*resp.QueueMessagesList[0].MessageText, testcommon.QueueDefaultData)
	_require.Equal(5, len(resp.QueueMessagesList))

	// dequeue other 5 messages
	resp, err = queueClient.DequeueMessages(context.Background(), &opts)
	_require.Nil(err)
	_require.Equal(*resp.QueueMessagesList[0].MessageText, testcommon.QueueDefaultData)
	_require.Equal(5, len(resp.QueueMessagesList))
}

func (s *RecordedTestSuite) TestPeekMessageBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// peek 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.PeekMessage(context.Background(), nil)
		_require.Nil(err)
		_require.Equal(1, len(resp.QueueMessagesList))
		_require.NotNil(resp.QueueMessagesList[0].MessageID)
		_require.Equal(*resp.QueueMessagesList[0].MessageText, testcommon.QueueDefaultData)
	}

	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	// should all still be there
	resp, err := queueClient.DequeueMessages(context.Background(), &opts)
	_require.Equal(4, len(resp.QueueMessagesList))
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestPeekMessagesBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// dequeue 4 messages
	opts := azqueue.PeekMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	resp, err := queueClient.PeekMessages(context.Background(), &opts)
	_require.Nil(err)
	_require.Equal(4, len(resp.QueueMessagesList))

	opts1 := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	// should all still be there
	resp1, err := queueClient.DequeueMessages(context.Background(), &opts1)
	_require.Equal(4, len(resp1.QueueMessagesList))
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestPeekMessagesDefault() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// should peek only 1 message (since default num of messages is 1 when not specified)
	resp, err := queueClient.PeekMessages(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(1, len(resp.QueueMessagesList))
	_require.Equal(*resp.QueueMessagesList[0].MessageText, testcommon.QueueDefaultData)
}

func (s *RecordedTestSuite) TestPeekMessagesWithNumMessagesLargerThan32() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 33 messages
	for i := 0; i < 33; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	opts := azqueue.PeekMessagesOptions{NumberOfMessages: to.Ptr(int32(33))}
	_, err = queueClient.PeekMessages(context.Background(), &opts)
	// should fail
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.OutOfRangeQueryParameterValue)
}

func (s *RecordedTestSuite) TestDeleteMessageBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	var popReceipts []string
	var messageIDs []string
	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
		popReceipts = append(popReceipts, *resp.QueueMessagesList[0].PopReceipt)
		messageIDs = append(messageIDs, *resp.QueueMessagesList[0].MessageID)
	}

	// delete 4 messages
	for i := 0; i < 4; i++ {
		opts := &azqueue.DeleteMessageOptions{}
		_, err := queueClient.DeleteMessage(context.Background(), messageIDs[i], popReceipts[i], opts)
		_require.Nil(err)
	}
	// should be 0 now
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Equal(0, len(resp.QueueMessagesList))
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestDeleteMessageNilOptions() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	var popReceipts []string
	var messageIDs []string
	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
		popReceipts = append(popReceipts, *resp.QueueMessagesList[0].PopReceipt)
		messageIDs = append(messageIDs, *resp.QueueMessagesList[0].MessageID)
	}

	// delete 4 messages
	for i := 0; i < 4; i++ {
		_, err := queueClient.DeleteMessage(context.Background(), messageIDs[i], popReceipts[i], nil)
		_require.Nil(err)
	}
	// should be 0 now
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Equal(0, len(resp.QueueMessagesList))
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestDeleteMessageDoesNotExist() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.Nil(err)
	popReceipt := *resp.QueueMessagesList[0].PopReceipt
	messageID := *resp.QueueMessagesList[0].MessageID

	opts := &azqueue.DeleteMessageOptions{}
	_, err = queueClient.DeleteMessage(context.Background(), messageID, popReceipt, opts)
	_require.Nil(err)

	// should fail since we already deleted it
	_, err = queueClient.DeleteMessage(context.Background(), messageID, popReceipt, opts)
	_require.NotNil(err)
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.MessageNotFound)
}

func (s *RecordedTestSuite) TestClearMessagesBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// delete the queue's messages
	opts := azqueue.ClearMessagesOptions{}
	_, err = queueClient.ClearMessages(context.Background(), &opts)
	_require.Nil(err)

	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(0, len(resp.QueueMessagesList))
}

func (s *RecordedTestSuite) TestClearMessagesNilOptions() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// delete the queue's messages
	_, err = queueClient.ClearMessages(context.Background(), nil)
	_require.Nil(err)

	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(0, len(resp.QueueMessagesList))
}

func (s *RecordedTestSuite) TestClearMessagesMoreThan32() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	// enqueue 33 messages
	for i := 0; i < 33; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.Nil(err)
	}

	// delete the queue's messages
	opts := azqueue.ClearMessagesOptions{}
	_, err = queueClient.ClearMessages(context.Background(), &opts)
	_require.Nil(err)

	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(0, len(resp.QueueMessagesList))
}

func (s *RecordedTestSuite) TestUpdateMessageBasic() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.Nil(err)
	popReceipt := *resp.QueueMessagesList[0].PopReceipt
	messageID := *resp.QueueMessagesList[0].MessageID

	opts := &azqueue.UpdateMessageOptions{}
	_, err = queueClient.UpdateMessage(context.Background(), messageID, popReceipt, "new content", opts)
	_require.Nil(err)

	resp1, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Nil(err)
	content := *resp1.QueueMessagesList[0].MessageText
	_require.Equal("new content", content)
}

func (s *RecordedTestSuite) TestUpdateMessageWithVisibilityTimeout() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.Nil(err)

	resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.Nil(err)
	popReceipt := *resp.QueueMessagesList[0].PopReceipt
	messageID := *resp.QueueMessagesList[0].MessageID

	opts := &azqueue.UpdateMessageOptions{VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.UpdateMessage(context.Background(), messageID, popReceipt, "new content", opts)
	_require.Nil(err)
	time.Sleep(time.Second * 2)
	resp1, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Nil(err)
	content := *resp1.QueueMessagesList[0].MessageText
	_require.Equal("new content", content)
}

//TODO: TestSAS
