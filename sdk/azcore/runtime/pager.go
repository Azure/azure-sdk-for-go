//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// PageProcessor contains the required data for constructing a Pager.
type PageProcessor[T any] struct {
	// Do sends the request to fetch the next page.
	Do func(*policy.Request) (*http.Response, error)

	// More returns a boolean indicating if there are more pages to fetch.
	// It uses the provided page to make the determination.
	More func(T) bool

	// Requester creates the request to fetch the first page.
	Requester func(context.Context) (*policy.Request, error)

	// Advancer creates the request to fetch subsequent pages.
	Advancer func(context.Context, T) (*policy.Request, error)

	// Responder handles the responses when fetching pages.
	Responder func(*http.Response) (T, error)
}

// Pager provides operations for iterating over paged responses.
type Pager[T any] struct {
	current   T
	handler   PageProcessor[T]
	lroSecond bool
}

// NewPager creates an instance of Pager using the specified PageProcessor.
func NewPager[T any](pageHandler PageProcessor[T]) *Pager[T] {
	return &Pager[T]{
		handler: pageHandler,
	}
}

// More returns true if there are more pages to retrieve.
func (p *Pager[T]) More() bool {
	if !reflect.ValueOf(p.current).IsZero() {
		return p.handler.More(p.current)
	}
	return true
}

// NextPage advances the pager to the next page.
func (p *Pager[T]) NextPage(ctx context.Context) (T, error) {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() && !p.lroSecond {
		// we get here if it's an LRO-pager, we already have the first page
		p.lroSecond = true
		return p.current, nil
	} else if !reflect.ValueOf(p.current).IsZero() {
		if !p.handler.More(p.current) {
			return *new(T), errors.New("no more pages")
		}
		req, err = p.handler.Advancer(ctx, p.current)
	} else {
		// non-LRO case, so skip the LRO branch above
		p.lroSecond = true
		req, err = p.handler.Requester(ctx)
	}
	if err != nil {
		return *new(T), err
	}
	resp, err := p.handler.Do(req)
	if err != nil {
		return *new(T), err
	}
	if !shared.HasStatusCode(resp, http.StatusOK) {
		return *new(T), shared.NewResponseError(resp)
	}
	result, err := p.handler.Responder(resp)
	if err != nil {
		return *new(T), err
	}
	p.current = result
	return p.current, nil
}
