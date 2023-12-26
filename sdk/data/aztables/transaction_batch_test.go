// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestBatchAdd(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, deleteAndStop := initClientTest(t, service, true, NewSpanValidator(t, SpanMatcher{
				Name: "Client.SubmitTransaction",
			}))
			defer deleteAndStop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)
			err = recording.AddGeneralRegexSanitizer("batch_00000000-0000-0000-0000-000000000000", "batch_[0-9A-Fa-f]{8}[-]([0-9A-Fa-f]{4}[-]?){3}[0-9a-fA-F]{12}", nil)
			require.NoError(t, err)

			entitiesToCreate := createComplexEntities(10, "partition")
			var batch []TransactionAction

			for _, e := range entitiesToCreate {
				marshalled, err := json.Marshal(e)
				require.NoError(t, err)
				batch = append(batch, TransactionAction{ActionType: TransactionTypeAdd, Entity: marshalled})
			}

			_, err = client.SubmitTransaction(ctx, batch, nil)
			require.NoError(t, err)

			pager := client.NewListEntitiesPager(nil)
			count := 0
			for pager.More() {
				response, err := pager.NextPage(ctx)
				require.NoError(t, err)
				count += len(response.Entities)
			}

			require.Equal(t, count, 10)

		})
	}
}

func TestBatchInsert(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, deleteAndStop := initClientTest(t, service, true, tracing.Provider{})
			defer deleteAndStop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)
			err = recording.AddGeneralRegexSanitizer("batch_00000000-0000-0000-0000-000000000000", "batch_[0-9A-Fa-f]{8}[-]([0-9A-Fa-f]{4}[-]?){3}[0-9a-fA-F]{12}", nil)
			require.NoError(t, err)

			entitiesToCreate := createComplexEntities(1, "partition")
			var batch []TransactionAction

			for _, e := range entitiesToCreate {
				marshalled, err := json.Marshal(e)
				require.NoError(t, err)
				batch = append(
					batch,
					TransactionAction{
						ActionType: TransactionTypeInsertMerge,
						Entity:     marshalled,
					},
				)
			}

			_, err = client.SubmitTransaction(ctx, batch, nil)
			require.NoError(t, err)

			pager := client.NewListEntitiesPager(nil)
			count := 0
			for pager.More() {
				response, err := pager.NextPage(ctx)
				require.NoError(t, err)
				count += len(response.Entities)
			}

			require.Equal(t, count, 1)
		})
	}
}

