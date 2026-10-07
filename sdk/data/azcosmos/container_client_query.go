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
	"sync"
)

// NewQueryItemsPager queries the supplied scope, returning raw JSON items in pages.
// Construction performs no network I/O; argument errors are reported by NextPage.
//
// Defer Close when stopping before exhaustion. Empty pages do not imply exhaustion.
// Use ContinuationToken to snapshot progress explicitly; not every query supports snapshots.
func (c *ContainerClient) NewQueryItemsPager(query Query, scope FeedScope, options *QueryOptions) *QueryItemsPager {
	req, err := newQueryRequest(query, scope, options)
	req.databaseID = c.database.id
	req.containerID = c.id
	return &QueryItemsPager{client: c.database.client, req: req, validationErr: err}
}

// QueryItemsPager owns a retained query plan. Close releases it when iteration stops early.
// Iteration and checkpoint calls must not be concurrent; Close safely waits for an active call.
// An empty token is not an exhaustion signal. Use More instead.
type QueryItemsPager struct {
	mu            sync.Mutex
	client        *Client
	req           queryRequest
	validationErr error
	cursor        queryCursor
	done          bool
}

type queryCursor interface {
	next(context.Context) (QueryItemsResponse, bool, error)
	checkpoint(context.Context) (string, error)
	close()
}

// More reports whether another page may be fetched. It does not perform I/O.
func (p *QueryItemsPager) More() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return !p.done
}

// NextPage fetches one page. A final empty page can signal exhaustion through More.
// Failures after execution starts terminate this pager; retries belong to the native driver.
// The context bounds the Go wait, not admitted native execution. Client.Close may wait
// for native work to finish after this method returns a context error.
func (p *QueryItemsPager) NextPage(ctx context.Context) (QueryItemsResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.validationErr != nil {
		return QueryItemsResponse{}, p.validationErr
	}
	release, err := p.client.acquire()
	if err != nil {
		return QueryItemsResponse{}, err
	}
	defer release()
	if p.done {
		return QueryItemsResponse{}, &Error{Code: CodeClientError, Message: "azcosmos: query pager is closed or exhausted"}
	}
	if err := ctx.Err(); err != nil {
		return QueryItemsResponse{}, err
	}
	ctx, cancel := contextWithEndToEndTimeout(ctx, p.req.options.Operation.EndToEndTimeout)
	defer cancel()
	var setup Response
	if p.cursor == nil {
		p.cursor, setup, err = p.client.openQuery(ctx, &p.req)
		if err != nil {
			p.finish()
			return QueryItemsResponse{}, addQuerySetupCharge(err, setup)
		}
	}
	page, end, err := p.cursor.next(ctx)
	if err != nil {
		p.finish()
		return QueryItemsResponse{}, addQuerySetupCharge(err, setup)
	}
	page.RequestCharge += setup.RequestCharge
	if page.ActivityID == "" && setup.RequestCharge != 0 {
		page.ActivityID = setup.ActivityID
	}
	if end {
		p.finish()
	}
	return page, nil
}

// ContinuationToken snapshots delivered progress without advancing the query.
// Call after a successful NextPage and before exhaustion or Close. Unsupported snapshots
// return an error but do not prevent further paging. Resume with the same query and scope.
// The context bounds the Go wait; it does not cancel admitted native checkpoint work.
func (p *QueryItemsPager) ContinuationToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.validationErr != nil {
		return "", p.validationErr
	}
	release, err := p.client.acquire()
	if err != nil {
		return "", err
	}
	defer release()
	if p.done || p.cursor == nil {
		return "", &Error{Code: CodeClientError, Message: "azcosmos: checkpoint requires an open query pager"}
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	ctx, cancel := contextWithEndToEndTimeout(ctx, p.req.options.Operation.EndToEndTimeout)
	defer cancel()
	token, err := p.cursor.checkpoint(ctx)
	if err != nil && !unsupportedQueryCheckpoint(err) {
		p.finish()
	}
	return token, err
}

func unsupportedQueryCheckpoint(err error) bool {
	var cosmosErr *Error
	return errors.As(err, &cosmosErr) && !cosmosErr.FromWire &&
		cosmosErr.StatusCode == 400 &&
		(cosmosErr.SubStatus == subStatusBufferedQueryContinuationUnsupported ||
			cosmosErr.SubStatus == subStatusContinuationTokenNonQueryOperation)
}

// Close releases the retained query plan. It is idempotent and safe alongside Client.Close.
func (p *QueryItemsPager) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.finish()
	return nil
}

