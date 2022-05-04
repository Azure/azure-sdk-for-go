// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func newFakeResponse(statusCode int, status string, contents string) *http.Response {
	var body io.ReadCloser = http.NoBody

	if contents != "" {
		body = &FakeReader{
			Reader: bytes.NewBufferString(contents),
		}
	}

	return &http.Response{
		Request: &http.Request{
			URL: &url.URL{},
		},
		StatusCode: statusCode,
		Status:     status,
		Body:       body,
	}
}

func TestResponseError(t *testing.T) {
	resp := newFakeResponse(http.StatusConflict, "statusString", "")
	require.Contains(t, NewResponseError(resp).Error(), "statusString")

	resp = newFakeResponse(http.StatusConflict, "statusString", "contents")
	require.Contains(t, NewResponseError(resp).Error(), "statusString")

	resp = newFakeResponse(http.StatusBadGateway, "statusString", "<Error><Code>401</Code><Detail>Manage,EntityRead claims required for this operation.</Detail></Error>")
	err := NewResponseError(resp)

	re, ok := err.(*azcore.ResponseError)
	require.True(t, ok)

	require.Contains(t, re.Error(), "statusString")
	require.EqualValues(t, http.StatusBadGateway, re.StatusCode)
}

type FakeReader struct {
	io.Reader
	closed   bool
	closeErr error
}

func (f *FakeReader) Close() error {
	f.closed = true
	return f.closeErr
}

func TestCloseRes(t *testing.T) {
	reader := strings.NewReader("hello")
	body := &FakeReader{Reader: reader}

	CloseRes(context.Background(), &http.Response{
		Body: body,
	})

	// check that we're at EOF (ie, was fully drained)
	n, err := reader.Read(nil)
	require.EqualValues(t, 0, n)
	require.EqualError(t, err, io.EOF.Error())

	require.True(t, body.closed)
}
