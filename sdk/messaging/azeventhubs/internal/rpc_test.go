// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRPCLink(t *testing.T) {
	initFn := func() *fakeAMQPClient {
		return &fakeAMQPClient{
			session: &FakeAMQPSession{
				NS: &FakeNSForPartClient{
					Receiver: &FakeAMQPReceiver{},
					Sender:   &FakeAMQPSender{},
				},
			},
		}
	}

	t.Run("everything works, RPCLink is created", func(t *testing.T) {
		fakeClient := initFn()

		rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
			Client:   fakeClient,
			Address:  "fake-address",
			LogEvent: log.Event("testing"),
		})

		require.NoError(t, err)
		require.NotNil(t, rpcLink)

		defer test.RequireClose(t, rpcLink)

		require.Zero(t, fakeClient.session.CloseCalled)
		require.Zero(t, fakeClient.session.NS.Receiver.CloseCalled)
		require.Zero(t, fakeClient.session.NS.Sender.CloseCalled)
	})

	t.Run("session created, sender fails", func(t *testing.T) {
		fakeClient := initFn()

		fakeClient.session.NS.NewSenderErr = errors.New("test error")

		rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
			Client:   fakeClient,
			Address:  "fake-address",
			LogEvent: log.Event("testing"),
		})
		require.EqualError(t, err, "test error")
		require.Nil(t, rpcLink)

		require.Equal(t, 1, fakeClient.session.CloseCalled, "session closed as part of cleanup")
		require.Equal(t, 1, fakeClient.session.NS.NewSenderCalled, "sender creation failed, but was called")
		require.Zero(t, fakeClient.session.NS.NewReceiverCalled, "receiver was never created")
	})

	t.Run("receiver fails to be created", func(t *testing.T) {
		// receiver is last in the list, so we'll have to close out sender and session.
		fakeClient := initFn()

		fakeClient.session.NS.NewReceiverErr = errors.New("test error")

		rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
			Client:   fakeClient,
			Address:  "fake-address",
			LogEvent: log.Event("testing"),
		})
		require.EqualError(t, err, "test error")
		require.Nil(t, rpcLink)

		require.Equal(t, 1, fakeClient.session.NS.NewSenderCalled, "sender creation failed, but was called")
		require.Equal(t, 1, fakeClient.session.CloseCalled, "session closed as part of cleanup")
	})
}

// TestRPCLinkNonErrorRequiresRecovery shows that an error, if it requires recovery,
// will cause the RPCLink to properly broadcast the failure so the caller can initiate
// a link recreation/connection recovery (or potentially just fail out)
func TestRPCLinkNonErrorRequiresRecovery(t *testing.T) {
	tester := NewRPCTester(t)
	messages := make(chan string, 10000)
	_ = test.CaptureLogsForTestWithChannel(messages)

	link, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client: &rpcTesterClient{
			session: tester,
		},
		Address:  "some-address",
		LogEvent: "rpctesting",
	})
	require.NoError(t, err)
	require.NotNil(t, link)

	defer func() { require.NoError(t, link.Close(context.Background())) }()

	responses := []*rpcTestResp{
		// this error requires recovery (in this case, connection but there's no
		// distinction between types in RPCLink)
		{E: &net.DNSError{}},
	}

	resp, err := link.RPC(context.Background(), &amqp.Message{
		ApplicationProperties: map[string]any{
			rpcTesterProperty: responses,
		},
	})
	require.Nil(t, resp)

	// (give the response router a teeny bit to shut down)
	time.Sleep(500 * time.Millisecond)

	var netOpError net.Error
	require.ErrorAs(t, err, &netOpError)

LogLoop:
	for {
		select {
		case msg := <-messages:
			if msg == "[rpctesting] "+responseRouterShutdownMessage {
				break LogLoop
			}
		default:
			require.Fail(t, "RPC router never shut down")
		}
	}
}

