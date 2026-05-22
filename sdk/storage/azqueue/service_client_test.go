// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/queueerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/sas"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running service Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &RecordedTestSuite{})
		suite.Run(t, &UnrecordedTestSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &RecordedTestSuite{})
	case recording.RecordingMode:
		suite.Run(t, &RecordedTestSuite{})
	}
}

func (s *RecordedTestSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *RecordedTestSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
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
	proxy *recording.TestProxyInstance
}

type UnrecordedTestSuite struct {
	suite.Suite
}

func (s *UnrecordedTestSuite) TestServiceClientFromConnectionString() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetAccountInfo(testcommon.TestAccountDefault)
	connectionString := testcommon.GetConnectionString(testcommon.TestAccountDefault)

	parsedConnStr, err := shared.ParseConnectionString(connectionString)
	_require.NoError(err)
	_require.Equal(parsedConnStr.ServiceURL, "https://"+accountName+".queue.core.windows.net/")

	sharedKeyCred, err := azqueue.NewSharedKeyCredential(parsedConnStr.AccountName, parsedConnStr.AccountKey)
	_require.NoError(err)

	svcClient, err := azqueue.NewServiceClientWithSharedKeyCredential(parsedConnStr.ServiceURL, sharedKeyCred, nil)
	_require.NoError(err)

	sProps, err := svcClient.GetServiceProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotZero(sProps)
}

func (s *UnrecordedTestSuite) TestServiceClientFromConnectionString1() {
	_require := require.New(s.T())

	connectionString := testcommon.GetConnectionString(testcommon.TestAccountDefault)

	svcClient, err := azqueue.NewServiceClientFromConnectionString(connectionString, nil)
	_require.NoError(err)

	sProps, err := svcClient.GetServiceProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotZero(sProps)
}

func (s *RecordedTestSuite) TestSetPropertiesLogging() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)

	loggingOpts := azqueue.Logging{
		Read: enabled, Write: enabled, Delete: enabled,
		RetentionPolicy: &azqueue.RetentionPolicy{Enabled: enabled, Days: days}}
	opts := azqueue.SetPropertiesOptions{Logging: &loggingOpts}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.NoError(err)
	resp1, err := svcClient.GetServiceProperties(context.Background(), nil)

	_require.NoError(err)
	_require.Equal(resp1.Logging.Write, enabled)
	_require.Equal(resp1.Logging.Read, enabled)
	_require.Equal(resp1.Logging.Delete, enabled)
	_require.Equal(resp1.Logging.RetentionPolicy.Days, days)
	_require.Equal(resp1.Logging.RetentionPolicy.Enabled, enabled)
}

func (s *RecordedTestSuite) TestSetPropertiesHourMetrics() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)

	metricsOpts := azqueue.Metrics{
		Enabled: enabled, IncludeAPIs: enabled, RetentionPolicy: &azqueue.RetentionPolicy{Enabled: enabled, Days: days}}
	opts := azqueue.SetPropertiesOptions{HourMetrics: &metricsOpts}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.NoError(err)
	resp1, err := svcClient.GetServiceProperties(context.Background(), nil)

	_require.NoError(err)
	_require.Equal(resp1.HourMetrics.Enabled, enabled)
	_require.Equal(resp1.HourMetrics.IncludeAPIs, enabled)
	_require.Equal(resp1.HourMetrics.RetentionPolicy.Days, days)
	_require.Equal(resp1.HourMetrics.RetentionPolicy.Enabled, enabled)
}

