// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/queueerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/sas"
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
	_require.NoError(err)
	_require.NotZero(resp)
}

func (s *UnrecordedTestSuite) TestQueueClientFromConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)
	connectionString := testcommon.GetConnectionString(testcommon.TestAccountDefault)

	parsedConnStr, err := shared.ParseConnectionString(connectionString)
	_require.NoError(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".queue.core.windows.net/")

	queueName := testcommon.GenerateQueueName(testName)

	sharedKeyCred, err := azqueue.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.NoError(err)

	qClient, err := azqueue.NewQueueClientWithSharedKeyCredential(
		runtime.JoinPaths(parsedConnStr.ServiceURL, queueName), sharedKeyCred, nil)
	_require.NoError(err)

	_, err = qClient.Create(context.Background(), nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestQueueClientFromConnectionString1() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)
	connectionString := testcommon.GetConnectionString(testcommon.TestAccountDefault)

	parsedConnStr, err := shared.ParseConnectionString(connectionString)
	_require.NoError(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".queue.core.windows.net/")

	queueName := testcommon.GenerateQueueName(testName)

	qClient, err := azqueue.NewQueueClientFromConnectionString(connectionString, queueName, nil)
	_require.NoError(err)

	_, err = qClient.Create(context.Background(), nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestQueueClientUsingOauth() {
	_require := require.New(s.T())
	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)

	queueURL := fmt.Sprintf("https://%s.queue.core.windows.net/%s", accountName, queueName)
	cred, err := credential.New(nil)
	_require.NoError(err)
	qClient, err := azqueue.NewQueueClient(queueURL, cred, nil)
	_require.NoError(err)

	_, err = qClient.Create(context.Background(), nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestQueueClientUsingOauthWithCustomAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)

	queueURL := fmt.Sprintf("https://%s.queue.core.windows.net/%s", accountName, queueName)
	cred, err := credential.New(nil)
	_require.NoError(err)

	options := &azqueue.ClientOptions{Audience: "https://" + accountName + ".queue.core.windows.net"}
	options.Logging.AllowedHeaders = append(options.Logging.AllowedHeaders, "X-Request-Mismatch", "X-Request-Mismatch-Error")
	transport, err := recording.NewRecordingHTTPClient(s.T(), nil)
	require.NoError(s.T(), err)

	options.Transport = transport

	qClient, err := azqueue.NewQueueClient(queueURL, cred, nil)
	_require.NoError(err)

	_, err = qClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
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
	_require.NoError(err)
	_require.NotZero(resp)

	delResp, err := queueClient.Delete(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	opts := azqueue.SetMetadataOptions{Metadata: testcommon.BasicMetadata}
	_, err = queueClient.SetMetadata(context.Background(), &opts)
	_require.NoError(err)

	resp, err := queueClient.GetProperties(context.Background(), nil)
	_require.Equal(resp.Metadata, testcommon.BasicMetadata)
	_require.NoError(err)
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
	_require.NoError(err)

	_, err = queueClient.SetMetadata(context.Background(), nil)
	_require.NoError(err)

	_, err = queueClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	opts := azqueue.SetAccessPolicyOptions{QueueACL: nil}
	_, err = queueClient.SetAccessPolicy(context.Background(), &opts)
	_require.NoError(err)
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
	_require.NoError(err)

	_, err = queueClient.SetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	sI := make([]*azqueue.SignedIdentifier, 0)
	opts := azqueue.SetAccessPolicyOptions{QueueACL: sI}
	_, err = queueClient.SetAccessPolicy(context.Background(), &opts)
	_require.NoError(err)
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
	_require.NoError(err)

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
	_require.NoError(err)
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
	_require.NoError(err)

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
	_require.NoError(err)

	// Make a Get to assert two access policies
	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
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

	_require.NoError(err)

	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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
	_require.Error(err)
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
	_require.NoError(err)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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
	_require.NoError(err)

	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := azqueue.SetAccessPolicyOptions{
		QueueACL: permissions,
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions1)
	_require.NoError(err)

	resp, err = queueClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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
	_require.NoError(err)

	resp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = azqueue.SetAccessPolicyOptions{
		QueueACL: []*azqueue.SignedIdentifier{},
	}
	_, err = queueClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_require.NoError(err)

	resp, err = queueClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)
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
	_require.NoError(err)
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
	_require.NoError(err)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.NoError(err)
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
	_require.Error(err)

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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)

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
	_require.NoError(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.NoError(err)
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
	_require.NoError(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.NoError(err)

	time.Sleep(time.Second * 2)
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(0, len(resp.Messages))
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
	_require.NoError(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(-1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.NoError(err)
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
	_require.NoError(err)

	opts := azqueue.EnqueueMessageOptions{VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.NoError(err)
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
	_require.NoError(err)

	opts := azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(2)), VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, &opts)
	_require.NoError(err)
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
	_require.NoError(err)

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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// dequeue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.DequeueMessage(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(1, len(resp.Messages))
		_require.NotNil(resp.Messages[0].MessageID)
	}
	// should be 0 now
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Equal(0, len(resp.Messages))
	_require.NoError(err)
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
	_require.NoError(err)

	opts := azqueue.DequeueMessageOptions{VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.NoError(err)

	resp, err := queueClient.DequeueMessage(context.Background(), &opts)
	_require.NoError(err)
	_require.NotNil(resp.Messages[0].TimeNextVisible)
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// dequeue 4 messages
	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	resp, err := queueClient.DequeueMessages(context.Background(), &opts)
	_require.NoError(err)
	_require.Equal(4, len(resp.Messages))
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// should dequeue only 1 message (since default num of messages is 1 when not specified)
	resp, err := queueClient.DequeueMessages(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(1, len(resp.Messages))
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// dequeue 4 messages
	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4)), VisibilityTimeout: to.Ptr(int32(2))}
	_, err = queueClient.DequeueMessages(context.Background(), &opts)
	_require.NoError(err)
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
	_require.NoError(err)

	// enqueue 33 messages
	for i := 0; i < 33; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
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
	_require.NoError(err)

	// enqueue 10 messages
	for i := 0; i < 10; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// dequeue 5 messages
	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(5))}
	resp, err := queueClient.DequeueMessages(context.Background(), &opts)
	_require.NoError(err)
	_require.Equal(*resp.Messages[0].MessageText, testcommon.QueueDefaultData)
	_require.Equal(5, len(resp.Messages))

	// dequeue other 5 messages
	resp, err = queueClient.DequeueMessages(context.Background(), &opts)
	_require.NoError(err)
	_require.Equal(*resp.Messages[0].MessageText, testcommon.QueueDefaultData)
	_require.Equal(5, len(resp.Messages))
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// peek 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.PeekMessage(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(1, len(resp.Messages))
		_require.NotNil(resp.Messages[0].MessageID)
		_require.Equal(*resp.Messages[0].MessageText, testcommon.QueueDefaultData)
	}

	opts := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	// should all still be there
	resp, err := queueClient.DequeueMessages(context.Background(), &opts)
	_require.Equal(4, len(resp.Messages))
	_require.NoError(err)
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// dequeue 4 messages
	opts := azqueue.PeekMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	resp, err := queueClient.PeekMessages(context.Background(), &opts)
	_require.NoError(err)
	_require.Equal(4, len(resp.Messages))

	opts1 := azqueue.DequeueMessagesOptions{NumberOfMessages: to.Ptr(int32(4))}
	// should all still be there
	resp1, err := queueClient.DequeueMessages(context.Background(), &opts1)
	_require.Equal(4, len(resp1.Messages))
	_require.NoError(err)
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// should peek only 1 message (since default num of messages is 1 when not specified)
	resp, err := queueClient.PeekMessages(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(1, len(resp.Messages))
	_require.Equal(*resp.Messages[0].MessageText, testcommon.QueueDefaultData)
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
	_require.NoError(err)

	// enqueue 33 messages
	for i := 0; i < 33; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
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
	_require.NoError(err)

	var popReceipts []string
	var messageIDs []string
	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
		popReceipts = append(popReceipts, *resp.Messages[0].PopReceipt)
		messageIDs = append(messageIDs, *resp.Messages[0].MessageID)
	}

	// delete 4 messages
	for i := 0; i < 4; i++ {
		opts := &azqueue.DeleteMessageOptions{}
		_, err := queueClient.DeleteMessage(context.Background(), messageIDs[i], popReceipts[i], opts)
		_require.NoError(err)
	}
	// should be 0 now
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Equal(0, len(resp.Messages))
	_require.NoError(err)
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
	_require.NoError(err)

	var popReceipts []string
	var messageIDs []string
	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
		popReceipts = append(popReceipts, *resp.Messages[0].PopReceipt)
		messageIDs = append(messageIDs, *resp.Messages[0].MessageID)
	}

	// delete 4 messages
	for i := 0; i < 4; i++ {
		_, err := queueClient.DeleteMessage(context.Background(), messageIDs[i], popReceipts[i], nil)
		_require.NoError(err)
	}
	// should be 0 now
	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.Equal(0, len(resp.Messages))
	_require.NoError(err)
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
	_require.NoError(err)

	resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.NoError(err)
	popReceipt := *resp.Messages[0].PopReceipt
	messageID := *resp.Messages[0].MessageID

	opts := &azqueue.DeleteMessageOptions{}
	_, err = queueClient.DeleteMessage(context.Background(), messageID, popReceipt, opts)
	_require.NoError(err)

	// should fail since we already deleted it
	_, err = queueClient.DeleteMessage(context.Background(), messageID, popReceipt, opts)
	_require.Error(err)
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// delete the queue's messages
	opts := azqueue.ClearMessagesOptions{}
	_, err = queueClient.ClearMessages(context.Background(), &opts)
	_require.NoError(err)

	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(0, len(resp.Messages))
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
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// delete the queue's messages
	_, err = queueClient.ClearMessages(context.Background(), nil)
	_require.NoError(err)

	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(0, len(resp.Messages))
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
	_require.NoError(err)

	// enqueue 33 messages
	for i := 0; i < 33; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)
	}

	// delete the queue's messages
	opts := azqueue.ClearMessagesOptions{}
	_, err = queueClient.ClearMessages(context.Background(), &opts)
	_require.NoError(err)

	resp, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(0, len(resp.Messages))
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
	_require.NoError(err)

	resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.NoError(err)
	popReceipt := *resp.Messages[0].PopReceipt
	messageID := *resp.Messages[0].MessageID

	opts := &azqueue.UpdateMessageOptions{}
	_, err = queueClient.UpdateMessage(context.Background(), messageID, popReceipt, "new content", opts)
	_require.NoError(err)

	resp1, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.NoError(err)
	content := *resp1.Messages[0].MessageText
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
	_require.NoError(err)

	resp, err := queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
	_require.NoError(err)
	popReceipt := *resp.Messages[0].PopReceipt
	messageID := *resp.Messages[0].MessageID

	opts := &azqueue.UpdateMessageOptions{VisibilityTimeout: to.Ptr(int32(1))}
	_, err = queueClient.UpdateMessage(context.Background(), messageID, popReceipt, "new content", opts)
	_require.NoError(err)
	time.Sleep(time.Second * 2)
	resp1, err := queueClient.DequeueMessage(context.Background(), nil)
	_require.NoError(err)
	content := *resp1.Messages[0].MessageText
	_require.Equal("new content", content)
}

