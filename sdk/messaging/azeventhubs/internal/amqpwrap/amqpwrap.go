// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package amqpwrap has some simple wrappers to make it easier to
// abstract the go-amqp types.
package amqpwrap

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
)

// AMQPReceiver is implemented by *amqp.Receiver
type AMQPReceiver interface {
	IssueCredit(credit uint32) error
	Receive(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error)
	Prefetched() *amqp.Message

	// settlement functions
	AcceptMessage(ctx context.Context, msg *amqp.Message) error
	RejectMessage(ctx context.Context, msg *amqp.Message, e *amqp.Error) error
	ReleaseMessage(ctx context.Context, msg *amqp.Message) error
	ModifyMessage(ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions) error

	LinkName() string
	LinkSourceFilterValue(name string) any

	// wrapper only functions,

	// Credits returns the # of credits still active on this link.
	Credits() uint32
}

// AMQPReceiverCloser is implemented by *amqp.Receiver
type AMQPReceiverCloser interface {
	AMQPReceiver
	Close(ctx context.Context) error
}

// AMQPSender is implemented by *amqp.Sender
type AMQPSender interface {
	Send(ctx context.Context, msg *amqp.Message, o *amqp.SendOptions) error
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
	NewReceiver(ctx context.Context, source string, opts *amqp.ReceiverOptions) (AMQPReceiverCloser, error)
	NewSender(ctx context.Context, target string, opts *amqp.SenderOptions) (AMQPSenderCloser, error)
}

type AMQPClient interface {
	Close() error
	NewSession(ctx context.Context, opts *amqp.SessionOptions) (AMQPSession, error)
}

// RPCLink is implemented by *rpc.Link
type RPCLink interface {
	Close(ctx context.Context) error
	RPC(ctx context.Context, msg *amqp.Message) (*RPCResponse, error)
	LinkName() string
}

// RPCResponse is the simplified response structure from an RPC like call
type RPCResponse struct {
	// Code is the response code - these originate from Service Bus. Some
	// common values are called out below, with the RPCResponseCode* constants.
	Code        int
	Description string
	Message     *amqp.Message
}

type goamqpConn interface {
	NewSession(ctx context.Context, opts *amqp.SessionOptions) (*amqp.Session, error)
	Close() error
}

type goamqpSession interface {
	Close(ctx context.Context) error
	NewReceiver(ctx context.Context, source string, opts *amqp.ReceiverOptions) (*amqp.Receiver, error)
	NewSender(ctx context.Context, target string, opts *amqp.SenderOptions) (*amqp.Sender, error)
}

type goamqpReceiver interface {
	IssueCredit(credit uint32) error
	Receive(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error)
	Prefetched() *amqp.Message

	// settlement functions
	AcceptMessage(ctx context.Context, msg *amqp.Message) error
	RejectMessage(ctx context.Context, msg *amqp.Message, e *amqp.Error) error
	ReleaseMessage(ctx context.Context, msg *amqp.Message) error
	ModifyMessage(ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions) error

	LinkName() string
	LinkSourceFilterValue(name string) any
	Close(ctx context.Context) error
}

// AMQPClientWrapper is a simple interface, implemented by *AMQPClientWrapper
// It exists only so we can return AMQPSession, which itself only exists so we can
// return interfaces for AMQPSender and AMQPReceiver from AMQPSession.
type AMQPClientWrapper struct {
	Inner goamqpConn
}

func (w *AMQPClientWrapper) Close() error {
	return w.Inner.Close()
}

func (w *AMQPClientWrapper) NewSession(ctx context.Context, opts *amqp.SessionOptions) (AMQPSession, error) {
	sess, err := w.Inner.NewSession(ctx, opts)

	if err != nil {
		return nil, HandleNewOrCloseError(err)
	}

	return &AMQPSessionWrapper{
		Inner: sess,
	}, nil
}

type AMQPSessionWrapper struct {
	Inner goamqpSession
}

