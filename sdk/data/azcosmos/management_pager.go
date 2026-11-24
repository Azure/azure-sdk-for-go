// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
)

// ManagementQueryOptions configures database or container discovery.
type ManagementQueryOptions struct {
	// Operation holds shared settings, applied separately to each page.
	Operation OperationOptions
	// PageSizeHint is a positive page-size hint. Zero uses the native default.
	PageSizeHint int32
	// ContinuationToken resumes the same query. Empty starts a new query.
	ContinuationToken string
}

// ManagementQueryResponse contains a page of management resources.
type ManagementQueryResponse[T any] struct {
	Response
	// Items contains the resources on this page. An empty page need not be the last.
	Items []T
}

// NewQueryDatabasesPager queries databases. Construction performs no network I/O.
// Queries must return resource objects, not scalar projections.
func (c *Client) NewQueryDatabasesPager(query Query, options *ManagementQueryOptions) *ManagementPager[DatabaseProperties] {
	return newManagementPager[DatabaseProperties](c, operationKindQueryDatabases, "", query, options)
}

// NewQueryContainersPager queries this database's containers without resolving any container.
// Queries must return resource objects, not scalar projections.
func (d *DatabaseClient) NewQueryContainersPager(query Query, options *ManagementQueryOptions) *ManagementPager[ContainerProperties] {
	return newManagementPager[ContainerProperties](d.client, operationKindQueryContainers, d.id, query, options)
}

// ManagementPager owns a native management query cursor. Close it when stopping early.
// Methods serialize iteration and cleanup; native retries and routing remain driver-owned.
type ManagementPager[T any] struct {
	mu            sync.Mutex
	client        *Client
	request       managementRequest
	options       ManagementQueryOptions
	validationErr error
	cursor        queryCursor
	done          bool
}

func newManagementPager[T any](client *Client, kind operationKind, databaseID string, query Query, options *ManagementQueryOptions) *ManagementPager[T] {
	p := &ManagementPager[T]{client: client, request: managementRequest{kind: kind, databaseID: databaseID}}
	p.request.body, p.validationErr = query.body()
	if options != nil {
		p.options = *options
		p.options.Operation = options.Operation.clone()
	}
	p.request.operation = p.options.Operation
	if p.validationErr == nil {
		p.validationErr = p.options.Operation.validate()
	}
	if p.validationErr == nil && (p.options.PageSizeHint < 0 || strings.ContainsRune(p.options.ContinuationToken, 0)) {
		p.validationErr = errors.New("azcosmos: management page size must not be negative and continuation must not contain NUL")
	}
	return p
}

// More reports whether another page may be fetched. Empty pages do not imply exhaustion.
func (p *ManagementPager[T]) More() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return !p.done
}

// NextPage fetches one page. Execution or decoding errors close the cursor.
// Argument and pre-submission cancellation errors leave the pager usable.
func (p *ManagementPager[T]) NextPage(ctx context.Context) (ManagementQueryResponse[T], error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.validationErr != nil {
		return ManagementQueryResponse[T]{}, p.validationErr
	}
	if err := validateManagementContext(ctx, p.options.Operation); err != nil {
		return ManagementQueryResponse[T]{}, err
	}
	release, err := p.client.acquire()
	if err != nil {
		return ManagementQueryResponse[T]{}, err
	}
	defer release()
	if p.done {
		return ManagementQueryResponse[T]{}, &Error{Code: CodeClientError, Message: "azcosmos: management pager is closed or exhausted"}
	}
	if err := ctx.Err(); err != nil {
		return ManagementQueryResponse[T]{}, err
	}
	ctx, cancel, err := p.client.queryContext(ctx, p.options.Operation)
	if err != nil {
		return ManagementQueryResponse[T]{}, err
	}
	defer cancel()
	if p.cursor == nil {
		p.cursor, err = p.client.openManagementQuery(ctx, p.request, p.options)
		if err != nil {
			p.finish()
			return ManagementQueryResponse[T]{}, err
		}
	}
	page, end, err := p.cursor.next(ctx)
	if err != nil {
		p.finish()
		return ManagementQueryResponse[T]{}, err
	}
	out := ManagementQueryResponse[T]{Response: page.Response, Items: make([]T, len(page.Items))}
	for i, body := range page.Items {
		if err := ctx.Err(); err != nil {
			p.finish()
			failure := queryResponseError(page.Response, page.SessionToken, err)
			failure.Code = CodeOperationCancelled
			return ManagementQueryResponse[T]{}, failure
		}
		if err := json.Unmarshal(body, &out.Items[i]); err != nil {
			p.finish()
			return ManagementQueryResponse[T]{}, newManagementSerializationError(ItemResponse{Response: page.Response, SessionToken: page.SessionToken}, body, "query resource", err)
		}
	}
	if err := ctx.Err(); err != nil {
		p.finish()
		failure := queryResponseError(page.Response, page.SessionToken, err)
		failure.Code = CodeOperationCancelled
		return ManagementQueryResponse[T]{}, failure
	}
	if end {
		p.finish()
	}
	return out, nil
}

// ContinuationToken returns the last delivered page's opaque server continuation.
// Empty means the last delivered service page had no continuation. Use More for iteration.
func (p *ManagementPager[T]) ContinuationToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.validationErr != nil {
		return "", p.validationErr
	}
	if err := validateManagementContext(ctx, p.options.Operation); err != nil {
		return "", err
	}
	release, err := p.client.acquire()
	if err != nil {
		return "", err
	}
	defer release()
	if p.cursor == nil || p.done {
		return "", &Error{Code: CodeClientError, Message: "azcosmos: checkpoint requires an open management pager"}
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	token, err := p.cursor.checkpoint(ctx)
	if err != nil && !unsupportedQueryCheckpoint(err) {
		p.finish()
	}
	return token, err
}

// Close releases the cursor. It is idempotent and safe alongside Client.Close.
func (p *ManagementPager[T]) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.finish()
	return nil
}

func (p *ManagementPager[T]) finish() {
	p.done = true
	if p.cursor != nil {
		p.cursor.close()
		p.cursor = nil
	}
}

type managementRequest struct {
	kind         operationKind
	databaseID   string
	containerID  string
	resourceLink string
	body         []byte
	operation    OperationOptions
	precondition Precondition
	throughput   *ThroughputProperties
}

type managementResult struct {
	response ItemResponse
	body     []byte
	pending  bool
}
