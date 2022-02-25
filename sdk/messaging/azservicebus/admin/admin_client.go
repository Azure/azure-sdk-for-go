// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
)

// Client allows you to administer resources in a Service Bus Namespace.
// For example, you can create queues, enabling capabilities like partitioning, duplicate detection, etc..
// NOTE: For sending and receiving messages you'll need to use the `azservicebus.Client` type instead.
type Client struct {
	em atom.EntityManager
}

// RetryOptions represent the options for retries.
type RetryOptions = utils.RetryOptions

// ClientOptions allows you to set optional configuration for `Client`.
type ClientOptions struct {
	// RetryOptions controls how often operations are retried from this client.
	RetryOptions *RetryOptions
}

// NewClientFromConnectionString creates a Client authenticating using a connection string.
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	em, err := atom.NewEntityManagerWithConnectionString(connectionString, internal.Version)

	if err != nil {
		return nil, err
	}

	return &Client{em: em}, nil
}

// NewClient creates a Client authenticating using a TokenCredential.
func NewClient(fullyQualifiedNamespace string, tokenCredential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	var retryOptions utils.RetryOptions

	if options != nil && options.RetryOptions != nil {
		retryOptions = *options.RetryOptions
	}

	em, err := atom.NewEntityManager(fullyQualifiedNamespace, tokenCredential, internal.Version, retryOptions)

	if err != nil {
		return nil, err
	}

	return &Client{em: em}, nil
}

type NamespaceProperties struct {
	CreatedTime  time.Time
	ModifiedTime time.Time

	SKU            string
	MessagingUnits *int64
	Name           string
}

type GetNamespacePropertiesResult struct {
	NamespaceProperties
}

type GetNamespacePropertiesResponse struct {
	GetNamespacePropertiesResult
	RawResponse *http.Response
}

type GetNamespacePropertiesOptions struct {
	// For future expansion
}

// GetNamespaceProperties gets the properties for the namespace, includings properties like SKU and CreatedTime.
func (ac *Client) GetNamespaceProperties(ctx context.Context, options *GetNamespacePropertiesOptions) (*GetNamespacePropertiesResponse, error) {
	var body *atom.NamespaceEntry
	resp, err := ac.em.Get(ctx, "/$namespaceinfo", &body)

	if err != nil {
		return nil, err
	}

	props := &GetNamespacePropertiesResponse{
		RawResponse: resp,
		GetNamespacePropertiesResult: GetNamespacePropertiesResult{
			NamespaceProperties: NamespaceProperties{
				Name:           body.NamespaceInfo.Name,
				SKU:            body.NamespaceInfo.MessagingSKU,
				MessagingUnits: body.NamespaceInfo.MessagingUnits,
			},
		},
	}

	if props.CreatedTime, err = atom.StringToTime(body.NamespaceInfo.CreatedTime); err != nil {
		return nil, err
	}

	if props.ModifiedTime, err = atom.StringToTime(body.NamespaceInfo.ModifiedTime); err != nil {
		return nil, err
	}
	return props, nil
}

type pagerFunc func(ctx context.Context, pv interface{}) (*http.Response, error)

// newPagerFunc gets a function that can be used to page sequentially through an ATOM resource
func (ac *Client) newPagerFunc(baseFragment string, maxPageSize int32, lenV func(pv interface{}) int) pagerFunc {
	eof := false
	skip := int32(0)

	return func(ctx context.Context, pv interface{}) (*http.Response, error) {
		if eof {
			return nil, nil
		}

		url := baseFragment + "?"
		if maxPageSize > 0 {
			url += fmt.Sprintf("&$top=%d", maxPageSize)
		}

		if skip > 0 {
			url += fmt.Sprintf("&$skip=%d", skip)
		}

		resp, err := ac.em.Get(ctx, url, pv)

		if err != nil {
			eof = true
			return nil, err
		}

		if lenV(pv) == 0 {
			eof = true
			return nil, nil
		}

		if lenV(pv) < int(maxPageSize) {
			eof = true
		}

		skip += int32(lenV(pv))
		return resp, nil
	}
}
