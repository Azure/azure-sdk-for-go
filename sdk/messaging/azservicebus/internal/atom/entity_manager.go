// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sbauth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/conn"
)

const (
	serviceBusSchema = "http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"
	atomSchema       = "http://www.w3.org/2005/Atom"
	applicationXML   = "application/xml"
)

type (
	EntityManager interface {
		Get(ctx context.Context, entityPath string, respObj interface{}, mw ...MiddlewareFunc) (*http.Response, error)
		Put(ctx context.Context, entityPath string, body interface{}, respObj interface{}, mw ...MiddlewareFunc) (*http.Response, error)
		Delete(ctx context.Context, entityPath string, mw ...MiddlewareFunc) (*http.Response, error)
		TokenProvider() auth.TokenProvider
	}

	// entityManager provides CRUD functionality for Service Bus entities (Queues, Topics, Subscriptions...)
	entityManager struct {
		tokenProvider auth.TokenProvider
		Host          string
		mwStack       []MiddlewareFunc
		version       string
		retryOptions  exported.RetryOptions
	}

	// BaseEntityDescription provides common fields which are part of Queues, Topics and Subscriptions
	BaseEntityDescription struct {
		InstanceMetadataSchema *string `xml:"xmlns:i,attr,omitempty"`
		ServiceBusSchema       *string `xml:"xmlns,attr,omitempty"`
	}

	// example: <Error><Code>401</Code><Detail>Manage,EntityRead claims required for this operation.</Detail></Error>
	ManagementError struct {
		XMLName xml.Name `xml:"Error"`
		Code    int      `xml:"Code"`
		Detail  string   `xml:"Detail"`
	}

	// CountDetails has current active (and other) messages for queue/topic.
	CountDetails struct {
		XMLName                        xml.Name `xml:"CountDetails"`
		ActiveMessageCount             *int32   `xml:"ActiveMessageCount,omitempty"`
		DeadLetterMessageCount         *int32   `xml:"DeadLetterMessageCount,omitempty"`
		ScheduledMessageCount          *int32   `xml:"ScheduledMessageCount,omitempty"`
		TransferDeadLetterMessageCount *int32   `xml:"TransferDeadLetterMessageCount,omitempty"`
		TransferMessageCount           *int32   `xml:"TransferMessageCount,omitempty"`
	}

	// EntityStatus enumerates the values for entity status.
	EntityStatus string

	// MiddlewareFunc allows a consumer of the entity manager to inject handlers within the request / response pipeline
	//
	// The example below adds the atom xml content type to the request, calls the next middleware and returns the
	// result.
	//
	// addAtomXMLContentType MiddlewareFunc = func(next RestHandler) RestHandler {
	//		return func(ctx context.Context, req *http.Request) (res *http.Response, e error) {
	//			if req.Method != http.MethodGet && req.Method != http.MethodHead {
	//				req.Header.Add("content-Type", "application/atom+xml;type=entry;charset=utf-8")
	//			}
	//			return next(ctx, req)
	//		}
	//	}
	MiddlewareFunc func(next RestHandler) RestHandler

	// RestHandler is used to transform a request and response within the http pipeline
	RestHandler func(ctx context.Context, req *http.Request) (*http.Response, error)
)

var (
	addAtomXMLContentType MiddlewareFunc = func(next RestHandler) RestHandler {
		return func(ctx context.Context, req *http.Request) (res *http.Response, e error) {
			if req.Method != http.MethodGet && req.Method != http.MethodHead {
				req.Header.Add("content-Type", "application/atom+xml;type=entry;charset=utf-8")
			}
			return next(ctx, req)
		}
	}

	addAPIVersion201704 MiddlewareFunc = func(next RestHandler) RestHandler {
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			q.Add("api-version", "2017-04")
			req.URL.RawQuery = q.Encode()
			return next(ctx, req)
		}
	}
)

const (
	// Active ...
	Active EntityStatus = "Active"
	// Creating ...
	Creating EntityStatus = "Creating"
	// Deleting ...
	Deleting EntityStatus = "Deleting"
	// Disabled ...
	Disabled EntityStatus = "Disabled"
	// ReceiveDisabled ...
	ReceiveDisabled EntityStatus = "ReceiveDisabled"
	// Renaming ...
	Renaming EntityStatus = "Renaming"
	// Restoring ...
	Restoring EntityStatus = "Restoring"
	// SendDisabled ...
	SendDisabled EntityStatus = "SendDisabled"
	// Unknown ...
	Unknown EntityStatus = "Unknown"
)

