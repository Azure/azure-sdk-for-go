package servicebus

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-amqp-common-go/auth"
	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest/date"
)

const (
	instanceMetadataSchema    = "http://www.w3.org/2001/XMLSchema-instance"
	serviceBusSchema          = "http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"
	dataServiceSchema         = "http://schemas.microsoft.com/ado/2007/08/dataservices"
	dataServiceMetadataSchema = "http://schemas.microsoft.com/ado/2007/08/dataservices/metadata"
	atomSchema                = "http://www.w3.org/2005/Atom"
	applicationXML            = "application/xml"
)

type (
	// EntityManager provides CRUD functionality for Service Bus entities (Queues, Topics, Subscriptions...)
	EntityManager struct {
		TokenProvider auth.TokenProvider
		Host          string
	}

	// Feed is an Atom feed which contains entries
	Feed struct {
		XMLName xml.Name   `xml:"feed"`
		ID      string     `xml:"id"`
		Title   string     `xml:"title"`
		Updated *date.Time `xml:"updated,omitempty"`
		Entries []Entry    `xml:"entry"`
	}

	// Entry is the Atom wrapper for a management request
	Entry struct {
		XMLName                   xml.Name   `xml:"entry"`
		ID                        string     `xml:"id"`
		Title                     string     `xml:"title"`
		Published                 *date.Time `xml:"published,omitempty"`
		Updated                   *date.Time `xml:"updated,omitempty"`
		Author                    *Author    `xml:"author,omitempty"`
		Link                      *Link      `xml:"link,omitempty"`
		Content                   *Content   `xml:"content"`
		DataServiceSchema         string     `xml:"xmlns:d,attr"`
		DataServiceMetadataSchema string     `xml:"xmlns:m,attr"`
		AtomSchema                string     `xml:"xmlns,attr"`
	}

	// Author is an Atom author used in an Entry
	Author struct {
		XMLName xml.Name `xml:"author"`
		Name    *string  `xml:"name,omitempty"`
	}

	// Link is an Atom link used in an Entry
	Link struct {
		XMLName xml.Name `xml:"link"`
		Rel     string   `xml:"rel,attr"`
		HREF    string   `xml:"href,attr"`
	}

	// Content is a generic body for an Atom Entry
	Content struct {
		XMLName xml.Name `xml:"content"`
		Type    string   `xml:"type,attr"`
		Body    string   `xml:",innerxml"`
	}

	// ReceiveBaseDescription provides common fields for Subscriptions and Queues
	ReceiveBaseDescription struct {
		LockDuration                     *string `xml:"LockDuration,omitempty"` // LockDuration - ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1 minute.
		RequiresSession                  *bool   `xml:"RequiresSession,omitempty"`
		DeadLetteringOnMessageExpiration *bool   `xml:"DeadLetteringOnMessageExpiration,omitempty"` // DeadLetteringOnMessageExpiration - A value that indicates whether this queue has dead letter support when a message expires.
		MaxDeliveryCount                 *int32  `xml:"MaxDeliveryCount,omitempty"`                 // MaxDeliveryCount - The maximum delivery count. A message is automatically deadlettered after this number of deliveries. default value is 10.
		MessageCount                     *int64  `xml:"MessageCount,omitempty"`                     // MessageCount - The number of messages in the queue.
	}

	// SendBaseDescription provides common fields for Queues and Topics
	SendBaseDescription struct {
		RequiresDuplicateDetection          *bool   `xml:"RequiresDuplicateDetection,omitempty"`          // RequiresDuplicateDetection - A value indicating if this queue requires duplicate detection.
		DuplicateDetectionHistoryTimeWindow *string `xml:"DuplicateDetectionHistoryTimeWindow,omitempty"` // DuplicateDetectionHistoryTimeWindow - ISO 8601 timeSpan structure that defines the duration of the duplicate detection history. The default value is 10 minutes.
		SizeInBytes                         *int64  `xml:"SizeInBytes,omitempty"`                         // SizeInBytes - The size of the queue, in bytes.
	}

	// BaseEntityDescription provides common fields which are part of Queues, Topics and Subscriptions
	BaseEntityDescription struct {
		InstanceMetadataSchema   string                   `xml:"xmlns:i,attr"`
		ServiceBusSchema         string                   `xml:"xmlns,attr"`
		MaxSizeInMegabytes       *int32                   `xml:"MaxSizeInMegabytes,omitempty"`      // MaxSizeInMegabytes - The maximum size of the queue in megabytes, which is the size of memory allocated for the queue. Default is 1024.
		EnableBatchedOperations  *bool                    `xml:"EnableBatchedOperations,omitempty"` // EnableBatchedOperations - Value that indicates whether server-side batched operations are enabled.
		IsAnonymousAccessible    *bool                    `xml:"IsAnonymousAccessible,omitempty"`
		Status                   *servicebus.EntityStatus `xml:"Status,omitempty"`
		CreatedAt                *date.Time               `xml:"CreatedAt,omitempty"`
		UpdatedAt                *date.Time               `xml:"UpdatedAt,omitempty"`
		SupportOrdering          *bool                    `xml:"SupportOrdering,omitempty"`
		AutoDeleteOnIdle         *string                  `xml:"AutoDeleteOnIdle,omitempty"`
		EnablePartitioning       *bool                    `xml:"EnablePartitioning,omitempty"`
		EnableExpress            *bool                    `xml:"EnableExpress,omitempty"`
		DefaultMessageTimeToLive *string                  `xml:"DefaultMessageTimeToLive,omitempty"` // DefaultMessageTimeToLive - ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.
	}

	managementError struct {
		XMLName xml.Name `xml:"Error"`
		Code    int      `xml:"Code"`
		Detail  string   `xml:"Detail"`
	}
)

