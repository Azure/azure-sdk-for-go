// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

func TestBatchAtomicity(t *testing.T) {
	t.Run("batch.atomicity", func(t *testing.T) {
		fx := newFixture(t, "batch.atomicity", "")
		ctx := t.Context()
		pkValue := uniqueID(t)
		pk := azcosmos.NewPartitionKeyString(pkValue)
		createdID := uniqueID(t) + "-created"
		replacedID := uniqueID(t) + "-replaced"
		trackItem(t, fx.Container, pk, createdID)
		trackItem(t, fx.Container, pk, replacedID)

		batch := fx.Container.NewTransactionalBatch(pk)
		require.NoError(t, batch.CreateItem(map[string]any{"id": createdID, "pk": pkValue, "value": "created"}, nil))
		require.NoError(t, batch.ReadItem(createdID, nil))
		require.NoError(t, batch.ReplaceItem(createdID,
			map[string]any{"id": createdID, "pk": pkValue, "value": "replaced"}, nil))
		require.NoError(t, batch.CreateItem(map[string]any{"id": replacedID, "pk": pkValue, "value": "final"}, nil))

		enabled := true
		response, err := fx.Container.ExecuteTransactionalBatch(ctx, batch, &azcosmos.TransactionalBatchOptions{
			Operation: azcosmos.OperationOptions{EnableContentResponseOnWrite: &enabled},
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Len(t, response.OperationResults, 4)
		require.Equal(t, []int{201, 200, 200, 201}, []int{
			response.OperationResults[0].StatusCode,
			response.OperationResults[1].StatusCode,
			response.OperationResults[2].StatusCode,
			response.OperationResults[3].StatusCode,
		})
		var read map[string]any
		require.NoError(t, json.Unmarshal(response.OperationResults[1].ResourceBody, &read))
		require.Equal(t, "created", read["value"], "the batch read observes the preceding create")
		persisted, err := fx.Container.ReadItem(ctx, pk, createdID, nil)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(persisted.Value, &read))
		require.Equal(t, "replaced", read["value"])
		persisted, err = fx.Container.ReadItem(ctx, pk, replacedID, nil)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(persisted.Value, &read))
		require.Equal(t, "final", read["value"])

		existingID := uniqueID(t) + "-existing"
		rolledBackID := uniqueID(t) + "-rolled-back"
		trackItem(t, fx.Container, pk, existingID)
		trackItem(t, fx.Container, pk, rolledBackID)
		original, err := fx.Container.CreateItem(ctx, pk, existingID,
			map[string]any{"id": existingID, "pk": pkValue, "value": "original"}, nil)
		require.NoError(t, err)

		failed := fx.Container.NewTransactionalBatch(pk)
		require.NoError(t, failed.CreateItem(
			map[string]any{"id": rolledBackID, "pk": pkValue, "value": "must roll back"}, nil))
		require.NoError(t, failed.CreateItem(
			map[string]any{"id": existingID, "pk": pkValue, "value": "duplicate"}, nil))
		require.NoError(t, failed.ReplaceItem(existingID,
			map[string]any{"id": existingID, "pk": pkValue, "value": "must not replace"}, nil))

		rolledBack, err := fx.Container.ExecuteTransactionalBatch(ctx, failed, nil)
		require.NoError(t, err, "a received rollback is returned with per-operation statuses")
		require.Equal(t, http.StatusMultiStatus, rolledBack.StatusCode)
		require.Len(t, rolledBack.OperationResults, 3)
		require.Equal(t, []int{424, http.StatusConflict, 424}, []int{
			rolledBack.OperationResults[0].StatusCode,
			rolledBack.OperationResults[1].StatusCode,
			rolledBack.OperationResults[2].StatusCode,
		})

		stillOriginal, err := fx.Container.ReadItem(ctx, pk, existingID, nil)
		require.NoError(t, err)
		require.Equal(t, original.ETag, stillOriginal.ETag)
		require.NoError(t, json.Unmarshal(stillOriginal.Value, &read))
		require.Equal(t, "original", read["value"])
		_, err = fx.Container.ReadItem(ctx, pk, rolledBackID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)
	})
}

