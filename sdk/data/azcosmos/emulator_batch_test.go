// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestEmulatorTransactionalBatchCommitsAllFiveOperations(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	key := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, key, id)
	other := id + "-other"
	trackEmulatorItem(t, container, key, other)
	item := func(itemID, value string) []byte {
		body, err := json.Marshal(map[string]any{"id": itemID, "pk": id, "value": value})
		require.NoError(t, err)
		return body
	}
	batch := NewTransactionalBatch(key)
	require.NoError(t, batch.CreateItem(item(id, "created"), nil))
	require.NoError(t, batch.ReadItem(id, nil))
	require.NoError(t, batch.UpsertItem(item(other, "upserted"), nil))
	require.NoError(t, batch.ReplaceItem(id, item(id, "replaced"), nil))
	require.NoError(t, batch.DeleteItem(other, nil))
	enabled := true
	response, err := container.ExecuteTransactionalBatch(t.Context(), batch, &TransactionalBatchOptions{
		Operation: OperationOptions{EnableContentResponseOnWrite: &enabled},
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Len(t, response.OperationResults, 5)
	require.Equal(t, []int{201, 200, 201, 200, 204}, []int{
		response.OperationResults[0].StatusCode, response.OperationResults[1].StatusCode,
		response.OperationResults[2].StatusCode, response.OperationResults[3].StatusCode, response.OperationResults[4].StatusCode,
	})
	require.Positive(t, response.RequestCharge)
	require.NotEmpty(t, response.ActivityID)
	require.NotEmpty(t, response.SessionToken)
	require.NotEmpty(t, response.OperationResults[0].ETag)
	require.NotEmpty(t, response.OperationResults[3].ETag)
	var read map[string]any
	require.NoError(t, json.Unmarshal(response.OperationResults[1].ResourceBody, &read))
	require.Equal(t, "created", read["value"], "the read must observe the prior operation in its transaction")
	persisted, err := container.ReadItem(t.Context(), key, id, &ReadItemOptions{SessionToken: response.SessionToken})
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(persisted.Value, &read))
	require.Equal(t, "replaced", read["value"])
	_, err = container.ReadItem(t.Context(), key, other, nil)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeNotFound, cosmosErr.Code)
}

func TestEmulatorTransactionalBatchReadIfNoneMatch(t *testing.T) {
	for _, test := range []struct {
		name    string
		matches bool
		status  int
	}{
		{"matching ETag", true, http.StatusNotModified},
		{"different ETag", false, http.StatusOK},
	} {
		t.Run(test.name, func(t *testing.T) {
			container := emulatorContainer(t)
			id := uniqueItemID(t)
			key := NewPartitionKeyString(id)
			trackEmulatorItem(t, container, key, id)
			body, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": "original"})
			require.NoError(t, err)
			original, err := container.CreateItem(t.Context(), key, id, body, nil)
			require.NoError(t, err)
			require.NotEmpty(t, original.ETag)
			etag := original.ETag
			if !test.matches {
				etag = `"different"`
			}
			batch := NewTransactionalBatch(key)
			require.NoError(t, batch.ReadItem(id, &TransactionalBatchReadItemOptions{IfNoneMatchETag: &etag}))
			disabled := false
			response, err := container.ExecuteTransactionalBatch(t.Context(), batch, &TransactionalBatchOptions{
				Operation: OperationOptions{EnableContentResponseOnWrite: &disabled}, SessionToken: original.SessionToken,
			})
			require.NoError(t, err)
			require.Len(t, response.OperationResults, 1)
			result := response.OperationResults[0]
			require.Equal(t, test.status, result.StatusCode)
			require.Equal(t, original.ETag, result.ETag)
			if test.matches {
				require.Contains(t, []string{"", "null"}, string(result.ResourceBody))
			} else {
				var document map[string]any
				require.NoError(t, json.Unmarshal(result.ResourceBody, &document))
				require.Equal(t, "original", document["value"])
			}
			persisted, err := container.ReadItem(t.Context(), key, id, nil)
			require.NoError(t, err)
			require.Equal(t, original.ETag, persisted.ETag, "the conditional batch read must not modify the item")
		})
	}
}

func TestEmulatorTransactionalBatchRollbackPreservesPersistedState(t *testing.T) {
	for _, failure := range []string{"missing-read", "duplicate-create", "conditional-replace"} {
		t.Run(failure, func(t *testing.T) {
			container := emulatorContainer(t)
			existing := uniqueItemID(t)
			created := existing + "-rolled-back"
			key := NewPartitionKeyString(existing)
			trackEmulatorItem(t, container, key, existing)
			trackEmulatorItem(t, container, key, created)
			item := func(id, value string) []byte {
				body, err := json.Marshal(map[string]any{"id": id, "pk": existing, "value": value})
				require.NoError(t, err)
				return body
			}
			original, err := container.CreateItem(t.Context(), key, existing, item(existing, "original"), nil)
			require.NoError(t, err)
			batch := NewTransactionalBatch(key)
			require.NoError(t, batch.CreateItem(item(created, "uncommitted"), nil))
			status := 404
			switch failure {
			case "missing-read":
				require.NoError(t, batch.ReadItem(existing+"-missing", nil))
			case "duplicate-create":
				status = 409
				require.NoError(t, batch.CreateItem(item(existing, "duplicate"), nil))
			case "conditional-replace":
				status = 412
				wrong := azcore.ETag(`"wrong"`)
				require.NoError(t, batch.ReplaceItem(existing, item(existing, "uncommitted"), &TransactionalBatchReplaceItemOptions{IfMatchETag: &wrong}))
			}
			require.NoError(t, batch.DeleteItem(existing, nil))
			response, err := container.ExecuteTransactionalBatch(t.Context(), batch, nil)
			require.NoError(t, err, "a received rollback is a transaction result, not an execution error")
			require.Equal(t, http.StatusMultiStatus, response.StatusCode)
			require.Len(t, response.OperationResults, 3)
			require.Equal(t, 424, response.OperationResults[0].StatusCode)
			require.Equal(t, status, response.OperationResults[1].StatusCode)
			require.Equal(t, 424, response.OperationResults[2].StatusCode)
			require.Positive(t, response.RequestCharge)

			persisted, err := container.ReadItem(t.Context(), key, existing, nil)
			require.NoError(t, err)
			require.Equal(t, original.ETag, persisted.ETag, "the failed transaction must not modify or delete the existing item")
			var document map[string]any
			require.NoError(t, json.Unmarshal(persisted.Value, &document))
			require.Equal(t, "original", document["value"])
			_, err = container.ReadItem(t.Context(), key, created, nil)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, CodeNotFound, cosmosErr.Code, "the preceding create must actually be rolled back")
		})
	}
}