func (s *RecordedTestSuite) TestSetPropertiesMinuteMetrics() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	days := to.Ptr[int32](5)
	enabled := to.Ptr(true)

	metricsOpts := azqueue.Metrics{
		Enabled: enabled, IncludeAPIs: enabled, RetentionPolicy: &azqueue.RetentionPolicy{Enabled: enabled, Days: days}}
	opts := azqueue.SetPropertiesOptions{MinuteMetrics: &metricsOpts}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.NoError(err)
	resp1, err := svcClient.GetServiceProperties(context.Background(), nil)

	_require.NoError(err)
	_require.Equal(resp1.MinuteMetrics.Enabled, enabled)
	_require.Equal(resp1.MinuteMetrics.IncludeAPIs, enabled)
	_require.Equal(resp1.MinuteMetrics.RetentionPolicy.Days, days)
	_require.Equal(resp1.MinuteMetrics.RetentionPolicy.Enabled, enabled)
}

func (s *RecordedTestSuite) TestSetPropertiesSetQueueCORS() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	defaultAge := to.Ptr[int32](500)
	defaultStr := to.Ptr[string]("")

	allowedOrigins1 := "www.xyz.com"
	allowedMethods1 := "GET"
	CORSOpts1 := &azqueue.CORSRule{AllowedOrigins: &allowedOrigins1, AllowedMethods: &allowedMethods1}

	allowedOrigins2 := "www.xyz.com,www.ab.com,www.bc.com"
	allowedMethods2 := "GET, PUT"
	maxAge2 := to.Ptr[int32](500)
	exposedHeaders2 := "x-ms-meta-data*,x-ms-meta-source*,x-ms-meta-abc,x-ms-meta-bcd"
	allowedHeaders2 := "x-ms-meta-data*,x-ms-meta-target*,x-ms-meta-xyz,x-ms-meta-foo"

	CORSOpts2 := &azqueue.CORSRule{
		AllowedOrigins: &allowedOrigins2, AllowedMethods: &allowedMethods2,
		MaxAgeInSeconds: maxAge2, ExposedHeaders: &exposedHeaders2, AllowedHeaders: &allowedHeaders2}

	CORSRules := []*azqueue.CORSRule{CORSOpts1, CORSOpts2}

	opts := azqueue.SetPropertiesOptions{CORS: CORSRules}
	_, err = svcClient.SetProperties(context.Background(), &opts)

	_require.NoError(err)
	resp, err := svcClient.GetServiceProperties(context.Background(), nil)
	for i := 0; i < len(resp.CORS); i++ {
		if resp.CORS[i].AllowedOrigins == &allowedOrigins1 {
			_require.Equal(resp.CORS[i].AllowedMethods, &allowedMethods1)
			_require.Equal(resp.CORS[i].MaxAgeInSeconds, defaultAge)
			_require.Equal(resp.CORS[i].ExposedHeaders, defaultStr)
			_require.Equal(resp.CORS[i].AllowedHeaders, defaultStr)

		} else if resp.CORS[i].AllowedOrigins == &allowedOrigins2 {
			_require.Equal(resp.CORS[i].AllowedMethods, &allowedMethods2)
			_require.Equal(resp.CORS[i].MaxAgeInSeconds, &maxAge2)
			_require.Equal(resp.CORS[i].ExposedHeaders, &exposedHeaders2)
			_require.Equal(resp.CORS[i].AllowedHeaders, &allowedHeaders2)
		}
	}
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestServiceGetProperties() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sProps, err := svcClient.GetServiceProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotZero(sProps)
}

func (s *RecordedTestSuite) TestServiceCreateQueue() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := svcClient.NewQueueClient(queueName)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)

	resp, err := svcClient.CreateQueue(context.Background(), queueName, nil)
	_require.NoError(err)
	_require.NotZero(resp)
}

func (s *RecordedTestSuite) TestServiceCreateQueueWithMetadata() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	queueClient := svcClient.NewQueueClient(queueName)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient)
	opts := azqueue.CreateOptions{Metadata: testcommon.BasicMetadata}

	resp, err := svcClient.CreateQueue(context.Background(), queueName, &opts)
	_require.NoError(err)
	_require.NotZero(resp)
}

