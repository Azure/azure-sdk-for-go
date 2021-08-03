// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type tableServiceClientLiveTests struct {
	suite.Suite
	endpointType EndpointType
	mode         recording.RecordMode
}

// Hookup to the testing framework
func TestServiceClient_Storage(t *testing.T) {
	storage := tableServiceClientLiveTests{endpointType: StorageEndpoint, mode: recording.Playback /* change to Record to re-record tests */}
	suite.Run(t, &storage)
}

// Hookup to the testing framework
func TestServiceClient_Cosmos(t *testing.T) {
	cosmos := tableServiceClientLiveTests{endpointType: CosmosEndpoint, mode: recording.Playback /* change to Record to re-record tests */}
	suite.Run(t, &cosmos)
}

func (s *tableServiceClientLiveTests) TestServiceErrors() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, err := getTableName(context)
	require.NoError(err)

	_, err = context.client.CreateTable(ctx, tableName)
	delete := func() {
		_, err := context.client.DeleteTable(ctx, tableName, nil)
		if err != nil {
			fmt.Printf("Error cleaning up test. %v\n", err.Error())
		}
	}
	defer delete()
	require.NoError(err)

	// Create a duplicate table to produce an error
	_, err = context.client.CreateTable(ctx, tableName)
	var svcErr *runtime.ResponseError
	errors.As(err, &svcErr)
	require.Equal(svcErr.RawResponse().StatusCode, http.StatusConflict)
}

func (s *tableServiceClientLiveTests) TestCreateTable() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, err := getTableName(context)
	require.NoError(err)

	resp, err := context.client.CreateTable(ctx, tableName)
	delete := func() {
		_, err := context.client.DeleteTable(ctx, tableName, nil)
		if err != nil {
			fmt.Printf("Error cleaning up test. %v\n", err.Error())
		}
	}
	defer delete()

	require.NoError(err)
	require.Equal(*resp.TableResponse.TableName, tableName)
}

func (s *tableServiceClientLiveTests) TestQueryTable() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	tableCount := 5
	tableNames := make([]string, tableCount)
	prefix1 := "zzza"
	prefix2 := "zzzb"

	defer cleanupTables(context, &tableNames)
	//create 10 tables with our exected prefix and 1 with a different prefix
	for i := 0; i < tableCount; i++ {
		if i < (tableCount - 1) {
			name, _ := getTableName(context, prefix1)
			tableNames[i] = name
		} else {
			name, _ := getTableName(context, prefix2)
			tableNames[i] = name
		}
		_, err := context.client.CreateTable(ctx, tableNames[i])
		require.NoError(err)
	}

	// Query for tables with no pagination. The filter should exclude one table from the results
	filter := fmt.Sprintf("TableName ge '%s' and TableName lt '%s'", prefix1, prefix2)
	pager := context.client.ListTables(&ListOptions{Filter: &filter})

	resultCount := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		resultCount += len(resp.TableQueryResponse.Value)
	}

	require.NoError(pager.Err())
	require.Equal(resultCount, tableCount-1)

	// Query for tables with pagination
	top := int32(2)
	pager = context.client.ListTables(&ListOptions{Filter: &filter, Top: &top})

	resultCount = 0
	pageCount := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		resultCount += len(resp.TableQueryResponse.Value)
		pageCount++
	}

	require.NoError(pager.Err())
	require.Equal(resultCount, tableCount-1)
	require.Equal(pageCount, int(top))
}

func clearAllTables(context *testContext) error {
	pager := context.client.ListTables(nil)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		for _, v := range resp.TableQueryResponse.Value {
			_, err := context.client.DeleteTable(ctx, *v.TableName, nil)
			if err != nil {
				return err
			}
		}
	}
	return pager.Err()
}

func (s *tableServiceClientLiveTests) TestListTables() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, err := getTableName(context)
	require.NoError(err)

	err = clearAllTables(context)
	require.NoError(err)

	for i := 0; i < 5; i++ {
		_, err := context.client.CreateTable(ctx, fmt.Sprintf("%v%v", tableName, i))
		require.NoError(err)
	}

	count := 0
	pager := context.client.ListTables(nil)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		count += len(resp.TableQueryResponse.Value)
	}

	require.NoError(pager.Err())
	require.Equal(5, count)
}

func (s *tableServiceClientLiveTests) TestGetStatistics() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	s.T().Skip() // TODO: need to change URL to -secondary https://docs.microsoft.com/en-us/rest/api/storageservices/get-table-service-stats
	resp, err := context.client.GetStatistics(ctx, nil)
	require.NoError(err)
	require.NotNil(resp)
}

func (s *tableServiceClientLiveTests) TestGetProperties() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	resp, err := context.client.GetProperties(ctx, nil)
	require.NoError(err)
	require.NotNil(resp)
}

func (s *tableServiceClientLiveTests) TestSetLogging() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	logging := Logging{
		Read:    to.BoolPtr(true),
		Write:   to.BoolPtr(true),
		Delete:  to.BoolPtr(true),
		Version: to.StringPtr("1.0"),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(5),
		},
	}
	props := TableServiceProperties{Logging: &logging}

	resp, err := context.client.SetProperties(ctx, props, nil)
	require.NoError(err)
	require.NotNil(resp)

	// time.Sleep(45 * time.Second)

	received, err := context.client.GetProperties(ctx, nil)
	require.NoError(err)

	require.Equal(*props.Logging.Read, *received.StorageServiceProperties.Logging.Read)
	require.Equal(*props.Logging.Write, *received.StorageServiceProperties.Logging.Write)
	require.Equal(*props.Logging.Delete, *received.StorageServiceProperties.Logging.Delete)
	require.Equal(*props.Logging.RetentionPolicy.Enabled, *received.StorageServiceProperties.Logging.RetentionPolicy.Enabled)
	require.Equal(*props.Logging.RetentionPolicy.Days, *received.StorageServiceProperties.Logging.RetentionPolicy.Days)
}

