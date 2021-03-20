// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/testframework"
	chk "gopkg.in/check.v1"
)

type tableClientLiveTests struct {
	endpointType EndpointType
	mode         testframework.RecordMode
}

// Hookup to the testing framework
func Test(t *testing.T) { chk.TestingT(t) }

// wire up chk to testing
var _ = chk.Suite(&tableClientLiveTests{endpointType: StorageEndpoint, mode: testframework.Playback /* change to Record to re-record tests */})
var _ = chk.Suite(&tableClientLiveTests{endpointType: CosmosEndpoint, mode: testframework.Playback /* change to Record to re-record tests */})

func (s *tableClientLiveTests) TestServiceErrors(c *chk.C) {
	client, delete := s.init(c, true)
	defer delete()

	// Create a duplicate table to produce an error
	_, err := client.Create(ctx)
	c.Assert(err.RawResponse().StatusCode, chk.Equals, http.StatusConflict)
}

func (s *tableClientLiveTests) TestCreateTable(c *chk.C) {
	client, delete := s.init(c, false)
	defer delete()

	resp, err := client.Create(ctx)

	c.Assert(err, chk.IsNil)
	c.Assert(*resp.TableResponse.TableName, chk.Equals, client.Name())
}

func (s *tableClientLiveTests) TestAddEntity(c *chk.C) {
	client, delete := s.init(c, true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	for _, e := range *entitiesToCreate {
		_, err := client.AddEntity(ctx, &e)
		c.Assert(err, chk.IsNil)
	}
}

func (s *tableClientLiveTests) TestQuerySimpleEntity(c *chk.C) {
	client, delete := s.init(c, true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createSimpleEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		_, err := client.AddEntity(ctx, &e)
		c.Assert(err, chk.IsNil)
	}

	filter := "RowKey lt '5'"
	expectedCount := 4
	var resp TableEntityQueryResponseResponse
	pager := client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		c.Assert(len(*resp.TableEntityQueryResponse.Value), chk.Equals, expectedCount)
	}
	resp = pager.PageResponse()
	c.Assert(pager.Err(), chk.IsNil)
	for _, e := range *resp.TableEntityQueryResponse.Value {
		_, ok := e[PartitionKey].(string)
		c.Assert(ok, chk.Equals, true)
		_, ok = e[RowKey].(string)
		c.Assert(ok, chk.Equals, true)
		_, ok = e[Timestamp].(string)
		c.Assert(ok, chk.Equals, true)
		_, ok = e[EtagOdata].(string)
		c.Assert(ok, chk.Equals, true)
		_, ok = e["StringProp"].(string)
		c.Assert(ok, chk.Equals, true)
		//TODO: fix when serialization is implemented
		_, ok = e["IntProp"].(float64)
		c.Assert(ok, chk.Equals, true)
		_, ok = e["BoolProp"].(bool)
		c.Assert(ok, chk.Equals, true)
	}
}

func (s *tableClientLiveTests) SetUpTest(c *chk.C) {
	// setup the test environment
	recordedTestSetup(c, testKey(c, s.endpointType), s.endpointType, s.mode)
}

func (s *tableClientLiveTests) TearDownTest(c *chk.C) {
	// teardown the test context
	recordedTestTeardown(testKey(c, s.endpointType))
}

func (s *tableClientLiveTests) init(c *chk.C, doCreate bool) (*TableClient, func()) {
	context := getTestContext(testKey(c, s.endpointType))
	tableName := getTableName(context)
	client := context.client.GetTableClient(tableName)
	if doCreate {
		_, err := client.Create(ctx)
		c.Assert(err, chk.IsNil)
	}
	return client, func() {
		client.Delete(ctx)
	}
}
