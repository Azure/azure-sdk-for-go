// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/testframework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type tableServiceClientLiveTests struct {
	suite.Suite
	endpointType EndpointType
	mode         testframework.RecordMode
}

// Hookup to the testing framework
func TestServiceClient_Storage(t *testing.T) {
	storage := tableServiceClientLiveTests{endpointType: StorageEndpoint, mode: testframework.Playback /* change to Record to re-record tests */}
	suite.Run(t, &storage)
}

// Hookup to the testing framework
func TestServiceClient_Cosmos(t *testing.T) {
	cosmos := tableServiceClientLiveTests{endpointType: CosmosEndpoint, mode: testframework.Playback /* change to Record to re-record tests */}
	suite.Run(t, &cosmos)
}

func (s *tableServiceClientLiveTests) TestServiceErrors() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, err := getTableName(context)

	_, err = context.client.Create(ctx, tableName)
	defer context.client.Delete(ctx, tableName)
	assert.Nil(err)

	// Create a duplicate table to produce an error
	_, svcErr := context.client.Create(ctx, tableName)
	assert.Equal(svcErr.RawResponse().StatusCode, http.StatusConflict)
}

func (s *tableServiceClientLiveTests) TestCreateTable() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, err := getTableName(context)

	resp, err := context.client.Create(ctx, tableName)
	defer context.client.Delete(ctx, tableName)

	assert.Nil(err)
	assert.Equal(*resp.TableResponse.TableName, tableName)
}

func (s *tableServiceClientLiveTests) TestQueryTable() {
	assert := assert.New(s.T())
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
		_, err := context.client.Create(ctx, tableNames[i])
		assert.Nil(err)
	}

	// Query for tables with no pagination. The filter should exclude one table from the results
	filter := fmt.Sprintf("TableName ge '%s' and TableName lt '%s'", prefix1, prefix2)
	pager := context.client.QueryTables(QueryOptions{Filter: &filter})

	resultCount := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		resultCount += len(*resp.TableQueryResponse.Value)
	}

	assert.Nil(pager.Err())
	assert.Equal(resultCount, tableCount-1)

	// Query for tables with pagination
	top := int32(2)
	pager = context.client.QueryTables(QueryOptions{Filter: &filter, Top: &top})

	resultCount = 0
	pageCount := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		resultCount += len(*resp.TableQueryResponse.Value)
		pageCount++
	}

	assert.Nil(pager.Err())
	assert.Equal(resultCount, tableCount-1)
	assert.Equal(pageCount, int(top))
}

func (s *tableServiceClientLiveTests) BeforeTest(suite string, test string) {
	// setup the test environment
	recordedTestSetup(s.T(), s.T().Name(), s.endpointType, s.mode)
}

func (s *tableServiceClientLiveTests) AfterTest(suite string, test string) {
	// teardown the test context
	recordedTestTeardown(s.T().Name())
}
