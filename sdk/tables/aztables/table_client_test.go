// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/testframework"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func (s *tableClientLiveTests) TestQueryComplexEntity() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createComplexMapEntities(context, 5, "partition")
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
		_, ok = e["SomeBinaryProperty"].([]byte)
		assert.True(ok)
		_, ok = e["SomeDateProperty"].(time.Time)
		assert.True(ok)
		_, ok = e["SomeDoubleProperty0"].(float64)
		assert.True(ok)
		_, ok = e["SomeDoubleProperty1"].(float64)
		assert.True(ok)
		_, ok = e["SomeGuidProperty"].(uuid.UUID)
		assert.True(ok)
		_, ok = e["SomeInt64Property"].(int64)
		assert.True(ok)
		//TODO: fix when serialization is implemented
		_, ok = e["SomeIntProperty"].(float64)
		assert.True(ok)
		_, ok = e["SomeStringProperty"].(string)
		assert.True(ok)
	}
}

func (s *tableClientLiveTests) TestBatchAdd() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexMapEntities(context, 10, "partition")
	batch := make([]TableTransactionAction, 10)

	for i, e := range *entitiesToCreate {
		batch[i] = TableTransactionAction{ActionType: Add, Entity: e}
	}

	resp, err := client.submitTransactionInternal(&batch, context.recording.UUID(), context.recording.UUID(), nil, ctx)
	assert.Nil(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(r.StatusCode, http.StatusNoContent)
	}
}

func (s *tableClientLiveTests) TestBatchMixed() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexMapEntities(context, 5, "partition")
	batch := make([]TableTransactionAction, 3)

	// Add the first 3 entities.
	for i, _ := range batch {
		batch[i] = TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[i]}
	}

	resp, err := client.submitTransactionInternal(&batch, context.recording.UUID(), context.recording.UUID(), nil, ctx)
	require.Nil(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(http.StatusNoContent, r.StatusCode)
	}

	// create a new batch slice.
	batch = make([]TableTransactionAction, 5)

	// create a merge action for the first added entity
	mergeProp := "MergeProperty"
	val := "foo"
	var mergeProperty = map[string]interface{}{
		PartitionKey: (*entitiesToCreate)[0][PartitionKey],
		RowKey:       (*entitiesToCreate)[0][RowKey],
		mergeProp:    val,
	}
	batch[0] = TableTransactionAction{ActionType: UpdateMerge, Entity: mergeProperty, ETag: (*resp.TransactionResponses)[0].Header.Get(ETag)}

	// create a delete action for the second added entity
	batch[1] = TableTransactionAction{ActionType: Delete, Entity: (*entitiesToCreate)[1]}

	// create an upsert action to replace the third added entity with a new value
	replaceProp := "ReplaceProperty"
	var replaceProperties = map[string]interface{}{
		PartitionKey: (*entitiesToCreate)[2][PartitionKey],
		RowKey:       (*entitiesToCreate)[2][RowKey],
		replaceProp:  val,
	}
	batch[2] = TableTransactionAction{ActionType: UpsertReplace, Entity: replaceProperties}

	// Add the remaining 2 entities.
	batch[3] = TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[3]}
	batch[4] = TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[4]}

	//batch = batch[1:]

	resp, err = client.submitTransactionInternal(&batch, context.recording.UUID(), context.recording.UUID(), nil, ctx)
	require.Nil(err)

	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(http.StatusNoContent, r.StatusCode)
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
