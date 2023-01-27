//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/sas"
	"time"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// URLParts object represents the components that make up an Azure Storage Queue URL.
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type URLParts = sas.URLParts

// ParseURL parses a URL initializing URLParts' fields including any SAS-related & snapshot query parameters. Any other
// query parameters remain in the UnparsedParams field. This method overwrites all fields in the URLParts object.
func ParseURL(u string) (URLParts, error) {
	return sas.ParseURL(u)
}

// ================================================================

// CorsRule - CORS is an HTTP feature that enables a web application running under one domain to access resources in another
// domain. Web browsers implement a security restriction known as same-origin policy that
// prevents a web page from calling APIs in a different domain; CORS provides a secure way to allow one domain (the origin
// domain) to call APIs in another domain
type CorsRule = generated.CorsRule

// GeoReplication - Geo-Replication information for the Secondary Storage Service
type GeoReplication = generated.GeoReplication

// RetentionPolicy - the retention policy which determines how long the associated data should persist
type RetentionPolicy = generated.RetentionPolicy

// Metrics - a summary of request statistics grouped by API in hour or minute aggregates for queues
type Metrics = generated.Metrics

// Logging - Azure Analytics Logging settings.
type Logging = generated.Logging

// StorageServiceProperties - Storage Service Properties.
type StorageServiceProperties = generated.StorageServiceProperties

// StorageServiceStats - Stats for the storage service.
type StorageServiceStats = generated.StorageServiceStats

// SignedIdentifier - signed identifier
type SignedIdentifier = generated.SignedIdentifier

// EnqueuedMessage - enqueued message
type EnqueuedMessage = generated.EnqueuedMessage

// DequeuedMessageItem - dequeued message
type DequeuedMessageItem = generated.DequeuedMessageItem

// PeekedMessageItem - peeked message
type PeekedMessageItem = generated.PeekedMessageItem

// ListQueuesSegmentResponse - response segment
type ListQueuesSegmentResponse = generated.ListQueuesSegmentResponse

// QueueItem - queue item
type QueueItem = generated.QueueItem

// ---------------------------------------------------------------------------------------------------------------------

// ListQueuesOptions provides set of configurations for ListQueues operation
type ListQueuesOptions struct {
	Include ListQueuesInclude

	// A string value that identifies the portion of the list of queues to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing operation did not return all queues
	// remaining to be listed with the current page. The NextMarker value can be used as the value for the marker parameter in
	// a subsequent call to request the next page of list items. The marker value is opaque to the client.
	Marker *string

	// Specifies the maximum number of queues to return. If the request does not specify max results, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the listing operation crosses a partition boundary,
	// then the service will return a continuation token for retrieving the remainder of the results. For this reason, it is possible
	// that the service will return fewer results than specified by max results, or than the default of 5000.
	MaxResults *int32

	// Filters the results to return only queues whose name begins with the specified prefix.
	Prefix *string
}

// ListQueuesInclude indicates what additional information the service should return with each queue.
type ListQueuesInclude struct {
	// Tells the service whether to return metadata for each queue.
	Metadata bool
}

// ---------------------------------------------------------------------------------------------------------------------

// SetPropertiesOptions provides set of options for ServiceClient.SetProperties
type SetPropertiesOptions struct {
	// The set of CORS rules.
	Cors []*CorsRule

	// a summary of request statistics grouped by API in hour or minute aggregates for queues
	HourMetrics *Metrics

	// Azure Analytics Logging settings.
	Logging *Logging

	// a summary of request statistics grouped by API in hour or minute aggregates for queues
	MinuteMetrics *Metrics
}

