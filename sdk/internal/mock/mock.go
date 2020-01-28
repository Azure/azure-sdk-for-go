// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package mock

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// Server is a wrapper around an httptest.Server.
// The serving of requests is not safe for concurrent use
// which is ok for right now as each test creates is own
// server and doesn't create additional go routines.
type Server struct {
	srv *httptest.Server

	// static is the static response, if this is not nil it's always returned.
	static *mockResponse

	// resp is the queue of responses.  each response is taken from the front.
	resp []mockResponse

	// count tracks the number of requests that have been made.
	count int
}

// NewServer creates a new Server object.
// The returned close func must be called when the server is no longer needed.
func NewServer() (*Server, func()) {
	s := Server{}
	s.srv = httptest.NewServer(http.HandlerFunc(s.serveHTTP))
	return &s, func() { s.srv.Close() }
}

// NewTLSServer creates a new Server object.
// The returned close func must be called when the server is no longer needed.
func NewTLSServer() (*Server, func()) {
	s := Server{}
	s.srv = httptest.NewTLSServer(http.HandlerFunc(s.serveHTTP))
	return &s, func() { s.srv.Close() }
}

// returns true if the next response is an error response
func (s *Server) isErrorResp() bool {
	if s.static == nil && len(s.resp) == 0 {
		panic("no more responses")
	}
	// always favor static response
	if s.static != nil && s.static.err != nil {
		return true
	}
	if len(s.resp) == 0 {
		return false
	}
	return s.resp[0].err != nil
}

// returns the static response or the next response in the queue
func (s *Server) getResponse() mockResponse {
	if s.static == nil && len(s.resp) == 0 {
		panic("no more responses")
	}
	// always favor static response
	if s.static != nil {
		return *s.static
	}
	// pop off first response and return it
	resp := s.resp[0]
	s.resp = s.resp[1:]
	return resp
}

// URL returns the endpoint of the test server in URL format.
func (s *Server) URL() url.URL {
	u, err := url.Parse(s.srv.URL)
	if err != nil {
		panic(err)
	}
	return *u
}

// Do implements the azcore.Transport interface on Server.
// Calling this when the response queue is empty and no static
// response has been set will cause a panic.
func (s *Server) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	s.count++
	// error responses are returned here
	if s.isErrorResp() {
		resp := s.getResponse()
		return nil, resp.err
	}
	return s.srv.Client().Do(req.WithContext(ctx))
}

func (s *Server) serveHTTP(w http.ResponseWriter, req *http.Request) {
	s.getResponse().write(w)
}

// Requests returns the number of times an HTTP request was made.
func (s *Server) Requests() int {
	return s.count
}

// AppendError appends the error to the end of the response queue.
func (s *Server) AppendError(err error) {
	s.resp = append(s.resp, mockResponse{err: err})
}

// RepeatError appends the error n number of times to the end of the response queue.
func (s *Server) RepeatError(n int, err error) {
	for i := 0; i < n; i++ {
		s.AppendError(err)
	}
}

// SetError indicates the same error should always be returned.
// Any responses set via other methods will be ignored.
func (s *Server) SetError(err error) {
	s.static = &mockResponse{err: err}
}

// AppendResponse appends the response to the end of the response queue.
// If no options are provided the default response is an http.StatusOK.
func (s *Server) AppendResponse(opts ...ResponseOption) {
	mr := mockResponse{code: http.StatusOK}
	for _, o := range opts {
		o.apply(&mr)
	}
	s.resp = append(s.resp, mr)
}

// RepeatResponse appends the response n number of times to the end of the response queue.
// If no options are provided the default response is an http.StatusOK.
func (s *Server) RepeatResponse(n int, opts ...ResponseOption) {
	for i := 0; i < n; i++ {
		s.AppendResponse(opts...)
	}
}

// SetResponse indicates the same response should always be returned.
// Any responses set via other methods will be ignored.
// If no options are provided the default response is an http.StatusOK.
func (s *Server) SetResponse(opts ...ResponseOption) {
	mr := mockResponse{code: http.StatusOK}
	for _, o := range opts {
		o.apply(&mr)
	}
	s.static = &mr
}

// ResponseOption is an abstraction for configuring a mock HTTP response.
type ResponseOption interface {
	apply(mr *mockResponse)
}

type fnRespOpt func(*mockResponse)

func (fn fnRespOpt) apply(mr *mockResponse) {
	fn(mr)
}

type mockResponse struct {
	code    int
	body    []byte
	headers http.Header
	err     error
}

func (mr mockResponse) write(w http.ResponseWriter) {
	w.WriteHeader(mr.code)
	if mr.body != nil {
		w.Write(mr.body)
	}
	if len(mr.headers) > 0 {
		for k, v := range mr.headers {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	}
}

// WithStatusCode sets the HTTP response's status code to the specified value.
func WithStatusCode(c int) ResponseOption {
	return fnRespOpt(func(mr *mockResponse) {
		mr.code = c
	})
}

// WithBody sets the HTTP response's body to the specified value.
func WithBody(b []byte) ResponseOption {
	return fnRespOpt(func(mr *mockResponse) {
		mr.body = b
	})
}

// WithHeader adds the specified header and value to the HTTP response.
func WithHeader(k, v string) ResponseOption {
	return fnRespOpt(func(mr *mockResponse) {
		mr.headers.Add(k, v)
	})
}
