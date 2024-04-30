// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package amqpwrap

import (
	"context"

	"github.com/Azure/go-amqp"
)

// RPCResponse is the simplified response structure from an RPC like call
type RPCResponse struct {
	// Code is the response code - these originate from Service Bus. Some
	// common values are called out below, with the RPCResponseCode* constants.
	//
	// NOTE: These status codes are intended to mirror HTTP status codes. For instance
	// peeking messages returns http.StatusOK, etc...
	//
	// See https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-amqp-request-response
	// for more details about the ins and outs of each operation.
	Code        int
	Description string
	Message     *amqp.Message
}

// RPCLink is implemented by *rpc.Link
type RPCLink interface {
	Close(ctx context.Context) error
	RPC(ctx context.Context, msg *amqp.Message) (*RPCResponse, error)
}
