// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package azcosmos implements the client to interact with the Azure Cosmos DB SQL API.
//
// # Status
//
// This is the v2 major version of the module and it is not usable yet. The v2 surface is being
// assembled incrementally so that it can be reviewed as it lands. This release covers the error
// and response model, partition keys, client construction, and creating, reading, replacing,
// upserting, deleting, and patching single items, plus retained item-query paging.
//
// Go 1.27 or later is required. Item writes accept models using encoding/json or ClientOptions.Serializer.
// ItemResponse.IntoModel decodes response bodies; generic query pagers decode each result as T.
// Use json.RawMessage for pre-encoded write payloads instead of plain []byte.
//
// QueryItemsPager supports logical-partition, hierarchical-prefix, and full-container scopes.
// Defer its Close method when stopping early. ContinuationToken captures resumable progress
// separately from paging; not every native query plan supports a checkpoint.
//
// v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so
// that routing, retries, session handling, failover behavior and query fan-out are consistent
// across the Cosmos DB SDKs.
//
// # Options and lifetime
//
// OperationOptions applies at runtime, client, and request scope, with requests taking precedence.
// Rust owns resolution and supported environment overrides; defaults are fixed at construction.
// Clients without an explicit Runtime use a process-wide runtime and share native account caches.
// Client.Close closes only that client; Runtime.Close releases ownership without closing attached clients.
// Context cancellation stops waiting, not submitted native work; writes may still commit.
// Close waits for native work to finish, including operations whose callers stopped waiting.
//
// # Advanced availability
//
// The native driver owns per-partition automatic failover and circuit-breaker routing; read
// hedging is configured through OperationOptions. See the module README's "Advanced availability"
// section for applicability, environment configuration, initialization timing, and diagnostics.
//
// # Response encoding
//
// Binary wire encoding is enabled by default; native code converts responses to text JSON.
// Set BinaryEncodingOptions.Enabled to false to disable binary wire encoding.
// ItemResponse.Value and QueryItemsResponse.Items contain text JSON, not decoded application models.
package azcosmos