func TestRPCLinkNonErrorRequiresNoRecovery(t *testing.T) {
	tester := NewRPCTester(t)

	getLogs := test.CaptureLogsForTest()

	link, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client: &rpcTesterClient{
			session: tester,
		},
		Address:  "some-address",
		LogEvent: "rpctesting",
	})
	require.NoError(t, err)
	require.NotNil(t, link)

	defer func() { require.NoError(t, link.Close(context.Background())) }()

	responses := []*rpcTestResp{
		// server busy is a "retry, no reconnect needed" type of error. The response router
		// will just immediately go back to receiving.
		{E: exampleServerBusyError},
		// uncorrelated message, will generate a warning but we'll continue on
		{M: exampleUncorrelatedMessage},
		// this is an actual response and it correlates to the message we sent. We'll get this
		// response back.
		{M: exampleMessageWithStatusCode(200)},
	}

	resp, err := link.RPC(context.Background(), &amqp.Message{
		ApplicationProperties: map[string]any{
			rpcTesterProperty: responses,
		},
		Properties: &amqp.MessageProperties{
			MessageID: "hello",
		},
	})

	require.NoError(t, err)
	require.Equal(t, 200, resp.Code)
	require.Equal(t, "response from service", resp.Message.Value)

	require.NoError(t, link.Close(context.Background()))

	logMessages := getLogs()

	require.Contains(t, logMessages, "[rpctesting] RPCLink had no response channel for correlation ID you've-never-seen-this", "exampleUncorrelatedMessage causes warning for uncorrelated message")
	require.Contains(t, logMessages, "[rpctesting] Non-fatal error in RPCLink, starting to receive again: *Error{Condition: com.microsoft:server-busy, Description: , Info: map[]}")
}

func TestRPCLinkNonErrorLockLostDoesNotBreakAnything(t *testing.T) {
	tester := NewRPCTester(t)

	link, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client: &rpcTesterClient{
			session: tester,
		},
		Address:  "some-address",
		LogEvent: "rpctesting",
	})
	require.NoError(t, err)
	require.NotNil(t, link)

	resp, err := link.RPC(context.Background(), &amqp.Message{
		ApplicationProperties: map[string]any{
			rpcTesterProperty: []*rpcTestResp{
				{M: exampleMessageWithStatusCode(400)},
			},
		},
	})

	// the 400 automatically gets translated into an RPC error. The response router should still be running.
	require.Nil(t, resp)
	var rpcErr RPCError
	require.ErrorAs(t, err, &rpcErr)
	require.Equal(t, 400, rpcErr.RPCCode())

	// validate that a normal error doesn't cause the response router to shut down
	resp, err = link.RPC(context.Background(), &amqp.Message{
		ApplicationProperties: map[string]any{
			rpcTesterProperty: []*rpcTestResp{
				{M: exampleMessageWithStatusCode(200)},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "response from service", resp.Message.Value)
}

func TestRPCLinkClosingClean_SessionCreationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mock.NewMockAMQPClient(ctrl)

	sessionErr := errors.New("failed to create session")

	conn.EXPECT().NewSession(test.NotCancelled, gomock.Any()).Return(nil, sessionErr)

	rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client:   conn,
		Address:  "rpcAddress",
		LogEvent: "Testing",
	})
	require.EqualError(t, err, sessionErr.Error())
	require.Nil(t, rpcLink)
}

func TestRPCLinkClosingClean_SenderCreationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mock.NewMockAMQPClient(ctrl)
	sess := mock.NewMockAMQPSession(ctrl)

	senderErr := errors.New("failed to create sender")

	conn.EXPECT().NewSession(test.NotCancelled, gomock.Any()).Return(sess, nil)
	sess.EXPECT().NewSender(test.NotCancelled, "rpcAddress", gomock.Any(), gomock.Any()).Return(nil, senderErr)
	sess.EXPECT().Close(test.NotCancelled).Return(nil)

	rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client:   conn,
		Address:  "rpcAddress",
		LogEvent: "Testing",
	})
	require.EqualError(t, err, senderErr.Error())
	require.Nil(t, rpcLink)
}

func TestRPCLinkClosingClean_ReceiverCreationFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mock.NewMockAMQPClient(ctrl)
	sess := mock.NewMockAMQPSession(ctrl)
	sender := mock.NewMockAMQPSenderCloser(ctrl)

	receiverErr := errors.New("failed to create receiver")

	conn.EXPECT().NewSession(test.NotCancelled, gomock.Any()).Return(sess, nil)
	sess.EXPECT().NewSender(test.NotCancelled, "rpcAddress", gomock.Any(), gomock.Any()).Return(sender, nil)
	sess.EXPECT().NewReceiver(test.NotCancelled, "rpcAddress", gomock.Any(), gomock.Any()).Return(nil, receiverErr)

	sess.EXPECT().Close(test.NotCancelled).Return(nil)

	rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client:   conn,
		Address:  "rpcAddress",
		LogEvent: "Testing",
	})
	require.EqualError(t, err, receiverErr.Error())
	require.Nil(t, rpcLink)
}

