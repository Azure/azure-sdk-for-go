// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package azcosmos implements the client to interact with the Azure Cosmos DB SQL API.
//
// # Status
//
// This is the v2 major version of the module and it is not usable yet. The v2 surface is being
// assembled incrementally so that it can be reviewed as it lands; this release contains no public
// API. Operations arrive in subsequent releases and will initially return an error reporting that
// they are not implemented.
//
// v2 replaces the v1 pure-Go implementation with a binding to the shared Rust Cosmos driver, so
// that routing, retries, session handling, failover behavior and query fan-out are consistent
// across the Cosmos DB SDKs. The decision record is docs/adr/0001-go-v2-uses-ffi.md.
package azcosmos
