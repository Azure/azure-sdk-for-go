// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
