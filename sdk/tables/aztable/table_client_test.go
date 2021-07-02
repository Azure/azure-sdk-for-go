// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type tableClientLiveTests struct {
	suite.Suite
	endpointType EndpointType
	mode         recording.RecordMode
}

// Hookup to the testing framework
func TestTableClient_Storage(t *testing.T) {
	storage := tableClientLiveTests{endpointType: StorageEndpoint, mode: recording.Playback /* change to Record to re-record tests */}
	suite.Run(t, &storage)
}

// Hookup to the testing framework
func TestTableClient_Cosmos(t *testing.T) {
	cosmos := tableClientLiveTests{endpointType: CosmosEndpoint, mode: recording.Playback /* change to Record to re-record tests */}
	suite.Run(t, &cosmos)
}

func (s *tableClientLiveTests) TestServiceErrors() {
	client, delete := s.init(true)
	defer delete()

	// Create a duplicate table to produce an error
	_, err := client.Create(ctx)
	var svcErr *runtime.ResponseError
	errors.As(err, &svcErr)
	assert.Equal(s.T(), svcErr.RawResponse().StatusCode, http.StatusConflict)
}

func (s *tableClientLiveTests) TestCreateTable() {
	assert := assert.New(s.T())
	client, delete := s.init(false)
	defer delete()

	resp, err := client.Create(ctx)

	assert.Nil(err)
	assert.Equal(*resp.TableResponse.TableName, client.Name)
}

func (s *tableClientLiveTests) TestAddEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	_, err := client.AddEntity(ctx, (*entitiesToCreate)[0])
	assert.Nil(err)
}

func (s *tableClientLiveTests) TestAddComplexEntity() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(context, 1, "partition")

	for _, e := range *entitiesToCreate {
		_, err := client.AddEntity(ctx, e)
		var svcErr *runtime.ResponseError
		errors.As(err, &svcErr)
		assert.Nilf(err, getStringFromBody(svcErr))
	}
}

func (s *tableClientLiveTests) TestDeleteEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	_, err := client.AddEntity(ctx, (*entitiesToCreate)[0])
	assert.Nil(err)
	_, delErr := client.DeleteEntity(ctx, (*entitiesToCreate)[0][partitionKey].(string), (*entitiesToCreate)[0][rowKey].(string), "*")
	assert.Nil(delErr)
}

func (s *tableClientLiveTests) TestMergeEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	_, err := client.AddEntity(ctx, (*entitiesToCreate)[0])
	assert.Nil(err)

	var qResp TableEntityQueryResponseResponse
	filter := "RowKey eq '1'"
	pager := client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	preMerge := qResp.TableEntityQueryResponse.Value[0]

	mergeProp := "MergeProperty"
	val := "foo"
	var mergeProperty = map[string]interface{}{
		partitionKey: (*entitiesToCreate)[0][partitionKey],
		rowKey:       (*entitiesToCreate)[0][rowKey],
		mergeProp:    val,
	}

	_, updateErr := client.UpdateEntity(ctx, mergeProperty, nil, Merge)
	assert.Nil(updateErr)

	pager = client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]

	// The merged entity has all its properties + the merged property
	assert.Equalf(len(preMerge)+1, len(postMerge), "postMerge should have one more property than preMerge")
	assert.Equalf(postMerge[mergeProp], val, "%s property should equal %s", mergeProp, val)
}

func (s *tableClientLiveTests) TestUpsertEntity() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createSimpleEntities(1, "partition")

	_, err := client.UpsertEntity(ctx, (*entitiesToCreate)[0], Replace)
	require.Nil(err)

	var qResp TableEntityQueryResponseResponse
	filter := "RowKey eq '1'"
	pager := client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	preMerge := qResp.TableEntityQueryResponse.Value[0]

	mergeProp := "MergeProperty"
	val := "foo"
	var mergeProperty = map[string]interface{}{
		partitionKey: (*entitiesToCreate)[0][partitionKey],
		rowKey:       (*entitiesToCreate)[0][rowKey],
		mergeProp:    val,
	}

	_, updateErr := client.UpsertEntity(ctx, mergeProperty, Replace)
	require.Nil(updateErr)

	pager = client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]

	// The merged entity has only the standard properties + the merged property
	assert.Greater(len(preMerge), len(postMerge), "postMerge should have fewer properties than preMerge")
	assert.Equalf(postMerge[mergeProp], val, "%s property should equal %s", mergeProp, val)
}

func (s *tableClientLiveTests) _TestGetEntity() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createSimpleEntities(1, "partition")
	for _, e := range *entitiesToCreate {
		_, err := client.AddEntity(ctx, e)
		assert.Nil(err)
	}

	resp, err := client.GetEntity(ctx, "partition", "1")
	require.Nil(err)
	e := resp.Value
	_, ok := e[partitionKey].(string)
	assert.True(ok)
	_, ok = e[rowKey].(string)
	assert.True(ok)
	_, ok = e[timestamp].(string)
	assert.True(ok)
	_, ok = e[etagOdata].(string)
	assert.True(ok)
	_, ok = e["StringProp"].(string)
	assert.True(ok)
	//TODO: fix when serialization is implemented
	_, ok = e["IntProp"].(float64)
	assert.True(ok)
	_, ok = e["BoolProp"].(bool)
	assert.True(ok)
}

