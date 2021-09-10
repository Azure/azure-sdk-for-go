package internal

import (
	"context"

	"github.com/Azure/go-amqp"
)

// FakeNamespace is a fake namespace suitable for unit testing.
type FakeNamespace struct {
	NextReceiver    *FakeInternalReceiver
	NewReceiverImpl func(ctx context.Context, entityPath string, opts ...ReceiverOption) (LegacyReceiver, error)
}

// NewReceiver creates a receiver using your fakens.NewReceiverImpl function.
func (ns *FakeNamespace) NewReceiver(ctx context.Context, entityPath string, opts ...ReceiverOption) (LegacyReceiver, error) {
	return ns.NewReceiverImpl(ctx, entityPath, opts...)
}

// NewFakeLegacyReceiver creates a `legacyReceiver` with the minimum required
// fields.
func NewFakeLegacyReceiver() *FakeInternalReceiver {
	return &FakeInternalReceiver{
		ListenerRegisteredChan: make(chan ListenerRegisteredEvent, 1),
	}
}

// ListenerRegisteredEvent is sent on the `FakeInternalReceiver.ListenerRegisteredChan`
// when the caller calls into `FakeInternalReceiver.Listen`. Useful for when you want to
// call into the registered handler (for instance, to pump a single message through).
type ListenerRegisteredEvent struct {
	Handler Handler
	Handle  ListenerHandle
	Cancel  context.CancelFunc
}

// FakeInternalReceiver is a basic mock that lets you stub in custom implementations for
// `Close`, `Listen`, `AbandonMessage`,  `CompleteMessage`, etc...
// It also provides complementary bools `<function>Called` that let you know if a
// function was called.
// Feel free to add more as you find them!
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

// Close mimics a `legacyReceiver.Close` call.
func (r *FakeInternalReceiver) Close(ctx context.Context) error {
	r.CloseCalled = true

	if r.CloseImpl == nil {
		return nil
	}

	return r.CloseImpl(ctx)
}

// Listen mimics a `legacyReceiver.Listen` call.
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

// AbandonMessage mimics a `legacyReceiver.AbandonMessage` call. By default returns nil.
func (r *FakeInternalReceiver) AbandonMessage(ctx context.Context, msg *Message) error {
	r.AbandonCalled = true

	if r.AbandonMessageImpl == nil {
		return nil
	}

	return r.AbandonMessageImpl(ctx, msg)
}

// CompleteMessage mimics a `legacyReceiver.CompleteMessage` call. By default returns nil.
func (r *FakeInternalReceiver) CompleteMessage(ctx context.Context, msg *Message) error {
	r.CompleteCalled = true

	if r.CompleteMessageImpl == nil {
		return nil
	}

	return r.CompleteMessageImpl(ctx, msg)
}

// IssueCredit mimics a `legacyReceiver.IssueCredit` call. By default returns nil.
func (r *FakeInternalReceiver) IssueCredit(credit uint32) error {
	r.CreditsAdded += credit
	return nil
}

// DrainCredit mimics a `legacyReceiver.DrainCredit` call. By default returns nil.
func (r *FakeInternalReceiver) DrainCredit(ctx context.Context) error {
	r.DrainCalled = true
	return nil
}

// ReceiveOne mimics a `legacyReceiver.ReceiveOne` call. By default returns nil.
func (r *FakeInternalReceiver) ReceiveOne(ctx context.Context, handler Handler) error {
	return nil
}

// Session mimics a `legacyReceiver.Session` call. By default returns nil.
func (r *FakeInternalReceiver) Session() *amqp.Session {
	return nil
}

// SessionID mimics a `legacyReceiver.SessionID` call. By default returns nil.
func (r *FakeInternalReceiver) SessionID() *string {
	return nil
}

// Client mimics a `legacyReceiver.Client` call. By default returns nil.
func (r *FakeInternalReceiver) Client() *amqp.Client {
	return nil
}

// FakeListenerHandle mimics a handle returned from a Listen call. You can control
// the important bits using `DoneChan` (to mark the listen handle as Done) or `ErrValue`,
// which lets you control the value returned from `Err`
type FakeListenerHandle struct {
	DoneChan chan struct{}
	ErrValue error
}

// Err returns the value in lh.ErrValue. Default will be nil.
func (lh *FakeListenerHandle) Err() error {
	return lh.ErrValue
}

// Done returns the value in lh.DoneChan. Default will be nil!
func (lh *FakeListenerHandle) Done() <-chan struct{} {
	return lh.DoneChan
}

// NewFakeNamespace creates a legacyNamespace compatible struct.
// By default, fakeNs.NewReceiverImpl returns the fakeNS.NextReceiver value.
// - If a single cached fake receiver is enough for your test replace fakeNS.NextReceiver.
// - If you want more control you can replace fakeNS.NewReceiverImpl.
func NewFakeNamespace() *FakeNamespace {
	receiver := &FakeInternalReceiver{}
	fakeNS := &FakeNamespace{}

	fakeNS.NextReceiver = receiver

	fakeNS.NewReceiverImpl = func(ctx context.Context, entityPath string, opts ...ReceiverOption) (LegacyReceiver, error) {
		return fakeNS.NextReceiver, nil
	}

	return fakeNS
}
