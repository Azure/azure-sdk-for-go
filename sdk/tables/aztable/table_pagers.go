// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"net/http"

	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// ListEntitiesPager is a Pager for Table entity query results.
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
//     fmt.Printf("The page contains %i results.\n", len(resp.Value))
// }
// err := pager.Err()
type ListEntitiesPager interface {

	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() ListEntitiesResponseEnvelope
	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
	// Err returns an error if there was an error on the last request
	Err() error
}

type tableEntityQueryResponsePager struct {
	tableClient       *TableClient
	current           *ListEntitiesResponseEnvelope
	tableQueryOptions *generated.TableQueryEntitiesOptions
	listOptions       *ListEntitiesOptions
	err               error
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
func (p *tableEntityQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextPartitionKey == nil && p.current.XMSContinuationNextRowKey == nil) {
		return false
	}
	var resp generated.TableQueryEntitiesResponse
	resp, p.err = p.tableClient.client.QueryEntities(ctx, p.tableClient.name, p.tableQueryOptions, p.listOptions.toQueryOptions())

	c, err := castToByteResponse(&resp)
	if err != nil {
		p.err = nil
	}

	p.current = &c
	p.tableQueryOptions.NextPartitionKey = resp.XMSContinuationNextPartitionKey
	p.tableQueryOptions.NextRowKey = resp.XMSContinuationNextRowKey
	return p.err == nil && len(resp.TableEntityQueryResponse.Value) > 0
}

// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.Value))
// }
// err := pager.Err()
func (p *tableEntityQueryResponsePager) PageResponse() ListEntitiesResponseEnvelope {
	return *p.current
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableEntityQueryResponsePager) Err() error {
	return p.err
}

// ListTablesPager is a Pager for Table List operations
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
//     fmt.Printf("The page contains %i results.\n", len(resp.Value))
// }
// err := pager.Err()
type ListTablesPager interface {
	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() ListTablesResponseEnvelope
	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
	// Err returns an error if there was an error on the last request
	Err() error
}

type ListTablesResponseEnvelope struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// XMSContinuationNextTableName contains the information returned from the x-ms-continuation-NextTableName header response.
	XMSContinuationNextTableName *string

	// The metadata response of the table.
	ODataMetadata *string `json:"odata.metadata,omitempty"`

	// List of tables.
	Value []*TableResponseProperties `json:"value,omitempty"`
}

func fromGeneratedTableQueryResponseEnvelope(g *generated.TableQueryResponseEnvelope) *ListTablesResponseEnvelope {
	if g == nil {
		return nil
	}

	var value []*TableResponseProperties

	for _, v := range g.Value {
		value = append(value, fromGeneratedTableResponseProperties(v))
	}

	return &ListTablesResponseEnvelope{
		RawResponse:                  g.RawResponse,
		XMSContinuationNextTableName: g.XMSContinuationNextTableName,
		ODataMetadata:                g.ODataMetadata,
		Value:                        value,
	}
}

type tableQueryResponsePager struct {
	client            *generated.TableClient
	current           *generated.TableQueryResponseEnvelope
	tableQueryOptions *generated.TableQueryOptions
	listOptions       *ListTablesOptions
	err               error
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaulated by calling PageResponse on this Pager.
func (p *tableQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextTableName == nil) {
		return false
	}
	var resp generated.TableQueryResponseEnvelope
	resp, p.err = p.client.Query(ctx, p.tableQueryOptions, p.listOptions.toQueryOptions())
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
func (p *tableQueryResponsePager) PageResponse() ListTablesResponseEnvelope {
	return *fromGeneratedTableQueryResponseEnvelope(p.current)
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableQueryResponsePager) Err() error {
	return p.err
}