func (s *RecordedTestSuite) TestServiceDeleteQueue() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	testName := s.T().Name()
	queueName := testcommon.GenerateQueueName(testName)
	createResp, err := svcClient.CreateQueue(context.Background(), queueName, nil)
	_require.NoError(err)
	_require.NotZero(createResp)

	delResp, err := svcClient.DeleteQueue(context.Background(), queueName, nil)
	_require.NoError(err)
	_require.NotZero(delResp)
}

func (s *RecordedTestSuite) TestServiceListQueuesWithMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
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
			_require.NoError(err)
		}
	}(queueClient, context.Background(), nil)
	_require.NoError(err)
	listOptions := azqueue.ListQueuesOptions{Include: azqueue.ListQueuesInclude{Metadata: true}}
	pager := svcClient.NewListQueuesPager(&listOptions)

	exists := false
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, queue := range resp.Queues {
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

	_require.NoError(err)
	_require.True(exists)
}

func (s *RecordedTestSuite) TestServiceListQueuesPagination() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)
	queueNames := []string{queueName + "1", queueName + "2", queueName + "3"}

	for i := 0; i < len(queueNames); i++ {
		_, err := svcClient.CreateQueue(context.Background(), queueNames[i], nil)
		if err != nil {
			_require.NoError(err)
		}
	}
	// cleanup created queues
	defer func(queueNames []string, ctx context.Context, options *azqueue.DeleteOptions) {
		for i := 0; i < len(queueNames); i++ {
			_, err := svcClient.DeleteQueue(ctx, queueNames[i], nil)
			if err != nil {
				_require.NoError(err)
			}
		}
	}(queueNames, context.Background(), nil)
	_require.NoError(err)
	pager := svcClient.NewListQueuesPager(&azqueue.ListQueuesOptions{MaxResults: to.Ptr(int32(1))})

	count := 0
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, queue := range resp.Queues {
			_require.NotNil(queue.Name)
			count += 1
		}
		if err != nil {
			break
		}
	}
	// greater or equal since storage account might have more than the 3 queues created above
	_require.GreaterOrEqual(count, 3)
}

func (s *RecordedTestSuite) TestServiceListQueuesPaginationEmptyPrefix() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	queueClient1 := testcommon.CreateNewQueue(context.Background(), _require, testcommon.GenerateQueueName(testName)+"1", svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient1)
	queueClient2 := testcommon.CreateNewQueue(context.Background(), _require, testcommon.GenerateQueueName(testName)+"2", svcClient)
	defer testcommon.DeleteQueue(context.Background(), _require, queueClient2)

	count := 0
	pager := svcClient.NewListQueuesPager(nil)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)

		for _, queue := range resp.Queues {
			count++
			_require.NotNil(queue.Name)
		}
		if err != nil {
			break
		}
	}
	_require.GreaterOrEqual(count, 2)
}

func (s *RecordedTestSuite) TestServiceListQueuesPaged() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	const numQueues = 6
	maxResults := int32(2)
	const pagedqueuesPrefix = "azqueuepaged"

	queues := make([]*azqueue.QueueClient, numQueues)
	expectedResults := make(map[string]bool)
	for i := 0; i < numQueues; i++ {
		queueName := pagedqueuesPrefix + testcommon.GenerateQueueName(testName) + fmt.Sprintf("%d", i)
		queueClient := testcommon.CreateNewQueue(context.Background(), _require, queueName, svcClient)
		queues[i] = queueClient
		expectedResults[queueName] = false
	}

	defer func() {
		for i := range queues {
			testcommon.DeleteQueue(context.Background(), _require, queues[i])
		}
	}()

	prefix := pagedqueuesPrefix + testcommon.QueuePrefix
	listOptions := azqueue.ListQueuesOptions{MaxResults: &maxResults, Prefix: &prefix, Include: azqueue.ListQueuesInclude{Metadata: true}}
	count := 0
	results := make([]azqueue.Queue, 0)
	pager := svcClient.NewListQueuesPager(&listOptions)

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, ctnr := range resp.Queues {
			_require.NotNil(ctnr.Name)
			results = append(results, *ctnr)
			count += 1
		}
	}

	_require.Equal(count, numQueues)
	_require.Equal(len(results), numQueues)

	// make sure each queue we see is expected
	for _, q := range results {
		_, ok := expectedResults[*q.Name]
		_require.Equal(ok, true)
		expectedResults[*q.Name] = true
	}

	// make sure every expected queue was seen
	for _, seen := range expectedResults {
		_require.Equal(seen, true)
	}
}