func (o *SetPropertiesOptions) format() (generated.StorageServiceProperties, *generated.ServiceClientSetPropertiesOptions) {
	if o == nil {
		return generated.StorageServiceProperties{}, nil
	}

	defaultVersion := to.Ptr[string]("1.0")
	defaultAge := to.Ptr[int32](0)
	emptyStr := to.Ptr[string]("")

	if o.Cors != nil {
		for i := 0; i < len(o.Cors); i++ {
			if o.Cors[i].AllowedHeaders == nil {
				o.Cors[i].AllowedHeaders = emptyStr
			}
			if o.Cors[i].ExposedHeaders == nil {
				o.Cors[i].ExposedHeaders = emptyStr
			}
			if o.Cors[i].MaxAgeInSeconds == nil {
				o.Cors[i].MaxAgeInSeconds = defaultAge
			}
		}
	}

	if o.HourMetrics != nil {
		if o.HourMetrics.Version == nil {
			o.HourMetrics.Version = defaultVersion
		}
	}

	if o.Logging != nil {
		if o.Logging.Version == nil {
			o.Logging.Version = defaultVersion
		}
	}

	if o.MinuteMetrics != nil {
		if o.MinuteMetrics.Version == nil {
			o.MinuteMetrics.Version = defaultVersion
		}

	}

	return generated.StorageServiceProperties{
		Cors:          o.Cors,
		HourMetrics:   o.HourMetrics,
		Logging:       o.Logging,
		MinuteMetrics: o.MinuteMetrics,
	}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetPropertiesOptions contains the optional parameters for the ServiceClient.GetProperties method.
type GetPropertiesOptions struct {
	// placeholder for future options
}

func (o *GetPropertiesOptions) format() *generated.ServiceClientGetPropertiesOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetStatisticsOptions provides set of options for ServiceClient.GetStatistics
type GetStatisticsOptions struct {
	// placeholder for future options
}

func (o *GetStatisticsOptions) format() *generated.ServiceClientGetStatisticsOptions {
	return nil
}

// -------------------------------------------------QUEUES--------------------------------------------------------------

// CreateOptions contains the optional parameters for creating a queue.
type CreateOptions struct {
	// Optional. Specifies a user-defined name-value pair associated with the queue.
	Metadata map[string]*string
}

func (o *CreateOptions) format() *generated.QueueClientCreateOptions {
	if o == nil {
		return nil
	}
	return &generated.QueueClientCreateOptions{Metadata: o.Metadata}
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteOptions contains the optional parameters for deleting a queue.
type DeleteOptions struct {
}

func (o *DeleteOptions) format() *generated.QueueClientDeleteOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// SetMetadataOptions contains the optional parameters for the QueueClient.SetMetadata method.
type SetMetadataOptions struct {
	Metadata map[string]*string
}

func (o *SetMetadataOptions) format() *generated.QueueClientSetMetadataOptions {
	if o == nil {
		return nil
	}

	return &generated.QueueClientSetMetadataOptions{Metadata: o.Metadata}
}

// ---------------------------------------------------------------------------------------------------------------------

// GetAccessPolicyOptions contains the optional parameters for the QueueClient.GetAccessPolicy method.
type GetAccessPolicyOptions struct {
}

func (o *GetAccessPolicyOptions) format() *generated.QueueClientGetAccessPolicyOptions {
	if o == nil {
		return nil
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// SetAccessPolicyOptions provides set of configurations for QueueClient.SetAccessPolicy operation
type SetAccessPolicyOptions struct {
	QueueACL []*SignedIdentifier
}

func (o *SetAccessPolicyOptions) format() (*generated.QueueClientSetAccessPolicyOptions, []*SignedIdentifier, error) {
	if o == nil {
		return nil, nil, nil
	}
	if o.QueueACL != nil {
		for _, c := range o.QueueACL {
			err := formatTime(c)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return &generated.QueueClientSetAccessPolicyOptions{}, o.QueueACL, nil
}

func formatTime(c *SignedIdentifier) error {
	if c.AccessPolicy == nil {
		return nil
	}

	if c.AccessPolicy.Start != nil {
		st, err := time.Parse(time.RFC3339, c.AccessPolicy.Start.UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
		c.AccessPolicy.Start = &st
	}
	if c.AccessPolicy.Expiry != nil {
		et, err := time.Parse(time.RFC3339, c.AccessPolicy.Expiry.UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
		c.AccessPolicy.Expiry = &et
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetQueuePropertiesOptions contains the optional parameters for the QueueClient.GetProperties method.
type GetQueuePropertiesOptions struct {
}

func (o *GetQueuePropertiesOptions) format() *generated.QueueClientGetPropertiesOptions {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// EnqueueMessageOptions contains the optional parameters for the QueueClient.EnqueueMessage method.
type EnqueueMessageOptions struct {
	TimeToLive        *int32
	VisibilityTimeout *int32
}

func (o *EnqueueMessageOptions) format() *generated.MessagesClientEnqueueOptions {
	if o == nil {
		return nil
	}

	return &generated.MessagesClientEnqueueOptions{MessageTimeToLive: o.TimeToLive,
		Visibilitytimeout: o.VisibilityTimeout}
}

// ---------------------------------------------------------------------------------------------------------------------

// DequeueMessageOptions contains the optional parameters for the QueueClient.EnqueueMessage method.
type DequeueMessageOptions struct {
	VisibilityTimeout *int32
}

func (o *DequeueMessageOptions) format() *generated.MessagesClientDequeueOptions {
	numberOfMessages := int32(1)
	if o == nil {
		return &generated.MessagesClientDequeueOptions{NumberOfMessages: &numberOfMessages}
	}

	return &generated.MessagesClientDequeueOptions{NumberOfMessages: &numberOfMessages,
		Visibilitytimeout: o.VisibilityTimeout}
}

// ---------------------------------------------------------------------------------------------------------------------

// DequeueMessagesOptions contains the optional parameters for the QueueClient.DequeueMessages method.
type DequeueMessagesOptions struct {
	NumberOfMessages  *int32
	VisibilityTimeout *int32
}

func (o *DequeueMessagesOptions) format() *generated.MessagesClientDequeueOptions {
	if o == nil {
		return nil
	}

	return &generated.MessagesClientDequeueOptions{NumberOfMessages: o.NumberOfMessages,
		Visibilitytimeout: o.VisibilityTimeout}
}

// ---------------------------------------------------------------------------------------------------------------------

// UpdateMessageOptions contains the optional parameters for the QueueClient.UpdateMessage method.
type UpdateMessageOptions struct {
	VisibilityTimeout *int32
}

func (o *UpdateMessageOptions) format() *generated.MessageIDClientUpdateOptions {
	if o == nil {
		return nil
	}

	return &generated.MessageIDClientUpdateOptions{Visibilitytimeout: o.VisibilityTimeout}
}

// ---------------------------------------------------------------------------------------------------------------------

// DeleteMessageOptions contains the optional parameters for the QueueClient.DeleteMessage method.
type DeleteMessageOptions struct {
}

func (o *DeleteMessageOptions) format() *generated.MessageIDClientDeleteOptions {
	if o == nil {
		return nil
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// PeekMessageOptions contains the optional parameters for the QueueClient.PeekMessage method.
type PeekMessageOptions struct {
}

func (o *PeekMessageOptions) format() *generated.MessagesClientPeekOptions {
	numberOfMessages := int32(1)
	if o == nil {
		return &generated.MessagesClientPeekOptions{NumberOfMessages: &numberOfMessages}
	}

	return &generated.MessagesClientPeekOptions{NumberOfMessages: &numberOfMessages}
}

// ---------------------------------------------------------------------------------------------------------------------

// PeekMessagesOptions contains the optional parameters for the QueueClient.PeekMessages method.
type PeekMessagesOptions struct {
	NumberOfMessages *int32
}

func (o *PeekMessagesOptions) format() *generated.MessagesClientPeekOptions {
	numberOfMessages := int32(1)
	if o == nil {
		return &generated.MessagesClientPeekOptions{NumberOfMessages: &numberOfMessages}
	}

	return &generated.MessagesClientPeekOptions{NumberOfMessages: o.NumberOfMessages}
}

// ---------------------------------------------------------------------------------------------------------------------

// ClearMessagesOptions contains the optional parameters for the QueueClient.ClearMessages method.
type ClearMessagesOptions struct {
}

func (o *ClearMessagesOptions) format() *generated.MessagesClientClearOptions {
	if o == nil {
		return nil
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// GetSASURLOptions contains the optional parameters for the Client.GetSASURL method.
type GetSASURLOptions struct {
	StartTime *time.Time
}

func (o *GetSASURLOptions) format() time.Time {
	var st time.Time
	if o.StartTime != nil {
		st = o.StartTime.UTC()
	} else {
		st = time.Time{}
	}
	return st
}