func (m *ManagementError) String() string {
	return fmt.Sprintf("Code: %d, Details: %s", m.Code, m.Detail)
}

// NewEntityManagerWithConnectionString creates an entity manager (a lower level HTTP client
// for the ATOM endpoint). This is typically wrapped by an entity specific client (like
// TopicManager, QueueManager or , SubscriptionManager).
func NewEntityManagerWithConnectionString(connectionString string, version string) (EntityManager, error) {
	parsed, err := conn.ParsedConnectionFromStr(connectionString)

	if err != nil {
		return nil, err
	}

	provider, err := sbauth.NewTokenProviderWithConnectionString(parsed)

	if err != nil {
		return nil, err
	}

	return &entityManager{
		Host:          fmt.Sprintf("https://%s/", parsed.Namespace),
		version:       version,
		tokenProvider: provider,
		mwStack: []MiddlewareFunc{
			addAPIVersion201704,
			addAtomXMLContentType,
			addAuthorization(provider),
		},
	}, nil
}

// NewEntityManager creates an entity manager using a TokenCredential.
func NewEntityManager(ns string, tokenCredential azcore.TokenCredential, version string, retryOptions exported.RetryOptions) (EntityManager, error) {
	return &entityManager{
		Host:          fmt.Sprintf("https://%s/", ns),
		version:       version,
		tokenProvider: sbauth.NewTokenProvider(tokenCredential),
		mwStack: []MiddlewareFunc{
			addAPIVersion201704,
			addAtomXMLContentType,
			addAuthorization(sbauth.NewTokenProvider(tokenCredential)),
		},
		retryOptions: retryOptions,
	}, nil
}

// Get performs an HTTP Get for a given entity path, deserializing the returned XML into `respObj`
func (em *entityManager) Get(ctx context.Context, entityPath string, respObj interface{}, mw ...MiddlewareFunc) (*http.Response, error) {
	resp, err := em.execute(ctx, http.MethodGet, entityPath, http.NoBody, mw...)
	defer CloseRes(ctx, resp)

	if err != nil {
		return resp, err
	}

	return deserializeBody(resp, respObj)
}

// Put performs an HTTP PUT for a given entity path and body, deserializing the returned XML into `respObj`
func (em *entityManager) Put(ctx context.Context, entityPath string, body interface{}, respObj interface{}, mw ...MiddlewareFunc) (*http.Response, error) {
	bodyBytes, err := xml.Marshal(body)

	if err != nil {
		return nil, err
	}

	resp, err := em.execute(ctx, http.MethodPut, entityPath, bytes.NewReader(bodyBytes), mw...)
	defer CloseRes(ctx, resp)

	if err != nil {
		return resp, err
	}

	return deserializeBody(resp, respObj)
}

// Delete performs an HTTP DELETE for a given entity path
func (em *entityManager) Delete(ctx context.Context, entityPath string, mw ...MiddlewareFunc) (*http.Response, error) {
	return em.execute(ctx, http.MethodDelete, entityPath, http.NoBody, mw...)
}

func (em *entityManager) execute(ctx context.Context, method string, entityPath string, body io.Reader, mw ...MiddlewareFunc) (*http.Response, error) {
	var finalResp *http.Response

	err := utils.Retry(ctx, exported.EventAdmin, fmt.Sprintf("%s %s", method, entityPath), func(ctx context.Context, args *utils.RetryFnArgs) error {
		req, err := http.NewRequest(method, em.Host+strings.TrimPrefix(entityPath, "/"), body)
		if err != nil {
			return err
		}

		final := func(_ RestHandler) RestHandler {
			return func(reqCtx context.Context, request *http.Request) (*http.Response, error) {
				client := &http.Client{
					Timeout: 60 * time.Second,
				}
				request = request.WithContext(reqCtx)
				return client.Do(request)
			}
		}

		mwStack := []MiddlewareFunc{final}
		sl := len(em.mwStack) - 1
		for i := sl; i >= 0; i-- {
			mwStack = append(mwStack, em.mwStack[i])
		}

		for i := len(mw) - 1; i >= 0; i-- {
			mwStack = append(mwStack, mw[i])
		}

		var h RestHandler
		for _, mw := range mwStack {
			h = mw(h)
		}

		resp, err := h(ctx, req)

		if err == nil {
			if resp.StatusCode >= http.StatusBadRequest {
				return NewResponseError(resp)
			}

			finalResp = resp
			return nil
		}

		if resp != nil {
			return NewResponseError(resp)
		}

		return err
	}, isFatalHTTPError, em.retryOptions)

	if err != nil {
		return nil, err
	}

	return finalResp, nil
}