func TestItemPatchState(t *testing.T) {
	t.Run("item.patch-state", func(t *testing.T) {
		fx := newFixture(t, "item.patch-state", "")
		ctx := t.Context()
		id := uniqueID(t)
		pk := azcosmos.NewPartitionKeyString("patch-" + id)
		trackItem(t, fx.Container, pk, id)
		original := map[string]any{
			"id": id, "pk": "patch-" + id, "count": 1,
			"profile": map[string]any{"name": "before", "remove": true},
			"stable":  "untouched",
		}
		_, err := fx.Container.CreateItem(ctx, pk, id, original, nil)
		require.NoError(t, err)

		var safePatch azcosmos.PatchOperations
		require.NoError(t, safePatch.AppendIncrement("/count", 2))
		require.NoError(t, safePatch.AppendRemove("/profile/remove"))
		safeResponse, err := fx.Container.PatchItem(ctx, pk, id, safePatch, &azcosmos.PatchItemOptions{
			Strategy: azcosmos.PatchStrategyServerSide,
		})
		require.NoError(t, err)
		require.Equal(t, 200, safeResponse.StatusCode)
		safePersisted, err := fx.Container.ReadItem(ctx, pk, id, nil)
		require.NoError(t, err)
		require.JSONEq(t, string(safeResponse.Value), string(safePersisted.Value))
		var safeImage map[string]any
		mustUnmarshalInto(t, safeResponse.Value, &safeImage)
		require.Equal(t, float64(3), safeImage["count"])
		require.Equal(t, "before", safeImage["profile"].(map[string]any)["name"])
		require.NotContains(t, safeImage["profile"], "remove")
		require.Equal(t, "untouched", safeImage["stable"])

		var rmwPatch azcosmos.PatchOperations
		require.NoError(t, rmwPatch.AppendSet("/profile/name", "after"))
		require.NoError(t, rmwPatch.AppendSet("/stable", "updated"))
		rmwResponse, err := fx.Container.PatchItem(ctx, pk, id, rmwPatch, &azcosmos.PatchItemOptions{
			Strategy: azcosmos.PatchStrategyClientSide,
		})
		require.NoError(t, err)
		require.Equal(t, 200, rmwResponse.StatusCode)
		rmwPersisted, err := fx.Container.ReadItem(ctx, pk, id, nil)
		require.NoError(t, err)
		require.JSONEq(t, string(rmwResponse.Value), string(rmwPersisted.Value))
		var rmwImage map[string]any
		mustUnmarshalInto(t, rmwResponse.Value, &rmwImage)
		require.Equal(t, id, rmwImage["id"])
		require.Equal(t, "patch-"+id, rmwImage["pk"])
		require.Equal(t, float64(3), rmwImage["count"])
		require.Equal(t, "after", rmwImage["profile"].(map[string]any)["name"])
		require.NotContains(t, rmwImage["profile"], "remove")
		require.Equal(t, "updated", rmwImage["stable"])

		var invalidPatch azcosmos.PatchOperations
		newID := uniqueID(t)
		require.NoError(t, invalidPatch.AppendSet("/id", newID))
		_, err = fx.Container.PatchItem(ctx, pk, id, invalidPatch, nil)
		invalidID := requireCode(t, err, azcosmos.CodeClientError)
		require.Equal(t, http.StatusBadRequest, invalidID.StatusCode)
		afterInvalid, err := fx.Container.ReadItem(ctx, pk, id, nil)
		require.NoError(t, err)
		require.JSONEq(t, string(rmwPersisted.Value), string(afterInvalid.Value))
		_, err = fx.Container.ReadItem(ctx, pk, newID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)
	})
}

