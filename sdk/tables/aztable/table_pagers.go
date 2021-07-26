// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// TableEntityQueryResponsePager is a Pager for Table entity query results.
//
// NextPage should be called first. It fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
//
// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
type TableEntityQueryResponsePager interface {
	azcore.Pager

	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() TableEntityQueryByteResponseResponse
}

type tableEntityQueryResponsePager struct {
	tableClient       *TableClient
	current           *TableEntityQueryByteResponseResponse
	tableQueryOptions *TableQueryEntitiesOptions
	queryOptions      *ListOptions
	err               error
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
func (p *tableEntityQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextPartitionKey == nil && p.current.XMSContinuationNextRowKey == nil) {
		return false
	}
	var resp TableEntityQueryResponseResponse
	resp, p.err = p.tableClient.client.QueryEntities(ctx, p.tableClient.Name, p.tableQueryOptions, p.queryOptions.toQueryOptions())

	c, err := castToByteResponse(&resp)
	if err != nil {
		p.err = nil
	}

	p.current = &c
	p.tableQueryOptions.NextPartitionKey = resp.XMSContinuationNextPartitionKey
	p.tableQueryOptions.NextRowKey = resp.XMSContinuationNextRowKey
	return p.err == nil && resp.TableEntityQueryResponse != nil && len(resp.TableEntityQueryResponse.Value) > 0
}

// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
func (p *tableEntityQueryResponsePager) PageResponse() TableEntityQueryByteResponseResponse {
	return *p.current
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableEntityQueryResponsePager) Err() error {
	return p.err
}

// TableQueryResponsePager is a Pager for Table Queries
//
// NextPage should be called first. It fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
//
// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
type TableQueryResponsePager interface {
	azcore.Pager

	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() TableQueryResponseResponse
}

type tableQueryResponsePager struct {
	client            *tableClient
	current           *TableQueryResponseResponse
	tableQueryOptions *TableQueryOptions
	queryOptions      *ListOptions
	err               error
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaulated by calling PageResponse on this Pager.
func (p *tableQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextTableName == nil) {
		return false
	}
	var resp TableQueryResponseResponse
	resp, p.err = p.client.Query(ctx, p.tableQueryOptions, p.queryOptions.toQueryOptions())
	p.current = &resp
	p.tableQueryOptions.NextTableName = resp.XMSContinuationNextTableName
	return p.err == nil && resp.TableQueryResponse.Value != nil && len(resp.TableQueryResponse.Value) > 0
}

// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.TableEntityQueryResponse.Value))
// }
func (p *tableQueryResponsePager) PageResponse() TableQueryResponseResponse {
	return *p.current
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableQueryResponsePager) Err() error {
	return p.err
}
