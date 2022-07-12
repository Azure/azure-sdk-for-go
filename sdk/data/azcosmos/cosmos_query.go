// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// QueryParameter represents a parameter for a parametrized query.
type QueryParameter struct {
	// Name represents the name of the parameter in the parametrized query.
	Name string `json:"name"`
	// Value represents the value of the parameter in the parametrized query.
	Value any `json:"value"`
}

// NewQueryParameter creates a new parameter for a parametrized query.
// name - The name of the parameter.
// value - The value of the parameter.
// See https://docs.microsoft.com/azure/cosmos-db/sql/sql-query-parameterized-queries
func NewQueryParameter(name string, value any) QueryParameter {
	return QueryParameter{name, value}
}

type queryBody struct {
	Query      string           `json:"query"`
	Parameters []QueryParameter `json:"parameters,omitempty"`
}
