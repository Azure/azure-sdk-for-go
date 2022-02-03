// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/devigned/tab"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/go-amqp"
)

const (
	replyPostfix           = "-reply-to-"
	statusCodeKey          = "status-code"
	descriptionKey         = "status-description"
	defaultReceiverCredits = 1000
)

type (
	// rpcLink is the bidirectional communication structure used for CBS negotiation
	rpcLink struct {
		session *amqp.Session

		receiver amqpReceiver // *amqp.Receiver
		sender   amqpSender   // *amqp.Sender

		clientAddress string
		sessionID     *string
		id            string

		responseMu              sync.Mutex
		startResponseRouterOnce *sync.Once
		responseMap             map[string]chan rpcResponse
		broadcastErr            error // the error that caused the responseMap to be nil'd

		// for unit tests
		uuidNewV4     func() (uuid.UUID, error)
		messageAccept func(ctx context.Context, message *amqp.Message) error
	}

	// RPCResponse is the simplified response structure from an RPC like call
	RPCResponse struct {
		Code        int
		Description string
		Message     *amqp.Message
	}

	// RPCLinkOption provides a way to customize the construction of a Link
	RPCLinkOption func(link *rpcLink) error

	rpcResponse struct {
		message *amqp.Message
		err     error
	}

	// Actually: *amqp.Receiver
	amqpReceiver interface {
		Receive(ctx context.Context) (*amqp.Message, error)
		Close(ctx context.Context) error
	}

	amqpSender interface {
		Send(ctx context.Context, msg *amqp.Message) error
		Close(ctx context.Context) error
	}
)

// NewLink will build a new request response link
func NewRPCLink(conn *amqp.Client, address string, opts ...RPCLinkOption) (*rpcLink, error) {
	authSession, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	return newRPCLinkWithSession(authSession, address, opts...)
}

// NewLinkWithSession will build a new request response link, but will reuse an existing AMQP session
func newRPCLinkWithSession(session *amqp.Session, address string, opts ...RPCLinkOption) (*rpcLink, error) {
	linkID, err := uuid.New()
	if err != nil {
		return nil, err
	}

	id := linkID.String()
	link := &rpcLink{
		session:       session,
		clientAddress: strings.Replace("$", "", address, -1) + replyPostfix + id,
		id:            id,

		uuidNewV4:               uuid.New,
		responseMap:             map[string]chan rpcResponse{},
		startResponseRouterOnce: &sync.Once{},
	}

	for _, opt := range opts {
		if err := opt(link); err != nil {
			return nil, err
		}
	}

	sender, err := session.NewSender(
		amqp.LinkTargetAddress(address),
	)
	if err != nil {
		return nil, err
	}

	receiverOpts := []amqp.LinkOption{
		amqp.LinkSourceAddress(address),
		amqp.LinkTargetAddress(link.clientAddress),
		amqp.LinkCredit(defaultReceiverCredits),
	}

	if link.sessionID != nil {
		const name = "com.microsoft:session-filter"
		const code = uint64(0x00000137000000C)
		if link.sessionID == nil {
			receiverOpts = append(receiverOpts, amqp.LinkSourceFilter(name, code, nil))
		} else {
			receiverOpts = append(receiverOpts, amqp.LinkSourceFilter(name, code, link.sessionID))
		}
	}

	receiver, err := session.NewReceiver(receiverOpts...)
	if err != nil {
		// make sure we close the sender
		clsCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = sender.Close(clsCtx)
		return nil, err
	}

	link.sender = sender
	link.receiver = receiver
	link.messageAccept = receiver.AcceptMessage

	return link, nil
}

