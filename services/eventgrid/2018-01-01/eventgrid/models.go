package eventgrid

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest/date"
)

// BlobCreatedEventData schema of the Data property of an EventGridEvent for an Microsoft.Storage.BlobCreated event.
type BlobCreatedEventData struct {
	// API - The name of the API/operation that triggered this event.
	API *string `json:"api,omitempty"`
	// ClientRequestID - A request id provided by the client of the storage API operation that triggered this event.
	ClientRequestID *string `json:"clientRequestId,omitempty"`
	// RequestID - The request id generated by the Storage service for the storage API operation that triggered this event.
	RequestID *string `json:"requestId,omitempty"`
	// ETag - The etag of the object at the time this event was triggered.
	ETag *string `json:"eTag,omitempty"`
	// ContentType - The content type of the blob. This is the same as what would be returned in the Content-Type header from the blob.
	ContentType *string `json:"contentType,omitempty"`
	// ContentLength - The size of the blob in bytes. This is the same as what would be returned in the Content-Length header from the blob.
	ContentLength *int32 `json:"contentLength,omitempty"`
	// BlobType - The type of blob.
	BlobType *string `json:"blobType,omitempty"`
	// URL - The path to the blob.
	URL *string `json:"url,omitempty"`
	// Sequencer - An opaque string value representing the logical sequence of events for any particular blob name. Users can use standard string comparison to understand the relative sequence of two events on the same blob name.
	Sequencer *string `json:"sequencer,omitempty"`
	// StorageDiagnostics - For service use only. Diagnostic data occasionally included by the Azure Storage service. This property should be ignored by event consumers.
	StorageDiagnostics *map[string]interface{} `json:"storageDiagnostics,omitempty"`
}

// BlobDeletedEventData schema of the Data property of an EventGridEvent for an Microsoft.Storage.BlobDeleted event.
type BlobDeletedEventData struct {
	// API - The name of the API/operation that triggered this event.
	API *string `json:"api,omitempty"`
	// ClientRequestID - A request id provided by the client of the storage API operation that triggered this event.
	ClientRequestID *string `json:"clientRequestId,omitempty"`
	// RequestID - The request id generated by the Storage service for the storage API operation that triggered this event.
	RequestID *string `json:"requestId,omitempty"`
	// ContentType - The content type of the blob. This is the same as what would be returned in the Content-Type header from the blob.
	ContentType *string `json:"contentType,omitempty"`
	// BlobType - The type of blob.
	BlobType *string `json:"blobType,omitempty"`
	// URL - The path to the blob.
	URL *string `json:"url,omitempty"`
	// Sequencer - An opaque string value representing the logical sequence of events for any particular blob name. Users can use standard string comparison to understand the relative sequence of two events on the same blob name.
	Sequencer *string `json:"sequencer,omitempty"`
	// StorageDiagnostics - For service use only. Diagnostic data occasionally included by the Azure Storage service. This property should be ignored by event consumers.
	StorageDiagnostics *map[string]interface{} `json:"storageDiagnostics,omitempty"`
}

// CaptureFileCreatedEventData schema of the Data property of an EventGridEvent for an
// Microsoft.EventHub.CaptureFileCreated event.
type CaptureFileCreatedEventData struct {
	// Fileurl - The path to the capture file.
	Fileurl *string `json:"fileurl,omitempty"`
	// FileType - The file type of the capture file.
	FileType *string `json:"fileType,omitempty"`
	// PartitionID - The shard ID.
	PartitionID *string `json:"partitionId,omitempty"`
	// SizeInBytes - The file size.
	SizeInBytes *int32 `json:"sizeInBytes,omitempty"`
	// EventCount - The number of events in the file.
	EventCount *int32 `json:"eventCount,omitempty"`
	// FirstSequenceNumber - The smallest sequence number from the queue.
	FirstSequenceNumber *int32 `json:"firstSequenceNumber,omitempty"`
	// LastSequenceNumber - The last sequence number from the queue.
	LastSequenceNumber *int32 `json:"lastSequenceNumber,omitempty"`
	// FirstEnqueueTime - The first time from the queue.
	FirstEnqueueTime *date.Time `json:"firstEnqueueTime,omitempty"`
	// LastEnqueueTime - The last time from the queue.
	LastEnqueueTime *date.Time `json:"lastEnqueueTime,omitempty"`
}

// Event properties of an event published to an Event Grid topic.
type Event struct {
	// ID - An unique identifier for the event.
	ID *string `json:"id,omitempty"`
	// Topic - The resource path of the event source.
	Topic *string `json:"topic,omitempty"`
	// Subject - A resource path relative to the topic path.
	Subject *string `json:"subject,omitempty"`
	// Data - Event data specific to the event type.
	Data *map[string]interface{} `json:"data,omitempty"`
	// EventType - The type of the event that occurred.
	EventType *string `json:"eventType,omitempty"`
	// EventTime - The time (in UTC) the event was generated.
	EventTime *date.Time `json:"eventTime,omitempty"`
	// MetadataVersion - The schema version of the event metadata.
	MetadataVersion *string `json:"metadataVersion,omitempty"`
	// DataVersion - The schema version of the data object.
	DataVersion *string `json:"dataVersion,omitempty"`
}