func TestItemValidationContracts(t *testing.T) {
	t.Run("item.validation-contracts", func(t *testing.T) {
		fx := newFixture(t, "item.validation-contracts", "")
		ctx := t.Context()
		pkValue := uniqueID(t)
		pk := azcosmos.NewPartitionKeyString(pkValue)

		missingID := uniqueID(t)
		_, err := fx.Container.CreateItem(ctx, pk, missingID,
			map[string]any{"pk": pkValue, "value": 1}, nil)
		require.Equal(t, http.StatusBadRequest, requireCode(t, err, azcosmos.CodeBadRequest).StatusCode)
		_, err = fx.Container.ReadItem(ctx, pk, missingID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)

		emptyID := uniqueID(t)
		_, err = fx.Container.CreateItem(ctx, pk, "",
			map[string]any{"id": emptyID, "pk": pkValue, "value": 1}, nil)
		require.ErrorContains(t, err, "item id must not be empty")
		_, err = fx.Container.ReadItem(ctx, pk, emptyID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)

		mismatchedID := uniqueID(t)
		mismatchedPK := uniqueID(t)
		_, err = fx.Container.CreateItem(ctx, pk, mismatchedID,
			map[string]any{"id": mismatchedID, "pk": mismatchedPK, "value": 1}, nil)
		require.Equal(t, http.StatusBadRequest, requireCode(t, err, azcosmos.CodeBadRequest).StatusCode)
		for _, key := range []azcosmos.PartitionKey{pk, azcosmos.NewPartitionKeyString(mismatchedPK)} {
			_, err = fx.Container.ReadItem(ctx, key, mismatchedID, nil)
			_ = requireCode(t, err, azcosmos.CodeNotFound)
		}

		oversizedID := uniqueID(t)
		oversized := json.RawMessage(fmt.Sprintf(
			`{"id":%q,"pk":%q,"padding":%q}`, oversizedID, pkValue, strings.Repeat("x", 2*1024*1024)))
		_, err = fx.Container.CreateItem(ctx, pk, oversizedID, oversized, nil)
		tooLarge := requireCode(t, err, azcosmos.CodeServiceError)
		require.Equal(t, http.StatusRequestEntityTooLarge, tooLarge.StatusCode)
		_, err = fx.Container.ReadItem(ctx, pk, oversizedID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)

		uniqueKeyContainer, _ := createTestContainer(t, fx.Database, azcosmos.ContainerProperties{
			PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{Paths: []string{"/pk"}},
			UniqueKeyPolicy: &azcosmos.UniqueKeyPolicy{
				UniqueKeys: []azcosmos.UniqueKey{{Paths: []string{"/email"}}},
			},
		})
		ownerID, duplicateID, otherPartitionID := uniqueID(t), uniqueID(t), uniqueID(t)
		ownerPK := azcosmos.NewPartitionKeyString("tenant-a")
		otherPK := azcosmos.NewPartitionKeyString("tenant-b")
		for _, pair := range []struct {
			pk azcosmos.PartitionKey
			id string
		}{{ownerPK, ownerID}, {ownerPK, duplicateID}, {otherPK, otherPartitionID}} {
			trackItem(t, uniqueKeyContainer, pair.pk, pair.id)
		}
		owner, err := uniqueKeyContainer.CreateItem(ctx, ownerPK, ownerID,
			map[string]any{"id": ownerID, "pk": "tenant-a", "email": "same@example.test", "value": "original"}, nil)
		require.NoError(t, err)
		_, err = uniqueKeyContainer.CreateItem(ctx, ownerPK, duplicateID,
			map[string]any{"id": duplicateID, "pk": "tenant-a", "email": "same@example.test", "value": "duplicate"}, nil)
		duplicate := requireCode(t, err, azcosmos.CodeConflict)
		require.Equal(t, http.StatusConflict, duplicate.StatusCode)
		persisted, err := uniqueKeyContainer.ReadItem(ctx, ownerPK, ownerID, nil)
		require.NoError(t, err)
		require.Equal(t, owner.ETag, persisted.ETag)
		var stored map[string]any
		mustUnmarshalInto(t, persisted.Value, &stored)
		require.Equal(t, "original", stored["value"])
		_, err = uniqueKeyContainer.ReadItem(ctx, ownerPK, duplicateID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)
		_, err = uniqueKeyContainer.CreateItem(ctx, otherPK, otherPartitionID,
			map[string]any{"id": otherPartitionID, "pk": "tenant-b", "email": "same@example.test"}, nil)
		require.NoError(t, err, "unique keys are scoped to a logical partition")

		oldID, newID := uniqueID(t), uniqueID(t)
		replacePK := azcosmos.NewPartitionKeyString("replace-" + oldID)
		trackItem(t, fx.Container, replacePK, oldID)
		trackItem(t, fx.Container, replacePK, newID)
		_, err = fx.Container.CreateItem(ctx, replacePK, oldID,
			map[string]any{"id": oldID, "pk": "replace-" + oldID, "value": "before"}, nil)
		require.NoError(t, err)
		beforeReplace, err := fx.Container.ReadItem(ctx, replacePK, oldID, nil)
		require.NoError(t, err)
		var beforeImage map[string]any
		mustUnmarshalInto(t, beforeReplace.Value, &beforeImage)
		replaced, err := fx.Container.ReplaceItem(ctx, replacePK, oldID,
			map[string]any{"id": newID, "pk": "replace-" + oldID, "value": "after"}, &azcosmos.ReplaceItemOptions{
				Operation: azcosmos.OperationOptions{EnableContentResponseOnWrite: to(true)},
			})
		require.NoError(t, err)
		require.Equal(t, 200, replaced.StatusCode)
		var replacement map[string]any
		mustUnmarshalInto(t, replaced.Value, &replacement)
		require.Equal(t, newID, replacement["id"])
		require.Equal(t, "after", replacement["value"])
		require.Equal(t, beforeImage["_rid"], replacement["_rid"],
			"name-addressed Replace retains the existing resource identity")
		afterReplace, err := fx.Container.ReadItem(ctx, replacePK, newID, nil)
		require.NoError(t, err)
		require.JSONEq(t, string(replaced.Value), string(afterReplace.Value))
		_, err = fx.Container.ReadItem(ctx, replacePK, oldID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)
	})
}

