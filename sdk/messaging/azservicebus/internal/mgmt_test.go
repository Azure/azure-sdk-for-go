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
			name: "deadline shorter than 60s returns floor of 60s",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			wantExact: ptrUint(60000),
		},
		{
			name: "very short deadline returns floor of 60s",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 1*time.Second)
			},
			wantExact: ptrUint(60000),
		},
		{
			name: "deadline longer than 60s respects user timeout",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 120*time.Second)
			},
			// time.Until shifts slightly between setup and call, so allow a range
			wantMin: 119000,
			wantMax: 120500,
		},
		{
			name: "deadline exactly 60s returns ~60000 ms",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 60*time.Second)
			},
			// With exactly 60s remaining, time.Until will be slightly under 60s,
			// so the floor (defaultServerTimeout) wins → exact 60000.
			wantExact: ptrUint(60000),
		},
		{
			name: "expired context returns floor of 60s",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
				return ctx, cancel
			},
			// time.Until returns negative, which is < defaultServerTimeout, so the floor wins.
			wantExact: ptrUint(60000),
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
