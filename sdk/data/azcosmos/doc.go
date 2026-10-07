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
// Initialization and query cancellation stop waiting, not submitted native work.
// Submitted point item operations await authoritative completion, even after context cancellation.
// Close waits for native work to finish, including operations whose callers stopped waiting.
//
// # Response encoding
//
// ItemResponse.Value can contain Cosmos binary JSON by default. Select BinaryEncodingOptions with
// RequestTextResponse true before decoding responses with encoding/json, or explicitly set Enabled
// to false. A nil Enabled defaults to true. Go does not deserialize item schemas.
// QueryItemsResponse.Items always contains text JSON values; queries default to text wire encoding.
package azcosmos