func (s *tableServiceClientLiveTests) TestSetHoursMetrics() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	metrics := Metrics{
		Enabled:     to.BoolPtr(true),
		IncludeAPIs: to.BoolPtr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(5),
		},
		Version: to.StringPtr("1.0"),
	}
	props := TableServiceProperties{HourMetrics: &metrics}

	resp, err := context.client.SetProperties(ctx, props, nil)
	require.NoError(err)
	require.NotNil(resp)

	// time.Sleep(45 * time.Second)

	received, err := context.client.GetProperties(ctx, nil)
	require.NoError(err)

	require.Equal(*props.HourMetrics.Enabled, *received.StorageServiceProperties.HourMetrics.Enabled)
	require.Equal(*props.HourMetrics.IncludeAPIs, *received.StorageServiceProperties.HourMetrics.IncludeAPIs)
	require.Equal(*props.HourMetrics.RetentionPolicy.Days, *received.StorageServiceProperties.HourMetrics.RetentionPolicy.Days)
	require.Equal(*props.HourMetrics.RetentionPolicy.Enabled, *received.StorageServiceProperties.HourMetrics.RetentionPolicy.Enabled)
}

func (s *tableServiceClientLiveTests) TestSetMinuteMetrics() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	metrics := Metrics{
		Enabled:     to.BoolPtr(true),
		IncludeAPIs: to.BoolPtr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(5),
		},
		Version: to.StringPtr("1.0"),
	}
	props := TableServiceProperties{MinuteMetrics: &metrics}

	resp, err := context.client.SetProperties(ctx, props, nil)
	require.NoError(err)
	require.NotNil(resp)

	// time.Sleep(45 * time.Second)

	received, err := context.client.GetProperties(ctx, nil)
	require.NoError(err)

	require.Equal(*props.MinuteMetrics.Enabled, *received.StorageServiceProperties.MinuteMetrics.Enabled)
	require.Equal(*props.MinuteMetrics.IncludeAPIs, *received.StorageServiceProperties.MinuteMetrics.IncludeAPIs)
	require.Equal(*props.MinuteMetrics.RetentionPolicy.Days, *received.StorageServiceProperties.MinuteMetrics.RetentionPolicy.Days)
	require.Equal(*props.MinuteMetrics.RetentionPolicy.Enabled, *received.StorageServiceProperties.MinuteMetrics.RetentionPolicy.Enabled)
}

func (s *tableServiceClientLiveTests) TestSetCors() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	corsRules1 := CorsRule{
		AllowedHeaders:  to.StringPtr("x-ms-meta-data*"),
		AllowedMethods:  to.StringPtr("PUT"),
		AllowedOrigins:  to.StringPtr("www.xyz.com"),
		ExposedHeaders:  to.StringPtr("x-ms-meta-source*"),
		MaxAgeInSeconds: to.Int32Ptr(500),
	}
	props := TableServiceProperties{Cors: []*CorsRule{&corsRules1}}

	resp, err := context.client.SetProperties(ctx, props, nil)
	require.NoError(err)
	require.NotNil(resp)

	// time.Sleep(45 * time.Second)

	received, err := context.client.GetProperties(ctx, nil)
	require.NoError(err)

	require.Equal(*props.Cors[0].AllowedHeaders, *received.StorageServiceProperties.Cors[0].AllowedHeaders)
	require.Equal(*props.Cors[0].AllowedMethods, *received.StorageServiceProperties.Cors[0].AllowedMethods)
	require.Equal(*props.Cors[0].AllowedOrigins, *received.StorageServiceProperties.Cors[0].AllowedOrigins)
	require.Equal(*props.Cors[0].ExposedHeaders, *received.StorageServiceProperties.Cors[0].ExposedHeaders)
	require.Equal(*props.Cors[0].MaxAgeInSeconds, *received.StorageServiceProperties.Cors[0].MaxAgeInSeconds)
}

func (s *tableServiceClientLiveTests) TestSetTooManyCors() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	corsRules1 := CorsRule{
		AllowedHeaders:  to.StringPtr("x-ms-meta-data*"),
		AllowedMethods:  to.StringPtr("PUT"),
		AllowedOrigins:  to.StringPtr("www.xyz.com"),
		ExposedHeaders:  to.StringPtr("x-ms-meta-source*"),
		MaxAgeInSeconds: to.Int32Ptr(500),
	}
	props := TableServiceProperties{Cors: make([]*CorsRule, 0)}
	for i := 0; i < 6; i++ {
		props.Cors = append(props.Cors, &corsRules1)
	}

	_, err := context.client.SetProperties(ctx, props, nil)
	require.Error(err)
}

func (s *tableServiceClientLiveTests) TestRetentionTooLong() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	if _, ok := cosmosTestsMap[s.T().Name()]; ok {
		s.T().Skip()
	}

	metrics := Metrics{
		Enabled:     to.BoolPtr(true),
		IncludeAPIs: to.BoolPtr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(366),
		},
		Version: to.StringPtr("1.0"),
	}
	props := TableServiceProperties{MinuteMetrics: &metrics}

	_, err := context.client.SetProperties(ctx, props, nil)
	require.Error(err)
}

func (s *tableServiceClientLiveTests) BeforeTest(suite string, test string) {
	// setup the test environment
	recordedTestSetup(s.T(), s.T().Name(), s.endpointType, s.mode)
}

func (s *tableServiceClientLiveTests) AfterTest(suite string, test string) {
	// teardown the test context
	recordedTestTeardown(s.T().Name())
}