func TestBatchMixed(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, deleteAndStop := initClientTest(t, service, true, tracing.Provider{})
			defer deleteAndStop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)
			err = recording.AddGeneralRegexSanitizer("batch_00000000-0000-0000-0000-000000000000", "batch_[0-9A-Fa-f]{8}[-]([0-9A-Fa-f]{4}[-]?){3}[0-9a-fA-F]{12}", nil)
			require.NoError(t, err)

			entitiesToCreate := createComplexEntities(5, "partition")
			var batch []TransactionAction

			for _, e := range entitiesToCreate {
				marshalled, err := json.Marshal(e)
				require.NoError(t, err)
				batch = append(batch, TransactionAction{
					ActionType: TransactionTypeAdd,
					Entity:     marshalled,
				})
			}

			_, err = client.SubmitTransaction(ctx, batch, nil)
			require.NoError(t, err)

			var qResp ListEntitiesResponse
			filter := "RowKey eq '1'"
			list := &ListEntitiesOptions{Filter: &filter}
			pager := client.NewListEntitiesPager(list)
			for pager.More() {
				qResp, err = pager.NextPage(ctx)
				require.NoError(t, err)
			}
			preMerge := qResp.Entities[0]
			var unMarshalledPreMerge map[string]any
			err = json.Unmarshal(preMerge, &unMarshalledPreMerge)
			require.NoError(t, err)

			// create a new batch slice.
			var batch2 []TransactionAction

			// create a merge action for the first added entity
			mergeProp := "MergeProperty"
			val := "foo"
			var mergeEntity = map[string]any{
				partitionKey: (entitiesToCreate)[0].PartitionKey,
				rowKey:       (entitiesToCreate)[0].RowKey,
				mergeProp:    val,
			}
			marshalledMergeEntity, err := json.Marshal(mergeEntity)
			require.NoError(t, err)
			batch2 = append(batch2, TransactionAction{
				ActionType: TransactionTypeUpdateMerge,
				Entity:     marshalledMergeEntity,
			})

			// create a delete action for the second added entity
			marshalledSecondEntity, err := json.Marshal((entitiesToCreate)[1])
			require.NoError(t, err)
			batch2 = append(batch2, TransactionAction{ActionType: TransactionTypeDelete, Entity: marshalledSecondEntity})

			// create an insert action to replace the third added entity with a new value
			replaceProp := "ReplaceProperty"
			var replaceProperties = map[string]any{
				partitionKey: (entitiesToCreate)[2].PartitionKey,
				rowKey:       (entitiesToCreate)[2].RowKey,
				replaceProp:  val,
			}
			marshalledThirdEntity, err := json.Marshal(replaceProperties)
			require.NoError(t, err)
			batch2 = append(batch2, TransactionAction{ActionType: TransactionTypeInsertReplace, Entity: marshalledThirdEntity})

			// Add the remaining 2 entities.
			marshalled4thEntity, err := json.Marshal((entitiesToCreate)[3])
			require.NoError(t, err)
			marshalled5thEntity, err := json.Marshal((entitiesToCreate)[4])
			require.NoError(t, err)
			batch2 = append(batch2, TransactionAction{ActionType: TransactionTypeUpdateMerge, Entity: marshalled4thEntity})
			batch2 = append(batch2, TransactionAction{ActionType: TransactionTypeInsertMerge, Entity: marshalled5thEntity})

			_, err = client.SubmitTransaction(ctx, batch2, nil)
			require.NoError(t, err)

			pager = client.NewListEntitiesPager(list)
			for pager.More() {
				qResp, err = pager.NextPage(ctx)
				require.NoError(t, err)
			}
			postMerge := qResp.Entities[0]
			var unMarshaledPostMerge map[string]any
			err = json.Unmarshal(postMerge, &unMarshaledPostMerge)
			require.NoError(t, err)

			// The merged entity has all its properties + the merged property
			require.Equalf(t, len(unMarshalledPreMerge)+1, len(unMarshaledPostMerge), "postMerge should have one more property than preMerge")
			require.Equalf(t, unMarshaledPostMerge[mergeProp], val, "%s property should equal %s", mergeProp, val)
		})
	}
}

func TestBatchError(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, deleteAndStop := initClientTest(t, service, true, tracing.Provider{})
			defer deleteAndStop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)
			err = recording.AddGeneralRegexSanitizer("batch_00000000-0000-0000-0000-000000000000", "batch_[0-9A-Fa-f]{8}[-]([0-9A-Fa-f]{4}[-]?){3}[0-9a-fA-F]{12}", nil)
			require.NoError(t, err)

			entitiesToCreate := createComplexEntities(3, "partition")

			// Create the batch.
			var batch []TransactionAction

			// Sending an empty batch throws.
			_, err = client.SubmitTransaction(ctx, batch, nil)
			require.Error(t, err)
			require.Equal(t, errEmptyTransaction.Error(), err.Error())

			// Add the last entity to the table prior to adding it as part of the batch to cause a batch failure.
			marshalledFinalEntity, err := json.Marshal((entitiesToCreate)[2])
			require.NoError(t, err)
			_, err = client.AddEntity(ctx, marshalledFinalEntity, nil)
			require.NoError(t, err)

			// Add the entities to the batch
			for _, e := range entitiesToCreate {
				marshalledEntity, err := json.Marshal(e)
				require.NoError(t, err)
				batch = append(batch, TransactionAction{ActionType: TransactionTypeAdd, Entity: marshalledEntity})
			}

			_, err = client.SubmitTransaction(ctx, batch, nil)
			require.Error(t, err)
			var httpErr *azcore.ResponseError
			require.ErrorAs(t, err, &httpErr)
		})
	}
}

func TestBatchErrorHandleResponse(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, deleteAndStop := initClientTest(t, service, true, tracing.Provider{})
			defer deleteAndStop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)
			err = recording.AddGeneralRegexSanitizer("batch_00000000-0000-0000-0000-000000000000", "batch_[0-9A-Fa-f]{8}[-]([0-9A-Fa-f]{4}[-]?){3}[0-9a-fA-F]{12}", nil)
			require.NoError(t, err)

			entitiesToCreate := createComplexEntities(3, "partition")

			// Create the batch.
			var batch []TransactionAction

			for _, e := range entitiesToCreate {
				marshalled, err := json.Marshal(e)
				require.NoError(t, err)
				batch = append(batch, TransactionAction{
					ActionType: TransactionTypeAdd,
					Entity:     marshalled,
				})
			}

			// Add the first entity a second type
			marshalled, err := json.Marshal(entitiesToCreate[0])
			require.NoError(t, err)
			batch = append(batch, TransactionAction{
				ActionType: TransactionTypeAdd,
				Entity:     marshalled,
			})

			// Sending a batch with two entities on the same row returns an error
			_, err = client.SubmitTransaction(ctx, batch, nil)
			require.Error(t, err)
			var httpErr *azcore.ResponseError
			require.ErrorAs(t, err, &httpErr)
		})
	}
}

