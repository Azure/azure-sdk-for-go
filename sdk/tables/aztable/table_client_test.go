// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
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

	marshalledEntity, err := json.Marshal((*entitiesToCreate)[0])
	assert.Nil(err)
	resp, err := client.AddEntity(ctx, marshalledEntity)
	assert.Nil(err)
	assert.NotNil(resp)
}

func (s *tableClientLiveTests) TestAddComplexEntity() {
	assert := assert.New(s.T())
	// context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(1, "partition")

	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		assert.Nil(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		var svcErr *runtime.ResponseError
		errors.As(err, &svcErr)
		assert.Nilf(err, getStringFromBody(svcErr))
	}
}

func (s *tableClientLiveTests) TestDeleteEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	simpleEntity := createSimpleEntity(1, "partition")

	marshalledEntity, err := json.Marshal(simpleEntity)
	assert.Nil(err)
	_, err = client.AddEntity(ctx, marshalledEntity)
	assert.Nil(err)
	_, delErr := client.DeleteEntity(ctx, simpleEntity.PartitionKey, simpleEntity.RowKey, nil)
	assert.Nil(delErr)
}

func (s *tableClientLiveTests) TestMergeEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	entityToCreate := createSimpleEntity(1, "partition")
	marshalled := marshalBasicEntity(entityToCreate, assert)

	_, err := client.AddEntity(ctx, *marshalled)
	assert.Nil(err)

	filter := "RowKey eq '1'"
	queryOptions := &QueryOptions{Filter: &filter}

	preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey)
	assert.Nil(err)

	var unMarshalledPreMerge map[string]interface{}
	err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
	assert.Nil(err)

	var mapEntity map[string]interface{}
	err = json.Unmarshal(*marshalled, &mapEntity)
	assert.Nil(err)
	mapEntity["MergeProperty"] = "foo"

	reMarshalled, err := json.Marshal(mapEntity)
	assert.Nil(err)

	_, updateErr := client.UpdateEntity(ctx, reMarshalled, nil, Merge)
	assert.Nil(updateErr)

	var qResp TableEntityQueryByteResponseResponse
	pager := client.Query(queryOptions)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unmarshalledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unmarshalledPostMerge)
	assert.Nil(err)

	assert.Equal(unmarshalledPostMerge["PartitionKey"], unMarshalledPreMerge["PartitionKey"])
	assert.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

	_, ok := unMarshalledPreMerge["MergeProperty"]
	assert.False(ok)
}

func (s *tableClientLiveTests) TestUpsertEntity() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// 1. Create Basic Entity
	entityToCreate := createSimpleEntity(1, "partition")
	marshalled := marshalBasicEntity(entityToCreate, assert)

	_, err := client.UpsertEntity(ctx, *marshalled, Replace)
	assert.Nil(err)

	filter := "RowKey eq '1'"
	query := &QueryOptions{Filter: &filter}

	// 2. Query for basic Entity
	preMerge, err := client.GetEntity(ctx, entityToCreate.PartitionKey, entityToCreate.RowKey)
	assert.Nil(err)

	var unMarshalledPreMerge map[string]interface{}
	err = json.Unmarshal(preMerge.Value, &unMarshalledPreMerge)
	assert.Nil(err)

	// 3. Create same entity without Bool property, add "MergeProperty" prop
	mapEntity := createSimpleEntityNoBool(1, "partition")
	mapEntity["MergeProperty"] = "foo"

	reMarshalled, err := json.Marshal(mapEntity)

	// 4. Replace Entity with "bool"-less entity
	_, err = client.UpsertEntity(ctx, reMarshalled, Replace)
	require.Nil(err)

	// 5. Query for new entity
	var qResp TableEntityQueryByteResponseResponse
	pager := client.Query(query)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unmarshalledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unmarshalledPostMerge)

	// 6. Make assertions
	assert.Less(len(unmarshalledPostMerge), len(unMarshalledPreMerge))
	assert.Equal(unmarshalledPostMerge["MergeProperty"], "foo")

	_, ok := unmarshalledPostMerge["Bool"]
	assert.Falsef(ok, "Bool property should not be available in the merged entity")
}

