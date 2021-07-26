// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *tableClientLiveTests) TestBatchAdd() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(10, "partition")
	batch := make([]TableTransactionAction, 10)

	for i, e := range *entitiesToCreate {
		marshalled, err := json.Marshal(e)
		require.NoError(err)
		batch[i] = TableTransactionAction{ActionType: Add, Entity: marshalled}
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.NoError(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		require.Equal(r.StatusCode, http.StatusNoContent)
	}

	pager := client.List(nil)
	count := 0
	for pager.NextPage(ctx) {
		response := pager.PageResponse()
		count += len(response.TableEntityQueryResponse.Value)
	}

	require.Equal(count, 10)
}

func (s *tableClientLiveTests) TestBatchMixed() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(5, "partition")
	batch := make([]TableTransactionAction, 3)

	for i := range batch {
		marshalled, err := json.Marshal((*entitiesToCreate)[i])
		require.NoError(err)
		batch[i] = TableTransactionAction{
			ActionType: Add,
			Entity:     marshalled,
		}
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.NoError(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		require.Equal(http.StatusNoContent, r.StatusCode)
	}

	var qResp TableEntityQueryByteResponseResponse
	filter := "RowKey eq '1'"
	list := &ListOptions{Filter: &filter}
	pager := client.List(list)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	preMerge := qResp.TableEntityQueryResponse.Value[0]
	var unMarshalledPreMerge map[string]interface{}
	err = json.Unmarshal(preMerge, &unMarshalledPreMerge)
	require.NoError(err)

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
	require.NoError(err)
	batch[0] = TableTransactionAction{ActionType: UpdateMerge, Entity: marshalledMergeEntity, ETag: (*resp.TransactionResponses)[0].Header.Get(etag)}

	// create a delete action for the second added entity
	marshalledSecondEntity, err := json.Marshal((*entitiesToCreate)[1])
	require.NoError(err)
	batch[1] = TableTransactionAction{ActionType: Delete, Entity: marshalledSecondEntity}

	// create an upsert action to replace the third added entity with a new value
	replaceProp := "ReplaceProperty"
	var replaceProperties = map[string]interface{}{
		partitionKey: (*entitiesToCreate)[2].PartitionKey,
		rowKey:       (*entitiesToCreate)[2].RowKey,
		replaceProp:  val,
	}
	marshalledThirdEntity, err := json.Marshal(replaceProperties)
	require.NoError(err)
	batch[2] = TableTransactionAction{ActionType: UpsertReplace, Entity: marshalledThirdEntity}

	// Add the remaining 2 entities.
	marshalled4thEntity, err := json.Marshal((*entitiesToCreate)[3])
	marshalled5thEntity, err := json.Marshal((*entitiesToCreate)[4])
	batch[3] = TableTransactionAction{ActionType: Add, Entity: marshalled4thEntity}
	batch[4] = TableTransactionAction{ActionType: Add, Entity: marshalled5thEntity}

	resp, err = client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.NoError(err)

	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		require.Equal(http.StatusNoContent, r.StatusCode)

	}

	pager = client.List(list)
	for pager.NextPage(ctx) {
		qResp = pager.PageResponse()
	}
	postMerge := qResp.TableEntityQueryResponse.Value[0]
	var unMarshaledPostMerge map[string]interface{}
	err = json.Unmarshal(postMerge, &unMarshaledPostMerge)
	require.NoError(err)

	// The merged entity has all its properties + the merged property
	require.Equalf(len(unMarshalledPreMerge)+1, len(unMarshaledPostMerge), "postMerge should have one more property than preMerge")
	require.Equalf(unMarshaledPostMerge[mergeProp], val, "%s property should equal %s", mergeProp, val)
}

func (s *tableClientLiveTests) TestBatchError() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	entitiesToCreate := createComplexEntities(3, "partition")

	// Create the batch.
	batch := make([]TableTransactionAction, 0, 3)

	// Sending an empty batch throws.
	_, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.NotNil(err)
	require.Equal(error_empty_transaction, err.Error())

	// Add the last entity to the table prior to adding it as part of the batch to cause a batch failure.
	marshalledFinalEntity, err := json.Marshal((*entitiesToCreate)[2])
	client.AddEntity(ctx, marshalledFinalEntity)

	// Add the entities to the batch
	for i := 0; i < cap(batch); i++ {
		marshalledEntity, err := json.Marshal((*entitiesToCreate)[i])
		require.NoError(err)
		batch = append(batch, TableTransactionAction{ActionType: Add, Entity: marshalledEntity})
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.NotNil(err)
	transactionError, ok := err.(*TableTransactionError)
	require.Truef(ok, "err should be of type TableTransactionError")
	require.Equal("EntityAlreadyExists", transactionError.OdataError.Code)
	require.Equal(2, transactionError.FailedEntityIndex)
	require.Equal(http.StatusConflict, (*resp.TransactionResponses)[0].StatusCode)
}