func (w *AMQPSessionWrapper) Close(ctx context.Context) error {
	if err := w.Inner.Close(ctx); err != nil {
		return HandleNewOrCloseError(err)
	}

	return nil
}

func (w *AMQPSessionWrapper) NewReceiver(ctx context.Context, source string, opts *amqp.ReceiverOptions) (AMQPReceiverCloser, error) {
	receiver, err := w.Inner.NewReceiver(ctx, source, opts)

	if err != nil {
		return nil, HandleNewOrCloseError(err)
	}

	return &AMQPReceiverWrapper{inner: receiver}, nil
}

func (w *AMQPSessionWrapper) NewSender(ctx context.Context, target string, opts *amqp.SenderOptions) (AMQPSenderCloser, error) {
	sender, err := w.Inner.NewSender(ctx, target, opts)

	if err != nil {
		return nil, HandleNewOrCloseError(err)
	}

	return sender, nil
}

type AMQPReceiverWrapper struct {
	inner   goamqpReceiver
	credits uint32
}

func (rw *AMQPReceiverWrapper) Credits() uint32 {
	return rw.credits
}

func (rw *AMQPReceiverWrapper) IssueCredit(credit uint32) error {
	err := rw.inner.IssueCredit(credit)

	if err == nil {
		rw.credits += credit
	}

	return err
}

func (rw *AMQPReceiverWrapper) Receive(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
	message, err := rw.inner.Receive(ctx, o)

	if err != nil {
		return nil, err
	}

	rw.credits--
	return message, nil
}

func (rw *AMQPReceiverWrapper) Prefetched() *amqp.Message {
	msg := rw.inner.Prefetched()

	if msg == nil {
		return nil
	}

	rw.credits--
	return msg
}

// settlement functions
func (rw *AMQPReceiverWrapper) AcceptMessage(ctx context.Context, msg *amqp.Message) error {
	return rw.inner.AcceptMessage(ctx, msg)
}

func (rw *AMQPReceiverWrapper) RejectMessage(ctx context.Context, msg *amqp.Message, e *amqp.Error) error {
	return rw.inner.RejectMessage(ctx, msg, e)
}

func (rw *AMQPReceiverWrapper) ReleaseMessage(ctx context.Context, msg *amqp.Message) error {
	return rw.inner.ReleaseMessage(ctx, msg)
}

func (rw *AMQPReceiverWrapper) ModifyMessage(ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions) error {
	return rw.inner.ModifyMessage(ctx, msg, options)
}

func (rw *AMQPReceiverWrapper) LinkName() string {
	return rw.inner.LinkName()
}

func (rw *AMQPReceiverWrapper) LinkSourceFilterValue(name string) any {
	return rw.inner.LinkSourceFilterValue(name)
}

func (rw *AMQPReceiverWrapper) Close(ctx context.Context) error {
	if err := rw.inner.Close(ctx); err != nil {
		return HandleNewOrCloseError(err)
	}

	return nil
}

type AMQPSenderWrapper struct {
	inner AMQPSenderCloser
}

func (sw *AMQPSenderWrapper) Send(ctx context.Context, msg *amqp.Message, o *amqp.SendOptions) error {
	return sw.inner.Send(ctx, msg, o)
}

func (sw *AMQPSenderWrapper) MaxMessageSize() uint64 {
	return sw.inner.MaxMessageSize()
}

func (sw *AMQPSenderWrapper) LinkName() string {
	return sw.inner.LinkName()
}

func (sw *AMQPSenderWrapper) Close(ctx context.Context) error {
	if err := sw.inner.Close(ctx); err != nil {
		return HandleNewOrCloseError(err)
	}

	return nil
}

func HandleNewOrCloseError(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		// we treat this a bit differently when creating or closing entities. The big concern here is
		// that failing to cleanup stray
		return ErrConnResetNeeded
	}

	return err
}

var ErrConnResetNeeded = errors.New("connection must be reset, link/connection state may be inconsistent")
