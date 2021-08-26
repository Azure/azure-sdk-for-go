// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// TableEntityListResponsePager is a Pager for Table entity query results.
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
type TableEntityListResponsePager interface {
	// azcore.Pager

	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() TableEntityListByteResponseResponse
	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
	// Err returns an error if there was an error on the last request
	Err() error
}

type tableEntityQueryResponsePager struct {
	tableClient       *TableClient
	current           *TableEntityListByteResponseResponse
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
func (p *tableEntityQueryResponsePager) PageResponse() TableEntityListByteResponseResponse {
	return *p.current
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableEntityQueryResponsePager) Err() error {
	return p.err
}

// TableListResponsePager is a Pager for Table Queries
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
type TableListResponsePager interface {
	// azcore.Pager

	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() TableListResponseResponse
	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
	// Err returns an error if there was an error on the last request
	Err() error
}

type tableQueryResponsePager struct {
	client            *tableClient
	current           *TableListResponseResponse
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
	p.current = listResponseFromQueryResponse(resp)
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
func (p *tableQueryResponsePager) PageResponse() TableListResponseResponse {
	return *p.current
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableQueryResponsePager) Err() error {
	return p.err
}

// ByteArrayResponse is the return type for a GetEntity operation. The entities properties are stored in the Value property
type ByteArrayResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// ContentType contains the information returned from the Content-Type header response.
	ContentType *string

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// ETag contains the information returned from the ETag header response.
	ETag *string

	// PreferenceApplied contains the information returned from the Preference-Applied header response.
	PreferenceApplied *string

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// The other properties of the table entity.
	Value []byte

	// Version contains the information returned from the x-ms-version header response.
	Version *string

	// XMSContinuationNextPartitionKey contains the information returned from the x-ms-continuation-NextPartitionKey header response.
	XMSContinuationNextPartitionKey *string

	// XMSContinuationNextRowKey contains the information returned from the x-ms-continuation-NextRowKey header response.
	XMSContinuationNextRowKey *string
}

// newByteArrayResponse converts a MapofInterfaceResponse from a map[string]interface{} to a []byte.
func newByteArrayResponse(m MapOfInterfaceResponse) (ByteArrayResponse, error) {
	marshalledValue, err := json.Marshal(m.Value)
	if err != nil {
		return ByteArrayResponse{}, err
	}
	return ByteArrayResponse{
		ClientRequestID:                 m.ClientRequestID,
		ContentType:                     m.ContentType,
		Date:                            m.Date,
		ETag:                            m.ETag,
		PreferenceApplied:               m.PreferenceApplied,
		RawResponse:                     m.RawResponse,
		RequestID:                       m.RequestID,
		Value:                           marshalledValue,
		Version:                         m.Version,
		XMSContinuationNextPartitionKey: m.XMSContinuationNextPartitionKey,
		XMSContinuationNextRowKey:       m.XMSContinuationNextRowKey,
	}, nil
}

// TableEntityListByteResponseResponse is the response envelope for operations that return a TableEntityQueryResponse type.
type TableEntityListByteResponseResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// The properties for the table entity query response.
	TableEntityQueryResponse *TableEntityQueryByteResponse

	// Version contains the information returned from the x-ms-version header response.
	Version *string

	// XMSContinuationNextPartitionKey contains the information returned from the x-ms-continuation-NextPartitionKey header response.
	XMSContinuationNextPartitionKey *string

	// XMSContinuationNextRowKey contains the information returned from the x-ms-continuation-NextRowKey header response.
	XMSContinuationNextRowKey *string
}

// TableEntityQueryByteResponse - The properties for the table entity query response.
type TableEntityQueryByteResponse struct {
	// The metadata response of the table.
	OdataMetadata *string

	// List of table entities.
	Value [][]byte
}

func castToByteResponse(resp *TableEntityQueryResponseResponse) (TableEntityListByteResponseResponse, error) {
	marshalledValue := make([][]byte, 0)
	for _, e := range resp.TableEntityQueryResponse.Value {
		m, err := json.Marshal(e)
		if err != nil {
			return TableEntityListByteResponseResponse{}, err
		}
		marshalledValue = append(marshalledValue, m)
	}

	t := TableEntityQueryByteResponse{
		OdataMetadata: resp.TableEntityQueryResponse.OdataMetadata,
		Value:         marshalledValue,
	}

	return TableEntityListByteResponseResponse{
		ClientRequestID:                 resp.ClientRequestID,
		Date:                            resp.Date,
		RawResponse:                     resp.RawResponse,
		RequestID:                       resp.RequestID,
		TableEntityQueryResponse:        &t,
		Version:                         resp.Version,
		XMSContinuationNextPartitionKey: resp.XMSContinuationNextPartitionKey,
		XMSContinuationNextRowKey:       resp.XMSContinuationNextRowKey,
	}, nil
}

type TableListResponse struct {
	// The metadata response of the table.
	OdataMetadata *string `json:"odata.metadata,omitempty"`

	// List of tables.
	Value []*TableResponseProperties `json:"value,omitempty"`
}

func tableListResponseFromQueryResponse(q *TableQueryResponse) *TableListResponse {
	return &TableListResponse{
		OdataMetadata: q.OdataMetadata,
		Value:         q.Value,
	}
}

// TableListResponseResponse stores the results of a ListTables operation
type TableListResponseResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// The properties for the table query response.
	TableListResponse *TableListResponse

	// Version contains the information returned from the x-ms-version header response.
	Version *string

	// XMSContinuationNextTableName contains the information returned from the x-ms-continuation-NextTableName header response.
	XMSContinuationNextTableName *string
}

func listResponseFromQueryResponse(q TableQueryResponseResponse) *TableListResponseResponse {
	return &TableListResponseResponse{
		ClientRequestID:              q.ClientRequestID,
		Date:                         q.Date,
		RawResponse:                  q.RawResponse,
		RequestID:                    q.RequestID,
		TableListResponse:            tableListResponseFromQueryResponse(q.TableQueryResponse),
		Version:                      q.Version,
		XMSContinuationNextTableName: q.XMSContinuationNextTableName,
	}
}
