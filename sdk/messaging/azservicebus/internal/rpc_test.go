// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

// TestRPCLinkNonErrorRequiresRecovery shows that an error, if it requires recovery,
// will cause the RPCLink to properly broadcast the failure so the caller can initiate
// a link recreation/connection recovery (or potentially just fail out)
func TestRPCLinkNonErrorRequiresRecovery(t *testing.T) {
	tester := &rpcTester{t: t, ResponsesCh: make(chan *rpcTestResp, 1000)}

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

	messagesCh := make(chan string, 10000)
	endCapture := test.CaptureLogsForTestWithChannel(messagesCh)
	defer endCapture()

	responses := []*rpcTestResp{
		// this error requires recovery (in this case, connection but there's no
		// distinction between types in RPCLink)
		{E: &net.DNSError{}},
	}

	resp, err := link.RPC(context.Background(), &amqp.Message{
		ApplicationProperties: map[string]interface{}{
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
		case msg := <-messagesCh:
			if msg == "[rpctesting] "+responseRouterShutdownMessage {
				break LogLoop
			}
		default:
			require.Fail(t, "RPC router never shut down")
		}
	}
}

func TestRPCLinkNonErrorRequiresNoRecovery(t *testing.T) {
	tester := &rpcTester{t: t, ResponsesCh: make(chan *rpcTestResp, 1000), Accepted: make(chan *amqp.Message, 1)}

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

	cleanupLogs := test.CaptureLogsForTest()
	defer cleanupLogs()

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
		ApplicationProperties: map[string]interface{}{
			rpcTesterProperty: responses,
		},
		Properties: &amqp.MessageProperties{
			MessageID: "hello",
		},
	})

	require.NoError(t, err)
	require.Equal(t, 200, resp.Code)
	require.Equal(t, "response from service", resp.Message.Value)

	acceptedMessage := <-tester.Accepted
	require.Equal(t, "response from service", acceptedMessage.Value, "successfully received message is accepted")

	logMessages := cleanupLogs()
	require.Contains(t, logMessages, "[rpctesting] RPCLink had no response channel for correlation ID you've-never-seen-this", "exampleUncorrelatedMessage causes warning for uncorrelated message")
	require.Contains(t, logMessages, "[rpctesting] Non-fatal error in RPCLink, starting to receive again: *Error{Condition: com.microsoft:server-busy, Description: , Info: map[]}")
}

func TestRPCLinkNonErrorLockLostDoesNotBreakAnything(t *testing.T) {
	tester := &rpcTester{t: t, ResponsesCh: make(chan *rpcTestResp, 1000), Accepted: make(chan *amqp.Message, 1)}

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
		ApplicationProperties: map[string]interface{}{
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

	acceptedMessage := <-tester.Accepted
	require.Equal(t, "response from service", acceptedMessage.Value, "successfully received message is accepted")

	// validate that a normal error doesn't cause the response router to shut down
	resp, err = link.RPC(context.Background(), &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			rpcTesterProperty: []*rpcTestResp{
				{M: exampleMessageWithStatusCode(200)},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "response from service", resp.Message.Value)
	acceptedMessage = <-tester.Accepted
	require.Equal(t, "response from service", acceptedMessage.Value, "successfully received message is accepted")
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

	// Accepted contains all the messages where we called AcceptMessage(msg)
	// We only call this when we
	Accepted    chan *amqp.Message
	ResponsesCh chan *rpcTestResp
	t           *testing.T
}

type rpcTestResp struct {
	M *amqp.Message
	E error
}

type rpcTesterClient struct {
	session amqpwrap.AMQPSession
}

func (c *rpcTesterClient) Name() string {
	return "rpcClientName"
}

func (c *rpcTesterClient) NewSession(ctx context.Context, opts *amqp.SessionOptions) (amqpwrap.AMQPSession, error) {
	return c.session, nil
}

func (c *rpcTesterClient) Close() error { return nil }

func (tester *rpcTester) NewReceiver(ctx context.Context, source string, opts *amqp.ReceiverOptions) (amqpwrap.AMQPReceiverCloser, error) {
	return tester, nil
}

func (tester *rpcTester) NewSender(ctx context.Context, target string, opts *amqp.SenderOptions) (amqpwrap.AMQPSenderCloser, error) {
	return tester, nil
}

func (tester *rpcTester) Close(ctx context.Context) error {
	return nil
}

func (tester *rpcTester) LinkName() string {
	return "hello"
}

// receiver functions

func (tester *rpcTester) AcceptMessage(ctx context.Context, msg *amqp.Message) error {
	require.NotNil(tester.t, tester.Accepted, "No messages should be AcceptMessage()'d since the tester.Accepted channel was nil")
	tester.Accepted <- msg
	return nil
}

func (tester *rpcTester) Receive(ctx context.Context) (*amqp.Message, error) {
	resp := <-tester.ResponsesCh
	return resp.M, resp.E
}

// sender functions

func (tester *rpcTester) Send(ctx context.Context, msg *amqp.Message) error {
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

var exampleServerBusyError error = &amqp.Error{Condition: amqp.ErrorCondition("com.microsoft:server-busy")}

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
		ApplicationProperties: map[string]interface{}{
			statusCodeKey: statusCode,
		},
	}
}
