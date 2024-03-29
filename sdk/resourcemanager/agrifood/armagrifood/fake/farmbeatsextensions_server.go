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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/agrifood/armagrifood"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// FarmBeatsExtensionsServer is a fake server for instances of the armagrifood.FarmBeatsExtensionsClient type.
type FarmBeatsExtensionsServer struct {
	// Get is the fake for method FarmBeatsExtensionsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, farmBeatsExtensionID string, options *armagrifood.FarmBeatsExtensionsClientGetOptions) (resp azfake.Responder[armagrifood.FarmBeatsExtensionsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method FarmBeatsExtensionsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armagrifood.FarmBeatsExtensionsClientListOptions) (resp azfake.PagerResponder[armagrifood.FarmBeatsExtensionsClientListResponse])
}

// NewFarmBeatsExtensionsServerTransport creates a new instance of FarmBeatsExtensionsServerTransport with the provided implementation.
// The returned FarmBeatsExtensionsServerTransport instance is connected to an instance of armagrifood.FarmBeatsExtensionsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewFarmBeatsExtensionsServerTransport(srv *FarmBeatsExtensionsServer) *FarmBeatsExtensionsServerTransport {
	return &FarmBeatsExtensionsServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armagrifood.FarmBeatsExtensionsClientListResponse]](),
	}
}

// FarmBeatsExtensionsServerTransport connects instances of armagrifood.FarmBeatsExtensionsClient to instances of FarmBeatsExtensionsServer.
// Don't use this type directly, use NewFarmBeatsExtensionsServerTransport instead.
type FarmBeatsExtensionsServerTransport struct {
	srv          *FarmBeatsExtensionsServer
	newListPager *tracker[azfake.PagerResponder[armagrifood.FarmBeatsExtensionsClientListResponse]]
}

// Do implements the policy.Transporter interface for FarmBeatsExtensionsServerTransport.
func (f *FarmBeatsExtensionsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "FarmBeatsExtensionsClient.Get":
		resp, err = f.dispatchGet(req)
	case "FarmBeatsExtensionsClient.NewListPager":
		resp, err = f.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *FarmBeatsExtensionsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if f.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/providers/Microsoft\.AgFoodPlatform/farmBeatsExtensionDefinitions/(?P<farmBeatsExtensionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	farmBeatsExtensionIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("farmBeatsExtensionId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.Get(req.Context(), farmBeatsExtensionIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).FarmBeatsExtension, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FarmBeatsExtensionsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if f.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := f.newListPager.get(req)
	if newListPager == nil {
		qp := req.URL.Query()
		farmBeatsExtensionIDsEscaped := qp["farmBeatsExtensionIds"]
		farmBeatsExtensionIDsParam := make([]string, len(farmBeatsExtensionIDsEscaped))
		for i, v := range farmBeatsExtensionIDsEscaped {
			u, unescapeErr := url.QueryUnescape(v)
			if unescapeErr != nil {
				return nil, unescapeErr
			}
			farmBeatsExtensionIDsParam[i] = u
		}
		farmBeatsExtensionNamesEscaped := qp["farmBeatsExtensionNames"]
		farmBeatsExtensionNamesParam := make([]string, len(farmBeatsExtensionNamesEscaped))
		for i, v := range farmBeatsExtensionNamesEscaped {
			u, unescapeErr := url.QueryUnescape(v)
			if unescapeErr != nil {
				return nil, unescapeErr
			}
			farmBeatsExtensionNamesParam[i] = u
		}
		extensionCategoriesEscaped := qp["extensionCategories"]
		extensionCategoriesParam := make([]string, len(extensionCategoriesEscaped))
		for i, v := range extensionCategoriesEscaped {
			u, unescapeErr := url.QueryUnescape(v)
			if unescapeErr != nil {
				return nil, unescapeErr
			}
			extensionCategoriesParam[i] = u
		}
		publisherIDsEscaped := qp["publisherIds"]
		publisherIDsParam := make([]string, len(publisherIDsEscaped))
		for i, v := range publisherIDsEscaped {
			u, unescapeErr := url.QueryUnescape(v)
			if unescapeErr != nil {
				return nil, unescapeErr
			}
			publisherIDsParam[i] = u
		}
		maxPageSizeUnescaped, err := url.QueryUnescape(qp.Get("$maxPageSize"))
		if err != nil {
			return nil, err
		}
		maxPageSizeParam, err := parseOptional(maxPageSizeUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armagrifood.FarmBeatsExtensionsClientListOptions
		if len(farmBeatsExtensionIDsParam) > 0 || len(farmBeatsExtensionNamesParam) > 0 || len(extensionCategoriesParam) > 0 || len(publisherIDsParam) > 0 || maxPageSizeParam != nil {
			options = &armagrifood.FarmBeatsExtensionsClientListOptions{
				FarmBeatsExtensionIDs:   farmBeatsExtensionIDsParam,
				FarmBeatsExtensionNames: farmBeatsExtensionNamesParam,
				ExtensionCategories:     extensionCategoriesParam,
				PublisherIDs:            publisherIDsParam,
				MaxPageSize:             maxPageSizeParam,
			}
		}
		resp := f.srv.NewListPager(options)
		newListPager = &resp
		f.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armagrifood.FarmBeatsExtensionsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		f.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		f.newListPager.remove(req)
	}
	return resp, nil
}
