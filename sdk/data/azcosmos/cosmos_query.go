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

type queryBody struct {
	Query      string           `json:"query"`
	Parameters []QueryParameter `json:"parameters,omitempty"`
}
