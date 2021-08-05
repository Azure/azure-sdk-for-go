package servicebus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	MockedDeadLetterBuilder struct {
		mock.Mock
		receiver *MockedReceiveOner
	}

	MockedReceiveOner struct {
		mock.Mock
	}

	MockedTransferDeadLetterBuilder struct {
		mock.Mock
		receiver *MockedReceiveOner
	}
)

func (ro *MockedReceiveOner) ReceiveOne(ctx context.Context, handler Handler) error {
	args := ro.Called(ctx, handler)
	return args.Error(0)
}

func (ro *MockedReceiveOner) Close(ctx context.Context) error {
	args := ro.Called(ctx)
	return args.Error(0)
}

func (m *MockedDeadLetterBuilder) NewDeadLetterReceiver(ctx context.Context, opts ...ReceiverOption) (ReceiveOner, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(ReceiveOner), args.Error(1)
}

func (mt *MockedTransferDeadLetterBuilder) NewTransferDeadLetterReceiver(ctx context.Context, opts ...ReceiverOption) (ReceiveOner, error) {
	args := mt.Called(ctx, opts)
	return args.Get(0).(ReceiveOner), args.Error(1)
}

func TestTransferDeadLetter_ReceiveOne(t *testing.T) {
	builder := setupDefaultDTLMock()
	dl := NewTransferDeadLetter(builder)
	err := dl.ReceiveOne(context.Background(), HandlerFunc(func(ctx context.Context, msg *Message) error {
		return nil
	}))
	assert.NoError(t, err)
	builder.AssertCalled(t, "NewTransferDeadLetterReceiver", context.Background(), mock.Anything)
	builder.receiver.AssertCalled(t, "ReceiveOne", context.Background(), mock.Anything)
}

func TestTransferDeadLetter_Close(t *testing.T) {
	builder := setupDefaultDTLMock()
	dl := NewTransferDeadLetter(builder)
	err := dl.ReceiveOne(context.Background(), HandlerFunc(func(ctx context.Context, msg *Message) error {
		return nil
	}))
	assert.NoError(t, err)
	err = dl.Close(context.Background())
	assert.NoError(t, err)
	builder.receiver.AssertCalled(t, "Close", context.Background())
}

func TestDeadLetter_ReceiveOne(t *testing.T) {
	builder := setupDefaultDLMock()
	dl := NewDeadLetter(builder)
	err := dl.ReceiveOne(context.Background(), HandlerFunc(func(ctx context.Context, msg *Message) error {
		return nil
	}))
	assert.NoError(t, err)
	builder.AssertCalled(t, "NewDeadLetterReceiver", context.Background(), mock.Anything)
	builder.receiver.AssertCalled(t, "ReceiveOne", context.Background(), mock.Anything)
}

func TestDeadLetter_Close(t *testing.T) {
	builder := setupDefaultDLMock()
	dl := NewDeadLetter(builder)
	err := dl.ReceiveOne(context.Background(), HandlerFunc(func(ctx context.Context, msg *Message) error {
		return nil
	}))
	assert.NoError(t, err)
	err = dl.Close(context.Background())
	assert.NoError(t, err)
	builder.receiver.AssertCalled(t, "Close", context.Background())
}

func setupDefaultDLMock() *MockedDeadLetterBuilder {
	receiver := new(MockedReceiveOner)
	receiver.Mock.On(
		"ReceiveOne",
		context.Background(),
		mock.AnythingOfType("servicebus.HandlerFunc"),
	).Return(nil)
	receiver.Mock.On(
		"Close",
		context.Background(),
	).Return(nil)

	builder := &MockedDeadLetterBuilder{
		receiver: receiver,
	}
	builder.Mock.On(
		"NewDeadLetterReceiver",
		context.Background(),
		mock.AnythingOfType("[]servicebus.ReceiverOption"),
	).Return(receiver, nil)
	return builder
}

func setupDefaultDTLMock() *MockedTransferDeadLetterBuilder {
	receiver := new(MockedReceiveOner)
	receiver.Mock.On(
		"ReceiveOne",
		context.Background(),
		mock.AnythingOfType("servicebus.HandlerFunc"),
	).Return(nil)
	receiver.Mock.On(
		"Close",
		context.Background(),
	).Return(nil)

	builder := &MockedTransferDeadLetterBuilder{
		receiver: receiver,
	}
	builder.Mock.On(
		"NewTransferDeadLetterReceiver",
		context.Background(),
		mock.AnythingOfType("[]servicebus.ReceiverOption"),
	).Return(receiver, nil)
	return builder
}
