package servicebus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	MockedBuilder struct {
		mock.Mock
	}
)

func (m *MockedBuilder) ManagementPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockedBuilder) NewReceiver(ctx context.Context, opts ...ReceiverOption) (*Receiver, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*Receiver), args.Error(1)
}

func (m *MockedBuilder) NewSender(ctx context.Context, opts ...SenderOption) (*Sender, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*Sender), args.Error(1)
}

func TestQueueSession_SessionID(t *testing.T) {
	builder := new(MockedBuilder)
	sessionID := "123"
	qs := NewQueueSession(builder, &sessionID)
	assert.Equal(t, sessionID, *qs.sessionID)
}
