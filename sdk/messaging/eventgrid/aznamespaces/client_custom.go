//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces/internal"
)

// ClientOptions contains optional settings for [Client]
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClientWithSharedKeyCredential creates a [Client] using a shared key.
func NewClientWithSharedKeyCredential(endpoint string, keyCred *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	azc, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewKeyCredentialPolicy(keyCred, "Authorization", &runtime.KeyCredentialPolicyOptions{
				Prefix: "SharedAccessKey ",
			}),
		},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azc,
		endpoint: endpoint,
	}, nil
}

// RejectCloudEvents - Reject batch of Cloud Events. The server responds with an HTTP 200 status code if the request is successfully
// accepted. The response body will include the set of successfully rejected lockTokens,
// along with other failed lockTokens with their corresponding error information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - RejectCloudEventsOptions contains the optional parameters for the Client.RejectCloudEvents method.
func (client *Client) RejectCloudEvents(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *RejectCloudEventsOptions) (RejectCloudEventsResponse, error) {
	return client.internalRejectCloudEvents(ctx, topicName, eventSubscriptionName, rejectOptions{LockTokens: lockTokens}, options)
}

// AcknowledgeCloudEvents - Acknowledge batch of Cloud Events. The server responds with an HTTP 200 status code if the request
// is successfully accepted. The response body will include the set of successfully acknowledged
// lockTokens, along with other failed lockTokens with their corresponding error information. Successfully acknowledged events
// will no longer be available to any consumer.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - AcknowledgeCloudEventsOptions contains the optional parameters for the Client.AcknowledgeCloudEvents method.
func (client *Client) AcknowledgeCloudEvents(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *AcknowledgeCloudEventsOptions) (AcknowledgeCloudEventsResponse, error) {
	return client.internalAcknowledgeCloudEvents(ctx, topicName, eventSubscriptionName, acknowledgeOptions{LockTokens: lockTokens}, options)
}

// ReleaseCloudEvents - Release batch of Cloud Events. The server responds with an HTTP 200 status code if the request is
// successfully accepted. The response body will include the set of successfully released lockTokens,
// along with other failed lockTokens with their corresponding error information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - ReleaseCloudEventsOptions contains the optional parameters for the Client.ReleaseCloudEvents method.
func (client *Client) ReleaseCloudEvents(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *ReleaseCloudEventsOptions) (ReleaseCloudEventsResponse, error) {
	return client.internalReleaseCloudEvents(ctx, topicName, eventSubscriptionName, releaseOptions{LockTokens: lockTokens}, options)
}

// RenewCloudEventLocks - Renew lock for batch of Cloud Events. The server responds with an HTTP 200 status code if the request
// is successfully accepted. The response body will include the set of successfully renewed
// lockTokens, along with other failed lockTokens with their corresponding error information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - topicName - Topic Name.
//   - eventSubscriptionName - Event Subscription Name.
//   - lockTokens - slice of lock tokens.
//   - options - RenewCloudEventLocksOptions contains the optional parameters for the Client.RenewCloudEventLocks method.
func (client *Client) RenewCloudEventLocks(ctx context.Context, topicName string, eventSubscriptionName string, lockTokens []string, options *RenewCloudEventLocksOptions) (RenewCloudEventLocksResponse, error) {
	return client.internalRenewCloudEventLocks(ctx, topicName, eventSubscriptionName, renewLockOptions{LockTokens: lockTokens}, options)
}

// PublishCloudEventOptions contains the optional parameters for the Client.PublishCloudEvent method.
type PublishCloudEventOptions struct {
	// BinaryMode sends a CloudEvent more efficiently by avoiding unnecessary encoding of the Body.
	// There are some caveats to be aware of:
	// - [CloudEvent.Data] must be of type []byte.
	// - [CloudEvent.DataContentType] will be used as the Content-Type for the HTTP request.
	// - [CloudEvent.Extensions] fields are converted to strings.
	BinaryMode bool
}

func (client *Client) publishCloudEventCreateRequest(ctx context.Context, topicName string, event messaging.CloudEvent, options *PublishCloudEventOptions) (*policy.Request, error) {
	if options != nil && options.BinaryMode {
		return client.publishCloudEventCreateRequestUsingBinaryContentMode(ctx, topicName, event, options)
	}

	return client.publishCloudEventCreateRequestUsingJSONEncoding(ctx, topicName, event, options)
}

// publishCloudEventCreateRequestUsingBinaryContentMode creates a request for sending a CloudEvent using [Binary Content mode]
//
// [Binary Content mode]: https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/http-protocol-binding.md#31-binary-content-mode
func (client *Client) publishCloudEventCreateRequestUsingBinaryContentMode(ctx context.Context, topicName string, event messaging.CloudEvent, _ *PublishCloudEventOptions) (*policy.Request, error) {
	urlPath := "/topics/{topicName}:publish"
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}

	req.Raw().Header.Set("ce-id", event.ID)
	req.Raw().Header.Set("ce-specversion", event.SpecVersion)
	req.Raw().Header.Set("ce-time", event.Time.UTC().Format(time.RFC3339))
	req.Raw().Header.Set("ce-source", event.Source)

	if event.Subject != nil {
		req.Raw().Header.Set("ce-subject", *event.Subject)
	}

	req.Raw().Header.Set("ce-type", event.Type)

	if event.DataSchema != nil {
		req.Raw().Header.Set("ce-dataschema", *event.DataSchema)
	}

	contentType := ""

	if event.DataContentType != nil {
		contentType = *event.DataContentType
	}

	asBytes, ok := event.Data.([]byte)

	if !ok {
		return nil, fmt.Errorf("CloudEvent.Data must be of type []byte, was type %T", event.Data)
	}

	bodyStream := streaming.NopCloser(bytes.NewReader(asBytes))

	// encode custom headers
	for k, v := range event.Extensions {
		// https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#type-system
		headerName := "ce-" + k

		str, err := stringize(v)

		if err != nil {
			return nil, err
		}

		req.Raw().Header.Set(headerName, str)
	}

	if err := req.SetBody(bodyStream, contentType); err != nil {
		return nil, err
	}

	return req, nil
}

// stringize converts an arbitrary value into an equivalent string, complying with the types
// in the CloudEvents spec: https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#type-system
func stringize(tmpV any) (string, error) {
	switch v := tmpV.(type) {
	case bool:
		return strconv.FormatBool(v), nil

	// signed integer
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(int64(v), 10), nil

	// unsigned integers
	case uint8: // (also handles `byte`)
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(v), 10), nil

	case []byte:
		return base64.StdEncoding.EncodeToString(v), nil
	case *url.URL:
		return v.String(), nil
	case string:
		return v, nil

	case time.Time:
		return v.Format(time.RFC3339), nil
	case *time.Time:
		return v.Format(time.RFC3339), nil

	case fmt.Stringer:
		return v.String(), nil

	default:
		return "", fmt.Errorf("type %T cannot be converted to a string", v)
	}
}
