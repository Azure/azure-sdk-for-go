//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package fake provides the building blocks for fake servers.
// This includes fakes for authentication, API responses, and more.
package fake

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
)

// TokenCredential is a fake credential that implements the azcore.TokenCredential interface.
type TokenCredential struct {
	err error
}

// SetError sets the specified error to be returned from GetToken().
// Use this to simulate an error during authentication.
func (t *TokenCredential) SetError(err error) {
	t.err = errorinfo.NonRetriableError(err)
}

// GetToken implements the azcore.TokenCredential for the TokenCredential type.
func (t *TokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if t.err != nil {
		return azcore.AccessToken{}, errorinfo.NonRetriableError(t.err)
	}
	return azcore.AccessToken{Token: "fake_token", ExpiresOn: time.Now().Add(24 * time.Hour)}, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Responder represents a scalar response.
type Responder[T any] exported.Responder[T]

// SetResponse sets the specified value to be returned.
//   - httpStatus is the HTTP status code to be returned
//   - resp is the response to be returned
//   - o contains optional values, pass nil to accept the defaults
func (r *Responder[T]) SetResponse(httpStatus int, resp T, o *SetResponseOptions) {
	(*exported.Responder[T])(r).SetResponse(httpStatus, resp, o)
}

// SetResponseOptions contains the optional values for Responder[T].SetResponse.
type SetResponseOptions = exported.SetResponseOptions

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ErrorResponder represents a scalar error response.
type ErrorResponder exported.ErrorResponder

// SetError sets the specified error to be returned.
// Use SetResponseError for returning an *azcore.ResponseError.
func (e *ErrorResponder) SetError(err error) {
	(*exported.ErrorResponder)(e).SetError(err)
}

// SetResponseError sets an *azcore.ResponseError with the specified values to be returned.
//   - errorCode is the value to be used in the ResponseError.Code field
//   - httpStatus is the HTTP status code
func (e *ErrorResponder) SetResponseError(httpStatus int, errorCode string) {
	(*exported.ErrorResponder)(e).SetResponseError(httpStatus, errorCode)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// PagerResponder represents a sequence of paged responses.
// Responses are consumed in the order in which they were added.
// If no pages or errors have been added, calls to Pager[T].NextPage
// will return an error.
type PagerResponder[T any] exported.PagerResponder[T]

// AddPage adds a page to the sequence of respones.
//   - page is the response page to be added
//   - o contains optional values, pass nil to accept the defaults
func (p *PagerResponder[T]) AddPage(httpStatus int, page T, o *AddPageOptions) {
	(*exported.PagerResponder[T])(p).AddPage(httpStatus, page, o)
}

// AddError adds an error to the sequence of responses.
// The error is returned from the call to runtime.Pager[T].NextPage().
func (p *PagerResponder[T]) AddError(err error) {
	(*exported.PagerResponder[T])(p).AddError(err)
}

// AddResponseError adds an *azcore.ResponseError to the sequence of responses.
// The error is returned from the call to runtime.Pager[T].NextPage().
func (p *PagerResponder[T]) AddResponseError(httpStatus int, errorCode string) {
	(*exported.PagerResponder[T])(p).AddResponseError(httpStatus, errorCode)
}

// AddPageOptions contains the optional values for PagerResponder[T].AddPage.
type AddPageOptions = exported.AddPageOptions

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

// PollerResponder represents a sequence of responses for a long-running operation.
// Any non-terminal responses are consumed in the order in which they were added.
// The terminal response, success or error, is always the final response.
// If no responses or errors have been added, the following method calls on Poller[T]
// will return an error: PollUntilDone, Poll, Result.
type PollerResponder[T any] exported.PollerResponder[T]

// AddNonTerminalResponse adds a non-terminal response to the sequence of responses.
func (p *PollerResponder[T]) AddNonTerminalResponse(httpStatus int, o *AddNonTerminalResponseOptions) {
	(*exported.PollerResponder[T])(p).AddNonTerminalResponse(httpStatus, o)
}

// AddPollingError adds an error to the sequence of responses.
// Use this to simulate an error durring polling.
// NOTE: adding this as the first response will cause the Begin* LRO API to return this error.
func (p *PollerResponder[T]) AddPollingError(err error) {
	(*exported.PollerResponder[T])(p).AddPollingError(err)
}

// SetTerminalResponse sets the provided value as the successful, terminal response.
func (p *PollerResponder[T]) SetTerminalResponse(httpStatus int, result T, o *SetTerminalResponseOptions) {
	(*exported.PollerResponder[T])(p).SetTerminalResponse(httpStatus, result, o)
}

// SetTerminalError sets an *azcore.ResponseError with the specified values as the failed terminal response.
func (p *PollerResponder[T]) SetTerminalError(httpStatus int, errorCode string) {
	(*exported.PollerResponder[T])(p).SetTerminalError(httpStatus, errorCode)
}

// AddNonTerminalResponseOptions contains the optional values for PollerResponder[T].AddNonTerminalResponse.
type AddNonTerminalResponseOptions = exported.AddNonTerminalResponseOptions

// SetTerminalResponseOptions contains the optional values for PollerResponder[T].SetTerminalResponse.
type SetTerminalResponseOptions = exported.SetTerminalResponseOptions
