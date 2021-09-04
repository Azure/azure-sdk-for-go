package internal

import (
	"context"

	"github.com/Azure/go-amqp"
)

type FakeNamespace struct {
	NextReceiver    *FakeInternalReceiver
	NewReceiverImpl func(ctx context.Context, entityPath string, opts ...ReceiverOption) (LegacyReceiver, error)
}

func (ns *FakeNamespace) NewReceiver(ctx context.Context, entityPath string, opts ...ReceiverOption) (LegacyReceiver, error) {
	return ns.NewReceiverImpl(ctx, entityPath, opts...)
}

func NewFakeLegacyReceiver() *FakeInternalReceiver {
	return &FakeInternalReceiver{
		ListenerRegisteredChan: make(chan ListenerRegisteredEvent, 1),
	}
}

type ListenerRegisteredEvent struct {
	Handler Handler
	Handle  ListenerHandle
	Cancel  context.CancelFunc
}

type FakeInternalReceiver struct {
	// this channel gets an event when `Listen` is called by the Receiver code.
	// It'll contain the latest handler function (passed by the user) as well as
	// as the listener handle and a function you can use to cancel the listen
	// operation.
	ListenerRegisteredChan chan ListenerRegisteredEvent
	CloseCalled            bool
	AbandonCalled          bool
	CompleteCalled         bool

	DrainCalled  bool
	CreditsAdded uint32

	CloseImpl           func(ctx context.Context) error
	ListenImpl          func(ctx context.Context, handler Handler) ListenerHandle
	AbandonMessageImpl  func(ctx context.Context, msg *Message) error
	CompleteMessageImpl func(ctx context.Context, msg *Message) error
}

func (r *FakeInternalReceiver) Close(ctx context.Context) error {
	r.CloseCalled = true

	if r.CloseImpl == nil {
		return nil
	}

	return r.CloseImpl(ctx)
}

func (r *FakeInternalReceiver) Listen(ctx context.Context, handler Handler) ListenerHandle {
	if r.ListenImpl != nil {
		return r.ListenImpl(ctx, handler)
	}

	// default listener just creates a cancellable context
	// and notifies you, via the listenerRegisteredChan, that
	// Listen() has been called.
	listenerCtx, cancel := context.WithCancel(context.Background())

	r.ListenerRegisteredChan <- ListenerRegisteredEvent{
		Handler: handler,
		Handle:  listenerCtx,
		Cancel:  cancel,
	}

	return listenerCtx
}

func (r *FakeInternalReceiver) AbandonMessage(ctx context.Context, msg *Message) error {
	r.AbandonCalled = true

	if r.AbandonMessageImpl == nil {
		return nil
	}

	return r.AbandonMessageImpl(ctx, msg)
}

func (r *FakeInternalReceiver) CompleteMessage(ctx context.Context, msg *Message) error {
	r.CompleteCalled = true

	if r.CompleteMessageImpl == nil {
		return nil
	}

	return r.CompleteMessageImpl(ctx, msg)
}

func (r *FakeInternalReceiver) IssueCredit(credit uint32) error {
	r.CreditsAdded += credit
	return nil
}

func (r *FakeInternalReceiver) DrainCredit(ctx context.Context) error {
	r.DrainCalled = true
	return nil
}

func (r *FakeInternalReceiver) ReceiveOne(ctx context.Context, handler Handler) error {
	return nil
}

func (r *FakeInternalReceiver) Session() *amqp.Session {
	return nil
}

func (r *FakeInternalReceiver) SessionID() *string {
	return nil
}

func (r *FakeInternalReceiver) Client() *amqp.Client {
	return nil
}

type FakeListenerHandle struct {
	DoneChan chan struct{}
	ErrValue error
}

func (lh *FakeListenerHandle) Err() error {
	return lh.ErrValue
}

func (lh *FakeListenerHandle) Done() <-chan struct{} {
	return lh.DoneChan
}

func NewFakeNamespace() *FakeNamespace {
	receiver := &FakeInternalReceiver{}
	fakeNS := &FakeNamespace{}

	fakeNS.NextReceiver = receiver

	fakeNS.NewReceiverImpl = func(ctx context.Context, entityPath string, opts ...ReceiverOption) (LegacyReceiver, error) {
		return fakeNS.NextReceiver, nil
	}

	return fakeNS
}