func (s *tableClientLiveTests) TestQuerySimpleEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createSimpleEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		assert.Nil(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		assert.Nil(err)
	}

	filter := "RowKey lt '5'"
	query := &QueryOptions{Filter: &filter}
	expectedCount := 4

	var resp TableEntityQueryByteResponseResponse
	pager := client.Query(query)
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		assert.Equal(len(resp.TableEntityQueryResponse.Value), expectedCount)
	}

	for i, e := range resp.TableEntityQueryResponse.Value {
		var mapModel map[string]interface{}
		err := json.Unmarshal(e, &mapModel)
		assert.Nil(err)

		_, ok := mapModel[timestamp]
		assert.True(ok)

		_, ok = mapModel[etagOdata]
		assert.True(ok)

		var b basicTestEntity
		err = json.Unmarshal(e, &b)
		assert.Nil(err)

		assert.Equal(b.PartitionKey, "partition")
		assert.Equal(b.RowKey, fmt.Sprint(i+1))
		assert.Equal(b.String, (*entitiesToCreate)[i].String)
		assert.Equal(b.Integer, (*entitiesToCreate)[i].Integer)
		assert.Equal(b.Bool, (*entitiesToCreate)[i].Bool)
	}
}

func (s *tableClientLiveTests) TestQueryComplexEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	// Add 5 entities
	entitiesToCreate := createComplexEntities(5, "partition")
	for _, e := range *entitiesToCreate {
		marshalledEntity, err := json.Marshal(e)
		assert.Nil(err)
		_, err = client.AddEntity(ctx, marshalledEntity)
		assert.Nil(err)
	}

	filter := "RowKey lt '5'"
	expectedCount := 4
	query := &QueryOptions{Filter: &filter}

	var resp TableEntityQueryByteResponseResponse
	pager := client.Query(query)
	for pager.NextPage(ctx) {
		resp = pager.PageResponse()
		assert.Equal(expectedCount, len(resp.TableEntityQueryResponse.Value))

		for idx, entity := range resp.TableEntityQueryResponse.Value {
			model := complexTestEntity{}
			err := json.Unmarshal(entity, &model)
			assert.Nil(err)

			assert.Equal(model.PartitionKey, "partition")
			assert.Equal(model.RowKey, (*entitiesToCreate)[idx].RowKey)
			assert.Equal(model.Integer, (*entitiesToCreate)[idx].Integer)
			assert.Equal(model.String, (*entitiesToCreate)[idx].String)
			assert.Equal(model.Bool, (*entitiesToCreate)[idx].Bool)
			assert.Equal(model.Float, (*entitiesToCreate)[idx].Float)
			assert.Equal(model.DateTime, (*entitiesToCreate)[idx].DateTime)
			assert.Equal(model.Byte, (*entitiesToCreate)[idx].Byte)
		}

	}
}

