// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azblob

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ListBlobsFlatSegmentResponsePager provides iteration over ListBlobsFlatSegmentResponse pages.
type ListBlobsFlatSegmentResponsePager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// Page returns the current ListBlobsFlatSegmentResponseResponse.
	PageResponse() *ListBlobsFlatSegmentResponseResponse

	// Err returns the last error encountered while paging.
	Err() error
}

type listBlobsFlatSegmentResponseHandleResponse func(*azcore.Response) (*ListBlobsFlatSegmentResponseResponse, error)

type listBlobsFlatSegmentResponseAdvancePage func(*ListBlobsFlatSegmentResponseResponse) (*azcore.Request, error)

type listBlobsFlatSegmentResponsePager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// contains the pending request
	request *azcore.Request
	// callback for handling the HTTP response
	responder listBlobsFlatSegmentResponseHandleResponse
	// callback for advancing to the next page
	advancer listBlobsFlatSegmentResponseAdvancePage
	// contains the current response
	current *ListBlobsFlatSegmentResponseResponse
	// any error encountered
	err error
}

func (p *listBlobsFlatSegmentResponsePager) Err() error {
	return p.err
}

func (p *listBlobsFlatSegmentResponsePager) NextPage(ctx context.Context) bool {
	if p.current != nil {
		if p.current.EnumerationResults.NextMarker == nil || len(*p.current.EnumerationResults.NextMarker) == 0 {
			return false
		}
		req, err := p.advancer(p.current)
		if err != nil {
			p.err = err
			return false
		}
		p.request = req
	}
	resp, err := p.pipeline.Do(ctx, p.request)
	if err != nil {
		p.err = err
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *listBlobsFlatSegmentResponsePager) PageResponse() *ListBlobsFlatSegmentResponseResponse {
	return p.current
}

// ListBlobsHierarchySegmentResponsePager provides iteration over ListBlobsHierarchySegmentResponse pages.
type ListBlobsHierarchySegmentResponsePager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// Page returns the current ListBlobsHierarchySegmentResponseResponse.
	PageResponse() *ListBlobsHierarchySegmentResponseResponse

	// Err returns the last error encountered while paging.
	Err() error
}

type listBlobsHierarchySegmentResponseHandleResponse func(*azcore.Response) (*ListBlobsHierarchySegmentResponseResponse, error)

type listBlobsHierarchySegmentResponseAdvancePage func(*ListBlobsHierarchySegmentResponseResponse) (*azcore.Request, error)

type listBlobsHierarchySegmentResponsePager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// contains the pending request
	request *azcore.Request
	// callback for handling the HTTP response
	responder listBlobsHierarchySegmentResponseHandleResponse
	// callback for advancing to the next page
	advancer listBlobsHierarchySegmentResponseAdvancePage
	// contains the current response
	current *ListBlobsHierarchySegmentResponseResponse
	// any error encountered
	err error
}

func (p *listBlobsHierarchySegmentResponsePager) Err() error {
	return p.err
}

func (p *listBlobsHierarchySegmentResponsePager) NextPage(ctx context.Context) bool {
	if p.current != nil {
		if p.current.EnumerationResults.NextMarker == nil || len(*p.current.EnumerationResults.NextMarker) == 0 {
			return false
		}
		req, err := p.advancer(p.current)
		if err != nil {
			p.err = err
			return false
		}
		p.request = req
	}
	resp, err := p.pipeline.Do(ctx, p.request)
	if err != nil {
		p.err = err
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *listBlobsHierarchySegmentResponsePager) PageResponse() *ListBlobsHierarchySegmentResponseResponse {
	return p.current
}

// ListContainersSegmentResponsePager provides iteration over ListContainersSegmentResponse pages.
type ListContainersSegmentResponsePager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// Page returns the current ListContainersSegmentResponseResponse.
	PageResponse() *ListContainersSegmentResponseResponse

	// Err returns the last error encountered while paging.
	Err() error
}

type listContainersSegmentResponseHandleResponse func(*azcore.Response) (*ListContainersSegmentResponseResponse, error)

type listContainersSegmentResponseAdvancePage func(*ListContainersSegmentResponseResponse) (*azcore.Request, error)

type listContainersSegmentResponsePager struct {
	// the pipeline for making the request
	pipeline azcore.Pipeline
	// contains the pending request
	request *azcore.Request
	// callback for handling the HTTP response
	responder listContainersSegmentResponseHandleResponse
	// callback for advancing to the next page
	advancer listContainersSegmentResponseAdvancePage
	// contains the current response
	current *ListContainersSegmentResponseResponse
	// any error encountered
	err error
}

func (p *listContainersSegmentResponsePager) Err() error {
	return p.err
}

func (p *listContainersSegmentResponsePager) NextPage(ctx context.Context) bool {
	if p.current != nil {
		if p.current.EnumerationResults.NextMarker == nil || len(*p.current.EnumerationResults.NextMarker) == 0 {
			return false
		}
		req, err := p.advancer(p.current)
		if err != nil {
			p.err = err
			return false
		}
		p.request = req
	}
	resp, err := p.pipeline.Do(ctx, p.request)
	if err != nil {
		p.err = err
		return false
	}
	result, err := p.responder(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *listContainersSegmentResponsePager) PageResponse() *ListContainersSegmentResponseResponse {
	return p.current
}
