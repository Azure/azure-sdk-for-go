// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"strings"
)

// Query is SQL query text with optional parameters. Values are immutable and can be shared.
// The zero value is invalid; use [NewQuery] to supply query text.
type Query struct {
	text       string
	parameters []queryParameter
}

type queryParameter struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

// NewQuery creates a query. Empty or whitespace-only text is rejected when fetching a page.
func NewQuery(text string) Query {
	return Query{text: text}
}

// WithParameter returns a new query with a JSON snapshot of value, replacing any parameter
// with the same name. Names must start with @; values must be JSON-serializable.
func (q Query) WithParameter(name string, value any) (Query, error) {
	if len(name) < 2 || name[0] != '@' || strings.ContainsAny(name, " \t\r\n\x00") {
		return Query{}, errors.New("azcosmos: query parameter name must start with @ and contain no whitespace or NUL")
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return Query{}, &Error{Code: CodeSerializationFailed, Message: "encoding query parameter", cause: err}
	}
	parameters := append([]queryParameter(nil), q.parameters...)
	for i := range parameters {
		if parameters[i].Name == name {
			parameters[i].Value = encoded
			q.parameters = parameters
			return q, nil
		}
	}
	parameters = append(parameters, queryParameter{Name: name, Value: encoded})
	q.parameters = parameters
	return q, nil
}

func (q Query) body() ([]byte, error) {
	if strings.TrimSpace(q.text) == "" {
		return nil, errors.New("azcosmos: query text must not be empty")
	}
	return json.Marshal(struct {
		Query      string           `json:"query"`
		Parameters []queryParameter `json:"parameters,omitempty"`
	}{q.text, q.parameters})
}

// FeedScope identifies the partition targeted by a query. The zero value is invalid.
// Only complete logical partition keys are supported; cross-partition and prefix scopes are not.
type FeedScope struct {
	partitionKey PartitionKey
}

// NewFeedScopeForPartitionKey targets one logical partition. Hierarchical keys must include
// every component; partial keys are rejected on the first page fetch.
func NewFeedScopeForPartitionKey(partitionKey PartitionKey) FeedScope {
	return FeedScope{partitionKey: partitionKey}
}

// FeedOptions configures pagination and resumption.
type FeedOptions struct {
	// PageSizeHint is a positive maximum-item-count hint, not a guaranteed page size.
	// Zero uses the driver default. Negative values are invalid.
	PageSizeHint int32

	// ContinuationToken resumes a previous query with the same query and scope.
	// Empty starts a new query. Tokens are opaque and must not be modified.
	ContinuationToken string
}

// QueryOptions configures item queries. A nil *QueryOptions selects defaults.
type QueryOptions struct {
	// Operation holds the settings every operation accepts. Its timeout applies to each page fetch.
	Operation OperationOptions

	// Feed controls page size and resumption.
	Feed FeedOptions

	// SessionToken overrides the client's captured session token. Empty uses the client's token.
	SessionToken SessionToken
}
