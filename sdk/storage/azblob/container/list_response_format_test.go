// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

// acceptTransport records the Accept header of each request and answers with an empty XML listing,
// as the service does when Arrow isn't available.
type acceptTransport struct{ accept []string }

func (a *acceptTransport) Do(req *http.Request) (*http.Response, error) {
	// the generated code may set the header under a non-canonical key
	var accept []string
	for k, v := range req.Header {
		if strings.EqualFold(k, "Accept") {
			accept = append(accept, v...)
		}
	}
	a.accept = append(a.accept, strings.Join(accept, ","))
	body := `<?xml version="1.0" encoding="utf-8"?><EnumerationResults ContainerName="container"><Blobs /><NextMarker /></EnumerationResults>`
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/xml"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}

// nil options select the default response format, Auto, exactly as empty options do.
func TestListBlobsNilOptionsUseDefaultResponseFormat(t *testing.T) {
	newClient := func(t *testing.T, tr *acceptTransport) *container.Client {
		c, err := container.NewClientWithNoCredential("https://account.blob.core.windows.net/container", &container.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
		})
		require.NoError(t, err)
		return c
	}
	isArrow := func(accept string) bool { return strings.Contains(accept, "application/vnd.apache.arrow.stream") }

	for _, tt := range []struct {
		name      string
		flat      *container.ListBlobsFlatOptions
		hierarchy *container.ListBlobsHierarchyOptions
		wantArrow bool
	}{
		{"Nil", nil, nil, true},
		{"Empty", &container.ListBlobsFlatOptions{}, &container.ListBlobsHierarchyOptions{}, true},
		{"XML", &container.ListBlobsFlatOptions{ResponseFormat: container.StorageResponseFormatXML}, &container.ListBlobsHierarchyOptions{ResponseFormat: container.StorageResponseFormatXML}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tr := &acceptTransport{}
			c := newClient(t, tr)
			_, err := c.NewListBlobsFlatPager(tt.flat).NextPage(context.Background())
			require.NoError(t, err)
			_, err = c.NewListBlobsHierarchyPager("/", tt.hierarchy).NextPage(context.Background())
			require.NoError(t, err)
			require.Len(t, tr.accept, 2)
			require.Equal(t, tt.wantArrow, isArrow(tr.accept[0]), "flat Accept %q", tr.accept[0])
			require.Equal(t, tt.wantArrow, isArrow(tr.accept[1]), "hierarchy Accept %q", tr.accept[1])
		})
	}
}