func TestRPCLinkClosingClean_CreationFailsButSessionCloseFailsToo(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mock.NewMockAMQPClient(ctrl)
	sess := mock.NewMockAMQPSession(ctrl)

	senderErr := errors.New("failed to create receiver")

	conn.EXPECT().NewSession(test.NotCancelled, gomock.Any()).Return(sess, nil)
	sess.EXPECT().NewSender(test.NotCancelled, "rpcAddress", gomock.Any(), gomock.Any()).Return(nil, senderErr)
	sess.EXPECT().Close(test.NotCancelled).Return(errors.New("session closing failed"))

	rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client:   conn,
		Address:  "rpcAddress",
		LogEvent: "Testing",
	})
	require.EqualError(t, err, senderErr.Error(), "original error is more relevant, so we favor it over session.Close()")
	require.Nil(t, rpcLink)
}

func TestRPCLinkClosingQuickly(t *testing.T) {
	tester := NewRPCTester(t)

	link, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client: &rpcTesterClient{
			session: tester,
		},
		Address:  "some-address",
		LogEvent: "rpctesting",
	})
	require.NoError(t, err)
	require.NotNil(t, link)
	require.NoError(t, link.Close(context.Background()))
}

func TestRPCLinkBroadcastErrorWhenClosed(t *testing.T) {
	tester := NewRPCTester(t)

	link, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client: &rpcTesterClient{
			session: tester,
		},
		Address:  "some-address",
		LogEvent: "rpctesting",
	})
	require.NoError(t, err)
	require.NotNil(t, link)

	ch := make(chan struct{}, 1)

	go func() {
		defer close(ch)
		_, err := link.RPC(context.Background(), &amqp.Message{
			ApplicationProperties: map[string]any{
				rpcTesterProperty: []*rpcTestResp{},
			},
		})
		require.ErrorIs(t, err, ErrRPCLinkClosed)
	}()

	<-tester.RPCLoopStarted

	require.NoError(t, link.Close(context.Background()))
	<-ch

	// and the error is cached so further calls also get ErrRPCLinkClosed
	// similar to what we do in go-amqp.
	_, err = link.RPC(context.Background(), &amqp.Message{
		ApplicationProperties: map[string]any{
			rpcTesterProperty: []*rpcTestResp{},
		},
	})
	require.ErrorIs(t, err, ErrRPCLinkClosed)
}

func TestRPCLinkCancelClientSideWait(t *testing.T) {
	tester := NewRPCTester(t)

	link, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client: &rpcTesterClient{
			session: tester,
		},
		Address:  "some-address",
		LogEvent: "rpctesting",
	})
	require.NoError(t, err)
	require.NotNil(t, link)

	ch := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer close(ch)
		_, err := link.RPC(ctx, &amqp.Message{
			ApplicationProperties: map[string]any{
				rpcTesterProperty: []*rpcTestResp{},
			},
		})
		require.ErrorIs(t, err, context.Canceled)
	}()

	<-tester.RPCLoopStarted
	cancel()
	<-ch

}

func TestRPCLinkUsesCorrectFlags(t *testing.T) {
	tester := NewRPCTester(t)

	link, err := NewRPCLink(context.Background(), RPCLinkArgs{
		Client: &rpcTesterClient{
			session: tester,
		},
		Address:  "some-address",
		LogEvent: "rpctesting",
	})
	require.NoError(t, err)
	require.NoError(t, link.Close(context.Background()))

	require.Equal(t, amqp.SenderSettleModeSettled, *tester.receiverOpts.RequestedSenderSettleMode)
	require.Equal(t, amqp.ReceiverSettleModeFirst, *tester.receiverOpts.SettlementMode)
}

func NewRPCTester(t *testing.T) *rpcTester {
	return &rpcTester{t: t,
		ResponsesCh:    make(chan *rpcTestResp, 1000),
		RPCLoopStarted: make(chan struct{}, 1),
	}
}

