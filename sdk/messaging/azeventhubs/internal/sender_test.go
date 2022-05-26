// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// conforms to amqpSender
type testAmqpSender struct {
	sendErrors []error
	sendCount  int
}

type recoveryCall struct {
	linkID  string
	err     error
	recover bool
}

func (s *testAmqpSender) LinkName() string {
	return "sender-id"
}

func (s *testAmqpSender) Send(ctx context.Context, msg *amqp.Message) error {
	var err error

	if len(s.sendErrors) > s.sendCount {
		err = s.sendErrors[s.sendCount]
	}

	s.sendCount++
	return err
}

func (s *testAmqpSender) Close(ctx context.Context) error {
	return nil
}

func TestSenderRetries(t *testing.T) {
	var recoverCalls []recoveryCall

	var sender *testAmqpSender

	getAmqpSender := func() amqpSender {
		return sender
	}

	recover := func(linkID string, err error, recover bool) {
		recoverCalls = append(recoverCalls, recoveryCall{linkID, err, recover})
	}

	t.Run("SendSucceedsOnFirstTry", func(t *testing.T) {
		recoverCalls = nil
		sender = &testAmqpSender{}

		err := sendMessage(context.TODO(), getAmqpSender, 3, nil, recover)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, sender.sendCount)
		assert.Empty(t, recoverCalls)
	})

	t.Run("SendExceedingRetries", func(*testing.T) {
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				&amqp.DetachError{},
				amqp.ErrSessionClosed,
				errors.New("We'll never attempt to use this one since we ran out of retries")},
		}

		actualErr := sendMessage(context.TODO(), getAmqpSender,
			1, // note we're only allowing 1 retry attempt - so we get the first send() and then 1 additional.
			nil, recover)

		assert.EqualValues(t, amqp.ErrSessionClosed, actualErr)
		assert.EqualValues(t, 2, sender.sendCount)
		assert.EqualValues(t, []recoveryCall{
			{
				linkID:  "sender-id",
				err:     &amqp.DetachError{},
				recover: true,
			},
			{
				linkID:  "sender-id",
				err:     amqp.ErrSessionClosed,
				recover: true,
			},
		}, recoverCalls)

	})

	t.Run("SendWithUnrecoverableAndNonRetryableError", func(*testing.T) {
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				errors.New("Anything not explicitly retryable kills all retries"),
				amqp.ErrConnClosed, // we'll never get here.
			},
		}

		actualErr := sendMessage(context.TODO(), getAmqpSender, 5, nil, recover)

		assert.EqualValues(t, errors.New("Anything not explicitly retryable kills all retries"), actualErr)
		assert.EqualValues(t, 1, sender.sendCount)
		assert.Empty(t, recoverCalls, "No recovery attempts should happen for non-recoverable errors")
	})

	t.Run("SendIsNotRecoverableIfLinkIsClosed", func(*testing.T) {
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				amqp.ErrLinkClosed, // this is no longer considered a retryable error (ErrLinkDetached is, however)
			},
		}

		actualErr := sendMessage(context.TODO(), getAmqpSender, 5, nil, recover)

		assert.EqualValues(t, amqp.ErrLinkClosed, actualErr)
		assert.EqualValues(t, 1, sender.sendCount)
		assert.Empty(t, recoverCalls, "No recovery attempts should happen for non-recoverable errors")
	})

	t.Run("SendWithAmqpErrors", func(*testing.T) {
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{&amqp.Error{
				// retry but doesn't attempt to recover the connection
				Condition: errorServerBusy,
			}, &amqp.Error{
				// retry but doesn't attempt to recover the connection
				Condition: errorTimeout,
			},
				&amqp.Error{
					// retry and will attempt to recover the connection
					Condition: amqp.ErrorNotImplemented,
				}},
		}

		err := sendMessage(context.TODO(), getAmqpSender, 6, nil, recover)
		assert.NoError(t, err)
		assert.EqualValues(t, 4, sender.sendCount)
		assert.EqualValues(t, []recoveryCall{
			{
				linkID: "sender-id",
				err: &amqp.Error{
					Condition: errorServerBusy,
				},
				recover: false,
			},
			{
				linkID: "sender-id",
				err: &amqp.Error{
					Condition: errorTimeout,
				},
				recover: false,
			},
			{
				linkID: "sender-id",
				err: &amqp.Error{
					Condition: amqp.ErrorNotImplemented,
				},
				recover: true,
			},
		}, recoverCalls)
	})

	t.Run("SendWithDetachOrNetErrors", func(*testing.T) {
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				&amqp.DetachError{},
				&net.DNSError{},
			},
		}

		err := sendMessage(context.TODO(), getAmqpSender, 6, nil, recover)
		assert.NoError(t, err)
		assert.EqualValues(t, 3, sender.sendCount)
		assert.EqualValues(t, []recoveryCall{
			{
				linkID:  "sender-id",
				err:     &amqp.DetachError{},
				recover: true,
			},
			{
				linkID:  "sender-id",
				err:     &net.DNSError{},
				recover: true,
			},
		}, recoverCalls)
	})

	t.Run("SendWithRecoverableCloseError", func(*testing.T) {
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				amqp.ErrConnClosed,
				&amqp.DetachError{},
				amqp.ErrSessionClosed,
			},
		}

		err := sendMessage(context.TODO(), getAmqpSender, 6, nil, recover)
		assert.NoError(t, err)
		assert.EqualValues(t, 4, sender.sendCount)
		assert.EqualValues(t, []recoveryCall{
			{
				linkID:  "sender-id",
				err:     amqp.ErrConnClosed,
				recover: true,
			},
			{
				linkID:  "sender-id",
				err:     &amqp.DetachError{},
				recover: true,
			},
			{
				linkID:  "sender-id",
				err:     amqp.ErrSessionClosed,
				recover: true,
			},
		}, recoverCalls)
	})

	t.Run("SendWithInfiniteRetries", func(*testing.T) {
		maxRetries := -1
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				// kind of silly but let's just make sure we would continue to retry.
				amqp.ErrConnClosed,
				amqp.ErrConnClosed,
				amqp.ErrConnClosed,
			},
		}

		err := sendMessage(context.TODO(), getAmqpSender, maxRetries, nil, recover)
		assert.NoError(t, err, "Last call succeeds")
		assert.EqualValues(t, 3+1, sender.sendCount)
		assert.EqualValues(t, recoverCalls, []recoveryCall{
			{linkID: "sender-id", err: amqp.ErrConnClosed, recover: true},
			{linkID: "sender-id", err: amqp.ErrConnClosed, recover: true},
			{linkID: "sender-id", err: amqp.ErrConnClosed, recover: true},
		})
	})

	t.Run("SendWithNoRetries", func(*testing.T) {
		maxRetries := 0
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				amqp.ErrConnClosed, // this is normally a retryable error _but_ we disabled retries.
			},
		}

		err := sendMessage(context.TODO(), getAmqpSender, maxRetries, nil, recover)
		assert.EqualValues(t, amqp.ErrConnClosed, err)
		assert.EqualValues(t, maxRetries+1, sender.sendCount)
		assert.EqualValues(t, recoverCalls, []recoveryCall{
			{linkID: "sender-id", err: amqp.ErrConnClosed, recover: true},
		})
	})

	t.Run("SendRespectsContextStatus", func(*testing.T) {
		maxRetries := 0
		recoverCalls = nil
		sender = &testAmqpSender{
			sendErrors: []error{
				amqp.ErrConnClosed, // this is normally a retryable error _but_ we disabled retries.
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := sendMessage(ctx, getAmqpSender, maxRetries, nil, recover)
		assert.EqualValues(t, context.Canceled, err)
		assert.EqualValues(t, 0, sender.sendCount)
		assert.Empty(t, recoverCalls)
	})
}

type FakeLocker struct {
	afterBlock1 func()
	mu          *sync.Mutex
}

func (l FakeLocker) Lock() {
	l.mu.Lock()
}
func (l FakeLocker) Unlock() {
	l.afterBlock1()
	l.mu.Unlock()
}

// TestRecoveryBlock1 tests recoverWithExpectedLinkID function's first "block" of code that
// decides if we are going to recover the link, ignore it, or wait for an in-progress recovery to
// complete.
func TestRecoveryBlock1(t *testing.T) {
	t.Run("Empty link ID skips link ID checking and just does recovery", func(t *testing.T) {
		cleanup, sender := createRecoveryBlock1Sender(t, func(s *sender) {
			require.True(t, s.recovering)
		})

		defer cleanup()

		err := sender.recoverWithExpectedLinkID(context.TODO(), "")
		require.NoError(t, err)
	})

	t.Run("Matching link ID does recovery", func(t *testing.T) {
		cleanup, sender := createRecoveryBlock1Sender(t, func(s *sender) {
			require.True(t, s.recovering, "s.recovering should be true since the lock is available and we our expected link ID matches")
		})

		defer cleanup()

		err := sender.recoverWithExpectedLinkID(context.TODO(), "the-actual-link-id")
		require.NoError(t, err)
	})

	t.Run("Non-matching link ID skips recovery", func(t *testing.T) {
		cleanup, sender := createRecoveryBlock1Sender(t, func(s *sender) {
			require.False(t, s.recovering, "s.recovering should be false - the link ID isn't current, so nothing needs to be closed/recovered")
		})

		defer cleanup()

		err := sender.recoverWithExpectedLinkID(context.TODO(), "non-matching-link-id")
		require.NoError(t, err)
	})

	// TODO: can't quite test this one
	// t.Run("Already recovering, should wait for condition variable", func(t *testing.T) {
	// 	cleanup, sender := createRecoveryBlock1Sender(t, func(s *sender) {
	// 	})

	// 	defer cleanup()

	// 	sender.recovering = true // oops, someone else is already recovering
	// 	sender.recoverWithExpectedLinkID(context.TODO(), "the-actual-link-id")
	// })
}

func TestAMQPSenderIsCompatibleWithInterface(t *testing.T) {
	var validateCompile amqpSender = &amqp.Sender{}
	require.NotNil(t, validateCompile)
}

type fakeSender struct {
	id     string
	closed bool
}

func (s *fakeSender) ID() string {
	return s.id
}

func (s *fakeSender) LinkName() string {
	return "the-actual-link-id"
}

func (s *fakeSender) Send(ctx context.Context, msg *amqp.Message) error {
	return nil
}
func (s *fakeSender) Close(ctx context.Context) error {
	s.closed = true
	return nil
}

func createRecoveryBlock1Sender(t *testing.T, afterBlock1 func(s *sender)) (func(), *sender) {
	s := &sender{
		partitionID: to.StringPtr("0"),
		hub: &Hub{
			namespace: &namespace{},
		},
	}

	s.sender.Store(&fakeSender{
		id: "the-actual-link-id",
	})

	s.cond = &sync.Cond{
		L: FakeLocker{
			mu: &sync.Mutex{},
			afterBlock1: func() {
				afterBlock1(s)
				panic("Panicking to exit before block 2")
			},
		}}

	return func() {
		val := recover()
		require.EqualValues(t, "Panicking to exit before block 2", val)
	}, s
}
