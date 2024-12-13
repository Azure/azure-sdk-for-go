//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"net/http"
	"net/url"
	"regexp"
)

// SoftDeletedResourceServer is a fake server for instances of the armcompute.SoftDeletedResourceClient type.
type SoftDeletedResourceServer struct {
	// NewListByArtifactNamePager is the fake for method SoftDeletedResourceClient.NewListByArtifactNamePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByArtifactNamePager func(resourceGroupName string, galleryName string, artifactType string, artifactName string, options *armcompute.SoftDeletedResourceClientListByArtifactNameOptions) (resp azfake.PagerResponder[armcompute.SoftDeletedResourceClientListByArtifactNameResponse])
}

// NewSoftDeletedResourceServerTransport creates a new instance of SoftDeletedResourceServerTransport with the provided implementation.
// The returned SoftDeletedResourceServerTransport instance is connected to an instance of armcompute.SoftDeletedResourceClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSoftDeletedResourceServerTransport(srv *SoftDeletedResourceServer) *SoftDeletedResourceServerTransport {
	return &SoftDeletedResourceServerTransport{
		srv:                        srv,
		newListByArtifactNamePager: newTracker[azfake.PagerResponder[armcompute.SoftDeletedResourceClientListByArtifactNameResponse]](),
	}
}

// SoftDeletedResourceServerTransport connects instances of armcompute.SoftDeletedResourceClient to instances of SoftDeletedResourceServer.
// Don't use this type directly, use NewSoftDeletedResourceServerTransport instead.
type SoftDeletedResourceServerTransport struct {
	srv                        *SoftDeletedResourceServer
	newListByArtifactNamePager *tracker[azfake.PagerResponder[armcompute.SoftDeletedResourceClientListByArtifactNameResponse]]
}

// Do implements the policy.Transporter interface for SoftDeletedResourceServerTransport.
func (s *SoftDeletedResourceServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "SoftDeletedResourceClient.NewListByArtifactNamePager":
		resp, err = s.dispatchNewListByArtifactNamePager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SoftDeletedResourceServerTransport) dispatchNewListByArtifactNamePager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByArtifactNamePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByArtifactNamePager not implemented")}
	}
	newListByArtifactNamePager := s.newListByArtifactNamePager.get(req)
	if newListByArtifactNamePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/galleries/(?P<galleryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/softDeletedArtifactTypes/(?P<artifactType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/artifacts/(?P<artifactName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/versions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		galleryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("galleryName")])
		if err != nil {
			return nil, err
		}
		artifactTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("artifactType")])
		if err != nil {
			return nil, err
		}
		artifactNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("artifactName")])
		if err != nil {
			return nil, err
		}
		resp := s.srv.NewListByArtifactNamePager(resourceGroupNameParam, galleryNameParam, artifactTypeParam, artifactNameParam, nil)
		newListByArtifactNamePager = &resp
		s.newListByArtifactNamePager.add(req, newListByArtifactNamePager)
		server.PagerResponderInjectNextLinks(newListByArtifactNamePager, req, func(page *armcompute.SoftDeletedResourceClientListByArtifactNameResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByArtifactNamePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByArtifactNamePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByArtifactNamePager) {
		s.newListByArtifactNamePager.remove(req)
	}
	return resp, nil
}