func (s *tableClientLiveTests) TestBatchComplex() {
	require := require.New(s.T())
	context := getTestContext(s.T().Name())
	client, delete := s.init(true)
	defer delete()

	edmEntity := createEdmEntity(1, "pk001")
	edmEntity2 := createEdmEntity(2, "pk001")
	edmEntity3 := createEdmEntity(3, "pk001")
	edmEntity4 := createEdmEntity(4, "pk001")
	edmEntity5 := createEdmEntity(5, "pk001")
	batch := make([]TableTransactionAction, 5)

	marshalled1, err := json.Marshal(edmEntity)
	require.NoError(err)
	batch[0] = TableTransactionAction{
		ActionType: Add,
		Entity:     marshalled1,
	}

	marshalled2, err := json.Marshal(edmEntity2)
	require.NoError(err)
	batch[1] = TableTransactionAction{
		ActionType: Add,
		Entity:     marshalled2,
	}

	marshalled3, err := json.Marshal(edmEntity3)
	require.NoError(err)
	batch[2] = TableTransactionAction{
		ActionType: Add,
		Entity:     marshalled3,
	}

	marshalled4, err := json.Marshal(edmEntity4)
	require.NoError(err)
	batch[3] = TableTransactionAction{
		ActionType: Add,
		Entity:     marshalled4,
	}

	marshalled5, err := json.Marshal(edmEntity5)
	require.NoError(err)
	batch[4] = TableTransactionAction{
		ActionType: Add,
		Entity:     marshalled5,
	}

	resp, err := client.submitTransactionInternal(ctx, &batch, context.recording.UUID(), context.recording.UUID(), nil)
	require.NoError(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		require.Equal(http.StatusNoContent, r.StatusCode)
	}

	batch2 := make([]TableTransactionAction, 3)
	edmEntity.Properties["Bool"] = false
	edmEntity2.Properties["Int32"] = int32(10)

	marshalled1, err = json.Marshal(edmEntity)
	require.NoError(err)
	batch2[0] = TableTransactionAction{
		ActionType: UpsertMerge,
		Entity:     marshalled1,
	}

	marshalled2, err = json.Marshal(edmEntity2)
	require.NoError(err)
	batch2[1] = TableTransactionAction{
		ActionType: UpsertReplace,
		Entity:     marshalled2,
	}

	marshalled3, err = json.Marshal(edmEntity3)
	require.NoError(err)
	batch2[2] = TableTransactionAction{
		ActionType: Delete,
		Entity:     marshalled3,
	}

	resp, err = client.submitTransactionInternal(ctx, &batch2, context.recording.UUID(), context.recording.UUID(), nil)
	require.NoError(err)
	for i := 0; i < len(*resp.TransactionResponses); i++ {
		r := (*resp.TransactionResponses)[i]
		require.Equal(http.StatusNoContent, r.StatusCode)
	}

	received, err := client.GetEntity(ctx, edmEntity.PartitionKey, edmEntity.RowKey, nil)
	require.NoError(err)

	var receivedEdm EdmEntity
	err = json.Unmarshal(received.Value, &receivedEdm)
	require.NoError(err)
	require.Equal(edmEntity.Properties["Bool"], receivedEdm.Properties["Bool"])

	received2, err := client.GetEntity(ctx, edmEntity2.PartitionKey, edmEntity2.RowKey, nil)
	require.NoError(err)

	var receivedEdm2 EdmEntity
	err = json.Unmarshal(received2.Value, &receivedEdm2)
	require.NoError(err)
	require.Equal(edmEntity2.Properties["Int32"], receivedEdm2.Properties["Int32"])

	_, err = client.GetEntity(ctx, edmEntity3.PartitionKey, edmEntity3.RowKey, nil)
	require.Error(err)
}