func TestItemNumericUniqueKeyEquivalence(t *testing.T) {
	t.Run("item.numeric-unique-key-equivalence", func(t *testing.T) {
		fx := newFixture(t, "item.numeric-unique-key-equivalence", "")
		container, _ := createTestContainer(t, fx.Database, azcosmos.ContainerProperties{
			PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{Paths: []string{"/pk"}},
			UniqueKeyPolicy: &azcosmos.UniqueKeyPolicy{
				UniqueKeys: []azcosmos.UniqueKey{{Paths: []string{"/number"}}},
			},
		})
		ctx := t.Context()
		pkValue := uniqueID(t)
		pk := azcosmos.NewPartitionKeyString(pkValue)
		firstID, equivalentID := uniqueID(t), uniqueID(t)
		trackItem(t, container, pk, firstID)
		trackItem(t, container, pk, equivalentID)
		first := json.RawMessage(fmt.Sprintf(
			`{"id":%q,"pk":%q,"number":9007199254740993}`, firstID, pkValue))
		equivalent := json.RawMessage(fmt.Sprintf(
			`{"id":%q,"pk":%q,"number":9007199254740992}`, equivalentID, pkValue))

		_, err := container.CreateItem(ctx, pk, firstID, first, nil)
		require.NoError(t, err)
		_, err = container.CreateItem(ctx, pk, equivalentID, equivalent, nil)
		duplicate := requireCode(t, err, azcosmos.CodeConflict)
		require.Equal(t, http.StatusConflict, duplicate.StatusCode)

		persisted, err := container.ReadItem(ctx, pk, firstID, nil)
		require.NoError(t, err)
		var stored map[string]any
		mustUnmarshalInto(t, persisted.Value, &stored)
		require.Equal(t, float64(9007199254740992), stored["number"])
		_, err = container.ReadItem(ctx, pk, equivalentID, nil)
		_ = requireCode(t, err, azcosmos.CodeNotFound)
	})
}
