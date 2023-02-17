// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sbauth"
	"github.com/golang/mock/gomock"
)

type MockData struct {
	nextID int64
	Ctrl   *gomock.Controller

	Events *Events

	cbsRouterOnce    sync.Once
	cbsContext       context.Context
	cancelCbsContext context.CancelFunc

	options *MockDataOptions

	queuesMu sync.Mutex
	queues   map[string]*Queue

	mocksMu sync.Mutex

	conns     []*MockConnection
	receivers map[string][]*MockReceiver
	senders   map[string][]*MockSender
}

type MockDataOptions struct {
	NewConnectionFn func(orig *mock.MockAMQPClient, ctx context.Context) error

	// PreReceiverMock is called with the mock Receiver instance before we instrument it with
	// our defaults.
	PreReceiverMock func(mr *MockReceiver, ctx context.Context) error

	// PreSenderMock is called with the mock Sender instance before we instrument it with
	// our defaults.
	PreSenderMock func(ms *MockSender, ctx context.Context) error

	// PreSessionMock is called with the mock Session instance before we instrument it with
	// our defaults.
	PreSessionMock func(msess *MockSession, ctx context.Context, opts *amqp.SessionOptions) error
}

func NewMockData(t *testing.T, options *MockDataOptions) *MockData {
	if options == nil {
		options = &MockDataOptions{}
	}

	if options.PreReceiverMock == nil {
		options.PreReceiverMock = func(orig *MockReceiver, ctx context.Context) error {
			return nil
		}
	}

	if options.PreSenderMock == nil {
		options.PreSenderMock = func(orig *MockSender, ctx context.Context) error {
			return nil
		}
	}

	if options.PreSessionMock == nil {
		options.PreSessionMock = func(orig *MockSession, ctx context.Context, opts *amqp.SessionOptions) error {
			return nil
		}
	}

	if options.NewConnectionFn == nil {
		options.NewConnectionFn = func(orig *mock.MockAMQPClient, ctx context.Context) error {
			return nil
		}
	}

	cbsContext, cancelCbsContext := context.WithCancel(context.Background())

	return &MockData{
		Ctrl:             gomock.NewController(t),
		queues:           map[string]*Queue{},
		cbsContext:       cbsContext,
		cancelCbsContext: cancelCbsContext,
		Events:           NewEvents(),
		options:          options,
		receivers:        map[string][]*MockReceiver{},
		senders:          map[string][]*MockSender{},
	}
}

func (md *MockData) Close() {
	log.Writef(EventEmulator, "MockData, shutting down")
	defer log.Writef(EventEmulator, "MockData shut down")

	md.cancelCbsContext()

	md.mocksMu.Lock()
	defer md.mocksMu.Unlock()

	for _, conn := range md.conns {
		err := conn.Close()

		if err != nil {
			panic(err)
		}
	}

	md.queuesMu.Lock()
	defer md.queuesMu.Unlock()

	for k := range md.queues {
		md.queues[k].Close()
	}
}

type MockConnection struct {
	Status *Status
	*mock.MockAMQPClient
}

func (md *MockData) NewConnection(ctx context.Context) (amqpwrap.AMQPClient, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	conn := &MockConnection{
		MockAMQPClient: mock.NewMockAMQPClient(md.Ctrl),
		Status:         NewStatus(nil),
	}

	md.mocksMu.Lock()
	md.conns = append(md.conns, conn)
	md.mocksMu.Unlock()

	connID := md.nextUniqueName("c")

	conn.EXPECT().Name().Return(connID).AnyTimes()
	conn.EXPECT().NewSession(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, opts *amqp.SessionOptions) (amqpwrap.AMQPSession, error) {
		return md.newSession(ctx, opts, conn)
	}).AnyTimes()

	conn.EXPECT().Close().DoAndReturn(func() error {
		md.Events.CloseConnection(conn.Name())
		conn.Status.CloseWithError(&amqp.ConnectionError{})
		return nil
	}).AnyTimes()

	md.Events.OpenConnection(conn.Name())
	return conn, nil
}

func (md *MockData) DetachSenders(entityName string) {
	md.mocksMu.Lock()
	defer md.mocksMu.Unlock()

	for _, ms := range md.senders[entityName] {
		ms.Status.CloseWithError(&amqp.DetachError{})
	}

	md.senders[entityName] = nil
}

func (md *MockData) DetachReceivers(entityName string) {
	md.mocksMu.Lock()
	defer md.mocksMu.Unlock()

	for _, mr := range md.receivers[entityName] {
		mr.Status.CloseWithError(&amqp.DetachError{})
	}

	md.receivers[entityName] = nil
}

func (md *MockData) upsertQueue(name string) *Queue {
	md.queuesMu.Lock()
	defer md.queuesMu.Unlock()

	q := md.queues[name]

	if q == nil {
		q = NewQueue(name, md.Events)
		md.queues[name] = q
	}

	return q
}

func (md *MockData) getQueue(name string) *Queue {
	md.queuesMu.Lock()
	defer md.queuesMu.Unlock()

	return md.queues[name]
}

func (md *MockData) AllQueues() map[string]*Queue {
	md.queuesMu.Lock()
	defer md.queuesMu.Unlock()

	m := map[string]*Queue{}

	for k, v := range md.queues {
		m[k] = v
	}

	return m
}

func (md *MockData) nextUniqueName(prefix string) string {
	nextID := atomic.AddInt64(&md.nextID, 1)
	return fmt.Sprintf("%s:%03X", prefix, nextID)
}

func (md *MockData) NewTokenProvider() auth.TokenProvider {
	tc := mock.NewMockTokenCredential(md.Ctrl)

	var tokenCounter int64

	tc.EXPECT().GetToken(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
		tc := atomic.AddInt64(&tokenCounter, 1)

		return azcore.AccessToken{
			Token:     fmt.Sprintf("Token-%d", tc),
			ExpiresOn: time.Now().Add(10 * time.Minute),
		}, nil
	}).AnyTimes()

	return sbauth.NewTokenProvider(tc)
}

func (md *MockData) cbsRouter(ctx context.Context, in *Queue, getQueue func(name string) *Queue) {
	log.Writef(EventEmulator, "cbsRouter starting...")
	defer log.Writef(EventEmulator, "cbsRouter done")

	_ = in.IssueCredit(math.MaxUint32, LinkEvent{
		ConnID: "none",
		Entity: in.name,
		Name:   "none",
		Role:   LinkRoleReceiver,
	}, nil)

	for {
		msg, err := in.Receive(ctx, LinkEvent{
			ConnID: "none",
			Entity: in.name,
			Name:   "none",
			Role:   LinkRoleReceiver,
		}, nil)

		if err != nil {
			log.Writef(EventEmulator, "cbsRouter exiting due to error from receive: %s", err)
			break
		}

		// route response to the right spot
		replyTo := *msg.Properties.ReplyTo

		out := getQueue(replyTo)

		// assume auth is valid.
		err = out.Send(ctx, &amqp.Message{
			Properties: &amqp.MessageProperties{
				CorrelationID: msg.Properties.MessageID,
			},
			ApplicationProperties: map[string]any{
				"statusCode": int32(200),
			},
		}, LinkEvent{
			ConnID: "none",
			Entity: out.name,
			Name:   "none",
			Role:   LinkRoleSender,
		}, nil)

		if err != nil {
			log.Writef(EventEmulator, "cbsRouter exiting due to error from send : %s", err)
			break
		}
	}
}