func TestEmulatorTransactionalBatchCompleteHierarchicalKeysAndIsolation(t *testing.T) {
	client, databaseID, _ := emulatorClient(t)
	for _, levels := range []int{2, 3} {
		containerID := "query-hierarchical"
		if levels == 3 {
			containerID = "batch-hierarchical"
		}
		container, err := client.NewContainer(databaseID, containerID)
		require.NoError(t, err)
		id := uniqueItemID(t) + containerID
		prefix := NewPartitionKeyString(id)
		for i, key := range []PartitionKey{prefix.AppendNull(), prefix.AppendUndefined(), prefix.AppendNumber(42)} {
			document := map[string]any{"id": id, "pk": id, "value": i}
			if i == 0 {
				document["child"] = nil
			} else if i == 2 {
				document["child"] = 42
			}
			if levels == 3 {
				key = key.AppendBool(false)
				document["leaf"] = false
			}
			trackEmulatorItem(t, container, key, id)
			body, err := json.Marshal(document)
			require.NoError(t, err)
			batch := NewTransactionalBatch(key)
			require.NoError(t, batch.CreateItem(body, nil))
			require.NoError(t, batch.ReadItem(id, nil))
			response, err := container.ExecuteTransactionalBatch(t.Context(), batch, nil)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response.StatusCode)
			require.Len(t, response.OperationResults, 2)
			require.Equal(t, http.StatusCreated, response.OperationResults[0].StatusCode)
			require.Equal(t, http.StatusOK, response.OperationResults[1].StatusCode)
			var read map[string]any
			require.NoError(t, json.Unmarshal(response.OperationResults[1].ResourceBody, &read))
			require.Equal(t, float64(i), read["value"], "null and undefined HPK components must route distinctly")
		}
		missing := id + "-prefix-rejected"
		for _, key := range []PartitionKey{prefix, prefix.AppendNull()} {
			if key.Len() == levels {
				continue
			}
			batch := NewTransactionalBatch(key)
			body, err := json.Marshal(map[string]any{"id": missing, "pk": id, "child": nil, "leaf": false})
			require.NoError(t, err)
			require.NoError(t, batch.CreateItem(body, nil))
			response, err := container.ExecuteTransactionalBatch(t.Context(), batch, nil)
			require.Zero(t, response)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, CodeBadRequest, cosmosErr.Code)
			require.Contains(t, string(cosmosErr.Body), `"partitionKey"`, "retain the fetched metadata that caused local prefix rejection")
		}
		complete := prefix.AppendNull()
		if levels == 3 {
			complete = complete.AppendBool(false)
		}
		_, err = container.ReadItem(t.Context(), complete, missing, nil)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeNotFound, cosmosErr.Code, "prefix validation must reject before any mutation")
	}
}

func TestEmulatorTransactionalBatchContentResponsesAndLifetime(t *testing.T) {
	disabled, enabled := false, true
	for _, setting := range []*bool{nil, &disabled, &enabled} {
		client, databaseID, containerID := emulatorClientWithOptions(t, &ClientOptions{EnableContentResponseOnWrite: &disabled})
		container, err := client.NewContainer(databaseID, containerID)
		require.NoError(t, err)
		id := uniqueItemID(t)
		body, err := json.Marshal(map[string]any{"id": id, "pk": id, "value": "owned"})
		require.NoError(t, err)
		batch := NewTransactionalBatch(NewPartitionKeyString(id))
		require.NoError(t, batch.CreateItem(body, nil))
		require.NoError(t, batch.ReadItem(id, nil))
		require.NoError(t, batch.DeleteItem(id, nil))
		for i := range body {
			body[i] = 'x'
		}
		response, err := container.ExecuteTransactionalBatch(t.Context(), batch, &TransactionalBatchOptions{
			Operation: OperationOptions{EnableContentResponseOnWrite: setting},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Len(t, response.OperationResults, 3)
		require.Equal(t, http.StatusCreated, response.OperationResults[0].StatusCode)
		require.Equal(t, http.StatusOK, response.OperationResults[1].StatusCode)
		require.Equal(t, http.StatusNoContent, response.OperationResults[2].StatusCode)
		require.NoError(t, client.Close())
		if setting != nil && *setting {
			require.NotEmpty(t, response.OperationResults[0].ResourceBody)
		} else {
			require.Contains(t, []string{"", "null"}, string(response.OperationResults[0].ResourceBody))
		}
		var document map[string]any
		require.NoError(t, json.Unmarshal(response.OperationResults[1].ResourceBody, &document))
		require.Equal(t, "owned", document["value"], "reads and Go-owned results must survive disabled write bodies and client shutdown")
		require.NotEmpty(t, response.Body)
	}
}
