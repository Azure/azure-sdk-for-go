// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestResponseError(t *testing.T) {
	// sanity check to make sure my error conforms to azcore's interface
	var err azcore.HTTPResponse = ResponseError{}
	require.NotNil(t, err)

	require.EqualValues(t, "this is now the error message: 409", NewResponseError(nil, &http.Response{
		StatusCode: http.StatusConflict,
		Status:     "this is now the error message",
	}).Error())

	require.EqualValues(t, "inner errors message takes precedence", NewResponseError(errors.New("inner errors message takes precedence"), &http.Response{
		StatusCode: http.StatusConflict,
		Status:     "going to be ignored",
	}).Error())
}

type FakeReader struct {
	io.Reader
	closed bool
}

func (f *FakeReader) Close() error {
	f.closed = true
	return nil
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
