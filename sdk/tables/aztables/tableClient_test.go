// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/testframework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type tableClientLiveTests struct {
	suite.Suite
	endpointType EndpointType
	mode         testframework.RecordMode
}

// Hookup to the testing framework
func TestTableClient_Storage(t *testing.T) {
	storage := tableClientLiveTests{endpointType: StorageEndpoint, mode: testframework.Playback /* change to Record to re-record tests */}
	suite.Run(t, &storage)
}

// Hookup to the testing framework
func TestTableClient_Cosmos(t *testing.T) {
	cosmos := tableClientLiveTests{endpointType: CosmosEndpoint, mode: testframework.Playback /* change to Record to re-record tests */}
	suite.Run(t, &cosmos)
}

func (s *tableClientLiveTests) TestServiceErrors() {
	client, delete := s.init(true)
	defer delete()

	// Create a duplicate table to produce an error
	_, err := client.Create(ctx)
	assert.Equal(s.T(), err.RawResponse().StatusCode, http.StatusConflict)
}

func (s *tableClientLiveTests) TestCreateTable() {
	assert := assert.New(s.T())
	client, delete := s.init(false)
	defer delete()

	resp, err := client.Create(ctx)

	assert.Nil(err)
	assert.Equal(*resp.TableResponse.TableName, client.Name())
}

func (s *tableClientLiveTests) TestAddEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	for _, e := range *entitiesToCreate {
		_, err := client.AddMapEntity(ctx, &e)
		assert.Nil(err)
	}
}

func (s *tableClientLiveTests) TestAddComplexEntity() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(context, 1, "partition")

	for _, e := range *entitiesToCreate {
		_, err := client.AddEntity(ctx, &e)
		assert.Nilf(err, getStringFromBody(err))
	}
}

func (s *tableClientLiveTests) TestQuerySimpleEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createSimpleEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		_, err := client.AddMapEntity(ctx, &e)
		assert.Nil(err)
	}

	filter := "RowKey lt '5'"
	expectedCount := 4
	var resp TableEntityQueryResponseResponse
	pager := client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		assert.Equal(len(*resp.TableEntityQueryResponse.Value), expectedCount)
	}
	resp = pager.PageResponse()
	assert.Nil(pager.Err())
	for _, e := range *resp.TableEntityQueryResponse.Value {
		_, ok := e[PartitionKey].(string)
		assert.True(ok)
		_, ok = e[RowKey].(string)
		assert.True(ok)
		_, ok = e[Timestamp].(string)
		assert.True(ok)
		_, ok = e[EtagOdata].(string)
		assert.True(ok)
		_, ok = e["StringProp"].(string)
		assert.True(ok)
		//TODO: fix when serialization is implemented
		_, ok = e["IntProp"].(float64)
		assert.True(ok)
		_, ok = e["BoolProp"].(bool)
		assert.True(ok)
	}
}

// setup the test environment
func (s *tableClientLiveTests) BeforeTest(suite string, test string) {
	recordedTestSetup(s.T(), s.T().Name(), s.endpointType, s.mode)
}

// teardown the test context
func (s *tableClientLiveTests) AfterTest(suite string, test string) {
	recordedTestTeardown(s.T().Name())
}

func (s *tableClientLiveTests) init(doCreate bool) (*TableClient, func()) {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, _ := getTableName(context)
	client := context.client.GetTableClient(tableName)
	if doCreate {
		_, err := client.Create(ctx)
		if err != nil {
			assert.FailNow(getStringFromBody(err))
		}
	}
	return client, func() {
		client.Delete(ctx)
	}
}

func getStringFromBody(e *runtime.ResponseError) string {
	if e == nil {
		return "Error is nil"
	}
	r := e.RawResponse()
	body := bytes.Buffer{}
	b := r.Body
	b.Close()
	if b != nil {
		_, err := body.ReadFrom(b)
		if err != nil {
			return "<emtpy body>"
		}
		b = ioutil.NopCloser(&body)
	}
	return body.String()
}
