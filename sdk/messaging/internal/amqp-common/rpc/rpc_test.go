package rpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestResponseRouterBasic(t *testing.T) {
	receiver := &fakeReceiver{
		Responses: []rpcResponse{
			{amqpMessageWithCorrelationId("my message id"), nil},
			{nil, amqp.ErrLinkClosed},
		},
	}

	link := &Link{
		responseMap: map[string]chan rpcResponse{
			"my message id": make(chan rpcResponse, 1),
		},
		receiver: receiver,
	}

	ch := link.responseMap["my message id"]

	link.startResponseRouter()
	result := <-ch
	require.EqualValues(t, result.message.Data[0], []byte("ID was my message id"))
	require.Empty(t, receiver.Responses)
	require.Nil(t, link.responseMap, "Response map is nil'd out after we get a closed error")
}

func TestResponseRouterMissingMessageID(t *testing.T) {
	receiver := &fakeReceiver{
		Responses: []rpcResponse{
			{amqpMessageWithCorrelationId("my message id"), nil},
			{nil, amqp.ErrLinkClosed},
		},
	}

	link := &Link{
		responseMap: map[string]chan rpcResponse{},
		receiver:    receiver,
	}

	link.startResponseRouter()
	require.Empty(t, receiver.Responses)
}

func TestResponseRouterBadCorrelationID(t *testing.T) {
	messageWithBadCorrelationID := &amqp.Message{
		Properties: &amqp.MessageProperties{
			CorrelationID: uint64(1),
		},
	}

	receiver := &fakeReceiver{
		Responses: []rpcResponse{
			{messageWithBadCorrelationID, nil},
			{nil, amqp.ErrLinkClosed},
		},
	}

	link := &Link{
		responseMap: map[string]chan rpcResponse{},
		receiver:    receiver,
	}

	link.startResponseRouter()
	require.Empty(t, receiver.Responses)
}

func TestResponseRouterFatalErrors(t *testing.T) {
	fatalErrors := []error{
		amqp.ErrLinkClosed,
		amqp.ErrLinkDetached,
		amqp.ErrConnClosed,
		amqp.ErrSessionClosed,
	}

	for _, fatalError := range fatalErrors {
		t.Run(fatalError.Error(), func(t *testing.T) {
			receiver := &fakeReceiver{
				Responses: []rpcResponse{
					{nil, fatalError},
				},
			}
			sentinelCh := make(chan rpcResponse, 1)

			link := &Link{
				responseMap: map[string]chan rpcResponse{
					"sentinel": sentinelCh,
				},
				receiver: receiver,
			}

			link.startResponseRouter()
			require.Empty(t, receiver.Responses)

			// also, we should have broadcasted that the link is closed to anyone else
			// that had not yet received a response but was still waiting.
			select {
			case rpcResponse := <-sentinelCh:
				require.Error(t, rpcResponse.err, fatalError.Error())
				require.Nil(t, rpcResponse.message)
			case <-time.After(time.Second * 5):
				require.Fail(t, "sentinel channel should have received a message")
			}
		})
	}
}

func TestResponseRouterNoResponse(t *testing.T) {
	receiver := &fakeReceiver{
		Responses: []rpcResponse{
			{nil, errors.New("Some other error that will get ignored since we can't route it to anyone (ie: no message ID)")},
			{nil, amqp.ErrConnClosed},
		},
	}

	link := &Link{
		responseMap: map[string]chan rpcResponse{},
		receiver:    receiver,
	}

	link.startResponseRouter()
	require.Empty(t, receiver.Responses)
}

func TestAddMessageID(t *testing.T) {
	message, id, err := addMessageID(&amqp.Message{}, uuid.NewV4)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.EqualValues(t, message.Properties.MessageID, id)

	m := &amqp.Message{
		Data: [][]byte{[]byte("hello world")},
		Properties: &amqp.MessageProperties{
			UserID:    []byte("my user ID"),
			MessageID: "is that will not be copied"},
	}
	message, id, err = addMessageID(m, uuid.NewV4)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.EqualValues(t, message.Properties.MessageID, id)
	require.EqualValues(t, message.Properties.UserID, []byte("my user ID"))
	require.EqualValues(t, message.Data[0], []byte("hello world"))
}