// this test ensures that our sas related methods work properly
func (s *UnrecordedTestSuite) TestQueueSignatureValues() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)

	_require.NoError(err)
	queueName := testcommon.GenerateQueueName(testName)

	permissions := sas.QueuePermissions{
		Read:   true,
		Add:    true,
		Update: true,
	}

	expiry := time.Now().Add(time.Hour)
	qsv := sas.QueueSignatureValues{
		Version:     sas.Version,
		Protocol:    sas.ProtocolHTTPS,
		StartTime:   time.Time{},
		ExpiryTime:  expiry,
		Permissions: permissions.String(),
		QueueName:   queueName,
	}
	_, err = qsv.SignWithSharedKey(cred)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestQueueGetSASURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)
	queueClient := serviceClient.NewQueueClient(queueName)

	permissions := sas.QueuePermissions{
		Read: true,
		Add:  true,
	}
	start := time.Now().Add(-5 * time.Minute).UTC()
	expiry := time.Now().Add(time.Hour)

	opts := azqueue.GetSASURLOptions{StartTime: &start}
	sasUrl, err := queueClient.GetSASURL(permissions, expiry, &opts)
	_require.NoError(err)

	queueClient2, err := azqueue.NewQueueClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	_, err = queueClient2.Create(context.Background(), &azqueue.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.Error(err)
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.AuthorizationFailure)
}

