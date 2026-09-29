// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package azcosmos implements the client to interact with the Azure Cosmos DB SQL API.
//
// # Status
//
// This is the v2 major version of the module and it is not usable yet. The v2 surface is being
// assembled incrementally so that it can be reviewed as it lands. This release covers the error
// and response model, partition keys, client construction, and creating, reading, replacing,
// upserting, deleting, and patching single items.
//
// v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so
// that routing, retries, session handling, failover behavior and query fan-out are consistent
// across the Cosmos DB SDKs.
//
// # Options and lifetime
//
// OperationOptions applies at runtime, client, and request scope, with requests taking precedence.
// Rust owns resolution and supported environment overrides. Runtime.SetOperationOptions replaces
// defaults atomically; admitted operations retain a native snapshot and one timeout budget.
// Client.Close closes only that client; Runtime.Close drains and closes all attached clients.
//
// # Response encoding
//
// ItemResponse.Value can contain Cosmos binary JSON by default. Select BinaryEncodingOptions with
// Enabled and RequestTextResponse both true before decoding responses with encoding/json, or
// explicitly disable binary encoding. Go does not deserialize item schemas.
package azcosmos
