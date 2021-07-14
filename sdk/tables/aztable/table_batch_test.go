// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func (s *tableClientLiveTests) TestBatchError() {
	assert := assert.New(s.T())
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(3, "partition")

	// Create the batch.
	batch := make([]TableTransactionAction, 0, 3)

	// Sending an empty batch throws.
	_, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	assert.NotNil(err)
	assert.Equal(error_empty_transaction, err.Error())

	// Add the last entity to the table prior to adding it as part of the batch to cause a batch failure.
	marshalledFinalEntity, err := json.Marshal((*entitiesToCreate)[2])
	client.AddEntity(ctx, marshalledFinalEntity)

	// Add the entities to the batch
	for i := 0; i < cap(batch); i++ {
		marshalledEntity, err := json.Marshal((*entitiesToCreate)[i])
		require.Nil(err)
		batch = append(batch, TableTransactionAction{ActionType: Add, Entity: marshalledEntity})
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	assert.NotNil(err)
	transactionError, ok := err.(*TableTransactionError)
	require.Truef(ok, "err should be of type TableTransactionError")
	assert.Equal("EntityAlreadyExists", transactionError.OdataError.Code)
	assert.Equal(2, transactionError.FailedEntityIndex)
	assert.Equal(http.StatusConflict, (*resp.TransactionResponses)[0].StatusCode)
}