func (s *UnrecordedTestSuite) TestQueueGetSASURL2() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)
	queueClient := serviceClient.NewQueueClient(queueName)
	start := time.Now().Add(-5 * time.Minute).UTC()
	opts := azqueue.GetSASURLOptions{StartTime: &start}

	sasUrlReadAdd, err := queueClient.GetSASURL(sas.QueuePermissions{Add: true}, time.Now().Add(time.Hour), &opts)
	_require.NoError(err)
	_, err = queueClient.Create(context.Background(), &azqueue.CreateOptions{Metadata: testcommon.BasicMetadata})
	_require.NoError(err)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	queueClient1, err := azqueue.NewQueueClientWithNoCredential(sasUrlReadAdd, nil)
	_require.NoError(err)

	// queue metadata and properties can't be read or written with SAS auth
	_, err = queueClient1.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.AuthorizationPermissionMismatch)

	start = time.Now().Add(-5 * time.Minute).UTC()
	opts = azqueue.GetSASURLOptions{StartTime: &start}

	// this should work now
	sasUrlRCWL, err := queueClient.GetSASURL(sas.QueuePermissions{Add: true,
		Read: true, Update: true, Process: true}, time.Now().Add(time.Hour), &opts)
	_require.NoError(err)

	queueClient2, err := azqueue.NewQueueClientWithNoCredential(sasUrlRCWL, nil)
	_require.NoError(err)

	// queues can't be created, deleted, or listed with service SAS auth
	_, err = queueClient2.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateQueueErrorCode(_require, err, queueerror.AuthorizationFailure)
}