func (s *tableClientLiveTests) TestQuerySimpleEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createSimpleEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		_, err := client.AddEntity(ctx, e)
		assert.Nil(err)
	}

	filter := "RowKey lt '5'"
	expectedCount := 4
	var resp TableEntityQueryResponseResponse
	var models []simpleEntity
	pager := client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		models = make([]simpleEntity, len(resp.TableEntityQueryResponse.Value))
		resp.TableEntityQueryResponse.AsModels(&models)
		assert.Equal(len(resp.TableEntityQueryResponse.Value), expectedCount)
	}
	resp = pager.PageResponse()
	assert.Nil(pager.Err())
	for i, e := range resp.TableEntityQueryResponse.Value {
		_, ok := e[partitionKey].(string)
		assert.True(ok)
		assert.Equal(e[partitionKey], models[i].PartitionKey)
		_, ok = e[rowKey].(string)
		assert.True(ok)
		assert.Equal(e[rowKey], models[i].RowKey)
		_, ok = e[timestamp].(string)
		assert.True(ok)
		_, ok = e[etagOdata].(string)
		assert.True(ok)
		assert.Equal(e[etagOdata], models[i].ETag)
		_, ok = e["StringProp"].(string)
		assert.True(ok)
		//TODO: fix when serialization is implemented
		_, ok = e["IntProp"].(float64)
		assert.Equal(int(e["IntProp"].(float64)), models[i].IntProp)
		assert.True(ok)
		_, ok = e["BoolProp"].(bool)
		assert.Equal((*entitiesToCreate)[i]["BoolProp"], e["BoolProp"])
		assert.Equal(e["BoolProp"], models[i].BoolProp)
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
		_, err := client.AddEntity(ctx, e)
		assert.Nil(err)
	}

	filter := "RowKey lt '5'"
	expectedCount := 4
	var resp TableEntityQueryResponseResponse
	pager := client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		assert.Equal(expectedCount, len(resp.TableEntityQueryResponse.Value))
	}
	resp = pager.PageResponse()
	assert.Nil(pager.Err())
	for _, e := range resp.TableEntityQueryResponse.Value {
		_, ok := e[partitionKey].(string)
		assert.True(ok)
		_, ok = e[rowKey].(string)
		assert.True(ok)
		_, ok = e[timestamp].(string)
		assert.True(ok)
		_, ok = e[etagOdata].(string)
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

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
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
	for i := range batch {
		batch[i] = TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[i]}
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.Nil(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(http.StatusNoContent, r.StatusCode)
	}

	var qResp TableEntityQueryResponseResponse
	filter := "RowKey eq '1'"
	pager := client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	preMerge := qResp.TableEntityQueryResponse.Value[0]

	// create a new batch slice.
	batch = make([]TableTransactionAction, 5)

	// create a merge action for the first added entity
	mergeProp := "MergeProperty"
	val := "foo"
	var mergeProperty = map[string]interface{}{
		partitionKey: (*entitiesToCreate)[0][partitionKey],
		rowKey:       (*entitiesToCreate)[0][rowKey],
		mergeProp:    val,
	}
	batch[0] = TableTransactionAction{ActionType: UpdateMerge, Entity: mergeProperty, ETag: (*resp.TransactionResponses)[0].Header.Get(etag)}

	// create a delete action for the second added entity
	batch[1] = TableTransactionAction{ActionType: Delete, Entity: (*entitiesToCreate)[1]}

	// create an upsert action to replace the third added entity with a new value
	replaceProp := "ReplaceProperty"
	var replaceProperties = map[string]interface{}{
		partitionKey: (*entitiesToCreate)[2][partitionKey],
		rowKey:       (*entitiesToCreate)[2][rowKey],
		replaceProp:  val,
	}
	batch[2] = TableTransactionAction{ActionType: UpsertReplace, Entity: replaceProperties}

	// Add the remaining 2 entities.
	batch[3] = TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[3]}
	batch[4] = TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[4]}

	//batch = batch[1:]

	resp, err = client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.Nil(err)

	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(http.StatusNoContent, r.StatusCode)

	}

	pager = client.Query(QueryOptions{Filter: &filter})
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]

	// The merged entity has all its properties + the merged property
	assert.Equalf(len(preMerge)+1, len(postMerge), "postMerge should have one more property than preMerge")
	assert.Equalf(postMerge[mergeProp], val, "%s property should equal %s", mergeProp, val)
}

func (s *tableClientLiveTests) TestBatchError() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexMapEntities(context, 3, "partition")

	// Create the batch.
	batch := make([]TableTransactionAction, 0, 3)

	// Sending an empty batch throws.
	_, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	assert.NotNil(err)
	assert.Equal(error_empty_transaction, err.Error())

	// Add the last entity to the table prior to adding it as part of the batch to cause a batch failure.
	client.AddEntity(ctx, (*entitiesToCreate)[2])

	// Add the entities to the batch
	for i := 0; i < cap(batch); i++ {
		batch = append(batch, TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[i]})
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	assert.NotNil(err)
	te, ok := err.(*TableTransactionError)
	require.Truef(ok, "err should be of type TableTransactionError")
	assert.Equal("EntityAlreadyExists", te.OdataError.Code)
	assert.Equal(2, te.FailedEntityIndex)
	assert.Equal(http.StatusConflict, (*resp.TransactionResponses)[0].StatusCode)
}

func (s *tableClientLiveTests) TestInvalidEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	badEntity := &map[string]interface{}{
		"Value":  10,
		"String": "stringystring",
	}

	_, err := client.AddEntity(ctx, *badEntity)

	assert.NotNil(err)
	assert.Contains(err.Error(), partitionKeyRowKeyError.Error())
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
	client := context.client.NewTableClient(tableName)
	if doCreate {
		_, err := client.Create(ctx)
		if err != nil {
			var svcErr *runtime.ResponseError
			errors.As(err, &svcErr)
			assert.FailNow(getStringFromBody(svcErr))
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
