// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryBuilder_IDPartitionKey_Simple(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/id"}}
	items := []indexedItem{
		{id: "a", pk: NewPartitionKeyString("a")},
		{id: "b", pk: NewPartitionKeyString("b")},
		{id: "c", pk: NewPartitionKeyString("c")},
	}

	require.True(t, qb.isIDPartitionKeyQuery(items, pkDef))

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE c.id IN (@param_id0, @param_id1, @param_id2)", query)
	require.Len(t, params, 3)
	require.Equal(t, "@param_id0", params[0].Name)
	require.Equal(t, "a", params[0].Value)
	require.Equal(t, "@param_id1", params[1].Name)
	require.Equal(t, "b", params[1].Value)
}

func TestQueryBuilder_IDPartitionKey_Mismatch(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/id"}}
	items := []indexedItem{
		{id: "a", pk: NewPartitionKeyString("a")},
		{id: "b", pk: NewPartitionKeyString("different")},
	}

	require.False(t, qb.isIDPartitionKeyQuery(items, pkDef))

	// Falls through to OR-of-conjunctions since PKs differ
	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Contains(t, query, "OR")
	require.Len(t, params, 4) // 2 ids + 2 pks
}

func TestQueryBuilder_SingleLogicalPartition(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/myPk"}}
	items := []indexedItem{
		{id: "1", pk: NewPartitionKeyString("samePK")},
		{id: "2", pk: NewPartitionKeyString("samePK")},
		{id: "3", pk: NewPartitionKeyString("samePK")},
		{id: "4", pk: NewPartitionKeyString("samePK")},
		{id: "5", pk: NewPartitionKeyString("samePK")},
	}

	require.True(t, qb.isSingleLogicalPartitionQuery(items))

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE c.myPk = @pk0 AND c.id IN (@param_id0, @param_id1, @param_id2, @param_id3, @param_id4)", query)
	require.Len(t, params, 6) // 1 pk + 5 ids
	require.Equal(t, "@pk0", params[0].Name)
	require.Equal(t, "samePK", params[0].Value)
}

func TestQueryBuilder_MultipleLogicalPartitions(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}
	items := []indexedItem{
		{id: "1", pk: NewPartitionKeyString("pkA")},
		{id: "2", pk: NewPartitionKeyString("pkB")},
		{id: "3", pk: NewPartitionKeyString("pkA")},
	}

	require.False(t, qb.isIDPartitionKeyQuery(items, pkDef))
	require.False(t, qb.isSingleLogicalPartitionQuery(items))

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Contains(t, query, "OR")
	require.Contains(t, query, "c.id = @param_id0")
	require.Contains(t, query, "c.pk = @param_pk00")
	require.Len(t, params, 6) // 3 ids + 3 pks
}

func TestQueryBuilder_NestedPKPath(t *testing.T) {
	qb := queryBuilder{}
	expr := qb.getFieldExpression("/address/zipCode")
	require.Equal(t, `c["address"]["zipCode"]`, expr)
}

func TestQueryBuilder_NonIdentifierPKPath(t *testing.T) {
	qb := queryBuilder{}
	expr := qb.getFieldExpression("/my-pk")
	require.Equal(t, `c["my-pk"]`, expr)
}

func TestQueryBuilder_NullPartitionKey(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}
	items := []indexedItem{
		{id: "1", pk: NullPartitionKey},
		{id: "2", pk: NullPartitionKey},
	}

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	// Single logical partition (both null), uses PK+ID IN shape
	require.Contains(t, query, "IS_DEFINED(c.pk) = false")
	require.Contains(t, query, "c.id IN")
	// Only id params, no pk params for null
	require.Len(t, params, 2)
}

func TestQueryBuilder_HierarchicalPK(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{
		Paths: []string{"/tenantId", "/userId"},
		Kind:  PartitionKeyKindMultiHash,
	}
	items := []indexedItem{
		{id: "1", pk: NewPartitionKeyString("t1").AppendString("u1")},
		{id: "2", pk: NewPartitionKeyString("t1").AppendString("u2")},
	}

	// Different PKs → OR-of-conjunctions
	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Contains(t, query, "OR")
	require.Contains(t, query, "c.tenantId = @param_pk00")
	require.Contains(t, query, "c.userId = @param_pk01")
	require.Contains(t, query, "c.tenantId = @param_pk10")
	require.Contains(t, query, "c.userId = @param_pk11")
	// 2 ids + 4 pk components
	require.Len(t, params, 6)
}

func TestQueryBuilder_EmptyItems(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}

	require.False(t, qb.isIDPartitionKeyQuery(nil, pkDef))
	require.False(t, qb.isSingleLogicalPartitionQuery(nil))
}

func TestQueryBuilder_SingleItem(t *testing.T) {
	qb := queryBuilder{}
	items := []indexedItem{
		{id: "1", pk: NewPartitionKeyString("pk1")},
	}

	// Single item → isSingleLogicalPartitionQuery returns false
	require.False(t, qb.isSingleLogicalPartitionQuery(items))

	// buildParameterizedQueryForItems still works (OR-of-conjunctions with 1 item)
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}
	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Contains(t, query, "c.id = @param_id0")
	require.Len(t, params, 2)
}

func TestQueryBuilder_GetFieldExpression(t *testing.T) {
	qb := queryBuilder{}

	// Simple path
	require.Equal(t, "c.pk", qb.getFieldExpression("/pk"))

	// Nested path
	require.Equal(t, `c["a"]["b"]["c"]`, qb.getFieldExpression("/a/b/c"))

	// Non-identifier
	require.Equal(t, `c["non-ident"]`, qb.getFieldExpression("/non-ident"))

	// Underscore-prefixed (valid identifier)
	require.Equal(t, "c._pk", qb.getFieldExpression("/_pk"))

	// id path
	require.Equal(t, "c.id", qb.getFieldExpression("/id"))
}
