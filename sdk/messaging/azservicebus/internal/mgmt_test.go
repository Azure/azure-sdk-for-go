// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerTimeoutMillis(t *testing.T) {
	tests := []struct {
		name      string
		ctx       func() (context.Context, context.CancelFunc)
		wantMin   uint
		wantMax   uint
		wantExact *uint // when non-nil, assert exact value instead of range
	}{
		{
			name: "no deadline returns default 60s",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			wantExact: ptrUint(60000),
		},
		{
			name: "deadline shorter than 60s uses remaining time",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			// The remaining deadline is used so server-timeout never exceeds the
			// client-side wait. time.Until can only be slightly under 10s, so the
			// upper bound is the deadline itself and the lower bound tolerates
			// scheduler drift without flaking.
			wantMin: 9500,
			wantMax: 10000,
		},
		{
			name: "very short deadline uses remaining time",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 1*time.Second)
			},
			wantMin: 800,
			wantMax: 1000,
		},
		{
			name: "deadline longer than 60s respects user timeout",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 120*time.Second)
			},
			// time.Until shifts slightly between setup and call, so allow a range
			wantMin: 119000,
			wantMax: 120000,
		},
		{
			name: "deadline exactly 60s uses remaining time",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 60*time.Second)
			},
			// With exactly 60s remaining, time.Until is slightly under 60s.
			wantMin: 59000,
			wantMax: 60000,
		},
		{
			name: "expired context returns 0",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
				return ctx, cancel
			},
			// Negative remaining is clamped to 0 to avoid an unsigned underflow; the
			// RPC returns at ctx.Done() immediately, so the value is not used.
			wantExact: ptrUint(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctx()
			defer cancel()

			got := serverTimeoutMillis(ctx)

			if tt.wantExact != nil {
				require.Equal(t, *tt.wantExact, got)
			} else {
				require.GreaterOrEqual(t, got, tt.wantMin, "timeout should be >= %d, got %d", tt.wantMin, got)
				require.LessOrEqual(t, got, tt.wantMax, "timeout should be <= %d, got %d", tt.wantMax, got)
			}
		})
	}
}

func ptrUint(v uint) *uint {
	return &v
}
