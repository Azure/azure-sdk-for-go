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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"net/http"
	"net/url"
	"regexp"
)

// ImportPipelinesServer is a fake server for instances of the armcontainerregistry.ImportPipelinesClient type.
type ImportPipelinesServer struct {
	// BeginCreate is the fake for method ImportPipelinesClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, importPipelineCreateParameters armcontainerregistry.ImportPipeline, options *armcontainerregistry.ImportPipelinesClientBeginCreateOptions) (resp azfake.PollerResponder[armcontainerregistry.ImportPipelinesClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method ImportPipelinesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *armcontainerregistry.ImportPipelinesClientBeginDeleteOptions) (resp azfake.PollerResponder[armcontainerregistry.ImportPipelinesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ImportPipelinesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *armcontainerregistry.ImportPipelinesClientGetOptions) (resp azfake.Responder[armcontainerregistry.ImportPipelinesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ImportPipelinesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, registryName string, options *armcontainerregistry.ImportPipelinesClientListOptions) (resp azfake.PagerResponder[armcontainerregistry.ImportPipelinesClientListResponse])
}

// NewImportPipelinesServerTransport creates a new instance of ImportPipelinesServerTransport with the provided implementation.
// The returned ImportPipelinesServerTransport instance is connected to an instance of armcontainerregistry.ImportPipelinesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewImportPipelinesServerTransport(srv *ImportPipelinesServer) *ImportPipelinesServerTransport {
	return &ImportPipelinesServerTransport{
		srv:          srv,
		beginCreate:  newTracker[azfake.PollerResponder[armcontainerregistry.ImportPipelinesClientCreateResponse]](),
		beginDelete:  newTracker[azfake.PollerResponder[armcontainerregistry.ImportPipelinesClientDeleteResponse]](),
		newListPager: newTracker[azfake.PagerResponder[armcontainerregistry.ImportPipelinesClientListResponse]](),
	}
}

// ImportPipelinesServerTransport connects instances of armcontainerregistry.ImportPipelinesClient to instances of ImportPipelinesServer.
// Don't use this type directly, use NewImportPipelinesServerTransport instead.
type ImportPipelinesServerTransport struct {
	srv          *ImportPipelinesServer
	beginCreate  *tracker[azfake.PollerResponder[armcontainerregistry.ImportPipelinesClientCreateResponse]]
	beginDelete  *tracker[azfake.PollerResponder[armcontainerregistry.ImportPipelinesClientDeleteResponse]]
	newListPager *tracker[azfake.PagerResponder[armcontainerregistry.ImportPipelinesClientListResponse]]
}

// Do implements the policy.Transporter interface for ImportPipelinesServerTransport.
func (i *ImportPipelinesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return i.dispatchToMethodFake(req, method)
}

func (i *ImportPipelinesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if importPipelinesServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = importPipelinesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "ImportPipelinesClient.BeginCreate":
				res.resp, res.err = i.dispatchBeginCreate(req)
			case "ImportPipelinesClient.BeginDelete":
				res.resp, res.err = i.dispatchBeginDelete(req)
			case "ImportPipelinesClient.Get":
				res.resp, res.err = i.dispatchGet(req)
			case "ImportPipelinesClient.NewListPager":
				res.resp, res.err = i.dispatchNewListPager(req)
			default:
				res.err = fmt.Errorf("unhandled API %s", method)
			}

		}
		select {
		case resultChan <- res:
		case <-req.Context().Done():
		}
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (i *ImportPipelinesServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if i.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := i.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/importPipelines/(?P<importPipelineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcontainerregistry.ImportPipeline](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		importPipelineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("importPipelineName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginCreate(req.Context(), resourceGroupNameParam, registryNameParam, importPipelineNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		i.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		i.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		i.beginCreate.remove(req)
	}

	return resp, nil
}

func (i *ImportPipelinesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if i.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := i.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/importPipelines/(?P<importPipelineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		importPipelineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("importPipelineName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginDelete(req.Context(), resourceGroupNameParam, registryNameParam, importPipelineNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		i.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		i.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		i.beginDelete.remove(req)
	}

	return resp, nil
}

func (i *ImportPipelinesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if i.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/importPipelines/(?P<importPipelineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
	if err != nil {
		return nil, err
	}
	importPipelineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("importPipelineName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Get(req.Context(), resourceGroupNameParam, registryNameParam, importPipelineNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ImportPipeline, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *ImportPipelinesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := i.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ContainerRegistry/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/importPipelines`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		resp := i.srv.NewListPager(resourceGroupNameParam, registryNameParam, nil)
		newListPager = &resp
		i.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armcontainerregistry.ImportPipelinesClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		i.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		i.newListPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to ImportPipelinesServerTransport
var importPipelinesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