// rpcTester has all the functions needed (for our RPC tests) to be:
// - an AMQPSession
// - an AMQPReceiverCloser
// - an AMQPSenderCloser
// This just makes it simpler since there's this request/response pattern that the tests need. Rather than
// spread it out we can do all the communicating here.
type rpcTester struct {
	amqpwrap.AMQPSenderCloser
	amqpwrap.AMQPReceiverCloser
	receiverOpts *amqp.ReceiverOptions

	ResponsesCh chan *rpcTestResp
	t           *testing.T

	connID uint64

	// RPCLoopStarted is closed when the first Receive() call starts,
	// which indicates that the RPC receiver loop has started.
	RPCLoopStarted      chan struct{}
	closeRPCLoopStarted sync.Once
}

func (c *rpcTester) ConnID() uint64 {
	return c.connID
}

type rpcTestResp struct {
	M *amqp.Message
	E error
}

type rpcTesterClient struct {
	session amqpwrap.AMQPSession
	connID  uint64
}

func (c *rpcTesterClient) ID() uint64 {
	return c.connID
}

func (c *rpcTesterClient) Name() string {
	return "rpcClientName"
}

func (c *rpcTesterClient) NewSession(ctx context.Context, opts *amqp.SessionOptions) (amqpwrap.AMQPSession, error) {
	return c.session, nil
}

func (c *rpcTesterClient) Close() error { return nil }

func (tester *rpcTester) NewReceiver(ctx context.Context, source string, partitionID string, opts *amqp.ReceiverOptions) (amqpwrap.AMQPReceiverCloser, error) {
	tester.receiverOpts = opts
	return tester, nil
}

func (tester *rpcTester) NewSender(ctx context.Context, target string, partitionID string, opts *amqp.SenderOptions) (amqpwrap.AMQPSenderCloser, error) {
	return tester, nil
}

func (tester *rpcTester) Close(ctx context.Context) error {
	return nil
}

func (tester *rpcTester) LinkName() string {
	return "hello"
}

// receiver functions

func (tester *rpcTester) Receive(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
	tester.closeRPCLoopStarted.Do(func() {
		close(tester.RPCLoopStarted)
	})

	select {
	case resp := <-tester.ResponsesCh:
		return resp.M, resp.E
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// sender functions

func (tester *rpcTester) Send(ctx context.Context, msg *amqp.Message, o *amqp.SendOptions) error {
	require.NotEmpty(tester.t, msg.Properties.MessageID)

	// we'll let the payload dictate the response
	if msg.ApplicationProperties["test-send-error"] != nil {
		sendErr := msg.ApplicationProperties["test-send-error"].(error)
		delete(msg.ApplicationProperties, "test-send-error")

		if sendErr != nil {
			return sendErr
		}
	}

	// okay, we're simulating a Send() that works. Let's enqueue the appropriate
	// test response.
	resps := msg.ApplicationProperties[rpcTesterProperty].([]*rpcTestResp)

	for _, resp := range resps {
		if resp.M != nil && resp.M.Properties.CorrelationID == nil {
			// auto-associate it since it's intended to be the response for this message
			resp.M.Properties.CorrelationID = msg.Properties.MessageID
		}

		tester.ResponsesCh <- resp
	}

	return nil
}

// rpcTesterProperty is the property we can shove some messages under that will get
// routed through our rpcTester. It's 100% a test only thing.
const rpcTesterProperty = "test-resps"

var exampleServerBusyError error = &amqp.Error{Condition: amqp.ErrCond("com.microsoft:server-busy")}

var exampleUncorrelatedMessage = &amqp.Message{
	Value: "response from service",
	Properties: &amqp.MessageProperties{
		// this message doesn't actually correlate to a message that was sent
		// it just gets logged and ignored
		CorrelationID: "you've-never-seen-this",
	},
}

func exampleMessageWithStatusCode(statusCode int32) *amqp.Message {
	return &amqp.Message{
		Value: "response from service",
		Properties: &amqp.MessageProperties{
			// will get auto-filled in by the test
			CorrelationID: nil,
		},
		ApplicationProperties: map[string]any{
			statusCodeKey: statusCode,
		},
	}
}
