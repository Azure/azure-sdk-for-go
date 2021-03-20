// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/testframework"
	chk "gopkg.in/check.v1"
)

type tableServiceClientLiveTests struct {
	endpointType EndpointType
	mode         testframework.RecordMode
}

// Hookup to the testing framework
var _ = chk.Suite(&tableServiceClientLiveTests{endpointType: StorageEndpoint, mode: testframework.Playback /* change to Record to re-record tests */})
var _ = chk.Suite(&tableServiceClientLiveTests{endpointType: CosmosEndpoint, mode: testframework.Playback /* change to Record to re-record tests */})

func (s *tableServiceClientLiveTests) TestServiceErrors(c *chk.C) {
	context := getTestContext(testKey(c, s.endpointType))
	tableName := getTableName(context)

	_, err := context.client.Create(ctx, tableName)
	defer context.client.Delete(ctx, tableName)
	c.Assert(err, chk.IsNil)

	// Create a duplicate table to produce an error
	_, err = context.client.Create(ctx, tableName)
	c.Assert(err.RawResponse().StatusCode, chk.Equals, http.StatusConflict)
}

func (s *tableServiceClientLiveTests) TestCreateTable(c *chk.C) {
	context := getTestContext(testKey(c, s.endpointType))
	tableName := getTableName(context)

	resp, err := context.client.Create(ctx, tableName)
	defer context.client.Delete(ctx, tableName)

	c.Assert(err, chk.IsNil)
	c.Assert(*resp.TableResponse.TableName, chk.Equals, tableName)
}

func (s *tableServiceClientLiveTests) TestQueryTable(c *chk.C) {
	context := getTestContext(testKey(c, s.endpointType))
	tableCount := 5
	tableNames := make([]string, tableCount)
	prefix1 := "zzza"
	prefix2 := "zzzb"

	defer cleanupTables(context, &tableNames)
	//create 10 tables with our exected prefix and 1 with a different prefix
	for i := 0; i < tableCount; i++ {
		if i < (tableCount - 1) {
			tableNames[i] = getTableName(context, prefix1)
		} else {
			tableNames[i] = getTableName(context, prefix2)
		}
		_, err := context.client.Create(ctx, tableNames[i])
		c.Assert(err, chk.IsNil)
	}

	// Query for tables with no pagination. The filter should exclude one table from the results
	filter := fmt.Sprintf("TableName ge '%s' and TableName lt '%s'", prefix1, prefix2)
	pager := context.client.QueryTables(QueryOptions{Filter: &filter})

	resultCount := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		resultCount += len(*resp.TableQueryResponse.Value)
	}

	c.Assert(pager.Err(), chk.IsNil)
	c.Assert(resultCount, chk.Equals, tableCount-1)

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

	c.Assert(pager.Err(), chk.IsNil)
	c.Assert(resultCount, chk.Equals, tableCount-1)
	c.Assert(pageCount, chk.Equals, int(top))
}

func (s *tableServiceClientLiveTests) SetUpTest(c *chk.C) {
	// setup the test environment
	recordedTestSetup(c, testKey(c, s.endpointType), s.endpointType, s.mode)
}

func (s *tableServiceClientLiveTests) TearDownTest(c *chk.C) {
	// teardown the test context
	recordedTestTeardown(testKey(c, s.endpointType))
}
