// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"strings"
	"unicode"
)

// queryBuilder builds parameterized SQL queries for read-many operations.
// It chooses one of three query shapes depending on the items in the batch:
//   - ID-only IN: when the PK path is /id and every PK value equals the item id
//   - PK + ID IN: when all items share the same logical partition key value
//   - OR-of-conjunctions: general case with multiple logical PKs
type queryBuilder struct{}

// isIDPartitionKeyQuery returns true when the partition key path is "/id" and
// every item's PK value equals its item ID. In that case, only an id IN (…)
// clause is needed.
func (qb queryBuilder) isIDPartitionKeyQuery(items []ItemIdentity, pkDef PartitionKeyDefinition) bool {
	if len(items) == 0 {
		return false
	}
	if len(pkDef.Paths) != 1 || pkDef.Paths[0] != "/id" {
		return false
	}
	for i := range items {
		if len(items[i].PartitionKey.values) != 1 {
			return false
		}
		pkStr, ok := items[i].PartitionKey.values[0].(string)
		if !ok || pkStr != items[i].ID {
			return false
		}
	}
	return true
}

// isSingleLogicalPartitionQuery returns true when all items share the same
// logical partition key value. The check compares JSON-serialised PK strings.
func (qb queryBuilder) isSingleLogicalPartitionQuery(items []ItemIdentity) bool {
	if len(items) <= 1 {
		return false
	}
	first, err := items[0].PartitionKey.toJsonString()
	if err != nil {
		return false
	}
	for i := 1; i < len(items); i++ {
		s, err := items[i].PartitionKey.toJsonString()
		if err != nil || s != first {
			return false
		}
	}
	return true
}

// buildIDInQuery builds: SELECT * FROM c WHERE c.id IN (@param_id0, @param_id1, …)
func (qb queryBuilder) buildIDInQuery(items []ItemIdentity) (string, []QueryParameter) {
	params := make([]QueryParameter, 0, len(items))
	placeholders := make([]string, 0, len(items))
	for i, item := range items {
		name := fmt.Sprintf("@param_id%d", i)
		params = append(params, QueryParameter{Name: name, Value: item.ID})
		placeholders = append(placeholders, name)
	}
	query := fmt.Sprintf("SELECT * FROM c WHERE c.id IN (%s)", strings.Join(placeholders, ", "))
	return query, params
}

// buildPKAndIDInQuery builds:
//
//	SELECT * FROM c WHERE <pkExpr> = @pk AND c.id IN (@id0, @id1, …)
func (qb queryBuilder) buildPKAndIDInQuery(items []ItemIdentity, pkDef PartitionKeyDefinition) (string, []QueryParameter) {
	params := make([]QueryParameter, 0, len(items)+len(pkDef.Paths))

	// Build PK equality conditions (one per path for hierarchical PKs)
	pkConditions := make([]string, 0, len(pkDef.Paths))
	for pathIdx, path := range pkDef.Paths {
		fieldExpr := qb.getFieldExpression(path)
		paramName := fmt.Sprintf("@pk%d", pathIdx)

		var pkVal interface{}
		if pathIdx < len(items[0].PartitionKey.values) {
			pkVal = items[0].PartitionKey.values[pathIdx]
		}

		if pkVal == nil {
			// null PK → IS_DEFINED check
			pkConditions = append(pkConditions, fmt.Sprintf("IS_DEFINED(%s) = false", fieldExpr))
		} else {
			params = append(params, QueryParameter{Name: paramName, Value: pkVal})
			pkConditions = append(pkConditions, fmt.Sprintf("%s = %s", fieldExpr, paramName))
		}
	}

	// Build id IN clause
	idPlaceholders := make([]string, 0, len(items))
	for i, item := range items {
		name := fmt.Sprintf("@param_id%d", i)
		params = append(params, QueryParameter{Name: name, Value: item.ID})
		idPlaceholders = append(idPlaceholders, name)
	}

	query := fmt.Sprintf("SELECT * FROM c WHERE %s AND c.id IN (%s)",
		strings.Join(pkConditions, " AND "),
		strings.Join(idPlaceholders, ", "))
	return query, params
}

// buildOrOfConjunctionsQuery builds:
//
//	SELECT * FROM c WHERE (c.id = @param_id0 AND <pk0> = @param_pk00)
//	  OR (c.id = @param_id1 AND <pk1> = @param_pk10) …
func (qb queryBuilder) buildOrOfConjunctionsQuery(items []ItemIdentity, pkDef PartitionKeyDefinition) (string, []QueryParameter) {
	params := make([]QueryParameter, 0, len(items)*(1+len(pkDef.Paths)))
	clauses := make([]string, 0, len(items))

	for i, item := range items {
		idParam := fmt.Sprintf("@param_id%d", i)
		params = append(params, QueryParameter{Name: idParam, Value: item.ID})

		conditions := []string{fmt.Sprintf("c.id = %s", idParam)}

		for pathIdx, path := range pkDef.Paths {
			fieldExpr := qb.getFieldExpression(path)

			var pkVal interface{}
			if pathIdx < len(item.PartitionKey.values) {
				pkVal = item.PartitionKey.values[pathIdx]
			}

			if pkVal == nil {
				conditions = append(conditions, fmt.Sprintf("IS_DEFINED(%s) = false", fieldExpr))
			} else {
				pkParam := fmt.Sprintf("@param_pk%d%d", i, pathIdx)
				params = append(params, QueryParameter{Name: pkParam, Value: pkVal})
				conditions = append(conditions, fmt.Sprintf("%s = %s", fieldExpr, pkParam))
			}
		}

		clauses = append(clauses, fmt.Sprintf("(%s)", strings.Join(conditions, " AND ")))
	}

	query := fmt.Sprintf("SELECT * FROM c WHERE %s", strings.Join(clauses, " OR "))
	return query, params
}

// buildParameterizedQueryForItems selects the appropriate query shape and
// returns the query string and parameters.
func (qb queryBuilder) buildParameterizedQueryForItems(items []ItemIdentity, pkDef PartitionKeyDefinition) (string, []QueryParameter) {
	if qb.isIDPartitionKeyQuery(items, pkDef) {
		return qb.buildIDInQuery(items)
	}
	if qb.isSingleLogicalPartitionQuery(items) {
		return qb.buildPKAndIDInQuery(items, pkDef)
	}
	return qb.buildOrOfConjunctionsQuery(items, pkDef)
}

// getFieldExpression converts a partition key path like "/pk" or "/address/zipCode"
// into a Cosmos SQL field expression.
// Simple identifiers (letters, digits, underscore) → c.pk
// Non-identifiers or nested paths → c["my-pk"] or c["address"]["zipCode"]
func (qb queryBuilder) getFieldExpression(path string) string {
	// Strip leading "/"
	trimmed := strings.TrimPrefix(path, "/")
	if trimmed == "" {
		return "c"
	}

	parts := strings.Split(trimmed, "/")

	// If single part and is a simple identifier, use dot notation
	if len(parts) == 1 && isSimpleIdentifier(parts[0]) {
		return "c." + parts[0]
	}

	// Use bracket notation for all parts
	var b strings.Builder
	b.WriteString("c")
	for _, p := range parts {
		fmt.Fprintf(&b, "[\"%s\"]", p)
	}
	return b.String()
}

// isSimpleIdentifier returns true if s contains only letters, digits, and underscores
// and starts with a letter or underscore.
func isSimpleIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}
