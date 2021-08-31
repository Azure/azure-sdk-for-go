// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// GetEntityResponse is the return type for a GetEntity operation. The individual entities are stored in the Value property
type GetEntityResponse struct {
	// ETag contains the information returned from the ETag header response.
	ETag azcore.ETag

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// The properties of the table entity.
	Value []byte
}

// newGetEntityResponse transforms a generated response to the GetEntityResponse type
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

// ListEntitiesResponseEnvelope is the response envelope for operations that return a list of entities.
type ListEntitiesResponseEnvelope struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
	// XMSContinuationNextPartitionKey contains the information returned from the x-ms-continuation-NextPartitionKey header response.
	XMSContinuationNextPartitionKey *string
	// XMSContinuationNextRowKey contains the information returned from the x-ms-continuation-NextRowKey header response.
	XMSContinuationNextRowKey *string
	// The metadata response of the table.
	ODataMetadata *string
	// List of table entities.
	Entities [][]byte
}

// ListEntitiesResponse - The properties for the table entity query response.
type ListEntitiesResponse struct {
	// The metadata response of the table.
	ODataMetadata *string
	// List of table entities stored as byte slices.
	Entities [][]byte
}

// transforms a generated query response into the ListEntitiesResponseEnveloped
func newListEntitiesResponseEnvelope(resp *generated.TableQueryEntitiesResponse) (ListEntitiesResponseEnvelope, error) {
	marshalledValue := make([][]byte, 0)
	for _, e := range resp.TableEntityQueryResponse.Value {
		m, err := json.Marshal(e)
		if err != nil {
			return ListEntitiesResponseEnvelope{}, err
		}
		marshalledValue = append(marshalledValue, m)
	}

	t := ListEntitiesResponse{
		ODataMetadata: resp.TableEntityQueryResponse.ODataMetadata,
		Entities:      marshalledValue,
	}

	return ListEntitiesResponseEnvelope{
		RawResponse:                     resp.RawResponse,
		XMSContinuationNextPartitionKey: resp.XMSContinuationNextPartitionKey,
		XMSContinuationNextRowKey:       resp.XMSContinuationNextRowKey,
		ODataMetadata:                   t.ODataMetadata,
		Entities:                           t.Entities,
	}, nil
}

// ListTablesResponse contains the properties for a list of tables.
type ListTablesResponse struct {
	// The metadata response of the table.
	OdataMetadata *string `json:"odata.metadata,omitempty"`

	// List of tables.
	Tables []*ResponseProperties `json:"value,omitempty"`
}

// ResponseProperties contains the properties for a single Table
type ResponseProperties struct {
	// The edit link of the table.
	ODataEditLink *string `json:"odata.editLink,omitempty"`

	// The id of the table.
	ODataID *string `json:"odata.id,omitempty"`

	// The odata type of the table.
	ODataType *string `json:"odata.type,omitempty"`

	// The name of the table.
	TableName *string `json:"TableName,omitempty"`
}

// Convets a generated TableResponseProperties to a ResponseProperties
func fromGeneratedTableResponseProperties(g *generated.TableResponseProperties) *ResponseProperties {
	if g == nil {
		return nil
	}

	return &ResponseProperties{
		TableName:     g.TableName,
		ODataEditLink: g.ODataEditLink,
		ODataID:       g.ODataID,
		ODataType:     g.ODataType,
	}
}
