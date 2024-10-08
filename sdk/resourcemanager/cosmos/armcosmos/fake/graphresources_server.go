//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v4"
	"net/http"
	"net/url"
	"regexp"
)

// GraphResourcesServer is a fake server for instances of the armcosmos.GraphResourcesClient type.
type GraphResourcesServer struct {
	// BeginCreateUpdateGraph is the fake for method GraphResourcesClient.BeginCreateUpdateGraph
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateUpdateGraph func(ctx context.Context, resourceGroupName string, accountName string, graphName string, createUpdateGraphParameters armcosmos.GraphResourceCreateUpdateParameters, options *armcosmos.GraphResourcesClientBeginCreateUpdateGraphOptions) (resp azfake.PollerResponder[armcosmos.GraphResourcesClientCreateUpdateGraphResponse], errResp azfake.ErrorResponder)

	// BeginDeleteGraphResource is the fake for method GraphResourcesClient.BeginDeleteGraphResource
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDeleteGraphResource func(ctx context.Context, resourceGroupName string, accountName string, graphName string, options *armcosmos.GraphResourcesClientBeginDeleteGraphResourceOptions) (resp azfake.PollerResponder[armcosmos.GraphResourcesClientDeleteGraphResourceResponse], errResp azfake.ErrorResponder)

	// GetGraph is the fake for method GraphResourcesClient.GetGraph
	// HTTP status codes to indicate success: http.StatusOK
	GetGraph func(ctx context.Context, resourceGroupName string, accountName string, graphName string, options *armcosmos.GraphResourcesClientGetGraphOptions) (resp azfake.Responder[armcosmos.GraphResourcesClientGetGraphResponse], errResp azfake.ErrorResponder)

	// NewListGraphsPager is the fake for method GraphResourcesClient.NewListGraphsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListGraphsPager func(resourceGroupName string, accountName string, options *armcosmos.GraphResourcesClientListGraphsOptions) (resp azfake.PagerResponder[armcosmos.GraphResourcesClientListGraphsResponse])
}

// NewGraphResourcesServerTransport creates a new instance of GraphResourcesServerTransport with the provided implementation.
// The returned GraphResourcesServerTransport instance is connected to an instance of armcosmos.GraphResourcesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewGraphResourcesServerTransport(srv *GraphResourcesServer) *GraphResourcesServerTransport {
	return &GraphResourcesServerTransport{
		srv:                      srv,
		beginCreateUpdateGraph:   newTracker[azfake.PollerResponder[armcosmos.GraphResourcesClientCreateUpdateGraphResponse]](),
		beginDeleteGraphResource: newTracker[azfake.PollerResponder[armcosmos.GraphResourcesClientDeleteGraphResourceResponse]](),
		newListGraphsPager:       newTracker[azfake.PagerResponder[armcosmos.GraphResourcesClientListGraphsResponse]](),
	}
}

// GraphResourcesServerTransport connects instances of armcosmos.GraphResourcesClient to instances of GraphResourcesServer.
// Don't use this type directly, use NewGraphResourcesServerTransport instead.
type GraphResourcesServerTransport struct {
	srv                      *GraphResourcesServer
	beginCreateUpdateGraph   *tracker[azfake.PollerResponder[armcosmos.GraphResourcesClientCreateUpdateGraphResponse]]
	beginDeleteGraphResource *tracker[azfake.PollerResponder[armcosmos.GraphResourcesClientDeleteGraphResourceResponse]]
	newListGraphsPager       *tracker[azfake.PagerResponder[armcosmos.GraphResourcesClientListGraphsResponse]]
}

// Do implements the policy.Transporter interface for GraphResourcesServerTransport.
func (g *GraphResourcesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "GraphResourcesClient.BeginCreateUpdateGraph":
		resp, err = g.dispatchBeginCreateUpdateGraph(req)
	case "GraphResourcesClient.BeginDeleteGraphResource":
		resp, err = g.dispatchBeginDeleteGraphResource(req)
	case "GraphResourcesClient.GetGraph":
		resp, err = g.dispatchGetGraph(req)
	case "GraphResourcesClient.NewListGraphsPager":
		resp, err = g.dispatchNewListGraphsPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *GraphResourcesServerTransport) dispatchBeginCreateUpdateGraph(req *http.Request) (*http.Response, error) {
	if g.srv.BeginCreateUpdateGraph == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateUpdateGraph not implemented")}
	}
	beginCreateUpdateGraph := g.beginCreateUpdateGraph.get(req)
	if beginCreateUpdateGraph == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/graphs/(?P<graphName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcosmos.GraphResourceCreateUpdateParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
		if err != nil {
			return nil, err
		}
		graphNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("graphName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginCreateUpdateGraph(req.Context(), resourceGroupNameParam, accountNameParam, graphNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateUpdateGraph = &respr
		g.beginCreateUpdateGraph.add(req, beginCreateUpdateGraph)
	}

	resp, err := server.PollerResponderNext(beginCreateUpdateGraph, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		g.beginCreateUpdateGraph.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateUpdateGraph) {
		g.beginCreateUpdateGraph.remove(req)
	}

	return resp, nil
}

func (g *GraphResourcesServerTransport) dispatchBeginDeleteGraphResource(req *http.Request) (*http.Response, error) {
	if g.srv.BeginDeleteGraphResource == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDeleteGraphResource not implemented")}
	}
	beginDeleteGraphResource := g.beginDeleteGraphResource.get(req)
	if beginDeleteGraphResource == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/graphs/(?P<graphName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
		if err != nil {
			return nil, err
		}
		graphNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("graphName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginDeleteGraphResource(req.Context(), resourceGroupNameParam, accountNameParam, graphNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDeleteGraphResource = &respr
		g.beginDeleteGraphResource.add(req, beginDeleteGraphResource)
	}

	resp, err := server.PollerResponderNext(beginDeleteGraphResource, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		g.beginDeleteGraphResource.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDeleteGraphResource) {
		g.beginDeleteGraphResource.remove(req)
	}

	return resp, nil
}

func (g *GraphResourcesServerTransport) dispatchGetGraph(req *http.Request) (*http.Response, error) {
	if g.srv.GetGraph == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetGraph not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/graphs/(?P<graphName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
	if err != nil {
		return nil, err
	}
	graphNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("graphName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.GetGraph(req.Context(), resourceGroupNameParam, accountNameParam, graphNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GraphResourceGetResults, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GraphResourcesServerTransport) dispatchNewListGraphsPager(req *http.Request) (*http.Response, error) {
	if g.srv.NewListGraphsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListGraphsPager not implemented")}
	}
	newListGraphsPager := g.newListGraphsPager.get(req)
	if newListGraphsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/graphs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
		if err != nil {
			return nil, err
		}
		resp := g.srv.NewListGraphsPager(resourceGroupNameParam, accountNameParam, nil)
		newListGraphsPager = &resp
		g.newListGraphsPager.add(req, newListGraphsPager)
	}
	resp, err := server.PagerResponderNext(newListGraphsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		g.newListGraphsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListGraphsPager) {
		g.newListGraphsPager.remove(req)
	}
	return resp, nil
}
