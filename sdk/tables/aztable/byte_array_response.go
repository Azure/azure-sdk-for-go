// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// GetEntityResponse is the return type for a GetEntity operation. The entities properties are stored in the Value property
type GetEntityResponse struct {
	// ETag contains the information returned from the ETag header response.
	ETag azcore.ETag

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// The other properties of the table entity.
	Value []byte
}

func newGetEntityResponse(m generated.TableQueryEntityWithPartitionAndRowKeyResponse) (GetEntityResponse, error) {
	marshalledValue, err := json.Marshal(m.Value)
	if err != nil {
		return GetEntityResponse{}, err
	}
	return GetEntityResponse{
		ETag:        azcore.ETag(*m.ETag),
		RawResponse: m.RawResponse,
		Value:       marshalledValue,
	}, nil
}

// ListEntitiesResponseEnvelope is the response envelope for operations that return a TableEntityQueryResponse type.
type ListEntitiesResponseEnvelope struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
	// The properties for the table entity query response.
	// TableEntityQueryResponse *TableEntityListResponse
	// XMSContinuationNextPartitionKey contains the information returned from the x-ms-continuation-NextPartitionKey header response.
	XMSContinuationNextPartitionKey *string
	// XMSContinuationNextRowKey contains the information returned from the x-ms-continuation-NextRowKey header response.
	XMSContinuationNextRowKey *string
	// The metadata response of the table.
	ODataMetadata *string

	// List of table entities.
	Value [][]byte
}

// TableEntityListResponse - The properties for the table entity query response.
type TableEntityListResponse struct {
	// The metadata response of the table.
	ODataMetadata *string

	// List of table entities.
	Value [][]byte
}

func castToByteResponse(resp *generated.TableQueryEntitiesResponse) (ListEntitiesResponseEnvelope, error) {
	marshalledValue := make([][]byte, 0)
	for _, e := range resp.TableEntityQueryResponse.Value {
		m, err := json.Marshal(e)
		if err != nil {
			return ListEntitiesResponseEnvelope{}, err
		}
		marshalledValue = append(marshalledValue, m)
	}

	t := TableEntityListResponse{
		ODataMetadata: resp.TableEntityQueryResponse.ODataMetadata,
		Value:         marshalledValue,
	}

	return ListEntitiesResponseEnvelope{
		RawResponse:                     resp.RawResponse,
		// TableEntityQueryResponse:        &t,
		XMSContinuationNextPartitionKey: resp.XMSContinuationNextPartitionKey,
		XMSContinuationNextRowKey:       resp.XMSContinuationNextRowKey,
		ODataMetadata:                   t.ODataMetadata,
		Value:                           t.Value,
	}, nil
}

type TableListResponse struct {
	// The metadata response of the table.
	OdataMetadata *string `json:"odata.metadata,omitempty"`

	// List of tables.
	Value []*TableResponseProperties `json:"value,omitempty"`
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