func isFatalHTTPError(err error) bool {
	var netErr net.Error

	if errors.As(err, &netErr) {
		return false
	}

	var respErr *azcore.ResponseError

	// TODO: this is very much temporary. We need to move this over to the azcore HTTP stack.
	if errors.As(err, &respErr) {
		if respErr.StatusCode == http.StatusRequestTimeout || // 408
			respErr.StatusCode == http.StatusTooManyRequests || // 429
			respErr.StatusCode == http.StatusInternalServerError || // 500
			respErr.StatusCode == http.StatusBadGateway || // 502
			respErr.StatusCode == http.StatusServiceUnavailable || // 503
			respErr.StatusCode == http.StatusGatewayTimeout { // 504	)
			return false
		}
	}

	return true
}

// Use adds middleware to the middleware mwStack
func (em *entityManager) Use(mw ...MiddlewareFunc) {
	em.mwStack = append(em.mwStack, mw...)
}

// TokenProvider generates authorization tokens for communicating with the Service Bus management API
func (em *entityManager) TokenProvider() auth.TokenProvider {
	return em.tokenProvider
}

func FormatManagementError(body []byte, origErr error) error {
	var mgmtError ManagementError
	unmarshalErr := xml.Unmarshal(body, &mgmtError)
	if unmarshalErr != nil {
		return origErr
	}

	return fmt.Errorf("error code: %d, Details: %s", mgmtError.Code, mgmtError.Detail)
}

func addAuthorization(tp auth.TokenProvider) MiddlewareFunc {
	return func(next RestHandler) RestHandler {
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			signature, err := tp.GetToken(req.URL.String())
			if err != nil {
				return nil, err
			}

			req.Header.Add("Authorization", signature.Token)
			return next(ctx, req)
		}
	}
}

func addSupplementalAuthorization(supplementalURI string, tp auth.TokenProvider) MiddlewareFunc {
	return func(next RestHandler) RestHandler {
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			signature, err := tp.GetToken(supplementalURI)
			if err != nil {
				return nil, err
			}

			req.Header.Add("ServiceBusSupplementaryAuthorization", signature.Token)
			return next(ctx, req)
		}
	}
}

func addDeadLetterSupplementalAuthorization(targetURI string, tp auth.TokenProvider) MiddlewareFunc {
	return func(next RestHandler) RestHandler {
		return func(ctx context.Context, req *http.Request) (response *http.Response, e error) {
			signature, err := tp.GetToken(targetURI)
			if err != nil {
				return nil, err
			}

			req.Header.Add("ServiceBusDlqSupplementaryAuthorization", signature.Token)
			return next(ctx, req)
		}
	}
}

// TraceReqAndResponseMiddleware will print the dump of the management request and response.
//
// This should only be used for debugging or educational purposes.
func TraceReqAndResponseMiddleware() MiddlewareFunc {
	return func(next RestHandler) RestHandler {
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			if dump, err := httputil.DumpRequest(req, true); err == nil {
				fmt.Println(string(dump))
			}

			res, err := next(ctx, req)

			if dump, err := httputil.DumpResponse(res, true); err == nil {
				fmt.Println(string(dump))
			}

			return res, err
		}
	}
}

var ErrFeedEmpty = errors.New("entity does not exist")

// deserializeBody deserializes the body of the response into the type specified by respObj
// (similar to xml.Unmarshal, which this func is calling).
// If an empty feed is found, it returns nil.
func deserializeBody(resp *http.Response, respObj interface{}) (*http.Response, error) {
	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return resp, err
	}

	if err := xml.Unmarshal(bytes, respObj); err != nil {
		// In ATOM when you request a specific entity (queue, topic, sub) you typically get an
		// <Entry>. However, if the entity is not found, instead of getting a 404 you actually
		// get a <Feed> XML object that is empty and an HTTP status code of 200.
		//
		// So the combination of "can't deserialize object" and "it's an empty feed" are enough
		// for us to note that we weren't expecting a feed (ie, GET /queue) and that the feed
		// itself is the special "empty feed".
		var emptyFeed QueueFeed
		feedErr := xml.Unmarshal(bytes, &emptyFeed)

		if feedErr == nil && emptyFeed.Title == "Publicly Listed Services" {
			return resp, ErrFeedEmpty
		}

		return resp, err
	}

	return resp, nil
}
