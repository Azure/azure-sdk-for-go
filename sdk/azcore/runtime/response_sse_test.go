// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package runtime_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

type sseResponseTestBody struct {
	io.Reader
	closeErr error
	reads    int
	closes   int
}

func (b *sseResponseTestBody) Read(p []byte) (int, error) {
	b.reads++
	return b.Reader.Read(p)
}

func (b *sseResponseTestBody) Close() error {
	b.closes++
	return b.closeErr
}

func TestSSEResponseIncomingError(t *testing.T) {
	connectionErr := errors.New("connection failed")
	for _, tc := range []struct {
		name string
		err  error
	}{
		{name: "connection failure", err: connectionErr},
		{name: "EOF", err: io.EOF},
		{name: "wrapped EOF", err: fmt.Errorf("receiving headers: %w", io.EOF)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &sseResponseTestBody{Reader: strings.NewReader("unused"), closeErr: errors.New("body must not be closed again")}
			got, err := runtime.SSEResponse(&http.Response{Body: body}, tc.err, http.StatusOK)
			require.Nil(t, got)
			require.Zero(t, body.reads)
			require.Zero(t, body.closes)
			if errors.Is(tc.err, io.EOF) {
				require.ErrorIs(t, err, io.ErrUnexpectedEOF)
				require.NotErrorIs(t, err, io.EOF)
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.Same(t, tc.err, err)
			}
		})
	}
}

func TestSSEResponseOwnership(t *testing.T) {
	for _, tc := range []struct {
		name        string
		status      int
		contentType string
		wantErr     string
	}{
		{name: "valid SSE", status: http.StatusOK, contentType: "text/event-stream"},
		{name: "parameterized SSE", status: http.StatusOK, contentType: "text/event-stream; charset=utf-8"},
		{name: "completion", status: http.StatusNoContent},
		{name: "HTTP error", status: http.StatusBadRequest},
		{name: "wrong content type", status: http.StatusOK, contentType: "application/json", wantErr: "unexpected SSE Content-Type"},
		{name: "malformed content type", status: http.StatusOK, contentType: "text/event-stream; charset", wantErr: "invalid SSE Content-Type"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &sseResponseTestBody{Reader: strings.NewReader(`{"error":{"code":"BadRequest"}}`)}
			resp := &http.Response{
				StatusCode: tc.status,
				Header:     http.Header{"Content-Type": []string{tc.contentType}},
				Body:       body,
			}
			got, err := runtime.SSEResponse(resp, nil, http.StatusOK)
			if tc.status == http.StatusOK && tc.wantErr == "" {
				require.NoError(t, err)
				require.Same(t, body, got)
				require.Zero(t, body.reads)
				require.Zero(t, body.closes)
				require.NoError(t, got.Close())
			} else {
				require.Nil(t, got)
				require.Equal(t, 1, body.closes)
				if tc.status == http.StatusNoContent {
					require.ErrorIs(t, err, io.EOF)
					require.Zero(t, body.reads)
				} else if tc.wantErr != "" {
					require.ErrorContains(t, err, tc.wantErr)
					require.Zero(t, body.reads)
				} else {
					var responseErr *azcore.ResponseError
					require.ErrorAs(t, err, &responseErr)
					require.Equal(t, tc.status, responseErr.StatusCode)
					require.Equal(t, "BadRequest", responseErr.ErrorCode)
				}
			}
		})
	}
}

func TestSSEResponseIgnoresCleanupError(t *testing.T) {
	closeErr := errors.New("close failed")
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{name: "completion close failure", err: closeErr, status: http.StatusNoContent},
		{name: "completion close EOF", err: io.EOF, status: http.StatusNoContent},
		{name: "completion wrapped close EOF", err: fmt.Errorf("close: %w", io.EOF), status: http.StatusNoContent},
		{name: "validation close failure", err: closeErr, status: http.StatusOK},
		{name: "validation close EOF", err: io.EOF, status: http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &sseResponseTestBody{Reader: strings.NewReader("unused"), closeErr: tc.err}
			got, err := runtime.SSEResponse(&http.Response{
				StatusCode: tc.status,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       body,
			}, nil, http.StatusOK)
			require.Nil(t, got)
			require.Error(t, err)
			if tc.status == http.StatusOK {
				require.NotErrorIs(t, err, io.EOF)
				require.NotErrorIs(t, err, tc.err)
				require.ErrorContains(t, err, "unexpected SSE Content-Type")
			} else {
				require.ErrorIs(t, err, io.EOF)
			}
			require.Zero(t, body.reads)
			require.Equal(t, 1, body.closes)
		})
	}
}

func TestSSEResponseRejectedBodyEOF(t *testing.T) {
	readErr := fmt.Errorf("reading response: %w", io.EOF)
	body := &sseResponseTestBody{Reader: iotest.ErrReader(readErr)}
	got, err := runtime.SSEResponse(&http.Response{
		StatusCode: http.StatusCreated,
		Header:     make(http.Header),
		Body:       body,
	}, nil, http.StatusOK)
	require.Nil(t, got)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.NotErrorIs(t, err, io.EOF)
	require.ErrorContains(t, err, readErr.Error())
	require.Equal(t, 1, body.closes)
}

func TestSSEResponseStatusCodes(t *testing.T) {
	for _, tc := range []struct {
		name        string
		status      int
		statusCodes []int
		wantSuccess bool
	}{
		{name: "authored 201", status: http.StatusCreated, statusCodes: []int{http.StatusCreated}, wantSuccess: true},
		{name: "multiple authored codes", status: http.StatusAccepted, statusCodes: []int{http.StatusOK, http.StatusAccepted}, wantSuccess: true},
		{name: "unexpected 200", status: http.StatusOK, statusCodes: []int{http.StatusCreated}},
		{name: "unexpected 201", status: http.StatusCreated, statusCodes: []int{http.StatusOK}},
		{name: "no implicit 200", status: http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &sseResponseTestBody{Reader: strings.NewReader("{}")}
			got, err := runtime.SSEResponse(&http.Response{
				StatusCode: tc.status,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
				Body:       body,
			}, nil, tc.statusCodes...)
			if tc.wantSuccess {
				require.NoError(t, err)
				require.Same(t, body, got)
				require.Zero(t, body.reads)
				require.Zero(t, body.closes)
				require.NoError(t, got.Close())
			} else {
				require.Nil(t, got)
				var responseErr *azcore.ResponseError
				require.ErrorAs(t, err, &responseErr)
				require.Equal(t, tc.status, responseErr.StatusCode)
				require.Equal(t, 1, body.closes)
			}
		})
	}
}
