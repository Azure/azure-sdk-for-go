// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// SessionToken is a marker of how far a client has read or written, used to guarantee a client
// observes its own writes under session consistency.
//
// Tokens are opaque and are produced by the service, so callers obtain one from a response rather
// than constructing it. Passing a token from one operation to another is how session consistency
// is carried across processes: take [ItemResponse.SessionToken] from a write and set it on the
// options of a later read that must observe it.
//
// The client captures and carries tokens on its own, so this is only needed when the two
// operations run in different processes.
//
// See https://learn.microsoft.com/azure/cosmos-db/consistency-levels#session-consistency.
type SessionToken string
