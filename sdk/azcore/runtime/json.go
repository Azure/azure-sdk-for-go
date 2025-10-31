//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/exported"
)

// Unmarshaler is an interface for custom JSON unmarshaling implementations.
// By default, the Azure SDK uses encoding/json from the standard library.
// You can replace this with high-performance alternatives like sonic or jsoniter.
//
// Example using sonic:
//
//	import "github.com/bytedance/sonic"
//
//	type SonicUnmarshaler struct{}
//
//	func (s *SonicUnmarshaler) Unmarshal(data []byte, v any) error {
//	    return sonic.Unmarshal(data, v)
//	}
//
//	func init() {
//	    runtime.SetUnmarshaler(&SonicUnmarshaler{})
//	}
//
// This will affect ALL JSON unmarshaling in the SDK including:
//   - Single response operations (Get, Create, Update, Delete)
//   - Pagination (Pager.NextPage)
//   - Long-running operations (Poller.PollUntilDone, Poller.Result)
//   - Resume tokens (NewPollerFromResumeToken)
type Unmarshaler interface {
	// Unmarshal parses the JSON-encoded data and stores the result
	// in the value pointed to by v.
	//
	// Unmarshal uses the same semantics as encoding/json.Unmarshal:
	// - If v is nil or not a pointer, Unmarshal returns an error
	// - Unmarshal uses the inverse of the encoding rules that Marshal uses
	Unmarshal(data []byte, v any) error
}

// SetUnmarshaler replaces the default JSON unmarshaler with a custom implementation.
// This affects all subsequent JSON unmarshaling operations in the Azure SDK.
//
// IMPORTANT: This is a global setting and should typically be called once during
// application initialization (e.g., in init() or early in main()).
// It is NOT goroutine-safe and should not be called concurrently with SDK operations.
//
// For thread-safe per-request customization, use WithUnmarshaler() instead.
//
// Example:
//
//	import "github.com/bytedance/sonic"
//
//	func init() {
//	    runtime.SetUnmarshaler(&SonicUnmarshaler{})
//	}
func SetUnmarshaler(u Unmarshaler) {
	exported.SetUnmarshaler(u)
}

// WithUnmarshaler returns a new context with the specified Unmarshaler attached.
// This is THREAD-SAFE and allows per-request unmarshaler customization without
// affecting global state or other concurrent operations.
//
// The unmarshaler is used for all JSON unmarshaling operations within the context,
// including single operations, pagination, and long-running operations.
//
// Example for concurrent operations with different unmarshalers:
//
//	import "github.com/bytedance/sonic"
//
//	type SonicUnmarshaler struct{}
//	func (s *SonicUnmarshaler) Unmarshal(data []byte, v any) error {
//	    return sonic.Unmarshal(data, v)
//	}
//
//	// This is thread-safe - each goroutine uses its own unmarshaler
//	go func() {
//	    ctx := runtime.WithUnmarshaler(context.Background(), &SonicUnmarshaler{})
//	    vmssClient.Get(ctx, "rg", "vmss", nil)  // Uses sonic
//	}()
//
//	go func() {
//	    vmClient.Get(context.Background(), "rg", "vm", nil)  // Uses default
//	}()
func WithUnmarshaler(ctx context.Context, u Unmarshaler) context.Context {
	return exported.WithUnmarshaler(ctx, u)
}