// startResponseRouter is responsible for taking any messages received on the 'response'
// link and forwarding it to the proper channel. The channel is being select'd by the
// original `RPC` call.
func (l *rpcLink) startResponseRouter() {
	for {
		res, err := l.receiver.Receive(context.Background())

		// You'll see this when the link is shutting down (either
		// service-initiated via 'detach' or a user-initiated shutdown)
		if isClosedError(err) {
			l.broadcastError(err)
			break
		}

		// I don't believe this should happen. The JS version of this same code
		// ignores errors as well since responses should always be correlated
		// to actual send requests. So this is just here for completeness.
		if res == nil {
			continue
		}

		autogenMessageId, ok := res.Properties.CorrelationID.(string)

		if !ok {
			// TODO: it'd be good to track these in some way. We don't have a good way to
			// forward this on at this point.
			continue
		}

		ch := l.deleteChannelFromMap(autogenMessageId)

		if ch != nil {
			ch <- rpcResponse{message: res, err: err}
		}
	}
}

// RPC sends a request and waits on a response for that request
func (l *rpcLink) RPC(ctx context.Context, msg *amqp.Message) (*RPCResponse, error) {
	l.startResponseRouterOnce.Do(func() {
		go l.startResponseRouter()
	})

	copiedMessage, messageID, err := addMessageID(msg, l.uuidNewV4)

	if err != nil {
		return nil, err
	}

	// use the copiedMessage from this point
	msg = copiedMessage

	const altStatusCodeKey, altDescriptionKey = "statusCode", "statusDescription"

	ctx, span := tab.StartSpan(ctx, "rpc.RPC")
	tracing.ApplyComponentInfo(span, Version)
	defer span.End()

	msg.Properties.ReplyTo = &l.clientAddress

	if msg.ApplicationProperties == nil {
		msg.ApplicationProperties = make(map[string]interface{})
	}

	if _, ok := msg.ApplicationProperties["server-timeout"]; !ok {
		if deadline, ok := ctx.Deadline(); ok {
			msg.ApplicationProperties["server-timeout"] = uint(time.Until(deadline) / time.Millisecond)
		}
	}

	responseCh := l.addChannelToMap(messageID)

	if responseCh == nil {
		return nil, l.broadcastErr
	}

	err = l.sender.Send(ctx, msg)

	if err != nil {
		l.deleteChannelFromMap(messageID)
		tab.For(ctx).Error(err)
		return nil, err
	}

	var res *amqp.Message

	select {
	case <-ctx.Done():
		l.deleteChannelFromMap(messageID)
		res, err = nil, ctx.Err()
	case resp := <-responseCh:
		// this will get triggered by the loop in 'startReceiverRouter' when it receives
		// a message with our autoGenMessageID set in the correlation_id property.
		res, err = resp.message, resp.err
	}

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	var statusCode int
	statusCodeCandidates := []string{statusCodeKey, altStatusCodeKey}
	for i := range statusCodeCandidates {
		if rawStatusCode, ok := res.ApplicationProperties[statusCodeCandidates[i]]; ok {
			if cast, ok := rawStatusCode.(int32); ok {
				statusCode = int(cast)
				break
			}

			err := errors.New("status code was not of expected type int32")
			tab.For(ctx).Error(err)
			return nil, err
		}
	}
	if statusCode == 0 {
		err := errors.New("status codes was not found on rpc message")
		tab.For(ctx).Error(err)
		return nil, err
	}

	var description string
	descriptionCandidates := []string{descriptionKey, altDescriptionKey}
	for i := range descriptionCandidates {
		if rawDescription, ok := res.ApplicationProperties[descriptionCandidates[i]]; ok {
			if description, ok = rawDescription.(string); ok || rawDescription == nil {
				break
			} else {
				return nil, errors.New("status description was not of expected type string")
			}
		}
	}

	span.AddAttributes(tab.StringAttribute("http.status_code", fmt.Sprintf("%d", statusCode)))

	response := &RPCResponse{
		Code:        int(statusCode),
		Description: description,
		Message:     res,
	}

	if err := l.messageAccept(ctx, res); err != nil {
		tab.For(ctx).Error(err)
		return response, err
	}

	return response, err
}