func (s *UnrecordedTestSuite) TestQueueUserDelegationSAS_WithSduoid() {
	_require := require.New(s.T())
	testName := s.T().Name()

	now := time.Now().UTC().Add(-time.Minute)
	expiry := now.Add(30 * time.Minute)
	serviceCode := "q"
	version := sas.Version
	oid := "00000000-0000-0000-0000-000000000000"
	tid := "11111111-1111-1111-1111-111111111111"
	val := to.Ptr("AAAAAAAAAAAAAAAAAAAAAA==")

	udk := exported.UserDelegationKey{
		SignedStart:   &now,
		SignedExpiry:  &expiry,
		SignedService: &serviceCode,
		SignedVersion: &version,
		SignedOID:     &oid,
		SignedTID:     &tid,
		Value:         val,
	}
	udc := exported.NewUserDelegationCredential("testaccount", udk)

	queueName := testcommon.GenerateQueueName(testName)
	sv := sas.QueueSignatureValues{
		Protocol:                    sas.ProtocolHTTPS,
		StartTime:                   now,
		ExpiryTime:                  expiry,
		Permissions:                 (&sas.QueuePermissions{Read: true, Add: true}).String(),
		QueueName:                   queueName,
		SignedDelegatedUserObjectID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	}
	qp, err := sv.SignWithUserDelegation(udc)
	_require.NoError(err)
	enc := qp.Encode()
	_require.Contains(enc, "sduoid=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
}

func (s *UnrecordedTestSuite) TestServiceSASEnqueueMessage() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(
		fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)
	queueClient := serviceClient.NewQueueClient(queueName)
	_, err = queueClient.Create(context.Background(), nil)
	_require.NoError(err)

	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	permissions := sas.QueuePermissions{
		Read:   true,
		Add:    true,
		Update: true,
	}

	expiry := time.Now().Add(time.Hour)

	sasUrl, err := queueClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	queueClient1, err := azqueue.NewQueueClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient1.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)

	}
}