func TestBatchComplex(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client, deleteAndStop := initClientTest(t, service, true, tracing.Provider{})
			defer deleteAndStop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)
			err = recording.AddGeneralRegexSanitizer("batch_00000000-0000-0000-0000-000000000000", "batch_[0-9A-Fa-f]{8}[-]([0-9A-Fa-f]{4}[-]?){3}[0-9a-fA-F]{12}", nil)
			require.NoError(t, err)

			edmEntity := createEdmEntity(1, "pk01")
			edmEntity2 := createEdmEntity(2, "pk01")
			edmEntity3 := createEdmEntity(3, "pk01")
			edmEntity4 := createEdmEntity(4, "pk01")
			edmEntity5 := createEdmEntity(5, "pk01")
			var batch []TransactionAction

			marshalled1, err := json.Marshal(edmEntity)
			require.NoError(t, err)
			batch = append(batch, TransactionAction{
				ActionType: TransactionTypeAdd,
				Entity:     marshalled1,
			})

			marshalled2, err := json.Marshal(edmEntity2)
			require.NoError(t, err)
			batch = append(batch, TransactionAction{
				ActionType: TransactionTypeAdd,
				Entity:     marshalled2,
			})

			marshalled3, err := json.Marshal(edmEntity3)
			require.NoError(t, err)
			batch = append(batch, TransactionAction{
				ActionType: TransactionTypeAdd,
				Entity:     marshalled3,
			})

			marshalled4, err := json.Marshal(edmEntity4)
			require.NoError(t, err)
			batch = append(batch, TransactionAction{
				ActionType: TransactionTypeAdd,
				Entity:     marshalled4,
			})

			marshalled5, err := json.Marshal(edmEntity5)
			require.NoError(t, err)
			batch = append(batch, TransactionAction{
				ActionType: TransactionTypeAdd,
				Entity:     marshalled5,
			})

			_, err = client.SubmitTransaction(ctx, batch, nil)
			require.NoError(t, err)

			var batch2 []TransactionAction
			edmEntity.Properties["Bool"] = false
			edmEntity2.Properties["Int32"] = int32(10)

			marshalled1, err = json.Marshal(edmEntity)
			require.NoError(t, err)
			batch2 = append(batch2, TransactionAction{
				ActionType: TransactionTypeInsertMerge,
				Entity:     marshalled1,
			})

			marshalled2, err = json.Marshal(edmEntity2)
			require.NoError(t, err)
			batch2 = append(batch2, TransactionAction{
				ActionType: TransactionTypeInsertReplace,
				Entity:     marshalled2,
			})

			marshalled3, err = json.Marshal(edmEntity3)
			require.NoError(t, err)
			batch2 = append(batch2, TransactionAction{
				ActionType: TransactionTypeDelete,
				Entity:     marshalled3,
			})

			_, err = client.SubmitTransaction(ctx, batch2, nil)
			require.NoError(t, err)

			received, err := client.GetEntity(ctx, edmEntity.PartitionKey, edmEntity.RowKey, nil)
			require.NoError(t, err)

			var receivedEdm EDMEntity
			err = json.Unmarshal(received.Value, &receivedEdm)
			require.NoError(t, err)
			require.Equal(t, edmEntity.Properties["Bool"], receivedEdm.Properties["Bool"])

			received2, err := client.GetEntity(ctx, edmEntity2.PartitionKey, edmEntity2.RowKey, nil)
			require.NoError(t, err)

			var receivedEdm2 EDMEntity
			err = json.Unmarshal(received2.Value, &receivedEdm2)
			require.NoError(t, err)
			require.Equal(t, edmEntity2.Properties["Int32"], receivedEdm2.Properties["Int32"])

			_, err = client.GetEntity(ctx, edmEntity3.PartitionKey, edmEntity3.RowKey, nil)
			require.Error(t, err)
			var httpErr *azcore.ResponseError
			require.ErrorAs(t, err, &httpErr)
			require.Equal(t, string(ResourceNotFound), httpErr.ErrorCode)
		})
	}
}
