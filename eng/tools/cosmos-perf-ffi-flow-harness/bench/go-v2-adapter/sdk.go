// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import "context"

type operationResponse struct {
	Body          []byte
	StatusCode    int
	Attempts      uint32
	TransportKind string
	HTTPVersion   string
}

type sdkClient interface {
	Read(context.Context, string, string) (operationResponse, error)
	Create(context.Context, string, string, []byte) (operationResponse, error)
	Upsert(context.Context, string, string, []byte) (operationResponse, error)
	Delete(context.Context, string, string) (operationResponse, error)
	Close() error
}