func (s *UnrecordedTestSuite) TestServiceSASDequeueMessage() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(
		fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)
	queueClient := serviceClient.NewQueueClient(queueName)
	_, err = queueClient.Create(context.Background(), nil)
	_require.NoError(err)

	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	permissions := sas.QueuePermissions{
		Read:    true,
		Add:     true,
		Update:  true,
		Process: true,
	}

	expiry := time.Now().Add(time.Hour)

	sasUrl, err := queueClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	queueClient1, err := azqueue.NewQueueClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient1.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)

	}
	// dequeue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClient1.DequeueMessage(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(1, len(resp.Messages))
		_require.NotNil(resp.Messages[0].MessageID)
	}
	// should be 0 now
	resp, err := queueClient1.DequeueMessage(context.Background(), nil)
	_require.Equal(0, len(resp.Messages))
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestQueueSASUsingAccessPolicy() {
	_require := require.New(s.T())

	cred, err := testcommon.GetGenericCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := testcommon.GetQueueClient(queueName, svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	_, err = queueClient.Create(context.Background(), nil)
	_require.NoError(err)

	id := "testAccessPolicy"
	ps := azqueue.AccessPolicyPermission{Read: true, Add: true, Update: true, Process: true}
	signedIdentifiers := make([]*azqueue.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &azqueue.SignedIdentifier{
		AccessPolicy: &azqueue.AccessPolicy{
			Expiry:     to.Ptr(time.Now().Add(1 * time.Hour)),
			Start:      to.Ptr(time.Now()),
			Permission: to.Ptr(ps.String()),
		},
		ID: &id,
	})

	_, err = queueClient.SetAccessPolicy(context.Background(), &azqueue.SetAccessPolicyOptions{
		QueueACL: signedIdentifiers,
	})
	_require.NoError(err)

	gResp, err := queueClient.GetAccessPolicy(context.Background(), nil)
	_require.NoError(err)
	_require.Len(gResp.SignedIdentifiers, 1)

	time.Sleep(30 * time.Second)

	sasQueryParams, err := sas.QueueSignatureValues{
		Protocol:   sas.ProtocolHTTPS,
		Identifier: id,
		QueueName:  queueName,
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	queueSAS := queueClient.URL() + "?" + sasQueryParams.Encode()
	queueClientSAS, err := azqueue.NewQueueClientWithNoCredential(queueSAS, nil)
	_require.NoError(err)

	_, err = queueClientSAS.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClientSAS.EnqueueMessage(context.Background(), fmt.Sprintf("%v : %v", testcommon.QueueDefaultData, i), nil)
		_require.NoError(err)
	}

	// dequeue 4 messages
	for i := 0; i < 4; i++ {
		resp, err := queueClientSAS.DequeueMessage(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(1, len(resp.Messages))
		_require.NotNil(resp.Messages[0].MessageText)
		_require.Equal(fmt.Sprintf("%v : %v", testcommon.QueueDefaultData, i), *resp.Messages[0].MessageText)
		_require.NotNil(resp.Messages[0].MessageID)
	}
}

func (s *UnrecordedTestSuite) TestQueueClientGetPropertiesApproximateMessagesCountInt64() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(
		fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)
	queueClient := serviceClient.NewQueueClient(queueName)

	_, err = queueClient.Create(context.Background(), nil)
	_require.NoError(err)
	defer func() {
		_, _ = queueClient.Delete(context.Background(), nil)
	}()

	const messageCount = 50
	for i := 0; i < messageCount; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), fmt.Sprintf("msg-%d", i), nil)
		_require.NoError(err)
	}
	time.Sleep(3 * time.Second)
	propsResp, err := queueClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	count := propsResp.ApproximateMessagesCount
	_require.NotNil(count)
	_require.GreaterOrEqual(*count, int64(1))
	_require.IsType(int64(0), *count)
}
