// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"net/http"
	"time"

	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// ByteArrayResponse is the return type for a GetEntity operation. The entities properties are stored in the Value property
type ByteArrayResponse struct {
	// ETag contains the information returned from the ETag header response.
	ETag *string

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// The other properties of the table entity.
	Value []byte
}

func newByteArrayResponse(m generated.TableQueryEntityWithPartitionAndRowKeyResponse) (ByteArrayResponse, error) {
	marshalledValue, err := json.Marshal(m.Value)
	if err != nil {
		return ByteArrayResponse{}, err
	}
	return ByteArrayResponse{
		ETag:        m.ETag,
		RawResponse: m.RawResponse,
		Value:       marshalledValue,
	}, nil
}

// ListEntitiesByteResponse is the response envelope for operations that return a TableEntityQueryResponse type.
type ListEntitiesByteResponse struct {
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
	ODataMetadata *string

	// List of table entities.
	Value [][]byte
}

func castToByteResponse(resp *generated.TableQueryEntitiesResponse) (ListEntitiesByteResponse, error) {
	marshalledValue := make([][]byte, 0)
	for _, e := range resp.TableEntityQueryResponse.Value {
		m, err := json.Marshal(e)
		if err != nil {
			return ListEntitiesByteResponse{}, err
		}
		marshalledValue = append(marshalledValue, m)
	}

	t := TableEntityQueryByteResponse{
		ODataMetadata: resp.TableEntityQueryResponse.ODataMetadata,
		Value:         marshalledValue,
	}

	return ListEntitiesByteResponse{
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

func tableListResponseFromQueryResponse(q *generated.TableQueryResponse) *TableListResponse {
	var value []*TableResponseProperties
	for _, t := range q.Value {
		value = append(value, fromGeneratedTableResponseProperties(t))
	}
	return &TableListResponse{
		OdataMetadata: q.ODataMetadata,
		Value:         value,
	}
}

type TableResponseProperties struct {
	// The edit link of the table.
	ODataEditLink *string `json:"odata.editLink,omitempty"`

	// The id of the table.
	ODataID *string `json:"odata.id,omitempty"`

	// The odata type of the table.
	ODataType *string `json:"odata.type,omitempty"`

	// The name of the table.
	TableName *string `json:"TableName,omitempty"`
}

func fromGeneratedTableResponseProperties(g *generated.TableResponseProperties) *TableResponseProperties {
	if g == nil {
		return nil
	}

	return &TableResponseProperties{
		TableName:     g.TableName,
		ODataEditLink: g.ODataEditLink,
		ODataID:       g.ODataID,
		ODataType:     g.ODataType,
	}
}