func TestRPCBasic(t *testing.T) {
	fakeUUID := uuid.UUID([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	replyMessage := amqpMessageWithCorrelationId(fakeUUID.String())
	replyMessage.ApplicationProperties = map[string]interface{}{
		"status-code": int32(200),
	}

	ch := make(chan struct{})
	sender := &fakeSender{ch: ch}
	receiver := &fakeReceiver{
		Responses: []rpcResponse{
			{replyMessage, nil},
			{nil, amqp.ErrConnClosed},
		},
		ch: ch,
	}

	l := &Link{
		receiver:                receiver,
		sender:                  sender,
		startResponseRouterOnce: &sync.Once{},
		responseMap:             map[string]chan rpcResponse{},

		uuidNewV4: func() (uuid.UUID, error) {
			return fakeUUID, nil
		},
		messageAccept: func(ctx context.Context, message *amqp.Message) error {
			return nil
		},
	}

	messageToSend := &amqp.Message{}
	resp, err := l.RPC(context.Background(), messageToSend)
	require.NoError(t, err)

	require.EqualValues(t, fakeUUID.String(), sender.Sent[0].Properties.MessageID, "Sent message contains a uniquely generated ID")
	require.EqualValues(t, fakeUUID.String(), resp.Message.Properties.CorrelationID, "Correlation ID matches our originally sent message")

	require.Nil(t, replyMessage.Properties.MessageID, "Original message not modified")
}

func TestRPCFailedSend(t *testing.T) {
	// important bit is that we clean up the channel we stored in the map
	// since we're no longer waiting for the response.
	fakeUUID := uuid.UUID([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	replyMessage := amqpMessageWithCorrelationId(fakeUUID.String())
	replyMessage.ApplicationProperties = map[string]interface{}{
		"status-code": int32(200),
	}

	ch := make(chan struct{})

	sender := &fakeSender{
		ch: ch,
	}
	receiver := &fakeReceiver{
		ch: ch,
		Responses: []rpcResponse{
			{nil, amqp.ErrConnClosed},
		},
	}

	l := &Link{
		receiver:                receiver,
		sender:                  sender,
		startResponseRouterOnce: &sync.Once{},
		responseMap:             map[string]chan rpcResponse{},

		uuidNewV4: func() (uuid.UUID, error) {
			return fakeUUID, nil
		},
		messageAccept: func(ctx context.Context, message *amqp.Message) error {
			panic("Should not be called")
		},
	}

	messageToSend := &amqp.Message{}
	cancelledContext, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := l.RPC(cancelledContext, messageToSend)
	require.Nil(t, resp)
	require.EqualError(t, err, context.Canceled.Error())

	require.EqualValues(t, fakeUUID.String(), sender.Sent[0].Properties.MessageID, "Sent message contains a uniquely generated ID")
}

func TestRPCNilMessageMap(t *testing.T) {
	fakeSender := &fakeSender{}
	fakeReceiver := &fakeReceiver{
		Responses: []rpcResponse{
			// this should let us see what deleteChannelFromMap does
			{amqpMessageWithCorrelationId("hello"), nil},
			{nil, amqp.ErrLinkClosed},
		},
	}

	link := &Link{
		sender:   fakeSender,
		receiver: fakeReceiver,
		// responseMap is nil if the broadcastError() function is called. Since this can be
		// at any time our individual map functions need to handle the map not being
		// there.
		responseMap:             nil,
		startResponseRouterOnce: &sync.Once{},
		uuidNewV4:               uuid.NewV4,
	}

	// sanity check - all the map/channel functions are returning nil
	require.Nil(t, link.addChannelToMap("hello"))
	require.Nil(t, link.deleteChannelFromMap("hello"))

	link.startResponseRouter()

	require.Empty(t, fakeReceiver.Responses, "All responses are used")

	// we're not testing the responseRouter for this second part, so just short-circuit
	// the running.
	link.startResponseRouterOnce.Do(func() {})

	// now check that sending can handle it.
	resp, err := link.RPC(context.Background(), &amqp.Message{})
	require.Error(t, err, amqp.ErrLinkClosed.Error())
	require.Nil(t, resp)
}

func amqpMessageWithCorrelationId(id string) *amqp.Message {
	return &amqp.Message{
		Data: [][]byte{[]byte(fmt.Sprintf("ID was %s", id))},
		Properties: &amqp.MessageProperties{
			CorrelationID: id,
		},
	}
}

type fakeReceiver struct {
	Responses []rpcResponse
	ch        <-chan struct{}
}

func (fr *fakeReceiver) Receive(ctx context.Context) (*amqp.Message, error) {
	// wait until the actual send if we're simulating request/response
	if fr.ch != nil {
		<-fr.ch
	}

	resp := fr.Responses[0]
	fr.Responses = fr.Responses[1:]
	return resp.message, resp.err
}

func (fr *fakeReceiver) Close(ctx context.Context) error {
	panic("Not used for this test")
}

type fakeSender struct {
	Sent []*amqp.Message
	ch   chan<- struct{}
}

func (s *fakeSender) Send(ctx context.Context, msg *amqp.Message) error {
	s.Sent = append(s.Sent, msg)

	if s.ch != nil {
		close(s.ch)
	}

	return nil
}

func (fs *fakeSender) Close(ctx context.Context) error {
	panic("Not used for this test")
}
