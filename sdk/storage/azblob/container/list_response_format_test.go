// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
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

// formatTransport answers list requests with one blob in Arrow format, or in XML format echoing
// the request the way the service does.
type formatTransport struct {
	t     *testing.T
	arrow bool
}

func (f *formatTransport) Do(req *http.Request) (*http.Response, error) {
	if f.arrow {
		schema := arrow.NewSchema([]arrow.Field{{Name: "Name", Type: arrow.BinaryTypes.String}}, nil)
		var buf bytes.Buffer
		w := ipc.NewWriter(&buf, ipc.WithSchema(schema), ipc.WithAllocator(memory.DefaultAllocator))
		b := array.NewRecordBuilder(memory.DefaultAllocator, schema)
		b.Field(0).(*array.StringBuilder).Append("blob")
		rec := b.NewRecordBatch()
		require.NoError(f.t, w.Write(rec))
		rec.Release()
		b.Release()
		require.NoError(f.t, w.Close())
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/vnd.apache.arrow.stream"}}, Body: io.NopCloser(&buf), Request: req}, nil
	}
	q := req.URL.Query()
	delimiter := ""
	if d := q.Get("delimiter"); d != "" {
		delimiter = "<Delimiter>" + d + "</Delimiter>"
	}
	body := `<?xml version="1.0" encoding="utf-8"?><EnumerationResults ServiceEndpoint="https://account.blob.core.windows.net/" ContainerName="container">` +
		`<Prefix>` + q.Get("prefix") + `</Prefix><Marker>` + q.Get("marker") + `</Marker><MaxResults>` + q.Get("maxresults") + `</MaxResults>` + delimiter +
		`<Blobs><Blob><Name>blob</Name><Properties /></Blob></Blobs><NextMarker /></EnumerationResults>`
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/xml"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
}

// An Arrow listing reports the same container and request values as an XML listing. The service
// endpoint isn't in an Arrow response and can't be derived reliably, so it stays nil.
func TestListBlobsArrowResponseReportsRequestValues(t *testing.T) {
	list := func(t *testing.T, arrowFormat bool) (container.ListBlobsFlatResponse, container.ListBlobsHierarchyResponse) {
		c, err := container.NewClientWithNoCredential("https://account.blob.core.windows.net/container?sv=2026-12-06&sig=x", &container.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: &formatTransport{t: t, arrow: arrowFormat}, Retry: policy.RetryOptions{MaxRetries: -1}},
		})
		require.NoError(t, err)
		flat, err := c.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{Prefix: to.Ptr("p"), Marker: to.Ptr("m"), MaxResults: to.Ptr[int32](5)}).NextPage(context.Background())
		require.NoError(t, err)
		hier, err := c.NewListBlobsHierarchyPager("/", &container.ListBlobsHierarchyOptions{Prefix: to.Ptr("p"), Marker: to.Ptr("m"), MaxResults: to.Ptr[int32](5)}).NextPage(context.Background())
		require.NoError(t, err)
		return flat, hier
	}
	xmlFlat, xmlHier := list(t, false)
	arrowFlat, arrowHier := list(t, true)

	require.Equal(t, "container", *arrowFlat.ContainerName)
	require.Nil(t, arrowFlat.ServiceEndpoint)
	require.Nil(t, arrowHier.ServiceEndpoint)
	for _, tt := range []struct {
		name       string
		xml, arrow any
	}{
		{"flat ContainerName", xmlFlat.ContainerName, arrowFlat.ContainerName},
		{"flat Prefix", xmlFlat.Prefix, arrowFlat.Prefix},
		{"flat Marker", xmlFlat.Marker, arrowFlat.Marker},
		{"flat MaxResults", xmlFlat.MaxResults, arrowFlat.MaxResults},
		{"hierarchy ContainerName", xmlHier.ContainerName, arrowHier.ContainerName},
		{"hierarchy Prefix", xmlHier.Prefix, arrowHier.Prefix},
		{"hierarchy Marker", xmlHier.Marker, arrowHier.Marker},
		{"hierarchy MaxResults", xmlHier.MaxResults, arrowHier.MaxResults},
		{"hierarchy Delimiter", xmlHier.Delimiter, arrowHier.Delimiter},
	} {
		require.NotNil(t, tt.xml, tt.name)
		require.Equal(t, tt.xml, tt.arrow, tt.name)
	}
	require.Len(t, arrowFlat.Segment.BlobItems, 1)
	require.Equal(t, "blob", *arrowFlat.Segment.BlobItems[0].Name)
}