// this test ensures that our sas related methods work properly
func (s *UnrecordedTestSuite) TestServiceSignatureValues() {
	_require := require.New(s.T())

	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Update: true,
		Delete: true,
	}

	expiry := time.Now().Add(time.Hour)

	qsv := sas.AccountSignatureValues{
		Version:       sas.Version,
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Time{},
		ExpiryTime:    expiry,
		ResourceTypes: resources.String(),
		Permissions:   permissions.String(),
	}
	_, err = qsv.SignWithSharedKey(cred)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestSASServiceClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)

	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Update: true,
		Delete: true,
	}

	expiry := time.Now().Add(time.Hour)

	sasUrl, err := serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.NoError(err)

	svcClient, err := azqueue.NewServiceClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	// create queue using SAS
	_, err = svcClient.CreateQueue(context.Background(), queueName, nil)
	_require.NoError(err)

	_, err = svcClient.DeleteQueue(context.Background(), queueName, nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestNoSharedKeyCredError() {
	_require := require.New(s.T())
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")

	// Creating service client without credentials
	serviceClient, err := azqueue.NewServiceClientWithNoCredential(fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), nil)
	_require.NoError(err)

	// Adding SAS and options
	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Update: true,
		Delete: true,
	}
	start := time.Now().Add(-time.Hour)
	expiry := start.Add(time.Hour)
	opts := azqueue.GetSASURLOptions{StartTime: &start}

	// GetSASURL fails (with MissingSharedKeyCredential) because service client is created without credentials
	_, err = serviceClient.GetSASURL(resources, permissions, expiry, &opts)
	_require.Equal(err, queueerror.MissingSharedKeyCredential)

}

func (s *UnrecordedTestSuite) TestAccountSASEnqueueMessage() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)

	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Update: true,
		Delete: true,
	}

	expiry := time.Now().Add(time.Hour)

	sasUrl, err := serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.NoError(err)

	svcClient, err := azqueue.NewServiceClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	// create queue using account SAS
	_, err = svcClient.CreateQueue(context.Background(), queueName, nil)
	_require.NoError(err)

	queueClient := svcClient.NewQueueClient(queueName)
	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		_, err = queueClient.EnqueueMessage(context.Background(), testcommon.QueueDefaultData, nil)
		_require.NoError(err)

	}
	_, err = svcClient.DeleteQueue(context.Background(), queueName, nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestAccountSASDequeueMessage() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(fmt.Sprintf("https://%s.queue.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	queueName := testcommon.GenerateQueueName(testName)

	resources := sas.AccountResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := sas.AccountPermissions{
		Read:    true,
		Add:     true,
		Write:   true,
		Create:  true,
		Update:  true,
		Delete:  true,
		Process: true,
	}

	expiry := time.Now().Add(time.Hour)

	sasUrl, err := serviceClient.GetSASURL(resources, permissions, expiry, nil)
	_require.NoError(err)

	svcClient, err := azqueue.NewServiceClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	// create queue using account SAS
	_, err = svcClient.CreateQueue(context.Background(), queueName, nil)
	_require.NoError(err)

	queueClient := svcClient.NewQueueClient(queueName)
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
	_, err = svcClient.DeleteQueue(context.Background(), queueName, nil)
	_require.NoError(err)
}
