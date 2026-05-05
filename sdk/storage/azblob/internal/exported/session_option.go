// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

// SessionMode specifies how session-based authentication is handled.
type SessionMode string

const (
	// SessionModeDefault is the default mode where sessions are disabled.
	SessionModeDefault SessionMode = ""
	// SessionModeDisabled explicitly disables session-based authentication.
	SessionModeDisabled SessionMode = "disabled"
	// SessionModeEnabled enables session-based authentication.
	SessionModeEnabled SessionMode = "enabled"
)

// PossibleSessionModeValues returns a slice of possible values for SessionMode.
func PossibleSessionModeValues() []SessionMode {
	return []SessionMode{
		SessionModeDefault,
		SessionModeDisabled,
		SessionModeEnabled,
	}
}

// SessionOptions configures session-based authentication behavior.
type SessionOptions struct {
	// Mode specifies the session authentication mode.
	Mode SessionMode

	// AccountName is the storage account name.
	AccountName string
}
