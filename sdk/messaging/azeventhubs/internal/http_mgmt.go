// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

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

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/devigned/tab"
)

const (
	serviceBusSchema = "http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"
	atomSchema       = "http://www.w3.org/2005/Atom"
	applicationXML   = "application/xml"
)

type (
	// entityManager provides CRUD functionality for Service Bus entities (Queues, Topics, Subscriptions...)
	entityManager struct {
		TokenProvider auth.TokenProvider
		Host          string
	}

	// BaseEntityDescription provides common fields which are part of Queues, Topics and Subscriptions
	BaseEntityDescription struct {
		InstanceMetadataSchema *string `xml:"xmlns:i,attr,omitempty"`
		ServiceBusSchema       *string `xml:"xmlns,attr,omitempty"`
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

// newEntityManager creates a new instance of an entityManager given a token provider and host
func newEntityManager(host string, tokenProvider auth.TokenProvider) *entityManager {
	return &entityManager{
		Host:          host,
		TokenProvider: tokenProvider,
	}
}

// Get performs an HTTP Get for a given entity path
func (em *entityManager) Get(ctx context.Context, entityPath string) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Get")
	defer span.End()

	return em.Execute(ctx, http.MethodGet, entityPath, http.NoBody)
}

// Put performs an HTTP PUT for a given entity path and body
func (em *entityManager) Put(ctx context.Context, entityPath string, body []byte) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Put")
	defer span.End()

	return em.Execute(ctx, http.MethodPut, entityPath, bytes.NewReader(body))
}

// Delete performs an HTTP DELETE for a given entity path
func (em *entityManager) Delete(ctx context.Context, entityPath string) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Delete")
	defer span.End()

	return em.Execute(ctx, http.MethodDelete, entityPath, http.NoBody)
}

// Post performs an HTTP POST for a given entity path and body
func (em *entityManager) Post(ctx context.Context, entityPath string, body []byte) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Post")
	defer span.End()

	return em.Execute(ctx, http.MethodPost, entityPath, bytes.NewReader(body))
}

// Execute performs an HTTP request given a http method, path and body
func (em *entityManager) Execute(ctx context.Context, method string, entityPath string, body io.Reader) (*http.Response, error) {
	span, ctx := em.startSpanFromContext(ctx, "sb.EntityManger.Execute")
	defer span.End()

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequest(method, em.Host+strings.TrimPrefix(entityPath, "/"), body)
	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	req = addAtomXMLContentType(req)
	req = addAPIVersion201704(req)
	applyRequestInfo(span, req)
	req, err = em.addAuthorization(req)
	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := client.Do(req)

	if err != nil {
		tab.For(ctx).Error(err)
	}

	if res != nil {
		applyResponseInfo(span, res)
	}

	return res, err
}

func (em *entityManager) addAuthorization(req *http.Request) (*http.Request, error) {
	signature, err := em.TokenProvider.GetToken(req.URL.String())
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", signature.Token)
	return req, nil
}

func addAtomXMLContentType(req *http.Request) *http.Request {
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		req.Header.Add("content-Type", "application/atom+xml;type=entry;charset=utf-8")
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

func formatManagementError(body []byte) error {
	var mgmtError managementError
	unmarshalErr := xml.Unmarshal(body, &mgmtError)
	if unmarshalErr != nil {
		return errors.New(string(body))
	}

	return fmt.Errorf("error code: %d, Details: %s", mgmtError.Code, mgmtError.Detail)
}
