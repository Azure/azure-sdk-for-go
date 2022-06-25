// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package amqpwrap has some simple wrappers to make it easier to
// abstract the go-amqp types.
package amqpwrap

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
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

// AMQPReceiverCloser is implemented by *amqp.Receiver
type AMQPReceiverCloser interface {
	AMQPReceiver
	Close(ctx context.Context) error
}

// AMQPSender is implemented by *amqp.Sender
type AMQPSender interface {
	Send(ctx context.Context, msg *amqp.Message) error
	MaxMessageSize() uint64
	LinkName() string
}

// AMQPSenderCloser is implemented by *amqp.Sender
type AMQPSenderCloser interface {
	AMQPSender
	Close(ctx context.Context) error
}

// AMQPSession is a simple interface, implemented by *AMQPSessionWrapper.
// It exists only so we can return AMQPReceiver/AMQPSender interfaces.
type AMQPSession interface {
	Close(ctx context.Context) error
	NewReceiver(opts ...amqp.LinkOption) (AMQPReceiverCloser, error)
	NewSender(opts ...amqp.LinkOption) (AMQPSenderCloser, error)
}

type AMQPClient interface {
	Close() error
	NewSession(opts ...amqp.SessionOption) (AMQPSession, error)
}

// AMQPClientWrapper is a simple interface, implemented by *AMQPClientWrapper
// It exists only so we can return AMQPSession, which itself only exists so we can
// return interfaces for AMQPSender and AMQPReceiver from AMQPSession.
type AMQPClientWrapper struct {
	Inner *amqp.Client
}

func (w *AMQPClientWrapper) Close() error {
	return w.Inner.Close()
}

func (w *AMQPClientWrapper) NewSession(opts ...amqp.SessionOption) (AMQPSession, error) {
	sess, err := w.Inner.NewSession(opts...)

	if err != nil {
		return nil, err
	}

	return &AMQPSessionWrapper{
		Inner: sess,
	}, nil
}

type AMQPSessionWrapper struct {
	Inner *amqp.Session
}

func (w *AMQPSessionWrapper) Close(ctx context.Context) error {
	return w.Inner.Close(ctx)
}

func (w *AMQPSessionWrapper) NewReceiver(opts ...amqp.LinkOption) (AMQPReceiverCloser, error) {
	receiver, err := w.Inner.NewReceiver(opts...)

	if err != nil {
		return nil, err
	}

	return receiver, nil
}

func (w *AMQPSessionWrapper) NewSender(opts ...amqp.LinkOption) (AMQPSenderCloser, error) {
	sender, err := w.Inner.NewSender(opts...)

	if err != nil {
		return nil, err
	}

	return sender, nil
}