func (s *tableClientLiveTests) TestBatchAdd() {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(10, "partition")
	batch := make([]TableTransactionAction, 10)

	for i, e := range *entitiesToCreate {
		marshalled, err := json.Marshal(e)
		assert.Nil(err)
		batch[i] = TableTransactionAction{ActionType: Add, Entity: marshalled}
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	assert.Nil(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(r.StatusCode, http.StatusNoContent)
	}

	pager := client.Query(nil)
	count := 0
	for pager.NextPage(ctx) {
		response := pager.PageResponse()
		count += len(response.TableEntityQueryResponse.Value)
	}

	assert.Equal(count, 10)
}

func (s *tableClientLiveTests) TestBatchMixed() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(5, "partition")
	batch := make([]TableTransactionAction, 3)

	for i := range batch {
		marshalled, err := json.Marshal((*entitiesToCreate)[i])
		require.Nil(err)
		batch[i] = TableTransactionAction{
			ActionType: Add,
			Entity:     marshalled,
		}
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.Nil(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(http.StatusNoContent, r.StatusCode)
	}

	var qResp TableEntityQueryByteResponseResponse
	filter := "RowKey eq '1'"
	query := &QueryOptions{Filter: &filter}
	pager := client.Query(query)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	preMerge := qResp.TableEntityQueryResponse.Value[0]
	var unMarshalledPreMerge map[string]interface{}
	err = json.Unmarshal(preMerge, &unMarshalledPreMerge)
	require.Nil(err)

	// create a new batch slice.
	batch = make([]TableTransactionAction, 5)

	// create a merge action for the first added entity
	mergeProp := "MergeProperty"
	val := "foo"
	var mergeEntity = map[string]interface{}{
		partitionKey: (*entitiesToCreate)[0].PartitionKey,
		rowKey:       (*entitiesToCreate)[0].RowKey,
		mergeProp:    val,
	}
	marshalledMergeEntity, err := json.Marshal(mergeEntity)
	require.Nil(err)
	batch[0] = TableTransactionAction{ActionType: UpdateMerge, Entity: marshalledMergeEntity, ETag: (*resp.TransactionResponses)[0].Header.Get(etag)}

	// create a delete action for the second added entity
	marshalledSecondEntity, err := json.Marshal((*entitiesToCreate)[1])
	require.Nil(err)
	batch[1] = TableTransactionAction{ActionType: Delete, Entity: marshalledSecondEntity}

	// create an upsert action to replace the third added entity with a new value
	replaceProp := "ReplaceProperty"
	var replaceProperties = map[string]interface{}{
		partitionKey: (*entitiesToCreate)[2].PartitionKey,
		rowKey:       (*entitiesToCreate)[2].RowKey,
		replaceProp:  val,
	}
	marshalledThirdEntity, err := json.Marshal(replaceProperties)
	require.Nil(err)
	batch[2] = TableTransactionAction{ActionType: UpsertReplace, Entity: marshalledThirdEntity}

	// Add the remaining 2 entities.
	marshalled4thEntity, err := json.Marshal((*entitiesToCreate)[3])
	marshalled5thEntity, err := json.Marshal((*entitiesToCreate)[4])
	batch[3] = TableTransactionAction{ActionType: Add, Entity: marshalled4thEntity}
	batch[4] = TableTransactionAction{ActionType: Add, Entity: marshalled5thEntity}

	resp, err = client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.Nil(err)

	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		assert.Equal(http.StatusNoContent, r.StatusCode)

	}

	pager = client.Query(query)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unMarshaledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unMarshaledPostMerge)
	require.Nil(err)

	// The merged entity has all its properties + the merged property
	assert.Equalf(len(unMarshalledPreMerge)+1, len(unMarshaledPostMerge), "postMerge should have one more property than preMerge")
	assert.Equalf(unMarshaledPostMerge[mergeProp], val, "%s property should equal %s", mergeProp, val)
}

// func (s *tableClientLiveTests) TestBatchError() {
// 	assert := assert.New(s.T())
// 	require := require.New(s.T())
// 	context := getTestContext(s.T().Name())
// 	client, delete := s.init(true)
// 	defer delete()

// 	entitiesToCreate := createComplexMapEntities(context, 3, "partition")

// 	// Create the batch.
// 	batch := make([]TableTransactionAction, 0, 3)

// 	// Sending an empty batch throws.
// 	_, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
// 	assert.NotNil(err)
// 	assert.Equal(error_empty_transaction, err.Error())

// 	// Add the last entity to the table prior to adding it as part of the batch to cause a batch failure.
// 	client.AddEntity(ctx, (*entitiesToCreate)[2])

// 	// Add the entities to the batch
// 	for i := 0; i < cap(batch); i++ {
// 		batch = append(batch, TableTransactionAction{ActionType: Add, Entity: (*entitiesToCreate)[i]})
// 	}

// 	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
// 	assert.NotNil(err)
// 	te, ok := err.(*TableTransactionError)
// 	require.Truef(ok, "err should be of type TableTransactionError")
// 	assert.Equal("EntityAlreadyExists", te.OdataError.Code)
// 	assert.Equal(2, te.FailedEntityIndex)
// 	assert.Equal(http.StatusConflict, (*resp.TransactionResponses)[0].StatusCode)
// }

func (s *tableClientLiveTests) TestInvalidEntity() {
	assert := assert.New(s.T())
	client, delete := s.init(true)
	defer delete()

	badEntity := map[string]interface{}{
		"Value":  10,
		"String": "stringystring",
	}

	badEntityMarshalled, err := json.Marshal(badEntity)
	_, err = client.AddEntity(ctx, badEntityMarshalled)

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

func (s *tableClientLiveTests) init(createTable bool) (*TableClient, func()) {
	assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	tableName, _ := getTableName(context)
	client := context.client.NewTableClient(tableName)
	if createTable {
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
