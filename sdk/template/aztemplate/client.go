// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/template/aztemplate/internal"
)

// ClientOptions contains optional parameters for NewClient
type ClientOptions struct {
	azcore.ClientOptions
}

// Client is the client to interact with.
// Don't use this type directly, use NewClient() instead.
type Client struct {
	client *internal.TemplateClient
	pl     runtime.Pipeline
}

// NewClient returns a pointer to a Client
func NewClient(cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	pl := runtime.NewPipeline(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(cred, []string{"service_scope"}, nil),
		},
	}, &options.ClientOptions)

	return &Client{client: internal.NewTemplateClient(pl), pl: pl}, nil

}

// SomeServiceActionOptions contains the optional values for the Client.SomeServiceAction method.
type SomeServiceActionOptions struct {
	// OptionalValue is some optional value to be sent to the service.  nil means nothing is sent.
	OptionalValue *string
}

// ClientSomeServiceActionResponse contains the response from method Client.SomeServiceAction.
type ClientSomeServiceActionResponse struct {
	Value *string
}

// SomeServiceAction does some service action
func (c *Client) SomeServiceAction(ctx context.Context, options *SomeServiceActionOptions) (ClientSomeServiceActionResponse, error) {
	resp, err := c.client.SomeAPI(ctx, nil)
	if err != nil {
		return ClientSomeServiceActionResponse{}, err
	}
	return ClientSomeServiceActionResponse{Value: resp.Value}, nil
}

// ClientListValuesOptions contains the optional values for the Client.NewListValuesPager method.
type ClientListValuesOptions struct {
	// PerPage is the optional number of items to return per page.
	PerPage *int32
}

// ClientListValuesResponse contains the response from method Client.ListValues.
type ClientListValuesResponse struct {
	// Next might be a URL to fetch the next page or a continuation token.
	Next *string

	// Values contains the contents of the page.
	Values []*string
}

// NewListValuesPager creates a pager to iterate over pages of results.
func (c *Client) NewListValuesPager(options *ClientListValuesOptions) *runtime.Pager[ClientListValuesResponse] {
	return runtime.NewPager(runtime.PagingHandler[ClientListValuesResponse]{
		More: func(resp ClientListValuesResponse) bool {
			// inspect resp to see if there are more pages
			return false
		},
		Fetcher: func(ctx context.Context, resp *ClientListValuesResponse) (ClientListValuesResponse, error) {
			// use resp to construct the request to fetch the next page
			return ClientListValuesResponse{}, nil
		},
	})
}

// ClientBeginLongRunningOperationOptions contains the optional values for the Client.BeginLongRunningOperation method.
type ClientBeginLongRunningOperationOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ClientLongRunningOperationResponse contains the response from method Client.LongRunningOperation.
type ClientLongRunningOperationResponse struct {
	Value *string
}

// BeginLongRunningOperation is a long-running operation that can take several seconds to complete.
func (c *Client) BeginLongRunningOperation(ctx context.Context, options *ClientBeginLongRunningOperationOptions) (*runtime.Poller[ClientLongRunningOperationResponse], error) {
	if options == nil {
		options = &ClientBeginLongRunningOperationOptions{}
	}

	if options.ResumeToken != "" {
		return runtime.NewPollerFromResumeToken[ClientLongRunningOperationResponse](options.ResumeToken, c.pl, nil)
	}

	// start the LRO
	req, err := runtime.NewRequest(ctx, http.MethodPost, "the LRO URL")
	if err != nil {
		return nil, err
	}

	resp, err := c.pl.Do(req)
	if err != nil {
		return nil, err
	}
	return runtime.NewPoller[ClientLongRunningOperationResponse](resp, c.pl, nil)
}
