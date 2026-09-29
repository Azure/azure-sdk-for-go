// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewQueryItemsPager queries a complete logical partition, returning raw JSON items in pages.
// Construction performs no network I/O; argument errors are reported by NextPage.
//
// The first fetch reads container metadata to reject incomplete partition keys. Its charge is
// included in that fetch's response or error. Resuming with a new pager repeats this validation.
// Empty pages may have a continuation. Pagers are not safe for concurrent use.
func (c *ContainerClient) NewQueryItemsPager(query Query, scope FeedScope, options *QueryOptions) *runtime.Pager[QueryItemsResponse] {
	req, err := newQueryRequest(query, scope, options)
	req.databaseID = c.database.id
	req.containerID = c.id
	return newQueryItemsPager(req, err, c.queryItems)
}

type queryRequest struct {
	databaseID     string
	containerID    string
	body           []byte
	partitionKey   PartitionKey
	options        QueryOptions
	scopeValidated bool
}

func newQueryRequest(query Query, scope FeedScope, options *QueryOptions) (queryRequest, error) {
	req := queryRequest{partitionKey: scope.partitionKey}
	var err error
	req.body, err = query.body()
	if err != nil {
		return req, err
	}
	if err := validateItemArguments(scope.partitionKey); err != nil {
		return req, err
	}
	if options != nil {
		req.options = *options
		req.options.Operation.ExcludedRegions = slices.Clone(options.Operation.ExcludedRegions)
		if value := options.Operation.EnableContentResponseOnWrite; value != nil {
			copied := *value
			req.options.Operation.EnableContentResponseOnWrite = &copied
		}
	}
	if err := req.options.Operation.ConsistencyStrategy.validate(); err != nil {
		return req, err
	}
	if err := req.options.SessionToken.validate(); err != nil {
		return req, err
	}
	if req.options.Feed.PageSizeHint < 0 {
		return req, errors.New("azcosmos: page size hint must not be negative")
	}
	if strings.IndexByte(req.options.Feed.ContinuationToken, 0) >= 0 {
		return req, errors.New("azcosmos: continuation token must not contain a NUL byte")
	}
	return req, nil
}

func newQueryItemsPager(req queryRequest, validationErr error, fetch func(context.Context, *queryRequest) (QueryItemsResponse, error)) *runtime.Pager[QueryItemsResponse] {
	return runtime.NewPager(runtime.PagingHandler[QueryItemsResponse]{
		More: func(page QueryItemsResponse) bool { return page.ContinuationToken != "" },
		Fetcher: func(ctx context.Context, previous *QueryItemsResponse) (QueryItemsResponse, error) {
			if validationErr != nil {
				return QueryItemsResponse{}, validationErr
			}
			if previous != nil {
				req.options.Feed.ContinuationToken = previous.ContinuationToken
			}
			return fetch(ctx, &req)
		},
	})
}

func (c *ContainerClient) queryItems(ctx context.Context, req *queryRequest) (QueryItemsResponse, error) {
	client := c.database.client
	release, err := client.acquire()
	if err != nil {
		return QueryItemsResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return QueryItemsResponse{}, err
	}
	return client.executeQuery(ctx, req)
}

func decodeQueryPage(body []byte, response Response, sessionToken SessionToken, continuation string, exhausted bool) (QueryItemsResponse, error) {
	page := QueryItemsResponse{Response: response, SessionToken: sessionToken, ContinuationToken: continuation}
	if exhausted {
		if len(body) != 0 || continuation != "" {
			return QueryItemsResponse{}, queryResponseError(response, sessionToken, errors.New("inconsistent exhausted query completion"))
		}
		return page, nil
	}
	var envelope struct {
		Documents []json.RawMessage `json:"Documents"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return QueryItemsResponse{}, queryResponseError(response, sessionToken, err)
	}
	if envelope.Documents == nil {
		return QueryItemsResponse{}, queryResponseError(response, sessionToken, errors.New("query response has no Documents array"))
	}
	page.Items = make([][]byte, len(envelope.Documents))
	for i, item := range envelope.Documents {
		page.Items[i] = item
	}
	return page, nil
}

func queryResponseError(response Response, sessionToken SessionToken, cause error) *Error {
	return &Error{
		Code: CodeSerializationFailed, Message: "decoding query response",
		RequestCharge: response.RequestCharge, ActivityID: response.ActivityID,
		SessionToken: sessionToken, cause: cause,
	}
}

func validateQueryPartitionKey(body []byte, partitionKey PartitionKey) error {
	var metadata struct {
		PartitionKey struct {
			Paths   []string `json:"paths"`
			Kind    string   `json:"kind"`
			Version *int     `json:"version"`
		} `json:"partitionKey"`
	}
	if err := json.Unmarshal(body, &metadata); err != nil {
		return &Error{Code: CodeSerializationFailed, Message: "decoding container partition key definition", cause: err}
	}
	definition := metadata.PartitionKey
	if len(definition.Paths) == 0 {
		return &Error{Code: CodeSerializationFailed, Message: "container has no partition key paths"}
	}
	if definition.Kind != "Hash" && definition.Kind != "MultiHash" {
		return &Error{Code: CodeBadRequest, Message: "query scope requires a Hash or MultiHash partition key definition"}
	}
	if (definition.Kind == "Hash" && len(definition.Paths) != 1) ||
		(definition.Kind == "MultiHash" && (definition.Version == nil || *definition.Version != 2)) ||
		(definition.Version != nil && *definition.Version != 1 && *definition.Version != 2) {
		return &Error{Code: CodeBadRequest, Message: "unsupported query partition key definition"}
	}
	if partitionKey.Len() != len(definition.Paths) {
		return &Error{Code: CodeBadRequest, Message: fmt.Sprintf(
			"query scope requires all %d partition key components; got %d (prefix queries are not supported)",
			len(definition.Paths), partitionKey.Len())}
	}
	return nil
}

func addQuerySetupCharge(err error, setup Response) error {
	if setup.RequestCharge == 0 {
		return err
	}
	var cosmosErr *Error
	if errors.As(err, &cosmosErr) {
		copied := *cosmosErr
		copied.RequestCharge += setup.RequestCharge
		if copied.ActivityID == "" {
			copied.ActivityID = setup.ActivityID
		}
		return &copied
	}
	code := CodeClientError
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		code = CodeOperationCancelled
	}
	return &Error{Code: code, Message: "fetching query page", RequestCharge: setup.RequestCharge,
		ActivityID: setup.ActivityID, cause: err}
}
