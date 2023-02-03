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

//TODO: TestPutMessage
//TODO: TestGetMessages
//TODO: TestPeekMessages
//TODO: TestDeleteMessage
//TODO: TestClearMessages
//TODO: TestUpdateMessage