func (m *managementError) String() string {
	return fmt.Sprintf("Code: %d, Details: %s", m.Code, m.Detail)
}

// NewEntityManager creates a new instance of an EntityManager given a token provider and host
func NewEntityManager(host string, tokenProvider auth.TokenProvider) *EntityManager {
	return &EntityManager{
		Host:          host,
		TokenProvider: tokenProvider,
	}
}

// Get performs an HTTP Get for a given entity path
func (em *EntityManager) Get(ctx context.Context, entityPath string) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Get")
	defer span.Finish()

	return em.Execute(ctx, http.MethodGet, entityPath, http.NoBody)
}

// Put performs an HTTP PUT for a given entity path and body
func (em *EntityManager) Put(ctx context.Context, entityPath string, body []byte) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Put")
	defer span.Finish()

	return em.Execute(ctx, http.MethodPut, entityPath, bytes.NewReader(body))
}

// Delete performs an HTTP DELETE for a given entity path
func (em *EntityManager) Delete(ctx context.Context, entityPath string) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Delete")
	defer span.Finish()

	return em.Execute(ctx, http.MethodDelete, entityPath, http.NoBody)
}

// Post performs an HTTP POST for a given entity path and body
func (em *EntityManager) Post(ctx context.Context, entityPath string, body []byte) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Post")
	defer span.Finish()

	return em.Execute(ctx, http.MethodPost, entityPath, bytes.NewReader(body))
}

// Execute performs an HTTP request given a http method, path and body
func (em *EntityManager) Execute(ctx context.Context, method string, entityPath string, body io.Reader) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Execute")
	defer span.Finish()

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequest(method, em.Host+strings.TrimPrefix(entityPath, "/"), body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}
	req = addAtomXMLContentType(req)
	req = addAPIVersion201704(req)
	applyRequestInfo(span, req)
	req, err = em.addAuthorization(req)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res, err := client.Do(req)
	applyResponseInfo(span, res)
	if err != nil {
		log.For(ctx).Error(err)
	}
	return res, err
}

func isEmptyFeed(b []byte) bool {
	var emptyFeed queueFeed
	feedErr := xml.Unmarshal(b, &emptyFeed)
	return feedErr == nil && emptyFeed.Title == "Publicly Listed Services"
}

func (em *EntityManager) addAuthorization(req *http.Request) (*http.Request, error) {
	signature, err := em.TokenProvider.GetToken(req.URL.String())
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", signature.Token)
	return req, nil
}

func addAtomXMLContentType(req *http.Request) *http.Request {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		req.Header.Add("Content-Type", "application/atom+xml;type=entry;charset=utf-8")
	}
	return req
}

func addAPIVersion201704(req *http.Request) *http.Request {
	q := req.URL.Query()
	q.Add("api-version", "2017-04")
	req.URL.RawQuery = q.Encode()
	return req
}

func xmlDoc(content []byte) []byte {
	return []byte(xml.Header + string(content))
}

// ptrBool takes a boolean and returns a pointer to that bool. For use in literal pointers, ptrBool(true) -> *bool
func ptrBool(toPtr bool) *bool {
	return &toPtr
}

// ptrString takes a string and returns a pointer to that string. For use in literal pointers,
// ptrString(fmt.Sprintf("..", foo)) -> *string
func ptrString(toPtr string) *string {
	return &toPtr
}

// durationTo8601Seconds takes a duration and returns a string period of whole seconds (int cast of float)
func durationTo8601Seconds(duration *time.Duration) *string {
	return ptrString(fmt.Sprintf("PT%dS", int(duration.Seconds())))
}

func formatManagementError(body []byte) error {
	var mgmtError managementError
	unmarshalErr := xml.Unmarshal(body, &mgmtError)
	if unmarshalErr != nil {
		return errors.New(string(body))
	}

	return fmt.Errorf("error code: %d, Details: %s", mgmtError.Code, mgmtError.Detail)
}
