// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestManagementOperationTimeouts(t *testing.T) {
	t.Run("RenewLocks respects default timeout when no context deadline", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRPCLink := mock.NewMockRPCLink(ctrl)
		
		// Create a context that hangs for a long time to verify our timeout logic works
		mockRPCLink.EXPECT().RPC(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
			// Verify that the context has a deadline
			_, hasDeadline := ctx.Deadline()
			require.True(t, hasDeadline, "Context should have a deadline")
			
			// Sleep for a short time to simulate some work, but not long enough to hit timeout
			time.Sleep(10 * time.Millisecond)
			
			return &amqpwrap.RPCResponse{
				Code:    200,
				Message: &amqp.Message{Value: map[string]any{"expirations": []time.Time{time.Now().Add(time.Hour)}}},
			}, nil
		})

		// Call with context that has no deadline
		ctx := context.Background()
		lockTokens := []amqp.UUID{{0x01}}
		
		result, err := RenewLocks(ctx, mockRPCLink, "test-link", lockTokens)
		require.NoError(t, err)
		require.Len(t, result, 1)
	})

	t.Run("RenewLocks preserves existing context timeout", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRPCLink := mock.NewMockRPCLink(ctrl)
		
		// Set up a shorter timeout than the default management timeout
		shortTimeout := 500 * time.Millisecond
		
		mockRPCLink.EXPECT().RPC(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
			// Verify that the context has the original deadline
			deadline, hasDeadline := ctx.Deadline()
			require.True(t, hasDeadline, "Context should have a deadline")
			
			// The deadline should be close to our original timeout (within a reasonable margin)
			timeUntilDeadline := time.Until(deadline)
			require.True(t, timeUntilDeadline <= shortTimeout, "Context should preserve original timeout")
			require.True(t, timeUntilDeadline > shortTimeout-100*time.Millisecond, "Context deadline should be close to original")
			
			return &amqpwrap.RPCResponse{
				Code:    200,
				Message: &amqp.Message{Value: map[string]any{"expirations": []time.Time{time.Now().Add(time.Hour)}}},
			}, nil
		})

		// Call with context that already has a deadline shorter than the default
		ctx, cancel := context.WithTimeout(context.Background(), shortTimeout)
		defer cancel()
		
		lockTokens := []amqp.UUID{{0x01}}
		
		result, err := RenewLocks(ctx, mockRPCLink, "test-link", lockTokens)
		require.NoError(t, err)
		require.Len(t, result, 1)
	})
}