func (p *QueryItemsPager) finish() {
	p.done = true
	if p.cursor != nil {
		p.cursor.close()
		p.cursor = nil
	}
}

type queryRequest struct {
	databaseID    string
	containerID   string
	body          []byte
	partitionKey  PartitionKey
	fullContainer bool
	options       QueryOptions
}

func newQueryRequest(query Query, scope FeedScope, options *QueryOptions) (queryRequest, error) {
	req := queryRequest{partitionKey: scope.partitionKey, fullContainer: scope.fullContainer}
	var err error
	req.body, err = query.body()
	if err != nil {
		return req, err
	}
	if !scope.fullContainer {
		if err := validateItemArguments(scope.partitionKey); err != nil {
			return req, err
		}
	}
	if options != nil {
		req.options = *options
		req.options.Operation.ExcludedRegions = slices.Clone(options.Operation.ExcludedRegions)
		if value := options.Operation.EnableContentResponseOnWrite; value != nil {
			copied := *value
			req.options.Operation.EnableContentResponseOnWrite = &copied
		}
		if value := options.PopulateIndexMetrics; value != nil {
			copied := *value
			req.options.PopulateIndexMetrics = &copied
		}
		if value := options.PopulateQueryMetrics; value != nil {
			copied := *value
			req.options.PopulateQueryMetrics = &copied
		}
	}
	switch req.options.QueryPlanMode {
	case "", QueryPlanModeLocalPreferred, QueryPlanModeGatewayOnly:
	default:
		return req, errors.New("azcosmos: invalid query plan mode")
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

func decodeQueryPage(body []byte, response Response, sessionToken SessionToken, continuation string, exhausted bool) (QueryItemsResponse, error) {
	page := QueryItemsResponse{Response: response, SessionToken: sessionToken}
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
		Diagnostics:   response.Diagnostics,
		StatusCode:    response.StatusCode,
		SubStatus:     response.SubStatus,
		AttemptCount:  response.AttemptCount,
		RequestCharge: response.RequestCharge,
		ActivityID:    response.ActivityID,
		SessionToken:  sessionToken,
		cause:         cause,
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
	if partitionKey.Len() == 0 || partitionKey.Len() > len(definition.Paths) {
		return &Error{Code: CodeBadRequest, Message: fmt.Sprintf(
			"query scope requires between 1 and %d partition key components; got %d",
			len(definition.Paths), partitionKey.Len())}
	}
	return nil
}

func addQuerySetupCharge(err error, setup Response) error {
	// RequestCharge alone is not a reliable "no setup" sentinel: a successful setup completion can
	// carry diagnostics/status/attempt metadata even when the charge header is absent or zero.
	// Only treat the zero Response itself as "nothing to add".
	if setup == (Response{}) {
		return err
	}
	var cosmosErr *Error
	if errors.As(err, &cosmosErr) {
		copied := *cosmosErr
		copied.RequestCharge += setup.RequestCharge
		if copied.ActivityID == "" {
			copied.ActivityID = setup.ActivityID
		}
		// err may be derived from the setup fetch itself, such as validateQueryPartitionKey's,
		// which has no response to copy diagnostics from and so carries none of its own. Backfill
		// them from setup rather than leaving them zero. A query-completion error already carries
		// its own non-zero values here, so this never overwrites them.
		if copied.Diagnostics == nil {
			copied.Diagnostics = setup.Diagnostics
		}
		if copied.StatusCode == 0 {
			copied.StatusCode = setup.StatusCode
		}
		if copied.SubStatus == 0 {
			copied.SubStatus = setup.SubStatus
		}
		if copied.AttemptCount == 0 {
			copied.AttemptCount = setup.AttemptCount
		}
		return &copied
	}
	code := CodeClientError
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		code = CodeOperationCancelled
	}
	// err here is a raw cause rather than an *Error, such as context.DeadlineExceeded from the
	// setup fetch itself: there is no existing *Error to backfill non-zero fields onto, so the
	// setup fetch's diagnostics/status/attempt metadata is copied in directly, the same metadata
	// the *Error branch above backfills.
	return &Error{
		Code:          code,
		Message:       "fetching query page",
		RequestCharge: setup.RequestCharge,
		ActivityID:    setup.ActivityID,
		Diagnostics:   setup.Diagnostics,
		StatusCode:    setup.StatusCode,
		SubStatus:     setup.SubStatus,
		AttemptCount:  setup.AttemptCount,
		cause:         err,
	}
}
