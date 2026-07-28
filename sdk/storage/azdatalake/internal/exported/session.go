// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// SessionMode specifies how session-based authentication is handled.
type SessionMode = azblob.SessionMode

const (
	// SessionModeDefault is the default mode where sessions are disabled.
	SessionModeDefault = azblob.SessionModeDefault
	// SessionModeDisabled explicitly disables session-based authentication.
	SessionModeDisabled = azblob.SessionModeDisabled
	// SessionModeEnabled enables session-based authentication.
	SessionModeEnabled = azblob.SessionModeEnabled
)

// SessionOptions configures session-based authentication behavior.
type SessionOptions = azblob.SessionOptions

// SessionCredential contains session authentication credentials.
type SessionCredential = azblob.SessionCredential

// SessionContext contains context for session operations.
type SessionContext = azblob.SessionContext

// SessionProvider is the interface for session-based authentication providers.
type SessionProvider = azblob.SessionProvider