// Close the link receiver, sender and session
func (l *rpcLink) Close(ctx context.Context) error {
	ctx, span := startRPCSpan(ctx, "rpc.Close")
	defer span.End()

	if err := l.closeReceiver(ctx); err != nil {
		_ = l.closeSender(ctx)
		_ = l.closeSession(ctx)
		return err
	}

	if err := l.closeSender(ctx); err != nil {
		_ = l.closeSession(ctx)
		return err
	}

	return l.closeSession(ctx)
}

func (l *rpcLink) closeReceiver(ctx context.Context) error {
	ctx, span := startRPCSpan(ctx, "rpc.closeReceiver")
	defer span.End()

	if l.receiver != nil {
		return l.receiver.Close(ctx)
	}
	return nil
}

func (l *rpcLink) closeSender(ctx context.Context) error {
	ctx, span := startRPCSpan(ctx, "rpc.closeSender")
	defer span.End()

	if l.sender != nil {
		return l.sender.Close(ctx)
	}
	return nil
}

func (l *rpcLink) closeSession(ctx context.Context) error {
	ctx, span := startRPCSpan(ctx, "rpc.closeSession")
	defer span.End()

	if l.session != nil {
		return l.session.Close(ctx)
	}
	return nil
}

// addChannelToMap adds a channel which will be used by the response router to
// notify when there is a response to the request.
// If l.responseMap is nil (for instance, via broadcastError) this function will
// return nil.
func (l *rpcLink) addChannelToMap(messageID string) chan rpcResponse {
	l.responseMu.Lock()
	defer l.responseMu.Unlock()

	if l.responseMap == nil {
		return nil
	}

	responseCh := make(chan rpcResponse, 1)
	l.responseMap[messageID] = responseCh

	return responseCh
}

// deleteChannelFromMap removes the message from our internal map and returns
// a channel that the corresponding RPC() call is waiting on.
// If l.responseMap is nil (for instance, via broadcastError) this function will
// return nil.
func (l *rpcLink) deleteChannelFromMap(messageID string) chan rpcResponse {
	l.responseMu.Lock()
	defer l.responseMu.Unlock()

	if l.responseMap == nil {
		return nil
	}

	ch := l.responseMap[messageID]
	delete(l.responseMap, messageID)

	return ch
}

// broadcastError notifies the anyone waiting for a response that the link/session/connection
// has closed.
func (l *rpcLink) broadcastError(err error) {
	l.responseMu.Lock()
	defer l.responseMu.Unlock()

	for _, ch := range l.responseMap {
		ch <- rpcResponse{err: err}
	}

	l.broadcastErr = err
	l.responseMap = nil
}

// addMessageID generates a unique UUID for the message. When the service
// responds it will fill out the correlation ID property of the response
// with this ID, allowing us to link the request and response together.
//
// NOTE: this function copies 'message', adding in a 'Properties' object
// if it does not already exist.
func addMessageID(message *amqp.Message, uuidNewV4 func() (uuid.UUID, error)) (*amqp.Message, string, error) {
	uuid, err := uuidNewV4()

	if err != nil {
		return nil, "", err
	}

	autoGenMessageID := uuid.String()

	// we need to modify the message so we'll make a copy
	copiedMessage := *message

	if message.Properties == nil {
		copiedMessage.Properties = &amqp.MessageProperties{
			MessageID: autoGenMessageID,
		}
	} else {
		// properties already exist, make a copy and then update
		// the message ID
		copiedProperties := *message.Properties
		copiedProperties.MessageID = autoGenMessageID

		copiedMessage.Properties = &copiedProperties
	}

	return &copiedMessage, autoGenMessageID, nil
}

func isClosedError(err error) bool {
	var detachError *amqp.DetachError

	return errors.Is(err, amqp.ErrLinkClosed) ||
		errors.As(err, &detachError) ||
		errors.Is(err, amqp.ErrConnClosed) ||
		errors.Is(err, amqp.ErrSessionClosed)
}

func startRPCSpan(ctx context.Context, operation string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operation)
	tracing.ApplyComponentInfo(span, Version)
	return ctx, span
}
