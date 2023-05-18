//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

// Responder represents a scalar response.
type Responder[T any] struct {
	httpStatus int
	resp       T
	opts       SetResponseOptions
}

// SetResponse sets the specified value to be returned.
//   - httpStatus is the HTTP status code to be returned
//   - resp is the response to be returned
//   - o contains optional values, pass nil to accept the defaults
func (r *Responder[T]) SetResponse(httpStatus int, resp T, o *SetResponseOptions) {
	r.httpStatus = httpStatus
	r.resp = resp
	if o != nil {
		r.opts = *o
	}
}

// SetResponseOptions contains the optional values for Responder[T].SetResponse.
type SetResponseOptions struct {
	// Header contains optional HTTP headers to include in the response.
	Header http.Header
}

// GetResponse returns the response associated with the Responder.
// This function is called by the fake server internals.
func (r Responder[T]) GetResponse() T {
	return r.resp
}

// GetResponseContent returns the ResponseContent associated with the Responder.
// This function is called by the fake server internals.
func (r Responder[T]) GetResponseContent() ResponseContent {
	return ResponseContent{HTTPStatus: r.httpStatus, Header: r.opts.Header}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ErrorResponder represents a scalar error response.
type ErrorResponder struct {
	err error
}

// SetError sets the specified error to be returned.
// Use SetResponseError for returning an *azcore.ResponseError.
func (e *ErrorResponder) SetError(err error) {
	e.err = shared.NonRetriableError(err)
}

// SetResponseError sets an *azcore.ResponseError with the specified values to be returned.
//   - errorCode is the value to be used in the ResponseError.Code field
//   - httpStatus is the HTTP status code
func (e *ErrorResponder) SetResponseError(httpStatus int, errorCode string) {
	e.err = shared.NonRetriableError(&exported.ResponseError{ErrorCode: errorCode, StatusCode: httpStatus})
}

// GetError returns the error for this responder.
// This function is called by the fake server internals.
func (e ErrorResponder) GetError(req *http.Request) error {
	if e.err == nil {
		return nil
	}

	var respErr *azcore.ResponseError
	if errors.As(e.err, &respErr) {
		// fix up the raw response
		rawResp, err := newErrorResponse(respErr.StatusCode, respErr.ErrorCode, req)
		if err != nil {
			return shared.NonRetriableError(err)
		}
		respErr.RawResponse = rawResp
	}
	return shared.NonRetriableError(e.err)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// PagerResponder represents a sequence of paged responses.
// Responses are replayed in the order in which they were added.
type PagerResponder[T any] struct {
	pages []any
}

// AddPage adds a page to the sequence of respones.
//   - page is the response page to be added
//   - o contains optional values, pass nil to accept the defaults
func (p *PagerResponder[T]) AddPage(httpStatus int, page T, o *AddPageOptions) {
	p.pages = append(p.pages, pageResp[T]{httpStatus: httpStatus, entry: page})
}

// AddError adds an error to the sequence of responses.
// The error is returned from the call to runtime.Pager[T].NextPage().
func (p *PagerResponder[T]) AddError(err error) {
	p.pages = append(p.pages, shared.NonRetriableError(err))
}

// AddResponseError adds an *azcore.ResponseError to the sequence of responses.
// The error is returned from the call to runtime.Pager[T].NextPage().
func (p *PagerResponder[T]) AddResponseError(httpStatus int, errorCode string) {
	p.pages = append(p.pages, shared.NonRetriableError(&exported.ResponseError{ErrorCode: errorCode, StatusCode: httpStatus}))
}

// AddPageOptions contains the optional values for PagerResponder[T].AddPage.
type AddPageOptions struct {
	// placeholder for future options
}

// Next returns the next response in the sequence (a T or an error).
// This function is called by the fake server internals.
func (p *PagerResponder[T]) Next(req *http.Request) (*http.Response, error) {
	if len(p.pages) == 0 {
		return nil, shared.NonRetriableError(errors.New("paged response has no pages"))
	}

	page := p.pages[0]
	p.pages = p.pages[1:]

	pageT, ok := page.(pageResp[T])
	if ok {
		body, err := json.Marshal(pageT.entry)
		if err != nil {
			return nil, shared.NonRetriableError(err)
		}
		content := ResponseContent{
			HTTPStatus: pageT.httpStatus,
			Header:     http.Header{},
		}
		resp, err := NewResponse(content, req)
		if err != nil {
			return nil, shared.NonRetriableError(err)
		}
		return SetResponseBody(resp, body, shared.ContentTypeAppJSON), nil
	}

	err := page.(error)
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		// fix up the raw response
		rawResp, err := newErrorResponse(respErr.StatusCode, respErr.ErrorCode, req)
		if err != nil {
			return nil, shared.NonRetriableError(err)
		}
		respErr.RawResponse = rawResp
	}
	return nil, shared.NonRetriableError(err)
}

// More returns true if there are more responses for consumption.
// This function is called by the fake server internals.
func (p *PagerResponder[T]) More() bool {
	return len(p.pages) > 0
}

type pageindex[T any] struct {
	i    int
	page pageResp[T]
}

// InjectNextLinks is used to populate the nextLink field.
// The inject callback is executed for every T in the sequence except for the last one.
// This function is called by the fake server internals.
func (p *PagerResponder[T]) InjectNextLinks(req *http.Request, inject func(page *T, createLink func() string)) {
	// first find all the actual pages in the list
	pages := make([]pageindex[T], 0, len(p.pages))
	for i := range p.pages {
		if pageT, ok := p.pages[i].(pageResp[T]); ok {
			pages = append(pages, pageindex[T]{
				i:    i,
				page: pageT,
			})
		}
	}

	// now populate the next links
	for i := range pages {
		if i+1 == len(pages) {
			// no nextLink for last page
			break
		}

		inject(&pages[i].page.entry, func() string {
			return fmt.Sprintf("%s://%s%s/page_%d", req.URL.Scheme, req.URL.Host, req.URL.Path, i+1)
		})

		// update the original slice with the modified page
		p.pages[pages[i].i] = pages[i].page
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// PollerResponder represents a sequence of responses for a long-running operation.
// Any non-terminal responses are replayed in the order in which they were added.
// The terminal response, success or error, is always the final response.
type PollerResponder[T any] struct {
	nonTermResps []nonTermResp
	httpStatus   int
	res          *T
	err          *exported.ResponseError
}

// AddNonTerminalResponse adds a non-terminal response to the sequence of responses.
func (p *PollerResponder[T]) AddNonTerminalResponse(httpStatus int, o *AddNonTerminalResponseOptions) {
	p.nonTermResps = append(p.nonTermResps, nonTermResp{httpStatus: httpStatus, status: "InProgress"})
}

// AddPollingError adds an error to the sequence of responses.
// Use this to simulate an error durring polling.
// NOTE: adding this as the first response will cause the Begin* LRO API to return this error.
func (p *PollerResponder[T]) AddPollingError(err error) {
	p.nonTermResps = append(p.nonTermResps, nonTermResp{err: err})
}

// SetTerminalResponse sets the provided value as the successful, terminal response.
func (p *PollerResponder[T]) SetTerminalResponse(httpStatus int, result T, o *SetTerminalResponseOptions) {
	p.httpStatus = httpStatus
	p.res = &result
}

// SetTerminalError sets an *azcore.ResponseError with the specified values as the failed terminal response.
func (p *PollerResponder[T]) SetTerminalError(httpStatus int, errorCode string) {
	p.err = &exported.ResponseError{ErrorCode: errorCode, StatusCode: httpStatus}
}

// AddNonTerminalResponseOptions contains the optional values for PollerResponder[T].AddNonTerminalResponse.
type AddNonTerminalResponseOptions struct {
	// place holder for future optional values
}

// SetTerminalResponseOptions contains the optional values for PollerResponder[T].SetTerminalResponse.
type SetTerminalResponseOptions struct {
	// place holder for future optional values
}

// More returns true if there are more responses for consumption.
// This function is called by the fake server internals.
func (p *PollerResponder[T]) More() bool {
	return len(p.nonTermResps) > 0 || p.err != nil || p.res != nil
}

// Next returns the next response in the sequence (a *http.Response or an error).
// This function is called by the fake server internals.
func (p *PollerResponder[T]) Next(req *http.Request) (*http.Response, error) {
	if len(p.nonTermResps) > 0 {
		resp := p.nonTermResps[0]
		p.nonTermResps = p.nonTermResps[1:]

		if resp.err != nil {
			return nil, shared.NonRetriableError(resp.err)
		}

		content := ResponseContent{
			HTTPStatus: resp.httpStatus,
			Header:     http.Header{},
		}
		httpResp, err := NewResponse(content, req)
		if err != nil {
			return nil, shared.NonRetriableError(err)
		}
		httpResp.Header.Set(shared.HeaderFakePollerStatus, resp.status)

		if resp.retryAfter > 0 {
			httpResp.Header.Add(shared.HeaderRetryAfter, strconv.Itoa(resp.retryAfter))
		}

		return httpResp, nil
	}

	if p.err != nil {
		respErr := p.err
		rawResp, err := newErrorResponse(p.err.StatusCode, p.err.ErrorCode, req)
		if err != nil {
			return nil, shared.NonRetriableError(err)
		}
		respErr.RawResponse = rawResp
		p.err = nil
		return nil, shared.NonRetriableError(respErr)
	} else if p.res != nil {
		body, err := json.Marshal(*p.res)
		if err != nil {
			return nil, shared.NonRetriableError(err)
		}
		p.res = nil
		content := ResponseContent{
			HTTPStatus: p.httpStatus,
			Header:     http.Header{},
		}
		resp, err := NewResponse(content, req)
		if err != nil {
			return nil, shared.NonRetriableError(err)
		}
		httpResp := SetResponseBody(resp, body, shared.ContentTypeAppJSON)
		httpResp.Header.Set(shared.HeaderFakePollerStatus, "Succeeded")
		return httpResp, nil
	} else {
		return nil, shared.NonRetriableError(fmt.Errorf("%T has no terminal response", p))
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ResponseContent is used when building the *http.Response.
// This type is used by the fake server internals.
type ResponseContent struct {
	// HTTPStatus is the HTTP status code to use in the response.
	HTTPStatus int

	// Header contains the headers from SetResponseOptions.Header to include in the HTTP response.
	Header http.Header
}

// ResponseOptions contains the optional values for NewResponse().
type ResponseOptions struct {
	// Body is the HTTP response body.
	Body io.ReadCloser

	// ContentType is the value for the Content-Type HTTP header.
	ContentType string
}

type pageResp[T any] struct {
	httpStatus int
	entry      T
}

type nonTermResp struct {
	httpStatus int
	status     string
	retryAfter int
	err        error
}

// SetResponseBody wraps body in a nop-closing bytes reader and assigned it to resp.Body.
// The Content-Type header will be added with the specified value.
func SetResponseBody(resp *http.Response, body []byte, contentType string) *http.Response {
	if l := int64(len(body)); l > 0 {
		resp.Header.Set(shared.HeaderContentType, contentType)
		resp.ContentLength = l
		resp.Body = io.NopCloser(bytes.NewReader(body))
	}
	return resp
}

// NewResponse creates a new *http.Response with the specified content and req as the response's request.
func NewResponse(content ResponseContent, req *http.Request) (*http.Response, error) {
	if content.HTTPStatus == 0 {
		return nil, errors.New("fake: no HTTP status code was specified")
	} else if content.Header == nil {
		content.Header = http.Header{}
	}
	return &http.Response{
		Body:       http.NoBody,
		Header:     content.Header,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    req,
		Status:     fmt.Sprintf("%d %s", content.HTTPStatus, http.StatusText(content.HTTPStatus)),
		StatusCode: content.HTTPStatus,
	}, nil
}

func newErrorResponse(statusCode int, errorCode string, req *http.Request) (*http.Response, error) {
	content := ResponseContent{
		HTTPStatus: statusCode,
		Header:     http.Header{},
	}
	resp, err := NewResponse(content, req)
	if err != nil {
		return nil, err
	}
	resp.Header.Set(shared.HeaderXMSErrorCode, errorCode)
	return resp, nil
}
