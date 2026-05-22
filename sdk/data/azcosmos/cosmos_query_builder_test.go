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
	items := []ItemIdentity{
		{ID: "a", PartitionKey: NewPartitionKeyString("a")},
		{ID: "b", PartitionKey: NewPartitionKeyString("b")},
		{ID: "c", PartitionKey: NewPartitionKeyString("c")},
	}

	require.True(t, qb.isIDPartitionKeyQuery(items, pkDef))

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE c.id IN (@param_id0, @param_id1, @param_id2)", query)
	require.Len(t, params, 3)
	require.Equal(t, QueryParameter{Name: "@param_id0", Value: "a"}, params[0])
	require.Equal(t, QueryParameter{Name: "@param_id1", Value: "b"}, params[1])
	require.Equal(t, QueryParameter{Name: "@param_id2", Value: "c"}, params[2])
}

func TestQueryBuilder_IDPartitionKey_Mismatch(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/id"}}
	items := []ItemIdentity{
		{ID: "a", PartitionKey: NewPartitionKeyString("a")},
		{ID: "b", PartitionKey: NewPartitionKeyString("different")},
	}

	require.False(t, qb.isIDPartitionKeyQuery(items, pkDef))

	// Falls through to OR-of-conjunctions since PKs differ
	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE (c.id = @param_id0 AND c.id = @param_pk00) OR (c.id = @param_id1 AND c.id = @param_pk10)", query)
	require.Len(t, params, 4)
	require.Equal(t, QueryParameter{Name: "@param_id0", Value: "a"}, params[0])
	require.Equal(t, QueryParameter{Name: "@param_pk00", Value: "a"}, params[1])
	require.Equal(t, QueryParameter{Name: "@param_id1", Value: "b"}, params[2])
	require.Equal(t, QueryParameter{Name: "@param_pk10", Value: "different"}, params[3])
}

func TestQueryBuilder_SingleLogicalPartition(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/myPk"}}
	items := []ItemIdentity{
		{ID: "1", PartitionKey: NewPartitionKeyString("samePK")},
		{ID: "2", PartitionKey: NewPartitionKeyString("samePK")},
		{ID: "3", PartitionKey: NewPartitionKeyString("samePK")},
		{ID: "4", PartitionKey: NewPartitionKeyString("samePK")},
		{ID: "5", PartitionKey: NewPartitionKeyString("samePK")},
	}

	require.True(t, qb.isSingleLogicalPartitionQuery(items))

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE c.myPk = @pk0 AND c.id IN (@param_id0, @param_id1, @param_id2, @param_id3, @param_id4)", query)
	require.Len(t, params, 6)
	require.Equal(t, QueryParameter{Name: "@pk0", Value: "samePK"}, params[0])
	require.Equal(t, QueryParameter{Name: "@param_id0", Value: "1"}, params[1])
	require.Equal(t, QueryParameter{Name: "@param_id1", Value: "2"}, params[2])
	require.Equal(t, QueryParameter{Name: "@param_id2", Value: "3"}, params[3])
	require.Equal(t, QueryParameter{Name: "@param_id3", Value: "4"}, params[4])
	require.Equal(t, QueryParameter{Name: "@param_id4", Value: "5"}, params[5])
}

func TestQueryBuilder_MultipleLogicalPartitions(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}
	items := []ItemIdentity{
		{ID: "1", PartitionKey: NewPartitionKeyString("pkA")},
		{ID: "2", PartitionKey: NewPartitionKeyString("pkB")},
		{ID: "3", PartitionKey: NewPartitionKeyString("pkA")},
	}

	require.False(t, qb.isIDPartitionKeyQuery(items, pkDef))
	require.False(t, qb.isSingleLogicalPartitionQuery(items))

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE (c.id = @param_id0 AND c.pk = @param_pk00) OR (c.id = @param_id1 AND c.pk = @param_pk10) OR (c.id = @param_id2 AND c.pk = @param_pk20)", query)
	require.Len(t, params, 6)
	require.Equal(t, QueryParameter{Name: "@param_id0", Value: "1"}, params[0])
	require.Equal(t, QueryParameter{Name: "@param_pk00", Value: "pkA"}, params[1])
	require.Equal(t, QueryParameter{Name: "@param_id1", Value: "2"}, params[2])
	require.Equal(t, QueryParameter{Name: "@param_pk10", Value: "pkB"}, params[3])
	require.Equal(t, QueryParameter{Name: "@param_id2", Value: "3"}, params[4])
	require.Equal(t, QueryParameter{Name: "@param_pk20", Value: "pkA"}, params[5])
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
	items := []ItemIdentity{
		{ID: "1", PartitionKey: NullPartitionKey},
		{ID: "2", PartitionKey: NullPartitionKey},
	}

	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	// Single logical partition (both null), uses PK+ID IN shape
	require.Equal(t, "SELECT * FROM c WHERE IS_NULL(c.pk) AND c.id IN (@param_id0, @param_id1)", query)
	require.Len(t, params, 2)
	require.Equal(t, QueryParameter{Name: "@param_id0", Value: "1"}, params[0])
	require.Equal(t, QueryParameter{Name: "@param_id1", Value: "2"}, params[1])
}

func TestQueryBuilder_HierarchicalPK(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{
		Paths: []string{"/tenantId", "/userId"},
		Kind:  PartitionKeyKindMultiHash,
	}
	items := []ItemIdentity{
		{ID: "1", PartitionKey: NewPartitionKeyString("t1").AppendString("u1")},
		{ID: "2", PartitionKey: NewPartitionKeyString("t1").AppendString("u2")},
	}

	// Different PKs → OR-of-conjunctions
	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE (c.id = @param_id0 AND c.tenantId = @param_pk00 AND c.userId = @param_pk01) OR (c.id = @param_id1 AND c.tenantId = @param_pk10 AND c.userId = @param_pk11)", query)
	require.Len(t, params, 6)
	require.Equal(t, QueryParameter{Name: "@param_id0", Value: "1"}, params[0])
	require.Equal(t, QueryParameter{Name: "@param_pk00", Value: "t1"}, params[1])
	require.Equal(t, QueryParameter{Name: "@param_pk01", Value: "u1"}, params[2])
	require.Equal(t, QueryParameter{Name: "@param_id1", Value: "2"}, params[3])
	require.Equal(t, QueryParameter{Name: "@param_pk10", Value: "t1"}, params[4])
	require.Equal(t, QueryParameter{Name: "@param_pk11", Value: "u2"}, params[5])
}

func TestQueryBuilder_EmptyItems(t *testing.T) {
	qb := queryBuilder{}
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}

	require.False(t, qb.isIDPartitionKeyQuery(nil, pkDef))
	require.False(t, qb.isSingleLogicalPartitionQuery(nil))
}

func TestQueryBuilder_SingleItem(t *testing.T) {
	qb := queryBuilder{}
	items := []ItemIdentity{
		{ID: "1", PartitionKey: NewPartitionKeyString("pk1")},
	}

	// Single item → isSingleLogicalPartitionQuery returns false
	require.False(t, qb.isSingleLogicalPartitionQuery(items))

	// buildParameterizedQueryForItems still works (OR-of-conjunctions with 1 item)
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}
	query, params := qb.buildParameterizedQueryForItems(items, pkDef)
	require.Equal(t, "SELECT * FROM c WHERE (c.id = @param_id0 AND c.pk = @param_pk00)", query)
	require.Len(t, params, 2)
	require.Equal(t, QueryParameter{Name: "@param_id0", Value: "1"}, params[0])
	require.Equal(t, QueryParameter{Name: "@param_pk00", Value: "pk1"}, params[1])
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
