// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"

	"github.com/Azure/go-amqp"
)

// AMQPReceiver is implemented by *amqp.Receiver
type AMQPReceiver interface {
	IssueCredit(credit uint32) error
	DrainCredit(ctx context.Context) error
	Receive(ctx context.Context) (*amqp.Message, error)
	Prefetched(ctx context.Context) (*amqp.Message, error)

	// settlement functions
	AcceptMessage(ctx context.Context, msg *amqp.Message) error
	RejectMessage(ctx context.Context, msg *amqp.Message, e *amqp.Error) error
	ReleaseMessage(ctx context.Context, msg *amqp.Message) error
	ModifyMessage(ctx context.Context, msg *amqp.Message, deliveryFailed, undeliverableHere bool, messageAnnotations amqp.Annotations) error

	LinkName() string
	LinkSourceFilterValue(name string) interface{}
}

// AMQPReceiver is implemented by *amqp.Receiver
type AMQPReceiverCloser interface {
	AMQPReceiver
	Close(ctx context.Context) error
}

// AMQPSession is implemented by *amqp.Session
type AMQPSession interface {
	NewReceiver(opts ...amqp.LinkOption) (*amqp.Receiver, error)
	NewSender(opts ...amqp.LinkOption) (*amqp.Sender, error)
}

// AMQPSessionCloser is implemented by *amqp.Session
type AMQPSessionCloser interface {
	AMQPSession
	Close(ctx context.Context) error
}

// AMQPSender is implemented by *amqp.Sender
type AMQPSender interface {
	Send(ctx context.Context, msg *amqp.Message) error
	MaxMessageSize() uint64
}

// AMQPSenderCloser is implemented by *amqp.Sender
type AMQPSenderCloser interface {
	AMQPSender
	Close(ctx context.Context) error
}

// RPCLink is implemented by *rpc.Link
type RPCLink interface {
	Close(ctx context.Context) error
	RPC(ctx context.Context, msg *amqp.Message) (*RPCResponse, error)
}

// Closeable is implemented by pretty much any AMQP link/client
// including our own higher level Receiver/Sender.
type Closeable interface {
	Close(ctx context.Context) error
}
