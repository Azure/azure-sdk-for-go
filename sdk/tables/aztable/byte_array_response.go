// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/json"
	"net/http"
	"time"
)

// ByteArrayResponse converts the MapOfInterfaceResponse.Value from a map[string]interface{} to a []byte
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

// TableEntityQueryByteResponseResponse is the response envelope for operations that return a TableEntityQueryResponse type.
type TableEntityQueryByteResponseResponse struct {
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

func castToByteResponse(resp *TableEntityQueryResponseResponse) (TableEntityQueryByteResponseResponse, error) {
	marshalledValue := make([][]byte, 0)
	for _, e := range resp.TableEntityQueryResponse.Value {
		m, err := json.Marshal(e)
		if err != nil {
			return TableEntityQueryByteResponseResponse{}, err
		}
		marshalledValue = append(marshalledValue, m)
	}

	t := TableEntityQueryByteResponse{
		OdataMetadata: resp.TableEntityQueryResponse.OdataMetadata,
		Value:         marshalledValue,
	}

	return TableEntityQueryByteResponseResponse{
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
